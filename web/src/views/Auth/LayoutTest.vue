<template>
  <AuthLayout>
    <div class="space-y-6">
      <!-- 页面标题 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-2">
          布局测试页面
        </h1>
        <p class="text-gray-600 dark:text-gray-400">
          测试AuthLayout组件的布局是否正确，内容是否被sidebar遮挡
        </p>
      </div>

      <!-- 布局测试区域 -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- 测试卡片1 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              左侧边距测试
            </h2>
            <div class="w-3 h-3 bg-green-500 rounded-full"></div>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            这个卡片应该完全可见，不被sidebar遮挡。
          </p>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            ✅ 如果你能完整看到这段文字，说明布局正常
          </div>
        </div>

        <!-- 测试卡片2 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              响应式测试
            </h2>
            <div class="w-3 h-3 bg-blue-500 rounded-full"></div>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            调整浏览器窗口大小测试响应式布局。
          </p>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            📱 移动端 | 💻 桌面端
          </div>
        </div>

        <!-- 测试卡片3 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <div class="flex items-center justify-between mb-4">
            <h2 class="text-lg font-semibold text-gray-900 dark:text-white">
              Sidebar交互测试
            </h2>
            <div class="w-3 h-3 bg-purple-500 rounded-full"></div>
          </div>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            点击sidebar切换按钮测试展开/收缩功能。
          </p>
          <div class="text-sm text-gray-500 dark:text-gray-400">
            🔄 展开/收缩 | 🖱️ 鼠标悬停
          </div>
        </div>
      </div>

      <!-- 宽度测试区域 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          内容宽度测试
        </h2>
        <div class="bg-gradient-to-r from-blue-500 to-purple-600 rounded-lg p-6 text-white">
          <h3 class="text-xl font-bold mb-2">全宽度内容区域</h3>
          <p class="mb-4">
            这个区域应该占据除sidebar外的全部宽度。当sidebar展开/收缩时，这个区域的宽度应该相应调整。
          </p>
          <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
            <div class="bg-white/20 rounded p-3 text-center">
              <div class="text-2xl font-bold">290px</div>
              <div class="text-sm">展开宽度</div>
            </div>
            <div class="bg-white/20 rounded p-3 text-center">
              <div class="text-2xl font-bold">90px</div>
              <div class="text-sm">收缩宽度</div>
            </div>
            <div class="bg-white/20 rounded p-3 text-center">
              <div class="text-2xl font-bold">Auto</div>
              <div class="text-sm">内容宽度</div>
            </div>
            <div class="bg-white/20 rounded p-3 text-center">
              <div class="text-2xl font-bold">100%</div>
              <div class="text-sm">响应式</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 边距测试 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          边距和间距测试
        </h2>
        <div class="space-y-4">
          <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <span class="text-gray-700 dark:text-gray-300">左边距 (lg:ml-[290px] / lg:ml-[90px])</span>
            <span class="text-green-600 dark:text-green-400 font-semibold">✓ 正常</span>
          </div>
          <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <span class="text-gray-700 dark:text-gray-300">内容padding (p-4 md:p-6)</span>
            <span class="text-green-600 dark:text-green-400 font-semibold">✓ 正常</span>
          </div>
          <div class="flex items-center justify-between p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <span class="text-gray-700 dark:text-gray-300">最大宽度 (max-w-(--breakpoint-2xl))</span>
            <span class="text-green-600 dark:text-green-400 font-semibold">✓ 正常</span>
          </div>
        </div>
      </div>

      <!-- 滚动测试 -->
      <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          滚动测试区域
        </h2>
        <p class="text-gray-600 dark:text-gray-400 mb-4">
          这里有足够的内容来测试页面滚动功能。
        </p>
        <div class="space-y-4">
          <div v-for="i in 10" :key="i" class="p-4 bg-gray-50 dark:bg-gray-700 rounded-lg">
            <h3 class="font-semibold text-gray-900 dark:text-white mb-2">
              测试内容块 {{ i }}
            </h3>
            <p class="text-gray-600 dark:text-gray-400">
              这是第{{ i }}个测试内容块。用于测试页面滚动时的布局稳定性。
              内容应该在sidebar下方正常滚动，不会被遮挡或产生布局问题。
            </p>
          </div>
        </div>
      </div>

      <!-- 底部测试 -->
      <div class="bg-gradient-to-r from-green-500 to-blue-600 rounded-lg shadow p-6 text-white">
        <h2 class="text-xl font-bold mb-2">
          🎉 布局测试完成
        </h2>
        <p class="mb-4">
          如果你能看到这个区域，并且所有内容都没有被sidebar遮挡，说明AuthLayout布局修复成功！
        </p>
        <div class="flex items-center justify-between">
          <div class="text-sm opacity-90">
            测试时间: {{ new Date().toLocaleString() }}
          </div>
          <div class="text-sm opacity-90">
            ✅ 布局正常
          </div>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup>
import { AuthLayout } from '@/components/auth'

// 页面标题
const currentPageTitle = '布局测试'
</script>

<style scoped>
/* 自定义样式 */
</style>
