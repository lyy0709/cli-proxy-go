/*
 * 文件作用：邮箱验证码数据模型
 * 负责功能：
 *   - 邮箱验证码存储结构
 *   - 验证码用途定义
 * 重要程度：⭐⭐⭐ 重要（邮箱验证功能）
 * 依赖模块：无
 */
package model

import (
	"time"
)

// EmailVerification 邮箱验证码
type EmailVerification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Email     string    `gorm:"size:255;index;not null" json:"email"`
	Code      string    `gorm:"size:6;not null" json:"-"`          // 6位数字验证码
	Purpose   string    `gorm:"size:20;not null" json:"purpose"`   // register, reset_password
	ExpiresAt time.Time `gorm:"not null" json:"expires_at"`
	Used      bool      `gorm:"default:false" json:"used"`
	CreatedAt time.Time `json:"created_at"`
}

func (e *EmailVerification) TableName() string {
	return "email_verifications"
}

// 验证码用途常量
const (
	EmailPurposeRegister      = "register"       // 注册
	EmailPurposeResetPassword = "reset_password" // 重置密码
)

// IsExpired 检查验证码是否已过期
func (e *EmailVerification) IsExpired() bool {
	return time.Now().After(e.ExpiresAt)
}

// IsValid 检查验证码是否有效（未过期且未使用）
func (e *EmailVerification) IsValid() bool {
	return !e.Used && !e.IsExpired()
}
