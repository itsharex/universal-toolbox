// 全局应用状态管理（Pinia Store）
// 管理主题、配置、全局通知等应用级状态

import { defineStore } from 'pinia'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { GetConfig, SaveConfig, GetTheme, SetTheme } from '../../wailsjs/go/config/Config'

// 应用配置接口定义
interface AppConfig {
  theme: string
  language: string
  fontSize: number
  density: string
  transparency: number
  alwaysOnTop: boolean
  minToTray: boolean
  autoBackup: boolean
}

// Toast 通知接口
export interface Toast {
  id: string
  type: 'success' | 'error' | 'warning' | 'info'
  message: string
  duration?: number
}

// 所有已知的主题类名，用于切换时彻底清除
const ALL_THEME_CLASSES = [
  'dark', 'light',
  'theme-blue', 'theme-green', 'theme-purple', 'theme-orange', 'theme-mica',
]

// 系统主题变化的媒体查询
let systemThemeQuery: MediaQueryList | null = null
let systemThemeHandler: ((e: MediaQueryListEvent) => void) | null = null

export const useAppStore = defineStore('app', () => {
  // ============ 状态定义 ============

  // 应用配置
  const config = ref<AppConfig>({
    theme: 'dark',
    language: 'zh-CN',
    fontSize: 14,
    density: 'normal',
    transparency: 0.95,
    alwaysOnTop: false,
    minToTray: true,
    autoBackup: true,
  })

  // 当前主题
  const theme = computed(() => config.value.theme)

  // 全局搜索显示状态
  const searchVisible = ref(false)

  // Toast 通知列表
  const toasts = ref<Toast[]>([])

  // 侧边栏折叠状态
  const sidebarCollapsed = ref(false)

  // ============ Actions ============

  // 从后端加载配置
  async function loadConfig() {
    try {
      const cfg = await GetConfig()
      if (cfg) {
        config.value = cfg as AppConfig
        applyTheme(config.value.theme)
      }
    } catch (err) {
      console.error('加载配置失败:', err)
    }
    // 初始化系统主题监听
    initSystemThemeListener()
  }

  // 保存配置到后端
  async function saveConfig(newConfig: Partial<AppConfig>) {
    config.value = { ...config.value, ...newConfig }
    try {
      await SaveConfig(JSON.stringify(config.value))
      applyTheme(config.value.theme)
    } catch (err) {
      showToast('error', '保存配置失败: ' + String(err))
    }
  }

  // 切换主题
  async function switchTheme(themeId: string) {
    config.value.theme = themeId
    applyTheme(themeId)
    try {
      await SetTheme(themeId)
    } catch (err) {
      console.error('切换主题失败:', err)
    }
  }

  // 快捷切换深色/浅色主题
  function toggleTheme() {
    const current = config.value.theme
    if (current === 'dark') {
      switchTheme('light')
    } else if (current === 'light') {
      switchTheme('dark')
    } else {
      // 如果当前是彩色主题（blue/green/purple/orange/mica），切换到 light
      switchTheme('light')
    }
  }

  // 应用主题到 DOM（增强版）
  function applyTheme(themeId: string) {
    const root = document.documentElement

    // 1. 清除所有主题类和效果类
    ALL_THEME_CLASSES.forEach(cls => root.classList.remove(cls))
    root.classList.remove('mica-effect')  // 确保清除云母效果

    // 2. 设置 data-theme 属性
    root.setAttribute('data-theme', themeId)

    // 3. 添加对应的类
    if (themeId === 'auto') {
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      root.classList.add(prefersDark ? 'dark' : 'light')
    } else if (themeId === 'light') {
      root.classList.add('light')
    } else {
      root.classList.add('dark')
      if (themeId !== 'dark') {
        root.classList.add(`theme-${themeId}`)
      }
      if (themeId === 'mica') {
        root.classList.add('mica-effect')
      }
    }

    // 4. 主题过渡动画
    root.classList.add('theme-transitioning')
    setTimeout(() => {
      root.classList.remove('theme-transitioning')
    }, 300)

    updateScrollbarStyle()
  }

  // 更新滚动条样式，使其跟随当前主题
  function updateScrollbarStyle() {
    let styleEl = document.getElementById('dynamic-scrollbar-style') as HTMLStyleElement | null
    if (!styleEl) {
      styleEl = document.createElement('style')
      styleEl.id = 'dynamic-scrollbar-style'
      document.head.appendChild(styleEl)
    }
    styleEl.textContent = `
      ::-webkit-scrollbar { width: 6px; height: 6px; }
      ::-webkit-scrollbar-track { background: var(--bg-secondary, transparent); }
      ::-webkit-scrollbar-thumb { background: var(--border-color, rgba(128,128,128,0.3)); border-radius: 3px; }
      ::-webkit-scrollbar-thumb:hover { background: var(--text-muted, rgba(128,128,128,0.5)); }
    `
  }

  // 初始化系统主题变化监听
  function initSystemThemeListener() {
    if (systemThemeQuery) return // 避免重复监听
    systemThemeQuery = window.matchMedia('(prefers-color-scheme: dark)')
    systemThemeHandler = (e: MediaQueryListEvent) => {
      // 仅在 auto 模式下响应系统主题变化
      if (config.value.theme === 'auto') {
        applyTheme('auto')
      }
    }
    systemThemeQuery.addEventListener('change', systemThemeHandler)
  }

  // 清理系统主题监听
  function cleanupSystemThemeListener() {
    if (systemThemeQuery && systemThemeHandler) {
      systemThemeQuery.removeEventListener('change', systemThemeHandler)
      systemThemeQuery = null
      systemThemeHandler = null
    }
  }

  // 显示 Toast 通知
  function showToast(type: Toast['type'], message: string, duration = 3000) {
    const id = Date.now().toString() + Math.random().toString(36).slice(2)
    const toast: Toast = { id, type, message, duration }
    toasts.value.push(toast)

    // 自动移除
    if (duration > 0) {
      setTimeout(() => removeToast(id), duration)
    }
  }

  // 移除 Toast
  function removeToast(id: string) {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index !== -1) toasts.value.splice(index, 1)
  }

  // 切换侧边栏折叠
  function toggleSidebar() {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  return {
    config,
    theme,
    searchVisible,
    toasts,
    sidebarCollapsed,
    loadConfig,
    saveConfig,
    switchTheme,
    toggleTheme,
    applyTheme,
    showToast,
    removeToast,
    toggleSidebar,
    cleanupSystemThemeListener,
  }
})
