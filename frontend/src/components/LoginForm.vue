<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-3">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>登入</v-toolbar-title>
          </v-toolbar>
          
          <!-- 新增 Alert 組件 -->
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
            <v-form @submit.prevent="handleSubmit" ref="form" v-model="valid">
              <!-- 選擇角色 -->
              <v-select
                v-model="selectedRole"
                :items="roles"
                item-title="text"
                item-value="value"
                label="選擇角色"
                prepend-icon="mdi-account-group"
                outlined
                dense
                :rules="[rules.required]"
                class="mb-4"
              />
              
              <!-- 身分證號/員工編號 -->
              <v-text-field
                v-model="username"
                label="身分證號/員工編號"
                prepend-icon="mdi-account"
                outlined
                dense
                :rules="[rules.required]"
                class="mb-4"
              />
              
              <!-- 密碼 -->
              <v-text-field
                v-model="password"
                label="密碼"
                prepend-icon="mdi-lock"
                :type="showPassword ? 'text' : 'password'"
                :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append="showPassword = !showPassword"
                outlined
                dense
                :rules="[rules.required, rules.minLength]"
                class="mb-6"
              />
              
              <!-- 提交按鈕 -->
              <v-btn
                :loading="loading"
                color="primary"
                block
                type="submit"
                :disabled="!valid"
                elevation="2"
                height="44"
                class="mb-4"
              >
                登入
              </v-btn>
              
              <!-- 導航按鈕 -->
              <div class="d-flex justify-space-between mb-4">
                <v-btn text color="primary" @click="goToHome">
                  返回首頁
                </v-btn>
                <v-btn text color="primary" @click="goToRegister">
                  註冊
                </v-btn>
              </div>
              
              <!-- 開發模式下顯示表單數據 -->
              <pre v-if="isDevelopment" class="debug-info mt-4 pa-2">{{ debugInfo }}</pre>
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
  { text: '其他用戶', value: 'other' }
];
const form = ref(null);
const valid = ref(false);

// 判斷是否為開發環境
const isDevelopment = import.meta.env.MODE === 'development';

// 用於調試的數據
const debugInfo = computed(() => {
  const roleText = roles.find(r => r.value === selectedRole.value)?.text || selectedRole.value;
  return {
    username: username.value,
    password: password.value ? '******' : '',
    role: selectedRole.value,
    roleText,
    valid: valid.value,
    loading: loading.value
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
  if (!form.value.validate()) {
    console.log('表單驗證失敗');
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
</script>

<style scoped>
.fill-height {
  min-height: calc(100vh - 64px);
  background-color: #f5f5f5;
}

.v-card {
  border-radius: 8px !important;
}

.v-text-field :deep(.v-input__slot) {
  min-height: 44px !important;
}

.v-select :deep(.v-input__slot) {
  min-height: 44px !important;
}

.debug-info {
  background-color: #f9f9f9;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-family: monospace;
  font-size: 12px;
  overflow: auto;
  max-height: 150px;
}
</style>