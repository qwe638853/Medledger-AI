<script setup>
import { onMounted } from 'vue';
import axios from 'axios';
import Home from './views/Home.vue';
import { useAuth } from './composables/useAuth';
import { useUser } from './composables/useUser';
import { useNavigation } from './composables/useNavigation';

// 設置 axios 基礎 URL（目前為空）
axios.defaults.baseURL = '';

// 使用 composables
const { 
    userRole, 
    isLoggedIn, 
    currentUser, 
    token, 
    showLoginForm, 
    initAuth, 
    login, 
    logout 
} = useAuth();

const { forgotPassword, register, fetchData } = useUser(token, currentUser);
const { showFooter, menuItems, goToHome } = useNavigation();

// 在組件掛載時初始化認證狀態
onMounted(() => {
    initAuth();
});

// 處理登入表單顯示切換
const toggleLoginForm = () => {
    showLoginForm.value = !showLoginForm.value;
};

// 處理頁腳顯示切換
const toggleFooter = () => {
    showFooter.value = !showFooter.value;
};
</script>

<template>
    <Home
        :user-role="userRole"
        :is-logged-in="isLoggedIn"
        :current-user="currentUser"
        :show-login-form="showLoginForm"
        :show-footer="showFooter"
        :menu-items="menuItems"
        @login="login"
        @logout="logout"
        @forgot-password="forgotPassword"
        @register="register"
        @go-home="() => goToHome(showLoginForm)"
        @toggle-login-form="toggleLoginForm"
        @toggle-footer="toggleFooter"
    />
</template>

<style>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.v-application {
  font-family: 'Noto Sans TC', sans-serif !important;
}

.v-footer {
  transition: all 0.3s ease;
}
</style>