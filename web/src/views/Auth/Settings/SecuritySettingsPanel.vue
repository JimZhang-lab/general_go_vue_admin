<template>
  <div class="space-y-6 animate-in fade-in zoom-in-95 duration-200">
    <div class="grid gap-6 md:grid-cols-2">
      <!-- 账号封锁策略 -->
      <div class="rounded-2xl border border-red-200 bg-red-50/30 p-5 dark:border-red-900/50 dark:bg-red-900/10">
        <div class="flex items-center gap-3 mb-4">
          <div class="p-2 bg-red-100 text-red-600 rounded-lg dark:bg-red-500/20 dark:text-red-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"></path></svg>
          </div>
          <h3 class="text-base font-semibold text-gray-900 dark:text-white/90">账号安全防护</h3>
        </div>

        <div class="space-y-5">
          <div class="flex items-center justify-between">
            <div>
              <label class="text-sm font-medium text-gray-900 dark:text-gray-100">允许连续登录失败次数</label>
              <p class="text-xs text-gray-500">超过此限制后，账号会被临时冻结防爆破</p>
            </div>
            <div class="w-24">
              <input v-model="formData['security.login_failed_limit']" type="number" min="1" max="20" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
            </div>
          </div>

          <div class="flex items-center justify-between border-t border-red-100/50 dark:border-red-900/40 pt-4">
            <div>
              <label class="text-sm font-medium text-gray-900 dark:text-gray-100">冻结锁定时间 (分钟)</label>
              <p class="text-xs text-gray-500">防爆破锁被触发后账号解封冷却时长</p>
            </div>
            <div class="w-24">
              <input v-model="formData['security.lock_minutes']" type="number" min="1" max="1440" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
            </div>
          </div>
        </div>
      </div>

      <!-- 密码与会话强度规则 -->
      <div class="rounded-2xl border border-gray-200 p-5 bg-white dark:border-gray-800 dark:bg-gray-900/30">
        <div class="flex items-center gap-3 mb-4">
          <div class="p-2 bg-indigo-50 text-indigo-600 rounded-lg dark:bg-indigo-500/20 dark:text-indigo-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z"></path></svg>
          </div>
          <h3 class="text-base font-semibold text-gray-900 dark:text-white/90">密码与会话强度</h3>
        </div>

        <div class="space-y-4">
           <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">全局 Session 过期时间 (小时)</label>
            <input v-model="formData['security.session_timeout_hours']" type="number" min="1" max="720" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
          </div>

          <label class="flex items-center justify-between rounded-xl border border-gray-200 px-4 py-3 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors cursor-pointer mt-2">
            <div>
              <span class="block text-sm font-medium text-gray-900 dark:text-gray-100">强制最高级强密码策略</span>
              <span class="block text-xs text-gray-500 dark:text-gray-400 mt-0.5">要求必须包含大小写字母、数字及特殊符号（只对新修改或注册生效）</span>
            </div>
            <input v-model="formData['security.force_strong_pwd']" type="checkbox" class="h-5 w-5 rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
          </label>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { inject } from 'vue'

const formData = inject<Record<string, any>>('settingsFormData', {})
</script>
