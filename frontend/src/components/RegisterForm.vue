<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="6" lg="4">
        <v-card class="elevation-3">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>註冊</v-toolbar-title>
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
            <v-form @submit.prevent="handleRegister" ref="form" v-model="valid">
              <!-- 選擇角色 -->
              <v-select
                v-model="registerForm.selectedRole"
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
                v-model="registerForm.idNumber"
                label="身分證號/員工編號 (選填)"
                prepend-icon="mdi-account"
                outlined
                dense
                class="mb-4"
              />
              
              <!-- 密碼 -->
              <v-text-field
                v-model="registerForm.password"
                label="密碼"
                prepend-icon="mdi-lock"
                :type="showPassword ? 'text' : 'password'"
                :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append="showPassword = !showPassword"
                outlined
                dense
                :rules="[rules.required, rules.minLength]"
                class="mb-4"
              />
              
              <!-- 確認密碼 -->
              <v-text-field
                v-model="registerForm.confirmPassword"
                label="確認密碼"
                prepend-icon="mdi-lock-check"
                type="password"
                outlined
                dense
                :rules="[rules.required, rules.matchPassword]"
                class="mb-4"
              />
              
              <!-- 全名 -->
              <v-text-field
                v-model="registerForm.fullName"
                label="全名"
                prepend-icon="mdi-account-box"
                outlined
                dense
                :rules="[rules.required]"
                class="mb-4"
              />
              
              <!-- 性別 -->
              <v-select
                v-model="registerForm.gender"
                :items="genders"
                label="性別"
                prepend-icon="mdi-gender-male-female"
                outlined
                dense
                :rules="[rules.required]"
                class="mb-4"
              />
              
              <!-- 出生日期 -->
              <v-text-field
                v-model="registerForm.birthDate"
                label="出生日期"
                prepend-icon="mdi-calendar"
                type="date"
                outlined
                dense
                :rules="[rules.required, rules.validDate]"
                class="mb-4"
              />
              
              <!-- 電話號碼 -->
              <v-text-field
                v-model="registerForm.phoneNumber"
                label="電話號碼"
                prepend-icon="mdi-phone"
                type="tel"
                outlined
                dense
                :rules="[rules.required, rules.phone]"
                class="mb-4"
              />
              
              <!-- 電子郵件 -->
              <v-text-field
                v-model="registerForm.email"
                label="電子郵件"
                prepend-icon="mdi-email"
                type="email"
                outlined
                dense
                :rules="[rules.required, rules.email]"
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
                提交註冊
              </v-btn>
              
              <!-- ===================== 測試功能開始 ===================== -->
              <!-- 
                測試註冊按鈕：
                1. 用於在開發階段測試註冊功能
                2. 不需要填寫表單即可測試
                3. 使用預設的測試數據進行註冊
                4. 註冊成功後會自動導向到登入頁面
                
                修改說明：
                1. 如需修改測試數據，請修改 handleTestRegister 函數中的測試數據
                2. 如需修改註冊成功後的跳轉邏輯，請修改 handleTestRegister 函數中的重定向部分
              -->
              <v-btn
                color="info"
                block
                @click="handleTestRegister"
                elevation="2"
                height="44"
                class="mb-4"
              >
                測試註冊
              </v-btn>
              <!-- ===================== 測試功能結束 ===================== -->
              
              <!-- 導航按鈕 -->
              <div class="d-flex justify-space-between mb-4">
                <v-btn text color="primary" @click="goToHome">
                  返回首頁
                </v-btn>
                <v-btn text color="primary" @click="goToLogin">
                  已有帳號？登入
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
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

// 表單數據
const registerForm = ref({
  selectedRole: '',
  idNumber: '',
  password: '',
  confirmPassword: '',
  fullName: '',
  gender: '',
  birthDate: '',
  phoneNumber: '',
  email: ''
});
const loading = ref(false);
const showPassword = ref(false);
const roles = [
  { text: '一般用戶', value: 'user' },
  { text: '醫療機構', value: 'medical' },
  { text: '其他用戶', value: 'other' }
];
const genders = ['男', '女', '其他', '不願透露'];
const form = ref(null);
const valid = ref(false);

// 判斷是否為開發環境
const isDevelopment = import.meta.env.MODE === 'development';

// 用於調試的數據
const debugInfo = computed(() => {
  const roleText = roles.find(r => r.value === registerForm.value.selectedRole)?.text || registerForm.value.selectedRole;
  return {
    formData: {
      ...registerForm.value,
      password: registerForm.value.password ? '******' : '',
      confirmPassword: registerForm.value.confirmPassword ? '******' : '',
      roleText
    },
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
  minLength: value => (value && value.length >= 6) || '密碼至少需要6個字符',
  matchPassword: value => value === registerForm.value.password || '密碼不一致',
  email: value => /.+@.+\..+/.test(value) || '請輸入有效的電子郵件',
  phone: value => /^\d{10}$/.test(value) || '請輸入有效的10位電話號碼',
  validDate: value => !value || new Date(value) <= new Date() || '出生日期不得為未來日期'
};

// 表單提交
const handleRegister = async () => {
  if (!form.value.validate()) {
    console.log('表單驗證失敗');
    return;
  }
  
  loading.value = true;
  
  try {
    console.group('註冊請求');
    console.log('註冊數據(部分敏感資訊已隱藏):', {
      username: registerForm.value.idNumber,
      name: registerForm.value.fullName,
      gender: registerForm.value.gender,
      birthDate: registerForm.value.birthDate,
      email: registerForm.value.email,
      phone: registerForm.value.phoneNumber,
      role: registerForm.value.selectedRole,
      hasPassword: !!registerForm.value.password
    });
    
    await authStore.register({
      username: registerForm.value.idNumber,
      password: registerForm.value.password,
      name: registerForm.value.fullName,
      date: registerForm.value.birthDate,
      email: registerForm.value.email,
      phone: registerForm.value.phoneNumber,
      role: registerForm.value.selectedRole
    });
    
    // 註冊成功顯示
    showAlert('success', '註冊成功！即將為您導向登入頁面', '成功');
    console.log('註冊成功！');
    
    // 重置表單
    registerForm.value = {
      selectedRole: '',
      idNumber: '',
      password: '',
      confirmPassword: '',
      fullName: '',
      gender: '',
      birthDate: '',
      phoneNumber: '',
      email: ''
    };
    
    // 延遲導航，讓用戶看到成功訊息
    setTimeout(() => {
      // 導航到登入頁
      console.log('正在導航到登入頁...');
      router.push('/login');
    }, 2000);
  } catch (error) {
    console.error('註冊處理錯誤:', error);
    
    // 根據錯誤類型顯示適當訊息
    let errorMsg = '註冊失敗，請稍後再試';
    
    if (error.response) {
      errorMsg = error.response.data?.message || errorMsg;
      console.log('伺服器回應:', error.response.data);
    } else if (error.message) {
      errorMsg = error.message;
    }
    
    showAlert('error', errorMsg, '錯誤');
  } finally {
    console.groupEnd();
    loading.value = false;
  }
};

// 導航
const goToHome = () => router.push('/');
const goToLogin = () => router.push('/login');

// ===================== 測試功能開始 =====================
/**
 * 測試註冊功能
 * 
 * 功能說明：
 * 1. 使用固定的測試數據進行註冊
 * 2. 模擬註冊成功的情況
 * 3. 自動導向到登入頁面
 * 
 * 測試數據：
 * - 一般用戶：
 *   username: 'test_user'
 *   password: 'test123'
 *   name: '測試用戶'
 *   date: '2000-01-01'
 *   email: 'test@example.com'
 *   phone: '0912345678'
 *   role: 'user'
 * 
 * - 醫療機構：
 *   username: 'test_hospital'
 *   password: 'test123'
 *   name: '測試醫院'
 *   date: '2000-01-01'
 *   email: 'hospital@example.com'
 *   phone: '0912345678'
 *   role: 'medical'
 * 
 * - 其他用戶：
 *   username: 'test_other'
 *   password: 'test123'
 *   name: '測試其他'
 *   date: '2000-01-01'
 *   email: 'other@example.com'
 *   phone: '0912345678'
 *   role: 'other'
 */
const handleTestRegister = async () => {
  try {
    console.group('測試註冊');
    
    // 根據選擇的角色使用對應的測試數據
    let testData = {
      username: 'test_user',
      password: 'test123',
      name: '測試用戶',
      date: '2000-01-01',
      email: 'test@example.com',
      phone: '0912345678',
      role: 'user'
    };
    
    if (registerForm.value.selectedRole === 'medical') {
      testData = {
        username: 'test_hospital',
        password: 'test123',
        name: '測試醫院',
        date: '2000-01-01',
        email: 'hospital@example.com',
        phone: '0912345678',
        role: 'medical'
      };
    } else if (registerForm.value.selectedRole === 'other') {
      testData = {
        username: 'test_other',
        password: 'test123',
        name: '測試其他',
        date: '2000-01-01',
        email: 'other@example.com',
        phone: '0912345678',
        role: 'other'
      };
    }
    
    console.log('使用測試數據:', testData);
    
    // 模擬註冊成功
    await authStore.register(testData);
    
    // 註冊成功顯示
    showAlert('success', '測試註冊成功！即將為您導向登入頁面', '成功');
    console.log('測試註冊成功！');
    
    // 延遲導航，讓用戶看到成功訊息
    setTimeout(() => {
      console.log('正在導航到登入頁...');
      router.push('/login');
    }, 2000);
  } catch (error) {
    console.error('測試註冊處理錯誤:', error);
    showAlert('error', '測試註冊失敗', '錯誤');
  } finally {
    console.groupEnd();
  }
};
// ===================== 測試功能結束 =====================
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