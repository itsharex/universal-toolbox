<template>
  <!-- 文件批量处理工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <FolderSearch :size="20" class="text-primary-400" />
        文件批量处理
      </div>
      <div class="page-desc">文件搜索 · 批量复制 · 批量移动 · 批量删除</div>
    </div>

    <!-- 操作工具栏 -->
    <div class="toolbar mb-4">
      <button @click="selectSourceDir" class="btn btn-primary">
        <FolderOpen :size="14" />
        选择目录
      </button>
      <div class="flex-1" />
      <select v-model="fileFilter" class="input-field w-32">
        <option value="*">全部文件</option>
        <option value="*.jpg">JPG 图片</option>
        <option value="*.png">PNG 图片</option>
        <option value="*.pdf">PDF 文档</option>
        <option value="*.doc*">Word 文档</option>
        <option value="*.txt">文本文件</option>
        <option value="*.zip">压缩文件</option>
      </select>
      <input v-model="searchPattern" class="input-field w-40" placeholder="搜索文件名..." />
      <button @click="scanFiles" :disabled="!sourceDir" class="btn btn-secondary">
        <Search :size="14" />
        扫描
      </button>
    </div>

    <!-- 主内容区 -->
    <div class="flex-1 flex gap-4 min-h-0">
      <!-- 左侧：文件列表 -->
      <div class="flex-1 flex flex-col min-h-0">
        <div class="flex items-center justify-between mb-2">
          <div class="label mb-0">
            <span>扫描结果</span>
            <span v-if="files.length" class="text-xs opacity-50 ml-2">{{ files.length }} 个文件</span>
          </div>
          <div v-if="files.length" class="flex items-center gap-2">
            <button @click="toggleSelectAll" class="text-xs opacity-60 hover:opacity-100">
              {{ allSelected ? '取消全选' : '全选' }}
            </button>
            <span class="text-xs opacity-50">|</span>
            <span class="text-xs opacity-50">已选 {{ selectedCount }} 个</span>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-if="!files.length" class="flex-1 flex items-center justify-center opacity-30">
          <div class="text-center">
            <FolderSearch :size="40" class="mx-auto mb-2 opacity-50" />
            <div class="text-sm">选择目录并扫描文件</div>
          </div>
        </div>

        <!-- 文件列表 -->
        <div v-else class="flex-1 overflow-auto card p-0">
          <table class="w-full text-sm">
            <thead class="sticky top-0 bg-[var(--bg-card)]">
              <tr class="border-b border-white/5">
                <th class="w-8 p-2"><input type="checkbox" :checked="allSelected" @change="toggleSelectAll" /></th>
                <th class="text-left p-2 opacity-60">文件名</th>
                <th class="text-right p-2 opacity-60 w-20">大小</th>
                <th class="text-right p-2 opacity-60 w-32">修改时间</th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="(file, idx) in filteredFiles"
                :key="idx"
                class="border-b border-white/5 hover:bg-white/5 cursor-pointer"
                :class="{ 'bg-primary-500/10': file.selected }"
                @click="file.selected = !file.selected"
              >
                <td class="p-2"><input type="checkbox" v-model="file.selected" @click.stop /></td>
                <td class="p-2 truncate">
                  <div class="flex items-center gap-2">
                    <FileIcon :size="14" class="opacity-50 shrink-0" />
                    <span class="truncate">{{ file.name }}</span>
                  </div>
                  <div class="text-xs opacity-40 truncate">{{ file.path }}</div>
                </td>
                <td class="p-2 text-right opacity-60">{{ formatSize(file.size) }}</td>
                <td class="p-2 text-right opacity-60">{{ formatDate(file.modTime) }}</td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 右侧：操作面板 -->
      <div class="w-64 flex flex-col gap-3 shrink-0">
        <!-- 操作类型 -->
        <div class="card">
          <div class="font-semibold mb-3 flex items-center gap-2"><Settings2 :size="14" />操作类型</div>
          <div class="space-y-2">
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="operation" value="copy" />
              <Copy :size="14" />
              <span>复制到</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input type="radio" v-model="operation" value="move" />
              <FolderInput :size="14" />
              <span>移动到</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer text-red-400">
              <input type="radio" v-model="operation" value="delete" />
              <Trash2 :size="14" />
              <span>删除</span>
            </label>
          </div>
        </div>

        <!-- 目标目录 -->
        <div v-if="operation !== 'delete'" class="card">
          <div class="font-semibold mb-2">目标目录</div>
          <div class="flex gap-2">
            <input v-model="targetDir" class="input-field flex-1" placeholder="选择目标..." readonly />
            <button @click="selectTargetDir" class="btn btn-secondary px-2">
              <FolderOpen :size="14" />
            </button>
          </div>
        </div>

        <!-- 统计 -->
        <div class="card">
          <div class="font-semibold mb-2">统计信息</div>
          <div class="space-y-1 text-sm">
            <div class="flex justify-between">
              <span class="opacity-60">选中文件</span>
              <span>{{ selectedCount }} 个</span>
            </div>
            <div class="flex justify-between">
              <span class="opacity-60">总大小</span>
              <span>{{ formatSize(selectedSize) }}</span>
            </div>
          </div>
        </div>

        <!-- 执行按钮 -->
        <button
          @click="executeOperation"
          :disabled="!canExecute"
          :class="['btn w-full', operation === 'delete' ? 'btn-danger' : 'btn-primary']"
        >
          <Loader2 v-if="isExecuting" :size="14" class="loading-spin" />
          <Play v-else :size="14" />
          {{ isExecuting ? '执行中...' : '执行操作' }}
        </button>

        <!-- 执行日志 -->
        <div v-if="logs.length" class="card flex-1 min-h-0 flex flex-col">
          <div class="flex items-center justify-between mb-2">
            <div class="font-semibold text-sm">执行日志</div>
            <button @click="logs = []" class="text-xs opacity-50">清空</button>
          </div>
          <div class="flex-1 overflow-auto text-xs space-y-1 font-mono">
            <div v-for="(log, idx) in logs" :key="idx" :class="log.type === 'error' ? 'text-red-400' : log.type === 'success' ? 'text-green-400' : 'opacity-60'">
              {{ log.message }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { FolderSearch, FolderOpen, Search, Settings2, Copy, FolderInput, Trash2, Play, Loader2, File as FileIcon } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { ListDirectory, CopyFile, MoveFile, DeleteFile } from '@/../wailsjs/go/sysinfo/SysInfo'

const appStore = useAppStore()

// 文件接口（匹配后端 FileItem）
interface FileItem {
  name: string
  path: string
  size: number
  isDir: boolean
  modTime: number
  ext: string
  selected: boolean
}

// 状态
const sourceDir = ref('')
const targetDir = ref('')
const files = ref<FileItem[]>([])
const fileFilter = ref('*')
const searchPattern = ref('')
const operation = ref<'copy' | 'move' | 'delete'>('copy')
const isExecuting = ref(false)
const logs = ref<Array<{ type: 'info' | 'success' | 'error'; message: string }>>([])

// 过滤后的文件（排除目录 + 按扩展名筛选）
const filteredFiles = computed(() => {
  let result = files.value.filter(f => !f.isDir)

  // 按文件类型过滤
  if (fileFilter.value && fileFilter.value !== '*') {
    const ext = fileFilter.value.replace('*', '')
    result = result.filter(f => f.name.toLowerCase().endsWith(ext.toLowerCase()))
  }

  // 按文件名搜索
  if (searchPattern.value) {
    const pattern = searchPattern.value.toLowerCase()
    result = result.filter(f => f.name.toLowerCase().includes(pattern))
  }

  return result
})

// 选中数量
const selectedCount = computed(() => files.value.filter(f => f.selected && !f.isDir).length)
const selectedSize = computed(() => files.value.filter(f => f.selected && !f.isDir).reduce((sum, f) => sum + f.size, 0))
const allSelected = computed(() => {
  const selectable = files.value.filter(f => !f.isDir)
  return selectable.length > 0 && selectable.every(f => f.selected)
})

// 是否可执行
const canExecute = computed(() => {
  if (isExecuting.value) return false
  if (selectedCount.value === 0) return false
  if (operation.value !== 'delete' && !targetDir.value) return false
  return true
})

// 选择源目录
async function selectSourceDir() {
  const dir = prompt('请输入目录路径（绝对路径）：', sourceDir.value || 'C:\\Users')
  if (!dir) return
  sourceDir.value = dir
  logs.value = []
  await scanFiles()
}

// 选择目标目录
function selectTargetDir() {
  const dir = prompt('请输入目标目录路径：', targetDir.value || sourceDir.value)
  if (dir) targetDir.value = dir
}

// 扫描文件（真实调用后端 API）
async function scanFiles() {
  if (!sourceDir.value) {
    appStore.showToast('warning', '请先选择目录')
    return
  }

  logs.value = []
  logs.value.push({ type: 'info', message: `正在扫描: ${sourceDir.value}` })

  try {
    const entries = await ListDirectory(sourceDir.value)
    files.value = entries.map((e: any) => ({
      name: e.name,
      path: e.path,
      size: e.size,
      isDir: e.isDir,
      modTime: e.modTime,
      ext: e.ext,
      selected: false,
    }))

    const count = filteredFiles.value.length
    logs.value.push({ type: 'success', message: `扫描完成，共 ${count} 个文件` })
    appStore.showToast('success', `已加载 ${count} 个文件`)
  } catch (err) {
    logs.value.push({ type: 'error', message: `扫描失败: ${err}` })
    appStore.showToast('error', `扫描失败：${err}`)
  }
}

// 全选/取消
function toggleSelectAll() {
  const newState = !allSelected.value
  files.value.forEach(f => {
    if (!f.isDir) f.selected = newState
  })
}

// 执行操作（复制/移动/删除）
async function executeOperation() {
  if (!canExecute.value) return

  const selectedFiles = files.value.filter(f => f.selected && !f.isDir)
  if (!confirm(`确定要${operation.value === 'copy' ? '复制' : operation.value === 'move' ? '移动' : '删除'}这 ${selectedFiles.length} 个文件吗？`)) {
    return
  }

  isExecuting.value = true
  logs.value = []

  let success = 0
  let failed = 0

  for (const file of selectedFiles) {
    try {
      if (operation.value === 'copy') {
        // 复制：源路径 → 目标完整路径
        const destPath = `${targetDir.value}\\${file.name}`.replace(/\\\\+/g, '\\')
        await CopyFile(file.path, destPath)
        logs.value.push({ type: 'success', message: `✓ 复制: ${file.name} → ${destPath}` })
      } else if (operation.value === 'move') {
        const destPath = `${targetDir.value}\\${file.name}`.replace(/\\\\+/g, '\\')
        await MoveFile(file.path, destPath)
        logs.value.push({ type: 'success', message: `✓ 移动: ${file.name} → ${destPath}` })
      } else {
        await DeleteFile(file.path)
        logs.value.push({ type: 'success', message: `✓ 删除: ${file.name}` })
      }
      success++
    } catch (err) {
      logs.value.push({ type: 'error', message: `✗ 失败: ${file.name} - ${err}` })
      failed++
    }
  }

  isExecuting.value = false

  if (failed === 0) {
    appStore.showToast('success', `操作完成，${success} 个文件处理成功`)
  } else {
    appStore.showToast('warning', `完成，${success} 成功，${failed} 失败`)
  }

  // 删除或移动后刷新列表
  if (operation.value === 'delete' || operation.value === 'move') {
    await scanFiles()
  }
}

// 格式化大小
function formatSize(bytes: number): string {
  if (bytes < 1024) return bytes + ' B'
  if (bytes < 1024 * 1024) return (bytes / 1024).toFixed(1) + ' KB'
  return (bytes / 1024 / 1024).toFixed(1) + ' MB'
}

// 格式化日期
function formatDate(timestamp: number): string {
  if (!timestamp) return '-'
  return new Date(timestamp * 1000).toLocaleDateString()
}

// 声明批量操作 API
declare global {
  interface Window {
    __batch: {
      BatchCopy: (pair: string) => Promise<void>
      BatchMove: (pair: string) => Promise<void>
      BatchDelete: (path: string) => Promise<void>
    }
  }
}
</script>
</script>

<style scoped>
/* 表格样式 */
table {
  border-collapse: collapse;
}

input[type="checkbox"] {
  width: 14px;
  height: 14px;
  cursor: pointer;
}

input[type="radio"] {
  cursor: pointer;
}
</style>
