<!--
 * 文件作用：用户仪表盘页面 - Apple HIG 风格
 * 负责功能：
 *   - 显示用户使用概览
 *   - 套餐余额和状态
 *   - API Key 统计
 * 重要程度：⭐⭐⭐⭐ 重要（用户首页）
-->
<template>
  <div class="dashboard">
    <!-- 欢迎区域 -->
    <div class="welcome-section">
      <h1 class="welcome-title">欢迎回来，{{ userStore.user?.username || '用户' }}</h1>
      <p class="welcome-subtitle">这里是您的使用概览</p>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon cost">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <line x1="12" y1="1" x2="12" y2="23"/>
            <path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">${{ formatNumber(stats.today_cost, 4) }}</div>
          <div class="stat-label">今日消费</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon tokens">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 6v12M6 12h12"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ formatLargeNumber(stats.today_tokens) }}</div>
          <div class="stat-label">今日 Token</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon requests">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <polyline points="22,12 18,12 15,21 9,3 6,12 2,12"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ formatNumber(stats.today_requests) }}</div>
          <div class="stat-label">今日请求</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon keys">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
          </svg>
        </div>
        <div class="stat-content">
          <div class="stat-value">{{ stats.api_key_count }}</div>
          <div class="stat-label">API Key</div>
        </div>
      </div>
    </div>

    <!-- 套餐和 API Key -->
    <div class="content-grid">
      <!-- 我的套餐 -->
      <div class="content-card">
        <div class="card-header">
          <div class="card-title">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z"/>
              <line x1="3" y1="6" x2="21" y2="6"/>
              <path d="M16 10a4 4 0 01-8 0"/>
            </svg>
            我的套餐
          </div>
          <button class="card-action" @click="$router.push('/user/packages')">
            查看全部
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9,18 15,12 9,6"/>
            </svg>
          </button>
        </div>
        <div class="card-body">
          <div v-if="packages.length > 0" class="package-list">
            <div v-for="pkg in packages.slice(0, 3)" :key="pkg.id" class="package-item">
              <div class="package-header">
                <div class="package-name">
                  <span :class="['package-badge', pkg.type]">
                    {{ pkg.type === 'subscription' ? '订阅' : '额度' }}
                  </span>
                  {{ pkg.name }}
                </div>
                <span :class="['status-badge', pkg.status]">
                  {{ pkg.status === 'active' ? '有效' : pkg.status === 'expired' ? '已过期' : pkg.status }}
                </span>
              </div>
              <div v-if="pkg.type === 'quota'" class="package-progress">
                <div class="progress-bar">
                  <div
                    class="progress-fill"
                    :style="{ width: getQuotaPercentage(pkg) + '%' }"
                    :class="getQuotaClass(pkg)"
                  ></div>
                </div>
                <div class="progress-text">
                  已用 ${{ pkg.quota_used?.toFixed(2) || '0' }} / ${{ pkg.quota_total?.toFixed(2) || '0' }}
                </div>
              </div>
              <div v-else class="package-subscription">
                <span>日额度: ${{ pkg.daily_used?.toFixed(2) || '0' }} / ${{ pkg.daily_quota?.toFixed(2) || '无限' }}</span>
                <span class="expire-info">到期: {{ formatDate(pkg.expire_time) }}</span>
              </div>
            </div>
          </div>
          <div v-else class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z"/>
              <line x1="3" y1="6" x2="21" y2="6"/>
            </svg>
            <span>暂无套餐</span>
          </div>
        </div>
      </div>

      <!-- 我的 API Key -->
      <div class="content-card">
        <div class="card-header">
          <div class="card-title">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
            </svg>
            我的 API Key
          </div>
          <button class="card-action" @click="$router.push('/user/api-keys')">
            管理
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <polyline points="9,18 15,12 9,6"/>
            </svg>
          </button>
        </div>
        <div class="card-body">
          <div v-if="apiKeys.length > 0" class="apikey-list">
            <div v-for="key in apiKeys.slice(0, 3)" :key="key.id" class="apikey-item">
              <div class="apikey-info">
                <div class="apikey-name">{{ key.name }}</div>
                <code class="apikey-prefix">{{ key.key_prefix }}...</code>
              </div>
              <span :class="['status-badge', key.status === 'active' ? 'active' : 'inactive']">
                {{ key.status === 'active' ? '正常' : '禁用' }}
              </span>
            </div>
          </div>
          <div v-else class="empty-state">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
            </svg>
            <span>暂无 API Key</span>
            <button class="empty-action" @click="$router.push('/user/api-keys')">
              创建 API Key
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 累计统计 -->
    <div class="total-stats-card">
      <div class="card-header">
        <div class="card-title">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M3 3v18h18"/>
            <path d="M18 9l-5 5-4-4-3 3"/>
          </svg>
          累计使用统计
        </div>
      </div>
      <div class="total-stats-grid">
        <div class="total-stat">
          <div class="total-value">${{ formatNumber(stats.total_cost, 2) }}</div>
          <div class="total-label">累计消费</div>
        </div>
        <div class="total-stat">
          <div class="total-value">{{ formatLargeNumber(stats.total_tokens) }}</div>
          <div class="total-label">累计 Token</div>
        </div>
        <div class="total-stat">
          <div class="total-value">{{ formatNumber(stats.total_requests) }}</div>
          <div class="total-label">累计请求</div>
        </div>
        <div class="total-stat">
          <div class="total-value">{{ stats.model_count }}</div>
          <div class="total-label">使用模型数</div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import api from '@/api'

const userStore = useUserStore()

const stats = reactive({
  today_cost: 0,
  today_tokens: 0,
  today_requests: 0,
  api_key_count: 0,
  total_cost: 0,
  total_tokens: 0,
  total_requests: 0,
  model_count: 0
})

const packages = ref([])
const apiKeys = ref([])

const formatNumber = (num, decimals = 0) => {
  if (num === undefined || num === null) return '0'
  if (decimals > 0) return num.toFixed(decimals)
  return num.toLocaleString()
}

const formatLargeNumber = (num) => {
  if (!num) return '0'
  if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toLocaleString()
}

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const getQuotaPercentage = (pkg) => {
  if (!pkg.quota_total || pkg.quota_total === 0) return 0
  return Math.min(100, (pkg.quota_used / pkg.quota_total) * 100)
}

const getQuotaClass = (pkg) => {
  const percentage = getQuotaPercentage(pkg)
  if (percentage < 60) return 'success'
  if (percentage < 80) return 'warning'
  return 'danger'
}

const fetchStats = async () => {
  try {
    const summaryRes = await api.getUserUsageSummary()
    if (summaryRes.data) {
      stats.today_cost = summaryRes.data.today?.total_cost || 0
      stats.today_tokens = summaryRes.data.today?.total_tokens || 0
      stats.today_requests = summaryRes.data.today?.request_count || 0
      stats.total_cost = summaryRes.data.total?.total_cost || 0
      stats.total_tokens = summaryRes.data.total?.total_tokens || 0
      stats.total_requests = summaryRes.data.total?.request_count || 0
      stats.model_count = summaryRes.data.model_count || 0
    }
  } catch (e) {
    console.error('Failed to fetch stats:', e)
  }
}

const fetchPackages = async () => {
  try {
    const res = await api.getMyActivePackages()
    packages.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch packages:', e)
  }
}

const fetchApiKeys = async () => {
  try {
    const res = await api.getApiKeys()
    apiKeys.value = res.data || []
    stats.api_key_count = apiKeys.value.length
  } catch (e) {
    console.error('Failed to fetch API keys:', e)
  }
}

onMounted(() => {
  fetchStats()
  fetchPackages()
  fetchApiKeys()
})
</script>

<style scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
}

/* 欢迎区域 */
.welcome-section {
  margin-bottom: var(--apple-spacing-xl);
}

.welcome-title {
  font-size: var(--apple-text-3xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-xs) 0;
}

.welcome-subtitle {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  margin: 0;
}

/* 统计卡片 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--apple-spacing-lg);
  margin-bottom: var(--apple-spacing-xl);
}

.stat-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  padding: var(--apple-spacing-xl);
  box-shadow: var(--apple-shadow-card);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-lg);
  transition: all var(--apple-duration-normal) var(--apple-ease-default);
}

.stat-card:hover {
  box-shadow: var(--apple-shadow-card-hover);
  transform: translateY(-2px);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: var(--apple-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 28px;
  height: 28px;
  color: white;
}

.stat-icon.cost {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.stat-icon.tokens {
  background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
}

.stat-icon.requests {
  background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
}

.stat-icon.keys {
  background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
}

.stat-label {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-top: var(--apple-spacing-xxs);
}

/* 内容区域 */
.content-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-lg);
  margin-bottom: var(--apple-spacing-xl);
}

.content-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-lg) var(--apple-spacing-xl);
  border-bottom: 1px solid var(--apple-separator);
}

.card-title {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.card-title svg {
  width: 20px;
  height: 20px;
  color: var(--apple-blue);
}

.card-action {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
  font-size: var(--apple-text-sm);
  color: var(--apple-blue);
  transition: opacity var(--apple-duration-fast) var(--apple-ease-default);
}

.card-action:hover {
  opacity: 0.7;
}

.card-action svg {
  width: 16px;
  height: 16px;
}

.card-body {
  padding: var(--apple-spacing-lg);
}

/* 套餐列表 */
.package-list, .apikey-list {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-md);
}

.package-item, .apikey-item {
  padding: var(--apple-spacing-md);
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
}

.package-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: var(--apple-spacing-sm);
}

.package-name {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.package-badge {
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.package-badge.subscription {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

.package-badge.quota {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.status-badge {
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.status-badge.active {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.status-badge.expired, .status-badge.inactive {
  background: var(--apple-fill-secondary);
  color: var(--apple-text-tertiary);
}

.package-progress {
  margin-top: var(--apple-spacing-sm);
}

.progress-bar {
  height: 6px;
  background: var(--apple-fill-tertiary);
  border-radius: var(--apple-radius-full);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: var(--apple-radius-full);
  transition: width var(--apple-duration-normal) var(--apple-ease-default);
}

.progress-fill.success { background: var(--apple-green); }
.progress-fill.warning { background: var(--apple-orange); }
.progress-fill.danger { background: var(--apple-red); }

.progress-text {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
  margin-top: var(--apple-spacing-xs);
  text-align: right;
}

.package-subscription {
  display: flex;
  justify-content: space-between;
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.expire-info {
  color: var(--apple-text-tertiary);
}

/* API Key 列表 */
.apikey-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.apikey-info {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xxs);
}

.apikey-name {
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.apikey-prefix {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  background: transparent;
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-xxl);
  color: var(--apple-text-tertiary);
}

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: var(--apple-spacing-md);
  opacity: 0.5;
}

.empty-action {
  margin-top: var(--apple-spacing-md);
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  background: var(--apple-blue);
  color: white;
  border-radius: var(--apple-radius-md);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.empty-action:hover {
  background: var(--apple-blue-hover);
}

/* 累计统计 */
.total-stats-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.total-stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  padding: var(--apple-spacing-xl);
  gap: var(--apple-spacing-lg);
}

.total-stat {
  text-align: center;
}

.total-value {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-blue);
}

.total-label {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-top: var(--apple-spacing-xs);
}

/* 响应式 */
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .content-grid {
    grid-template-columns: 1fr;
  }

  .total-stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 640px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .total-stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stat-card {
    padding: var(--apple-spacing-lg);
  }
}
</style>
