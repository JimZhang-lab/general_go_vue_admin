/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/views/Auth/SystemSettings.vue
 * @Description: 
 * 
 */
<template>
  <AuthLayout>
    <PageBreadcrumb :pageTitle="pageTitle" />

    <div class="space-y-5 sm:space-y-6">
      <AuthCard :title="groupMeta.title" :subtitle="groupMeta.description">
        <div class="flex flex-wrap gap-3">
          <router-link
            v-for="item in groups"
            :key="item.key"
            :to="item.path"
            :class="[
              'inline-flex items-center rounded-xl border px-4 py-2 text-sm font-medium transition-colors',
              route.params.group === item.key
                ? 'border-brand-500 bg-brand-50 text-brand-700 dark:border-brand-400 dark:bg-brand-500/15 dark:text-brand-200'
                : 'border-gray-200 text-gray-600 hover:border-brand-300 hover:text-brand-600 dark:border-gray-700 dark:text-gray-300 dark:hover:border-brand-500'
            ]"
          >
            {{ item.label }}
          </router-link>
        </div>
      </AuthCard>

      <AuthCard title="配置项">
        <div v-if="loading" class="grid gap-4 md:grid-cols-2">
          <div
            v-for="index in 4"
            :key="index"
            class="h-36 animate-pulse rounded-2xl border border-gray-200 bg-gray-50 dark:border-gray-800 dark:bg-gray-800/40"
          ></div>
        </div>

        <div v-else-if="settings.length === 0" class="rounded-2xl border border-dashed border-gray-300 px-6 py-10 text-center text-sm text-gray-500 dark:border-gray-700 dark:text-gray-400">
          当前分组暂无可配置项。
        </div>

        <div v-else class="grid gap-4 md:grid-cols-2">
          <div
            v-for="item in settings"
            :key="item.settingKey"
            class="rounded-2xl border border-gray-200 p-5 dark:border-gray-800"
          >
            <div class="flex items-start justify-between gap-4">
              <div>
                <h3 class="text-sm font-semibold text-gray-900 dark:text-white/90">{{ item.settingName }}</h3>
                <p class="mt-1 text-xs text-gray-500 dark:text-gray-400">{{ item.settingKey }}</p>
              </div>
              <span class="rounded-full bg-gray-100 px-2.5 py-1 text-xs font-medium text-gray-600 dark:bg-gray-800 dark:text-gray-300">
                {{ valueTypeLabel(item.valueType) }}
              </span>
            </div>

            <p v-if="item.remark" class="mt-3 text-sm text-gray-600 dark:text-gray-400">
              {{ item.remark }}
            </p>

            <div class="mt-4">
              <label v-if="item.valueType === 'switch'" class="flex items-center justify-between rounded-xl border border-gray-200 px-4 py-3 dark:border-gray-700">
                <span class="text-sm text-gray-700 dark:text-gray-300">启用该配置</span>
                <input
                  v-model="item.boolValue"
                  type="checkbox"
                  class="h-5 w-5 rounded border-gray-300 text-brand-600 focus:ring-brand-500"
                />
              </label>

              <input
                v-else-if="item.valueType === 'number'"
                v-model="item.settingValue"
                type="number"
                class="h-11 w-full rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              />

              <textarea
                v-else-if="item.settingValue.length > 80"
                v-model="item.settingValue"
                rows="4"
                class="w-full rounded-xl border border-gray-300 bg-white px-4 py-3 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              ></textarea>

              <input
                v-else
                v-model="item.settingValue"
                type="text"
                class="h-11 w-full rounded-xl border border-gray-300 bg-white px-4 text-sm text-gray-800 shadow-theme-xs focus:border-brand-300 focus:outline-hidden focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-900 dark:text-white/90"
              />
            </div>
          </div>
        </div>

        <template #footer>
          <div class="flex flex-wrap items-center justify-between gap-3">
            <p class="text-sm text-gray-500 dark:text-gray-400">
              更改将立即保存到数据库，并影响后续系统行为。
            </p>
            <div class="flex gap-3">
              <AuthButton
                text="重新加载"
                variant="secondary"
                :disabled="loading"
                @click="loadSettings"
              />
              <AuthButton
                text="保存设置"
                variant="primary"
                :loading="saving"
                loadingText="保存中..."
                @click="saveSettings"
              />
            </div>
          </div>
        </template>
      </AuthCard>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, ref, watch } from 'vue'
import { useRoute } from 'vue-router'
import { AuthButton, AuthCard, AuthLayout } from '@/components/auth'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'

interface SettingItem {
  id?: number
  groupKey: string
  settingKey: string
  settingName: string
  settingValue: string
  valueType: string
  optionsJson?: string
  isEncrypted?: boolean
  sort?: number
  remark?: string
  boolValue?: boolean
}

const route = useRoute()
const loading = ref(false)
const saving = ref(false)
const settings = ref<SettingItem[]>([])

const groups = [
  { key: 'basic', label: '基础设置', path: '/auth/settings/basic', title: '基础设置', description: '维护站点名称、注册策略等基础配置。' },
  { key: 'security', label: '安全设置', path: '/auth/settings/security', title: '安全设置', description: '控制登录失败锁定策略等运行时安全参数。' },
  { key: 'notification', label: '通知设置', path: '/auth/settings/notification', title: '通知设置', description: '维护站内通知、邮件通知与保留策略。' }
] as const

const currentGroup = computed(() => String(route.params.group || 'basic'))
const groupMeta = computed(() => groups.find((item) => item.key === currentGroup.value) || groups[0])
const pageTitle = computed(() => groupMeta.value.title)

const normalizeSettings = (payload: unknown): SettingItem[] => {
  const list = Array.isArray(payload) ? payload : []
  return list.map((item) => {
    const normalized = {
      ...(item as SettingItem),
      settingValue: String((item as SettingItem).settingValue ?? ''),
      valueType: (item as SettingItem).valueType || 'text'
    }
    normalized.boolValue = normalized.settingValue === 'true'
    return normalized
  })
}

const valueTypeLabel = (valueType?: string) => {
  switch (valueType) {
    case 'switch':
      return '开关'
    case 'number':
      return '数字'
    default:
      return '文本'
  }
}

const loadSettings = async () => {
  try {
    loading.value = true
    const { data: res } = await adminApi.getSettingList(currentGroup.value)
    if (res.code !== 200) {
      throw new Error(res.message || '获取系统设置失败')
    }
    settings.value = normalizeSettings(res.data)
  } catch (error) {
    settings.value = []
    ToastAlert.error({
      title: '加载失败',
      message: error instanceof Error ? error.message : '系统设置加载失败'
    })
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  if (settings.value.length === 0) {
    return
  }

  try {
    saving.value = true
    const payload = settings.value.map((item) => ({
      ...item,
      settingValue: item.valueType === 'switch' ? String(Boolean(item.boolValue)) : String(item.settingValue ?? '')
    }))

    const { data: res } = await adminApi.batchUpdateSettings(payload)
    if (res.code !== 200) {
      throw new Error(res.message || '系统设置保存失败')
    }

    ToastAlert.success({
      title: '保存成功',
      message: `${groupMeta.value.title}已更新`
    })
    await loadSettings()
  } catch (error) {
    ToastAlert.error({
      title: '保存失败',
      message: error instanceof Error ? error.message : '系统设置保存失败'
    })
  } finally {
    saving.value = false
  }
}

watch(
  () => route.params.group,
  () => {
    loadSettings()
  },
  { immediate: true }
)
</script>
