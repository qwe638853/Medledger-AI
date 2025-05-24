<template>
  <div class="landing-page">
    <!-- Hero Section -->
    <section class="hero-section-codingbio">
      <v-container>
        <v-row align="center">
          <v-col cols="12" md="6" class="hero-left d-flex flex-column align-center justify-center">
            <h1 class="hero-title-codingbio mb-4">智慧守護，鏈接健康未來</h1>
            <div class="hero-desc-codingbio mb-6">
              結合 AI 與區塊鏈技術，打造安全可信的健康數據平台
            </div>
            <v-btn color="primary" class="hero-btn-codingbio" large :to="{ path: '/register' }">立即使用</v-btn>
          </v-col>
          <v-col cols="12" md="6" class="hero-right">
            <!-- 幾何圖形SVG或插圖，或用背景圖 -->
            <div class="hero-geometry-bg"></div>
          </v-col>
        </v-row>
      </v-container>
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
.hero-section-codingbio {
  background: #f9f7f4;
  min-height: 520px;
  display: flex;
  align-items: center;
  position: relative;
  overflow: hidden;
}
.hero-title-codingbio {
  font-size: 3rem;
  font-weight: 800;
  color: #222;
  letter-spacing: 1.5px;
  line-height: 1.1;
}
.hero-desc-codingbio {
  font-size: 1.25rem;
  color: #6b7280;
  font-weight: 400;
  margin-bottom: 2rem;
}
.hero-btn-codingbio {
  background: #f7e85a !important;
  color: #222 !important;
  font-weight: 700;
  border-radius: 8px;
  box-shadow: none;
  padding: 0.8rem 2.5rem;
  font-size: 1.1rem;
  letter-spacing: 1px;
}
.hero-right {
  position: relative;
  min-height: 320px;
}
.hero-geometry-bg {
  width: 100%;
  height: 420px;
  background: url('@/assets/geometry-bg.svg') right center/contain no-repeat;
  /* 你可以用SVG或PNG，或直接寫SVG在這裡 */
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
  color: #080808;
  font-weight: 700;
}
.feature-desc {
  font-size: 1rem;
  color: #080808;
}
.about-section {
  background: linear-gradient(90deg, #f9f7f4 60%);
  padding: 64px 0 48px 0;
}
.about-title {
  font-size: 1.5rem;
  color: #080808;
  font-weight: 800;
  margin-bottom: 1.5rem;
}
.about-desc {
  font-size: 1.1rem;
  color: #080808;
  line-height: 1.8;
}
.about-illustration {
  border-radius: 20px;
  box-shadow: 0 2px 12px 0 rgba(33, 150, 243, 0.10);
}
.hero-left {
  min-height: 420px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}
@media (max-width: 960px) {
  .hero-title-codingbio { font-size: 2rem; }
  .hero-geometry-bg { height: 220px; }
  .about-title {
    font-size: 1.2rem;
  }
  .feature-title {
    font-size: 1rem;
  }
}
</style>