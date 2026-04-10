<template>
  <div class="tool-page">
    <!-- 面包屑导航（可选） -->
    <div class="breadcrumb" v-if="showBreadcrumb && breadcrumbs.length > 0">
      <router-link to="/" class="breadcrumb-item breadcrumb-home">
        <Home :size="12" />
        <span>首页</span>
      </router-link>
      <ChevronRight :size="12" class="breadcrumb-separator" v-for="(_, i) in breadcrumbs" :key="i" />
      <router-link
        v-for="(crumb, index) in breadcrumbs"
        :key="crumb.path"
        :to="crumb.path"
        class="breadcrumb-item"
        :class="{ active: index === breadcrumbs.length - 1 }"
      >
        {{ crumb.label }}
      </router-link>
    </div>

    <!-- 顶部标题区 -->
    <div class="tool-header" v-if="title">
      <div class="tool-header-left">
        <h2 class="tool-title">{{ title }}</h2>
        <p class="tool-desc" v-if="description">{{ description }}</p>
      </div>
      <div class="tool-header-right">
        <!-- 返回首页按钮（当不在首页时显示） -->
        <button
          v-if="!isHomePage"
          class="btn btn-ghost tool-back-btn"
          @click="goHome"
          title="返回首页"
        >
          <ArrowLeft :size="14" />
          <span class="hidden sm:inline">返回首页</span>
        </button>
        <!-- actions 插槽 -->
        <template v-if="$slots.actions">
          <slot name="actions" />
        </template>
      </div>
    </div>
    <!-- 内容区 -->
    <div class="tool-body">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { Home, ChevronRight, ArrowLeft } from 'lucide-vue-next'

const props = withDefaults(defineProps<{
  title?: string
  description?: string
  showBreadcrumb?: boolean
}>(), {
  showBreadcrumb: false,
})

const route = useRoute()
const router = useRouter()

// 是否在首页
const isHomePage = computed(() => route.path === '/')

// 面包屑数据
const breadcrumbs = computed(() => {
  const pathSegments = route.path.split('/').filter(Boolean)
  const crumbs: { label: string; path: string }[] = []

  // 分类映射
  const categoryMap: Record<string, string> = {
    devtools: '开发工具',
    sysinfo: '系统工具',
    daily: '日常工具',
    network: '网络工具',
    settings: '设置',
  }

  let currentPath = ''
  for (const segment of pathSegments) {
    currentPath += '/' + segment
    const label = categoryMap[segment] || (route.meta?.title as string) || segment
    crumbs.push({ label, path: currentPath })
  }

  return crumbs
})

// 返回首页
function goHome() {
  router.push('/')
}
</script>

<style scoped>
.tool-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 24px;
  overflow: hidden;
}

/* 面包屑导航 */
.breadcrumb {
  display: flex;
  align-items: center;
  gap: 4px;
  margin-bottom: 12px;
  flex-shrink: 0;
}

.breadcrumb-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 12px;
  color: var(--text-muted);
  text-decoration: none;
  padding: 2px 6px;
  border-radius: 4px;
  transition: all 0.15s ease;
  cursor: pointer;
}

.breadcrumb-item:hover {
  color: var(--text-secondary);
  background: var(--bg-hover);
}

.breadcrumb-item.active {
  color: var(--text-secondary);
  cursor: default;
}

.breadcrumb-item.active:hover {
  background: transparent;
}

.breadcrumb-home {
  opacity: 0.6;
}

.breadcrumb-separator {
  color: var(--text-muted);
  opacity: 0.4;
  flex-shrink: 0;
}

.tool-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 16px;
  flex-shrink: 0;
}

.tool-header-left {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.tool-header-right {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-shrink: 0;
}

.tool-title {
  font-size: 18px;
  font-weight: 600;
  color: var(--text-primary, #e2e8f0);
  margin: 0;
  line-height: 1.4;
}

.tool-desc {
  font-size: 13px;
  color: var(--text-muted, #94a3b8);
  margin: 0;
  line-height: 1.4;
}

/* 返回首页按钮 */
.tool-back-btn {
  font-size: 12px;
  padding: 5px 10px;
  gap: 4px;
  opacity: 0.7;
}

.tool-back-btn:hover {
  opacity: 1;
}

.tool-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
</style>
