<template>
  <!-- 根组件：四区域布局 -->
  <div
    :class="[
      'app-root h-screen w-screen overflow-hidden flex flex-col',
      themeClass,
      { 'mica-effect': isMicaTheme }
    ]"
  >
    <!-- ====== 顶部功能栏 (48px) ====== -->
    <header
      class="topbar flex items-center justify-between px-3 shrink-0 select-none border-b"
      style="height: 48px; --wails-draggable: drag"
    >
      <!-- 左侧：Logo + 名称 -->
      <div class="flex items-center gap-2.5" style="--wails-draggable: no-drag">
        <div class="logo-icon w-7 h-7 rounded-lg flex items-center justify-center">
          <span class="text-white text-sm font-bold tracking-tight">X</span>
        </div>
        <span class="text-sm font-semibold tracking-wide opacity-90">XTool</span>
      </div>

      <!-- 中间：全局搜索框 -->
      <div
        class="search-trigger flex items-center gap-2 px-3 py-1.5 rounded-lg cursor-pointer mx-4"
        style="--wails-draggable: no-drag"
        @click="appStore.searchVisible = true"
      >
        <Search :size="14" class="opacity-40" />
        <span class="text-xs opacity-35 whitespace-nowrap">搜索工具...</span>
        <kbd class="text-[10px] opacity-30 px-1.5 py-0.5 rounded border border-current/20 ml-4">Ctrl+K</kbd>
      </div>

      <!-- 右侧：主题切换 + 窗口控制 -->
      <div class="flex items-center gap-1" style="--wails-draggable: no-drag">
        <!-- 主题切换下拉 -->
        <div class="relative" ref="themeDropdownRef">
          <button
            @click="showThemeDropdown = !showThemeDropdown"
            class="theme-btn flex items-center gap-1.5 px-2 py-1 rounded-md transition-colors"
            :title="'当前主题: ' + currentThemeLabel"
          >
            <component :is="currentThemeIcon" :size="14" class="opacity-70" />
            <span class="text-xs opacity-60 hidden sm:inline">{{ currentThemeLabel }}</span>
            <ChevronDown :size="12" class="opacity-40" />
          </button>
          <!-- 下拉面板 -->
          <Transition name="dropdown">
            <div
              v-if="showThemeDropdown"
              class="theme-dropdown absolute right-0 top-full mt-1 py-1 rounded-lg shadow-xl z-50 min-w-[140px]"
            >
              <button
                v-for="t in themeOptions"
                :key="t.id"
                @click="selectTheme(t.id)"
                class="theme-option flex items-center gap-2 w-full px-3 py-1.5 text-left transition-colors"
                :class="{ active: appStore.theme === t.id }"
              >
                <component :is="t.icon" :size="14" />
                <span class="text-xs">{{ t.label }}</span>
                <Check v-if="appStore.theme === t.id" :size="12" class="ml-auto" />
              </button>
            </div>
          </Transition>
        </div>

        <!-- 分隔线 -->
        <div class="w-px h-4 mx-1 opacity-10 bg-current" />

        <!-- 最小化 -->
        <button @click="minimizeWindow" class="win-btn hover:bg-white/10" title="最小化">
          <Minus :size="13" />
        </button>
        <!-- 最大化/还原 -->
        <button @click="toggleMaximize" class="win-btn hover:bg-white/10" title="最大化">
          <Square :size="12" />
        </button>
        <!-- 关闭 -->
        <button @click="closeWindow" class="win-btn hover:bg-red-500/80" title="关闭">
          <X :size="13" />
        </button>
      </div>
    </header>

    <!-- ====== 中间区域 (flex 横向布局) ====== -->
    <div class="flex flex-1 overflow-hidden">
      <!-- 左侧：SideNav 组件 -->
      <SideNav />

      <!-- 右侧：RouterView 内容区 -->
      <main class="flex-1 overflow-hidden relative">
        <RouterView v-slot="{ Component, route }">
          <Transition name="page" mode="out-in">
            <component :is="Component" :key="route.path" />
          </Transition>
        </RouterView>
      </main>
    </div>

    <!-- ====== 底部状态栏 (28px) ====== -->
    <footer
      class="statusbar flex items-center justify-between px-3 shrink-0 select-none border-t"
      style="height: 28px"
    >
      <!-- 左侧：当前工具名称 -->
      <div class="flex items-center gap-1.5">
        <component :is="currentToolIcon" :size="12" class="opacity-50" />
        <span class="text-[11px] opacity-50">{{ currentToolName }}</span>
      </div>
      <!-- 右侧：版本号 -->
      <div class="flex items-center gap-2">
        <span class="text-[10px] opacity-30">Wails v2</span>
        <span class="text-[10px] opacity-30">|</span>
        <span class="text-[10px] opacity-30">v{{ appVersion }}</span>
      </div>
    </footer>

    <!-- 全局通知 Toast -->
    <ToastContainer />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, markRaw } from 'vue'
import { RouterView, useRoute } from 'vue-router'
import {
  Minus, Square, X, Search, ChevronDown, Check,
  Moon, Sun, Monitor, Palette, Leaf, Sparkles, Layers,
  Braces, Code2, Binary, Link, Hash, FileText, QrCode, GitCompare,
  Fingerprint, Clock, Regex, BookMarked, Shield, Keyboard,
  Cpu, Network, MonitorIcon, FilePen, Pipette, Image, FolderSearch,
  Calculator, ArrowLeftRight, StickyNote, KeyRound,
  Signal, Radar, Globe, Globe2, Server, Settings,
  Database, FileCode, Lock, Trash2, Clipboard, SearchCode, Gauge,
} from 'lucide-vue-next'
import SideNav from '@/components/SideNav.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { useAppStore } from '@/stores/app'
import { WindowMinimise, WindowToggleMaximise, Quit } from '../wailsjs/runtime/runtime'

const appStore = useAppStore()
const route = useRoute()

// ============ 主题相关 ============

// 主题选项配置
const themeOptions = [
  { id: 'dark',   label: '深色',   icon: markRaw(Moon) },
  { id: 'light',  label: '浅色',   icon: markRaw(Sun) },
  { id: 'blue',   label: '蓝色',   icon: markRaw(Palette) },
  { id: 'green',  label: '绿色',   icon: markRaw(Leaf) },
  { id: 'purple', label: '紫色',   icon: markRaw(Sparkles) },
  { id: 'mica',   label: '云母',   icon: markRaw(Layers) },
  { id: 'auto',   label: '跟随系统', icon: markRaw(Monitor) },
]

// 当前主题标签和图标
const currentThemeLabel = computed(() => {
  const found = themeOptions.find(t => t.id === appStore.theme)
  return found ? found.label : '深色'
})

const currentThemeIcon = computed(() => {
  const found = themeOptions.find(t => t.id === appStore.theme)
  return found ? found.icon : Moon
})

// 根据主题计算类名
const themeClass = computed(() => {
  const theme = appStore.theme
  if (theme === 'dark' || theme === 'blue' || theme === 'purple' || theme === 'mica') return 'dark'
  if (theme === 'light') return 'light'
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
})

// 是否启用云母效果
const isMicaTheme = computed(() => appStore.theme === 'mica')

// 主题下拉菜单状态
const showThemeDropdown = ref(false)
const themeDropdownRef = ref<HTMLElement | null>(null)

// 选择主题
function selectTheme(themeId: string) {
  appStore.switchTheme(themeId)
  showThemeDropdown.value = false
}

// 点击外部关闭下拉菜单
function handleClickOutside(e: MouseEvent) {
  if (themeDropdownRef.value && !themeDropdownRef.value.contains(e.target as Node)) {
    showThemeDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside)
})

// ============ 当前工具信息 ============

// 图标映射表（与 SideNav 保持一致）
const iconMap: Record<string, any> = {
  Braces, Code2, Binary, Link, Hash, FileText, QrCode, GitCompare,
  Fingerprint, Clock, Regex, BookMarked, Shield, Palette, Keyboard,
  Cpu, Network, Monitor: MonitorIcon, FilePen, Pipette, Image, FolderSearch,
  Calculator, ArrowLeftRight, StickyNote,
  Signal, Radar, Globe, Globe2, Server, Settings,
  Database, FileCode, Lock, Trash2, Clipboard, KeyRound, SearchCode, Gauge,
}

// 当前工具名称
const currentToolName = computed(() => {
  return (route.meta?.title as string) || 'XTool'
})

// 当前工具图标
const currentToolIcon = computed(() => {
  const iconName = route.meta?.icon as string
  return iconName ? iconMap[iconName] || FileText : Braces
})

// ============ 版本号 ============

const appVersion = ref('1.0.0')

// ============ 窗口控制 ============

const minimizeWindow = () => WindowMinimise()
const toggleMaximize = () => WindowToggleMaximise()
const closeWindow = () => {
  if (appStore.config.minToTray) {
    WindowMinimise()
  } else {
    Quit()
  }
}

// ============ 初始化 ============

onMounted(async () => {
  await appStore.loadConfig()
  // 尝试从 wails.json 获取版本号
  try {
    // @ts-ignore - Wails runtime global
    if (window.runtime) {
      // 版本号已嵌入构建，这里使用默认值
    }
  } catch {
    // ignore
  }
})
</script>

<style>
/* ============================================================
   CSS 变量驱动主题系统
   ============================================================ */

:root {
  /* 过渡时长 */
  --transition-theme: 0.3s ease;

  /* 深色主题变量（默认） */
  --bg-primary: #0f0f14;
  --bg-secondary: #16161d;
  --bg-tertiary: #1c1c26;
  --bg-hover: rgba(255, 255, 255, 0.06);
  --bg-active: rgba(99, 102, 241, 0.15);
  --text-primary: #e2e8f0;
  --text-secondary: #94a3b8;
  --text-muted: #475569;
  --border-color: rgba(255, 255, 255, 0.06);
  --accent: #6366f1;
  --accent-hover: #818cf8;
  --accent-bg: rgba(99, 102, 241, 0.12);
  --shadow: 0 4px 24px rgba(0, 0, 0, 0.4);
  --statusbar-bg: rgba(255, 255, 255, 0.02);
}

/* 浅色主题 */
.light {
  --bg-primary: #f8fafc;
  --bg-secondary: #f1f5f9;
  --bg-tertiary: #e2e8f0;
  --bg-hover: rgba(0, 0, 0, 0.05);
  --bg-active: rgba(79, 70, 229, 0.1);
  --text-primary: #1e293b;
  --text-secondary: #475569;
  --text-muted: #94a3b8;
  --border-color: rgba(0, 0, 0, 0.08);
  --accent: #4f46e5;
  --accent-hover: #6366f1;
  --accent-bg: rgba(79, 70, 229, 0.08);
  --shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  --statusbar-bg: rgba(0, 0, 0, 0.02);
}

/* 蓝色主题 */
.theme-blue {
  --accent: #3b82f6;
  --accent-hover: #60a5fa;
  --accent-bg: rgba(59, 130, 246, 0.12);
  --bg-primary: #0a0e1a;
  --bg-secondary: #0f1525;
  --bg-tertiary: #141c30;
}

/* 绿色主题 */
.theme-green {
  --accent: #22c55e;
  --accent-hover: #4ade80;
  --accent-bg: rgba(34, 197, 94, 0.12);
  --bg-primary: #0a140e;
  --bg-secondary: #0f1a13;
  --bg-tertiary: #14201a;
}

/* 紫色主题 */
.theme-purple {
  --accent: #a855f7;
  --accent-hover: #c084fc;
  --accent-bg: rgba(168, 85, 247, 0.12);
  --bg-primary: #120a1a;
  --bg-secondary: #170f22;
  --bg-tertiary: #1c1428;
}

/* 云母主题 */
.theme-mica {
  --accent: #6366f1;
  --accent-hover: #818cf8;
  --accent-bg: rgba(99, 102, 241, 0.15);
}

/* ============================================================
   全局样式
   ============================================================ */

.app-root {
  font-family: "Microsoft YaHei", "PingFang SC", system-ui, -apple-system, sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: background var(--transition-theme), color var(--transition-theme);
}

/* ============================================================
   顶部功能栏
   ============================================================ */

.topbar {
  background: var(--bg-secondary);
  border-color: var(--border-color);
  transition: background var(--transition-theme), border-color var(--transition-theme);
}

/* Logo 图标 */
.logo-icon {
  background: linear-gradient(135deg, var(--accent), var(--accent-hover));
  box-shadow: 0 2px 8px rgba(99, 102, 241, 0.3);
}

/* 搜索触发器 */
.search-trigger {
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  color: var(--text-secondary);
  transition: all 0.2s ease;
  min-width: 200px;
  max-width: 320px;
}

.search-trigger:hover {
  background: var(--bg-active);
  border-color: var(--accent);
}

.search-trigger kbd {
  color: var(--text-muted);
}

/* 主题按钮 */
.theme-btn {
  background: transparent;
  border: none;
  color: inherit;
  cursor: pointer;
  border-radius: 6px;
  padding: 4px 8px;
  transition: background 0.15s ease;
}

.theme-btn:hover {
  background: var(--bg-hover);
}

/* 主题下拉面板 */
.theme-dropdown {
  background: var(--bg-tertiary);
  border: 1px solid var(--border-color);
  box-shadow: var(--shadow);
  backdrop-filter: blur(12px);
}

.theme-option {
  background: transparent;
  border: none;
  color: var(--text-secondary);
  cursor: pointer;
  font-size: 12px;
  transition: all 0.15s ease;
}

.theme-option:hover {
  background: var(--bg-hover);
  color: var(--text-primary);
}

.theme-option.active {
  color: var(--accent);
}

/* 窗口控制按钮 */
.win-btn {
  width: 32px;
  height: 28px;
  border-radius: 6px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  cursor: pointer;
  color: inherit;
  transition: background 0.15s ease;
}

.light .win-btn:hover {
  background: rgba(0, 0, 0, 0.06);
}

/* ============================================================
   底部状态栏
   ============================================================ */

.statusbar {
  background: var(--statusbar-bg);
  border-color: var(--border-color);
  color: var(--text-muted);
  transition: background var(--transition-theme), border-color var(--transition-theme);
}

/* ============================================================
   云母半透明效果
   ============================================================ */

.mica-effect {
  background: transparent !important;
}

.mica-effect .topbar,
.mica-effect .statusbar {
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(20px) saturate(180%);
}

.light .mica-effect .topbar,
.light .mica-effect .statusbar {
  background: rgba(255, 255, 255, 0.7);
}

/* ============================================================
   页面切换动画（淡入淡出 + 微位移）
   ============================================================ */

.page-enter-active,
.page-leave-active {
  transition: all 0.25s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateY(6px);
}

.page-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* ============================================================
   下拉菜单动画
   ============================================================ */

.dropdown-enter-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.dropdown-leave-active {
  transition: all 0.15s ease-in;
}

.dropdown-enter-from {
  opacity: 0;
  transform: translateY(-4px) scale(0.96);
}

.dropdown-leave-to {
  opacity: 0;
  transform: translateY(-2px) scale(0.98);
}
</style>
