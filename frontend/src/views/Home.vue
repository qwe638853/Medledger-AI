<template>
  <div>
    <!-- Hero Section -->
    <section class="hero-section">
      <v-container class="py-16">
        <v-row align="center">
          <v-col cols="12" md="6">
            <h1 class="display-1 font-weight-bold mb-6 hero-title-xl">
              健康檢查數據平台
            </h1>
            <p class="text-h5 mb-8 hero-subtitle">
              安全、簡單地管理您的健康數據，守護您的健康生活
            </p>
          </v-col>
          <v-col cols="12" md="6">
            <v-img
              :src="require('@/assets/hero-image.jpg')"
              class="rounded-lg hero-img"
              max-height="350"
              contain
              :lazy-src="require('@/assets/placeholder.jpg')"
            />
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Features Section -->
    <section class="features-section">
      <v-container class="py-12">
        <v-row justify="center" align="stretch">
          <v-col cols="12" sm="6" md="3" v-for="feature in features" :key="feature.title">
            <v-card class="pa-8 text-center feature-card" elevation="2">
              <v-icon size="56" color="info" class="mb-4">{{ feature.icon }}</v-icon>
              <h2 class="text-h5 font-weight-bold mb-3 feature-title">{{ feature.title }}</h2>
              <p class="text-h6 feature-desc">{{ feature.desc }}</p>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Vision/Mission Section -->
    <section class="vision-section">
      <v-container class="py-12">
        <v-row justify="center">
          <v-col cols="12" md="8">
            <v-card class="pa-8 text-center vision-card" elevation="1">
              <h2 class="text-h4 font-weight-bold mb-4 vision-title">平台願景與任務</h2>
              <p class="text-h6 vision-desc mb-0">
                我們致力於打造一個安全、易用的健康數據管理平台，協助每一位用戶輕鬆管理個人健康資訊，促進健康生活，並推動台灣健康大數據的永續發展。
              </p>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <v-container class="py-12 text-center">
        <h2 class="text-h4 font-weight-bold mb-6 cta-title">立即加入我們，開始守護您的健康！</h2>
        <v-btn color="info" class="cta-btn mr-4" size="x-large" @click="openRegisterDialog">註冊</v-btn>
        <v-btn color="primary" class="cta-btn" size="x-large" :to="{ path: '/login' }">登入</v-btn>

        <!-- 註冊對話框 -->
        <v-dialog v-model="registerDialog" max-width="500px">
          <v-card>
            <v-toolbar color="primary" dark flat>
              <v-toolbar-title>註冊</v-toolbar-title>
            </v-toolbar>
            <v-card-text class="pt-6">
              <v-form ref="registerForm" v-model="valid" @submit.prevent="handleRegister">
                <v-select
                  v-model="registerForm.selectedRole"
                  :items="roles"
                  label="選擇角色"
                  outlined
                  prepend-icon="mdi-account-group"
                  :rules="[rules.required]"
                  class="mb-4"
                />
                <v-text-field
                  v-model="registerForm.idNumber"
                  label="身分證號/員工編號 (選填)"
                  prepend-icon="mdi-account"
                  type="text"
                  outlined
                  class="mb-4"
                />
                <v-text-field
                  v-model="registerForm.password"
                  label="密碼"
                  prepend-icon="mdi-lock"
                  :type="showPassword ? 'text' : 'password'"
                  :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                  @click:append="showPassword = !showPassword"
                  outlined
                  :rules="[rules.required, rules.minLength]"
                  class="mb-4"
                />
                <v-text-field
                  v-model="registerForm.confirmPassword"
                  label="確認密碼"
                  prepend-icon="mdi-lock"
                  type="password"
                  outlined
                  :rules="[rules.required, rules.matchPassword]"
                  class="mb-4"
                />
                <v-text-field
                  v-model="registerForm.fullName"
                  label="全名"
                  prepend-icon="mdi-account-box"
                  type="text"
                  outlined
                  :rules="[rules.required]"
                  class="mb-4"
                  placeholder="請輸入您的全名"
                />
                <v-select
                  v-model="registerForm.gender"
                  :items="['男', '女', '其他', '不願透露']"
                  label="性別"
                  outlined
                  :rules="[rules.required]"
                  class="mb-4"
                />
                <v-text-field
                  v-model="registerForm.birthDate"
                  label="出生日期"
                  prepend-icon="mdi-calendar"
                  type="date"
                  outlined
                  :rules="[rules.required, rules.validDate]"
                  class="mb-4"
                />
                <v-text-field
                  v-model="registerForm.phoneNumber"
                  label="電話號碼"
                  prepend-icon="mdi-phone"
                  type="tel"
                  outlined
                  :rules="[rules.required, rules.phone]"
                  class="mb-4"
                />
                <v-text-field
                  v-model="registerForm.email"
                  label="電子郵件"
                  prepend-icon="mdi-email"
                  type="email"
                  outlined
                  :rules="[rules.required, rules.email]"
                  class="mb-4"
                />
                <v-btn
                  :loading="loading"
                  color="primary"
                  block
                  type="submit"
                  class="mt-6"
                  :disabled="!valid"
                >
                  提交註冊
                </v-btn>
                <v-btn
                  text
                  @click="resetAndCloseDialog"
                  class="mt-2"
                >
                  取消
                </v-btn>
                <v-alert
                  v-if="errorMessage"
                  type="error"
                  class="mt-4"
                >
                  {{ errorMessage }}
                </v-alert>
              </v-form>
            </v-card-text>
          </v-card>
        </v-dialog>
      </v-container>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuth } from '../composables/useAuth';

const props = defineProps({
  userRole: { type: String, default: null },
  isLoggedIn: { type: Boolean, default: false },
  currentUser: { type: String, default: '' },
  showLoginForm: { type: Boolean, default: false },
  showFooter: { type: Boolean, default: false },
  menuItems: { type: Array, default: () => [] }
});

const emit = defineEmits([
  'login', 'logout', 'forgot-password', 'register', 'go-home', 'toggle-login-form', 'toggle-footer', 'show-snackbar'
]);

const features = ref([
  { icon: 'mdi-shield-check', title: '數據安全', desc: '採用高等級加密技術，確保您的健康數據安全無虞。' },
  { icon: 'mdi-chart-line', title: '即時追蹤', desc: '隨時查看與追蹤您的健康檢查數據，掌握健康趨勢。' },
  { icon: 'mdi-brain', title: '智能分析', desc: '運用智能分析工具，協助您理解健康數據，預防疾病。' },
  { icon: 'mdi-account-circle', title: '簡易操作', desc: '介面簡單明瞭，適合所有年齡層輕鬆上手。' }
]);

const registerDialog = ref(false);
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
const errorMessage = ref(''); // 添加錯誤訊息狀態

const roles = ['一般用戶', '醫療機構', '其他機構'];

const rules = {
  required: value => !!value || '此欄位為必填',
  minLength: value => (value && value.length >= 6) || '密碼至少需要6個字符',
  matchPassword: value => value === registerForm.value.password || '密碼不一致',
  email: value => /.+@.+\..+/.test(value) || '請輸入有效的電子郵件',
  phone: value => /^\d{10}$/.test(value) || '請輸入有效的10位電話號碼',
  validDate: value => !value || new Date(value) <= new Date() || '出生日期不得為未來日期'
};

const { register } = useAuth();

const openRegisterDialog = () => {
  registerDialog.value = true;
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
  errorMessage.value = ''; // 清空錯誤訊息
};

const resetAndCloseDialog = () => {
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
  errorMessage.value = '';
  registerDialog.value = false;
};

const handleRegister = async () => {
  const { valid: isValid } = await form.value.validate();
  if (!isValid) {
    emit('show-snackbar', '請修正表單中的錯誤', 'error');
    return;
  }
  loading.value = true;
  errorMessage.value = ''; // 清空之前的錯誤訊息
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
    emit('show-snackbar', '註冊成功！請登入', 'success');
    // 不關閉對話框，允許連續註冊
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
    }; // 清空表單以便下次使用
  } catch (error) {
    errorMessage.value = error.response?.data?.detail || '註冊失敗，請稍後再試';
    emit('show-snackbar', errorMessage.value, 'error');
  } finally {
    loading.value = false;
  }
};
</script>

<style scoped>
.hero-section {
  background: linear-gradient(135deg, #e3f2fd 60%, #ffffff 100%);
  animation: fadeIn 1s ease-in-out;
}
.hero-title-xl {
  font-size: 4rem;
  letter-spacing: 2px;
  animation: slideInLeft 0.8s ease-out;
}
.hero-subtitle {
  font-size: 1.5rem;
  color: #1976d2;
  animation: slideInRight 0.8s ease-out;
}
.hero-img {
  box-shadow: 0 4px 24px 0 rgba(33, 150, 243, 0.15);
  animation: zoomIn 1s ease-out;
}
.features-section {
  background-color: #f4faff;
}
.feature-card {
  border-radius: 18px;
  background: #fff;
  box-shadow: 0 2px 8px 0 rgba(33, 150, 243, 0.08);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}
.feature-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 24px 0 rgba(33, 150, 243, 0.18);
}
.feature-title {
  font-size: 1.5rem;
}
.feature-desc {
  font-size: 1.15rem;
  color: #1976d2;
}
.vision-section {
  background: #e3f2fd;
}
.vision-card {
  border-radius: 18px;
  background: #fff;
  box-shadow: 0 2px 8px 0 rgba(33, 150, 243, 0.08);
  animation: fadeIn 1s ease-in-out;
}
.vision-title {
  font-size: 1.5rem;
}
.vision-desc {
  font-size: 1.15rem;
  color: #1976d2;
}
.cta-section {
  background: linear-gradient(135deg, #e3f2fd 60%, #ffffff 100%);
}
.cta-title {
  font-size: 1.5rem;
}
.cta-btn {
  font-size: 1.25rem;
  font-weight: bold;
  border-radius: 8px;
  transition: transform 0.2s ease;
}
.cta-btn:hover {
  transform: scale(1.05);
}

.v-form {
  padding: 16px;
}
.mb-4 {
  margin-bottom: 24px !important;
}

@keyframes fadeIn {
  from { opacity: 0; }
  to { opacity: 1; }
}
@keyframes slideInLeft {
  from { transform: translateX(-50px); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}
@keyframes slideInRight {
  from { transform: translateX(50px); opacity: 0; }
  to { transform: translateX(0); opacity: 1; }
}
@keyframes zoomIn {
  from { transform: scale(0.8); opacity: 0; }
  to { transform: scale(1); opacity: 1; }
}
</style>