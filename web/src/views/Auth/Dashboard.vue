<template>
  <AuthLayout>
    <PageBreadcrumb :pageTitle="currentPageTitle" />
    
    <div class="grid grid-cols-12 gap-4 md:gap-6">
      <!-- 统计卡片 -->
      <div class="col-span-12 space-y-6 xl:col-span-8">
        <div class="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
          <!-- 管理员总数 -->
          <div class="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-800 dark:bg-white/[0.03]">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-600 dark:text-gray-400">管理员总数</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white/90">{{ stats.adminCount }}</p>
              </div>
              <div class="flex h-12 w-12 items-center justify-center rounded-full bg-blue-50 dark:bg-blue-500/15">
                <svg class="h-6 w-6 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- 角色总数 -->
          <div class="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-800 dark:bg-white/[0.03]">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-600 dark:text-gray-400">角色总数</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white/90">{{ stats.roleCount }}</p>
              </div>
              <div class="flex h-12 w-12 items-center justify-center rounded-full bg-green-50 dark:bg-green-500/15">
                <svg class="h-6 w-6 text-green-600 dark:text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.196-2.121M9 6a3 3 0 106 0 3 3 0 00-6 0zm9 13h4v-2a3 3 0 00-3-3m-4-8a5 5 0 11-10 0 5 5 0 0110 0z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- 权限总数 -->
          <div class="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-800 dark:bg-white/[0.03]">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-600 dark:text-gray-400">权限总数</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white/90">{{ stats.permissionCount }}</p>
              </div>
              <div class="flex h-12 w-12 items-center justify-center rounded-full bg-purple-50 dark:bg-purple-500/15">
                <svg class="h-6 w-6 text-purple-600 dark:text-purple-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
              </div>
            </div>
          </div>

          <!-- 在线用户 -->
          <div class="rounded-2xl border border-gray-200 bg-white p-6 dark:border-gray-800 dark:bg-white/[0.03]">
            <div class="flex items-center justify-between">
              <div>
                <p class="text-sm font-medium text-gray-600 dark:text-gray-400">在线用户</p>
                <p class="text-2xl font-bold text-gray-900 dark:text-white/90">{{ stats.onlineCount }}</p>
              </div>
              <div class="flex h-12 w-12 items-center justify-center rounded-full bg-orange-50 dark:bg-orange-500/15">
                <svg class="h-6 w-6 text-orange-600 dark:text-orange-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.636 18.364a9 9 0 010-12.728m12.728 0a9 9 0 010 12.728m-9.9-2.829a5 5 0 010-7.07m7.072 0a5 5 0 010 7.07M13 12a1 1 0 11-2 0 1 1 0 012 0z" />
                </svg>
              </div>
            </div>
          </div>
        </div>

        <!-- 最近活动 -->
        <ComponentCard title="最近活动">
          <div class="space-y-4">
            <div v-for="activity in recentActivities" :key="activity.id" class="flex items-center gap-4 p-4 rounded-lg bg-gray-50 dark:bg-gray-800/50">
              <div class="flex h-10 w-10 items-center justify-center rounded-full bg-blue-50 dark:bg-blue-500/15">
                <svg class="h-5 w-5 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                </svg>
              </div>
              <div class="flex-1">
                <p class="text-sm font-medium text-gray-900 dark:text-white/90">{{ activity.action }}</p>
                <p class="text-xs text-gray-500 dark:text-gray-400">{{ activity.user }} • {{ activity.time }}</p>
              </div>
            </div>
          </div>
        </ComponentCard>
      </div>

      <!-- 右侧面板 -->
      <div class="col-span-12 xl:col-span-4">
        <!-- 快速操作 -->
        <ComponentCard title="快速操作">
          <div class="grid grid-cols-2 gap-4">
            <router-link
              to="/auth/admin"
              class="flex flex-col items-center p-4 rounded-lg border border-gray-200 hover:border-blue-300 hover:bg-blue-50 transition-colors dark:border-gray-700 dark:hover:border-blue-600 dark:hover:bg-blue-500/10"
            >
              <svg class="h-8 w-8 text-blue-600 dark:text-blue-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
              </svg>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">管理员</span>
            </router-link>

            <router-link
              to="/auth/role"
              class="flex flex-col items-center p-4 rounded-lg border border-gray-200 hover:border-green-300 hover:bg-green-50 transition-colors dark:border-gray-700 dark:hover:border-green-600 dark:hover:bg-green-500/10"
            >
              <svg class="h-8 w-8 text-green-600 dark:text-green-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.196-2.121M9 6a3 3 0 106 0 3 3 0 00-6 0zm9 13h4v-2a3 3 0 00-3-3m-4-8a5 5 0 11-10 0 5 5 0 0110 0z" />
              </svg>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">角色</span>
            </router-link>

            <router-link
              to="/auth/permission"
              class="flex flex-col items-center p-4 rounded-lg border border-gray-200 hover:border-purple-300 hover:bg-purple-50 transition-colors dark:border-gray-700 dark:hover:border-purple-600 dark:hover:bg-purple-500/10"
            >
              <svg class="h-8 w-8 text-purple-600 dark:text-purple-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">权限</span>
            </router-link>

            <router-link
              to="/auth/logs"
              class="flex flex-col items-center p-4 rounded-lg border border-gray-200 hover:border-orange-300 hover:bg-orange-50 transition-colors dark:border-gray-700 dark:hover:border-orange-600 dark:hover:bg-orange-500/10"
            >
              <svg class="h-8 w-8 text-orange-600 dark:text-orange-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span class="text-sm font-medium text-gray-700 dark:text-gray-300">日志</span>
            </router-link>
          </div>
        </ComponentCard>

        <!-- 系统状态 -->
        <ComponentCard title="系统状态" class="mt-6">
          <div class="space-y-4">
            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600 dark:text-gray-400">CPU使用率</span>
              <span class="text-sm font-medium text-gray-900 dark:text-white/90">{{ systemStatus.cpu }}%</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
              <div class="bg-blue-600 h-2 rounded-full" :style="{ width: systemStatus.cpu + '%' }"></div>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600 dark:text-gray-400">内存使用率</span>
              <span class="text-sm font-medium text-gray-900 dark:text-white/90">{{ systemStatus.memory }}%</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
              <div class="bg-green-600 h-2 rounded-full" :style="{ width: systemStatus.memory + '%' }"></div>
            </div>

            <div class="flex items-center justify-between">
              <span class="text-sm text-gray-600 dark:text-gray-400">磁盘使用率</span>
              <span class="text-sm font-medium text-gray-900 dark:text-white/90">{{ systemStatus.disk }}%</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2 dark:bg-gray-700">
              <div class="bg-orange-600 h-2 rounded-full" :style="{ width: systemStatus.disk + '%' }"></div>
            </div>
          </div>
        </ComponentCard>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import AuthLayout from '@/components/auth/AuthLayout.vue'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import ComponentCard from '@/components/common/ComponentCard.vue'

const currentPageTitle = ref('权限管理总览')

// 统计数据
const stats = ref({
  adminCount: 12,
  roleCount: 5,
  permissionCount: 28,
  onlineCount: 8
})

// 最近活动
const recentActivities = ref([
  {
    id: 1,
    action: '管理员 admin 登录系统',
    user: 'admin',
    time: '2分钟前'
  },
  {
    id: 2,
    action: '创建了新角色 "编辑员"',
    user: 'admin',
    time: '10分钟前'
  },
  {
    id: 3,
    action: '修改了用户 "test" 的权限',
    user: 'admin',
    time: '1小时前'
  },
  {
    id: 4,
    action: '删除了角色 "临时用户"',
    user: 'admin',
    time: '2小时前'
  }
])

// 系统状态
const systemStatus = ref({
  cpu: 45,
  memory: 68,
  disk: 32
})

onMounted(() => {
  // 这里可以调用API获取真实数据
  console.log('Auth Dashboard mounted')
})
</script>
