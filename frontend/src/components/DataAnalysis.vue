<template>
    <div class="analysis-container">
      <h3>您的健康數據分析</h3>
      <p class="description">這裡將顯示您的健康檢查數據（待實現）</p>
      <div class="placeholder-data">
        <p>示例數據：心率 - 72 bpm</p>
        <p>示例數據：血壓 - 120/80 mmHg</p>
      </div>
    </div>
  </template>
  
  <script setup>
  </script>
  
  <style scoped>
  .analysis-container {
    text-align: center;
  }
  
  h3 {
    font-size: 24px;
    margin-bottom: 10px;
  }
  
  .description {
    color: var(--muted-color);
    margin-bottom: 20px;
  }
  
  .placeholder-data {
    background: #f8f9fa;
    padding: 20px;
    border-radius: var(--border-radius);
    border: 1px solid #e0e0e0;
  }
  
  .placeholder-data p {
    margin: 10px 0;
    color: var(--muted-color);
  }
  </style>