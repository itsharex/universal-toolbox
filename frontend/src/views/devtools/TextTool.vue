<template>
  <!-- 文本工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <FileText :size="20" class="text-primary-400" />
        文本工具
      </div>
      <div class="page-desc">字符统计 · 大小写转换 · 查找替换 · 文本对比</div>
    </div>

    <!-- Tab 栏 -->
    <div class="tab-bar mb-4">
      <button
        v-for="t in tabs"
        :key="t.id"
        :class="['tab-item', tab === t.id && 'active']"
        @click="tab = t.id"
      >
        <component :is="t.icon" :size="14" />
        {{ t.label }}
      </button>
    </div>

    <!-- 字符统计 -->
    <template v-if="tab === 'stats'">
      <div class="flex-1 flex flex-col min-h-0">
        <textarea
          v-model="statsInput"
          class="textarea-field flex-1 min-h-0 mb-3"
          placeholder="输入或粘贴文本进行统计..."
          @input="autoCalcStats"
        />
        <div class="grid grid-cols-4 gap-3">
          <div class="stat-card">
            <div class="stat-value">{{ stats.chars }}</div>
            <div class="stat-label">字符数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.charsNoSpace }}</div>
            <div class="stat-label">不含空格</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.words }}</div>
            <div class="stat-label">单词数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.lines }}</div>
            <div class="stat-label">行数</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.chinese }}</div>
            <div class="stat-label">中文字符</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.english }}</div>
            <div class="stat-label">英文字母</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.digits }}</div>
            <div class="stat-label">数字</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ stats.punctuation }}</div>
            <div class="stat-label">标点符号</div>
          </div>
        </div>
      </div>
    </template>

    <!-- 大小写转换 -->
    <template v-if="tab === 'case'">
      <div class="flex-1 flex flex-col min-h-0">
        <textarea
          v-model="caseInput"
          class="textarea-field flex-1 min-h-0 mb-3"
          placeholder="输入或粘贴文本..."
        />
        <div class="flex gap-2 mb-3 flex-wrap">
          <button @click="toUpper" class="btn btn-secondary">
            <ArrowUp :size="14" />
            全部大写
          </button>
          <button @click="toLower" class="btn btn-secondary">
            <ArrowDown :size="14" />
            全部小写
          </button>
          <button @click="toCapitalize" class="btn btn-secondary">
            <CaseSensitive :size="14" />
            首字母大写
          </button>
          <button @click="toToggle" class="btn btn-secondary">
            <RefreshCw :size="14" />
            大小写互换
          </button>
          <button @click="toCamelCase" class="btn btn-secondary">驼峰命名</button>
          <button @click="toSnakeCase" class="btn btn-secondary">下划线命名</button>
        </div>
        <div class="card flex-1 min-h-0 overflow-auto font-mono text-sm">
          <div v-if="!caseOutput" class="opacity-30">转换结果显示在这里...</div>
          <pre v-else class="whitespace-pre-wrap">{{ caseOutput }}</pre>
        </div>
      </div>
    </template>

    <!-- 查找替换 -->
    <template v-if="tab === 'replace'">
      <div class="flex-1 flex flex-col min-h-0">
        <div class="card mb-3">
          <div class="grid grid-cols-2 gap-3 mb-3">
            <div>
              <div class="label mb-1">查找内容</div>
              <input v-model="search" class="input-field" placeholder="要查找的内容..." />
            </div>
            <div>
              <div class="label mb-1">替换为</div>
              <input v-model="replace" class="input-field" placeholder="替换后的内容..." />
            </div>
          </div>
          <div class="flex items-center gap-4">
            <label class="flex items-center gap-2 cursor-pointer text-sm">
              <input type="checkbox" v-model="useRegex" />
              使用正则表达式
            </label>
            <label class="flex items-center gap-2 cursor-pointer text-sm">
              <input type="checkbox" v-model="caseSensitive" />
              区分大小写
            </label>
            <div class="flex-1" />
            <button @click="doReplace(false)" class="btn btn-secondary">替换第一个</button>
            <button @click="doReplace(true)" class="btn btn-primary">
              <Replace :size="14" />
              全部替换
            </button>
          </div>
        </div>
        <div class="flex-1 min-h-0">
          <textarea
            v-model="replaceInput"
            class="textarea-field h-full"
            placeholder="输入要处理的文本..."
          />
        </div>
        <div v-if="replaceCount !== null" class="text-sm opacity-60 mt-2">
          已替换 {{ replaceCount }} 处
        </div>
      </div>
    </template>

    <!-- 文本对比 -->
    <template v-if="tab === 'compare'">
      <div class="flex-1 flex flex-col min-h-0">
        <div class="two-col flex-1 min-h-0 mb-3">
          <div class="flex flex-col min-h-0">
            <div class="label mb-1 text-blue-400">原文</div>
            <textarea v-model="text1" class="textarea-field flex-1 min-h-0" placeholder="粘贴原文..." />
          </div>
          <div class="flex flex-col min-h-0">
            <div class="label mb-1 text-green-400">新文</div>
            <textarea v-model="text2" class="textarea-field flex-1 min-h-0" placeholder="粘贴新文..." />
          </div>
        </div>
        <div class="flex gap-2 mb-3">
          <button @click="swapTexts" class="btn btn-secondary">
            <ArrowLeftRight :size="14" />
            交换
          </button>
          <button @click="compare" class="btn btn-primary">
            <GitCompare :size="14" />
            开始对比
          </button>
        </div>
        <div class="card flex-1 min-h-0 overflow-auto">
          <div v-if="!compareResult" class="opacity-30 text-center py-8">对比结果显示在这里...</div>
          <div v-else class="space-y-1 font-mono text-sm">
            <div
              v-for="(line, idx) in compareResult.split('\n')"
              :key="idx"
              :class="[
                'px-2 py-1 rounded',
                line.startsWith('+') ? 'bg-green-500/10 text-green-400' :
                line.startsWith('-') ? 'bg-red-500/10 text-red-400' :
                'opacity-50'
              ]"
            >
              {{ line }}
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import {
  FileText, BarChart2, Replace, GitCompare,
  ArrowUp, ArrowDown, CaseSensitive, RefreshCw, ArrowLeftRight
} from 'lucide-vue-next'
import { TextCompare } from '../../../wailsjs/go/devtools/DevTools'

// Tab 配置
const tab = ref('stats')
const tabs = [
  { id: 'stats', label: '字符统计', icon: BarChart2 },
  { id: 'case', label: '大小写转换', icon: CaseSensitive },
  { id: 'replace', label: '查找替换', icon: Replace },
  { id: 'compare', label: '文本对比', icon: GitCompare },
]

// 字符统计
const statsInput = ref('')
const stats = ref({
  chars: 0, charsNoSpace: 0, words: 0, lines: 0,
  chinese: 0, english: 0, digits: 0, punctuation: 0
})

function autoCalcStats() {
  const text = statsInput.value
  stats.value = {
    chars: text.length,
    charsNoSpace: text.replace(/\s/g, '').length,
    words: text.trim() ? text.trim().split(/\s+/).length : 0,
    lines: text ? text.split('\n').length : 0,
    chinese: (text.match(/[\u4e00-\u9fa5]/g) || []).length,
    english: (text.match(/[a-zA-Z]/g) || []).length,
    digits: (text.match(/\d/g) || []).length,
    punctuation: (text.match(/[,.!?;:'"()（）。，！？；：""''、]/g) || []).length,
  }
}

// 大小写转换
const caseInput = ref('')
const caseOutput = ref('')

function toUpper() { caseOutput.value = caseInput.value.toUpperCase() }
function toLower() { caseOutput.value = caseInput.value.toLowerCase() }
function toCapitalize() {
  caseOutput.value = caseInput.value.replace(/\b\w/g, c => c.toUpperCase())
}
function toToggle() {
  caseOutput.value = caseInput.value.split('').map(c =>
    c === c.toUpperCase() ? c.toLowerCase() : c.toUpperCase()
  ).join('')
}
function toCamelCase() {
  caseOutput.value = caseInput.value
    .toLowerCase()
    .replace(/[-_\s]+(.)/g, (_, c) => c.toUpperCase())
    .replace(/^(.)/, c => c.toLowerCase())
}
function toSnakeCase() {
  caseOutput.value = caseInput.value
    .replace(/([a-z])([A-Z])/g, '$1_$2')
    .replace(/[-\s]+/g, '_')
    .toLowerCase()
}

// 查找替换
const search = ref('')
const replace = ref('')
const replaceInput = ref('')
const useRegex = ref(false)
const caseSensitive = ref(false)
const replaceCount = ref<number | null>(null)

function doReplace(all: boolean) {
  if (!search.value) return
  replaceCount.value = null

  try {
    const flags = caseSensitive.value ? 'g' : 'gi'
    if (!all) flags.replace('g', '')

    if (useRegex.value) {
      const regex = new RegExp(search.value, all ? flags : flags.replace('g', ''))
      const matches = replaceInput.value.match(new RegExp(search.value, flags))
      replaceCount.value = matches ? matches.length : 0
      replaceInput.value = replaceInput.value.replace(regex, replace.value)
    } else {
      const escaped = search.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
      const regex = new RegExp(escaped, all ? flags : flags.replace('g', ''))
      const matches = replaceInput.value.match(new RegExp(escaped, flags))
      replaceCount.value = matches ? matches.length : 0
      replaceInput.value = replaceInput.value.replace(regex, replace.value)
    }
  } catch {
    // 忽略正则错误
  }
}

// 文本对比
const text1 = ref('')
const text2 = ref('')
const compareResult = ref('')

function swapTexts() {
  const temp = text1.value
  text1.value = text2.value
  text2.value = temp
}

async function compare() {
  if (!text1.value || !text2.value) return

  const res = await TextCompare(text1.value, text2.value)
  compareResult.value = res.data || ''

  // 简单差异标记
  const lines1 = text1.value.split('\n')
  const lines2 = text2.value.split('\n')
  const result: string[] = []

  const maxLen = Math.max(lines1.length, lines2.length)
  for (let i = 0; i < maxLen; i++) {
    const l1 = lines1[i] || ''
    const l2 = lines2[i] || ''

    if (l1 === l2) {
      result.push(`  ${l1}`)
    } else {
      if (l1) result.push(`- ${l1}`)
      if (l2) result.push(`+ ${l2}`)
    }
  }

  compareResult.value = result.join('\n')
}
</script>

<style scoped>
.stat-card {
  background: var(--bg-card);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 16px;
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: var(--accent);
}

.stat-label {
  font-size: 12px;
  opacity: 0.6;
  margin-top: 4px;
}
</style>
