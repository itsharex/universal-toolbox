<template>
  <ToolPage title="HTTP 测试" description="发送 HTTP 请求，查看响应结果">
    <!-- 请求配置 -->
    <div class="card mb-3">
      <div class="flex gap-2 mb-3">
        <!-- 方法选择 -->
        <select v-model="method" class="input-field w-28">
          <option v-for="m in methods" :key="m" :value="m">{{ m }}</option>
        </select>
        <!-- URL 输入 -->
        <input v-model="url" class="input-field flex-1"
          placeholder="输入请求 URL，如 https://api.example.com/v1/users"
          @keyup.enter="sendRequest" />
        <!-- 发送按钮 -->
        <button @click="sendRequest" class="btn btn-primary" :disabled="loading">
          <Send :size="14" :class="loading ? 'loading-spin' : ''"/>
          {{ loading ? '发送中...' : '发送' }}
        </button>
      </div>

      <!-- 自定义 Headers -->
      <div class="mb-3">
        <div class="flex items-center justify-between mb-2">
          <div class="label mb-0">请求头 (Headers)</div>
          <button @click="addHeader" class="btn btn-secondary py-0.5 px-2 text-xs"><Plus :size="11"/>添加</button>
        </div>
        <div class="space-y-1">
          <div v-for="(h, idx) in headers" :key="idx" class="flex gap-2 items-center">
            <input v-model="h.key" class="input-field flex-1" placeholder="Header 名称"/>
            <input v-model="h.value" class="input-field flex-1" placeholder="Header 值"/>
            <button @click="removeHeader(idx)" class="btn btn-danger py-0.5 px-2 text-xs"><X :size="11"/></button>
          </div>
          <div v-if="headers.length === 0" class="text-xs opacity-30">暂无自定义请求头</div>
        </div>
      </div>

      <!-- 请求体（POST/PUT/PATCH 显示） -->
      <div v-if="['POST','PUT','PATCH'].includes(method)">
        <div class="label">请求体 (JSON)</div>
        <textarea v-model="body" class="textarea-field" rows="3"
          placeholder='{"key": "value"}' spellcheck="false" />
      </div>
    </div>

    <!-- 响应结果 -->
    <div v-if="result" class="flex-1 overflow-auto">
      <div class="flex items-center gap-3 mb-3">
        <span :class="['badge text-sm', result.statusCode < 300 ? 'badge-success' : result.statusCode < 400 ? 'badge-warning' : 'badge-error']">
          {{ result.statusCode }} {{ result.status }}
        </span>
        <span class="text-xs opacity-50">{{ result.latencyMs.toFixed(0) }} ms</span>
        <span class="text-xs opacity-50">{{ result.contentType }}</span>
      </div>

      <div class="tab-bar mb-3">
        <button :class="['tab-item', tab==='body' && 'active']" @click="tab='body'">响应体</button>
        <button :class="['tab-item', tab==='headers' && 'active']" @click="tab='headers'">响应头</button>
      </div>

      <div v-if="tab==='body'" class="code-output overflow-auto" style="max-height: 300px">
        {{ formattedBody }}
      </div>
      <div v-if="tab==='headers'" class="code-output overflow-auto" style="max-height: 300px">
        <div v-for="(v, k) in result.headers" :key="k" class="flex gap-3">
          <span class="text-primary-400 shrink-0">{{ k }}:</span>
          <span class="opacity-80 break-all">{{ v }}</span>
        </div>
      </div>
    </div>

    <div v-if="error" class="card border-red-500/30 text-red-400 mt-3">{{ error }}</div>
  </ToolPage>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { Send, Plus, X } from 'lucide-vue-next'
import ToolPage from '@/components/ToolPage.vue'
import { useAppStore } from '@/stores/app'
import { HTTPRequest } from '../../../wailsjs/go/network/NetworkTools'

const appStore = useAppStore()

interface HttpResult {
  statusCode: number; status: string; latencyMs: number
  body: string; headers: Record<string, string>; contentType: string; error: string
}

const methods = ['GET', 'POST', 'PUT', 'DELETE', 'PATCH', 'HEAD', 'OPTIONS']
const method  = ref('GET')
const url     = ref('')
const body    = ref('')
const loading = ref(false)
const result  = ref<HttpResult | null>(null)
const error   = ref('')
const tab     = ref<'body' | 'headers'>('body')
const headers = ref<Array<{key: string; value: string}>>([])

const formattedBody = computed(() => {
  if (!result.value?.body) return ''
  try {
    return JSON.stringify(JSON.parse(result.value.body), null, 2)
  } catch {
    return result.value.body
  }
})

function addHeader() {
  headers.value.push({ key: '', value: '' })
}

function removeHeader(idx: number) {
  headers.value.splice(idx, 1)
}

// 构建请求头对象
function buildHeaders(): Record<string, string> {
  const h: Record<string, string> = {}
  for (const header of headers.value) {
    if (header.key.trim()) {
      h[header.key.trim()] = header.value
    }
  }
  return h
}

async function sendRequest() {
  if (!url.value.trim()) return
  loading.value = true
  error.value = ''
  result.value = null
  try {
    const headersObj = buildHeaders()
    const res = await HTTPRequest(method.value, url.value, body.value, headersObj) as HttpResult
    if (res.error) {
      error.value = res.error
    } else {
      result.value = res
    }
  } catch (e) {
    error.value = String(e)
    appStore.showToast('error', '请求失败: ' + String(e))
  } finally {
    loading.value = false
  }
}
</script>
