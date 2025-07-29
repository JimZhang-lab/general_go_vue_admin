<template>
  <div class="relative flex-1 max-w-md mx-4">
    <!-- 搜索输入框 -->
    <div class="relative">
      <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
        <svg
          class="w-4 h-4 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
          />
        </svg>
      </div>
      <input
        v-model="searchQuery"
        @input="handleSearch"
        @focus="showResults = true"
        @keydown.escape="hideResults"
        @keydown.enter="selectResult(selectedIndex)"
        @keydown.arrow-down.prevent="navigateResults(1)"
        @keydown.arrow-up.prevent="navigateResults(-1)"
        type="text"
        placeholder="搜索管理员、角色、权限..."
        class="w-full pl-10 pr-4 py-2 text-sm border border-gray-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-brand-500 focus:border-transparent dark:bg-gray-800 dark:border-gray-600 dark:text-white dark:placeholder-gray-400"
      />
      <!-- 清除按钮 -->
      <button
        v-if="searchQuery"
        @click="clearSearch"
        class="absolute inset-y-0 right-0 flex items-center pr-3"
      >
        <svg
          class="w-4 h-4 text-gray-400 hover:text-gray-600"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M6 18L18 6M6 6l12 12"
          />
        </svg>
      </button>
    </div>

    <!-- 搜索结果下拉框 -->
    <div
      v-if="showResults && (searchResults.length > 0 || searchQuery)"
      class="absolute top-full left-0 right-0 mt-1 bg-white border border-gray-200 rounded-lg shadow-lg z-50 max-h-96 overflow-y-auto dark:bg-gray-800 dark:border-gray-600"
    >
      <!-- 搜索结果 -->
      <div v-if="searchResults.length > 0">
        <div
          v-for="(result, index) in searchResults"
          :key="result.id"
          @click="selectResult(index)"
          @mouseenter="selectedIndex = index"
          :class="[
            'px-4 py-3 cursor-pointer border-b border-gray-100 dark:border-gray-700 last:border-b-0',
            {
              'bg-brand-50 dark:bg-brand-500/15': selectedIndex === index,
              'hover:bg-gray-50 dark:hover:bg-gray-700': selectedIndex !== index,
            },
          ]"
        >
          <div class="flex items-center" style="gap: 0.75rem;">
            <!-- 结果类型图标 -->
            <div
              :class="[
                'flex h-8 w-8 items-center justify-center rounded-lg',
                getTypeStyle(result.type),
              ]"
            >
              <component :is="getTypeIcon(result.type)" class="h-4 w-4" />
            </div>
            <!-- 结果信息 -->
            <div class="flex-1 min-w-0">
              <div class="flex items-center" style="gap: 0.5rem;">
                <h4 class="text-sm font-medium text-gray-900 dark:text-white truncate">
                  {{ result.title }}
                </h4>
                <span
                  :class="[
                    'px-2 py-0.5 text-xs font-medium rounded-full',
                    getTypeBadgeStyle(result.type),
                  ]"
                >
                  {{ getTypeLabel(result.type) }}
                </span>
              </div>
              <p class="text-xs text-gray-500 dark:text-gray-400 truncate">
                {{ result.description }}
              </p>
            </div>
            <!-- 快捷键提示 -->
            <div class="text-xs text-gray-400">
              <kbd class="px-1.5 py-0.5 bg-gray-100 rounded dark:bg-gray-700">
                {{ index === selectedIndex ? '↵' : '' }}
              </kbd>
            </div>
          </div>
        </div>
      </div>

      <!-- 无结果提示 -->
      <div v-else-if="searchQuery" class="px-4 py-6 text-center">
        <svg
          class="w-8 h-8 mx-auto text-gray-400 mb-2"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M9.172 16.172a4 4 0 015.656 0M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"
          />
        </svg>
        <p class="text-sm text-gray-500 dark:text-gray-400">
          没有找到相关结果
        </p>
        <p class="text-xs text-gray-400 dark:text-gray-500 mt-1">
          尝试使用不同的关键词搜索
        </p>
      </div>

      <!-- 搜索提示 -->
      <div v-else class="px-4 py-3 border-t border-gray-100 dark:border-gray-700">
        <div class="flex items-center text-xs text-gray-500 dark:text-gray-400" style="gap: 0.5rem;">
          <kbd class="px-1.5 py-0.5 bg-gray-100 rounded dark:bg-gray-700">↑↓</kbd>
          <span>导航</span>
          <kbd class="px-1.5 py-0.5 bg-gray-100 rounded dark:bg-gray-700">↵</kbd>
          <span>选择</span>
          <kbd class="px-1.5 py-0.5 bg-gray-100 rounded dark:bg-gray-700">Esc</kbd>
          <span>关闭</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import {
  UserGroupIcon,
  SettingsIcon,
  PlugInIcon,
  DocsIcon,
  UserCircleIcon,
} from '@/icons'

const router = useRouter()

// 响应式数据
const searchQuery = ref('')
const showResults = ref(false)
const selectedIndex = ref(0)

// 模拟搜索数据
const searchData = [
  // 管理员
  { id: 1, type: 'admin', title: 'admin', description: '系统管理员', path: '/auth/admin' },
  { id: 2, type: 'admin', title: 'test', description: '测试用户', path: '/auth/admin' },
  
  // 角色
  { id: 3, type: 'role', title: '超级管理员', description: '拥有所有权限的角色', path: '/auth/role' },
  { id: 4, type: 'role', title: '普通用户', description: '基础用户角色', path: '/auth/role' },
  { id: 5, type: 'role', title: '运营人员', description: '运营相关权限', path: '/auth/role' },
  
  // 权限
  { id: 6, type: 'permission', title: '用户管理', description: '管理系统用户', path: '/auth/permission' },
  { id: 7, type: 'permission', title: '角色管理', description: '管理系统角色', path: '/auth/permission' },
  { id: 8, type: 'permission', title: '权限管理', description: '管理系统权限', path: '/auth/permission' },
  
  // 日志
  { id: 9, type: 'log', title: '操作日志', description: '查看系统操作记录', path: '/auth/logs' },
  { id: 10, type: 'log', title: '登录日志', description: '查看用户登录记录', path: '/auth/logs' },
  
  // 个人资料
  { id: 11, type: 'profile', title: '个人资料', description: '管理个人信息', path: '/auth/profile' },
]

// 搜索结果
const searchResults = computed(() => {
  if (!searchQuery.value.trim()) {
    return []
  }
  
  const query = searchQuery.value.toLowerCase()
  return searchData.filter(item =>
    item.title.toLowerCase().includes(query) ||
    item.description.toLowerCase().includes(query)
  ).slice(0, 8) // 限制结果数量
})

// 搜索处理
const handleSearch = () => {
  selectedIndex.value = 0
  showResults.value = true
}

// 清除搜索
const clearSearch = () => {
  searchQuery.value = ''
  showResults.value = false
  selectedIndex.value = 0
}

// 隐藏结果
const hideResults = () => {
  showResults.value = false
}

// 导航结果
const navigateResults = (direction) => {
  if (searchResults.value.length === 0) return
  
  selectedIndex.value += direction
  if (selectedIndex.value < 0) {
    selectedIndex.value = searchResults.value.length - 1
  } else if (selectedIndex.value >= searchResults.value.length) {
    selectedIndex.value = 0
  }
}

// 选择结果
const selectResult = (index) => {
  if (searchResults.value[index]) {
    const result = searchResults.value[index]
    router.push(result.path)
    clearSearch()
  }
}

// 获取类型图标
const getTypeIcon = (type) => {
  const icons = {
    admin: UserGroupIcon,
    role: SettingsIcon,
    permission: PlugInIcon,
    log: DocsIcon,
    profile: UserCircleIcon,
  }
  return icons[type] || UserGroupIcon
}

// 获取类型样式
const getTypeStyle = (type) => {
  const styles = {
    admin: 'bg-blue-50 text-blue-600 dark:bg-blue-500/15 dark:text-blue-400',
    role: 'bg-green-50 text-green-600 dark:bg-green-500/15 dark:text-green-400',
    permission: 'bg-purple-50 text-purple-600 dark:bg-purple-500/15 dark:text-purple-400',
    log: 'bg-orange-50 text-orange-600 dark:bg-orange-500/15 dark:text-orange-400',
    profile: 'bg-gray-50 text-gray-600 dark:bg-gray-500/15 dark:text-gray-400',
  }
  return styles[type] || styles.admin
}

// 获取类型标签样式
const getTypeBadgeStyle = (type) => {
  const styles = {
    admin: 'bg-blue-100 text-blue-700 dark:bg-blue-500/15 dark:text-blue-400',
    role: 'bg-green-100 text-green-700 dark:bg-green-500/15 dark:text-green-400',
    permission: 'bg-purple-100 text-purple-700 dark:bg-purple-500/15 dark:text-purple-400',
    log: 'bg-orange-100 text-orange-700 dark:bg-orange-500/15 dark:text-orange-400',
    profile: 'bg-gray-100 text-gray-700 dark:bg-gray-500/15 dark:text-gray-400',
  }
  return styles[type] || styles.admin
}

// 获取类型标签
const getTypeLabel = (type) => {
  const labels = {
    admin: '管理员',
    role: '角色',
    permission: '权限',
    log: '日志',
    profile: '资料',
  }
  return labels[type] || '未知'
}

// 点击外部关闭
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showResults.value = false
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
/* 自定义样式 */
</style>
