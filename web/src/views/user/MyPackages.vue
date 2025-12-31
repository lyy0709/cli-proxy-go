<!--
 * 文件作用：我的套餐页面 - Apple HIG 风格
 * 负责功能：
 *   - 显示用户的所有套餐
 *   - 套餐详情和使用情况
 *   - 额度/订阅进度展示
 * 重要程度：⭐⭐⭐⭐ 重要（用户核心功能）
-->
<template>
  <div class="packages-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">我的套餐</h1>
        <p class="page-subtitle">查看和管理您的套餐使用情况</p>
      </div>
      <button class="btn btn-outline" @click="fetchPackages" :disabled="loading">
        <svg v-if="!loading" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <polyline points="23,4 23,10 17,10"/>
          <path d="M20.49 15a9 9 0 11-2.12-9.36L23 10"/>
        </svg>
        <span v-if="loading" class="btn-loading"></span>
        刷新
      </button>
    </div>

    <!-- 套餐卡片列表 -->
    <div v-if="loading" class="loading-state">
      <div class="loading-spinner"></div>
      <span>加载中...</span>
    </div>

    <div v-else-if="packages.length > 0" class="packages-grid">
      <div
        v-for="pkg in packages"
        :key="pkg.id"
        :class="['package-card', { 'expired': pkg.status !== 'active' }]"
      >
        <!-- 卡片头部 -->
        <div class="card-header">
          <div class="package-title">
            <span :class="['type-badge', pkg.type === 'subscription' ? 'badge-primary' : 'badge-success']">
              {{ pkg.type === 'subscription' ? '订阅' : '额度' }}
            </span>
            <span class="package-name">{{ pkg.name }}</span>
          </div>
          <span :class="['status-badge', getStatusClass(pkg.status)]">
            {{ getStatusText(pkg.status) }}
          </span>
        </div>

        <!-- 卡片内容 -->
        <div class="card-body">
          <!-- 额度套餐 -->
          <template v-if="pkg.type === 'quota'">
            <div class="quota-section">
              <div class="quota-header">
                <span class="quota-label">额度使用</span>
                <span class="quota-percentage">{{ getQuotaPercentage(pkg).toFixed(1) }}%</span>
              </div>
              <div class="progress-bar large">
                <div
                  class="progress-fill"
                  :style="{ width: getQuotaPercentage(pkg) + '%' }"
                  :class="getProgressClass(getQuotaPercentage(pkg))"
                ></div>
              </div>
              <div class="quota-detail">
                <span>已用: ${{ (pkg.quota_used || 0).toFixed(2) }}</span>
                <span>总额: ${{ (pkg.quota_total || 0).toFixed(2) }}</span>
              </div>
              <div class="quota-remaining">
                <span class="remaining-label">剩余额度</span>
                <span class="remaining-value">${{ ((pkg.quota_total || 0) - (pkg.quota_used || 0)).toFixed(2) }}</span>
              </div>
            </div>
          </template>

          <!-- 订阅套餐 -->
          <template v-else>
            <div class="subscription-section">
              <!-- 日额度 -->
              <div class="limit-item" v-if="pkg.daily_quota > 0">
                <div class="limit-header">
                  <span class="limit-label">日额度</span>
                  <span class="limit-value">${{ (pkg.daily_used || 0).toFixed(2) }} / ${{ pkg.daily_quota.toFixed(2) }}</span>
                </div>
                <div class="progress-bar">
                  <div
                    class="progress-fill"
                    :style="{ width: getUsagePercentage(pkg.daily_used, pkg.daily_quota) + '%' }"
                    :class="getProgressClass(getUsagePercentage(pkg.daily_used, pkg.daily_quota))"
                  ></div>
                </div>
              </div>

              <!-- 周额度 -->
              <div class="limit-item" v-if="pkg.weekly_quota > 0">
                <div class="limit-header">
                  <span class="limit-label">周额度</span>
                  <span class="limit-value">${{ (pkg.weekly_used || 0).toFixed(2) }} / ${{ pkg.weekly_quota.toFixed(2) }}</span>
                </div>
                <div class="progress-bar">
                  <div
                    class="progress-fill"
                    :style="{ width: getUsagePercentage(pkg.weekly_used, pkg.weekly_quota) + '%' }"
                    :class="getProgressClass(getUsagePercentage(pkg.weekly_used, pkg.weekly_quota))"
                  ></div>
                </div>
              </div>

              <!-- 月额度 -->
              <div class="limit-item" v-if="pkg.monthly_quota > 0">
                <div class="limit-header">
                  <span class="limit-label">月额度</span>
                  <span class="limit-value">${{ (pkg.monthly_used || 0).toFixed(2) }} / ${{ pkg.monthly_quota.toFixed(2) }}</span>
                </div>
                <div class="progress-bar">
                  <div
                    class="progress-fill"
                    :style="{ width: getUsagePercentage(pkg.monthly_used, pkg.monthly_quota) + '%' }"
                    :class="getProgressClass(getUsagePercentage(pkg.monthly_used, pkg.monthly_quota))"
                  ></div>
                </div>
              </div>

              <!-- 无限额 -->
              <div v-if="!pkg.daily_quota && !pkg.weekly_quota && !pkg.monthly_quota" class="no-limit">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <path d="M22 11.08V12a10 10 0 11-5.93-9.14"/>
                  <polyline points="22,4 12,14.01 9,11.01"/>
                </svg>
                <span>无额度限制</span>
              </div>
            </div>
          </template>
        </div>

        <!-- 卡片底部 -->
        <div class="card-footer">
          <div class="info-row">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
              <line x1="16" y1="2" x2="16" y2="6"/>
              <line x1="8" y1="2" x2="8" y2="6"/>
              <line x1="3" y1="10" x2="21" y2="10"/>
            </svg>
            <span>开始: {{ formatDate(pkg.start_time) }}</span>
          </div>
          <div class="info-row">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="12" cy="12" r="10"/>
              <polyline points="12,6 12,12 16,14"/>
            </svg>
            <span>到期: {{ formatDate(pkg.expire_time) }}</span>
          </div>
          <div class="info-row highlight" v-if="pkg.status === 'active'">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <path d="M12 2v10l4.5 4.5"/>
              <circle cx="12" cy="12" r="10"/>
            </svg>
            <span>剩余 <strong>{{ getDaysRemaining(pkg.expire_time) }}</strong> 天</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else class="empty-state">
      <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
        <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z"/>
        <line x1="3" y1="6" x2="21" y2="6"/>
        <path d="M16 10a4 4 0 01-8 0"/>
      </svg>
      <span>暂无套餐</span>
      <p class="empty-tip">您还没有任何套餐，请联系管理员开通</p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const packages = ref([])

const formatDate = (date) => {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('zh-CN')
}

const getStatusClass = (status) => {
  switch (status) {
    case 'active': return 'badge-success'
    case 'expired': return 'badge-secondary'
    case 'exhausted': return 'badge-warning'
    default: return 'badge-secondary'
  }
}

const getStatusText = (status) => {
  switch (status) {
    case 'active': return '有效'
    case 'expired': return '已过期'
    case 'exhausted': return '已耗尽'
    default: return status
  }
}

const getQuotaPercentage = (pkg) => {
  if (!pkg.quota_total || pkg.quota_total === 0) return 0
  return Math.min(100, (pkg.quota_used / pkg.quota_total) * 100)
}

const getUsagePercentage = (used, total) => {
  if (!total || total === 0) return 0
  return Math.min(100, ((used || 0) / total) * 100)
}

const getProgressClass = (percentage) => {
  if (percentage < 60) return 'progress-safe'
  if (percentage < 80) return 'progress-warning'
  return 'progress-danger'
}

const getDaysRemaining = (expireTime) => {
  if (!expireTime) return 0
  const now = new Date()
  const expire = new Date(expireTime)
  const diff = expire - now
  return Math.max(0, Math.ceil(diff / (1000 * 60 * 60 * 24)))
}

const fetchPackages = async () => {
  loading.value = true
  try {
    const res = await api.getMyPackages()
    packages.value = res.data || []
  } catch (e) {
    ElMessage.error('获取套餐列表失败')
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchPackages()
})
</script>

<style scoped>
.packages-page {
  max-width: 1200px;
  margin: 0 auto;
}

/* 页面标题 */
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--apple-spacing-xl);
}

.header-content {
  flex: 1;
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
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  cursor: pointer;
}

.btn svg {
  width: 16px;
  height: 16px;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-outline {
  background: transparent;
  color: var(--apple-blue);
  border: 1px solid var(--apple-blue);
}

.btn-outline:hover:not(:disabled) {
  background: var(--apple-blue-light);
}

.btn-loading {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(0, 122, 255, 0.3);
  border-top-color: var(--apple-blue);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* 套餐网格 */
.packages-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(360px, 1fr));
  gap: var(--apple-spacing-xl);
}

/* 套餐卡片 */
.package-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
  transition: all var(--apple-duration-normal) var(--apple-ease-default);
}

.package-card:hover {
  transform: translateY(-4px);
  box-shadow: var(--apple-shadow-elevated);
}

.package-card.expired {
  opacity: 0.7;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-lg);
  border-bottom: 1px solid var(--apple-separator);
}

.package-title {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.package-name {
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

/* 徽章 */
.type-badge, .status-badge {
  display: inline-block;
  padding: 3px 10px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.badge-primary { background: var(--apple-blue-light); color: var(--apple-blue); }
.badge-success { background: var(--apple-green-light); color: var(--apple-green); }
.badge-warning { background: var(--apple-orange-light); color: var(--apple-orange); }
.badge-secondary { background: var(--apple-fill-tertiary); color: var(--apple-text-tertiary); }

/* 卡片内容 */
.card-body {
  padding: var(--apple-spacing-lg);
}

/* 额度套餐样式 */
.quota-section {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-md);
}

.quota-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.quota-label {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.quota-percentage {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-bold);
  color: var(--apple-blue);
}

.quota-detail {
  display: flex;
  justify-content: space-between;
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
}

.quota-remaining {
  background: linear-gradient(135deg, rgba(52, 199, 89, 0.1) 0%, rgba(48, 209, 88, 0.1) 100%);
  border-radius: var(--apple-radius-md);
  padding: var(--apple-spacing-md);
  text-align: center;
}

.remaining-label {
  display: block;
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-bottom: var(--apple-spacing-xxs);
}

.remaining-value {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-green);
}

/* 订阅套餐样式 */
.subscription-section {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-lg);
}

.limit-item {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-sm);
}

.limit-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.limit-label {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.limit-value {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

/* 进度条 */
.progress-bar {
  height: 8px;
  background: var(--apple-fill-tertiary);
  border-radius: var(--apple-radius-full);
  overflow: hidden;
}

.progress-bar.large {
  height: 12px;
}

.progress-fill {
  height: 100%;
  border-radius: var(--apple-radius-full);
  transition: width var(--apple-duration-normal) var(--apple-ease-default);
}

.progress-safe { background: var(--apple-green); }
.progress-warning { background: var(--apple-orange); }
.progress-danger { background: var(--apple-red); }

/* 无限额 */
.no-limit {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-xl);
  background: linear-gradient(135deg, rgba(52, 199, 89, 0.1) 0%, rgba(48, 209, 88, 0.1) 100%);
  border-radius: var(--apple-radius-md);
  color: var(--apple-green);
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-medium);
}

.no-limit svg {
  width: 24px;
  height: 24px;
}

/* 卡片底部 */
.card-footer {
  padding: var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
  background: var(--apple-bg-secondary);
}

.info-row {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
  margin-bottom: var(--apple-spacing-sm);
}

.info-row:last-child {
  margin-bottom: 0;
}

.info-row svg {
  width: 16px;
  height: 16px;
  flex-shrink: 0;
}

.info-row.highlight {
  color: var(--apple-blue);
}

.info-row.highlight strong {
  font-weight: var(--apple-font-bold);
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-xxxl);
  color: var(--apple-text-tertiary);
}

.loading-spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--apple-fill-tertiary);
  border-top-color: var(--apple-blue);
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: var(--apple-spacing-md);
}

/* 空状态 */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-xxxl);
  color: var(--apple-text-tertiary);
}

.empty-state svg {
  width: 64px;
  height: 64px;
  margin-bottom: var(--apple-spacing-lg);
  opacity: 0.5;
}

.empty-state span {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
}

.empty-tip {
  margin-top: var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
}

/* 响应式 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-lg);
  }

  .packages-grid {
    grid-template-columns: 1fr;
  }
}
</style>
