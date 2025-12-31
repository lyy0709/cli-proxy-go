<!--
 * 文件作用：套餐管理页面 - Apple HIG 风格
 * 负责功能：
 *   - 套餐列表和CRUD
 *   - 订阅类型配置（日/周/月额度）
 *   - 额度类型配置
 *   - 模型限制配置
 * 重要程度：⭐⭐⭐⭐ 重要（套餐配置）
-->
<template>
  <div class="packages-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">套餐管理</h1>
        <p class="page-subtitle">管理订阅和额度套餐配置</p>
      </div>
      <button class="btn btn-primary" @click="showCreateDialog">
        <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
          <line x1="12" y1="5" x2="12" y2="19"/>
          <line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        创建套餐
      </button>
    </div>

    <!-- 套餐列表 -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th>ID</th>
              <th>名称</th>
              <th>类型</th>
              <th>价格</th>
              <th>有效期</th>
              <th>额度限制</th>
              <th>模型限制</th>
              <th>状态</th>
              <th>操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="pkg in packages" :key="pkg.id">
              <td><code class="id-code">#{{ pkg.id }}</code></td>
              <td>
                <span class="pkg-name">{{ pkg.name }}</span>
              </td>
              <td>
                <span :class="['type-badge', pkg.type]">
                  {{ pkg.type === 'subscription' ? '订阅' : '额度' }}
                </span>
              </td>
              <td>
                <span class="price">${{ pkg.price.toFixed(2) }}</span>
              </td>
              <td>{{ pkg.duration }}天</td>
              <td>
                <div v-if="pkg.type === 'subscription'" class="quota-tags">
                  <span v-if="pkg.daily_quota > 0" class="quota-tag">日: ${{ pkg.daily_quota }}</span>
                  <span v-if="pkg.weekly_quota > 0" class="quota-tag">周: ${{ pkg.weekly_quota }}</span>
                  <span v-if="pkg.monthly_quota > 0" class="quota-tag">月: ${{ pkg.monthly_quota }}</span>
                  <span v-if="!pkg.daily_quota && !pkg.weekly_quota && !pkg.monthly_quota" class="text-muted">无限制</span>
                </div>
                <span v-else>总额度: ${{ pkg.quota_amount }}</span>
              </td>
              <td>
                <span v-if="pkg.allowed_models" class="model-badge">
                  {{ pkg.allowed_models.split(',').length }}个模型
                </span>
                <span v-else class="text-muted">全部</span>
              </td>
              <td>
                <span :class="['status-badge', pkg.status === 'active' ? 'success' : 'muted']">
                  {{ pkg.status === 'active' ? '启用' : '禁用' }}
                </span>
              </td>
              <td>
                <div class="action-group">
                  <button class="action-btn" @click="handleEdit(pkg)" title="编辑">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                      <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                  </button>
                  <button class="action-btn danger" @click="confirmDelete(pkg)" title="删除">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <polyline points="3,6 5,6 21,6"/>
                      <path d="M19 6v14a2 2 0 01-2 2H7a2 2 0 01-2-2V6m3 0V4a2 2 0 012-2h4a2 2 0 012 2v2"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <div v-if="!loading && packages.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <rect x="3" y="3" width="18" height="18" rx="2" ry="2"/>
            <line x1="3" y1="9" x2="21" y2="9"/>
            <line x1="9" y1="21" x2="9" y2="9"/>
          </svg>
          <span>暂无套餐配置</span>
        </div>
      </div>
    </div>

    <!-- 创建/编辑弹窗 -->
    <Teleport to="body">
      <div v-if="dialogVisible" class="modal-overlay" @click.self="dialogVisible = false">
        <div class="modal modal-lg">
          <div class="modal-header">
            <h2>{{ editMode ? '编辑套餐' : '创建套餐' }}</h2>
            <button class="modal-close" @click="dialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">名称 <span class="required">*</span></label>
              <input v-model="form.name" type="text" class="form-input" placeholder="套餐名称" />
            </div>

            <div class="form-group">
              <label class="form-label">类型 <span class="required">*</span></label>
              <div class="type-selector">
                <label :class="['type-option', { active: form.type === 'subscription', disabled: editMode }]">
                  <input type="radio" v-model="form.type" value="subscription" :disabled="editMode" />
                  <div class="type-content">
                    <span class="type-title">订阅 (包月)</span>
                    <span class="type-desc">按周期限制额度</span>
                  </div>
                </label>
                <label :class="['type-option', { active: form.type === 'quota', disabled: editMode }]">
                  <input type="radio" v-model="form.type" value="quota" :disabled="editMode" />
                  <div class="type-content">
                    <span class="type-title">额度</span>
                    <span class="type-desc">一次性消耗额度</span>
                  </div>
                </label>
              </div>
            </div>

            <div class="form-row">
              <div class="form-group">
                <label class="form-label">价格 ($) <span class="required">*</span></label>
                <input v-model.number="form.price" type="number" min="0" step="0.01" class="form-input" />
              </div>
              <div class="form-group">
                <label class="form-label">有效期 (天) <span class="required">*</span></label>
                <input v-model.number="form.duration" type="number" min="1" class="form-input" />
              </div>
            </div>

            <!-- 订阅类型额度限制 -->
            <template v-if="form.type === 'subscription'">
              <div class="form-section">
                <span class="section-title">周期额度限制 (0=不限)</span>
              </div>
              <div class="form-row three-col">
                <div class="form-group">
                  <label class="form-label">每日额度</label>
                  <div class="input-with-unit">
                    <input v-model.number="form.daily_quota" type="number" min="0" class="form-input" />
                    <span class="unit">$</span>
                  </div>
                </div>
                <div class="form-group">
                  <label class="form-label">每周额度</label>
                  <div class="input-with-unit">
                    <input v-model.number="form.weekly_quota" type="number" min="0" class="form-input" />
                    <span class="unit">$</span>
                  </div>
                </div>
                <div class="form-group">
                  <label class="form-label">每月额度</label>
                  <div class="input-with-unit">
                    <input v-model.number="form.monthly_quota" type="number" min="0" class="form-input" />
                    <span class="unit">$</span>
                  </div>
                </div>
              </div>
            </template>

            <!-- 额度类型总额度 -->
            <template v-if="form.type === 'quota'">
              <div class="form-group">
                <label class="form-label">总额度 <span class="required">*</span></label>
                <div class="input-with-unit">
                  <input v-model.number="form.quota_amount" type="number" min="0" class="form-input" />
                  <span class="unit">$</span>
                </div>
              </div>
            </template>

            <div class="form-group">
              <label class="form-label">允许的模型</label>
              <div class="model-select-wrap">
                <div class="model-chips">
                  <span v-for="model in selectedModels" :key="model" class="model-chip">
                    {{ model }}
                    <button @click="removeModel(model)" class="chip-remove">×</button>
                  </span>
                  <span v-if="selectedModels.length === 0" class="text-muted">留空表示全部模型</span>
                </div>
                <select v-model="modelToAdd" @change="addModel" class="model-select">
                  <option value="">选择模型...</option>
                  <option v-for="model in availableModels" :key="model.id" :value="model.name">
                    {{ model.name }}
                  </option>
                </select>
              </div>
              <p class="form-tip">限制该套餐可使用的模型列表</p>
            </div>

            <div class="form-group">
              <label class="form-label">状态</label>
              <div class="status-selector">
                <label :class="['status-option', { active: form.status === 'active' }]">
                  <input type="radio" v-model="form.status" value="active" />
                  <span>启用</span>
                </label>
                <label :class="['status-option', { active: form.status === 'disabled' }]">
                  <input type="radio" v-model="form.status" value="disabled" />
                  <span>禁用</span>
                </label>
              </div>
            </div>

            <div class="form-group">
              <label class="form-label">描述</label>
              <textarea v-model="form.description" class="form-textarea" rows="2" placeholder="套餐描述"></textarea>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="dialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="submitting" @click="handleSubmit">
              <span v-if="submitting" class="btn-loading"></span>
              {{ submitting ? '保存中...' : (editMode ? '保存' : '创建') }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 删除确认弹窗 -->
    <Teleport to="body">
      <div v-if="deleteDialogVisible" class="modal-overlay" @click.self="deleteDialogVisible = false">
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
            <p class="delete-message">确定要删除套餐 "{{ deleteTarget?.name }}" 吗？</p>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="deleteDialogVisible = false">取消</button>
            <button class="btn btn-danger" :disabled="deleting" @click="handleDelete">
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
import { ref, computed, onMounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const packages = ref([])
const modelList = ref([])

const dialogVisible = ref(false)
const editMode = ref(false)
const submitting = ref(false)

const deleteDialogVisible = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

const modelToAdd = ref('')

const form = ref({
  id: 0,
  name: '',
  type: 'subscription',
  price: 0,
  duration: 30,
  daily_quota: 0,
  weekly_quota: 0,
  monthly_quota: 0,
  quota_amount: 0,
  allowed_models: '',
  status: 'active',
  description: ''
})

const selectedModels = computed({
  get() {
    if (!form.value.allowed_models) return []
    return form.value.allowed_models.split(',').filter(m => m.trim())
  },
  set(val) {
    form.value.allowed_models = val.join(',')
  }
})

const availableModels = computed(() => {
  return modelList.value.filter(m => !selectedModels.value.includes(m.name))
})

function addModel() {
  if (modelToAdd.value && !selectedModels.value.includes(modelToAdd.value)) {
    selectedModels.value = [...selectedModels.value, modelToAdd.value]
  }
  modelToAdd.value = ''
}

function removeModel(model) {
  selectedModels.value = selectedModels.value.filter(m => m !== model)
}

async function fetchPackages() {
  loading.value = true
  try {
    const res = await api.getPackages()
    packages.value = res.data || []
  } catch (e) {
    console.error('Failed to fetch packages:', e)
  } finally {
    loading.value = false
  }
}

async function fetchModels() {
  try {
    const res = await api.getModels()
    modelList.value = (res.data || []).filter(m => m.enabled)
  } catch (e) {
    console.error('Failed to fetch models:', e)
  }
}

function showCreateDialog() {
  editMode.value = false
  form.value = {
    id: 0,
    name: '',
    type: 'subscription',
    price: 0,
    duration: 30,
    daily_quota: 0,
    weekly_quota: 0,
    monthly_quota: 0,
    quota_amount: 0,
    allowed_models: '',
    status: 'active',
    description: ''
  }
  dialogVisible.value = true
}

function handleEdit(row) {
  editMode.value = true
  form.value = { ...row }
  dialogVisible.value = true
}

async function handleSubmit() {
  if (!form.value.name) {
    ElMessage.warning('请输入名称')
    return
  }

  submitting.value = true
  try {
    const data = {
      ...form.value,
      price: parseFloat(form.value.price) || 0,
      duration: parseInt(form.value.duration) || 30,
      daily_quota: parseFloat(form.value.daily_quota) || 0,
      weekly_quota: parseFloat(form.value.weekly_quota) || 0,
      monthly_quota: parseFloat(form.value.monthly_quota) || 0,
      quota_amount: parseFloat(form.value.quota_amount) || 0
    }

    if (editMode.value) {
      await api.updatePackage(form.value.id, data)
      ElMessage.success('更新成功')
    } else {
      await api.createPackage(data)
      ElMessage.success('创建成功')
    }
    dialogVisible.value = false
    fetchPackages()
  } catch (e) {
    ElMessage.error(e.message || '操作失败')
  } finally {
    submitting.value = false
  }
}

function confirmDelete(pkg) {
  deleteTarget.value = pkg
  deleteDialogVisible.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.deletePackage(deleteTarget.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    deleteTarget.value = null
    fetchPackages()
  } catch (e) {
    ElMessage.error(e.message || '删除失败')
  } finally {
    deleting.value = false
  }
}

onMounted(() => {
  fetchPackages()
  fetchModels()
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

.id-code {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
}

.pkg-name {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.price {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-green);
}

/* 类型徽章 */
.type-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
}

.type-badge.subscription {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

.type-badge.quota {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

/* 额度标签 */
.quota-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--apple-spacing-xxs);
}

.quota-tag {
  display: inline-block;
  padding: 2px 6px;
  background: var(--apple-fill-tertiary);
  border-radius: var(--apple-radius-xs);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

.model-badge {
  display: inline-block;
  padding: 4px 10px;
  background: var(--apple-orange-light);
  color: var(--apple-orange);
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

/* 状态徽章 */
.status-badge {
  display: inline-block;
  padding: 4px 10px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.status-badge.success {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.status-badge.muted {
  background: var(--apple-fill-tertiary);
  color: var(--apple-text-tertiary);
}

.text-muted {
  color: var(--apple-text-tertiary);
  font-size: var(--apple-text-xs);
}

/* 操作按钮 */
.action-group {
  display: flex;
  gap: var(--apple-spacing-xxs);
}

.action-btn {
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

.action-btn svg {
  width: 16px;
  height: 16px;
}

.action-btn:hover {
  background: var(--apple-blue);
  color: white;
}

.action-btn.danger:hover {
  background: var(--apple-red);
}

/* 加载和空状态 */
.loading-state, .empty-state {
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

.empty-state svg {
  width: 48px;
  height: 48px;
  margin-bottom: var(--apple-spacing-md);
  opacity: 0.5;
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
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  border-top-color: white;
  border-radius: 50%;
  animation: spin 1s linear infinite;
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
  max-width: 500px;
  max-height: 90vh;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.modal.modal-sm { max-width: 400px; }
.modal.modal-lg { max-width: 600px; }

.modal-header {
  padding: var(--apple-spacing-xl);
  border-bottom: 1px solid var(--apple-separator);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header.danger {
  flex-direction: column;
  text-align: center;
  gap: var(--apple-spacing-md);
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

.delete-message {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  text-align: center;
  margin: 0;
}

/* 表单 */
.form-group {
  margin-bottom: var(--apple-spacing-lg);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-md);
}

.form-row.three-col {
  grid-template-columns: repeat(3, 1fr);
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

.form-input, .form-textarea, .model-select {
  width: 100%;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-textarea {
  resize: vertical;
  min-height: 60px;
}

.form-input:focus, .form-textarea:focus, .model-select:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.form-tip {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: var(--apple-spacing-xxs);
}

.form-section {
  margin: var(--apple-spacing-lg) 0 var(--apple-spacing-md) 0;
  padding-top: var(--apple-spacing-md);
  border-top: 1px solid var(--apple-separator);
}

.section-title {
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* 输入框带单位 */
.input-with-unit {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
}

.input-with-unit .form-input {
  flex: 1;
}

.unit {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-tertiary);
}

/* 类型选择器 */
.type-selector {
  display: flex;
  gap: var(--apple-spacing-md);
}

.type-option {
  flex: 1;
  display: flex;
  align-items: flex-start;
  padding: var(--apple-spacing-md);
  border: 2px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-md);
  cursor: pointer;
  transition: all var(--apple-duration-fast);
}

.type-option input {
  display: none;
}

.type-option:hover:not(.disabled) {
  border-color: var(--apple-blue);
}

.type-option.active {
  background: var(--apple-blue-light);
  border-color: var(--apple-blue);
}

.type-option.disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.type-content {
  display: flex;
  flex-direction: column;
}

.type-title {
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
}

.type-desc {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: 2px;
}

/* 状态选择器 */
.status-selector {
  display: flex;
  gap: var(--apple-spacing-xs);
}

.status-option {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  cursor: pointer;
  transition: all var(--apple-duration-fast);
}

.status-option input {
  display: none;
}

.status-option:hover {
  border-color: var(--apple-blue);
  color: var(--apple-blue);
}

.status-option.active {
  background: var(--apple-green);
  border-color: var(--apple-green);
  color: white;
}

/* 模型选择 */
.model-select-wrap {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-sm);
}

.model-chips {
  display: flex;
  flex-wrap: wrap;
  gap: var(--apple-spacing-xs);
  min-height: 32px;
  padding: var(--apple-spacing-xs);
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-sm);
}

.model-chip {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 8px;
  background: var(--apple-blue);
  color: white;
  border-radius: var(--apple-radius-sm);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.chip-remove {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 16px;
  height: 16px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 50%;
  font-size: 12px;
  cursor: pointer;
}

.chip-remove:hover {
  background: rgba(255, 255, 255, 0.4);
}

/* 响应式 */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-lg);
  }

  .form-row, .form-row.three-col {
    grid-template-columns: 1fr;
  }

  .type-selector {
    flex-direction: column;
  }
}
</style>
