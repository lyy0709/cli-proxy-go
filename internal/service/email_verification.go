/*
 * 文件作用：邮箱验证码服务，处理验证码生成和验证
 * 负责功能：
 *   - 验证码生成
 *   - 验证码发送
 *   - 验证码验证
 *   - 频率限制检查
 * 重要程度：⭐⭐⭐ 重要（邮箱验证功能）
 * 依赖模块：model, repository, email service
 */
package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"sync"
	"time"

	"cli-proxy/internal/model"
	"cli-proxy/internal/repository"
	"cli-proxy/pkg/logger"
)

// EmailVerificationService 邮箱验证码服务
type EmailVerificationService struct {
	repo      *repository.EmailVerificationRepository
	configSvc *ConfigService
	emailSvc  *EmailService
}

var emailVerificationService *EmailVerificationService
var emailVerificationOnce sync.Once

// GetEmailVerificationService 获取邮箱验证码服务单例
func GetEmailVerificationService() *EmailVerificationService {
	emailVerificationOnce.Do(func() {
		emailVerificationService = &EmailVerificationService{
			repo:      repository.NewEmailVerificationRepository(),
			configSvc: GetConfigService(),
			emailSvc:  GetEmailService(),
		}
	})
	return emailVerificationService
}

func getEmailVerificationLog() *logger.Logger {
	return logger.GetLogger("email_verification")
}

// SendCodeRequest 发送验证码请求
type SendCodeRequest struct {
	Email   string `json:"email" binding:"required,email"`
	Purpose string `json:"purpose" binding:"required"`
}

// SendCodeResponse 发送验证码响应
type SendCodeResponse struct {
	ExpireSeconds int `json:"expire_seconds"`
}

// GenerateAndSend 生成并发送验证码
func (s *EmailVerificationService) GenerateAndSend(email, purpose string) (*SendCodeResponse, error) {
	log := getEmailVerificationLog()

	// 检查邮件服务是否已配置
	if !s.emailSvc.IsConfigured() {
		return nil, fmt.Errorf("邮件服务未配置，请先配置 SMTP 服务器")
	}

	// 检查发送频率
	interval := s.configSvc.GetEmailSendInterval()
	lastSendTime, err := s.repo.GetLastSendTime(email, purpose)
	if err != nil {
		log.Error("[email_verification] 获取上次发送时间失败: %v", err)
	}
	if lastSendTime != nil {
		elapsed := time.Since(*lastSendTime)
		if elapsed < time.Duration(interval)*time.Second {
			waitSeconds := int(time.Duration(interval)*time.Second - elapsed) / int(time.Second)
			return nil, fmt.Errorf("发送过于频繁，请 %d 秒后再试", waitSeconds)
		}
	}

	// 检查每日发送上限
	dailyLimit := s.configSvc.GetEmailDailyLimit()
	todayCount, err := s.repo.CountTodayByEmail(email)
	if err != nil {
		log.Error("[email_verification] 统计今日发送次数失败: %v", err)
	}
	if todayCount >= int64(dailyLimit) {
		return nil, fmt.Errorf("今日发送次数已达上限（%d 次）", dailyLimit)
	}

	// 使之前的验证码失效
	if err := s.repo.InvalidatePreviousCodes(email, purpose); err != nil {
		log.Warn("[email_verification] 使之前验证码失效失败: %v", err)
	}

	// 生成验证码
	code, err := s.generateCode(6)
	if err != nil {
		return nil, fmt.Errorf("生成验证码失败: %v", err)
	}

	// 计算过期时间
	expireMinutes := s.configSvc.GetEmailCodeExpireMinutes()
	expiresAt := time.Now().Add(time.Duration(expireMinutes) * time.Minute)

	// 保存验证码
	verification := &model.EmailVerification{
		Email:     email,
		Code:      code,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		Used:      false,
	}
	if err := s.repo.Create(verification); err != nil {
		return nil, fmt.Errorf("保存验证码失败: %v", err)
	}

	// 发送邮件
	if err := s.emailSvc.SendVerificationCode(email, code, expireMinutes); err != nil {
		return nil, err
	}

	log.Info("[email_verification] 验证码已发送 | Email: %s | Purpose: %s", email, purpose)

	return &SendCodeResponse{
		ExpireSeconds: expireMinutes * 60,
	}, nil
}

// Verify 验证验证码
func (s *EmailVerificationService) Verify(email, code, purpose string) bool {
	log := getEmailVerificationLog()

	verification, err := s.repo.GetValidByEmailAndCode(email, code, purpose)
	if err != nil {
		log.Warn("[email_verification] 验证码验证失败 | Email: %s | Error: %v", email, err)
		return false
	}

	if verification == nil {
		log.Warn("[email_verification] 验证码不存在或已失效 | Email: %s", email)
		return false
	}

	// 标记为已使用
	if err := s.repo.MarkAsUsed(verification.ID); err != nil {
		log.Error("[email_verification] 标记验证码已使用失败: %v", err)
	}

	log.Info("[email_verification] 验证码验证成功 | Email: %s | Purpose: %s", email, purpose)
	return true
}

// generateCode 生成指定长度的数字验证码
func (s *EmailVerificationService) generateCode(length int) (string, error) {
	code := ""
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			return "", err
		}
		code += fmt.Sprintf("%d", n.Int64())
	}
	return code, nil
}

// GetStatus 获取邮箱验证状态
func (s *EmailVerificationService) GetStatus() map[string]interface{} {
	return map[string]interface{}{
		"email_verification_enabled": s.configSvc.GetEmailVerificationEnabled(),
		"send_interval":              s.configSvc.GetEmailSendInterval(),
		"configured":                 s.emailSvc.IsConfigured(),
	}
}

// CleanExpiredCodes 清理过期验证码
func (s *EmailVerificationService) CleanExpiredCodes() (int64, error) {
	return s.repo.CleanExpired()
}
