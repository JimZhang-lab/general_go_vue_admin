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
                : 'border-gray-200 bg-white text-gray-600 hover:border-brand-300 hover:text-brand-600 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-300 dark:hover:border-brand-500'
            ]"
          >
            {{ item.label }}
          </router-link>
        </div>
      </AuthCard>

      <div class="rounded-2xl border border-gray-200 bg-white dark:border-gray-800 dark:bg-gray-900 shadow-sm relative overflow-hidden">

        <div v-if="loading" class="absolute inset-0 z-10 bg-white/50 dark:bg-gray-900/50 flex items-center justify-center backdrop-blur-sm">
          <div class="flex flex-col items-center">
            <svg class="animate-spin h-8 w-8 text-brand-600 mb-2" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"><circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle><path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path></svg>
            <span class="text-sm font-medium text-gray-500">加载配置中...</span>
          </div>
        </div>

        <div class="p-5 md:p-8">
           <component :is="activePanelComponent" />
        </div>

        <div class="bg-gray-50 border-t border-gray-200 p-4 sm:px-8 dark:bg-gray-800/40 dark:border-gray-800 flex flex-wrap items-center justify-between gap-3">
          <p class="text-sm text-gray-500 dark:text-gray-400">
            修改配置将立即生效，如有必要，安全设置更替将重置当前活动的所有会话。
          </p>
          <div class="flex gap-3 w-full sm:w-auto">
            <AuthButton
              text="还原默认"
              variant="secondary"
              :disabled="loading || saving"
              class="flex-1 sm:flex-none"
              @click="loadSettings"
            />
            <AuthButton
              text="保存设置"
              variant="primary"
              :loading="saving"
              loadingText="执行中..."
              class="flex-1 sm:flex-none"
              @click="saveSettings"
            />
          </div>
        </div>
      </div>
    </div>
  </AuthLayout>
</template>

<script setup lang="ts">
import { computed, ref, reactive, watch, markRaw, provide } from 'vue'
import { useRoute } from 'vue-router'
import { AuthButton, AuthCard, AuthLayout } from '@/components/auth'
import PageBreadcrumb from '@/components/common/PageBreadcrumb.vue'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'

// Import bespoke config panels
import BasicSettingsPanel from './Settings/BasicSettingsPanel.vue'
import SecuritySettingsPanel from './Settings/SecuritySettingsPanel.vue'
import NotificationSettingsPanel from './Settings/NotificationSettingsPanel.vue'

const route = useRoute()
const loading = ref(false)
const saving = ref(false)

// We hold a local flat key-value state for the bespoke panels to easily bind
const formData = reactive<Record<string, any>>({})

// Provide the reactive store to child panels uniformly
provide('settingsFormData', formData)

const groups = [
  { key: 'basic', label: '通用基础设置', path: '/auth/settings/basic', title: '基础应用配置', description: '维护站点的跨域通信、法律备案、基础注册等最基础的信息网。', component: BasicSettingsPanel },
  { key: 'security', label: '引擎与安全', path: '/auth/settings/security', title: '边界与安全参数', description: '控制高危用户的重试封锁、JWT有效期、全局强制双因素认证。' },
  { key: 'notification', label: '节点与通知', path: '/auth/settings/notification', title: '通知架构与网关', description: '配置站点外发的物理邮箱、短信通知及内部事件负载留存期。' }
] as const

const currentGroup = computed(() => String(route.params.group || 'basic'))
const groupMeta = computed(() => groups.find((item) => item.key === currentGroup.value) || groups[0])
const pageTitle = computed(() => groupMeta.value.title)

const activePanelComponent = computed(() => {
  switch (currentGroup.value) {
    case 'security': return markRaw(SecuritySettingsPanel)
    case 'notification': return markRaw(NotificationSettingsPanel)
    case 'basic':
    default:
      return markRaw(BasicSettingsPanel)
  }
})

// Definition lists to map bespoke fields to their Database representations
const fieldDefaults: Record<string, Record<string, {name: string, type: string, default: any}>> = {
  'basic': {
    'basic.site_name': { name: '系统名称', type: 'text', default: '企业级后台' },
    'basic.site_slogan': { name: '系统副标题', type: 'text', default: '中台基座' },
    'basic.site_domain': { name: '主域名', type: 'text', default: 'admin.example.com' },
    'basic.company_name': { name: '企业名称', type: 'text', default: '' },
    'basic.icp_beian': { name: 'ICP备案', type: 'text', default: '' },
    'basic.gongan_beian': { name: '公安联网备案', type: 'text', default: '' },
    'basic.allow_register': { name: '开放注册', type: 'switch', default: false },
    'basic.force_watermark': { name: '全站水印', type: 'switch', default: false },
  },
  'security': {
    'security.login_failed_limit': { name: '最大登录失败', type: 'number', default: 5 },
    'security.lock_minutes': { name: '冻结时长', type: 'number', default: 30 },
    'security.session_timeout_hours': { name: '凭据有效期', type: 'number', default: 24 },
    'security.force_strong_pwd': { name: '强制强密码', type: 'switch', default: true },
  },
  'notification': {
    'notification.enable_site_notice': { name: '启用全局系统通知', type: 'switch', default: true },
    'notification.smtp_enable': { name: '启用SMTP', type: 'switch', default: false },
    'notification.smtp_host': { name: 'SMTP主机', type: 'text', default: '' },
    'notification.smtp_port': { name: '端口', type: 'number', default: 465 },
    'notification.smtp_user': { name: 'SMTP账号', type: 'text', default: '' },
    'notification.smtp_pass': { name: 'SMTP密码', type: 'text', default: '' },
    'notification.auto_clean_read': { name: '自动瘦身已读通知', type: 'switch', default: true },
    'notification.route_alerts_to_admin': { name: '致命告警抄送短信', type: 'switch', default: false },
  }
}

const loadSettings = async () => {
  try {
    loading.value = true
    const { data: res } = await adminApi.getSettingList(currentGroup.value)

    // Reset local formData block
    for(const k in formData) delete formData[k]

    const defs = fieldDefaults[currentGroup.value] || {}
    // Seed defaults
    for(const key in defs) {
      formData[key] = defs[key].default
    }

    if (res.code === 200 && Array.isArray(res.data) && res.data.length > 0) {
      // Overlay database loaded values
      res.data.forEach(item => {
        if (item.valueType === 'switch') {
          formData[item.settingKey] = (item.settingValue === 'true')
        } else if (item.valueType === 'number') {
          formData[item.settingKey] = Number(item.settingValue) || 0
        } else {
          formData[item.settingKey] = item.settingValue
        }
      })
    }
  } catch (error) {
    ToastAlert.error({
      title: '加载底座环境失败',
      message: error instanceof Error ? error.message : '服务掉线'
    })
  } finally {
    loading.value = false
  }
}

const saveSettings = async () => {
  try {
    saving.value = true
    const defs = fieldDefaults[currentGroup.value] || {}

    // Map formData to SettingItem array
    const payload = Object.keys(defs).map(key => {
      const fieldDef = defs[key]
      let stringValue = String(formData[key] || '')
      if (fieldDef.type === 'switch') stringValue = Boolean(formData[key]) ? 'true' : 'false'

      return {
        groupKey: currentGroup.value,
        settingKey: key,
        settingName: fieldDef.name,
        settingValue: stringValue,
        valueType: fieldDef.type,
      }
    })

    const { data: res } = await adminApi.batchUpdateSettings(payload)
    if (res.code !== 200) throw new Error(res.message || '环境快照封存失败')

    ToastAlert.success({
      title: '部署架构应用成功',
      message: `${groupMeta.value.title} 参数已热刷新并在后方集群生效`
    })

  } catch (error) {
    ToastAlert.error({
      title: '部署失败',
      message: error instanceof Error ? error.message : '通信中断'
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
