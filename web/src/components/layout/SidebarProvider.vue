/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/components/layout/SidebarProvider.vue
 * @Description: 
 * 
 */
<template>
  <slot></slot>
</template>

<script setup lang="ts">
import { useSidebarProvider } from '@/composables/useSidebar'

useSidebarProvider()
</script>
