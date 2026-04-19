/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/components/auth/AuthNotificationMenu.vue
 * @Description:
 *
 */
<template>
  <div class="relative" ref="dropdownRef">
    <button
      @click="toggleDropdown"
      class="relative flex h-9 w-9 items-center justify-center rounded-lg text-gray-500 hover:bg-gray-100 hover:text-gray-700 dark:text-gray-400 dark:hover:bg-gray-800 dark:hover:text-gray-300"
    >
      <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M15 17h5l-5 5v-5zM11 19H6.5A2.5 2.5 0 014 16.5v-9A2.5 2.5 0 016.5 5h11A2.5 2.5 0 0120 7.5v3.5"
        />
      </svg>
      <span
        v-if="unreadCount > 0"
        class="absolute -right-1 -top-1 flex h-4 min-w-4 items-center justify-center rounded-full bg-red-500 px-1 text-[10px] font-medium text-white"
      >
        {{ unreadCount > 9 ? '9+' : unreadCount }}
      </span>
    </button>

    <div
      v-if="dropdownOpen"
      class="absolute right-0 z-50 mt-2 w-[22rem] rounded-2xl border border-gray-200 bg-white shadow-lg dark:border-gray-700 dark:bg-gray-800"
    >
      <div class="flex items-center justify-between border-b border-gray-200 p-4 dark:border-gray-700">
        <div>
          <h3 class="text-sm font-semibold text-gray-900 dark:text-white">通知中心</h3>
          <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">实时读取当前用户的站内通知</p>
        </div>
        <div class="flex items-center gap-3">
          <button class="text-xs text-gray-500 hover:text-brand-600 dark:text-gray-400" @click="loadNotifications">
            刷新
          </button>
          <button class="text-xs text-brand-600 hover:text-brand-700 dark:text-brand-400" @click="markAllAsRead">
            全部已读
          </button>
        </div>
      </div>

      <div v-if="loading" class="space-y-3 p-4">
        <div
          v-for="index in 3"
          :key="index"
          class="h-18 animate-pulse rounded-xl bg-gray-100 dark:bg-gray-700/60"
        ></div>
      </div>

      <div v-else-if="notifications.length === 0" class="p-8 text-center">
        <p class="text-sm text-gray-500 dark:text-gray-400">暂无通知</p>
      </div>

      <div v-else class="max-h-96 overflow-y-auto">
        <button
          v-for="notification in notifications"
          :key="notification.id"
          class="w-full border-b border-gray-100 p-4 text-left transition-colors last:border-b-0 dark:border-gray-700"
          :class="notification.read ? 'hover:bg-gray-50 dark:hover:bg-gray-700/60' : 'bg-brand-50/60 hover:bg-brand-100/80 dark:bg-brand-500/10 dark:hover:bg-brand-500/15'"
          @click="handleNotificationClick(notification)"
        >
          <div class="flex items-start gap-3">
            <div
              class="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-xl text-xs font-semibold"
              :class="badgeClass(notification.noticeLevel)"
            >
              {{ levelShortLabel(notification.noticeLevel) }}
            </div>
            <div class="min-w-0 flex-1">
              <div class="flex items-center gap-2">
                <h4 class="truncate text-sm font-medium text-gray-900 dark:text-white">{{ notification.title }}</h4>
                <span
                  v-if="!notification.read"
                  class="h-2 w-2 shrink-0 rounded-full bg-brand-500"
                ></span>
              </div>
              <p class="mt-1 line-clamp-2 text-xs leading-5 text-gray-600 dark:text-gray-300">
                {{ notification.content }}
              </p>
              <div class="mt-2 flex flex-wrap items-center gap-2 text-[11px] text-gray-500 dark:text-gray-400">
                <span>{{ noticeTypeLabel(notification.noticeType) }}</span>
                <span>{{ formatTime(notification.publishTime || notification.createTime) }}</span>
              </div>
            </div>
          </div>
        </button>
      </div>

      <div class="border-t border-gray-200 p-3 dark:border-gray-700">
        <router-link
          to="/auth/notifications"
          class="block rounded-xl py-2 text-center text-sm text-brand-600 hover:bg-brand-50 dark:text-brand-400 dark:hover:bg-brand-500/10"
          @click="dropdownOpen = false"
        >
          查看全部通知
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { useRouter } from 'vue-router'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'

interface CurrentNotice {
  id: number
  title: string
  content: string
  noticeType: string
  noticeLevel: string
  publishTime?: string
  createTime?: string
  read?: boolean
}

const router = useRouter()
const dropdownRef = ref<HTMLElement | null>(null)
const dropdownOpen = ref(false)
const loading = ref(false)
const notifications = ref<CurrentNotice[]>([])

const unreadCount = computed(() => notifications.value.filter((item) => !item.read).length)

const resolveList = (payload: unknown): CurrentNotice[] => {
  if (Array.isArray(payload)) {
    return payload as CurrentNotice[]
  }
  if (payload && typeof payload === 'object') {
    const list = (payload as { list?: unknown }).list
    if (Array.isArray(list)) {
      return list as CurrentNotice[]
    }
  }
  return []
}

const noticeTypeLabel = (type?: string) => {
  switch (type) {
    case 'security':
      return '安全通知'
    case 'maintenance':
      return '维护通知'
    default:
      return '系统通知'
  }
}

const levelShortLabel = (level?: string) => {
  switch (level) {
    case 'critical':
      return '紧急'
    case 'warning':
      return '提醒'
    default:
      return '通知'
  }
}

const badgeClass = (level?: string) => {
  switch (level) {
    case 'critical':
      return 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300'
    case 'warning':
      return 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300'
    default:
      return 'bg-blue-100 text-blue-700 dark:bg-blue-500/15 dark:text-blue-300'
  }
}

const formatTime = (value?: string) => {
  if (!value) {
    return '刚刚'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  const diff = Date.now() - date.getTime()
  const minute = 60 * 1000
  const hour = 60 * minute
  const day = 24 * hour
  if (diff < minute) return '刚刚'
  if (diff < hour) return `${Math.floor(diff / minute)} 分钟前`
  if (diff < day) return `${Math.floor(diff / hour)} 小时前`
  return `${Math.floor(diff / day)} 天前`
}

const loadNotifications = async () => {
  try {
    loading.value = true
    const { data: res } = await adminApi.getCurrentNotices({
      limit: 8
    })
    if (res.code !== 200) {
      throw new Error(res.message || '通知读取失败')
    }
    notifications.value = resolveList(res.data)
  } catch (error) {
    notifications.value = []
    ToastAlert.error({
      title: '通知读取失败',
      message: error instanceof Error ? error.message : '无法读取当前通知'
    })
  } finally {
    loading.value = false
  }
}

const toggleDropdown = async () => {
  dropdownOpen.value = !dropdownOpen.value
  if (dropdownOpen.value) {
    await loadNotifications()
  }
}

const markAllAsRead = async () => {
  try {
    const { data: res } = await adminApi.markAllNoticeRead()
    if (res.code !== 200) {
      throw new Error(res.message || '全部已读失败')
    }
    await loadNotifications()
    ToastAlert.success({
      title: '操作成功',
      message: '当前通知已全部标记为已读'
    })
  } catch (error) {
    ToastAlert.error({
      title: '操作失败',
      message: error instanceof Error ? error.message : '全部已读失败'
    })
  }
}

const handleNotificationClick = async (notification: CurrentNotice) => {
  try {
    if (!notification.read) {
      const { data: res } = await adminApi.markNoticeRead(notification.id)
      if (res.code !== 200) {
        throw new Error(res.message || '标记已读失败')
      }
    }
    dropdownOpen.value = false
    await router.push('/auth/notifications')
  } catch (error) {
    ToastAlert.error({
      title: '跳转失败',
      message: error instanceof Error ? error.message : '通知跳转失败'
    })
  }
}

let pollingInterval: ReturnType<typeof setInterval> | null = null

const startPolling = () => {
  // 按照企业级架构，每 30 秒静默轮询一次通知
  pollingInterval = setInterval(async () => {
    try {
      const { data: res } = await adminApi.getCurrentNotices({ limit: 8 })
      if (res.code === 200) {
        notifications.value = resolveList(res.data)
      }
    } catch (e) {
      // 静默轮询失败不提示，避免打扰用户
      console.warn("通知静默轮询失败", e)
    }
  }, 30000)
}

const handleClickOutside = (event: MouseEvent) => {
  const target = event.target as Node
  if (dropdownRef.value && !dropdownRef.value.contains(target)) {
    dropdownOpen.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  loadNotifications()
  startPolling()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
  if (pollingInterval) clearInterval(pollingInterval)
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
