<template>
  <Teleport to="body">
    <div
      v-if="visible"
      class="toast-overlay"
      :class="positionClasses[position]"
    >
      <transition name="backdrop">
        <div
          v-if="visible"
          class="toast-backdrop fixed inset-0 bg-transparent"
          :class="{ 'backdrop-blur-sm': blurBackground }"
          @click="close"
        ></div>
      </transition>
      
      <transition 
        :name="position === 'top' ? 'slide-down' : 'fade'"
        appear
      >
        <div
          v-if="visible"
          class="toast-content rounded-xl border p-4 backdrop-blur-lg bg-white/95 dark:bg-gray-800/95 w-auto min-w-[20rem] max-w-2xl shadow-2xl"
          :class="[variantClasses[variant].container, sizeClasses]"
        >
          <div class="flex items-start gap-3">
            <div :class="['-mt-0.5', variantClasses[variant].icon]">
              <component :is="icons[variant]" />
            </div>
            <div class="flex-1">
              <h4 class="mb-1 text-sm font-semibold text-gray-800 dark:text-white/90">
                {{ title }}
              </h4>
              <p v-if="message" class="text-sm text-gray-500 dark:text-gray-400 break-words whitespace-normal">{{ message }}</p>

              <!-- 错误列表 -->
              <div v-if="errors && errors.length > 0" class="mt-2 space-y-1">
                <div
                  v-for="(error, index) in errors"
                  :key="index"
                  class="text-sm text-gray-500 dark:text-gray-400 break-words whitespace-normal pl-3 border-l-2 border-gray-300 dark:border-gray-600"
                >
                  {{ error }}
                </div>
              </div>
            </div>
          </div>
          <div class="mt-4 flex justify-end gap-2">
            <button
              v-if="showCancel"
              @click="cancel"
              class="px-4 py-2 text-sm font-medium rounded-lg border border-gray-300 text-gray-700 hover:bg-gray-100 dark:border-gray-600 dark:text-gray-300 dark:hover:bg-gray-800"
            >
              取消
            </button>
            <button
              @click="confirm"
              :class="confirmButtonClasses"
            >
              确定
            </button>
          </div>
        </div>
      </transition>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { ref, watch, onUnmounted, computed } from 'vue'
import { SuccessIcon, ErrorIcon, WarningIcon, InfoCircleIcon } from '@/icons'

interface ToastProps {
  modelValue: boolean
  variant: 'success' | 'error' | 'warning' | 'info'
  title: string
  message?: string
  showCancel?: boolean
  autoClose?: boolean
  duration?: number
  position?: 'center' | 'top'
  blurBackground?: boolean
  errors?: string[]
}

const props = withDefaults(defineProps<ToastProps>(), {
  message: '',
  showCancel: false,
  autoClose: true,
  duration: 2000,
  position: 'center',
  blurBackground: false,
  errors: () => []
})

const emit = defineEmits<{
  (e: 'update:modelValue', value: boolean): void
  (e: 'confirm'): void
  (e: 'cancel'): void
}>()

// 单例模式：确保同一时间只能显示一个弹窗
let currentInstance: any = null

const visible = ref(false)
let timer: number | null = null

type VariantClass = {
  container: string
  icon: string
}

const variantClasses: Record<string, VariantClass> = {
  success: {
    container: 'border-success-500 dark:border-success-500/30',
    icon: 'text-success-500',
  },
  error: {
    container: 'border-error-500 dark:border-error-500/30',
    icon: 'text-error-500',
  },
  warning: {
    container: 'border-warning-500 dark:border-warning-500/30',
    icon: 'text-warning-500',
  },
  info: {
    container: 'border-blue-light-500 dark:border-blue-light-500/30',
    icon: 'text-blue-light-500',
  },
}

const positionClasses = {
  center: 'flex items-center justify-center',
  top: 'flex items-start justify-center pt-10'
}

// 确认按钮样式
const confirmButtonClasses = computed(() => {
  const baseClasses = 'px-4 py-2 text-sm font-medium rounded-lg text-white'

  switch (props.variant) {
    case 'success':
      return `${baseClasses} bg-green-600 hover:bg-green-700 dark:bg-green-500 dark:hover:bg-green-600`
    case 'error':
      return `${baseClasses} bg-red-600 hover:bg-red-700 dark:bg-red-500 dark:hover:bg-red-600`
    case 'warning':
      return `${baseClasses} bg-orange-600 hover:bg-orange-700 dark:bg-orange-500 dark:hover:bg-orange-600`
    case 'info':
    default:
      return `${baseClasses} bg-blue-600 hover:bg-blue-700 dark:bg-blue-500 dark:hover:bg-blue-600`
  }
})

const icons = {
  success: SuccessIcon,
  error: ErrorIcon,
  warning: WarningIcon,
  info: InfoCircleIcon,
}

// 根据内容计算弹窗大小
const sizeClasses = computed(() => {
  const hasErrors = props.errors && props.errors.length > 0
  const hasLongMessage = props.message && props.message.length > 100
  const hasMultipleErrors = hasErrors && props.errors!.length > 3

  if (hasMultipleErrors || hasLongMessage) {
    return 'max-w-3xl'
  } else if (hasErrors || (props.message && props.message.length > 50)) {
    return 'max-w-xl'
  }
  return 'max-w-md'
})

const close = () => {
  visible.value = false
  emit('update:modelValue', false)
  clearTimer()
  // 清除当前实例引用
  if (currentInstance === visible) {
    currentInstance = null
  }
}

const confirm = () => {
  emit('confirm')
  close()
}

const cancel = () => {
  emit('cancel')
  close()
}

const clearTimer = () => {
  if (timer) {
    clearTimeout(timer)
    timer = null
  }
}

const startAutoCloseTimer = () => {
  if (props.autoClose && props.duration > 0) {
    clearTimer()
    timer = setTimeout(() => {
      close()
    }, props.duration) as unknown as number
  }
}

watch(
  () => props.modelValue,
  (newVal) => {
    // 单例模式：如果已经有实例显示且尝试显示新实例，则关闭当前实例
    if (newVal && currentInstance && currentInstance !== visible) {
      currentInstance.value = false
    }
    
    // 如果要显示弹窗，设置当前实例
    if (newVal) {
      currentInstance = visible
    }
    
    visible.value = newVal
    if (newVal) {
      startAutoCloseTimer()
    } else {
      clearTimer()
      // 清除当前实例引用
      if (currentInstance === visible) {
        currentInstance = null
      }
    }
  },
  { immediate: true }
)

onUnmounted(() => {
  clearTimer()
  // 组件销毁时清除当前实例引用
  if (currentInstance === visible) {
    currentInstance = null
  }
})
</script>

<style scoped>
/* 淡入淡出动画 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 背景遮罩动画 */
.backdrop-enter-active,
.backdrop-leave-active {
  transition: opacity 0.3s ease;
}

.backdrop-enter-from,
.backdrop-leave-to {
  opacity: 0;
}

/* 顶部滑入滑出动画 */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: transform 0.3s ease;
}

.slide-down-enter-from {
  transform: translateY(-100%);
}

.slide-down-leave-to {
  transform: translateY(-100%);
}

/* 确保Toast始终在最上层 */
.toast-container {
  position: fixed !important;
  z-index: var(--z-index-toast) !important;
  pointer-events: none;
}

.toast-container > * {
  pointer-events: auto;
}

/* 防止被其他元素遮挡 */
.toast-backdrop {
  backdrop-filter: blur(4px);
  -webkit-backdrop-filter: blur(4px);
}
</style>