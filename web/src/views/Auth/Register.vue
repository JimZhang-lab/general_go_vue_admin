/*
 * @Author: JimZhang
 * @Date: 2026-04-19 14:26:57
 * @LastEditors: JimZhang
 * @LastEditTime: 2026-04-19 14:26:57
 * @FilePath: /web/src/views/Auth/Register.vue
 * @Description: 
 * 
 */
<template>
  <div class="relative min-h-screen overflow-hidden bg-gradient-to-br from-[#ecfbff] via-[#f3f8ff] to-[#e7f1ff]">
    <div class="pointer-events-none absolute inset-0">
      <div class="water-gradient-layer absolute inset-0"></div>
      <div class="water-orb orb-one absolute -left-20 -top-24 h-80 w-80 rounded-full"></div>
      <div class="water-orb orb-two absolute -right-24 top-1/3 h-96 w-96 rounded-full"></div>
      <div class="water-orb orb-three absolute left-1/3 bottom-[-8rem] h-[22rem] w-[22rem] rounded-full"></div>
      <div class="water-wave wave-one absolute left-[-18%] top-[15%] h-44 w-[140%]"></div>
      <div class="water-wave wave-two absolute left-[-22%] top-[42%] h-56 w-[145%]"></div>
      <div class="water-wave wave-three absolute left-[-14%] top-[72%] h-40 w-[132%]"></div>
      <svg
        class="sine-line line-one absolute left-[-24%] top-[8%] h-20 w-[156%]"
        viewBox="0 0 1600 120"
        preserveAspectRatio="none"
        aria-hidden="true"
      >
        <path class="sine-path" d="M0 60 Q40 20 80 60 T160 60 T240 60 T320 60 T400 60 T480 60 T560 60 T640 60 T720 60 T800 60 T880 60 T960 60 T1040 60 T1120 60 T1200 60 T1280 60 T1360 60 T1440 60 T1520 60 T1600 60" />
      </svg>
      <svg
        class="sine-line line-two absolute left-[-30%] top-[31%] h-28 w-[170%]"
        viewBox="0 0 1600 120"
        preserveAspectRatio="none"
        aria-hidden="true"
      >
        <path class="sine-path" d="M0 60 Q40 20 80 60 T160 60 T240 60 T320 60 T400 60 T480 60 T560 60 T640 60 T720 60 T800 60 T880 60 T960 60 T1040 60 T1120 60 T1200 60 T1280 60 T1360 60 T1440 60 T1520 60 T1600 60" />
      </svg>
      <svg
        class="sine-line line-three absolute left-[-18%] top-[57%] h-24 w-[148%]"
        viewBox="0 0 1600 120"
        preserveAspectRatio="none"
        aria-hidden="true"
      >
        <path class="sine-path" d="M0 60 Q40 20 80 60 T160 60 T240 60 T320 60 T400 60 T480 60 T560 60 T640 60 T720 60 T800 60 T880 60 T960 60 T1040 60 T1120 60 T1200 60 T1280 60 T1360 60 T1440 60 T1520 60 T1600 60" />
      </svg>
      <svg
        class="sine-line line-four absolute left-[-22%] top-[84%] h-20 w-[160%]"
        viewBox="0 0 1600 120"
        preserveAspectRatio="none"
        aria-hidden="true"
      >
        <path class="sine-path" d="M0 60 Q40 20 80 60 T160 60 T240 60 T320 60 T400 60 T480 60 T560 60 T640 60 T720 60 T800 60 T880 60 T960 60 T1040 60 T1120 60 T1200 60 T1280 60 T1360 60 T1440 60 T1520 60 T1600 60" />
      </svg>
      <div class="water-ripple ripple-one absolute left-[14%] top-[22%] h-44 w-44 rounded-full"></div>
      <div class="water-ripple ripple-two absolute right-[18%] top-[28%] h-64 w-64 rounded-full"></div>
      <div class="water-ripple ripple-three absolute left-[52%] bottom-[10%] h-52 w-52 rounded-full"></div>
      <span class="bubble bubble-1"></span>
      <span class="bubble bubble-2"></span>
      <span class="bubble bubble-3"></span>
      <span class="bubble bubble-4"></span>
      <span class="bubble bubble-5"></span>
      <span class="bubble bubble-6"></span>
      <span class="bubble bubble-7"></span>
      <span class="bubble bubble-8"></span>
    </div>

    <div class="relative grid min-h-screen lg:grid-cols-[1.05fr_0.95fr]">
      <section class="hidden px-10 py-12 text-slate-800 lg:flex lg:flex-col lg:justify-between xl:px-16">
        <div>
          <div class="inline-flex items-center gap-3 rounded-full border border-cyan-100 bg-white/70 px-4 py-2 text-sm font-semibold text-cyan-700 backdrop-blur-xl">
            <span class="h-2.5 w-2.5 rounded-full bg-cyan-500 shadow-[0_0_0_5px_rgba(6,182,212,0.18)]"></span>
            安全后台工作台
          </div>
          <div class="mt-10 max-w-xl">
            <p class="text-sm font-semibold uppercase tracking-[0.35em] text-cyan-600/80">General Admin</p>
            <h1 class="mt-5 text-5xl font-semibold leading-tight text-slate-900">
              像水一样流动的登录界面，
              <span class="bg-gradient-to-r from-cyan-600 to-blue-600 bg-clip-text text-transparent">让后台入口更有生命力。</span>
            </h1>
            <p class="mt-6 text-lg leading-8 text-slate-600">
              统一管理员登录、验证码校验和菜单权限分发，减少误操作，也让系统入口更可靠。
            </p>
          </div>
        </div>

        <div class="grid gap-4 xl:grid-cols-3">
          <article
            v-for="item in showcaseItems"
            :key="item.title"
            class="rounded-3xl border border-cyan-100/80 bg-white/70 p-5 shadow-[0_18px_34px_rgba(2,132,199,0.12)] backdrop-blur-xl"
          >
            <div class="flex items-center gap-3">
              <div class="flex h-11 w-11 items-center justify-center rounded-2xl bg-cyan-100 text-cyan-700">
                <CheckCircleIcon class="h-6 w-6" />
              </div>
              <h2 class="text-base font-semibold text-slate-800">{{ item.title }}</h2>
            </div>
            <p class="mt-4 text-sm leading-7 text-slate-600">{{ item.description }}</p>
          </article>
        </div>
      </section>

      <section class="flex items-center justify-center px-6 py-10 sm:px-8 lg:px-10">
        <div class="w-full max-w-xl">
          <div class="mb-6 lg:hidden">
            <div class="inline-flex items-center gap-2 rounded-full border border-cyan-100 bg-white/80 px-4 py-2 text-sm font-semibold text-cyan-700 backdrop-blur-xl">
              <span class="h-2.5 w-2.5 rounded-full bg-cyan-500 shadow-[0_0_0_5px_rgba(6,182,212,0.18)]"></span>
              General Admin
            </div>
          </div>

          <div class="rounded-[32px] border border-cyan-100/80 bg-white/86 p-7 shadow-[0_34px_80px_rgba(2,132,199,0.24)] backdrop-blur-2xl sm:p-8">
            <div class="flex items-start justify-between gap-4">
              <div>
                <p class="text-sm font-semibold uppercase tracking-[0.32em] text-cyan-600">管理员入口</p>
                <h2 class="mt-3 text-3xl font-semibold text-slate-900">创建账号</h2>
                <p class="mt-3 text-sm leading-7 text-slate-500">
                  仅需填写必要信息即可注册，注册成功后将自动回到登录页。
                </p>
              </div>
              <div class="hidden h-14 w-14 items-center justify-center rounded-3xl bg-gradient-to-br from-cyan-500 to-blue-600 text-white shadow-[0_10px_25px_rgba(14,116,144,0.35)] sm:flex">
                <UserPlusIcon class="h-7 w-7" />
              </div>
            </div>

            <form class="mt-8 space-y-5" novalidate @submit.prevent="submitRegister">
              <div>
                <label for="register-username" class="mb-2 block text-sm font-medium text-slate-700">账号</label>
                <div class="relative">
                  <UserCircleIcon class="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400" />
                  <input
                    id="register-username"
                    v-model="registerForm.username"
                    autocomplete="username"
                    type="text"
                    placeholder="4-20 位字母、数字或下划线"
                    :class="inputClass('username')"
                    @blur="validateField('username')"
                    @input="handleTextInput('username')"
                  />
                </div>
                <p v-if="fieldErrors.username" class="mt-2 text-xs font-medium text-rose-500">
                  {{ fieldErrors.username }}
                </p>
              </div>

              <div class="grid gap-5 sm:grid-cols-2">
                <div>
                  <label for="register-password" class="mb-2 block text-sm font-medium text-slate-700">密码</label>
                  <div class="relative">
                    <KeyIcon class="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400" />
                    <input
                      id="register-password"
                      v-model="registerForm.password"
                      autocomplete="new-password"
                      :type="showPassword ? 'text' : 'password'"
                      placeholder="至少 6 位，包含字母和数字"
                      :class="inputClass('password')"
                      @blur="validateField('password')"
                      @input="handlePasswordInput"
                    />
                    <button
                      type="button"
                      class="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 transition hover:text-slate-700"
                      @click="showPassword = !showPassword"
                    >
                      <EyeIcon v-if="!showPassword" class="h-5 w-5" />
                      <EyeSlashIcon v-else class="h-5 w-5" />
                    </button>
                  </div>
                  <p v-if="fieldErrors.password" class="mt-2 text-xs font-medium text-rose-500">
                    {{ fieldErrors.password }}
                  </p>
                </div>

                <div>
                  <label for="register-confirm-password" class="mb-2 block text-sm font-medium text-slate-700">确认密码</label>
                  <div class="relative">
                    <ShieldCheckIcon class="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400" />
                    <input
                      id="register-confirm-password"
                      v-model="registerForm.confirmPassword"
                      autocomplete="new-password"
                      :type="showConfirmPassword ? 'text' : 'password'"
                      placeholder="请再次输入密码"
                      :class="inputClass('confirmPassword')"
                      @blur="validateField('confirmPassword')"
                      @input="handleTextInput('confirmPassword')"
                    />
                    <button
                      type="button"
                      class="absolute right-4 top-1/2 -translate-y-1/2 text-slate-400 transition hover:text-slate-700"
                      @click="showConfirmPassword = !showConfirmPassword"
                    >
                      <EyeIcon v-if="!showConfirmPassword" class="h-5 w-5" />
                      <EyeSlashIcon v-else class="h-5 w-5" />
                    </button>
                  </div>
                  <p v-if="fieldErrors.confirmPassword" class="mt-2 text-xs font-medium text-rose-500">
                    {{ fieldErrors.confirmPassword }}
                  </p>
                </div>
              </div>

              <div class="rounded-2xl border border-cyan-100 bg-cyan-50/50 px-4 py-4">
                <div class="flex items-center justify-between gap-4">
                  <div>
                    <p class="text-sm font-medium text-slate-700">密码强度</p>
                    <p class="mt-1 text-xs text-slate-500">{{ passwordStrength.label }}</p>
                  </div>
                  <span class="rounded-full px-3 py-1 text-xs font-semibold" :class="passwordStrength.badgeClass">
                    {{ passwordStrength.tag }}
                  </span>
                </div>
                <div class="mt-3 grid grid-cols-4 gap-2">
                  <span
                    v-for="index in 4"
                    :key="index"
                    class="h-2 rounded-full"
                    :class="index <= passwordStrength.score ? passwordStrength.barClass : 'bg-slate-200'"
                  ></span>
                </div>
              </div>

              <div>
                <div class="mb-2 flex items-center justify-between gap-3">
                  <label for="register-captcha" class="block text-sm font-medium text-slate-700">验证码</label>
                  <button
                    type="button"
                    class="inline-flex items-center gap-1 text-xs font-medium text-cyan-600 transition hover:text-cyan-700"
                    @click="getCaptcha()"
                  >
                    <ArrowPathIcon class="h-4 w-4" />
                    刷新验证码
                  </button>
                </div>

                <div class="grid gap-3 sm:grid-cols-[1fr_160px]">
                  <div class="relative">
                    <SparklesIcon class="pointer-events-none absolute left-4 top-1/2 h-5 w-5 -translate-y-1/2 text-slate-400" />
                    <input
                      id="register-captcha"
                      v-model="registerForm.captcha"
                      autocomplete="off"
                      type="text"
                      placeholder="请输入验证码"
                      :class="inputClass('captcha')"
                      @blur="validateField('captcha')"
                      @input="handleTextInput('captcha')"
                    />
                  </div>

                  <button
                    type="button"
                    class="group relative h-14 overflow-hidden rounded-2xl border border-cyan-100 bg-white/90 transition hover:border-cyan-300 hover:bg-cyan-50"
                    @click="getCaptcha()"
                  >
                    <span
                      v-if="loadingCaptcha"
                      class="flex h-full items-center justify-center text-sm font-medium text-slate-500"
                    >
                      正在加载...
                    </span>
                    <img
                      v-else-if="captchaImage"
                      :src="captchaImage"
                      alt="注册验证码"
                      class="h-full w-full object-contain px-3"
                    />
                    <span
                      v-else
                      class="flex h-full items-center justify-center text-sm font-medium text-slate-500"
                    >
                      点击获取
                    </span>
                    <span class="absolute inset-x-0 bottom-1 text-[11px] font-medium text-slate-400 transition group-hover:text-cyan-600">
                      点击更换
                    </span>
                  </button>
                </div>

                <p v-if="fieldErrors.captcha" class="mt-2 text-xs font-medium text-rose-500">
                  {{ fieldErrors.captcha }}
                </p>
              </div>

              <div class="rounded-2xl border border-cyan-100 bg-white/90 px-4 py-3">
                <label class="inline-flex cursor-pointer items-start gap-3 text-sm leading-6 text-slate-600">
                  <input
                    v-model="registerForm.agreeToTerms"
                    type="checkbox"
                    class="peer sr-only"
                    @change="validateField('agreeToTerms')"
                  />
                  <span class="mt-1 flex h-5 w-5 shrink-0 items-center justify-center rounded border border-slate-300 bg-white shadow-sm transition peer-checked:border-cyan-500">
                    <svg
                      class="h-3.5 w-3.5 text-cyan-600 opacity-0 transition peer-checked:opacity-100"
                      viewBox="0 0 14 14"
                      fill="none"
                      xmlns="http://www.w3.org/2000/svg"
                    >
                      <path
                        d="M11.6666 3.5L5.24992 9.91667L2.33325 7"
                        stroke="currentColor"
                        stroke-width="1.7"
                        stroke-linecap="round"
                        stroke-linejoin="round"
                      />
                    </svg>
                  </span>
                  <span>
                    我已确认注册后将使用系统默认角色/部门/岗位，并遵守后台账号使用规范。
                  </span>
                </label>
                <p v-if="fieldErrors.agreeToTerms" class="mt-2 text-xs font-medium text-rose-500">
                  {{ fieldErrors.agreeToTerms }}
                </p>
              </div>

              <button
                type="submit"
                class="flex h-14 w-full items-center justify-center rounded-2xl bg-gradient-to-r from-cyan-600 via-sky-600 to-blue-600 text-sm font-semibold text-white shadow-[0_16px_28px_rgba(2,132,199,0.35)] transition hover:from-cyan-500 hover:via-sky-500 hover:to-blue-500 disabled:cursor-not-allowed disabled:from-slate-300 disabled:to-slate-300 disabled:shadow-none"
                :disabled="submitting"
              >
                {{ submitting ? '注册中...' : '创建账号' }}
              </button>
            </form>

            <div class="mt-5 flex flex-col gap-3 border-t border-slate-100 pt-5 sm:flex-row sm:items-center sm:justify-between">
              <p class="text-sm text-slate-500">已经有账号了？直接返回登录继续。</p>
              <router-link
                to="/adminLogin"
                class="inline-flex items-center justify-center rounded-2xl border border-cyan-100 px-4 py-2 text-sm font-medium text-cyan-700 transition hover:border-cyan-300 hover:text-cyan-800"
              >
                返回登录
              </router-link>
            </div>
          </div>
        </div>
      </section>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'
import {
  ArrowPathIcon,
  CheckCircleIcon,
  EyeIcon,
  EyeSlashIcon,
  KeyIcon,
  ShieldCheckIcon,
  SparklesIcon,
  UserCircleIcon,
  UserPlusIcon
} from '@heroicons/vue/24/outline'
import adminApi from '@/api/system'
import ToastAlert from '@/composables/ToastAlert'
import {
  CAPTCHA_PATTERN,
  PASSWORD_PATTERN,
  USERNAME_PATTERN,
  getRequestErrorMessage
} from '@/utils/authForms'

type RegisterField =
  | 'username'
  | 'password'
  | 'confirmPassword'
  | 'captcha'
  | 'agreeToTerms'

const router = useRouter()

const showPassword = ref(false)
const showConfirmPassword = ref(false)
const submitting = ref(false)
const loadingCaptcha = ref(false)
const captchaImage = ref('')

const showcaseItems = [
  {
    title: '统一入口',
    description: '管理员登录、会话保持和重定向逻辑集中在一个入口，排查问题更直接。'
  },
  {
    title: '验证码防护',
    description: '登录失败会立即刷新验证码，避免旧验证码反复重试影响判断。'
  },
  {
    title: '权限分发',
    description: '登录成功后同步拉取菜单与权限列表，进入后台后即可按角色工作。'
  }
]

const registerForm = reactive({
  username: '',
  password: '',
  confirmPassword: '',
  captcha: '',
  idKey: '',
  agreeToTerms: false
})

const fieldErrors = reactive<Record<RegisterField, string>>({
  username: '',
  password: '',
  confirmPassword: '',
  captcha: '',
  agreeToTerms: ''
})

const inputClass = (field: Exclude<RegisterField, 'agreeToTerms'>) => {
  const baseClass =
    'h-14 w-full rounded-2xl border bg-white pl-12 pr-4 text-sm text-slate-900 outline-none transition placeholder:text-slate-400 focus:ring-4'

  return fieldErrors[field]
    ? `${baseClass} border-rose-300 focus:border-rose-400 focus:ring-rose-100`
    : `${baseClass} border-slate-200 focus:border-cyan-400 focus:ring-cyan-100`
}

const passwordStrength = computed(() => {
  const value = registerForm.password
  let score = 0

  if (value.length >= 6) score += 1
  if (/[A-Za-z]/.test(value) && /\d/.test(value)) score += 1
  if (/[^\w]/.test(value)) score += 1
  if (value.length >= 10) score += 1

  if (score <= 1) {
    return {
      score,
      label: '建议至少 6 位，并同时包含字母和数字。',
      tag: '较弱',
      barClass: 'bg-rose-400',
      badgeClass: 'bg-rose-50 text-rose-600'
    }
  }

  if (score === 2) {
    return {
      score,
      label: '基础可用，继续增加长度或符号会更稳妥。',
      tag: '一般',
      barClass: 'bg-amber-400',
      badgeClass: 'bg-amber-50 text-amber-600'
    }
  }

  if (score === 3) {
    return {
      score,
      label: '已经比较可靠，适合作为后台账号密码。',
      tag: '良好',
      barClass: 'bg-cyan-500',
      badgeClass: 'bg-cyan-50 text-cyan-700'
    }
  }

  return {
    score,
    label: '强度较高，请妥善保管账号凭据。',
    tag: '很强',
    barClass: 'bg-emerald-500',
    badgeClass: 'bg-emerald-50 text-emerald-700'
  }
})

const validateField = (field: RegisterField) => {
  const username = registerForm.username.trim()
  const captcha = registerForm.captcha.trim()

  const validators: Record<RegisterField, () => string> = {
    username: () => {
      if (!username) return '请输入账号'
      if (!USERNAME_PATTERN.test(username)) return '账号需为 4-20 位字母、数字或下划线'
      return ''
    },
    password: () => {
      if (!registerForm.password) return '请输入密码'
      if (!PASSWORD_PATTERN.test(registerForm.password)) return '密码需至少 6 位，且包含字母和数字'
      return ''
    },
    confirmPassword: () => {
      if (!registerForm.confirmPassword) return '请再次输入密码'
      if (registerForm.confirmPassword !== registerForm.password) return '两次输入的密码不一致'
      return ''
    },
    captcha: () => {
      if (!captcha) return '请输入验证码'
      if (!CAPTCHA_PATTERN.test(captcha)) return '验证码需为 4-6 位字母或数字'
      return ''
    },
    agreeToTerms: () => {
      if (!registerForm.agreeToTerms) return '请先确认注册规范'
      return ''
    }
  }

  const errorMessage = validators[field]()
  fieldErrors[field] = errorMessage
  return !errorMessage
}

const validateForm = () => {
  const fields: RegisterField[] = [
    'username',
    'password',
    'confirmPassword',
    'captcha',
    'agreeToTerms'
  ]

  return fields.every((field) => validateField(field))
}

const handleTextInput = (field: Exclude<RegisterField, 'agreeToTerms'>) => {
  if (field === 'username' || field === 'captcha') {
    registerForm[field] = registerForm[field].replace(/\s+/g, '')
  }

  if (fieldErrors[field]) {
    validateField(field)
  }
}

const handlePasswordInput = () => {
  if (fieldErrors.password) {
    validateField('password')
  }
  if (registerForm.confirmPassword) {
    validateField('confirmPassword')
  }
}

const getCaptcha = async (clearValue = true) => {
  loadingCaptcha.value = true

  try {
    const { data: res } = await adminApi.captcha()

    if (res.code === 200) {
      captchaImage.value = res.data.image
      registerForm.idKey = res.data.idKey

      if (clearValue) {
        registerForm.captcha = ''
        fieldErrors.captcha = ''
      }
      return
    }

    ToastAlert.error({
      title: '获取验证码失败',
      message: res.message || '请稍后重试',
      blurBackground: false
    })
  } catch (error) {
    ToastAlert.error({
      title: '获取验证码失败',
      message: getRequestErrorMessage(error, '网络异常，请稍后重试'),
      blurBackground: false
    })
  } finally {
    loadingCaptcha.value = false
  }
}

const submitRegister = async () => {
  if (submitting.value) {
    return
  }

  if (!validateForm()) {
    ToastAlert.error({
      title: '表单校验未通过',
      message: '请先修正标红字段',
      blurBackground: false
    })
    return
  }

  submitting.value = true

  try {
    const { data: res } = await adminApi.register({
      username: registerForm.username.trim(),
      password: registerForm.password,
      confirmPassword: registerForm.confirmPassword,
      image: registerForm.captcha.trim(),
      idKey: registerForm.idKey
    })

    if (res.code !== 200) {
      throw new Error(res.message || '注册失败，请稍后重试')
    }

    ToastAlert.success({
      title: '注册成功',
      message: '即将返回登录页',
      duration: 1500,
      blurBackground: false
    })

    setTimeout(() => {
      router.push({
        path: '/adminLogin',
        query: {
          username: registerForm.username.trim()
        }
      })
    }, 900)
  } catch (error) {
    ToastAlert.error({
      title: '注册失败',
      message: getRequestErrorMessage(error, '注册失败，请稍后重试'),
      blurBackground: false
    })
    await getCaptcha()
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  getCaptcha(false)
})
</script>

<style scoped>
.water-gradient-layer {
  background:
    radial-gradient(circle at 20% 20%, rgba(56, 189, 248, 0.2), transparent 45%),
    radial-gradient(circle at 80% 18%, rgba(14, 165, 233, 0.18), transparent 40%),
    linear-gradient(120deg, rgba(240, 249, 255, 0.8), rgba(219, 234, 254, 0.86), rgba(224, 242, 254, 0.8));
  background-size: 200% 200%;
  animation: waterShift 20s ease-in-out infinite;
}

.water-orb {
  background: radial-gradient(circle at 30% 30%, rgba(186, 230, 253, 0.82), rgba(56, 189, 248, 0.2));
  filter: blur(8px);
  opacity: 0.72;
  animation: orbFloat 12s ease-in-out infinite;
}

.orb-one {
  animation-delay: -1s;
}

.orb-two {
  animation-duration: 16s;
  animation-delay: -5s;
}

.orb-three {
  animation-duration: 14s;
  animation-delay: -8s;
}

.water-wave {
  background: linear-gradient(90deg, rgba(191, 219, 254, 0.15), rgba(125, 211, 252, 0.26), rgba(191, 219, 254, 0.15));
  border-radius: 999px;
  filter: blur(2px);
  opacity: 0.45;
  animation: waveDrift 15s linear infinite;
}

.wave-two {
  opacity: 0.35;
  animation-duration: 22s;
  animation-direction: reverse;
}

.wave-three {
  opacity: 0.28;
  animation-duration: 26s;
}

.sine-line {
  overflow: visible;
  opacity: 0.58;
  animation: lineFlow 18s ease-in-out infinite;
}

.sine-path {
  fill: none;
  stroke: rgba(14, 165, 233, 0.42);
  stroke-width: 2;
  stroke-linecap: round;
  stroke-linejoin: round;
  stroke-dasharray: 14 18;
  animation: dashFlow 10s linear infinite, waveBreath 6.2s ease-in-out infinite;
}

.line-one {
  animation-delay: -2s;
}

.line-two {
  opacity: 0.44;
  animation-duration: 24s;
  animation-delay: -11s;
}

.line-two .sine-path {
  stroke: rgba(14, 165, 233, 0.34);
  stroke-width: 1.8;
  stroke-dasharray: 10 18;
  animation-duration: 12s, 7.4s;
}

.line-three {
  opacity: 0.36;
  animation-duration: 20s;
  animation-delay: -7s;
}

.line-three .sine-path {
  stroke: rgba(56, 189, 248, 0.28);
  stroke-width: 1.6;
  stroke-dasharray: 9 16;
  animation-duration: 8s, 5.6s;
}

.line-four {
  opacity: 0.3;
  animation-duration: 27s;
  animation-delay: -15s;
}

.line-four .sine-path {
  stroke: rgba(125, 211, 252, 0.25);
  stroke-width: 1.5;
  stroke-dasharray: 8 15;
  animation-duration: 14s, 8s;
}

.water-ripple {
  border: 1px solid rgba(14, 165, 233, 0.24);
  box-shadow: inset 0 0 0 2px rgba(255, 255, 255, 0.35);
  opacity: 0.45;
  animation: ripplePulse 8s ease-in-out infinite;
}

.ripple-two {
  animation-duration: 11s;
  animation-delay: -2s;
}

.ripple-three {
  animation-duration: 10s;
  animation-delay: -4s;
}

.bubble {
  position: absolute;
  bottom: -120px;
  border-radius: 999px;
  background: radial-gradient(circle at 35% 30%, rgba(255, 255, 255, 0.9), rgba(186, 230, 253, 0.45));
  box-shadow: inset 0 0 0 1px rgba(14, 165, 233, 0.14);
  opacity: 0.56;
  animation: bubbleRise linear infinite;
}

.bubble-1 {
  left: 7%;
  width: 14px;
  height: 14px;
  animation-duration: 14s;
  animation-delay: -5s;
}

.bubble-2 {
  left: 18%;
  width: 19px;
  height: 19px;
  animation-duration: 17s;
  animation-delay: -9s;
}

.bubble-3 {
  left: 32%;
  width: 12px;
  height: 12px;
  animation-duration: 13s;
  animation-delay: -2s;
}

.bubble-4 {
  left: 50%;
  width: 16px;
  height: 16px;
  animation-duration: 15s;
  animation-delay: -7s;
}

.bubble-5 {
  left: 63%;
  width: 11px;
  height: 11px;
  animation-duration: 12s;
  animation-delay: -10s;
}

.bubble-6 {
  left: 76%;
  width: 18px;
  height: 18px;
  animation-duration: 19s;
  animation-delay: -4s;
}

.bubble-7 {
  left: 88%;
  width: 14px;
  height: 14px;
  animation-duration: 16s;
  animation-delay: -12s;
}

.bubble-8 {
  left: 94%;
  width: 9px;
  height: 9px;
  animation-duration: 11s;
  animation-delay: -6s;
}

@keyframes waterShift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}

@keyframes orbFloat {
  0%,
  100% {
    transform: translate3d(0, 0, 0) scale(1);
  }
  50% {
    transform: translate3d(10px, -22px, 0) scale(1.08);
  }
}

@keyframes waveDrift {
  0% {
    transform: translateX(0) translateY(0);
  }
  50% {
    transform: translateX(-6%) translateY(8px);
  }
  100% {
    transform: translateX(0) translateY(0);
  }
}

@keyframes lineFlow {
  0% {
    transform: translate3d(0, 0, 0);
  }
  50% {
    transform: translate3d(-8%, 9px, 0) scaleX(1.015);
  }
  100% {
    transform: translate3d(0, 0, 0);
  }
}

@keyframes dashFlow {
  0% {
    stroke-dashoffset: 0;
  }
  100% {
    stroke-dashoffset: -360;
  }
}

@keyframes waveBreath {
  0%,
  100% {
    transform: scaleY(1);
    opacity: 0.9;
  }
  50% {
    transform: scaleY(1.3);
    opacity: 1;
  }
}

@keyframes ripplePulse {
  0%,
  100% {
    transform: scale(1);
    opacity: 0.38;
  }
  50% {
    transform: scale(1.08);
    opacity: 0.7;
  }
}

@keyframes bubbleRise {
  0% {
    transform: translate3d(0, 0, 0) scale(0.9);
    opacity: 0;
  }
  12% {
    opacity: 0.6;
  }
  100% {
    transform: translate3d(-26px, -120vh, 0) scale(1.15);
    opacity: 0;
  }
}

@media (prefers-reduced-motion: reduce) {
  .water-gradient-layer,
  .water-orb,
  .water-wave,
  .sine-line,
  .sine-path,
  .water-ripple,
  .bubble {
    animation: none;
  }
}
</style>
