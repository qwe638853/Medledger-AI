<script setup>
import { ref } from 'vue';
import LoginForm from './components/LoginForm.vue';
import UserDashboard from './UserDashboard.vue';
import HospitalDashboard from './HospitalDashboard.vue';

const userRole = ref(null); // 'hospital' or 'patient'
const isLoggedIn = ref(false);
const currentUser = ref('');

const login = (data) => {
  userRole.value = data.role;
  currentUser.value = data.username;
  isLoggedIn.value = true;
};

const logout = () => {
  userRole.value = null;
  currentUser.value = '';
  isLoggedIn.value = false;
};
</script>

<template>
  <div class="app-container">
    <!-- 封面圖片區塊 -->
    <div class="cover-container">
      <div class="cover-overlay">
        <h1>健康檢查數據平台</h1>
        <p>安全管理您的健康數據</p>
      </div>
    </div>

    <!-- 主要內容 -->
    <div class="content-container">
      <LoginForm v-if="!isLoggedIn" @login="login" />
      <UserDashboard 
        v-if="isLoggedIn && userRole === 'patient'" 
        :username="currentUser" 
        @logout="logout"
      />
      <HospitalDashboard 
        v-if="isLoggedIn && userRole === 'hospital'" 
        :username="currentUser" 
        @logout="logout"
      />
    </div>
  </div>
</template>

<style scoped>
.app-container {
  min-height: 100vh;
}

.cover-container {
  position: relative;
  background: url('https://i.pinimg.com/736x/25/3c/11/253c118ba1720173da5cb8e010eed930.jpg') 
    no-repeat center center;
  background-size: cover;
  height: 450px; /* 適合桌面端的高度 */
  width: 100%;
}

.cover-overlay {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(149, 169, 190, 0.6); /* 使用藍色半透明遮罩 */
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  color: var(--white);
  text-align: center;
  padding: 30px;
}

.cover-overlay h1 {
  font-size: 48px;
  font-weight: 700;
  margin-bottom: 15px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.3);
}

.cover-overlay p {
  font-size: 22px;
  font-weight: 300;
  text-shadow: 1px 1px 2px rgba(0, 0, 0, 0.3);
}

.content-container {
  max-width: 1200px; /* 適合桌面端的寬度 */
  margin: 0 auto;
  padding: 40px 20px;
}
</style>