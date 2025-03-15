<template>
    <div class="dashboard-container">
      <header class="dashboard-header">
        <h2>👤 使用者儀表板</h2>
        <p>歡迎，{{ username }}！這裡顯示您的健康檢查數據與分析</p>
        <button class="logout-btn" @click="$emit('logout')">登出</button>
      </header>
      <section class="data-section card">
        <DataAnalysis />
      </section>
    </div>
  </template>
  
  <script setup>
  import DataAnalysis from './components/DataAnalysis.vue';
  
  defineProps({
    username: {
      type: String,
      default: '使用者'
    }
  });
  
  defineEmits(['logout']);
  </script>
  
  <style scoped>
  .dashboard-container {
    max-width: 1200px; /* 適合桌面端 */
    margin: 40px auto;
    padding: 20px;
  }
  
  .dashboard-header {
    text-align: center;
    margin-bottom: 30px;
  }
  
  .dashboard-header h2 {
    font-size: 32px;
    margin-bottom: 10px;
  }
  
  .dashboard-header p {
    color: var(--muted-color);
    font-size: 18px;
  }
  
  .logout-btn {
    margin-top: 15px;
    padding: 10px 20px;
    background-color: #dc3545;
    color: var(--white);
    border: none;
    border-radius: var(--border-radius);
    font-size: 16px;
  }
  
  .card {
    background: var(--white);
    padding: 30px;
    border-radius: var(--border-radius);
    box-shadow: var(--shadow);
    transition: transform 0.3s ease;
  }
  
  .card:hover {
    transform: translateY(-5px);
  }
  
  .data-section {
    text-align: center;
  }
  </style>