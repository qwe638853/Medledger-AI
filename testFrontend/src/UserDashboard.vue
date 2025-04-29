<script setup>
// 引入 Vue 的響應式 API 和生命週期鉤子
import { ref, onMounted } from 'vue';
// 引入 axios 用於發送 HTTP 請求
import axios from 'axios';

// 定義組件的 props
const props = defineProps(['username', 'data']);
// 定義組件的事件
const emit = defineEmits(['logout']);
// 創建健康數據的響應式變量
const healthData = ref([]);
// 從 localStorage 獲取用戶角色
const userRole = ref(localStorage.getItem('role'));

// 組件掛載時執行
onMounted(async () => {
    if (userRole.value === 'user') {
        // 病患角色：根據 id_number 查詢健康檢查資料
        try {
            // 發送請求獲取健康檢查數據
            const response = await axios.get(`https://7aa9-140-124-249-9.ngrok-free.app/default/health-check/user/${props.username}`, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    'Accept': 'application/json'
                }
            });
            // 更新健康數據
            healthData.value = response.data || [];
        } catch (error) {
            // 錯誤處理
            console.error('Fetch health data failed:', error.response?.data || error.message);
            if (error.response?.status === 404) {
                healthData.value = []; // 後端返回 404 表示無資料
            } else if (error.response?.status === 401) {
                alert('認證過期，請重新登入');
                emit('logout');
            } else {
                alert('獲取健康檢查數據失敗：' + (error.response?.data?.detail || error.message));
                healthData.value = [];
            }
        }
    } else if (userRole.value === 'other') {
        // 其他使用者角色：使用原有的 data 函數查詢資料
        healthData.value = await props.data();
    }
});
</script>

<template>
    <!-- 儀表板容器 -->
    <div class="dashboard-container">
        <!-- 根據用戶角色顯示不同的標題 -->
        <h2>{{ userRole === 'user' ? '👤 病人儀表板' : '👥 其他使用者儀表板' }}</h2>
        <!-- 顯示歡迎訊息 -->
        <p>歡迎，{{ username }}</p>
        <!-- 如果有健康數據則顯示 -->
        <div v-if="healthData.length">
            <h3>健康檢查數據</h3>
            <ul>
                <!-- 遍歷並顯示健康數據 -->
                <li v-for="item in healthData" :key="item.id">
                    {{ item.content || item }}
                </li>
            </ul>
        </div>
        <!-- 如果沒有數據則顯示提示訊息 -->
        <p v-else>暫無資料</p>
        <!-- 登出按鈕 -->
        <button @click="$emit('logout')">登出</button>
    </div>
</template>

<style scoped>
/* 儀表板容器樣式 */
.dashboard-container { padding: 20px; }
/* 標題樣式 */
h2 { font-size: 28px; margin-bottom: 20px; }
h3 { font-size: 20px; margin-top: 20px; }
/* 列表樣式 */
ul { list-style-type: none; padding: 0; }
/* 列表項目樣式 */
li { padding: 10px 0; border-bottom: 1px solid #ddd; }
</style>