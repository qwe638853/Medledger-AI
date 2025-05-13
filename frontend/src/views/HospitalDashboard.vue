<script setup>
import { ref, onMounted } from 'vue';
import { useAuth } from '../composables/useAuth';
import axios from 'axios';

const { currentUser, token, logout } = useAuth();
const uploadedData = ref([]);
const loading = ref(false);
const file = ref(null);
const uploadMessage = ref('');

onMounted(async () => {
  await fetchUploadedData();
});

const fetchUploadedData = async () => {
  loading.value = true;
  try {
    const response = await axios.get(
      `https://7aa9-140-124-249-9.ngrok-free.app/default/health-check/uploaded/${currentUser.value}`,
      {
        headers: {
          Authorization: `Bearer ${token.value}`,
          Accept: 'application/json'
        },
        timeout: 10000
      }
    );
    uploadedData.value = response.data || [];
  } catch (error) {
    document.dispatchEvent(new CustomEvent('show-snackbar', {
      detail: { message: `獲取已上傳數據失敗：${error.message}`, color: 'error' }
    }));
    uploadedData.value = [];
  } finally {
    loading.value = false;
  }
};

const handleFileUpload = async () => {
  if (!file.value) {
    document.dispatchEvent(new CustomEvent('show-snackbar', {
      detail: { message: '請選擇檔案', color: 'error' }
    }));
    return;
  }

  const formData = new FormData();
  formData.append('file', file.value);
  formData.append('id_number', currentUser.value);

  loading.value = true;
  try {
    const response = await axios.post(
      'https://7aa9-140-124-249-9.ngrok-free.app/default/health-check/upload',
      formData,
      {
        headers: {
          Authorization: `Bearer ${token.value}`,
          'Content-Type': 'multipart/form-data'
        },
        timeout: 10000
      }
    );
    uploadedData.value.push(response.data);
    uploadMessage.value = '上傳成功！';
    document.dispatchEvent(new CustomEvent('show-snackbar', {
      detail: { message: '上傳成功！', color: 'success' }
    }));
  } catch (error) {
    document.dispatchEvent(new CustomEvent('show-snackbar', {
      detail: { message: `上傳失敗：${error.message}`, color: 'error' }
    }));
    uploadMessage.value = '';
  } finally {
    loading.value = false;
  }
};

const handleLogout = () => {
  logout();
};
</script>

<template>
  <div class="dashboard-container">
    <h2>🏥 健檢中心儀表板</h2>
    <p>歡迎，{{ currentUser }}</p>
    <div class="upload-section">
      <h3>上傳健康檢查數據</h3>
      <input type="file" @change="e => file = e.target.files[0]" />
      <v-btn color="primary" @click="handleFileUpload" :loading="loading">上傳</v-btn>
      <p v-if="uploadMessage">{{ uploadMessage }}</p>
    </div>
    <div v-if="loading" class="loading">加載中...</div>
    <div v-else-if="uploadedData.length">
      <h3>已上傳數據</h3>
      <ul>
        <li v-for="item in uploadedData" :key="item.id">
          用戶：{{ item.user_id }} - 數據：{{ item.content || item }}
        </li>
      </ul>
    </div>
    <p v-else>暫無已上傳資料</p>
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
.upload-section { margin-bottom: 20px; }
button { padding: 10px 20px; background-color: #1976d2; color: white; border: none; border-radius: 4px; cursor: pointer; }
button:hover { background-color: #1565c0; }
</style>