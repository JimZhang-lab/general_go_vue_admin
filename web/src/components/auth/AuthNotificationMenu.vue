<template>
  <div class="relative" ref="dropdownRef">
    <!-- 通知按钮 -->
    <button
      @click="toggleDropdown"
      class="relative flex items-center justify-center w-9 h-9 text-gray-500 rounded-lg hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-300"
    >
      <svg
        class="w-5 h-5"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M15 17h5l-5 5v-5zM11 19H6.5A2.5 2.5 0 014 16.5v-9A2.5 2.5 0 016.5 5h11A2.5 2.5 0 0120 7.5v3.5"
        />
      </svg>
      <!-- 未读通知数量 -->
      <span
        v-if="unreadCount > 0"
        class="absolute -top-1 -right-1 flex h-4 w-4 items-center justify-center rounded-full bg-red-500 text-xs font-medium text-white"
      >
        {{ unreadCount > 9 ? '9+' : unreadCount }}
      </span>
    </button>

    <!-- 通知下拉菜单 -->
    <div
      v-if="dropdownOpen"
      class="absolute right-0 mt-2 w-80 bg-white border border-gray-200 rounded-2xl shadow-lg z-50 dark:bg-gray-800 dark:border-gray-700"
    >
      <!-- 头部 -->
      <div class="flex items-center justify-between p-4 border-b border-gray-200 dark:border-gray-700">
        <h3 class="text-sm font-semibold text-gray-900 dark:text-white">
          通知中心
        </h3>
        <div class="flex items-center" style="gap: 0.5rem;">
          <button
            @click="markAllAsRead"
            class="text-xs text-brand-600 hover:text-brand-700 dark:text-brand-400"
          >
            全部已读
          </button>
          <button
            @click="clearAll"
            class="text-xs text-gray-500 hover:text-gray-700 dark:text-gray-400"
          >
            清空
          </button>
        </div>
      </div>

      <!-- 通知列表 -->
      <div class="max-h-96 overflow-y-auto">
        <div v-if="notifications.length > 0">
          <div
            v-for="notification in notifications"
            :key="notification.id"
            @click="handleNotificationClick(notification)"
            :class="[
              'p-4 border-b border-gray-100 dark:border-gray-700 last:border-b-0 cursor-pointer transition-colors',
              {
                'bg-brand-50 dark:bg-brand-500/10': !notification.read,
                'hover:bg-gray-50 dark:hover:bg-gray-700': notification.read,
                'hover:bg-brand-100 dark:hover:bg-brand-500/20': !notification.read,
              },
            ]"
          >
            <div class="flex items-start" style="gap: 0.75rem;">
              <!-- 通知图标 -->
              <div
                :class="[
                  'flex h-8 w-8 items-center justify-center rounded-lg flex-shrink-0',
                  getNotificationStyle(notification.type),
                ]"
              >
                <component :is="getNotificationIcon(notification.type)" class="h-4 w-4" />
              </div>
              
              <!-- 通知内容 -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center" style="gap: 0.5rem;">
                  <h4 class="text-sm font-medium text-gray-900 dark:text-white truncate">
                    {{ notification.title }}
                  </h4>
                  <div
                    v-if="!notification.read"
                    class="w-2 h-2 bg-brand-500 rounded-full flex-shrink-0"
                  ></div>
                </div>
                <p class="text-xs text-gray-600 dark:text-gray-400 mt-1 line-clamp-2">
                  {{ notification.message }}
                </p>
                <div class="flex items-center justify-between mt-2">
                  <span class="text-xs text-gray-500 dark:text-gray-400">
                    {{ formatTime(notification.time) }}
                  </span>
                  <span
                    :class="[
                      'px-2 py-0.5 text-xs font-medium rounded-full',
                      getTypeBadgeStyle(notification.type),
                    ]"
                  >
                    {{ getTypeLabel(notification.type) }}
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 空状态 -->
        <div v-else class="p-8 text-center">
          <svg
            class="w-12 h-12 mx-auto text-gray-400 mb-3"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M15 17h5l-5 5v-5zM11 19H6.5A2.5 2.5 0 014 16.5v-9A2.5 2.5 0 016.5 5h11A2.5 2.5 0 0120 7.5v3.5"
            />
          </svg>
          <p class="text-sm text-gray-500 dark:text-gray-400">
            暂无通知
          </p>
        </div>
      </div>

      <!-- 底部 -->
      <div class="p-3 border-t border-gray-200 dark:border-gray-700">
        <router-link
          to="/auth/notifications"
          class="block w-full text-center py-2 text-sm text-brand-600 hover:text-brand-700 dark:text-brand-400"
        >
          查看所有通知
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import ToastAlert from '@/composables/ToastAlert'

const dropdownRef = ref(null)
const dropdownOpen = ref(false)

// 模拟通知数据
const notifications = ref([
  {
    id: 1,
    type: 'security',
    title: '安全警告',
    message: '检测到异常登录尝试，来自IP: 192.168.1.100',
    time: new Date(Date.now() - 5 * 60 * 1000), // 5分钟前
    read: false,
  },
  {
    id: 2,
    type: 'user',
    title: '新用户注册',
    message: '用户 "张三" 已成功注册并等待审核',
    time: new Date(Date.now() - 30 * 60 * 1000), // 30分钟前
    read: false,
  },
  {
    id: 3,
    type: 'system',
    title: '系统更新',
    message: '权限管理系统已更新到版本 v2.1.0',
    time: new Date(Date.now() - 2 * 60 * 60 * 1000), // 2小时前
    read: true,
  },
])

// 未读通知数量
const unreadCount = computed(() => {
  return notifications.value.filter(n => !n.read).length
})

// 切换下拉菜单
const toggleDropdown = () => {
  dropdownOpen.value = !dropdownOpen.value
}

// 标记所有为已读
const markAllAsRead = () => {
  notifications.value.forEach(n => n.read = true)
  ToastAlert.success({
    title: '操作成功',
    message: '所有通知已标记为已读'
  })
}

// 清空所有通知
const clearAll = () => {
  notifications.value = []
  ToastAlert.success({
    title: '操作成功',
    message: '所有通知已清空'
  })
}

// 处理通知点击
const handleNotificationClick = (notification) => {
  notification.read = true
}

// 获取通知图标
const getNotificationIcon = (type) => {
  // 返回简单的SVG图标组件
  return 'div'
}

// 获取通知样式
const getNotificationStyle = (type) => {
  const styles = {
    security: 'bg-red-50 text-red-600 dark:bg-red-500/15 dark:text-red-400',
    user: 'bg-blue-50 text-blue-600 dark:bg-blue-500/15 dark:text-blue-400',
    system: 'bg-purple-50 text-purple-600 dark:bg-purple-500/15 dark:text-purple-400',
  }
  return styles[type] || 'bg-gray-50 text-gray-600 dark:bg-gray-500/15 dark:text-gray-400'
}

// 获取类型标签样式
const getTypeBadgeStyle = (type) => {
  const styles = {
    security: 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-400',
    user: 'bg-blue-100 text-blue-700 dark:bg-blue-500/15 dark:text-blue-400',
    system: 'bg-purple-100 text-purple-700 dark:bg-purple-500/15 dark:text-purple-400',
  }
  return styles[type] || 'bg-gray-100 text-gray-700 dark:bg-gray-500/15 dark:text-gray-400'
}

// 获取类型标签
const getTypeLabel = (type) => {
  const labels = {
    security: '安全',
    user: '用户',
    system: '系统',
  }
  return labels[type] || '通知'
}

// 格式化时间
const formatTime = (time) => {
  const now = new Date()
  const diff = now - time
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  if (hours < 24) return `${hours}小时前`
  if (days < 7) return `${days}天前`
  
  return time.toLocaleDateString('zh-CN')
}

// 点击外部关闭
const handleClickOutside = (event) => {
  if (dropdownRef.value && !dropdownRef.value.contains(event.target)) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
