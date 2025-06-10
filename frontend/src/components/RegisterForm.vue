<template>
  <div class="auth-page">
    <v-container class="fill-height">
      <v-row align="center" justify="center">
        <v-col cols="12" sm="10" md="10" lg="10" xl="8" style="display: flex; justify-content: center;">
          <v-slide-y-transition>
            <v-card class="auth-card" elevation="0">
              <!-- 頂部標題區 -->
              <div class="header-section">
                <h1 class="header-title">註冊</h1>
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
              
              <v-card-text class="pt-2">
                <!-- 步驟指示器 -->
                <div class="steps-indicator mb-4">
                  <div 
                    v-for="(step, idx) in steps" 
                    :key="step.value"
                    class="step"
                    :class="{
                      'active': currentStep === String(step.value),
                      'completed': Number(currentStep) > step.value
                    }"
                  >
                    <div class="step-circle">
                      <template v-if="Number(currentStep) > step.value">
                        <v-icon size="16">mdi-check</v-icon>
                      </template>
                      <template v-else>
                        {{ step.value }}
                      </template>
                    </div>
                    <div class="step-label">{{ step.title }}</div>
                    <div v-if="idx < steps.length - 1" class="step-line"></div>
                  </div>
                </div>

                <!-- 步驟 1: 角色選擇 -->
                <div v-if="currentStep === '1'" class="step-container">
                  <RoleSelector
                    v-model="registerForm.selectedRole"
                    @next="goToStep2"
                  />
                </div>

                <!-- 步驟 2: 基本資料 -->
                <div v-if="currentStep === '2'">
                  <v-card flat class="mt-6 pa-4 rounded-lg">
                    <v-form ref="basicForm" v-model="basicFormValid" lazy-validation>
                      <div class="section-title">
                        <h3 class="text-h5 font-weight-bold">
                          <v-icon color="primary" class="me-2">mdi-card-account-details</v-icon>
                          {{ isInsurerRole ? '企業基本資料' : '個人基本資料' }}
                        </h3>
                        <div class="text-body-2 text-grey">
                          {{ isInsurerRole ? '請填寫貴公司的基本資訊' : '請填寫您的個人基本資訊' }}
                        </div>
                        <!-- 調試用：顯示當前角色 
                        <v-chip 
                          v-if="registerForm.selectedRole" 
                          class="ma-1" 
                          :color="isInsurerRole ? 'primary' : 'secondary'"
                          size="small"
                        >
                          當前角色：{{ roles.find(r => r.value === registerForm.selectedRole)?.text || registerForm.selectedRole }}
                        </v-chip>
                         -->
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
                              class="form-field"
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
                              class="form-field"
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
                              class="form-field"
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
                              class="form-field"
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
                              class="form-field"
                            ></v-text-field>
                          </v-col>
                        </v-row>
                      </template>
                      <!-- 共同欄位：聯絡資訊 -->
                      <div class="section-title mt-6">
                        <h3 class="text-h5 font-weight-bold">
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
                            class="form-field"
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
                            class="form-field"
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
                          <v-icon class="ms-2">mdi-arrow-right-circle</v-icon>
                        </v-btn>
                      </div>
                    </v-form>
                  </v-card>
                </div>

                <!-- 步驟 3: 帳號設定 -->
                <div v-if="currentStep === '3'">
                  <v-card flat class="mt-6 pa-4 rounded-lg">
                    <v-form ref="accountForm" v-model="accountFormValid" lazy-validation>
                      <div class="section-title">
                        <h3 class="text-h5 font-weight-bold">
                          <v-icon color="primary" class="me-2">mdi-account-key</v-icon>
                          帳號安全設定
                        </h3>
                        <div class="text-body-2 text-grey">請設定您的身分識別碼與密碼</div>
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
                            class="form-field"
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
                            class="form-field"
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
                            class="form-field"
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
                          <v-icon class="ms-2">mdi-arrow-right-circle</v-icon>
                        </v-btn>
                      </div>
                    </v-form>
                  </v-card>
                </div>

                <!-- 步驟 4: 身分證上傳 -->
                <div v-if="currentStep === '4'">
                  <v-card flat class="mt-6 pa-4 rounded-lg">
                    <div class="section-title">
                      <h3 class="text-h5 font-weight-bold">
                        <v-icon color="primary" class="me-2">mdi-card-account-details-outline</v-icon>
                        身分證上傳
                      </h3>
                      <div class="text-body-2 text-grey">請上傳身分證正反面照片（JPG/PNG, 5MB以內）</div>
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
                            class="form-field"
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
                            class="form-field"
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
                        <v-icon class="ms-2">mdi-check-circle</v-icon>
                      </v-btn>
                    </div>
                  </v-card>
                </div>
                
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
                    @click="goToLogin"
                  >
                    <v-icon class="me-2">mdi-login-variant</v-icon>
                    已有帳號？登入
                  </v-btn>
                </div>
                
                <!-- 測試按鈕 
                <v-btn
                  class="secondary-btn mb-4"
                  block
                  @click="handleTestRegister"
                  prepend-icon="mdi-test-tube-outline"
                >
                  測試註冊
                </v-btn>
                -->
                
                <!-- 可展開的調試資訊 
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
                -->
              </v-card-text>
            </v-card>
          </v-slide-y-transition>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '../stores/auth';
import RoleSelector from './common/RoleSelector.vue';
import '@/styles/auth-forms.css';

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
    icon: 'mdi-account'
  },
  { 
    text: '醫療機構',
    value: 'medical',
    icon: 'mdi-hospital-building'
  },
  { 
    text: '保險業者',
    value: 'insurer',
    icon: 'mdi-shield-account'
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

// 角色選擇相關
const handleRoleSelect = (roleValue) => {
  registerForm.value.selectedRole = roleValue;
};

const goToStep2 = () => {
  if (registerForm.value.selectedRole) {
    console.log('選擇的角色:', registerForm.value.selectedRole);
    console.log('是否為保險業者:', isInsurerRole.value);
    nextStep();
  }
};

// 監聽角色選擇變化
watch(() => registerForm.value.selectedRole, (newRole) => {
  console.log('角色已變更為:', newRole);
  console.log('isInsurerRole 計算屬性值:', isInsurerRole.value);
}, { immediate: true });
</script>

<style scoped>
/* 全局樣式 */
.auth-page {
  background-color: #F9F7F4;
  min-height: calc(100vh - 64px);
  padding: 1rem;
}

/* 註冊卡片 */
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
  color: #6B7280 !important;
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

/* 上傳區域樣式 */
.upload-area {
  border: 2px dashed #e5e7eb;
  border-radius: 16px;
  padding: 1.5rem;
  background: white;
  transition: all 0.3s ease;
  min-height: 120px;
}

.upload-area:hover {
  border-color: #00B8D9;
  background-color: rgba(0, 184, 217, 0.05);
}

.upload-area-active {
  border-color: #00B8D9;
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