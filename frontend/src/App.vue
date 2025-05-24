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
  background: #f9f7f4 !important;
  color: #222 !important;
  box-shadow: none;
  border-bottom: 1px solid #222;
}
.brand-title {
  font-size: 1.8rem;
  font-weight: 900;
  color: #080808;
  letter-spacing: 1.5px;
  line-height: 1.2;
  text-shadow: none;
  font-family: 'Montserrat', 'Noto Sans TC', 'Arial', sans-serif;
}
.brand-slogan {
  color: #080808;
  font-size: 1.1rem;
  font-weight: 500;
  letter-spacing: 1.2px;
  line-height: 1.2;
}
.footer-brand {
  background: linear-gradient(90deg, #f9f7f4);
  color: #fff;
  font-weight: 600;
  font-size: 1.1rem;
  padding: 0;
}
.footer-title {
  color: #080808;
  font-weight: bold;
  letter-spacing: 1.2px;
  font-size: 1.1rem;
}

.v-icon {
  color: #1565c0 !important;
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
  background: #f9f7f4;
}
.main-content {
  background: #f9f7f4;
}
.nav-link {
  font-size: 1.1rem;
  color: #6b7280 !important;
  font-weight: 700;
  letter-spacing: 1.5px;
  text-transform: uppercase;
  margin-right: 24px;
  font-family: 'Montserrat', 'Noto Sans TC', 'Arial', sans-serif;
}
.nav-link:last-child {
  margin-right: 0;
}
.v-btn.nav-link {
  min-width: unset;
  padding: 0 8px;
}
.v-img {
  filter: grayscale(100%) contrast(1.1);
}
</style>