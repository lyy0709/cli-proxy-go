<!--
 * 文件作用：使用记录页面 - Apple HIG 风格
 * 负责功能：
 *   - 显示用户的请求历史
 *   - 按日期/模型筛选
 *   - Token 和费用统计
 * 重要程度：⭐⭐⭐⭐ 重要（用户核心功能）
-->
<template>
  <div class="records-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">使用记录</h1>
        <p class="page-subtitle">查看您的 API 请求历史和费用明细</p>
      </div>
    </div>

    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon requests">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M22 12h-4l-3 9L9 3l-3 9H2"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">请求次数</span>
          <span class="stat-value">{{ formatNumber(summary.request_count) }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon tokens">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <circle cx="12" cy="12" r="10"/>
            <polygon points="16.24,7.76 14.12,14.12 7.76,16.24 9.88,9.88"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">总 Token</span>
          <span class="stat-value">{{ formatLargeNumber(summary.total_tokens) }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon cost">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="1" x2="12" y2="23"/>
            <path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">总费用</span>
          <span class="stat-value">${{ (summary.total_cost || 0).toFixed(4) }}</span>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon models">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <rect x="3" y="3" width="7" height="7"/>
            <rect x="14" y="3" width="7" height="7"/>
            <rect x="14" y="14" width="7" height="7"/>
            <rect x="3" y="14" width="7" height="7"/>
          </svg>
        </div>
        <div class="stat-content">
          <span class="stat-label">使用模型数</span>
          <span class="stat-value">{{ summary.model_count || 0 }}</span>
        </div>
      </div>
    </div>

    <!-- 筛选条件 -->
    <div class="filter-card">
      <div class="filter-row">
        <div class="filter-group">
          <label class="filter-label">时间范围</label>
          <div class="date-inputs">
            <input
              type="date"
              v-model="filters.startDate"
              class="date-input"
            />
            <span class="date-separator">至</span>
            <input
              type="date"
              v-model="filters.endDate"
              class="date-input"
            />
          </div>
        </div>

        <div class="filter-group">
          <label class="filter-label">模型</label>
          <select v-model="filters.model" class="filter-select">
            <option value="">全部模型</option>
            <option v-for="m in models" :key="m" :value="m">{{ m }}</option>
          </select>
        </div>

        <div class="filter-actions">
          <button class="btn btn-primary" @click="fetchRecords">
            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
              <circle cx="11" cy="11" r="8"/>
              <line x1="21" y1="21" x2="16.65" y2="16.65"/>
            </svg>
            查询
          </button>
          <button class="btn btn-secondary" @click="resetFilters">重置</button>
        </div>
      </div>

      <div class="date-shortcuts">
        <button
          v-for="shortcut in dateShortcuts"
          :key="shortcut.text"
          class="shortcut-btn"
          @click="applyShortcut(shortcut)"
        >
          {{ shortcut.text }}
        </button>
      </div>
    </div>

    <!-- 记录列表 -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>时间</th>
              <th>模型</th>
              <th class="col-right">输入</th>
              <th class="col-right">输出</th>
              <th class="col-right">缓存创建</th>
              <th class="col-right">缓存读取</th>
              <th class="col-right">总 Token</th>
              <th class="col-right">费用</th>
              <th>IP</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="(row, idx) in records" :key="idx">
              <td>{{ formatTime(row.request_time || row.timestamp) }}</td>
              <td>
                <code class="model-name">{{ row.model }}</code>
              </td>
              <td class="col-right">{{ formatNumber(row.input_tokens) }}</td>
              <td class="col-right">{{ formatNumber(row.output_tokens) }}</td>
              <td class="col-right">
                <span :class="{ 'cache-create': row.cache_creation_input_tokens > 0 }">
                  {{ formatNumber(row.cache_creation_input_tokens) }}
                </span>
              </td>
              <td class="col-right">
                <span :class="{ 'cache-read': row.cache_read_input_tokens > 0 }">
                  {{ formatNumber(row.cache_read_input_tokens) }}
                </span>
              </td>
              <td class="col-right">{{ formatNumber(row.total_tokens) }}</td>
              <td class="col-right cost-cell">${{ (row.total_cost || 0).toFixed(4) }}</td>
              <td>
                <span class="ip-cell">{{ row.request_ip || '-' }}</span>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && records.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
            <polyline points="14,2 14,8 20,8"/>
            <line x1="16" y1="13" x2="8" y2="13"/>
            <line x1="16" y1="17" x2="8" y2="17"/>
          </svg>
          <span>暂无记录</span>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrap" v-if="pagination.total > 0">
        <div class="pagination-info">
          共 {{ pagination.total }} 条记录
        </div>
        <div class="pagination-controls">
          <select v-model="pagination.pageSize" @change="fetchRecords" class="page-size-select">
            <option :value="20">20 条/页</option>
            <option :value="50">50 条/页</option>
            <option :value="100">100 条/页</option>
          </select>
          <div class="page-btns">
            <button
              class="page-btn"
              :disabled="pagination.page <= 1"
              @click="pagination.page--; fetchRecords()"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15,18 9,12 15,6"/>
              </svg>
            </button>
            <span class="page-current">{{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) || 1 }}</span>
            <button
              class="page-btn"
              :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)"
              @click="pagination.page++; fetchRecords()"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9,18 15,12 9,6"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const records = ref([])
const models = ref([])

const filters = reactive({
  startDate: '',
  endDate: '',
  model: ''
})

const summary = reactive({
  request_count: 0,
  total_tokens: 0,
  total_cost: 0,
  model_count: 0
})

const pagination = reactive({
  page: 1,
  pageSize: 20,
  total: 0
})

const dateShortcuts = [
  { text: '今天', days: 0 },
  { text: '最近7天', days: 7 },
  { text: '最近30天', days: 30 },
  { text: '本月', type: 'month' }
]

function formatDateForInput(date) {
  return date.toISOString().split('T')[0]
}

function applyShortcut(shortcut) {
  const end = new Date()
  const start = new Date()

  if (shortcut.type === 'month') {
    start.setDate(1)
  } else {
    start.setDate(start.getDate() - shortcut.days)
  }

  filters.startDate = formatDateForInput(start)
  filters.endDate = formatDateForInput(end)
  fetchRecords()
}

const formatTime = (time) => {
  if (!time) return '-'
  return new Date(time).toLocaleString('zh-CN')
}

const formatNumber = (num) => {
  if (!num) return '0'
  return num.toLocaleString()
}

const formatLargeNumber = (num) => {
  if (!num) return '0'
  if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toLocaleString()
}

const fetchRecords = async () => {
  loading.value = true
  try {
    const params = {
      page: pagination.page,
      page_size: pagination.pageSize
    }

    if (filters.startDate) {
      params.start_date = filters.startDate
    }

    if (filters.endDate) {
      params.end_date = filters.endDate
    }

    if (filters.model) {
      params.model = filters.model
    }

    // 使用正确的用户 API（不是管理员 API）
    const res = await api.getMyUsageRecords(params)
    if (res.data) {
      records.value = res.data.items || []
      pagination.total = res.data.total || 0
    }
  } catch (e) {
    ElMessage.error('获取记录失败')
  } finally {
    loading.value = false
  }
}

const fetchSummary = async () => {
  try {
    // 使用正确的用户 API（不是管理员 API）
    const res = await api.getMyUsageSummary()
    if (res.data) {
      summary.total_cost = res.data.total_cost || 0
      summary.total_tokens = res.data.total_tokens || 0
      summary.request_count = res.data.total_requests || 0
      summary.model_count = res.data.model_count || 0
    }
  } catch (e) {
    console.error('Failed to fetch summary:', e)
  }
}

const fetchModels = async () => {
  try {
    // 使用正确的用户 API（不是管理员 API）
    const res = await api.getMyModelUsage()
    if (res.data && res.data.models) {
      models.value = res.data.models.map(m => m.model).filter(Boolean)
    }
  } catch (e) {
    console.error('Failed to fetch models:', e)
  }
}

const resetFilters = () => {
  filters.startDate = ''
  filters.endDate = ''
  filters.model = ''
  pagination.page = 1
  fetchRecords()
}

onMounted(() => {
  fetchRecords()
  fetchSummary()
  fetchModels()
})
</script>

<style scoped>
.records-page {
  max-width: 1400px;
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

/* 统计卡片网格 */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--apple-spacing-md);
  margin-bottom: var(--apple-spacing-xl);
}

.stat-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  padding: var(--apple-spacing-lg);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
}

.stat-icon {
  width: 48px;
  height: 48px;
  border-radius: var(--apple-radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-icon svg {
  width: 24px;
  height: 24px;
  color: white;
}

.stat-icon.requests { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
.stat-icon.tokens { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
.stat-icon.cost { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
.stat-icon.models { background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%); }

.stat-content {
  flex: 1;
  min-width: 0;
}

.stat-label {
  display: block;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-bottom: 4px;
}

.stat-value {
  font-size: var(--apple-text-xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
}

/* 筛选卡片 */
.filter-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  padding: var(--apple-spacing-lg);
  margin-bottom: var(--apple-spacing-xl);
}

.filter-row {
  display: flex;
  align-items: flex-end;
  gap: var(--apple-spacing-lg);
  flex-wrap: wrap;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
}

.filter-label {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
}

.date-inputs {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.date-input {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
}

.date-input:focus {
  outline: none;
  border-color: var(--apple-blue);
}

.date-separator {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
}

.filter-select {
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  min-width: 160px;
}

.filter-select:focus {
  outline: none;
  border-color: var(--apple-blue);
}

.filter-actions {
  display: flex;
  gap: var(--apple-spacing-sm);
}

.date-shortcuts {
  display: flex;
  gap: var(--apple-spacing-xs);
  margin-top: var(--apple-spacing-md);
  padding-top: var(--apple-spacing-md);
  border-top: 1px solid var(--apple-separator);
}

.shortcut-btn {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs);
  color: var(--apple-blue);
  background: var(--apple-blue-light);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast);
}

.shortcut-btn:hover {
  background: var(--apple-blue);
  color: white;
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

.btn-primary {
  background: var(--apple-blue);
  color: white;
}

.btn-primary:hover {
  background: var(--apple-blue-hover);
}

.btn-secondary {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-primary);
}

.btn-secondary:hover {
  background: var(--apple-fill-secondary);
}

/* 表格卡片 */
.table-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
}

.table-container {
  overflow-x: auto;
}

.data-table {
  width: 100%;
  border-collapse: collapse;
  font-size: var(--apple-text-sm);
}

.data-table th,
.data-table td {
  padding: var(--apple-spacing-md);
  text-align: left;
  border-bottom: 1px solid var(--apple-separator);
}

.data-table th {
  background: var(--apple-bg-secondary);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary);
  white-space: nowrap;
}

.data-table tbody tr:hover {
  background: var(--apple-bg-secondary);
}

.col-right {
  text-align: right;
}

.model-name {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-primary);
  background: var(--apple-bg-tertiary);
  padding: 2px 8px;
  border-radius: var(--apple-radius-xs);
}

.cache-create {
  color: var(--apple-orange);
  font-weight: var(--apple-font-semibold);
}

.cache-read {
  color: var(--apple-green);
  font-weight: var(--apple-font-semibold);
}

.cost-cell {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-blue);
}

.ip-cell {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

/* 分页 */
.pagination-wrap {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: var(--apple-spacing-md) var(--apple-spacing-lg);
  border-top: 1px solid var(--apple-separator);
}

.pagination-info {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.pagination-controls {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
}

.page-size-select {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-sm);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  background: var(--apple-bg-primary);
}

.page-btns {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.page-btn {
  width: 32px;
  height: 32px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-secondary);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.page-btn svg {
  width: 16px;
  height: 16px;
}

.page-btn:hover:not(:disabled) {
  background: var(--apple-blue);
  color: white;
}

.page-btn:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

.page-current {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  min-width: 60px;
  text-align: center;
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

@keyframes spin {
  to { transform: rotate(360deg); }
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
  width: 48px;
  height: 48px;
  margin-bottom: var(--apple-spacing-md);
  opacity: 0.5;
}

/* 响应式 */
@media (max-width: 1024px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .filter-row {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-group {
    width: 100%;
  }

  .date-inputs {
    width: 100%;
  }

  .date-input {
    flex: 1;
  }

  .filter-select {
    width: 100%;
  }

  .filter-actions {
    width: 100%;
  }

  .filter-actions .btn {
    flex: 1;
  }
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .date-shortcuts {
    flex-wrap: wrap;
  }
}
</style>
