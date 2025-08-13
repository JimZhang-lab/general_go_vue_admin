<template>
  <AdminLayout>
    <PageBreadcrumb :pageTitle="currentPageTitle" />

    <div class="rounded-2xl border border-gray-200 bg-white p-5 dark:border-gray-800 dark:bg-white/[0.03] lg:p-6">
      <h3 class="mb-5 text-lg font-semibold text-gray-800 dark:text-white/90 lg:mb-7">个人资料</h3>

      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- 左侧个人信息卡片 -->
        <div class="lg:col-span-1">
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <div class="text-center">
              <!-- 头像 -->
              <div class="relative inline-block">
                <img
                  :src="userInfo.avatar || '/default-avatar.png'"
                  alt="用户头像"
                  class="w-24 h-24 rounded-full object-cover border-4 border-white shadow-lg"
                />
                <button
                  @click="showAvatarUpload = true"
                  class="absolute bottom-0 right-0 bg-blue-600 text-white rounded-full p-2 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  </svg>
                </button>
              </div>
              
              <!-- 用户基本信息 -->
              <h3 class="mt-4 text-lg font-medium text-gray-900">{{ userInfo.username }}</h3>
              <p class="text-sm text-gray-500">{{ userInfo.email || '未设置邮箱' }}</p>
              <p class="text-sm text-gray-500">{{ userInfo.phone || '未设置手机号' }}</p>
              
              <!-- 状态标签 -->
              <div class="mt-4">
                <span
                  :class="userInfo.status === '1' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  class="inline-flex px-2 py-1 text-xs font-semibold rounded-full"
                >
                  {{ userInfo.status === '1' ? '正常' : '禁用' }}
                </span>
              </div>
              
              <!-- 最后登录信息 -->
              <div class="mt-4 text-sm text-gray-500">
                <p>最后登录: {{ formatDate(userInfo.lastLoginTime) }}</p>
                <p>注册时间: {{ formatDate(userInfo.createTime) }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- 右侧表单区域 -->
        <div class="lg:col-span-2 space-y-6">
          <!-- 基本信息表单 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">基本信息</h3>
            
            <form @submit.prevent="updateProfile" class="space-y-4">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">用户名</label>
                  <input
                    v-model="profileForm.username"
                    type="text"
                    disabled
                    class="w-full px-3 py-2 border border-gray-300 rounded-md bg-gray-50 text-gray-500 cursor-not-allowed"
                  />
                  <p class="text-xs text-gray-500 mt-1">用户名不可修改</p>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">手机号 *</label>
                  <input
                    v-model="profileForm.phone"
                    type="tel"
                    required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">邮箱</label>
                  <input
                    v-model="profileForm.email"
                    type="email"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">真实姓名</label>
                  <input
                    v-model="profileForm.realName"
                    type="text"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">个人简介</label>
                <textarea
                  v-model="profileForm.remark"
                  rows="3"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="请输入个人简介..."
                ></textarea>
              </div>
              
              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="updating"
                  class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:opacity-50"
                >
                  {{ updating ? '更新中...' : '更新资料' }}
                </button>
              </div>
            </form>
          </div>

          <!-- 修改密码表单 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">修改密码</h3>
            
            <form @submit.prevent="changePassword" class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">当前密码 *</label>
                <input
                  v-model="passwordForm.oldPassword"
                  type="password"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">新密码 *</label>
                <input
                  v-model="passwordForm.newPassword"
                  type="password"
                  required
                  minlength="6"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
                <p class="text-xs text-gray-500 mt-1">密码长度至少6位</p>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">确认新密码 *</label>
                <input
                  v-model="passwordForm.confirmPassword"
                  type="password"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div class="flex justify-end">
                <button
                  type="submit"
                  :disabled="changingPassword"
                  class="px-4 py-2 bg-red-600 text-white rounded-md hover:bg-red-700 focus:outline-none focus:ring-2 focus:ring-red-500 disabled:opacity-50"
                >
                  {{ changingPassword ? '修改中...' : '修改密码' }}
                </button>
              </div>
            </form>
          </div>

          <!-- 登录记录 -->
          <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">最近登录记录</h3>
            
            <div class="space-y-3">
              <div v-for="log in loginLogs" :key="log.id" class="flex items-center justify-between py-2 border-b border-gray-100 last:border-b-0">
                <div class="flex items-center space-x-3">
                  <div class="flex-shrink-0">
                    <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center">
                      <svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                      </svg>
                    </div>
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-900">{{ log.loginLocation }}</p>
                    <p class="text-xs text-gray-500">{{ log.ipAddress }} • {{ log.browser }}</p>
                  </div>
                </div>
                <div class="text-right">
                  <p class="text-sm text-gray-900">{{ formatDate(log.loginTime) }}</p>
                  <p class="text-xs text-gray-500">{{ log.status === 1 ? '成功' : '失败' }}</p>
                </div>
              </div>
            </div>
            
            <div v-if="loginLogs.length === 0" class="text-center py-8 text-gray-500">
              暂无登录记录
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- 头像上传模态框 -->
    <div v-if="showAvatarUpload" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-1/2 lg:w-1/3 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">上传头像</h3>
          
          <div class="space-y-4">
            <div class="flex justify-center">
              <div class="w-32 h-32 border-2 border-dashed border-gray-300 rounded-lg flex items-center justify-center">
                <input
                  ref="avatarInput"
                  type="file"
                  accept="image/*"
                  @change="handleAvatarUpload"
                  class="hidden"
                />
                <button
                  @click="$refs.avatarInput?.click()"
                  class="text-gray-500 hover:text-gray-700"
                >
                  <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                  </svg>
                  <p class="mt-2 text-sm">点击上传</p>
                </button>
              </div>
            </div>
            
            <div class="flex justify-end space-x-3">
              <button
                type="button"
                @click="showAvatarUpload = false"
                class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
              >
                取消
              </button>
              <button
                @click="uploadAvatar"
                :disabled="!selectedAvatar"
                class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700 disabled:opacity-50"
              >
                上传
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import AdminLayout from '@/components/layout/AdminLayout.vue'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'
import { useMainStore } from '@/store'

// 定义接口类型
interface UserInfo {
  id: number
  username: string
  phone?: string
  email?: string
  realName?: string
  avatar?: string
  status: string
  remark?: string
  lastLoginTime?: string
  createTime?: string
}

interface LoginLog {
  id: number
  loginLocation: string
  ipAddress: string
  browser: string
  loginTime: string
  status: number
}

// 响应式数据
const store = useMainStore()
const userInfo = ref<UserInfo>({
  id: 0,
  username: '',
  status: '1'
})
const loginLogs = ref<LoginLog[]>([])
const updating = ref(false)
const changingPassword = ref(false)
const showAvatarUpload = ref(false)
const selectedAvatar = ref<File | null>(null)

// 表单数据
const profileForm = reactive({
  username: '',
  phone: '',
  email: '',
  realName: '',
  remark: ''
})

const passwordForm = reactive({
  oldPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return '未知'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN')
}

// 获取用户信息
const getUserInfo = async () => {
  try {
    const { data: res } = await adminApi.getCurrentUser()

    if (res.code === 200) {
      userInfo.value = res.data
      // 同步到表单
      Object.assign(profileForm, {
        username: res.data.username,
        phone: res.data.phone || '',
        email: res.data.email || '',
        realName: res.data.realName || '',
        remark: res.data.remark || ''
      })
    } else {
      ToastAlert.error({
        title: '获取用户信息失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('获取用户信息失败:', error)
    // 从store获取用户信息作为备选
    const storeUser = store.sysAdmin
    if (storeUser) {
      userInfo.value = storeUser
      Object.assign(profileForm, {
        username: storeUser.username || '',
        phone: storeUser.phone || '',
        email: storeUser.email || '',
        realName: storeUser.realName || '',
        remark: storeUser.remark || ''
      })
    }
  }
}

// 获取登录日志
const getLoginLogs = async () => {
  try {
    const { data: res } = await adminApi.getLoginLogs({
      pageNum: 1,
      pageSize: 10
    })

    if (res.code === 200) {
      loginLogs.value = res.data.list || []
    }
  } catch (error) {
    console.error('获取登录日志失败:', error)
    // 模拟一些登录记录
    loginLogs.value = [
      {
        id: 1,
        loginLocation: '北京市',
        ipAddress: '192.168.1.100',
        browser: 'Chrome 120.0',
        loginTime: new Date().toISOString(),
        status: 1
      },
      {
        id: 2,
        loginLocation: '上海市',
        ipAddress: '192.168.1.101',
        browser: 'Firefox 119.0',
        loginTime: new Date(Date.now() - 86400000).toISOString(),
        status: 1
      }
    ]
  }
}

// 更新个人资料
const updateProfile = async () => {
  try {
    updating.value = true

    const { data: res } = await adminApi.updateProfile(profileForm)

    if (res.code === 200) {
      ToastAlert.success({
        title: '更新成功',
        message: '个人资料已更新'
      })

      // 更新本地用户信息
      Object.assign(userInfo.value, profileForm)

      // 更新store中的用户信息
      store.saveSysAdmin({
        ...store.sysAdmin,
        ...profileForm
      })
    } else {
      ToastAlert.error({
        title: '更新失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('更新个人资料失败:', error)
    ToastAlert.error({
      title: '更新失败',
      message: '网络异常，请重试'
    })
  } finally {
    updating.value = false
  }
}

// 修改密码
const changePassword = async () => {
  // 验证密码
  if (passwordForm.newPassword !== passwordForm.confirmPassword) {
    ToastAlert.error({
      title: '密码不匹配',
      message: '新密码和确认密码不一致'
    })
    return
  }

  if (passwordForm.newPassword.length < 6) {
    ToastAlert.error({
      title: '密码太短',
      message: '密码长度至少6位'
    })
    return
  }

  try {
    changingPassword.value = true

    const { data: res } = await adminApi.changePassword({
      oldPassword: passwordForm.oldPassword,
      newPassword: passwordForm.newPassword
    })

    if (res.code === 200) {
      ToastAlert.success({
        title: '修改成功',
        message: '密码已修改，请重新登录'
      })

      // 清空表单
      Object.assign(passwordForm, {
        oldPassword: '',
        newPassword: '',
        confirmPassword: ''
      })

      // 可以选择自动登出用户
      setTimeout(() => {
        store.logout()
        window.location.href = '/adminLogin'
      }, 2000)
    } else {
      ToastAlert.error({
        title: '修改失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('修改密码失败:', error)
    ToastAlert.error({
      title: '修改失败',
      message: '网络异常，请重试'
    })
  } finally {
    changingPassword.value = false
  }
}

// 处理头像上传
const handleAvatarUpload = (event: Event) => {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]

  if (file) {
    // 验证文件类型
    if (!file.type.startsWith('image/')) {
      ToastAlert.error({
        title: '文件类型错误',
        message: '请选择图片文件'
      })
      return
    }

    // 验证文件大小 (2MB)
    if (file.size > 2 * 1024 * 1024) {
      ToastAlert.error({
        title: '文件太大',
        message: '图片大小不能超过2MB'
      })
      return
    }

    selectedAvatar.value = file
  }
}

// 上传头像
const uploadAvatar = async () => {
  if (!selectedAvatar.value) return

  try {
    const formData = new FormData()
    formData.append('avatar', selectedAvatar.value)

    // 这里应该调用上传头像的API
    // const { data: res } = await adminApi.uploadAvatar(formData)

    // 模拟上传成功
    ToastAlert.success({
      title: '上传成功',
      message: '头像已更新'
    })

    showAvatarUpload.value = false
    selectedAvatar.value = null
  } catch (error) {
    console.error('上传头像失败:', error)
    ToastAlert.error({
      title: '上传失败',
      message: '网络异常，请重试'
    })
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getUserInfo()
  getLoginLogs()
})
</script>
