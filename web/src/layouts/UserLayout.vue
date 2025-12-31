<!--
 * 文件作用：用户中心布局组件 - Apple HIG 风格
 * 负责功能：
 *   - 用户中心页面框架
 *   - 毛玻璃顶部导航栏
 *   - 侧边菜单
 *   - 内容区域
 * 重要程度：⭐⭐⭐⭐ 重要（用户界面框架）
-->
<template>
  <div class="user-layout">
    <!-- 顶部导航栏 -->
    <header class="top-navbar">
      <div class="navbar-left">
        <div class="logo">
          <svg class="logo-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <path d="M12 2L2 7l10 5 10-5-10-5zM2 17l10 5 10-5M2 12l10 5 10-5"/>
          </svg>
          <span class="logo-text">Cli-Proxy</span>
        </div>
      </div>

      <div class="navbar-right">
        <div class="user-dropdown" @click="showDropdown = !showDropdown" v-click-outside="closeDropdown">
          <div class="user-avatar">
            {{ userStore.user?.username?.charAt(0)?.toUpperCase() || 'U' }}
          </div>
          <span class="user-name">{{ userStore.user?.username || '用户' }}</span>
          <svg class="dropdown-arrow" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2">
            <polyline points="6,9 12,15 18,9"/>
          </svg>

          <!-- 下拉菜单 -->
          <div v-if="showDropdown" class="dropdown-menu">
            <div class="dropdown-item" @click="goToProfile">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M20 21v-2a4 4 0 00-4-4H8a4 4 0 00-4 4v2"/>
                <circle cx="12" cy="7" r="4"/>
              </svg>
              个人资料
            </div>
            <div v-if="userStore.user?.role === 'admin'" class="dropdown-item" @click="goToAdmin">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <circle cx="12" cy="12" r="3"/>
                <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/>
              </svg>
              管理后台
            </div>
            <div class="dropdown-divider"></div>
            <div class="dropdown-item danger" @click="handleLogout">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M9 21H5a2 2 0 01-2-2V5a2 2 0 012-2h4"/>
                <polyline points="16,17 21,12 16,7"/>
                <line x1="21" y1="12" x2="9" y2="12"/>
              </svg>
              退出登录
            </div>
          </div>
        </div>
      </div>
    </header>

    <div class="main-wrapper">
      <!-- 侧边栏 -->
      <aside :class="['sidebar', { 'is-collapsed': isCollapse }]">
        <nav class="sidebar-nav">
          <router-link
            to="/user/dashboard"
            :class="['nav-item', { active: isActive('/user/dashboard') }]"
            @mouseenter="prefetchFor('/user/dashboard')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <rect x="3" y="3" width="7" height="7"/>
                <rect x="14" y="3" width="7" height="7"/>
                <rect x="14" y="14" width="7" height="7"/>
                <rect x="3" y="14" width="7" height="7"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">仪表盘</span>
          </router-link>

          <router-link
            to="/user/api-keys"
            :class="['nav-item', { active: isActive('/user/api-keys') }]"
            @mouseenter="prefetchFor('/user/api-keys')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M21 2l-2 2m-7.61 7.61a5.5 5.5 0 11-7.778 7.778 5.5 5.5 0 017.777-7.777zm0 0L15.5 7.5m0 0l3 3L22 7l-3-3m-3.5 3.5L19 4"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">我的 API Key</span>
          </router-link>

          <router-link
            to="/user/packages"
            :class="['nav-item', { active: isActive('/user/packages') }]"
            @mouseenter="prefetchFor('/user/packages')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M6 2L3 6v14a2 2 0 002 2h14a2 2 0 002-2V6l-3-4z"/>
                <line x1="3" y1="6" x2="21" y2="6"/>
                <path d="M16 10a4 4 0 01-8 0"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">我的套餐</span>
          </router-link>

          <router-link
            to="/user/records"
            :class="['nav-item', { active: isActive('/user/records') }]"
            @mouseenter="prefetchFor('/user/records')"
          >
            <span class="nav-icon">
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
                <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                <polyline points="14,2 14,8 20,8"/>
                <line x1="16" y1="13" x2="8" y2="13"/>
                <line x1="16" y1="17" x2="8" y2="17"/>
                <polyline points="10,9 9,9 8,9"/>
              </svg>
            </span>
            <span v-if="!isCollapse" class="nav-label">使用记录</span>
          </router-link>
        </nav>

        <!-- 折叠按钮 -->
        <div class="sidebar-footer">
          <button class="collapse-toggle" @click="isCollapse = !isCollapse">
            <svg v-if="isCollapse" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <polyline points="9,18 15,12 9,6"/>
            </svg>
            <svg v-else viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5">
              <polyline points="15,18 9,12 15,6"/>
            </svg>
          </button>
        </div>
      </aside>

      <!-- 主内容区 -->
      <main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { prefetchChunk } from '@/prefetch'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const isCollapse = ref(false)
const showDropdown = ref(false)

function isActive(path) {
  return route.path === path
}

function prefetchFor(path) {
  const loaders = {
    '/user/dashboard': () => import('@/views/user/UserDashboard.vue'),
    '/user/api-keys': () => import('@/views/user/MyAPIKeys.vue'),
    '/user/packages': () => import('@/views/user/MyPackages.vue'),
    '/user/records': () => import('@/views/user/MyUsageRecords.vue'),
    '/user/profile': () => import('@/views/Profile.vue'),
    '/admin/system-monitor': () => import('@/views/SystemMonitor.vue')
  }
  const loader = loaders[path]
  if (!loader) return
  prefetchChunk(path, loader)
}

function closeDropdown() {
  showDropdown.value = false
}

function goToProfile() {
  showDropdown.value = false
  prefetchFor('/user/profile')
  router.push('/user/profile')
}

function goToAdmin() {
  showDropdown.value = false
  prefetchFor('/admin/system-monitor')
  router.push('/admin/system-monitor')
}

function handleLogout() {
  showDropdown.value = false
  userStore.logout()
  router.push('/login')
}

// 点击外部关闭下拉菜单指令
const vClickOutside = {
  mounted(el, binding) {
    el._clickOutside = (e) => {
      if (!el.contains(e.target)) {
        binding.value()
      }
    }
    document.addEventListener('click', el._clickOutside)
  },
  unmounted(el) {
    document.removeEventListener('click', el._clickOutside)
  }
}
</script>

<style scoped>
.user-layout {
  height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--apple-bg-secondary);
}

/* 顶部导航栏 */
.top-navbar {
  height: 56px;
  background: linear-gradient(135deg, #007aff 0%, #5856d6 100%);
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0 var(--apple-spacing-xl);
  position: relative;
  z-index: 100;
}

.navbar-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
}

.logo-icon {
  width: 28px;
  height: 28px;
  color: white;
}

.logo-text {
  font-size: var(--apple-text-lg);
  font-weight: var(--apple-font-bold);
  color: white;
}

.navbar-right {
  display: flex;
  align-items: center;
}

/* 用户下拉菜单 */
.user-dropdown {
  position: relative;
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-xs);
  padding: var(--apple-spacing-xs) var(--apple-spacing-sm);
  border-radius: var(--apple-radius-md);
  cursor: pointer;
  transition: background var(--apple-duration-fast) var(--apple-ease-default);
}

.user-dropdown:hover {
  background: rgba(255, 255, 255, 0.15);
}

.user-avatar {
  width: 32px;
  height: 32px;
  border-radius: var(--apple-radius-full);
  background: rgba(255, 255, 255, 0.25);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--apple-text-sm);
  font-weight: var(--apple-font-semibold);
}

.user-name {
  font-size: var(--apple-text-sm);
  color: white;
  margin: 0 var(--apple-spacing-xs);
}

.dropdown-arrow {
  width: 16px;
  height: 16px;
  color: rgba(255, 255, 255, 0.8);
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 8px);
  right: 0;
  min-width: 180px;
  background: var(--apple-bg-primary);
  border-radius: var(--apple-radius-md);
  box-shadow: var(--apple-shadow-lg);
  border: 1px solid var(--apple-separator);
  padding: var(--apple-spacing-xs);
  animation: apple-fade-in-down var(--apple-duration-fast) var(--apple-ease-default);
  z-index: 200;
}

.dropdown-item {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-primary);
  border-radius: var(--apple-radius-sm);
  cursor: pointer;
  transition: background var(--apple-duration-fast) var(--apple-ease-default);
}

.dropdown-item:hover {
  background: var(--apple-fill-quaternary);
}

.dropdown-item svg {
  width: 16px;
  height: 16px;
  color: var(--apple-text-secondary);
}

.dropdown-item.danger {
  color: var(--apple-red);
}

.dropdown-item.danger svg {
  color: var(--apple-red);
}

.dropdown-divider {
  height: 1px;
  background: var(--apple-separator);
  margin: var(--apple-spacing-xs) 0;
}

/* 主体区域 */
.main-wrapper {
  flex: 1;
  display: flex;
  overflow: hidden;
}

/* 侧边栏 */
.sidebar {
  width: 220px;
  background: var(--apple-bg-primary);
  border-right: 1px solid var(--apple-separator);
  display: flex;
  flex-direction: column;
  transition: width var(--apple-duration-normal) var(--apple-ease-default);
}

.sidebar.is-collapsed {
  width: 72px;
}

.sidebar-nav {
  flex: 1;
  padding: var(--apple-spacing-md);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--apple-spacing-sm);
  padding: var(--apple-spacing-sm) var(--apple-spacing-md);
  font-size: var(--apple-text-sm);
  color: var(--apple-text-secondary);
  border-radius: var(--apple-radius-md);
  cursor: pointer;
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
  text-decoration: none;
  margin-bottom: var(--apple-spacing-xs);
}

.nav-item:hover {
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-primary);
}

.nav-item.active {
  background: var(--apple-blue-light);
  color: var(--apple-blue);
}

.nav-icon {
  width: 20px;
  height: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.nav-icon svg {
  width: 18px;
  height: 18px;
}

.nav-label {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 侧边栏底部 */
.sidebar-footer {
  padding: var(--apple-spacing-md);
  border-top: 1px solid var(--apple-separator);
}

.collapse-toggle {
  width: 100%;
  padding: var(--apple-spacing-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--apple-text-tertiary);
  border-radius: var(--apple-radius-md);
  transition: all var(--apple-duration-fast) var(--apple-ease-default);
}

.collapse-toggle:hover {
  background: var(--apple-fill-quaternary);
  color: var(--apple-text-primary);
}

.collapse-toggle svg {
  width: 18px;
  height: 18px;
}

/* 折叠状态 */
.is-collapsed .nav-item {
  justify-content: center;
  padding: var(--apple-spacing-sm);
}

/* 主内容区 */
.main-content {
  flex: 1;
  padding: var(--apple-spacing-xl);
  overflow-y: auto;
}

/* 页面切换动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity var(--apple-duration-fast) var(--apple-ease-default);
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
