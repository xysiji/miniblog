<template>
    <div class="max-w-md mx-auto bg-gray-50 min-h-screen">
      <header class="bg-white p-4 sticky top-0 z-10 flex items-center space-x-4 border-b">
        <button @click="$router.back()" class="text-gray-600">← 返回</button>
        <h1 class="text-lg font-bold">个人主页</h1>
      </header>
  
      <div v-if="userInfo" class="bg-white p-6 mb-2">
        <div class="flex items-center space-x-4 mb-4">
          <img :src="userInfo.avatar || 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'" class="w-20 h-20 rounded-full border">
          <div>
            <h2 class="text-xl font-bold">{{ userInfo.username || '未知用户' }}</h2>
            <p class="text-sm text-gray-500 mt-1">ID: {{ userInfo.userId }}</p>
          </div>
        </div>
        <p class="text-sm text-gray-700 mb-4">{{ userInfo.bio || '这个人很懒，什么都没写~' }}</p>
        <div class="flex space-x-6 text-sm">
          <div class="text-center"><span class="font-bold">{{ userInfo.followingCount || 0 }}</span> <span class="text-gray-500">关注</span></div>
          <div class="text-center"><span class="font-bold">{{ userInfo.followerCount || 0 }}</span> <span class="text-gray-500">粉丝</span></div>
        </div>
      </div>
      <div v-else-if="loading" class="p-10 text-center text-gray-400">加载中...</div>
      <div v-else class="p-10 text-center text-gray-400">用户不存在</div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue'
  import { useRoute } from 'vue-router'
  import api from '../utils/api'
  
  const route = useRoute()
  const userInfo = ref(null)
  const loading = ref(true)
  
  const fetchProfile = async () => {
    const targetId = parseInt(route.params.id)
    try {
      const res = await api.post('/v1/user/profile', { targetUserId: targetId })
      userInfo.value = res.data
    } catch (err) {
      console.error('获取主页失败', err)
    } finally {
      loading.value = false
    }
  }
  
  onMounted(() => {
    fetchProfile()
  })
  </script>