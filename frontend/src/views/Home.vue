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
            <p class="start-desc mt-2">開始管理你的健康數據</p>
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
                </v-card-text>
              </v-card>
            </v-col>
          </v-row>
        </div>
      </v-container>
    </section>

    <!-- Who Should Use Section -->
    <section class="who-section">
      <v-container>
        <h2 class="who-title mb-8 text-center">誰適合使用智慧鏈？</h2>
        <v-row justify="center" align="stretch" class="who-row">
          <v-col cols="12" sm="4" class="mb-4 mb-sm-0">
            <v-card class="who-card" elevation="1">
              <v-card-text class="d-flex flex-column align-center text-center">
                <v-icon size="38" color="#1976D2" class="mb-3">mdi-account-heart</v-icon>
                <div class="who-role">一般民眾</div>
                <div class="who-desc">追蹤與管理自己的健康紀錄。</div>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" sm="4" class="mb-4 mb-sm-0">
            <v-card class="who-card" elevation="1">
              <v-card-text class="d-flex flex-column align-center text-center">
                <v-icon size="38" color="#43AA8B" class="mb-3">mdi-hospital-building</v-icon>
                <div class="who-role">醫療機構</div>
                <div class="who-desc">集中管理病患體檢資料並保持隱私。</div>
              </v-card-text>
            </v-card>
          </v-col>
          <v-col cols="12" sm="4">
            <v-card class="who-card" elevation="1">
              <v-card-text class="d-flex flex-column align-center text-center">
                <v-icon size="38" color="#F9A825" class="mb-3">mdi-shield-account</v-icon>
                <div class="who-role">保險公司</div>
                <div class="who-desc">查閱經授權的歷史健康數據以輔助理賠。</div>
              </v-card-text>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- About Us Section -->
    <section class="about-section" id="about-section">
      <v-container>
        <v-row align="center">
          <v-col cols="12" md="6" class="text-center mb-4 mb-md-0">
          </v-col>
          <v-col cols="12" md="6">
            <h2 class="about-title mb-4">關於智療鏈</h2>
            <p class="about-desc mb-3">
              智療鏈致力於結合 AI 與區塊鏈，讓每一位用戶都能安全、便利的管理個人健康數據。
            </p>
            <p class="about-desc mb-3">
              數據的自主權與隱私保護是健康科技的核心。平台以高標準加密與權限控管，確保您的資訊安全無虞。
            </p>
            <p class="about-desc">
              智療鏈團隊持續創新，推動健康數據上鏈，攜手用戶共創智慧健康新未來。
            </p>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- FAQ Section -->
    <section class="faq-section">
      <div class="faq-container">
        <h2 class="faq-title">常見問題 FAQ</h2>
        <div v-for="(item, idx) in faqList" :key="item.q" class="faq-item">
          <button class="faq-question" @click="toggleFAQ(idx)">
            <span>{{ item.q }}</span>
            <svg :class="['faq-arrow', { 'open': openFAQ === idx }]" width="24" height="24" viewBox="0 0 24 24"><path d="M7 10l5 5 5-5" stroke="#111827" stroke-width="2" fill="none" stroke-linecap="round" stroke-linejoin="round"/></svg>
          </button>
          <transition name="faq-fade">
            <div v-if="openFAQ === idx" class="faq-answer">
              {{ item.a }}
            </div>
          </transition>
        </div>
      </div>
    </section>

    <!-- Contact Section -->
    <section class="contact-section" id="contact-section">
      <v-container>
        <v-row align="center" justify="center">
          <v-col cols="12" md="8" class="text-center">
            <h2 class="contact-title mb-4">聯絡資訊</h2>
            <p class="contact-desc mb-2">Email：t111AB0009@ntut.org.tw</p>
            <p class="contact-desc mb-2">電話：(886-2) 2771-2171</p>
            <p class="contact-desc">地址：台北市忠孝東路三段一號</p>
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

const faqList = [
  {
    q: '我的健康資料會不會外洩？',
    a: '我們使用高階加密與區塊鏈追蹤技術，保障每筆資料的私密性與安全性。'
  },
  {
    q: '需要專業背景才能使用嗎？',
    a: '不需要，平台設計直觀，任何人都能輕鬆上手紀錄與查詢健康資訊。'
  },
  {
    q: '可以和醫院或其他健康平台同步嗎？',
    a: '支援多家合作醫院同步，未來會持續擴充與 Apple Health、Google Fit 的整合。'
  },
  {
    q: '資料可以刪除或下載嗎？',
    a: '您可隨時下載或刪除自己的健康資料，完全掌控個人數據。'
  }
];

const openFAQ = ref(null);
function toggleFAQ(idx) {
  openFAQ.value = openFAQ.value === idx ? null : idx;
}

function scrollToSection(sectionId) {
  const el = document.getElementById(sectionId);
  if (el) {
    el.scrollIntoView({ behavior: 'smooth' });
  }
}
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
  background: #e6e000 !important;
  color: #111827 !important;
  box-shadow: 0 6px 20px rgba(248, 244, 65, 0.25) !important;
  transform: translateY(-2px) scale(1.04);
}

.start-desc {
  font-size: 1rem;
  color: #888;
  margin-top: 0.5rem;
  margin-bottom: 0;
  letter-spacing: 0.5px;
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
  background: #f9f7f4;
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
  font-size: 2rem;
  font-weight: 700;
  color: #111827;
  letter-spacing: -0.5px;
  line-height: 1.3;
}

.feature-desc {
  font-size: 1.25rem;
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
  font-size: 2.7rem;
  color: var(--gray-dark);
  font-weight: 700;
  margin-bottom: 2rem;
  letter-spacing: -0.5px;
}

.about-desc {
  font-size: 1.25rem;
  color: var(--gray-medium);
  line-height: 1.8;
  margin-bottom: 1.5rem;
}

.about-illustration {
  border-radius: 20px;
  box-shadow: 0 2px 12px 0 rgba(33, 150, 243, 0.10);
}

.nav-btn {
  background: #fff !important;
  color: #111827 !important;
  font-weight: 600 !important;
  border-radius: 9999px !important;
  margin-left: 12px !important;
  padding: 12px 28px !important;
  font-size: 1.05rem !important;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04) !important;
  transition: all 0.2s;
}

.nav-btn:hover {
  background: #F8F441 !important;
  color: #333 !important;
}

.contact-section {
  background: #F9F7F4;
  padding: 6rem 0 4rem 0;
}

.contact-title {
  font-family: 'Inter', sans-serif;
  font-size: 2.2rem;
  color: #111827;
  font-weight: 700;
  margin-bottom: 1.5rem;
}

.contact-desc {
  font-size: 1.22rem;
  color: #666;
  margin-bottom: 0.5rem;
}

.who-section {
  background: #f9f7f4;
  padding: 5rem 0 3rem 0;
}
.who-title {
  font-family: 'Inter', sans-serif;
  font-size: 2.3rem;
  color: #111827;
  font-weight: 800;
  letter-spacing: -0.5px;
}
.who-row {
  gap: 24px 0;
}
.who-card {
  border-radius: 24px !important;
  padding: 2rem 1.5rem !important;
  min-height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f9f7f4;
  box-shadow: 0 2px 12px rgba(33, 150, 243, 0.06);
  transition: box-shadow 0.2s;
}
.who-card:hover {
  box-shadow: 0 8px 24px rgba(33, 150, 243, 0.12);
}
.who-role {
  font-size: 1.35rem;
  font-weight: 700;
  color: #1976D2;
  margin-bottom: 0.5rem;
}
.who-desc {
  font-size: 1.18rem;
  color: #666;
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

.faq-section {
  background: #f9f7f4;
  padding: 6rem 0 4rem 0;
}
.faq-container {
  max-width: 720px;
  margin: 0 auto;
  padding: 0 1.5rem;
}
.faq-title {
  font-family: 'Inter', sans-serif;
  font-size: 2.5rem;
  font-weight: 900;
  color: #111827;
  margin-bottom: 2.5rem;
  letter-spacing: -1px;
  text-align: left;
}
.faq-item + .faq-item {
  margin-top: 2rem;
}
.faq-question {
  width: 100%;
  background: none;
  border: none;
  outline: none;
  text-align: left;
  font-size: 1.35rem;
  font-weight: 700;
  color: #111827;
  display: flex;
  align-items: center;
  justify-content: space-between;
  cursor: pointer;
  padding: 0;
  transition: color 0.2s;
}
.faq-question:hover {
  color: #3B82F6;
}
.faq-arrow {
  margin-left: 1rem;
  transition: transform 0.3s cubic-bezier(0.4,0,0.2,1);
}
.faq-arrow.open {
  transform: rotate(180deg);
}
.faq-answer {
  font-size: 1.15rem;
  color: #666;
  font-weight: 400;
  margin-top: 1.2rem;
  line-height: 1.8;
  padding-left: 2px;
  letter-spacing: 0.1px;
}
.faq-fade-enter-active, .faq-fade-leave-active {
  transition: all 0.35s cubic-bezier(0.4,0,0.2,1);
}
.faq-fade-enter-from, .faq-fade-leave-to {
  opacity: 0;
  max-height: 0;
  transform: translateY(-8px);
}
.faq-fade-enter-to, .faq-fade-leave-from {
  opacity: 1;
  max-height: 200px;
  transform: translateY(0);
}
@media (max-width: 600px) {
  .faq-title {
    font-size: 1.7rem;
  }
  .faq-container {
    padding: 0 0.5rem;
  }
  .faq-question {
    font-size: 1.05rem;
  }
  .faq-answer {
    font-size: 0.98rem;
  }
}
</style>