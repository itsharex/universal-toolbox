<template>
  <!-- 左侧导航栏组件 -->
  <nav
    :class="[
      'sidenav flex flex-col shrink-0 transition-all duration-300 ease-in-out relative',
      collapsed ? 'w-14' : 'w-56'
    ]"
  >
    <!-- 搜索框 -->
    <div v-if="!collapsed" class="px-3 py-2">
      <div
        class="search-box flex items-center gap-2 px-3 py-1.5 rounded-lg cursor-pointer"
        @click="appStore.searchVisible = true"
      >
        <Search :size="13" class="opacity-50 shrink-0" />
        <span class="text-xs opacity-40">搜索工具... <kbd class="search-kbd">Ctrl+K</kbd></span>
      </div>
    </div>
    <!-- 折叠模式搜索按钮 -->
    <div v-else class="px-2 py-2 flex justify-center">
      <button
        class="search-btn-collapsed flex items-center justify-center w-8 h-8 rounded-lg cursor-pointer"
        @click="appStore.searchVisible = true"
        v-tooltip="'搜索工具 (Ctrl+K)'"
      >
        <Search :size="15" class="opacity-60" />
      </button>
    </div>

    <!-- 导航分组 -->
    <div class="flex-1 overflow-y-auto overflow-x-hidden px-2 py-1 space-y-0.5 scrollbar-thin">
      <template v-for="group in navGroups" :key="group.id">
        <!-- 分组标题 -->
        <div
          v-if="!collapsed"
          class="nav-group-title px-2 py-1 mt-3 first:mt-1 flex items-center gap-2"
        >
          <span
            class="w-1.5 h-1.5 rounded-full shrink-0"
            :style="{ background: getCategoryColor(group.id) }"
          />
          {{ group.label }}
        </div>
        <div v-else class="h-px mx-2 my-2 opacity-10 bg-current" />

        <!-- 导航项 -->
        <RouterLink
          v-for="item in group.items"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{
            'active': isActive(item.path),
            'collapsed': collapsed,
          }"
          :style="getActiveStyle(group.id, item.path)"
          v-tooltip="collapsed ? item.label : ''"
        >
          <component
            :is="getIcon(item.icon)"
            :size="16"
            class="shrink-0 nav-icon"
            :style="isActive(item.path) ? { color: getCategoryColor(group.id) } : {}"
          />
          <span v-if="!collapsed" class="truncate text-sm">{{ item.label }}</span>
        </RouterLink>
      </template>
    </div>

    <!-- 底部区域 -->
    <div class="sidenav-footer px-2 py-2 border-t border-white/5 space-y-0.5">
      <!-- 设置 -->
      <RouterLink
        to="/settings"
        class="nav-item"
        :class="{ 'active': isActive('/settings'), 'collapsed': collapsed }"
        :style="getActiveStyle('settings', '/settings')"
      >
        <Settings
          :size="16"
          class="shrink-0 nav-icon"
          :style="isActive('/settings') ? { color: getCategoryColor('settings') } : {}"
        />
        <span v-if="!collapsed" class="text-sm">设置</span>
      </RouterLink>
      <!-- 主题编辑器 -->
      <RouterLink
        to="/settings/theme"
        class="nav-item"
        :class="{ 'active': isActive('/settings/theme'), 'collapsed': collapsed }"
        :style="getActiveStyle('settings', '/settings/theme')"
      >
        <Palette
          :size="16"
          class="shrink-0 nav-icon"
          :style="isActive('/settings/theme') ? { color: getCategoryColor('settings') } : {}"
        />
        <span v-if="!collapsed" class="text-sm">主题编辑</span>
      </RouterLink>
      <!-- 快捷键管理 -->
      <RouterLink
        to="/settings/shortcuts"
        class="nav-item"
        :class="{ 'active': isActive('/settings/shortcuts'), 'collapsed': collapsed }"
        :style="getActiveStyle('settings', '/settings/shortcuts')"
      >
        <Keyboard
          :size="16"
          class="shrink-0 nav-icon"
          :style="isActive('/settings/shortcuts') ? { color: getCategoryColor('settings') } : {}"
        />
        <span v-if="!collapsed" class="text-sm">快捷键</span>
      </RouterLink>

      <!-- 折叠按钮 -->
      <button
        @click="appStore.toggleSidebar()"
        class="nav-item w-full"
      >
        <PanelLeftClose v-if="!collapsed" :size="16" class="shrink-0" />
        <PanelLeftOpen v-else :size="16" class="shrink-0" />
        <span v-if="!collapsed" class="text-sm">收起侧栏</span>
      </button>

      <!-- 版本号 + 版权信息 -->
      <div v-if="!collapsed" class="px-2 pt-2 pb-1">
        <div class="text-[10px] opacity-30 leading-relaxed">
          <div>Universal Toolbox</div>
          <div>v1.0.0</div>
          <div class="mt-0.5">&copy; 2024 All rights reserved</div>
        </div>
      </div>
      <div v-else class="flex justify-center pt-1 pb-1">
        <span class="text-[9px] opacity-25">v1.0</span>
      </div>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted } from 'vue'
import { RouterLink, useRoute } from 'vue-router'
import {
  Search, Settings, PanelLeftClose, PanelLeftOpen,
  Braces, Code2, Binary, Link, Hash, FileText, QrCode, GitCompare,
  Fingerprint, Clock, Regex, BookMarked, Shield, Palette, Keyboard,
  Cpu, Network, Monitor, FilePen, Pipette, Image, FolderSearch,
  Calculator, ArrowLeftRight, StickyNote,
  Signal, Radar, Globe, Globe2, Server,
  Database, FileCode, Lock, Trash2, Clipboard, KeyRound, SearchCode, Gauge,
} from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()
const route = useRoute()

// 侧边栏折叠状态
const collapsed = computed(() => appStore.sidebarCollapsed)

// 判断路由是否激活
const isActive = (path: string) => route.path.startsWith(path)

// 分类配色映射
const categoryColors: Record<string, string> = {
  devtools: '#00ADD8',   // 开发工具 - 蓝色
  sysinfo: '#42b883',    // 系统工具 - 绿色
  daily: '#FF9C41',      // 日常工具 - 橙色
  network: '#8B5CF6',    // 网络工具 - 紫色
  settings: '#6B6B7B',   // 设置 - 灰色
}

// 获取分类颜色
const getCategoryColor = (groupId: string): string => {
  return categoryColors[groupId] || '#6B6B7B'
}

// 获取激活状态的内联样式（左侧指示条颜色）
const getActiveStyle = (groupId: string, path: string) => {
  if (!isActive(path)) return {}
  const color = getCategoryColor(groupId)
  return {
    '--cat-color': color,
  }
}

// 图标映射表
const iconMap: Record<string, any> = {
  Braces, Code2, Binary, Link, Hash, FileText, QrCode, GitCompare,
  Fingerprint, Clock, Regex, BookMarked, Shield, Palette, Keyboard,
  Cpu, Network, Monitor, FilePen, Pipette, Image, FolderSearch,
  Calculator, ArrowLeftRight, StickyNote,
  Signal, Radar, Globe, Globe2, Server, Settings,
  Database, FileCode, Lock, Trash2, Clipboard, KeyRound, SearchCode, Gauge,
}
const getIcon = (name: string) => iconMap[name] || FileText

// 导航分组配置
const navGroups = [
  {
    id: 'devtools',
    label: '开发工具',
    items: [
      { path: '/devtools/json',      label: 'JSON 工具',   icon: 'Braces' },
      { path: '/devtools/jsondiff',  label: 'JSON 对比',   icon: 'GitCompare' },
      { path: '/devtools/xml',       label: 'XML 工具',    icon: 'Code2' },
      { path: '/devtools/base64',    label: 'Base64',      icon: 'Binary' },
      { path: '/devtools/url',       label: 'URL 编解码',  icon: 'Link' },
      { path: '/devtools/hash',      label: '哈希计算',    icon: 'Hash' },
      { path: '/devtools/crypto',    label: '加密解密',    icon: 'Shield' },
      { path: '/devtools/text',      label: '文本工具',    icon: 'FileText' },
      { path: '/devtools/qrcode',    label: '二维码',      icon: 'QrCode' },
      { path: '/devtools/uuid',      label: 'UUID 生成',   icon: 'Fingerprint' },
      { path: '/devtools/timestamp', label: '时间戳',      icon: 'Clock' },
      { path: '/devtools/regex',     label: '正则测试',    icon: 'Regex' },
      { path: '/devtools/snippets',  label: '代码片段',    icon: 'BookMarked' },
      { path: '/devtools/radix',     label: '进制转换',    icon: 'Binary' },
      { path: '/devtools/dummydata', label: '占位文本生成', icon: 'Database' },
      { path: '/devtools/apidoc',    label: '接口文档生成', icon: 'FileCode' },
      { path: '/devtools/obfuscator',label: '代码混淆',    icon: 'Lock' },
    ],
  },
  {
    id: 'sysinfo',
    label: '系统工具',
    items: [
      { path: '/sysinfo/process', label: '进程管理',   icon: 'Cpu' },
      { path: '/sysinfo/ports',   label: '端口查看',   icon: 'Network' },
      { path: '/sysinfo/sysinfo', label: '系统信息',   icon: 'Monitor' },
      { path: '/sysinfo/rename',  label: '批量重命名', icon: 'FilePen' },
      { path: '/sysinfo/colorpicker', label: '取色器', icon: 'Pipette' },
      { path: '/sysinfo/imagetool', label: '图片工具', icon: 'Image' },
      { path: '/sysinfo/filebatch',   label: '文件批量处理', icon: 'FolderSearch' },
      { path: '/sysinfo/diskcleaner', label: '磁盘清理',     icon: 'Trash2' },
      { path: '/sysinfo/clipboard',   label: '剪贴板管理',   icon: 'Clipboard' },
    ],
  },
  {
    id: 'daily',
    label: '日常工具',
    items: [
      { path: '/daily/calculator', label: '计算器',   icon: 'Calculator' },
      { path: '/daily/converter',  label: '单位换算', icon: 'ArrowLeftRight' },
      { path: '/daily/notes',    label: '备忘录',       icon: 'StickyNote' },
      { path: '/daily/password', label: '密码生成',     icon: 'KeyRound' },
      { path: '/daily/qrbatch',  label: '二维码批量处理', icon: 'QrCode' },
    ],
  },
  {
    id: 'network',
    label: '网络工具',
    items: [
      { path: '/network/ping',  label: 'Ping 测试', icon: 'Signal' },
      { path: '/network/scan',  label: '内网扫描',  icon: 'Radar' },
      { path: '/network/http',  label: 'HTTP 测试', icon: 'Globe' },
      { path: '/network/dns',   label: 'DNS 查询',  icon: 'Globe2' },
      { path: '/network/hosts',     label: 'Hosts 编辑',  icon: 'Server' },
      { path: '/network/whois',     label: 'WHOIS 查询', icon: 'SearchCode' },
      { path: '/network/speedtest', label: '网络测速',   icon: 'Gauge' },
    ],
  },
]

// Ctrl+K 快捷键触发搜索
const handleKeydown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    appStore.searchVisible = true
  }
}

onMounted(() => {
  document.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
/* 导航栏背景 */
.sidenav {
  background: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
}

/* 分组标题 */
.nav-group-title {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  color: var(--text-muted);
  user-select: none;
}

/* 导航项基础样式 */
.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: var(--radius-md);
  text-decoration: none;
  color: var(--text-secondary);
  transition: all var(--transition-fast);
  position: relative;
  overflow: hidden;
  border: none;
  background: transparent;
  width: 100%;
  cursor: pointer;
  font-size: 13px;
}

.nav-item.collapsed {
  justify-content: center;
  padding: 7px;
}

/* hover 效果：背景渐变 + 微上浮 */
.nav-item:hover {
  color: var(--text-primary);
  background: linear-gradient(
    135deg,
    var(--bg-hover) 0%,
    var(--bg-active) 100%
  );
  transform: translateY(-1px);
}

/* 激活状态：左侧指示条 + 背景高亮 */
.nav-item.active {
  color: var(--text-primary);
  background: var(--bg-hover);
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 15%;
  bottom: 15%;
  width: 3px;
  border-radius: 0 3px 3px 0;
  background: var(--cat-color, var(--accent));
  box-shadow: 0 0 8px var(--cat-color, var(--accent-glow));
}

.nav-item.active:hover {
  background: var(--bg-active);
}

/* 搜索框样式 */
.search-box {
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  color: inherit;
  transition: all var(--transition-fast);
}

.search-box:hover {
  background: var(--bg-active);
  border-color: var(--accent);
  box-shadow: 0 0 0 2px var(--accent-glow);
}

.search-kbd {
  display: inline-block;
  padding: 0 4px;
  font-size: 10px;
  font-family: inherit;
  background: var(--bg-active);
  border: 1px solid var(--border-color);
  border-radius: 3px;
  line-height: 1.6;
  margin-left: 2px;
}

/* 折叠模式搜索按钮 */
.search-btn-collapsed {
  background: transparent;
  border: 1px solid var(--border-color);
  color: inherit;
  transition: all var(--transition-fast);
}

.search-btn-collapsed:hover {
  background: var(--bg-hover);
  border-color: var(--accent);
}

/* 底部区域 */
.sidenav-footer {
  border-top-color: var(--border-color);
}

/* 细滚动条 */
.scrollbar-thin::-webkit-scrollbar { width: 3px; }
.scrollbar-thin::-webkit-scrollbar-track { background: transparent; }
.scrollbar-thin::-webkit-scrollbar-thumb {
  background: rgba(255, 255, 255, 0.08);
  border-radius: 2px;
}

.light .scrollbar-thin::-webkit-scrollbar-thumb {
  background: rgba(0, 0, 0, 0.08);
}
</style>
