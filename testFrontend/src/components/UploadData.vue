<script setup>
import { defineEmits } from 'vue';

const emit = defineEmits(['upload']);

// 處理上傳
const handleUpload = () => {
  const patientId = prompt('請輸入病患 ID:');
  const data = prompt('請輸入健康數據 (例如: blood_pressure: 120/80):');
  const timestamp = new Date().toISOString();

  if (patientId && data) {
    emit('upload', { patientId, data, timestamp });
  } else {
    alert('請輸入所有必要資訊');
  }
};
</script>

<template>
  <div class="upload-container">
    <h3>上傳健康檢查數據</h3>
    <p class="description">請手動輸入數據</p>
    <button class="upload-btn" @click="handleUpload">上傳數據</button>
  </div>
</template>

<style scoped>
.upload-container { text-align: center; }
h3 { font-size: 24px; margin-bottom: 10px; }
.description { color: var(--muted-color); margin-bottom: 20px; }
.upload-btn { background: var(--primary-color); color: var(--white); padding: 10px 20px; border: none; border-radius: var(--border-radius); font-size: 16px; }
</style>