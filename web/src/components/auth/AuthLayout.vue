<template>
  <div class="min-h-screen xl:flex bg-gray-50 dark:bg-gray-900">
    <!-- 侧边栏 -->
    <AuthSidebar />

    <!-- 背景遮罩 (移动端) -->
    <Backdrop />

    <!-- 主内容区域 -->
    <div
      class="flex-1 transition-all duration-300 ease-in-out"
      :class="[isExpanded || isHovered ? 'lg:ml-[290px]' : 'lg:ml-[90px]']"
    >
      <!-- 头部 -->
      <AuthHeader />

      <!-- 主内容 -->
      <div class="p-4 mx-auto max-w-(--breakpoint-2xl) md:p-6">
        <slot />
      </div>
    </div>

    <!-- 背景遮罩 (移动端) - 备用 -->
    <div
      v-if="isMobileOpen"
      @click="closeMobileSidebar"
      class="fixed inset-0 z-20 bg-black bg-opacity-50 transition-opacity lg:hidden"
    ></div>
  </div>
</template>

<script setup>
import { useSidebar } from '@/composables/useSidebar'
import AuthSidebar from './AuthSidebar.vue'
import AuthHeader from './AuthHeader.vue'
import Backdrop from '@/components/layout/Backdrop.vue'

const { isExpanded, isHovered, isMobileOpen, closeMobileSidebar } = useSidebar()
</script>

<style scoped>
/* 自定义样式 */
</style>
