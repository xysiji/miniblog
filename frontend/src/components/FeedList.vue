<template>
  <div class="space-y-4 pb-10">
    <div class="bg-white p-4 rounded-xl shadow-sm border border-gray-100">
      <textarea class="w-full bg-gray-50 rounded-lg p-3 outline-none resize-none text-sm focus:ring-1 ring-gray-200" rows="3" placeholder="有什么新鲜事想分享给大家？"></textarea>
      <div class="flex justify-between items-center mt-3">
        <div class="flex space-x-4 text-gray-500">
          <button class="hover:text-blue-500 flex items-center space-x-1"><span class="text-lg">📷</span><span class="text-sm">图片</span></button>
          <button class="hover:text-blue-500 flex items-center space-x-1"><span class="text-lg">😃</span><span class="text-sm">表情</span></button>
        </div>
        <button class="bg-red-500 text-white px-5 py-1.5 rounded-full text-sm font-medium hover:bg-red-600 opacity-80 cursor-not-allowed">发布 (待完善)</button>
      </div>
    </div>

    <div v-for="post in postList" :key="post.id" class="bg-white p-5 rounded-xl shadow-sm border border-gray-100">
      
      <div class="flex items-center justify-between mb-3 cursor-pointer" @click="goToProfile(post.userId)">
        <div class="flex items-center space-x-3">
          <img :src="post.authorAvatar || defaultAvatar" class="w-12 h-12 rounded-full object-cover border border-gray-100">
          <div>
            <p class="font-bold text-gray-800">{{ post.authorName || '匿名用户' }}</p>
            <p class="text-xs text-gray-400 mt-0.5">{{ post.createAt }} 来自 Web</p>
          </div>
        </div>
        <button class="text-sm text-red-500 border border-red-500 bg-red-50 px-4 py-1 rounded-full hover:bg-red-500 hover:text-white transition">
          + 关注
        </button>
      </div>

      <div class="mb-4 text-[15px] text-gray-800 leading-relaxed whitespace-pre-wrap pl-[60px]">{{ post.content }}</div>
      
      <div v-if="post.images && post.images.length > 0" class="pl-[60px] grid grid-cols-3 gap-2 mb-4 w-4/5">
        <img v-for="(img, idx) in post.images" :key="idx" :src="img" class="aspect-square object-cover rounded-lg border">
      </div>

      <div class="flex justify-around items-center text-gray-500 text-sm border-t border-gray-100 pt-3 mt-2 pl-[60px]">
        <button class="flex-1 flex justify-center items-center space-x-2 hover:text-gray-700 transition">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.684 13.342C8.886 12.938 9 12.482 9 12c0-.482-.114-.938-.316-1.342m0 2.684a3 3 0 110-2.684m0 2.684l6.632 3.316m-6.632-6l6.632-3.316m0 0a3 3 0 105.367-2.684 3 3 0 00-5.367 2.684zm0 9.316a3 3 0 105.368 2.684 3 3 0 00-5.368-2.684z"></path></svg>
          <span>转发</span>
        </button>
        
        <button @click="toggleComments(post)" class="flex-1 flex justify-center items-center space-x-2 transition" :class="post.showComments ? 'text-red-500' : 'hover:text-gray-700'">
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path></svg>
          <span>{{ post.comment_count > 0 ? post.comment_count : '评论' }}</span>
        </button>
        
        <button @click="toggleLike(post)" class="flex-1 flex justify-center items-center space-x-2 transition" :class="post.is_liked ? 'text-red-500' : 'hover:text-red-500'">
          <svg v-if="post.is_liked" class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20"><path fill-rule="evenodd" d="M3.172 5.172a4 4 0 015.656 0L10 6.343l1.172-1.171a4 4 0 115.656 5.656L10 17.657l-6.828-6.829a4 4 0 010-5.656z" clip-rule="evenodd"></path></svg>
          <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"></path></svg>
          <span>{{ post.like_count > 0 ? post.like_count : '赞' }}</span>
        </button>
      </div>

      <div v-if="post.showComments" class="mt-4 pl-[60px]">
        <div class="bg-gray-50 rounded-xl p-4">
          <div class="flex space-x-3 mb-4">
            <img :src="defaultAvatar" class="w-8 h-8 rounded-full">
            <div class="flex-1 flex">
              <input type="text" placeholder="发布你的评论..." class="flex-1 bg-white border border-gray-200 rounded-l-lg px-3 py-1.5 text-sm outline-none focus:border-red-300">
              <button class="bg-red-500 text-white px-4 text-sm font-medium rounded-r-lg hover:bg-red-600">评论</button>
            </div>
          </div>
          
          <div v-if="post.commentsLoading" class="text-center text-xs text-gray-400 py-2">努力加载中...</div>
          <div v-else-if="post.comments && post.comments.length > 0" class="space-y-4">
            <div v-for="comment in post.comments" :key="comment.id" class="flex space-x-3">
              <img :src="comment.avatar || defaultAvatar" class="w-8 h-8 rounded-full object-cover">
              <div class="flex-1 text-sm border-b border-gray-100 pb-3">
                <span class="font-medium text-gray-800">{{ comment.username }}</span>
                <span v-if="comment.reply_to_user_id" class="mx-1 text-gray-400">回复</span>
                <span v-if="comment.reply_to_user_id" class="font-medium text-blue-500">@{{ comment.reply_to_name }}</span>
                <p class="text-gray-700 mt-1 leading-relaxed">{{ comment.content }}</p>
                <div class="mt-1 text-xs text-gray-400 flex space-x-4">
                  <span>{{ formatTime(comment.create_time) }}</span>
                  <button class="hover:text-gray-600">回复</button>
                </div>
              </div>
            </div>
          </div>
          <div v-else class="text-center text-xs text-gray-400 py-4">还没有人评论，快来抢沙发~</div>
        </div>
      </div>

    </div>

    <div v-if="loading" class="text-center py-4 text-gray-400 text-sm">正在加载更多...</div>
    <div v-if="noMore" class="text-center py-6 text-gray-400 text-sm border-t">- 没有更多动态了 -</div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import api from '../utils/api'
import { globalState } from '../store' // 引入全局状态

const defaultAvatar = 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png'
const postList = ref([])
const page = ref(1)
const pageSize = ref(10)
const loading = ref(false)
const noMore = ref(false)
const router = useRouter()

// 格式化时间戳 (后端返回的是秒级时间戳)
const formatTime = (ts) => {
  if (!ts) return ''
  const date = new Date(ts * 1000)
  return `${date.getMonth()+1}-${date.getDate()} ${date.getHours()}:${date.getMinutes()}`
}

const fetchPosts = async () => {
  if (loading.value || noMore.value) return
  loading.value = true
  try {
    const res = await api.post('/v1/post/list', { page: page.value, pageSize: pageSize.value })
    const list = res.data?.list || []
    if (list.length > 0) {
      const mappedList = list.map(item => ({
        ...item,
        showComments: false,
        commentsLoading: false,
        comments: []
      }))
      postList.value.push(...mappedList)
      page.value += 1
      if (list.length < pageSize.value) noMore.value = true
    } else {
      noMore.value = true
    }
  } catch (err) {
    console.error('获取列表失败:', err)
  } finally {
    loading.value = false
  }
}

const toggleComments = async (post) => {
  post.showComments = !post.showComments
  
  if (post.showComments && post.comments.length === 0) {
    post.commentsLoading = true
    try {
      // 严格对齐后端字段
      const res = await api.post('/v1/interaction/comment/list', {
        post_id: post.id,
        page: 1,
        page_size: 20
      })
      post.comments = res.data?.list || []
    } catch (err) {
      console.error('加载评论失败', err)
    } finally {
      post.commentsLoading = false
    }
  }
}

const toggleLike = async (post) => {
  // 优雅的登录拦截：未登录直接弹出 Modal，不跳转打断体验
  if (!globalState.isLoggedIn) {
    globalState.showLoginModal = true
    return
  }
  
  const originalLiked = post.is_liked
  const originalCount = post.like_count
  post.is_liked = !post.is_liked
  post.like_count += post.is_liked ? 1 : -1
  
  try {
    if (post.is_liked) await api.post('/v1/interaction/like', { post_id: post.id })
    else await api.post('/v1/interaction/unlike', { post_id: post.id })
  } catch (err) {
    post.is_liked = originalLiked
    post.like_count = originalCount
  }
}

const goToProfile = (userId) => {
  if (userId) router.push(`/profile/${userId}`)
}

const handleScroll = () => {
  if (window.innerHeight + window.scrollY >= document.documentElement.offsetHeight - 200) {
    fetchPosts()
  }
}

onMounted(() => {
  fetchPosts()
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>