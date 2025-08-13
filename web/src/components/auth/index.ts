/*
 * @Author: JimZhang
 * @Date: 2025-07-27
 * @Description: Auth组件统一导出
 */

// 导出Auth相关组件
export { default as AuthCard } from './AuthCard.vue'
export { default as AuthButton } from './AuthButton.vue'
export { default as AuthInput } from './AuthInput.vue'

// 布局组件
export { default as AuthLayout } from './AuthLayout.vue'
export { default as AuthSidebar } from './AuthSidebar.vue'
export { default as AuthHeader } from './AuthHeader.vue'

// Header子组件
export { default as AuthHeaderLogo } from './AuthHeaderLogo.vue'
export { default as AuthSearchBar } from './AuthSearchBar.vue'
export { default as AuthNotificationMenu } from './AuthNotificationMenu.vue'
export { default as AuthUserMenu } from './AuthUserMenu.vue'

// Sidebar子组件
export { default as AuthSidebarWidget } from './AuthSidebarWidget.vue'

// 组件类型定义
export interface AuthCardProps {
  title?: string
  subtitle?: string
  cardClass?: string
  contentClass?: string
}

export interface AuthButtonProps {
  text: string
  type?: 'button' | 'submit' | 'reset'
  variant?: 'primary' | 'secondary' | 'success' | 'danger' | 'warning' | 'info'
  size?: 'sm' | 'md' | 'lg'
  disabled?: boolean
  loading?: boolean
  loadingText?: string
  icon?: any
  block?: boolean
}

export interface AuthInputProps {
  modelValue: string | number
  type?: 'text' | 'email' | 'password' | 'number' | 'tel' | 'url'
  label?: string
  placeholder?: string
  name?: string
  disabled?: boolean
  readonly?: boolean
  required?: boolean
  autocomplete?: string
  leftIcon?: any
  rightIcon?: any
  error?: string
  help?: string
  size?: 'sm' | 'md' | 'lg'
}

// Auth组件常量
export const AUTH_COMPONENT_CONSTANTS = {
  // 按钮变体
  BUTTON_VARIANTS: {
    PRIMARY: 'primary',
    SECONDARY: 'secondary',
    SUCCESS: 'success',
    DANGER: 'danger',
    WARNING: 'warning',
    INFO: 'info'
  } as const,

  // 组件尺寸
  SIZES: {
    SMALL: 'sm',
    MEDIUM: 'md',
    LARGE: 'lg'
  } as const,

  // 输入框类型
  INPUT_TYPES: {
    TEXT: 'text',
    EMAIL: 'email',
    PASSWORD: 'password',
    NUMBER: 'number',
    TEL: 'tel',
    URL: 'url'
  } as const
}

// 默认导出
export default {
  AuthCard: () => import('./AuthCard.vue'),
  AuthButton: () => import('./AuthButton.vue'),
  AuthInput: () => import('./AuthInput.vue')
}
