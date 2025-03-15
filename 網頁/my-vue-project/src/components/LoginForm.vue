<script setup>
import { ref } from 'vue';

const username = ref('');
const password = ref('');
const errorMessage = ref('');
const roleSelected = ref(false);
const selectedRole = ref('');

const emit = defineEmits(['login']);

const selectRole = (role) => {
  selectedRole.value = role;
  roleSelected.value = true;
};

const goBack = () => {
  roleSelected.value = false;
  selectedRole.value = '';
  username.value = '';
  password.value = '';
  errorMessage.value = '';
};

const login = () => {
  if (username.value === 'admin' && password.value === '123456') {
    emit('login', { username: username.value, role: selectedRole.value });
    roleSelected.value = false;
    selectedRole.value = '';
    username.value = '';
    password.value = '';
    errorMessage.value = '';
  } else {
    errorMessage.value = '帳號或密碼錯誤，請再試一次';
  }
};
</script>

<template>
  <div class="login-container">
    <!-- 角色選擇區域 -->
    <div v-if="!roleSelected" class="card role-selection">
      <h2>請選擇您的角色</h2>
      <div class="role-buttons">
        <button class="role-btn hospital-btn" @click="selectRole('hospital')">
          <span>🏥</span> 醫院員工
        </button>
        <button class="role-btn patient-btn" @click="selectRole('patient')">
          <span>👤</span> 使用者
        </button>
      </div>
    </div>

    <!-- 登入表單區域 -->
    <div v-if="roleSelected" class="card login-form">
      <h2>{{ selectedRole === 'hospital' ? '🏥 醫院員工登入' : '👤 使用者登入' }}</h2>
      <div class="form-group">
        <input v-model="username" type="text" placeholder="帳號" />
      </div>
      <div class="form-group">
        <input v-model="password" type="password" placeholder="密碼" />
      </div>
      <div class="button-group">
        <button class="login-btn" @click="login">登入</button>
        <button class="back-btn" @click="goBack">返回選擇角色</button>
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
  max-width: 500px; /* 適合桌面端的寬度 */
  text-align: center;
  transition: transform 0.3s ease;
}

.card:hover {
  transform: translateY(-5px);
}

h2 {
  font-size: 28px;
  margin-bottom: 20px;
}

.role-buttons {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-top: 20px;
}

.role-btn {
  padding: 12px 24px;
  font-size: 16px;
  border: none;
  border-radius: var(--border-radius);
  color: var(--white);
  display: flex;
  align-items: center;
  gap: 8px;
}

.hospital-btn {
  background: var(--primary-color);
}

.patient-btn {
  background: var(--secondary-color);
}

.form-group {
  margin-bottom: 20px;
}

.button-group {
  display: flex;
  justify-content: center;
  gap: 15px;
  margin-top: 20px;
}

.login-btn {
  background: var(--primary-color);
  color: var(--white);
  padding: 12px 24px;
  border: none;
  border-radius: var(--border-radius);
  font-size: 16px;
}

.back-btn {
  background: #6c757d;
  color: var(--white);
  padding: 12px 24px;
  border: none;
  border-radius: var(--border-radius);
  font-size: 16px;
}

.error {
  color: #dc3545;
  margin-top: 15px;
  font-size: 14px;
}
</style>