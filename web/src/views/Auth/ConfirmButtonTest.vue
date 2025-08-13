<template>
  <div class="p-6 space-y-6">
    <div class="bg-white rounded-lg shadow-lg p-6 dark:bg-gray-800">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        确认按钮测试页面
      </h1>
      
      <!-- 测试按钮区域 -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <button
          @click="testWarningConfirm"
          class="px-4 py-2 bg-orange-500 text-white rounded-lg hover:bg-orange-600 transition-colors"
        >
          警告确认对话框
        </button>
        
        <button
          @click="testErrorConfirm"
          class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
        >
          错误确认对话框
        </button>
        
        <button
          @click="testSuccessConfirm"
          class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
        >
          成功确认对话框
        </button>
        
        <button
          @click="testInfoConfirm"
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          信息确认对话框
        </button>
      </div>

      <!-- 登出测试 -->
      <div class="mb-8 p-4 border border-gray-200 rounded-lg dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
          登出确认测试
        </h3>
        <p class="text-gray-600 dark:text-gray-300 mb-4">
          测试与AuthUserMenu中相同的登出确认对话框
        </p>
        <button
          @click="testLogoutConfirm"
          class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors"
        >
          测试登出确认
        </button>
      </div>

      <!-- 测试结果显示 -->
      <div class="bg-gray-50 rounded-lg p-4 dark:bg-gray-900">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
          测试结果
        </h3>
        <div class="space-y-2">
          <div v-for="(result, index) in testResults" :key="index" class="flex items-center gap-2">
            <span :class="result.success ? 'text-green-600' : 'text-red-600'">
              {{ result.success ? '✅' : '❌' }}
            </span>
            <span class="text-gray-700 dark:text-gray-300">{{ result.message }}</span>
            <span class="text-xs text-gray-500">{{ result.timestamp }}</span>
          </div>
        </div>
      </div>

      <!-- 问题排查说明 -->
      <div class="mt-8 p-4 bg-yellow-50 rounded-lg dark:bg-yellow-900/20">
        <h3 class="text-lg font-semibold text-yellow-900 dark:text-yellow-300 mb-2">
          确认按钮问题排查
        </h3>
        <ul class="text-sm text-yellow-800 dark:text-yellow-200 space-y-1">
          <li>• <strong>样式问题</strong>: 检查确认按钮的CSS样式是否正确应用</li>
          <li>• <strong>颜色定义</strong>: 确保使用的颜色类在Tailwind中有定义</li>
          <li>• <strong>z-index层级</strong>: 确认按钮是否被其他元素遮挡</li>
          <li>• <strong>显示逻辑</strong>: 检查showCancel属性和按钮显示条件</li>
          <li>• <strong>事件绑定</strong>: 确认按钮的点击事件是否正确绑定</li>
        </ul>
      </div>

      <!-- 修复说明 -->
      <div class="mt-6 p-4 bg-green-50 rounded-lg dark:bg-green-900/20">
        <h3 class="text-lg font-semibold text-green-900 dark:text-green-300 mb-2">
          修复内容
        </h3>
        <ul class="text-sm text-green-800 dark:text-green-200 space-y-1">
          <li>• <strong>颜色修复</strong>: 将bg-primary-500改为bg-blue-600等标准颜色</li>
          <li>• <strong>动态样式</strong>: 根据variant类型设置不同的确认按钮颜色</li>
          <li>• <strong>暗色模式</strong>: 添加dark模式下的按钮样式</li>
          <li>• <strong>悬停效果</strong>: 添加hover状态的颜色变化</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ToastAlert from '@/composables/ToastAlert'

// 测试结果记录
const testResults = ref<Array<{
  message: string
  success: boolean
  timestamp: string
}>>([])

// 添加测试结果
const addTestResult = (message: string, success: boolean) => {
  testResults.value.unshift({
    message,
    success,
    timestamp: new Date().toLocaleTimeString()
  })
}

// 测试警告确认对话框
const testWarningConfirm = async () => {
  try {
    const confirmed = await ToastAlert.confirm({
      title: '警告确认',
      message: '这是一个警告样式的确认对话框，请检查是否有确定和取消按钮。',
      variant: 'warning'
    })

    if (confirmed) {
      addTestResult('警告确认对话框 - 用户点击确定', true)
      ToastAlert.success({
        title: '确定',
        message: '您点击了确定按钮'
      })
    } else {
      addTestResult('警告确认对话框 - 用户点击取消', true)
    }
  } catch (error) {
    addTestResult('警告确认对话框测试失败: ' + error, false)
  }
}

// 测试错误确认对话框
const testErrorConfirm = async () => {
  try {
    const confirmed = await ToastAlert.confirm({
      title: '错误确认',
      message: '这是一个错误样式的确认对话框。',
      variant: 'error'
    })

    if (confirmed) {
      addTestResult('错误确认对话框 - 用户点击确定', true)
    } else {
      addTestResult('错误确认对话框 - 用户点击取消', true)
    }
  } catch (error) {
    addTestResult('错误确认对话框测试失败: ' + error, false)
  }
}

// 测试成功确认对话框
const testSuccessConfirm = async () => {
  try {
    const confirmed = await ToastAlert.confirm({
      title: '成功确认',
      message: '这是一个成功样式的确认对话框。',
      variant: 'success'
    })

    if (confirmed) {
      addTestResult('成功确认对话框 - 用户点击确定', true)
    } else {
      addTestResult('成功确认对话框 - 用户点击取消', true)
    }
  } catch (error) {
    addTestResult('成功确认对话框测试失败: ' + error, false)
  }
}

// 测试信息确认对话框
const testInfoConfirm = async () => {
  try {
    const confirmed = await ToastAlert.confirm({
      title: '信息确认',
      message: '这是一个信息样式的确认对话框。',
      variant: 'info'
    })

    if (confirmed) {
      addTestResult('信息确认对话框 - 用户点击确定', true)
    } else {
      addTestResult('信息确认对话框 - 用户点击取消', true)
    }
  } catch (error) {
    addTestResult('信息确认对话框测试失败: ' + error, false)
  }
}

// 测试登出确认（与AuthUserMenu相同）
const testLogoutConfirm = async () => {
  try {
    const confirmed = await ToastAlert.confirm({
      title: '确认退出',
      message: '您确定要退出系统吗？',
      variant: 'warning'
    })

    if (confirmed) {
      addTestResult('登出确认对话框 - 用户点击确定（模拟登出）', true)
      ToastAlert.success({
        title: '退出成功',
        message: '您已安全退出系统（模拟）'
      })
    } else {
      addTestResult('登出确认对话框 - 用户点击取消', true)
    }
  } catch (error) {
    addTestResult('登出确认对话框测试失败: ' + error, false)
  }
}
</script>
