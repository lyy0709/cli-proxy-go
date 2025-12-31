<!--
 * 文件作用：系统首页，极简版
 * 负责功能：
 *   - API接口配置说明
 *   - 登录后台入口
 *   - 联系方式
 * 重要程度：⭐⭐⭐⭐⭐ 核心
 * 依赖组件：Element Plus, Vue Router
-->
<template>
  <div class="home-container">
    <!-- 头部 -->
    <header class="home-header">
      <div class="logo">
        <i class="fa-solid fa-atom"></i>
        <span>Go-AIProxy</span>
      </div>
      <router-link to="/login" class="btn-login">
        <i class="fa-solid fa-right-to-bracket"></i>
        登录后台
      </router-link>
    </header>

    <!-- 主内容 -->
    <main class="main-content">
      <div class="content-box">
        <h1>API 接口配置</h1>
        <p class="subtitle">多平台 AI API 统一接入，支持 Claude、OpenAI、Gemini</p>

        <!-- API 端点 -->
        <section class="api-section">
          <h2>服务端点</h2>

          <div class="endpoint-card">
            <div class="endpoint-header">
              <span class="platform-name">Claude</span>
              <button class="copy-btn" @click="copyText(baseUrl + '/claude')">
                <i class="fa-solid fa-copy"></i> 复制
              </button>
            </div>
            <code class="endpoint-url">{{ baseUrl }}/claude</code>
          </div>

          <div class="endpoint-card">
            <div class="endpoint-header">
              <span class="platform-name">OpenAI</span>
              <button class="copy-btn" @click="copyText(baseUrl + '/openai')">
                <i class="fa-solid fa-copy"></i> 复制
              </button>
            </div>
            <code class="endpoint-url">{{ baseUrl }}/openai</code>
          </div>

          <div class="endpoint-card">
            <div class="endpoint-header">
              <span class="platform-name">Gemini</span>
              <button class="copy-btn" @click="copyText(baseUrl + '/gemini')">
                <i class="fa-solid fa-copy"></i> 复制
              </button>
            </div>
            <code class="endpoint-url">{{ baseUrl }}/gemini</code>
          </div>
        </section>

        <!-- 配置示例 -->
        <section class="config-section">
          <h2>配置示例</h2>
          <div class="code-block">
            <pre><code># Claude 配置
export ANTHROPIC_BASE_URL="{{ baseUrl }}/claude"
export ANTHROPIC_API_KEY="你的API_Key"

# OpenAI 配置
export OPENAI_API_BASE="{{ baseUrl }}/openai/v1"
export OPENAI_API_KEY="你的API_Key"</code></pre>
          </div>
        </section>

        <!-- 登录提示 -->
        <section class="login-section">
          <p>需要 API Key？</p>
          <router-link to="/login" class="btn-primary">
            <i class="fa-solid fa-key"></i>
            登录后台创建 API Key
          </router-link>
        </section>
      </div>
    </main>

    <!-- 页脚 -->
    <footer class="home-footer">
      <div class="contact-section">
        <h3>联系我们</h3>
        <p class="wechat">
          <i class="fa-brands fa-weixin"></i>
          作者微信：suiyue_creation
        </p>
        <a href="https://qm.qq.com/q/iJ4bHLlMEa" target="_blank" class="qq-link">
          <i class="fa-brands fa-qq"></i>
          点击加入 QQ 交流群【go-proxy-pro】
        </a>
        <p class="footer-note">
          <a href="https://github.com/suiyuebaobao/go-proxy-pro" target="_blank">
            <i class="fa-brands fa-github"></i> GitHub
          </a>
        </p>
      </div>
      <p class="copyright">&copy; 2025 Go-AIProxy. All rights reserved.</p>
    </footer>
  </div>
</template>

<script setup>
import { ensureFontAwesomeLoaded } from '@/utils/fontawesome'
ensureFontAwesomeLoaded()

import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'

// 获取当前 Base URL
const baseUrl = ref(window.location.origin)

// 复制文本
const copyText = (text) => {
  navigator.clipboard.writeText(text).then(() => {
    ElMessage.success('已复制到剪贴板')
  }).catch(() => {
    ElMessage.error('复制失败')
  })
}
</script>

<style scoped>
.home-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: #0f172a;
  color: #e5e7eb;
}

/* 头部 */
.home-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem 2rem;
  border-bottom: 1px solid rgba(99, 102, 241, 0.2);
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 1.5rem;
  font-weight: 700;
  color: white;
}

.logo i {
  color: #a855f7;
  font-size: 1.75rem;
}

.btn-login {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.625rem 1.25rem;
  background: linear-gradient(135deg, #a855f7 0%, #6366f1 100%);
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-login:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.4);
}

/* 主内容 */
.main-content {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 3rem 2rem;
}

.content-box {
  max-width: 800px;
  width: 100%;
}

.content-box > h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  background: linear-gradient(135deg, #a855f7 0%, #6366f1 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.subtitle {
  color: #9ca3af;
  font-size: 1.125rem;
  margin-bottom: 3rem;
}

/* API 端点 */
.api-section {
  margin-bottom: 3rem;
}

.api-section h2 {
  font-size: 1.25rem;
  margin-bottom: 1.5rem;
  color: white;
}

.endpoint-card {
  background: rgba(30, 41, 59, 0.5);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 12px;
  padding: 1rem 1.25rem;
  margin-bottom: 1rem;
}

.endpoint-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

.platform-name {
  font-weight: 600;
  color: #e5e7eb;
}

.copy-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.375rem;
  padding: 0.375rem 0.75rem;
  background: rgba(99, 102, 241, 0.2);
  border: 1px solid rgba(99, 102, 241, 0.3);
  color: #a855f7;
  border-radius: 6px;
  font-size: 0.875rem;
  cursor: pointer;
  transition: all 0.2s;
}

.copy-btn:hover {
  background: rgba(99, 102, 241, 0.3);
  border-color: rgba(99, 102, 241, 0.5);
}

.endpoint-url {
  display: block;
  background: rgba(15, 23, 42, 0.5);
  padding: 0.75rem 1rem;
  border-radius: 8px;
  color: #a855f7;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.938rem;
  overflow-x: auto;
}

/* 配置示例 */
.config-section {
  margin-bottom: 3rem;
}

.config-section h2 {
  font-size: 1.25rem;
  margin-bottom: 1rem;
  color: white;
}

.code-block {
  background: rgba(15, 23, 42, 0.8);
  border: 1px solid rgba(99, 102, 241, 0.2);
  border-radius: 12px;
  padding: 1.25rem;
  overflow-x: auto;
}

.code-block pre {
  margin: 0;
}

.code-block code {
  color: #e5e7eb;
  font-family: 'Monaco', 'Menlo', monospace;
  font-size: 0.875rem;
  line-height: 1.7;
}

/* 登录提示 */
.login-section {
  text-align: center;
  padding: 2rem;
  background: rgba(30, 41, 59, 0.3);
  border-radius: 12px;
}

.login-section > p {
  color: #9ca3af;
  margin-bottom: 1rem;
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #a855f7 0%, #6366f1 100%);
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 500;
  transition: all 0.2s;
}

.btn-primary:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.4);
}

/* 页脚 */
.home-footer {
  border-top: 1px solid rgba(99, 102, 241, 0.2);
  padding: 3rem 2rem 2rem;
  text-align: center;
}

.contact-section h3 {
  color: white;
  margin-bottom: 1rem;
}

.wechat {
  color: #9ca3af;
  font-size: 1rem;
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.wechat i {
  color: #07c160;
  font-size: 1.25rem;
}

.qq-link {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.75rem 1.5rem;
  background: linear-gradient(135deg, #a855f7 0%, #6366f1 100%);
  color: white;
  text-decoration: none;
  border-radius: 8px;
  font-weight: 500;
  margin-bottom: 1rem;
  transition: all 0.2s;
}

.qq-link:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(168, 85, 247, 0.4);
}

.footer-note {
  margin: 1rem 0;
}

.footer-note a {
  color: #6b7280;
  text-decoration: none;
  transition: color 0.2s;
}

.footer-note a:hover {
  color: #a855f7;
}

.copyright {
  color: #6b7280;
  font-size: 0.875rem;
  margin-top: 1rem;
}

/* 响应式 */
@media (max-width: 768px) {
  .home-header {
    padding: 1rem;
  }

  .logo {
    font-size: 1.25rem;
  }

  .btn-login {
    padding: 0.5rem 1rem;
    font-size: 0.875rem;
  }

  .content-box > h1 {
    font-size: 1.75rem;
  }

  .main-content {
    padding: 2rem 1rem;
  }
}
</style>
