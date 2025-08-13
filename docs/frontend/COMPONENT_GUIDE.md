# 🎨 Vue3组件开发指南

Vue3组件开发规范和最佳实践，包含组合式API、TypeScript、组件设计等完整指南。

## 📋 组件开发规范

### 1. 组件命名规范
```typescript
// ✅ 正确：使用PascalCase
export default defineComponent({
  name: 'UserProfile'
})

// ✅ 文件命名：使用kebab-case
// user-profile.vue
// user-profile-card.vue

// ❌ 错误：不要使用camelCase
// userProfile.vue
```

### 2. 组件结构规范
```vue
<template>
  <!-- 模板内容 -->
</template>

<script setup lang="ts">
// 1. 导入依赖
import { ref, computed, onMounted } from 'vue'
import type { User } from '@/types/user'

// 2. 定义Props
interface Props {
  user: User
  readonly?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  readonly: false
})

// 3. 定义Emits
interface Emits {
  update: [user: User]
  delete: [id: number]
}

const emit = defineEmits<Emits>()

// 4. 响应式数据
const isLoading = ref(false)
const formData = ref<User>({ ...props.user })

// 5. 计算属性
const isValid = computed(() => {
  return formData.value.name && formData.value.email
})

// 6. 方法
const handleSubmit = () => {
  if (isValid.value) {
    emit('update', formData.value)
  }
}

// 7. 生命周期
onMounted(() => {
  console.log('Component mounted')
})

// 8. 暴露给父组件的方法
defineExpose({
  reset: () => {
    formData.value = { ...props.user }
  }
})
</script>

<style scoped>
/* 组件样式 */
</style>
```

## 🔧 组合式API最佳实践

### 1. 使用Composables
```typescript
// composables/useUser.ts
import { ref, computed } from 'vue'
import type { User } from '@/types/user'
import { userApi } from '@/api/user'

export function useUser() {
  const users = ref<User[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  const activeUsers = computed(() => 
    users.value.filter(user => user.status === 1)
  )

  const fetchUsers = async () => {
    try {
      loading.value = true
      error.value = null
      const response = await userApi.getList()
      users.value = response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : '获取用户失败'
    } finally {
      loading.value = false
    }
  }

  const createUser = async (userData: Omit<User, 'id'>) => {
    try {
      loading.value = true
      const response = await userApi.create(userData)
      users.value.push(response.data)
      return response.data
    } catch (err) {
      error.value = err instanceof Error ? err.message : '创建用户失败'
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    users: readonly(users),
    loading: readonly(loading),
    error: readonly(error),
    activeUsers,
    fetchUsers,
    createUser
  }
}
```

### 2. 在组件中使用
```vue
<script setup lang="ts">
import { useUser } from '@/composables/useUser'

const { users, loading, error, fetchUsers, createUser } = useUser()

onMounted(() => {
  fetchUsers()
})
</script>
```

## 📝 表单组件开发

### 1. 基础表单组件
```vue
<!-- components/form/BaseForm.vue -->
<template>
  <el-form
    ref="formRef"
    :model="modelValue"
    :rules="rules"
    :label-width="labelWidth"
    @submit.prevent="handleSubmit"
  >
    <slot />
    
    <el-form-item v-if="showActions">
      <el-button @click="handleReset">重置</el-button>
      <el-button 
        type="primary" 
        :loading="loading"
        @click="handleSubmit"
      >
        {{ submitText }}
      </el-button>
    </el-form-item>
  </el-form>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import type { FormInstance, FormRules } from 'element-plus'

interface Props {
  modelValue: Record<string, any>
  rules?: FormRules
  labelWidth?: string
  loading?: boolean
  submitText?: string
  showActions?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  labelWidth: '120px',
  loading: false,
  submitText: '提交',
  showActions: true
})

interface Emits {
  'update:modelValue': [value: Record<string, any>]
  submit: [value: Record<string, any>]
  reset: []
}

const emit = defineEmits<Emits>()

const formRef = ref<FormInstance>()

const handleSubmit = async () => {
  if (!formRef.value) return
  
  try {
    await formRef.value.validate()
    emit('submit', props.modelValue)
  } catch (error) {
    console.error('表单验证失败:', error)
  }
}

const handleReset = () => {
  formRef.value?.resetFields()
  emit('reset')
}

defineExpose({
  validate: () => formRef.value?.validate(),
  resetFields: () => formRef.value?.resetFields(),
  clearValidate: () => formRef.value?.clearValidate()
})
</script>
```

### 2. 使用表单组件
```vue
<template>
  <BaseForm
    v-model="formData"
    :rules="formRules"
    :loading="loading"
    @submit="handleSubmit"
    @reset="handleReset"
  >
    <el-form-item label="用户名" prop="username">
      <el-input v-model="formData.username" />
    </el-form-item>
    
    <el-form-item label="邮箱" prop="email">
      <el-input v-model="formData.email" type="email" />
    </el-form-item>
  </BaseForm>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import BaseForm from '@/components/form/BaseForm.vue'
import type { FormRules } from 'element-plus'

const formData = ref({
  username: '',
  email: ''
})

const formRules: FormRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' }
  ],
  email: [
    { required: true, message: '请输入邮箱', trigger: 'blur' },
    { type: 'email', message: '请输入正确的邮箱格式', trigger: 'blur' }
  ]
}

const loading = ref(false)

const handleSubmit = async (data: typeof formData.value) => {
  loading.value = true
  try {
    // 提交逻辑
    console.log('提交数据:', data)
  } finally {
    loading.value = false
  }
}

const handleReset = () => {
  formData.value = {
    username: '',
    email: ''
  }
}
</script>
```

## 📊 表格组件开发

### 1. 基础表格组件
```vue
<!-- components/table/BaseTable.vue -->
<template>
  <div class="base-table">
    <el-table
      :data="data"
      :loading="loading"
      v-bind="$attrs"
      @selection-change="handleSelectionChange"
    >
      <el-table-column
        v-if="showSelection"
        type="selection"
        width="55"
      />
      
      <el-table-column
        v-for="column in columns"
        :key="column.prop"
        v-bind="column"
      >
        <template #default="scope" v-if="column.slot">
          <slot 
            :name="column.slot" 
            :row="scope.row" 
            :column="column"
            :index="scope.$index"
          />
        </template>
      </el-table-column>
      
      <el-table-column
        v-if="showActions"
        label="操作"
        :width="actionWidth"
        fixed="right"
      >
        <template #default="scope">
          <slot 
            name="actions" 
            :row="scope.row" 
            :index="scope.$index"
          />
        </template>
      </el-table-column>
    </el-table>
    
    <el-pagination
      v-if="showPagination"
      v-model:current-page="currentPage"
      v-model:page-size="pageSize"
      :total="total"
      :page-sizes="pageSizes"
      layout="total, sizes, prev, pager, next, jumper"
      @size-change="handleSizeChange"
      @current-change="handleCurrentChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'

interface Column {
  prop: string
  label: string
  width?: string | number
  minWidth?: string | number
  slot?: string
  sortable?: boolean
  [key: string]: any
}

interface Props {
  data: any[]
  columns: Column[]
  loading?: boolean
  showSelection?: boolean
  showActions?: boolean
  actionWidth?: string | number
  showPagination?: boolean
  total?: number
  pageSize?: number
  currentPage?: number
  pageSizes?: number[]
}

const props = withDefaults(defineProps<Props>(), {
  loading: false,
  showSelection: false,
  showActions: true,
  actionWidth: 200,
  showPagination: true,
  total: 0,
  pageSize: 10,
  currentPage: 1,
  pageSizes: () => [10, 20, 50, 100]
})

interface Emits {
  'update:currentPage': [page: number]
  'update:pageSize': [size: number]
  'selection-change': [selection: any[]]
  'page-change': [page: number, size: number]
}

const emit = defineEmits<Emits>()

const currentPage = computed({
  get: () => props.currentPage,
  set: (value) => emit('update:currentPage', value)
})

const pageSize = computed({
  get: () => props.pageSize,
  set: (value) => emit('update:pageSize', value)
})

const handleSelectionChange = (selection: any[]) => {
  emit('selection-change', selection)
}

const handleSizeChange = (size: number) => {
  emit('page-change', currentPage.value, size)
}

const handleCurrentChange = (page: number) => {
  emit('page-change', page, pageSize.value)
}
</script>
```

## 🎯 业务组件示例

### 1. 用户管理组件
```vue
<!-- views/admin/UserManagement.vue -->
<template>
  <div class="user-management">
    <div class="toolbar">
      <el-button type="primary" @click="showCreateDialog = true">
        新增用户
      </el-button>
      <el-button 
        :disabled="!selectedUsers.length"
        @click="handleBatchDelete"
      >
        批量删除
      </el-button>
    </div>
    
    <BaseTable
      :data="users"
      :columns="tableColumns"
      :loading="loading"
      :total="total"
      :current-page="pagination.page"
      :page-size="pagination.pageSize"
      show-selection
      @selection-change="selectedUsers = $event"
      @page-change="handlePageChange"
    >
      <template #status="{ row }">
        <el-tag :type="row.status === 1 ? 'success' : 'danger'">
          {{ row.status === 1 ? '启用' : '禁用' }}
        </el-tag>
      </template>
      
      <template #actions="{ row }">
        <el-button size="small" @click="handleEdit(row)">
          编辑
        </el-button>
        <el-button 
          size="small" 
          type="danger" 
          @click="handleDelete(row)"
        >
          删除
        </el-button>
      </template>
    </BaseTable>
    
    <!-- 创建/编辑对话框 -->
    <UserDialog
      v-model="showCreateDialog"
      :user="currentUser"
      @success="handleDialogSuccess"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import BaseTable from '@/components/table/BaseTable.vue'
import UserDialog from './components/UserDialog.vue'
import { useUser } from '@/composables/useUser'
import type { User } from '@/types/user'

const { users, loading, total, fetchUsers, deleteUser } = useUser()

const selectedUsers = ref<User[]>([])
const showCreateDialog = ref(false)
const currentUser = ref<User | null>(null)

const pagination = ref({
  page: 1,
  pageSize: 10
})

const tableColumns = [
  { prop: 'id', label: 'ID', width: 80 },
  { prop: 'username', label: '用户名', minWidth: 120 },
  { prop: 'nickname', label: '昵称', minWidth: 120 },
  { prop: 'email', label: '邮箱', minWidth: 180 },
  { prop: 'phone', label: '手机号', minWidth: 120 },
  { prop: 'status', label: '状态', width: 100, slot: 'status' },
  { prop: 'createTime', label: '创建时间', minWidth: 160 }
]

const handlePageChange = (page: number, pageSize: number) => {
  pagination.value = { page, pageSize }
  fetchUsers(pagination.value)
}

const handleEdit = (user: User) => {
  currentUser.value = user
  showCreateDialog.value = true
}

const handleDelete = async (user: User) => {
  try {
    await ElMessageBox.confirm(
      `确定要删除用户 "${user.username}" 吗？`,
      '确认删除',
      { type: 'warning' }
    )
    
    await deleteUser(user.id)
    ElMessage.success('删除成功')
    fetchUsers(pagination.value)
  } catch (error) {
    // 用户取消删除
  }
}

const handleBatchDelete = async () => {
  // 批量删除逻辑
}

const handleDialogSuccess = () => {
  showCreateDialog.value = false
  currentUser.value = null
  fetchUsers(pagination.value)
}

onMounted(() => {
  fetchUsers(pagination.value)
})
</script>
```

## 🎨 样式规范

### 1. CSS变量使用
```vue
<style scoped>
.user-management {
  --primary-color: #409eff;
  --success-color: #67c23a;
  --warning-color: #e6a23c;
  --danger-color: #f56c6c;
  
  padding: 20px;
}

.toolbar {
  margin-bottom: 20px;
  display: flex;
  gap: 12px;
}
</style>
```

### 2. 响应式设计
```vue
<style scoped>
.user-management {
  padding: 20px;
}

@media (max-width: 768px) {
  .user-management {
    padding: 10px;
  }
  
  .toolbar {
    flex-direction: column;
  }
}
</style>
```

## 🧪 组件测试

### 1. 单元测试
```typescript
// tests/components/BaseTable.test.ts
import { mount } from '@vue/test-utils'
import BaseTable from '@/components/table/BaseTable.vue'

describe('BaseTable', () => {
  const mockData = [
    { id: 1, name: 'John', email: 'john@example.com' },
    { id: 2, name: 'Jane', email: 'jane@example.com' }
  ]

  const mockColumns = [
    { prop: 'id', label: 'ID' },
    { prop: 'name', label: '姓名' },
    { prop: 'email', label: '邮箱' }
  ]

  it('renders table with data', () => {
    const wrapper = mount(BaseTable, {
      props: {
        data: mockData,
        columns: mockColumns
      }
    })

    expect(wrapper.find('.el-table').exists()).toBe(true)
    expect(wrapper.text()).toContain('John')
    expect(wrapper.text()).toContain('Jane')
  })

  it('emits page-change event', async () => {
    const wrapper = mount(BaseTable, {
      props: {
        data: mockData,
        columns: mockColumns,
        showPagination: true,
        total: 100
      }
    })

    await wrapper.find('.el-pagination .btn-next').trigger('click')
    expect(wrapper.emitted('page-change')).toBeTruthy()
  })
})
```

## 📚 组件文档

### 1. 组件注释
```vue
<script setup lang="ts">
/**
 * 基础表格组件
 * 
 * @description 封装了Element Plus的表格组件，提供了分页、选择、操作列等常用功能
 * @author 开发团队
 * @version 1.0.0
 * 
 * @example
 * <BaseTable
 *   :data="tableData"
 *   :columns="tableColumns"
 *   :loading="loading"
 *   show-selection
 *   @selection-change="handleSelectionChange"
 * />
 */

interface Props {
  /** 表格数据 */
  data: any[]
  /** 表格列配置 */
  columns: Column[]
  /** 加载状态 */
  loading?: boolean
}
</script>
```

### 2. README文档
```markdown
# BaseTable 组件

基础表格组件，封装了Element Plus的表格功能。

## Props

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| data | Array | [] | 表格数据 |
| columns | Array | [] | 表格列配置 |
| loading | Boolean | false | 加载状态 |

## Events

| 事件名 | 参数 | 说明 |
|--------|------|------|
| selection-change | selection | 选择项改变 |
| page-change | page, size | 分页改变 |

## Slots

| 插槽名 | 参数 | 说明 |
|--------|------|------|
| actions | row, index | 操作列内容 |
| [column.slot] | row, column, index | 自定义列内容 |
```

---

**最后更新**: 2025-07-29  
**维护者**: 前端开发团队  
**技术支持**: 查看Vue3官方文档或联系开发团队
