<script setup>
import { ref } from 'vue';
import LoginForm from './components/LoginForm.vue';
import UserDashboard from './UserDashboard.vue';
import HospitalDashboard from './HospitalDashboard.vue';
import axios from 'axios';

// 設置 axios 基礎 URL（使用 ngrok URL）
axios.defaults.baseURL = 'https://2b67-140-124-249-9.ngrok-free.app';

const userRole = ref(null);
const isLoggedIn = ref(false);
const currentUser = ref('');
const token = ref('');

// 頁面載入時檢查 Token
if (localStorage.getItem('token')) {
    token.value = localStorage.getItem('token');
    userRole.value = localStorage.getItem('role');
    currentUser.value = localStorage.getItem('id_number');
    isLoggedIn.value = true;
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
}

// 登入處理函數
const login = async (data) => {
    try {
        const response = await axios.post('/default/login', {
            id_number: data.username,
            password: data.password,
            role: data.role
        });
        token.value = response.data.token;
        localStorage.setItem('token', token.value); // 儲存 Token
        localStorage.setItem('role', response.data.role); // 儲存角色
        localStorage.setItem('id_number', response.data.id_number); // 儲存 id_number
        userRole.value = response.data.role;
        currentUser.value = response.data.id_number;
        isLoggedIn.value = true;
        axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
    } catch (error) {
        console.error('Login failed:', error.response?.data || error.message);
        const errorMessage = error.response?.data?.detail || error.message;
        if (error.response?.status === 400) {
            if (errorMessage.includes("無效的角色")) {
                alert('角色無效，請選擇有效的角色（hospital、patient 或 other）');
            } else if (errorMessage.includes("身分證號格式錯誤")) {
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
    }
};

// 忘記密碼處理函數
const forgotPassword = async (data) => {
    try {
        const response = await axios.post('/forget-password', {
            id_number: data.id_number,
            role: data.role
        });
        alert(response.data.message || '已發送重設密碼郵件，請檢查您的電子郵件');
    } catch (error) {
        const errorMessage = error.response?.data?.detail || error.message;
        alert(`重設密碼失敗：${errorMessage}`);
    }
};

// 註冊處理函數
const register = async (data) => {
    try {
        const response = await axios.post('/register', {
            id_number: data.id_number,
            password: data.password,
            role: data.role
        });
        alert(response.data.message || '註冊成功，請登入');
    } catch (error) {
        const errorMessage = error.response?.data?.detail || error.message;
        alert(`註冊失敗：${errorMessage}`);
    }
};

// 登出處理函數
const logout = () => {
    userRole.value = null;
    currentUser.value = '';
    token.value = '';
    localStorage.removeItem('token'); // 移除 Token
    localStorage.removeItem('role'); // 移除角色
    localStorage.removeItem('id_number'); // 移除 id_number
    isLoggedIn.value = false;
    delete axios.defaults.headers.common['Authorization'];
};

// 獲取數據的函數（僅用於 other 角色）
const fetchData = async () => {
    try {
        const response = await axios.get(`/health-check/other/${currentUser.value}`, {
            headers: {
                'Authorization': `Bearer ${token.value}`,
                'Accept': 'application/json'
            }
        });
        return response.data || [];
    } catch (error) {
        if (error.response?.status === 401) {
            alert('認證過期，請重新登入');
            logout();
        } else if (error.response?.status === 404) {
            alert('找不到數據端點，請確認後端服務是否正常運行');
        } else {
            console.error('Fetch data failed:', error);
            alert('獲取數據失敗，請重試');
        }
        return [];
    }
};
</script>

<template>
    <div class="app-container">
        <!-- 封面圖片區塊 -->
        <div class="cover-container">
            <div class="cover-overlay">
                <h1>健康檢查數據平台</h1>
                <p>安全管理您的健康數據</p>
            </div>
        </div>

        <!-- 主要內容區域 -->
        <div class="content-container">
            <LoginForm 
                v-if="!isLoggedIn" 
                @login="login" 
                @forgot-password="forgotPassword" 
                @register="register" 
            />
            <UserDashboard 
                v-if="isLoggedIn && userRole === 'patient'" 
                :username="currentUser" 
                :data="fetchData"
                @logout="logout"
            />
            <HospitalDashboard 
                v-if="isLoggedIn && userRole === 'hospital'" 
                :username="currentUser" 
                :data="fetchData"
                @logout="logout"
            />
            <UserDashboard 
                v-if="isLoggedIn && userRole === 'other'" 
                :username="currentUser" 
                :data="fetchData"
                @logout="logout"
            />
        </div>
    </div>
</template>

<style scoped>
.app-container { min-height: 100vh; }
.cover-container { position: relative; background: url('https://img.freepik.com/free-vector/hand-drawn-patient-taking-medical-examination_23-2148843031.jpg?t=st=1741877699~exp=1741881299~hmac=5e0611af12a15c1f641536eedfe1f76b2e23b513816c05f6b1e6b78a46f7a1f4&w=826') no-repeat center center; background-size: cover; height: 450px; width: 100%; }
.cover-overlay { position: absolute; top: 0; left: 0; width: 100%; height: 100%; background: rgba(0, 123, 255, 0.6); display: flex; flex-direction: column; justify-content: center; align-items: center; color: var(--white); text-align: center; padding: 30px; }
.cover-overlay h1 { font-size: 48px; font-weight: 700; margin-bottom: 15px; text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3); }
.cover-overlay p { font-size: 22px; font-weight: 300; text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3); }
.content-container { max-width: 1200px; margin: 0 auto; padding: 40px 20px; }
</style>