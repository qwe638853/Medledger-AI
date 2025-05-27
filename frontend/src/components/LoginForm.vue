<template>
  <div class="login-page">
    <v-container class="fill-height">
      <v-row align="center" justify="center">
                  <v-col cols="12" sm="8" md="6" lg="4">
            <v-slide-y-transition>
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
              <!-- 步驟指示器 -->
              <div class="steps-indicator mb-6">
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
                  <div class="role-cards-wrapper">
                    <v-slide-y-transition group>
                      <div
                        v-for="role in roles"
                        :key="role.value"
                        class="role-card"
                        :class="{ 'role-card--selected': selectedRole === role.value }"
                        @click="handleRoleSelect(role.value)"
                      >
                        <div class="role-card__icon-wrapper">
                          <v-icon 
                            :color="selectedRole === role.value ? 'primary' : 'grey'" 
                            size="32"
                          >{{ role.icon }}</v-icon>
                        </div>
                        <div class="role-card__content">
                          <h3 class="role-card__title">{{ role.text }}</h3>
                          <p class="role-card__description">{{ role.description }}</p>
                        </div>
                        <v-scale-transition>
                          <div v-if="selectedRole === role.value" class="role-card__check">
                            <v-icon color="primary" size="24">mdi-check-circle</v-icon>
                          </div>
                        </v-scale-transition>
                      </div>
                    </v-slide-y-transition>
                  </div>
                  
                  <v-slide-y-transition>
                    <div class="step-actions" v-if="selectedRole">
                      <v-btn
                        class="next-btn"
                        block
                        @click="goToStep2"
                        elevation="0"
                        height="48"
                      >
                        下一步
                        <v-icon class="ms-2">mdi-arrow-right</v-icon>
                      </v-btn>
                    </div>
                  </v-slide-y-transition>
                </div>
              </v-fade-transition>

              <!-- 步驟 2：帳號密碼輸入 -->
              <v-fade-transition>
                <div v-if="currentStep === 2" class="step-container">
                  <!-- 返回按鈕 -->
                  <v-btn
                    class="back-btn mb-6"
                    variant="text"
                    @click="currentStep = 1"
                    prepend-icon="mdi-arrow-left"
                  >
                    返回選擇角色
                  </v-btn>

                  <v-slide-y-transition>
                    <div class="login-form-fields">
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

                  <v-slide-y-transition>
                    <div class="step-actions">
                      <!-- 登入按鈕 -->
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
                        class="test-btn mt-3"
                        block
                        @click="handleTestLogin"
                        elevation="0"
                        height="48"
                        prepend-icon="mdi-test-tube-outline"
                      >
                        測試登入
                      </v-btn>
                    </div>
                  </v-slide-y-transition>
                </div>
              </v-fade-transition>
              
              <!-- 導航按鈕組 -->
              <div class="nav-buttons mt-6">
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
    description: '個人用戶登入，查看健康數據'
  },
  { 
    text: '醫療機構', 
    value: 'medical',
    icon: 'mdi-hospital-building',
    description: '醫院、診所等醫療服務提供者'
  },
  { 
    text: '保險業者', 
    value: 'insurer',
    icon: 'mdi-shield-account',
    description: '保險公司、保險服務提供商'
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

/* 角色選擇卡片容器 */
.role-cards-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  width: 100%;
  flex: 1;
}

/* 角色卡片基本樣式 */
.role-card {
  position: relative;
  display: flex;
  align-items: center;
  padding: 1.5rem;
  border: 2px solid #e5e7eb;
  border-radius: 16px;
  background: white;
  cursor: pointer;
  transition: all 0.3s ease;
}

.role-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}

/* 選中狀態樣式 */
.role-card--selected {
  border-color: var(--v-theme-primary);
  background-color: var(--v-theme-primary-lighten-5);
  box-shadow: 0 4px 12px rgba(var(--v-theme-primary), 0.1);
}

/* 圖標包裝器 */
.role-card__icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  min-width: 64px;
  min-height: 64px;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #f3f4f6;
  margin-right: 1rem;
  transition: all 0.3s ease;
  position: relative;
  overflow: visible;
}

.role-card--selected .role-card__icon-wrapper {
  background: var(--v-theme-primary-lighten-4);
}

.role-card__icon-wrapper :deep(.v-icon) {
  font-size: 32px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.role-card--selected .role-card__icon-wrapper :deep(.v-icon) {
  transform: scale(1.1);
}

/* 內容區域 */
.role-card__content {
  flex: 1;
}

.role-card__title {
  font-size: 1.125rem;
  font-weight: 600;
  color: #111827;
  margin: 0 0 0.5rem;
}

.role-card__description {
  font-size: 0.875rem;
  color: #6b7280;
  margin: 0;
  line-height: 1.4;
}

/* RWD 適配 */
@media (max-width: 768px) {
  .role-cards-container {
    grid-template-columns: 1fr;
  }
  
  .role-card {
    padding: 1.25rem;
  }
  
  .role-card__icon-wrapper {
    width: 56px;
    height: 56px;
  }
}

/* 步驟指示器樣式 */
.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1rem 0;
  margin-bottom: 2rem;
}

.step {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  flex: 1;
  max-width: 120px;
}

.step-circle {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  border: 2px solid #e5e7eb;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.5rem;
  position: relative;
  transition: all 0.3s ease;
  font-weight: 600;
  color: #888;
}

.step.active .step-circle {
  border-color: #F8F441;
  background: #F8F441;
  color: #111827;
}

.step.completed .step-circle {
  border-color: #463F3A;
  background: #463F3A;
  color: white;
}

.step-label {
  font-size: 0.875rem;
  color: #888;
  font-weight: 500;
  transition: all 0.3s ease;
}

.step.active .step-label {
  color: #111827;
  font-weight: 600;
}

.step.completed .step-label {
  color: #111827;
  font-weight: 600;
}

.step-line {
  flex: 1;
  height: 2px;
  background: #e5e7eb;
  margin: 0 0.5rem;
  margin-bottom: 2rem;
  max-width: 100px;
  transition: all 0.3s ease;
}

.step.completed + .step-line {
  background: #111827;
}

/* RWD 適配 */
@media (max-width: 600px) {
  .steps-indicator {
    padding: 0.5rem 0;
  }
  
  .step-circle {
    width: 32px;
    height: 32px;
    font-size: 0.875rem;
  }
  
  .step-label {
    font-size: 0.75rem;
  }
  
  .step-line {
    max-width: 60px;
  }
}

/* 返回按鈕樣式 */
.back-btn {
  color: #6b7280 !important;
  font-weight: 500;
}

.back-btn:hover {
  color: #111827 !important;
  background: rgba(0, 0, 0, 0.04);
}

/* 步驟容器 */
.step-container {
  display: flex;
  flex-direction: column;
  min-height: 400px; /* 確保容器有足夠的高度 */
}

/* 步驟操作按鈕容器 */
.step-actions {
  position: sticky;
  bottom: 0;
  background: white;
  padding: 1rem 0;
  margin-top: auto;
  width: 100%;
}

/* 下一步按鈕樣式 */
.next-btn {
  background-color: #F8F441 !important;
  color: #111827 !important;
  font-weight: 600 !important;
  border-radius: 16px !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  height: 48px !important;
  letter-spacing: 0.5px !important;
  text-transform: none !important;
  font-size: 1rem !important;
  border: 2px solid transparent !important;
}

.next-btn:hover {
  transform: translateY(-2px);
  background-color: #F9F650 !important;
  box-shadow: 0 8px 24px rgba(248, 244, 65, 0.25) !important;
  border-color: #F8F441 !important;
}

.next-btn:active {
  transform: translateY(0);
  background-color: #F7F332 !important;
  box-shadow: 0 4px 12px rgba(248, 244, 65, 0.15) !important;
}

.next-btn:disabled {
  background-color: #F3F4F6 !important;
  color: #9CA3AF !important;
  transform: none;
  box-shadow: none;
  border-color: transparent !important;
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
  
  .steps-indicator {
    padding: 0.5rem 0;
  }
  
  .step-circle {
    width: 32px;
    height: 32px;
  }
  
  .step-label {
    font-size: 0.75rem;
  }
  
  .step-line {
    max-width: 60px;
  }
}
</style>