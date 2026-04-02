<template>
  <!-- 正则测试工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <Regex :size="20" class="text-primary-400" />
        正则测试
      </div>
      <div class="page-desc">实时测试正则表达式 · 匹配高亮 · 常用模板</div>
    </div>

    <!-- 正则输入 -->
    <div class="card mb-3">
      <div class="flex gap-2 items-center">
        <span class="opacity-50 font-mono text-lg">/</span>
        <input
          v-model="pattern"
          class="input-field flex-1 font-mono"
          placeholder="输入正则表达式，如 \d{4}-\d{2}-\d{2}"
          @input="debouncedTest"
        />
        <div class="flex gap-1">
          <button
            v-for="flag in ['g', 'i', 'm', 's']"
            :key="flag"
            @click="toggleFlag(flag)"
            :class="[
              'px-2 py-1 rounded text-xs font-mono transition-all',
              flags.includes(flag) ? 'bg-primary-500 text-white' : 'bg-white/5 opacity-50'
            ]"
          >
            {{ flag }}
          </button>
        </div>
        <button @click="test" class="btn btn-primary">
          <Play :size="14" />
          测试
        </button>
      </div>
      <!-- 正则错误提示 -->
      <div v-if="regexError" class="text-xs text-red-400 mt-2 flex items-center gap-1">
        <AlertCircle :size="12" />
        {{ regexError }}
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="flex-1 flex gap-4 min-h-0">
      <!-- 左侧：测试文本 -->
      <div class="flex-1 flex flex-col min-h-0">
        <div class="label mb-1">
          测试文本
          <span class="text-xs opacity-50">{{ testText.length }} 字符</span>
        </div>
        <textarea
          v-model="testText"
          class="textarea-field flex-1 min-h-0 font-mono text-sm"
          placeholder="输入要测试的文本..."
          @input="debouncedTest"
        />
      </div>

      <!-- 右侧：结果和常用正则 -->
      <div class="w-72 flex flex-col gap-3 shrink-0">
        <!-- 匹配结果 -->
        <div class="card flex-1 min-h-0 flex flex-col">
          <div class="font-semibold mb-2 flex items-center gap-2">
            <CheckCircle :size="14" class="text-green-400" />
            匹配结果
            <span v-if="matches.length" class="text-xs opacity-50">{{ matches.length }} 个</span>
          </div>
          <div class="flex-1 overflow-auto">
            <div v-if="!matches.length" class="text-sm opacity-30">无匹配结果</div>
            <div v-else class="space-y-2">
              <div
                v-for="(m, idx) in matches"
                :key="idx"
                class="p-2 rounded bg-white/5 font-mono text-sm"
              >
                <div class="flex items-center gap-2 mb-1">
                  <span class="text-xs px-1.5 py-0.5 rounded bg-primary-500/20 text-primary-400">#{{ idx + 1 }}</span>
                  <span class="flex-1 truncate">{{ m[0] }}</span>
                  <button @click="copyMatch(m[0])" class="opacity-50 hover:opacity-100">
                    <Copy :size="12" />
                  </button>
                </div>
                <div v-if="m.length > 1" class="text-xs opacity-50 pl-6">
                  分组: {{ m.slice(1).join(', ') }}
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 常用正则 -->
        <div class="card">
          <div class="font-semibold mb-2">常用正则</div>
          <div class="space-y-1 max-h-48 overflow-auto">
            <button
              v-for="r in commonRegex"
              :key="r.name"
              @click="applyRegex(r.pattern)"
              class="w-full text-left px-2 py-1.5 rounded hover:bg-white/5 text-sm flex items-center justify-between group"
            >
              <span>{{ r.name }}</span>
              <span class="text-xs opacity-0 group-hover:opacity-50 font-mono">{{ r.pattern.slice(0, 15) }}...</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- 匹配高亮预览 -->
    <div v-if="highlightedText" class="card mt-3">
      <div class="font-semibold mb-2">匹配高亮预览</div>
      <div class="p-3 rounded bg-white/5 font-mono text-sm whitespace-pre-wrap break-all" v-html="highlightedText" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Regex, Play, AlertCircle, CheckCircle, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 状态
const pattern = ref('')
const testText = ref('')
const flags = ref<string[]>(['g'])
const matches = ref<string[][]>([])
const regexError = ref('')

// 常用正则列表
const commonRegex = [
  { name: '手机号', pattern: '1[3-9]\\d{9}' },
  { name: '邮箱', pattern: '[\\w.-]+@[\\w.-]+\\.\\w+' },
  { name: 'URL', pattern: 'https?://[^\\s]+' },
  { name: 'IP 地址', pattern: '\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}' },
  { name: '日期', pattern: '\\d{4}[-/]\\d{2}[-/]\\d{2}' },
  { name: '时间', pattern: '\\d{2}:\\d{2}(:\\d{2})?' },
  { name: '身份证号', pattern: '\\d{17}[\\dXx]' },
  { name: '中文字符', pattern: '[\\u4e00-\\u9fa5]+' },
  { name: '数字', pattern: '-?\\d+(\\.\\d+)?' },
  { name: '空白字符', pattern: '\\s+' },
  { name: 'HTML 标签', pattern: '<[^>]+>' },
  { name: '十六进制颜色', pattern: '#[0-9a-fA-F]{6}' },
]

// 高亮文本
const highlightedText = computed(() => {
  if (!testText.value || !matches.value.length) return ''

  let result = testText.value
  let offset = 0

  // 按位置排序匹配
  const sortedMatches = [...matches.value].sort((a, b) => {
    const posA = testText.value.indexOf(a[0])
    const posB = testText.value.indexOf(b[0])
    return posA - posB
  })

  for (const match of sortedMatches) {
    const text = match[0]
    const pos = result.indexOf(text, offset > 0 ? offset - 20 : 0)
    if (pos >= 0) {
      const highlighted = `<mark class="bg-yellow-400/30 text-inherit rounded px-0.5">${escapeHtml(text)}</mark>`
      result = result.slice(0, pos) + highlighted + result.slice(pos + text.length)
      offset = pos + highlighted.length
    }
  }

  return result
})

// HTML 转义
function escapeHtml(text: string) {
  return text
    .replace(/&/g, '&amp;')
    .replace(/</g, '&lt;')
    .replace(/>/g, '&gt;')
}

// 切换标志
function toggleFlag(flag: string) {
  const idx = flags.value.indexOf(flag)
  if (idx >= 0) {
    flags.value.splice(idx, 1)
  } else {
    flags.value.push(flag)
  }
  test()
}

// 应用常用正则
function applyRegex(p: string) {
  pattern.value = p
  test()
}

// 复制匹配
async function copyMatch(text: string) {
  await navigator.clipboard.writeText(text)
  appStore.showToast('success', '已复制')
}

// 防抖测试
let debounceTimer: number
function debouncedTest() {
  clearTimeout(debounceTimer)
  debounceTimer = window.setTimeout(test, 300)
}

// 执行测试
function test() {
  regexError.value = ''
  matches.value = []

  if (!pattern.value || !testText.value) return

  try {
    const flagStr = flags.value.join('')
    const regex = new RegExp(pattern.value, flagStr)
    const allMatches: string[][] = []

    let match
    while ((match = regex.exec(testText.value)) !== null) {
      allMatches.push([...match])
      // 防止零宽匹配死循环
      if (match[0].length === 0) regex.lastIndex++
    }

    matches.value = allMatches
  } catch (err) {
    regexError.value = String(err)
  }
}
</script>

<style scoped>
mark {
  background: rgba(250, 204, 21, 0.3);
  border-radius: 2px;
  padding: 0 2px;
}
</style>
