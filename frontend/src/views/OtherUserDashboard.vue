<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import axios from 'axios';

const authStore = useAuthStore();
const currentUser = computed(() => authStore.currentUser);
const token = ref(authStore.token);
const clientHealthData = ref([]);
const loading = ref(false);

// AI 分析相關
const llmLoading = ref(false);
const llmSummary = ref('');
const riskAssessment = ref(null);

// ===================== 保險業者客戶健康風險評估 =====================
const showGauges = ref(false);
const healthGauges = ref([
  { label: 'BMI', value: 25.9, min: 10, max: 40, unit: '', color: 'blue-darken-1', suggestion: '體重過重風險', icon: 'mdi-scale-bathroom' },
  { label: '體脂率', value: 23.8, min: 5, max: 40, unit: '%', color: 'green-darken-1', suggestion: '體脂率正常', icon: 'mdi-percent' },
  { label: '腰圍', value: 88.0, min: 50, max: 120, unit: 'cm', color: 'amber-darken-2', suggestion: '代謝症候群風險', icon: 'mdi-tape-measure' },
  { label: '腰臀圍比', value: 0.86, min: 0.6, max: 1.2, unit: '', color: 'green-darken-1', suggestion: '正常範圍', icon: 'mdi-human' },
  { label: '安靜心率', value: 56, min: 40, max: 120, unit: '', color: 'green-darken-1', suggestion: '心臟健康', icon: 'mdi-heart-pulse' },
  { label: '最大攝氧量', value: 41, min: 20, max: 60, unit: '', color: 'green-darken-1', suggestion: '有氧能力佳', icon: 'mdi-run-fast' },
  { label: '收縮壓', value: 118, min: 80, max: 180, unit: 'mmHg', color: 'green-darken-1', suggestion: '血壓風險低', icon: 'mdi-diabetes' },
  { label: '空腹血糖', value: 90, min: 60, max: 200, unit: 'mg/dL', color: 'green-darken-1', suggestion: '糖尿病風險低', icon: 'mdi-water-opacity' },
  { label: '總膽固醇', value: 191, min: 100, max: 300, unit: 'mg/dL', color: 'amber-darken-2', suggestion: '心血管疾病風險', icon: 'mdi-heart-circle' },
  { label: '尿酸', value: 6.3, min: 3, max: 10, unit: 'mg/dL', color: 'green-darken-1', suggestion: '痛風風險低', icon: 'mdi-water-check' },
  { label: '全天壓力', value: 2, min: 0, max: 10, unit: '', color: 'green-darken-1', suggestion: '壓力相關疾病風險低', icon: 'mdi-head-cog' },
  { label: '睡眠品質', value: 3, min: 0, max: 5, unit: '', color: 'green-darken-1', suggestion: '睡眠品質佳', icon: 'mdi-sleep' },
  { label: '血氧濃度', value: 94.7, min: 80, max: 100, unit: '%', color: 'green-darken-1', suggestion: '呼吸系統健康', icon: 'mdi-lungs' },
]);

// 保險風險評估
const insuranceRisks = ref({
  overall: '低風險',
  healthScore: 85,
  recommendations: [
    '適合標準費率人壽保險',
    '建議提供健康促進獎勵計劃',
    '無需額外承保限制'
  ]
});

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
    clientHealthData.value = (response.data || []).map(report => ({
      id: report.reportId || report.id,
      content: report.content || JSON.stringify(report.testResults),
      date: report.timestamp || report.date || new Date().toISOString(),
      clientId: report.clientId || '未指定',
      clientName: report.clientName || '未指定客戶'
    }));
    // ===================== 假資料區塊 =====================
    if (!clientHealthData.value.length) {
      clientHealthData.value = [
        {
          id: 1,
          clientId: 'C10023',
          clientName: '王小明',
          content: "身高: 175 cm, 體重: 70 kg, 身體質量指數(BMI): 22.9, 血壓: 118/78 mmHg, 心率: 68 bpm, 血氧: 97.2%, 睡眠品質: 5, 空腹血糖: 88 mg/dL, 總膽固醇: 190 mg/dL, 尿酸: 6.0 mg/dL",
          date: "2025-05-15"
        },
        {
          id: 2,
          clientId: 'C10045',
          clientName: '林小華',
          content: "身高: 165 cm, 體重: 58 kg, 身體質量指數(BMI): 21.3, 血壓: 122/82 mmHg, 心率: 72 bpm, 血氧: 98.1%, 睡眠品質: 4, 空腹血糖: 92 mg/dL, 總膽固醇: 205 mg/dL, 尿酸: 5.8 mg/dL",
          date: "2025-05-10"
        }
      ];
    }
    // ===================== 假資料區塊結束 =====================
  } catch (error) {
    clientHealthData.value = [
      {
        id: 1,
        clientId: 'C10023',
        clientName: '王小明',
        content: "身高: 175 cm, 體重: 70 kg, 身體質量指數(BMI): 22.9, 血壓: 118/78 mmHg, 心率: 68 bpm, 血氧: 97.2%, 睡眠品質: 5, 空腹血糖: 88 mg/dL, 總膽固醇: 190 mg/dL, 尿酸: 6.0 mg/dL",
        date: "2025-05-15"
      },
      {
        id: 2,
        clientId: 'C10045',
        clientName: '林小華',
        content: "身高: 165 cm, 體重: 58 kg, 身體質量指數(BMI): 21.3, 血壓: 122/82 mmHg, 心率: 72 bpm, 血氧: 98.1%, 睡眠品質: 4, 空腹血糖: 92 mg/dL, 總膽固醇: 205 mg/dL, 尿酸: 5.8 mg/dL",
        date: "2025-05-10"
      }
    ];
  } finally {
    loading.value = false;
  }
});

const handleLogout = () => {
  authStore.logout();
};

// 選中的客戶資料
const selectedClient = ref(null);

// 處理客戶選擇
const handleSelectClient = (client) => {
  selectedClient.value = client;
};

// 處理 LLM 風險評估
const handleLLMRiskAssessment = async () => {
  if (!selectedClient.value) {
    return;
  }
  llmLoading.value = true;
  try {
    await new Promise(resolve => setTimeout(resolve, 1000));
    llmSummary.value = `【AI 風險評估】客戶 ${selectedClient.value.clientName} 的健康狀況良好，建議標準費率承保，無需額外限制。健康指標普遍正常，總膽固醇稍高但在可接受範圍內。`;
    riskAssessment.value = {
      riskLevel: '低風險',
      approvalRecommendation: '建議承保',
      premiumAdjustment: '標準費率',
      conditions: '無特殊條件'
    };
    showGauges.value = true;
  } finally {
    llmLoading.value = false;
  }
};

// 保費試算
const calculatePremium = () => {
  return {
    baseAmount: 12000,
    adjustmentFactor: riskAssessment.value?.riskLevel === '低風險' ? 1.0 : 1.2,
    finalPremium: riskAssessment.value?.riskLevel === '低風險' ? 12000 : 14400
  };
};
</script>

<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="11" md="10" lg="9">
        <!-- 頂部導航卡片 -->
        <v-card class="mb-6 rounded-lg" elevation="3" color="blue-lighten-5">
          <v-row align="center" justify="space-between" no-gutters>
            <v-col cols="8" class="pa-6">
              <div class="d-flex align-center">
                <v-icon size="42" color="blue-darken-2" class="me-4">mdi-shield-account</v-icon>
                <div>
                  <h1 class="text-h5 font-weight-bold text-blue-darken-3 mb-1">保險業者健康風險評估平台</h1>
                  <div class="text-subtitle-1 text-blue-darken-1">歡迎，保險顧問 {{ currentUser }}</div>
                </div>
              </div>
            </v-col>
            <v-col cols="4" class="pa-6 d-flex justify-end align-center">
              <v-btn 
                color="blue-darken-2" 
                @click="handleLogout" 
                elevation="2" 
                prepend-icon="mdi-logout"
                class="font-weight-medium"
              >
                登出系統
              </v-btn>
            </v-col>
          </v-row>
        </v-card>

        <!-- 客戶健康報告列表區塊 -->
        <v-card class="mb-6 rounded-lg" elevation="2">
          <v-card-title class="py-4 px-6 bg-blue-lighten-4">
            <v-icon size="24" color="blue-darken-3" class="me-2">mdi-account-group</v-icon>
            <span class="text-h6 font-weight-bold text-blue-darken-3">客戶健康檢查報告</span>
          </v-card-title>
          <v-card-text class="pa-0">
            <v-data-table
              :headers="[
                { title: '報告編號', key: 'id', align: 'start', width: '80px' },
                { title: '客戶ID', key: 'clientId', align: 'start', width: '100px' },
                { title: '客戶姓名', key: 'clientName', align: 'start', width: '120px' },
                { title: '健康數據', key: 'content', align: 'start' },
                { title: '日期', key: 'date', align: 'center', width: '110px' },
                { title: '操作', key: 'actions', align: 'center', width: '100px' }
              ]"
              :items="clientHealthData"
              :loading="loading"
              loading-text="資料載入中..."
              class="elevation-0"
              hover
              item-value="id"
              density="comfortable"
            >
              <template v-slot:item.content="{ item }">
                <div style="max-width: 350px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                  {{ item.content }}
                </div>
              </template>
              <template v-slot:item.date="{ item }">
                {{ item.date || '-' }}
              </template>
              <template v-slot:item.actions="{ item }">
                <v-btn 
                  color="blue" 
                  variant="flat"
                  size="small"
                  @click="handleSelectClient(item)"
                  :disabled="selectedClient && selectedClient.id === item.id"
                  :prepend-icon="selectedClient && selectedClient.id === item.id ? 'mdi-check' : 'mdi-chart-box'"
                  rounded
                >
                  {{ selectedClient && selectedClient.id === item.id ? '已選擇' : '風險分析' }}
                </v-btn>
              </template>
            </v-data-table>
          </v-card-text>
        </v-card>

        <!-- 風險評估區塊 -->
        <v-card v-if="selectedClient" class="mb-6 rounded-lg" elevation="2">
          <!-- 風險評估標題與按鈕 -->
          <v-card-title class="py-4 px-6 bg-blue-lighten-4 d-flex justify-space-between align-center">
            <div class="d-flex align-center">
              <v-icon size="24" color="blue-darken-3" class="me-2">mdi-clipboard-pulse</v-icon>
              <span class="text-h6 font-weight-bold text-blue-darken-3">
                保險風險評估 - 客戶：{{ selectedClient.clientName }}
              </span>
            </div>
            <v-btn
              color="deep-purple"
              size="large"
              :loading="llmLoading"
              :disabled="llmLoading"
              @click="handleLLMRiskAssessment"
              prepend-icon="mdi-robot"
              rounded="pill"
              elevation="3"
              class="font-weight-medium"
            >
              產生 AI 風險評估
            </v-btn>
          </v-card-title>
          
          <!-- AI 摘要區塊 -->
          <v-card-text v-if="llmSummary" class="pa-6">
            <v-alert
              color="blue-lighten-5"
              icon="mdi-robot"
              border="start"
              border-color="blue-darken-2"
              elevation="1"
              class="mb-6"
              density="comfortable"
            >
              <div class="text-blue-darken-4">{{ llmSummary }}</div>
            </v-alert>
            
            <!-- 雙欄排版 - 承保建議與保費試算 -->
            <v-row v-if="riskAssessment" class="mb-6">
              <v-col cols="12" md="6">
                <v-card
                  class="rounded-lg"
                  elevation="1"
                  border
                >
                  <v-card-title class="py-3 px-4 bg-green-lighten-5 d-flex align-center">
                    <v-icon size="24" color="green-darken-1" class="me-2">mdi-shield-check</v-icon>
                    <span class="text-h6 font-weight-medium text-green-darken-3">承保建議</span>
                  </v-card-title>
                  <v-card-text class="pa-4">
                    <v-list density="compact" class="bg-transparent pa-0">
                      <v-list-item>
                        <template v-slot:prepend>
                          <v-icon color="green-darken-1" class="me-2">mdi-alert-circle</v-icon>
                        </template>
                        <v-list-item-title class="text-body-1 font-weight-medium">
                          風險等級：<span class="text-green-darken-2 font-weight-bold">{{ riskAssessment.riskLevel }}</span>
                        </v-list-item-title>
                      </v-list-item>
                      <v-list-item>
                        <template v-slot:prepend>
                          <v-icon color="green-darken-1" class="me-2">mdi-check-circle</v-icon>
                        </template>
                        <v-list-item-title class="text-body-1 font-weight-medium">
                          建議決策：{{ riskAssessment.approvalRecommendation }}
                        </v-list-item-title>
                      </v-list-item>
                      <v-list-item>
                        <template v-slot:prepend>
                          <v-icon color="green-darken-1" class="me-2">mdi-currency-usd</v-icon>
                        </template>
                        <v-list-item-title class="text-body-1 font-weight-medium">
                          保費調整：{{ riskAssessment.premiumAdjustment }}
                        </v-list-item-title>
                      </v-list-item>
                      <v-list-item>
                        <template v-slot:prepend>
                          <v-icon color="green-darken-1" class="me-2">mdi-file-document</v-icon>
                        </template>
                        <v-list-item-title class="text-body-1 font-weight-medium">
                          特殊條件：{{ riskAssessment.conditions }}
                        </v-list-item-title>
                      </v-list-item>
                    </v-list>
                  </v-card-text>
                </v-card>
              </v-col>
              
              <v-col cols="12" md="6">
                <v-card
                  class="rounded-lg"
                  elevation="1"
                  border
                >
                  <v-card-title class="py-3 px-4 bg-blue-lighten-5 d-flex align-center">
                    <v-icon size="24" color="blue-darken-2" class="me-2">mdi-calculator</v-icon>
                    <span class="text-h6 font-weight-medium text-blue-darken-3">保費試算</span>
                  </v-card-title>
                  <v-card-text class="pa-4 py-4">
                    <v-row align="center" class="mb-2">
                      <v-col cols="6" class="text-body-1 font-weight-medium">基本保費：</v-col>
                      <v-col cols="6" class="text-end text-body-1">{{ calculatePremium().baseAmount }} 元/年</v-col>
                    </v-row>
                    <v-row align="center" class="mb-2">
                      <v-col cols="6" class="text-body-1 font-weight-medium">風險調整係數：</v-col>
                      <v-col cols="6" class="text-end text-body-1">× {{ calculatePremium().adjustmentFactor }}</v-col>
                    </v-row>
                    <v-divider class="my-3"></v-divider>
                    <v-row align="center">
                      <v-col cols="6" class="text-h6 font-weight-bold text-blue-darken-3">最終保費：</v-col>
                      <v-col cols="6" class="text-end text-h6 font-weight-bold text-blue-darken-3">{{ calculatePremium().finalPremium }} 元/年</v-col>
                    </v-row>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
            
            <!-- 健康指標分析區塊 -->
            <div v-if="showGauges">
              <div class="d-flex align-center mb-4">
                <v-icon size="24" color="blue-darken-1" class="me-2">mdi-chart-areaspline</v-icon>
                <span class="text-h6 font-weight-medium text-blue-darken-3">健康指標分析</span>
              </div>
              <v-row>
                <v-col v-for="gauge in healthGauges" :key="gauge.label" cols="6" sm="4" md="3" class="mb-4">
                  <v-card variant="flat" class="text-center py-2 rounded-lg">
                    <v-progress-circular
                      :model-value="((gauge.value - gauge.min) / (gauge.max - gauge.min)) * 100"
                      :size="90"
                      :width="12"
                      :color="gauge.color"
                      class="mb-2"
                    >
                      <v-icon v-if="gauge.icon" :color="gauge.color" size="24">{{ gauge.icon }}</v-icon>
                    </v-progress-circular>
                    <div class="mt-2 d-flex justify-center align-center">
                      <span class="text-h6 font-weight-bold me-1">{{ gauge.value }}</span>
                      <span v-if="gauge.unit" class="text-caption">{{ gauge.unit }}</span>
                    </div>
                    <div class="font-weight-medium">{{ gauge.label }}</div>
                    <div class="text-caption mt-1">{{ gauge.suggestion }}</div>
                  </v-card>
                </v-col>
              </v-row>
            </div>
          </v-card-text>
          
          <!-- 無評估資料時的提示 -->
          <v-card-text v-else class="pa-6 text-center">
            <v-icon size="48" color="blue-lighten-2" class="mb-3">mdi-information-outline</v-icon>
            <div class="text-h6 font-weight-medium text-blue-darken-1 mb-1">請點擊「產生 AI 風險評估」按鈕</div>
            <div class="text-body-2 text-grey">系統將分析客戶 {{ selectedClient.clientName }} 的健康數據並提供保險風險評估</div>
          </v-card-text>
        </v-card>
        
        <!-- 尚未選擇客戶的提示卡片 -->
        <v-card v-if="!selectedClient" class="mb-6 rounded-lg text-center py-8" elevation="1" color="grey-lighten-4">
          <v-icon size="64" color="blue-lighten-3" class="mb-4">mdi-account-search</v-icon>
          <div class="text-h6 font-weight-medium text-blue-grey-darken-1 mb-2">請從上方列表選擇客戶</div>
          <div class="text-body-2 text-grey px-6">
            選擇一位客戶後，系統將提供該客戶的健康風險評估和保險建議
          </div>
        </v-card>
      </v-col>
    </v-row>
    
    <!-- 右下角浮動按鈕組 -->
    <v-speed-dial
      v-if="selectedClient"
      location="bottom right"
      direction="top"
      transition="slide-y-reverse-transition"
      class="mb-8 me-8"
    >
      <template v-slot:activator>
        <v-btn
          color="deep-purple-darken-1"
          icon
          size="large"
          elevation="4"
        >
          <v-icon>mdi-plus</v-icon>
        </v-btn>
      </template>
      
      <v-btn
        color="green-darken-1"
        icon
        size="small"
        elevation="2"
      >
        <v-icon>mdi-email-fast-outline</v-icon>
        <v-tooltip activator="parent" location="left">傳送評估報告</v-tooltip>
      </v-btn>
      
      <v-btn
        color="blue-darken-1"
        icon
        size="small"
        elevation="2"
      >
        <v-icon>mdi-file-pdf-box</v-icon>
        <v-tooltip activator="parent" location="left">下載PDF報告</v-tooltip>
      </v-btn>
    </v-speed-dial>
  </v-container>
</template>

<style scoped>
.fill-height {
  min-height: 100vh;
  background-color: #f5f8fb;
}

:deep(.v-data-table) {
  background-color: transparent !important;
}

:deep(.v-data-table-header th) {
  font-weight: bold !important;
  color: rgba(0, 0, 0, 0.7) !important;
  background-color: #f0f7ff !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
}

:deep(.v-data-table-row:hover) {
  background-color: #f0f7ff !important;
}

.v-card {
  transition: transform 0.2s, box-shadow 0.2s;
}

.v-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1) !important;
}

.v-progress-circular {
  transition: transform 0.3s;
}

.v-card:hover .v-progress-circular {
  transform: scale(1.05);
}
</style>