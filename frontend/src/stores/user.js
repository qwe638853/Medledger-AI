import { defineStore } from 'pinia';
import { ref } from 'vue';
import { userService, healthCheckService } from '../services';
import { useAuthStore } from './auth';
import router from '../router';

export const useUserStore = defineStore('user', () => {
  // 狀態
  const healthData = ref([]);
  const currentReport = ref(null);
  const loading = ref(false);

  // 設置當前報告
  const setCurrentReport = (report) => {
    currentReport.value = report;
  };

  // 忘記密碼
  const forgotPassword = async (data) => {
    loading.value = true;
    try {
      await userService.forgotPassword(data);
      router.push('/');
    } catch (error) {
      console.error('忘記密碼失敗:', error);
    } finally {
      loading.value = false;
    }
  };

  // 獲取健康檢查數據
  const fetchHealthData = async () => {
    const authStore = useAuthStore();
    loading.value = true;
    
    try {
      const response = await healthCheckService.fetchOtherHealthData(authStore.currentUser);
      healthData.value = response || [];
      return response;
    } catch (error) {
      if (error.response?.status === 401) {
        authStore.logout();
        router.push('/login');
        return null;
      }
      console.error('獲取數據失敗:', error);
      return [];
    } finally {
      loading.value = false;
    }
  };

  // 獲取用戶個人資料
  const fetchUserProfile = async () => {
    const authStore = useAuthStore();
    loading.value = true;
    
    try {
      const response = await userService.getUserProfile(authStore.currentUser);
      return response;
    } catch (error) {
      console.error('獲取用戶資料失敗:', error);
      return null;
    } finally {
      loading.value = false;
    }
  };

  // 更新用戶資料
  const updateProfile = async (profileData) => {
    const authStore = useAuthStore();
    loading.value = true;
    
    try {
      const response = await userService.updateUserProfile(authStore.currentUser, profileData);
      return response;
    } catch (error) {
      console.error('更新用戶資料失敗:', error);
      return null;
    } finally {
      loading.value = false;
    }
  };

  return {
    healthData,
    currentReport,
    loading,
    setCurrentReport,
    forgotPassword,
    fetchHealthData,
    fetchUserProfile,
    updateProfile
  };
}); 