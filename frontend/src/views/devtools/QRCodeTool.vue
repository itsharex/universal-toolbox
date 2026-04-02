<template>
  <!-- 二维码工具页面 -->
  <div class="page-container">
    <!-- 页面标题 -->
    <div>
      <div class="page-title">
        <QrCode :size="20" class="text-primary-400" />
        二维码工具
      </div>
      <div class="page-desc">生成二维码 · 自定义样式 · 下载导出</div>
    </div>

    <div class="flex-1 flex gap-4 min-h-0">
      <!-- 左侧：输入配置 -->
      <div class="w-80 flex flex-col gap-3 shrink-0">
        <!-- 内容输入 -->
        <div class="card">
          <div class="font-semibold mb-3">输入内容</div>
          <textarea
            v-model="text"
            class="textarea-field mb-3"
            rows="3"
            placeholder="输入文字、URL、电话、邮箱..."
          />
          <button @click="generate" :disabled="!text.trim()" class="btn btn-primary w-full">
            <Wand2 :size="14" />
            生成二维码
          </button>
        </div>

        <!-- 样式配置 -->
        <div class="card">
          <div class="font-semibold mb-3">样式配置</div>
          <div class="space-y-3">
            <!-- 尺寸 -->
            <div>
              <div class="label mb-1">尺寸</div>
              <div class="flex gap-2">
                <button
                  v-for="s in [128, 256, 384, 512]"
                  :key="s"
                  @click="size = s"
                  :class="['btn text-xs', size === s ? 'btn-primary' : 'btn-secondary']"
                >
                  {{ s }}px
                </button>
              </div>
            </div>

            <!-- 前景色 -->
            <div>
              <div class="label mb-1">前景色</div>
              <div class="flex gap-2 items-center">
                <input type="color" v-model="darkColor" class="w-10 h-8 rounded cursor-pointer" />
                <input v-model="darkColor" class="input-field text-xs font-mono flex-1" />
              </div>
            </div>

            <!-- 背景色 -->
            <div>
              <div class="label mb-1">背景色</div>
              <div class="flex gap-2 items-center">
                <input type="color" v-model="lightColor" class="w-10 h-8 rounded cursor-pointer" />
                <input v-model="lightColor" class="input-field text-xs font-mono flex-1" />
              </div>
            </div>

            <!-- 容错级别 -->
            <div>
              <div class="label mb-1">容错级别</div>
              <select v-model="errorLevel" class="input-field">
                <option value="L">L - 7%</option>
                <option value="M">M - 15%</option>
                <option value="Q">Q - 25%</option>
                <option value="H">H - 30%</option>
              </select>
            </div>
          </div>
        </div>

        <!-- 快捷模板 -->
        <div class="card">
          <div class="font-semibold mb-2">快捷模板</div>
          <div class="flex flex-wrap gap-2">
            <button @click="text = 'https://'" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">URL</button>
            <button @click="text = 'tel:'" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">电话</button>
            <button @click="text = 'mailto:'" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">邮箱</button>
            <button @click="text = 'WIFI:T:WPA;S:;P:;;'" class="text-xs px-2 py-1 rounded bg-white/5 hover:bg-white/10">WiFi</button>
          </div>
        </div>
      </div>

      <!-- 右侧：预览和操作 -->
      <div class="flex-1 flex flex-col min-h-0">
        <!-- 预览区 -->
        <div class="flex-1 flex items-center justify-center card">
          <div v-if="!generated" class="text-center opacity-30">
            <QrCode :size="48" class="mx-auto mb-2" />
            <div class="text-sm">输入内容后点击生成</div>
          </div>
          <div v-else class="text-center">
            <canvas ref="canvasRef" class="rounded-xl shadow-2xl"></canvas>
            <div class="text-xs opacity-50 mt-2">{{ size }} × {{ size }} px</div>
          </div>
        </div>

        <!-- 操作按钮 -->
        <div v-if="generated" class="flex gap-2 mt-3">
          <button @click="copyImage" class="btn btn-secondary flex-1">
            <Copy :size="14" />
            复制图片
          </button>
          <button @click="downloadPng" class="btn btn-secondary flex-1">
            <Download :size="14" />
            下载 PNG
          </button>
          <button @click="downloadSvg" class="btn btn-primary flex-1">
            <FileImage :size="14" />
            下载 SVG
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { QrCode, Wand2, Copy, Download, FileImage } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 状态
const text = ref('')
const size = ref(256)
const darkColor = ref('#000000')
const lightColor = ref('#ffffff')
const errorLevel = ref<'L' | 'M' | 'Q' | 'H'>('M')
const generated = ref(false)
const canvasRef = ref<HTMLCanvasElement | null>(null)

let QRCodeLib: any = null

onMounted(async () => {
  try {
    const mod = await import('qrcode')
    QRCodeLib = mod.default || mod
  } catch {
    appStore.showToast('error', 'qrcode 库未安装')
  }
})

// 生成二维码
async function generate() {
  if (!text.value.trim() || !QRCodeLib || !canvasRef.value) return

  try {
    await QRCodeLib.toCanvas(canvasRef.value, text.value, {
      width: size.value,
      margin: 2,
      color: {
        dark: darkColor.value,
        light: lightColor.value
      },
      errorCorrectionLevel: errorLevel.value
    })
    generated.value = true
    appStore.showToast('success', '二维码生成成功')
  } catch (err) {
    appStore.showToast('error', '生成失败：' + String(err))
  }
}

// 复制图片
async function copyImage() {
  if (!canvasRef.value) return
  try {
    const blob = await new Promise<Blob>((resolve) => {
      canvasRef.value!.toBlob((b) => resolve(b!), 'image/png')
    })
    await navigator.clipboard.write([new ClipboardItem({ 'image/png': blob })])
    appStore.showToast('success', '已复制到剪贴板')
  } catch {
    appStore.showToast('error', '复制失败')
  }
}

// 下载 PNG
function downloadPng() {
  if (!canvasRef.value) return
  const link = document.createElement('a')
  link.download = 'qrcode.png'
  link.href = canvasRef.value.toDataURL('image/png')
  link.click()
  appStore.showToast('success', '已下载 PNG')
}

// 下载 SVG
async function downloadSvg() {
  if (!text.value.trim() || !QRCodeLib) return
  try {
    const svg = await QRCodeLib.toString(text.value, {
      type: 'svg',
      width: size.value,
      margin: 2,
      color: {
        dark: darkColor.value,
        light: lightColor.value
      }
    })
    const blob = new Blob([svg], { type: 'image/svg+xml' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.download = 'qrcode.svg'
    link.href = url
    link.click()
    URL.revokeObjectURL(url)
    appStore.showToast('success', '已下载 SVG')
  } catch (err) {
    appStore.showToast('error', '导出失败：' + String(err))
  }
}
</script>

<style scoped>
input[type="color"] {
  border: none;
  padding: 0;
  background: transparent;
}

input[type="color"]::-webkit-color-swatch-wrapper {
  padding: 0;
}

input[type="color"]::-webkit-color-swatch {
  border-radius: 4px;
  border: 1px solid rgba(255, 255, 255, 0.1);
}
</style>
