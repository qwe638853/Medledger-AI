// 引入 Vue 的核心功能
import { createApp } from 'vue'
// 引入 Pinia 狀態管理庫
import { createPinia } from 'pinia'
// 引入 Vue Router 相關功能
import { createRouter, createWebHistory } from 'vue-router'
// 引入根組件
import App from './App.vue'
// 引入 Vuetify 的基礎樣式
import 'vuetify/styles'
// 引入 Vuetify 的創建函數
import { createVuetify } from 'vuetify'
// 引入 Vuetify 的所有組件
import * as components from 'vuetify/components'
// 引入 Vuetify 的所有指令
import * as directives from 'vuetify/directives'
// 引入 Material Design Icons 字體
import '@mdi/font/css/materialdesignicons.css'
// 引入自定義樣式
import './style.css'

// 創建 Vuetify 實例，配置主題和組件
const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light' // 設置默認主題為亮色
  }
})

// 創建路由實例，配置路由規則
const router = createRouter({
  history: createWebHistory(), // 使用 HTML5 History 模式
  routes: [
    {
      path: '/', // 首頁路由
      name: 'Home',
      component: () => import('./views/Home.vue') // 懶加載首頁組件
    },
    {
      path: '/blog', // 部落格頁面路由
      name: 'Blog',
      component: () => import('./views/Blog.vue') // 懶加載部落格組件
    },
    {
      path: '/contact', // 聯絡我們頁面路由
      name: 'Contact',
      component: () => import('./views/Contact.vue') // 懶加載聯絡我們組件
    }
  ]
})

// 創建 Vue 應用實例
const app = createApp(App)
// 使用 Pinia 進行狀態管理
app.use(createPinia())
// 使用路由
app.use(router)
// 使用 Vuetify
app.use(vuetify)
// 掛載應用到 DOM
app.mount('#app')



