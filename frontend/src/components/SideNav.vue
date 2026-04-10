<template>
  <!-- 左侧导航栏组件 -->
  <nav
    :class="[
      'sidenav flex flex-col shrink-0 transition-all duration-300 ease-in-out',
      collapsed ? 'w-14' : 'w-56'
    ]"
  >
    <!-- 搜索框 -->
    <div v-if="!collapsed" class="px-3 py-2">
      <div class="search-box flex items-center gap-2 px-3 py-1.5 rounded-lg cursor-pointer"
           @click="appStore.searchVisible = true">
        <Search :size="13" class="opacity-50" />
        <span class="text-xs opacity-40">搜索工具... Ctrl+K</span>
      </div>
    </div>

    <!-- 导航分组 -->
    <div class="flex-1 overflow-y-auto overflow-x-hidden px-2 py-1 space-y-0.5 scrollbar-thin">
      <template v-for="group in navGroups" :key="group.id">
        <!-- 分组标题 -->
        <div v-if="!collapsed"
             class="nav-group-title px-2 py-1 mt-3 first:mt-1">
          {{ group.label }}
        </div>
        <div v-else class="h-px mx-2 my-2 opacity-10 bg-current" />

        <!-- 导航项 -->
        <RouterLink
          v-for="item in group.items"
          :key="item.path"
          :to="item.path"
          class="nav-item"
          :class="{ 'active': isActive(item.path), 'collapsed': collapsed }"
          v-tooltip="collapsed ? item.label : ''"
        >
          <component :is="getIcon(item.icon)" :size="16" class="shrink-0" />
          <span v-if="!collapsed" class="truncate text-sm">{{ item.label }}</span>
        </RouterLink>
      </template>
    </div>

    <!-- 底部操作区 -->
    <div class="px-2 py-2 border-t border-white/5 space-y-0.5">
      <!-- 设置 -->
      <RouterLink to="/settings" class="nav-item" :class="{ 'active': isActive('/settings'), 'collapsed': collapsed }">
        <Settings :size="16" class="shrink-0" />
        <span v-if="!collapsed" class="text-sm">设置</span>
      </RouterLink>
      <!-- 主题编辑器 -->
      <RouterLink to="/settings/theme" class="nav-item" :class="{ 'active': isActive('/settings/theme'), 'collapsed': collapsed }">
        <Palette :size="16" class="shrink-0" />
        <span v-if="!collapsed" class="text-sm">主题编辑</span>
      </RouterLink>
      <!-- 快捷键管理 -->
      <RouterLink to="/settings/shortcuts" class="nav-item" :class="{ 'active': isActive('/settings/shortcuts'), 'collapsed': collapsed }">
        <Keyboard :size="16" class="shrink-0" />
        <span v-if="!collapsed" class="text-sm">快捷键</span>
      </RouterLink>

      <!-- 折叠按钮 -->
      <button @click="appStore.toggleSidebar()"
              class="nav-item w-full">
        <PanelLeftClose v-if="!collapsed" :size="16" class="shrink-0" />
        <PanelLeftOpen v-else :size="16" class="shrink-0" />
        <span v-if="!collapsed" class="text-sm">收起侧栏</span>
      </button>
    </div>
  </nav>
</template>

<script setup lang="ts">
import { computed } from 'vue'
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
</script>

<style scoped>
/* 导航栏背景 */
.sidenav {
  background: rgba(255, 255, 255, 0.02);
  border-right: 1px solid rgba(255, 255, 255, 0.06);
}

.light .sidenav {
  background: rgba(0, 0, 0, 0.02);
  border-right: 1px solid rgba(0, 0, 0, 0.08);
}

/* 分组标题 */
.nav-group-title {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.08em;
  opacity: 0.4;
  color: var(--nav-text);
}

/* 导航项基础样式 */
.nav-item {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 10px;
  border-radius: 8px;
  text-decoration: none;
  color: inherit;
  opacity: 0.65;
  transition: all 0.15s ease;
  position: relative;
  overflow: hidden;
}

.nav-item.collapsed {
  justify-content: center;
  padding: 7px;
}

.nav-item:hover {
  opacity: 1;
  background: rgba(255, 255, 255, 0.06);
}

.light .nav-item:hover {
  background: rgba(0, 0, 0, 0.05);
}

/* 激活状态 */
.nav-item.active {
  opacity: 1;
  background: rgba(99, 102, 241, 0.15);
  color: #818cf8;
}

.light .nav-item.active {
  background: rgba(79, 70, 229, 0.1);
  color: #4f46e5;
}

.nav-item.active::before {
  content: '';
  position: absolute;
  left: 0;
  top: 20%;
  bottom: 20%;
  width: 3px;
  border-radius: 0 3px 3px 0;
  background: #6366f1;
}

/* 搜索框样式 */
.search-box {
  background: rgba(255, 255, 255, 0.04);
  border: 1px solid rgba(255, 255, 255, 0.07);
  color: inherit;
  transition: all 0.15s ease;
}

.search-box:hover {
  background: rgba(255, 255, 255, 0.07);
  border-color: rgba(255, 255, 255, 0.12);
}

.light .search-box {
  background: rgba(0, 0, 0, 0.04);
  border-color: rgba(0, 0, 0, 0.08);
}

/* 细滚动条 */
.scrollbar-thin::-webkit-scrollbar { width: 3px; }
.scrollbar-thin::-webkit-scrollbar-track { background: transparent; }
.scrollbar-thin::-webkit-scrollbar-thumb { background: rgba(255,255,255,0.1); border-radius: 2px; }
</style>
