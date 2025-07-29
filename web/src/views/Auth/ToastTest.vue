<template>
  <div class="p-6 space-y-6">
    <div class="bg-white rounded-lg shadow-lg p-6 dark:bg-gray-800">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        Toast Z-Index 测试页面
      </h1>
      
      <!-- Toast测试按钮 -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <button
          @click="showSuccessToast"
          class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors"
        >
          成功提示
        </button>
        <button
          @click="showErrorToast"
          class="px-4 py-2 bg-red-500 text-white rounded-lg hover:bg-red-600 transition-colors"
        >
          错误提示
        </button>
        <button
          @click="showWarningToast"
          class="px-4 py-2 bg-yellow-500 text-white rounded-lg hover:bg-yellow-600 transition-colors"
        >
          警告提示
        </button>
        <button
          @click="showConfirmToast"
          class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
        >
          确认对话框
        </button>
      </div>

      <!-- 高z-index元素测试 -->
      <div class="space-y-6">
        <h2 class="text-xl font-semibold text-gray-900 dark:text-white">
          高Z-Index元素测试
        </h2>
        
        <!-- 模拟模态框 -->
        <div class="relative">
          <button
            @click="showModal = !showModal"
            class="px-4 py-2 bg-purple-500 text-white rounded-lg hover:bg-purple-600 transition-colors"
          >
            显示模态框 (z-99999)
          </button>
          
          <div
            v-if="showModal"
            class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-99999"
            @click="showModal = false"
          >
            <div
              class="bg-white rounded-lg p-6 max-w-md mx-4 dark:bg-gray-800"
              @click.stop
            >
              <h3 class="text-lg font-semibold mb-4 text-gray-900 dark:text-white">
                模态框测试
              </h3>
              <p class="text-gray-600 dark:text-gray-300 mb-4">
                这是一个z-index为99999的模态框。Toast应该能够显示在它上面。
              </p>
              <div class="flex gap-2">
                <button
                  @click="showSuccessToast"
                  class="px-3 py-1 bg-green-500 text-white rounded text-sm hover:bg-green-600"
                >
                  测试Toast
                </button>
                <button
                  @click="showModal = false"
                  class="px-3 py-1 bg-gray-500 text-white rounded text-sm hover:bg-gray-600"
                >
                  关闭
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- 模拟下拉菜单 -->
        <div class="relative inline-block">
          <button
            @click="showDropdown = !showDropdown"
            class="px-4 py-2 bg-indigo-500 text-white rounded-lg hover:bg-indigo-600 transition-colors"
          >
            下拉菜单 (z-50)
          </button>
          
          <div
            v-if="showDropdown"
            class="absolute top-full left-0 mt-1 w-48 bg-white border border-gray-200 rounded-lg shadow-lg z-50 dark:bg-gray-800 dark:border-gray-600"
          >
            <div class="p-2">
              <button
                @click="showWarningToast"
                class="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded dark:text-gray-300 dark:hover:bg-gray-700"
              >
                测试Toast
              </button>
              <button
                @click="showDropdown = false"
                class="w-full text-left px-3 py-2 text-sm text-gray-700 hover:bg-gray-100 rounded dark:text-gray-300 dark:hover:bg-gray-700"
              >
                关闭菜单
              </button>
            </div>
          </div>
        </div>

        <!-- 固定定位元素 -->
        <div class="relative">
          <button
            @click="showFixedElement = !showFixedElement"
            class="px-4 py-2 bg-teal-500 text-white rounded-lg hover:bg-teal-600 transition-colors"
          >
            固定元素 (z-9999)
          </button>
          
          <div
            v-if="showFixedElement"
            class="fixed top-20 right-4 bg-orange-500 text-white p-4 rounded-lg shadow-lg z-9999"
          >
            <div class="flex items-center justify-between">
              <span class="text-sm">固定定位元素</span>
              <button
                @click="showFixedElement = false"
                class="ml-2 text-white hover:text-gray-200"
              >
                ×
              </button>
            </div>
            <button
              @click="showErrorToast"
              class="mt-2 px-2 py-1 bg-red-600 text-white rounded text-xs hover:bg-red-700"
            >
              测试Toast
            </button>
          </div>
        </div>
      </div>

      <!-- 测试说明 -->
      <div class="mt-8 p-4 bg-blue-50 rounded-lg dark:bg-blue-900/20">
        <h3 class="text-lg font-semibold text-blue-900 dark:text-blue-300 mb-2">
          测试说明
        </h3>
        <ul class="text-sm text-blue-800 dark:text-blue-200 space-y-1">
          <li>• Toast现在使用最高z-index值 (2147483647)，确保始终在最上层显示</li>
          <li>• 点击各种按钮测试Toast是否能正确显示在所有元素之上</li>
          <li>• 模态框使用z-99999，下拉菜单使用z-50，固定元素使用z-9999</li>
          <li>• Toast应该能够显示在所有这些元素之上</li>
          <li>• 使用Teleport确保Toast挂载到body元素</li>
          <li>• 使用专用CSS类和最大z-index值保证层级优先级</li>
        </ul>
      </div>

      <!-- 实时测试结果 -->
      <div class="mt-6 p-4 bg-green-50 rounded-lg dark:bg-green-900/20">
        <h3 class="text-lg font-semibold text-green-900 dark:text-green-300 mb-2">
          优化效果
        </h3>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
          <div class="space-y-2">
            <h4 class="font-medium text-green-800 dark:text-green-200">技术实现:</h4>
            <ul class="text-green-700 dark:text-green-300 space-y-1">
              <li>✅ 使用Teleport挂载到body</li>
              <li>✅ 最大z-index值 (2147483647)</li>
              <li>✅ 专用CSS工具类</li>
              <li>✅ 指针事件优化</li>
            </ul>
          </div>
          <div class="space-y-2">
            <h4 class="font-medium text-green-800 dark:text-green-200">显示效果:</h4>
            <ul class="text-green-700 dark:text-green-300 space-y-1">
              <li>✅ 始终在最上层显示</li>
              <li>✅ 不被任何组件遮挡</li>
              <li>✅ 背景模糊效果正常</li>
              <li>✅ 动画过渡流畅</li>
            </ul>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ToastAlert from '@/composables/ToastAlert'

// 控制各种元素的显示状态
const showModal = ref(false)
const showDropdown = ref(false)
const showFixedElement = ref(false)

// Toast测试方法
const showSuccessToast = () => {
  ToastAlert.success({
    title: '操作成功',
    message: '这是一个成功提示，应该显示在所有元素之上！',
    duration: 3000
  })
}

const showErrorToast = () => {
  ToastAlert.error({
    title: '操作失败',
    message: '这是一个错误提示，测试z-index是否正确。',
    duration: 3000
  })
}

const showWarningToast = () => {
  ToastAlert.warning({
    title: '警告提示',
    message: '这是一个警告提示，检查是否被其他元素遮挡。',
    duration: 3000
  })
}

const showConfirmToast = async () => {
  const confirmed = await ToastAlert.confirm({
    title: '确认操作',
    message: '这是一个确认对话框，应该显示在最上层。您确定要继续吗？'
  })
  
  if (confirmed) {
    ToastAlert.success({
      title: '已确认',
      message: '您点击了确认按钮！'
    })
  } else {
    ToastAlert.info({
      title: '已取消',
      message: '您点击了取消按钮！'
    })
  }
}
</script>
