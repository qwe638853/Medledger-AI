<script setup>
import { ref, onMounted } from 'vue';
import axios from 'axios';

const props = defineProps(['username', 'data']);
const emit = defineEmits(['logout']);
const healthData = ref([]);
const userRole = ref(localStorage.getItem('role')); // 從 localStorage 獲取角色

onMounted(async () => {
    if (userRole.value === 'patient') {
        // 病患角色：根據 id_number 查詢健康檢查資料
        try {
            const response = await axios.get(`https://7aa9-140-124-249-9.ngrok-free.app/health-check/user/${props.username}`, {
                headers: {
                    'Authorization': `Bearer ${localStorage.getItem('token')}`,
                    'Accept': 'application/json'
                }
            });
            healthData.value = response.data || [];
        } catch (error) {
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
    <div class="dashboard-container">
        <h2>{{ userRole === 'patient' ? '👤 病人儀表板' : '👥 其他使用者儀表板' }}</h2>
        <p>歡迎，{{ username }}</p>
        <div v-if="healthData.length">
            <h3>健康檢查數據</h3>
            <ul>
                <li v-for="item in healthData" :key="item.id">
                    <!-- 假設後端返回的資料有 id 和 content 字段 -->
                    {{ item.content || item }}
                </li>
            </ul>
        </div>
        <p v-else>暫無資料</p>
        <button @click="$emit('logout')">登出</button>
    </div>
</template>

<style scoped>
.dashboard-container { padding: 20px; }
h2 { font-size: 28px; margin-bottom: 20px; }
h3 { font-size: 20px; margin-top: 20px; }
ul { list-style-type: none; padding: 0; }
li { padding: 10px 0; border-bottom: 1px solid #ddd; }
</style>