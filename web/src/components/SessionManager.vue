<template>
  <!-- 会话管理组件，无UI，仅处理会话逻辑 -->
</template>

<script setup lang="ts">
import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { AuthUtils } from '@/utils/auth'
import ToastAlert from '@/composables/ToastAlert'

const router = useRouter()

let sessionCheckInterval: number | null = null
let activityTimer: number | null = null

/**
 * 检查会话状态
 */
const checkSession = () => {
  if (!AuthUtils.isAuthenticated() || !AuthUtils.isSessionActive()) {
    handleSessionExpired()
  }
}

/**
 * 处理会话过期
 */
const handleSessionExpired = () => {
  // 清除定时器
  clearIntervals()
  
  // 显示会话过期提示
  ToastAlert.warning({
    title: '会话已过期',
    message: '您的登录会话已过期，请重新登录',
    duration: 3000
  })
  
  // 清除用户数据
  AuthUtils.logout()
  
  // 延迟跳转到登录页
  setTimeout(() => {
    router.push({
      path: '/adminLogin',
      query: { redirect: router.currentRoute.value.fullPath }
    })
  }, 2000)
}

/**
 * 处理用户活动
 */
const handleUserActivity = () => {
  // 如果用户已登录，刷新会话
  if (AuthUtils.isAuthenticated()) {
    AuthUtils.refreshSession()
  }
}

/**
 * 清除所有定时器
 */
const clearIntervals = () => {
  if (sessionCheckInterval) {
    clearInterval(sessionCheckInterval)
    sessionCheckInterval = null
  }
  if (activityTimer) {
    clearTimeout(activityTimer)
    activityTimer = null
  }
}

/**
 * 初始化会话管理
 */
const initSessionManager = () => {
  // 只在用户已登录时启动会话管理
  if (!AuthUtils.isAuthenticated()) {
    return
  }
  
  // 每分钟检查一次会话状态
  sessionCheckInterval = setInterval(checkSession, 60 * 1000) as unknown as number
  
  // 监听用户活动事件
  const events = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart', 'click']
  
  const throttledActivityHandler = throttle(handleUserActivity, 30000) // 30秒内最多触发一次
  
  events.forEach(event => {
    document.addEventListener(event, throttledActivityHandler, true)
  })
  
  // 监听页面可见性变化
  document.addEventListener('visibilitychange', () => {
    if (!document.hidden && AuthUtils.isAuthenticated()) {
      checkSession()
    }
  })
}

/**
 * 节流函数
 */
function throttle(func: Function, limit: number) {
  let inThrottle: boolean
  return function(this: any, ...args: any[]) {
    if (!inThrottle) {
      func.apply(this, args)
      inThrottle = true
      setTimeout(() => inThrottle = false, limit)
    }
  }
}

/**
 * 清理会话管理
 */
const cleanupSessionManager = () => {
  clearIntervals()
  
  // 移除事件监听器
  const events = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart', 'click']
  events.forEach(event => {
    document.removeEventListener(event, handleUserActivity, true)
  })
  
  document.removeEventListener('visibilitychange', checkSession)
}

// 组件挂载时初始化
onMounted(() => {
  initSessionManager()
})

// 组件卸载时清理
onUnmounted(() => {
  cleanupSessionManager()
})

// 监听路由变化，重新初始化会话管理
router.afterEach(() => {
  // 延迟一点时间，确保路由切换完成
  setTimeout(() => {
    cleanupSessionManager()
    initSessionManager()
  }, 100)
})
</script>
