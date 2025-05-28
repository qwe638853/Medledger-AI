<template>
  <div class="landing-page">
    <!-- Hero Section -->
    <section class="hero-section">
      <v-container>
        <v-row align="center" justify="center">
          <v-col cols="12" md="8" class="text-center">
            <h1 class="hero-title mb-6">MedledgerAI</h1>
            <h1 class="hero-title mb-6">智慧守護，鏈接健康未來</h1>
            <p class="hero-desc mb-8">
              結合 AI 與區塊鏈技術，打造安全可信的健康數據平台
            </p>
            <div class="hero-btn-group">
              <v-btn 
                class="start-btn mr-btn"
                elevation="0"
                :to="{ path: '/login' }"
              >
                直接登入
                <v-icon right class="ml-2">mdi-login</v-icon>
              </v-btn>
            <v-btn 
              class="start-btn" 
              elevation="0"
              :to="{ path: '/register' }"
            >
                立即註冊
              <v-icon right class="ml-2">mdi-arrow-right</v-icon>
            </v-btn>
            </div>
            
            <!-- 區塊鏈動畫 -->
            <div class="blockchain-animation">
              <div class="chain-container">
                <div v-for="n in 5" :key="n" class="block-item" :style="{ animationDelay: `${(n-1) * 0.2}s` }">
                  <div class="block">
                    <div class="block-content">
                      <v-icon size="32" :color="blockIcons[n-1].color">{{ blockIcons[n-1].icon }}</v-icon>
                    </div>
                    <div class="block-hash">
                      <span class="hash-text">{{ generateHash() }}</span>
                    </div>
                  </div>
                  <div class="chain-link" v-if="n < 5">
                    <div class="link-dot"></div>
                    <div class="link-line"></div>
                    <div class="link-dot"></div>
                  </div>
                </div>
              </div>
              
              <!-- 浮動元素 -->
              <div class="floating-elements">
                <div v-for="n in 6" :key="`float-${n}`" 
                     class="float-item"
                     :class="`float-${n}`">
                  <v-icon size="24" :color="floatingIcons[n-1].color">
                    {{ floatingIcons[n-1].icon }}
                  </v-icon>
                </div>
              </div>
            </div>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Features Section -->
    <section class="features-section" id="features">
      <v-container>
        <div class="features-two-col">
          <div class="features-list-col">
            <div class="feature-list">
              <div v-for="(feature, index) in features" :key="feature.title" class="feature-list-item">
                <v-icon :size="36" :color="feature.iconColor" class="feature-list-icon">{{ feature.icon }}</v-icon>
                <div class="feature-list-text">
                  <div class="feature-list-title">{{ feature.title }}</div>
                  <div class="feature-list-desc">{{ feature.desc }}</div>
                </div>
              </div>
            </div>
          </div>
          <div class="features-illustration-col">
            <img :src="featureImg" alt="醫療科技插圖" class="features-illustration" />
                  </div>
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
            <img :src="aboutImg" alt="關於智療鏈" class="about-illustration" />
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
import { ref, computed } from 'vue';
import { useRouter } from 'vue-router';
import aboutImg from '@/assets/image.jpg';
import featureImg from '@/assets/image2.svg';

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
    desc: '高等級加密與權限控管，守護您的健康數據隱私。',
    iconColor: '#4CAF50'
  },
  { 
    icon: 'mdi-clock-outline',
    title: '即時評估', 
    desc: '隨時掌握健康數據，動態追蹤健康趨勢。',
    iconColor: '#2196F3'
  },
  { 
    icon: 'mdi-gesture-tap-button',
    title: '簡單操作', 
    desc: '直覺介面設計，所有年齡層都能輕鬆上手。',
    iconColor: '#FF9800'
  },
  { 
    icon: 'mdi-trending-up', 
    title: '健康趨勢掌握', 
    desc: 'AI 智能分析，協助您洞察健康風險與改善建議。',
    iconColor: '#9C27B0'
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

// 區塊鏈動畫相關數據
const blockIcons = ref([
  { icon: 'mdi-hospital-building', color: '#4CAF50' },
  { icon: 'mdi-account-heart', color: '#2196F3' },
  { icon: 'mdi-shield-check', color: '#FFC107' },
  { icon: 'mdi-database-lock', color: '#9C27B0' },
  { icon: 'mdi-medical-bag', color: '#FF5252' }
]);

const floatingIcons = ref([
  { icon: 'mdi-heart-pulse', color: '#FF5252' },
  { icon: 'mdi-dna', color: '#2196F3' },
  { icon: 'mdi-pill', color: '#4CAF50' },
  { icon: 'mdi-hospital-box', color: '#FFC107' },
  { icon: 'mdi-shield-check', color: '#9C27B0' },
  { icon: 'mdi-chart-line', color: '#FF9800' }
]);

// 生成隨機哈希值
function generateHash() {
  const chars = '0123456789abcdef';
  let hash = '0x';
  for (let i = 0; i < 8; i++) {
    hash += chars[Math.floor(Math.random() * chars.length)];
  }
  return hash;
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
  font-size: clamp(2.5rem, 5vw, 4.5rem);
  font-weight: 900;
  color: #111827;
  letter-spacing: -1.5px;
  line-height: 1.1;
  max-width: 1200px;
  margin: 0 auto;
}

.hero-desc {
  font-size: clamp(1.5rem, 2vw, 2rem);
  color: #6b7280;
  font-weight: 400;
  max-width: 800px;
  margin: 0 auto;
  line-height: 1.6;
}

.start-btn {
  background: #F8F441 !important;
  color: #111827 !important;
  font-weight: 700 !important;
  font-size: 1.6rem !important;
  border-radius: 48px !important;
  padding: 22px 56px !important;
  height: auto !important;
  min-width: 220px;
  min-height: 64px;
  box-shadow: 0 4px 24px rgba(248, 244, 65, 0.13) !important;
}

.start-btn .v-icon {
  font-size: 2.1rem !important;
  margin-left: 1.2rem !important;
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
  padding: clamp(4rem, 8vw, 10rem) 0;
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
  border-radius: 8px;
  background: linear-gradient(to bottom, #ffffff, #f5f5f5);
  padding: 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  border: 1px solid #e5e7eb;
}

.feature-card:hover {
  transform: scale(1.02);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.feature-title {
  font-family: sans-serif;
  font-size: 1.5rem;
  font-weight: 600;
  color: #111827;
  letter-spacing: -0.5px;
  line-height: 1.3;
  margin-bottom: 0;
}

.feature-desc {
  font-family: sans-serif;
  font-size: 1.1rem;
  color: #666;
  line-height: 1.6;
  margin-left: 24px;
  max-width: 340px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  text-align: left;
}

.feature-icon {
  margin-right: 12px;
}

.about-section {
  background: var(--background-cream);
  padding: clamp(4rem, 8vw, 10rem) 0;
}

.about-title {
  font-family: 'Inter', sans-serif;
  font-size: clamp(2.2rem, 3vw, 3.2rem);
  color: var(--gray-dark);
  font-weight: 700;
  margin-bottom: 2rem;
  letter-spacing: -0.5px;
}

.about-desc {
  font-size: clamp(1.1rem, 1.3vw, 1.4rem);
  color: var(--gray-medium);
  line-height: 1.8;
  margin-bottom: 1.5rem;
  max-width: 800px;
}

.about-illustration {
  max-width: clamp(420px, 40vw, 600px);
  width: 100%;
  height: auto;
  border-radius: 20px;
  box-shadow: 0 2px 12px 0 rgba(33, 150, 243, 0.10);
  margin: 0 auto;
  display: block;
}

.nav-btn {
  background: #fff !important;
  color: #111827 !important;
  font-weight: 600 !important;
  border-radius: 9999px !important;
  margin-left: 12px !important;
  padding: 16px 32px !important;
  font-size: 1.2rem !important;
  box-shadow: 0 2px 8px rgba(0,0,0,0.04) !important;
  transition: all 0.2s;
  height: auto !important;
  min-height: 50px;
}

.nav-btn:hover {
  background: #F8F441 !important;
  color: #333 !important;
}

.contact-section {
  background: #F9F7F4;
  padding: clamp(4rem, 6vw, 8rem) 0;
}

.contact-title {
  font-size: clamp(2rem, 2.5vw, 2.8rem);
  max-width: 1200px;
  margin: 0 auto 2rem;
}

.contact-desc {
  font-size: clamp(1.1rem, 1.2vw, 1.3rem);
  max-width: 800px;
  margin: 0 auto 0.5rem;
}

.who-section {
  background: #f9f7f4;
  padding: clamp(3rem, 6vw, 8rem) 0;
}
.who-title {
  font-family: 'Inter', sans-serif;
  font-size: clamp(2rem, 3vw, 2.8rem);
  color: #111827;
  font-weight: 800;
  letter-spacing: -0.5px;
  max-width: 1200px;
  margin: 0 auto 2rem;
}
.who-row {
  gap: 24px;
  max-width: 1400px;
  margin: 0 auto;
}
.who-card {
  border-radius: 24px !important;
  padding: clamp(2rem, 3vw, 3rem) 1.5rem !important;
  min-height: 220px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f9f7f4;
  box-shadow: 0 2px 12px rgba(33, 150, 243, 0.06);
  transition: all 0.3s ease;
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

@media (min-width: 1921px) {
  .hero-section,
  .features-section,
  .who-section,
  .about-section,
  .faq-section,
  .contact-section {
    padding-left: calc((100vw - 1920px) / 2);
    padding-right: calc((100vw - 1920px) / 2);
  }
  
  .features-two-col,
  .who-row,
  .about-illustration,
  .faq-container {
    max-width: 1920px;
    margin-left: auto;
    margin-right: auto;
  }
}

@media (max-width: 1440px) {
  .features-two-col {
    padding: 0 1.5rem;
  }
  
  .features-illustration-col {
    max-width: 700px;
  }
}

@media (max-width: 960px) {
  .hero-section {
    padding: 4rem 1rem;
  }

  .features-two-col {
    flex-direction: column;
    align-items: stretch;
    gap: 2.5rem;
    padding: 0 1rem;
  }

  .features-illustration-col {
    max-width: 100%;
  }

  .who-card {
    padding: 1.5rem !important;
  }
}

@media (max-width: 600px) {
  .hero-section {
    padding: 3rem 1rem;
  }

  .features-section,
  .who-section,
  .about-section,
  .faq-section,
  .contact-section {
    padding: 3rem 1rem;
  }

  .features-two-col {
    padding: 0 0.5rem;
  }

  .nav-btn {
    font-size: 1rem !important;
    padding: 10px 20px !important;
    min-height: 40px;
  }
}

.faq-section {
  background: #f9f7f4;
  padding: clamp(4rem, 8vw, 10rem) 0;
}
.faq-container {
  max-width: 1000px;
  margin: 0 auto;
  padding: 0 clamp(1.5rem, 3vw, 3rem);
}
.faq-title {
  font-size: clamp(2rem, 3vw, 2.8rem);
  max-width: 1200px;
  margin: 0 auto 3rem;
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

.feature-info-row {
  display: flex;
  align-items: center;
  min-height: 80px;
}

.feature-list {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin: 2rem 0 3rem 0;
}
.feature-list-item {
  display: flex;
  align-items: flex-start;
  gap: 1.25rem;
  padding: 1.5rem 1.8rem;
  border-radius: 12px;
  background: transparent;
  transition: background 0.18s, box-shadow 0.18s, transform 0.18s;
  cursor: pointer;
}
.feature-list-item:hover {
  background: #f3f3f3;
}
.feature-list-icon {
  min-width: 48px;
  min-height: 48px;
  transition: transform 0.18s;
}
.feature-list-item:hover .feature-list-icon {
  transform: translateX(4px);
}
.feature-list-text {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
}
.feature-list-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
  margin-bottom: 0.5rem;
}
.feature-list-desc {
  font-size: 1.25rem;
  color: #4B5563;
  font-weight: 400;
  line-height: 1.6;
}

@media (max-width: 768px) {
  .feature-list-item {
    padding: 1.25rem;
  }
  
  .feature-list-icon {
    min-width: 40px;
    min-height: 40px;
  }
  
  .feature-list-title {
    font-size: 1.35rem;
  }
  
  .feature-list-desc {
    font-size: 1.15rem;
  }
}

.features-two-col {
  display: flex;
  flex-direction: row;
  align-items: center;
  justify-content: space-between;
  gap: 3rem;
  max-width: 1600px;
  margin: 0 auto;
  padding: 0 2rem;
}
.features-list-col {
  flex: 1 1 0%;
  max-width: 600px;
}
.features-illustration-col {
  flex: 0 0 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  max-width: 900px;
}
.features-illustration {
  max-width: 100%;
  max-height: 800px;
  width: auto;
  height: auto;
  display: block;
  margin: 0 auto;
  border-radius: 18px;
}
@media (max-width: 900px) {
  .features-two-col {
    flex-direction: column;
    align-items: stretch;
    gap: 2.5rem;
  }
  .features-illustration-col {
    max-width: 100%;
    margin-top: 2rem;
    justify-content: flex-start;
  }
  .features-illustration {
    margin: 0 auto;
  }
}

.hero-btn-group {
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 1.2rem;
  margin-bottom: 0.5rem;
}
.mr-btn {
  margin-right: 0 !important;
}
@media (max-width: 600px) {
  .hero-btn-group {
    flex-direction: column;
    gap: 0.7rem;
  }
}

/* 區塊鏈動畫樣式 */
.blockchain-animation {
  position: relative;
  width: 100%;
  max-width: 1200px;
  height: 400px;
  margin: 4rem auto 0;
  overflow: hidden;
  background: transparent;
  display: flex;
  align-items: center;
  justify-content: center;
}

.chain-container {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  animation: infiniteChain 20s linear infinite;
  width: max-content;
  position: relative;
  z-index: 2;
}

.block-item {
  display: flex;
  align-items: center;
  opacity: 1;
  transform: scale(1);
}

.block {
  background: white;
  border-radius: 16px;
  padding: 1.5rem;
  box-shadow: 0 4px 15px rgba(0, 0, 0, 0.06);
  min-width: 150px;
  transition: all 0.3s ease;
  border: 2px solid #e5e7eb;
  animation: blockBreath 3s ease-in-out infinite;
  backdrop-filter: blur(8px);
}

.block:hover {
  transform: translateY(-5px) scale(1.05);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.08);
  border-color: #2196F3;
  z-index: 1;
}

.block-content {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 0.8rem;
}

.block-content :deep(.v-icon) {
  font-size: 32px;
}

.block-hash {
  font-size: 0.85rem;
  color: #666;
  text-align: center;
  overflow: hidden;
  animation: textBreath 3s ease-in-out infinite;
}

.hash-text {
  display: inline-block;
}

.chain-link {
  display: flex;
  align-items: center;
  margin: 0 1.5rem;
  animation: pulseLink 2s ease-in-out infinite;
  position: relative;
}

.link-dot {
  width: 8px;
  height: 8px;
  background: #2196F3;
  border-radius: 50%;
  position: relative;
}

.link-line {
  width: 40px;
  height: 2px;
  background: #2196F3;
  margin: 0 3px;
  position: relative;
}

.link-dot::after,
.link-line::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: inherit;
  border-radius: inherit;
  animation: flowEffect 2s linear infinite;
}

/* 浮動元素樣式 */
.floating-elements {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  z-index: 1;
}

.float-item {
  position: absolute;
  animation: float 6s ease-in-out infinite;
  opacity: 0.5;
}

.float-item :deep(.v-icon) {
  font-size: 24px;
}

/* 浮動元素位置調整 */
.float-1 { top: 10%; left: 15%; animation-delay: 0s; }
.float-2 { top: 15%; right: 20%; animation-delay: 1s; }
.float-3 { bottom: 20%; left: 25%; animation-delay: 2s; }
.float-4 { bottom: 25%; right: 15%; animation-delay: 3s; }
.float-5 { top: 40%; left: 10%; animation-delay: 4s; }
.float-6 { top: 35%; right: 25%; animation-delay: 5s; }

/* 動畫定義 */
@keyframes fadeInScale {
  to {
    opacity: 1;
    transform: scale(1);
  }
}

@keyframes blockBreath {
  0%, 100% {
    box-shadow: 0 4px 15px rgba(33, 150, 243, 0.06);
    border-color: #e5e7eb;
  }
  50% {
    box-shadow: 0 8px 25px rgba(33, 150, 243, 0.12);
    border-color: #2196F3;
  }
}

@keyframes textBreath {
  0%, 100% {
    opacity: 0.7;
    color: #666;
  }
  50% {
    opacity: 1;
    color: #333;
  }
}

@keyframes float {
  0%, 100% {
    transform: translateY(0) rotate(0deg) scale(0.8);
  }
  50% {
    transform: translateY(-15px) rotate(5deg) scale(0.9);
  }
}

@keyframes pulseLink {
  0%, 100% {
    opacity: 0.5;
    transform: scaleX(0.95);
  }
  50% {
    opacity: 1;
    transform: scaleX(1.05);
  }
}

@keyframes flowEffect {
  0% {
    opacity: 0;
    transform: translateX(-100%);
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0;
    transform: translateX(100%);
  }
}

@keyframes infiniteChain {
  0% {
    transform: translateX(calc(50% + 600px)) scale(1.2);
  }
  100% {
    transform: translateX(calc(-50% - 600px)) scale(1.2);
  }
}

/* RWD 適配 */
@media (max-width: 1200px) {
  .blockchain-animation {
    max-width: 100%;
    padding: 0 1rem;
  }
  
  .chain-container {
    animation: infiniteChain 20s linear infinite;
  }
  
  .block {
    min-width: 130px;
  }
}

@media (max-width: 960px) {
  .blockchain-animation {
    height: 350px;
    padding: 0 0.5rem;
  }
  
  .chain-container {
    animation: infiniteChain 15s linear infinite;
  }
  
  .block {
    min-width: 120px;
    padding: 1rem;
  }
  
  .link-line {
    width: 30px;
  }
}

@media (max-width: 600px) {
  .blockchain-animation {
    height: 300px;
    margin-top: 2rem;
    padding: 0;
  }
  
  .chain-container {
    animation: infiniteChain 12s linear infinite;
  }
  
  .block {
    min-width: 120px;
    padding: 1rem;
  }
  
  .link-line {
    width: 30px;
  }
}
</style>