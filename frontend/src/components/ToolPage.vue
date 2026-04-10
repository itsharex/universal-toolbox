<template>
  <div class="tool-page">
    <!-- 顶部标题区 -->
    <div class="tool-header" v-if="title">
      <div class="tool-header-left">
        <h2 class="tool-title">{{ title }}</h2>
        <p class="tool-desc" v-if="description">{{ description }}</p>
      </div>
      <div class="tool-header-right" v-if="$slots.actions">
        <slot name="actions" />
      </div>
    </div>
    <!-- 内容区 -->
    <div class="tool-body">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  title?: string
  description?: string
}>()
</script>

<style scoped>
.tool-page {
  display: flex;
  flex-direction: column;
  height: 100%;
  padding: 20px 24px;
  overflow: hidden;
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

.tool-body {
  flex: 1;
  overflow: hidden;
  display: flex;
  flex-direction: column;
}
</style>
