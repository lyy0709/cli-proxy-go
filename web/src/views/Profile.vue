<!--
 * 文件作用：个人设置页面 - Apple HIG 风格
 * 负责功能：
 *   - 基本信息修改（邮箱）
 *   - 密码修改
 *   - 用户配置保存
 * 重要程度：⭐⭐ 辅助（用户配置）
-->
<template>
  <div class="profile-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <h1 class="page-title">个人设置</h1>
      <p class="page-subtitle">管理您的账户信息和安全设置</p>
    </div>

    <div class="profile-grid">
      <!-- 基本信息 -->
      <div class="profile-card">
        <div class="card-header">
          <div class="card-icon">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
          </div>
          <h3>基本信息</h3>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label class="form-label">用户名</label>
            <input type="text" :value="userStore.user?.username" class="form-input" disabled />
            <p class="form-tip">用户名不可修改</p>
          </div>
          <div class="form-group">
            <label class="form-label">邮箱</label>
            <input v-model="profileForm.email" type="email" class="form-input" placeholder="your@email.com" />
          </div>
          <div class="card-actions">
            <button class="btn btn-primary" :disabled="savingProfile" @click="handleSaveProfile">
              <span v-if="savingProfile" class="btn-loading"></span>
              {{ savingProfile ? '保存中...' : '保存信息' }}
            </button>
          </div>
        </div>
      </div>

      <!-- 修改密码 -->
      <div class="profile-card">
        <div class="card-header">
          <div class="card-icon warning">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4"/>
            </svg>
          </div>
          <h3>修改密码</h3>
        </div>
        <div class="card-body">
          <div class="form-group">
            <label class="form-label">原密码 <span class="required">*</span></label>
            <div class="password-input">
              <input
                v-model="pwdForm.old_password"
                :type="showOldPwd ? 'text' : 'password'"
                class="form-input"
                placeholder="请输入原密码"
              />
              <button class="toggle-visibility" @click="showOldPwd = !showOldPwd">
                <svg v-if="showOldPwd" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              </button>
            </div>
          </div>
          <div class="form-group">
            <label class="form-label">新密码 <span class="required">*</span></label>
            <div class="password-input">
              <input
                v-model="pwdForm.new_password"
                :type="showNewPwd ? 'text' : 'password'"
                class="form-input"
                placeholder="请输入新密码（至少6位）"
              />
              <button class="toggle-visibility" @click="showNewPwd = !showNewPwd">
                <svg v-if="showNewPwd" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                  <circle cx="12" cy="12" r="3"/>
                </svg>
                <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                  <line x1="1" y1="1" x2="23" y2="23"/>
                </svg>
              </button>
            </div>
            <div v-if="pwdForm.new_password" class="password-strength">
              <div :class="['strength-bar', passwordStrength.level]">
                <span></span><span></span><span></span>
              </div>
              <span class="strength-text">{{ passwordStrength.text }}</span>
            </div>
          </div>
          <div class="card-actions">
            <button class="btn btn-primary" :disabled="savingPwd || !pwdForm.old_password || !pwdForm.new_password" @click="handleChangePassword">
              <span v-if="savingPwd" class="btn-loading"></span>
              {{ savingPwd ? '修改中...' : '修改密码' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const userStore = useUserStore()

const profileForm = reactive({ email: '' })
const savingProfile = ref(false)

const pwdForm = reactive({ old_password: '', new_password: '' })
const savingPwd = ref(false)
const showOldPwd = ref(false)
const showNewPwd = ref(false)

const passwordStrength = computed(() => {
  const pwd = pwdForm.new_password
  if (!pwd) return { level: '', text: '' }

  let score = 0
  if (pwd.length >= 6) score++
  if (pwd.length >= 10) score++
  if (/[A-Z]/.test(pwd)) score++
  if (/[a-z]/.test(pwd)) score++
  if (/[0-9]/.test(pwd)) score++
  if (/[^A-Za-z0-9]/.test(pwd)) score++

  if (score <= 2) return { level: 'weak', text: '弱' }
  if (score <= 4) return { level: 'medium', text: '中等' }
  return { level: 'strong', text: '强' }
})

onMounted(() => {
  profileForm.email = userStore.user?.email || ''
})

async function handleSaveProfile() {
  savingProfile.value = true
  try {
    await api.updateProfile({ email: profileForm.email })
    await userStore.fetchProfile()
    ElMessage.success('保存成功')
  } catch (e) {
    ElMessage.error('保存失败')
  } finally {
    savingProfile.value = false
  }
}

async function handleChangePassword() {
  if (!pwdForm.old_password || !pwdForm.new_password) {
    ElMessage.warning('请填写完整信息')
    return
  }

  if (pwdForm.new_password.length < 6) {
    ElMessage.warning('新密码至少6位')
    return
  }

  savingPwd.value = true
  try {
    await api.changePassword(pwdForm)
    ElMessage.success('密码修改成功')
    pwdForm.old_password = ''
    pwdForm.new_password = ''
  } catch (e) {
    ElMessage.error(e.message || '密码修改失败')
  } finally {
    savingPwd.value = false
  }
}
</script>

<style scoped>
.profile-page {
  max-width: 900px;
  margin: 0 auto;
}

/* 页面标题 */
.page-header {
  margin-bottom: var(--apple-spacing-xl);
}

.page-title {
  font-size: var(--apple-text-3xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-xs) 0;
}

.page-subtitle {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  margin: 0;
}

/* 卡片网格 */
.profile-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-lg);
}

/* 卡片 */
.profile-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.card-header {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}

.card-icon {
  width: 40px;
  height: 40px;
  border-radius: var(--apple-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.card-icon.warning {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.card-icon svg {
  width: 20px;
  height: 20px;
  color: white;
}

.card-header h3 {
  margin: 0;
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.card-body {
  padding: var(--apple-spacing-lg);
}

.card-actions {
  margin-top: var(--apple-spacing-lg);
  padding-top: var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
}

/* 表单 */
.form-group {
  margin-bottom: var(--apple-spacing-lg);
}

.form-label {
  display: block;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xs);
}

.form-label .required {
  color: var(--apple-red);
}

.form-input {
  width: 100%;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast);
}

.form-input:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.form-input:disabled {
  background: var(--apple-bg-secondary);
  color: var(--apple-text-tertiary);
  cursor: not-allowed;
}

.form-tip {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: var(--apple-spacing-xxs);
}

/* 密码输入 */
.password-input {
  position: relative;
}

.password-input .form-input {
  padding-right: 44px;
}

.toggle-visibility {
  position: absolute;
  right: var(--apple-spacing-sm);
  top: 50%;
  transform: translateY(-50%);
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast);
}

.toggle-visibility:hover {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-secondary);
}

.toggle-visibility svg {
  width: 18px;
  height: 18px;
}

/* 密码强度 */
.password-strength {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  margin-top: var(--apple-spacing-sm);
}

.strength-bar {
  display: flex;
  gap: 4px;
}

.strength-bar span {
  width: 40px;
  height: 4px;
  background: var(--apple-fill-tertiary);
  border-radius: 2px;
  transition: background var(--apple-duration-fast);
}

.strength-bar.weak span:nth-child(1) {
  background: var(--apple-red);
}

.strength-bar.medium span:nth-child(1),
.strength-bar.medium span:nth-child(2) {
  background: var(--apple-orange);
}

.strength-bar.strong span {
  background: var(--apple-green);
}

.strength-text {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

.strength-bar.weak + .strength-text { color: var(--apple-red); }
.strength-bar.medium + .strength-text { color: var(--apple-orange); }
.strength-bar.strong + .strength-text { color: var(--apple-green); }

/* 按钮 */
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast);
  cursor: pointer;
  width: 100%;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-primary {
  background: var(--apple-blue);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--apple-blue-hover);
}

.btn-loading {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 响应式 */
@media (max-width: 768px) {
  .profile-grid {
    grid-template-columns: 1fr;
  }
}
</style>
