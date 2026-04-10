<template>
  <ToolPage title="加密解密" description="AES/DES/RSA 加密解密">
    <!-- 算法选择 -->
    <div class="toolbar mb-4">
      <div class="tab-bar">
        <button
          v-for="algo in algorithms"
          :key="algo.id"
          @click="currentAlgo = algo.id"
          :class="['tab-item', currentAlgo === algo.id && 'active']"
        >
          {{ algo.name }}
        </button>
      </div>
    </div>

    <!-- 密钥输入（对称加密需要） -->
    <div v-if="needsKey" class="card mb-4">
      <div class="flex items-center gap-3">
        <div class="flex-1">
          <label class="label">密钥</label>
          <input
            v-model="secretKey"
            type="password"
            class="input-field"
            placeholder="请输入密钥（8-32位）"
          />
        </div>
        <div v-if="currentAlgo === 'base64'" class="w-20">
          <label class="label">模式</label>
          <select v-model="base64Mode" class="input-field">
            <option value="encode">编码</option>
            <option value="decode">解码</option>
          </select>
        </div>
      </div>
    </div>

    <!-- 输入输出区 -->
    <div class="flex-1 flex flex-col gap-3 min-h-0">
      <!-- 输入 -->
      <div class="flex flex-col gap-2 flex-1 min-h-0">
        <div class="label">
          <span>{{ currentAlgo === 'md5' ? '原文（MD5 只支持单行）' : '输入内容' }}</span>
          <button v-if="inputText" @click="inputText = ''" class="text-xs opacity-50 hover:opacity-100">
            清空
          </button>
        </div>
        <textarea
          v-model="inputText"
          class="textarea-field flex-1 min-h-0 font-mono text-sm"
          :placeholder="inputPlaceholder"
          spellcheck="false"
        />
      </div>

      <!-- 操作按钮 -->
      <div class="flex justify-center gap-3 py-2">
        <button @click="encrypt" :disabled="!canEncrypt" class="btn btn-primary px-8">
          <Lock :size="14" />
          {{ currentAlgo === 'base64' && base64Mode === 'decode' ? '解码' : '加密' }}
        </button>
        <button @click="decrypt" :disabled="!canDecrypt" class="btn btn-secondary px-8">
          <Unlock :size="14" />
          解密
        </button>
        <button @click="swapInputOutput" :disabled="!outputText" class="btn btn-secondary px-4">
          <ArrowUpDown :size="14" />
          交换
        </button>
      </div>

      <!-- 输出 -->
      <div class="flex flex-col gap-2 flex-1 min-h-0">
        <div class="label">
          <span>输出结果</span>
          <div class="flex items-center gap-2" v-if="outputText">
            <button @click="copyOutput" class="btn btn-secondary text-xs py-1 px-2">
              <Copy :size="12" />
              复制
            </button>
          </div>
        </div>
        <div class="code-output flex-1 min-h-0 overflow-auto whitespace-pre-wrap">
          <span v-if="!outputText" class="opacity-30">结果将显示在这里...</span>
          <span v-else :class="{ 'text-green-400': isSuccess, 'text-red-400': !isSuccess }">
            {{ outputText }}
          </span>
        </div>
      </div>
    </div>

    <!-- 说明 -->
    <div class="mt-4 p-3 bg-blue-500/10 border border-blue-500/20 rounded-lg">
      <div class="text-xs opacity-70">
        <strong class="text-blue-400">说明：</strong>
        <template v-if="currentAlgo === 'md5'">
          MD5 是单向哈希算法，不可逆。用于密码存储、文件校验等。
        </template>
        <template v-else-if="currentAlgo === 'base64'">
          Base64 是编码方式，非加密。可逆，用于数据传输、URL安全字符等。
        </template>
        <template v-else>
          对称加密算法，使用相同密钥加密和解密。请妥善保管密钥，丢失将无法解密。
        </template>
      </div>
    </div>
  </div>
  </ToolPage>
</template>

<script setup lang="ts">
import ToolPage from '@/components/ToolPage.vue'
import { ref, computed } from 'vue'
import { Shield, Lock, Unlock, ArrowUpDown, Copy } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import CryptoJS from 'crypto-js'

const appStore = useAppStore()

// 算法列表
const algorithms = [
  { id: 'aes', name: 'AES' },
  { id: 'des', name: 'DES' },
  { id: 'md5', name: 'MD5' },
  { id: 'sha256', name: 'SHA256' },
  { id: 'base64', name: 'Base64' },
]

// 状态
const currentAlgo = ref('aes')
const secretKey = ref('')
const inputText = ref('')
const outputText = ref('')
const isSuccess = ref(true)
const base64Mode = ref<'encode' | 'decode'>('encode')

// 需要密钥的算法
const needsKey = computed(() => ['aes', 'des'].includes(currentAlgo.value))

// 输入提示
const inputPlaceholder = computed(() => {
  if (currentAlgo.value === 'md5') return '请输入要哈希的内容（支持批量，每行一个）'
  if (currentAlgo.value === 'base64') return base64Mode.value === 'encode' ? '请输入要编码的内容' : '请输入要解码的 Base64 字符串'
  return '请输入要加密的内容'
})

// 可以加密
const canEncrypt = computed(() => {
  if (currentAlgo.value === 'md5' || currentAlgo.value === 'sha256') return !!inputText.value.trim()
  if (currentAlgo.value === 'base64') return !!inputText.value
  return !!inputText.value && !!secretKey.value
})

// 可以解密
const canDecrypt = computed(() => {
  if (currentAlgo.value === 'md5' || currentAlgo.value === 'sha256') return false
  if (currentAlgo.value === 'base64') return !!inputText.value
  return !!inputText.value && !!secretKey.value
})

// 加密
function encrypt() {
  outputText.value = ''
  isSuccess.value = true

  try {
    const text = inputText.value
    const key = secretKey.value

    switch (currentAlgo.value) {
      case 'aes':
        outputText.value = CryptoJS.AES.encrypt(text, key).toString()
        break
      case 'des':
        outputText.value = CryptoJS.DES.encrypt(text, key).toString()
        break
      case 'md5':
        outputText.value = text.split('\n').map(t => t.trim()).filter(Boolean).map(t => CryptoJS.MD5(t).toString()).join('\n')
        break
      case 'sha256':
        outputText.value = text.split('\n').map(t => t.trim()).filter(Boolean).map(t => CryptoJS.SHA256(t).toString()).join('\n')
        break
      case 'base64':
        outputText.value = btoa(unescape(encodeURIComponent(text)))
        break
    }
    appStore.showToast('success', '加密成功')
  } catch (err) {
    isSuccess.value = false
    outputText.value = `加密失败: ${err}`
    appStore.showToast('error', '加密失败')
  }
}

// 解密
function decrypt() {
  outputText.value = ''
  isSuccess.value = true

  try {
    const text = inputText.value
    const key = secretKey.value

    switch (currentAlgo.value) {
      case 'aes':
        const decrypted = CryptoJS.AES.decrypt(text, key)
        outputText.value = decrypted.toString(CryptoJS.enc.Utf8)
        if (!outputText.value) throw new Error('密钥错误或数据损坏')
        break
      case 'des':
        const desDecrypted = CryptoJS.DES.decrypt(text, key)
        outputText.value = desDecrypted.toString(CryptoJS.enc.Utf8)
        if (!outputText.value) throw new Error('密钥错误或数据损坏')
        break
      case 'base64':
        outputText.value = decodeURIComponent(escape(atob(text)))
        break
    }
    appStore.showToast('success', '解密成功')
  } catch (err) {
    isSuccess.value = false
    outputText.value = `解密失败: ${err}`
    appStore.showToast('error', '解密失败')
  }
}

// 交换输入输出
function swapInputOutput() {
  if (outputText.value) {
    inputText.value = outputText.value
    outputText.value = ''
  }
}

// 复制输出
async function copyOutput() {
  await navigator.clipboard.writeText(outputText.value)
  appStore.showToast('success', '已复制到剪贴板')
}
</script>
