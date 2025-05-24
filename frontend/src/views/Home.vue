<template>
  <div class="landing-page">
    <!-- Hero Section -->
    <section class="hero-section">
      <v-container>
        <v-row align="center">
          <v-col cols="12" md="6" class="hero-content">
            <v-img src="@/assets/logo.png" alt="logo" max-width="64" class="mb-4" />
            <h1 class="hero-title mb-2">智慧守護，健康鏈結</h1>
            <div class="hero-slogan mb-4">結合 AI 與區塊鏈技術，打造安全可信的健康數據平台</div>
          </v-col>
          <v-col cols="12" md="6" class="text-center">
          </v-col>
        </v-row>
      </v-container>
      <!-- 浮動幾何背景元素 -->
      <div class="floating-bg"></div>
    </section>

    <!-- Features Section -->
    <section class="features-section" id="features">
      <v-container>
        <v-row justify="center" align="stretch" class="feature-row">
          <v-col cols="12" sm="6" md="3" v-for="feature in features" :key="feature.title">
            <v-card class="feature-card" elevation="2">
              <v-card-text class="text-center">
                <v-icon size="40" color="primary" class="mb-3">{{ feature.icon }}</v-icon>
                <h3 class="feature-title mb-2">{{ feature.title }}</h3>
                <p class="feature-desc">{{ feature.desc }}</p>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- About Us Section -->
    <section class="about-section" id="about">
      <v-container>
        <v-row align="center">
          <v-col cols="12" md="6" class="text-center mb-4 mb-md-0">
          </v-col>
          <v-col cols="12" md="6">
            <h2 class="about-title mb-4">關於智療鏈</h2>
            <p class="about-desc mb-3">
              智療鏈致力於結合 AI 與區塊鏈，讓每一位用戶都能安全、便捷地管理個人健康數據。
            </p>
            <p class="about-desc mb-3">
              我們相信，數據的自主權與隱私保護是健康科技的核心。平台以高標準加密與權限控管，確保您的資訊安全無虞。
            </p>
            <p class="about-desc">
              智療鏈團隊持續創新，推動健康數據上鏈，攜手用戶共創智慧健康新未來。
            </p>
          </v-col>
        </v-row>
      </v-container>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue';

const props = defineProps({
  userRole: { type: String, default: null },
  isLoggedIn: { type: Boolean, default: false },
  currentUser: { type: String, default: '' },
  showLoginForm: { type: Boolean, default: false },
  showFooter: { type: Boolean, default: false },
  menuItems: { type: Array, default: () => [] }
});

const emit = defineEmits([
  'login', 'logout', 'forgot-password', 'register', 'go-home', 'toggle-login-form', 'toggle-footer', 'show-snackbar'
]);

const features = ref([
  { icon: 'mdi-shield-check', title: '資料安全', desc: '高等級加密與權限控管，守護您的健康數據隱私。' },
  { icon: 'mdi-chart-timeline-variant', title: '即時記錄', desc: '隨時掌握健康數據，動態追蹤健康趨勢。' },
  { icon: 'mdi-gesture-tap-button', title: '簡單操作', desc: '直覺介面設計，所有年齡層都能輕鬆上手。' },
  { icon: 'mdi-trending-up', title: '健康趨勢掌握', desc: 'AI 智能分析，協助您洞察健康風險與改善建議。' }
]);
</script>

<style scoped>
.landing-page {
  background: #A7E2E3;
  min-height: 100vh;
}
.nav-bar {
  box-shadow: 0 2px 12px 0 rgba(33, 150, 243, 0.08);
  border-radius: 0 0 24px 24px;
}
.brand-title {
  font-size: 1.6rem;
  font-weight: 900;
  color: #A7E2E3;
  letter-spacing: 1.5px;
}
.nav-link {
  font-size: 1.1rem;
  color: #1565c0 !important;
  font-weight: 500;
  letter-spacing: 1px;
}
.hero-section {
  background: linear-gradient(120deg, #A7E2E3 30%, #ffffff 110%);
  position: relative;
  padding: 64px 0 48px 0;
  overflow: hidden;
}
.hero-content {
  z-index: 2;
}
.hero-title {
  font-size: 2.6rem;
  font-weight: 800;
  color: #1565c0;
  letter-spacing: 2px;
  text-shadow: 0 2px 8px rgba(33, 150, 243, 0.08);
}
.hero-slogan {
  font-size: 1.2rem;
  color: #1976d2;
  font-weight: 500;
  letter-spacing: 1px;
}
.hero-cta {
  font-size: 1.15rem;
  font-weight: bold;
  border-radius: 10px;
  box-shadow: 0 2px 8px 0 rgba(33, 150, 243, 0.10);
}
.hero-illustration {
  z-index: 2;
}
.floating-bg {
  position: absolute;
  top: -60px;
  right: -120px;
  width: 420px;
  height: 420px;
  background: radial-gradient(circle at 60% 40%, #42a5f5 0%, #e3f2fd 100%);
  opacity: 0.25;
  border-radius: 50%;
  z-index: 1;
  filter: blur(12px);
}
.features-section {
  background: #f9f7f4;
  padding: 64px 0 48px 0;
}
.feature-row {
  row-gap: 32px;
}
.feature-card {
  border-radius: 20px;
  background: #f8fafc;
  box-shadow: 0 2px 12px 0 rgba(33, 150, 243, 0.10);
  transition: transform 0.3s, box-shadow 0.3s;
  cursor: pointer;
}
.feature-card:hover {
  transform: translateY(-8px) scale(1.03);
  box-shadow: 0 8px 32px 0 rgba(33, 150, 243, 0.18);
}
.feature-title {
  font-size: 1.2rem;
  color: #1565c0;
  font-weight: 700;
}
.feature-desc {
  font-size: 1rem;
  color: #1976d2;
}
.about-section {
  background: linear-gradient(90deg, #f9f7f4 60%);
  padding: 64px 0 48px 0;
}
.about-title {
  font-size: 1.5rem;
  color: #1565c0;
  font-weight: 800;
  margin-bottom: 1.5rem;
}
.about-desc {
  font-size: 1.1rem;
  color: #1976d2;
  line-height: 1.8;
}
.about-illustration {
  border-radius: 20px;
  box-shadow: 0 2px 12px 0 rgba(33, 150, 243, 0.10);
}
@media (max-width: 960px) {
  .hero-title {
    font-size: 2rem;
  }
  .about-title {
    font-size: 1.2rem;
  }
  .feature-title {
    font-size: 1rem;
  }
}
</style>