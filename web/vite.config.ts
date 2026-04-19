/*
 * @Author: JimZhang
 * @Date: 2025-05-15 15:06:52
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2026-04-19 11:51:17
 * @FilePath: /general_go_vue_admin/web/vite.config.ts
 * @Description: 
 * 
 */
import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import vueJsx from '@vitejs/plugin-vue-jsx'
import vueDevTools from 'vite-plugin-vue-devtools'
import tailwindcss from "@tailwindcss/vite";
// https://vite.dev/config/
export default defineConfig({
  plugins: [
    tailwindcss(),
    vue(),
    vueJsx(),
    vueDevTools(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    },
  },
  server: {
    // 端口
    port: 3000,
    // 启动时是否自动打开浏览器
    open: false,
    // 允许跨域
    cors: true,
    // 允许代理
    proxy: {
      '/api': {
        target: 'http://localhost:8368',
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '')
      }
    },
  },
  build: {
    // charts-core（ApexCharts）体积较大，提升告警阈值避免CI噪音
    chunkSizeWarningLimit: 650,
    rollupOptions: {
      output: {
        manualChunks(id) {
          if (!id.includes('node_modules')) return

          if (id.includes('/vue/') || id.includes('/vue-router/') || id.includes('/pinia/')) {
            return 'framework'
          }

          if (id.includes('/@fullcalendar/')) {
            return 'calendar'
          }

          if (id.includes('/apexcharts/')) {
            return 'charts-core'
          }

          if (id.includes('/vue3-apexcharts/')) {
            return 'charts-vue'
          }

          if (id.includes('/axios/') || id.includes('/qs/')) {
            return 'network'
          }

          if (id.includes('/@heroicons/') || id.includes('/lucide-vue-next/')) {
            return 'icons'
          }

          if (id.includes('/dropzone/') || id.includes('/flatpickr/') || id.includes('/swiper/')) {
            return 'ui-kit'
          }
        }
      }
    }
  }
})
