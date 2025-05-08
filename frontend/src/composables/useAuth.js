import { ref } from 'vue';
import axios from 'axios';
import { useRouter } from 'vue-router';

// 處理所有認證相關的邏輯
export function useAuth() {
    const userRole = ref(null);         // 儲存當前使用者的角色
    const isLoggedIn = ref(false);      // 標記是否已登入
    const currentUser = ref('');        // 儲存當前使用者的身分證號/員工編號
    const token = ref('');              // 儲存認證 token
    const showLoginForm = ref(false);   // 控制登入表單顯示
    const loading = ref(false);         // 標記載入狀態
    const router = useRouter();         // 引入路由實例

    // 錯誤訊息通知事件
    const emitError = (message) => {
        document.dispatchEvent(new CustomEvent('show-snackbar', { detail: { message, color: 'error' } }));
    };

    const emitSuccess = (message) => {
        document.dispatchEvent(new CustomEvent('show-snackbar', { detail: { message, color: 'success' } }));
    };

    // 初始化認證狀態，從 localStorage 載入
    const initAuth = () => {
        const storedToken = localStorage.getItem('token');
        if (storedToken) {
            token.value = storedToken;
            userRole.value = localStorage.getItem('role');
            currentUser.value = localStorage.getItem('id_number');
            isLoggedIn.value = true;
            axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
        }
    };

    // 登入處理
    const login = async (data) => {
        loading.value = true;
        try {
            const response = await axios.post('/default/login', {
                username: data.username,
                password: data.password,
                scope: `role:${data.role}`
            }, {
                headers: { 'Content-Type': 'application/x-www-form-urlencoded' },
                transformRequest: [(data) => {
                    const params = new URLSearchParams();
                    for (const key in data) {
                        params.append(key, data[key]);
                    }
                    return params.toString();
                }]
            });
            token.value = response.data.access_token;
            localStorage.setItem('token', token.value);
            localStorage.setItem('role', data.role);
            localStorage.setItem('id_number', data.username);
            userRole.value = data.role;
            currentUser.value = data.username;
            isLoggedIn.value = true;
            showLoginForm.value = false;
            axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
            emitSuccess('登入成功！');
            router.push('/');
        } catch (error) {
            handleAuthError(error, 'login');
        } finally {
            loading.value = false;
        }
    };

    // 註冊處理
    const register = async (data) => {
        loading.value = true;
        try {
            const response = await axios.post('/default/register', {
                id_number: data.id_number,
                password: data.password,
                full_name: data.full_name,
                gender: data.gender,
                birth_date: data.birth_date,
                phone_number: data.phone_number,
                email: data.email,
                role: data.role
            }, {
                headers: { 'Content-Type': 'application/json' }
            });
            emitSuccess('註冊成功，請登入！');
            showLoginForm.value = true;
            router.push('/login');
        } catch (error) {
            handleAuthError(error, 'register');
        } finally {
            loading.value = false;
        }
    };

    // 登出處理
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

    // 統一處理認證錯誤
    const handleAuthError = (error, action) => {
        console.error(`${action} failed:`, error);
        const errorMessage = error.response?.data?.detail || error.message;
        const status = error.response?.status;

        const errorMap = {
            400: {
                '未提供角色': '請選擇角色',
                '身分證字號格式不正確': '身分證號/員工編號格式錯誤，請檢查輸入',
                default: '請求格式錯誤，請檢查輸入'
            },
            404: '找不到端點，請確認後端服務是否正常運行',
            422: '請求數據格式錯誤，請檢查輸入',
            default: {
                '密碼錯誤': '密碼錯誤，請重新輸入',
                '帳號不存在': '帳號不存在，請確認身分證號/員工編號',
                '角色不匹配': '角色不匹配，請確認選擇的角色',
                default: `${action === 'login' ? '登入' : '註冊'}失敗：${errorMessage}`
            }
        };

        const getMessage = (status, message) => {
            if (status && errorMap[status]) {
                return errorMap[status][message] || errorMap[status].default || errorMap.default.default;
            }
            return errorMap.default[message] || errorMap.default.default;
        };

        emitError(getMessage(status, errorMessage));
    };

    return {
        userRole,
        isLoggedIn,
        currentUser,
        token,
        showLoginForm,
        loading,
        initAuth,
        login,
        register,
        logout,
        handleAuthError
    };
}