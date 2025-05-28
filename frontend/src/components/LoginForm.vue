<template>
  <div class="auth-page">
    <v-container class="fill-height">
      <v-row align="center" justify="center">
        <v-col cols="12" sm="10" md="10" lg="10" xl="8" style="display: flex; justify-content: center;">
          <v-slide-y-transition>
            <!-- 主要登入卡片 -->
            <v-card class="auth-card" elevation="0">
              <!-- 頂部標題區 -->
              <div class="header-section">
                <h1 class="header-title">登入</h1>
              </div>
            
            <!-- 系統訊息提示 -->
            <v-alert
              v-if="alertInfo.show"
              :type="alertInfo.type"
              :title="alertInfo.title"
              :icon="alertInfo.icon"
              closable
              border
              class="ma-2"
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
              class="auth-form pt-2"
            >
              <!-- 步驟指示器 -->
              <div class="steps-indicator mb-4">
                <div 
                  class="step"
                  :class="{ 'active': currentStep === 1, 'completed': currentStep > 1 }"
                >
                  <div class="step-circle">
                    <span v-if="currentStep <= 1">1</span>
                    <v-icon v-else>mdi-check</v-icon>
                  </div>
                  <div class="step-label">選擇角色</div>
                </div>
                <div class="step-line"></div>
                <div 
                  class="step"
                  :class="{ 'active': currentStep === 2 }"
                >
                  <div class="step-circle">2</div>
                  <div class="step-label">帳號登入</div>
                </div>
              </div>

              <!-- 步驟 1：角色選擇 -->
              <v-fade-transition>
                <div v-if="currentStep === 1" class="step-container">
                  <RoleSelector
                    v-model="selectedRole"
                    @next="goToStep2"
                  />
                </div>
              </v-fade-transition>

              <!-- 步驟 2：帳號密碼輸入 -->
              <v-fade-transition>
                <div v-if="currentStep === 2" class="step-container">
                  <v-card flat class="mt-6 pa-4 rounded-lg">
                    <div class="section-title">
                      <h3 class="text-h5 font-weight-bold">
                        <v-icon color="primary" class="me-2">mdi-account-key</v-icon>
                        帳號登入
                      </h3>
                      <div class="text-body-2 text-grey">請輸入您的帳號密碼</div>
                      <v-divider class="mt-2"></v-divider>
                    </div>

                    <v-slide-y-transition>
                      <div class="login-form-fields mt-4">
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
                      </div>
                    </v-slide-y-transition>

                    <div class="mt-6 d-flex justify-space-between">
                      <v-btn class="secondary-btn" @click="currentStep = 1">
                        <v-icon class="me-2">mdi-arrow-left</v-icon>
                        返回
                      </v-btn>
                      <v-btn
                        :loading="loading"
                        class="primary-btn"
                        type="submit"
                        :disabled="!valid"
                      >
                        登入系統
                        <v-icon class="ms-2">mdi-login-variant</v-icon>
                      </v-btn>
                    </div>
                  </v-card>
                </div>
              </v-fade-transition>
              
              <!-- 導航按鈕組 -->
              <div class="nav-links">
                <v-btn
                  class="nav-btn"
                  @click="goToHome"
                >
                  <v-icon class="me-2">mdi-home-outline</v-icon>
                  返回首頁
                </v-btn>
                <v-btn
                  class="nav-btn"
                  @click="goToRegister"
                >
                  <v-icon class="me-2">mdi-account-plus-outline</v-icon>
                  還沒有帳號？註冊
                </v-btn>
              </div>
              
              <!-- 開發調試資訊 
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
            -->
            </v-form>
          </v-card>
          </v-slide-y-transition>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores';
import RoleSelector from './common/RoleSelector.vue';
import '@/styles/auth-forms.css';

const router = useRouter();
const authStore = useAuthStore();

// 表單數據
const username = ref('');
const password = ref('');
const showPassword = ref(false);
const selectedRole = ref('');
const loading = computed(() => authStore.loading);
const roles = [
  { 
    text: '一般用戶', 
    value: 'user',
    icon: 'mdi-account',
  },
  { 
    text: '醫療機構', 
    value: 'medical',
    icon: 'mdi-hospital-building',
  },
  { 
    text: '保險業者', 
    value: 'insurer',
    icon: 'mdi-shield-account',
  }
];
const form = ref(null);
const valid = ref(false);

// 添加步驟控制
const currentStep = ref(1);

// 處理角色選擇
const handleRoleSelect = (roleValue) => {
  selectedRole.value = roleValue;
  console.log('Selected role:', roleValue); // 用於調試
};

// 前往第二步
const goToStep2 = () => {
  if (selectedRole.value) {
    currentStep.value = 2;
    console.log('Moving to step 2'); // 用於調試
  }
};

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
body, #app, .v-application, .v-app {
  background: #f9f7f4 !important;
  font-size: 1.18rem !important;
  line-height: 1.7 !important;
}

.auth-page {
  background-color: #F9F7F4;
  min-height: calc(100vh - 64px);
  padding: 1rem;
}

/* 登入卡片 */
.auth-card {
  max-width: 1200px;
  width: 100%;
  border-radius: 32px !important;
  background: #fff !important;
  padding: 3rem 2.5rem !important;
  border: 1px solid #e5e7eb !important;
  box-shadow: 0 8px 32px rgba(0,0,0,0.10) !important;
  font-size: 1.18rem !important;
}

/* 頂部標題區 */
.header-section {
  text-align: center;
  margin-bottom: 2.5rem;
}

.header-title {
  font-size: 3rem !important;
  font-weight: 900;
  color: #111827;
  margin: 0;
  letter-spacing: -0.5px;
}

/* 表單樣式 */
.auth-form {
  max-width: 1000px;
  margin: 0 auto;
  font-size: 1.18rem;
}

.form-field {
  margin-bottom: 1.5rem !important;
}

/* 步驟指示器樣式 */
.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 0 auto 3rem;
  padding: 0;
  width: 100%;
  max-width: 600px;
  position: relative;
}

.steps-indicator::before {
  content: '';
  position: absolute;
  top: 28px;
  left: 80px;
  right: 80px;
  height: 2px;
  background: #e5e7eb;
  z-index: 1;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  z-index: 2;
  flex: 0 1 auto;
  margin: 0 2rem;
}

.step-circle {
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: 2px solid #e5e7eb;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.75rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 600;
  color: #888;
  font-size: 1.7rem;
  position: relative;
  z-index: 2;
}

.step.active .step-circle {
  border-color: #00B8D9;
  background: #00B8D9;
  color: white;
  transform: scale(1.1);
  box-shadow: 0 4px 12px rgba(0, 184, 217, 0.25);
}

.step.completed .step-circle {
  border-color: #9CA3AF;
  background: #9CA3AF;
  color: white;
}

.step-label {
  font-size: 1.18rem;
  color: #888;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  white-space: nowrap;
  margin-top: 0.25rem;
}

.step.active .step-label {
  color: #111827;
  font-weight: 600;
}

.step.completed .step-label {
  color: #6B7280;
  font-weight: 600;
}

/* 按鈕樣式 */
.v-btn {
  border-radius: 20px !important;
  text-transform: none !important;
  font-size: 1.5rem !important;
  font-weight: 600 !important;
  letter-spacing: 0 !important;
  height: 48px !important;
  min-width: 140px !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.primary-btn {
  background-color: #00B8D9 !important;
  color: white !important;
  border: none !important;
  box-shadow: 0 4px 16px rgba(0, 184, 217, 0.25) !important;
}

.primary-btn:hover {
  background-color: #0095B0 !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.35) !important;
}

.secondary-btn {
  background-color: white !important;
  color: #6B7280 !important;
  border: 1px solid #e5e7eb !important;
}

.secondary-btn:hover {
  background-color: #f9fafb !important;
  color: #00B8D9 !important;
  border-color: #00B8D9 !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05) !important;
}

/* 區塊標題樣式 */
.section-title {
  margin-bottom: 2rem;
}

.section-title h3 {
  font-size: 2rem !important;
  font-weight: 700 !important;
  color: #111827 !important;
  margin-bottom: 0.5rem !important;
}

.section-title .text-body-2 {
  font-size: 1.1rem !important;
  color: #6B7280 !important;
  line-height: 1.5;
}

.section-title .v-divider {
  margin-top: 1rem !important;
  opacity: 0.1;
}

/* 導航按鈕組 */
.nav-links {
  margin-top: 3rem;
  display: flex;
  justify-content: space-between;
  gap: 1.5rem;
  padding: 0 1rem;
}

.nav-btn {
  flex: 1;
  background: transparent !important;
  color: #6B7280 !important;
  border: 1px solid #E5E7EB !important;
  border-radius: 16px !important;
  padding: 1rem 1.5rem !important;
  font-size: 1.1rem !important;
  font-weight: 500 !important;
  transition: all 0.3s ease !important;
  height: 48px !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
}

.nav-btn .v-icon {
  margin-right: 0.5rem;
}

.nav-btn:hover {
  color: #00B8D9 !important;
  border-color: #00B8D9 !important;
  background: rgba(0, 184, 217, 0.05) !important;
  text-decoration: none;
}

/* RWD 適配 */
@media (max-width: 900px) {
  .auth-card {
    padding: 2rem 1.5rem !important;
  }
  
  .header-title {
    font-size: 2.5rem !important;
  }
  
  .steps-indicator {
    max-width: 100%;
    margin-bottom: 2rem;
  }
  
  .steps-indicator::before {
    left: 60px;
    right: 60px;
  }
  
  .step {
    margin: 0 1rem;
  }
  
  .step-circle {
    width: 48px;
    height: 48px;
    font-size: 1.4rem;
  }
  
  .step-label {
    font-size: 1rem;
  }
  
  .nav-links {
    margin-top: 2rem;
    padding: 0;
  }
}

@media (max-width: 600px) {
  .auth-card {
    padding: 1.5rem 1rem !important;
  }
  
  .header-section {
    margin-bottom: 2rem;
  }
  
  .steps-indicator {
    margin-bottom: 1.5rem;
  }
  
  .steps-indicator::before {
    left: 40px;
    right: 40px;
  }
  
  .step {
    margin: 0 0.5rem;
  }
  
  .step-circle {
    width: 40px;
    height: 40px;
    font-size: 1.2rem;
    margin-bottom: 0.5rem;
  }
  
  .nav-links {
    flex-direction: column;
    margin-top: 1.5rem;
    gap: 1rem;
  }
  
  .nav-btn {
    width: 100%;
  }
}
</style>