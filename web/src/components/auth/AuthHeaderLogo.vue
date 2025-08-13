<template>
  <div class="flex items-center" style="gap: 0.75rem;">
    <!-- Logo -->
    <router-link to="/auth/dashboard" class="flex items-center" style="gap: 0.5rem;">
      <img
        class="h-8 w-auto dark:hidden"
        src="/images/logo/logo-icon.svg"
        alt="Logo"
      />
      <img
        class="hidden h-8 w-auto dark:block"
        src="/images/logo/logo-icon.svg"
        alt="Logo"
      />
      <div class="hidden sm:block">
        <h1 class="text-lg font-bold text-gray-900 dark:text-white/90">
          权限管理系统
        </h1>
        <p class="text-xs text-gray-500 dark:text-gray-400">
          Auth Management System
        </p>
      </div>
    </router-link>

    <!-- 分隔符 -->
    <div class="hidden lg:block h-6 w-px bg-gray-300 dark:bg-gray-700"></div>

    <!-- 当前页面信息 -->
    <div class="hidden lg:flex items-center" style="gap: 0.5rem;">
      <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-brand-50 dark:bg-brand-500/15">
        <component 
          :is="currentPageIcon" 
          class="h-4 w-4 text-brand-600 dark:text-brand-400"
        />
      </div>
      <div>
        <h2 class="text-sm font-semibold text-gray-900 dark:text-white/90">
          {{ currentPageTitle }}
        </h2>
        <p class="text-xs text-gray-500 dark:text-gray-400">
          {{ currentPageDescription }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import {
  BarChartIcon,
  UserGroupIcon,
  SettingsIcon,
  PlugInIcon,
  DocsIcon,
  UserCircleIcon,
  SettingsIcon as CogIcon,
  GridIcon,
} from '@/icons'

const route = useRoute()

// 页面配置映射
const pageConfig = {
  '/auth/dashboard': {
    title: '权限总览',
    description: '系统权限管理概览',
    icon: BarChartIcon
  },
  '/auth/admin': {
    title: '管理员管理',
    description: '管理系统管理员账户',
    icon: UserGroupIcon
  },
  '/auth/role': {
    title: '角色管理',
    description: '管理系统角色和权限',
    icon: SettingsIcon
  },
  '/auth/permission': {
    title: '权限管理',
    description: '管理系统菜单和权限',
    icon: PlugInIcon
  },
  '/auth/logs': {
    title: '系统日志',
    description: '查看系统操作日志',
    icon: DocsIcon
  },
  '/auth/profile': {
    title: '个人资料',
    description: '管理个人信息',
    icon: UserCircleIcon
  },
  '/auth/settings/basic': {
    title: '基础设置',
    description: '系统基础配置',
    icon: CogIcon
  },
  '/auth/settings/security': {
    title: '安全设置',
    description: '系统安全配置',
    icon: CogIcon
  },
  '/auth/settings/notification': {
    title: '通知设置',
    description: '系统通知配置',
    icon: CogIcon
  }
}

// 当前页面标题
const currentPageTitle = computed(() => {
  const config = pageConfig[route.path]
  return config ? config.title : '权限管理系统'
})

// 当前页面描述
const currentPageDescription = computed(() => {
  const config = pageConfig[route.path]
  return config ? config.description : '管理系统权限和用户'
})

// 当前页面图标
const currentPageIcon = computed(() => {
  const config = pageConfig[route.path]
  return config ? config.icon : GridIcon
})
</script>

<style scoped>
/* 自定义样式 */
</style>
