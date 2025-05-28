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
      <div class="px-4 w-100">
        <div class="d-flex align-center">
          <router-link to="/" class="d-flex align-center text-decoration-none">
            <v-img
              :src="logo"
              alt="智療鏈 Logo"
              width="40"
              height="40"
              class="logo-image"
            />
            <span class="brand-title">Medledger AI</span>
          </router-link>
          
          <nav class="d-flex align-center ml-4">
            <v-btn text class="nav-link mx-2" @click="scrollToSection('about-section')">關於我們</v-btn>
            <v-btn text class="nav-link mx-2" @click="scrollToSection('contact-section')">聯絡資訊</v-btn>
          </nav>

          <v-spacer />
        </div>
      </div>
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>

    <v-footer app class="footer-brand" height="56">
      <v-container class="d-flex align-center justify-center fill-height">
        <span class="footer-title">© {{ new Date().getFullYear() }} Medledger | 智慧守護，鏈接健康未來</span>
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
import logo from '@/assets/logo.png';

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

function scrollToSection(sectionId) {
  const el = document.getElementById(sectionId);
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' });
  }
}
</script>

<style scoped>
body, #app, .v-application, .v-app {
  background: #f9f7f4 !important;
  font-size: 1.18rem !important;
  line-height: 1.7 !important;
}

.brand-app-bar {
  background: #F9F7F4 !important;
  border-bottom: 1px solid #e5e7eb;
  height: 80px !important;
}

.brand-app-bar :deep(.v-toolbar__content) {
  padding: 0;
}

.brand-app-bar :deep(.v-container) {
  padding-left: 16px !important;
  padding-right: 16px !important;
}

.logo-image {
  border-radius: 8px;
  transition: transform 0.2s ease;
  width: 60px !important;
  height: 60px !important;
  filter: none !important;
}

.brand-title {
  font-family: 'Inter', sans-serif;
  font-size: 2.2rem;
  font-weight: 900;
  color: #111827;
  letter-spacing: -0.5px;
}

.nav-link {
  color: #111827 !important;
  font-weight: 900;
  font-size: 2.2rem;
  text-transform: none;
  letter-spacing: -0.5px;
}

.nav-link:hover {
  opacity: 0.8;
}

.footer-brand {
  background: #f3f2ef !important;
  border-top: 1px solid #e5e7eb;
}

.footer-title {
  color: #6b7280;
  font-size: 1.1rem !important;
  font-weight: 400;
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

.v-btn.nav-link {
  min-width: unset;
  padding: 0 8px;
}

.v-img {
  filter: grayscale(100%) contrast(1.1);
}

h1, .text-h1 { font-size: 3.2rem !important; }
h2, .text-h2 { font-size: 2.5rem !important; }
h3, .text-h3 { font-size: 2rem !important; }
h4, .text-h4 { font-size: 1.7rem !important; }
h5, .text-h5 { font-size: 1.4rem !important; }
h6, .text-h6 { font-size: 1.2rem !important; }

.v-btn, .v-btn__content, .v-btn .v-icon {
  font-size: 1.18rem !important;
}

.v-list-item, .v-list-item__title {
  font-size: 1.15rem !important;
}

input, textarea, .v-input__slot, .v-label {
  font-size: 1.15rem !important;
}
</style>