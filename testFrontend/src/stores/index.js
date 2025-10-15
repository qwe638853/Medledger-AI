import { useAuthStore } from './auth';
import { useUserStore } from './user';
import { useNavigationStore } from './navigation';

export {
  useAuthStore,
  useUserStore,
  useNavigationStore
};

// 初始化程序，例如在 main.js 中調用
export const initializeApp = () => {
  const authStore = useAuthStore();
  authStore.initAuth();
}; 