<template>
  <!--
    This example requires updating your template:

    ```
    <html class="h-full bg-white">
    <body class=""h-full">
    ```
  -->
  <div class="flex min-h-full flex-1 flex-col justify-center px-6 py-12 lg:px-8 bg-gray-50 dark:bg-gray-900">
    <div class="mt-12 sm:mx-auto sm:w-full sm:max-w-md">
      <!-- Logo -->
      <div class="flex justify-center mb-6">
        <div class="flex items-center justify-center w-16 h-16 bg-blue-600 rounded-2xl shadow-lg">
          <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
      </div>
      <h2 class="text-center text-3xl font-bold tracking-tight text-gray-900 dark:text-white/90">
        通用管理系统
      </h2>
      <p class="mt-2 text-center text-sm text-gray-600 dark:text-gray-400">
        请输入您的账户信息以访问管理系统
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <div class="rounded-2xl border border-gray-200 bg-white px-6 py-8 shadow-lg dark:border-gray-800 dark:bg-white/[0.03]">
        <form class="space-y-6"
        ref="loginFormRef"
        :rules="rules"
        :model="loginForm"
        @submit.prevent="loginBtn">
        <!-- Username field -->
        <div>
          <label for="adminUsername" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2"
            >用户名</label
          >
          <div class="mt-2 relative">
            
            <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
              <UserCircleIcon class="h-5 w-5 text-gray-400" />
            </div>
            <input
              type="adminUsername"
              name="adminUsername"
              id="adminUsername"
              autocomplete="adminUsername"
              v-model="loginForm.adminUsername"
              placeholder="请输入用户名或管理员ID"
              class="h-12 w-full rounded-xl border border-gray-300 bg-transparent py-3 pl-11 pr-4 text-sm text-gray-800 placeholder:text-gray-400 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-gray-700 dark:bg-gray-800 dark:text-white/90 dark:placeholder:text-gray-500"
            />
          </div>
        </div>

        <!-- Password field -->
        <div>
          <div class="flex items-center justify-between">
            <label for="password" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">密码</label>
          </div>

          <div class="mt-2">
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <KeyIcon class="h-5 w-5 text-gray-400" />
              </div>
              <input
                name="password"
                id="password"
                autocomplete="current-password"
                v-model="loginForm.password"
                :type="showPassword ? 'text' : 'password'"
                placeholder="请输入密码"
                class="h-12 w-full rounded-xl border border-gray-300 bg-transparent py-3 pl-11 pr-12 text-sm text-gray-800 placeholder:text-gray-400 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-gray-700 dark:bg-gray-800 dark:text-white/90 dark:placeholder:text-gray-500"
              />
              <span
                @click="togglePasswordVisibility"
                class="absolute inset-y-0 right-0 flex items-center pr-3 cursor-pointer"
              >
                <svg
                  v-if="!showPassword"
                  class="fill-current"
                  width="20"
                  height="20"
                  viewBox="0 0 20 20"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M10.0002 13.8619C7.23361 13.8619 4.86803 12.1372 3.92328 9.70241C4.86804 7.26761 7.23361 5.54297 10.0002 5.54297C12.7667 5.54297 15.1323 7.26762 16.0771 9.70243C15.1323 12.1372 12.7667 13.8619 10.0002 13.8619ZM10.0002 4.04297C6.48191 4.04297 3.49489 6.30917 2.4155 9.4593C2.3615 9.61687 2.3615 9.78794 2.41549 9.94552C3.49488 13.0957 6.48191 15.3619 10.0002 15.3619C13.5184 15.3619 16.5055 13.0957 17.5849 9.94555C17.6389 9.78797 17.6389 9.6169 17.5849 9.45932C16.5055 6.30919 13.5184 4.04297 10.0002 4.04297ZM9.99151 7.84413C8.96527 7.84413 8.13333 8.67606 8.13333 9.70231C8.13333 10.7286 8.96527 11.5605 9.99151 11.5605H10.0064C11.0326 11.5605 11.8646 10.7286 11.8646 9.70231C11.8646 8.67606 11.0326 7.84413 10.0064 7.84413H9.99151Z"
                    fill="#98A2B3"
                  />
                </svg>
                <svg
                  v-else
                  class="fill-current"
                  width="20"
                  height="20"
                  viewBox="0 0 20 20"
                  fill="none"
                  xmlns="http://www.w3.org/2000/svg"
                >
                  <path
                    fill-rule="evenodd"
                    clip-rule="evenodd"
                    d="M4.63803 3.57709C4.34513 3.2842 3.87026 3.2842 3.57737 3.57709C3.28447 3.86999 3.28447 4.34486 3.57737 4.63775L4.85323 5.91362C3.74609 6.84199 2.89363 8.06395 2.4155 9.45936C2.3615 9.61694 2.3615 9.78801 2.41549 9.94558C3.49488 13.0957 6.48191 15.3619 10.0002 15.3619C11.255 15.3619 12.4422 15.0737 13.4994 14.5598L15.3625 16.4229C15.6554 16.7158 16.1302 16.7158 16.4231 16.4229C16.716 16.13 16.716 15.6551 16.4231 15.3622L4.63803 3.57709ZM12.3608 13.4212L10.4475 11.5079C10.3061 11.5423 10.1584 11.5606 10.0064 11.5606H9.99151C8.96527 11.5606 8.13333 10.7286 8.13333 9.70237C8.13333 9.5461 8.15262 9.39434 8.18895 9.24933L5.91885 6.97923C5.03505 7.69015 4.34057 8.62704 3.92328 9.70247C4.86803 12.1373 7.23361 13.8619 10.0002 13.8619C10.8326 13.8619 11.6287 13.7058 12.3608 13.4212ZM16.0771 9.70249C15.7843 10.4569 15.3552 11.1432 14.8199 11.7311L15.8813 12.7925C16.6329 11.9813 17.2187 11.0143 17.5849 9.94561C17.6389 9.78803 17.6389 9.61696 17.5849 9.45938C16.5055 6.30925 13.5184 4.04303 10.0002 4.04303C9.13525 4.04303 8.30244 4.17999 7.52218 4.43338L8.75139 5.66259C9.1556 5.58413 9.57311 5.54303 10.0002 5.54303C12.7667 5.54303 15.1323 7.26768 16.0771 9.70249Z"
                    fill="#98A2B3"
                  />
                </svg>
              </span>
            </div>
            
          </div>
        </div>

        <!-- captcha -->
        <div>
          <div class="flex items-center justify-between">
            <label for="captcha" class="block text-sm font-medium text-gray-700 dark:text-gray-300 mb-2">验证码</label>
          </div>
          <div class="mt-2 flex items-center gap-3">
            <div class="relative flex-1">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <ShieldCheckIcon class="h-5 w-5 text-gray-400" />
              </div>
              <input
                name="captcha"
                id="captcha"
                autocomplete="current-captcha"
                placeholder="请输入验证码"
                v-model="loginForm.captcha"
                class="h-12 w-full rounded-xl border border-gray-300 bg-transparent py-3 pl-11 pr-4 text-sm text-gray-800 placeholder:text-gray-400 focus:border-blue-500 focus:outline-none focus:ring-2 focus:ring-blue-500/20 dark:border-gray-700 dark:bg-gray-800 dark:text-white/90 dark:placeholder:text-gray-500"
              />
            </div>
            <img 
            
            @click="getCaptcha()"
            :src="image"
            class="h-12 w-32 object-contain rounded-xl border border-gray-300 cursor-pointer hover:border-blue-500 transition-colors dark:border-gray-700" />
          </div>
          <div class="mt-2">
              <label
                for="keepLoggedIn"
                class="flex items-center text-sm font-normal text-gray-700 cursor-pointer select-none dark:text-gray-300"
              >
                <div class="relative">
                  <input v-model="keepLoggedIn" type="checkbox" id="keepLoggedIn" class="sr-only" />
                  <div
                    :class="
                      keepLoggedIn
                        ? 'border-blue-500 bg-blue-500'
                        : 'bg-transparent border-gray-300 dark:border-gray-700'
                    "
                    class="mr-3 flex h-5 w-5 items-center justify-center rounded-md border-[1.25px] transition-colors"
                  >
                    <span :class="keepLoggedIn ? '' : 'opacity-0'">
                      <svg
                        width="14"
                        height="14"
                        viewBox="0 0 14 14"
                        fill="none"
                        xmlns="http://www.w3.org/2000/svg"
                      >
                        <path
                          d="M11.6666 3.5L5.24992 9.91667L2.33325 7"
                          stroke="white"
                          stroke-width="1.94437"
                          stroke-linecap="round"
                          stroke-linejoin="round"
                        />
                      </svg>
                    </span>
                  </div>
                </div>
                记住我
              </label>
            </div>
        </div>
        <div>
          <button
            type="submit"
            class="flex w-full justify-center rounded-xl bg-blue-600 px-4 py-3 text-sm font-semibold text-white shadow-lg hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
          >
            登录
          </button>
        </div>
      </form>

      <!-- 底部信息 -->
      <div class="mt-6 text-center">
        <p class="text-xs text-gray-500 dark:text-gray-400">
          © 2025 通用管理系统. 保留所有权利.
        </p>
      </div>
    </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref,reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { UserCircleIcon } from '@/icons'
import { KeyIcon, ShieldCheckIcon } from '@heroicons/vue/24/outline'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'
import { useMainStore } from '@/store'
import { AuthUtils } from '@/utils/auth'

const showPassword = ref(false)
const keepLoggedIn = ref(false)
const router = useRouter()
const store = useMainStore()

// 登录表单
const loginFormRef = ref()

const image = ref('')

// 表单数据
const loginForm = reactive({
    adminUsername: '',
    password: '',
    captcha: '',
    idKey: ''
})

// 获取验证码
const getCaptcha = async () => {
  try {
    const {data: res} = await adminApi.captcha()
    console.log(res)
    if (res.code === 200) {
      image.value = res.data.image
      loginForm.idKey = res.data.idKey
    } else {
      ToastAlert.error({
        "title": "获取验证码失败",
        "message": res.message
      })
    }
  }catch (error) {
    console.log("获取验证码失败")
    ToastAlert.error({
      "title": "获取验证码失败",
      "message": "网络异常，请重试"
    })
  }
}


// 表单验证规则
const rules = {
    adminUsername: [
        {
            required: true, message: "请输入账号", trigger: "blur"
        }
    ],
    password: [
        {
            required: true, message: "请输入密码", trigger: "blur"
        }
    ],
    captcha: [
        {
            required: true, message: "请输入验证码", trigger: "blur"
        }
    ]
}


const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

// 登录方法
const loginBtn = async () => {
    // 移除不存在的validate方法调用，直接进行表单验证
    // 验证必填字段
    if (!loginForm.adminUsername || !loginForm.password || !loginForm.captcha) {
        ToastAlert.error({
            "title": "表单验证失败",
            "message": "请填写所有必填字段"
        });
        return;
    }
    
    try {
        // 使用已定义的API方法而不是直接使用fetch
        const loginData: {
            username: string;
            password: string;
            image: string;
            idKey: string;
        } = {
            username: loginForm.adminUsername,
            password: loginForm.password,
            image: loginForm.captcha,
            idKey: loginForm.idKey
        }

        const res = await adminApi.login(loginData)

        if (res.data.code !== 200) {
            ToastAlert.error({
                "title": "登录失败",
                "message": res.data.message || "登录失败，请重试"
            })
            // 登录失败时刷新验证码
            await getCaptcha()
            console.log("登录失败")
        } else {
            // 先保存登录信息到store
            store.saveSysAdmin(res.data.data.sysAdmin)
            store.saveToken(res.data.data.token)
            store.saveLeftMenuList(res.data.data.leftMenuList)
            store.savePermissionList(res.data.data.permissionList)

            // 设置"保持登录"状态
            AuthUtils.setKeepLoggedIn(keepLoggedIn.value)

            // 设置最后活动时间
            AuthUtils.refreshSession()

            // 显示成功提示
            ToastAlert.success({
                "title": "登录成功",
                "message": "正在跳转到仪表板...",
                "duration": 1500
            })

            // 延迟跳转，确保用户看到成功提示
            setTimeout(() => {
                // 检查是否有重定向参数
                const redirect = router.currentRoute.value.query.redirect as string
                router.push(redirect || "/auth")
            }, 1000)

            console.log("登录成功")
        }
    } catch (error) {
        // ToastAlert.error({
        //     "title": "登录失败",
        //     "message": String(error.message)
        // })
        console.log(error)
    }
}


onMounted(() => {
  // console.log(getCaptcha())
  getCaptcha()
})
</script>
