<script setup>
import { ref } from 'vue';
import LoginForm from './components/LoginForm.vue';
import UserDashboard from './UserDashboard.vue';
import HospitalDashboard from './HospitalDashboard.vue';
import axios from 'axios';

// 設置 axios 基礎 URL（目前為空）
axios.defaults.baseURL = 'http://localhost:8080';

// 創建用戶角色的響應式變量
const userRole = ref(null);
// 創建登入狀態的響應式變量
const isLoggedIn = ref(false);
// 創建當前用戶的響應式變量
const currentUser = ref('');
// 創建 token 的響應式變量
const token = ref('');
// 創建是否顯示登入表單的響應式變量
const showLoginForm = ref(false);
// 創建是否顯示頁腳的響應式變量
const showFooter = ref(false);

// 定義導航菜單項目的響應式數組
const menuItems = ref([
  { title: '首頁', path: '/' },
  { title: '健康紀錄', path: '/records' },
  { title: '聯絡我們', path: '/contact' }
]);

// 在頁面載入時檢查本地存儲中是否有 token
if (localStorage.getItem('token')) {
    // 如果有 token，則設置相關的響應式變量
    token.value = localStorage.getItem('token');
    userRole.value = localStorage.getItem('role');
    currentUser.value = localStorage.getItem('id_number');
    isLoggedIn.value = true;
    // 設置 axios 的默認請求頭，加入 token
    axios.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
}

// 定義登入處理函數
const login = async (data) => {
    try {
        // 發送登入請求到後端
        const response = await axios.post('/v1/login', {
            username: data.username,
            password: data.password,
        }, {
            headers: {
                'Content-Type': 'application/json'
            },
            // 轉換請求數據格式
            transformRequest: [(data) => {
                const params = new URLSearchParams();
                for (const key in data) {
                    params.append(key, data[key]);
                }
                return params.toString();
            }]
        });
        // 保存登入成功後的數據
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
        // 錯誤處理
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
    }
};

// 定義忘記密碼處理函數
const forgotPassword = async (data) => {
    try {
        // 發送忘記密碼請求
        const response = await axios.post('/default/forget-password', {
            id_number: data.id_number,
            role: data.role
        });
        alert(response.data.message || '已發送重設密碼郵件，請檢查您的電子郵件');
    } catch (error) {
        const errorMessage = error.response?.data?.detail || error.message;
        alert(`重設密碼失敗：${errorMessage}`);
    }
};

// 定義註冊處理函數
const register = async (data) => {
    try {
        // 發送註冊請求
        const response = await axios.post('/v1/register', {
            id_number: data.id_number,
            password: data.password,
        }, {
            headers: {
                'Content-Type': 'application/json'
            }
        });
        alert(response.data || '註冊成功，請登入');
    } catch (error) {
        const errorMessage = error.response?.data?.detail || error.message;
        alert(`註冊失敗：${errorMessage}`);
    }
};

// 定義登出處理函數
const logout = () => {
    // 清除所有用戶相關數據
    userRole.value = null;
    currentUser.value = '';
    token.value = '';
    localStorage.removeItem('token');
    localStorage.removeItem('role');
    localStorage.removeItem('id_number');
    isLoggedIn.value = false;
    delete axios.defaults.headers.common['Authorization'];
};

// 定義獲取數據的函數
const fetchData = async () => {
    try {
        // 發送獲取健康檢查數據的請求
        const response = await axios.get(`/default/health-check/other/${currentUser.value}`, {
            headers: {
                'Authorization': `Bearer ${token.value}`,
                'Accept': 'application/json'
            }
        });
        return response.data || [];
    } catch (error) {
        // 錯誤處理
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
  <v-app>
    <v-app-bar app color="white" elevation="1">
      <v-container class="d-flex align-center">
        <router-link to="/" class="text-decoration-none">
          <div class="d-flex align-center">
            <v-img src="@/assets/logo.svg" width="40" class="mr-2" />
            <span class="text-h5 font-weight-bold primary--text">健康檢查平台</span>
          </div>
        </router-link>

        <v-spacer />

        <template v-if="isLoggedIn">
          <nav>
            <v-btn
              v-for="item in menuItems"
              :key="item.title"
              :to="item.path"
              text
              class="mx-2"
            >
              {{ item.title }}
            </v-btn>
          </nav>
          <v-btn
            color="error"
            class="ml-4"
            rounded
            elevation="0"
            @click="logout"
          >
            登出
          </v-btn>
        </template>
        <template v-else>
          <v-btn
            color="primary"
            class="ml-4"
            rounded
            elevation="0"
            @click="showLoginForm = true"
          >
            登入
          </v-btn>
        </template>
      </v-container>
    </v-app-bar>

    <v-main>
      <template v-if="!isLoggedIn">
        <LoginForm 
          v-if="showLoginForm"
          @login="login" 
          @forgot-password="forgotPassword" 
          @register="register"
        />
        <router-view v-else />
      </template>
      <template v-else>
        <UserDashboard 
          v-if="userRole === 'user'" 
          :username="currentUser" 
          :data="fetchData"
        />
        <HospitalDashboard 
          v-if="userRole === 'health_center'" 
          :username="currentUser" 
          :data="fetchData"
        />
        <UserDashboard 
          v-if="userRole === 'other'" 
          :username="currentUser" 
          :data="fetchData"
        />
      </template>
    </v-main>

    <v-footer app color="primary" dark class="py-2">
      <v-container>
        <v-row no-gutters align="center">
          <v-col cols="12" md="6">
            <div class="d-flex align-center">
              <v-btn icon @click="showFooter = !showFooter">
                <v-icon>{{ showFooter ? 'mdi-chevron-up' : 'mdi-chevron-down' }}</v-icon>
              </v-btn>
              <span class="text-subtitle-2 ml-2">健康檢查平台</span>
            </div>
          </v-col>
          <v-col cols="12" md="6" class="text-right">
            <v-btn
              v-if="!isLoggedIn"
              color="white"
              class="text-primary"
              small
              @click="showLoginForm = true"
            >
              立即註冊
            </v-btn>
          </v-col>
        </v-row>
        <v-expand-transition>
          <div v-show="showFooter">
            <v-row class="mt-2">
              <v-col cols="12" md="6">
                <p class="text-body-2 mb-0">
                  提供安全、便捷的健康檢查數據管理平台，讓您的健康數據一目了然。
                </p>
              </v-col>
              <v-col cols="12" md="6" class="text-right">
                <p class="text-body-2 mb-0">
                  Email: service@health-platform.com<br>
                  電話: (02) 1234-5678
                </p>
              </v-col>
            </v-row>
          </div>
        </v-expand-transition>
      </v-container>
    </v-footer>
  </v-app>
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.v-application {
  font-family: 'Noto Sans TC', sans-serif !important;
}

.v-footer {
  transition: all 0.3s ease;
}
</style>