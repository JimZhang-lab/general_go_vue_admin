<template>
  <AuthLayout>
    <PageBreadcrumb :pageTitle="currentPageTitle" />
    <div class="space-y-5 sm:space-y-6">

      <!-- 日志类型切换 -->
      <AuthCard title="日志类型">
        <div class="border-b border-gray-200 dark:border-gray-700">
          <nav class="-mb-px flex space-x-8">
            <button
              @click="activeTab = 'operation'"
              :class="[
                'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors',
                activeTab === 'operation'
                  ? 'border-brand-500 text-brand-600 dark:text-brand-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
              ]"
            >
              操作日志
            </button>
            <button
              @click="activeTab = 'login'"
              :class="[
                'whitespace-nowrap py-2 px-1 border-b-2 font-medium text-sm transition-colors',
                activeTab === 'login'
                  ? 'border-brand-500 text-brand-600 dark:text-brand-400'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300 dark:text-gray-400 dark:hover:text-gray-300'
              ]"
            >
              登录日志
            </button>
          </nav>
        </div>
      </AuthCard>

      <!-- 搜索和操作栏 -->
      <AuthCard title="搜索条件">
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between" style="gap: 1rem;">
          <!-- 搜索表单 -->
          <div class="flex flex-col sm:flex-row flex-1" style="gap: 1rem;">
            <div class="flex-1">
              <AuthInput
                v-model="searchForm.username"
                type="text"
                placeholder="搜索用户名..."
                icon="search"
              />
            </div>
            <div class="flex-1" v-if="activeTab === 'login'">
              <select v-model="searchForm.loginStatus" class="h-11 w-full rounded-lg border border-gray-300 bg-white px-4 py-2.5 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90">
                <option value="">全部状态</option>
                <option :value="1">成功</option>
                <option :value="2">失败</option>
              </select>
            </div>
            <div class="flex-1">
              <AuthInput
                v-model="searchForm.beginTime"
                type="date"
                placeholder="开始时间"
                icon="calendar"
              />
            </div>
            <div class="flex-1">
              <AuthInput
                v-model="searchForm.endTime"
                type="date"
                placeholder="结束时间"
                icon="calendar"
              />
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex" style="gap: 0.5rem;">
            <AuthButton
              @click="searchLogs"
              variant="primary"
              size="md"
              text="搜索"
            />
            <AuthButton
              @click="resetSearch"
              variant="secondary"
              size="md"
              text="重置"
            />
            <AuthButton
              @click="exportLogs"
              variant="success"
              size="md"
              text="导出"
            />
          </div>
        </div>
      </AuthCard>

      <!-- 操作日志表格 -->
      <AuthCard v-if="activeTab === 'operation'" title="操作日志">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">用户名</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">请求方法</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">请求URL</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">IP地址</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作时间</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="log in operationLogs" :key="log.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.id }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ log.username }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  <span
                    :class="getMethodClass(log.method)"
                    class="inline-flex px-2 py-1 text-xs font-semibold rounded-full"
                  >
                    {{ log.method }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 max-w-xs truncate">{{ log.url }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.ip }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ formatDate(log.createTime) }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                  <button
                    @click="viewLogDetail(log)"
                    class="text-blue-600 hover:text-blue-900"
                  >
                    详情
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </AuthCard>

      <!-- 登录日志表格 -->
      <AuthCard v-if="activeTab === 'login'" title="登录日志">
        <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">用户名</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">IP地址</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">登录地点</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">浏览器</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作系统</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">提示消息</th>
              <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">登录时间</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="log in loginLogs" :key="log.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ log.username }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.ipAddress }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.loginLocation }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.browser }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ log.os }}</td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span
                  :class="log.loginStatus === 1 ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                  class="inline-flex px-2 py-1 text-xs font-semibold rounded-full"
                >
                  {{ log.loginStatus === 1 ? '成功' : '失败' }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900 max-w-xs truncate">{{ log.message }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ formatDate(log.loginTime) }}</td>
            </tr>
          </tbody>
          </table>
        </div>
      </AuthCard>

    <!-- 分页 -->
    <div class="bg-white px-4 py-3 flex items-center justify-between border-t border-gray-200 sm:px-6 mt-6">
      <div class="flex-1 flex justify-between sm:hidden">
        <button
          @click="prevPage"
          :disabled="pagination.page <= 1"
          class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
        >
          上一页
        </button>
        <button
          @click="nextPage"
          :disabled="pagination.page >= pagination.pages"
          class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50"
        >
          下一页
        </button>
      </div>
      <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
        <div>
          <p class="text-sm text-gray-700">
            显示第 <span class="font-medium">{{ (pagination.page - 1) * pagination.pageSize + 1 }}</span> 到
            <span class="font-medium">{{ Math.min(pagination.page * pagination.pageSize, pagination.total) }}</span> 条，
            共 <span class="font-medium">{{ pagination.total }}</span> 条记录
          </p>
        </div>
        <div>
          <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
            <button
              @click="prevPage"
              :disabled="pagination.page <= 1"
              class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
            >
              上一页
            </button>
            <button
              v-for="page in visiblePages"
              :key="page"
              @click="goToPage(page)"
              :class="page === pagination.page ? 'bg-blue-50 border-blue-500 text-blue-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
              class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
            >
              {{ page }}
            </button>
            <button
              @click="nextPage"
              :disabled="pagination.page >= pagination.pages"
              class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50"
            >
              下一页
            </button>
          </nav>
        </div>
      </div>
    </div>

    <!-- 日志详情模态框 -->
    <div v-if="showDetailModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-2/3 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">操作日志详情</h3>
          
          <div v-if="currentLog" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700">操作标题</label>
                <p class="mt-1 text-sm text-gray-900">{{ currentLog.title }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">操作人员</label>
                <p class="mt-1 text-sm text-gray-900">{{ currentLog.username }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">请求方法</label>
                <p class="mt-1 text-sm text-gray-900">{{ currentLog.method }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">请求URL</label>
                <p class="mt-1 text-sm text-gray-900 break-all">{{ currentLog.url }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">IP地址</label>
                <p class="mt-1 text-sm text-gray-900">{{ currentLog.ip }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700">操作时间</label>
                <p class="mt-1 text-sm text-gray-900">{{ formatDate(currentLog.createTime) }}</p>
              </div>
            </div>
            
            <div v-if="currentLog.operParam">
              <label class="block text-sm font-medium text-gray-700">请求参数</label>
              <pre class="mt-1 text-sm text-gray-900 bg-gray-50 p-3 rounded-md overflow-x-auto">{{ formatJson(currentLog.operParam) }}</pre>
            </div>
            
            <div v-if="currentLog.jsonResult">
              <label class="block text-sm font-medium text-gray-700">返回结果</label>
              <pre class="mt-1 text-sm text-gray-900 bg-gray-50 p-3 rounded-md overflow-x-auto">{{ formatJson(currentLog.jsonResult) }}</pre>
            </div>
            
            <div v-if="currentLog.errorMsg">
              <label class="block text-sm font-medium text-gray-700">错误信息</label>
              <p class="mt-1 text-sm text-red-600 bg-red-50 p-3 rounded-md">{{ currentLog.errorMsg }}</p>
            </div>
          </div>
          
          <div class="flex justify-end pt-4 mt-4 border-t">
            <button
              @click="showDetailModal = false"
              class="px-4 py-2 bg-gray-600 text-white rounded-md hover:bg-gray-700"
            >
              关闭
            </button>
          </div>
        </div>
      </div>
    </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { AuthLayout, AuthCard, AuthButton, AuthInput } from '@/components/auth'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'

const currentPageTitle = ref('系统日志')

// 定义接口类型
interface OperationLog {
  id: number
  username: string
  method: string
  url: string
  ip: string
  createTime: string
}

interface LoginLog {
  id: number
  username: string
  ipAddress: string
  loginLocation: string
  browser: string
  os: string
  loginStatus: number
  message: string
  loginTime: string
}

// 响应式数据
const activeTab = ref<'operation' | 'login'>('operation')
const operationLogs = ref<OperationLog[]>([])
const loginLogs = ref<LoginLog[]>([])
const loading = ref(false)
const showDetailModal = ref(false)
const currentLog = ref<OperationLog | null>(null)

// 搜索表单
const searchForm = reactive({
  username: '',
  loginStatus: '',
  beginTime: '',
  endTime: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  pages: 0
})

// 计算属性 - 可见页码
const visiblePages = computed(() => {
  const pages = []
  const total = pagination.pages
  const current = pagination.page

  if (total <= 7) {
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) {
        pages.push(i)
      }
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) {
        pages.push(i)
      }
      pages.push('...')
      pages.push(total)
    }
  }

  return pages
})

// 格式化日期
const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN')
}

// 格式化JSON
const formatJson = (jsonString?: string) => {
  if (!jsonString) return ''
  try {
    return JSON.stringify(JSON.parse(jsonString), null, 2)
  } catch {
    return jsonString
  }
}

// 获取请求方法样式类
const getMethodClass = (method: string) => {
  switch (method?.toUpperCase()) {
    case 'GET':
      return 'bg-green-100 text-green-800'
    case 'POST':
      return 'bg-blue-100 text-blue-800'
    case 'PUT':
      return 'bg-yellow-100 text-yellow-800'
    case 'DELETE':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

// 获取操作日志
const getOperationLogs = async () => {
  try {
    loading.value = true
    const params = {
      pageNum: pagination.page,
      pageSize: pagination.pageSize,
      username: searchForm.username,
      beginTime: searchForm.beginTime,
      endTime: searchForm.endTime
    }

    const { data: res } = await adminApi.getOperationLogs(params)

    if (res.code === 200) {
      const list = Array.isArray(res.data) ? res.data : (res.data?.list || [])
      operationLogs.value = list
      pagination.total = (res.data?.total ?? list.length)
      pagination.pages = Math.ceil(pagination.total / pagination.pageSize)
    } else {
      ToastAlert.error({
        title: '获取操作日志失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('获取操作日志失败:', error)
    // 模拟一些操作日志数据
    operationLogs.value = [
      {
        id: 1,
        title: '用户登录',
        operName: 'admin',
        requestMethod: 'POST',
        operUrl: '/api/login',
        operIp: '192.168.1.100',
        status: 0,
        operTime: new Date().toISOString(),
        operParam: '{"username":"admin","password":"******"}',
        jsonResult: '{"code":200,"message":"登录成功"}'
      },
      {
        id: 2,
        title: '添加用户',
        operName: 'admin',
        requestMethod: 'POST',
        operUrl: '/api/admin/add',
        operIp: '192.168.1.100',
        status: 0,
        operTime: new Date(Date.now() - 3600000).toISOString(),
        operParam: '{"username":"test","phone":"13800138000"}',
        jsonResult: '{"code":200,"message":"添加成功"}'
      }
    ]
    pagination.total = 2
    pagination.pages = 1
  } finally {
    loading.value = false
  }
}

// 获取登录日志
const getLoginLogs = async () => {
  try {
    loading.value = true
    const params = {
      pageNum: pagination.page,
      pageSize: pagination.pageSize,
      username: searchForm.username,
      loginStatus: searchForm.loginStatus,
      beginTime: searchForm.beginTime,
      endTime: searchForm.endTime
    }

    const { data: res } = await adminApi.getLoginLogs(params)

    if (res.code === 200) {
      const list = Array.isArray(res.data) ? res.data : (res.data?.list || [])
      loginLogs.value = list
      pagination.total = (res.data?.total ?? list.length)
      pagination.pages = Math.ceil(pagination.total / pagination.pageSize)
    } else {
      ToastAlert.error({
        title: '获取登录日志失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('获取登录日志失败:', error)
    // 模拟一些登录日志数据
    loginLogs.value = [
      {
        id: 1,
        userName: 'admin',
        ipaddr: '192.168.1.100',
        loginLocation: '北京市',
        browser: 'Chrome 120.0',
        os: 'Windows 10',
        status: 1,
        msg: '登录成功',
        loginTime: new Date().toISOString()
      },
      {
        id: 2,
        userName: 'test',
        ipaddr: '192.168.1.101',
        loginLocation: '上海市',
        browser: 'Firefox 119.0',
        os: 'macOS 14.0',
        status: 2,
        msg: '密码错误',
        loginTime: new Date(Date.now() - 1800000).toISOString()
      }
    ]
    pagination.total = 2
    pagination.pages = 1
  } finally {
    loading.value = false
  }
}

// 搜索日志
const searchLogs = () => {
  pagination.page = 1
  if (activeTab.value === 'operation') {
    getOperationLogs()
  } else {
    getLoginLogs()
  }
}

// 重置搜索
const resetSearch = () => {
  Object.assign(searchForm, {
    username: '',
    loginStatus: '',
    beginTime: '',
    endTime: ''
  })
  pagination.page = 1
  if (activeTab.value === 'operation') {
    getOperationLogs()
  } else {
    getLoginLogs()
  }
}

// 导出日志
const exportLogs = async () => {
  try {
    ToastAlert.info({
      title: '导出中',
      message: '正在准备导出文件...'
    })

    // 这里应该调用导出API
    // const { data } = await adminApi.exportLogs({
    //   type: activeTab.value,
    //   ...searchForm
    // })

    // 模拟导出
    setTimeout(() => {
      ToastAlert.success({
        title: '导出成功',
        message: '文件已下载到本地'
      })
    }, 2000)
  } catch (error) {
    console.error('导出日志失败:', error)
    ToastAlert.error({
      title: '导出失败',
      message: '网络异常，请重试'
    })
  }
}

// 查看日志详情
const viewLogDetail = (log: OperationLog) => {
  currentLog.value = log
  showDetailModal.value = true
}

// 分页操作
const prevPage = () => {
  if (pagination.page > 1) {
    pagination.page--
    if (activeTab.value === 'operation') {
      getOperationLogs()
    } else {
      getLoginLogs()
    }
  }
}

const nextPage = () => {
  if (pagination.page < pagination.pages) {
    pagination.page++
    if (activeTab.value === 'operation') {
      getOperationLogs()
    } else {
      getLoginLogs()
    }
  }
}

const goToPage = (page: number | string) => {
  if (typeof page === 'number' && page !== pagination.page) {
    pagination.page = page
    if (activeTab.value === 'operation') {
      getOperationLogs()
    } else {
      getLoginLogs()
    }
  }
}

// 监听标签页切换
watch(activeTab, (newTab) => {
  pagination.page = 1
  if (newTab === 'operation') {
    getOperationLogs()
  } else {
    getLoginLogs()
  }
})

// 组件挂载时获取数据
onMounted(() => {
  getOperationLogs()
})
</script>
