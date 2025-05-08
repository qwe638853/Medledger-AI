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
              <v-select
                v-model="selectedRole"
                :items="roles"
                item-title="text"
                item-value="value"
                label="選擇角色"
                outlined
                :rules="[v => !!v || '請選擇角色']"
                prepend-icon="mdi-account-group"
                @update:model-value="selectRole"
              />
              <v-text-field
                v-model="username"
                label="身分證號/員工編號"
                prepend-icon="mdi-account"
                type="text"
                outlined
                :rules="[v => !!v || '請輸入身分證號/員工編號']"
                class="mb-2"
              />
              <v-text-field
                v-model="password"
                label="密碼"
                prepend-icon="mdi-lock"
                :type="showPassword ? 'text' : 'password'"
                :append-icon="showPassword ? 'mdi-eye-off' : 'mdi-eye'"
                @click:append="showPassword = !showPassword"
                outlined
                :rules="[v => !!v || '請輸入密碼']"
                class="mb-2"
              />
              <v-btn
                :loading="loading"
                color="primary"
                block
                type="submit"
                class="mb-4"
              >
                登入
              </v-btn>
              <v-btn
                text
                @click="emit('go-home')"
                class="mr-4"
              >
                返回首頁
              </v-btn>
              <v-btn
                text
                @click="emit('forgot-password')"
              >
                忘記密碼？
              </v-btn>
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

const emit = defineEmits(['login', 'logout', 'forgot-password', 'go-home']);

const { login, loading } = useAuth();
const router = useRouter();

const selectedRole = ref('');
const username = ref('');
const password = ref('');
const showPassword = ref(false);
const roles = [
  { text: '一般用戶', value: 'general' },
  { text: '醫療機構', value: 'medical' },
  { text: '其他用戶', value: 'other' }
];

const selectRole = (value) => {
  selectedRole.value = value;
};

const handleSubmit = async () => {
  if (!selectedRole.value || !username.value || !password.value) {
    return;
  }
  try {
    await login({
      username: username.value,
      password: password.value,
      role: selectedRole.value
    });
    router.push('/');
  } catch (error) {
    console.error('Login error:', error);
  }
};
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
</style>