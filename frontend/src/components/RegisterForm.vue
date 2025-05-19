<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="10" md="8" lg="7" xl="6">
        <v-card class="elevation-3 register-card">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title class="text-h5 font-weight-bold">企業帳號註冊</v-toolbar-title>
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
            <!-- 自訂步驟條 -->
            <div class="custom-stepper">
              <div
                v-for="(step, idx) in steps"
                :key="step.value"
                class="step-btn"
                :class="{
                  active: currentStep === String(step.value),
                  completed: Number(currentStep) > step.value
                }"
              >
                <span class="step-circle">{{ step.value }}</span>
                <span class="step-title">{{ step.title }}</span>
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
                  <v-radio
                    v-for="role in roles"
                    :key="role.value"
                    :label="role.text"
                    :value="role.value"
                    color="primary"
                  >
                    <template v-slot:label>
                      <div class="d-flex align-center">
                        <v-icon :color="role.color" class="me-2">{{ role.icon }}</v-icon>
                        <span>{{ role.text }}</span>
                      </div>
                    </template>
                  </v-radio>
                </v-radio-group>
                <div class="mt-6 text-center">
                  <v-btn
                    color="primary"
                    variant="flat"
                    size="large"
                    width="180"
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
                    <v-btn variant="tonal" @click="previousStep" prepend-icon="mdi-arrow-left">返回</v-btn>
                    <v-btn color="primary" variant="flat" @click="validateAndGoNext" :disabled="!basicFormValid">
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
                    <v-btn variant="tonal" @click="previousStep" prepend-icon="mdi-arrow-left">返回</v-btn>
                    <v-btn color="primary" variant="flat" @click="validateAndGoNextAccount" :disabled="!accountFormValid">
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
                  <h3 class="text-subtitle-1 text-primary font-weight-bold mb-0">
                    <v-icon color="primary" class="me-2">mdi-credit-card</v-icon>
                    身分證上傳
                  </h3>
                  <div class="text-caption text-grey mt-1">請上傳身分證正反面照片（JPG/PNG, 5MB以內）</div>
                  <v-divider class="mt-2"></v-divider>
                </div>
                <v-row>
                  <v-col cols="12" md="6">
                    <v-file-input
                      v-model="registerForm.idCardFront"
                      label="上傳身分證正面"
                      accept="image/jpeg, image/png"
                      prepend-icon="mdi-upload"
                      show-size
                      :rules="[v => !v || (v && ['image/jpeg','image/png'].includes(v.type)) || '只接受 JPG/PNG']"
                      @change="onFileInputChange('idCardFront')"
                      outlined
                      dense
                    ></v-file-input>
                    <div v-if="registerForm.idCardFront">
                      <div class="mt-2 text-caption">已選檔案：{{ registerForm.idCardFront.name }}</div>
                      <v-img
                        :src="frontPreviewUrl"
                        aspect-ratio="1.6"
                        max-height="180"
                        class="mt-2 id-card-preview"
                        cover
                      ></v-img>
                    </div>
                  </v-col>
                  <v-col cols="12" md="6">
                    <v-file-input
                      v-model="registerForm.idCardBack"
                      label="上傳身分證反面"
                      accept="image/jpeg, image/png"
                      prepend-icon="mdi-upload"
                      show-size
                      :rules="[v => !v || (v && ['image/jpeg','image/png'].includes(v.type)) || '只接受 JPG/PNG']"
                      @change="onFileInputChange('idCardBack')"
                      outlined
                      dense
                    ></v-file-input>
                    <div v-if="registerForm.idCardBack">
                      <div class="mt-2 text-caption">已選檔案：{{ registerForm.idCardBack.name }}</div>
                      <v-img
                        :src="backPreviewUrl"
                        aspect-ratio="1.6"
                        max-height="180"
                        class="mt-2 id-card-preview"
                        cover
                      ></v-img>
                    </div>
                  </v-col>
                </v-row>
                <div class="mt-6 d-flex justify-space-between">
                  <v-btn variant="tonal" @click="previousStep" prepend-icon="mdi-arrow-left">上一步</v-btn>
                  <v-btn color="success" variant="flat" :loading="loading" @click="handleRegister" :disabled="!registerForm.idCardFront || !registerForm.idCardBack">
                    完成註冊
                    <v-icon class="ms-2">mdi-check</v-icon>
                  </v-btn>
                </div>
              </v-card>
            </div>
            
            <!-- 導航按鈕 -->
            <div class="d-flex justify-space-between mb-4 mt-4">
              <v-btn text color="primary" @click="goToHome" prepend-icon="mdi-home">
                返回首頁
              </v-btn>
              <v-btn text color="info" @click="goToLogin" prepend-icon="mdi-login">
                已有帳號？登入
              </v-btn>
            </div>
            
            <!-- 測試按鈕 -->
            <v-btn
              color="info"
              block
              @click="handleTestRegister"
              variant="outlined"
              class="mb-4"
              prepend-icon="mdi-test-tube"
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
    color: 'blue'
  },
  { 
    text: '醫療機構',
    value: 'medical',
    icon: 'mdi-hospital-building',
    color: 'green'
  },
  { 
    text: '保險業者',
    value: 'insurer',
    icon: 'mdi-shield-account',
    color: 'purple'
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
function goToStep(val) {
  // 可加驗證，這裡直接切換
  currentStep.value = String(val);
}
</script>

<style scoped>
.fill-height {
  min-height: calc(100vh - 64px);
  background-color: #f8f9fa;
}

.register-card {
  max-width: 720px;
  margin: auto;
}

.register-stepper {
  min-width: 0;
}

:deep(.v-stepper-header) {
  display: flex !important;
  flex-direction: row !important;
  flex-wrap: nowrap !important;
  min-width: 0;
  overflow-x: hidden;
  width: 100%;
}
:deep(.v-stepper-item) {
  flex: 1 1 0%;
  min-width: 0;
  max-width: 100%;
  text-align: center;
}
:deep(.v-divider) {
  margin: 0 0.5rem;
}

@media (max-width: 900px) {
  .register-card {
    max-width: 98vw;
  }
  :deep(.v-stepper-header) {
    font-size: 15px;
  }
}
@media (max-width: 600px) {
  .register-card {
    max-width: 100vw;
    border-radius: 0 !important;
  }
  :deep(.v-stepper-header) {
    font-size: 13px;
    padding-left: 0.5rem;
    padding-right: 0.5rem;
  }
  :deep(.v-stepper-item) {
    font-size: 13px;
    padding-left: 0.1rem;
    padding-right: 0.1rem;
  }
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

/* 統一按鈕高度與間距 */
.v-btn {
  letter-spacing: 0.5px;
  transition: all 0.2s ease;
}

.v-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

/* Stepper 樣式優化 */
:deep(.v-stepper) {
  box-shadow: none !important;
  border: none !important;
}

:deep(.v-stepper-header) {
  box-shadow: none !important;
}

.custom-stepper {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 32px;
  gap: 0.5rem;
  flex-wrap: wrap;
}
.step-btn {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex: 1 1 0;
  background: #f5f7fa;
  border: 2px solid #b0bec5;
  border-radius: 12px;
  padding: 12px 8px 8px 8px;
  transition: border 0.2s, background 0.2s;
  min-width: 80px;
  max-width: 160px;
  margin: 0 2px;
  user-select: none;
}
.step-btn .step-circle {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  background: #b0bec5;
  color: #fff;
  font-weight: bold;
  font-size: 18px;
  margin-bottom: 6px;
  transition: background 0.2s;
}
.step-btn .step-title {
  font-size: 15px;
  color: #607d8b;
  font-weight: 500;
  text-align: center;
  white-space: nowrap;
}
.step-btn.active {
  border-color: #1976d2;
  background: #e3f2fd;
}
.step-btn.active .step-circle {
  background: #1976d2;
}
.step-btn.completed {
  border-color: #4caf50;
  background: #e8f5e9;
}
.step-btn.completed .step-circle {
  background: #4caf50;
}
@media (max-width: 700px) {
  .custom-stepper {
    gap: 0.2rem;
  }
  .step-btn {
    min-width: 60px;
    max-width: 100px;
    padding: 8px 2px 6px 2px;
  }
  .step-btn .step-title {
    font-size: 13px;
  }
  .step-btn .step-circle {
    width: 26px;
    height: 26px;
    font-size: 15px;
  }
}
</style>