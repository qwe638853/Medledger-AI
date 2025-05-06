<script setup>
import { ref, watch } from 'vue';

const username = ref('');  // 實際上是 id_number
const password = ref('');
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

const emit = defineEmits(['login', 'forgot-password', 'register']);

const roles = [
    { text: '一般用戶', value: 'user' },
    { text: '醫療機構', value: 'health_center' },
    { text: '其他機構', value: 'other' }
];

const selectRole = (role) => {
    selectedRole.value = role;
    roleSelected.value = true;
};

const goBack = () => {
    roleSelected.value = false;
    selectedRole.value = '';
    username.value = '';
    password.value = '';
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

const register = () => {
    if (!username.value || !password.value || !fullName.value || !gender.value || !birthDate.value || !phoneNumber.value || !email.value) {
        errorMessage.value = '請填寫所有必填字段';
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
</script>

<template>
    <v-container class="fill-height">
        <v-row justify="center" align="center">
            <v-col cols="12" sm="8" md="6" lg="4">
                <!-- 角色選擇區域 -->
                <v-card v-if="!roleSelected" class="pa-6">
                    <h2 class="text-h5 font-weight-bold mb-6 text-center">請選擇您的角色</h2>
                    <v-row justify="center" class="mb-4">
                        <v-col cols="12">
                            <v-btn
                                block
                                color="primary"
                                class="mb-3"
                                height="56"
                                @click="selectRole('health_center')"
                            >
                                <v-icon left>mdi-hospital-building</v-icon>
                                健康中心
                            </v-btn>
                            <v-btn
                                block
                                color="primary"
                                class="mb-3"
                                height="56"
                                @click="selectRole('user')"
                            >
                                <v-icon left>mdi-account</v-icon>
                                病人
                            </v-btn>
                            <v-btn
                                block
                                color="primary"
                                height="56"
                                @click="selectRole('other')"
                            >
                                <v-icon left>mdi-account-group</v-icon>
                                其他使用者
                            </v-btn>
                        </v-col>
                    </v-row>
                </v-card>

                <!-- 登入表單區域 -->
                <v-card v-if="roleSelected && !showForgotPassword && !showRegister" class="pa-6">
                    <div class="d-flex align-center mb-6">
                        <v-btn icon class="mr-2" @click="goBack">
                            <v-icon>mdi-arrow-left</v-icon>
                        </v-btn>
                        <h2 class="text-h5 font-weight-bold mb-0">
                            {{ selectedRole === 'health_center' ? '健康中心登入' : selectedRole === 'user' ? '病人登入' : '其他使用者登入' }}
                        </h2>
                    </div>

                    <v-form @submit.prevent="handleSubmit">
                        <v-text-field
                            v-model="formData.username"
                            label="身分證號/員工編號"
                            required
                            :rules="[v => !!v || '請輸入身分證號/員工編號']"
                            outlined
                            dense
                        ></v-text-field>

                        <v-text-field
                            v-model="formData.password"
                            label="密碼"
                            type="password"
                            required
                            :rules="[v => !!v || '請輸入密碼']"
                            outlined
                            dense
                        ></v-text-field>

                        <v-btn
                            color="primary"
                            block
                            x-large
                            type="submit"
                            class="mb-4"
                            :loading="loading"
                            height="44"
                        >
                            登入
                        </v-btn>

                        <div class="d-flex justify-space-between align-center mb-4">
                            <v-btn
                                text
                                color="primary"
                                @click="showForgotPassword = true"
                                small
                            >
                                忘記密碼？
                            </v-btn>
                            <v-btn
                                text
                                color="primary"
                                @click="showRegister = true"
                                small
                            >
                                註冊新帳號
                            </v-btn>
                        </div>
                    </v-form>
                </v-card>

                <!-- 忘記密碼表單 -->
                <v-card v-if="showForgotPassword" class="pa-6">
                    <div class="d-flex align-center mb-6">
                        <v-btn icon class="mr-2" @click="showForgotPassword = false">
                            <v-icon>mdi-arrow-left</v-icon>
                        </v-btn>
                        <h2 class="text-h5 font-weight-bold mb-0">忘記密碼</h2>
                    </div>

                    <v-form @submit.prevent="forgotPassword">
                        <v-text-field
                            v-model="username"
                            label="身分證號/員工編號"
                            required
                            outlined
                            dense
                        ></v-text-field>

                        <v-btn
                            color="primary"
                            block
                            x-large
                            type="submit"
                            height="44"
                        >
                            提交
                        </v-btn>
                    </v-form>
                </v-card>

                <!-- 註冊表單 -->
                <v-card v-if="showRegister" class="pa-6">
                    <div class="d-flex align-center mb-6">
                        <v-btn icon class="mr-2" @click="showRegister = false">
                            <v-icon>mdi-arrow-left</v-icon>
                        </v-btn>
                        <h2 class="text-h5 font-weight-bold mb-0">註冊新帳號</h2>
                    </div>

                    <v-form @submit.prevent="register">
                        <v-text-field
                            v-model="fullName"
                            label="姓名"
                            required
                            outlined
                            dense
                        ></v-text-field>

                        <v-select
                            v-model="gender"
                            :items="['男', '女', '其他']"
                            label="性別"
                            required
                            outlined
                            dense
                        ></v-select>

                        <v-text-field
                            v-model="birthDate"
                            label="出生日期"
                            type="date"
                            required
                            outlined
                            dense
                        ></v-text-field>

                        <v-text-field
                            v-model="username"
                            label="身分證號/員工編號"
                            required
                            outlined
                            dense
                        ></v-text-field>

                        <v-text-field
                            v-model="password"
                            label="密碼"
                            type="password"
                            required
                            outlined
                            dense
                        ></v-text-field>

                        <v-text-field
                            v-model="phoneNumber"
                            label="電話號碼"
                            required
                            outlined
                            dense
                        ></v-text-field>

                        <v-text-field
                            v-model="email"
                            label="電子郵件"
                            type="email"
                            required
                            outlined
                            dense
                        ></v-text-field>

                        <v-btn
                            color="primary"
                            block
                            x-large
                            type="submit"
                            height="44"
                        >
                            註冊
                        </v-btn>
                    </v-form>
                </v-card>

                <v-snackbar
                    v-model="showError"
                    color="error"
                    timeout="3000"
                    @update:model-value="errorMessage = ''"
                >
                    {{ errorMessage }}
                </v-snackbar>
            </v-col>
        </v-row>
    </v-container>
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
</style>