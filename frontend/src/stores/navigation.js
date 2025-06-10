import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import { useAuthStore } from './auth';

export const useNavigationStore = defineStore('navigation', () => {
  // 狀態
  const drawer = ref(false);
  const authStore = useAuthStore();
  
  // 計算屬性
  const navItems = computed(() => {
    const items = [
      { title: '首頁', icon: 'mdi-home', path: '/' }
    ];
    
    if (authStore.isLoggedIn) {
      // 根據用戶角色添加導航項目
      if (authStore.userRole === 'user') {
        items.push(
          { title: '我的儀表板', icon: 'mdi-view-dashboard', path: '/user-dashboard' },
          { title: '上傳健檢報告', icon: 'mdi-file-upload', path: '/upload-report' },
          { title: '個人資料', icon: 'mdi-account', path: '/profile' }
        );
      } else if (authStore.userRole === 'medical') {
        items.push(
          { title: '醫療機構儀表板', icon: 'mdi-view-dashboard', path: '/hospital-dashboard' },
          { title: '查詢患者', icon: 'mdi-account-search', path: '/patient-search' }
        );
      } else if (authStore.userRole === 'other') {
        items.push(
          { title: '其他用戶儀表板', icon: 'mdi-view-dashboard', path: '/other-user-dashboard' }
        );
      }
      
      // 添加登出選項
      items.push({ title: '登出', icon: 'mdi-logout', path: '/logout', action: authStore.logout });
    } else {
      // 未登入用戶的選項
      items.push(
        { title: '登入', icon: 'mdi-login', path: '/login' },
        { title: '註冊', icon: 'mdi-account-plus', path: '/register' }
      );
    }
    
    return items;
  });
  
  // 切換導航抽屜
  const toggleDrawer = () => {
    drawer.value = !drawer.value;
  };

  return {
    drawer,
    navItems,
    toggleDrawer
  };
}); 