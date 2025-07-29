<template>
  <div class="p-8 space-y-4">
    <h1 class="text-2xl font-bold mb-6">Toast 测试页面</h1>
    
    <div class="grid grid-cols-2 gap-4">
      <button 
        @click="showSuccess"
        class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
      >
        显示成功提示
      </button>
      
      <button 
        @click="showError"
        class="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
      >
        显示错误提示
      </button>
      
      <button 
        @click="showMultipleErrors"
        class="px-4 py-2 bg-orange-500 text-white rounded hover:bg-orange-600"
      >
        显示多个错误（快速连续）
      </button>
      
      <button 
        @click="showLongMessage"
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        显示长消息
      </button>
      
      <button 
        @click="showWarning"
        class="px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600"
      >
        显示警告
      </button>
      
      <button 
        @click="showInfo"
        class="px-4 py-2 bg-indigo-500 text-white rounded hover:bg-indigo-600"
      >
        显示信息
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import ToastAlert from '@/composables/ToastAlert'

const showSuccess = () => {
  ToastAlert.success({
    title: '操作成功',
    message: '数据已成功保存到数据库中'
  })
}

const showError = () => {
  ToastAlert.error({
    title: '操作失败',
    message: '网络连接超时，请检查您的网络连接后重试'
  })
}

const showMultipleErrors = () => {
  // 快速连续显示多个错误，测试错误收集功能
  ToastAlert.error({
    title: '验证失败',
    message: '用户名不能为空'
  })
  
  setTimeout(() => {
    ToastAlert.error({
      title: '验证失败',
      message: '密码长度至少为8位'
    })
  }, 100)
  
  setTimeout(() => {
    ToastAlert.error({
      title: '验证失败',
      message: '邮箱格式不正确'
    })
  }, 200)
  
  setTimeout(() => {
    ToastAlert.error({
      title: '验证失败',
      message: '手机号码格式不正确'
    })
  }, 300)
}

const showLongMessage = () => {
  ToastAlert.error({
    title: '系统错误',
    message: '这是一个非常长的错误消息，用来测试弹窗的自适应大小功能。当消息内容很长时，弹窗应该能够自动调整大小以适应内容，同时保持良好的用户体验。这个消息包含了很多详细的错误信息，比如错误代码、错误原因、可能的解决方案等等。'
  })
}

const showWarning = () => {
  ToastAlert.warning({
    title: '警告',
    message: '您即将删除重要数据，此操作不可撤销'
  })
}

const showInfo = () => {
  ToastAlert.info({
    title: '提示信息',
    message: '系统将在5分钟后进行维护，请及时保存您的工作'
  })
}
</script>
