<template>
  <ThemeProvider>
    <SidebarProvider>
      <RouterView />
      <!-- 会话管理组件 -->
      <SessionManager />
    </SidebarProvider>
  </ThemeProvider>
</template>

<script setup lang="ts">
import ThemeProvider from './components/layout/ThemeProvider.vue'
import SidebarProvider from './components/layout/SidebarProvider.vue'
import SessionManager from './components/SessionManager.vue'
</script>
