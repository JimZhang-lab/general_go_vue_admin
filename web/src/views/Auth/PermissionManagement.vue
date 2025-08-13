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
              <AuthInput v-model="searchForm.menuName" type="text" placeholder="搜索菜单名称..." icon="search" />
            </div>
            <div class="flex-1">
              <select v-model="searchForm.status"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent dark:border-gray-700 dark:bg-gray-800 dark:text-white dark:placeholder-gray-500">
                <option value="">全部状态</option>
                <option value="1">启用</option>
                <option value="2">禁用</option>
              </select>
            </div>
          </div>

          <!-- 操作按钮 -->
          <div class="flex" style="gap: 0.5rem;">
            <AuthButton @click="searchMenus" variant="primary" size="md" text="搜索" />
            <AuthButton @click="resetSearch" variant="secondary" size="md" text="重置" />
            <AuthButton @click="showAddModal = true" variant="success" size="md" text="添加菜单" />
            <AuthButton @click="expandAll" variant="info" size="md" :text="allExpanded ? '收起全部' : '展开全部'" />
          </div>
        </div>
      </AuthCard>

      <!-- 菜单树形表格 -->
      <div class="bg-white rounded-lg shadow-sm border border-gray-200">
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">菜单名称</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">图标</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">类型</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">排序</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">权限标识</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">组件路径</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">状态</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">操作</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="row in visibleRows" :key="row.menu.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  <div class="flex items-center">
                    <button v-if="row.menu.children && row.menu.children.length"
                            @click="toggleMenu(row.menu.id)"
                            class="mr-2 text-gray-500 hover:text-gray-700">
                      <span>{{ expandedMenus.has(row.menu.id) ? '▾' : '▸' }}</span>
                    </button>
                    <span :style="{ paddingLeft: `${row.level * 16}px` }">{{ row.menu.menuName }}</span>
                  </div>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ row.menu.icon || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ row.menu.menuType === 'M' ? '目录' : row.menu.menuType === 'C' ? '菜单' : '按钮' }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ row.menu.sort }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ row.menu.perms || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">{{ row.menu.component || '-' }}</td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="row.menu.status === '1' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'"
                        class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                    {{ row.menu.status === '1' ? '启用' : '禁用' }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium space-x-2">
                  <button @click="editMenu(row.menu)" class="text-blue-600 hover:text-blue-900">编辑</button>
                  <button @click="addChildMenu(row.menu)" class="text-purple-600 hover:text-purple-900">添加子菜单</button>
                  <button @click="deleteMenu(row.menu)" class="text-red-600 hover:text-red-900">删除</button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>

      <!-- 添加/编辑菜单模态框 -->
      <div v-if="showAddModal || showEditModal" class="fixed inset-0 z-99999 flex items-center justify-center p-4">
        <div class="fixed inset-0 bg-black/50" @click="closeModal"></div>
        <div class="relative top-10 mx-auto p-5 border w-11/12 md:w-3/4 lg:w-2/3 shadow-lg rounded-md bg-white">
          <div class="mt-3">
            <div class="flex justify-between items-center">
              <h3 class="text-lg font-medium text-gray-900 mb-4">
                {{ showAddModal ? '添加菜单' : '编辑菜单' }}
              </h3>
              <button @click="closeModal" class="text-gray-400 hover:text-gray-600 dark:hover:text-gray-300">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>

            <form @submit.prevent="submitForm" class="space-y-4">
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">菜单名称 *</label>
                  <input v-model="menuForm.menuName" type="text" required
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">上级菜单</label>
                  <select v-model="menuForm.parentId"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
                    <option value="0">顶级菜单</option>
                    <option v-for="menu in flatMenuList" :key="menu.id" :value="menu.id">
                      {{ '　'.repeat(menu.level) }}{{ menu.menuName }}
                    </option>
                  </select>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">菜单类型 *</label>
                  <select v-model="menuForm.menuType"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
                    <option value="M">目录</option>
                    <option value="C">菜单</option>
                    <option value="F">按钮</option>
                  </select>
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">排序</label>
                  <input v-model.number="menuForm.sort" type="number"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>

                <div v-if="menuForm.menuType !== 'F'">
                  <label class="block text-sm font-medium text-gray-700 mb-1">图标</label>
                  <input v-model="menuForm.icon" type="text" placeholder="如: user, setting"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>

                <div v-if="menuForm.menuType === 'C'">
                  <label class="block text-sm font-medium text-gray-700 mb-1">组件路径</label>
                  <input v-model="menuForm.component" type="text" placeholder="如: system/user/index"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">权限标识</label>
                  <input v-model="menuForm.perms" type="text" placeholder="如: system:user:list"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500" />
                </div>

                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">状态</label>
                  <select v-model="menuForm.status"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500">
                    <option value="1">启用</option>
                    <option value="2">禁用</option>
                  </select>
                </div>
              </div>

              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">备注</label>
                <textarea v-model="menuForm.remark" rows="3"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"></textarea>
              </div>

              <div class="flex justify-end space-x-3 pt-4">
                <button type="button" @click="closeModal"
                  class="px-4 py-2 border border-gray-300 rounded-md text-gray-700 hover:bg-gray-50">
                  取消
                </button>
                <button type="submit" class="px-4 py-2 bg-blue-600 text-white rounded-md hover:bg-blue-700">
                  {{ showAddModal ? '添加' : '更新' }}
                </button>
              </div>
            </form>
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

const currentPageTitle = ref('权限管理')

// 定义接口类型
interface Menu {
  id: number
  menuName: string
  parentId: number
  menuType: string
  sort: number
  component?: string
  perms?: string
  icon?: string
  status: string
  remark?: string
  createTime?: string
  children?: Menu[]
  level?: number
}

// 响应式数据
const menuList = ref<Menu[]>([])
const expandedMenus = ref(new Set<number>())
const allExpanded = ref(false)
const loading = ref(false)

// 搜索表单
const searchForm = reactive({
  menuName: '',
  status: ''
})

// 模态框状态
const showAddModal = ref(false)
const showEditModal = ref(false)


// 计算属性 - 展示行（扁平化树，受展开状态控制）
const visibleRows = computed(() => {
  const rows: { menu: Menu; level: number }[] = []
  const walk = (menus: Menu[], level = 0) => {
    menus.forEach(menu => {
      rows.push({ menu, level })
      if (menu.children && menu.children.length > 0 && expandedMenus.value.has(menu.id)) {
        walk(menu.children, level + 1)
      }
    })
  }
  walk(menuTree.value, 0)
  return rows
})

// 表单数据
const menuForm = reactive({
  id: undefined as number | undefined,
  menuName: '',
  parentId: 0,
  menuType: 'M',
  sort: 0,
  component: '',
  perms: '',
  icon: '',
  status: '1',
  remark: ''
})

// 计算属性 - 菜单树
const menuTree = computed(() => {
  return buildMenuTree(menuList.value.filter(menu => {
    if (searchForm.menuName && !menu.menuName.includes(searchForm.menuName)) {
      return false
    }
    if (searchForm.status && menu.status !== searchForm.status) {
      return false
    }
    return true
  }))
})

// 计算属性 - 扁平化菜单列表（用于父级菜单选择）
const flatMenuList = computed(() => {
  const result: (Menu & { level: number })[] = []

  const flatten = (menus: Menu[], level = 0) => {
    menus.forEach(menu => {
      result.push({ ...menu, level })
      if (menu.children && menu.children.length > 0) {
        flatten(menu.children, level + 1)
      }
    })
  }

  flatten(buildMenuTree(menuList.value))
  return result
})

// 构建菜单树
const buildMenuTree = (menus: Menu[]): Menu[] => {
  const menuMap = new Map<number, Menu>()
  const rootMenus: Menu[] = []

  // 创建菜单映射
  menus.forEach(menu => {
    menuMap.set(menu.id, { ...menu, children: [] })
  })

  // 构建树结构
  menus.forEach(menu => {
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

// 获取菜单列表
const getMenuList = async () => {
  try {
    loading.value = true
    const { data: res } = await adminApi.getMenuList({})
    if (res.code === 200) {
      const list = Array.isArray(res.data) ? res.data : (res.data?.list || [])
      // 将后端字段映射成前端表单模型（用于编辑）
      menuList.value = list.map((m: any) => mapBackendToForm(m))
    } else {
      ToastAlert.error({ title: '获取菜单列表失败', message: res.message })
    }
  } catch (error) {
    console.error('获取菜单列表失败:', error)
    ToastAlert.error({
      title: '获取菜单列表失败',
      message: '网络异常，请重试'
    })
  } finally {
    loading.value = false
  }
}

// 搜索菜单
const searchMenus = () => {
  // 搜索时重新计算菜单树
}

// 重置搜索
const resetSearch = () => {
  Object.assign(searchForm, {
    menuName: '',
    status: ''
  })
}

// 展开/收起所有菜单
const expandAll = () => {
  if (allExpanded.value) {
    expandedMenus.value.clear()
  } else {
    const addAllIds = (menus: Menu[]) => {
      menus.forEach(menu => {
        expandedMenus.value.add(menu.id)
        if (menu.children && menu.children.length > 0) {
          addAllIds(menu.children)
        }
      })
    }
    addAllIds(menuTree.value)
  }
  allExpanded.value = !allExpanded.value
}

// 切换菜单展开状态
const toggleMenu = (menuId: number) => {
  if (expandedMenus.value.has(menuId)) {
    expandedMenus.value.delete(menuId)
  } else {
    expandedMenus.value.add(menuId)
  }
}

// 编辑菜单
const editMenu = (menu: Menu) => {
  Object.assign(menuForm, {
    id: menu.id,
    menuName: menu.menuName,
    parentId: menu.parentId,
    menuType: menu.menuType,
    sort: menu.sort,
    component: menu.component || '',
    perms: menu.perms || '',
    icon: menu.icon || '',
    status: menu.status,
    remark: menu.remark || ''
  })
  showEditModal.value = true
}

// 添加子菜单
const addChildMenu = (parentMenu: Menu) => {
  Object.assign(menuForm, {
    id: undefined,
    menuName: '',
    parentId: parentMenu.id,
    menuType: 'C',
    sort: 0,
    component: '',
    perms: '',
    icon: '',
    status: '1',
    remark: ''
  })
  showAddModal.value = true
}

// 删除菜单
const deleteMenu = async (menu: Menu) => {
  if (!confirm(`确定要删除菜单 "${menu.menuName}" 吗？`)) {
    return
  }

  try {
    const { data: res } = await adminApi.deleteMenu(menu.id)

    if (res.code === 200) {
      ToastAlert.success({
        title: '删除成功',
        message: '菜单已删除'
      })
      getMenuList()
    } else {
      ToastAlert.error({
        title: '删除失败',
        message: res.message
      })
    }
  } catch (error) {
    console.error('删除菜单失败:', error)
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
    const apiMethod = isAdd ? adminApi.addMenu : adminApi.updateMenu
    const payload = mapFormToBackend(menuForm)

    const { data: res } = await apiMethod(payload as any)

    if (res.code === 200) {
      ToastAlert.success({
        title: '操作成功',
        message: `菜单已${isAdd ? '添加' : '更新'}`
      })
      closeModal()
      getMenuList()
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

// 关闭模态框
const closeModal = () => {
  showAddModal.value = false
  showEditModal.value = false
  Object.assign(menuForm, {
    id: undefined,
    menuName: '',
    parentId: 0,
    menuType: 'M',
    sort: 0,
    component: '',
    perms: '',
    icon: '',
    status: '1',
    remark: ''
  })
}

// 后端字段双向映射
const mapFormToBackend = (form: any) => {
  const typeMap: Record<string, number> = { M: 1, C: 2, F: 3 }
  return {
    id: form.id,
    parentId: Number(form.parentId || 0),
    menuName: form.menuName,
    menuType: typeMap[form.menuType] || 1,
    url: form.menuType === 'C' ? form.component : (form.menuType === 'M' ? (form.component || '') : ''),
    value: form.menuType === 'F' ? form.perms : '',
    icon: form.icon || '',
    menuStatus: Number(form.status || '1'),
    sort: Number(form.sort || 0),
    remark: form.remark || ''
  }
}

const mapBackendToForm = (item: any) => {
  const typeMap: Record<number, string> = { 1: 'M', 2: 'C', 3: 'F' }
  return {
    id: item.id,
    parentId: item.parentId ?? 0,
    menuName: item.menuName || '',
    menuType: typeMap[item.menuType] || 'M',
    sort: item.sort || 0,
    component: item.url || '',
    perms: item.value || '',
    icon: item.icon || '',
    status: String(item.menuStatus ?? '1'),
    remark: item.remark || ''
  }
}


// 组件挂载时获取数据
onMounted(() => {
  getMenuList()
})
</script>
