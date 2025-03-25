<script setup>
import { ref } from 'vue';

const username = ref('');  // 實際上是 id_number
const password = ref('');
const errorMessage = ref('');
const showLogin = ref(false);
const showRegister = ref(false);
const showForgotPassword = ref(false);
const selectedRole = ref('');

const emit = defineEmits(['login', 'forgot-password', 'register']);

const goBack = () => {
    showLogin.value = false;
    showRegister.value = false;
    showForgotPassword.value = false;
    username.value = '';
    password.value = '';
    selectedRole.value = '';
    errorMessage.value = '';
};

const login = () => {
    if (!username.value || !password.value) {
        errorMessage.value = '請輸入身分證號/員工編號和密碼';
        return;
    }
    if (!selectedRole.value) {
        errorMessage.value = '請選擇角色';
        return;
    }
    emit('login', {
        username: username.value,
        password: password.value,
        role: selectedRole.value
    });
    username.value = '';
    password.value = '';
    selectedRole.value = '';
    errorMessage.value = '';
    showLogin.value = false;
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
    selectedRole.value = '';
    errorMessage.value = '';
    showForgotPassword.value = false;
    showLogin.value = true;
};

const register = () => {
    if (!username.value || !password.value) {
        errorMessage.value = '請輸入身分證號/員工編號和密碼';
        return;
    }
    if (!selectedRole.value) {
        errorMessage.value = '請選擇角色';
        return;
    }
    emit('register', {
        id_number: username.value,
        password: password.value,
        role: selectedRole.value
    });
    username.value = '';
    password.value = '';
    selectedRole.value = '';
    errorMessage.value = '';
    showRegister.value = false;
};
</script>

<template>
    <div class="login-container">
        <!-- 主選單 -->
        <div v-if="!showLogin && !showRegister && !showForgotPassword" class="card main-menu">
            <h2>歡迎使用健康檢查數據平台</h2>
            <div class="button-group">
                <button class="login-btn" @click="showLogin = true">登入</button>
                <button class="register-btn" @click="showRegister = true">註冊</button>
            </div>
        </div>

        <!-- 登入表單區域 -->
        <div v-if="showLogin" class="card login-form">
            <h2>登入</h2>
            <div class="form-group">
                <label>身分證號/員工編號</label>
                <input v-model="username" type="text" placeholder="輸入身分證號或員工編號" />
            </div>
            <div class="form-group">
                <label>密碼</label>
                <input v-model="password" type="password" placeholder="輸入密碼" />
            </div>
            <div class="form-group">
                <label>選擇角色</label>
                <select v-model="selectedRole" class="role-select">
                    <option value="">請選擇角色</option>
                    <option value="hospital">健檢中心</option>
                    <option value="patient">用戶</option>
                    <option value="other">其他</option>
                </select>
            </div>
            <div class="button-group">
                <button class="login-btn" @click="login">登入</button>
                <button class="back-btn" @click="goBack">返回</button>
            </div>
            <div class="extra-links">
                <button class="link-btn" @click="showForgotPassword = true; showLogin = false">忘記密碼？</button>
            </div>
            <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
        </div>

        <!-- 忘記密碼表單 -->
        <div v-if="showForgotPassword" class="card login-form">
            <h2>忘記密碼</h2>
            <div class="form-group">
                <label>身分證號/員工編號</label>
                <input v-model="username" type="text" placeholder="輸入身分證號或員工編號" />
            </div>
            <div class="form-group">
                <label>選擇角色</label>
                <select v-model="selectedRole" class="role-select">
                    <option value="">請選擇角色</option>
                    <option value="hospital">健檢中心</option>
                    <option value="patient">用戶</option>
                    <option value="other">其他</option>
                </select>
            </div>
            <div class="button-group">
                <button class="login-btn" @click="forgotPassword">重設密碼</button>
                <button class="back-btn" @click="goBack">返回</button>
            </div>
            <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
        </div>

        <!-- 註冊表單 -->
        <div v-if="showRegister" class="card login-form">
            <h2>註冊新帳號</h2>
            <div class="form-group">
                <label>身分證號/員工編號</label>
                <input v-model="username" type="text" placeholder="輸入身分證號或員工編號" />
            </div>
            <div class="form-group">
                <label>密碼</label>
                <input v-model="password" type="password" placeholder="輸入密碼" />
            </div>
            <div class="form-group">
                <label>選擇角色</label>
                <select v-model="selectedRole" class="role-select">
                    <option value="">請選擇角色</option>
                    <option value="hospital">健檢中心</option>
                    <option value="patient">用戶</option>
                    <option value="other">其他</option>
                </select>
            </div>
            <div class="button-group">
                <button class="login-btn" @click="register">註冊</button>
                <button class="back-btn" @click="goBack">返回</button>
            </div>
            <p v-if="errorMessage" class="error">{{ errorMessage }}</p>
        </div>
    </div>
</template>

<style scoped>
.login-container { 
    display: flex; 
    justify-content: center; 
    align-items: center; 
    padding: 40px 20px; 
}

.card { 
    background: var(--white); 
    padding: 40px; 
    border-radius: var(--border-radius); 
    box-shadow: var(--shadow); 
    width: 100%; 
    max-width: 500px; 
    text-align: center; 
    transition: transform 0.3s ease; 
}

.card:hover { 
    transform: translateY(-5px); 
}

.main-menu {
    text-align: center;
}

.main-menu h2 {
    margin-bottom: 30px;
    font-size: 24px;
}

.form-group { 
    margin-bottom: 20px; 
}

.form-group label { 
    display: block; 
    margin-bottom: 5px; 
    font-weight: bold; 
}

.role-select {
    width: 100%;
    padding: 10px;
    border: 1px solid #d1d5db;
    border-radius: var(--border-radius);
    font-size: 16px;
    background-color: white;
}

.button-group { 
    display: flex; 
    justify-content: center; 
    gap: 15px; 
    margin-top: 20px; 
}

.login-btn, .register-btn { 
    background: var(--primary-color); 
    color: var(--white); 
    padding: 12px 24px; 
    border: none; 
    border-radius: var(--border-radius); 
    font-size: 16px;
    min-width: 120px;
}

.register-btn {
    background: var(--secondary-color);
}

.back-btn { 
    background: #6c757d; 
    color: var(--white); 
    padding: 12px 24px; 
    border: none; 
    border-radius: var(--border-radius); 
    font-size: 16px;
    min-width: 120px;
}

.extra-links { 
    margin-top: 15px; 
}

.link-btn { 
    background: none; 
    border: none; 
    color: var(--primary-color); 
    text-decoration: underline; 
    cursor: pointer; 
}

.error { 
    color: #dc3545; 
    margin-top: 15px; 
    font-size: 14px; 
}
</style>