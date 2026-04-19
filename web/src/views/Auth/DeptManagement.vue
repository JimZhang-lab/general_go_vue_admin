/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/views/Auth/DeptManagement.vue
 * @Description: 
 * 
 */
<template>
  <AuthLayout>
    <PageBreadcrumb :pageTitle="currentPageTitle" />
    <div class="space-y-5 sm:space-y-6">

      <!-- 搜索和操作栏 -->
      <AuthCard title="搜索条件">
        <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between" style="gap: 1rem;">
          <!-- 搜索表单 -->
          <div class="flex flex-col sm:flex-row flex-1" style="gap: 1rem;">
            <div class="flex-1">
              <AuthInput
                v-model="searchForm.deptName"
                type="text"
                placeholder="搜索部门名称..."
                icon="search"
              />
            </div>
            <div class="flex-1">
              <select
                v-model="searchForm.status"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:border-gray-700 dark:bg-gray-800 dark:text-white dark:placeholder-gray-500"
              >
                <option value="">全部状态</option>
                <option value="1">启用</option>
                <option value="2">禁用</option>
              </select>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex" style="gap: 0.5rem;">
            <AuthButton
              @click="searchDepts"
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
              @click="showAddModal = true"
              variant="success"
              size="md"
              text="添加部门"
            />
          </div>
        </div>
      </AuthCard>

      <!-- 部门列表 -->
      <AuthCard title="部门列表">
        <div class="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-white/[0.03]">
          <div class="max-w-full overflow-x-auto custom-scrollbar">
            <table class="min-w-full">
              <thead>
                <tr class="border-b border-gray-200 dark:border-gray-700">
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">ID</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">部门名称</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">父级ID</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">排序</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">状态</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">记录时间</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">操作</p>
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                <tr v-if="loading"><td colspan="7" class="text-center py-4 text-gray-500">加载中...</td></tr>
                <tr v-else-if="deptList.length === 0"><td colspan="7" class="text-center py-4 text-gray-500">无数据</td></tr>
                <tr v-else v-for="dept in deptList" :key="dept.id" class="border-t border-gray-100 dark:border-gray-800">
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-500 text-theme-sm dark:text-gray-400">#{{ dept.id }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div class="flex items-center gap-3">
                      <div class="w-10 h-10 overflow-hidden rounded-full bg-brand-50 dark:bg-brand-500/15 flex items-center justify-center">
                        <svg class="w-5 h-5 text-brand-600 dark:text-brand-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                        </svg>
                      </div>
                      <div>
                        <span class="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
                          {{ dept.deptName }}
                        </span>
                      </div>
                    </div>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-800 text-theme-sm dark:text-white/90">{{ dept.parentId === 0 ? '顶级部门' : '#' + dept.parentId }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-800 text-theme-sm dark:text-white/90">{{ dept.sort }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <span
                      :class="[
                        'rounded-full px-2 py-0.5 text-theme-xs font-medium',
                        {
                          'bg-success-50 text-success-700 dark:bg-success-500/15 dark:text-success-500':
                            dept.status === '1',
                          'bg-error-50 text-error-700 dark:bg-error-500/15 dark:text-error-500':
                            dept.status === '2',
                        },
                      ]"
                    >
                      {{ dept.status === '1' ? '启用' : '禁用' }}
                    </span>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ formatDate(dept.createTime) }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div class="flex items-center gap-2">
                      <button
                        @click="editDept(dept)"
                        class="inline-flex items-center px-2 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 dark:bg-blue-500/15 dark:text-blue-400 dark:hover:bg-blue-500/25"
                      >
                        编辑
                      </button>
                      <button
                        @click="deleteDept(dept)"
                        class="inline-flex items-center px-2 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 dark:bg-red-500/15 dark:text-red-400 dark:hover:bg-red-500/25"
                      >
                        删除
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>
      </AuthCard>

      <!-- 分页 -->
      <div class="flex items-center justify-between px-4 py-3 bg-white border border-gray-200 rounded-2xl dark:border-gray-800 dark:bg-white/[0.03]">
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
                @click="goToPage(Number(page))"
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
    </div>

    <!-- 添加/编辑部门模态框 -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 overflow-y-auto h-full w-full z-99999">
      <div class="fixed inset-0 bg-black/50" @click="closeModal"></div>
      <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-medium text-gray-900 mb-4">
              {{ showAddModal ? '添加部门' : '编辑部门' }}
            </h3>
            <button
                @click="closeModal"
                class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300"
              >
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
            </button>
          </div>

          <form @submit.prevent="submitForm" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">部门名称 *</label>
                <input
                  v-model="deptForm.deptName"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">上级部门ID</label>
                <input
                  v-model.number="deptForm.parentId"
                  type="number"
                  placeholder="顶级部门填0"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
                <input
                  v-model.number="deptForm.sort"
                  type="number"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
                <select
                  v-model="deptForm.status"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  <option value="1">启用</option>
                  <option value="2">禁用</option>
                </select>
              </div>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">备注</label>
              <textarea
                v-model="deptForm.remark"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              ></textarea>
            </div>

            <div class="flex justify-end space-x-3 pt-4">
              <button
                type="button"
                @click="closeModal"
                class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
              >
                取消
              </button>
              <button
                type="submit"
                class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
              >
                {{ showAddModal ? '添加' : '更新' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { AuthLayout, AuthCard, AuthButton, AuthInput } from '@/components/auth'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'

const currentPageTitle = ref('部门管理')

interface Dept {
  id: number
  parentId: number
  deptName: string
  sort: number
  status: string
  remark?: string
  createTime?: string
}

const deptList = ref<Dept[]>([])
const loading = ref(false)

const searchForm = reactive({
  deptName: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  pages: 0
})

const showAddModal = ref(false)
const showEditModal = ref(false)

const deptForm = reactive({
  id: undefined as number | undefined,
  parentId: 0,
  deptName: '',
  sort: 0,
  status: '1',
  remark: ''
})

const visiblePages = computed(() => {
  const pages = []
  const total = pagination.pages || 1
  const current = pagination.page

  if (total <= 7) {
    for (let i = 1; i <= total; i++) pages.push(i)
  } else {
    if (current <= 4) {
      for (let i = 1; i <= 5; i++) pages.push(i)
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 3) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) pages.push(i)
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) pages.push(i)
      pages.push('...')
      pages.push(total)
    }
  }
  return pages
})

const formatDate = (dateString?: string) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('zh-CN') + ' ' + date.toLocaleTimeString('zh-CN')
}

const getDeptList = async () => {
  try {
    loading.value = true
    const params = {
      pageNum: pagination.page,
      pageSize: pagination.pageSize,
      deptName: searchForm.deptName,
      status: searchForm.status
    }
    const { data: res } = await adminApi.getDeptList(params)
    if (res.code === 200) {
      deptList.value = res.data.list || []
      pagination.total = res.data.total || 0
      pagination.pages = Math.ceil((res.data.total || 0) / pagination.pageSize) || 1
    } else {
      ToastAlert.error({ title: '获取数据失败', message: res.message })
    }
  } catch (error) {
    ToastAlert.error({ title: '失败', message: '网络异常' })
  } finally {
    loading.value = false
  }
}

const searchDepts = () => {
  pagination.page = 1
  getDeptList()
}

const resetSearch = () => {
  Object.assign(searchForm, { deptName: '', status: '' })
  searchDepts()
}

const editDept = (dept: Dept) => {
  Object.assign(deptForm, {
    id: dept.id,
    parentId: dept.parentId,
    deptName: dept.deptName,
    sort: dept.sort,
    status: dept.status,
    remark: dept.remark || ''
  })
  showEditModal.value = true
}

const deleteDept = async (dept: Dept) => {
  if (!confirm(`确定要删除部门 "${dept.deptName}" 吗？`)) return
  try {
    const { data: res } = await adminApi.deleteDept(dept.id)
    if (res.code === 200) {
      ToastAlert.success({ title: '删除成功', message: '部门已删除' })
      getDeptList()
    } else {
      ToastAlert.error({ title: '删除失败', message: res.message })
    }
  } catch (error) {
    ToastAlert.error({ title: '删除失败', message: '网络异常' })
  }
}

const submitForm = async () => {
  try {
    const isAdd = showAddModal.value
    const apiMethod = isAdd ? adminApi.addDept : adminApi.updateDept
    const { data: res } = await apiMethod(deptForm as any)

    if (res.code === 200) {
      ToastAlert.success({ title: '成功', message: `部门已${isAdd ? '添加' : '更新'}` })
      closeModal()
      getDeptList()
    } else {
      ToastAlert.error({ title: '失败', message: res.message })
    }
  } catch (error) {
    ToastAlert.error({ title: '失败', message: '网络异常' })
  }
}

const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  Object.assign(deptForm, {
    id: undefined,
    parentId: 0,
    deptName: '',
    sort: 0,
    status: '1',
    remark: ''
  })
}

const prevPage = () => { if (pagination.page > 1) { pagination.page--; getDeptList() } }
const nextPage = () => { if (pagination.page < pagination.pages) { pagination.page++; getDeptList() } }
const goToPage = (page: number) => { if (!isNaN(page) && page !== pagination.page) { pagination.page = page; getDeptList() } }

onMounted(() => {
  getDeptList()
})
</script>
