/*
 * 文件作用：邮箱验证码处理器，处理邮箱验证相关 API
 * 负责功能：
 *   - 发送验证码
 *   - 获取验证状态
 *   - 测试邮件发送
 * 重要程度：⭐⭐⭐ 重要（邮箱验证功能）
 * 依赖模块：service
 */
package handler

import (
	"cli-proxy/internal/model"
	"cli-proxy/internal/service"
	"cli-proxy/pkg/response"

	"github.com/gin-gonic/gin"
)

type EmailVerificationHandler struct {
	service   *service.EmailVerificationService
	emailSvc  *service.EmailService
	configSvc *service.ConfigService
}

func NewEmailVerificationHandler() *EmailVerificationHandler {
	return &EmailVerificationHandler{
		service:   service.GetEmailVerificationService(),
		emailSvc:  service.GetEmailService(),
		configSvc: service.GetConfigService(),
	}
}

// SendCode 发送验证码
// POST /api/auth/email/send-code
func (h *EmailVerificationHandler) SendCode(c *gin.Context) {
	var req service.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	// 验证 purpose 是否合法
	if req.Purpose != model.EmailPurposeRegister && req.Purpose != model.EmailPurposeResetPassword {
		response.BadRequest(c, "无效的验证类型")
		return
	}

	// 如果是注册，检查邮箱是否已存在
	if req.Purpose == model.EmailPurposeRegister {
		userSvc := service.NewUserService()
		if userSvc.ExistsByEmail(req.Email) {
			response.BadRequest(c, "该邮箱已被注册")
			return
		}
	}

	result, err := h.service.GenerateAndSend(req.Email, req.Purpose)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Success(c, result)
}

// GetStatus 获取邮箱验证状态
// GET /api/auth/email/status
func (h *EmailVerificationHandler) GetStatus(c *gin.Context) {
	status := h.service.GetStatus()
	response.Success(c, status)
}

// TestSend 测试邮件发送（管理员）
// POST /api/admin/email/test
func (h *EmailVerificationHandler) TestSend(c *gin.Context) {
	var req struct {
		ToEmail string `json:"to_email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.emailSvc.SendTestEmail(req.ToEmail); err != nil {
		response.InternalError(c, err.Error())
		return
	}

	response.Success(c, gin.H{
		"message": "测试邮件已发送",
	})
}
