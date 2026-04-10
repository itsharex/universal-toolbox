<template>
  <ToolPage title="批量重命名" description="批量文件重命名">
<!-- 目录和文件匹配 -->
    <div class="card mb-4">
      <div class="label mb-2">文件目录</div>
      <div class="flex gap-2 mb-3">
        <input v-model="directory" class="input-field flex-1" placeholder="选择文件所在目录..." readonly />
        <button @click="selectDirectory" class="btn btn-secondary"><FolderOpen :size="14"/>选择目录</button>
      </div>
      <div class="label mb-2">文件匹配模式</div>
      <input v-model="filePattern" class="input-field mb-3" placeholder="如 * 或 *.txt" />
      <button @click="loadFiles" class="btn btn-secondary" :disabled="!directory">
        <Upload :size="14"/>加载文件
      </button>
    </div>

    <!-- 重命名模式选择 -->
    <div class="card mb-4">
      <div class="label mb-2">重命名模式</div>
      <div class="flex gap-2 flex-wrap">
        <button :class="['btn', renameMode==='regex'?'btn-primary':'btn-secondary']" @click="renameMode='regex'">正则替换</button>
        <button :class="['btn', renameMode==='prefix'?'btn-primary':'btn-secondary']" @click="renameMode='prefix'">添加前缀</button>
        <button :class="['btn', renameMode==='suffix'?'btn-primary':'btn-secondary']" @click="renameMode='suffix'">添加后缀</button>
        <button :class="['btn', renameMode==='sequence'?'btn-primary':'btn-secondary']" @click="renameMode='sequence'">序号命名</button>
      </div>
    </div>

    <!-- 模式参数 -->
    <div class="card mb-4">
      <!-- 正则替换 -->
      <template v-if="renameMode==='regex'">
        <div class="flex gap-2 mb-3">
          <input v-model="pattern" class="input-field flex-1" placeholder="查找（支持正则）..."/>
          <input v-model="replacement" class="input-field flex-1" placeholder="替换为（$1 引用捕获组）..."/>
        </div>
      </template>
      <!-- 前缀 -->
      <template v-if="renameMode==='prefix'">
        <div class="flex gap-2 mb-3">
          <input v-model="prefixText" class="input-field flex-1" placeholder="输入前缀文本..."/>
        </div>
      </template>
      <!-- 后缀 -->
      <template v-if="renameMode==='suffix'">
        <div class="flex gap-2 mb-3">
          <input v-model="suffixText" class="input-field flex-1" placeholder="输入后缀文本（在扩展名之前）..."/>
        </div>
      </template>
      <!-- 序号 -->
      <template v-if="renameMode==='sequence'">
        <div class="flex gap-2 mb-3">
          <input v-model="seqBaseName" class="input-field flex-1" placeholder="基础名称（如 photo）..."/>
          <input v-model.number="seqStart" class="input-field w-24" placeholder="起始序号" type="number" min="0"/>
          <input v-model="seqSeparator" class="input-field w-16" placeholder="分隔符" value="_"/>
        </div>
      </template>

      <div class="flex gap-2">
        <button @click="preview" class="btn btn-secondary">预览</button>
        <button @click="doRename" class="btn btn-primary">执行重命名</button>
        <button @click="clearFiles" class="btn btn-secondary">清空列表</button>
      </div>
    </div>

    <!-- 文件列表预览 -->
    <div v-if="files.length" class="card">
      <div class="label mb-2">文件列表 ({{ files.length }} 个文件)</div>
      <div class="space-y-1 max-h-64 overflow-auto">
        <div v-for="(f,i) in files" :key="i" class="flex gap-3 text-xs py-1">
          <span class="opacity-60 flex-1 truncate">{{ f.original }}</span>
          <span class="text-primary-400">-></span>
          <span class="flex-1 truncate" :class="f.renamed !== f.original ? 'text-green-400' : ''">{{ f.renamed || f.original }}</span>
        </div>
      </div>
    </div>
  </div>
  </ToolPage>
</template>
<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref } from 'vue'
import { FilePen, Upload, FolderOpen } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { BatchRename } from '../../../wailsjs/go/sysinfo/SysInfo'
const appStore = useAppStore()

const files = ref<{original:string,renamed:string,path:string}[]>([])
const directory = ref('')
const filePattern = ref('*')
const renameMode = ref<'regex'|'prefix'|'suffix'|'sequence'>('prefix')
const pattern = ref(''), replacement = ref('')
const prefixText = ref(''), suffixText = ref('')
const seqBaseName = ref(''), seqStart = ref(1), seqSeparator = ref('_')

// 选择目录
async function selectDirectory() {
  try {
    const dir = prompt('请输入文件所在目录路径:')
    if (dir) {
      directory.value = dir
      appStore.showToast('success', `已设置目录: ${dir}`)
    }
  } catch (e) {
    appStore.showToast('error', '设置目录失败: ' + String(e))
  }
}

// 加载文件列表（从目录路径提取文件名作为预览）
async function loadFiles() {
  if (!directory.value) return
  files.value = []
  // 由于后端没有扫描方法，这里仅做目录确认提示
  // 实际重命名由后端根据 directory + filePattern 执行
  appStore.showToast('success', `已设置目录: ${directory.value}，模式: ${filePattern.value}`)
  // 添加占位提示
  files.value = [{ original: `(${filePattern.value} 匹配的文件)`, renamed: '将在执行时由后端处理', path: '' }]
}

function getFileNameWithoutExt(name: string): { base: string; ext: string } {
  const lastDot = name.lastIndexOf('.')
  if (lastDot === -1) return { base: name, ext: '' }
  return { base: name.substring(0, lastDot), ext: name.substring(lastDot) }
}

function preview() {
  files.value = files.value.map((f, idx) => {
    const { base, ext } = getFileNameWithoutExt(f.original)
    let newBase = base
    try {
      switch (renameMode.value) {
        case 'regex':
          if (pattern.value) {
            const re = new RegExp(pattern.value, 'g')
            newBase = base.replace(re, replacement.value)
          }
          break
        case 'prefix':
          newBase = prefixText.value + base
          break
        case 'suffix':
          newBase = base + suffixText.value
          break
        case 'sequence':
          newBase = seqBaseName.value + seqSeparator.value + String(seqStart.value + idx).padStart(3, '0')
          break
      }
    } catch {
      // 正则错误时保持原名
    }
    return { ...f, renamed: newBase + ext }
  })
}

async function doRename() {
  if (!directory.value) {
    appStore.showToast('warning', '请先选择文件目录')
    return
  }

  // 根据模式构建 value 和 replace 参数
  let mode = renameMode.value
  let value = ''
  let replace = ''

  switch (renameMode.value) {
    case 'regex':
      value = pattern.value
      replace = replacement.value
      break
    case 'prefix':
      value = prefixText.value
      replace = ''
      break
    case 'suffix':
      value = suffixText.value
      replace = ''
      break
    case 'sequence':
      value = seqBaseName.value + seqSeparator.value + String(seqStart.value)
      replace = ''
      break
  }

  try {
    const result = await BatchRename(directory.value, filePattern.value, mode, value, replace) as any
    if (result) {
      const count = result.RenamedCount || result.renamedCount || result.count || 0
      appStore.showToast('success', `成功重命名 ${count} 个文件`)
      files.value = []
    } else {
      appStore.showToast('error', '重命名失败')
    }
  } catch (e) {
    appStore.showToast('error', '重命名失败: ' + String(e))
  }
}

function clearFiles() {
  files.value = []
}
</script>
<style scoped>
.drop-zone {
  border: 2px dashed var(--border-color);
  border-radius: 10px; padding: 32px;
  text-align: center; cursor: pointer;
  transition: border-color 0.2s;
}
.drop-zone:hover { border-color: var(--accent); }
</style>
