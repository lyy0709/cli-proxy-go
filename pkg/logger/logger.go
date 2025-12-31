/*
 * 文件作用：日志系统封装，基于zap提供模块化日志功能
 * 负责功能：
 *   - 多模块日志支持
 *   - 结构化日志（JSON格式）
 *   - Context日志追踪（request_id）
 *   - 输出到 stdout（Docker 友好）
 * 重要程度：⭐⭐⭐⭐ 重要（日志核心工具）
 * 依赖模块：zap
 */
package logger

import (
	"context"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// 日志级别常量（保持向后兼容）
const (
	LevelDebug = iota
	LevelInfo
	LevelWarn
	LevelError
)

// Context key for request_id
type contextKey string

const RequestIDKey contextKey = "request_id"

// Logger 日志器封装
type Logger struct {
	zap    *zap.Logger
	sugar  *zap.SugaredLogger
	module string
	level  zapcore.Level
}

var (
	defaultLogger *Logger
	loggers       = make(map[string]*Logger)
	mu            sync.RWMutex
	globalLevel   zap.AtomicLevel
)

// Init 初始化日志系统（dir 参数保留兼容，但不再使用）
func Init(dir string, level int) error {
	globalLevel = zap.NewAtomicLevelAt(intToZapLevel(level))

	var err error
	defaultLogger, err = newLogger("main")
	if err != nil {
		return err
	}

	return nil
}

// newLogger 创建新的日志器（输出到 stdout）
func newLogger(module string) (*Logger, error) {
	// 编码器配置 - 控制台友好格式
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "module",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder, // 彩色级别
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建 core - 输出到 stdout
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		globalLevel,
	)

	// 创建 logger
	zapLogger := zap.New(core,
		zap.AddCaller(),
		zap.AddCallerSkip(2), // 跳过封装层
	).Named(module)

	return &Logger{
		zap:    zapLogger,
		sugar:  zapLogger.Sugar(),
		module: module,
		level:  globalLevel.Level(),
	}, nil
}

// GetLogger 获取指定模块的日志器
func GetLogger(module string) *Logger {
	mu.RLock()
	if l, ok := loggers[module]; ok {
		mu.RUnlock()
		return l
	}
	mu.RUnlock()

	mu.Lock()
	defer mu.Unlock()

	// 双重检查
	if l, ok := loggers[module]; ok {
		return l
	}

	l, err := newLogger(module)
	if err != nil {
		// 如果创建失败，返回默认logger
		if defaultLogger != nil {
			return defaultLogger
		}
		// 极端情况：创建一个不写文件的logger
		zapLogger, _ := zap.NewProduction()
		return &Logger{
			zap:    zapLogger,
			sugar:  zapLogger.Sugar(),
			module: module,
		}
	}

	loggers[module] = l
	return l
}

// intToZapLevel 转换日志级别
func intToZapLevel(level int) zapcore.Level {
	switch level {
	case LevelDebug:
		return zapcore.DebugLevel
	case LevelInfo:
		return zapcore.InfoLevel
	case LevelWarn:
		return zapcore.WarnLevel
	case LevelError:
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// ============ 兼容旧API（格式化字符串） ============

// Debug 调试日志（格式化字符串）
func (l *Logger) Debug(format string, args ...interface{}) {
	l.sugar.Debugf(format, args...)
}

// Info 信息日志（格式化字符串）
func (l *Logger) Info(format string, args ...interface{}) {
	l.sugar.Infof(format, args...)
}

// Warn 警告日志（格式化字符串）
func (l *Logger) Warn(format string, args ...interface{}) {
	l.sugar.Warnf(format, args...)
}

// Error 错误日志（格式化字符串）
func (l *Logger) Error(format string, args ...interface{}) {
	l.sugar.Errorf(format, args...)
}

// ============ 新API（结构化日志） ============

// Field 日志字段类型
type Field = zap.Field

// 常用字段构造函数
var (
	String   = zap.String
	Int      = zap.Int
	Int64    = zap.Int64
	Uint     = zap.Uint
	Uint64   = zap.Uint64
	Float64  = zap.Float64
	Bool     = zap.Bool
	Duration = zap.Duration
	Time     = zap.Time
	Any      = zap.Any
	Err      = zap.Error
)

// DebugZ 结构化调试日志
func (l *Logger) DebugZ(msg string, fields ...Field) {
	l.zap.Debug(msg, fields...)
}

// InfoZ 结构化信息日志
func (l *Logger) InfoZ(msg string, fields ...Field) {
	l.zap.Info(msg, fields...)
}

// WarnZ 结构化警告日志
func (l *Logger) WarnZ(msg string, fields ...Field) {
	l.zap.Warn(msg, fields...)
}

// ErrorZ 结构化错误日志
func (l *Logger) ErrorZ(msg string, fields ...Field) {
	l.zap.Error(msg, fields...)
}

// ============ Context支持（带request_id） ============

// Ctx 从context获取request_id并返回带字段的logger
func (l *Logger) Ctx(ctx context.Context) *Logger {
	if ctx == nil {
		return l
	}

	requestID, ok := ctx.Value(RequestIDKey).(string)
	if !ok || requestID == "" {
		return l
	}

	// 返回带request_id字段的新logger
	return &Logger{
		zap:    l.zap.With(zap.String("request_id", requestID)),
		sugar:  l.sugar.With("request_id", requestID),
		module: l.module,
		level:  l.level,
	}
}

// With 添加固定字段
func (l *Logger) With(fields ...Field) *Logger {
	return &Logger{
		zap:    l.zap.With(fields...),
		sugar:  l.zap.With(fields...).Sugar(),
		module: l.module,
		level:  l.level,
	}
}

// ============ 默认日志器快捷方法 ============

func Debug(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Debug(format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Info(format, args...)
	}
}

func Warn(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Warn(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	if defaultLogger != nil {
		defaultLogger.Error(format, args...)
	}
}

// DebugZ 结构化调试日志（默认logger）
func DebugZ(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.DebugZ(msg, fields...)
	}
}

// InfoZ 结构化信息日志（默认logger）
func InfoZ(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.InfoZ(msg, fields...)
	}
}

// WarnZ 结构化警告日志（默认logger）
func WarnZ(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.WarnZ(msg, fields...)
	}
}

// ErrorZ 结构化错误日志（默认logger）
func ErrorZ(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.ErrorZ(msg, fields...)
	}
}

// ============ 配置和管理 ============

// Close 关闭日志系统
func Close() {
	mu.Lock()
	defer mu.Unlock()

	if defaultLogger != nil {
		defaultLogger.zap.Sync()
	}

	for _, l := range loggers {
		l.zap.Sync()
	}
}

// SetLevel 设置全局日志级别
func SetLevel(level int) {
	globalLevel.SetLevel(intToZapLevel(level))
}

// ParseLevel 解析日志级别字符串
func ParseLevel(levelStr string) int {
	switch levelStr {
	case "debug", "DEBUG":
		return LevelDebug
	case "info", "INFO":
		return LevelInfo
	case "warn", "WARN", "warning", "WARNING":
		return LevelWarn
	case "error", "ERROR":
		return LevelError
	default:
		return LevelInfo
	}
}

// GetZapLogger 获取底层zap.Logger（用于第三方库集成）
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.zap
}

// SetRequestID 设置context中的request_id
func SetRequestID(ctx context.Context, requestID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, requestID)
}

// GetRequestID 从context获取request_id
func GetRequestID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if id, ok := ctx.Value(RequestIDKey).(string); ok {
		return id
	}
	return ""
}
