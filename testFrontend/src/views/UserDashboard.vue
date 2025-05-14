<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores';
import { healthCheckService, notifyError } from '../services';

const authStore = useAuthStore();
const currentUser = computed(() => authStore.currentUser);
const healthData = ref([]);
const loading = ref(false);

onMounted(async () => {
  loading.value = true;
  try {
    const response = await healthCheckService.fetchUserHealthData(currentUser.value);
    healthData.value = response || [];
  } catch (error) {
    notifyError(`獲取健康檢查數據失敗：${error.message}`);
    healthData.value = [];
  } finally {
    loading.value = false;
  }
});

const handleLogout = () => {
  authStore.logout();
};
</script>

<template>
  <div class="dashboard-container">
    <h2>👤 使用者儀表板</h2>
    <p>歡迎，{{ currentUser }}</p>
    <div v-if="loading" class="loading">加載中...</div>
    <div v-else-if="healthData.length">
      <h3>您的健康檢查數據</h3>
      <ul>
        <li v-for="item in healthData" :key="item.id">
          {{ item.content || item }}
        </li>
      </ul>
    </div>
    <p v-else>暫無資料</p>
    <button @click="handleLogout">登出</button>
  </div>
</template>

<style scoped>
.dashboard-container { padding: 20px; }
h2 { font-size: 28px; margin-bottom: 20px; }
h3 { font-size: 20px; margin-top: 20px; }
ul { list-style-type: none; padding: 0; }
li { padding: 10px 0; border-bottom: 1px solid #ddd; }
.loading { color: #1976d2; font-style: italic; }
button { padding: 10px 20px; background-color: #1976d2; color: white; border: none; border-radius: 4px; cursor: pointer; }
button:hover { background-color: #1565c0; }
</style>