<script setup>
// 引入 Vue 的響應式 API
import { ref } from 'vue';

// 創建響應式變量用於存儲病患身分證號
const id_number = ref('');
// 創建響應式變量用於存儲上傳的檔案
const file = ref(null);
// 創建響應式變量用於存儲訊息
const message = ref('');
// 定義組件的事件
const emit = defineEmits(['logout']);

// 處理表單提交的函數
const handleSubmit = async (e) => {
    e.preventDefault();

    // 驗證身分證號是否已輸入
    if (!id_number.value) {
        message.value = '請輸入病患身分證號';
        message.classList.add('error');
        return;
    }

    // 驗證是否已選擇檔案
    if (!file.value) {
        message.value = '請選擇一個檔案';
        message.classList.add('error');
        return;
    }

    // 從 localStorage 獲取 token
    const token = localStorage.getItem('token');

    // 驗證是否已登入
    if (!token) {
        message.value = '請先登入';
        message.classList.add('error');
        return;
    }

    // 創建 FormData 對象用於上傳檔案
    const formData = new FormData();
    formData.append('file', file.value);

    try {
        // 發送上傳請求
        const response = await fetch(`https://7aa9-140-124-249-9.ngrok-free.app/default/health-check/upload/${id_number.value}`, {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`,
                'Accept': 'application/json'
            },
            body: formData
        });

        // 解析響應數據
        const data = await response.json();

        // 根據響應狀態顯示不同的訊息
        if (response.ok) {
            message.value = '上傳成功！' + (data.message || '');
            message.classList.add('success');
        } else {
            message.value = data.detail || '上傳失敗';
            message.classList.add('error');
        }
    } catch (error) {
        // 錯誤處理
        message.value = '發生錯誤：' + error.message;
        message.classList.add('error');
    }

    // 重置表單
    id_number.value = '';
    file.value = null;
};
</script>

<template>
    <!-- 儀表板容器 -->
    <div class="dashboard-container">
        <!-- 標題 -->
        <h2>🏥 健康中心儀表板</h2>
        <!-- 歡迎訊息 -->
        <p>歡迎，{{ username }}</p>
        <!-- 上傳表單 -->
        <form id="upload-form" @submit="handleSubmit">
            <!-- 身分證號輸入欄位 -->
            <div class="form-group">
                <label>病患身分證號</label>
                <input v-model="id_number" type="text" placeholder="輸入病患身分證號" />
            </div>
            <!-- 檔案上傳欄位 -->
            <div class="form-group">
                <label>上傳健康檢查檔案</label>
                <input type="file" @change="e => file = e.target.files[0]" />
            </div>
            <!-- 提交按鈕 -->
            <button type="submit">上傳</button>
        </form>
        <!-- 訊息顯示區域 -->
        <p v-if="message" :class="message.includes('成功') ? 'success' : 'error'">
            {{ message }}
        </p>
        <!-- 登出按鈕 -->
        <button @click="$emit('logout')">登出</button>
    </div>
</template>

<style scoped>
/* 儀表板容器樣式 */
.dashboard-container { padding: 20px; }
/* 標題樣式 */
h2 { font-size: 28px; margin-bottom: 20px; }
/* 表單組樣式 */
.form-group { margin-bottom: 20px; }
/* 表單標籤樣式 */
.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
/* 成功訊息樣式 */
.success { color: green; }
/* 錯誤訊息樣式 */
.error { color: red; }
</style>