import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import tailwindcss from '@tailwindcss/vite' // 引入 v4 专属插件

export default defineConfig({
  plugins: [
    vue(),
    tailwindcss(), // 激活插件
  ],
})