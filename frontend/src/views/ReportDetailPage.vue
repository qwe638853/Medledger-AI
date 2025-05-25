<script setup>
// === 健檢指標中英文對照與參考值 ===
const METRIC_NAME_MAP = {
  'Glu-AC': '空腹血糖',
  'HbA1c': '糖化血色素',
  'Glu-PC': '飯後血糖',
  'Alb': '白蛋白',
  'TP': '血清蛋白總量',
  'AST（GOT）': '天門冬胺酸轉胺酶',
  'ALT（GPT）': '丙胺酸轉胺酶',
  'D-Bil': '直接膽紅素',
  'ALP': '鹼性磷酸酯酶',
  'T-Bil': '總膽紅素',
  'UN': '尿素氮',
  'CRE': '肌酸酐',
  'U.A': '尿酸',
  'T-CHO': '總膽固醇',
  'LDL-C': '低密度脂蛋白',
  'HDL-C': '高密度脂蛋白',
  'TG': '三酸甘油酯',
  'Hb': '血色素',
  'Hct': '血比容',
  'PLT': '血小板',
  'WBC': '白血球',
  'RBC': '紅血球',
  'hsCRP': '高敏感度C-反應蛋白',
  'AFP': '甲型胎兒蛋白',
  'CEA': '癌胚胎抗原',
  'CA-125': '癌症抗原125',
  'CA19-9': '癌症抗原19-9',
  // 血液常規
  'Hb': '血紅素',
  'RBC': '紅血球',
  'WBC': '白血球',
  'Hct': '血比容',
  'PLT': '血小板',
  'Platelet': '血小板',
  'MCV': '平均血球容積',
  'MCH': '平均血色素蛋白',
  'MCHC': '平均血球血紅素濃度',
  'Reticulocyte': '網狀紅血球',
  'Neutrophils (seg)': '嗜中性球',
  'Lymphocytes (Lym)': '淋巴球',
  'Monocytes (Mono)': '單核球',
  'Eosinophils (Eso)': '嗜酸性球',
  'Basophils (Baso)': '嗜鹼性球',
  'PT': '凝血酵素原時間',
  'aPTT': '活化部分凝血激素時間',
  'E.S.R.': '紅血球沉澱速率',
  'RDW-CV': '紅血球分佈寬度',
  // 尿液
  'Specific Gravity': '尿比重',
  'Specific Gravity (Dipstick)': '比重(尿液試紙)',
  'Color-Appearance': '外觀',
  'PH': '酸鹼度',
  'PH (Dipstick)': '酸鹼度(尿液試紙)',
  'Protein (Dipstick)': '尿蛋白(尿液試紙)',
  'Glucose (Dipstick)': '尿糖(尿液試紙)',
  'Bilirubin (Dipstick)': '膽紅素(尿液試紙)',
  'Urobilinogen (Dipstick)': '尿膽素原(尿液試紙)',
  'RBC (Urine)': '尿紅血球',
  'RBC (Sediment)': '尿沉渣紅血球',
  'WBC (Urine)': '尿白血球',
  'WBC (Sediment)': '尿沉渣白血球',
  'Epith Cell (Sediment)': '尿沉渣上皮細胞',
  'Casts (Sediment)': '尿沉渣圓柱體',
  'Ketone': '酮體',
  'Ketone(Dipstick)': '酮體(尿液試紙)',
  'Crystal (Sediment)': '尿沉渣結晶體',
  'Bacteria (Sediment)': '尿沉渣細菌',
  'Albumin (Dipstick)': '白蛋白(尿液試紙)',
  'Creatinine (Dipstick)': '肌酸酐(尿液試紙)',
  'Albumin / Creatinine Ratio (Dipstick)': '白蛋白對肌酸酐比值(尿液試紙)',
  'Nitrite': '亞硝酸',
  'Nitrite(Dipstick)': '亞硝酸(尿液試紙)',
  'Occult Blood': '潛血',
  'O.B.(Dipstick)': '潛血(尿液試紙)',
  'WBC Esterase': '白血球酯酶',
  'WBC esterase(Dipstick)': '白血球酯酶(尿液試紙)',
};
const METRIC_REF_RANGE = {
  'Glu-AC': { min: 70, max: 100, unit: 'mg/dL' },
  'HbA1c': { min: 4.0, max: 6.0, unit: '%' },
  'Glu-PC': { max: 140, unit: 'mg/dL' },
  'LDL-C': { max: 130, unit: 'mg/dL' },
  'HDL-C': { min: 40, unit: 'mg/dL' },
  'TG': { max: 150, unit: 'mg/dL' },
  'T-CHO': { max: 200, unit: 'mg/dL' },
  'Hb': { min: 11, max: 17.2, unit: 'g/dL' },
  // 血液
  'Hb': { min: 11.0, max: 17.2, unit: 'g/dL' }, // 男13.1~17.2 女11.0~15.2
  'RBC': { min: 3.78, max: 5.9, unit: '10^6/uL' }, // 男4.21~5.9 女3.78~5.25
  'WBC': { min: 3.25, max: 9.16, unit: '10^3/uL' },
  'Hct': { min: 34.8, max: 51.5, unit: '%' }, // 男39.6~51.5 女34.8~46.3
  'PLT': { min: 150, max: 378, unit: '10^3/uL' },
  'Platelet': { min: 150, max: 378, unit: '10^3/uL' },
  'MCV': { min: 80.9, max: 99.3, unit: 'fL' },
  'MCH': { min: 25.5, max: 33.2, unit: 'pg' },
  'MCHC': { min: 31.0, max: 34.9, unit: 'g/dL' },
  'Reticulocyte': { min: 0.87, max: 2.50, unit: '%' }, // 男1.05~2.50 女0.87~2.48
  'Neutrophils (seg)': { min: 41.6, max: 74.4, unit: '%' },
  'Lymphocytes (Lym)': { min: 18.0, max: 48.8, unit: '%' },
  'Monocytes (Mono)': { min: 3.3, max: 8.9, unit: '%' },
  'Eosinophils (Eso)': { min: 0.3, max: 7.9, unit: '%' },
  'Basophils (Baso)': { min: 0.2, max: 1.6, unit: '%' },
  'PT': { min: 9.7, max: 11.8, unit: 'sec' },
  'aPTT': { min: 25.6, max: 32.6, unit: 'sec' },
  'E.S.R.': { min: 2, max: 15, unit: 'mm/hr' }, // 男2~10 女2~15
  'RDW-CV': { min: 11.6, max: 15.0, unit: '%' },
  // 尿液
  'Specific Gravity (Dipstick)': { min: 1.003, max: 1.035 },
  'PH (Dipstick)': { min: 5.0, max: 8.0 },
  'Urobilinogen (Dipstick)': { max: 1.5, unit: 'mg/dL' },
  'RBC (Sediment)': { min: 0, max: 2, unit: '/HPF' },
  'WBC (Sediment)': { min: 0, max: 5, unit: '/HPF' },
};
function getMetricColor(key, value) {
  const ref = METRIC_REF_RANGE[key];
  if (!ref) return 'grey';
  const match = (value || '').toString().match(/-?\d+(\.\d+)?/);
  const num = match ? parseFloat(match[0]) : NaN;
  if (isNaN(num)) return 'grey';
  if (ref.min !== undefined && num < ref.min) return 'red';
  if (ref.max !== undefined && num > ref.max) return 'orange';
  return 'green';
}

import { ref, onMounted, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { healthCheckService } from '../services';
import { useAuthStore } from '../stores';
import { useUserStore } from '../stores/user';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const userStore = useUserStore();
const loading = ref(true);
const errorMsg = ref('');
const report = ref(null);

const reportId = route.params.report_id;
const patientId = route.params.patient_id;
const userRole = computed(() => route.query.role || authStore.userRole || 'patient');

// 彈窗控制
const showAISummary = ref(false);
const showRisk = ref(false);
const aiSummary = ref('這是 AI 分析摘要的假資料。');
const riskLevel = ref('低風險');
const riskAdvice = ref('您的主要指標均在正常範圍，請持續保持健康生活。');

// 風險評估簡單規則
function evaluateRisk(metrics = {}) {
  let high = 0, mid = 0;
  for (const [key, value] of Object.entries(metrics || {})) {
    const ref = METRIC_REF_RANGE[key];
    if (!ref) continue;
    const match = (value || '').toString().match(/-?\d+(\.\d+)?/);
    const num = match ? parseFloat(match[0]) : NaN;
    if (isNaN(num)) continue;
    if (ref.max !== undefined && num > ref.max) high++;
    else if (ref.min !== undefined && num < ref.min) mid++;
  }
  if (high > 0) {
    riskLevel.value = '高風險';
    riskAdvice.value = '部分指標超出正常範圍，建議儘速諮詢醫師。';
  } else if (mid > 0) {
    riskLevel.value = '中風險';
    riskAdvice.value = '部分指標偏低，建議定期追蹤。';
  } else {
    riskLevel.value = '低風險';
    riskAdvice.value = '您的主要指標均在正常範圍，請持續保持健康生活。';
  }
}

// 根據不同角色使用不同的 API endpoint
const fetchReportData = async () => {
  loading.value = true;
  errorMsg.value = '';
  
  try {
    let response;
    console.log("Role",userRole.value);
    // 只有保險公司角色需要調用 API
    if (userRole.value === 'insurer') {
      response = await healthCheckService.fetchReportContent(reportId, patientId);
      if (!response) {
        throw new Error('無法獲取報告數據');
      }
      
      report.value = {
        id: reportId,
        patient_id: patientId,
        date: new Date().toISOString(),
        rawData: JSON.parse(response.resultJson) || {}
      };
      console.log("response",report.value);
    } else {
      // 一般用戶直接使用 store 中的數據
      report.value = userStore.currentReport;
      console.log("response",report.value);
    }
    
    // 如果有數據，進行風險評估
    if (report.value?.rawData) {
      evaluateRisk(report.value.rawData);
    }
  } catch (error) {
    console.error('獲取報告詳情失敗:', error);
    errorMsg.value = error.message || '獲取報告詳情失敗';
  } finally {
    loading.value = false;
  }
};

// 計算屬性：數值型指標
const numericMetrics = computed(() => {
  if (!report.value?.rawData) return {};
  
  const metrics = {};
  Object.entries(report.value.rawData).forEach(([key, value]) => {
    if (typeof value === 'number' || !isNaN(parseFloat(value))) {
      metrics[key] = value;
    }
  });
  return metrics;
});

// 計算屬性：文字型指標
const textMetrics = computed(() => {
  if (!report.value?.rawData) return {};
  
  const metrics = {};
  Object.entries(report.value.rawData).forEach(([key, value]) => {
    if (typeof value === 'string' && isNaN(parseFloat(value))) {
      metrics[key] = value;
    }
  });
  return metrics;
});

function isNumericMetric(value) {
  // 只要能擷取出數字就算數值型
  if (typeof value === 'number') return true;
  if (typeof value === 'string') {
    const match = value.match(/-?\d+(\.\d+)?/);
    return !!match;
  }
  return false;
}

function getMetricPercent(key, value) {
  const ref = METRIC_REF_RANGE[key];
  if (!ref) return 0;
  const match = (value || '').toString().match(/-?\d+(\.\d+)?/);
  const num = match ? parseFloat(match[0]) : NaN;
  if (isNaN(num)) return 0;
  if (ref.max !== undefined && ref.min !== undefined) {
    return Math.min(Math.max(((num - ref.min) / (ref.max - ref.min)) * 100, 0), 100);
  } else if (ref.max !== undefined) {
    return Math.min((num / ref.max) * 100, 100);
  } else if (ref.min !== undefined) {
    return num > ref.min ? 100 : 0;
  }
  return 0;
}

function getMetricNumber(value) {
  // 只取第一個數字
  const match = (value || '').toString().match(/-?\d+(\.\d+)?/);
  return match ? match[0] : value;
}

const aiBtnHover = ref(false);
const riskBtnHover = ref(false);

onMounted(() => {
  fetchReportData();
});
</script>

<template>
  <div class="report-page">
    <v-container class="py-8">
      <!-- 返回按鈕 -->
      <v-btn
        @click="router.back()"
        class="back-btn mb-8"
        elevation="0"
      >
        <v-icon start size="20">mdi-arrow-left</v-icon>
        返回
      </v-btn>

      <v-card v-if="loading" class="loader-card">
        <v-progress-circular indeterminate color="#111827" />
      </v-card>
      
      <v-alert v-else-if="errorMsg" type="error" class="error-alert">
        {{ errorMsg }}
      </v-alert>

      <template v-else>
        <!-- 底部操作按鈕 -->
        <div class="action-buttons">
          <v-btn
            class="action-btn ai-btn"
            elevation="0"
            @click="showAISummary = !showAISummary"
          >
            <v-icon start size="20">mdi-robot-outline</v-icon>
            AI 分析摘要
          </v-btn>
          
          <v-btn
            class="action-btn risk-btn"
            elevation="0"
            @click="() => { if (!showRisk) evaluateRisk(numericMetrics); showRisk = !showRisk; }"
          >
            <v-icon start size="20">mdi-shield-outline</v-icon>
            風險評估
          </v-btn>
        </div>

        <!-- 分析結果區域 -->
        <div class="analysis-section" v-if="showAISummary || showRisk">
          <v-card v-if="showAISummary" class="analysis-card ai-card" elevation="0">
            <h3 class="analysis-title">AI 分析摘要</h3>
            <p class="analysis-content">{{ aiSummary }}</p>
          </v-card>

          <v-card v-if="showRisk" class="analysis-card risk-card" elevation="0">
            <h3 class="analysis-title">風險評估結果</h3>
            <div class="risk-level" :class="riskLevel === '高風險' ? 'high' : riskLevel === '中風險' ? 'medium' : 'low'">
              {{ riskLevel }}
            </div>
            <p class="analysis-content">{{ riskAdvice }}</p>
          </v-card>
        </div>
        <!-- 報告總覽卡片 -->
        <v-card class="overview-card mb-8" elevation="0">
          <div class="d-flex flex-column">
            <h1 class="report-title">健康檢查報告</h1>
            <p class="report-subtitle">{{ formatDate(report?.date) || '尚未設定日期' }}</p>
          </div>
          
          <v-divider class="my-6" />
          
          <v-row class="report-meta">
            <v-col cols="12" sm="4">
              <div class="meta-label">報告編號</div>
              <div class="meta-value">{{ reportId }}</div>
            </v-col>
            <v-col cols="12" sm="4">
              <div class="meta-label">檢查對象</div>
              <div class="meta-value">{{ patientId }}</div>
            </v-col>
            <v-col cols="12" sm="4">
              <div class="meta-label">檢查類型</div>
              <div class="meta-value">常規健康檢查</div>
            </v-col>
          </v-row>
        </v-card>

        <!-- 主要指標區域 -->
        <section class="metrics-section mb-12">
          <h2 class="section-title mb-6">主要健康指標</h2>
          <v-row>
            <v-col
              v-for="(value, key) in numericMetrics"
              :key="key"
              cols="12" sm="6" md="4" lg="3"
              class="metric-col"
            >
              <v-card class="metric-card" elevation="0">
                <div class="metric-content">
                  <!-- 圓環進度指示器 -->
                  <div class="metric-ring">
                    <svg class="ring" viewBox="0 0 100 100">
                      <!-- 背景圓環 -->
                      <circle
                        class="ring-bg"
                        cx="50"
                        cy="50"
                        r="40"
                        fill="none"
                        stroke="#f1f1f1"
                        stroke-width="8"
                      />
                      <!-- 進度圓環 -->
                      <circle
                        class="ring-progress"
                        cx="50"
                        cy="50"
                        r="40"
                        fill="none"
                        :stroke="getMetricColor(key, value)"
                        stroke-width="8"
                        :stroke-dasharray="`${getMetricPercent(key, value) * 2.51} 251`"
                        transform="rotate(-90 50 50)"
                      />
                      <!-- 中心數值 -->
                      <g transform="rotate(90 50 50)">
                        <text
                          x="50"
                          y="50"
                          text-anchor="middle"
                          dominant-baseline="central"
                          :fill="getMetricColor(key, value)"
                          class="ring-value"
                        >
                          {{ getMetricNumber(value) }}
                        </text>
                      </g>
                    </svg>
                  </div>
                  
                  <!-- 單位顯示 -->
                  <div class="metric-unit">
                    {{ METRIC_REF_RANGE[key]?.unit || '' }}
                  </div>
                  
                  <h3 class="metric-name">{{ METRIC_NAME_MAP[key] || key }}</h3>
                  <p class="metric-range">
                    參考值：
                    <template v-if="METRIC_REF_RANGE[key]">
                      {{ METRIC_REF_RANGE[key].min || '0' }} - 
                      {{ METRIC_REF_RANGE[key].max || '∞' }}
                      {{ METRIC_REF_RANGE[key].unit }}
                    </template>
                    <template v-else>--</template>
                  </p>
                </div>
              </v-card>
            </v-col>
          </v-row>
        </section>

        <!-- 其他指標區域 -->
        <section class="other-metrics-section mb-12">
          <h2 class="section-title mb-6">其他檢查項目</h2>
          <v-row>
            <v-col
              v-for="(value, key) in textMetrics"
              :key="key"
              cols="12" sm="6" md="4"
            >
              <v-card class="text-metric-card" elevation="0">
                <div class="text-metric-content">
                  <span class="text-metric-name">{{ METRIC_NAME_MAP[key] || key }}</span>
                  <span class="text-metric-value">{{ value }}</span>
                </div>
              </v-card>
            </v-col>
          </v-row>
        </section>

        
      </template>
    </v-container>
  </div>
</template>

<style scoped>
/* 全局樣式 */
.report-page {
  background-color: #F9F7F4;
  min-height: 100vh;
}

/* 返回按鈕 */
.back-btn {
  background-color: #F8F441 !important;
  color: #111827 !important;
  border-radius: 16px !important;
  font-weight: 600 !important;
  padding: 0 24px !important;
  height: 44px !important;
  transition: all 0.2s ease !important;
}

.back-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

/* 卡片基礎樣式 */
:deep(.v-card) {
  border-radius: 28px !important;
  background: white !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03) !important;
}

/* 報告總覽區域 */
.overview-card {
  padding: 2.5rem !important;
}

.report-title {
  font-size: 2rem;
  font-weight: 900;
  color: #111827;
  letter-spacing: -0.5px;
  margin: 0;
}

.report-subtitle {
  font-size: 1rem;
  color: #888;
  margin: 0.5rem 0 0;
}

.meta-label {
  font-size: 0.875rem;
  color: #888;
  margin-bottom: 0.25rem;
}

.meta-value {
  font-size: 1.125rem;
  color: #111827;
  font-weight: 500;
}

/* 分隔線 */
:deep(.v-divider) {
  border-color: rgba(0, 0, 0, 0.05) !important;
}

/* 區塊標題 */
.section-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
  letter-spacing: -0.5px;
}

/* 指標卡片 */
.metric-col {
  padding: 12px;
}

.metric-card {
  height: 100%;
  padding: 2rem !important;
  transition: all 0.2s ease;
  background: white;
  border-radius: 24px !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05) !important;
}

.metric-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.1) !important;
}

.metric-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
}

/* SVG 圓環樣式 */
.metric-ring {
  width: 140px;
  height: 140px;
  margin-bottom: 0.5rem;  /* 減少底部間距，為單位預留空間 */
  position: relative;
}

.ring {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.ring-bg {
  stroke: #f1f1f1;
}

.ring-progress {
  transition: stroke-dasharray 0.5s ease;
  stroke-linecap: round;
}

.ring-value {
  font-size: 28px;
  font-weight: 700;
  font-family: system-ui, -apple-system, sans-serif;
}

/* 單位樣式 */
.metric-unit {
  font-size: 14px;
  color: #888;
  font-weight: 500;
  margin-bottom: 1rem;  /* 與下方標題保持間距 */
  line-height: 1;
}

.metric-name {
  font-size: 1.125rem;
  color: #222;
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.metric-range {
  font-size: 0.875rem;
  color: #666;
  margin: 0;
  line-height: 1.5;
}

/* 文字型指標卡片 */
.text-metric-card {
  padding: 1.5rem !important;
}

.text-metric-content {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.text-metric-name {
  font-size: 0.875rem;
  color: #888;
}

.text-metric-value {
  font-size: 1rem;
  color: #111827;
}

/* 操作按鈕 */
.action-buttons {
  display: flex;
  gap: 1rem;
  justify-content: center;
  margin-bottom: 2rem;
}

.action-btn {
  border-radius: 16px !important;
  font-weight: 600 !important;
  padding: 0 32px !important;
  height: 48px !important;
  transition: all 0.2s ease !important;
}

.ai-btn {
  background-color: #F8F441 !important;
  color: #111827 !important;
}

.risk-btn {
  background-color: white !important;
  color: #111827 !important;
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
}

.action-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

/* 分析結果卡片 */
.analysis-section {
  max-width: 600px;
  margin: 0 auto;
}

.analysis-card {
  padding: 2rem !important;
  margin-bottom: 1rem;
}

.analysis-title {
  font-size: 1.25rem;
  font-weight: 700;
  color: #111827;
  margin-bottom: 1rem;
  text-align: center;
}

.analysis-content {
  color: #666;
  line-height: 1.6;
  margin: 0;
}

.risk-level {
  text-align: center;
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 1rem;
}

.risk-level.high { color: #ef4444; }
.risk-level.medium { color: #f59e0b; }
.risk-level.low { color: #10b981; }

/* RWD 適配 */
@media (max-width: 960px) {
  .overview-card {
    padding: 1.5rem !important;
  }
  
  .report-title {
    font-size: 1.75rem;
  }
  
  .section-title {
    font-size: 1.25rem;
  }
  
  .metric-ring {
    width: 120px;
    height: 120px;
    margin-bottom: 0.375rem;
  }
  
  .ring-value {
    font-size: 24px;
  }
  
  .metric-unit {
    font-size: 12px;
    margin-bottom: 0.75rem;
  }
}

@media (max-width: 600px) {
  .action-buttons {
    flex-direction: column;
  }
  
  .action-btn {
    width: 100%;
  }
  
  .metric-card {
    padding: 1.5rem !important;
  }
  
  .metric-ring {
    width: 100px;
    height: 100px;
    margin-bottom: 0.25rem;
  }
  
  .ring-value {
    font-size: 20px;
  }
  
  .metric-unit {
    font-size: 11px;
    margin-bottom: 0.5rem;
  }
}
</style> 