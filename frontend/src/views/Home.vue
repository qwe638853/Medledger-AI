<template>
  <div>
    <!-- Hero Section -->
    <section class="hero-section">
      <v-container class="py-16">
        <v-row align="center">
          <v-col cols="12" md="6">
            <h1 class="display-1 font-weight-bold mb-6 hero-title-xl">
              健康檢查數據平台
            </h1>
            <p class="text-h5 mb-8 hero-subtitle">
              安全、簡單地管理您的健康數據，守護您的健康生活
            </p>
          </v-col>
          <v-col cols="12" md="6">
            <v-img
              src="@/assets/hero-image.jpg"
              class="rounded-lg hero-img"
              max-height="350"
              contain
            />
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Features Section -->
    <section class="features-section">
      <v-container class="py-12">
        <v-row justify="center" align="stretch">
          <v-col cols="12" md="3" v-for="feature in features" :key="feature.title">
            <v-card class="pa-8 text-center feature-card" elevation="2">
              <v-icon size="56" color="info" class="mb-4">{{ feature.icon }}</v-icon>
              <h2 class="text-h5 font-weight-bold mb-3 feature-title">{{ feature.title }}</h2>
              <p class="text-h6 feature-desc">{{ feature.description }}</p>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Vision/Mission Section -->
    <section class="vision-section">
      <v-container class="py-12">
        <v-row justify="center">
          <v-col cols="12" md="8">
            <v-card class="pa-8 text-center vision-card" elevation="1">
              <h2 class="text-h4 font-weight-bold mb-4 vision-title">平台願景與任務</h2>
              <p class="text-h6 vision-desc mb-0">
                我們致力於打造一個安全、易用的健康數據管理平台，協助每一位用戶輕鬆管理個人健康資訊，促進健康生活，並推動台灣健康大數據的永續發展。
              </p>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <v-container class="py-12 text-center">
        <h2 class="text-h4 font-weight-bold mb-6 cta-title">立即加入我們，開始守護您的健康！</h2>
        <v-btn color="info" class="cta-btn mr-4" size="x-large" to="/register">註冊</v-btn>
        <v-btn color="primary" class="cta-btn" size="x-large" to="/login">登入</v-btn>
      </v-container>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
    userRole: {
        type: String,
        default: null
    },
    isLoggedIn: {
        type: Boolean,
        default: false
    },
    currentUser: {
        type: String,
        default: ''
    },
    showLoginForm: {
        type: Boolean,
        default: false
    },
    showFooter: {
        type: Boolean,
        default: false
    },
    menuItems: {
        type: Array,
        default: () => []
    }
});

const emit = defineEmits([
    'login',
    'logout',
    'forgot-password',
    'register',
    'go-home',
    'toggle-login-form',
    'toggle-footer'
]);

const features = ref([
  {
    icon: 'mdi-shield-check',
    title: '數據安全',
    description: '採用高等級加密技術，確保您的健康數據安全無虞。'
  },
  {
    icon: 'mdi-chart-line',
    title: '即時追蹤',
    description: '隨時查看與追蹤您的健康檢查數據，掌握健康趨勢。'
  },
  {
    icon: 'mdi-brain',
    title: '智能分析',
    description: '運用智能分析工具，協助您理解健康數據，預防疾病。'
  },
  {
    icon: 'mdi-account-circle',
    title: '簡易操作',
    description: '介面簡單明瞭，適合所有年齡層輕鬆上手。'
  }
]);
</script>

<style scoped>
.hero-section {
  background: linear-gradient(135deg, #e3f2fd 60%, #ffffff 100%);
}
.hero-title-xl {
  font-size: 4rem;
  letter-spacing: 2px;
}
.hero-subtitle {
  font-size: 1.5rem;
  color: #1976d2;
}
.hero-btn {
  font-size: 1.25rem;
  font-weight: bold;
  border-radius: 8px;
}
.hero-img {
  box-shadow: 0 4px 24px 0 rgba(33, 150, 243, 0.15);
}
.features-section {
  background-color: #f4faff;
}
.feature-card {
  border-radius: 18px;
  background: #fff;
  box-shadow: 0 2px 8px 0 rgba(33, 150, 243, 0.08);
  transition: box-shadow 0.3s;
}
.feature-card:hover {
  box-shadow: 0 8px 24px 0 rgba(33, 150, 243, 0.18);
}
.feature-title {
  font-size: 1.5rem;
}
.feature-desc {
  font-size: 1.15rem;
  color: #1976d2;
}
.vision-section {
  background: #e3f2fd;
}
.vision-card {
  border-radius: 18px;
  background: #fff;
  box-shadow: 0 2px 8px 0 rgba(33, 150, 243, 0.08);
}
.vision-title {
  font-size: 1.5rem;
}
.vision-desc {
  font-size: 1.15rem;
  color: #1976d2;
}
.cta-section {
  background: linear-gradient(135deg, #e3f2fd 60%, #ffffff 100%);
}
.cta-title {
  font-size: 1.5rem;
}
.cta-btn {
  font-size: 1.25rem;
  font-weight: bold;
  border-radius: 8px;
}
</style>