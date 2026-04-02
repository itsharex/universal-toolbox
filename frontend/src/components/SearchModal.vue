<template>
  <!-- 全局搜索弹窗 -->
  <Teleport to="body">
    <Transition name="search-fade">
      <div v-if="appStore.searchVisible" class="search-overlay" @click.self="close">
        <div class="search-modal" @keydown.esc="close">
          <!-- 搜索输入框 -->
          <div class="search-input-wrapper">
            <Search :size="18" class="search-icon" />
            <input
              ref="inputRef"
              v-model="query"
              type="text"
              placeholder="搜索工具... (按 ESC 关闭)"
              class="search-input"
              @keydown.down.prevent="moveDown"
              @keydown.up.prevent="moveUp"
              @keydown.enter="selectCurrent"
              @keydown.esc="close"
            />
            <kbd class="search-kbd">Ctrl+K</kbd>
          </div>

          <!-- 搜索结果 -->
          <div class="search-results" v-if="query || recentTools.length > 0">
            <!-- 最近使用 -->
            <div v-if="!query && recentTools.length > 0" class="search-section">
              <div class="search-section-title">
                <Clock :size="12" />
                最近使用
              </div>
              <div
                v-for="(tool, index) in recentTools"
                :key="tool.path"
                class="search-item"
                :class="{ active: index === selectedIndex }"
                @click="navigateTo(tool)"
                @mouseenter="selectedIndex = index"
              >
                <component :is="getIcon(tool.icon)" :size="14" class="item-icon" />
                <div class="item-content">
                  <span class="item-label">{{ tool.label }}</span>
                  <span class="item-category">{{ tool.category }}</span>
                </div>
              </div>
            </div>

            <!-- 搜索结果 -->
            <div v-if="filteredTools.length > 0" class="search-section">
              <div v-if="query" class="search-section-title">
                <Sparkles :size="12" />
                搜索结果 ({{ filteredTools.length }})
              </div>
              <div
                v-for="(tool, index) in filteredTools"
                :key="tool.path"
                class="search-item"
                :class="{ active: getGlobalIndex(index) === selectedIndex }"
                @click="navigateTo(tool)"
                @mouseenter="selectedIndex = getGlobalIndex(index)"
              >
                <component :is="getIcon(tool.icon)" :size="14" class="item-icon" />
                <div class="item-content">
                  <span class="item-label">{{ tool.label }}</span>
                  <span class="item-category">{{ tool.category }}</span>
                </div>
                <span v-if="tool.hot" class="item-hot">🔥</span>
              </div>
            </div>

            <!-- 无结果 -->
            <div v-else-if="query" class="search-empty">
              <SearchX :size="32" class="empty-icon" />
              <p>未找到匹配的工具</p>
              <p class="empty-hint">尝试其他关键词</p>
            </div>
          </div>

          <!-- 空状态提示 -->
          <div v-else class="search-hints">
            <div class="hint-title">快捷键提示</div>
            <div class="hint-list">
              <div class="hint-item">
                <kbd>Ctrl</kbd> + <kbd>K</kbd>
                <span>打开搜索</span>
              </div>
              <div class="hint-item">
                <kbd>↑</kbd> <kbd>↓</kbd>
                <span>选择</span>
              </div>
              <div class="hint-item">
                <kbd>Enter</kbd>
                <span>跳转</span>
              </div>
              <div class="hint-item">
                <kbd>Esc</kbd>
                <span>关闭</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, computed, watch, nextTick, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  Search, SearchX, Clock, Sparkles,
  Braces, Code2, Binary, Link, Hash, FileText, QrCode, GitCompare,
  Fingerprint, Regex, BookMarked, Shield, Palette, Keyboard,
  Cpu, Network, Monitor, FilePen, Pipette, Image, FolderSearch,
  Calculator, ArrowLeftRight, StickyNote,
  Signal, Radar, Globe, Globe2, Server,
} from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

// 工具索引数据（包含搜索关键词）
interface ToolItem {
  path: string
  label: string
  icon: string
  category: string
  keywords: string[]
  hot?: boolean
}

const router = useRouter()
const appStore = useAppStore()

const inputRef = ref<HTMLInputElement | null>(null)
const query = ref('')
const selectedIndex = ref(0)
const recentTools = ref<ToolItem[]>([])

// 图标映射
const iconMap: Record<string, any> = {
  Braces, Code2, Binary, Link, Hash, FileText, QrCode, GitCompare,
  Fingerprint, Clock, Regex, BookMarked, Shield, Palette, Keyboard,
  Cpu, Network, Monitor, FilePen, Pipette, Image, FolderSearch,
  Calculator, ArrowLeftRight, StickyNote,
  Signal, Radar, Globe, Globe2, Server,
}
const getIcon = (name: string) => iconMap[name] || FileText

// 所有工具索引（包含搜索别名）
const allTools: ToolItem[] = [
  // 开发工具
  { path: '/devtools/json', label: 'JSON 工具', icon: 'Braces', category: '开发工具', keywords: ['json', '格式化', '验证', '压缩', '美化'] },
  { path: '/devtools/jsondiff', label: 'JSON 对比', icon: 'GitCompare', category: '开发工具', keywords: ['json', 'diff', '对比', '差异', '比较'], hot: true },
  { path: '/devtools/xml', label: 'XML 工具', icon: 'Code2', category: '开发工具', keywords: ['xml', '格式化', '验证'] },
  { path: '/devtools/base64', label: 'Base64 编解码', icon: 'Binary', category: '开发工具', keywords: ['base64', '编码', '解码'], hot: true },
  { path: '/devtools/url', label: 'URL 编解码', icon: 'Link', category: '开发工具', keywords: ['url', '编码', '解码', '百分号'], hot: true },
  { path: '/devtools/hash', label: '哈希计算', icon: 'Hash', category: '开发工具', keywords: ['hash', 'md5', 'sha', '加密', '散列'] },
  { path: '/devtools/crypto', label: '加密解密', icon: 'Shield', category: '开发工具', keywords: ['aes', 'rsa', '加密', '解密', 'crypto'], hot: true },
  { path: '/devtools/text', label: '文本工具', icon: 'FileText', category: '开发工具', keywords: ['文本', '查找', '替换', '正则', '大小写'] },
  { path: '/devtools/qrcode', label: '二维码工具', icon: 'QrCode', category: '开发工具', keywords: ['qr', '二维码', '生成', '解码', 'scan'] },
  { path: '/devtools/uuid', label: 'UUID 生成', icon: 'Fingerprint', category: '开发工具', keywords: ['uuid', 'guid', '唯一', '生成'] },
  { path: '/devtools/timestamp', label: '时间戳转换', icon: 'Clock', category: '开发工具', keywords: ['时间戳', 'timestamp', 'unix', '日期'], hot: true },
  { path: '/devtools/regex', label: '正则测试', icon: 'Regex', category: '开发工具', keywords: ['正则', 'regex', '匹配', 'test'], hot: true },
  { path: '/devtools/snippets', label: '代码片段', icon: 'BookMarked', category: '开发工具', keywords: ['代码', '片段', 'snippet', '代码库'] },

  // 系统工具
  { path: '/sysinfo/process', label: '进程管理', icon: 'Cpu', category: '系统工具', keywords: ['进程', 'process', '任务管理器', 'kill'] },
  { path: '/sysinfo/ports', label: '端口查看', icon: 'Network', category: '系统工具', keywords: ['端口', 'port', '网络', '占用'] },
  { path: '/sysinfo/sysinfo', label: '系统信息', icon: 'Monitor', category: '系统工具', keywords: ['系统', '系统信息', 'info', '硬件', '配置'] },
  { path: '/sysinfo/rename', label: '批量重命名', icon: 'FilePen', category: '系统工具', keywords: ['批量', '重命名', 'rename', '批量改名', '文件'], hot: true },
  { path: '/sysinfo/colorpicker', label: '取色器', icon: 'Pipette', category: '系统工具', keywords: ['取色', '颜色', 'color', '吸管', '屏幕取色'], hot: true },
  { path: '/sysinfo/imagetool', label: '图片工具', icon: 'Image', category: '系统工具', keywords: ['图片', '压缩', 'resize', '格式', 'webp', 'png'] },
  { path: '/sysinfo/filebatch', label: '文件批量处理', icon: 'FolderSearch', category: '系统工具', keywords: ['文件', '批量', '复制', '移动', '删除', 'batch'] },

  // 日常工具
  { path: '/daily/calculator', label: '计算器', icon: 'Calculator', category: '日常工具', keywords: ['计算', '计算器', 'calc', '数学'] },
  { path: '/daily/converter', label: '单位换算', icon: 'ArrowLeftRight', category: '日常工具', keywords: ['单位', '换算', 'convert', '长度', '重量', '温度'] },
  { path: '/daily/notes', label: '备忘录', icon: 'StickyNote', category: '日常工具', keywords: ['备忘录', '笔记', 'note', '便签'] },

  // 网络工具
  { path: '/network/ping', label: 'Ping 测试', icon: 'Signal', category: '网络工具', keywords: ['ping', '延迟', '网络', '测试'], hot: true },
  { path: '/network/scan', label: '内网扫描', icon: 'Radar', category: '网络工具', keywords: ['扫描', '内网', 'ip', '局域网', 'scan'] },
  { path: '/network/http', label: 'HTTP 测试', icon: 'Globe', category: '网络工具', keywords: ['http', 'api', '请求', 'rest', 'curl', '测试'], hot: true },
  { path: '/network/dns', label: 'DNS 查询', icon: 'Globe2', category: '网络工具', keywords: ['dns', '域名', '解析', 'lookup'] },
  { path: '/network/hosts', label: 'Hosts 编辑', icon: 'Server', category: '网络工具', keywords: ['hosts', '域名', '绑定', '解析'] },
]

// 根据查询过滤工具
const filteredTools = computed(() => {
  if (!query.value.trim()) return []

  const q = query.value.toLowerCase().trim()
  return allTools.filter(tool => {
    // 匹配标签
    if (tool.label.toLowerCase().includes(q)) return true
    // 匹配关键词
    if (tool.keywords.some(k => k.toLowerCase().includes(q))) return true
    // 匹配分类
    if (tool.category.toLowerCase().includes(q)) return true
    return false
  }).sort((a, b) => {
    // 热门工具优先
    if (a.hot && !b.hot) return -1
    if (!a.hot && b.hot) return 1
    // 精确匹配优先
    if (a.label.toLowerCase().includes(q) && !b.label.toLowerCase().includes(q)) return -1
    return 0
  })
})

// 获取全局索引
const getGlobalIndex = (localIndex: number) => {
  return localIndex
}

// 监听搜索显示状态，自动聚焦输入框
watch(() => appStore.searchVisible, async (visible) => {
  if (visible) {
    query.value = ''
    selectedIndex.value = 0
    // 从 localStorage 加载最近使用
    try {
      const recent = localStorage.getItem('xtool_recent_tools')
      if (recent) {
        const paths = JSON.parse(recent) as string[]
        recentTools.value = paths
          .map(p => allTools.find(t => t.path === p))
          .filter((t): t is ToolItem => t !== undefined)
          .slice(0, 5)
      }
    } catch {}
    await nextTick()
    inputRef.value?.focus()
  }
})

// 关闭搜索
const close = () => {
  appStore.searchVisible = false
  query.value = ''
}

// 移动选择
const moveDown = () => {
  const max = filteredTools.value.length - 1
  if (selectedIndex.value < max) selectedIndex.value++
}

const moveUp = () => {
  if (selectedIndex.value > 0) selectedIndex.value--
}

// 选择当前项并导航
const selectCurrent = () => {
  const tools = filteredTools.value
  if (tools.length > 0 && selectedIndex.value < tools.length) {
    navigateTo(tools[selectedIndex.value])
  }
}

// 导航到工具页
const navigateTo = (tool: ToolItem) => {
  // 保存到最近使用
  try {
    let recent = JSON.parse(localStorage.getItem('xtool_recent_tools') || '[]') as string[]
    recent = [tool.path, ...recent.filter(p => p !== tool.path)].slice(0, 10)
    localStorage.setItem('xtool_recent_tools', JSON.stringify(recent))
  } catch {}

  router.push(tool.path)
  close()
}

// 键盘事件处理（Ctrl+K 全局）
const handleKeydown = (e: KeyboardEvent) => {
  if ((e.ctrlKey || e.metaKey) && e.key === 'k') {
    e.preventDefault()
    appStore.searchVisible = !appStore.searchVisible
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})
</script>

<style scoped>
/* 遮罩层 */
.search-overlay {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: flex-start;
  justify-content: center;
  padding-top: 15vh;
  z-index: 9999;
}

/* 搜索弹窗 */
.search-modal {
  width: 100%;
  max-width: 560px;
  background: var(--bg-secondary, #1e1e2e);
  border: 1px solid var(--border-color, rgba(255,255,255,0.1));
  border-radius: 12px;
  box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  overflow: hidden;
}

.light .search-modal {
  background: #ffffff;
  border-color: rgba(0, 0, 0, 0.1);
}

/* 搜索输入框 */
.search-input-wrapper {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  border-bottom: 1px solid var(--border-color, rgba(255,255,255,0.08));
}

.search-icon {
  color: var(--text-secondary, #888);
  flex-shrink: 0;
}

.search-input {
  flex: 1;
  background: transparent;
  border: none;
  outline: none;
  font-size: 16px;
  color: var(--text-primary, #e2e8f0);
}

.search-input::placeholder {
  color: var(--text-secondary, #888);
}

.search-kbd {
  padding: 4px 8px;
  font-size: 11px;
  font-family: monospace;
  background: rgba(255,255,255,0.1);
  border-radius: 4px;
  color: var(--text-secondary, #888);
}

.light .search-kbd {
  background: rgba(0,0,0,0.06);
}

/* 搜索结果 */
.search-results {
  max-height: 400px;
  overflow-y: auto;
}

.search-section {
  padding: 8px;
}

.search-section-title {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 6px 8px;
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--text-secondary, #888);
}

.search-item {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 10px 12px;
  border-radius: 8px;
  cursor: pointer;
  transition: background 0.1s ease;
}

.search-item:hover,
.search-item.active {
  background: rgba(99, 102, 241, 0.15);
}

.light .search-item:hover,
.light .search-item.active {
  background: rgba(99, 102, 241, 0.1);
}

.item-icon {
  color: var(--text-secondary, #888);
  flex-shrink: 0;
}

.item-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.item-label {
  font-size: 14px;
  color: var(--text-primary, #e2e8f0);
}

.item-category {
  font-size: 11px;
  color: var(--text-secondary, #888);
}

.item-hot {
  font-size: 12px;
}

/* 无结果 */
.search-empty {
  padding: 40px 20px;
  text-align: center;
  color: var(--text-secondary, #888);
}

.empty-icon {
  margin: 0 auto 12px;
  opacity: 0.5;
}

.search-empty p {
  margin: 0;
  font-size: 14px;
}

.empty-hint {
  font-size: 12px !important;
  margin-top: 4px !important;
  opacity: 0.7;
}

/* 快捷键提示 */
.search-hints {
  padding: 20px;
}

.hint-title {
  font-size: 12px;
  font-weight: 600;
  color: var(--text-secondary, #888);
  margin-bottom: 12px;
}

.hint-list {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 8px;
}

.hint-item {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: var(--text-secondary, #888);
}

.hint-item kbd {
  padding: 3px 6px;
  font-size: 11px;
  font-family: monospace;
  background: rgba(255,255,255,0.1);
  border-radius: 4px;
}

.light .hint-item kbd {
  background: rgba(0,0,0,0.06);
}

/* 过渡动画 */
.search-fade-enter-active,
.search-fade-leave-active {
  transition: opacity 0.15s ease;
}

.search-fade-enter-active .search-modal,
.search-fade-leave-active .search-modal {
  transition: transform 0.15s ease, opacity 0.15s ease;
}

.search-fade-enter-from,
.search-fade-leave-to {
  opacity: 0;
}

.search-fade-enter-from .search-modal,
.search-fade-leave-to .search-modal {
  transform: scale(0.95) translateY(-10px);
  opacity: 0;
}
</style>
