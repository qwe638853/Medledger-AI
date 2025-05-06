<script setup>
import { ref, watch } from 'vue';
import Home from './views/Home.vue'
import LoginForm from './components/LoginForm.vue'
import { useAuth } from './composables/useAuth'

const username = ref('');  // 實際上是 id_number
const password = ref('');
const confirmPassword = ref('');  // 新增密碼確認字段
const showPassword = ref(false);  // 新增密碼顯示控制
const showConfirmPassword = ref(false);  // 新增確認密碼顯示控制
const fullName = ref('');
const gender = ref('');
const birthDate = ref('');
const phoneNumber = ref('');
const email = ref('');
const errorMessage = ref('');
const showError = ref(false);  // 新增控制 snackbar 顯示的變數
const roleSelected = ref(false);
const selectedRole = ref('');
const showForgotPassword = ref(false);
const showRegister = ref(false);
const loading = ref(false);
const formData = ref({
    username: '',
    password: '',
    role: ''
});

const emit = defineEmits(['login', 'forgot-password', 'register', 'go-home', 'show-register', 'show-login']);

const roles = [
    { text: '一般用戶', value: 'user' },
    { text: '醫療機構', value: 'health_center' },
    { text: '其他機構', value: 'other' }
];

const { login, register } = useAuth()

const selectRole = (role) => {
    selectedRole.value = role;
    roleSelected.value = true;
};

const goBack = () => {
    roleSelected.value = false;
    selectedRole.value = '';
    username.value = '';
    password.value = '';
    confirmPassword.value = '';
    fullName.value = '';
    gender.value = '';
    birthDate.value = '';
    phoneNumber.value = '';
    email.value = '';
    errorMessage.value = '';
    showForgotPassword.value = false;
    showRegister.value = false;
};

const handleSubmit = async () => {
    loading.value = true;
    try {
        await emit('login', formData.value);
    } catch (error) {
        console.error('Login error:', error);
    } finally {
        loading.value = false;
    }
};

const forgotPassword = () => {
    if (!username.value) {
        errorMessage.value = '請輸入身分證號/員工編號';
        return;
    }
    if (!selectedRole.value) {
        errorMessage.value = '請選擇角色';
        return;
    }
    emit('forgot-password', {
        id_number: username.value,
        role: selectedRole.value
    });
    username.value = '';
    errorMessage.value = '';
};

const handleRegister = () => {
    if (!username.value || !password.value || !confirmPassword.value || !fullName.value || !gender.value || !birthDate.value || !phoneNumber.value || !email.value) {
        errorMessage.value = '請填寫所有必填字段';
        return;
    }
    if (password.value !== confirmPassword.value) {
        errorMessage.value = '兩次輸入的密碼不一致';
        return;
    }
    if (!selectedRole.value) {
        errorMessage.value = '請選擇角色';
        return;
    }
    emit('register', {
        id_number: username.value,
        password: password.value,
        full_name: fullName.value,
        gender: gender.value,
        birth_date: birthDate.value,
        phone_number: phoneNumber.value,
        email: email.value,
        role: selectedRole.value
    });
    username.value = '';
    password.value = '';
    confirmPassword.value = '';  // 清空確認密碼
    fullName.value = '';
    gender.value = '';
    birthDate.value = '';
    phoneNumber.value = '';
    email.value = '';
    errorMessage.value = '';
};

// 監聽 errorMessage 的變化
watch(errorMessage, (newValue) => {
    if (newValue) {
        showError.value = true;
    }
});

const showLoginForm = ref(false)
const showRegisterForm = ref(false)

const openLogin = () => {
  showLoginForm.value = true
  showRegisterForm.value = false
}
const openRegister = () => {
  showRegisterForm.value = true
  showLoginForm.value = false
}
const closeAll = () => {
  showLoginForm.value = false
  showRegisterForm.value = false
}
</script>

<template>
  <Home
    @show-login="openLogin"
    @show-register="openRegister"
  />
  <LoginForm
    v-if="showLoginForm"
    @close="closeAll"
    @go-register="openRegister"
    @login="login"
  />
  <RegisterForm
    v-if="showRegisterForm"
    @close="closeAll"
    @go-login="openLogin"
    @register="handleRegister"
  />
</template>

<style scoped>
.fill-height {
    min-height: calc(100vh - 64px);
}

.v-card {
    border-radius: 8px;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
}

.v-btn {
    text-transform: none;
    letter-spacing: 0;
}

.flex-grow-1 {
    flex-grow: 1;
}
</style>