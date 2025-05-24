<template>
  <div class="login-page">
    <v-container class="fill-height">
      <v-row align="center" justify="center">
        <v-col cols="12" sm="8" md="6" lg="4">
          <!-- 主要登入卡片 -->
          <v-card class="login-card" elevation="0">
            <!-- 頂部標題區 -->
            <div class="header-section">
              <h1 class="header-title">企業健康數據管理平台</h1>
              <p class="header-subtitle">安全、便捷的醫療數據管理解決方案</p>
            </div>
            
            <!-- 系統訊息提示 -->
            <v-alert
              v-if="alertInfo.show"
              :type="alertInfo.type"
              :title="alertInfo.title"
              density="comfortable"
              variant="tonal"
              class="alert-message"
              closable
              @click:close="alertInfo.show = false"
            >
              {{ alertInfo.message }}
            </v-alert>
            
            <!-- 登入表單 -->
            <v-form
              @submit.prevent="handleSubmit"
              ref="form"
              v-model="valid"
              lazy-validation
              class="login-form"
            >
              <!-- 角色選擇 -->
              <v-select
                v-model="selectedRole"
                :items="roles"
                item-title="text"
                item-value="value"
                label="選擇角色"
                prepend-inner-icon="mdi-account-outline"
                variant="outlined"
                :rules="[rules.required]"
                class="form-field"
                bg-color="white"
                density="comfortable"
              />
              
              <!-- 帳號輸入 -->
              <v-text-field
                v-model="username"
                :label="usernameLabel"
                prepend-inner-icon="mdi-identifier"
                variant="outlined"
                :rules="[rules.required]"
                class="form-field"
                bg-color="white"
                density="comfortable"
                :placeholder="usernamePlaceholder"
                clearable
              />
              
              <!-- 密碼輸入 -->
              <v-text-field
                v-model="password"
                label="密碼"
                prepend-inner-icon="mdi-lock-outline"
                :type="showPassword ? 'text' : 'password'"
                :append-inner-icon="showPassword ? 'mdi-eye-off-outline' : 'mdi-eye-outline'"
                @click:append-inner="showPassword = !showPassword"
                variant="outlined"
                :rules="[rules.required, rules.minLength]"
                class="form-field"
                bg-color="white"
                density="comfortable"
                clearable
              />
              
              <!-- 主要登入按鈕 -->
              <v-btn
                :loading="loading"
                class="login-btn"
                block
                type="submit"
                :disabled="!valid"
                elevation="0"
                height="48"
                prepend-icon="mdi-login-variant"
              >
                登入系統
              </v-btn>
              
              <!-- 測試登入按鈕 -->
              <v-btn
                class="test-btn"
                block
                @click="handleTestLogin"
                elevation="0"
                height="48"
                prepend-icon="mdi-test-tube-outline"
              >
                測試登入
              </v-btn>
              
              <!-- 導航按鈕組 -->
              <div class="nav-buttons">
                <v-btn
                  class="nav-btn"
                  elevation="0"
                  @click="goToHome"
                  prepend-icon="mdi-home-outline"
                >
                  返回首頁
                </v-btn>
                <v-btn
                  class="nav-btn"
                  elevation="0"
                  @click="goToRegister"
                  prepend-icon="mdi-account-plus-outline"
                >
                  註冊帳號
                </v-btn>
              </div>
              
              <!-- 開發調試資訊 -->
              <v-expand-transition>
                <div v-if="isDevelopment" class="debug-section">
                  <v-divider class="divider"></v-divider>
                  <v-expansion-panels>
                    <v-expansion-panel class="debug-panel">
                      <v-expansion-panel-title class="debug-title">
                        開發調試資訊
                      </v-expansion-panel-title>
                      <v-expansion-panel-text>
                        <pre class="debug-info">{{ JSON.stringify(debugInfo, null, 2) }}</pre>
                      </v-expansion-panel-text>
                    </v-expansion-panel>
                  </v-expansion-panels>
                </div>
              </v-expand-transition>
            </v-form>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores';

const router = useRouter();
const authStore = useAuthStore();

// 表單數據
const username = ref('');
const password = ref('');
const showPassword = ref(false);
const selectedRole = ref('');
const loading = computed(() => authStore.loading);
const roles = [
  { text: '一般用戶', value: 'user' },
  { text: '醫療機構', value: 'medical' },
  { text: '保險業者', value: 'insurer' }
];
const form = ref(null);
const valid = ref(false);

// 根據角色動態顯示不同的輸入框標籤和提示文字
const usernameLabel = computed(() => {
  if (selectedRole.value === 'insurer') {
    return '保險業者ID';
  } else if (selectedRole.value === 'medical') {
    return '醫療機構ID';
  } else {
    return '身分證號碼';
  }
});

const usernamePlaceholder = computed(() => {
  if (selectedRole.value === 'insurer') {
    return '請輸入您的保險業者ID';
  } else if (selectedRole.value === 'medical') {
    return '請輸入您的醫療機構ID';
  } else {
    return '請輸入您的身分證號碼';
  }
});

// 判斷是否為開發環境
const isDevelopment = import.meta.env.MODE === 'development';

// 用於調試的數據
const debugInfo = computed(() => {
  return {
    formState: {
    username: username.value,
    password: password.value ? '******' : '',
    role: selectedRole.value,
      roleText: roles.find(r => r.value === selectedRole.value)?.text || '',
    },
    validation: {
    valid: valid.value,
      loading: loading.value,
      errors: form.value?.errors || []
    },
    authStore: {
      isLoggedIn: authStore.isLoggedIn,
      currentUser: authStore.currentUser,
      userRole: authStore.userRole
    }
  };
});

// Alert 相關數據
const alertInfo = ref({
  show: false,
  type: 'error',
  title: '',
  message: '',
  icon: '',
});

// 顯示 Alert
const showAlert = (type, message, title = '') => {
  alertInfo.value = {
    show: true,
    type,
    title,
    message,
    icon: type === 'success' ? 'mdi-check-circle' : 'mdi-alert-circle',
  };
  
  // 5秒後自動關閉
  setTimeout(() => {
    alertInfo.value.show = false;
  }, 5000);
};

// 驗證規則
const rules = {
  required: value => !!value || '此欄位為必填',
  minLength: value => (value && value.length >= 6) || '密碼至少需要6個字符'
};

// 表單提交
const handleSubmit = async () => {
  const isValid = form.value?.validate();
  
  if (!isValid) {
    showAlert('error', '請確認表單欄位已正確填寫', '驗證錯誤');
    return;
  }
  
  try {
    console.group('登入表單提交');
    console.log('表單數據:', {
      username: username.value,
      role: selectedRole.value,
      // 不記錄密碼，但顯示長度
      passwordLength: password.value ? password.value.length : 0
    });
    
    // 檢查角色是否已選擇
    if (!selectedRole.value) {
      showAlert('error', '請選擇用戶角色', '錯誤');
      console.error('未選擇角色');
      return;
    }
    
    // 打印當前路由信息
    console.log('當前路由:', router.currentRoute.value.path);
    
    await authStore.login({
      username: username.value,
      password: password.value,
      role: selectedRole.value
    });
    
    // 登入成功顯示
    showAlert('success', '登入成功！正在為您導向...', '成功');
    console.log('登入成功！');
    
    // 檢查是否需要手動重定向（如果 store 中的重定向失敗）
    setTimeout(() => {
      // 檢查是否仍在登入頁面
      if (router.currentRoute.value.path === '/login') {
        console.log('檢測到仍在登入頁面，嘗試手動重定向');
        authStore.redirectToDashboard();
      }
    }, 1000);
  } catch (error) {
    console.error('登入處理錯誤:', error);
    
    // 根據錯誤類型顯示適當訊息
    let errorMsg = '登入失敗，請檢查帳號密碼或角色';
    
    if (error.response) {
      errorMsg = error.response.data?.message || errorMsg;
      console.log('伺服器回應:', error.response.data);
    } else if (error.message) {
      errorMsg = error.message;
    }
    
    showAlert('error', errorMsg, '錯誤');
  } finally {
    console.groupEnd();
  }
};

// 導航
const goToHome = () => router.push('/');
const goToRegister = () => router.push('/register');

// 測試登入功能
const handleTestLogin = async () => {
  let role = selectedRole.value || '';
  if (!role) {
    showAlert('error', '請先選擇角色', '錯誤');
    return;
  }
  // 設置登入狀態與 localStorage
  authStore.isLoggedIn = true;
  authStore.userRole = role;
  authStore.currentUser = role === 'user' ? 'test_user' : 
                         (role === 'medical' ? 'test_hospital' : 'test_insurer');
  localStorage.setItem('token', 'testtoken');
  localStorage.setItem('role', role);
  localStorage.setItem('id_number', authStore.currentUser);

  showAlert('success', '測試登入成功！正在為您導向...', '成功');
  setTimeout(() => {
    if (role === 'user') {
      router.push('/user-dashboard');
    } else if (role === 'medical') {
      router.push('/hospital-dashboard');
    } else if (role === 'insurer') {
      router.push('/other-user-dashboard');
    }
    // 強制刷新頁面，確保 dashboard 正確顯示
    window.location.reload();
  }, 1000);
};
</script>

<style scoped>
/* 全局樣式 */
.login-page {
  background-color: #F9F7F4;
  min-height: 100vh;
}

/* 登入卡片 */
.login-card {
  border-radius: 24px !important;
  background: white !important;
  padding: 2.5rem !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03) !important;
}

/* 頂部標題區 */
.header-section {
  text-align: center;
  margin-bottom: 2rem;
}

.header-title {
  font-size: 1.75rem;
  font-weight: 900;
  color: #111827;
  margin: 0;
  letter-spacing: -0.5px;
}

.header-subtitle {
  font-size: 1rem;
  color: #6B7280;
  margin: 0.5rem 0 0;
}

/* 表單樣式 */
.login-form {
  max-width: 400px;
  margin: 0 auto;
}

.form-field {
  margin-bottom: 1.25rem !important;
}

.form-field :deep(.v-field) {
  border-radius: 16px !important;
  border-color: #E5E7EB !important;
}

.form-field :deep(.v-field__outline) {
  border-color: #E5E7EB !important;
}

.form-field :deep(.v-field--focused) {
  border-color: #111827 !important;
}

/* 按鈕樣式 */
.login-btn {
  background-color: #F8F441 !important;
  color: #111827 !important;
  border-radius: 16px !important;
  font-weight: 600 !important;
  margin-bottom: 1rem !important;
  transition: all 0.2s ease !important;
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(248, 244, 65, 0.2) !important;
  background-color: #F9F650 !important;
}

.test-btn {
  background-color: #F3F4F6 !important;
  color: #111827 !important;
  border-radius: 16px !important;
  font-weight: 500 !important;
  margin-bottom: 1.5rem !important;
  transition: all 0.2s ease !important;
}

.test-btn:hover {
  transform: translateY(-2px);
  background-color: #E5E7EB !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05) !important;
}

/* 導航按鈕 */
.nav-buttons {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
}

.nav-btn {
  flex: 1;
  background-color: white !important;
  color: #6B7280 !important;
  border: 1px solid #E5E7EB !important;
  border-radius: 16px !important;
  font-weight: 500 !important;
  transition: all 0.2s ease !important;
}

.nav-btn:hover {
  transform: translateY(-2px);
  color: #111827 !important;
  border-color: #9CA3AF !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05) !important;
}

/* 提示訊息 */
.alert-message {
  margin-bottom: 1.5rem !important;
  border-radius: 16px !important;
}

/* 分隔線 */
.divider {
  border-color: #E5E7EB !important;
  margin: 1.5rem 0 !important;
}

/* 調試面板 */
.debug-section {
  margin-top: 1.5rem;
}

.debug-panel {
  border-radius: 16px !important;
  overflow: hidden !important;
  border: 1px solid #E5E7EB !important;
}

.debug-title {
  font-size: 0.875rem !important;
  color: #6B7280 !important;
}

.debug-info {
  background-color: #F9FAFB !important;
  border-radius: 8px !important;
  padding: 1rem !important;
  font-family: monospace !important;
  font-size: 12px !important;
  color: #374151 !important;
}

/* RWD 適配 */
@media (max-width: 600px) {
  .login-card {
    padding: 1.5rem !important;
    margin: 1rem;
  }
  
  .header-title {
    font-size: 1.5rem;
  }
  
  .nav-buttons {
    flex-direction: column;
  }
  
  .nav-btn {
    width: 100%;
  }
}
</style>