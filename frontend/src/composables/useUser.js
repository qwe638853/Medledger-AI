import axios from 'axios';
import { useAuth } from './useAuth';
import { useRouter } from 'vue-router';

// 處理用戶相關的操作
export function useUser() {
    const { token, currentUser, handleAuthError, logout } = useAuth();
    const router = useRouter();

    // 錯誤訊息通知事件
    const emitSuccess = (message) => {
        document.dispatchEvent(new CustomEvent('show-snackbar', { detail: { message, color: 'success' } }));
    };

    // 忘記密碼
    const forgotPassword = async (data) => {
        try {
            const response = await axios.post('/default/forget-password', {
                id_number: data.id_number,
                role: data.role
            });
            emitSuccess(response.data.message || '已發送重設密碼郵件，請檢查您的電子郵件');
            router.push('/');
        } catch (error) {
            handleAuthError(error, 'forgot-password');
        }
    };

    // 獲取健康檢查數據
    const fetchData = async () => {
        try {
            const response = await axios.get(`/default/health-check/other/${currentUser.value}`, {
                headers: {
                    'Authorization': `Bearer ${token.value}`,
                    'Accept': 'application/json'
                }
            });
            return response.data || [];
        } catch (error) {
            if (error.response?.status === 401) {
                logout();
                document.dispatchEvent(new CustomEvent('show-snackbar', { detail: { message: '認證過期，請重新登入', color: 'error' } }));
                router.push('/login');
                return null;
            }
            handleAuthError(error, 'fetch-data');
            return [];
        }
    };

    return {
        forgotPassword,
        fetchData
    };
}