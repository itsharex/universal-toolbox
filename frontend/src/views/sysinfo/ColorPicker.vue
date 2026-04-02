<template>
  <!-- 取色器工具页面：屏幕取色、颜色转换、历史记录 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <Pipette :size="20" class="text-primary-400" />
        取色器
      </div>
      <div class="page-desc">屏幕取色 · 颜色转换 · 历史记录</div>
    </div>

    <!-- 主内容区 -->
    <div class="flex-1 flex flex-col gap-4 min-h-0">
      <!-- 当前颜色展示 -->
      <div class="card flex items-center gap-6">
        <!-- 颜色预览块 -->
        <div
          class="w-24 h-24 rounded-xl shadow-lg cursor-pointer transition-all hover:scale-105"
          :style="{ background: currentColor.hex, boxShadow: `0 8px 32px ${currentColor.hex}40` }"
          @click="copyColor(currentColor.hex)"
          v-tooltip="'点击复制 HEX'"
        />

        <!-- 颜色信息 -->
        <div class="flex-1 space-y-3">
          <!-- HEX -->
          <div class="flex items-center gap-3">
            <span class="w-12 text-xs opacity-50">HEX</span>
            <code class="color-value" @click="copyColor(currentColor.hex)">{{ currentColor.hex }}</code>
          </div>
          <!-- RGB -->
          <div class="flex items-center gap-3">
            <span class="w-12 text-xs opacity-50">RGB</span>
            <code class="color-value" @click="copyColor(currentColor.rgb)">{{ currentColor.rgb }}</code>
          </div>
          <!-- HSL -->
          <div class="flex items-center gap-3">
            <span class="w-12 text-xs opacity-50">HSL</span>
            <code class="color-value" @click="copyColor(currentColor.hsl)">{{ currentColor.hsl }}</code>
          </div>
        </div>

        <!-- 取色按钮 -->
        <button @click="startPick" :disabled="isPicking" class="btn btn-primary text-base px-6 py-3">
          <Pipette v-if="!isPicking" :size="18" />
          <Loader2 v-else :size="18" class="loading-spin" />
          {{ isPicking ? '取色中...' : '开始取色' }}
        </button>
      </div>

      <!-- 颜色格式转换表 -->
      <div class="card">
        <div class="font-semibold mb-3 flex items-center gap-2">
          <Palette :size="16" />
          颜色格式
        </div>
        <div class="grid grid-cols-2 gap-3">
          <div v-for="format in colorFormats" :key="format.label" class="format-item">
            <span class="text-xs opacity-50">{{ format.label }}</span>
            <code class="flex-1 text-right" @click="copyColor(format.value)">{{ format.value }}</code>
          </div>
        </div>
      </div>

      <!-- 历史记录 -->
      <div class="card flex-1 min-h-0 flex flex-col">
        <div class="flex items-center justify-between mb-3">
          <div class="font-semibold flex items-center gap-2">
            <History :size="16" />
            历史记录
          </div>
          <button v-if="history.length" @click="clearHistory" class="btn btn-secondary text-xs py-1">
            <Trash2 :size="12" />
            清空
          </button>
        </div>

        <!-- 空状态 -->
        <div v-if="!history.length" class="flex-1 flex items-center justify-center opacity-30">
          <div class="text-center">
            <Pipette :size="32" class="mx-auto mb-2 opacity-50" />
            <div class="text-sm">暂无取色记录</div>
          </div>
        </div>

        <!-- 历史列表 -->
        <div v-else class="flex-1 overflow-auto">
          <div class="grid grid-cols-8 gap-2">
            <div
              v-for="(item, idx) in history"
              :key="idx"
              class="history-item"
              :style="{ background: item.hex }"
              @click="selectHistory(item)"
              v-tooltip="item.hex"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { Pipette, Palette, History, Trash2, Loader2 } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 取色状态
const isPicking = ref(false)

// 当前颜色
const currentColor = ref({
  hex: '#6366f1',
  rgb: 'rgb(99, 102, 241)',
  hsl: 'hsl(239, 84%, 67%)',
  r: 99, g: 102, b: 241
})

// 历史记录
const history = ref<Array<{ hex: string; rgb: string; hsl: string; r: number; g: number; b: number }>>([])

// 所有颜色格式
const colorFormats = computed(() => {
  const { r, g, b } = currentColor.value
  const hsl = rgbToHsl(r, g, b)
  return [
    { label: 'HEX', value: currentColor.value.hex },
    { label: 'HEX (大写)', value: currentColor.value.hex.toUpperCase() },
    { label: 'RGB', value: currentColor.value.rgb },
    { label: 'RGBA', value: `rgba(${r}, ${g}, ${b}, 1.0)` },
    { label: 'HSL', value: currentColor.value.hsl },
    { label: 'CSS 变量', value: `--color: ${currentColor.value.hex};` },
    { label: '整数 RGB', value: `${r}, ${g}, ${b}` },
    { label: 'RGB 二进制', value: `#${toHex(r)}${toHex(g)}${toHex(b)}` },
  ]
})

// 开始取色（使用浏览器原生 EyeDropper API 或回退到输入）
async function startPick() {
  // 尝试使用浏览器原生 API
  if ('EyeDropper' in window) {
    try {
      isPicking.value = true
      // @ts-ignore - EyeDropper 是实验性 API
      const eyeDropper = new window.EyeDropper()
      const result = await eyeDropper.open()
      updateColor(result.sRGBHex)
    } catch {
      // 用户取消或出错
    } finally {
      isPicking.value = false
    }
  } else {
    // 回退：使用颜色选择器
    const input = document.createElement('input')
    input.type = 'color'
    input.value = currentColor.value.hex
    input.onchange = () => updateColor(input.value)
    input.click()
  }
}

// 更新当前颜色
function updateColor(hex: string) {
  const { r, g, b } = hexToRgb(hex)
  const { h, s, l } = rgbToHsl(r, g, b)

  currentColor.value = {
    hex: hex.toLowerCase(),
    rgb: `rgb(${r}, ${g}, ${b})`,
    hsl: `hsl(${h}, ${s}%, ${l}%)`,
    r, g, b
  }

  // 添加到历史记录（避免重复）
  const exists = history.value.some(h => h.hex === hex)
  if (!exists) {
    history.value.unshift({ ...currentColor.value })
    if (history.value.length > 32) history.value.pop()
    saveHistory()
  }
}

// 选择历史颜色
function selectHistory(item: typeof currentColor.value) {
  currentColor.value = item
}

// 清空历史
function clearHistory() {
  history.value = []
  localStorage.removeItem('colorPickerHistory')
  appStore.showToast('success', '历史记录已清空')
}

// 复制颜色值
async function copyColor(value: string) {
  await navigator.clipboard.writeText(value)
  appStore.showToast('success', `已复制: ${value}`)
}

// 工具函数：HEX 转 RGB
function hexToRgb(hex: string): { r: number; g: number; b: number } {
  const result = /^#?([a-f\d]{2})([a-f\d]{2})([a-f\d]{2})$/i.exec(hex)
  return result
    ? { r: parseInt(result[1], 16), g: parseInt(result[2], 16), b: parseInt(result[3], 16) }
    : { r: 0, g: 0, b: 0 }
}

// 工具函数：RGB 转 HSL
function rgbToHsl(r: number, g: number, b: number): { h: number; s: number; l: number } {
  r /= 255; g /= 255; b /= 255
  const max = Math.max(r, g, b), min = Math.min(r, g, b)
  let h = 0, s = 0
  const l = (max + min) / 2

  if (max !== min) {
    const d = max - min
    s = l > 0.5 ? d / (2 - max - min) : d / (max + min)
    switch (max) {
      case r: h = ((g - b) / d + (g < b ? 6 : 0)) / 6; break
      case g: h = ((b - r) / d + 2) / 6; break
      case b: h = ((r - g) / d + 4) / 6; break
    }
  }

  return { h: Math.round(h * 360), s: Math.round(s * 100), l: Math.round(l * 100) }
}

// 工具函数：数字转 HEX
function toHex(n: number): string {
  return n.toString(16).padStart(2, '0')
}

// 持久化历史记录
function saveHistory() {
  localStorage.setItem('colorPickerHistory', JSON.stringify(history.value))
}

// 加载历史记录
function loadHistory() {
  const saved = localStorage.getItem('colorPickerHistory')
  if (saved) {
    try {
      history.value = JSON.parse(saved)
    } catch {
      // 解析失败，忽略
    }
  }
}

onMounted(loadHistory)
</script>

<style scoped>
/* 颜色值展示 */
.color-value {
  font-family: "JetBrains Mono", monospace;
  font-size: 13px;
  background: var(--bg-hover);
  padding: 4px 10px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.color-value:hover {
  background: var(--accent);
  color: #fff;
}

/* 格式项 */
.format-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 8px 12px;
  background: var(--bg-hover);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
}

.format-item:hover {
  background: rgba(99, 102, 241, 0.15);
}

.format-item code {
  font-family: "JetBrains Mono", monospace;
  font-size: 12px;
  color: var(--accent);
}

/* 历史项 */
.history-item {
  aspect-ratio: 1;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.15s ease;
  border: 2px solid transparent;
}

.history-item:hover {
  transform: scale(1.1);
  border-color: var(--accent);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
}
</style>
