<template>
  <AuthLayout>
    <PageBreadcrumb pageTitle="服务监控" />
    <div class="space-y-6">

      <div class="flex items-center justify-between">
        <h2 class="text-lg font-bold text-gray-800 dark:text-white">运行时监控</h2>
        <button 
          @click="fetchServerInfo" 
          class="flex items-center gap-2 px-4 py-2 bg-blue-50 text-blue-600 rounded-xl hover:bg-blue-100 transition shadow-sm dark:bg-blue-900/30 dark:text-blue-400"
          :class="{'opacity-50 cursor-not-allowed': loading}"
          :disabled="loading"
        >
          <svg class="w-4 h-4" :class="{'animate-spin': loading}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/></svg>
          重新获取
        </button>
      </div>

      <div v-if="loading && !serverInfo" class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div class="animate-pulse bg-white/50 h-64 rounded-2xl"></div>
        <div class="animate-pulse bg-white/50 h-64 rounded-2xl"></div>
      </div>

      <div v-if="serverInfo" class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        
        <!-- 环境信息 -->
        <AuthCard title="系统环境与 Go 运行态">
          <div class="space-y-4">
            <div class="flex justify-between items-center p-3 sm:p-4 rounded-xl bg-gray-50 border border-gray-100 dark:bg-gray-800/50 dark:border-gray-700">
              <span class="text-sm text-gray-500 font-medium">操作系统</span>
              <span class="font-mono text-sm text-gray-900 dark:text-white font-semibold">{{ serverInfo.os.os }} ({{ serverInfo.os.arch }})</span>
            </div>
            <div class="flex justify-between items-center p-3 sm:p-4 rounded-xl bg-gray-50 border border-gray-100 dark:bg-gray-800/50 dark:border-gray-700">
              <span class="text-sm text-gray-500 font-medium">逻辑核心数</span>
              <span class="font-mono text-sm text-gray-900 dark:text-white font-semibold">{{ serverInfo.os.num_cpu }} Cores</span>
            </div>
            <div class="flex justify-between items-center p-3 sm:p-4 rounded-xl bg-gray-50 border border-gray-100 dark:bg-gray-800/50 dark:border-gray-700">
              <span class="text-sm text-gray-500 font-medium">Go 语言版本</span>
              <span class="font-mono text-sm px-2 py-0.5 bg-blue-100 text-blue-700 rounded-lg">{{ serverInfo.os.go_version }}</span>
            </div>
            <div class="flex justify-between items-center p-3 sm:p-4 rounded-xl bg-gray-50 border border-gray-100 dark:bg-gray-800/50 dark:border-gray-700">
              <span class="text-sm text-gray-500 font-medium">运行时长</span>
              <span class="font-mono text-sm text-green-600 font-semibold">{{ serverInfo.os.uptime }}</span>
            </div>
          </div>
        </AuthCard>

        <!-- 内存信息 -->
        <AuthCard title="内存占用与性能池">
          <div class="space-y-4">
            <div class="relative p-5 rounded-2xl bg-gradient-to-br from-indigo-50 to-blue-50 border border-blue-100 dark:from-indigo-900/30 dark:to-blue-900/30 dark:border-blue-800">
              <div class="flex justify-between mb-2">
                <span class="text-sm font-semibold text-indigo-800 dark:text-indigo-200">已分配内存量 (Alloc)</span>
                <span class="font-mono font-bold text-indigo-600 dark:text-indigo-400">{{ serverInfo.mem.alloc }}</span>
              </div>
              <div class="w-full bg-indigo-200/50 rounded-full h-2.5 dark:bg-indigo-900/50">
                <div class="bg-indigo-500 h-2.5 rounded-full" style="width: 45%"></div>
              </div>
              <p class="text-xs text-indigo-500/80 mt-2">当前存活对象累积占用，代表应用实际使用的常驻内存。</p>
            </div>

            <div class="flex justify-between items-center p-3 sm:p-4 rounded-xl bg-gray-50 border border-gray-100 dark:bg-gray-800/50 dark:border-gray-700">
              <span class="text-sm text-gray-500 font-medium">向系统申请总内存 (Sys)</span>
              <span class="font-mono text-sm text-gray-900 dark:text-white">{{ serverInfo.mem.sys }}</span>
            </div>
            <div class="flex justify-between items-center p-3 sm:p-4 rounded-xl bg-gray-50 border border-gray-100 dark:bg-gray-800/50 dark:border-gray-700">
              <span class="text-sm text-gray-500 font-medium">累计分配内存 (Total Alloc)</span>
              <span class="font-mono text-sm text-gray-900 dark:text-white">{{ serverInfo.mem.total_alloc }}</span>
            </div>
            
            <div class="grid grid-cols-2 gap-4">
              <div class="p-4 rounded-xl border border-gray-200 bg-white shadow-sm flex flex-col justify-center items-center dark:bg-gray-800 dark:border-gray-700">
                <span class="text-3xl font-black text-rose-500 font-mono">{{ serverInfo.mem.num_gc }}</span>
                <span class="text-xs text-gray-500 tracking-widest mt-1 uppercase">GC 执行次数</span>
              </div>
              <div class="p-4 rounded-xl border border-gray-200 bg-white shadow-sm flex flex-col justify-center items-center dark:bg-gray-800 dark:border-gray-700">
                <span class="text-3xl font-black text-teal-500 font-mono">{{ serverInfo.os.num_goroutine }}</span>
                <span class="text-xs text-gray-500 tracking-widest mt-1 uppercase">当前携程数</span>
              </div>
            </div>
          </div>
        </AuthCard>

      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { AuthLayout, AuthCard } from '@/components/auth'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import { monitorApi } from '@/api/system/monitor'
import ToastAlert from '@/composables/ToastAlert'

const serverInfo = ref<any>(null)
const loading = ref(false)
let timer: any = null

const fetchServerInfo = async () => {
  loading.value = true
  try {
    const { data: res } = await monitorApi.getServerInfo()
    if (res.code === 200) {
      serverInfo.value = res.data
    } else {
      ToastAlert.error({ title: '获取状态失败', message: res.message })
    }
  } catch (error) {
    console.error('获取系统状态失败:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchServerInfo()
  // 每 5 秒刷新一次
  timer = setInterval(fetchServerInfo, 5000)
})

onUnmounted(() => {
  if (timer) clearInterval(timer)
})
</script>
