<script setup>
import { defineProps } from 'vue';

const props = defineProps({
  data: { type: Array, default: () => [] }
});
</script>

<template>
  <div class="analysis-container">
    <h3>您的健康數據分析</h3>
    <p class="description">這裡顯示您的健康檢查數據</p>
    <div v-if="data.length" class="placeholder-data">
      <p v-for="item in data" :key="item.id">
        {{ item.patientId }} - {{ item.data }} ({{ item.timestamp }})
      </p>
    </div>
    <p v-else>暫無數據</p>
  </div>
</template>

<style scoped>
.analysis-container { text-align: center; }
h3 { font-size: 24px; margin-bottom: 10px; }
.description { color: var(--muted-color); margin-bottom: 20px; }
.placeholder-data { background: #f8f9fa; padding: 20px; border-radius: var(--border-radius); border: 1px solid #e0e0e0; }
.placeholder-data p { margin: 10px 0; color: var(--muted-color); }
</style>