<template>
  <button
    @click="handleLogout"
    class="flex items-center px-4 py-2 text-sm text-gray-700 hover:bg-gray-100 dark:text-gray-300 dark:hover:bg-gray-700 rounded-md transition-colors duration-200"
    :disabled="isLoggingOut"
  >
    <svg
      class="w-4 h-4 mr-2"
      fill="none"
      stroke="currentColor"
      viewBox="0 0 24 24"
      xmlns="http://www.w3.org/2000/svg"
    >
      <path
        stroke-linecap="round"
        stroke-linejoin="round"
        stroke-width="2"
        d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"
      ></path>
    </svg>
    {{ isLoggingOut ? '登出中...' : '登出' }}
  </button>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useMainStore } from '@/store'
import { AuthUtils } from '@/utils/auth'
import ToastAlert from '@/composables/ToastAlert'

const router = useRouter()
const store = useMainStore()
const isLoggingOut = ref(false)

/**
 * 处理登出操作
 */
const handleLogout = async () => {
  if (isLoggingOut.value) return
  
  try {
    isLoggingOut.value = true
    
    // 显示确认对话框
    const confirmed = await ToastAlert.confirm({
      title: '确认登出',
      message: '您确定要登出系统吗？',
      variant: 'warning'
    })
    
    if (confirmed) {
      // 执行登出操作
      AuthUtils.logout()
      
      // 显示成功提示
      ToastAlert.success({
        title: '登出成功',
        message: '您已成功登出系统',
        duration: 1500
      })
      
      // 延迟跳转到登录页
      setTimeout(() => {
        router.push('/adminLogin')
      }, 1000)
    }
  } catch (error) {
    console.error('登出失败:', error)
    ToastAlert.error({
      title: '登出失败',
      message: '登出过程中发生错误，请重试'
    })
  } finally {
    isLoggingOut.value = false
  }
}
</script>
