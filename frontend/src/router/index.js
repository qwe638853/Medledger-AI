import { createRouter, createWebHistory } from 'vue-router';
import { useAuthStore } from '../stores/auth';

// 定義路由
const routes = [
  {
    path: '/',
    name: 'Home',
    component: () => import('../views/Home.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/login',
    name: 'Login',
    component: () => import('../components/LoginForm.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/register',
    name: 'Register',
    component: () => import('../components/RegisterForm.vue'),
    meta: { requiresAuth: false }
  },
  {
    path: '/user-dashboard',
    name: 'UserDashboard',
    component: () => import('../views/UserDashboard.vue'),
    meta: { requiresAuth: true, roles: ['user'] }
  },
  {
    path: '/hospital-dashboard',
    name: 'HospitalDashboard',
    component: () => import('../views/HospitalDashboard.vue'),
    meta: { requiresAuth: true, roles: ['medical'] }
  },
  {
    path: '/other-user-dashboard',
    name: 'OtherUserDashboard',
    component: () => import('../views/OtherUserDashboard.vue'),
    meta: { requiresAuth: true, roles: ['other'] }
  },
  {
    path: '/:pathMatch(.*)*',
    redirect: '/'
  }
];

// 創建路由實例
const router = createRouter({
  history: createWebHistory(),
  routes
});

// 檢查本地儲存的授權狀態
function checkLocalAuthState() {
  const token = localStorage.getItem('token');
  const role = localStorage.getItem('role');
  const id_number = localStorage.getItem('id_number');
  
  return {
    isAuthenticated: !!token,
    userRole: role,
    userId: id_number
  };
}

// 導航守衛
router.beforeEach((to, from, next) => {
  console.group('路由導航守衛');
  console.log('從:', from.path);
  console.log('到:', to.path);
  console.log('路由元數據:', to.meta);
  
  // 從 Pinia store 獲取狀態
  const authStore = useAuthStore();
  console.log('Store 狀態:', { isLoggedIn: authStore.isLoggedIn, userRole: authStore.userRole });
  
  // 直接從本地儲存檢查狀態
  const localAuthState = checkLocalAuthState();
  console.log('本地儲存狀態:', localAuthState);
  
  // 組合兩種狀態來決定用戶是否已登入
  const isAuthenticated = authStore.isLoggedIn || localAuthState.isAuthenticated;
  const activeRole = authStore.userRole || localAuthState.userRole;
  
  console.log('最終認證狀態:', { isAuthenticated, activeRole });
  
  if (to.meta.requiresAuth && !isAuthenticated) {
    console.log('需要登入但用戶未登入，重定向到登入頁面');
    next('/login');
  } else if (to.meta.requiresAuth && to.meta.roles && !to.meta.roles.includes(activeRole)) {
    console.log('用戶角色不匹配，重定向到首頁或登入頁面');
    next(isAuthenticated ? '/' : '/login');
  } else {
    console.log('允許導航');
    next();
  }
  
  console.groupEnd();
});

// 導出路由實例
export default router; 