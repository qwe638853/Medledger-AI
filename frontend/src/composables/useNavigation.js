import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuth } from './useAuth';

export function useNavigation() {
  const router = useRouter();
  const { userRole, isLoggedIn } = useAuth();
  const showFooter = ref(false);
  const menuItems = ref([
    { title: '首頁', path: '/' },
    { title: '註冊', path: '/register' },
    { title: '登入', path: '/login' }
  ]);

  const filteredMenuItems = computed(() => {
    return isLoggedIn.value
      ? menuItems.value.filter(item => item.path === '/')
      : menuItems.value;
  });

  const goToHome = (showLoginForm) => {
    if (showLoginForm && typeof showLoginForm.value !== 'undefined') {
      showLoginForm.value = false;
    }
    router.push('/');
  };

  return {
    showFooter,
    menuItems: filteredMenuItems,
    goToHome
  };
}