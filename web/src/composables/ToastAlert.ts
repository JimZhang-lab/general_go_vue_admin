import { createVNode, render } from 'vue'
import type { VNode } from 'vue'
import Toast from '@/components/ui/Toast.vue'

type ToastVariant = 'success' | 'error' | 'warning' | 'info'
type ToastPosition = 'center' | 'top'

interface ToastOptions {
  title: string
  message?: string
  variant?: ToastVariant
  showCancel?: boolean
  onConfirm?: () => void
  onCancel?: () => void
  autoClose?: boolean
  duration?: number
  position?: ToastPosition
  blurBackground?: boolean
  errors?: string[] // 支持多个错误消息
}

class ToastAlert {
  // 用于跟踪当前是否已有弹窗显示
  private static currentToast: HTMLElement | null = null
  // 错误队列，用于收集连续的错误
  private static errorQueue: string[] = []
  // 错误收集定时器
  private static errorTimer: number | null = null

  private static mountToast(options: ToastOptions): Promise<boolean> {
    // 如果已有弹窗显示，根据类型决定处理方式
    if (this.currentToast) {
      // 如果是错误类型，加入队列
      if (options.variant === 'error' && options.message) {
        this.addToErrorQueue(options.message)
        return Promise.resolve(false)
      }
      // 其他类型直接返回
      return Promise.resolve(false)
    }

    return new Promise((resolve) => {
      const container = document.createElement('div')
      this.currentToast = container
      
      const defaultOptions = {
        variant: 'info' as ToastVariant,
        showCancel: false,
        autoClose: true,
        duration: 2000,
        position: 'top' as ToastPosition,
        blurBackground: false,
        ...options
      }

      const handleConfirm = () => {
        options.onConfirm?.()
        resolve(true)
      }

      const handleCancel = () => {
        options.onCancel?.()
        resolve(false)
      }

      const handleClose = () => {
        render(null, container)
        if (container.parentNode) {
          container.parentNode.removeChild(container)
        }
        this.currentToast = null

        // 如果有错误队列，显示收集的错误
        if (this.errorQueue.length > 0) {
          const errors = [...this.errorQueue]
          this.errorQueue = []
          this.clearErrorTimer()

          // 延迟一点时间再显示错误弹窗，避免重叠
          setTimeout(() => {
            this.showCollectedErrors(errors)
          }, 100)
        }
      }

      const vnode: VNode = createVNode(Toast, {
        modelValue: true,
        ...defaultOptions,
        onConfirm: handleConfirm,
        onCancel: handleCancel,
        'onUpdate:modelValue': (val: boolean) => {
          if (!val) {
            handleClose()
          }
        }
      })

      document.body.appendChild(container)
      render(vnode, container)
    })
  }

  static success(options: Omit<ToastOptions, 'variant'>) {
    return this.mountToast({ ...options, variant: 'success' })
  }

  static error(options: Omit<ToastOptions, 'variant'>) {
    return this.mountToast({ ...options, variant: 'error' })
  }

  static warning(options: Omit<ToastOptions, 'variant'>) {
    return this.mountToast({ ...options, variant: 'warning' })
  }

  static info(options: Omit<ToastOptions, 'variant'>) {
    return this.mountToast({ ...options, variant: 'info' })
  }

  static confirm(options: Omit<ToastOptions, 'showCancel'>) {
    return this.mountToast({ ...options, showCancel: true, autoClose: false })
  }

  // 添加错误到队列
  private static addToErrorQueue(message: string) {
    if (!this.errorQueue.includes(message)) {
      this.errorQueue.push(message)
    }

    // 重置定时器，延迟显示错误
    this.clearErrorTimer()
    this.errorTimer = setTimeout(() => {
      if (this.errorQueue.length > 0 && !this.currentToast) {
        const errors = [...this.errorQueue]
        this.errorQueue = []
        this.showCollectedErrors(errors)
      }
    }, 1000) as unknown as number
  }

  // 清除错误定时器
  private static clearErrorTimer() {
    if (this.errorTimer) {
      clearTimeout(this.errorTimer)
      this.errorTimer = null
    }
  }

  // 显示收集的错误
  private static showCollectedErrors(errors: string[]) {
    const title = errors.length === 1 ? '错误' : `发现 ${errors.length} 个错误`
    const message = errors.length === 1 ? errors[0] : undefined

    this.mountToast({
      title,
      message,
      errors: errors.length > 1 ? errors : undefined,
      variant: 'error',
      autoClose: false,
      showCancel: false
    })
  }
}

export default ToastAlert