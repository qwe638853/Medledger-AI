import { defineStore } from 'pinia';
import { ref } from 'vue';
import { authService, setupAuthInterceptor, notifySuccess } from '../services';
import router from '../router';

export const useAuthStore = defineStore('auth', () => {
  // 狀態
  const userRole = ref(null);
  const isLoggedIn = ref(false);
  const currentUser = ref('');
  const token = ref('');
  const loading = ref(false);

  // 初始化驗證
  const initAuth = () => {
    const storedToken = localStorage.getItem('token');
    if (storedToken) {
      token.value = storedToken;
      userRole.value = localStorage.getItem('role');
      currentUser.value = localStorage.getItem('id_number');
      isLoggedIn.value = true;
      setupAuthInterceptor(token.value);
      redirectToDashboard();
    }
  };

  // 登入
  const login = async (data) => {
    loading.value = true;
    try {
      console.group('登入處理流程');
      
      const response = await authService.login(data);
      
      if (response.success) {
        // 存儲令牌和用戶信息
        token.value = response.token;
        
        // 角色映射處理
        let mappedRole = data.role;
        if (typeof data.role === 'string') {
          const roleMappings = {
            '一般用戶': 'user',
            '醫療機構': 'medical',
            '其他用戶': 'other'
          };
          
          if (roleMappings[data.role]) {
            mappedRole = roleMappings[data.role];
          }
        }
        
        console.log('角色映射:', { original: data.role, mapped: mappedRole });
        
        // 存儲到本地儲存和應用狀態
        localStorage.setItem('token', token.value);
        localStorage.setItem('role', mappedRole);
        localStorage.setItem('id_number', data.username);
        
        userRole.value = mappedRole;
        currentUser.value = data.username;
        isLoggedIn.value = true;
        
        setupAuthInterceptor(token.value);
        notifySuccess('登入成功！');
        
        // 重定向到儀表板
        setTimeout(() => {
          redirectToDashboard();
        }, 100);
      } else {
        throw new Error(response.message || '登入失敗');
      }
    } catch (error) {
      throw error; // 重新拋出錯誤
    } finally {
      console.groupEnd();
      loading.value = false;
    }
  };

  // 註冊
  const register = async (data) => {
    loading.value = true;
    try {
      const response = await authService.register(data);
      
      if (response.success) {
        notifySuccess('註冊成功！請登入您的帳號。');
        return true;
      } else {
        throw new Error(response.message || '註冊失敗');
      }
    } catch (error) {
      return false;
    } finally {
      loading.value = false;
    }
  };

  // 跳轉到儀表板
  const redirectToDashboard = () => {
    console.group('重定向到儀表板');
    console.log('當前用戶角色:', userRole.value);
    
    try {
      if (userRole.value === 'user') {
        router.push('/user-dashboard');
      } else if (userRole.value === 'medical') {
        router.push('/hospital-dashboard');
      } else if (userRole.value === 'other') {
        router.push('/other-user-dashboard');
      } else {
        router.push('/');
      }
    } catch (error) {
      console.error('路由導航錯誤:', error);
    } finally {
      console.groupEnd();
    }
  };

  // 登出
  const logout = () => {
    userRole.value = null;
    currentUser.value = '';
    token.value = '';
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    localStorage.removeItem('id_number');
    isLoggedIn.value = false;
    notifySuccess('登出成功！');
    router.push('/');
  };

  return {
    userRole,
    isLoggedIn,
    currentUser,
    token,
    loading,
    initAuth,
    login,
    register,
    logout,
    redirectToDashboard
  };
}); 