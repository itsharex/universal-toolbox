<template>
  <ToolPage title="计算器" description="标准计算与科学计算">
    <div class="tab-bar mb-4">
      <button :class="['tab-item', mode==='standard'&&'active']" @click="mode='standard'">标准</button>
      <button :class="['tab-item', mode==='scientific'&&'active']" @click="mode='scientific'">科学</button>
    </div>
    <div class="calc-screen card mb-4">
      <div class="text-xs opacity-40 h-4">{{ expression }}</div>
      <div class="text-3xl font-mono font-light mt-1 truncate">{{ display }}</div>
    </div>
    <div class="calc-grid">
      <template v-if="mode==='scientific'">
        <button v-for="b in sciButtons" :key="b" @click="pressSci(b)" class="calc-btn calc-btn-fn">{{ b }}</button>
      </template>
      <button @click="clear"           class="calc-btn calc-btn-clear col-span-2">AC</button>
      <button @click="backspace"       class="calc-btn calc-btn-fn">DEL</button>
      <button @click="press('/')"      class="calc-btn calc-btn-op">/</button>
      <button v-for="n in [7,8,9]"  :key="n" @click="press(String(n))" class="calc-btn">{{ n }}</button>
      <button @click="press('*')"      class="calc-btn calc-btn-op">x</button>
      <button v-for="n in [4,5,6]"  :key="n" @click="press(String(n))" class="calc-btn">{{ n }}</button>
      <button @click="press('-')"      class="calc-btn calc-btn-op">-</button>
      <button v-for="n in [1,2,3]"  :key="n" @click="press(String(n))" class="calc-btn">{{ n }}</button>
      <button @click="press('+')"      class="calc-btn calc-btn-op">+</button>
      <button @click="press('0')"      class="calc-btn col-span-2">0</button>
      <button @click="press('.')"      class="calc-btn">.</button>
      <button @click="calculate()"       class="calc-btn calc-btn-eq">=</button>
    </div>

    <!-- 计算历史 -->
    <div v-if="history.length" class="mt-4">
      <div class="flex items-center justify-between mb-2">
        <div class="text-sm opacity-60">计算历史</div>
        <button @click="history=[]" class="text-xs opacity-50 hover:opacity-100">清空</button>
      </div>
      <div class="space-y-1 max-h-40 overflow-auto">
        <div v-for="(h,i) in history" :key="i" class="flex justify-between text-xs opacity-50 hover:opacity-80 cursor-pointer" @click="useHistory(h)">
          <span class="font-mono">{{ h.expr }}</span>
          <span class="font-mono text-primary-400">= {{ h.result }}</span>
        </div>
      </div>
    </div>
  </ToolPage>
</template>
<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import ToolPage from '@/components/ToolPage.vue'
import { CalcBasic, CalcScientific } from '../../../wailsjs/go/daily/DailyTools'

const mode = ref<'standard'|'scientific'>('standard')
const display = ref('0'), expression = ref(''), current = ref('')
const history = ref<Array<{expr:string, result:string}>>([])
const sciButtons = ['sin','cos','tan','sqrt','x^2','ln','log','pi']

// 运算符映射
const opMap: Record<string, string> = {
  '+': '+', '-': '-', '*': '*', '/': '/'
}

function press(val: string) {
  if (['+','-','*','/','%'].includes(val)) {
    // 如果当前有表达式且有值，先计算之前的
    if (current.value && expression.value) {
      calculate(true)
    }
    expression.value = display.value + ' ' + val + ' '
    current.value = ''
  } else {
    if (current.value === '0' && val !== '.') {
      current.value = val
    } else {
      current.value += val
    }
    display.value = current.value
  }
}

async function calculate(chained: boolean = false) {
  try {
    const leftStr = expression.value.replace(/[+\-*/]\s*$/, '').trim()
    const rightStr = display.value
    const op = expression.value.match(/([+\-*/])\s*$/)?.[1]

    if (!leftStr || !rightStr) return

    const a = parseFloat(leftStr)
    const b = parseFloat(rightStr)

    if (isNaN(a) || isNaN(b)) {
      display.value = 'Error'
      return
    }

    const exprStr = `${leftStr} ${op} ${rightStr}`

    let result: number
    if (op === '%') {
      result = a % b
    } else {
      result = await CalcBasic(a, b, opMap[op || '+'] || '+')
    }

    const resultStr = String(parseFloat(result.toFixed(10)))
    display.value = resultStr

    if (!chained) {
      expression.value = exprStr + ' ='
      history.value.unshift({ expr: exprStr, result: resultStr })
      if (history.value.length > 50) history.value.pop()
    } else {
      expression.value = resultStr + ' '
    }
    current.value = resultStr
  } catch {
    display.value = 'Error'
  }
}

function pressSci(fn: string) {
  const v = parseFloat(display.value)
  const fnMap: Record<string, string> = {
    'sin': 'sin', 'cos': 'cos', 'tan': 'tan',
    'sqrt': 'sqrt', 'x^2': 'pow', 'ln': 'ln',
    'log': 'log', 'pi': 'pi'
  }
  const realFn = fnMap[fn] || fn
  const param = fn === 'x^2' ? 2 : fn === 'pi' ? 1 : 0

  CalcScientific(v, param, realFn).then((result: number) => {
    const resultStr = String(parseFloat(result.toFixed(10)))
    expression.value = `${fn}(${display.value}) =`
    display.value = resultStr
    current.value = resultStr
    history.value.unshift({ expr: `${fn}(${v})`, result: resultStr })
  }).catch(() => {
    display.value = 'Error'
  })
}

function backspace() {
  if (current.value.length > 1) {
    current.value = current.value.slice(0, -1)
    display.value = current.value
  } else {
    current.value = '0'
    display.value = '0'
  }
}

function clear() { display.value='0'; expression.value=''; current.value='' }

function useHistory(h: {expr:string, result:string}) {
  display.value = h.result
  current.value = h.result
}

// 键盘快捷键支持
function handleKeyDown(e: KeyboardEvent) {
  // 数字键
  if (/^[0-9]$/.test(e.key)) {
    e.preventDefault()
    press(e.key)
    return
  }
  // 运算符
  if (['+','-','*','/','%'].includes(e.key)) {
    e.preventDefault()
    press(e.key)
    return
  }
  // 小数点
  if (e.key === '.') {
    e.preventDefault()
    press('.')
    return
  }
  // 等号或回车
  if (e.key === '=' || e.key === 'Enter') {
    e.preventDefault()
    calculate()
    return
  }
  // 退格键
  if (e.key === 'Backspace') {
    e.preventDefault()
    backspace()
    return
  }
  // Escape 清空
  if (e.key === 'Escape') {
    e.preventDefault()
    clear()
    return
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeyDown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeyDown)
})
</script>
<style scoped>
.calc-screen { padding: 12px 16px; }
.calc-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 8px; }
.calc-btn {
  background: var(--bg-hover); border: 1px solid var(--border-color);
  border-radius: 10px; padding: 14px 8px; font-size: 16px; cursor: pointer;
  color: var(--text-primary); transition: all 0.1s; text-align: center;
}
.calc-btn:hover  { background: rgba(255,255,255,0.1); }
.calc-btn:active { transform: scale(0.95); }
.calc-btn-op    { color: #818cf8; font-size: 18px; }
.calc-btn-eq    { background: #6366f1; color: #fff; }
.calc-btn-eq:hover { background: #818cf8; }
.calc-btn-clear { background: rgba(239,68,68,0.15); color: #f87171; }
.calc-btn-fn    { color: #94a3b8; font-size: 12px; }
</style>
