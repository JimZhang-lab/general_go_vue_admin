<template>
  <div class="p-6">
    <h2 class="text-xl font-bold mb-4">Toast 弹窗示例</h2>
    
    <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
      <button 
        @click="showSuccessToast" 
        class="px-4 py-2 bg-green-500 text-white rounded hover:bg-green-600"
      >
        显示成功弹窗
      </button>
      
      <button 
        @click="showErrorToast" 
        class="px-4 py-2 bg-red-500 text-white rounded hover:bg-red-600"
      >
        显示错误弹窗
      </button>
      
      <button 
        @click="showWarningToast" 
        class="px-4 py-2 bg-yellow-500 text-white rounded hover:bg-yellow-600"
      >
        显示警告弹窗
      </button>
      
      <button 
        @click="showInfoToast" 
        class="px-4 py-2 bg-blue-500 text-white rounded hover:bg-blue-600"
      >
        显示信息弹窗
      </button>
    </div>
    
    <!-- Toast 组件 -->
    <Toast 
      v-model="toastVisible"
      :variant="toastVariant"
      :title="toastTitle"
      :message="toastMessage"
      :show-cancel="showCancel"
      @confirm="handleConfirm"
      @cancel="handleCancel"
    />
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import Toast from './Toast.vue'

// Toast 状态
const toastVisible = ref(false)
const toastVariant = ref<'success' | 'error' | 'warning' | 'info'>('success')
const toastTitle = ref('')
const toastMessage = ref('')
const showCancel = ref(false)

// 显示不同类型 Toast 的方法
const showSuccessToast = () => {
  toastVariant.value = 'success'
  toastTitle.value = '操作成功'
  toastMessage.value = '您的操作已成功完成'
  showCancel.value = false
  toastVisible.value = true
}

const showErrorToast = () => {
  toastVariant.value = 'error'
  toastTitle.value = '操作失败'
  toastMessage.value = '操作过程中发生错误，请重试'
  showCancel.value = false
  toastVisible.value = true
}

const showWarningToast = () => {
  toastVariant.value = 'warning'
  toastTitle.value = '警告'
  toastMessage.value = '请注意，此操作不可逆'
  showCancel.value = true
  toastVisible.value = true
}

const showInfoToast = () => {
  toastVariant.value = 'info'
  toastTitle.value = '提示信息'
  toastMessage.value = '这是系统提供的相关信息'
  showCancel.value = true
  toastVisible.value = true
}

// Toast 事件处理
const handleConfirm = () => {
  console.log('用户点击了确认')
}

const handleCancel = () => {
  console.log('用户点击了取消')
}
</script>