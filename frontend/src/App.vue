<template>
  <v-app>
    <v-navigation-drawer app v-model="drawer">
      <v-list>
        <v-list-item
          v-for="item in menuItems"
          :key="item.title"
          :to="item.path"
          :title="item.title"
        />
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app color="primary" dark>
      <v-app-bar-nav-icon @click="drawer = !drawer" />
      <v-toolbar-title>健康檢查數據平台</v-toolbar-title>
      <v-spacer />
      <v-btn v-if="isLoggedIn" text @click="logout">登出</v-btn>
      <v-progress-linear v-if="loading" indeterminate color="white" />
    </v-app-bar>

    <v-main>
      <router-view
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
        @go-home="goToHome(showLoginForm)"
        @toggle-login-form="showLoginForm = !showLoginForm"
        @toggle-footer="showFooter = !showFooter"
        @show-snackbar="showSnackbar"
      />
    </v-main>

    <v-footer app v-if="showFooter" color="primary" dark>
      <v-row justify="center">
        <v-col class="text-center" cols="12">
          © {{ new Date().getFullYear() }} 健康檢查數據平台
        </v-col>
      </v-row>
    </v-footer>

    <v-snackbar v-model="snackbar" :color="snackbarColor" timeout="3000">
      {{ snackbarMessage }}
      <template v-slot:actions>
        <v-btn color="white" text @click="snackbar = false">關閉</v-btn>
      </template>
    </v-snackbar>
  </v-app>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import { useAuth } from './composables/useAuth';
import { useNavigation } from './composables/useNavigation';
import { useUser } from './composables/useUser';

const { userRole, isLoggedIn, currentUser, showLoginForm, loading, initAuth, login, logout, register } = useAuth();
const { showFooter, menuItems, goToHome } = useNavigation();
const { forgotPassword } = useUser();

const drawer = ref(false);
const snackbar = ref(false);
const snackbarMessage = ref('');
const snackbarColor = ref('error');

const showSnackbar = (message, color = 'error') => {
  snackbarMessage.value = message;
  snackbarColor.value = color;
  snackbar.value = true;
};

const handleSnackbarEvent = (event) => {
  const { message, color } = event.detail;
  showSnackbar(message, color);
};

onMounted(() => {
  initAuth();
  window.addEventListener('show-snackbar', handleSnackbarEvent);
});

onUnmounted(() => {
  window.removeEventListener('show-snackbar', handleSnackbarEvent);
});
</script>