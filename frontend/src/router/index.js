import { createRouter, createWebHistory } from 'vue-router'
import FeedList from '../components/FeedList.vue'
import Login from '../views/Login.vue'
import Profile from '../views/Profile.vue'

const routes = [
  { path: '/', name: 'Home', component: FeedList },
  { path: '/login', name: 'Login', component: Login },
  { path: '/profile/:id', name: 'Profile', component: Profile }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router