<template>
  <div>
    <!-- Contact Header -->
    <section class="contact-header">
      <v-container class="py-16">
        <v-row justify="center">
          <v-col cols="12" md="8" class="text-center">
            <h1 class="text-h3 font-weight-bold mb-4">
              聯絡我們
            </h1>
            <p class="text-h6">
              我們期待聽到您的聲音，為您提供專業的健康數據服務
            </p>
          </v-col>
        </v-row>
      </v-container>
    </section>

    <!-- Contact Content -->
    <section class="contact-content">
      <v-container class="py-16">
        <v-row>
          <!-- Contact Form -->
          <v-col cols="12" md="6">
            <v-card class="pa-6" elevation="2">
              <h2 class="text-h4 mb-6">發送訊息</h2>
              <v-form @submit.prevent="submitForm">
                <v-text-field
                  v-model="form.name"
                  label="姓名"
                  required
                  variant="outlined"
                  class="mb-4"
                />
                <v-text-field
                  v-model="form.email"
                  label="電子郵件"
                  type="email"
                  required
                  variant="outlined"
                  class="mb-4"
                />
                <v-text-field
                  v-model="form.subject"
                  label="主旨"
                  required
                  variant="outlined"
                  class="mb-4"
                />
                <v-textarea
                  v-model="form.message"
                  label="訊息內容"
                  required
                  variant="outlined"
                  rows="5"
                  class="mb-6"
                />
                <v-btn
                  type="submit"
                  color="primary"
                  size="large"
                  block
                >
                  送出訊息
                </v-btn>
              </v-form>
            </v-card>
          </v-col>

          <!-- Contact Information -->
          <v-col cols="12" md="6">
            <v-card class="pa-6 h-100" elevation="2">
              <h2 class="text-h4 mb-6">聯絡資訊</h2>
              <div class="contact-info">
                <div class="mb-6">
                  <h3 class="text-h6 mb-2">地址</h3>
                  <p class="text-body-1">
                    台北市南港區園區街3號10樓<br>
                    (南港軟體園區F棟)
                  </p>
                </div>
                <div class="mb-6">
                  <h3 class="text-h6 mb-2">電話</h3>
                  <p class="text-body-1">
                    (02) 2653-4401
                  </p>
                </div>
                <div class="mb-6">
                  <h3 class="text-h6 mb-2">電子郵件</h3>
                  <p class="text-body-1">
                    info@healis.com
                  </p>
                </div>
                <div>
                  <h3 class="text-h6 mb-2">營業時間</h3>
                  <p class="text-body-1">
                    週一至週五: 9:00 - 18:00<br>
                    週六、週日: 休息
                  </p>
                </div>
              </div>
            </v-card>
          </v-col>
        </v-row>
      </v-container>
    </section>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const form = ref({
  name: '',
  email: '',
  subject: '',
  message: ''
})

const submitForm = () => {
  // 處理表單提交邏輯
  console.log('Form submitted:', form.value)
}
</script>

<style scoped>
.contact-header {
  background: linear-gradient(135deg, var(--v-primary-base) 0%, var(--v-primary-darken2) 100%);
  color: white;
}

.contact-content {
  background-color: #f8f9fa;
}

.contact-info {
  color: #666;
}

.v-card {
  transition: transform 0.3s ease;
}

.v-card:hover {
  transform: translateY(-5px);
}
</style> 