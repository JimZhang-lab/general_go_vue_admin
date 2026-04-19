<template>
  <div class="space-y-6 animate-in fade-in zoom-in-95 duration-200">
    <div class="grid gap-6 md:grid-cols-2">
      <!-- 站点基本信息 -->
      <div class="rounded-2xl border border-gray-200 p-5 dark:border-gray-800 bg-white dark:bg-gray-900/30">
        <div class="flex items-center gap-3 mb-4">
          <div class="p-2 bg-blue-50 text-blue-600 rounded-lg dark:bg-blue-500/15 dark:text-blue-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9a9 9 0 01-9-9m9 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9"></path></svg>
          </div>
          <h3 class="text-base font-semibold text-gray-900 dark:text-white/90">跨域与站点标识</h3>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">系统名称 (Site Name)</label>
            <input v-model="formData['basic.site_name']" type="text" placeholder="例如: 企业级权限管理平台" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
            <p class="mt-1 text-xs text-gray-500">将显示在浏览器标签页、顶部导航栏及登录页面。</p>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">系统副标题 (Slogan)</label>
            <input v-model="formData['basic.site_slogan']" type="text" placeholder="例如: 高效、安全、稳定的中台基座" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">主域名 (Primary Domain)</label>
            <div class="flex">
              <span class="inline-flex items-center px-3 rounded-l-xl border border-r-0 border-gray-300 bg-gray-50 text-gray-500 text-sm dark:bg-gray-800 dark:border-gray-700">https://</span>
              <input v-model="formData['basic.site_domain']" type="text" placeholder="admin.example.com" class="w-full rounded-r-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
            </div>
          </div>
        </div>
      </div>

      <!-- 企业与法律信息 -->
      <div class="rounded-2xl border border-gray-200 p-5 dark:border-gray-800 bg-white dark:bg-gray-900/30">
        <div class="flex items-center gap-3 mb-4">
          <div class="p-2 bg-purple-50 text-purple-600 rounded-lg dark:bg-purple-500/15 dark:text-purple-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path></svg>
          </div>
          <h3 class="text-base font-semibold text-gray-900 dark:text-white/90">合规与备案说明</h3>
        </div>

        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">公司/持有者名称</label>
            <input v-model="formData['basic.company_name']" type="text" placeholder="例如: 某某科技股份有限公司" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">ICP 备案号</label>
            <input v-model="formData['basic.icp_beian']" type="text" placeholder="例如: 京ICP备12345678号-1" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-1">公安联网备案</label>
            <input v-model="formData['basic.gongan_beian']" type="text" placeholder="例如: 京公网安备 11000002000001号" class="w-full rounded-xl border border-gray-300 px-4 py-2 text-sm focus:border-brand-300 focus:ring-3 focus:ring-brand-500/10 dark:border-gray-700 dark:bg-gray-800 dark:text-gray-100" />
          </div>
        </div>
      </div>

      <!-- 运行时策略 -->
      <div class="rounded-2xl border border-gray-200 p-5 dark:border-gray-800 bg-white dark:bg-gray-900/30 md:col-span-2">
        <div class="flex items-center gap-3 mb-4">
          <div class="p-2 bg-green-50 text-green-600 rounded-lg dark:bg-green-500/15 dark:text-green-400">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 110-4m0 4v2m0-6V4"></path></svg>
          </div>
          <h3 class="text-base font-semibold text-gray-900 dark:text-white/90">业务逻辑控制</h3>
        </div>

        <div class="grid gap-6 sm:grid-cols-2">
          <label class="flex items-center justify-between rounded-xl border border-gray-200 px-4 py-3 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors w-full cursor-pointer">
            <div>
              <span class="block text-sm font-medium text-gray-900 dark:text-gray-100">开放注册接口</span>
              <span class="block text-xs text-gray-500 dark:text-gray-400 mt-0.5">允许外部人员通过 /auth/register 自行请求账号</span>
            </div>
            <input v-model="formData['basic.allow_register']" type="checkbox" class="h-5 w-5 rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
          </label>
          <label class="flex items-center justify-between rounded-xl border border-gray-200 px-4 py-3 dark:border-gray-700 hover:bg-gray-50 dark:hover:bg-gray-800/50 transition-colors w-full cursor-pointer">
            <div>
              <span class="block text-sm font-medium text-gray-900 dark:text-gray-100">强制全系统水印</span>
              <span class="block text-xs text-gray-500 dark:text-gray-400 mt-0.5">在后台所有页面浮层渲染当前登录者账号水印防泄露</span>
            </div>
            <input v-model="formData['basic.force_watermark']" type="checkbox" class="h-5 w-5 rounded border-gray-300 text-brand-600 focus:ring-brand-500" />
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
