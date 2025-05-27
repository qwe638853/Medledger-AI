<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="10" md="8" lg="7" xl="6">
        <v-card class="register-card">
          <!-- 移除藍色頂部，改為簡潔的標題 -->
          <div class="header-section">
            <h1 class="header-title">企業帳號註冊</h1>
            <p class="header-subtitle">安全、便捷的醫療數據管理解決方案</p>
          </div>
          
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
            <!-- 優化步驟條設計 -->
            <div class="stepper-container">
              <div
                v-for="(step, idx) in steps"
                :key="step.value"
                class="step-item"
                :class="{
                  'active': currentStep === String(step.value),
                  'completed': Number(currentStep) > step.value
                }"
              >
                <div class="step-circle">
                  <span class="step-number">{{ step.value }}</span>
                  <span v-if="Number(currentStep) > step.value" class="step-check">✓</span>
                </div>
                <span class="step-title">{{ step.title }}</span>
                <div v-if="idx < steps.length - 1" class="step-line"></div>
              </div>
            </div>

            <!-- 步驟 1: 角色選擇 -->
            <div v-if="currentStep === '1'">
              <v-card flat class="mt-6 pa-4 rounded-lg">
                <div class="section-title mb-4">
                  <h3 class="text-subtitle-1 text-primary font-weight-bold mb-0">
                    <v-icon color="primary" class="me-2">mdi-account-multiple</v-icon>
                    角色選擇
                  </h3>
                  <div class="text-caption text-grey mt-1">請選擇您的帳號類型</div>
                  <v-divider class="mt-2"></v-divider>
                </div>
                <v-radio-group
                  v-model="registerForm.selectedRole"
                  inline
                  class="mt-4"
                  :rules="[rules.required]"
                >
                  <div class="role-cards-container">
                    <div
                      v-for="role in roles"
                      :key="role.value"
                      class="role-card"
                      :class="{ 'role-card--selected': registerForm.selectedRole === role.value }"
                      @click="registerForm.selectedRole = role.value"
                    >
                      <div class="role-card__icon-wrapper">
                        <v-icon :color="registerForm.selectedRole === role.value ? 'primary' : 'grey'" size="32">
                          {{ role.icon }}
                        </v-icon>
                      </div>
                      <div class="role-card__content">
                        <h3 class="role-card__title">{{ role.text }}</h3>
                        <p class="role-card__description">{{ role.description }}</p>
                      </div>
                    </div>
                  </div>
                </v-radio-group>
                <div class="mt-6 text-center">
                  <v-btn
                    class="primary-btn"
                    size="large"
                    @click="nextStep"
                    :disabled="!registerForm.selectedRole"
                  >
                    下一步
                    <v-icon class="ms-2">mdi-arrow-right</v-icon>
                  </v-btn>
                </div>
              </v-card>
            </div>

            <!-- 步驟 2: 基本資料 -->
            <div v-if="currentStep === '2'">
              <v-card flat class="mt-6 pa-4 rounded-lg">
                <v-form ref="basicForm" v-model="basicFormValid" lazy-validation>
                  <div class="section-title mb-4">
                    <h3 class="text-subtitle-1 text-primary font-weight-bold mb-0">
                      <v-icon color="primary" class="me-2">mdi-card-account-details</v-icon>
                      {{ isInsurerRole ? '企業基本資料' : '個人基本資料' }}
                    </h3>
                    <div class="text-caption text-grey mt-1">
                      {{ isInsurerRole ? '請填寫貴公司的基本資訊' : '請填寫您的個人基本資訊' }}
                    </div>
                    <v-divider class="mt-2"></v-divider>
                  </div>
                  <!-- 保險業者專屬欄位 -->
                  <template v-if="isInsurerRole">
                    <v-row>
                      <v-col cols="12">
                        <v-text-field
                          v-model="registerForm.companyName"
                          label="公司名稱"
                          variant="outlined"
                          density="comfortable"
                          prepend-inner-icon="mdi-office-building"
                          :rules="[rules.required]"
                          hint="請輸入完整的公司登記名稱"
                          persistent-hint
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12">
                        <v-text-field
                          v-model="registerForm.contactPerson"
                          label="聯絡人姓名"
                          variant="outlined"
                          density="comfortable"
                          prepend-inner-icon="mdi-account"
                          :rules="[rules.required]"
                        ></v-text-field>
                      </v-col>
                    </v-row>
                  </template>
                  <!-- 一般用戶/醫療機構欄位 -->
                  <template v-else>
                    <v-row>
                      <v-col cols="12">
                        <v-text-field
                          v-model="registerForm.fullName"
                          label="姓名"
                          variant="outlined"
                          density="comfortable"
                          prepend-inner-icon="mdi-account-box"
                          :rules="[rules.required]"
                        ></v-text-field>
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-select
                          v-model="registerForm.gender"
                          :items="genders"
                          label="性別"
                          variant="outlined"
                          density="comfortable"
                          prepend-inner-icon="mdi-gender-male-female"
                          :rules="[rules.required]"
                        ></v-select>
                      </v-col>
                      <v-col cols="12" md="6">
                        <v-text-field
                          v-model="registerForm.birthDate"
                          label="出生日期"
                          variant="outlined"
                          density="comfortable"
                          prepend-inner-icon="mdi-calendar"
                          type="date"
                          :rules="[rules.required, rules.validDate]"
                        ></v-text-field>
                      </v-col>
                    </v-row>
                  </template>
                  <!-- 共同欄位：聯絡資訊 -->
                  <div class="section-title mt-6 mb-4">
                    <h3 class="text-subtitle-1 text-primary font-weight-bold mb-0">
                      <v-icon color="primary" class="me-2">mdi-contacts</v-icon>
                      聯絡資訊
                    </h3>
                    <v-divider class="mt-2"></v-divider>
                  </div>
                  <v-row>
                    <v-col cols="12" md="6">
                      <v-text-field
                        v-model="registerForm.phoneNumber"
                        label="電話號碼"
                        variant="outlined"
                        density="comfortable"
                        prepend-inner-icon="mdi-phone"
                        type="tel"
                        :rules="[rules.required, rules.phone]"
                        placeholder="例如：0912345678"
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12" md="6">
                      <v-text-field
                        v-model="registerForm.email"
                        label="電子郵件"
                        variant="outlined"
                        density="comfortable"
                        prepend-inner-icon="mdi-email"
                        type="email"
                        :rules="[rules.required, rules.email]"
                        placeholder="例如：example@email.com"
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  <div class="mt-6 d-flex justify-space-between">
                    <v-btn class="secondary-btn" @click="previousStep">
                      <v-icon class="me-2">mdi-arrow-left</v-icon>
                      返回
                    </v-btn>
                    <v-btn
                      class="primary-btn"
                      @click="validateAndGoNext"
                      :disabled="!basicFormValid"
                    >
                      下一步
                      <v-icon class="ms-2">mdi-arrow-right</v-icon>
                    </v-btn>
                  </div>
                </v-form>
              </v-card>
            </div>

            <!-- 步驟 3: 帳號設定 -->
            <div v-if="currentStep === '3'">
              <v-card flat class="mt-6 pa-4 rounded-lg">
                <v-form ref="accountForm" v-model="accountFormValid" lazy-validation>
                  <div class="section-title mb-4">
                    <h3 class="text-subtitle-1 text-primary font-weight-bold mb-0">
                      <v-icon color="primary" class="me-2">mdi-account-key</v-icon>
                      帳號安全設定
                    </h3>
                    <div class="text-caption text-grey mt-1">請設定您的身分識別碼與密碼</div>
                    <v-divider class="mt-2"></v-divider>
                  </div>
                  <v-row>
                    <v-col cols="12">
                      <v-text-field
                        v-model="registerForm.idNumber"
                        :label="idLabel"
                        variant="outlined"
                        density="comfortable"
                        prepend-inner-icon="mdi-identifier"
                        :rules="[rules.required, rules.idFormat]"
                        :hint="idHint"
                        persistent-hint
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12">
                      <v-text-field
                        v-model="registerForm.password"
                        label="密碼"
                        variant="outlined"
                        density="comfortable"
                        prepend-inner-icon="mdi-lock"
                        :type="showPassword ? 'text' : 'password'"
                        :append-inner-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                        @click:append-inner="showPassword = !showPassword"
                        :rules="[rules.required, rules.minLength]"
                        hint="密碼至少需要6個字符"
                        persistent-hint
                      ></v-text-field>
                    </v-col>
                    <v-col cols="12">
                      <v-text-field
                        v-model="registerForm.confirmPassword"
                        label="確認密碼"
                        variant="outlined"
                        density="comfortable"
                        prepend-inner-icon="mdi-lock-check"
                        :type="showConfirmPassword ? 'text' : 'password'"
                        :append-inner-icon="showConfirmPassword ? 'mdi-eye-off' : 'mdi-eye'"
                        @click:append-inner="showConfirmPassword = !showConfirmPassword"
                        :rules="[rules.required, rules.matchPassword]"
                        hint="請再次輸入您的密碼"
                        persistent-hint
                      ></v-text-field>
                    </v-col>
                  </v-row>
                  <div class="mt-6 d-flex justify-space-between">
                    <v-btn class="secondary-btn" @click="previousStep">
                      <v-icon class="me-2">mdi-arrow-left</v-icon>
                      返回
                    </v-btn>
                    <v-btn
                      class="primary-btn"
                      @click="validateAndGoNextAccount"
                      :disabled="!accountFormValid"
                    >
                      下一步
                      <v-icon class="ms-2">mdi-arrow-right</v-icon>
                    </v-btn>
                  </div>
                </v-form>
              </v-card>
            </div>

            <!-- 步驟 4: 身分證上傳 -->
            <div v-if="currentStep === '4'">
              <v-card flat class="mt-6 pa-4 rounded-lg">
                <div class="section-title mb-4">
                  <h3 class="text-subtitle-1 font-weight-bold mb-0">
                    <v-icon class="me-2">mdi-card-account-details-outline</v-icon>
                    身分證上傳
                  </h3>
                  <div class="text-caption text-grey mt-1">請上傳身分證正反面照片（JPG/PNG, 5MB以內）</div>
                  <v-divider class="mt-2"></v-divider>
                </div>
                <v-row>
                  <v-col cols="12" md="6">
                    <div class="upload-area" :class="{ 'upload-area-active': registerForm.idCardFront }">
                      <v-file-input
                        v-model="registerForm.idCardFront"
                        label="上傳身分證正面"
                        accept="image/jpeg, image/png"
                        prepend-icon="mdi-upload-outline"
                        show-size
                        :rules="[v => !v || (v && ['image/jpeg','image/png'].includes(v.type)) || '只接受 JPG/PNG']"
                        @change="onFileInputChange('idCardFront')"
                        hide-details
                      ></v-file-input>
                      <div v-if="registerForm.idCardFront" class="preview-container">
                        <div class="mt-2 text-caption">已選檔案：{{ registerForm.idCardFront.name }}</div>
                        <v-img
                          :src="frontPreviewUrl"
                          aspect-ratio="1.6"
                          max-height="180"
                          class="mt-2 preview-image"
                          cover
                        ></v-img>
                      </div>
                    </div>
                  </v-col>
                  <v-col cols="12" md="6">
                    <div class="upload-area" :class="{ 'upload-area-active': registerForm.idCardBack }">
                      <v-file-input
                        v-model="registerForm.idCardBack"
                        label="上傳身分證反面"
                        accept="image/jpeg, image/png"
                        prepend-icon="mdi-upload-outline"
                        show-size
                        :rules="[v => !v || (v && ['image/jpeg','image/png'].includes(v.type)) || '只接受 JPG/PNG']"
                        @change="onFileInputChange('idCardBack')"
                        hide-details
                      ></v-file-input>
                      <div v-if="registerForm.idCardBack" class="preview-container">
                        <div class="mt-2 text-caption">已選檔案：{{ registerForm.idCardBack.name }}</div>
                        <v-img
                          :src="backPreviewUrl"
                          aspect-ratio="1.6"
                          max-height="180"
                          class="mt-2 preview-image"
                          cover
                        ></v-img>
                      </div>
                    </div>
                  </v-col>
                </v-row>
                <div class="mt-6 d-flex justify-space-between">
                  <v-btn class="secondary-btn" @click="previousStep">
                    <v-icon class="me-2">mdi-arrow-left</v-icon>
                    返回
                  </v-btn>
                  <v-btn
                    class="primary-btn"
                    :loading="loading"
                    @click="handleRegister"
                    :disabled="!registerForm.idCardFront || !registerForm.idCardBack"
                  >
                    完成註冊
                    <v-icon class="ms-2">mdi-check</v-icon>
                  </v-btn>
                </div>
              </v-card>
            </div>
            
            <!-- 導航按鈕 -->
            <div class="d-flex justify-space-between mb-4 mt-4">
              <v-btn
                class="secondary-btn"
                @click="goToHome"
              >
                <v-icon class="me-2">mdi-home-outline</v-icon>
                返回首頁
              </v-btn>
              <v-btn
                class="secondary-btn"
                @click="goToLogin"
              >
                <v-icon class="me-2">mdi-login-variant</v-icon>
                已有帳號？登入
              </v-btn>
            </div>
            
            <!-- 測試按鈕 -->
            <v-btn
              class="secondary-btn mb-4"
              block
              @click="handleTestRegister"
              prepend-icon="mdi-test-tube-outline"
            >
              測試註冊
            </v-btn>
            
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
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';

const router = useRouter();
const authStore = useAuthStore();

// 步驟控制
const currentStep = ref('1');
const basicForm = ref(null);
const accountForm = ref(null);
const basicFormValid = ref(false);
const accountFormValid = ref(false);

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
  email: '',
  // 保險業者專用欄位
  companyName: '',
  contactPerson: '',
  idCardFront: null,
  idCardBack: null
});

const loading = ref(false);
const showPassword = ref(false);
const showConfirmPassword = ref(false);

const roles = [
  { 
    text: '一般用戶',
    value: 'user',
    icon: 'mdi-account',
    color: 'blue',
    description: '個人用戶註冊，管理個人健康數據'
  },
  { 
    text: '醫療機構',
    value: 'medical',
    icon: 'mdi-hospital-building',
    color: 'green',
    description: '醫院、診所等醫療服務提供者'
  },
  { 
    text: '保險業者',
    value: 'insurer',
    icon: 'mdi-shield-account',
    color: 'purple',
    description: '保險公司、保險服務提供商'
  }
];

const genders = ['男', '女', '其他', '不願透露'];

// 角色相關計算屬性
const isInsurerRole = computed(() => {
  return registerForm.value.selectedRole === 'insurer';
});

// ID 標籤與提示
const idLabel = computed(() => {
  switch (registerForm.value.selectedRole) {
    case 'insurer':
      return '保險業者ID';
    case 'medical':
      return '醫療機構ID';
    default:
      return '身分證號碼';
  }
});

const idHint = computed(() => {
  switch (registerForm.value.selectedRole) {
    case 'insurer':
      return '請輸入您的保險業者ID（英數字組合，至少6位）';
    case 'medical':
      return '請輸入您的醫療機構ID（英數字組合，至少6位）';
    default:
      return '請輸入有效的身分證號碼（例如：A123456789）';
  }
});

// 判斷是否為開發環境
const isDevelopment = import.meta.env.MODE === 'development';

// 用於調試的數據
const debugInfo = computed(() => {
  return {
    formState: {
      ...registerForm.value,
      password: registerForm.value.password ? '******' : '',
      confirmPassword: registerForm.value.confirmPassword ? '******' : '',
      roleText: roles.find(r => r.value === registerForm.value.selectedRole)?.text || ''
    },
    validation: {
      basicFormValid: basicFormValid.value,
      accountFormValid: accountFormValid.value,
      currentStep: currentStep.value,
      loading: loading.value
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
  minLength: value => (value && value.length >= 6) || '密碼至少需要6個字符',
  matchPassword: value => value === registerForm.value.password || '密碼不一致',
  email: value => /.+@.+\..+/.test(value) || '請輸入有效的電子郵件',
  phone: value => /^\d{10}$/.test(value) || '請輸入有效的10位電話號碼',
  validDate: value => !value || new Date(value) <= new Date() || '出生日期不得為未來日期',
  idFormat: (value) => {
    if (!value) return '此欄位為必填';
    
    if (registerForm.value.selectedRole === 'user') {
      // 台灣身分證號驗證邏輯：一個英文字母加9個數字
      return /^[A-Za-z][12]\d{8}$/.test(value) || '請輸入有效的身分證號碼';
    }
    
    // 其他類型的ID至少需要6位
    return value.length >= 6 || 'ID至少需要6位字符';
  }
};

// 步驟控制功能
const nextStep = () => {
  const nextStepValue = String(Number(currentStep.value) + 1);
  currentStep.value = nextStepValue;
};

const previousStep = () => {
  const prevStepValue = String(Number(currentStep.value) - 1);
  currentStep.value = prevStepValue;
};

const validateAndGoNext = () => {
  if (basicForm.value.validate()) {
    nextStep();
  }
};

// 預覽網址
const frontPreviewUrl = ref('');
const backPreviewUrl = ref('');

function onFileInputChange(side) {
  if (side === 'idCardFront') {
    if (registerForm.value.idCardFront) {
      frontPreviewUrl.value = URL.createObjectURL(registerForm.value.idCardFront);
    } else {
      frontPreviewUrl.value = '';
    }
  } else if (side === 'idCardBack') {
    if (registerForm.value.idCardBack) {
      backPreviewUrl.value = URL.createObjectURL(registerForm.value.idCardBack);
    } else {
      backPreviewUrl.value = '';
    }
  }
}
// 清理 URL 物件（避免記憶體洩漏）
watch(() => registerForm.value.idCardFront, (newVal, oldVal) => {
  if (oldVal && frontPreviewUrl.value) URL.revokeObjectURL(frontPreviewUrl.value);
  if (newVal) frontPreviewUrl.value = URL.createObjectURL(newVal);
  else frontPreviewUrl.value = '';
});
watch(() => registerForm.value.idCardBack, (newVal, oldVal) => {
  if (oldVal && backPreviewUrl.value) URL.revokeObjectURL(backPreviewUrl.value);
  if (newVal) backPreviewUrl.value = URL.createObjectURL(newVal);
  else backPreviewUrl.value = '';
});

// 驗證與流程控制
const validateAndGoNextAccount = () => {
  if (accountForm.value.validate()) {
    nextStep();
  }
};

// 表單提交
const handleRegister = async () => {
  // 修改驗證邏輯，避免在最後步驟重新驗證前面步驟的表單
  if (currentStep.value !== '4') {
    // 如果不是最後一步，需要驗證當前表單
    if (currentStep.value === '3' && accountForm.value && !accountForm.value.validate()) {
      showAlert('error', '請確認所有欄位都已正確填寫', '表單驗證失敗');
      return;
    } else if (currentStep.value === '2' && basicForm.value && !basicForm.value.validate()) {
      showAlert('error', '請確認所有欄位都已正確填寫', '表單驗證失敗');
      return;
    }
  }
  
  // 在最後一步，檢查身份證是否已上傳
  if (currentStep.value === '4' && (!registerForm.value.idCardFront || !registerForm.value.idCardBack)) {
    showAlert('error', '請上傳身分證正反面照片', '表單驗證失敗');
    return;
  }
  
  loading.value = true;
  
  try {
    console.group('註冊請求');
    console.log('註冊角色:', registerForm.value.selectedRole);
    
    // 根據角色選擇註冊方式
    if (registerForm.value.selectedRole === 'insurer') {
      console.log('保險業者註冊:', {
        insurerId: registerForm.value.idNumber,
        companyName: registerForm.value.companyName,
        contactPerson: registerForm.value.contactPerson,
        email: registerForm.value.email,
        phone: registerForm.value.phoneNumber
      });
      
      await authStore.registerInsurer({
        insurerId: registerForm.value.idNumber,
        password: registerForm.value.password,
        companyName: registerForm.value.companyName,
        contactPerson: registerForm.value.contactPerson,
        email: registerForm.value.email,
        phone: registerForm.value.phoneNumber
      });
    } else {
      console.log('一般用戶/醫療機構註冊:', {
        username: registerForm.value.idNumber,
        name: registerForm.value.fullName,
        gender: registerForm.value.gender,
        birthDate: registerForm.value.birthDate,
        email: registerForm.value.email,
        phone: registerForm.value.phoneNumber,
        role: registerForm.value.selectedRole
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
    }
    
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
      email: '',
      companyName: '',
      contactPerson: '',
      idCardFront: null,
      idCardBack: null
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

// 測試註冊功能
const handleTestRegister = async () => {
  if (!registerForm.value.selectedRole) {
    showAlert('error', '請先選擇角色', '錯誤');
    return;
  }
  
  try {
    console.group('測試註冊');
    
    // 根據選擇的角色使用對應的測試數據
    if (registerForm.value.selectedRole === 'insurer') {
      const testData = {
        insurerId: 'test_insurer',
        password: 'test123',
        companyName: '測試保險公司',
        contactPerson: '保險聯絡人',
        email: 'insurer@example.com',
        phone: '0912345678'
      };
      
      console.log('使用測試數據 (保險業者):', testData);
      
      // 保險業者註冊
      await authStore.registerInsurer(testData);
    } else {
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
      }
      
      console.log('使用測試數據:', testData);
      
      // 一般用戶/醫療機構註冊
      await authStore.register(testData);
    }
    
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

const steps = [
  { value: 1, title: '角色選擇' },
  { value: 2, title: '基本資料' },
  { value: 3, title: '帳號設定' },
  { value: 4, title: '身分證上傳' }
];
</script>

<style scoped>
/* 全局樣式 */
.fill-height {
  min-height: calc(100vh - 64px);
  background-color: #F9F7F4;
  padding: 2rem 1rem;
}

/* 註冊卡片 */
.register-card {
  border-radius: 24px !important;
  background: white !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03) !important;
  padding: 2rem !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
}

/* 頂部標題區 */
.header-section {
  text-align: center;
  margin-bottom: 3rem;
  padding: 1rem;
}

.header-title {
  font-size: 2rem;
  font-weight: 900;
  color: #111827;
  margin-bottom: 0.5rem;
  letter-spacing: -0.5px;
}

.header-subtitle {
  font-size: 1rem;
  color: #888;
  margin: 0;
}

/* 步驟條設計 */
.stepper-container {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin: 2rem 0 3rem;
  padding: 0 1rem;
  position: relative;
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  position: relative;
  flex: 1;
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
}

.step-number {
  color: #888;
  font-weight: 600;
}

.step-check {
  color: #111827;
  font-weight: 600;
}

.step-title {
  font-size: 0.875rem;
  color: #888;
  text-align: center;
  font-weight: 500;
}

.step-line {
  position: absolute;
  top: 18px;
  right: calc(-50% + 18px);
  width: calc(100% - 36px);
  height: 2px;
  background: #e5e7eb;
  z-index: 0;
}

/* 活動步驟樣式 */
.step-item.active .step-circle {
  border-color: #F8F441;
  background: #F8F441;
}

.step-item.active .step-number {
  color: #111827;
}

.step-item.active .step-title {
  color: #111827;
  font-weight: 600;
}

/* 已完成步驟樣式 */
.step-item.completed .step-circle {
  border-color: #463F3A;
  background: #463F3A;
}

.step-item.completed .step-check {
  color: white;
}

.step-item.completed .step-line {
  background: #111827;
}

/* 表單樣式 */
:deep(.v-text-field) {
  border-radius: 16px !important;
}

:deep(.v-text-field .v-field) {
  border-radius: 16px !important;
  background: white !important;
  border: 1px solid #e5e7eb !important;
}

:deep(.v-text-field .v-field--focused) {
  border-color: #111827 !important;
}

:deep(.v-text-field .v-label) {
  color: #888 !important;
}

/* Radio 按鈕樣式 */
:deep(.v-radio) {
  margin-right: 1rem;
}

:deep(.v-radio .v-selection-control) {
  border-color: #e5e7eb !important;
}

:deep(.v-radio .v-selection-control--active) {
  color: #F8F441 !important;
  border-color: #F8F441 !important;
}

/* 按鈕樣式 */
.v-btn {
  border-radius: 16px !important;
  text-transform: none !important;
  font-weight: 600 !important;
  letter-spacing: 0 !important;
  height: 48px !important;
  min-width: 120px !important;
}

.v-btn.primary-btn {
  background-color: #F8F441 !important;
  color: #111827 !important;
  border: none !important;
}

.v-btn.primary-btn:hover {
  background-color: #f9f650 !important;
  transform: translateY(-2px);
}

.v-btn.secondary-btn {
  background-color: white !important;
  color: #111827 !important;
  border: 1px solid #e5e7eb !important;
}

.v-btn.secondary-btn:hover {
  background-color: #f9fafb !important;
  border-color: #d1d5db !important;
  transform: translateY(-2px);
}

/* 文件上傳區域 */
:deep(.v-file-input) {
  border-radius: 16px !important;
}

.upload-area {
  border: 2px dashed #e5e7eb;
  border-radius: 16px;
  padding: 1.5rem;
  background: white;
  transition: all 0.3s ease;
  min-height: 120px;
}

.upload-area:hover {
  border-color: #F8F441;
  background-color: #FEFEF5;
}

.upload-area-active {
  border-color: #111827;
  border-style: solid;
}

.preview-container {
  margin-top: 1rem;
}

.preview-image {
  border-radius: 8px;
  border: 1px solid #e5e7eb;
  overflow: hidden;
}

:deep(.v-file-input) {
  border: none;
  padding: 0;
}

:deep(.v-file-input .v-field) {
  border: none !important;
  background: transparent !important;
  box-shadow: none !important;
}

:deep(.v-file-input .v-field__field) {
  padding: 0;
}

:deep(.v-file-input .v-label) {
  opacity: 0.7;
}

/* RWD 適配 */
@media (max-width: 768px) {
  .register-card {
    padding: 1.5rem !important;
  }
  
  .header-title {
    font-size: 1.5rem;
  }
  
  .stepper-container {
    flex-direction: column;
    gap: 1.5rem;
  }
  
  .step-line {
    width: 2px;
    height: 24px;
    top: 36px;
    right: auto;
    left: 50%;
    transform: translateX(-50%);
  }
  
  .v-btn {
    width: 100%;
  }
}

/* Alert 訊息樣式 */
:deep(.v-alert) {
  border-radius: 16px !important;
  margin-bottom: 1.5rem !important;
}

/* 分隔線樣式 */
:deep(.v-divider) {
  border-color: #e5e7eb !important;
  margin: 1.5rem 0 !important;
  opacity: 0.5;
}

/* 角色選擇卡片容器 */
.role-cards-container {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  width: 100%;
  margin: 1rem 0;
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
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #f3f4f6;
  margin-right: 1rem;
  transition: all 0.3s ease;
}

.role-card--selected .role-card__icon-wrapper {
  background: var(--v-theme-primary-lighten-4);
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
</style>