import { ref } from 'vue';
import axios from 'axios';
//處理所有認證相關的邏輯
//包含登入、登出、token 管理等功能
export function useAuth() {
    const userRole = ref(null);
    const isLoggedIn = ref(false);
    const currentUser = ref('');
    const token = ref('');
    const showLoginForm = ref(false);

    // 初始化認證狀態
    const initAuth = () => {
        if (localStorage.getItem('token')) {
            token.value = localStorage.getItem('token');
            userRole.value = localStorage.getItem('role');
            currentUser.value = localStorage.getItem('id_number');
            isLoggedIn.value = true;
            axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
        }
    };

    // 登入處理
    const login = async (data) => {
        try {
            const response = await axios.post('/default/login', {
                username: data.username,
                password: data.password,
                scope: `role:${data.role}`
            }, {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded'
                },
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
        } catch (error) {
            handleLoginError(error);
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
    };

    // 錯誤處理
    const handleLoginError = (error) => {
        console.error('Login failed:', error);
        const errorMessage = error.response?.data?.detail || error.message;
        if (error.response?.status === 400) {
            if (errorMessage.includes("未提供角色")) {
                alert('請選擇角色');
            } else if (errorMessage.includes("身分證字號格式不正確")) {
                alert('身分證號/員工編號格式錯誤，請檢查輸入');
            } else {
                alert('請求格式錯誤，請檢查輸入的身分證號/員工編號、密碼和角色');
            }
        } else if (error.response?.status === 404) {
            alert('找不到登入端點，請確認後端服務是否正常運行');
        } else if (error.response?.status === 422) {
            alert('請求數據格式錯誤，請檢查輸入的身分證號/員工編號和密碼');
        } else if (errorMessage === '密碼錯誤') {
            alert('密碼錯誤，請重新輸入');
        } else if (errorMessage === '帳號不存在') {
            alert('帳號不存在，請確認身分證號/員工編號是否正確');
        } else if (errorMessage === '角色不匹配') {
            alert('角色不匹配，請確認選擇的角色是否正確');
        } else {
            alert(`登入失敗：${errorMessage}`);
        }
    };

    return {
        userRole,
        isLoggedIn,
        currentUser,
        token,
        showLoginForm,
        initAuth,
        login,
        logout
    };
} 