/*
 * 文件作用：邮箱验证码仓库，处理验证码数据持久化
 * 负责功能：
 *   - 验证码创建
 *   - 验证码查询
 *   - 验证码验证和标记使用
 *   - 统计和清理
 * 重要程度：⭐⭐⭐ 重要（邮箱验证功能）
 * 依赖模块：model, gorm
 */
package repository

import (
	"time"

	"cli-proxy/internal/model"

	"gorm.io/gorm"
)

type EmailVerificationRepository struct {
	db *gorm.DB
}

func NewEmailVerificationRepository() *EmailVerificationRepository {
	return &EmailVerificationRepository{db: DB}
}

// Create 创建验证码记录
func (r *EmailVerificationRepository) Create(verification *model.EmailVerification) error {
	return r.db.Create(verification).Error
}

// GetLatestByEmail 获取指定邮箱的最新验证码
func (r *EmailVerificationRepository) GetLatestByEmail(email, purpose string) (*model.EmailVerification, error) {
	var verification model.EmailVerification
	err := r.db.Where("email = ? AND purpose = ?", email, purpose).
		Order("created_at DESC").
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// GetValidByEmailAndCode 获取有效的验证码（未过期且未使用）
func (r *EmailVerificationRepository) GetValidByEmailAndCode(email, code, purpose string) (*model.EmailVerification, error) {
	var verification model.EmailVerification
	err := r.db.Where("email = ? AND code = ? AND purpose = ? AND used = ? AND expires_at > ?",
		email, code, purpose, false, time.Now()).
		First(&verification).Error
	if err != nil {
		return nil, err
	}
	return &verification, nil
}

// MarkAsUsed 标记验证码为已使用
func (r *EmailVerificationRepository) MarkAsUsed(id uint) error {
	return r.db.Model(&model.EmailVerification{}).
		Where("id = ?", id).
		Update("used", true).Error
}

// CountTodayByEmail 统计今日发送次数
func (r *EmailVerificationRepository) CountTodayByEmail(email string) (int64, error) {
	var count int64
	today := time.Now().Format("2006-01-02")
	err := r.db.Model(&model.EmailVerification{}).
		Where("email = ? AND DATE(created_at) = ?", email, today).
		Count(&count).Error
	return count, err
}

// GetLastSendTime 获取最后发送时间
func (r *EmailVerificationRepository) GetLastSendTime(email, purpose string) (*time.Time, error) {
	var verification model.EmailVerification
	err := r.db.Where("email = ? AND purpose = ?", email, purpose).
		Order("created_at DESC").
		First(&verification).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &verification.CreatedAt, nil
}

// CleanExpired 清理过期的验证码
func (r *EmailVerificationRepository) CleanExpired() (int64, error) {
	result := r.db.Where("expires_at < ? OR used = ?", time.Now(), true).
		Delete(&model.EmailVerification{})
	return result.RowsAffected, result.Error
}

// InvalidatePreviousCodes 使之前的验证码失效
func (r *EmailVerificationRepository) InvalidatePreviousCodes(email, purpose string) error {
	return r.db.Model(&model.EmailVerification{}).
		Where("email = ? AND purpose = ? AND used = ?", email, purpose, false).
		Update("used", true).Error
}
