<template>
  <div class="min-h-screen bg-gray-100 font-sans text-gray-800">
    <header class="bg-white shadow-sm fixed w-full top-0 z-50 h-14">
      <div class="max-w-5xl mx-auto flex justify-between items-center h-full px-4">
        <div class="flex items-center space-x-2 text-red-500 font-bold text-2xl cursor-pointer" @click="$router.push('/')">
          <svg class="w-8 h-8" fill="currentColor" viewBox="0 0 24 24"><path d="M22.23 6.13c-.8.36-1.66.6-2.56.71.92-.55 1.63-1.42 1.96-2.46-.86.51-1.81.88-2.82 1.08-.81-.86-1.97-1.4-3.24-1.4-2.44 0-4.42 1.98-4.42 4.42 0 .35.04.68.11 1.01-3.67-.18-6.93-1.94-9.11-4.61-.38.65-.6 1.41-.6 2.22 0 1.53.78 2.88 1.97 3.68-.73-.02-1.42-.22-2.02-.56v.05c0 2.16 1.54 3.96 3.58 4.37-.37.1-.76.15-1.16.15-.28 0-.56-.03-.83-.08.57 1.77 2.21 3.06 4.16 3.09-1.53 1.2-3.46 1.91-5.55 1.91-.36 0-.71-.02-1.06-.06 1.97 1.26 4.31 2 6.82 2 8.18 0 12.65-6.78 12.65-12.65 0-.19 0-.38-.01-.57.87-.63 1.63-1.41 2.23-2.31z"/></svg>
          <span>MiniBlog</span>
        </div>

        <nav class="flex space-x-8 h-full">
          <router-link to="/" class="nav-item">首页</router-link>
          <a @click.prevent="requireAuth('/following')" class="nav-item cursor-pointer">关注</a>
          <a @click.prevent="requireAuth('/messages')" class="nav-item cursor-pointer">消息</a>
          <a @click.prevent="requireAuth(`/profile/${globalState.userId}`)" class="nav-item cursor-pointer">我的</a>
        </nav>

        <div class="flex items-center space-x-4">
          <div class="relative">
            <input type="text" placeholder="搜索博文/用户" class="bg-gray-100 rounded-full py-1.5 px-4 text-sm outline-none focus:bg-white focus:ring-1 ring-gray-300 w-48 transition-all">
          </div>
          <button v-if="!globalState.isLoggedIn" @click="globalState.showLoginModal = true" class="bg-red-500 hover:bg-red-600 text-white px-4 py-1.5 rounded-full text-sm font-medium transition">
            登录 / 注册
          </button>
          <button v-else @click="handleLogout" class="text-gray-500 hover:text-red-500 text-sm font-medium">
            退出
          </button>
        </div>
      </div>
    </header>

    <main class="pt-20 pb-10 max-w-5xl mx-auto flex gap-6 px-4">
      <div class="w-8/12">
        <router-view />
      </div>

      <aside class="w-4/12 hidden md:block space-y-4">
        <div class="bg-white rounded-xl p-4 shadow-sm" v-if="globalState.isLoggedIn">
          <div class="flex items-center space-x-3">
            <div class="w-12 h-12 bg-gray-200 rounded-full"></div>
            <div>
              <p class="font-bold">你好！用户 {{ globalState.userId }}</p>
              <p class="text-xs text-gray-500">快来分享今天的新鲜事吧~</p>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-xl p-4 shadow-sm">
          <h3 class="font-bold border-b pb-2 mb-3 text-sm">热门话题</h3>
          <ul class="space-y-3 text-sm text-gray-600">
            <li class="hover:text-red-500 cursor-pointer"># 毕业答辩倒计时 #</li>
            <li class="hover:text-red-500 cursor-pointer"># Go-Zero微服务实战 #</li>
            <li class="hover:text-red-500 cursor-pointer"># Vue3前端开发 #</li>
          </ul>
        </div>
      </aside>
    </main>

    <div v-if="globalState.showLoginModal" class="fixed inset-0 z-[100] flex items-center justify-center bg-black bg-opacity-50">
      <div class="bg-white w-96 rounded-2xl shadow-xl overflow-hidden">
        <div class="px-6 py-4 border-b flex justify-between items-center bg-gray-50">
          <h3 class="font-bold text-lg">账号登录</h3>
          <button @click="globalState.showLoginModal = false" class="text-gray-400 hover:text-gray-600 text-xl">&times;</button>
        </div>
        <div class="p-6">
          <div class="space-y-4">
            <input v-model="loginForm.username" type="text" placeholder="用户名" class="w-full px-4 py-2 border rounded-lg focus:ring-2 ring-red-200 outline-none">
            <input v-model="loginForm.password" type="password" placeholder="密码" class="w-full px-4 py-2 border rounded-lg focus:ring-2 ring-red-200 outline-none">
            <button @click="doLogin" class="w-full bg-red-500 text-white py-2 rounded-lg hover:bg-red-600 font-medium transition">登 录</button>
            <p class="text-xs text-center text-gray-400">目前仅演示登录，无账号请先通过后端接口注册</p>
          </div>
        </div>
      </div>
    </div>

  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { globalState } from './store'
import api from './utils/api'

const router = useRouter()
const loginForm = ref({ username: '', password: '' })

// 路由拦截守护：未登录触发弹窗
const requireAuth = (path) => {
  if (!globalState.isLoggedIn) {
    globalState.showLoginModal = true
  } else {
    router.push(path)
  }
}

// 执行登录
const doLogin = async () => {
  try {
    const res = await api.post('/v1/user/login', loginForm.value)
    if (res.data && res.data.token) {
      globalState.login(res.data.token, res.data.userId) // 更新全局状态
      alert('登录成功！')
    }
  } catch (err) {
    alert('登录失败，请检查账号或服务状态')
  }
}

// 退出登录
const handleLogout = () => {
  if(confirm('确定要退出登录吗？')) {
    globalState.logout()
    router.push('/')
  }
}
</script>

<style scoped>
/* 引入 Tailwind v4 的引用，让 @apply 能认识工具类 */
@reference "tailwindcss";

.nav-item {
  @apply flex items-center h-full px-1 border-b-2 border-transparent text-gray-600 hover:text-red-500 font-medium transition-colors;
}
.router-link-active {
  @apply border-red-500 text-red-500 font-bold;
}
</style>