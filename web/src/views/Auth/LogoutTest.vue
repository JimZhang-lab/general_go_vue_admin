<template>
  <div class="p-6 space-y-6">
    <div class="bg-white rounded-lg shadow-lg p-6 dark:bg-gray-800">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        登出功能测试页面
      </h1>
      
      <!-- 测试按钮区域 -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
        <!-- 标准登出测试 -->
        <div class="p-4 border border-gray-200 rounded-lg dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
            标准登出测试
          </h3>
          <p class="text-gray-600 dark:text-gray-300 mb-4">
            测试完整的登出流程，包括确认对话框、API调用和页面跳转。
          </p>
          <button
            @click="testStandardLogout"
            class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
          >
            测试标准登出
          </button>
        </div>

        <!-- 确认对话框测试 -->
        <div class="p-4 border border-gray-200 rounded-lg dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
            确认对话框测试
          </h3>
          <p class="text-gray-600 dark:text-gray-300 mb-4">
            单独测试确认对话框的显示和交互功能。
          </p>
          <button
            @click="testConfirmDialog"
            class="px-4 py-2 bg-yellow-500 text-white rounded-lg hover:bg-yellow-600 transition-colors"
          >
            测试确认对话框
          </button>
        </div>

        <!-- 取消登出测试 -->
        <div class="p-4 border border-gray-200 rounded-lg dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
            取消登出测试
          </h3>
          <p class="text-gray-600 dark:text-gray-300 mb-4">
            测试用户点击取消按钮时的处理逻辑。
          </p>
          <button
            @click="testCancelLogout"
            class="px-4 py-2 bg-gray-500 text-white rounded-lg hover:bg-gray-600 transition-colors"
          >
            测试取消登出
          </button>
        </div>

        <!-- 路由重定向测试 -->
        <div class="p-4 border border-gray-200 rounded-lg dark:border-gray-700">
          <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
            路由重定向测试
          </h3>
          <p class="text-gray-600 dark:text-gray-300 mb-4">
            测试登出后的页面跳转和路由重定向功能。
          </p>
          <button
            @click="testRouteRedirect"
            class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            测试路由重定向
          </button>
        </div>
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

      <!-- 功能说明 -->
      <div class="mt-8 p-4 bg-blue-50 rounded-lg dark:bg-blue-900/20">
        <h3 class="text-lg font-semibold text-blue-900 dark:text-blue-300 mb-2">
          登出功能说明
        </h3>
        <ul class="text-sm text-blue-800 dark:text-blue-200 space-y-1">
          <li>• <strong>确认对话框</strong>: 显示警告样式的确认对话框，包含"确定"和"取消"按钮</li>
          <li>• <strong>API调用</strong>: 调用后端 /api/logout 接口记录登出操作</li>
          <li>• <strong>数据清理</strong>: 清除localStorage中的所有用户数据</li>
          <li>• <strong>页面跳转</strong>: 登出成功后跳转到 /adminLogin 页面</li>
          <li>• <strong>错误处理</strong>: API失败时仍能正常清理数据并跳转</li>
        </ul>
      </div>

      <!-- 路由重定向说明 -->
      <div class="mt-6 p-4 bg-green-50 rounded-lg dark:bg-green-900/20">
        <h3 class="text-lg font-semibold text-green-900 dark:text-green-300 mb-2">
          路由重定向优化
        </h3>
        <ul class="text-sm text-green-800 dark:text-green-200 space-y-1">
          <li>• <strong>根路径重定向</strong>: "/" 现在重定向到 "/adminLogin"</li>
          <li>• <strong>已登录重定向</strong>: 已登录用户访问登录页时重定向到 "/auth/dashboard"</li>
          <li>• <strong>未登录重定向</strong>: 未登录用户访问受保护页面时重定向到 "/adminLogin"</li>
          <li>• <strong>登出跳转</strong>: 登出成功后跳转到 "/adminLogin" 页面</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import ToastAlert from '@/composables/ToastAlert'

const router = useRouter()

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

// 测试标准登出流程
const testStandardLogout = async () => {
  try {
    // 显示确认对话框
    const confirmed = await ToastAlert.confirm({
      title: '确认退出',
      message: '您确定要退出系统吗？',
      variant: 'warning'
    })

    if (!confirmed) {
      addTestResult('用户取消了登出操作', true)
      return
    }

    // 模拟API调用
    addTestResult('用户确认登出，开始执行登出流程', true)
    
    // 这里不实际执行登出，只是测试
    ToastAlert.success({
      title: '测试成功',
      message: '登出确认对话框工作正常'
    })
    
    addTestResult('登出确认对话框测试通过', true)
  } catch (error) {
    addTestResult('登出测试失败: ' + error, false)
  }
}

// 测试确认对话框
const testConfirmDialog = async () => {
  try {
    const confirmed = await ToastAlert.confirm({
      title: '测试确认对话框',
      message: '这是一个测试确认对话框，请点击确定或取消。',
      variant: 'info'
    })

    if (confirmed) {
      addTestResult('用户点击了确定按钮', true)
      ToastAlert.success({
        title: '确定',
        message: '您点击了确定按钮'
      })
    } else {
      addTestResult('用户点击了取消按钮', true)
      ToastAlert.info({
        title: '取消',
        message: '您点击了取消按钮'
      })
    }
  } catch (error) {
    addTestResult('确认对话框测试失败: ' + error, false)
  }
}

// 测试取消登出
const testCancelLogout = async () => {
  try {
    const confirmed = await ToastAlert.confirm({
      title: '取消登出测试',
      message: '请点击取消按钮来测试取消功能。',
      variant: 'warning'
    })

    if (confirmed) {
      addTestResult('用户意外点击了确定（应该点击取消）', false)
    } else {
      addTestResult('取消登出功能测试通过', true)
      ToastAlert.info({
        title: '取消成功',
        message: '登出操作已取消'
      })
    }
  } catch (error) {
    addTestResult('取消登出测试失败: ' + error, false)
  }
}

// 测试路由重定向
const testRouteRedirect = () => {
  addTestResult('开始测试路由重定向', true)
  
  // 测试跳转到根路径（应该重定向到登录页）
  ToastAlert.info({
    title: '路由测试',
    message: '3秒后将测试根路径重定向功能'
  })
  
  setTimeout(() => {
    router.push('/')
    addTestResult('已跳转到根路径，检查是否重定向到登录页', true)
  }, 3000)
}
</script>
