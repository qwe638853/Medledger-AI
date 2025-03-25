<script setup>
import { ref } from 'vue';

const id_number = ref('');
const file = ref(null);
const message = ref('');
const emit = defineEmits(['upload', 'logout']);

const handleSubmit = (e) => {
    e.preventDefault();

    if (!id_number.value) {
        message.value = '請輸入病患身分證號';
        return;
    }

    if (!file.value) {
        message.value = '請選擇一個檔案';
        return;
    }

    const formData = new FormData();
    formData.append("file", file.value);

    emit('upload', { id_number: id_number.value, formData });
    id_number.value = '';
    file.value = null;
    message.value = '';
};
</script>

<template>
    <div class="dashboard-container">
        <h2>🏥 醫院員工儀表板</h2>
        <p>歡迎，{{ username }}</p>
        <form id="upload-form" @submit="handleSubmit">
            <div class="form-group">
                <label>病患身分證號</label>
                <input v-model="id_number" type="text" placeholder="輸入病患身分證號" />
            </div>
            <div class="form-group">
                <label>上傳健康檢查檔案</label>
                <input type="file" @change="e => file = e.target.files[0]" />
            </div>
            <button type="submit">上傳</button>
        </form>
        <p v-if="message" :class="{ 'success': message.includes('成功'), 'error': message.includes('失敗') }">
            {{ message }}
        </p>
        <button @click="$emit('logout')">登出</button>
    </div>
</template>

<style scoped>
.dashboard-container { padding: 20px; }
h2 { font-size: 28px; margin-bottom: 20px; }
.form-group { margin-bottom: 20px; }
.form-group label { display: block; margin-bottom: 5px; font-weight: bold; }
.success { color: green; }
.error { color: red; }
</style>