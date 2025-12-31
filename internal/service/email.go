/*
 * 文件作用：邮件发送服务，通过 SMTP 发送邮件
 * 负责功能：
 *   - SMTP 邮件发送
 *   - HTML 邮件模板
 *   - 验证码邮件发送
 *   - 测试邮件发送
 * 重要程度：⭐⭐⭐ 重要（邮箱验证功能）
 * 依赖模块：model, gomail
 */
package service

import (
	"fmt"
	"sync"

	"cli-proxy/internal/model"
	"cli-proxy/pkg/logger"

	"gopkg.in/gomail.v2"
)

// EmailService 邮件发送服务
type EmailService struct {
	configSvc *ConfigService
}

var emailService *EmailService
var emailOnce sync.Once

// GetEmailService 获取邮件服务单例
func GetEmailService() *EmailService {
	emailOnce.Do(func() {
		emailService = &EmailService{
			configSvc: GetConfigService(),
		}
	})
	return emailService
}

func getEmailLog() *logger.Logger {
	return logger.GetLogger("email")
}

// SendVerificationCode 发送验证码邮件
func (s *EmailService) SendVerificationCode(to, code string, expireMinutes int) error {
	subject := "Cli-Proxy 邮箱验证码"
	body := s.buildVerificationCodeHTML(code, expireMinutes)
	return s.sendHTML(to, subject, body)
}

// SendTestEmail 发送测试邮件
func (s *EmailService) SendTestEmail(to string) error {
	subject := "Cli-Proxy 邮件配置测试"
	body := s.buildTestEmailHTML()
	return s.sendHTML(to, subject, body)
}

// sendHTML 发送 HTML 邮件
func (s *EmailService) sendHTML(to, subject, body string) error {
	log := getEmailLog()

	host := s.configSvc.GetString(model.ConfigSMTPHost)
	port := s.configSvc.GetInt(model.ConfigSMTPPort)
	username := s.configSvc.GetString(model.ConfigSMTPUsername)
	password := s.configSvc.GetString(model.ConfigSMTPPassword)
	fromEmail := s.configSvc.GetString(model.ConfigSMTPFromEmail)
	fromName := s.configSvc.GetString(model.ConfigSMTPFromName)

	if host == "" || fromEmail == "" {
		return fmt.Errorf("SMTP 配置不完整，请先配置邮件服务器")
	}

	if port == 0 {
		port = 587
	}

	if fromName == "" {
		fromName = "Cli-Proxy"
	}

	m := gomail.NewMessage()
	m.SetHeader("From", m.FormatAddress(fromEmail, fromName))
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(host, port, username, password)

	// 根据加密方式配置决定
	// ssl: 使用隐式 SSL/TLS（端口 465）
	// starttls: 使用 STARTTLS（端口 587，gomail 默认）
	// none: 不使用加密
	encryption := s.configSvc.GetString(model.ConfigSMTPEncryption)
	log.Info("[email] SMTP 配置 | Host: %s | Port: %d | Encryption: %s", host, port, encryption)

	switch encryption {
	case "ssl":
		d.SSL = true
	case "starttls":
		d.SSL = false // gomail 默认使用 STARTTLS
	case "none":
		d.SSL = false
	default:
		// 默认使用 STARTTLS
		log.Warn("[email] 未知的加密方式: %s，使用默认 STARTTLS", encryption)
		d.SSL = false
	}

	log.Info("[email] 发送邮件 | To: %s | Subject: %s", to, subject)

	if err := d.DialAndSend(m); err != nil {
		log.Error("[email] 邮件发送失败 | To: %s | Error: %v", to, err)
		return fmt.Errorf("邮件发送失败: %v", err)
	}

	log.Info("[email] 邮件发送成功 | To: %s", to)
	return nil
}

// buildVerificationCodeHTML 构建验证码邮件 HTML
func (s *EmailService) buildVerificationCodeHTML(code string, expireMinutes int) string {
	return fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>验证码</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f7;">
    <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" width="100%%" cellspacing="0" cellpadding="0" border="0" style="max-width: 480px; background-color: #ffffff; border-radius: 16px; box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center;">
                            <h1 style="margin: 0; font-size: 24px; font-weight: 600; color: #1d1d1f;">邮箱验证码</h1>
                        </td>
                    </tr>
                    <!-- Code -->
                    <tr>
                        <td style="padding: 20px 40px; text-align: center;">
                            <div style="display: inline-block; padding: 20px 40px; background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); border-radius: 12px;">
                                <span style="font-size: 36px; font-weight: 700; letter-spacing: 8px; color: #ffffff; font-family: 'SF Mono', Monaco, 'Courier New', monospace;">%s</span>
                            </div>
                        </td>
                    </tr>
                    <!-- Info -->
                    <tr>
                        <td style="padding: 20px 40px; text-align: center;">
                            <p style="margin: 0 0 12px; font-size: 15px; color: #86868b; line-height: 1.6;">
                                请在 <strong style="color: #1d1d1f;">%d 分钟</strong>内完成验证
                            </p>
                            <p style="margin: 0; font-size: 13px; color: #aeaeb2; line-height: 1.5;">
                                如果这不是您的操作，请忽略此邮件
                            </p>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px 40px; text-align: center; border-top: 1px solid #f0f0f0;">
                            <p style="margin: 0; font-size: 12px; color: #aeaeb2;">
                                此邮件由 Cli-Proxy 系统自动发送，请勿回复
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`, code, expireMinutes)
}

// buildTestEmailHTML 构建测试邮件 HTML
func (s *EmailService) buildTestEmailHTML() string {
	return `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>邮件配置测试</title>
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, 'Helvetica Neue', Arial, sans-serif; background-color: #f5f5f7;">
    <table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0">
        <tr>
            <td align="center" style="padding: 40px 20px;">
                <table role="presentation" width="100%" cellspacing="0" cellpadding="0" border="0" style="max-width: 480px; background-color: #ffffff; border-radius: 16px; box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);">
                    <!-- Header -->
                    <tr>
                        <td style="padding: 40px 40px 20px; text-align: center;">
                            <div style="width: 60px; height: 60px; margin: 0 auto 20px; background: linear-gradient(135deg, #34c759 0%, #30d158 100%); border-radius: 50%; display: flex; align-items: center; justify-content: center;">
                                <span style="font-size: 30px;">✓</span>
                            </div>
                            <h1 style="margin: 0; font-size: 24px; font-weight: 600; color: #1d1d1f;">邮件配置成功</h1>
                        </td>
                    </tr>
                    <!-- Content -->
                    <tr>
                        <td style="padding: 20px 40px; text-align: center;">
                            <p style="margin: 0; font-size: 15px; color: #86868b; line-height: 1.6;">
                                恭喜！您的邮件服务器配置正确，可以正常发送邮件。
                            </p>
                        </td>
                    </tr>
                    <!-- Footer -->
                    <tr>
                        <td style="padding: 30px 40px 40px; text-align: center; border-top: 1px solid #f0f0f0;">
                            <p style="margin: 0; font-size: 12px; color: #aeaeb2;">
                                此邮件由 Cli-Proxy 系统自动发送，请勿回复
                            </p>
                        </td>
                    </tr>
                </table>
            </td>
        </tr>
    </table>
</body>
</html>
`
}

// IsConfigured 检查邮件服务是否已配置
func (s *EmailService) IsConfigured() bool {
	host := s.configSvc.GetString(model.ConfigSMTPHost)
	fromEmail := s.configSvc.GetString(model.ConfigSMTPFromEmail)
	return host != "" && fromEmail != ""
}
