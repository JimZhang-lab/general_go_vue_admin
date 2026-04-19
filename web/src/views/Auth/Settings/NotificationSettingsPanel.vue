<template>
  <div class="space-y-6 animate-in fade-in zoom-in-95 duration-200">
    <div class="grid gap-6 md:grid-cols-2">
      <!-- 邮件 SMTP 服务配置 -->
      <div class="rounded-2xl border border-gray-200 p-5 bg-white dark:border-gray-800 dark:bg-gray-900/30">
        <div class="flex items-center justify-between mb-5">
          <div class="flex items-center gap-3">
            <div class="p-2 bg-yellow-50 text-yellow-600 rounded-lg dark:bg-yellow-500/20 dark:text-yellow-400">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path></svg>
            </div>
            <h3 class="text-base font-semibold text-gray-900 dark:text-white/90">邮件发送 SMTP 配置</h3>
          </div>
          <label class="flex items-center cursor-pointer">
            <div class="relative">
              <input v-model="formData['notification.smtp_enable']" type="checkbox" class="sr-only" />
              <div class="block bg-gray-200 dark:bg-gray-700 w-10 h-6 pl-1 pr-1 rounded-full border border-transparent transition-colors duration-200 ease-in-out" :class="{'bg-brand-500': formData['notification.smtp_enable']}"></div>
              <div class="dot absolute left-1 top-1 bg-white w-4 h-4 rounded-full transition transform duration-200 ease-in-out shadow-sm" :class="{'translate-x-4': formData['notification.smtp_enable']}"></div>
            </div>
          </label>
        </div>

        <div class="space-y-4" :class="{'opacity-50 pointer-events-none': !formData['notification.smtp_enable']}">
          <div class="grid gap-4 sm:grid-cols-3">
            <div class="sm:col-span-2">
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">SMTP 主机</label>
              <input v-model="formData['notification.smtp_host']" type="text" placeholder="smtp.example.com" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">端口</label>
              <input v-model.number="formData['notification.smtp_port']" type="number" placeholder="465" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">发件人账号/邮箱</label>
            <input v-model="formData['notification.smtp_user']" type="email" placeholder="no-reply@example.com" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">授权码 / 密码</label>
            <input v-model="formData['notification.smtp_pass']" type="password" placeholder="••••••••••••" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
          </div>
        </div>
      </div>

       <!-- 全局通知流控 -->
       <div class="rounded-2xl border border-gray-200 p-5 bg-white dark:border-gray-800 dark:bg-gray-900/30">
        <div class="flex items-center gap-3 mb-5">
          <div class="p-2 bg-teal-50 text-teal-600 rounded-lg dark:bg-teal-500/20 dark:text-teal-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"></path></svg>
          </div>
          <h3 class="text-base font-semibold text-gray-900 dark:text-white/90">站内通知流控</h3>
        </div>

        <div class="space-y-4">
          <label class="flex items-center justify-between rounded-xl border border-gray-200 px-4 py-3 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors w-full cursor-pointer">
            <div>
              <span class="block text-sm font-medium text-gray-900 dark:text-gray-100">启用全局系统通知</span>
              <span class="block text-xs text-gray-500 dark:text-gray-400 mt-0.5">当关闭此项后，系统将彻底拦截所有站内通知消息推送接口</span>
            </div>
            <input v-model="formData['notification.enable_site_notice']" type="checkbox" class="h-5 w-5 rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
          </label>

          <label class="flex items-center justify-between rounded-xl border border-gray-200 px-4 py-3 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors w-full cursor-pointer mt-2">
            <div>
              <span class="block text-sm font-medium text-gray-900 dark:text-gray-100">通知清理策略</span>
              <span class="block text-xs text-gray-500 dark:text-gray-400 mt-0.5">当勾选后，系统会自动清理 30 天前已读的常规通知以释放数据库负载</span>
            </div>
            <input v-model="formData['notification.auto_clean_read']" type="checkbox" class="h-5 w-5 rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
          </label>

          <label class="flex items-center justify-between rounded-xl border border-gray-200 px-4 py-3 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors w-full cursor-pointer mt-2">
            <div>
              <span class="block text-sm font-medium text-gray-900 dark:text-gray-100">告警强路由系统短信</span>
              <span class="block text-xs text-gray-500 dark:text-gray-400 mt-0.5">极端资源异常或被爆破时，将站内信平移发送给超管</span>
            </div>
            <input v-model="formData['notification.route_alerts_to_admin']" type="checkbox" class="h-5 w-5 rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
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
