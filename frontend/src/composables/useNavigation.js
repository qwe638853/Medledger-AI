import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuth } from './useAuth';

export function useNavigation() {
  const router = useRouter();
  const { userRole } = useAuth();
  const showFooter = ref(false);
  const menuItems = ref([
    { title: '首頁', path: '/' },
    { title: '登入', path: '/login' }
  ]);

  const goToHome = (showLoginForm) => {
    if (showLoginForm && typeof showLoginForm.value !== 'undefined') {
      showLoginForm.value = false;
    }
    router.push('/');
  };

  return {
    showFooter,
    menuItems,
    goToHome
  };
}