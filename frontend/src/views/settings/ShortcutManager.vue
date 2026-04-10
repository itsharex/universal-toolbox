<template>
  <ToolPage title="快捷键管理" description="自定义工具快捷键">
<!-- 快捷键管理页面 -->
  
    <!-- 搜索 -->
    <div class="toolbar mb-4">
      <input v-model="searchQuery" class="input-field flex-1" placeholder="搜索快捷键..." />
      <button @click="resetAll" class="btn btn-secondary">
        <RotateCcw :size="14" />
        重置全部
      </button>
    </div>

    <!-- 快捷键列表 -->
    <div class="flex-1 overflow-auto">
      <div class="space-y-2">
        <div
          v-for="group in filteredGroups"
          :key="group.id"
          class="card"
        >
          <div class="font-semibold mb-3 flex items-center gap-2">
            <component :is="group.icon" :size="16" class="opacity-60" />
            {{ group.name }}
          </div>
          <div class="space-y-1">
            <div
              v-for="shortcut in group.shortcuts"
              :key="shortcut.id"
              class="flex items-center justify-between p-2 rounded-lg hover:bg-white/5"
            >
              <div class="flex-1">
                <div class="text-sm">{{ shortcut.name }}</div>
                <div class="text-xs opacity-50">{{ shortcut.description }}</div>
              </div>
              <div class="flex items-center gap-2">
                <!-- 快捷键显示/编辑 -->
                <div
                  @click="editShortcut(shortcut)"
                  :class="[
                    'shortcut-keys flex items-center gap-1 cursor-pointer transition-all',
                    editingId === shortcut.id ? 'ring-2 ring-primary-500' : '',
                    hasConflict(shortcut) ? 'border-red-500/50' : ''
                  ]"
                >
                  <template v-if="editingId === shortcut.id">
                    <span class="text-xs opacity-50 animate-pulse">按下快捷键...</span>
                  <template v-else-if="shortcut.keys">
                    <kbd v-for="key in shortcut.keys" :key="key" class="kbd">{{ key }}</kbd>
                  </template>
                  <template v-else>
                    <span class="text-xs opacity-30">未设置</span>
                  </template>
                </div>
                <!-- 冲突提示 -->
                <span v-if="hasConflict(shortcut)" class="text-xs text-red-400">冲突</span>
                <!-- 启用开关 -->
                <label class="toggle-mini">
                  <input type="checkbox" v-model="shortcut.enabled" @change="saveShortcut(shortcut)" />
                  <span class="toggle-track-mini" />
                </label>
                <!-- 清除按钮 -->
                <button
                  v-if="shortcut.keys"
                  @click="clearShortcut(shortcut)"
                  class="p-1 opacity-30 hover:opacity-100 hover:text-red-400"
                >
                  <X :size="12" />
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 编辑提示 -->
    <div v-if="editingId" class="fixed inset-0 bg-black/50 flex items-center justify-center z-50" @click="cancelEdit">
      <div class="card p-6 text-center" @click.stop>
        <Keyboard :size="32" class="mx-auto mb-3 text-primary-400" />
        <div class="font-semibold mb-2">按下新的快捷键</div>
        <div class="text-sm opacity-60 mb-4">或按 ESC 取消</div>
        <div class="flex gap-2 justify-center">
          <button @click="cancelEdit" class="btn btn-secondary">取消</button>
        </div>
      </div>
    </div>
  </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Keyboard, RotateCcw, X, Settings, Code, Cpu, Globe, Search } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 快捷键接口
interface Shortcut {
  id: string
  name: string
  description: string
  keys: string[]
  enabled: boolean
  action: string
}

interface ShortcutGroup {
  id: string
  name: string
  icon: any
  shortcuts: Shortcut[]
}

// 图标映射（反序列化后恢复图标引用）
const iconMap: Record<string, any> = {
  Settings,
  Code,
  Cpu,
  Globe,
  Search
}

// 默认快捷键配置
function getDefaultShortcuts(): ShortcutGroup[] {
  return [
    {
      id: 'general',
      name: '常规',
      icon: Settings,
      shortcuts: [
        { id: 'search', name: '搜索工具', description: '打开功能搜索', keys: ['Ctrl', 'K'], enabled: true, action: 'search' },
        { id: 'theme', name: '切换主题', description: '在深色/浅色主题间切换', keys: ['Ctrl', 'Shift', 'T'], enabled: true, action: 'toggleTheme' },
        { id: 'ontop', name: '窗口置顶', description: '切换窗口置顶状态', keys: ['Ctrl', 'Shift', 'O'], enabled: true, action: 'toggleOnTop' },
        { id: 'minimize', name: '最小化', description: '最小化到托盘', keys: ['Ctrl', 'M'], enabled: true, action: 'minimize' },
      ]
    },
    {
      id: 'devtools',
      name: '开发工具',
      icon: Code,
      shortcuts: [
        { id: 'json', name: 'JSON 工具', description: '打开 JSON 格式化工具', keys: ['Ctrl', 'Shift', 'J'], enabled: true, action: 'openJson' },
        { id: 'base64', name: 'Base64 工具', description: '打开 Base64 编解码', keys: ['Ctrl', 'Shift', 'B'], enabled: true, action: 'openBase64' },
        { id: 'hash', name: '哈希计算', description: '打开哈希计算工具', keys: ['Ctrl', 'Shift', 'H'], enabled: true, action: 'openHash' },
        { id: 'regex', name: '正则测试', description: '打开正则测试工具', keys: ['Ctrl', 'Shift', 'R'], enabled: true, action: 'openRegex' },
      ]
    },
    {
      id: 'system',
      name: '系统工具',
      icon: Cpu,
      shortcuts: [
        { id: 'process', name: '进程管理', description: '打开进程管理器', keys: ['Ctrl', 'Shift', 'P'], enabled: false, action: 'openProcess' },
        { id: 'port', name: '端口查看', description: '打开端口查看工具', keys: [], enabled: false, action: 'openPort' },
      ]
    },
    {
      id: 'network',
      name: '网络工具',
      icon: Globe,
      shortcuts: [
        { id: 'ping', name: 'Ping 测试', description: '打开 Ping 工具', keys: [], enabled: false, action: 'openPing' },
        { id: 'http', name: 'HTTP 测试', description: '打开 HTTP 测试工具', keys: [], enabled: false, action: 'openHttp' },
      ]
    }
  ]
}

// 状态
const searchQuery = ref('')
const editingId = ref<string | null>(null)
const shortcuts = ref<ShortcutGroup[]>(getDefaultShortcuts())

// 过滤后的分组
const filteredGroups = computed(() => {
  if (!searchQuery.value) return shortcuts.value

  const query = searchQuery.value.toLowerCase()
  return shortcuts.value
    .map(group => ({
      ...group,
      shortcuts: group.shortcuts.filter(s =>
        s.name.toLowerCase().includes(query) ||
        s.description.toLowerCase().includes(query)
      )
    }))
    .filter(group => group.shortcuts.length > 0)
})

// 检测快捷键冲突
function hasConflict(shortcut: Shortcut): boolean {
  if (!shortcut.keys || shortcut.keys.length === 0) return false
  const keyStr = shortcut.keys.join('+')
  for (const group of shortcuts.value) {
    for (const s of group.shortcuts) {
      if (s.id !== shortcut.id && s.keys.length > 0 && s.keys.join('+') === keyStr) {
        return true
      }
    }
  }
  return false
}

// 编辑快捷键
function editShortcut(shortcut: Shortcut) {
  editingId.value = shortcut.id
}

// 取消编辑
function cancelEdit() {
  editingId.value = null
}

// 清除快捷键
function clearShortcut(shortcut: Shortcut) {
  shortcut.keys = []
  saveShortcut(shortcut)
}

// 保存快捷键
function saveShortcut(shortcut: Shortcut) {
  // 检测冲突
  if (shortcut.keys.length > 0 && hasConflict(shortcut)) {
    appStore.showToast('warning', `快捷键 ${shortcut.keys.join('+')} 与其他快捷键冲突`)
  }
  localStorage.setItem('shortcuts', JSON.stringify(shortcuts.value))
  appStore.showToast('success', '快捷键已保存')
}

// 重置全部
function resetAll() {
  shortcuts.value = getDefaultShortcuts()
  localStorage.removeItem('shortcuts')
  appStore.showToast('success', '已重置为默认快捷键')
}

// 键盘事件处理
function handleKeyDown(e: KeyboardEvent) {
  if (!editingId.value) return

  e.preventDefault()
  e.stopPropagation()

  // ESC 取消
  if (e.key === 'Escape') {
    cancelEdit()
    return
  }

  // 忽略单独的修饰键
  if (['Control', 'Alt', 'Shift', 'Meta'].includes(e.key)) return

  // 构建快捷键
  const keys: string[] = []
  if (e.ctrlKey) keys.push('Ctrl')
  if (e.altKey) keys.push('Alt')
  if (e.shiftKey) keys.push('Shift')
  if (e.metaKey) keys.push('Meta')

  // 添加主键
  const mainKey = e.key.toUpperCase()
  if (!keys.includes(mainKey)) {
    keys.push(mainKey)
  }

  // 检测冲突
  const keyStr = keys.join('+')
  let conflictName = ''
  for (const group of shortcuts.value) {
    for (const s of group.shortcuts) {
      if (s.id !== editingId.value && s.keys.length > 0 && s.keys.join('+') === keyStr) {
        conflictName = s.name
        break
      }
    }
    if (conflictName) break
  }

  // 更新快捷键
  for (const group of shortcuts.value) {
    const shortcut = group.shortcuts.find(s => s.id === editingId.value)
    if (shortcut) {
      shortcut.keys = keys
      saveShortcut(shortcut)
      if (conflictName) {
        appStore.showToast('warning', `快捷键 ${keyStr} 与"${conflictName}"冲突`)
      }
      break
    }
  }

  editingId.value = null
}

// 恢复图标引用（反序列化后 icon 丢失的修复）
function restoreIcons(data: any[]): ShortcutGroup[] {
  return data.map((group: any) => ({
    ...group,
    icon: iconMap[group.id] || Settings
  }))
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)

  // 加载保存的快捷键
  const saved = localStorage.getItem('shortcuts')
  if (saved) {
    try {
      const parsed = JSON.parse(saved)
      shortcuts.value = restoreIcons(parsed)
    } catch {
      // 解析失败使用默认值
      shortcuts.value = getDefaultShortcuts()
    }
  }
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>

<style scoped>
/* 快捷键显示 */
.shortcut-keys {
  background: var(--bg-hover);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 4px 8px;
  min-width: 80px;
  justify-content: center;
}

.shortcut-keys:hover {
  border-color: var(--accent);
}

/* 键盘按键样式 */
.kbd {
  background: var(--bg-hover);
  border-radius: 4px;
  padding: 2px 6px;
  font-size: 11px;
  font-family: monospace;
  border: 1px solid rgba(255, 255, 255, 0.15);
  box-shadow: 0 2px 0 rgba(0, 0, 0, 0.3);
}

/* 小开关 */
.toggle-mini {
  position: relative;
  cursor: pointer;
}

.toggle-mini input {
  display: none;
}

.toggle-track-mini {
  display: block;
  width: 28px;
  height: 16px;
  background: var(--bg-hover);
  border-radius: 8px;
  transition: background 0.2s;
  position: relative;
}

.toggle-track-mini::after {
  content: '';
  position: absolute;
  top: 2px;
  left: 2px;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  background: #fff;
  transition: transform 0.2s;
}

.toggle-mini input:checked + .toggle-track-mini {
  background: var(--accent);
}

.toggle-mini input:checked + .toggle-track-mini::after {
  transform: translateX(12px);
}
</style>
