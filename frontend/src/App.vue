<template>
  <v-app>
    <v-navigation-drawer app v-model="navigationStore.drawer">
      <v-list>
        <v-list-item
          v-for="item in navigationStore.navItems"
          :key="item.title"
          :to="item.path"
          :title="item.title"
          @click="item.action && item.action()"
        />
      </v-list>
    </v-navigation-drawer>

    <v-app-bar app class="brand-app-bar" flat>
      <v-container class="d-flex align-center pa-0">
        <div>
          <div class="brand-title">智療鏈</div>
          <div class="brand-slogan">智慧守護，鏈接健康未來</div>
        </div>
        <v-spacer />
        <v-btn v-if="!authStore.isLoggedIn" color="primary" variant="outlined" class="mx-1" :to="{ path: '/register' }">註冊</v-btn>
        <v-btn v-if="!authStore.isLoggedIn" color="primary" variant="flat" class="mx-1" :to="{ path: '/login' }">登入</v-btn>
        <v-btn v-if="authStore.isLoggedIn" color="primary" variant="text" class="mx-1" @click="authStore.logout">登出</v-btn>
      </v-container>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>

    <v-footer app class="footer-brand" height="56">
      <v-container class="d-flex align-center justify-center fill-height">
        <span class="footer-title">© {{ new Date().getFullYear() }} 智療鏈 | 智慧守護，鏈接健康未來</span>
      </v-container>
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
import { useAuthStore, useUserStore, useNavigationStore } from './stores';

const authStore = useAuthStore();
const userStore = useUserStore();
const navigationStore = useNavigationStore();

const snackbar = ref(false);
const snackbarMessage = ref('');
const snackbarColor = ref('info');

const showSnackbar = (message, color = 'info') => {
  snackbarMessage.value = message;
  snackbarColor.value = color;
  snackbar.value = true;
};

const handleSnackbarEvent = (event) => {
  const { message, color } = event.detail;
  showSnackbar(message, color);
};

onMounted(() => {
  window.addEventListener('show-snackbar', handleSnackbarEvent);
});

onUnmounted(() => {
  window.removeEventListener('show-snackbar', handleSnackbarEvent);
});
</script>

<style scoped>
.brand-app-bar {
  background: linear-gradient(90deg, #1565c0 60%, #42a5f5 100%) !important;
  color: #fff !important;
  box-shadow: 0 2px 8px 0 rgba(21,101,192,0.08);
}
.brand-title {
  font-size: 1.8rem;
  font-weight: 900;
  color: #fff;
  letter-spacing: 1.5px;
  line-height: 1.2;
  text-shadow: 0 1px 2px rgba(0,0,0,0.1);
}
.brand-slogan {
  color: #e3f2fd;
  font-size: 1.1rem;
  font-weight: 500;
  letter-spacing: 1.2px;
  line-height: 1.2;
}
.footer-brand {
  background: linear-gradient(90deg, #1565c0 60%, #42a5f5 100%);
  color: #fff;
  font-weight: 600;
  font-size: 1.1rem;
  padding: 0;
}
.footer-title {
  font-weight: bold;
  letter-spacing: 1.2px;
  font-size: 1.1rem;
}
.v-btn {
  background: #1565c0 !important;
  color: #fff !important;
  border: none;
}
.v-btn[variant="outlined"], .v-btn[outlined] {
  background: #fff !important;
  color: #1565c0 !important;
  border: 2px solid #1565c0 !important;
}
.v-icon {
  color: #1565c0 !important;
}
.sidebar {
  background: rgba(255, 255, 255, 0.9) !important;
  backdrop-filter: blur(10px);
  border-right: 1px solid rgba(255, 255, 255, 0.3);
}
.sidebar-list {
  padding: 8px;
}
.brand-item {
  margin-bottom: 8px;
}
.nav-item {
  margin: 4px 0;
  transition: all 0.3s ease;
  color: #1565c0 !important;
}
.nav-item:hover {
  background: rgba(33, 150, 243, 0.1);
}
.main-content {
  background: #f8fafc;
}
</style>