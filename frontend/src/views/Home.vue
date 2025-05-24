<template>
  <div class="landing-page">
    <!-- Hero Section -->
    <section class="hero-section">
      <v-container>
        <v-row align="center" justify="center">
          <v-col cols="12" md="8" class="text-center">
            <h1 class="hero-title mb-6">智慧守護，鏈接健康未來</h1>
            <p class="hero-desc mb-8">
              結合 AI 與區塊鏈技術，打造安全可信的健康數據平台
            </p>
            <v-btn 
              class="start-btn" 
              elevation="0"
              :to="{ path: '/register' }"
            >
              立即開始
              <v-icon right class="ml-2">mdi-arrow-right</v-icon>
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Features Section -->
    <section class="features-section" id="features">
      <v-container>
        <div class="features-scroll-container">
          <v-row justify="center" align="stretch" class="feature-row">
            <v-col cols="12" sm="6" md="3" v-for="(feature, index) in features" :key="feature.title">
              <v-card class="feature-card" elevation="0">
                <v-card-text class="text-center d-flex flex-column align-center">
                  <div class="icon-wrapper">
                    <v-icon size="40" class="feature-icon">{{ feature.icon }}</v-icon>
                  </div>
                  <h3 class="feature-title mt-8 mb-4">{{ feature.title }}</h3>
                  <p class="feature-desc mb-8">{{ feature.desc }}</p>
                  <div class="learn-more">
                    <span class="learn-more-text">了解更多</span>
                    <span class="learn-more-arrow">→</span>
                  </div>
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </div>
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
  { 
    icon: 'mdi-shield-outline', 
    title: '資料安全', 
    desc: '高等級加密與權限控管，守護您的健康數據隱私。'
  },
  { 
    icon: 'mdi-clock-outline',
    title: '即時記錄', 
    desc: '隨時掌握健康數據，動態追蹤健康趨勢。'
  },
  { 
    icon: 'mdi-gesture-tap-button',
    title: '簡單操作', 
    desc: '直覺介面設計，所有年齡層都能輕鬆上手。'
  },
  { 
    icon: 'mdi-trending-up', 
    title: '健康趨勢掌握', 
    desc: 'AI 智能分析，協助您洞察健康風險與改善建議。'
  }
]);
</script>

<style scoped>
.hero-section {
  background: #F9F7F4;
  min-height: 600px;
  display: flex;
  align-items: center;
  padding: 6rem 0;
}

.hero-title {
  font-family: 'Inter', sans-serif;
  font-size: 4rem;
  font-weight: 900;
  color: #111827;
  letter-spacing: -1.5px;
  line-height: 1.1;
}

.hero-desc {
  font-size: 1.25rem;
  color: #6b7280;
  font-weight: 400;
  max-width: 600px;
  margin: 0 auto;
  line-height: 1.6;
}

.start-btn {
  background: #F8F441 !important;
  color: #111827 !important;
  font-weight: 700 !important;
  font-size: 1.125rem !important;
  border-radius: 9999px !important;
  padding: 12px 32px !important;
  height: auto !important;
  letter-spacing: 0;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.start-btn:hover {
  transform: scale(1.05);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.features-section {
  background: #F9F7F4;
  padding: 8rem 0;
  overflow: hidden;
}

.features-scroll-container {
  overflow-x: auto;
  padding: 1rem 0;
  margin: -1rem 0;
  scrollbar-width: none;
  -ms-overflow-style: none;
}

.features-scroll-container::-webkit-scrollbar {
  display: none;
}

.feature-row {
  row-gap: 48px;
  column-gap: 32px;
  flex-wrap: nowrap;
  padding: 0 1rem;
}

.feature-card {
  border-radius: 32px;
  background: #fff;
  padding: 3rem 2rem;
  height: 100%;
  min-height: 400px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #e5e7eb;
}

.feature-card:hover {
  transform: scale(1.02);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.08);
}

.icon-wrapper {
  width: 80px;
  height: 80px;
  border-radius: 24px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #F3F2EF;
  transition: all 0.3s ease;
}

.feature-card:hover .icon-wrapper {
  background: #E5E7EB;
}

.feature-icon {
  font-size: 36px !important;
  color: #111827 !important;
  transition: all 0.3s ease;
}

.feature-title {
  font-family: 'Inter', sans-serif;
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
  letter-spacing: -0.5px;
  line-height: 1.3;
}

.feature-desc {
  font-size: 1.1rem;
  color: #888888;
  line-height: 1.6;
}

.learn-more {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #111827;
  font-weight: 400;
  transition: all 0.3s ease;
  margin-top: auto;
}

.learn-more-text {
  font-size: 1rem;
}

.learn-more-arrow {
  font-family: system-ui;
  font-size: 1.25rem;
  line-height: 1;
}

.about-section {
  background: var(--background-cream);
  padding: 6rem 0;
}

.about-title {
  font-family: 'Inter', sans-serif;
  font-size: 2.5rem;
  color: var(--gray-dark);
  font-weight: 700;
  margin-bottom: 2rem;
  letter-spacing: -0.5px;
}

.about-desc {
  font-size: 1.125rem;
  color: var(--gray-medium);
  line-height: 1.8;
  margin-bottom: 1.5rem;
}

.about-illustration {
  border-radius: 20px;
  box-shadow: 0 2px 12px 0 rgba(33, 150, 243, 0.10);
}

@media (max-width: 960px) {
  .hero-section {
    padding: 4rem 0;
  }

  .hero-title {
    font-size: 3rem;
  }
  
  .features-section {
    padding: 6rem 0;
  }
  
  .feature-row {
    margin: 0 -1rem;
  }
  
  .feature-card {
    min-width: 300px;
    margin: 0 1rem;
    padding: 2rem;
  }
}

@media (max-width: 600px) {
  .hero-title {
    font-size: 2.5rem;
  }
  
  .hero-desc {
    font-size: 1.1rem;
  }
  
  .features-section {
    padding: 4rem 0;
  }
  
  .feature-card {
    min-width: 280px;
    padding: 2rem 1.5rem;
    min-height: 360px;
  }
  
  .feature-title {
    font-size: 1.25rem;
  }
  
  .feature-desc {
    font-size: 1rem;
  }
  
  .icon-wrapper {
    width: 64px;
    height: 64px;
  }
  
  .feature-icon {
    font-size: 28px !important;
  }
}
</style>