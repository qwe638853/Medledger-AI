<template>
  <div class="login-page">
    <v-container class="fill-height">
      <v-row align="center" justify="center">
        <v-col cols="12" sm="10" md="10" lg="10" xl="8" style="display: flex; justify-content: center;">
          <v-slide-y-transition>
            <!-- 主要登入卡片 -->
            <v-card class="login-card" elevation="0">
              <!-- 頂部標題區 -->
              <div class="header-section">
                <h1 class="header-title">登入</h1>
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
  min-height: calc(100vh - 64px);
  padding: 2rem 1rem;
}

/* 登入卡片樣式與註冊一致 */
.login-card {
  max-width: 1200px;
  width: 100%;
  border-radius: 32px !important;
  background: #fff !important;
  padding: 3.5rem 2.5rem !important;
  border: 1px solid #e5e7eb !important;
  box-shadow: 0 8px 32px rgba(0,0,0,0.10) !important;
  font-size: 1.5rem !important;
}

/* 頂部標題區 */
.header-section {
  text-align: center;
  margin-bottom: 2rem;
}

.header-title {
  font-size: 3rem !important;
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
  max-width: 1000px;
  margin: 0 auto;
  font-size: 1.5rem;
}

.form-field {
  margin-bottom: 2rem !important;
}

.form-field :deep(.v-field),
.form-field :deep(.v-label) {
  font-size: 1.3rem !important;
  min-height: 56px !important;
}

.form-field :deep(.v-field--focused) {
  border-color: #111827 !important;
}

.form-field :deep(.v-label) {
  color: #6B7280 !important;
}

/* 按鈕樣式 */
.v-btn, .login-btn, .test-btn, .next-btn, .nav-btn {
  font-size: 1.5rem !important;
  height: 56px !important;
  min-width: 140px !important;
  border-radius: 20px !important;
}

.login-btn {
  background-color: #F8F441 !important;
  color: #111827 !important;
  font-weight: 700 !important;
  margin-bottom: 1.5rem !important;
}

.login-btn:hover {
  background-color: #f9f650 !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(248, 244, 65, 0.25) !important;
}

.test-btn {
  background-color: #F3F4F6 !important;
  color: #6B7280 !important;
  font-weight: 500 !important;
  margin-bottom: 1.5rem !important;
}

.test-btn:hover {
  background-color: #E5E7EB !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05) !important;
  color: #111827 !important;
}

/* 導航按鈕 */
.nav-buttons {
  display: flex;
  justify-content: space-between;
  gap: 1rem;
  margin-top: 1.5rem;
}

.nav-btn {
  flex: 1;
  background-color: white !important;
  color: #6B7280 !important;
  border: 1px solid #e5e7eb !important;
  font-weight: 500 !important;
}

.nav-btn:hover {
  transform: translateY(-2px);
  color: #111827 !important;
  border-color: #9CA3AF !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05) !important;
}

/* 角色卡片樣式完全與註冊頁圖片一致 (藍色選中狀態) */
.role-cards-wrapper {
  display: flex;
  flex-direction: row;
  gap: 1.5rem;
  width: 100%;
  margin: 1rem 0; /* 與註冊頁圖片排版一致 */
  justify-content: center;
  flex-wrap: nowrap; /* 確保橫向排列不換行 */
  overflow-x: unset; /* 桌面不顯示滾動條 */
}
.role-card {
  position: relative;
  display: flex;
  align-items: center;
  min-width: 280px; /* 與註冊頁圖片卡片寬度接近 */
  max-width: 340px; /* 與註冊頁圖片卡片寬度接近 */
  flex: 1 1 0;
  padding: 2.2rem; /* 與註冊頁圖片 padding 一致 */
  font-size: 1.15rem; /* 與註冊頁圖片字體大小一致 */
  border: 2px solid #e5e7eb; /* 未選中時的邊框，與註冊頁圖片一致 */
  border-radius: 20px; /* 圓角，與註冊頁圖片一致 */
  background: white; /* 未選中時背景 */
  cursor: pointer;
  transition: all 0.3s ease;
  height: auto;
  box-sizing: border-box;
}
.role-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
}
.role-card--selected {
  border-color: #00B8D9; /* 選中時的邊框顏色，與模板圖片一致 (藍色) */
  background-color: #E6FDFF; /* 選中時的背景顏色，與模板圖片一致 (淡藍色) */
  box-shadow: 0 6px 24px rgba(0, 184, 217, 0.18); /* 選中時的陰影，與模板圖片一致 */
}
.role-card__icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 80px; /* icon 容器大小，與註冊頁圖片一致 */
  height: 80px; /* icon 容器大小，與註冊頁圖片一致 */
  border-radius: 50%;
  background: #f3f4f6; /* icon 容器背景 */
  margin-right: 1rem; /* icon 與文字間距，與註冊頁圖片一致 */
  transition: all 0.3s ease;
  font-size: 2.2rem; /* icon 大小，與註冊頁圖片一致 */
}
.role-card--selected .role-card__icon-wrapper {
  background: #B2EBf2; /* 選中時 icon 容器背景，與模板圖片一致 */
}
.role-card .v-icon {
  /* icon 顏色通過 template 中的 :color 綁定控制 */
}
.role-card__content {
  flex: 1;
}
.role-card__title {
  font-size: 1.35rem; /* 標題字體大小，與註冊頁圖片一致 */
  font-weight: 600;
  color: #111827; /* 標題字體顏色，與模板圖片一致 */
  margin: 0 0 0.5rem;
}
.role-card__description {
  font-size: 1.1rem; /* 描述字體大小，與註冊頁圖片一致 */
  color: #6b7280; /* 描述字體顏色，與模板圖片一致 */
  margin: 0;
  line-height: 1.4;
}
.role-card__check {
  position: absolute;
  top: 1rem;
  right: 1rem;
}
.role-card__check .v-icon {
  font-size: 36px !important; /* 打勾 Icon 大小，與模板圖片一致 */
  color: #00B8D9 !important; /* 打勾 Icon 顏色，與模板圖片一致 (藍色) */
}

/* Alert 訊息樣式 */
.alert-message {
  border-radius: 16px !important;
  margin-bottom: 1.5rem !important;
}

/* RWD 適配 */
@media (max-width: 900px) {
  .login-card {
    padding: 2rem 1rem !important;
    font-size: 1.1rem;
  }
  .header-title {
    font-size: 1.5rem !important;
  }
  .login-form {
    max-width: 98vw;
    font-size: 1rem;
  }
  .form-field :deep(.v-field),
  .form-field :deep(.v-label) {
    font-size: 1rem !important;
    min-height: 48px !important;
  }
  .v-btn, .login-btn, .test-btn, .next-btn, .nav-btn {
    font-size: 1rem !important;
    height: 48px !important;
    min-width: 120px !important;
  }
  .role-card {
    padding: 1.2rem;
    font-size: 1rem;
  }
  .role-card__icon-wrapper {
    width: 56px;
    height: 56px;
    font-size: 1.3rem;
  }
  .role-card__title {
    font-size: 1.1rem;
  }
  .role-card__description {
    font-size: 0.95rem;
  }
  .step-circle {
    width: 32px;
    height: 32px;
    font-size: 1rem;
  }
  .step-label {
    font-size: 0.85rem;
  }
}

/* 步驟指示器樣式 */
.steps-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  margin: 2.5rem 0;
  padding: 0;
  font-size: 1.6rem;
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
  width: 56px;
  height: 56px;
  border-radius: 50%;
  border: 2px solid #e5e7eb;
  background: white;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.5rem;
  position: relative;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  font-weight: 600;
  color: #888;
  font-size: 1.7rem;
}

.step.active .step-circle {
  border-color: #F8F441;
  background: #F8F441;
  color: #111827;
  transform: scale(1.1);
}

.step.completed .step-circle {
  border-color: #463F3A;
  background: #463F3A;
  color: white;
  transform: scale(1);
}

.step-label {
  font-size: 1.5rem;
  color: #888;
  font-weight: 500;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  opacity: 0.8;
}

.step.active .step-label {
  color: #111827;
  font-weight: 600;
  opacity: 1;
  transform: translateY(-2px);
}

.step.completed .step-label {
  color: #111827;
  font-weight: 600;
  opacity: 1;
}

.step-line {
  flex: 1;
  height: 2px;
  background: #e5e7eb;
  margin: 0 0.5rem;
  margin-bottom: 2rem;
  max-width: 100px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  transform-origin: left center;
}

.step.completed + .step-line {
  background: #111827;
  transform: scaleX(1.1);
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
  font-size: 1.15rem !important;
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

</style>