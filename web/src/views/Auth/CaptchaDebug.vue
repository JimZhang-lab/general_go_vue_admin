<template>
  <div class="p-6 space-y-6">
    <div class="bg-white rounded-lg shadow-lg p-6 dark:bg-gray-800">
      <h1 class="text-2xl font-bold text-gray-900 dark:text-white mb-6">
        验证码调试工具
      </h1>
      
      <!-- 验证码获取测试 -->
      <div class="mb-8 p-4 border border-gray-200 rounded-lg dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
          验证码获取测试
        </h3>
        
        <div class="flex items-center gap-4 mb-4">
          <button
            @click="getCaptcha"
            class="px-4 py-2 bg-blue-500 text-white rounded-lg hover:bg-blue-600 transition-colors"
          >
            获取验证码
          </button>
          
          <div v-if="captchaData.image" class="flex items-center gap-2">
            <img 
              :src="captchaData.image" 
              class="h-12 w-32 object-contain rounded border border-gray-300"
              alt="验证码"
            />
            <div class="text-sm text-gray-600 dark:text-gray-300">
              <div>ID: {{ captchaData.idKey }}</div>
              <div>时间: {{ captchaData.timestamp }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 验证码验证测试 -->
      <div class="mb-8 p-4 border border-gray-200 rounded-lg dark:border-gray-700">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
          验证码验证测试
        </h3>
        
        <div class="flex items-center gap-4 mb-4">
          <input
            v-model="testCaptcha"
            type="text"
            placeholder="输入验证码"
            class="px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-blue-500 dark:bg-gray-700 dark:border-gray-600 dark:text-white"
          />
          
          <button
            @click="testLogin"
            :disabled="!captchaData.idKey || !testCaptcha"
            class="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 disabled:bg-gray-400 disabled:cursor-not-allowed transition-colors"
          >
            测试登录
          </button>
        </div>
        
        <div class="text-sm text-gray-600 dark:text-gray-300">
          <p>使用测试账号: admin / admin123</p>
          <p>请先获取验证码，然后输入验证码进行测试</p>
        </div>
      </div>

      <!-- 测试结果显示 -->
      <div class="bg-gray-50 rounded-lg p-4 dark:bg-gray-900">
        <h3 class="text-lg font-semibold text-gray-900 dark:text-white mb-3">
          测试结果
        </h3>
        <div class="space-y-2 max-h-96 overflow-y-auto">
          <div v-for="(result, index) in testResults" :key="index" class="flex items-start gap-2">
            <span :class="result.success ? 'text-green-600' : 'text-red-600'">
              {{ result.success ? '✅' : '❌' }}
            </span>
            <div class="flex-1">
              <div class="text-gray-700 dark:text-gray-300">{{ result.message }}</div>
              <div v-if="result.details" class="text-xs text-gray-500 mt-1">
                {{ result.details }}
              </div>
              <div class="text-xs text-gray-500">{{ result.timestamp }}</div>
            </div>
          </div>
        </div>
      </div>

      <!-- 调试信息 -->
      <div class="mt-8 p-4 bg-yellow-50 rounded-lg dark:bg-yellow-900/20">
        <h3 class="text-lg font-semibold text-yellow-900 dark:text-yellow-300 mb-2">
          调试信息
        </h3>
        <div class="text-sm text-yellow-800 dark:text-yellow-200 space-y-1">
          <div><strong>后端地址:</strong> http://localhost:8080</div>
          <div><strong>验证码接口:</strong> GET /api/captcha</div>
          <div><strong>登录接口:</strong> POST /api/login</div>
          <div><strong>Redis配置:</strong> 127.0.0.1:6379</div>
          <div><strong>验证码有效期:</strong> 5分钟</div>
        </div>
      </div>

      <!-- 常见问题排查 -->
      <div class="mt-6 p-4 bg-blue-50 rounded-lg dark:bg-blue-900/20">
        <h3 class="text-lg font-semibold text-blue-900 dark:text-blue-300 mb-2">
          常见问题排查
        </h3>
        <ul class="text-sm text-blue-800 dark:text-blue-200 space-y-1">
          <li>• <strong>401错误</strong>: 检查验证码是否正确输入</li>
          <li>• <strong>验证码过期</strong>: 重新获取验证码</li>
          <li>• <strong>Redis连接</strong>: 确保Redis服务正常运行</li>
          <li>• <strong>大小写敏感</strong>: 验证码不区分大小写</li>
          <li>• <strong>网络问题</strong>: 检查前后端服务是否正常</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import ToastAlert from '@/composables/ToastAlert'

// 验证码数据
const captchaData = ref({
  idKey: '',
  image: '',
  timestamp: ''
})

// 测试验证码输入
const testCaptcha = ref('')

// 测试结果记录
const testResults = ref<Array<{
  message: string
  success: boolean
  timestamp: string
  details?: string
}>>([])

// 添加测试结果
const addTestResult = (message: string, success: boolean, details?: string) => {
  testResults.value.unshift({
    message,
    success,
    timestamp: new Date().toLocaleTimeString(),
    details
  })
}

// 获取验证码
const getCaptcha = async () => {
  try {
    addTestResult('开始获取验证码...', true)
    
    const response = await fetch('http://localhost:8080/api/captcha')
    const result = await response.json()
    
    console.log('验证码响应:', result)
    
    if (result.code === 200) {
      captchaData.value = {
        idKey: result.data.idKey,
        image: result.data.image,
        timestamp: new Date().toLocaleTimeString()
      }
      
      addTestResult('验证码获取成功', true, `ID: ${result.data.idKey}`)
    } else {
      addTestResult('验证码获取失败', false, `错误: ${result.message}`)
    }
  } catch (error) {
    console.error('获取验证码失败:', error)
    addTestResult('验证码获取异常', false, `网络错误: ${error}`)
  }
}

// 测试登录
const testLogin = async () => {
  if (!captchaData.value.idKey || !testCaptcha.value) {
    addTestResult('测试失败', false, '请先获取验证码并输入验证码')
    return
  }

  try {
    addTestResult('开始测试登录...', true, `验证码: ${testCaptcha.value}`)
    
    const loginData = {
      adminUsername: 'admin',
      password: 'admin123',
      idKey: captchaData.value.idKey,
      image: testCaptcha.value
    }
    
    console.log('登录请求数据:', loginData)
    
    const response = await fetch('http://localhost:8080/api/login', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify(loginData)
    })
    
    const result = await response.json()
    console.log('登录响应:', result)
    
    if (response.ok && result.code === 200) {
      addTestResult('登录测试成功', true, '验证码验证通过')
      ToastAlert.success({
        title: '测试成功',
        message: '验证码验证正常'
      })
    } else {
      addTestResult('登录测试失败', false, `错误: ${result.message || '未知错误'}`)
      
      if (result.message?.includes('验证码')) {
        ToastAlert.error({
          title: '验证码错误',
          message: result.message
        })
      }
    }
  } catch (error) {
    console.error('登录测试失败:', error)
    addTestResult('登录测试异常', false, `网络错误: ${error}`)
  }
}

// 页面加载时自动获取验证码
getCaptcha()
</script>
