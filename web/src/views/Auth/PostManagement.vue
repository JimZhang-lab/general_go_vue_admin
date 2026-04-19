/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/views/Auth/PostManagement.vue
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
                v-model="searchForm.postName"
                type="text"
                placeholder="搜索岗位名称..."
                icon="search"
              />
            </div>
            <div class="flex-1">
              <AuthInput
                v-model="searchForm.postCode"
                type="text"
                placeholder="搜索岗位编号..."
                icon="key"
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
              @click="searchPosts"
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
              text="添加岗位"
            />
          </div>
        </div>
      </AuthCard>

      <!-- 岗位列表 -->
      <AuthCard title="岗位列表">
        <div class="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-white/[0.03]">
          <div class="max-w-full overflow-x-auto custom-scrollbar">
            <table class="min-w-full">
              <thead>
                <tr class="border-b border-gray-200 dark:border-gray-700">
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">ID</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">岗位名称</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">岗位编号</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">排序</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">状态</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">创建时间</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">操作</p>
                  </th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200 dark:divide-gray-700">
                <tr v-if="loading"><td colspan="7" class="text-center py-4 text-gray-500">加载中...</td></tr>
                <tr v-else-if="postList.length === 0"><td colspan="7" class="text-center py-4 text-gray-500">无数据</td></tr>
                <tr v-else v-for="post in postList" :key="post.id" class="border-t border-gray-100 dark:border-gray-800">
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-500 text-theme-sm dark:text-gray-400">#{{ post.id }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div class="flex items-center gap-3">
                      <div class="w-10 h-10 overflow-hidden rounded-full bg-blue-50 dark:bg-blue-500/15 flex items-center justify-center">
                        <svg class="w-5 h-5 text-blue-600 dark:text-blue-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                        </svg>
                      </div>
                      <div>
                        <span class="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
                          {{ post.postName }}
                        </span>
                      </div>
                    </div>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-800 text-theme-sm dark:text-white/90">{{ post.postCode }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-800 text-theme-sm dark:text-white/90">{{ post.sort }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <span
                      :class="[
                        'rounded-full px-2 py-0.5 text-theme-xs font-medium',
                        {
                          'bg-success-50 text-success-700 dark:bg-success-500/15 dark:text-success-500':
                            post.status === '1',
                          'bg-error-50 text-error-700 dark:bg-error-500/15 dark:text-error-500':
                            post.status === '2',
                        },
                      ]"
                    >
                      {{ post.status === '1' ? '启用' : '禁用' }}
                    </span>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ formatDate(post.createTime) }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div class="flex items-center gap-2">
                      <button
                        @click="editPost(post)"
                        class="inline-flex items-center px-2 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 dark:bg-blue-500/15 dark:text-blue-400 dark:hover:bg-blue-500/25"
                      >
                        编辑
                      </button>
                      <button
                        @click="deletePost(post)"
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

    <!-- 添加/编辑岗位模态框 -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 overflow-y-auto h-full w-full z-99999">
      <div class="fixed inset-0 bg-black/50" @click="closeModal"></div>
      <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <div class="flex justify-between items-center">
            <h3 class="text-lg font-medium text-gray-900 mb-4">
              {{ showAddModal ? '添加岗位' : '编辑岗位' }}
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
                <label class="block text-sm font-medium text-gray-700 mb-1">岗位名称 *</label>
                <input
                  v-model="postForm.postName"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">岗位编号 *</label>
                <input
                  v-model="postForm.postCode"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
                <input
                  v-model.number="postForm.sort"
                  type="number"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
                <select
                  v-model="postForm.status"
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
                v-model="postForm.remark"
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

const currentPageTitle = ref('岗位管理')

interface Post {
  id: number
  postName: string
  postCode: string
  sort: number
  status: string
  remark?: string
  createTime?: string
}

const postList = ref<Post[]>([])
const loading = ref(false)

const searchForm = reactive({
  postName: '',
  postCode: '',
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

const postForm = reactive({
  id: undefined as number | undefined,
  postName: '',
  postCode: '',
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

const getPostList = async () => {
  try {
    loading.value = true
    const params = {
      pageNum: pagination.page,
      pageSize: pagination.pageSize,
      postName: searchForm.postName,
      postCode: searchForm.postCode,
      status: searchForm.status
    }
    const { data: res } = await adminApi.getPostList(params)
    if (res.code === 200) {
      postList.value = res.data.list || []
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

const searchPosts = () => {
  pagination.page = 1
  getPostList()
}

const resetSearch = () => {
  Object.assign(searchForm, { postName: '', postCode: '', status: '' })
  searchPosts()
}

const editPost = (post: Post) => {
  Object.assign(postForm, {
    id: post.id,
    postName: post.postName,
    postCode: post.postCode,
    sort: post.sort,
    status: post.status,
    remark: post.remark || ''
  })
  showEditModal.value = true
}

const deletePost = async (post: Post) => {
  if (!confirm(`确定要删除岗位 "${post.postName}" 吗？`)) return
  try {
    const { data: res } = await adminApi.deletePost(post.id)
    if (res.code === 200) {
      ToastAlert.success({ title: '删除成功', message: '岗位已删除' })
      getPostList()
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
    const apiMethod = isAdd ? adminApi.addPost : adminApi.updatePost
    const { data: res } = await apiMethod(postForm as any)

    if (res.code === 200) {
      ToastAlert.success({ title: '成功', message: `岗位已${isAdd ? '添加' : '更新'}` })
      closeModal()
      getPostList()
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
  Object.assign(postForm, {
    id: undefined,
    postName: '',
    postCode: '',
    sort: 0,
    status: '1',
    remark: ''
  })
}

const prevPage = () => { if (pagination.page > 1) { pagination.page--; getPostList() } }
const nextPage = () => { if (pagination.page < pagination.pages) { pagination.page++; getPostList() } }
const goToPage = (page: number) => { if (!isNaN(page) && page !== pagination.page) { pagination.page = page; getPostList() } }

onMounted(() => {
  getPostList()
})
</script>
