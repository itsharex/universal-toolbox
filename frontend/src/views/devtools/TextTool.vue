<template>
  <ToolPage title="文本处理" description="文本查找替换、大小写转换、去重">

    <!-- Tab 导航 -->
    <div class="tab-bar mb-4">
      <button
        v-for="t in tabs"
        :key="t.id"
        :class="['tab-item', tab === t.id && 'active']"
        @click="tab = t.id"
      >
        {{ t.label }}
      </button>
    </div>

    <!-- 字符统计 -->
    <template v-if="tab === 'stats'">
      <textarea
        v-model="statsInput"
        class="textarea-field mb-3"
        rows="8"
        placeholder="输入文本进行统计..."
        @input="autoStats"
      />
      <div class="flex gap-2 mb-3">
        <button @click="clearStats" class="btn btn-secondary">
          <Trash2 :size="14" />
          清空
        </button>
      </div>
      <!-- 统计结果 -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3 mb-3">
        <div class="card text-center">
          <div class="text-2xl font-bold text-primary-400">{{ statResult.chars || 0 }}</div>
          <div class="text-xs opacity-50">字符数</div>
        </div>
        <div class="card text-center">
          <div class="text-2xl font-bold text-green-400">{{ statResult.charsNoSpace || 0 }}</div>
          <div class="text-xs opacity-50">非空格字符</div>
        </div>
        <div class="card text-center">
          <div class="text-2xl font-bold text-yellow-400">{{ statResult.lines || 0 }}</div>
          <div class="text-xs opacity-50">行数</div>
        </div>
        <div class="card text-center">
          <div class="text-2xl font-bold text-cyan-400">{{ statResult.words || 0 }}</div>
          <div class="text-xs opacity-50">单词数</div>
        </div>
      </div>
      <div class="grid grid-cols-2 md:grid-cols-4 gap-3">
        <div class="card text-center">
          <div class="text-xl font-bold text-orange-400">{{ statResult.chinese || 0 }}</div>
          <div class="text-xs opacity-50">中文字符</div>
        </div>
        <div class="card text-center">
          <div class="text-xl font-bold text-purple-400">{{ statResult.numbers || 0 }}</div>
          <div class="text-xs opacity-50">数字</div>
        </div>
        <div class="card text-center">
          <div class="text-xl font-bold text-pink-400">{{ statResult.puncts || 0 }}</div>
          <div class="text-xs opacity-50">标点符号</div>
        </div>
        <div class="card text-center">
          <div class="text-xl font-bold text-blue-400">{{ statResult.bytes || 0 }} B</div>
          <div class="text-xs opacity-50">字节大小</div>
        </div>
      </div>
    </template>

    <!-- 查找替换 -->
    <template v-if="tab === 'replace'">
      <div class="card mb-3">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-3 mb-3">
          <div>
            <div class="label mb-1">查找</div>
            <input
              v-model="search"
              class="input-field"
              placeholder="查找内容..."
              @keyup.enter="doReplace"
            />
          </div>
          <div>
            <div class="label mb-1">替换为</div>
            <input
              v-model="replace"
              class="input-field"
              placeholder="替换为..."
              @keyup.enter="doReplace"
            />
          </div>
        </div>
        <div class="flex items-center gap-4 flex-wrap">
          <label class="flex items-center gap-2 text-sm cursor-pointer">
            <input type="checkbox" v-model="useRegex" class="checkbox-field" />
            使用正则
          </label>
          <label class="flex items-center gap-2 text-sm cursor-pointer">
            <input type="checkbox" v-model="caseSensitive" class="checkbox-field" />
            区分大小写
          </label>
          <label class="flex items-center gap-2 text-sm cursor-pointer">
            <input type="checkbox" v-model="globalReplace" class="checkbox-field" />
            全局替换
          </label>
          <div class="flex-1"></div>
          <button @click="doReplace" class="btn btn-primary">
            <Replace :size="14" />
            替换 ({{ matchCount }})
          </button>
          <button @click="copyResult" class="btn btn-secondary">
            <Copy :size="14" />
            复制
          </button>
        </div>
      </div>

      <!-- 匹配结果 -->
      <div v-if="search && replaceInput" class="mb-2 text-xs">
        <span class="text-green-400">找到 {{ matchCount }} 处匹配</span>
      </div>

      <!-- 原文输入 -->
      <div class="flex-1 flex flex-col min-h-0">
        <div class="label mb-1">原文</div>
        <textarea
          v-model="replaceInput"
          class="textarea-field flex-1 min-h-[200px]"
          placeholder="输入或粘贴要查找替换的文本..."
          spellcheck="false"
        />
      </div>

      <!-- 高亮预览 -->
      <div v-if="replaceInput" class="flex-1 flex flex-col min-h-0 mt-3">
        <div class="label mb-1">高亮预览</div>
        <div class="code-output flex-1 min-h-[100px] whitespace-pre-wrap">
          <span v-for="(part, idx) in highlightedText" :key="idx" :class="part.highlight ? 'bg-yellow-500/30 text-yellow-300' : ''">{{ part.text }}</span>
        </div>
      </div>

      <!-- 结果 -->
      <div v-if="replacedResult !== null" class="mt-3">
        <div class="label mb-1">替换结果</div>
        <div class="code-output flex-1 min-h-[100px] whitespace-pre-wrap text-green-400">{{ replacedResult }}</div>
      </div>
    </template>

    <!-- 文本对比 -->
    <template v-if="tab === 'compare'">
      <div class="grid grid-cols-2 gap-3 mb-3">
        <div class="flex-1">
          <div class="flex items-center justify-between mb-1">
            <div class="label">原文</div>
            <button @click="text1 = ''" class="text-xs opacity-50">清空</button>
          </div>
          <textarea
            v-model="text1"
            class="textarea-field h-48"
            placeholder="输入原文..."
          />
        </div>
        <div class="flex-1">
          <div class="flex items-center justify-between mb-1">
            <div class="label">新文</div>
            <button @click="text2 = ''" class="text-xs opacity-50">清空</button>
          </div>
          <textarea
            v-model="text2"
            class="textarea-field h-48"
            placeholder="输入新文..."
          />
        </div>
      </div>

      <div class="flex gap-2 mb-3">
        <button @click="compare" class="btn btn-primary">
          <GitCompare :size="14" />
          对比
        </button>
        <button @click="swapText" class="btn btn-secondary">
          <ArrowLeftRight :size="14" />
          交换
        </button>
      </div>

      <!-- 对比结果 -->
      <div v-if="compareResult" class="flex-1 overflow-auto">
        <div class="flex items-center gap-2 mb-2 text-sm">
          <span class="text-green-400">+{{ compareResult.added }}</span>
          <span class="text-red-400">-{{ compareResult.removed }}</span>
          <span class="opacity-50">{{ compareResult.similarity }}% 相似</span>
        </div>
        <!-- 差异展示 -->
        <div class="code-output p-3 text-sm font-mono whitespace-pre-wrap max-h-64 overflow-auto">
          <span v-for="(line, idx) in compareResult.diffLines" :key="idx" :class="line.type === 'add' ? 'text-green-400 bg-green-500/10' : line.type === 'remove' ? 'text-red-400 bg-red-500/10' : 'opacity-70'">{{ line.prefix }}{{ line.text }}
</span>
        </div>
      </div>
    </template>

    <!-- 行号处理 -->
    <template v-if="tab === 'lines'">
      <div class="flex gap-2 mb-3">
        <button @click="addLineNumbers" class="btn btn-primary">
          <Hash :size="14" />
          添加行号
        </button>
        <button @click="removeLineNumbers" class="btn btn-secondary">
          <MinusCircle :size="14" />
          去除行号
        </button>
        <button @click="sortLines" class="btn btn-secondary">
          <ArrowDownAZ :size="14" />
          排序
        </button>
        <button @click="removeDuplicateLines" class="btn btn-secondary">
          <Copy :size="14" />
          去重
        </button>
        <button @click="copyLinesResult" class="btn btn-secondary">
          <Copy :size="14" />
          复制
        </button>
      </div>
      <textarea
        v-model="linesInput"
        class="textarea-field flex-1 min-h-[300px]"
        placeholder="输入多行文本..."
      />
      <div v-if="linesResult" class="mt-3">
        <div class="label mb-1">结果</div>
        <div class="code-output flex-1 min-h-[100px] whitespace-pre">{{ linesResult }}</div>
      </div>
    </template>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, computed, watch } from 'vue'
import { FileText, BarChart2, Replace, GitCompare, Trash2, Copy, ArrowLeftRight, Hash, MinusCircle, ArrowDownAZ } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 状态
const tab = ref('stats')
const tabs = [
  { id: 'stats', label: '字符统计' },
  { id: 'replace', label: '查找替换' },
  { id: 'compare', label: '文本对比' },
  { id: 'lines', label: '行号处理' }
]

// 统计
const statsInput = ref('')
const statResult = ref<{
  chars: number
  charsNoSpace: number
  lines: number
  words: number
  chinese: number
  numbers: number
  puncts: number
  bytes: number
}>({ chars: 0, charsNoSpace: 0, lines: 0, words: 0, chinese: 0, numbers: 0, puncts: 0, bytes: 0 })

// 查找替换
const search = ref('')
const replace = ref('')
const replaceInput = ref('')
const useRegex = ref(false)
const caseSensitive = ref(true)
const globalReplace = ref(true)
const replacedResult = ref('')

// 文本对比
const text1 = ref('')
const text2 = ref('')
const compareResult = ref<any>(null)

// 行号处理
const linesInput = ref('')
const linesResult = ref('')

// 自动统计
function autoStats() {
  const text = statsInput.value
  if (!text) {
    statResult.value = { chars: 0, charsNoSpace: 0, lines: 0, words: 0, chinese: 0, numbers: 0, puncts: 0, bytes: 0 }
    return
  }

  statResult.value = {
    chars: text.length,
    charsNoSpace: text.replace(/\s/g, '').length,
    lines: text.split('\n').length,
    words: text.split(/\s+/).filter(w => w).length,
    chinese: (text.match(/[\u4e00-\u9fa5]/g) || []).length,
    numbers: (text.match(/\d/g) || []).length,
    puncts: (text.match(/[.,;:!?，。；！？""''【】《》()（）]/g) || []).length,
    bytes: new TextEncoder().encode(text).length
  }
}

function clearStats() {
  statsInput.value = ''
  autoStats()
}

// 匹配计数
const matchCount = computed(() => {
  if (!search.value || !replaceInput.value) return 0
  try {
    const flags = (caseSensitive.value ? '' : 'i') + (globalReplace.value ? 'g' : '')
    const regex = useRegex.value ? new RegExp(search.value, flags) : new RegExp(search.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), flags)
    const matches = replaceInput.value.match(regex)
    return matches ? matches.length : 0
  } catch {
    return 0
  }
})

// 高亮文本
const highlightedText = computed(() => {
  if (!search.value || !replaceInput.value) return [{ text: replaceInput.value, highlight: false }]

  try {
    const flags = (caseSensitive.value ? '' : 'i') + 'g'
    const regex = useRegex.value ? new RegExp(search.value, flags) : new RegExp(search.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), flags)

    const parts: Array<{ text: string; highlight: boolean }> = []
    let lastIndex = 0
    let match

    while ((match = regex.exec(replaceInput.value)) !== null) {
      if (match.index > lastIndex) {
        parts.push({ text: replaceInput.value.slice(lastIndex, match.index), highlight: false })
      }
      parts.push({ text: match[0], highlight: true })
      lastIndex = regex.lastIndex

      if (!globalReplace.value) break
      if (match[0] === '') regex.lastIndex++
    }

    if (lastIndex < replaceInput.value.length) {
      parts.push({ text: replaceInput.value.slice(lastIndex), highlight: false })
    }

    return parts.length ? parts : [{ text: replaceInput.value, highlight: false }]
  } catch {
    return [{ text: replaceInput.value, highlight: false }]
  }
})

// 执行替换
function doReplace() {
  if (!search.value || !replaceInput.value) return

  try {
    const flags = (caseSensitive.value ? '' : 'i') + (globalReplace.value ? 'g' : '')
    const regex = useRegex.value ? new RegExp(search.value, flags) : new RegExp(search.value.replace(/[.*+?^${}()|[\]\\]/g, '\\$&'), flags)
    replacedResult.value = replaceInput.value.replace(regex, replace.value)
    appStore.showToast('success', `已替换 ${matchCount.value} 处`)
  } catch (e) {
    appStore.showToast('error', String(e))
  }
}

async function copyResult() {
  if (!replacedResult.value) return
  await navigator.clipboard.writeText(replacedResult.value)
  appStore.showToast('success', '已复制')
}

// 对比
function compare() {
  if (!text1.value || !text2.value) return

  const lines1 = text1.value.split('\n')
  const lines2 = text2.value.split('\n')

  const added = lines2.filter(l => !lines1.includes(l)).length
  const removed = lines1.filter(l => !lines2.includes(l)).length
  const common = lines1.filter(l => lines2.includes(l)).length
  const total = Math.max(lines1.length, lines2.length)
  const similarity = total > 0 ? Math.round((common / total) * 100) : 100

  // 简化的差异计算
  const diffLines: Array<{ type: string; prefix: string; text: string }> = []
  const maxLen = Math.max(lines1.length, lines2.length)

  for (let i = 0; i < maxLen; i++) {
    const l1 = lines1[i]
    const l2 = lines2[i]

    if (l1 === l2) {
      diffLines.push({ type: 'same', prefix: '  ', text: l1 || '' })
    } else {
      if (l1 !== undefined && !lines2.includes(l1)) {
        diffLines.push({ type: 'remove', prefix: '- ', text: l1 })
      }
      if (l2 !== undefined && !lines1.includes(l2)) {
        diffLines.push({ type: 'add', prefix: '+ ', text: l2 })
      }
    }
  }

  compareResult.value = { added, removed, similarity, diffLines }
  appStore.showToast('success', `相似度 ${similarity}%`)
}

function swapText() {
  const temp = text1.value
  text1.value = text2.value
  text2.value = temp
  if (compareResult.value) compare()
}

// 行号处理
function addLineNumbers() {
  const lines = linesInput.value.split('\n')
  linesResult.value = lines.map((l, i) => `${i + 1}\t${l}`).join('\n')
  appStore.showToast('success', `已添加 ${lines.length} 个行号`)
}

function removeLineNumbers() {
  const lines = linesInput.value.split('\n')
  linesResult.value = lines.map(l => l.replace(/^\d+\s+/, '')).join('\n')
  appStore.showToast('success', '已去除行号')
}

function sortLines() {
  const lines = linesInput.value.split('\n')
  linesResult.value = lines.sort().join('\n')
  appStore.showToast('success', `已排序 ${lines.length} 行`)
}

function removeDuplicateLines() {
  const lines = linesInput.value.split('\n')
  const unique = [...new Set(lines)]
  linesResult.value = unique.join('\n')
  appStore.showToast('success', `去重后 ${unique.length} 行`)
}

async function copyLinesResult() {
  if (!linesResult.value) return
  await navigator.clipboard.writeText(linesResult.value)
  appStore.showToast('success', '已复制')
}
</script>