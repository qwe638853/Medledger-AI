<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-3">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title class="text-h5 font-weight-bold">企業健康數據管理平台</v-toolbar-title>
          </v-toolbar>
          
          <!-- 系統訊息提示 -->
          <v-alert
            v-if="alertInfo.show"
            :type="alertInfo.type"
            :title="alertInfo.title"
            :icon="alertInfo.icon"
            closable
            border
            class="ma-4"
            @click:close="alertInfo.show = false"
          >
            {{ alertInfo.message }}
          </v-alert>
          
          <v-card-text class="pt-6">
            <v-form
              @submit.prevent="handleSubmit"
              ref="form"
              v-model="valid"
              lazy-validation
            >
              <!-- 帳號資訊區塊 -->
              <div class="section-title mb-3">
                <h3 class="text-subtitle-1 text-blue-grey-darken-1 font-weight-bold mb-0">
                  <v-icon color="primary" class="mr-1">mdi-account-key</v-icon>
                  帳號資訊
                </h3>
                <v-divider class="mt-2"></v-divider>
              </div>
              
              <!-- 選擇角色 -->
              <v-select
                v-model="selectedRole"
                :items="roles"
                item-title="text"
                item-value="value"
                label="選擇角色"
                prepend-inner-icon="mdi-account-group"
                variant="outlined"
                :rules="[rules.required]"
                class="mb-4"
                density="comfortable"
              />
              
              <!-- 身分證號/醫療機構ID/保險業者ID -->
              <v-text-field
                v-model="username"
                :label="usernameLabel"
                prepend-inner-icon="mdi-account"
                variant="outlined"
                :rules="[rules.required]"
                class="mb-4"
                density="comfortable"
                :placeholder="usernamePlaceholder"
                clearable
              />
              
              <!-- 密碼 -->
              <v-text-field
                v-model="password"
                label="密碼"
                prepend-inner-icon="mdi-lock"
                :type="showPassword ? 'text' : 'password'"
                :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append-inner="showPassword = !showPassword"
                variant="outlined"
                :rules="[rules.required, rules.minLength]"
                class="mb-6"
                density="comfortable"
                clearable
              />
              
              <!-- 操作按鈕區 -->
              <div class="action-section mb-4">
                <!-- 登入按鈕 -->
              <v-btn
                :loading="loading"
                color="primary"
                block
                type="submit"
                :disabled="!valid"
                elevation="2"
                  height="48"
                  class="mb-4 text-body-1 font-weight-bold"
                  prepend-icon="mdi-login"
              >
                  登入系統
              </v-btn>
              
                <!-- 測試功能按鈕 -->
              <v-btn
                color="info"
                block
                @click="handleTestLogin"
                elevation="2"
                  height="48"
                class="mb-4"
                  prepend-icon="mdi-test-tube"
                  variant="outlined"
              >
                測試登入
              </v-btn>
              </div>
              
              <!-- 導航按鈕 -->
              <div class="d-flex justify-space-between mb-4">
                <v-btn text color="primary" @click="goToHome" prepend-icon="mdi-home">
                  返回首頁
                </v-btn>
                <v-btn text color="success" @click="goToRegister" prepend-icon="mdi-account-plus">
                  註冊帳號
                </v-btn>
              </div>
              
              <!-- 可展開的調試資訊 -->
              <v-expand-transition>
                <div v-if="isDevelopment">
                  <v-divider class="my-3"></v-divider>
                  <v-expansion-panels variant="accordion" class="mt-4">
                    <v-expansion-panel
                      title="開發調試資訊"
                      bg-color="grey-lighten-4"
                    >
                      <v-expansion-panel-text>
                        <pre class="debug-info pa-2">{{ JSON.stringify(debugInfo, null, 2) }}</pre>
                      </v-expansion-panel-text>
                    </v-expansion-panel>
                  </v-expansion-panels>
                </div>
              </v-expand-transition>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
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
.fill-height {
  min-height: calc(100vh - 64px);
  background-color: #f8f9fa;
}

.v-card {
  border-radius: 12px !important;
  overflow: hidden;
}

.debug-info {
  background-color: #f8f9fa;
  border: 1px solid #e9ecef;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
  overflow: auto;
  max-height: 300px;
  white-space: pre-wrap;
  word-break: break-all;
}

.section-title {
  position: relative;
}

.action-section {
  margin-top: 12px;
}

/* 統一按鈕高度與間距 */
.v-btn {
  letter-spacing: 0.5px;
  transition: all 0.2s ease;
}

.v-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}
</style>