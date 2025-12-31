<!--
 * 文件作用：我的 API Key 页面 - Apple HIG 风格
 * 负责功能：
 *   - 查看用户的 API Key 列表
 *   - 创建新的 API Key
 *   - 删除 API Key
 *   - 查看 API Key 使用情况
 * 重要程度：⭐⭐⭐⭐ 重要（用户核心功能）
-->
<template>
  <div class="my-api-keys">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">API Key 管理</h1>
        <p class="page-subtitle">创建和管理您的 API Key，用于访问 AI 服务</p>
      </div>
      <button class="create-btn" @click="showCreateDialog = true">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        创建 API Key
      </button>
    </div>

    <!-- API Key 列表 -->
    <div class="keys-container">
      <!-- 加载状态 -->
      <div v-if="loading" class="loading-state">
        <div class="loading-spinner"></div>
        <span>加载中...</span>
      </div>

      <!-- 空状态 -->
      <div v-else-if="apiKeys.length === 0" class="empty-state">
        <div class="empty-icon">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
          </svg>
        </div>
        <h3>暂无 API Key</h3>
        <p>创建您的第一个 API Key 开始使用 AI 服务</p>
        <button class="empty-action" @click="showCreateDialog = true">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          创建 API Key
        </button>
      </div>

      <!-- Key 列表 -->
      <div v-else class="keys-list">
        <div v-for="key in apiKeys" :key="key.id" class="key-card">
          <div class="key-main">
            <div class="key-header">
              <div class="key-info">
                <h3 class="key-name">{{ key.name }}</h3>
                <div class="key-meta">
                  <code class="key-prefix">{{ key.key_prefix }}...</code>
                  <span :class="['status-badge', key.status]">
                    {{ key.status === 'active' ? '正常' : '禁用' }}
                  </span>
                </div>
              </div>
              <div class="key-actions">
                <button class="action-btn view" @click="viewUsage(key)" title="查看详情">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="12" y1="16" x2="12" y2="12"/>
                    <line x1="12" y1="8" x2="12.01" y2="8"/>
                  </svg>
                </button>
                <button class="action-btn delete" @click="confirmDelete(key)" title="删除">
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                    <polyline points="3,6 5,6 21,6"/>
                    <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                  </svg>
                </button>
              </div>
            </div>

            <div class="key-stats">
              <div class="stat-item">
                <span class="stat-label">请求次数</span>
                <span class="stat-value">{{ formatNumber(key.request_count) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">Token 使用</span>
                <span class="stat-value">{{ formatLargeNumber(key.tokens_used) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">累计消费</span>
                <span class="stat-value cost">${{ (key.cost_used || 0).toFixed(4) }}</span>
              </div>
              <div class="stat-item">
                <span class="stat-label">绑定套餐</span>
                <span class="stat-value">
                  <span v-if="key.user_package" class="package-badge">{{ key.user_package.name }}</span>
                  <span v-else class="no-package">未绑定</span>
                </span>
              </div>
            </div>

            <div class="key-footer">
              <span class="create-time">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <circle cx="12" cy="12" r="10"/>
                  <polyline points="12,6 12,12 16,14"/>
                </svg>
                创建于 {{ formatTime(key.created_at) }}
              </span>
              <span v-if="key.last_used_at" class="last-used">
                最后使用: {{ formatTime(key.last_used_at) }}
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 创建 API Key 弹窗 -->
    <Teleport to="body">
      <div v-if="showCreateDialog" class="modal-overlay" @click.self="showCreateDialog = false">
        <div class="modal">
          <div class="modal-header">
            <h2>创建 API Key</h2>
            <button class="modal-close" @click="showCreateDialog = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">名称 <span class="required">*</span></label>
              <input
                v-model="createForm.name"
                type="text"
                class="form-input"
                placeholder="请输入名称，如：开发测试"
              />
            </div>
            <div class="form-group">
              <label class="form-label">绑定套餐</label>
              <select v-model="createForm.user_package_id" class="form-select">
                <option :value="null">不绑定套餐</option>
                <option
                  v-for="pkg in activePackages"
                  :key="pkg.id"
                  :value="pkg.id"
                >
                  {{ pkg.name }} ({{ pkg.type === 'subscription' ? '订阅' : '额度' }})
                </option>
              </select>
              <p class="form-tip">不绑定套餐则使用默认计费</p>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="showCreateDialog = false">取消</button>
            <button class="btn btn-primary" :disabled="creating || !createForm.name" @click="createKey">
              <span v-if="creating" class="btn-loading"></span>
              {{ creating ? '创建中...' : '创建' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 新创建的 Key 显示弹窗 -->
    <Teleport to="body">
      <div v-if="showNewKeyDialog" class="modal-overlay">
        <div class="modal">
          <div class="modal-header success">
            <div class="success-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="20,6 9,17 4,12"/>
              </svg>
            </div>
            <h2>API Key 创建成功</h2>
          </div>
          <div class="modal-body">
            <div class="warning-alert">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <path d="M10.29 3.86L1.82 18a2 2 0 001.71 3h16.94a2 2 0 001.71-3L13.71 3.86a2 2 0 00-3.42 0z"/>
                <line x1="12" y1="9" x2="12" y2="13"/>
                <line x1="12" y1="17" x2="12.01" y2="17"/>
              </svg>
              <span>请立即复制并保存此 API Key，关闭后将无法再次查看完整内容！</span>
            </div>
            <div class="key-display">
              <code>{{ newKeyValue }}</code>
              <button class="copy-btn" @click="copyKey(newKeyValue)">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
                复制
              </button>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-primary" @click="showNewKeyDialog = false">我已保存</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 使用详情弹窗 -->
    <Teleport to="body">
      <div v-if="showUsageDialog" class="modal-overlay" @click.self="showUsageDialog = false">
        <div class="modal modal-lg">
          <div class="modal-header">
            <h2>{{ selectedKey?.name }} - 使用详情</h2>
            <button class="modal-close" @click="showUsageDialog = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body" v-if="selectedKey">
            <div class="detail-grid">
              <div class="detail-item">
                <span class="detail-label">名称</span>
                <span class="detail-value">{{ selectedKey.name }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">状态</span>
                <span :class="['status-badge', selectedKey.status]">
                  {{ selectedKey.status === 'active' ? '正常' : '禁用' }}
                </span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Key 前缀</span>
                <code class="detail-code">{{ selectedKey.key_prefix }}...</code>
              </div>
              <div class="detail-item">
                <span class="detail-label">创建时间</span>
                <span class="detail-value">{{ formatTime(selectedKey.created_at) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">请求次数</span>
                <span class="detail-value highlight">{{ formatNumber(selectedKey.request_count) }}</span>
              </div>
              <div class="detail-item">
                <span class="detail-label">Token 使用</span>
                <span class="detail-value highlight">{{ formatLargeNumber(selectedKey.tokens_used) }}</span>
              </div>
              <div class="detail-item full-width">
                <span class="detail-label">累计消费</span>
                <span class="detail-value cost-large">${{ (selectedKey.cost_used || 0).toFixed(4) }}</span>
              </div>
              <div class="detail-item full-width">
                <span class="detail-label">最后使用</span>
                <span class="detail-value">{{ selectedKey.last_used_at ? formatTime(selectedKey.last_used_at) : '从未使用' }}</span>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="showUsageDialog = false">关闭</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 删除确认弹窗 -->
    <Teleport to="body">
      <div v-if="showDeleteDialog" class="modal-overlay" @click.self="showDeleteDialog = false">
        <div class="modal modal-sm">
          <div class="modal-header danger">
            <div class="danger-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <circle cx="12" cy="12" r="10"/>
                <line x1="15" y1="9" x2="9" y2="15"/>
                <line x1="9" y1="9" x2="15" y2="15"/>
              </svg>
            </div>
            <h2>确认删除</h2>
          </div>
          <div class="modal-body">
            <p class="delete-message">确定要删除 API Key "{{ deleteTarget?.name }}" 吗？此操作不可撤销。</p>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="showDeleteDialog = false">取消</button>
            <button class="btn btn-danger" :disabled="deleting" @click="deleteKey">
              <span v-if="deleting" class="btn-loading"></span>
              {{ deleting ? '删除中...' : '删除' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const creating = ref(false)
const deleting = ref(false)
const apiKeys = ref([])
const activePackages = ref([])

const showCreateDialog = ref(false)
const showNewKeyDialog = ref(false)
const showUsageDialog = ref(false)
const showDeleteDialog = ref(false)

const newKeyValue = ref('')
const selectedKey = ref(null)
const deleteTarget = ref(null)

const createForm = ref({
  name: '',
  user_package_id: null
})

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

const copyKey = async (key) => {
  try {
    await navigator.clipboard.writeText(key)
    ElMessage.success('已复制到剪贴板')
  } catch (e) {
    ElMessage.error('复制失败')
  }
}

const fetchApiKeys = async () => {
  loading.value = true
  try {
    const res = await api.getAPIKeys()
    apiKeys.value = res.data || []
  } catch (e) {
    ElMessage.error('获取 API Key 列表失败')
  } finally {
    loading.value = false
  }
}

const fetchPackages = async () => {
  try {
    const res = await api.getMyActivePackages()
    activePackages.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch packages:', e)
  }
}

const createKey = async () => {
  if (!createForm.value.name) {
    ElMessage.warning('请输入名称')
    return
  }

  creating.value = true
  try {
    const res = await api.createAPIKey(createForm.value)
    if (res.code === 0 && res.data) {
      newKeyValue.value = res.data.key
      showCreateDialog.value = false
      showNewKeyDialog.value = true
      createForm.value = { name: '', user_package_id: null }
      fetchApiKeys()
    } else {
      ElMessage.error(res.message || '创建失败')
    }
  } catch (e) {
    ElMessage.error('创建失败')
  } finally {
    creating.value = false
  }
}

const confirmDelete = (key) => {
  deleteTarget.value = key
  showDeleteDialog.value = true
}

const deleteKey = async () => {
  if (!deleteTarget.value) return

  deleting.value = true
  try {
    await api.deleteAPIKey(deleteTarget.value.id)
    ElMessage.success('删除成功')
    showDeleteDialog.value = false
    deleteTarget.value = null
    fetchApiKeys()
  } catch (e) {
    ElMessage.error('删除失败')
  } finally {
    deleting.value = false
  }
}

const viewUsage = (key) => {
  selectedKey.value = key
  showUsageDialog.value = true
}

onMounted(() => {
  fetchApiKeys()
  fetchPackages()
})
</script>

<style scoped>
.my-api-keys {
  max-width: 1000px;
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

.create-btn {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  background: var(--apple-blue);
  color: white;
  border-radius: var(--apple-radius-md);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.create-btn svg {
  width: 18px;
  height: 18px;
}

.create-btn:hover {
  background: var(--apple-blue-hover);
  transform: translateY(-1px);
}

/* 加载状态 */
.loading-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-xxxl);
  color: var(--apple-text-secondary);
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
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  box-shadow: var(--apple-shadow-card);
}

.empty-icon {
  width: 80px;
  height: 80px;
  background: var(--apple-fill-quaternary);
  border-radius: var(--apple-radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: var(--apple-spacing-lg);
}

.empty-icon svg {
  width: 40px;
  height: 40px;
  color: var(--apple-text-tertiary);
}

.empty-state h3 {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-xs) 0;
}

.empty-state p {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  margin: 0 0 var(--apple-spacing-xl) 0;
}

.empty-action {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  background: var(--apple-blue);
  color: white;
  border-radius: var(--apple-radius-md);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.empty-action svg {
  width: 18px;
  height: 18px;
}

.empty-action:hover {
  background: var(--apple-blue-hover);
}

/* Key 列表 */
.keys-list {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-md);
}

.key-card {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-lg);
  box-shadow: var(--apple-shadow-card);
  overflow: hidden;
  transition: all var(--apple-duration-normal) var(--apple-ease-default);
}

.key-card:hover {
  box-shadow: var(--apple-shadow-card-hover);
  transform: translateY(-2px);
}

.key-main {
  padding: var(--apple-spacing-xl);
}

.key-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: var(--apple-spacing-lg);
}

.key-info {
  flex: 1;
}

.key-name {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-xs) 0;
}

.key-meta {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.key-prefix {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  background: var(--apple-bg-tertiary);
  padding: 2px 8px;
  border-radius: var(--apple-radius-xs);
}

.status-badge {
  padding: 2px 10px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.status-badge.active {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.status-badge.disabled {
  background: var(--apple-red-light);
  color: var(--apple-red);
}

.key-actions {
  display: flex;
  gap: var(--apple-spacing-xs);
}

.action-btn {
  width: 36px;
  height: 36px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.action-btn svg {
  width: 18px;
  height: 18px;
}

.action-btn.view {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

.action-btn.view:hover {
  background: var(--apple-blue);
  color: white;
}

.action-btn.delete {
  background: var(--apple-red-light);
  color: var(--apple-red);
}

.action-btn.delete:hover {
  background: var(--apple-red);
  color: white;
}

/* Key 统计 */
.key-stats {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--apple-spacing-md);
  padding: var(--apple-spacing-lg);
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
  margin-bottom: var(--apple-spacing-lg);
}

.stat-item {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xxs);
}

.stat-label {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

.stat-value {
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.stat-value.cost {
  color: var(--apple-blue);
}

.package-badge {
  display: inline-block;
  padding: 2px 8px;
  background: var(--apple-blue-light);
  color: var(--apple-blue);
  border-radius: var(--apple-radius-xs);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.no-package {
  color: var(--apple-text-tertiary);
  font-size: var(--apple-text-sm);
}

/* Key 页脚 */
.key-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

.create-time {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
}

.create-time svg {
  width: 14px;
  height: 14px;
}

/* 模态框 */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: var(--apple-z-modal);
  padding: var(--apple-spacing-xl);
}

.modal {
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-xl);
  box-shadow: var(--apple-shadow-modal);
  width: 100%;
  max-width: 480px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal.modal-lg {
  max-width: 600px;
}

.modal.modal-sm {
  max-width: 400px;
}

.modal-header {
  padding: var(--apple-spacing-xl);
  border-bottom: 1px solid var(--apple-separator);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header.success {
  flex-direction: column;
  text-align: center;
  gap: var(--apple-spacing-md);
}

.modal-header.danger {
  flex-direction: column;
  text-align: center;
  gap: var(--apple-spacing-md);
}

.success-icon {
  width: 56px;
  height: 56px;
  background: var(--apple-green-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.success-icon svg {
  width: 28px;
  height: 28px;
  color: var(--apple-green);
}

.danger-icon {
  width: 56px;
  height: 56px;
  background: var(--apple-red-light);
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.danger-icon svg {
  width: 28px;
  height: 28px;
  color: var(--apple-red);
}

.modal-header h2 {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0;
}

.modal-close {
  width: 32px;
  height: 32px;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.modal-close svg {
  width: 18px;
  height: 18px;
}

.modal-close:hover {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-primary);
}

.modal-body {
  padding: var(--apple-spacing-xl);
  overflow-y: auto;
  flex: 1;
}

.modal-footer {
  padding: var(--apple-spacing-lg) var(--apple-spacing-xl);
  border-top: 1px solid var(--apple-separator);
  display: flex;
  justify-content: flex-end;
  gap: var(--apple-spacing-sm);
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

.form-input,
.form-select {
  width: 100%;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-input:focus,
.form-select:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.form-tip {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: var(--apple-spacing-xxs);
}

/* 按钮 */
.btn {
  padding: var(--apple-spacing-sm) var(--apple-spacing-lg);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
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

.btn-secondary {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-primary);
}

.btn-secondary:hover:not(:disabled) {
  background: var(--apple-fill-secondary);
}

.btn-danger {
  background: var(--apple-red);
  color: white;
}

.btn-danger:hover:not(:disabled) {
  background: #e6362d;
}

.btn-loading {
  width: 16px;
  height: 16px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

/* 警告提示 */
.warning-alert {
  display: flex;
  align-items: flex-start;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-md);
  background: var(--apple-orange-light);
  border-radius: var(--apple-radius-md);
  margin-bottom: var(--apple-spacing-lg);
}

.warning-alert svg {
  width: 20px;
  height: 20px;
  color: var(--apple-orange);
  flex-shrink: 0;
  margin-top: 2px;
}

.warning-alert span {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  line-height: 1.5;
}

/* Key 显示 */
.key-display {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-md);
  background: var(--apple-bg-tertiary);
  border-radius: var(--apple-radius-md);
}

.key-display code {
  flex: 1;
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  word-break: break-all;
}

.copy-btn {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xxs);
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  background: var(--apple-blue);
  color: white;
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
  white-space: nowrap;
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.copy-btn svg {
  width: 14px;
  height: 14px;
}

.copy-btn:hover {
  background: var(--apple-blue-hover);
}

/* 详情网格 */
.detail-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-lg);
}

.detail-item {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xxs);
}

.detail-item.full-width {
  grid-column: span 2;
}

.detail-label {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.detail-value {
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
}

.detail-value.highlight {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-blue);
}

.detail-code {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  background: var(--apple-bg-tertiary);
  padding: 4px 8px;
  border-radius: var(--apple-radius-xs);
}

.cost-large {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-blue);
}

/* 删除确认 */
.delete-message {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  text-align: center;
  margin: 0;
}

/* 响应式 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-lg);
  }

  .create-btn {
    width: 100%;
    justify-content: center;
  }

  .key-stats {
    grid-template-columns: repeat(2, 1fr);
  }

  .key-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--apple-spacing-xs);
  }

  .detail-grid {
    grid-template-columns: 1fr;
  }

  .detail-item.full-width {
    grid-column: span 1;
  }
}
</style>
