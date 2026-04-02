<template>
  <div class="page-container">
    <div>
      <div class="page-title"><Settings2 :size="20" class="text-primary-400"/>设置</div>
      <div class="page-desc">主题、外观、快捷键、关于</div>
    </div>

    <div class="flex-1 overflow-auto space-y-4">
      <!-- 主题设置 -->
      <div class="card">
        <div class="font-semibold mb-3 flex items-center gap-2"><Palette :size="16"/>主题外观</div>
        <div class="label mb-3">选择主题</div>
        <div class="grid grid-cols-4 gap-2 mb-4">
          <button
            v-for="t in themes" :key="t.id"
            @click="switchTheme(t.id)"
            :class="['theme-card', appStore.theme === t.id && 'active']"
          >
            <div class="theme-preview" :style="{ background: t.preview }" />
            <div class="text-xs mt-1 text-center">{{ t.name }}</div>
          </button>
        </div>

        <div class="divider" />
        <div class="flex items-center justify-between">
          <div>
            <div class="text-sm font-medium">字体大小</div>
            <div class="text-xs opacity-50">调整界面文字大小</div>
          </div>
          <div class="flex items-center gap-2">
            <button @click="decreaseFontSize" class="btn btn-secondary py-1 px-2"><Minus :size="12"/></button>
            <span class="w-10 text-center text-sm font-mono">{{ config.fontSize }}px</span>
            <button @click="increaseFontSize" class="btn btn-secondary py-1 px-2"><Plus :size="12"/></button>
          </div>
        </div>
      </div>

      <!-- 窗口行为 -->
      <div class="card">
        <div class="font-semibold mb-3 flex items-center gap-2"><AppWindow :size="16"/>窗口行为</div>
        <div class="space-y-3">
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm">最小化到系统托盘</div>
              <div class="text-xs opacity-50">关闭按钮时最小化而非退出</div>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="config.minToTray" @change="saveConfig" />
              <span class="toggle-track" />
            </label>
          </div>
          <div class="flex items-center justify-between">
            <div>
              <div class="text-sm">窗口置顶</div>
              <div class="text-xs opacity-50">保持在其他窗口上层</div>
            </div>
            <label class="toggle-switch">
              <input type="checkbox" v-model="config.alwaysOnTop" @change="saveConfig" />
              <span class="toggle-track" />
            </label>
          </div>
        </div>
      </div>

      <!-- 关于 -->
      <div class="card">
        <div class="font-semibold mb-3 flex items-center gap-2"><Info :size="16"/>关于</div>
        <div class="space-y-1 text-sm opacity-70">
          <div>版本：<span class="font-mono text-primary-400">v1.1.0</span></div>
          <div>作者：<span class="text-primary-400">MasterPick</span></div>
          <div>技术栈：<span class="font-mono">Wails v2 + Go + Vue3 + TypeScript</span></div>
          <div>开源：
            <a href="https://github.com/MasterPick/xtool" class="text-primary-400 underline" target="_blank">GitHub</a>
            · <a href="https://github.com/MasterPick/xtool/releases" class="text-primary-400 underline" target="_blank">Releases</a>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Settings2, Palette, AppWindow, Info, Minus, Plus } from 'lucide-vue-next'
import { useAppStore } from '@/stores/app'
import { GetAvailableThemes } from '../../../wailsjs/go/advanced/AdvancedTools'

const appStore = useAppStore()
const config = ref({ ...appStore.config })

interface ThemeOption { id: string; name: string; preview: string }
const themes = ref<ThemeOption[]>([])

async function switchTheme(id: string) {
  await appStore.switchTheme(id)
  config.value.theme = id
  appStore.showToast('success', '主题已切换')
}

function saveConfig() {
  appStore.saveConfig(config.value)
}

function decreaseFontSize() {
  if (config.value.fontSize <= 12) return
  config.value.fontSize -= 1
  document.documentElement.style.fontSize = config.value.fontSize + 'px'
  saveConfig()
}

function increaseFontSize() {
  if (config.value.fontSize >= 18) return
  config.value.fontSize += 1
  document.documentElement.style.fontSize = config.value.fontSize + 'px'
  saveConfig()
}

onMounted(async () => {
  themes.value = await GetAvailableThemes() as ThemeOption[]
})
</script>

<style scoped>
.theme-card {
  padding: 8px 6px;
  border-radius: 10px;
  border: 2px solid transparent;
  background: var(--bg-hover);
  cursor: pointer;
  transition: all 0.2s ease;
}
.theme-card:hover { border-color: rgba(99,102,241,0.4); }
.theme-card.active { border-color: var(--accent); }
.theme-preview {
  height: 36px;
  border-radius: 6px;
}

.toggle-switch { position: relative; display: inline-block; cursor: pointer; }
.toggle-switch input { display: none; }
.toggle-track {
  display: block; width: 38px; height: 22px;
  background: rgba(255,255,255,0.1);
  border-radius: 11px;
  transition: background 0.2s;
  position: relative;
}
.toggle-track::after {
  content: '';
  position: absolute;
  top: 3px; left: 3px;
  width: 16px; height: 16px;
  border-radius: 50%;
  background: #fff;
  transition: transform 0.2s;
}
.toggle-switch input:checked + .toggle-track { background: var(--accent); }
.toggle-switch input:checked + .toggle-track::after { transform: translateX(16px); }
</style>
