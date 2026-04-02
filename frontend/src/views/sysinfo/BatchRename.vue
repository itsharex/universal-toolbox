<template>
  <div class="page-container">
    <div class="page-title">
      <FilePen :size="20" class="text-primary-400" />
      批量重命名
    </div>
    <div class="page-desc">批量重命名文件，支持正则匹配和序号替换</div>

    <!-- 目录选择区 -->
    <div class="card mb-4">
      <div class="flex items-center gap-3 mb-3">
        <FolderOpen :size="16" class="opacity-60" />
        <span class="text-sm opacity-70">选择目录</span>
      </div>
      <div class="flex gap-2">
        <input
          v-model="currentDir"
          class="input-field flex-1"
          placeholder="点击「浏览」选择文件目录..."
          readonly
        />
        <button @click="browseDir" class="btn btn-secondary flex items-center gap-1">
          <FolderSearch :size="14" />
          浏览
        </button>
      </div>
    </div>

    <!-- 文件列表区 -->
    <div v-if="files.length" class="card mb-4">
      <div class="flex items-center justify-between mb-3">
        <span class="text-sm opacity-70">共 {{ files.length }} 个文件</span>
        <button @click="files = []; results = []" class="text-xs opacity-50 hover:opacity-100">清空</button>
      </div>
      <!-- 重命名规则 -->
      <div class="rule-grid mb-3">
        <div class="flex gap-2">
          <input
            v-model="pattern"
            class="input-field flex-1"
            placeholder="查找模式（如：IMG_(\d+)，支持正则）"
          />
          <input
            v-model="replacement"
            class="input-field flex-1"
            placeholder="替换为（如：photo_$1.jpg，使用 $1 $2 引用捕获组）"
          />
          <input
            v-model.number="startNum"
            class="input-field w-20"
            type="number"
            min="0"
            placeholder="起始"
            title="起始序号"
          />
        </div>
        <div class="flex gap-4 mt-2">
          <label class="flex items-center gap-2 text-xs opacity-70 cursor-pointer">
            <input type="checkbox" v-model="useRegex" class="checkbox" />
            使用正则
          </label>
          <label class="flex items-center gap-2 text-xs opacity-70 cursor-pointer">
            <input type="checkbox" v-model="caseSensitive" class="checkbox" />
            区分大小写
          </label>
          <label class="flex items-center gap-2 text-xs opacity-70 cursor-pointer">
            <input type="checkbox" v-model="previewOnly" class="checkbox" />
            仅预览
          </label>
        </div>
      </div>
      <!-- 操作按钮 -->
      <div class="flex gap-2 mb-3">
        <button @click="preview" class="btn btn-secondary">
          <Eye :size="14" />
          预览
        </button>
        <button @click="rename" class="btn btn-primary" :disabled="isRenaming">
          <Loader2 v-if="isRenaming" :size="14" class="animate-spin" />
          <FilePen v-else :size="14" />
          {{ isRenaming ? '重命名中...' : '执行重命名' }}
        </button>
      </div>
      <!-- 文件列表 -->
      <div class="file-list max-h-64 overflow-auto">
        <div
          v-for="(f, i) in files"
          :key="i"
          class="file-row"
          :class="{ changed: f.renamed !== f.original, error: !!f.error }"
        >
          <div class="flex items-center gap-2 flex-1 min-w-0">
            <FileIcon :size="14" class="opacity-40 shrink-0" />
            <span class="original truncate">{{ f.original }}</span>
            <span v-if="f.renamed !== f.original" class="text-primary-400 shrink-0">→</span>
            <span v-if="f.renamed !== f.original" class="renamed truncate text-primary-400">{{ f.renamed }}</span>
          </div>
          <span v-if="f.error" class="error-tag">{{ f.error }}</span>
        </div>
      </div>
    </div>

    <!-- 空状态 -->
    <div v-else-if="currentDir" class="card text-center py-12 opacity-50">
      <FolderOpen :size="40" class="mx-auto mb-3 opacity-30" />
      <p class="text-sm">该目录下没有文件</p>
    </div>

    <!-- 结果统计 -->
    <div v-if="results.length" class="card">
      <div class="result-header mb-3">
        <CheckCircle v-if="allSuccess" :size="20" class="text-green-400" />
        <AlertCircle v-else :size="20" class="text-yellow-400" />
        <span class="text-sm">
          {{ successCount }} / {{ results.length }} 个文件重命名成功
        </span>
      </div>
      <div class="space-y-1 max-h-48 overflow-auto">
        <div
          v-for="(r, i) in results"
          :key="i"
          class="text-xs py-1 px-2 rounded"
          :class="r.success ? 'bg-green-500/10 text-green-400' : 'bg-red-500/10 text-red-400'"
        >
          {{ r.success ? '✓' : '✗' }}
          {{ r.oldPath }} → {{ r.newPath }}
          <span v-if="r.error" class="ml-2 opacity-70">{{ r.error }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  FilePen, FolderOpen, FolderSearch, Eye, FileIcon,
  Loader2, CheckCircle, AlertCircle
} from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { ListDirectory, BatchRenameFiles } from '@/../wailsjs/go/sysinfo/SysInfo'

const appStore = useAppStore()

// 文件条目类型
interface FileEntry {
  original: string
  renamed: string
  path: string
  error?: string
}

// 重命名结果类型
interface RenameResult {
  oldPath: string
  newPath: string
  success: boolean
  error: string
}

// 状态
const currentDir = ref('')
const files = ref<FileEntry[]>([])
const pattern = ref('')
const replacement = ref('')
const startNum = ref(1)
const useRegex = ref(true)
const caseSensitive = ref(false)
const previewOnly = ref(false)
const isRenaming = ref(false)
const results = ref<RenameResult[]>([])

// 计算属性
const successCount = computed(() => results.value.filter(r => r.success).length)
const allSuccess = computed(() => successCount.value === results.value.length && results.value.length > 0)

// 浏览目录（使用 prompt 输入路径，Wails 环境中可通过 dialog 选择）
async function browseDir() {
  const dir = prompt('请输入目录路径（绝对路径）：', currentDir.value || 'C:\\Users')
  if (!dir) return

  currentDir.value = dir
  results.value = []

  try {
    const entries = await ListDirectory(dir)
    files.value = entries
      .filter((e: any) => !e.isDir)
      .map((e: any) => ({
        original: e.name,
        renamed: e.name,
        path: e.path,
      }))

    if (files.value.length === 0) {
      appStore.showToast('info', '该目录下没有文件')
    } else {
      appStore.showToast('success', `已加载 ${files.value.length} 个文件`)
    }
  } catch (err) {
    appStore.showToast('error', `读取目录失败：${err}`)
    files.value = []
  }
}

// 预览重命名结果
function preview() {
  if (!pattern.value) {
    appStore.showToast('warning', '请输入查找模式')
    return
  }

  let counter = startNum.value
  files.value = files.value.map(f => {
    try {
      let newName: string
      if (useRegex.value) {
        const flags = caseSensitive.value ? 'g' : 'gi'
        const re = new RegExp(pattern.value, flags)
        newName = f.original.replace(re, replacement.value)
      } else {
        // 简单字符串替换（支持 /g 全局替换语义）
        const idx = caseSensitive.value
          ? f.original.indexOf(pattern.value)
          : f.original.toLowerCase().indexOf(pattern.value.toLowerCase())

        if (idx === -1) {
          newName = f.original
        } else {
          newName = f.original.substring(0, idx) + replacement.value + f.original.substring(idx + pattern.value.length)
        }
      }

      // 替换序号占位符 $N / $n / $NN
      newName = newName.replace(/\$(\d+)/g, (_, digits) => {
        return String(counter - 1 + parseInt(digits)).padStart(digits.length, '0')
      })
      if (newName.includes('$N')) {
        newName = newName.replace(/\$N/g, String(counter++).padStart(2, '0'))
      }
      if (newName.includes('$n')) {
        newName = newName.replace(/\$n/g, String(counter++))
      }

      return { ...f, renamed: newName, error: undefined }
    } catch (err) {
      return { ...f, renamed: f.original, error: `正则错误` }
    }
  })

  const changed = files.value.filter(f => f.renamed !== f.original).length
  appStore.showToast('info', `预览完成，${changed} 个文件将重命名`)
}

// 执行重命名
async function rename() {
  if (previewOnly.value) {
    appStore.showToast('info', '当前为「仅预览」模式，取消勾选后可执行')
    return
  }

  const changed = files.value.filter(f => f.renamed !== f.original && !f.error)
  if (changed.length === 0) {
    appStore.showToast('warning', '没有需要重命名的文件')
    return
  }

  // 确认操作
  if (!confirm(`确定要重命名 ${changed.length} 个文件吗？`)) return

  isRenaming.value = true
  results.value = []

  try {
    // 构建重命名对（格式：原路径|新路径）
    const pairs = changed.map(f => {
      // 新路径 = 原路径的前缀目录 + 新文件名
      const lastSep = Math.max(f.path.lastIndexOf('/'), f.path.lastIndexOf('\\'))
      const dirPrefix = f.path.substring(0, lastSep + 1)
      return `${f.path}|${dirPrefix}${f.renamed}`
    })

    const res = await BatchRenameFiles(pairs)
    results.value = res.map((r: any) => ({
      oldPath: r.oldPath,
      newPath: r.newPath,
      success: r.success,
      error: r.error || '',
    }))

    const ok = results.value.filter(r => r.success).length
    if (ok === changed.length) {
      appStore.showToast('success', `成功重命名 ${ok} 个文件`)
    } else {
      appStore.showToast('warning', `完成，${ok}/${changed.length} 个成功`)
    }

    // 刷新文件列表
    if (currentDir.value) {
      const entries = await ListDirectory(currentDir.value)
      files.value = entries
        .filter((e: any) => !e.isDir)
        .map((e: any) => ({
          original: e.name,
          renamed: e.name,
          path: e.path,
        }))
    }
  } catch (err) {
    appStore.showToast('error', `重命名失败：${err}`)
  } finally {
    isRenaming.value = false
  }
}
</script>

<style scoped>
.rule-grid {
  background: rgba(255, 255, 255, 0.02);
  border-radius: 8px;
  padding: 12px;
  border: 1px solid rgba(255, 255, 255, 0.06);
}

.light .rule-grid {
  background: rgba(0, 0, 0, 0.02);
  border-color: rgba(0, 0, 0, 0.06);
}

.file-list {
  border-top: 1px solid rgba(255, 255, 255, 0.06);
  padding-top: 8px;
}

.light .file-list {
  border-color: rgba(0, 0, 0, 0.06);
}

.file-row {
  display: flex;
  align-items: center;
  gap: 8px;
  padding: 6px 8px;
  border-radius: 6px;
  font-size: 13px;
  transition: background 0.1s;
  overflow: hidden;
}

.file-row:hover {
  background: rgba(255, 255, 255, 0.03);
}

.light .file-row:hover {
  background: rgba(0, 0, 0, 0.02);
}

.file-row.changed {
  color: var(--text-secondary, #888);
}

.file-row.error {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.05);
}

.original {
  color: var(--text-secondary, #888);
  text-decoration: line-through;
  opacity: 0.6;
  max-width: 200px;
}

.renamed {
  font-weight: 500;
  max-width: 200px;
}

.error-tag {
  font-size: 11px;
  padding: 2px 6px;
  background: rgba(239, 68, 68, 0.15);
  border-radius: 4px;
  color: #ef4444;
  white-space: nowrap;
}

.result-header {
  display: flex;
  align-items: center;
  gap: 8px;
}
</style>
