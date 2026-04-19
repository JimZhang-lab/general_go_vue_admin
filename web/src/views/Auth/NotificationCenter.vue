<template>
  <AuthLayout>
    <PageBreadcrumb :pageTitle="currentPageTitle" />

    <div class="space-y-5 sm:space-y-6">
      <div class="grid gap-4 md:grid-cols-3">
        <AuthCard title="未读通知">
          <p class="text-3xl font-semibold text-gray-900 dark:text-white/90">{{ unreadCount }}</p>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">需要处理的站内通知数量</p>
        </AuthCard>
        <AuthCard title="已发布通知">
          <p class="text-3xl font-semibold text-gray-900 dark:text-white/90">{{ publishedCount }}</p>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">当前系统已发布的通知总量</p>
        </AuthCard>
        <AuthCard title="当前列表总数">
          <p class="text-3xl font-semibold text-gray-900 dark:text-white/90">{{ pagination.total }}</p>
          <p class="mt-2 text-sm text-gray-500 dark:text-gray-400">通知管理查询结果总数</p>
        </AuthCard>
      </div>

      <AuthCard title="我的通知" subtitle="支持未读筛选、逐条已读和全部已读">
        <template #headerActions>
          <div class="flex flex-wrap gap-2">
            <AuthButton
              :text="currentUnreadOnly ? '查看全部' : '仅看未读'"
              variant="secondary"
              size="sm"
              @click="toggleUnreadFilter"
            />
            <AuthButton
              text="全部已读"
              variant="primary"
              size="sm"
              :loading="readingAll"
              loadingText="处理中..."
              @click="markAllCurrentAsRead"
            />
          </div>
        </template>

        <div v-if="currentLoading" class="space-y-3">
          <div
            v-for="index in 3"
            :key="index"
            class="h-24 animate-pulse rounded-2xl border border-gray-200 bg-gray-50 dark:border-gray-800 dark:bg-gray-800/40"
          ></div>
        </div>

        <div v-else-if="currentNotices.length === 0" class="rounded-2xl border border-dashed border-gray-300 px-6 py-10 text-center text-sm text-gray-500 dark:border-gray-700 dark:text-gray-400">
          当前没有可展示的站内通知。
        </div>

        <div v-else class="space-y-3">
          <div
            v-for="notice in currentNotices"
            :key="notice.id"
            class="rounded-2xl border p-4 transition-colors"
            :class="notice.read ? 'border-gray-200 dark:border-gray-800' : 'border-brand-200 bg-brand-50/60 dark:border-brand-500/30 dark:bg-brand-500/10'"
          >
            <div class="flex flex-col gap-3 lg:flex-row lg:items-start lg:justify-between">
              <div class="min-w-0">
                <div class="flex flex-wrap items-center gap-2">
                  <h3 class="text-sm font-semibold text-gray-900 dark:text-white/90">{{ notice.title }}</h3>
                  <span class="rounded-full px-2.5 py-1 text-xs font-medium" :class="levelBadgeClass(notice.noticeLevel)">
                    {{ levelLabel(notice.noticeLevel) }}
                  </span>
                  <span
                    v-if="!notice.read"
                    class="rounded-full bg-brand-100 px-2.5 py-1 text-xs font-medium text-brand-700 dark:bg-brand-500/15 dark:text-brand-300"
                  >
                    未读
                  </span>
                </div>
                <p class="mt-2 text-sm leading-6 text-gray-600 dark:text-gray-300">{{ notice.content }}</p>
                <div class="mt-3 flex flex-wrap items-center gap-3 text-xs text-gray-500 dark:text-gray-400">
                  <span>类型：{{ noticeTypeLabel(notice.noticeType) }}</span>
                  <span>发布人：{{ notice.createdByName || '系统' }}</span>
                  <span>发布时间：{{ formatDate(notice.publishTime || notice.createTime) }}</span>
                </div>
              </div>
              <div class="flex gap-2">
                <AuthButton
                  v-if="!notice.read"
                  text="标记已读"
                  variant="secondary"
                  size="sm"
                  @click="markCurrentAsRead(notice.id)"
                />
              </div>
            </div>
          </div>
        </div>
      </AuthCard>

      <AuthCard title="通知管理" subtitle="查询、发布、归档、编辑和删除系统通知">
        <div class="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
          <div class="grid flex-1 gap-3 md:grid-cols-4">
            <AuthInput v-model="searchForm.title" type="text" placeholder="搜索通知标题" />
            <select v-model="searchForm.noticeType" class="h-11 rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90">
              <option value="">全部类型</option>
              <option value="system">系统</option>
              <option value="security">安全</option>
              <option value="maintenance">维护</option>
            </select>
            <select v-model="searchForm.status" class="h-11 rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90">
              <option value="">全部状态</option>
              <option value="1">草稿</option>
              <option value="2">已发布</option>
              <option value="3">已归档</option>
            </select>
            <div class="flex gap-3">
              <AuthButton text="搜索" variant="primary" @click="searchNotices" />
              <AuthButton text="重置" variant="secondary" @click="resetSearch" />
            </div>
          </div>

          <div class="flex flex-wrap gap-3">
            <AuthButton text="新增通知" variant="success" @click="openCreateModal" />
            <AuthButton
              :text="`批量删除${selectedIds.length ? `(${selectedIds.length})` : ''}`"
              variant="danger"
              :disabled="selectedIds.length === 0"
              :loading="deleting"
              loadingText="删除中..."
              @click="batchDeleteSelected"
            />
          </div>
        </div>

        <div class="mt-5 overflow-hidden rounded-2xl border border-gray-200 dark:border-gray-800">
          <div class="max-w-full overflow-x-auto">
            <table class="min-w-full">
              <thead class="bg-gray-50 dark:bg-gray-900/40">
                <tr class="border-b border-gray-200 dark:border-gray-800">
                  <th class="px-5 py-3 text-left">
                    <input
                      type="checkbox"
                      class="h-4 w-4 rounded border-gray-300 text-brand-600"
                      :checked="noticeList.length > 0 && selectedIds.length === noticeList.length"
                      @change="toggleAllSelection"
                    />
                  </th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">标题</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">类型</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">等级</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">状态</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">发布时间</th>
                  <th class="px-5 py-3 text-left text-xs font-medium uppercase tracking-wider text-gray-500">操作</th>
                </tr>
              </thead>
              <tbody v-if="listLoading">
                <tr v-for="index in 4" :key="index" class="border-t border-gray-100 dark:border-gray-800">
                  <td colspan="7" class="px-5 py-4">
                    <div class="h-12 animate-pulse rounded-xl bg-gray-100 dark:bg-gray-800"></div>
                  </td>
                </tr>
              </tbody>
              <tbody v-else-if="noticeList.length === 0">
                <tr>
                  <td colspan="7" class="px-5 py-10 text-center text-sm text-gray-500 dark:text-gray-400">
                    当前筛选条件下没有通知记录。
                  </td>
                </tr>
              </tbody>
              <tbody v-else class="divide-y divide-gray-200 dark:divide-gray-800">
                <tr v-for="notice in noticeList" :key="notice.id">
                  <td class="px-5 py-4">
                    <input
                      type="checkbox"
                      class="h-4 w-4 rounded border-gray-300 text-brand-600"
                      :checked="selectedIds.includes(notice.id)"
                      @change="toggleSelection(notice.id)"
                    />
                  </td>
                  <td class="px-5 py-4">
                    <div class="max-w-sm">
                      <p class="text-sm font-medium text-gray-900 dark:text-white/90">{{ notice.title }}</p>
                      <p class="mt-1 line-clamp-2 text-xs text-gray-500 dark:text-gray-400">{{ notice.content }}</p>
                    </div>
                  </td>
                  <td class="px-5 py-4 text-sm text-gray-600 dark:text-gray-300">{{ noticeTypeLabel(notice.noticeType) }}</td>
                  <td class="px-5 py-4">
                    <span class="rounded-full px-2.5 py-1 text-xs font-medium" :class="levelBadgeClass(notice.noticeLevel)">
                      {{ levelLabel(notice.noticeLevel) }}
                    </span>
                  </td>
                  <td class="px-5 py-4">
                    <span class="rounded-full px-2.5 py-1 text-xs font-medium" :class="statusBadgeClass(notice.status)">
                      {{ statusLabel(notice.status) }}
                    </span>
                  </td>
                  <td class="px-5 py-4 text-sm text-gray-500 dark:text-gray-400">{{ formatDate(notice.publishTime || notice.createTime) }}</td>
                  <td class="px-5 py-4">
                    <div class="flex flex-wrap gap-2">
                      <button class="text-sm font-medium text-brand-600 hover:text-brand-700" @click="editNotice(notice)">编辑</button>
                      <button class="text-sm font-medium text-emerald-600 hover:text-emerald-700" @click="toggleNoticeStatus(notice)">
                        {{ notice.status === 2 ? '归档' : '发布' }}
                      </button>
                      <button class="text-sm font-medium text-red-600 hover:text-red-700" @click="deleteSingleNotice(notice.id)">删除</button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <template #footer>
          <div class="flex flex-wrap items-center justify-between gap-3">
            <p class="text-sm text-gray-500 dark:text-gray-400">
              第 {{ pagination.page }} / {{ Math.max(pagination.pages, 1) }} 页，共 {{ pagination.total }} 条通知
            </p>
            <div class="flex gap-3">
              <AuthButton text="上一页" variant="secondary" :disabled="pagination.page <= 1" @click="prevPage" />
              <AuthButton text="下一页" variant="secondary" :disabled="pagination.page >= pagination.pages" @click="nextPage" />
            </div>
          </div>
        </template>
      </AuthCard>
    </div>

    <div v-if="showModal" class="fixed inset-0 z-99999 flex items-center justify-center p-4">
      <div class="absolute inset-0 bg-black/50" @click="closeModal"></div>
      <div class="relative w-full max-w-3xl rounded-2xl border border-gray-200 bg-white p-6 shadow-2xl dark:border-gray-800 dark:bg-gray-900">
        <div class="flex items-center justify-between">
          <div>
            <h3 class="text-lg font-semibold text-gray-900 dark:text-white/90">{{ isEditing ? '编辑通知' : '新增通知' }}</h3>
            <p class="mt-1 text-sm text-gray-500 dark:text-gray-400">通知内容将持久化保存到数据库，可直接发布到站内通知中心。</p>
          </div>
          <button class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300" @click="closeModal">关闭</button>
        </div>

        <div class="mt-6 grid gap-4 md:grid-cols-2">
          <div class="md:col-span-2">
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">标题</label>
            <input v-model="noticeForm.title" type="text" class="h-11 w-full rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90" />
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">通知类型</label>
            <select v-model="noticeForm.noticeType" class="h-11 w-full rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90">
              <option value="system">系统</option>
              <option value="security">安全</option>
              <option value="maintenance">维护</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">通知等级</label>
            <select v-model="noticeForm.noticeLevel" class="h-11 w-full rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90">
              <option value="info">普通</option>
              <option value="warning">提醒</option>
              <option value="critical">紧急</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">目标范围</label>
            <select v-model="noticeForm.targetType" class="h-11 w-full rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90">
              <option value="all">全部管理员</option>
              <option value="admin">仅管理员</option>
            </select>
          </div>
          <div>
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">保存状态</label>
            <select v-model="noticeForm.status" class="h-11 w-full rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90">
              <option :value="1">草稿</option>
              <option :value="2">立即发布</option>
              <option :value="3">归档</option>
            </select>
          </div>
          <div class="md:col-span-2">
            <label class="mb-2 block text-sm font-medium text-gray-700 dark:text-gray-300">通知内容</label>
            <textarea v-model="noticeForm.content" rows="6" class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"></textarea>
          </div>
        </div>

        <div class="mt-6 flex justify-end gap-3">
          <AuthButton text="取消" variant="secondary" @click="closeModal" />
          <AuthButton
            :text="isEditing ? '保存变更' : '创建通知'"
            variant="primary"
            :loading="submitting"
            loadingText="提交中..."
            @click="submitNotice"
          />
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { AuthButton, AuthCard, AuthInput, AuthLayout } from '@/components/auth'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'

interface NoticeItem {
  id: number
  title: string
  content: string
  noticeType: string
  noticeLevel: string
  status: number
  targetType: string
  createdByName?: string
  publishTime?: string
  createTime?: string
  read?: boolean
}

const currentPageTitle = ref('通知中心')
const currentNotices = ref<NoticeItem[]>([])
const noticeList = ref<NoticeItem[]>([])
const selectedIds = ref<number[]>([])
const currentUnreadOnly = ref(false)
const currentLoading = ref(false)
const listLoading = ref(false)
const submitting = ref(false)
const deleting = ref(false)
const readingAll = ref(false)
const showModal = ref(false)
const isEditing = ref(false)

const searchForm = reactive({
  title: '',
  noticeType: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  pages: 0
})

const noticeForm = reactive({
  id: 0,
  title: '',
  content: '',
  noticeType: 'system',
  noticeLevel: 'info',
  status: 1,
  targetType: 'all'
})

const unreadCount = computed(() => currentNotices.value.filter((item) => !item.read).length)
const publishedCount = computed(() => noticeList.value.filter((item) => item.status === 2).length)

const resolveListPayload = <T>(payload: unknown): T[] => {
  if (Array.isArray(payload)) {
    return payload as T[]
  }
  if (payload && typeof payload === 'object') {
    const list = (payload as { list?: unknown }).list
    if (Array.isArray(list)) {
      return list as T[]
    }
  }
  return []
}

const formatDate = (value?: string) => {
  if (!value) {
    return '-'
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return value
  }
  return `${date.toLocaleDateString('zh-CN')} ${date.toLocaleTimeString('zh-CN')}`
}

const noticeTypeLabel = (type?: string) => {
  switch (type) {
    case 'security':
      return '安全通知'
    case 'maintenance':
      return '维护通知'
    default:
      return '系统通知'
  }
}

const levelLabel = (level?: string) => {
  switch (level) {
    case 'critical':
      return '紧急'
    case 'warning':
      return '提醒'
    default:
      return '普通'
  }
}

const levelBadgeClass = (level?: string) => {
  switch (level) {
    case 'critical':
      return 'bg-red-100 text-red-700 dark:bg-red-500/15 dark:text-red-300'
    case 'warning':
      return 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300'
    default:
      return 'bg-blue-100 text-blue-700 dark:bg-blue-500/15 dark:text-blue-300'
  }
}

const statusLabel = (status: number) => {
  switch (status) {
    case 2:
      return '已发布'
    case 3:
      return '已归档'
    default:
      return '草稿'
  }
}

const statusBadgeClass = (status: number) => {
  switch (status) {
    case 2:
      return 'bg-emerald-100 text-emerald-700 dark:bg-emerald-500/15 dark:text-emerald-300'
    case 3:
      return 'bg-gray-100 text-gray-700 dark:bg-gray-500/15 dark:text-gray-300'
    default:
      return 'bg-amber-100 text-amber-700 dark:bg-amber-500/15 dark:text-amber-300'
  }
}

const resetForm = () => {
  Object.assign(noticeForm, {
    id: 0,
    title: '',
    content: '',
    noticeType: 'system',
    noticeLevel: 'info',
    status: 1,
    targetType: 'all'
  })
}

const loadCurrentNotices = async () => {
  try {
    currentLoading.value = true
    const { data: res } = await adminApi.getCurrentNotices({
      limit: 8,
      unreadOnly: currentUnreadOnly.value
    })
    if (res.code !== 200) {
      throw new Error(res.message || '获取当前通知失败')
    }
    currentNotices.value = resolveListPayload<NoticeItem>(res.data)
  } catch (error) {
    currentNotices.value = []
    ToastAlert.error({
      title: '通知加载失败',
      message: error instanceof Error ? error.message : '获取当前通知失败'
    })
  } finally {
    currentLoading.value = false
  }
}

const loadNoticeList = async () => {
  try {
    listLoading.value = true
    const { data: res } = await adminApi.getNoticeList({
      pageNum: pagination.page,
      pageSize: pagination.pageSize,
      title: searchForm.title,
      noticeType: searchForm.noticeType,
      status: searchForm.status
    })
    if (res.code !== 200) {
      throw new Error(res.message || '获取通知列表失败')
    }
    noticeList.value = resolveListPayload<NoticeItem>(res.data)
    pagination.total = Number((res.data as { total?: number }).total || 0)
    pagination.pages = Math.ceil(pagination.total / pagination.pageSize)
    selectedIds.value = []
  } catch (error) {
    noticeList.value = []
    pagination.total = 0
    pagination.pages = 0
    ToastAlert.error({
      title: '列表加载失败',
      message: error instanceof Error ? error.message : '获取通知列表失败'
    })
  } finally {
    listLoading.value = false
  }
}

const toggleUnreadFilter = async () => {
  currentUnreadOnly.value = !currentUnreadOnly.value
  await loadCurrentNotices()
}

const markCurrentAsRead = async (id: number) => {
  try {
    const { data: res } = await adminApi.markNoticeRead(id)
    if (res.code !== 200) {
      throw new Error(res.message || '标记已读失败')
    }
    await loadCurrentNotices()
  } catch (error) {
    ToastAlert.error({
      title: '操作失败',
      message: error instanceof Error ? error.message : '标记已读失败'
    })
  }
}

const markAllCurrentAsRead = async () => {
  try {
    readingAll.value = true
    const { data: res } = await adminApi.markAllNoticeRead()
    if (res.code !== 200) {
      throw new Error(res.message || '全部已读失败')
    }
    ToastAlert.success({
      title: '操作成功',
      message: '当前通知已全部标记为已读'
    })
    await loadCurrentNotices()
  } catch (error) {
    ToastAlert.error({
      title: '操作失败',
      message: error instanceof Error ? error.message : '全部已读失败'
    })
  } finally {
    readingAll.value = false
  }
}

const searchNotices = async () => {
  pagination.page = 1
  await loadNoticeList()
}

const resetSearch = async () => {
  Object.assign(searchForm, {
    title: '',
    noticeType: '',
    status: ''
  })
  pagination.page = 1
  await loadNoticeList()
}

const openCreateModal = () => {
  isEditing.value = false
  resetForm()
  showModal.value = true
}

const editNotice = (notice: NoticeItem) => {
  isEditing.value = true
  Object.assign(noticeForm, {
    id: notice.id,
    title: notice.title,
    content: notice.content,
    noticeType: notice.noticeType,
    noticeLevel: notice.noticeLevel,
    status: notice.status,
    targetType: notice.targetType
  })
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  resetForm()
}

const submitNotice = async () => {
  if (noticeForm.title.trim().length < 2 || noticeForm.content.trim().length < 2) {
    ToastAlert.warning({
      title: '内容不完整',
      message: '请完整填写通知标题和内容'
    })
    return
  }

  try {
    submitting.value = true
    const payload = {
      id: noticeForm.id || undefined,
      title: noticeForm.title.trim(),
      content: noticeForm.content.trim(),
      noticeType: noticeForm.noticeType,
      noticeLevel: noticeForm.noticeLevel,
      status: Number(noticeForm.status),
      targetType: noticeForm.targetType
    }
    const apiCall = isEditing.value ? adminApi.updateNotice(payload) : adminApi.addNotice(payload)
    const { data: res } = await apiCall
    if (res.code !== 200) {
      throw new Error(res.message || '通知保存失败')
    }
    ToastAlert.success({
      title: '保存成功',
      message: `通知已${isEditing.value ? '更新' : '创建'}`
    })
    closeModal()
    await Promise.all([loadCurrentNotices(), loadNoticeList()])
  } catch (error) {
    ToastAlert.error({
      title: '保存失败',
      message: error instanceof Error ? error.message : '通知保存失败'
    })
  } finally {
    submitting.value = false
  }
}

const toggleNoticeStatus = async (notice: NoticeItem) => {
  const nextStatus = notice.status === 2 ? 3 : 2
  try {
    const { data: res } = await adminApi.updateNoticeStatus({
      id: notice.id,
      status: nextStatus
    })
    if (res.code !== 200) {
      throw new Error(res.message || '通知状态更新失败')
    }
    ToastAlert.success({
      title: '操作成功',
      message: `通知已${nextStatus === 2 ? '发布' : '归档'}`
    })
    await Promise.all([loadCurrentNotices(), loadNoticeList()])
  } catch (error) {
    ToastAlert.error({
      title: '操作失败',
      message: error instanceof Error ? error.message : '通知状态更新失败'
    })
  }
}

const deleteSingleNotice = async (id: number) => {
  if (!window.confirm('确定要删除这条通知吗？')) {
    return
  }
  try {
    const { data: res } = await adminApi.deleteNotice(id)
    if (res.code !== 200) {
      throw new Error(res.message || '通知删除失败')
    }
    ToastAlert.success({
      title: '删除成功',
      message: '通知已删除'
    })
    await Promise.all([loadCurrentNotices(), loadNoticeList()])
  } catch (error) {
    ToastAlert.error({
      title: '删除失败',
      message: error instanceof Error ? error.message : '通知删除失败'
    })
  }
}

const toggleSelection = (id: number) => {
  selectedIds.value = selectedIds.value.includes(id)
    ? selectedIds.value.filter((item) => item !== id)
    : [...selectedIds.value, id]
}

const toggleAllSelection = () => {
  selectedIds.value = selectedIds.value.length === noticeList.value.length ? [] : noticeList.value.map((item) => item.id)
}

const batchDeleteSelected = async () => {
  if (selectedIds.value.length === 0) {
    return
  }
  if (!window.confirm(`确定删除已选中的 ${selectedIds.value.length} 条通知吗？`)) {
    return
  }

  try {
    deleting.value = true
    const { data: res } = await adminApi.batchDeleteNotices(selectedIds.value)
    if (res.code !== 200) {
      throw new Error(res.message || '批量删除失败')
    }
    ToastAlert.success({
      title: '删除成功',
      message: '通知已批量删除'
    })
    await Promise.all([loadCurrentNotices(), loadNoticeList()])
  } catch (error) {
    ToastAlert.error({
      title: '删除失败',
      message: error instanceof Error ? error.message : '批量删除失败'
    })
  } finally {
    deleting.value = false
  }
}

const prevPage = async () => {
  if (pagination.page <= 1) {
    return
  }
  pagination.page -= 1
  await loadNoticeList()
}

const nextPage = async () => {
  if (pagination.page >= pagination.pages) {
    return
  }
  pagination.page += 1
  await loadNoticeList()
}

onMounted(async () => {
  await Promise.all([loadCurrentNotices(), loadNoticeList()])
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}
</style>
