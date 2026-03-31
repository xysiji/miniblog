import axios from 'axios'

// 不再设置写死的 baseURL
const api = axios.create({
  timeout: 5000
})

// 请求拦截器：根据路径自动转发到正确的微服务端口，并携带 Token
api.interceptors.request.use(config => {
  const token = localStorage.getItem('token')
  if (token) {
    config.headers.Authorization = `Bearer ${token}`
  }

  // 动态路由转发逻辑
  const url = config.url
  if (url.startsWith('/v1/user/')) {
    // 用户相关接口转发到 user-api (8888端口)
    config.baseURL = 'http://127.0.0.1:8888'
  } else if (url.startsWith('/v1/interaction/')) {
    // 互动相关接口转发到 interaction-api (8890端口)
    config.baseURL = 'http://127.0.0.1:8890'
  } else if (url.startsWith('/v1/post/')) {
    // 博文相关接口转发到 post-api (8889端口)
    config.baseURL = 'http://127.0.0.1:8889'
  }
  
  return config
})

export default api