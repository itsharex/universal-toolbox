<template>
  <!-- 根组件：四区域布局 -->
  <div
    :class="[
      'app-root h-screen w-screen overflow-hidden flex flex-col',
      themeClass,
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

    <!-- ====== 全局搜索弹窗 ====== -->
    <Teleport to="body">
      <Transition name="search-modal">
        <div
          v-if="searchVisible"
          class="search-modal-overlay"
          @mousedown.self="closeSearch"
        >
          <div class="search-modal">
            <!-- 搜索输入框 -->
            <div class="search-input-wrapper">
              <Search :size="18" class="search-input-icon" />
              <input
                ref="searchInputRef"
                v-model="searchQuery"
                class="search-input"
                placeholder="搜索工具名称或描述..."
                @keydown.down.prevent="moveDown"
                @keydown.up.prevent="moveUp"
                @keydown.enter.prevent="selectCurrent"
                @keydown.esc.prevent="closeSearch"
              />
              <kbd class="kbd">ESC</kbd>
            </div>

            <!-- 搜索结果列表 -->
            <div class="search-results" v-if="searchQuery.trim()">
              <div
                v-for="(item, index) in filteredItems"
                :key="item.path"
                class="search-result-item"
                :class="{ active: index === activeIndex }"
                @click="navigateTo(item)"
                @mouseenter="activeIndex = index"
              >
                <component :is="getIcon(item.icon)" :size="16" class="search-result-icon" />
                <div class="search-result-info">
                  <span class="search-result-name">{{ item.name }}</span>
                  <span class="search-result-category">{{ item.category }}</span>
                </div>
                <component :is="ArrowRight" :size="14" class="search-result-arrow" />
              </div>
              <!-- 无结果 -->
              <div v-if="filteredItems.length === 0" class="empty-state">
                <Search :size="32" class="empty-state-icon" />
                <p class="empty-state-text">未找到匹配的工具</p>
                <p class="empty-state-hint">试试其他关键词</p>
              </div>
            </div>

            <!-- 默认提示（无搜索词时） -->
            <div class="search-results" v-else>
              <div class="search-hint">
                <span class="search-hint-text">输入关键词搜索所有工具</span>
              </div>
              <!-- 快捷操作 -->
              <div class="search-quick-list">
                <div
                  v-for="item in recentItems"
                  :key="item.path"
                  class="search-result-item"
                  @click="navigateTo(item)"
                >
                  <component :is="getIcon(item.icon)" :size="16" class="search-result-icon" />
                  <div class="search-result-info">
                    <span class="search-result-name">{{ item.name }}</span>
                    <span class="search-result-category">{{ item.category }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- 底部提示 -->
            <div class="search-footer">
              <span class="search-footer-item"><kbd class="kbd"><ArrowUp :size="10" /></kbd><kbd class="kbd"><ArrowDown :size="10" /></kbd> 导航</span>
              <span class="search-footer-item"><kbd class="kbd">Enter</kbd> 选择</span>
              <span class="search-footer-item"><kbd class="kbd">ESC</kbd> 关闭</span>
            </div>
          </div>
        </div>
      </Transition>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted, watch, nextTick, markRaw } from 'vue'
import { RouterView, useRoute, useRouter } from 'vue-router'
import {
  Minus, Square, X, Search, ChevronDown, Check,
  Moon, Sun, Monitor, Palette, Leaf, Sparkles, Layers,
  Braces, Code2, Binary, Link, Hash, FileText, QrCode, GitCompare,
  Fingerprint, Clock, Regex, BookMarked, Shield, Keyboard,
  Cpu, Network, MonitorIcon, FilePen, Pipette, Image, FolderSearch,
  Calculator, ArrowLeftRight, StickyNote, KeyRound,
  Signal, Radar, Globe, Globe2, Server, Settings,
  Database, FileCode, Lock, Trash2, Clipboard, SearchCode, Gauge,
  ArrowUp, ArrowDown, ArrowRight, Home, AlignLeft, Type, Timer,
  FileEdit, FileInput, Activity, HardDrive, Files, Info, Wifi,
} from 'lucide-vue-next'
import SideNav from '@/components/SideNav.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import { useAppStore } from '@/stores/app'
import { WindowMinimise, WindowToggleMaximise, Quit } from '../wailsjs/runtime/runtime'

const appStore = useAppStore()
const route = useRoute()
const router = useRouter()

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

const appVersion = ref('2.0.0')

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

// ============ 全局搜索 ============

// 搜索数据源（与 SideNav navGroups 保持同步）
const searchItems = [
  // devtools
  { name: 'JSON 工具', path: '/devtools/json', category: '开发工具', icon: 'Braces' },
  { name: 'Base64 编解码', path: '/devtools/base64', category: '开发工具', icon: 'Binary' },
  { name: '哈希计算', path: '/devtools/hash', category: '开发工具', icon: 'Hash' },
  { name: '正则测试', path: '/devtools/regex', category: '开发工具', icon: 'Regex' },
  { name: '代码格式化', path: '/devtools/formatter', category: '开发工具', icon: 'AlignLeft' },
  { name: '代码混淆', path: '/devtools/obfuscator', category: '开发工具', icon: 'Shield' },
  { name: 'UUID 生成', path: '/devtools/uuid', category: '开发工具', icon: 'Fingerprint' },
  { name: '时间戳转换', path: '/devtools/timestamp', category: '开发工具', icon: 'Clock' },
  { name: 'URL 编解码', path: '/devtools/url', category: '开发工具', icon: 'Link' },
  { name: '加密解密', path: '/devtools/crypto', category: '开发工具', icon: 'Lock' },
  { name: 'XML 工具', path: '/devtools/xml', category: '开发工具', icon: 'FileCode' },
  { name: '二维码', path: '/devtools/qrcode', category: '开发工具', icon: 'QrCode' },
  { name: '文本处理', path: '/devtools/text', category: '开发工具', icon: 'Type' },
  { name: '进制转换', path: '/devtools/radix', category: '开发工具', icon: 'Binary' },
  { name: '占位文本生成', path: '/devtools/dummydata', category: '开发工具', icon: 'Database' },
  { name: '接口文档生成', path: '/devtools/apidoc', category: '开发工具', icon: 'FileText' },
  { name: '代码片段', path: '/devtools/snippets', category: '开发工具', icon: 'Code2' },
  { name: 'JSON 对比', path: '/devtools/jsondiff', category: '开发工具', icon: 'GitCompare' },
  // sysinfo
  { name: '系统信息', path: '/sysinfo/system', category: '系统工具', icon: 'Monitor' },
  { name: '进程管理', path: '/sysinfo/process', category: '系统工具', icon: 'Activity' },
  { name: '端口管理', path: '/sysinfo/ports', category: '系统工具', icon: 'Network' },
  { name: '磁盘清理', path: '/sysinfo/diskcleaner', category: '系统工具', icon: 'HardDrive' },
  { name: '批量重命名', path: '/sysinfo/batchrename', category: '系统工具', icon: 'FilePen' },
  { name: '文件批量处理', path: '/sysinfo/filebatch', category: '系统工具', icon: 'Files' },
  { name: '图片工具', path: '/sysinfo/image', category: '系统工具', icon: 'Image' },
  { name: '取色器', path: '/sysinfo/colorpicker', category: '系统工具', icon: 'Pipette' },
  { name: '剪贴板管理', path: '/sysinfo/clipboard', category: '系统工具', icon: 'Clipboard' },
  { name: '定时任务', path: '/sysinfo/cron', category: '系统工具', icon: 'Timer' },
  // daily
  { name: '计算器', path: '/daily/calculator', category: '日常工具', icon: 'Calculator' },
  { name: '单位换算', path: '/daily/converter', category: '日常工具', icon: 'ArrowLeftRight' },
  { name: '密码生成', path: '/daily/password', category: '日常工具', icon: 'KeyRound' },
  { name: '备忘录', path: '/daily/notes', category: '日常工具', icon: 'StickyNote' },
  { name: '二维码批量处理', path: '/daily/qrbatch', category: '日常工具', icon: 'QrCode' },
  { name: '文档转换', path: '/daily/docconverter', category: '日常工具', icon: 'FileInput' },
  // network
  { name: 'HTTP 测试', path: '/network/http', category: '网络工具', icon: 'Globe' },
  { name: 'Ping 测试', path: '/network/ping', category: '网络工具', icon: 'Wifi' },
  { name: 'DNS 查询', path: '/network/dns', category: '网络工具', icon: 'Search' },
  { name: 'WHOIS 查询', path: '/network/whois', category: '网络工具', icon: 'Info' },
  { name: '内网扫描', path: '/network/scan', category: '网络工具', icon: 'Radar' },
  { name: '网络测速', path: '/network/speedtest', category: '网络工具', icon: 'Gauge' },
  { name: '端口转发', path: '/network/portforward', category: '网络工具', icon: 'ArrowLeftRight' },
  { name: 'Hosts 编辑', path: '/network/hosts', category: '网络工具', icon: 'FileEdit' },
  // settings
  { name: '设置', path: '/settings', category: '设置', icon: 'Settings' },
  { name: '快捷键管理', path: '/settings/shortcuts', category: '设置', icon: 'Keyboard' },
  { name: '主题编辑', path: '/settings/theme', category: '设置', icon: 'Palette' },
]

// 搜索图标映射（扩展 iconMap）
const searchIconMap: Record<string, any> = {
  ...iconMap,
  AlignLeft, Type, Timer, FileEdit, FileInput, Activity, HardDrive, Files, Info, Wifi,
}

// 获取搜索项图标组件
function getIcon(iconName: string) {
  return searchIconMap[iconName] || FileText
}

// 搜索弹窗状态
const searchVisible = ref(false)
const searchQuery = ref('')
const activeIndex = ref(0)
const searchInputRef = ref<HTMLInputElement | null>(null)

// 打开搜索弹窗
function openSearch() {
  searchVisible.value = true
  searchQuery.value = ''
  activeIndex.value = 0
  nextTick(() => {
    searchInputRef.value?.focus()
  })
}

// 关闭搜索弹窗
function closeSearch() {
  searchVisible.value = false
  searchQuery.value = ''
  activeIndex.value = 0
}

// 搜索过滤结果
const filteredItems = computed(() => {
  const query = searchQuery.value.trim().toLowerCase()
  if (!query) return []
  return searchItems.filter(item =>
    item.name.toLowerCase().includes(query) ||
    item.category.toLowerCase().includes(query) ||
    item.path.toLowerCase().includes(query)
  )
})

// 默认显示的常用工具（无搜索词时）
const recentItems = computed(() => {
  return searchItems.slice(0, 6)
})

// 监听过滤结果变化，重置 activeIndex
watch(filteredItems, () => {
  activeIndex.value = 0
})

// 键盘导航：向下
function moveDown() {
  if (filteredItems.value.length > 0) {
    activeIndex.value = (activeIndex.value + 1) % filteredItems.value.length
  }
}

// 键盘导航：向上
function moveUp() {
  if (filteredItems.value.length > 0) {
    activeIndex.value = (activeIndex.value - 1 + filteredItems.value.length) % filteredItems.value.length
  }
}

// 选择当前项
function selectCurrent() {
  const query = searchQuery.value.trim()
  if (query && filteredItems.value.length > 0) {
    navigateTo(filteredItems.value[activeIndex.value])
  }
}

// 跳转到对应工具
function navigateTo(item: typeof searchItems[0]) {
  closeSearch()
  router.push(item.path)
}

// Ctrl+K 全局快捷键
function handleGlobalKeydown(e: KeyboardEvent) {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    if (searchVisible.value) {
      closeSearch()
    } else {
      openSearch()
    }
  }
  // ESC 关闭搜索弹窗
  if (e.key === 'Escape' && searchVisible.value) {
    closeSearch()
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleGlobalKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleGlobalKeydown)
})
</script>

<style>
/* ============================================================
   全局样式
   ============================================================ */

.app-root {
  font-family: "Microsoft YaHei", "PingFang SC", system-ui, -apple-system, sans-serif;
  background: var(--bg-primary);
  color: var(--text-primary);
  transition: background var(--transition-theme), color var(--transition-theme);
  border: 1px solid var(--border-color);
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.05), 0 8px 32px rgba(0, 0, 0, 0.12);
}

.light .app-root {
  border: 1px solid rgba(0, 0, 0, 0.12);
  box-shadow: 0 0 0 1px rgba(0, 0, 0, 0.08), 0 8px 32px rgba(0, 0, 0, 0.15);
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

/* ============================================================
   搜索弹窗动画
   ============================================================ */

.search-modal-enter-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.search-modal-leave-active {
  transition: all 0.15s ease-in;
}

.search-modal-enter-from {
  opacity: 0;
}

.search-modal-enter-from .search-modal {
  transform: scale(0.95) translateY(-8px);
  opacity: 0;
}

.search-modal-leave-to {
  opacity: 0;
}

.search-modal-leave-to .search-modal {
  transform: scale(0.97) translateY(-4px);
  opacity: 0;
}
</style>
