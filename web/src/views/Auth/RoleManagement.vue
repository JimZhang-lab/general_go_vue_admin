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
                v-model="searchForm.roleName"
                type="text"
                placeholder="搜索角色名称..."
                icon="search"
              />
            </div>
            <div class="flex-1">
              <AuthInput
                v-model="searchForm.roleKey"
                type="text"
                placeholder="搜索角色标识..."
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
              @click="searchRoles"
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
              text="添加角色"
            />
          </div>
        </div>
      </AuthCard>

      <!-- 角色列表 -->
      <AuthCard title="角色列表">
        <div class="overflow-hidden rounded-xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-white/[0.03]">
          <div class="max-w-full overflow-x-auto custom-scrollbar">
            <table class="min-w-full">
              <thead>
                <tr class="border-b border-gray-200 dark:border-gray-700">
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">ID</p>
                  </th>
                  <th class="px-5 py-3 text-left sm:px-6">
                    <p class="font-medium text-gray-500 text-theme-xs dark:text-gray-400">角色信息</p>
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
                <tr v-for="role in roleList" :key="role.id" class="border-t border-gray-100 dark:border-gray-800">
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-500 text-theme-sm dark:text-gray-400">#{{ role.id }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div class="flex items-center gap-3">
                      <div class="w-10 h-10 overflow-hidden rounded-full bg-brand-50 dark:bg-brand-500/15 flex items-center justify-center">
                        <svg class="w-5 h-5 text-brand-600 dark:text-brand-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z" />
                        </svg>
                      </div>
                      <div>
                        <span class="block font-medium text-gray-800 text-theme-sm dark:text-white/90">
                          {{ role.roleName }}
                        </span>
                        <span class="block text-gray-500 text-theme-xs dark:text-gray-400">
                          {{ role.roleKey }}
                        </span>
                      </div>
                    </div>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-800 text-theme-sm dark:text-white/90">{{ role.sort }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <span
                      :class="[
                        'rounded-full px-2 py-0.5 text-theme-xs font-medium',
                        {
                          'bg-success-50 text-success-700 dark:bg-success-500/15 dark:text-success-500':
                            role.status === '1',
                          'bg-error-50 text-error-700 dark:bg-error-500/15 dark:text-error-500':
                            role.status === '2',
                        },
                      ]"
                    >
                      {{ role.status === '1' ? '启用' : '禁用' }}
                    </span>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <p class="text-gray-500 text-theme-sm dark:text-gray-400">{{ formatDate(role.createTime) }}</p>
                  </td>
                  <td class="px-5 py-4 sm:px-6">
                    <div class="flex items-center gap-2">
                      <button
                        @click="editRole(role)"
                        class="inline-flex items-center px-2 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 dark:bg-blue-500/15 dark:text-blue-400 dark:hover:bg-blue-500/25"
                      >
                        编辑
                      </button>
                      <button
                        @click="assignPermissions(role)"
                        class="inline-flex items-center px-2 py-1 text-xs font-medium text-purple-600 bg-purple-50 rounded hover:bg-purple-100 dark:bg-purple-500/15 dark:text-purple-400 dark:hover:bg-purple-500/25"
                      >
                        权限
                      </button>
                      <button
                        @click="deleteRole(role)"
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
    </div>

    <!-- 添加/编辑角色模态框 -->
    <div v-if="showAddModal || showEditModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-20 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-1/2 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">
            {{ showAddModal ? '添加角色' : '编辑角色' }}
          </h3>
          
          <form @submit.prevent="submitForm" class="space-y-4">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">角色名称 *</label>
                <input
                  v-model="roleForm.roleName"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">角色标识 *</label>
                <input
                  v-model="roleForm.roleKey"
                  type="text"
                  required
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
                <input
                  v-model.number="roleForm.sort"
                  type="number"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
                <select
                  v-model="roleForm.status"
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
                v-model="roleForm.remark"
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

    <!-- 权限分配模态框 -->
    <div v-if="showPermissionModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50">
      <div class="relative top-10 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-2/3 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <h3 class="text-lg font-medium text-gray-900 mb-4">
            分配权限 - {{ currentRole?.roleName }}
          </h3>
          
          <div class="max-h-96 overflow-y-auto">
            <div class="space-y-2">
              <div v-for="menu in menuTree" :key="menu.id" class="border rounded-lg p-3">
                <div class="flex items-center space-x-2">
                  <input
                    :id="`menu-${menu.id}`"
                    type="checkbox"
                    :checked="isMenuChecked(menu.id)"
                    @change="toggleMenu(menu.id, $event)"
                    class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                  />
                  <label :for="`menu-${menu.id}`" class="text-sm font-medium text-gray-900">
                    {{ menu.menuName }}
                  </label>
                </div>
                
                <div v-if="menu.children && menu.children.length > 0" class="ml-6 mt-2 space-y-1">
                  <div v-for="child in menu.children" :key="child.id" class="flex items-center space-x-2">
                    <input
                      :id="`menu-${child.id}`"
                      type="checkbox"
                      :checked="isMenuChecked(child.id)"
                      @change="toggleMenu(child.id, $event)"
                      class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded"
                    />
                    <label :for="`menu-${child.id}`" class="text-sm text-gray-700">
                      {{ child.menuName }}
                    </label>
                  </div>
                </div>
              </div>
            </div>
          </div>
          
          <div class="flex justify-end space-x-3 pt-4 mt-4 border-t">
            <button
              type="button"
              @click="showPermissionModal = false"
              class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50"
            >
              取消
            </button>
            <button
              @click="savePermissions"
              class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700"
            >
              保存
            </button>
          </div>
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

const currentPageTitle = ref('角色管理')

// 定义接口类型
interface Role {
  id: number
  roleName: string
  roleKey: string
  sort: number
  status: string
  remark?: string
  createTime?: string
}

interface Menu {
  id: number
  menuName: string
  parentId: number
  children?: Menu[]
}

// 响应式数据
const roleList = ref<Role[]>([])
const menuTree = ref<Menu[]>([])
const loading = ref(false)
const currentRole = ref<Role | null>(null)
const selectedMenuIds = ref<number[]>([])

// 搜索表单
const searchForm = reactive({
  roleName: '',
  roleKey: '',
  status: ''
})

// 分页信息
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  pages: 0
})

// 模态框状态
const showAddModal = ref(false)
const showEditModal = ref(false)
const showPermissionModal = ref(false)

// 表单数据
const roleForm = reactive({
  id: undefined as number | undefined,
  roleName: '',
  roleKey: '',
  sort: 0,
  status: '1',
  remark: ''
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

// 获取角色列表
const getRoleList = async () => {
  try {
    loading.value = true

    // 暂时使用模拟数据
    const mockRoles = [
      {
        id: 1,
        roleName: '超级管理员',
        roleKey: 'admin',
        sort: 1,
        status: '1',
        remark: '超级管理员角色',
        createTime: '2025-01-01 10:00:00'
      },
      {
        id: 2,
        roleName: '普通用户',
        roleKey: 'user',
        sort: 2,
        status: '1',
        remark: '普通用户角色',
        createTime: '2025-01-02 10:00:00'
      },
      {
        id: 3,
        roleName: '运营人员',
        roleKey: 'operator',
        sort: 3,
        status: '1',
        remark: '运营人员角色',
        createTime: '2025-01-03 10:00:00'
      }
    ]

    roleList.value = mockRoles
    pagination.total = mockRoles.length
    pagination.pages = 1

    ToastAlert.success({
      title: '获取成功',
      message: '角色列表已加载（模拟数据）'
    })
  } catch (error) {
    console.error('获取角色列表失败:', error)
    ToastAlert.error({
      title: '获取角色列表失败',
      message: '网络异常，请重试'
    })
  } finally {
    loading.value = false
  }
}

// 获取菜单树
const getMenuTree = async () => {
  try {
    const { data: res } = await adminApi.getMenuList({})
    if (res.code === 200) {
      menuTree.value = buildMenuTree(res.data.list || [])
    }
  } catch (error) {
    console.error('获取菜单列表失败:', error)
  }
}

// 构建菜单树
const buildMenuTree = (menuList: Menu[]): Menu[] => {
  const menuMap = new Map<number, Menu>()
  const rootMenus: Menu[] = []

  // 创建菜单映射
  menuList.forEach(menu => {
    menuMap.set(menu.id, { ...menu, children: [] })
  })

  // 构建树结构
  menuList.forEach(menu => {
    const menuItem = menuMap.get(menu.id)!
    if (menu.parentId === 0) {
      rootMenus.push(menuItem)
    } else {
      const parent = menuMap.get(menu.parentId)
      if (parent) {
        parent.children = parent.children || []
        parent.children.push(menuItem)
      }
    }
  })

  return rootMenus
}

// 搜索角色
const searchRoles = () => {
  pagination.page = 1
  getRoleList()
}

// 重置搜索
const resetSearch = () => {
  Object.assign(searchForm, {
    roleName: '',
    roleKey: '',
    status: ''
  })
  pagination.page = 1
  getRoleList()
}

// 编辑角色
const editRole = (role: Role) => {
  Object.assign(roleForm, {
    id: role.id,
    roleName: role.roleName,
    roleKey: role.roleKey,
    sort: role.sort,
    status: role.status,
    remark: role.remark || ''
  })
  showEditModal.value = true
}

// 分配权限
const assignPermissions = async (role: Role) => {
  currentRole.value = role
  // 这里应该获取角色已有的权限
  selectedMenuIds.value = []
  showPermissionModal.value = true
}

// 删除角色
const deleteRole = async (role: Role) => {
  if (!confirm(`确定要删除角色 "${role.roleName}" 吗？`)) {
    return
  }

  try {
    const { data: res } = await adminApi.deleteRole(role.id)

    if (res.code === 200) {
      ToastAlert.success({
        title: '删除成功',
        message: '角色已删除'
      })
      getRoleList()
    } else {
      ToastAlert.error({
        title: '删除失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('删除角色失败:', error)
    ToastAlert.error({
      title: '删除失败',
      message: '网络异常，请重试'
    })
  }
}

// 提交表单
const submitForm = async () => {
  try {
    const isAdd = showAddModal.value
    const apiMethod = isAdd ? adminApi.addRole : adminApi.updateRole

    const { data: res } = await apiMethod(roleForm)

    if (res.code === 200) {
      ToastAlert.success({
        title: '操作成功',
        message: `角色已${isAdd ? '添加' : '更新'}`
      })
      closeModal()
      getRoleList()
    } else {
      ToastAlert.error({
        title: '操作失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('提交表单失败:', error)
    ToastAlert.error({
      title: '操作失败',
      message: '网络异常，请重试'
    })
  }
}

// 检查菜单是否被选中
const isMenuChecked = (menuId: number): boolean => {
  return selectedMenuIds.value.includes(menuId)
}

// 切换菜单选择状态
const toggleMenu = (menuId: number, event: Event) => {
  const target = event.target as HTMLInputElement
  if (target.checked) {
    if (!selectedMenuIds.value.includes(menuId)) {
      selectedMenuIds.value.push(menuId)
    }
  } else {
    const index = selectedMenuIds.value.indexOf(menuId)
    if (index > -1) {
      selectedMenuIds.value.splice(index, 1)
    }
  }
}

// 保存权限
const savePermissions = async () => {
  try {
    // 这里应该调用保存角色权限的API
    ToastAlert.success({
      title: '保存成功',
      message: '权限已保存'
    })
    showPermissionModal.value = false
  } catch (error) {
    console.error('保存权限失败:', error)
    ToastAlert.error({
      title: '保存失败',
      message: '网络异常，请重试'
    })
  }
}

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  Object.assign(roleForm, {
    id: undefined,
    roleName: '',
    roleKey: '',
    sort: 0,
    status: '1',
    remark: ''
  })
}

// 分页操作
const prevPage = () => {
  if (pagination.page > 1) {
    pagination.page--
    getRoleList()
  }
}

const nextPage = () => {
  if (pagination.page < pagination.pages) {
    pagination.page++
    getRoleList()
  }
}

const goToPage = (page: number | string) => {
  if (typeof page === 'number' && page !== pagination.page) {
    pagination.page = page
    getRoleList()
  }
}

// 组件挂载时获取数据
onMounted(() => {
  getRoleList()
  getMenuTree()
})
</script>
