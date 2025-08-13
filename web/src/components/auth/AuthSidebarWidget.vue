<template>
  <div class="mt-auto mb-6">
    <!-- 系统状态卡片 -->
    <div class="rounded-2xl border border-gray-200 bg-gradient-to-br from-brand-50 to-brand-100 p-4 dark:border-gray-800 dark:from-brand-500/10 dark:to-brand-600/10">
      <div class="flex items-center mb-3" style="gap: 0.75rem;">
        <div class="flex h-10 w-10 items-center justify-center rounded-full bg-brand-500 dark:bg-brand-600">
          <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
          </svg>
        </div>
        <div>
          <h4 class="text-sm font-semibold text-gray-900 dark:text-white/90">系统状态</h4>
          <p class="text-xs text-gray-600 dark:text-gray-400">权限系统运行正常</p>
        </div>
      </div>
      
      <!-- 状态指标 -->
      <div class="space-y-2">
        <div class="flex items-center justify-between text-xs">
          <span class="text-gray-600 dark:text-gray-400">在线用户</span>
          <span class="font-medium text-gray-900 dark:text-white/90">{{ systemStats.onlineUsers }}</span>
        </div>
        <div class="flex items-center justify-between text-xs">
          <span class="text-gray-600 dark:text-gray-400">活跃会话</span>
          <span class="font-medium text-gray-900 dark:text-white/90">{{ systemStats.activeSessions }}</span>
        </div>
        <div class="flex items-center justify-between text-xs">
          <span class="text-gray-600 dark:text-gray-400">系统负载</span>
          <span class="font-medium text-green-600 dark:text-green-400">正常</span>
        </div>
      </div>

      <!-- 快速操作按钮 -->
      <div class="mt-4 flex" style="gap: 0.5rem;">
        <router-link
          to="/auth/dashboard"
          class="flex-1 rounded-lg bg-white/80 px-3 py-2 text-center text-xs font-medium text-brand-700 transition-colors hover:bg-white dark:bg-white/10 dark:text-brand-400 dark:hover:bg-white/20"
        >
          查看详情
        </router-link>
        <button
          @click="refreshStats"
          class="rounded-lg bg-white/80 px-3 py-2 text-xs font-medium text-brand-700 transition-colors hover:bg-white dark:bg-white/10 dark:text-brand-400 dark:hover:bg-white/20"
        >
          刷新
        </button>
      </div>
    </div>

    <!-- 快速链接 -->
    <div class="mt-4 rounded-2xl border border-gray-200 bg-white p-4 dark:border-gray-800 dark:bg-white/[0.03]">
      <h4 class="mb-3 text-sm font-semibold text-gray-900 dark:text-white/90">快速链接</h4>
      <div class="space-y-2">
        <router-link
          v-for="link in quickLinks"
          :key="link.path"
          :to="link.path"
          class="flex items-center rounded-lg px-2 py-1.5 text-xs text-gray-600 transition-colors hover:bg-gray-100 hover:text-gray-900 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-300"
          style="gap: 0.5rem;"
        >
          <component :is="link.icon" class="h-4 w-4" />
          {{ link.name }}
        </router-link>
      </div>
    </div>

    <!-- 帮助信息 -->
    <div class="mt-4 rounded-2xl border border-gray-200 bg-white p-4 dark:border-gray-800 dark:bg-white/[0.03]">
      <div class="flex items-center mb-2" style="gap: 0.5rem;">
        <svg class="h-4 w-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h4 class="text-sm font-semibold text-gray-900 dark:text-white/90">帮助</h4>
      </div>
      <p class="text-xs text-gray-600 dark:text-gray-400 mb-3">
        需要帮助？查看我们的文档或联系技术支持。
      </p>
      <div class="flex" style="gap: 0.5rem;">
        <button
          @click="openHelp"
          class="flex-1 rounded-lg bg-blue-50 px-2 py-1.5 text-xs font-medium text-blue-700 transition-colors hover:bg-blue-100 dark:bg-blue-500/15 dark:text-blue-400 dark:hover:bg-blue-500/25"
        >
          查看文档
        </button>
        <button
          @click="contactSupport"
          class="rounded-lg bg-gray-100 px-2 py-1.5 text-xs font-medium text-gray-700 transition-colors hover:bg-gray-200 dark:bg-gray-800 dark:text-gray-300 dark:hover:bg-gray-700"
        >
          联系支持
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import {
  PlusIcon,
  DocsIcon,
  SettingsIcon,
  InfoCircleIcon,
} from '@/icons'
import ToastAlert from '@/composables/ToastAlert'

// 系统统计数据
const systemStats = ref({
  onlineUsers: 8,
  activeSessions: 12,
  systemLoad: 'normal'
})

// 快速链接配置
const quickLinks = [
  {
    name: '添加管理员',
    path: '/auth/admin',
    icon: PlusIcon
  },
  {
    name: '查看日志',
    path: '/auth/logs',
    icon: DocsIcon
  },
  {
    name: '系统设置',
    path: '/auth/settings/basic',
    icon: SettingsIcon
  },
  {
    name: '帮助文档',
    path: '/help',
    icon: InfoCircleIcon
  }
]

// 刷新统计数据
const refreshStats = async () => {
  try {
    // 这里可以调用API获取最新的系统状态
    // const response = await api.getSystemStats()
    // systemStats.value = response.data
    
    // 模拟刷新
    systemStats.value = {
      onlineUsers: Math.floor(Math.random() * 20) + 5,
      activeSessions: Math.floor(Math.random() * 30) + 10,
      systemLoad: 'normal'
    }
    
    ToastAlert.success({
      title: '刷新成功',
      message: '系统状态已更新'
    })
  } catch (error) {
    ToastAlert.error({
      title: '刷新失败',
      message: '无法获取系统状态'
    })
  }
}

// 打开帮助文档
const openHelp = () => {
  // 可以打开新窗口或跳转到帮助页面
  window.open('/help/auth-management', '_blank')
}

// 联系技术支持
const contactSupport = () => {
  ToastAlert.info({
    title: '联系支持',
    message: '请发送邮件至 support@example.com 或拨打 400-123-4567'
  })
}

// 组件挂载时获取初始数据
onMounted(() => {
  // 可以在这里获取初始的系统状态数据
})
</script>

<style scoped>
/* 自定义样式可以在这里添加 */
</style>
