/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/App.vue
 * @Description: 
 * 
 */
<template>
  <ThemeProvider>
    <SidebarProvider>
      <router-view v-slot="{ Component, route }">
        <keep-alive :include="cachedViews">
          <component v-if="Component" :is="Component" :key="route.fullPath" />
        </keep-alive>
      </router-view>
      <!-- 会话管理组件 -->
      <SessionManager />
    </SidebarProvider>
  </ThemeProvider>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import ThemeProvider from './components/layout/ThemeProvider.vue'
import SidebarProvider from './components/layout/SidebarProvider.vue'
import SessionManager from './components/SessionManager.vue'
import { useTagsViewStore } from './stores/tagsView'

const tagsViewStore = useTagsViewStore()
const cachedViews = computed(() => tagsViewStore.cachedViews)
</script>
