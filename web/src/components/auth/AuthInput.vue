<template>
  <div class="space-y-2">
    <!-- 标签 -->
    <label v-if="label" :for="inputId" class="block text-sm font-medium text-gray-700 dark:text-gray-300">
      {{ label }}
      <span v-if="required" class="text-red-500 ml-1">*</span>
    </label>

    <!-- 输入框容器 -->
    <div class="relative">
      <!-- 左侧图标 -->
      <div v-if="leftIcon" class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <component :is="leftIcon" class="h-5 w-5 text-gray-400" />
      </div>

      <!-- 输入框 -->
      <input
        :id="inputId"
        :type="inputType"
        :name="name"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :required="required"
        :autocomplete="autocomplete"
        :value="modelValue"
        :class="inputClasses"
        @input="handleInput"
        @blur="handleBlur"
        @focus="handleFocus"
      />

      <!-- 右侧图标/按钮 -->
      <div v-if="rightIcon || type === 'password'" class="absolute inset-y-0 right-0 pr-3 flex items-center">
        <!-- 密码显示/隐藏按钮 -->
        <button
          v-if="type === 'password'"
          type="button"
          @click="togglePasswordVisibility"
          class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300 transition-colors"
        >
          <svg
            v-if="showPassword"
            class="h-5 w-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
          <svg
            v-else
            class="h-5 w-5"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L3 3m6.878 6.878L21 21" />
          </svg>
        </button>

        <!-- 自定义右侧图标 -->
        <component v-else-if="rightIcon" :is="rightIcon" class="h-5 w-5 text-gray-400" />
      </div>
    </div>

    <!-- 错误信息 -->
    <p v-if="error" class="text-sm text-red-600 dark:text-red-400">{{ error }}</p>

    <!-- 帮助文本 -->
    <p v-if="help && !error" class="text-sm text-gray-500 dark:text-gray-400">{{ help }}</p>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Props {
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

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  label: '',
  placeholder: '',
  name: '',
  disabled: false,
  readonly: false,
  required: false,
  autocomplete: '',
  leftIcon: null,
  rightIcon: null,
  error: '',
  help: '',
  size: 'md'
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
  blur: [event: FocusEvent]
  focus: [event: FocusEvent]
}>()

// 生成唯一ID
const inputId = computed(() => props.name || `input-${Math.random().toString(36).substr(2, 9)}`)

// 密码显示状态
const showPassword = ref(false)

// 计算输入框类型
const inputType = computed(() => {
  if (props.type === 'password') {
    return showPassword.value ? 'text' : 'password'
  }
  return props.type
})

// 输入框样式
const inputClasses = computed(() => {
  const baseClasses = [
    'block',
    'w-full',
    'border',
    'rounded-xl',
    'placeholder-gray-400',
    'focus:outline-none',
    'focus:ring-2',
    'focus:ring-blue-500',
    'focus:border-transparent',
    'transition-colors',
    'dark:bg-gray-800',
    'dark:text-white',
    'dark:placeholder-gray-500',
    'dark:border-gray-700'
  ]

  // 尺寸样式
  const sizeClasses = {
    sm: ['px-3', 'py-2', 'text-sm'],
    md: ['px-4', 'py-3', 'text-sm'],
    lg: ['px-5', 'py-4', 'text-base']
  }

  // 左侧图标间距
  if (props.leftIcon) {
    sizeClasses[props.size] = sizeClasses[props.size].map(cls => 
      cls.startsWith('px-') ? `pl-11 pr-${cls.split('-')[1]}` : cls
    )
  }

  // 右侧图标间距
  if (props.rightIcon || props.type === 'password') {
    sizeClasses[props.size] = sizeClasses[props.size].map(cls => 
      cls.startsWith('px-') || cls.startsWith('pr-') ? `${cls.split(' ')[0]} pr-12` : cls
    )
  }

  // 错误状态样式
  if (props.error) {
    baseClasses.push('border-red-300', 'focus:ring-red-500', 'focus:border-red-500')
  } else {
    baseClasses.push('border-gray-300')
  }

  // 禁用状态样式
  if (props.disabled) {
    baseClasses.push('bg-gray-50', 'cursor-not-allowed', 'dark:bg-gray-900')
  }

  return [
    ...baseClasses,
    ...sizeClasses[props.size]
  ].join(' ')
})

// 事件处理
const handleInput = (event: Event) => {
  const target = event.target as HTMLInputElement
  const value = props.type === 'number' ? Number(target.value) : target.value
  emit('update:modelValue', value)
}

const handleBlur = (event: FocusEvent) => {
  emit('blur', event)
}

const handleFocus = (event: FocusEvent) => {
  emit('focus', event)
}

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}
</script>
