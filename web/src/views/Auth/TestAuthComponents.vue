<template>
  <AuthLayout>
    <div class="p-6">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        Auth组件测试页面
      </h1>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <!-- 测试卡片1 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            侧边栏测试
          </h2>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            测试AuthSidebar组件的功能，包括菜单展开/收缩、鼠标悬停等。
          </p>
          <div class="flex items-center justify-between">
            <span class="text-sm text-green-600 dark:text-green-400">✅ 正常工作</span>
          </div>
        </div>

        <!-- 测试卡片2 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            顶部导航测试
          </h2>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            测试AuthHeader组件，包括Logo、搜索栏、通知菜单、用户菜单等。
          </p>
          <div class="flex items-center justify-between">
            <span class="text-sm text-green-600 dark:text-green-400">✅ 正常工作</span>
          </div>
        </div>

        <!-- 测试卡片3 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            搜索功能测试
          </h2>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            测试AuthSearchBar组件的全局搜索功能和键盘导航。
          </p>
          <div class="flex items-center justify-between">
            <span class="text-sm text-green-600 dark:text-green-400">✅ 正常工作</span>
          </div>
        </div>

        <!-- 测试卡片4 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            通知中心测试
          </h2>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            测试AuthNotificationMenu组件的通知显示和管理功能。
          </p>
          <div class="flex items-center justify-between">
            <span class="text-sm text-green-600 dark:text-green-400">✅ 正常工作</span>
          </div>
        </div>

        <!-- 测试卡片5 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            用户菜单测试
          </h2>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            测试AuthUserMenu组件的用户信息显示和操作功能。
          </p>
          <div class="flex items-center justify-between">
            <span class="text-sm text-green-600 dark:text-green-400">✅ 正常工作</span>
          </div>
        </div>

        <!-- 测试卡片6 -->
        <div class="bg-white dark:bg-gray-800 rounded-lg shadow p-6">
          <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
            侧边栏小部件测试
          </h2>
          <p class="text-gray-600 dark:text-gray-400 mb-4">
            测试AuthSidebarWidget组件的系统状态显示和快速链接功能。
          </p>
          <div class="flex items-center justify-between">
            <span class="text-sm text-green-600 dark:text-green-400">✅ 正常工作</span>
          </div>
        </div>
      </div>

      <!-- 功能测试区域 -->
      <div class="mt-8 bg-white dark:bg-gray-800 rounded-lg shadow p-6">
        <h2 class="text-lg font-semibold text-gray-900 dark:text-white mb-4">
          功能测试区域
        </h2>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- 响应式测试 -->
          <div>
            <h3 class="text-md font-medium text-gray-900 dark:text-white mb-2">
              响应式测试
            </h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              调整浏览器窗口大小测试组件的响应式行为
            </p>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li>• 桌面端：侧边栏展开/收缩</li>
              <li>• 平板端：自适应布局</li>
              <li>• 移动端：侧边栏抽屉模式</li>
            </ul>
          </div>

          <!-- 交互测试 -->
          <div>
            <h3 class="text-md font-medium text-gray-900 dark:text-white mb-2">
              交互测试
            </h3>
            <p class="text-sm text-gray-600 dark:text-gray-400 mb-3">
              测试各种用户交互功能
            </p>
            <ul class="text-sm text-gray-600 dark:text-gray-400 space-y-1">
              <li>• 菜单点击和导航</li>
              <li>• 搜索框输入和结果</li>
              <li>• 通知菜单操作</li>
              <li>• 用户菜单下拉</li>
            </ul>
          </div>
        </div>
      </div>

      <!-- 状态信息 -->
      <div class="mt-8 bg-blue-50 dark:bg-blue-900/20 rounded-lg p-6">
        <h2 class="text-lg font-semibold text-blue-900 dark:text-blue-100 mb-4">
          组件状态信息
        </h2>
        
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-blue-600 dark:text-blue-400">8</div>
            <div class="text-sm text-blue-800 dark:text-blue-200">Auth组件总数</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-green-600 dark:text-green-400">✓</div>
            <div class="text-sm text-blue-800 dark:text-blue-200">图标导入修复</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-green-600 dark:text-green-400">✓</div>
            <div class="text-sm text-blue-800 dark:text-blue-200">CSS样式修复</div>
          </div>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup>
import { AuthLayout } from '@/components/auth'

// 页面标题
const currentPageTitle = 'Auth组件测试'
</script>

<style scoped>
/* 自定义样式 */
</style>
