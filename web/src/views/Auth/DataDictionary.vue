<template>
  <AuthLayout>
    <PageBreadcrumb pageTitle="字典管理" />
    <div class="space-y-4 sm:space-y-5 lg:flex lg:space-y-0 lg:space-x-5 h-[calc(100vh-140px)]">

      <!-- 左侧：字典类型 -->
      <AuthCard title="字典分类" class="lg:w-1/3 flex flex-col h-full overflow-hidden">
        <div class="mb-4 flex space-x-2">
          <AuthInput
            v-model="typeSearchForm.dictName"
            type="text"
            placeholder="搜索字典分类..."
            icon="search"
            class="flex-1"
            @keyup.enter="searchDictTypes"
          />
          <button @click="showAddTypeModal = true" class="flex-none bg-blue-600 text-white px-3 py-2 rounded-xl border border-blue-600 hover:bg-blue-700 transition">
            <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/></svg>
          </button>
        </div>

        <div class="flex-1 overflow-y-auto custom-scrollbar -mx-2 px-2 relative">
          <div v-if="loadingTypes" class="flex justify-center py-8 text-gray-500">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-blue-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg> 拼命加载中...
          </div>
          <div v-else-if="dictTypeList.length === 0" class="text-center py-8 text-gray-500 text-sm">
            暂无字典分类
          </div>
          <ul v-else class="space-y-2">
            <li
              v-for="type in dictTypeList"
              :key="type.id"
              @click="selectDictType(type)"
              :class="[ 'p-3 rounded-xl border cursor-pointer transition-all duration-200 group relative', selectedDictType?.id === type.id ? 'border-blue-500 bg-blue-50 dark:bg-blue-500/10' : 'border-gray-100 hover:border-blue-200 hover:bg-gray-50 dark:border-gray-800 dark:hover:bg-gray-800/50' ]"
            >
              <div class="font-medium text-gray-800 dark:text-gray-200 text-sm flex justify-between items-center">
                <span>{{ type.dictName }}</span>
                <span :class="type.status === '1' ? 'bg-success-100 text-success-700' : 'bg-error-100 text-error-700'" class="text-[10px] px-1.5 rounded-full">{{ type.status === '1' ? '正常' : '停用' }}</span>
              </div>
              <div class="text-xs text-gray-500 dark:text-gray-400 mt-1 font-mono">{{ type.dictType }}</div>

              <!-- Hover Actions -->
              <div class="absolute right-2 top-1/2 -translate-y-1/2 flex space-x-1 opacity-0 group-hover:opacity-100 transition-opacity bg-white/80 dark:bg-gray-900/80 backdrop-blur-sm p-1 rounded-lg shadow-sm" v-if="selectedDictType?.id !== type.id">
                <button @click.stop="editDictType(type)" class="text-blue-600 hover:bg-blue-50 p-1 rounded"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg></button>
                <button @click.stop="deleteDictType(type)" class="text-red-500 hover:bg-red-50 p-1 rounded"><svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg></button>
              </div>
            </li>
          </ul>
        </div>
      </AuthCard>

      <!-- 右侧：字典数据 -->
      <AuthCard :title="selectedDictType ? `[${selectedDictType.dictName}] 数据键值` : '请选择字典字典分类'" class="lg:w-2/3 flex flex-col h-full overflow-hidden">
        <template v-if="selectedDictType">
          <div class="mb-4 flex flex-col sm:flex-row justify-between gap-3">
            <div class="flex-1 max-w-sm relative">
              <AuthInput
                v-model="dataSearchForm.dictLabel"
                type="text"
                placeholder="搜索标签名称..."
                icon="search"
                @keyup.enter="searchDictData"
              />
            </div>
            <div class="flex gap-2">
              <AuthButton @click="showAddDataModal = true" variant="primary" size="md" text="新增字典配置" />
            </div>
          </div>

          <div class="flex-1 overflow-hidden rounded-xl border border-gray-200 dark:border-gray-800">
            <div class="h-full overflow-auto custom-scrollbar">
              <table class="min-w-full relative">
                <thead class="sticky top-0 bg-gray-50/90 dark:bg-gray-800/90 backdrop-blur z-10 shadow-sm">
                  <tr class="border-b border-gray-200 dark:border-gray-700">
                    <th class="px-5 py-3 text-left font-medium text-gray-500 text-theme-xs dark:text-gray-400">标签展现 (Label)</th>
                    <th class="px-5 py-3 text-left font-medium text-gray-500 text-theme-xs dark:text-gray-400">底层键值 (Value)</th>
                    <th class="px-5 py-3 text-left font-medium text-gray-500 text-theme-xs dark:text-gray-400">排序</th>
                    <th class="px-5 py-3 text-left font-medium text-gray-500 text-theme-xs dark:text-gray-400">状态</th>
                    <th class="px-5 py-3 text-left font-medium text-gray-500 text-theme-xs dark:text-gray-400">操作</th>
                  </tr>
                </thead>
                <tbody class="divide-y divide-gray-100 dark:divide-gray-800 bg-white dark:bg-transparent">
                  <tr v-if="loadingData"><td colspan="5" class="py-8 text-center text-gray-400">正在获取数据...</td></tr>
                  <tr v-else-if="dictDataList.length === 0"><td colspan="5" class="py-8 text-center text-gray-400">该分类下暂无数据项配置</td></tr>
                  <tr v-else v-for="item in dictDataList" :key="item.id" class="hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors">
                    <td class="px-5 py-4 text-theme-sm text-gray-800 dark:text-gray-200 font-medium">{{ item.dictLabel }}</td>
                    <td class="px-5 py-4 text-theme-sm font-mono text-blue-600 dark:text-blue-400">
                      <span class="bg-blue-50 dark:bg-blue-900/30 px-2 py-0.5 rounded">{{ item.dictValue }}</span>
                    </td>
                    <td class="px-5 py-4 text-theme-sm text-gray-500 dark:text-gray-400">{{ item.dictSort }}</td>
                    <td class="px-5 py-4 text-theme-sm">
                      <span :class="item.status === '1' ? 'bg-success-50 text-success-700 dark:bg-success-500/15 dark:text-success-500' : 'bg-error-50 text-error-700 dark:bg-error-500/15 dark:text-error-500'" class="rounded-full px-2 py-0.5 text-theme-xs font-medium">{{ item.status === '1' ? '启用' : '禁用' }}</span>
                    </td>
                    <td class="px-5 py-4 text-theme-sm">
                      <div class="flex gap-2">
                        <button @click="editDictData(item)" class="inline-flex items-center px-2 py-1 text-xs font-medium text-blue-600 bg-blue-50 rounded hover:bg-blue-100 dark:bg-blue-500/15 dark:text-blue-400 dark:hover:bg-blue-500/25">编辑</button>
                        <button @click="deleteDictData(item)" class="inline-flex items-center px-2 py-1 text-xs font-medium text-red-600 bg-red-50 rounded hover:bg-red-100 dark:bg-red-500/15 dark:text-red-400 dark:hover:bg-red-500/25">删除</button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </template>

        <!-- 空状态 SVG -->
        <div v-else class="flex-1 flex flex-col items-center justify-center text-gray-400 opacity-60">
          <svg class="w-32 h-32 mb-4" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1" stroke-linecap="round" stroke-linejoin="round">
            <path d="M4 22h14a2 2 0 0 0 2-2V7.5L14.5 2H6a2 2 0 0 0-2 2v4" />
            <path d="M14 2v6h6" />
            <path d="M2.93 14c.32-.93 1.25-1.5 2.27-1.5h2.6M2.93 14C2.31 16 3.1 18 5.2 18h2.6M2.93 14l3.1-9h9M5.2 18L8.3 9h9" />
          </svg>
          <p class="text-lg">点击左侧分类来配置具体字典</p>
        </div>
      </AuthCard>

    </div>

    <!-- Dict Type 模态框 -->
    <div v-if="showAddTypeModal || showEditTypeModal" class="fixed inset-0 overflow-y-auto h-full w-full z-99999">
      <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" @click="closeTypeModal"></div>
      <div class="relative top-20 mx-auto p-6 border w-11/12 md:w-3/4 lg:w-1/3 shadow-2xl rounded-2xl bg-white dark:bg-gray-800 dark:border-gray-700">
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6">
          {{ showAddTypeModal ? '新增体系分类' : '修改字典分类' }}
        </h3>
        <form @submit.prevent="submitTypeForm" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">字典名称 *</label>
            <input v-model="typeForm.dictName" required class="w-full px-4 py-2 border rounded-xl" placeholder="例：用户性别" />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">字典标识类型 *</label>
            <input v-model="typeForm.dictType" required class="w-full px-4 py-2 border rounded-xl" placeholder="例：sys_user_sex" />
            <p class="text-xs text-gray-400 mt-1">全局唯一的英文键名，如 sys_normal_disable</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">状态</label>
            <select v-model="typeForm.status" class="w-full px-4 py-2 border rounded-xl">
              <option value="1">正常启用</option>
              <option value="2">停用冻结</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">备注说明</label>
            <textarea v-model="typeForm.remark" rows="3" class="w-full px-4 py-2 border rounded-xl"></textarea>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="closeTypeModal" class="px-5 py-2.5 rounded-xl text-gray-600 hover:bg-gray-100 transition">取消</button>
            <button type="submit" class="px-5 py-2.5 bg-blue-600 text-white rounded-xl shadow-lg shadow-blue-600/30 hover:bg-blue-700 transition">保存分类</button>
          </div>
        </form>
      </div>
    </div>

    <!-- Dict Data 模态框 -->
    <div v-if="showAddDataModal || showEditDataModal" class="fixed inset-0 overflow-y-auto h-full w-full z-99999">
      <div class="fixed inset-0 bg-black/50 backdrop-blur-sm" @click="closeDataModal"></div>
      <div class="relative top-20 mx-auto p-6 border w-11/12 md:w-3/4 lg:w-1/3 shadow-2xl rounded-2xl bg-white dark:bg-gray-800 dark:border-gray-700">
        <h3 class="text-xl font-bold text-gray-900 dark:text-white mb-6">
          {{ showAddDataModal ? '新增数据键值项' : '修改数据项' }} <span class="text-xs font-mono bg-gray-100 p-1 rounded text-gray-500">{{ selectedDictType?.dictType }}</span>
        </h3>
        <form @submit.prevent="submitDataForm" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">字典标签 (Label) *</label>
            <input v-model="dataForm.dictLabel" required class="w-full px-4 py-2 border rounded-xl" placeholder="例：男、女、未知" />
            <p class="text-xs text-gray-400 mt-1">前端展示的中文释义</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">字典键值 (Value) *</label>
            <input v-model="dataForm.dictValue" required class="w-full px-4 py-2 border rounded-xl" placeholder="例：1、2、0" />
            <p class="text-xs text-gray-400 mt-1">数据库中实际交互存储的底层值</p>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">展示排序</label>
              <input v-model.number="dataForm.dictSort" type="number" class="w-full px-4 py-2 border rounded-xl" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">状态</label>
              <select v-model="dataForm.status" class="w-full px-4 py-2 border rounded-xl">
                <option value="1">正常使用</option>
                <option value="2">禁止使用</option>
              </select>
            </div>
          </div>
          <div class="flex justify-end space-x-3 pt-4">
            <button type="button" @click="closeDataModal" class="px-5 py-2.5 rounded-xl text-gray-600 hover:bg-gray-100 transition">取消</button>
            <button type="submit" class="px-5 py-2.5 bg-blue-600 text-white rounded-xl shadow-lg shadow-blue-600/30 hover:bg-blue-700 transition">落库保存</button>
          </div>
        </form>
      </div>
    </div>

  </AuthLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { AuthLayout, AuthCard, AuthButton, AuthInput } from '@/components/auth'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import { dictApi } from '@/api/system/dict'
import ToastAlert from '@/composables/ToastAlert'

// =============== 字典分类 (左侧) ===============
const dictTypeList = ref<any[]>([])
const loadingTypes = ref(false)
const selectedDictType = ref<any>(null)

const typeSearchForm = reactive({ dictName: '' })
const showAddTypeModal = ref(false)
const showEditTypeModal = ref(false)

const typeForm = reactive({
  id: undefined,
  dictName: '',
  dictType: '',
  status: '1',
  remark: ''
})

const getDictTypeList = async () => {
  loadingTypes.value = true
  try {
    const res: any = await dictApi.getDictTypeList({ pageNum: 1, pageSize: 999, dictName: typeSearchForm.dictName })
    if (res.data.code === 200) {
      dictTypeList.value = res.data.data.list || []
      // Auto select first
      if (dictTypeList.value.length > 0 && !selectedDictType.value) {
        selectDictType(dictTypeList.value[0])
      }
    } else {
      ToastAlert.error({ title: '加载失败', message: res.data.message })
    }
  } finally {
    loadingTypes.value = false
  }
}

const searchDictTypes = () => getDictTypeList()
const selectDictType = (type: any) => {
  selectedDictType.value = type
  getDictDataList()
}

const submitTypeForm = async () => {
  try {
    const apiMethod = showAddTypeModal.value ? dictApi.addDictType : dictApi.updateDictType
    const res: any = await apiMethod(typeForm)
    if (res.data.code === 200) {
      ToastAlert.success({ title: '操作成功', message: '字典分类已更新' })
      closeTypeModal()
      getDictTypeList()
    } else {
      ToastAlert.error({ title: '操作失败', message: res.data.message })
    }
  } catch(e) {
    ToastAlert.error({ title: '网络异常' })
  }
}

const editDictType = (type: any) => {
  Object.assign(typeForm, { ...type })
  showEditTypeModal.value = true
}

const deleteDictType = async (type: any) => {
  if (!confirm(`将删除 "${type.dictName}" 的所有关联数据，千万小心！确认？`)) return
  const res: any = await dictApi.deleteDictType({ ids: [type.id] })
  if (res.data.code === 200) {
    ToastAlert.success({ title: '删除成功' })
    if (selectedDictType.value?.id === type.id) selectedDictType.value = null
    getDictTypeList()
  } else {
    ToastAlert.error({ title: '删除失败', message: res.data.message })
  }
}

const closeTypeModal = () => {
  showAddTypeModal.value = false
  showEditTypeModal.value = false
  Object.assign(typeForm, { id: undefined, dictName: '', dictType: '', status: '1', remark: '' })
}

// =============== 字典数据 (右侧) ===============
const dictDataList = ref<any[]>([])
const loadingData = ref(false)
const dataSearchForm = reactive({ dictLabel: '' })
const showAddDataModal = ref(false)
const showEditDataModal = ref(false)

const dataForm = reactive({
  id: undefined,
  dictLabel: '',
  dictValue: '',
  dictType: '',
  dictSort: 0,
  status: '1'
})

const getDictDataList = async () => {
  if (!selectedDictType.value) return
  loadingData.value = true
  try {
    const res: any = await dictApi.getDictDataList({
      pageNum: 1, pageSize: 999,
      dictType: selectedDictType.value.dictType,
      dictLabel: dataSearchForm.dictLabel
    })
    if (res.data.code === 200) {
      dictDataList.value = res.data.data.list || []
    }
  } finally {
    loadingData.value = false
  }
}

const searchDictData = () => getDictDataList()

const submitDataForm = async () => {
  try {
    if (!selectedDictType.value) return
    dataForm.dictType = selectedDictType.value.dictType

    const apiMethod = showAddDataModal.value ? dictApi.addDictData : dictApi.updateDictData
    const res: any = await apiMethod(dataForm)
    if (res.data.code === 200) {
      ToastAlert.success({ title: '操作成功', message: '键值已更新' })
      closeDataModal()
      getDictDataList()
    } else {
      ToastAlert.error({ title: '异常拦截', message: res.data.message })
    }
  } catch(e) {
    ToastAlert.error({ title: '网络异常' })
  }
}

const editDictData = (item: any) => {
  Object.assign(dataForm, { ...item })
  showEditDataModal.value = true
}

const deleteDictData = async (item: any) => {
  if (!confirm(`确认删除数据项 "${item.dictLabel}" 吗？`)) return
  const res: any = await dictApi.deleteDictData({ ids: [item.id] })
  if (res.data.code === 200) {
    ToastAlert.success({ title: '已抹除' })
    getDictDataList()
  } else {
    ToastAlert.error({ title: '错误', message: res.data.message })
  }
}

const closeDataModal = () => {
  showAddDataModal.value = false
  showEditDataModal.value = false
  Object.assign(dataForm, { id: undefined, dictLabel: '', dictValue: '', dictType: '', dictSort: 0, status: '1' })
}

onMounted(() => {
  getDictTypeList()
})
</script>
