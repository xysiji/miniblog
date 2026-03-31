<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-50 px-4">
    <div class="max-w-md w-full bg-white p-8 rounded-2xl shadow-sm border border-gray-100 relative">
      <button @click="$router.push('/')" class="absolute top-4 left-4 text-gray-400 hover:text-gray-600">
        ← 返回
      </button>

      <h2 class="text-2xl font-bold text-center text-gray-800 mb-8 mt-4">
        {{ isLoginMode ? '欢迎回来' : '注册新账号' }}
      </h2>
      
      <form @submit.prevent="handleSubmit" class="space-y-6">
        <div>
          <label class="block text-sm font-medium text-gray-700">用户名</label>
          <input v-model="form.username" type="text" required class="mt-1 block w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:ring-red-500 focus:border-red-500 outline-none transition-colors">
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700">密码</label>
          <input v-model="form.password" type="password" required class="mt-1 block w-full px-4 py-2 bg-gray-50 border border-gray-200 rounded-lg focus:ring-red-500 focus:border-red-500 outline-none transition-colors">
        </div>
        
        <button type="submit" :disabled="loading" class="w-full flex justify-center py-2.5 px-4 border border-transparent rounded-lg shadow-sm text-sm font-medium text-white bg-red-500 hover:bg-red-600 focus:outline-none transition-colors">
          {{ loading ? '处理中...' : (isLoginMode ? '登 录' : '注 册') }}
        </button>
      </form>
      
      <div class="mt-6 text-center text-sm text-gray-600">
        {{ isLoginMode ? '还没有账号？' : '已有账号？' }}
        <button @click="toggleMode" class="text-red-500 hover:text-red-600 font-medium transition-colors">
          {{ isLoginMode ? '立即注册' : '去登录' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'

const router = useRouter()
const form = ref({ username: '', password: '' })
const loading = ref(false)
const isLoginMode = ref(true) // 控制显示登录还是注册

const toggleMode = () => {
  isLoginMode.value = !isLoginMode.value
  form.value = { username: '', password: '' } // 切换时清空表单
}

const handleSubmit = async () => {
  loading.value = true
  try {
    if (isLoginMode.value) {
      // 执行登录逻辑
      const res = await api.post('/v1/user/login', form.value)
      if (res.data && res.data.token) {
        localStorage.setItem('token', res.data.token)
        localStorage.setItem('userId', res.data.userId)
        alert('登录成功！')
        router.push('/')
      }
    } else {
      // 执行注册逻辑
      const res = await api.post('/v1/user/register', form.value)
      if (res.data && res.data.userId) {
        alert('注册成功，请使用新账号登录！')
        isLoginMode.value = true // 注册成功后自动切回登录页面
      }
    }
  } catch (err) {
    alert(isLoginMode.value ? '登录失败，请检查账号密码' : '注册失败，用户名可能已存在')
    console.error(err)
  } finally {
    loading.value = false
  }
}
</script>