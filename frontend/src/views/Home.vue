<template>
  <div>
    <!-- Hero Section -->
    <section class="hero-section">
      <v-container class="py-16">
        <v-row align="center">
          <v-col cols="12" md="6">
            <h1 class="text-h2 font-weight-bold mb-6">
              健康檢查數據平台
            </h1>
            <p class="text-h6 mb-8">
              安全管理您的健康數據，讓健康管理更簡單
            </p>
            <v-btn
              color="primary"
              x-large
              rounded
              elevation="2"
              class="mr-4"
              @click="$router.push('/login')"
            >
              立即開始
            </v-btn>
            <v-btn
              text
              x-large
              rounded
              class="primary--text"
              @click="$router.push('/about')"
            >
              了解更多
            </v-btn>
          </v-col>
          <v-col cols="12" md="6">
            <v-img
              src="@/assets/hero-image.jpg"
              class="rounded-lg"
              max-height="500"
              contain
            />
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Features Section -->
    <section class="features-section bg-grey-lighten-4">
      <v-container class="py-16">
        <v-row justify="center" class="mb-12">
          <v-col cols="12" md="8" class="text-center">
            <h2 class="text-h3 font-weight-bold mb-6">
              為什麼選擇我們？
            </h2>
            <p class="text-h6">
              專業的健康檢查數據管理平台，為您的健康把關
            </p>
          </v-col>
        </v-row>

        <v-row>
          <v-col
            v-for="feature in features"
            :key="feature.title"
            cols="12"
            md="4"
          >
            <v-card
              elevation="0"
              class="pa-6 h-100 rounded-lg"
              color="transparent"
            >
              <v-icon
                size="48"
                color="primary"
                class="mb-4"
              >
                {{ feature.icon }}
              </v-icon>
              <h3 class="text-h5 font-weight-bold mb-4">
                {{ feature.title }}
              </h3>
              <p class="text-body-1">
                {{ feature.description }}
              </p>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- CTA Section -->
    <section class="cta-section">
      <v-container class="py-16">
        <v-row justify="center">
          <v-col cols="12" md="8" class="text-center">
            <h2 class="text-h3 font-weight-bold mb-6">
              準備開始使用了嗎？
            </h2>
            <p class="text-h6 mb-8">
              加入數千名用戶的行列，開始管理您的健康數據
            </p>
            <v-btn
              color="white"
              x-large
              rounded
              elevation="2"
              @click="$router.push('/register')"
            >
              立即註冊
            </v-btn>
          </v-col>
        </v-row>
      </v-container>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const features = ref([
  {
    icon: 'mdi-shield-check',
    title: '資料安全',
    description: '採用最高等級的加密技術，確保您的健康數據安全無虞。'
  },
  {
    icon: 'mdi-chart-line',
    title: '數據追蹤',
    description: '即時查看和追蹤您的健康檢查數據，掌握健康趨勢。'
  },
  {
    icon: 'mdi-account-group',
    title: '多方共享',
    description: '安全地與醫療機構共享數據，獲得更好的醫療照護。'
  }
])
</script>

<style scoped>
.hero-section {
  background: linear-gradient(to right, #f8f9fa 50%, #e9ecef 50%);
}

.features-section {
  background-color: #f8f9fa;
}

.cta-section {
  background: linear-gradient(135deg, var(--v-primary-base) 0%, var(--v-primary-darken2) 100%);
  color: white;
}
</style> 