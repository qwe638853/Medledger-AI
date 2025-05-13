import { ref } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

export function useAuth() {
  const userRole = ref(null);
  const isLoggedIn = ref(false);
  const currentUser = ref('');
  const token = ref('');
  const loading = ref(false);
  const router = useRouter();

  const emitSuccess = (message) => {
    document.dispatchEvent(new CustomEvent('show-snackbar', { detail: { message, color: 'success' } }));
  };

  const emitError = (message) => {
    document.dispatchEvent(new CustomEvent('show-snackbar', { detail: { message, color: 'error' } }));
  };

  const initAuth = () => {
    const storedToken = localStorage.getItem('token');
    if (storedToken) {
      token.value = storedToken;
      userRole.value = localStorage.getItem('role');
      currentUser.value = localStorage.getItem('id_number');
      isLoggedIn.value = true;
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
      redirectToDashboard();
    }
  };

  const login = async (data) => {
    loading.value = true;
    try {
      const response = await axios.post('https://7aa9-140-124-249-9.ngrok-free.app/default/login', {
        username: data.username,
        password: data.password,
        scope: `role:${data.role}`
      }, {
        headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
        transformRequest: [(data) => new URLSearchParams(data).toString()],
        timeout: 10000
      });
      token.value = response.data.access_token;
      localStorage.setItem('token', token.value);
      localStorage.setItem('role', data.role);
      localStorage.setItem('id_number', data.username);
      userRole.value = data.role;
      currentUser.value = data.username;
      isLoggedIn.value = true;
      axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
      emitSuccess('登入成功！');
      redirectToDashboard();
    } catch (error) {
      emitError(`登入失敗：${error.response?.data?.detail || error.message}`);
    } finally {
      loading.value = false;
    }
  };

  const redirectToDashboard = () => {
    if (userRole.value === 'user') {
      router.push('/user-dashboard');
    } else if (userRole.value === 'medical') {
      router.push('/hospital-dashboard');
    } else if (userRole.value === 'other') {
      router.push('/other-user-dashboard');
    } else {
      router.push('/');
    }
  };

  const logout = () => {
    userRole.value = null;
    currentUser.value = '';
    token.value = '';
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    localStorage.removeItem('id_number');
    isLoggedIn.value = false;
    delete axios.defaults.headers.common['Authorization'];
    emitSuccess('登出成功！');
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
    logout
  };
}