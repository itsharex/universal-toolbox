<template>
  <div class="page-container">
    <!-- 标题 -->
    <div class="mb-4">
      <div class="page-title">
        <ArrowLeftRight :size="20" class="text-primary-400" />
        单位换算
      </div>
      <div class="page-desc">长度 · 重量 · 温度 · 速度 · 数据 · 面积 · 体积</div>
    </div>

    <!-- 类型选择 -->
    <div class="tab-bar mb-4">
      <button
        v-for="t in types"
        :key="t.id"
        :class="['tab-item', type === t.id && 'active']"
        @click="switchType(t.id)"
      >
        {{ t.label }}
      </button>
    </div>

    <!-- 输入区 -->
    <div class="card mb-4">
      <div class="flex gap-3 items-end flex-wrap">
        <div class="flex-1 min-w-[120px]">
          <div class="label mb-1">数值</div>
          <input
            v-model.number="value"
            type="number"
            class="input-field"
            placeholder="输入数值..."
            @keyup.enter="convert"
          />
        </div>
        <div class="min-w-[140px]">
          <div class="label mb-1">从</div>
          <select v-model="fromUnit" class="input-field">
            <option v-for="u in currentUnits" :key="u.id" :value="u.id">{{ u.label }}</option>
          </select>
        </div>
        <button @click="swap" class="btn btn-secondary p-2" title="交换">
          <ArrowLeftRight :size="14" />
        </button>
        <div class="min-w-[140px]">
          <div class="label mb-1">到</div>
          <select v-model="toUnit" class="input-field">
            <option v-for="u in currentUnits" :key="u.id" :value="u.id">{{ u.label }}</option>
          </select>
        </div>
        <button @click="convert" class="btn btn-primary">
          <ArrowLeftRight :size="14" />
          转换
        </button>
      </div>
    </div>

    <!-- 结果 -->
    <div v-if="result" class="card mb-4 text-center">
      <div class="text-3xl font-bold text-primary-400 mb-1">{{ result }}</div>
      <div class="text-xs opacity-50">{{ formula }}</div>
    </div>

    <!-- 快捷转换 -->
    <div class="card mb-4">
      <div class="label mb-2">快捷换算</div>
      <div class="grid grid-cols-2 md:grid-cols-3 gap-2">
        <button
          v-for="q in quickConversions"
          :key="q.label"
          @click="applyQuick(q)"
          class="text-xs p-2 rounded bg-white/5 hover:bg-white/10 text-left"
        >
          <div class="opacity-70">{{ q.label }}</div>
          <div class="font-mono text-primary-400">{{ q.desc }}</div>
        </button>
      </div>
    </div>

    <!-- 换算表 -->
    <div v-if="value" class="card">
      <div class="label mb-2">{{ value }} {{ fromUnitLabel }} =</div>
      <div class="space-y-1 max-h-48 overflow-auto">
        <div
          v-for="u in currentUnits"
          :key="u.id"
          class="flex justify-between text-sm py-1 border-b border-white/5"
        >
          <span class="opacity-70">{{ u.label }}</span>
          <span class="font-mono" :class="u.id === toUnit ? 'text-primary-400' : ''">
            {{ convertTo(u.id) }}
          </span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { ArrowLeftRight } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'

const appStore = useAppStore()

// 状态
const type = ref('length')
const value = ref(1)
const fromUnit = ref('m')
const toUnit = ref('km')
const result = ref('')
const formula = ref('')

// 类型定义
const types = [
  { id: 'length', label: '长度' },
  { id: 'weight', label: '重量' },
  { id: 'temperature', label: '温度' },
  { id: 'speed', label: '速度' },
  { id: 'data', label: '数据' },
  { id: 'area', label: '面积' },
  { id: 'volume', label: '体积' },
  { id: 'time', label: '时间' }
]

// 单位定义（带基准换算）
const unitDefs: Record<string, { id: string; label: string; toBase: number }> = {
  length: [
    { id: 'mm', label: '毫米', toBase: 0.001 },
    { id: 'cm', label: '厘米', toBase: 0.01 },
    { id: 'm', label: '米', toBase: 1 },
    { id: 'km', label: '千米', toBase: 1000 },
    { id: 'inch', label: '英寸', toBase: 0.0254 },
    { id: 'foot', label: '英尺', toBase: 0.3048 },
    { id: 'mile', label: '英里', toBase: 1609.344 },
    { id: 'nautical', label: '海里', toBase: 1852 }
  ],
  weight: [
    { id: 'mg', label: '毫克', toBase: 0.001 },
    { id: 'g', label: '克', toBase: 1 },
    { id: 'kg', label: '千克', toBase: 1000 },
    { id: 't', label: '吨', toBase: 1000000 },
    { id: 'oz', label: '盎司', toBase: 28.3495 },
    { id: 'lb', label: '磅', toBase: 453.592 },
    { id: 'jin', label: '斤', toBase: 500 }
  ],
  temperature: [
    { id: 'C', label: '摄氏度 °C', toBase: 1 },
    { id: 'F', label: '华氏度 °F', toBase: 1 },
    { id: 'K', label: '开尔文 K', toBase: 1 }
  ],
  speed: [
    { id: 'ms', label: '米/秒', toBase: 1 },
    { id: 'kmh', label: '千米/时', toBase: 0.277778 },
    { id: 'mph', label: '英里/时', toBase: 0.44704 },
    { id: 'knot', label: '节', toBase: 0.514444 },
    { id: 'mach', label: '马赫', toBase: 340.29 }
  ],
  data: [
    { id: 'B', label: '字节 B', toBase: 1 },
    { id: 'KB', label: '千字节 KB', toBase: 1024 },
    { id: 'MB', label: '兆字节 MB', toBase: 1048576 },
    { id: 'GB', label: '吉字节 GB', toBase: 1073741824 },
    { id: 'TB', label: '太字节 TB', toBase: 1099511627776 },
    { id: 'bit', label: '比特 bit', toBase: 0.125 }
  ],
  area: [
    { id: 'mm2', label: '平方毫米', toBase: 0.000001 },
    { id: 'cm2', label: '平方厘米', toBase: 0.0001 },
    { id: 'm2', label: '平方米', toBase: 1 },
    { id: 'km2', label: '平方千米', toBase: 1000000 },
    { id: 'ha', label: '公顷', toBase: 10000 },
    { id: 'acre', label: '英亩', toBase: 4046.86 },
    { id: 'mu', label: '亩', toBase: 666.667 }
  ],
  volume: [
    { id: 'ml', label: '毫升', toBase: 0.001 },
    { id: 'L', label: '升', toBase: 1 },
    { id: 'm3', label: '立方米', toBase: 1000 },
    { id: 'gal', label: '加仑(美)', toBase: 3.78541 },
    { id: 'oz_vol', label: '盎司(液)', toBase: 0.0295735 }
  ],
  time: [
    { id: 'ms', label: '毫秒', toBase: 0.001 },
    { id: 's', label: '秒', toBase: 1 },
    { id: 'min', label: '分钟', toBase: 60 },
    { id: 'h', label: '小时', toBase: 3600 },
    { id: 'd', label: '天', toBase: 86400 },
    { id: 'w', label: '周', toBase: 604800 },
    { id: 'mon', label: '月(30天)', toBase: 2592000 },
    { id: 'y', label: '年(365天)', toBase: 31536000 }
  ]
}

// 当前单位列表
const currentUnits = computed(() => unitDefs[type.value] || [])

// 源单位标签
const fromUnitLabel = computed(() => {
  const u = currentUnits.value.find(u => u.id === fromUnit.value)
  return u?.label || fromUnit.value
})

// 快捷转换
const quickConversions = computed(() => {
  const items: Array<{ label: string; desc: string; type: string; from: string; to: string; value: number }> = []
  switch (type.value) {
    case 'length':
      items.push(
        { label: '英寸到厘米', desc: '1 inch = 2.54 cm', type: 'length', from: 'inch', to: 'cm', value: 1 },
        { label: '英尺到米', desc: '1 foot = 0.3048 m', type: 'length', from: 'foot', to: 'm', value: 1 },
        { label: '英里到千米', desc: '1 mile = 1.609 km', type: 'length', from: 'mile', to: 'km', value: 1 }
      )
      break
    case 'weight':
      items.push(
        { label: '磅到千克', desc: '1 lb = 0.4536 kg', type: 'weight', from: 'lb', to: 'kg', value: 1 },
        { label: '盎司到克', desc: '1 oz = 28.35 g', type: 'weight', from: 'oz', to: 'g', value: 1 },
        { label: '斤到千克', desc: '1斤 = 0.5 kg', type: 'weight', from: 'jin', to: 'kg', value: 1 }
      )
      break
    case 'temperature':
      items.push(
        { label: '摄氏转华氏', desc: '°C → °F', type: 'temperature', from: 'C', to: 'F', value: 25 },
        { label: '华氏转摄氏', desc: '°F → °C', type: 'temperature', from: 'F', to: 'C', value: 100 },
        { label: '摄氏转开尔文', desc: '°C → K', type: 'temperature', from: 'C', to: 'K', value: 0 }
      )
      break
    case 'data':
      items.push(
        { label: 'GB到MB', desc: '1 GB = 1024 MB', type: 'data', from: 'GB', to: 'MB', value: 1 },
        { label: 'TB到GB', desc: '1 TB = 1024 GB', type: 'data', from: 'TB', to: 'GB', value: 1 },
        { label: 'MB到KB', desc: '1 MB = 1024 KB', type: 'data', from: 'MB', to: 'KB', value: 1 }
      )
      break
  }
  return items
})

// 切换类型时重置单位
function switchType(newType: string) {
  type.value = newType
  const units = currentUnits.value
  if (units.length >= 2) {
    fromUnit.value = units[0].id
    toUnit.value = units[1].id
  }
  result.value = ''
}

// 交换单位
function swap() {
  const temp = fromUnit.value
  fromUnit.value = toUnit.value
  toUnit.value = temp
  if (value.value) convert()
}

// 应用快捷转换
function applyQuick(q: typeof quickConversions.value[0]) {
  type.value = q.type
  fromUnit.value = q.from
  toUnit.value = q.to
  value.value = q.value
  convert()
}

// 纯前端转换
function doConvert(val: number, from: string, to: string): number {
  if (type.value === 'temperature') {
    return convertTemperature(val, from, to)
  }

  const fromDef = currentUnits.value.find(u => u.id === from)
  const toDef = currentUnits.value.find(u => u.id === to)
  if (!fromDef || !toDef) return val

  const baseValue = val * fromDef.toBase
  return baseValue / toDef.toBase
}

// 温度特殊处理
function convertTemperature(val: number, from: string, to: string): number {
  // 先转成摄氏度
  let celsius = val
  if (from === 'F') celsius = (val - 32) * 5 / 9
  else if (from === 'K') celsius = val - 273.15

  // 再从摄氏度转成目标单位
  if (to === 'C') return celsius
  else if (to === 'F') return celsius * 9 / 5 + 32
  else if (to === 'K') return celsius + 273.15
  return celsius
}

// 执行转换
function convert() {
  if (!value.value && value.value !== 0) return

  try {
    const converted = doConvert(value.value, fromUnit.value, toUnit.value)
    const precision = Math.abs(converted) < 0.01 ? 6 : Math.abs(converted) < 1 ? 4 : 2
    result.value = formatNumber(converted, precision)
    formula.value = `${value.value} ${fromUnitLabel.value} = ${result.value}`
    appStore.showToast('success', '转换完成')
  } catch (e) {
    appStore.showToast('error', String(e))
  }
}

// 转换到指定单位（用于换算表）
function convertTo(targetUnit: string): string {
  const converted = doConvert(value.value, fromUnit.value, targetUnit)
  return formatNumber(converted, 4)
}

// 格式化数字
function formatNumber(num: number, precision: number = 2): string {
  if (num === 0) return '0'
  if (Math.abs(num) >= 1000000 || Math.abs(num) < 0.0001) {
    return num.toExponential(4)
  }
  return num.toLocaleString('zh-CN', { maximumFractionDigits: precision })
}

// 初始化
watch(type, () => {
  const units = currentUnits.value
  if (units.length >= 2) {
    fromUnit.value = units[0].id
    toUnit.value = units[1].id
  }
}, { immediate: true })
</script>
