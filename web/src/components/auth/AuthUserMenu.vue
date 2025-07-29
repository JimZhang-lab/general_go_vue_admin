<template>
  <div class="relative" ref="dropdownRef">
    <!-- 用户按钮 -->
    <button
      class="flex items-center text-gray-700 dark:text-gray-400"
      @click.prevent="toggleDropdown"
    >
      <span class="mr-3 overflow-hidden rounded-full h-11 w-11">
        <img 
          :src="userInfo.avatar || '/images/user/owner.jpg'" 
          :alt="userInfo.name || '用户头像'"
          class="w-full h-full object-cover"
        />
      </span>

      <span class="block mr-1 font-medium text-theme-sm">
        {{ userInfo.name || '管理员' }}
      </span>

      <svg
        :class="{ 'rotate-180': dropdownOpen }"
        class="w-4 h-4 transition-transform duration-200"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M19 9l-7 7-7-7"
        />
      </svg>
    </button>

    <!-- 用户下拉菜单 -->
    <div
      v-if="dropdownOpen"
      class="absolute right-0 mt-[17px] flex w-[280px] flex-col rounded-2xl border border-gray-200 bg-white p-3 shadow-theme-lg dark:border-gray-800 dark:bg-gray-dark"
    >
      <!-- 用户信息 -->
      <div class="p-3 border-b border-gray-200 dark:border-gray-700">
        <div class="flex items-center" style="gap: 0.75rem;">
          <div class="overflow-hidden rounded-full h-12 w-12">
            <img 
              :src="userInfo.avatar || '/images/user/owner.jpg'" 
              :alt="userInfo.name || '用户头像'"
              class="w-full h-full object-cover"
            />
          </div>
          <div class="flex-1 min-w-0">
            <h4 class="font-medium text-gray-700 text-theme-sm dark:text-gray-400 truncate">
              {{ userInfo.name || '系统管理员' }}
            </h4>
            <p class="mt-0.5 text-theme-xs text-gray-500 dark:text-gray-400 truncate">
              {{ userInfo.email || 'admin@example.com' }}
            </p>
            <div class="flex items-center mt-1" style="gap: 0.5rem;">
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800 dark:bg-green-500/15 dark:text-green-400">
                <div class="w-1.5 h-1.5 bg-green-500 rounded-full mr-1"></div>
                在线
              </span>
              <span class="text-xs text-gray-500 dark:text-gray-400">
                {{ userInfo.role || '超级管理员' }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- 菜单项 -->
      <ul class="flex flex-col py-3" style="gap: 0.25rem;">
        <li v-for="item in menuItems" :key="item.href">
          <router-link
            :to="item.href"
            @click="closeDropdown"
            class="flex items-center px-3 py-2.5 font-medium text-gray-700 rounded-lg group text-theme-sm hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-white/5 dark:hover:text-gray-300 transition-colors"
            style="gap: 0.75rem;"
          >
            <div
              :class="[
                'flex h-8 w-8 items-center justify-center rounded-lg',
                item.iconBg,
              ]"
            >
              <component
                :is="item.icon"
                class="h-4 w-4"
                :class="item.iconColor"
              />
            </div>
            <div class="flex-1">
              <div class="font-medium">{{ item.text }}</div>
              <div class="text-xs text-gray-500 dark:text-gray-400">{{ item.description }}</div>
            </div>
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </router-link>
        </li>
      </ul>

      <!-- 分隔线 -->
      <div class="border-t border-gray-200 dark:border-gray-700 my-2"></div>

      <!-- 系统信息 -->
      <div class="px-3 py-2 bg-gray-50 dark:bg-gray-800/50 rounded-lg mb-3">
        <div class="flex items-center justify-between text-xs">
          <span class="text-gray-600 dark:text-gray-400">系统版本</span>
          <span class="font-medium text-gray-900 dark:text-white/90">v2.1.0</span>
        </div>
        <div class="flex items-center justify-between text-xs mt-1">
          <span class="text-gray-600 dark:text-gray-400">最后登录</span>
          <span class="font-medium text-gray-900 dark:text-white/90">{{ lastLoginTime }}</span>
        </div>
      </div>

      <!-- 退出登录 -->
      <button
        @click="signOut"
        class="flex items-center px-3 py-2.5 font-medium text-red-600 rounded-lg group text-theme-sm hover:bg-red-50 hover:text-red-700 dark:text-red-400 dark:hover:bg-red-500/10 dark:hover:text-red-300 transition-colors"
        style="gap: 0.75rem;"
      >
        <div class="flex h-8 w-8 items-center justify-center rounded-lg bg-red-50 dark:bg-red-500/15">
          <svg class="h-4 w-4 text-red-600 dark:text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
          </svg>
        </div>
        <div class="flex-1 text-left">
          <div class="font-medium">退出登录</div>
          <div class="text-xs text-red-500 dark:text-red-400">安全退出系统</div>
        </div>
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import ToastAlert from '@/composables/ToastAlert'

const router = useRouter()
const dropdownRef = ref(null)
const dropdownOpen = ref(false)

// 用户信息
const userInfo = ref({
  name: '系统管理员',
  email: 'admin@example.com',
  avatar: '/images/user/owner.jpg',
  role: '超级管理员',
  lastLogin: new Date(Date.now() - 2 * 60 * 60 * 1000) // 2小时前
})

// 菜单项配置
const menuItems = [
  {
    href: '/auth/profile',
    text: '个人资料',
    description: '查看和编辑个人信息',
    icon: 'div', // UserCircleIcon
    iconBg: 'bg-blue-50 dark:bg-blue-500/15',
    iconColor: 'text-blue-600 dark:text-blue-400'
  },
  {
    href: '/auth/settings/security',
    text: '安全设置',
    description: '密码和安全选项',
    icon: 'div', // ShieldCheckIcon
    iconBg: 'bg-green-50 dark:bg-green-500/15',
    iconColor: 'text-green-600 dark:text-green-400'
  },
  {
    href: '/auth/settings/notification',
    text: '通知设置',
    description: '管理通知偏好',
    icon: 'div', // BellIcon
    iconBg: 'bg-purple-50 dark:bg-purple-500/15',
    iconColor: 'text-purple-600 dark:text-purple-400'
  },
  {
    href: '/auth/logs',
    text: '活动日志',
    description: '查看操作记录',
    icon: 'div', // DocumentTextIcon
    iconBg: 'bg-orange-50 dark:bg-orange-500/15',
    iconColor: 'text-orange-600 dark:text-orange-400'
  },
  {
    href: '/help',
    text: '帮助中心',
    description: '获取帮助和支持',
    icon: 'div', // QuestionMarkCircleIcon
    iconBg: 'bg-gray-50 dark:bg-gray-500/15',
    iconColor: 'text-gray-600 dark:text-gray-400'
  }
]

// 最后登录时间
const lastLoginTime = computed(() => {
  const time = userInfo.value.lastLogin
  const now = new Date()
  const diff = now - time
  const hours = Math.floor(diff / (1000 * 60 * 60))
  
  if (hours < 1) return '刚刚'
  if (hours < 24) return `${hours}小时前`
  
  return time.toLocaleDateString('zh-CN')
})

// 切换下拉菜单
const toggleDropdown = () => {
  dropdownOpen.value = !dropdownOpen.value
}

// 关闭下拉菜单
const closeDropdown = () => {
  dropdownOpen.value = false
}

// 退出登录
const signOut = async () => {
  try {
    // 显示确认对话框
    const confirmed = await ToastAlert.confirm({
      title: '确认退出',
      message: '您确定要退出系统吗？',
      variant: 'warning'
    })

    if (!confirmed) return

    // 调用后端退出登录API
    const adminApi = (await import('@/api/system')).default
    await adminApi.logout()

    // 使用AuthUtils清除所有用户数据
    const { AuthUtils } = await import('@/utils/auth')
    AuthUtils.logout()

    ToastAlert.success({
      title: '退出成功',
      message: '您已安全退出系统'
    })

    // 延迟跳转到登录页面
    setTimeout(() => {
      router.push('/adminLogin')
    }, 1000)
  } catch (error) {
    console.error('退出登录失败:', error)

    // 即使API调用失败，也要清除本地数据
    const { AuthUtils } = await import('@/utils/auth')
    AuthUtils.logout()

    ToastAlert.error({
      title: '退出失败',
      message: '退出登录时发生错误，但已清除本地数据'
    })

    // 跳转到登录页面
    setTimeout(() => {
      router.push('/adminLogin')
    }, 1000)
  }
}

// 点击外部关闭
const handleClickOutside = (event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  
  // 从本地存储或API获取用户信息
  const storedUserInfo = localStorage.getItem('user_info')
  if (storedUserInfo) {
    try {
      const parsed = JSON.parse(storedUserInfo)
      userInfo.value = { ...userInfo.value, ...parsed }
    } catch (error) {
      console.error('解析用户信息失败:', error)
    }
  }
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
/* 自定义样式 */
</style>
