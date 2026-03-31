// frontend/src/store.js
import { reactive } from 'vue'

export const globalState = reactive({
  isLoggedIn: !!localStorage.getItem('token'),
  userId: localStorage.getItem('userId'),
  showLoginModal: false, // 控制全局登录弹窗的显示

  login(token, userId) {
    localStorage.setItem('token', token)
    localStorage.setItem('userId', userId)
    this.isLoggedIn = true
    this.userId = userId
    this.showLoginModal = false
  },
  
  logout() {
    localStorage.removeItem('token')
    localStorage.removeItem('userId')
    this.isLoggedIn = false
    this.userId = null
  }
})