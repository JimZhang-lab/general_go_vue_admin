/*
 * @Author: JimZhang
 * @Date: 2025-07-26 20:21:36
 * @LastEditors: 很拉风的James
 * @LastEditTime: 2025-07-27 14:44:43
 * @FilePath: /go-vue-general-admin/web/src/utils/request.ts
 * @Description: axios封装
 * 
 */

import axios, {type AxiosRequestConfig, type AxiosResponse } from 'axios'
import { createApp } from 'vue'
import Toast from '@/components/ui/Toast.vue'


// 创建Storage类来替代原有的storage对象
class Storage {
  getStorage() {
    return JSON.parse(window.localStorage.getItem(import.meta.env.VITE_NAME_SPACE as string) || "{}")
  }

  setItem(key: string, val: any) {
    let storage = this.getStorage()
    storage[key] = val
    window.localStorage.setItem(import.meta.env.VITE_NAME_SPACE as string, JSON.stringify(storage))
  }

  getItem(key: string) {
    return this.getStorage()[key]
  }

  clearItem(key: string) {
    let storage = this.getStorage()
    delete storage[key]
    window.localStorage.setItem(import.meta.env.VITE_NAME_SPACE as string, JSON.stringify(storage))
  }

  clearAll() {
    window.localStorage.clear()
  }
}

const storage = new Storage()



// 创建axios实例，添加全局配置
const service = axios.create({
  baseURL: import.meta.env.VITE_BASE_API,
  timeout: 8000, // 请求超时时间，单位 ms
  headers: {
    'Content-Type': 'application/json;charset=utf-8'
  }
})

// 请求拦截
service.interceptors.request.use(
  (config) => {
    const headers = config.headers

    // 只有在需要token的时候才添加Authorization头
    // 登录和验证码接口不需要token
    if (headers.isToken !== false && !headers.Authorization) {
        const token = storage.getItem('token')
        if (token) {
            headers.Authorization = `Bearer ${token}`
        }
    }
    return config
}
)

// 创建一个用于显示 Toast 的函数
const showToast = (variant: 'success' | 'error' | 'warning' | 'info', title: string, message: string) => {
  const container = document.createElement('div')
  const app = createApp(Toast, {
    modelValue: true,
    variant,
    title,
    message,
    showCancel: false
  })
  
  document.body.appendChild(container)
  const instance = app.mount(container)
  
  // 3秒后自动关闭
  setTimeout(() => {
    app.unmount()
    document.body.removeChild(container)
  }, 3000)
}

// 响应拦截
service.interceptors.response.use(
  (response: AxiosResponse): AxiosResponse | Promise<any> => {
    const { code, message } = response.data
    if (code === 403) {
      showToast('error', '错误', message)
      setTimeout(() => {
        storage.clearAll()
        window.location.href = '/adminLogin'
      }, 1500)
      return Promise.reject(response.data)
    } else if (code === 406) {
      showToast('error', '错误', message)
      setTimeout(() => {
        storage.clearAll()
        window.location.href = '/adminLogin'
      }, 1500)
      return Promise.reject(response.data)
    } else {
      return response
    }
  },
  (error) => {
    return Promise.reject(error)
  }
)
// service.interceptors.response.use(
//   (response) => {
//     const (code, data, message) = response.data
//     const res = response.data
//     if (code !== 200) {
//       showToast('error', '错误', res.message)
//       return Promise.reject(res.message)
//     } else {
//       return res
//     }
//   },
//   (error) => {
//       showToast('error', '错误', error.message)
//     return Promise.reject(error)
//   }
// )

export interface RequestOptions extends AxiosRequestConfig {
  method?: string;
  data?: any;
  params?: any;
}

// 请求核心函数
function request(options: RequestOptions): Promise<AxiosResponse> {
  options.method = options.method || 'get';
  if (options.method.toLowerCase() === 'get') {
    options.params = options.data;
  }
  service.defaults.baseURL = import.meta.env.VITE_BASE_API as string;
  return service(options);
}

export default request;