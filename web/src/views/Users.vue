<!--
 * 文件作用：用户管理页面 - Apple HIG 风格
 * 负责功能：
 *   - 用户列表和CRUD
 *   - 用户套餐分配管理
 *   - 用户API Key管理
 *   - 批量设置价格倍率
 *   - 用户使用统计查看
 * 重要程度：⭐⭐⭐⭐⭐ 核心（用户管理）
-->
<template>
  <div class="users-page">
    <!-- 页面标题 -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">用户管理</h1>
        <p class="page-subtitle">管理系统用户、权限和套餐分配</p>
      </div>
      <div class="header-actions">
        <button
          class="btn btn-outline"
          :disabled="selectedUsers.length === 0"
          @click="showBatchRateDialog"
        >
          批量设置倍率 ({{ selectedUsers.length }})
        </button>
        <button class="btn btn-primary" @click="showCreateUserDialog">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <line x1="12" y1="5" x2="12" y2="19"/>
            <line x1="5" y1="12" x2="19" y2="12"/>
          </svg>
          创建用户
        </button>
      </div>
    </div>

    <!-- 用户表格 -->
    <div class="table-card">
      <div class="table-container">
        <table class="data-table">
          <thead>
            <tr>
              <th class="col-checkbox">
                <input
                  type="checkbox"
                  :checked="selectedUsers.length === users.length && users.length > 0"
                  @change="toggleSelectAll"
                />
              </th>
              <th class="col-id">ID</th>
              <th class="col-username">用户名</th>
              <th class="col-email">邮箱</th>
              <th class="col-role">角色</th>
              <th class="col-status">状态</th>
              <th class="col-rate">价格倍率</th>
              <th class="col-concurrency">并发限制</th>
              <th class="col-balance">额度余额</th>
              <th class="col-daily">订阅日余额</th>
              <th class="col-time">创建时间</th>
              <th class="col-actions">操作</th>
            </tr>
          </thead>
          <tbody v-if="!loading">
            <tr v-for="user in users" :key="user.id" :class="{ selected: isSelected(user) }">
              <td class="col-checkbox">
                <input type="checkbox" :checked="isSelected(user)" @change="toggleSelect(user)" />
              </td>
              <td class="col-id">{{ user.id }}</td>
              <td class="col-username">
                <div class="user-info">
                  <div class="user-avatar">{{ user.username.charAt(0).toUpperCase() }}</div>
                  <span>{{ user.username }}</span>
                </div>
              </td>
              <td class="col-email">{{ user.email || '-' }}</td>
              <td class="col-role">
                <span :class="['badge', user.role === 'admin' ? 'badge-danger' : 'badge-info']">
                  {{ user.role === 'admin' ? '管理员' : '普通用户' }}
                </span>
              </td>
              <td class="col-status">
                <span :class="['badge', user.status === 'active' ? 'badge-success' : 'badge-warning']">
                  {{ user.status === 'active' ? '正常' : '禁用' }}
                </span>
              </td>
              <td class="col-rate">
                <span :class="['rate-badge', getPriceRateClass(user.price_rate)]">
                  {{ user.price_rate === 0 ? '免费' : `${user.price_rate}x` }}
                </span>
              </td>
              <td class="col-concurrency">{{ user.max_concurrency || 10 }}</td>
              <td class="col-balance">
                <span :class="{ 'low-balance': user.quota_key_balance < 1 && user.quota_key_balance > 0 }">
                  ${{ (user.quota_key_balance || 0).toFixed(2) }}
                </span>
              </td>
              <td class="col-daily">
                <span v-if="user.subscription_daily_remain === -1" class="unlimited">无限</span>
                <span v-else :class="{ 'low-balance': user.subscription_daily_remain < 1 && user.subscription_daily_remain > 0 }">
                  ${{ (user.subscription_daily_remain || 0).toFixed(2) }}
                </span>
              </td>
              <td class="col-time">{{ formatDate(user.created_at) }}</td>
              <td class="col-actions">
                <div class="action-group">
                  <button class="action-btn" @click="handleEdit(user)" title="编辑">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M11 4H4a2 2 0 00-2 2v14a2 2 0 002 2h14a2 2 0 002-2v-7"/>
                      <path d="M18.5 2.5a2.121 2.121 0 013 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                  </button>
                  <button class="action-btn success" @click="viewAPIKeys(user)" title="API Key">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
                    </svg>
                  </button>
                  <button class="action-btn warning" @click="viewPackages(user)" title="套餐">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z"/>
                      <line x1="3" y1="6" x2="21" y2="6"/>
                      <path d="M16 10a4 4 0 01-8 0"/>
                    </svg>
                  </button>
                  <button class="action-btn info" @click="viewUsage(user)" title="统计">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                      <path d="M3 3v18h18"/>
                      <path d="M18 9l-5 5-4-4-3 3"/>
                    </svg>
                  </button>
                  <button class="action-btn danger" @click="confirmDelete(user)" title="删除">
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

        <!-- 加载状态 -->
        <div v-if="loading" class="loading-state">
          <div class="loading-spinner"></div>
          <span>加载中...</span>
        </div>

        <!-- 空状态 -->
        <div v-if="!loading && users.length === 0" class="empty-state">
          <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
            <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
            <circle cx="9" cy="7" r="4"/>
            <path d="M23 21v-2a4 4 0 00-3-3.87"/>
            <path d="M16 3.13a4 4 0 010 7.75"/>
          </svg>
          <span>暂无用户</span>
        </div>
      </div>

      <!-- 分页 -->
      <div class="pagination-wrap" v-if="pagination.total > 0">
        <div class="pagination-info">
          共 {{ pagination.total }} 条记录
        </div>
        <div class="pagination-controls">
          <select v-model="pagination.pageSize" @change="fetchUsers" class="page-size-select">
            <option :value="10">10 条/页</option>
            <option :value="20">20 条/页</option>
            <option :value="50">50 条/页</option>
          </select>
          <div class="page-btns">
            <button
              class="page-btn"
              :disabled="pagination.page <= 1"
              @click="pagination.page--; fetchUsers()"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="15,18 9,12 15,6"/>
              </svg>
            </button>
            <span class="page-current">{{ pagination.page }} / {{ Math.ceil(pagination.total / pagination.pageSize) }}</span>
            <button
              class="page-btn"
              :disabled="pagination.page >= Math.ceil(pagination.total / pagination.pageSize)"
              @click="pagination.page++; fetchUsers()"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <polyline points="9,18 15,12 9,6"/>
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑用户弹窗 -->
    <Teleport to="body">
      <div v-if="dialogVisible" class="modal-overlay" @click.self="dialogVisible = false">
        <div class="modal">
          <div class="modal-header">
            <h2>编辑用户</h2>
            <button class="modal-close" @click="dialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">用户名</label>
              <input :value="editForm.username" disabled class="form-input disabled" />
            </div>
            <div class="form-group">
              <label class="form-label">邮箱</label>
              <input v-model="editForm.email" type="email" class="form-input" placeholder="可选" />
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">角色</label>
                <select v-model="editForm.role" class="form-select">
                  <option value="user">普通用户</option>
                  <option value="admin">管理员</option>
                </select>
              </div>
              <div class="form-group">
                <label class="form-label">状态</label>
                <select v-model="editForm.status" class="form-select">
                  <option value="active">正常</option>
                  <option value="disabled">禁用</option>
                </select>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">价格倍率</label>
                <input v-model.number="editForm.price_rate" type="number" min="0" max="10" step="0.1" class="form-input" />
                <p class="form-tip">0 = 免费，1 = 原价，1.5 = 1.5倍</p>
              </div>
              <div class="form-group">
                <label class="form-label">最大并发</label>
                <input v-model.number="editForm.max_concurrency" type="number" min="1" max="100" class="form-input" />
                <p class="form-tip">用户同时进行的最大请求数</p>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="dialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="submitting" @click="handleSubmit">
              <span v-if="submitting" class="btn-loading"></span>
              {{ submitting ? '保存中...' : '保存' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 创建用户弹窗 -->
    <Teleport to="body">
      <div v-if="createUserDialogVisible" class="modal-overlay" @click.self="createUserDialogVisible = false">
        <div class="modal">
          <div class="modal-header">
            <h2>创建用户</h2>
            <button class="modal-close" @click="createUserDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">用户名 <span class="required">*</span></label>
              <input v-model="createUserForm.username" type="text" class="form-input" placeholder="请输入用户名" />
            </div>
            <div class="form-group">
              <label class="form-label">密码 <span class="required">*</span></label>
              <input v-model="createUserForm.password" type="password" class="form-input" placeholder="请输入密码" />
            </div>
            <div class="form-group">
              <label class="form-label">邮箱</label>
              <input v-model="createUserForm.email" type="email" class="form-input" placeholder="可选" />
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">角色</label>
                <select v-model="createUserForm.role" class="form-select">
                  <option value="user">普通用户</option>
                  <option value="admin">管理员</option>
                </select>
              </div>
              <div class="form-group">
                <label class="form-label">状态</label>
                <select v-model="createUserForm.status" class="form-select">
                  <option value="active">正常</option>
                  <option value="disabled">禁用</option>
                </select>
              </div>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">价格倍率</label>
                <input v-model.number="createUserForm.price_rate" type="number" min="0" max="10" step="0.1" class="form-input" />
              </div>
              <div class="form-group">
                <label class="form-label">最大并发</label>
                <input v-model.number="createUserForm.max_concurrency" type="number" min="1" max="100" class="form-input" />
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="createUserDialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="creatingUser" @click="handleCreateUser">
              <span v-if="creatingUser" class="btn-loading"></span>
              {{ creatingUser ? '创建中...' : '创建' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 批量设置倍率弹窗 -->
    <Teleport to="body">
      <div v-if="batchRateDialogVisible" class="modal-overlay" @click.self="batchRateDialogVisible = false">
        <div class="modal modal-sm">
          <div class="modal-header">
            <h2>批量设置价格倍率</h2>
            <button class="modal-close" @click="batchRateDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">选中用户</label>
              <div class="selected-users-tags">
                <span v-for="user in selectedUsers" :key="user.id" class="user-tag">
                  {{ user.username }}
                </span>
              </div>
            </div>
            <div class="form-group">
              <label class="form-label">价格倍率</label>
              <input v-model.number="batchRate" type="number" min="0" max="10" step="0.1" class="form-input" />
              <p class="form-tip">0 = 免费，1 = 原价，1.5 = 1.5倍</p>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="batchRateDialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="batchSubmitting" @click="submitBatchRate">
              <span v-if="batchSubmitting" class="btn-loading"></span>
              {{ batchSubmitting ? '设置中...' : '确认设置' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 使用统计弹窗 -->
    <Teleport to="body">
      <div v-if="usageDialogVisible" class="modal-overlay" @click.self="usageDialogVisible = false">
        <div class="modal modal-lg">
          <div class="modal-header">
            <h2>{{ usageUser?.username }} 的使用统计</h2>
            <button class="modal-close" @click="usageDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div v-if="usageLoading" class="loading-state">
              <div class="loading-spinner"></div>
              <span>加载中...</span>
            </div>
            <template v-else>
              <!-- 统计概览 -->
              <div class="stats-summary">
                <div class="stat-card">
                  <div class="stat-value">{{ formatNumber(usageSummary.totalRequests) }}</div>
                  <div class="stat-label">总请求数</div>
                </div>
                <div class="stat-card">
                  <div class="stat-value">{{ formatLargeNumber(usageSummary.totalTokens) }}</div>
                  <div class="stat-label">总 Token</div>
                </div>
                <div class="stat-card">
                  <div class="stat-value cost">${{ usageSummary.totalCost.toFixed(4) }}</div>
                  <div class="stat-label">总消费</div>
                </div>
              </div>

              <!-- 套餐使用情况 -->
              <div class="usage-section" v-if="packageUsage.length > 0">
                <h3 class="section-title">套餐使用情况</h3>
                <div class="package-usage-list">
                  <div v-for="pkg in packageUsage" :key="pkg.id" class="package-usage-item">
                    <div class="package-info">
                      <span class="package-name">{{ pkg.name }}</span>
                      <span :class="['badge', pkg.type === 'subscription' ? 'badge-primary' : 'badge-success']">
                        {{ pkg.type === 'subscription' ? '订阅' : '额度' }}
                      </span>
                      <span :class="['badge', pkg.status === 'active' ? 'badge-success' : 'badge-secondary']">
                        {{ pkg.status === 'active' ? '有效' : pkg.status === 'expired' ? '已过期' : '已用尽' }}
                      </span>
                    </div>
                    <div class="package-usage">
                      <template v-if="pkg.type === 'subscription'">
                        <span v-if="pkg.daily_quota > 0">日: ${{ pkg.daily_used?.toFixed(2) || '0.00' }}/${{ pkg.daily_quota }}</span>
                        <span v-if="pkg.weekly_quota > 0">周: ${{ pkg.weekly_used?.toFixed(2) || '0.00' }}/${{ pkg.weekly_quota }}</span>
                        <span v-if="pkg.monthly_quota > 0">月: ${{ pkg.monthly_used?.toFixed(2) || '0.00' }}/${{ pkg.monthly_quota }}</span>
                        <span v-if="!pkg.daily_quota && !pkg.weekly_quota && !pkg.monthly_quota" class="unlimited">无限额</span>
                      </template>
                      <template v-else>
                        <div class="quota-progress">
                          <div class="progress-bar">
                            <div
                              class="progress-fill"
                              :style="{ width: pkg.quota_total > 0 ? Math.min(100, (pkg.quota_used / pkg.quota_total) * 100) + '%' : '0%' }"
                              :class="{ danger: pkg.quota_used >= pkg.quota_total }"
                            ></div>
                          </div>
                          <span>${{ pkg.quota_used?.toFixed(2) || '0.00' }} / ${{ pkg.quota_total?.toFixed(2) || '0.00' }}</span>
                        </div>
                      </template>
                    </div>
                  </div>
                </div>
              </div>

              <!-- 按模型统计 -->
              <div class="usage-section" v-if="modelUsage.length > 0">
                <h3 class="section-title">按模型统计</h3>
                <div class="model-usage-table">
                  <div class="model-row header">
                    <span>模型</span>
                    <span>请求数</span>
                    <span>Token 数</span>
                    <span>费用</span>
                  </div>
                  <div v-for="item in modelUsage" :key="item.model" class="model-row">
                    <span class="model-name">{{ item.model }}</span>
                    <span>{{ formatNumber(item.requests) }}</span>
                    <span>{{ formatLargeNumber(item.tokens) }}</span>
                    <span class="cost">${{ item.cost.toFixed(4) }}</span>
                  </div>
                </div>
              </div>
            </template>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="usageDialogVisible = false">关闭</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- API Key 管理弹窗 -->
    <Teleport to="body">
      <div v-if="apiKeyDialogVisible" class="modal-overlay" @click.self="apiKeyDialogVisible = false">
        <div class="modal modal-xl">
          <div class="modal-header">
            <h2>{{ apiKeyUser?.username }} 的 API Key</h2>
            <button class="modal-close" @click="apiKeyDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="apikey-header">
              <button class="btn btn-primary btn-sm" @click="showCreateAPIKeyDialog">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                创建 API Key
              </button>
            </div>

            <div v-if="apiKeyLoading" class="loading-state">
              <div class="loading-spinner"></div>
              <span>加载中...</span>
            </div>

            <div v-else-if="apiKeys.length === 0" class="empty-state small">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
              </svg>
              <span>暂无 API Key</span>
            </div>

            <div v-else class="apikey-list">
              <div v-for="key in apiKeys" :key="key.id" class="apikey-item">
                <div class="apikey-main">
                  <div class="apikey-info">
                    <span class="apikey-name">{{ key.name }}</span>
                    <code class="apikey-prefix">{{ key.key_prefix }}...</code>
                  </div>
                  <div class="apikey-meta">
                    <span v-if="key.user_package" class="badge badge-primary">{{ key.user_package.name }}</span>
                    <span v-else class="badge badge-secondary">未绑定</span>
                    <span :class="['badge', key.status === 'active' ? 'badge-success' : 'badge-danger']">
                      {{ key.status === 'active' ? '正常' : '禁用' }}
                    </span>
                  </div>
                </div>
                <div class="apikey-stats">
                  <span>{{ key.rate_limit }}/分钟</span>
                  <span>{{ key.request_count }} 请求</span>
                  <span class="cost">${{ (key.cost_used || 0).toFixed(4) }}</span>
                </div>
                <div class="apikey-actions">
                  <button
                    :class="['btn btn-sm', key.status === 'active' ? 'btn-warning' : 'btn-success']"
                    @click="handleToggleAPIKey(key)"
                  >
                    {{ key.status === 'active' ? '禁用' : '启用' }}
                  </button>
                  <button class="btn btn-sm btn-danger" @click="handleDeleteAPIKey(key)">删除</button>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="apiKeyDialogVisible = false">关闭</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 创建 API Key 弹窗 -->
    <Teleport to="body">
      <div v-if="createKeyDialogVisible" class="modal-overlay" @click.self="createKeyDialogVisible = false">
        <div class="modal">
          <div class="modal-header">
            <h2>创建 API Key</h2>
            <button class="modal-close" @click="createKeyDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">名称 <span class="required">*</span></label>
              <input v-model="createKeyForm.name" type="text" class="form-input" placeholder="为这个 Key 起个名字" />
            </div>
            <div class="form-group">
              <label class="form-label">绑定套餐 <span class="required">*</span></label>
              <select v-model="createKeyForm.user_package_id" class="form-select">
                <option :value="null" disabled>请选择要绑定的套餐</option>
                <option v-for="pkg in userPackagesForKey" :key="pkg.id" :value="pkg.id">
                  {{ pkg.name }} ({{ pkg.type === 'subscription' ? '订阅' : '额度' }})
                </option>
              </select>
              <p class="form-tip">API Key 将从绑定的套餐扣费</p>
            </div>
            <div class="form-group">
              <label class="form-label">允许的平台</label>
              <select v-model="createKeyForm.allowed_platforms" class="form-select">
                <option value="all">全部平台</option>
                <option value="claude">仅 Claude</option>
                <option value="openai">仅 OpenAI</option>
                <option value="gemini">仅 Gemini</option>
              </select>
            </div>
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">每分钟限制</label>
                <input v-model.number="createKeyForm.rate_limit" type="number" min="1" max="1000" class="form-input" />
              </div>
              <div class="form-group">
                <label class="form-label">每日限制 (0=不限)</label>
                <input v-model.number="createKeyForm.daily_limit" type="number" min="0" class="form-input" />
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="createKeyDialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="creatingKey" @click="handleCreateAPIKey">
              <span v-if="creatingKey" class="btn-loading"></span>
              {{ creatingKey ? '创建中...' : '创建' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 新创建的 API Key 显示 -->
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
              <span>请立即复制并妥善保存 API Key，关闭后无法再次查看！</span>
            </div>
            <div class="key-display">
              <code>{{ newKey }}</code>
              <button class="copy-btn" @click="copyKey">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <rect x="9" y="9" width="13" height="13" rx="2" ry="2"/>
                  <path d="M5 15H4a2 2 0 01-2-2V4a2 2 0 012-2h9a2 2 0 012 2v1"/>
                </svg>
                {{ copied ? '已复制' : '复制' }}
              </button>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-primary" @click="showNewKeyDialog = false">我已保存，关闭</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 用户套餐管理弹窗 -->
    <Teleport to="body">
      <div v-if="packageDialogVisible" class="modal-overlay" @click.self="packageDialogVisible = false">
        <div class="modal modal-lg">
          <div class="modal-header">
            <h2>{{ packageUser?.username }} 的套餐</h2>
            <button class="modal-close" @click="packageDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="package-header">
              <button class="btn btn-primary btn-sm" @click="showAssignPackageDialog">
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                  <line x1="12" y1="5" x2="12" y2="19"/>
                  <line x1="5" y1="12" x2="19" y2="12"/>
                </svg>
                分配套餐
              </button>
            </div>

            <div v-if="packageLoading" class="loading-state">
              <div class="loading-spinner"></div>
              <span>加载中...</span>
            </div>

            <div v-else-if="userPackages.length === 0" class="empty-state small">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z"/>
                <line x1="3" y1="6" x2="21" y2="6"/>
                <path d="M16 10a4 4 0 01-8 0"/>
              </svg>
              <span>暂无套餐</span>
            </div>

            <div v-else class="user-packages-list">
              <div v-for="pkg in userPackages" :key="pkg.id" class="user-package-item">
                <div class="package-main">
                  <div class="package-title">
                    <span class="package-name">{{ pkg.name }}</span>
                    <span :class="['badge', pkg.type === 'subscription' ? 'badge-primary' : 'badge-success']">
                      {{ pkg.type === 'subscription' ? '订阅' : '额度' }}
                    </span>
                    <span :class="['badge', pkg.status === 'active' ? 'badge-success' : 'badge-secondary']">
                      {{ pkg.status === 'active' ? '有效' : pkg.status === 'expired' ? '已过期' : '已用尽' }}
                    </span>
                  </div>
                  <div class="package-details">
                    <template v-if="pkg.type === 'subscription'">
                      <div class="subscription-info">
                        <span class="expire-info">有效期: {{ formatDate(pkg.expire_time) }}</span>
                        <div class="quota-tags">
                          <span v-if="pkg.daily_quota > 0" class="quota-tag">
                            日: ${{ (pkg.daily_used || 0).toFixed(2) }}/${{ pkg.daily_quota }}
                          </span>
                          <span v-if="pkg.weekly_quota > 0" class="quota-tag">
                            周: ${{ (pkg.weekly_used || 0).toFixed(2) }}/${{ pkg.weekly_quota }}
                          </span>
                          <span v-if="pkg.monthly_quota > 0" class="quota-tag">
                            月: ${{ (pkg.monthly_used || 0).toFixed(2) }}/${{ pkg.monthly_quota }}
                          </span>
                          <span v-if="!pkg.daily_quota && !pkg.weekly_quota && !pkg.monthly_quota" class="quota-tag unlimited">
                            无限额
                          </span>
                        </div>
                      </div>
                    </template>
                    <template v-else>
                      <div class="quota-progress">
                        <span>额度: ${{ (pkg.quota_used || 0).toFixed(2) }} / ${{ (pkg.quota_total || 0).toFixed(2) }}</span>
                        <div class="progress-bar">
                          <div
                            class="progress-fill"
                            :style="{ width: pkg.quota_total > 0 ? Math.min(100, (pkg.quota_used / pkg.quota_total) * 100) + '%' : '0%' }"
                            :class="{ danger: pkg.quota_used >= pkg.quota_total }"
                          ></div>
                        </div>
                      </div>
                    </template>
                    <div v-if="pkg.allowed_models" class="model-limit">
                      <span class="badge badge-warning">限制 {{ pkg.allowed_models.split(',').length }} 个模型</span>
                    </div>
                  </div>
                </div>
                <div class="package-actions">
                  <button class="btn btn-sm btn-outline" @click="handleEditUserPackage(pkg)">编辑</button>
                  <button class="btn btn-sm btn-danger" @click="handleDeleteUserPackage(pkg.id)">删除</button>
                </div>
              </div>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="packageDialogVisible = false">关闭</button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 分配套餐弹窗 -->
    <Teleport to="body">
      <div v-if="assignPackageDialogVisible" class="modal-overlay" @click.self="assignPackageDialogVisible = false">
        <div class="modal modal-sm">
          <div class="modal-header">
            <h2>分配套餐</h2>
            <button class="modal-close" @click="assignPackageDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-group">
              <label class="form-label">选择套餐</label>
              <select v-model="selectedPackageId" class="form-select">
                <option :value="null" disabled>请选择套餐</option>
                <option v-for="pkg in availablePackages" :key="pkg.id" :value="pkg.id">
                  {{ pkg.name }} ({{ pkg.type === 'subscription' ? '订阅' : '额度' }})
                </option>
              </select>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="assignPackageDialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="assigningPackage || !selectedPackageId" @click="handleAssignPackage">
              <span v-if="assigningPackage" class="btn-loading"></span>
              {{ assigningPackage ? '分配中...' : '分配' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- 编辑用户套餐弹窗 -->
    <Teleport to="body">
      <div v-if="editPackageDialogVisible" class="modal-overlay" @click.self="editPackageDialogVisible = false">
        <div class="modal">
          <div class="modal-header">
            <h2>编辑用户套餐</h2>
            <button class="modal-close" @click="editPackageDialogVisible = false">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
                <line x1="18" y1="6" x2="6" y2="18"/>
                <line x1="6" y1="6" x2="18" y2="18"/>
              </svg>
            </button>
          </div>
          <div class="modal-body">
            <div class="form-row">
              <div class="form-group">
                <label class="form-label">状态</label>
                <select v-model="editPackageForm.status" class="form-select">
                  <option value="active">有效</option>
                  <option value="expired">已过期</option>
                  <option value="exhausted">已用尽</option>
                </select>
              </div>
              <div class="form-group">
                <label class="form-label">过期时间</label>
                <input v-model="editPackageForm.expire_time" type="datetime-local" class="form-input" />
              </div>
            </div>

            <!-- 订阅类型的周期额度 -->
            <template v-if="editPackageForm.type === 'subscription'">
              <div class="form-section-title">周期额度限制 (0=不限)</div>
              <div class="form-row three">
                <div class="form-group">
                  <label class="form-label">日限额($)</label>
                  <input v-model.number="editPackageForm.daily_quota" type="number" min="0" step="0.01" class="form-input" />
                </div>
                <div class="form-group">
                  <label class="form-label">周限额($)</label>
                  <input v-model.number="editPackageForm.weekly_quota" type="number" min="0" step="0.01" class="form-input" />
                </div>
                <div class="form-group">
                  <label class="form-label">月限额($)</label>
                  <input v-model.number="editPackageForm.monthly_quota" type="number" min="0" step="0.01" class="form-input" />
                </div>
              </div>
              <div class="form-section-title">当前周期已用 (可重置)</div>
              <div class="form-row three">
                <div class="form-group">
                  <label class="form-label">日已用($)</label>
                  <input v-model.number="editPackageForm.daily_used" type="number" min="0" step="0.01" class="form-input" />
                </div>
                <div class="form-group">
                  <label class="form-label">周已用($)</label>
                  <input v-model.number="editPackageForm.weekly_used" type="number" min="0" step="0.01" class="form-input" />
                </div>
                <div class="form-group">
                  <label class="form-label">月已用($)</label>
                  <input v-model.number="editPackageForm.monthly_used" type="number" min="0" step="0.01" class="form-input" />
                </div>
              </div>
            </template>

            <!-- 额度类型字段 -->
            <template v-if="editPackageForm.type === 'quota'">
              <div class="form-row">
                <div class="form-group">
                  <label class="form-label">总额度($)</label>
                  <input v-model.number="editPackageForm.quota_total" type="number" min="0" step="0.01" class="form-input" />
                </div>
                <div class="form-group">
                  <label class="form-label">已用额度($)</label>
                  <input v-model.number="editPackageForm.quota_used" type="number" min="0" step="0.01" class="form-input" />
                </div>
              </div>
            </template>

            <div class="form-group">
              <label class="form-label">允许的模型</label>
              <select v-model="editPackageSelectedModels" multiple class="form-select multi">
                <option v-for="model in modelList" :key="model.id" :value="model.name">
                  {{ model.name }}
                </option>
              </select>
              <p class="form-tip">留空表示全部模型，按住 Ctrl/Cmd 可多选</p>
            </div>
          </div>
          <div class="modal-footer">
            <button class="btn btn-secondary" @click="editPackageDialogVisible = false">取消</button>
            <button class="btn btn-primary" :disabled="savingPackage" @click="handleSaveUserPackage">
              <span v-if="savingPackage" class="btn-loading"></span>
              {{ savingPackage ? '保存中...' : '保存' }}
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
            <p class="delete-message">确定要删除用户 "{{ deleteTarget?.username }}" 吗？此操作不可撤销。</p>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from '@/utils/toast'
import api from '@/api'

const loading = ref(false)
const users = ref([])
const pagination = reactive({ page: 1, pageSize: 20, total: 0 })
const modelList = ref([])

const dialogVisible = ref(false)
const submitting = ref(false)
const editForm = reactive({ id: 0, username: '', email: '', role: '', status: '', price_rate: 1.0, max_concurrency: 10 })

// 批量设置倍率相关
const selectedUsers = ref([])
const batchRateDialogVisible = ref(false)
const batchRate = ref(1.0)
const batchSubmitting = ref(false)

// 用户使用统计相关
const usageDialogVisible = ref(false)
const usageLoading = ref(false)
const usageUser = ref(null)
const usageSummary = reactive({ totalRequests: 0, totalTokens: 0, totalCost: 0 })
const modelUsage = ref([])
const packageUsage = ref([])

// API Key 管理相关
const apiKeyDialogVisible = ref(false)
const apiKeyLoading = ref(false)
const apiKeyUser = ref(null)
const apiKeys = ref([])

// 创建 API Key 相关
const createKeyDialogVisible = ref(false)
const creatingKey = ref(false)
const userPackagesForKey = ref([])
const createKeyForm = ref({
  name: '',
  user_package_id: null,
  allowed_platforms: 'all',
  rate_limit: 60,
  daily_limit: 0
})

// 显示新 Key
const showNewKeyDialog = ref(false)
const newKey = ref('')
const copied = ref(false)

// 套餐管理相关
const packageDialogVisible = ref(false)
const packageLoading = ref(false)
const packageUser = ref(null)
const userPackages = ref([])
const availablePackages = ref([])

// 分配套餐相关
const assignPackageDialogVisible = ref(false)
const selectedPackageId = ref(null)
const assigningPackage = ref(false)

// 编辑用户套餐相关
const editPackageDialogVisible = ref(false)
const savingPackage = ref(false)
const editPackageForm = ref({
  id: 0,
  type: '',
  status: '',
  expire_time: null,
  daily_quota: 0,
  weekly_quota: 0,
  monthly_quota: 0,
  daily_used: 0,
  weekly_used: 0,
  monthly_used: 0,
  quota_total: 0,
  quota_used: 0,
  allowed_models: ''
})

const editPackageSelectedModels = computed({
  get() {
    if (!editPackageForm.value.allowed_models) return []
    return editPackageForm.value.allowed_models.split(',').filter(m => m.trim())
  },
  set(val) {
    editPackageForm.value.allowed_models = val.join(',')
  }
})

// 创建用户相关
const createUserDialogVisible = ref(false)
const creatingUser = ref(false)
const createUserForm = ref({
  username: '',
  password: '',
  email: '',
  role: 'user',
  status: 'active',
  price_rate: 1.0,
  max_concurrency: 10
})

// 删除相关
const deleteDialogVisible = ref(false)
const deleteTarget = ref(null)
const deleting = ref(false)

function getPriceRateClass(rate) {
  if (rate === 0) return 'free'
  if (rate < 1) return 'discount'
  if (rate > 1) return 'premium'
  return 'normal'
}

function formatDate(str) {
  if (!str) return '-'
  return new Date(str).toLocaleString('zh-CN')
}

function formatNumber(num) {
  if (!num) return '0'
  return num.toLocaleString()
}

function formatLargeNumber(num) {
  if (!num) return '0'
  if (num >= 1000000) return (num / 1000000).toFixed(2) + 'M'
  if (num >= 1000) return (num / 1000).toFixed(1) + 'K'
  return num.toLocaleString()
}

function isSelected(user) {
  return selectedUsers.value.some(u => u.id === user.id)
}

function toggleSelect(user) {
  const index = selectedUsers.value.findIndex(u => u.id === user.id)
  if (index >= 0) {
    selectedUsers.value.splice(index, 1)
  } else {
    selectedUsers.value.push(user)
  }
}

function toggleSelectAll() {
  if (selectedUsers.value.length === users.value.length) {
    selectedUsers.value = []
  } else {
    selectedUsers.value = [...users.value]
  }
}

async function fetchUsers() {
  loading.value = true
  try {
    const res = await api.getUsers({ page: pagination.page, page_size: pagination.pageSize })
    users.value = res.data.items || []
    pagination.total = res.data.total || 0
  } catch (e) {
    // handled
  } finally {
    loading.value = false
  }
}

async function fetchModels() {
  try {
    const res = await api.getModels()
    modelList.value = (res.data || []).filter(m => m.enabled)
  } catch (e) {
    // handled
  }
}

function handleEdit(row) {
  Object.assign(editForm, row)
  dialogVisible.value = true
}

async function handleSubmit() {
  submitting.value = true
  try {
    await api.updateUser(editForm.id, {
      email: editForm.email,
      role: editForm.role,
      status: editForm.status,
      price_rate: editForm.price_rate,
      max_concurrency: editForm.max_concurrency
    })
    ElMessage.success('更新成功')
    dialogVisible.value = false
    fetchUsers()
  } catch (e) {
    // handled
  } finally {
    submitting.value = false
  }
}

function confirmDelete(user) {
  deleteTarget.value = user
  deleteDialogVisible.value = true
}

async function handleDelete() {
  if (!deleteTarget.value) return
  deleting.value = true
  try {
    await api.deleteUser(deleteTarget.value.id)
    ElMessage.success('删除成功')
    deleteDialogVisible.value = false
    deleteTarget.value = null
    fetchUsers()
  } catch (e) {
    // handled
  } finally {
    deleting.value = false
  }
}

function showBatchRateDialog() {
  batchRate.value = 1.0
  batchRateDialogVisible.value = true
}

async function submitBatchRate() {
  if (selectedUsers.value.length === 0) return

  batchSubmitting.value = true
  try {
    const userIds = selectedUsers.value.map(u => u.id)
    await api.batchUpdateUserRates({
      user_ids: userIds,
      price_rate: batchRate.value
    })
    ElMessage.success(`已更新 ${userIds.length} 个用户的价格倍率`)
    batchRateDialogVisible.value = false
    selectedUsers.value = []
    fetchUsers()
  } catch (e) {
    // handled
  } finally {
    batchSubmitting.value = false
  }
}

async function viewUsage(row) {
  usageUser.value = row
  usageDialogVisible.value = true
  usageLoading.value = true

  usageSummary.totalRequests = 0
  usageSummary.totalTokens = 0
  usageSummary.totalCost = 0
  modelUsage.value = []
  packageUsage.value = []

  try {
    const res = await api.getUserUsageSummary(row.id)
    if (res.data) {
      usageSummary.totalRequests = res.data.total_requests || 0
      usageSummary.totalTokens = res.data.total_tokens || 0
      usageSummary.totalCost = res.data.total_cost || 0
      modelUsage.value = res.data.model_usage || []
      packageUsage.value = res.data.package_usage || []
    }
  } catch (e) {
    // handled
  } finally {
    usageLoading.value = false
  }
}

// ========== API Key 管理 ==========

async function viewAPIKeys(row) {
  apiKeyUser.value = row
  apiKeyDialogVisible.value = true
  apiKeyLoading.value = true
  apiKeys.value = []

  try {
    const res = await api.adminGetUserAPIKeys(row.id)
    apiKeys.value = res.data || []
  } catch (e) {
    // handled
  } finally {
    apiKeyLoading.value = false
  }
}

async function showCreateAPIKeyDialog() {
  createKeyForm.value = {
    name: '',
    user_package_id: null,
    allowed_platforms: 'all',
    rate_limit: 60,
    daily_limit: 0
  }
  try {
    const res = await api.getUserPackages(apiKeyUser.value.id)
    userPackagesForKey.value = (res.data || []).filter(p => p.status === 'active')
    if (userPackagesForKey.value.length === 0) {
      ElMessage.warning('该用户没有活跃的套餐，请先分配套餐')
      return
    }
  } catch (e) {
    ElMessage.error('获取用户套餐失败')
    return
  }
  createKeyDialogVisible.value = true
}

async function handleCreateAPIKey() {
  if (!createKeyForm.value.name || !createKeyForm.value.user_package_id) {
    ElMessage.warning('请填写完整信息')
    return
  }

  creatingKey.value = true
  try {
    const res = await api.adminCreateUserAPIKey(apiKeyUser.value.id, createKeyForm.value)
    newKey.value = res.data.key
    createKeyDialogVisible.value = false
    showNewKeyDialog.value = true
    viewAPIKeys(apiKeyUser.value)
  } catch (e) {
    // handled
  } finally {
    creatingKey.value = false
  }
}

async function handleToggleAPIKey(row) {
  try {
    await api.adminToggleUserAPIKey(apiKeyUser.value.id, row.id)
    ElMessage.success('状态已更新')
    viewAPIKeys(apiKeyUser.value)
  } catch (e) {
    // handled
  }
}

async function handleDeleteAPIKey(row) {
  try {
    await api.adminDeleteUserAPIKey(apiKeyUser.value.id, row.id)
    ElMessage.success('删除成功')
    viewAPIKeys(apiKeyUser.value)
  } catch (e) {
    // handled
  }
}

async function copyKey() {
  try {
    await navigator.clipboard.writeText(newKey.value)
    copied.value = true
    ElMessage.success('已复制到剪贴板')
    setTimeout(() => { copied.value = false }, 2000)
  } catch (e) {
    const input = document.createElement('input')
    input.value = newKey.value
    document.body.appendChild(input)
    input.select()
    document.execCommand('copy')
    document.body.removeChild(input)
    copied.value = true
    ElMessage.success('已复制到剪贴板')
    setTimeout(() => { copied.value = false }, 2000)
  }
}

// ========== 创建用户 ==========

function showCreateUserDialog() {
  createUserForm.value = {
    username: '',
    password: '',
    email: '',
    role: 'user',
    status: 'active',
    price_rate: 1.0,
    max_concurrency: 10
  }
  createUserDialogVisible.value = true
}

async function handleCreateUser() {
  if (!createUserForm.value.username || !createUserForm.value.password) {
    ElMessage.warning('请填写用户名和密码')
    return
  }

  creatingUser.value = true
  try {
    await api.createUser(createUserForm.value)
    ElMessage.success('创建成功')
    createUserDialogVisible.value = false
    fetchUsers()
  } catch (e) {
    // handled by api
  } finally {
    creatingUser.value = false
  }
}

// ========== 套餐管理 ==========

async function viewPackages(row) {
  packageUser.value = row
  packageDialogVisible.value = true
  packageLoading.value = true
  userPackages.value = []

  try {
    const res = await api.getUserPackages(row.id)
    userPackages.value = res.data || []
  } catch (e) {
    // handled
  } finally {
    packageLoading.value = false
  }
}

async function showAssignPackageDialog() {
  selectedPackageId.value = null
  try {
    const res = await api.getPackages()
    availablePackages.value = (res.data || []).filter(p => p.status === 'active')
  } catch (e) {
    // handled
  }
  assignPackageDialogVisible.value = true
}

async function handleAssignPackage() {
  if (!selectedPackageId.value) {
    ElMessage.warning('请选择套餐')
    return
  }

  assigningPackage.value = true
  try {
    await api.assignUserPackage(packageUser.value.id, { package_id: selectedPackageId.value })
    ElMessage.success('分配成功')
    assignPackageDialogVisible.value = false
    viewPackages(packageUser.value)
  } catch (e) {
    // handled
  } finally {
    assigningPackage.value = false
  }
}

function handleEditUserPackage(row) {
  editPackageForm.value = {
    id: row.id,
    type: row.type,
    status: row.status,
    expire_time: row.expire_time ? new Date(row.expire_time).toISOString().slice(0, 16) : '',
    daily_quota: row.daily_quota || 0,
    weekly_quota: row.weekly_quota || 0,
    monthly_quota: row.monthly_quota || 0,
    daily_used: row.daily_used || 0,
    weekly_used: row.weekly_used || 0,
    monthly_used: row.monthly_used || 0,
    quota_total: row.quota_total || 0,
    quota_used: row.quota_used || 0,
    allowed_models: row.allowed_models || ''
  }
  editPackageDialogVisible.value = true
}

async function handleSaveUserPackage() {
  savingPackage.value = true
  try {
    const data = {
      status: editPackageForm.value.status,
      expire_time: editPackageForm.value.expire_time,
      allowed_models: editPackageForm.value.allowed_models
    }

    if (editPackageForm.value.type === 'subscription') {
      data.daily_quota = editPackageForm.value.daily_quota
      data.weekly_quota = editPackageForm.value.weekly_quota
      data.monthly_quota = editPackageForm.value.monthly_quota
      data.daily_used = editPackageForm.value.daily_used
      data.weekly_used = editPackageForm.value.weekly_used
      data.monthly_used = editPackageForm.value.monthly_used
    } else {
      data.quota_total = editPackageForm.value.quota_total
      data.quota_used = editPackageForm.value.quota_used
    }

    await api.updateUserPackage(editPackageForm.value.id, data)
    ElMessage.success('保存成功')
    editPackageDialogVisible.value = false
    viewPackages(packageUser.value)
  } catch (e) {
    // handled
  } finally {
    savingPackage.value = false
  }
}

async function handleDeleteUserPackage(id) {
  try {
    await api.deleteUserPackage(id)
    ElMessage.success('删除成功')
    viewPackages(packageUser.value)
  } catch (e) {
    // handled
  }
}

onMounted(() => {
  fetchUsers()
  fetchModels()
})
</script>

<style scoped>
.users-page {
  max-width: 1400px;
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

.header-actions {
  display: flex;
  gap: var(--apple-spacing-sm);
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

.data-table tbody tr.selected {
  background: var(--apple-blue-light);
}

.col-checkbox { width: 40px; }
.col-id { width: 60px; }
.col-role, .col-status { width: 90px; }
.col-rate, .col-concurrency, .col-balance, .col-daily { width: 100px; }
.col-time { width: 160px; }
.col-actions { width: 200px; }

/* 用户信息 */
.user-info {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.user-avatar {
  width: 32px;
  height: 32px;
  background: linear-gradient(135deg, var(--apple-blue) 0%, var(--apple-purple) 100%);
  color: white;
  border-radius: var(--apple-radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: var(--apple-font-semibold);
  font-size: var(--apple-text-sm);
}

/* 徽章 */
.badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-medium);
}

.badge-primary { background: var(--apple-blue-light); color: var(--apple-blue); }
.badge-success { background: var(--apple-green-light); color: var(--apple-green); }
.badge-warning { background: var(--apple-orange-light); color: var(--apple-orange); }
.badge-danger { background: var(--apple-red-light); color: var(--apple-red); }
.badge-info { background: var(--apple-fill-tertiary); color: var(--apple-text-secondary); }
.badge-secondary { background: var(--apple-fill-quaternary); color: var(--apple-text-tertiary); }

/* 价格倍率标签 */
.rate-badge {
  display: inline-block;
  padding: 2px 8px;
  border-radius: var(--apple-radius-full);
  font-size: var(--apple-text-xs);
  font-weight: var(--apple-font-semibold);
}

.rate-badge.free { background: var(--apple-green-light); color: var(--apple-green); }
.rate-badge.discount { background: var(--apple-orange-light); color: var(--apple-orange); }
.rate-badge.normal { background: var(--apple-fill-tertiary); color: var(--apple-text-secondary); }
.rate-badge.premium { background: var(--apple-red-light); color: var(--apple-red); }

.low-balance { color: var(--apple-red); font-weight: var(--apple-font-semibold); }
.unlimited { color: var(--apple-green); font-weight: var(--apple-font-medium); }

/* 操作按钮组 */
.action-group {
  display: flex;
  gap: var(--apple-spacing-xxs);
}

.action-btn {
  width: 28px;
  height: 28px;
  border-radius: var(--apple-radius-xs);
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-secondary);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.action-btn svg {
  width: 14px;
  height: 14px;
}

.action-btn:hover {
  background: var(--apple-blue);
  color: white;
}

.action-btn.success:hover { background: var(--apple-green); }
.action-btn.warning:hover { background: var(--apple-orange); }
.action-btn.info:hover { background: var(--apple-teal); }
.action-btn.danger:hover { background: var(--apple-red); }

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

.empty-state.small {
  padding: var(--apple-spacing-xl);
}

.empty-state.small svg {
  width: 32px;
  height: 32px;
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

.btn-warning {
  background: var(--apple-orange);
  color: white;
}

.btn-success {
  background: var(--apple-green);
  color: white;
}

.btn-outline {
  background: transparent;
  color: var(--apple-blue);
  border: 1px solid var(--apple-blue);
}

.btn-outline:hover:not(:disabled) {
  background: var(--apple-blue-light);
}

.btn-sm {
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  font-size: var(--apple-text-xs);
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
.modal.modal-lg { max-width: 700px; }
.modal.modal-xl { max-width: 900px; }

.modal-header {
  padding: var(--apple-spacing-xl);
  border-bottom: 1px solid var(--apple-separator);
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.modal-header.success, .modal-header.danger {
  flex-direction: column;
  text-align: center;
  gap: var(--apple-spacing-md);
}

.success-icon, .danger-icon {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
}

.success-icon { background: var(--apple-green-light); }
.success-icon svg { width: 28px; height: 28px; color: var(--apple-green); }

.danger-icon { background: var(--apple-red-light); }
.danger-icon svg { width: 28px; height: 28px; color: var(--apple-red); }

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

.modal-close svg { width: 18px; height: 18px; }

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

.form-row {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: var(--apple-spacing-md);
}

.form-row.three {
  grid-template-columns: repeat(3, 1fr);
}

.form-label {
  display: block;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
  margin-bottom: var(--apple-spacing-xs);
}

.form-label .required { color: var(--apple-red); }

.form-input, .form-select {
  width: 100%;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-base);
  color: var(--apple-text-primary);
  background: var(--apple-bg-primary);
  border: 1px solid var(--apple-separator-opaque);
  border-radius: var(--apple-radius-sm);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.form-input:focus, .form-select:focus {
  outline: none;
  border-color: var(--apple-blue);
  box-shadow: 0 0 0 3px var(--apple-blue-light);
}

.form-input.disabled {
  background: var(--apple-bg-secondary);
  color: var(--apple-text-tertiary);
  cursor: not-allowed;
}

.form-select.multi {
  min-height: 100px;
}

.form-tip {
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  margin-top: var(--apple-spacing-xxs);
}

.form-section-title {
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-secondary);
  margin: var(--apple-spacing-lg) 0 var(--apple-spacing-sm) 0;
  padding-bottom: var(--apple-spacing-xs);
  border-bottom: 1px solid var(--apple-separator);
}

/* 选中用户标签 */
.selected-users-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--apple-spacing-xs);
  max-height: 100px;
  overflow-y: auto;
}

.user-tag {
  padding: 2px 8px;
  background: var(--apple-blue-light);
  color: var(--apple-blue);
  border-radius: var(--apple-radius-xs);
  font-size: var(--apple-text-xs);
}

/* 统计概览 */
.stats-summary {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--apple-spacing-md);
  margin-bottom: var(--apple-spacing-xl);
}

.stat-card {
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
  padding: var(--apple-spacing-lg);
  text-align: center;
}

.stat-card .stat-value {
  font-size: var(--apple-text-2xl);
  font-weight: var(--apple-font-bold);
  color: var(--apple-text-primary);
}

.stat-card .stat-value.cost {
  color: var(--apple-blue);
}

.stat-card .stat-label {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  margin-top: var(--apple-spacing-xs);
}

/* 使用详情部分 */
.usage-section {
  margin-top: var(--apple-spacing-xl);
}

.section-title {
  font-size: var(--apple-text-md);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-primary);
  margin: 0 0 var(--apple-spacing-md) 0;
}

.package-usage-list {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-sm);
}

.package-usage-item {
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
  padding: var(--apple-spacing-md);
}

.package-info {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  margin-bottom: var(--apple-spacing-sm);
}

.package-name {
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.package-usage {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  display: flex;
  gap: var(--apple-spacing-md);
  flex-wrap: wrap;
}

.quota-progress {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
}

.progress-bar {
  height: 6px;
  background: var(--apple-fill-tertiary);
  border-radius: var(--apple-radius-full);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  background: var(--apple-green);
  border-radius: var(--apple-radius-full);
  transition: width var(--apple-duration-normal) var(--apple-ease-default);
}

.progress-fill.danger {
  background: var(--apple-red);
}

/* 模型使用表格 */
.model-usage-table {
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
  overflow: hidden;
}

.model-row {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  border-bottom: 1px solid var(--apple-separator);
}

.model-row:last-child {
  border-bottom: none;
}

.model-row.header {
  background: var(--apple-fill-quaternary);
  font-weight: var(--apple-font-semibold);
  color: var(--apple-text-secondary);
  font-size: var(--apple-text-xs);
}

.model-name {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-sm);
}

.cost {
  color: var(--apple-blue);
  font-weight: var(--apple-font-medium);
}

/* API Key 列表 */
.apikey-header, .package-header {
  margin-bottom: var(--apple-spacing-lg);
}

.apikey-list {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-sm);
}

.apikey-item {
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
  padding: var(--apple-spacing-md);
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-md);
}

.apikey-main {
  flex: 1;
}

.apikey-info {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  margin-bottom: var(--apple-spacing-xs);
}

.apikey-name {
  font-weight: var(--apple-font-medium);
  color: var(--apple-text-primary);
}

.apikey-prefix {
  font-family: var(--apple-font-mono);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-tertiary);
  background: var(--apple-bg-tertiary);
  padding: 2px 6px;
  border-radius: var(--apple-radius-xs);
}

.apikey-meta {
  display: flex;
  gap: var(--apple-spacing-xs);
}

.apikey-stats {
  display: flex;
  gap: var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
}

.apikey-actions {
  display: flex;
  gap: var(--apple-spacing-xs);
}

/* 用户套餐列表 */
.user-packages-list {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-sm);
}

.user-package-item {
  background: var(--apple-bg-secondary);
  border-radius: var(--apple-radius-md);
  padding: var(--apple-spacing-md);
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.package-main {
  flex: 1;
}

.package-title {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  margin-bottom: var(--apple-spacing-sm);
}

.package-details {
  font-size: var(--apple-text-sm);
}

.subscription-info {
  display: flex;
  flex-direction: column;
  gap: var(--apple-spacing-xs);
}

.expire-info {
  color: var(--apple-text-secondary);
}

.quota-tags {
  display: flex;
  gap: var(--apple-spacing-sm);
  flex-wrap: wrap;
}

.quota-tag {
  padding: 2px 8px;
  background: var(--apple-fill-quaternary);
  border-radius: var(--apple-radius-xs);
  font-size: var(--apple-text-xs);
  color: var(--apple-text-secondary);
}

.quota-tag.unlimited {
  background: var(--apple-green-light);
  color: var(--apple-green);
}

.model-limit {
  margin-top: var(--apple-spacing-sm);
}

.package-actions {
  display: flex;
  gap: var(--apple-spacing-xs);
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
}

.warning-alert span {
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
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
}

.copy-btn svg {
  width: 14px;
  height: 14px;
}

.copy-btn:hover {
  background: var(--apple-blue-hover);
}

/* 删除确认 */
.delete-message {
  font-size: var(--apple-text-base);
  color: var(--apple-text-secondary);
  text-align: center;
  margin: 0;
}

/* 响应式 */
@media (max-width: 1024px) {
  .page-header {
    flex-direction: column;
    gap: var(--apple-spacing-lg);
  }

  .header-actions {
    width: 100%;
    justify-content: flex-end;
  }

  .form-row, .form-row.three {
    grid-template-columns: 1fr;
  }

  .stats-summary {
    grid-template-columns: 1fr;
  }
}
</style>
