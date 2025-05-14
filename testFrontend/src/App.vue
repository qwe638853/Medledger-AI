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

    <v-app-bar app color="primary" dark>
      <v-toolbar-title>健康檢查數據平台</v-toolbar-title>
      <v-spacer />
      <v-btn v-if="authStore.isLoggedIn" text @click="authStore.logout">登出</v-btn>
      <v-progress-linear v-if="authStore.loading || userStore.loading" indeterminate color="white" />
    </v-app-bar>

    <v-main>
      <router-view />
    </v-main>

    <v-footer app color="primary" dark>
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
import { useAuthStore, useUserStore, useNavigationStore } from './stores';

const authStore = useAuthStore();
const userStore = useUserStore();
const navigationStore = useNavigationStore();

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
  window.addEventListener('show-snackbar', handleSnackbarEvent);
});

onUnmounted(() => {
  window.removeEventListener('show-snackbar', handleSnackbarEvent);
});
</script>