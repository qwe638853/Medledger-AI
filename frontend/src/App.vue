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
        <v-img
          src="@/assets/logo.png"
          alt="智療鏈 Logo"
          class="mr-3"
          max-width="44"
          max-height="44"
          style="border-radius: 12px; background: #e3f2fd;"
        />
        <span class="brand-title mr-8">智療鏈</span>
        <v-btn text class="nav-link mr-2" to="#about">關於我們</v-btn>
        <v-btn text class="nav-link" to="#contact">聯絡資訊</v-btn>
        <v-spacer />
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
body, #app, .v-application, .v-app {
  background: #f9f7f4 !important;

}

.brand-app-bar {
  background: #F9F7F4 !important;
  color: #222 !important;
  box-shadow: none;
  border-bottom: 1px solid #e5e7eb;
  height: 80px !important;
}

>>>>>>> 56a3e8d34a9523f27a697194858ce2fe641a2be5
.brand-title {
  font-size: 1.8rem;
  font-weight: 900;
  color: #111827;
  letter-spacing: -0.5px;
  line-height: 1.2;
  font-family: 'Inter', 'Noto Sans TC', sans-serif;
}

.brand-slogan {
  color: #080808;
  font-size: 1.1rem;
  font-weight: 500;
  letter-spacing: 1.2px;
  line-height: 1.2;
}

.footer-brand {
  background: #f3f2ef !important;
  border-top: 1px solid #e5e7eb;
  color: #6b7280;
  font-weight: 400;
  font-size: 0.875rem;
  padding: 0;
}

.footer-title {

  color: #6b7280;
  font-weight: 400;
  letter-spacing: 0;
  font-size: 0.875rem;

}

.v-icon {
  color: #111827 !important;


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
  font-size: 1rem;
  color: #6b7280 !important;
  font-weight: 500;
  letter-spacing: 0;
  text-transform: none;
  margin-right: 32px;
  font-family: 'Inter', 'Noto Sans TC', sans-serif;
  transition: all 0.2s ease;
  position: relative;
}

.nav-link::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 0;
  height: 2px;
  background: #111827;
  transition: width 0.2s ease;
}

.nav-link:hover::after {
  width: 100%;
}

.v-btn.nav-link {
  min-width: unset;
  padding: 0 8px;
}

.v-img {
  filter: grayscale(100%) contrast(1.1);
}
</style>