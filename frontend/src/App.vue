<template>
  <!-- 根组件：应用主容器 -->
  <div
    :class="[
      'app-root h-screen w-screen overflow-hidden flex flex-col',
      themeClass,
      { 'mica-effect': isMicaTheme }
    ]"
  >
    <!-- 自定义标题栏（--wails-draggable 允许拖动） -->
    <div
      class="titlebar flex items-center justify-between px-4 h-10 shrink-0 select-none"
      style="--wails-draggable: drag"
    >
      <!-- 左侧：Logo + 标题 -->
      <div class="flex items-center gap-2" style="--wails-draggable: no-drag">
        <div class="w-5 h-5 rounded bg-primary-500 flex items-center justify-center">
          <span class="text-white text-xs font-bold">X</span>
        </div>
        <span class="text-sm font-medium opacity-80">XTool</span>
      </div>

      <!-- 右侧：窗口控制按钮 -->
      <div class="flex items-center gap-1" style="--wails-draggable: no-drag">
        <!-- 最小化 -->
        <button @click="minimizeWindow" class="win-btn hover:bg-white/10">
          <Minus :size="12" />
        </button>
        <!-- 最大化/还原 -->
        <button @click="toggleMaximize" class="win-btn hover:bg-white/10">
          <Square :size="11" />
        </button>
        <!-- 关闭 -->
        <button @click="closeWindow" class="win-btn hover:bg-red-500/80">
          <X :size="12" />
        </button>
      </div>
    </div>

    <!-- 主内容区 -->
    <div class="flex flex-1 overflow-hidden">
      <!-- 左侧导航栏 -->
      <SideNav />

      <!-- 右侧内容区（路由视图） -->
      <main class="flex-1 overflow-hidden relative">
        <RouterView v-slot="{ Component }">
          <Transition name="page" mode="out-in">
            <component :is="Component" />
          </Transition>
        </RouterView>
      </main>
    </div>

    <!-- 全局通知 Toast -->
    <ToastContainer />

    <!-- 全局搜索弹窗 -->
    <SearchModal />
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { RouterView } from 'vue-router'
import { Minus, Square, X } from 'lucide-vue-next'
import SideNav from '@/components/SideNav.vue'
import ToastContainer from '@/components/ToastContainer.vue'
import SearchModal from '@/components/SearchModal.vue'
import { useAppStore } from '@/stores/app'
import { WindowMinimise, WindowToggleMaximise, Quit } from '../wailsjs/runtime/runtime'

const appStore = useAppStore()

// 根据主题计算类名
const themeClass = computed(() => {
  const theme = appStore.theme
  if (theme === 'dark' || theme === 'blue' || theme === 'purple' || theme === 'mica') return 'dark'
  if (theme === 'light') return 'light'
  // auto: 跟随系统
  return window.matchMedia('(prefers-color-scheme: dark)').matches ? 'dark' : 'light'
})

// 是否启用云母效果
const isMicaTheme = computed(() => appStore.theme === 'mica')

// 最小化窗口
const minimizeWindow = () => WindowMinimise()

// 切换最大化
const toggleMaximize = () => WindowToggleMaximise()

// 关闭窗口（退出应用）
const closeWindow = () => {
  if (appStore.config.minToTray) {
    WindowMinimise()
  } else {
    Quit()
  }
}

onMounted(async () => {
  // 初始化应用配置
  await appStore.loadConfig()
})
</script>

<style>
/* 全局样式 */
.app-root {
  font-family: "Microsoft YaHei", "PingFang SC", system-ui, sans-serif;
}

/* 标题栏样式 */
.titlebar {
  background: rgba(255, 255, 255, 0.03);
  border-bottom: 1px solid rgba(255, 255, 255, 0.06);
}

.dark .titlebar {
  color: #e2e8f0;
}

.light .titlebar {
  color: #1e293b;
  background: rgba(0, 0, 0, 0.02);
  border-bottom: 1px solid rgba(0, 0, 0, 0.08);
}

/* 窗口控制按钮 */
.win-btn {
  width: 28px;
  height: 24px;
  border-radius: 4px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: none;
  background: transparent;
  cursor: pointer;
  color: inherit;
  transition: background 0.15s ease;
}

/* 云母半透明效果 */
.mica-effect {
  background: transparent !important;
}

/* 路由页面切换动画 */
.page-enter-active,
.page-leave-active {
  transition: all 0.2s cubic-bezier(0.16, 1, 0.3, 1);
}

.page-enter-from {
  opacity: 0;
  transform: translateX(8px);
}

.page-leave-to {
  opacity: 0;
  transform: translateX(-8px);
}
</style>
