<!--
 * 文件作用：用户注册页面 - Apple HIG 风格
 * 负责功能：
 *   - 用户名/邮箱/密码输入
 *   - 验证码获取和验证
 *   - 注册请求和状态管理
 * 重要程度：⭐⭐⭐⭐ 重要（用户入口）
-->
<template>
  <div class="register-page">
    <!-- 背景装饰 -->
    <div class="bg-decoration">
      <div class="bg-circle bg-circle-1"></div>
      <div class="bg-circle bg-circle-2"></div>
      <div class="bg-circle bg-circle-3"></div>
    </div>

    <!-- 注册卡片 -->
    <div class="register-card">
      <!-- Logo -->
      <div class="register-header">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
        </div>
        <h1 class="register-title">创建账户</h1>
        <p class="register-subtitle">加入 Cli-Proxy 平台</p>
      </div>

      <!-- 注册表单 -->
      <form class="register-form" @submit.prevent="handleRegister">
        <!-- 用户名 -->
        <div class="form-group">
          <label class="form-label">用户名 <span class="required">*</span></label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/>
              <circle cx="12" cy="7" r="4"/>
            </svg>
            <input
              v-model="form.username"
              type="text"
              class="form-input"
              placeholder="请输入用户名（3-50字符）"
              autocomplete="username"
            />
          </div>
          <span v-if="errors.username" class="form-error">{{ errors.username }}</span>
        </div>

        <!-- 邮箱 -->
        <div class="form-group">
          <label class="form-label">邮箱 <span class="required">*</span></label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
              <polyline points="22,6 12,13 2,6"/>
            </svg>
            <input
              v-model="form.email"
              type="email"
              class="form-input"
              placeholder="请输入邮箱地址"
              autocomplete="email"
            />
          </div>
          <span v-if="errors.email" class="form-error">{{ errors.email }}</span>
        </div>

        <!-- 密码 -->
        <div class="form-group">
          <label class="form-label">密码 <span class="required">*</span></label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4"/>
            </svg>
            <input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请输入密码（至少6位）"
              autocomplete="new-password"
            />
            <button type="button" class="password-toggle" @click="showPassword = !showPassword">
              <svg v-if="!showPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
            </button>
          </div>
          <span v-if="errors.password" class="form-error">{{ errors.password }}</span>
        </div>

        <!-- 确认密码 -->
        <div class="form-group">
          <label class="form-label">确认密码 <span class="required">*</span></label>
          <div class="input-wrapper">
            <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M9 12l2 2 4-4"/>
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2"/>
              <path d="M7 11V7a5 5 0 0110 0v4"/>
            </svg>
            <input
              v-model="form.confirmPassword"
              :type="showConfirmPassword ? 'text' : 'password'"
              class="form-input"
              placeholder="请再次输入密码"
              autocomplete="new-password"
            />
            <button type="button" class="password-toggle" @click="showConfirmPassword = !showConfirmPassword">
              <svg v-if="!showConfirmPassword" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z"/>
                <circle cx="12" cy="12" r="3"/>
              </svg>
              <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M17.94 17.94A10.07 10.07 0 0112 20c-7 0-11-8-11-8a18.45 18.45 0 015.06-5.94M9.9 4.24A9.12 9.12 0 0112 4c7 0 11 8 11 8a18.5 18.5 0 01-2.16 3.19m-6.72-1.07a3 3 0 11-4.24-4.24"/>
                <line x1="1" y1="1" x2="23" y2="23"/>
              </svg>
            </button>
          </div>
          <span v-if="errors.confirmPassword" class="form-error">{{ errors.confirmPassword }}</span>
        </div>

        <!-- 验证码 -->
        <div v-if="captchaEnabled" class="form-group">
          <label class="form-label">验证码 <span class="required">*</span></label>
          <div class="captcha-row">
            <div class="input-wrapper captcha-input">
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
                <line x1="9" y1="9" x2="15" y2="15"/>
                <line x1="15" y1="9" x2="9" y2="15"/>
              </svg>
              <input
                v-model="form.captcha"
                type="text"
                class="form-input"
                placeholder="请输入验证码"
                @keyup.enter="handleRegister"
              />
            </div>
            <div class="captcha-image-wrapper" @click="refreshCaptcha">
              <img v-if="captchaImage" :src="captchaImage" class="captcha-image" alt="验证码" />
              <div v-else class="captcha-placeholder">
                <svg class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
                  <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
                </svg>
              </div>
            </div>
          </div>
          <span v-if="errors.captcha" class="form-error">{{ errors.captcha }}</span>
        </div>

        <!-- 邮箱验证码 -->
        <div v-if="emailVerificationEnabled" class="form-group">
          <label class="form-label">邮箱验证码 <span class="required">*</span></label>
          <div class="email-code-row">
            <div class="input-wrapper email-code-input">
              <svg class="input-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M4 4h16c1.1 0 2 .9 2 2v12c0 1.1-.9 2-2 2H4c-1.1 0-2-.9-2-2V6c0-1.1.9-2 2-2z"/>
                <path d="M9 12l2 2 4-4"/>
              </svg>
              <input
                v-model="form.emailCode"
                type="text"
                class="form-input"
                placeholder="请输入6位验证码"
                maxlength="6"
                @keyup.enter="handleRegister"
              />
            </div>
            <button
              type="button"
              class="send-code-btn"
              :disabled="!canSendEmailCode"
              @click="sendEmailCode"
            >
              <svg v-if="emailSending" class="spinner-small" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
                <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
              </svg>
              <span v-else-if="emailCountdown > 0">{{ emailCountdown }}s</span>
              <span v-else>获取验证码</span>
            </button>
          </div>
          <span v-if="errors.emailCode" class="form-error">{{ errors.emailCode }}</span>
        </div>

        <!-- 注册按钮 -->
        <button type="submit" class="register-btn" :disabled="loading">
          <svg v-if="loading" class="spinner" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10" stroke-opacity="0.25"/>
            <path d="M12 2a10 10 0 019.95 9" stroke-linecap="round"/>
          </svg>
          <span v-else>立即注册</span>
        </button>
      </form>

      <!-- 底部链接 -->
      <div class="register-footer">
        <span class="footer-text">已有账户？</span>
        <router-link to="/login" class="login-link">立即登录</router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const router = useRouter()

const loading = ref(false)
const showPassword = ref(false)
const showConfirmPassword = ref(false)
const captchaImage = ref('')
const captchaId = ref('')
const captchaEnabled = ref(true)

// 邮箱验证码状态
const emailVerificationEnabled = ref(false)
const emailSendInterval = ref(60)
const emailCountdown = ref(0)
const emailSending = ref(false)
let countdownTimer = null

const form = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  captcha: '',
  emailCode: ''
})

const errors = reactive({
  username: '',
  email: '',
  password: '',
  confirmPassword: '',
  captcha: '',
  emailCode: ''
})

// 是否可以发送邮箱验证码
const canSendEmailCode = computed(() => {
  if (emailSending.value || emailCountdown.value > 0) return false
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return form.email && emailRegex.test(form.email)
})

function validate() {
  let valid = true
  errors.username = ''
  errors.email = ''
  errors.password = ''
  errors.confirmPassword = ''
  errors.captcha = ''
  errors.emailCode = ''

  // 用户名验证
  if (!form.username.trim()) {
    errors.username = '请输入用户名'
    valid = false
  } else if (form.username.length < 3) {
    errors.username = '用户名至少3个字符'
    valid = false
  } else if (form.username.length > 50) {
    errors.username = '用户名不能超过50个字符'
    valid = false
  }

  // 邮箱验证（必填）
  if (!form.email.trim()) {
    errors.email = '请输入邮箱地址'
    valid = false
  } else {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
    if (!emailRegex.test(form.email)) {
      errors.email = '请输入有效的邮箱地址'
      valid = false
    }
  }

  // 密码验证
  if (!form.password) {
    errors.password = '请输入密码'
    valid = false
  } else if (form.password.length < 6) {
    errors.password = '密码至少6位'
    valid = false
  }

  // 确认密码验证
  if (!form.confirmPassword) {
    errors.confirmPassword = '请确认密码'
    valid = false
  } else if (form.password !== form.confirmPassword) {
    errors.confirmPassword = '两次输入的密码不一致'
    valid = false
  }

  // 验证码验证
  if (captchaEnabled.value && !form.captcha.trim()) {
    errors.captcha = '请输入验证码'
    valid = false
  }

  // 邮箱验证码验证
  if (emailVerificationEnabled.value && !form.emailCode.trim()) {
    errors.emailCode = '请输入邮箱验证码'
    valid = false
  }

  return valid
}

async function refreshCaptcha() {
  try {
    const res = await api.getCaptcha()
    captchaEnabled.value = res.data.enabled !== false
    if (captchaEnabled.value) {
      captchaId.value = res.data.captcha_id
      captchaImage.value = res.data.image
    }
  } catch {
    captchaEnabled.value = false
  }
}

async function fetchEmailStatus() {
  try {
    const res = await api.getEmailStatus()
    emailVerificationEnabled.value = res.data.email_verification_enabled === true
    emailSendInterval.value = res.data.send_interval || 60
  } catch {
    emailVerificationEnabled.value = false
  }
}

async function sendEmailCode() {
  if (!canSendEmailCode.value) return

  errors.emailCode = ''
  emailSending.value = true

  try {
    const res = await api.sendEmailCode(form.email, 'register')
    ElMessage.success('验证码已发送到您的邮箱')
    // 开始倒计时
    emailCountdown.value = emailSendInterval.value
    countdownTimer = setInterval(() => {
      emailCountdown.value--
      if (emailCountdown.value <= 0) {
        clearInterval(countdownTimer)
        countdownTimer = null
      }
    }, 1000)
  } catch (error) {
    // 错误信息已由 api 拦截器处理
  } finally {
    emailSending.value = false
  }
}

async function handleRegister() {
  if (!validate()) return

  loading.value = true
  try {
    const registerData = {
      username: form.username,
      password: form.password,
      email: form.email
    }
    if (captchaEnabled.value) {
      registerData.captcha_id = captchaId.value
      registerData.captcha = form.captcha
    }
    if (emailVerificationEnabled.value) {
      registerData.email_code = form.emailCode
    }

    await api.register(registerData)
    ElMessage.success('注册成功，请登录')
    router.push('/login')
  } catch (error) {
    // 刷新验证码
    if (captchaEnabled.value) {
      refreshCaptcha()
      form.captcha = ''
    }
    // 错误信息已由 api 拦截器处理
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  refreshCaptcha()
  fetchEmailStatus()
})
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: var(--apple-spacing-xl);
  position: relative;
  overflow: hidden;
}

/* 背景装饰 */
.bg-decoration {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.bg-circle {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
}

.bg-circle-1 {
  width: 600px;
  height: 600px;
  top: -200px;
  right: -100px;
}

.bg-circle-2 {
  width: 400px;
  height: 400px;
  bottom: -100px;
  left: -50px;
}

.bg-circle-3 {
  width: 200px;
  height: 200px;
  top: 50%;
  left: 20%;
}

/* 注册卡片 */
.register-card {
  width: 100%;
  max-width: 420px;
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xxl);
  box-shadow: var(--apple-shadow-xl);
  padding: var(--apple-spacing-xxl);
  position: relative;
  z-index: 1;
  animation: apple-scale-in var(--apple-duration-normal) var(--apple-ease-spring);
}

/* 头部 */
.register-header {
  text-align: center;
  margin-bottom: var(--apple-spacing-xl);
}

.logo {
  width: 64px;
  height: 64px;
  margin: 0 auto var(--apple-spacing-lg);
  background: linear-gradient(135deg, var(--apple-green) 0%, var(--apple-teal) 100%);
  border-radius: var(--apple-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
}

.logo-icon {
  width: 36px;
  height: 36px;
  color: white;
}

.register-title {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xs);
}

.register-subtitle {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

/* 表单 */
.register-form {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-md);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
}

.form-label {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.required {
  color: var(--apple-red);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: var(--apple-spacing-md);
  width: 18px;
  height: 18px;
  color: var(--apple-text-tertiary);
  pointer-events: none;
}

.form-input {
  width: 100%;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  padding-left: 44px;
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
  background: var(--apple-fill-quaternary);
  border: 1px solid transparent;
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-input::placeholder {
  color: var(--apple-text-placeholder);
}

.form-input:hover {
  background: var(--apple-fill-tertiary);
}

.form-input:focus {
  background: var(--apple-bg-primary);
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.password-toggle {
  position: absolute;
  right: var(--apple-spacing-sm);
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.password-toggle:hover {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-secondary);
}

.password-toggle svg {
  width: 18px;
  height: 18px;
}

.form-error {
  font-size: var(--apple-text-xs);
  color: var(--apple-red);
}

/* 验证码 */
.captcha-row {
  display: flex;
  gap: var(--apple-spacing-sm);
}

.captcha-input {
  flex: 1;
}

.captcha-image-wrapper {
  width: 120px;
  height: 44px;
  border-radius: var(--apple-radius-sm);
  cursor: pointer;
  background: var(--apple-fill-quaternary);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  border: 1px solid var(--apple-separator);
}

.captcha-image {
  width: 100%;
  height: 100%;
  object-fit: contain;
  border-radius: var(--apple-radius-xs);
}

.captcha-placeholder {
  display: flex;
  align-items: center;
  justify-content: center;
}

/* 邮箱验证码 */
.email-code-row {
  display: flex;
  gap: var(--apple-spacing-sm);
}

.email-code-input {
  flex: 1;
}

.send-code-btn {
  min-width: 110px;
  height: 44px;
  padding: 0 var(--apple-spacing-md);
  background: linear-gradient(135deg, var(--apple-blue) 0%, #5856d6 100%);
  color: white;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.send-code-btn:hover:not(:disabled) {
  transform: scale(1.02);
  box-shadow: 0 4px 12px rgba(0, 122, 255, 0.3);
}

.send-code-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.send-code-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-tertiary);
}

.spinner-small {
  width: 16px;
  height: 16px;
  animation: spin 1s linear infinite;
}

/* 注册按钮 */
.register-btn {
  width: 100%;
  padding: var(--apple-spacing-md);
  background: linear-gradient(135deg, var(--apple-green) 0%, var(--apple-teal) 100%);
  color: white;
  font-size: var(--apple-text-base);
  font-weight: var(--apple-font-semibold);
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  margin-top: var(--apple-spacing-sm);
}

.register-btn:hover:not(:disabled) {
  transform: scale(1.02);
  box-shadow: 0 4px 16px rgba(52, 199, 89, 0.4);
}

.register-btn:active:not(:disabled) {
  transform: scale(0.98);
}

.register-btn:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

/* Spinner */
.spinner {
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  from { transform: rotate(0deg); }
  to { transform: rotate(360deg); }
}

/* 底部 */
.register-footer {
  margin-top: var(--apple-spacing-xl);
  text-align: center;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--apple-spacing-xs);
}

.footer-text {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.login-link {
  font-size: var(--apple-text-sm);
  color: var(--apple-blue);
  text-decoration: none;
  font-weight: var(--apple-font-medium);
  transition: color var(--apple-duration-fast) var(--apple-ease-default);
}

.login-link:hover {
  color: var(--apple-blue-hover);
  text-decoration: underline;
}

/* 响应式 */
@media (max-width: 480px) {
  .register-card {
    padding: var(--apple-spacing-xl);
  }

  .captcha-image-wrapper {
    width: 100px;
  }
}
</style>
