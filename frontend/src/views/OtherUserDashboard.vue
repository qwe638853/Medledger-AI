<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import axios from 'axios';

const authStore = useAuthStore();
const currentUser = computed(() => authStore.currentUser);
const token = ref(authStore.token);
const healthData = ref([]);
const loading = ref(false);

// AI 分析相關
const llmLoading = ref(false);
const llmSummary = ref('');

// ===================== AI 健康指標儀表板（簡易版）這邊也是先塞假資料 =====================
const showGauges = ref(false);
const healthGauges = ref([
  { label: 'BMI', value: 25.9, min: 10, max: 40, unit: '', color: 'orange', suggestion: '建議多運動' },
  { label: '體脂率', value: 23.8, min: 5, max: 40, unit: '%', color: 'green', suggestion: '體脂率正常' },
  { label: '腰圍', value: 88.0, min: 50, max: 120, unit: 'cm', color: 'orange', suggestion: '腰圍偏高' },
  { label: '腰臀圍比', value: 0.86, min: 0.6, max: 1.2, unit: '', color: 'green', suggestion: '正常' },
  { label: '安靜心率', value: 56, min: 40, max: 120, unit: '', color: 'green', suggestion: '心率良好' },
  { label: '最大攝氧量', value: 41, min: 20, max: 60, unit: '', color: 'green', suggestion: '運動能力佳' },
  { label: '收縮壓', value: 118, min: 80, max: 180, unit: 'mmHg', color: 'green', suggestion: '血壓正常' },
  { label: '空腹血糖', value: 90, min: 60, max: 200, unit: 'mg/dL', color: 'green', suggestion: '血糖正常' },
  { label: '總膽固醇', value: 191, min: 100, max: 300, unit: 'mg/dL', color: 'orange', suggestion: '膽固醇偏高' },
  { label: '尿酸', value: 6.3, min: 3, max: 10, unit: 'mg/dL', color: 'green', suggestion: '尿酸正常' },
  { label: '全天壓力', value: 2, min: 0, max: 10, unit: '', color: 'green', suggestion: '壓力低' },
  { label: '睡眠品質', value: 3, min: 0, max: 5, unit: '', color: 'green', suggestion: '睡眠良好' },
  { label: '血氧濃度', value: 94.7, min: 80, max: 100, unit: '%', color: 'green', suggestion: '血氧正常' },
]);

onMounted(async () => {
  loading.value = true;
  try {
    // 取得授權健康檢查數據
    const response = await axios.get(
      `https://7aa9-140-124-249-9.ngrok-free.app/default/health-check/authorized/${currentUser.value}`,
      {
        headers: {
          Authorization: `Bearer ${token.value}`,
          Accept: 'application/json'
        },
        timeout: 10000
      }
    );
    healthData.value = (response.data || []).map(report => ({
      id: report.reportId || report.id,
      content: report.content || JSON.stringify(report.testResults),
      date: report.timestamp || report.date || new Date().toISOString()
    }));
    // ===================== 假資料區塊 =====================
    if (!healthData.value.length) {
      healthData.value = [
        {
          id: 1,
          content: "身高: 175 cm, 體重: 70 kg, 身體質量指數(BMI): 22.9, 血壓: 118/78 mmHg, 心率: 68 bpm, 血氧: 97.2%, 睡眠品質: 5, 空腹血糖: 88 mg/dL, 總膽固醇: 190 mg/dL, 尿酸: 6.0 mg/dL",
          date: "2025-05-15"
        }
      ];
    }
    // ===================== 假資料區塊結束 =====================
  } catch (error) {
    healthData.value = [
        {
          id: 1,
          content: "身高: 175 cm, 體重: 70 kg, 身體質量指數(BMI): 22.9, 血壓: 118/78 mmHg, 心率: 68 bpm, 血氧: 97.2%, 睡眠品質: 5, 空腹血糖: 88 mg/dL, 總膽固醇: 190 mg/dL, 尿酸: 6.0 mg/dL",
          date: "2025-05-15"
        }
      ];
  } finally {
    loading.value = false;
  }
});

const handleLogout = () => {
  authStore.logout();
};

// 處理 LLM 分析
const handleLLMSummary = async () => {
  if (!healthData.value.length) {
    return;
  }
  llmLoading.value = true;
  try {
    await new Promise(resolve => setTimeout(resolve, 1000));
    llmSummary.value = "【AI 健康摘要】這邊之後再連接LLM去讀取分析";
    showGauges.value = true;
  } finally {
    llmLoading.value = false;
  }
};
</script>

<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="10" md="8" lg="7">
        <v-card class="pa-6 mb-6" elevation="8">
          <v-row align="center" justify="space-between" class="mb-4">
            <v-col cols="8">
              <h2 class="mb-1">👥 其他使用者儀表板</h2>
              <div class="subtitle-1">歡迎，{{ currentUser }}</div>
            </v-col>
            <v-col cols="4" class="d-flex justify-end align-center">
              <v-btn color="primary" @click="handleLogout" elevation="2">登出</v-btn>
            </v-col>
          </v-row>
        </v-card>

        <!-- 健檢報告列表 -->
        <v-card class="pa-4 mb-6" elevation="4">
          <h3 class="mb-3">健康檢查報告</h3>
          <v-data-table
            :headers="[
              { text: '報告編號', value: 'id', width: 120 },
              { text: '內容', value: 'content', width: 300 },
              { text: '日期', value: 'date', width: 150 }
            ]"
            :items="healthData"
            :loading="loading"
            loading-text="資料載入中..."
            class="elevation-0"
            dense
            hide-default-footer
            :no-data-text="'暫無資料'"
          >
            <template #item.content="{ item }">
              {{ item.content || item }}
            </template>
            <template #item.date="{ item }">
              {{ item.date || '-' }}
            </template>
          </v-data-table>
        </v-card>

        <!-- LLM 分析區塊 -->
        <v-card class="pa-4 mb-6" elevation="4">
          <h3 class="mb-3">AI 健檢摘要分析</h3>
          <v-btn
            color="deep-purple accent-4"
            :loading="llmLoading"
            :disabled="llmLoading || !healthData.length"
            @click="handleLLMSummary"
            elevation="2"
            class="mb-3"
          >
            <v-icon left>mdi-robot</v-icon> 產生 AI 摘要
          </v-btn>
          <v-alert
            v-if="llmSummary"
            type="info"
            class="mt-3"
            border="left"
            colored-border
            elevation="1"
          >
            {{ llmSummary }}
          </v-alert>
          <div v-if="showGauges" class="mt-6">
            <v-row>
              <v-col v-for="g in healthGauges" :key="g.label" cols="6" md="3" class="text-center mb-4">
                <v-progress-circular
                  :value="((g.value - g.min) / (g.max - g.min)) * 100"
                  :size="90"
                  :width="12"
                  :color="g.color"
                >
                  <span style="font-size:1.2em">{{ g.value }}</span>
                </v-progress-circular>
                <div class="mt-2 font-weight-bold">{{ g.label }}</div>
                <div class="text-caption">{{ g.suggestion }}</div>
              </v-col>
            </v-row>
          </div>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
.fill-height {
  min-height: 100vh;
  background: #f5f6fa;
}
h2 {
  font-size: 28px;
  font-weight: bold;
}
h3 {
  font-size: 20px;
  margin-top: 0;
  font-weight: 500;
}
.subtitle-1 {
  color: #666;
}
</style>