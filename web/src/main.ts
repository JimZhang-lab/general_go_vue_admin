import './assets/main.css'
// Import Swiper styles
import 'swiper/css'
import 'swiper/css/navigation'
import 'swiper/css/pagination'
import 'jsvectormap/dist/jsvectormap.css'
import 'flatpickr/dist/flatpickr.css'

import { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import VueApexCharts from 'vue3-apexcharts'
import { createPinia } from 'pinia'
// import api from '@/api/system'
// import request from '@/utils/request'
// import storage from '@/utils/storage'
// import useMainStore from '@/store'


// 获取工作环境根目录下的的环境变量 .env.dev 中的 VITE_BASE_API
// console.log("Environment variable: ", import.meta.env.VITE_BASE_API)


const pinia = createPinia()
const app = createApp(App)
// app.use(request) - 错误：request不是Vue插件
// app.use(storage)
// app.use(useMainStore)
// app.use(api)


app.use(router)
// @ts-ignore
app.use(VueApexCharts)
app.use(pinia)

app.mount('#app')