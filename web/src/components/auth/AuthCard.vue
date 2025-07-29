<template>
  <div class="rounded-2xl border border-gray-200 bg-white shadow-lg dark:border-gray-800 dark:bg-white/[0.03]" :class="cardClass">
    <!-- 卡片头部 -->
    <div v-if="$slots.header || title" class="border-b border-gray-200 dark:border-gray-700 px-6 py-4">
      <slot name="header">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white/90">{{ title }}</h3>
            <p v-if="subtitle" class="text-sm text-gray-600 dark:text-gray-400 mt-1">{{ subtitle }}</p>
          </div>
          <slot name="headerActions"></slot>
        </div>
      </slot>
    </div>

    <!-- 卡片内容 -->
    <div class="p-6" :class="contentClass">
      <slot></slot>
    </div>

    <!-- 卡片底部 -->
    <div v-if="$slots.footer" class="border-t border-gray-200 dark:border-gray-700 px-6 py-4">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup lang="ts">
interface Props {
  title?: string
  subtitle?: string
  cardClass?: string
  contentClass?: string
}

withDefaults(defineProps<Props>(), {
  title: '',
  subtitle: '',
  cardClass: '',
  contentClass: ''
})
</script>
