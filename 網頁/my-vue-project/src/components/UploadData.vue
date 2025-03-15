<template>
    <div class="upload-container">
      <h3>上傳健康檢查數據</h3>
      <p class="description">請選擇文件（支援 PDF、JPG、PNG 格式）</p>
      <input type="file" @change="handleFileUpload" accept=".pdf,.jpg,.png" />
      <button class="upload-btn">上傳</button>
    </div>
  </template>
  
  <script setup>
  const handleFileUpload = (event) => {
    const file = event.target.files[0];
    if (file) {
      console.log('上傳的文件:', file.name);
      // 這裡可以添加後端上傳邏輯
    }
  };
  </script>
  
  <style scoped>
  .upload-container {
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
  
  input[type="file"] {
    display: block;
    margin: 0 auto 20px;
    padding: 10px;
    border: 1px solid #d1d5db;
    border-radius: var(--border-radius);
    font-size: 16px;
  }
  
  .upload-btn {
    background: var(--primary-color);
    color: var(--white);
    padding: 10px 20px;
    border: none;
    border-radius: var(--border-radius);
    font-size: 16px;
  }
  </style>