<template>
  <div class="page-container">
    <!-- 标题 -->
    <div class="mb-4">
      <div class="page-title">
        <Regex :size="20" class="text-primary-400" />
        正则测试
      </div>
      <div class="page-desc">测试正则表达式 · 匹配高亮 · 常用模板</div>
    </div>

    <!-- 正则输入 -->
    <div class="card mb-3">
      <div class="flex gap-2 items-center flex-wrap">
        <span class="opacity-50 font-mono text-lg">/</span>
        <input
          v-model="pattern"
          class="input-field flex-1 min-w-[200px]"
          placeholder="输入正则表达式，如 \d{4}-\d{2}-\d{2}"
          @keyup.enter="test"
        />
        <div class="flex gap-1">
          <button
            v-for="f in flags"
            :key="f.flag"
            @click="toggleFlag(f.flag)"
            :class="['px-2 py-1 text-xs font-mono rounded transition-colors', activeFlags.includes(f.flag) ? 'bg-primary-500/30 text-primary-300' : 'bg-white/5 opacity-50']"
            :title="f.desc"
          >
            {{ f.flag }}
          </button>
        </div>
        <button @click="test" class="btn btn-primary">
          <Play :size="14" />
          测试
        </button>
      </div>
      <!-- 正则错误提示 -->
      <div v-if="regexError" class="mt-2 text-sm text-red-400 flex items-center gap-2">
        <AlertCircle :size="14" />
        {{ regexError }}
      </div>
    </div>

    <!-- 常用模板 -->
    <div class="mb-3">
      <div class="flex items-center gap-2 mb-2">
        <span class="text-xs opacity-50">常用模板</span>
        <div class="flex-1 h-px bg-white/5"></div>
      </div>
      <div class="flex gap-2 flex-wrap">
        <button
          v-for="t in templates"
          :key="t.name"
          @click="applyTemplate(t)"
          class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10 transition-colors"
        >
          {{ t.name }}
        </button>
      </div>
    </div>

    <!-- 测试文本 -->
    <div class="mb-3">
      <div class="flex items-center justify-between mb-1">
        <div class="label">测试文本</div>
        <div class="text-xs opacity-50">{{ testText.length }} 字符</div>
      </div>
      <textarea
        v-model="testText"
        class="textarea-field"
        rows="6"
        placeholder="输入要测试的文本..."
      />
    </div>

    <!-- 匹配结果 -->
    <div class="flex-1 flex flex-col min-h-0">
      <div class="flex items-center justify-between mb-1">
        <div class="label">匹配结果</div>
        <div v-if="matches.length" class="text-xs">
          <span class="text-green-400">{{ matches.length }}</span> 个匹配
        </div>
      </div>

      <!-- 空状态 -->
      <div v-if="!matches.length && !regexError" class="flex-1 flex items-center justify-center opacity-30">
        <div class="text-center">
          <Regex :size="32" class="mx-auto mb-2 opacity-50" />
          <div class="text-sm">输入正则和文本进行测试</div>
        </div>
      </div>

      <!-- 匹配列表 -->
      <div v-else-if="matches.length" class="flex-1 overflow-auto space-y-2">
        <div
          v-for="(m, idx) in matches"
          :key="idx"
          class="card bg-white/5"
        >
          <div class="flex items-center gap-2 mb-2">
            <span class="w-6 h-6 rounded-full bg-primary-500/20 text-primary-400 flex items-center justify-center text-xs">
              {{ idx + 1 }}
            </span>
            <span class="font-mono text-sm flex-1 truncate">{{ m.match }}</span>
            <span class="text-xs opacity-50">位置 {{ m.index }}</span>
          </div>
          <!-- 分组 -->
          <div v-if="m.groups && m.groups.length" class="ml-8 space-y-1">
            <div
              v-for="(g, gIdx) in m.groups"
              :key="gIdx"
              class="flex items-center gap-2 text-xs"
            >
              <span class="opacity-50">${{ gIdx + 1 }}</span>
              <span class="font-mono text-cyan-400">{{ g || '(空)' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 高亮预览 -->
    <div v-if="testText && matches.length" class="mt-3">
      <div class="label mb-1">高亮预览</div>
      <div class="code-output text-sm max-h-32 overflow-auto">
        <span v-for="(part, idx) in highlightedParts" :key="idx">
          <span v-if="part.highlight" class="bg-yellow-500/30 text-yellow-300 px-0.5 rounded">{{ part.text }}</span>
          <span v-else>{{ part.text }}</span>
        </span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Regex, Play, AlertCircle } from 'lucide-vue-next'

// 状态
const pattern = ref('')
const testText = ref('')
const activeFlags = ref<string[]>(['g'])
const matches = ref<Array<{ match: string; index: number; groups: string[] }>>([])
const regexError = ref('')

// 正则标志
const flags = [
  { flag: 'g', desc: '全局匹配' },
  { flag: 'i', desc: '忽略大小写' },
  { flag: 'm', desc: '多行模式' },
  { flag: 's', desc: '点匹配换行' }
]

// 常用模板
const templates = [
  { name: '邮箱', pattern: '[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}', flags: ['g', 'i'] },
  { name: '手机号', pattern: '1[3-9]\\d{9}', flags: ['g'] },
  { name: '日期', pattern: '\\d{4}[-/]\\d{1,2}[-/]\\d{1,2}', flags: ['g'] },
  { name: '时间', pattern: '\\d{1,2}:\\d{2}(:\\d{2})?', flags: ['g'] },
  { name: 'URL', pattern: 'https?://[^\\s]+', flags: ['g', 'i'] },
  { name: 'IP地址', pattern: '\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}\\.\\d{1,3}', flags: ['g'] },
  { name: '数字', pattern: '-?\\d+(\\.\\d+)?', flags: ['g'] },
  { name: '中文', pattern: '[\\u4e00-\\u9fa5]+', flags: ['g'] },
  { name: '身份证', pattern: '\\d{17}[\\dXx]', flags: ['g'] },
  { name: '十六进制', pattern: '#?[0-9a-fA-F]{6}', flags: ['g'] }
]

// 切换标志
function toggleFlag(flag: string) {
  const idx = activeFlags.value.indexOf(flag)
  if (idx >= 0) {
    activeFlags.value.splice(idx, 1)
  } else {
    activeFlags.value.push(flag)
  }
}

// 应用模板
function applyTemplate(t: typeof templates[0]) {
  pattern.value = t.pattern
  activeFlags.value = [...t.flags]
  test()
}

// 测试正则
function test() {
  matches.value = []
  regexError.value = ''

  if (!pattern.value || !testText.value) return

  try {
    const flags = activeFlags.value.join('')
    const regex = new RegExp(pattern.value, flags)

    let match
    while ((match = regex.exec(testText.value)) !== null) {
      matches.value.push({
        match: match[0],
        index: match.index,
        groups: match.slice(1)
      })

      // 避免无限循环（非全局模式）
      if (!activeFlags.value.includes('g')) break
      // 零宽匹配时前进一位
      if (match[0] === '') regex.lastIndex++
    }
  } catch (e) {
    regexError.value = String(e)
  }
}

// 高亮部分（用于预览）
const highlightedParts = computed(() => {
  if (!testText.value || !matches.value.length) return []

  const parts: Array<{ text: string; highlight: boolean }> = []
  const highlightRanges = matches.value.map(m => ({ start: m.index, end: m.index + m.match.length }))
  highlightRanges.sort((a, b) => a.start - b.start)

  let lastEnd = 0
  for (const range of highlightRanges) {
    if (range.start > lastEnd) {
      parts.push({ text: testText.value.slice(lastEnd, range.start), highlight: false })
    }
    if (range.end > lastEnd) {
      parts.push({ text: testText.value.slice(Math.max(lastEnd, range.start), range.end), highlight: true })
    }
    lastEnd = Math.max(lastEnd, range.end)
  }
  if (lastEnd < testText.value.length) {
    parts.push({ text: testText.value.slice(lastEnd), highlight: false })
  }

  return parts
})
</script>
