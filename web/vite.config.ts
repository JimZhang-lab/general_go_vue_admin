/*
 * @Author: JimZhang
 * @Date: 2025-05-15 15:06:52
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-27 21:21:24
 * @FilePath: /go-vue-general-admin/web/vite.config.ts
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
        target: 'http://localhost:8080',
        changeOrigin: true,
        // rewrite: (path) => path.replace(/^\/api/, '')
      }
    },
  },
})
