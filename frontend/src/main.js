import { createApp } from 'vue';
import { createPinia } from 'pinia';
import { createRouter, createWebHistory } from 'vue-router';
import App from './App.vue';
import 'vuetify/styles';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import '@mdi/font/css/materialdesignicons.css';
import './style.css';
import { useAuth } from './composables/useAuth';

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light'
  }
});

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: '/',
      name: 'Home',
      component: () => import('./views/Home.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/login',
      name: 'Login',
      component: () => import('./components/LoginForm.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/:pathMatch(.*)*', // 404 路由
      redirect: '/'
    }
  ]
});

router.beforeEach((to, from, next) => {
  const { isLoggedIn } = useAuth();
  if (to.meta.requiresAuth && !isLoggedIn.value) {
    next('/login');
  } else {
    next();
  }
});

const app = createApp(App);
app.use(createPinia());
app.use(router);
app.use(vuetify);

// 移除 alert，改為靜默記錄並觸發 snackbar
app.config.errorHandler = (err, vm, info) => {
  console.error('Vue error:', err, info);
  document.dispatchEvent(new CustomEvent('show-snackbar', {
    detail: { message: '應用程式發生錯誤，請稍後再試或聯繫管理員。', color: 'error' }
  }));
};

app.mount('#app');