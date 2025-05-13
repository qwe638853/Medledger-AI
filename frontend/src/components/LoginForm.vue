<template>
  <v-container class="fill-height">
    <v-row align="center" justify="center">
      <v-col cols="12" sm="8" md="4">
        <v-card class="elevation-12">
          <v-toolbar color="primary" dark flat>
            <v-toolbar-title>登入</v-toolbar-title>
          </v-toolbar>
          <v-card-text>
            <v-form @submit.prevent="handleSubmit">
              <!-- 選擇角色 -->
              <v-select
                v-model="selectedRole"
                :items="roles"
                label="選擇角色"
                outlined
                dense
                prepend-icon="mdi-account-group"
                :rules="[v => !!v || '請選擇角色']"
                class="mb-4"
                style="width: 100%; max-width: 100%;"
              />

              <!-- 性別（假設存在，根據圖片添加） -->
              <v-select
                v-model="gender"
                :items="genders"
                label="性別"
                outlined
                dense
                prepend-icon="mdi-gender-male-female"
                :rules="[v => !!v || '請選擇性別']"
                class="mb-4"
                style="width: 100%; max-width: 100%;"
              />

              <!-- 用戶名 -->
              <v-text-field
                v-model="username"
                label="身分證號/員工編號"
                prepend-icon="mdi-account"
                type="text"
                outlined
                dense
                :rules="[v => !!v || '請輸入身分證號/員工編號']"
                class="mb-4"
                style="width: 100%; max-width: 100%;"
              />

              <!-- 密碼 -->
              <v-text-field
                v-model="password"
                label="密碼"
                prepend-icon="mdi-lock"
                :type="showPassword ? 'text' : 'password'"
                :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append="showPassword = !showPassword"
                outlined
                dense
                :rules="[v => !!v || '請輸入密碼', v => (v && v.length >= 6) || '密碼至少6位']"
                class="mb-4"
                style="width: 100%; max-width: 100%;"
              />

              <!-- 提交按鈕 -->
              <v-btn
                :loading="loading"
                color="primary"
                block
                type="submit"
                class="mb-4"
              >
                登入
              </v-btn>

              <!-- 其他導航按鈕 -->
              <v-btn text @click="goToHome" class="mr-4">返回首頁</v-btn>
              <v-btn text @click="goToRegister">註冊</v-btn>
            </v-form>
          </v-card-text>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import { useAuth } from '../composables/useAuth';

const router = useRouter();
const { login, loading } = useAuth();

const username = ref('');
const password = ref('');
const showPassword = ref(false);
const selectedRole = ref('');
const gender = ref(''); // 添加性別字段
const roles = ['健檢中心', '使用者', '其他使用者'];
const genders = ['男', '女', '其他']; // 假設的性別選項

const handleSubmit = async () => {
  if (!username.value || !password.value || !selectedRole.value || !gender.value) {
    document.dispatchEvent(new CustomEvent('show-snackbar', {
      detail: { message: '請填寫所有欄位', color: 'error' }
    }));
    return;
  }

  const roleMap = {
    '健檢中心': 'medical',
    '使用者': 'user',
    '其他使用者': 'other'
  };

  const role = roleMap[selectedRole.value];
  try {
    await login({
      username: username.value,
      password: password.value,
      role: role,
      gender: gender.value // 假設後端需要性別
    });
  } catch (error) {
    document.dispatchEvent(new CustomEvent('show-snackbar', {
      detail: { message: '登入失敗，請檢查帳號密碼或角色', color: 'error' }
    }));
  }
};

const goToHome = () => router.push('/');
const goToRegister = () => router.push('/register');
</script>

<style scoped>
.fill-height {
  min-height: calc(100vh - 64px);
}
.v-card {
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
}
.v-btn {
  text-transform: none;
  letter-spacing: 0;
  transition: transform 0.2s ease;
}
.v-btn:hover {
  transform: scale(1.05);
}
.v-form {
  padding: 16px;
}
.mb-4 {
  margin-bottom: 24px !important;
}
</style>