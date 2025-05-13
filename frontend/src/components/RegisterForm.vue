<template>
    <v-container class="fill-height">
      <v-row align="center" justify="center">
        <v-col cols="12" sm="8" md="4">
          <v-card class="elevation-12">
            <v-toolbar color="primary" dark flat>
              <v-toolbar-title>註冊</v-toolbar-title>
            </v-toolbar>
            <v-card-text class="pt-6">
              <v-form ref="registerForm" v-model="valid" @submit.prevent="handleRegister">
                <!-- 選擇角色 -->
                <v-select
                  v-model="registerForm.selectedRole"
                  :items="roles"
                  label="選擇角色"
                  outlined
                  dense
                  prepend-icon="mdi-account-group"
                  :rules="[rules.required]"
                  class="mb-4"
                  style="width: 100%; max-width: 100%;"
                />
  
                <!-- 身分證號/員工編號 -->
                <v-text-field
                  v-model="registerForm.idNumber"
                  label="身分證號/員工編號 (選填)"
                  prepend-icon="mdi-account"
                  type="text"
                  outlined
                  dense
                  class="mb-4"
                  style="width: 100%; max-width: 100%;"
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
                  style="width: 100%; max-width: 100%;"
                />
  
                <!-- 確認密碼 -->
                <v-text-field
                  v-model="registerForm.confirmPassword"
                  label="確認密碼"
                  prepend-icon="mdi-lock"
                  type="password"
                  outlined
                  dense
                  :rules="[rules.required, rules.matchPassword]"
                  class="mb-4"
                  style="width: 100%; max-width: 100%;"
                />
  
                <!-- 全名 -->
                <v-text-field
                  v-model="registerForm.fullName"
                  label="全名"
                  prepend-icon="mdi-account-box"
                  type="text"
                  outlined
                  dense
                  :rules="[rules.required]"
                  class="mb-4"
                  placeholder="請輸入您的全名"
                  style="width: 100%; max-width: 100%;"
                />
  
                <!-- 性別 -->
                <v-select
                  v-model="registerForm.gender"
                  :items="['男', '女', '其他', '不願透露']"
                  label="性別"
                  outlined
                  dense
                  :rules="[rules.required]"
                  class="mb-4"
                  style="width: 100%; max-width: 100%;"
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
                  style="width: 100%; max-width: 100%;"
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
                  style="width: 100%; max-width: 100%;"
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
                  class="mb-4"
                  style="width: 100%; max-width: 100%;"
                />
  
                <!-- 提交按鈕 -->
                <v-btn
                  :loading="loading"
                  color="primary"
                  block
                  type="submit"
                  class="mb-4"
                  :disabled="!valid"
                >
                  提交註冊
                </v-btn>
  
                <!-- 其他導航按鈕 -->
                <v-btn text @click="goToHome" class="mr-4">
                  返回首頁
                </v-btn>
                <v-btn text @click="goToLogin">
                  已有帳號？登入
                </v-btn>
  
                <!-- 錯誤訊息 -->
                <v-alert v-if="errorMessage" type="error" class="mt-4">
                  {{ errorMessage }}
                </v-alert>
              </v-form>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  import { useRouter } from 'vue-router';
  import { useAuth } from '../composables/useAuth';
  
  const router = useRouter();
  const { register } = useAuth();
  
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
  const valid = ref(false);
  const showPassword = ref(false);
  const form = ref(null);
  const errorMessage = ref('');
  const roles = ['一般用戶', '醫療機構', '其他用戶'];
  
  const rules = {
    required: value => !!value || '此欄位為必填',
    minLength: value => (value && value.length >= 6) || '密碼至少需要6個字符',
    matchPassword: value => value === registerForm.value.password || '密碼不一致',
    email: value => /.+@.+\..+/.test(value) || '請輸入有效的電子郵件',
    phone: value => /^\d{10}$/.test(value) || '請輸入有效的10位電話號碼',
    validDate: value => !value || new Date(value) <= new Date() || '出生日期不得為未來日期'
  };
  
  const handleRegister = async () => {
    const { valid: isValid } = await form.value.validate();
    if (!isValid) {
      document.dispatchEvent(new CustomEvent('show-snackbar', {
        detail: { message: '請修正表單中的錯誤', color: 'error' }
      }));
      return;
    }
    loading.value = true;
    errorMessage.value = '';
    try {
      await register({
        id_number: registerForm.value.idNumber,
        password: registerForm.value.password,
        full_name: registerForm.value.fullName,
        gender: registerForm.value.gender,
        birth_date: registerForm.value.birthDate,
        phone_number: registerForm.value.phoneNumber,
        email: registerForm.value.email,
        role: registerForm.value.selectedRole.toLowerCase()
      });
      document.dispatchEvent(new CustomEvent('show-snackbar', {
        detail: { message: '註冊成功！請登入', color: 'success' }
      }));
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
    } catch (error) {
      errorMessage.value = error.response?.data?.detail || '註冊失敗，請稍後再試';
      document.dispatchEvent(new CustomEvent('show-snackbar', {
        detail: { message: errorMessage.value, color: 'error' }
      }));
    } finally {
      loading.value = false;
    }
  };
  
  const goToHome = () => {
    router.push('/');
  };
  
  const goToLogin = () => {
    router.push('/login');
  };
  </script>
  
  <style scoped>
  .fill-height {
    min-height: calc(100vh - 64px);
  }
  
  .v-card {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
    padding: 16px; /* 添加內邊距，確保內容不會貼邊 */
  }
  
  .v-btn {
    text-transform: none;
    letter-spacing: 0;
    transition: transform 0.2s ease;
  }
  
  .v-btn:hover {
    transform: scale(1.05);
  }
  
  .v-form {
    padding: 0; /* 移除 v-form 的內邊距，避免重複 */
  }
  
  .mb-4 {
    margin-bottom: 24px !important;
  }
  </style>c