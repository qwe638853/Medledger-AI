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

// 保險分析相關資料
const insuranceAnalysis = ref({
  riskScore: 75,
  overallRiskLevel: 'medium', // low, medium, high
  riskCategories: {
    cardiovascular: {
      score: 68,
      level: 'medium',
      factors: ['總膽固醇偏高', 'HDL-C略低', '年齡因素'],
      impact: 'medium',
      description: '心血管風險較一般人略高，建議定期追蹤'
    },
    diabetes: {
      score: 45,
      level: 'low',
      factors: ['空腹血糖正常', '糖化血色素正常'],
      impact: 'low',
      description: '糖尿病風險低，維持現狀即可'
    },
    liver: {
      score: 72,
      level: 'medium',
      factors: ['ALT略高', 'AST正常', '工作壓力'],
      impact: 'medium',
      description: '肝功能需要關注，建議改善生活習慣'
    },
    cancer: {
      score: 35,
      level: 'low',
      factors: ['腫瘤標記正常', '無家族史'],
      impact: 'low',
      description: '癌症風險相對較低'
    }
  },
  healthMetrics: [
    { name: '血壓', value: 125, unit: 'mmHg', status: 'normal', weight: 0.25 },
    { name: '血糖', value: 95, unit: 'mg/dL', status: 'normal', weight: 0.20 },
    { name: '膽固醇', value: 215, unit: 'mg/dL', status: 'elevated', weight: 0.30 },
    { name: '肝功能', value: 45, unit: 'U/L', status: 'elevated', weight: 0.25 }
  ],
  ageRiskFactors: {
    currentAge: 35,
    riskIncrease: {
      '5years': 15,
      '10years': 35,
      '15years': 60
    }
  },
  riskMitigation: [
    {
      category: '生活習慣改善',
      actions: ['戒菸', '規律運動', '控制飲酒'],
      riskReduction: 25,
      timeframe: '6個月'
    },
    {
      category: '醫療追蹤',
      actions: ['定期健檢', '專科諮詢', '藥物治療'],
      riskReduction: 15,
      timeframe: '持續進行'
    },
    {
      category: '飲食調整',
      actions: ['低脂飲食', '減少糖分', '增加蔬果'],
      riskReduction: 20,
      timeframe: '3個月'
    }
  ],
  recommendations: [
    {
      type: 'immediate',
      title: '立即建議',
      items: ['承保標準體保費', '建議加強肝功能檢查', '心血管風險評估']
    },
    {
      type: 'monitoring',
      title: '持續監控',
      items: ['每6個月追蹤膽固醇', '年度完整健檢', '生活習慣改善追蹤']
    },
    {
      type: 'assessment',
      title: '風險評估',
      items: ['定期重新評估風險等級', '追蹤健康指標變化', '調整承保策略']
    }
  ]
});

// AI 分析和保單推薦
const aiAnalysis = ref({
  summary: '根據您的健康檢查結果，整體健康狀況良好。血壓和血糖指標在正常範圍內，但膽固醇指標略高，建議注意飲食控制。',
  healthScore: 85,
  riskLevel: 'low', // low, medium, high
  diseaseRisks: [
    {
      name: '心血管疾病',
      risk: 25,
      level: 'low',
      factors: ['膽固醇略高', '血壓正常'],
      prevention: '控制飲食，增加運動'
    },
    {
      name: '糖尿病',
      risk: 15,
      level: 'low',
      factors: ['血糖正常', '體重適中'],
      prevention: '維持現有生活習慣'
    },
    {
      name: '肝功能異常',
      risk: 35,
      level: 'medium',
      factors: ['ALT略高', '工作壓力大'],
      prevention: '減少熬夜，定期追蹤'
    }
  ],
  healthTrends: [
    { metric: '血糖', trend: 'stable', change: 0 },
    { metric: '血壓', trend: 'improving', change: -5 },
    { metric: '膽固醇', trend: 'concern', change: 8 }
  ],
  recommendations: [
    {
      type: 'diet',
      title: '飲食建議',
      items: ['減少高膽固醇食物攝入', '增加高纖維食物', '控制總熱量攝取']
    },
    {
      type: 'exercise',
      title: '運動建議',
      items: ['每週至少150分鐘中等強度運動', '加強心肺功能訓練', '增加阻力訓練']
    },
    {
      type: 'lifestyle',
      title: '生活習慣',
      items: ['保持充足睡眠', '減少工作壓力', '定期健康檢查']
    }
  ],
  insuranceRecommendations: [
    {
      name: '優質醫療保障計劃',
      coverage: '住院醫療、手術、門診',
      features: ['免等待期', '一年一約', '可選擇醫院'],
      monthlyPremium: 2500,
      suitability: 90
    },
    {
      name: '重大疾病防護計劃',
      coverage: '癌症、心臟病、中風等重大疾病',
      features: ['保額最高500萬', '保費豁免權益', '理賠無等待期'],
      monthlyPremium: 3200,
      suitability: 75
    }
  ]
});

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
const activeTab = ref('diet');
const insuranceActiveTab = ref('immediate');

// 競爭對手分析表格標題
const competitorHeaders = [
  { title: '公司', key: 'company', align: 'start' },
  { title: '保費', key: 'premium', align: 'center' },
  { title: '保額', key: 'coverage', align: 'center' },
  { title: '風險調整', key: 'riskAdjustment', align: 'center' },
  { title: '競爭力評分', key: 'score', align: 'center' }
];

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
        color="#00B8D9"
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
          <!-- 只有病患看得到 AI 分析按鈕 -->
          <v-btn
            v-if="userRole === 'user'"
            class="action-btn ai-btn"
            elevation="0"
            color="#00B8D9"
            @click="showAISummary = true"
          >
            <v-icon start size="20">mdi-robot-outline</v-icon>
            AI 智能分析
          </v-btn>
          
          <v-btn
            v-if="userRole === 'insurer'"
            class="action-btn risk-btn"
            elevation="0"
            color="#00B8D9"
            @click="() => { if (!showRisk) evaluateRisk(numericMetrics); showRisk = true; }"
          >
            <v-icon start size="20">mdi-shield-outline</v-icon>
            專業風險評估
          </v-btn>
        </div>

        <!-- 保險業者風險分析區域 -->
        <div class="analysis-section" v-if="showRisk && userRole === 'insurer'">
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

    <!-- AI 智能分析彈窗 -->
    <v-dialog v-model="showAISummary" max-width="900" scrollable>
      <v-card class="ai-dialog-card">
        <!-- 對話框標題 -->
        <v-card-title class="ai-dialog-header">
          <div class="d-flex align-center">
            <v-avatar class="ai-avatar mr-3" color="gradient">
              <v-icon color="white" size="24">mdi-robot-outline</v-icon>
            </v-avatar>
            <div>
              <div class="ai-dialog-title">AI 智能健康分析</div>
              <div class="ai-dialog-subtitle">基於您的健康數據進行專業分析</div>
            </div>
          </div>
          <v-spacer></v-spacer>
          <v-btn
            icon
            variant="text"
            @click="showAISummary = false"
            class="close-btn"
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-card-text class="ai-dialog-content">
          <!-- 健康評分區域 -->
          <div class="health-score-section">
            <div class="health-score-container">
              <v-progress-circular
                :model-value="aiAnalysis.healthScore"
                :color="aiAnalysis.healthScore > 80 ? '#4CAF50' : aiAnalysis.healthScore > 60 ? '#FF9800' : '#F44336'"
                size="120"
                width="12"
                class="health-score-circle"
              >
                <div class="health-score-content">
                  <div class="health-score-number">{{ aiAnalysis.healthScore }}</div>
                  <div class="health-score-label">健康評分</div>
                </div>
              </v-progress-circular>
              <div class="health-level">
                <v-chip
                  :color="aiAnalysis.riskLevel === 'low' ? 'success' : aiAnalysis.riskLevel === 'medium' ? 'warning' : 'error'"
                  size="large"
                  class="health-level-chip"
                >
                  {{ aiAnalysis.riskLevel === 'low' ? '健康狀況良好' : aiAnalysis.riskLevel === 'medium' ? '需要注意' : '需要關注' }}
                </v-chip>
              </div>
            </div>
          </div>

          <!-- AI 分析摘要 -->
          <div class="ai-summary-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-brain</v-icon>
              AI 分析摘要
            </h3>
            <div class="ai-summary-content">
              {{ aiAnalysis.summary }}
            </div>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 疾病風險分析 -->
          <div class="disease-risk-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-shield-alert</v-icon>
              疾病風險評估
            </h3>
            <div class="disease-risk-grid">
              <v-card
                v-for="(disease, index) in aiAnalysis.diseaseRisks"
                :key="index"
                class="disease-risk-card"
                elevation="2"
              >
                <div class="disease-risk-header">
                  <div class="disease-name">{{ disease.name }}</div>
                  <v-chip
                    :color="disease.level === 'low' ? 'success' : disease.level === 'medium' ? 'warning' : 'error'"
                    size="small"
                  >
                    {{ disease.risk }}% 風險
                  </v-chip>
                </div>
                
                <v-progress-linear
                  :model-value="disease.risk"
                  :color="disease.level === 'low' ? 'success' : disease.level === 'medium' ? 'warning' : 'error'"
                  height="8"
                  rounded
                  class="risk-progress"
                ></v-progress-linear>

                <div class="disease-details">
                  <div class="risk-factors">
                    <div class="detail-label">影響因子：</div>
                    <div class="risk-factor-tags">
                      <v-chip
                        v-for="factor in disease.factors"
                        :key="factor"
                        size="x-small"
                        color="grey-lighten-3"
                        class="mr-1 mb-1"
                      >
                        {{ factor }}
                      </v-chip>
                    </div>
                  </div>
                  <div class="prevention">
                    <div class="detail-label">預防建議：</div>
                    <div class="prevention-text">{{ disease.prevention }}</div>
                  </div>
                </div>
              </v-card>
            </div>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 健康趨勢 -->
          <div class="health-trends-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-trending-up</v-icon>
              健康趨勢分析
            </h3>
            <div class="trends-grid">
              <div
                v-for="trend in aiAnalysis.healthTrends"
                :key="trend.metric"
                class="trend-item"
              >
                <div class="trend-metric">{{ trend.metric }}</div>
                <div class="trend-indicator">
                  <v-icon
                    :color="trend.trend === 'improving' ? 'success' : trend.trend === 'concern' ? 'error' : 'grey'"
                    size="20"
                  >
                    {{ trend.trend === 'improving' ? 'mdi-trending-up' : trend.trend === 'concern' ? 'mdi-trending-down' : 'mdi-trending-neutral' }}
                  </v-icon>
                  <span class="trend-text">
                    {{ trend.trend === 'improving' ? '改善中' : trend.trend === 'concern' ? '需關注' : '穩定' }}
                  </span>
                </div>
              </div>
            </div>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 健康建議 -->
          <div class="recommendations-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-lightbulb</v-icon>
              個人化健康建議
            </h3>
            <div class="recommendations-tabs">
              <v-tabs v-model="activeTab" color="#00B8D9" align-tabs="center">
                <v-tab
                  v-for="rec in aiAnalysis.recommendations"
                  :key="rec.type"
                  :value="rec.type"
                >
                  <v-icon start>
                    {{ rec.type === 'diet' ? 'mdi-food-apple' : rec.type === 'exercise' ? 'mdi-run' : 'mdi-heart' }}
                  </v-icon>
                  {{ rec.title }}
                </v-tab>
              </v-tabs>
              
              <v-window v-model="activeTab" class="recommendations-content">
                <v-window-item
                  v-for="rec in aiAnalysis.recommendations"
                  :key="rec.type"
                  :value="rec.type"
                >
                  <div class="recommendation-items">
                    <div
                      v-for="(item, index) in rec.items"
                      :key="index"
                      class="recommendation-item"
                    >
                      <v-icon color="#00B8D9" size="16" class="mr-2">mdi-check-circle</v-icon>
                      {{ item }}
                    </div>
                  </div>
                </v-window-item>
              </v-window>
            </div>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 保險推薦 -->
          <div class="insurance-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-shield-account</v-icon>
              智能保險推薦
            </h3>
            <div class="insurance-grid">
              <v-card
                v-for="(plan, index) in aiAnalysis.insuranceRecommendations"
                :key="index"
                class="insurance-card"
                elevation="3"
              >
                <div class="insurance-header">
                  <div class="insurance-name">{{ plan.name }}</div>
                  <div class="insurance-premium">
                    <span class="premium-amount">NT$ {{ plan.monthlyPremium }}</span>
                    <span class="premium-period">/月</span>
                  </div>
                </div>
                
                <div class="insurance-coverage">{{ plan.coverage }}</div>
                
                <div class="suitability-section">
                  <div class="suitability-label">適合度</div>
                  <v-progress-linear
                    :model-value="plan.suitability"
                    color="#00B8D9"
                    height="6"
                    rounded
                  ></v-progress-linear>
                  <div class="suitability-percentage">{{ plan.suitability }}%</div>
                </div>

                <div class="insurance-features">
                  <v-chip
                    v-for="feature in plan.features"
                    :key="feature"
                    size="small"
                    color="blue-grey-lighten-4"
                    class="mr-1 mb-1"
                  >
                    {{ feature }}
                  </v-chip>
                </div>
              </v-card>
            </div>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- 保險風險分析彈窗 -->
    <v-dialog v-model="showRisk" max-width="1000" scrollable>
      <v-card class="insurance-dialog-card">
        <!-- 對話框標題 -->
        <v-card-title class="insurance-dialog-header">
          <div class="d-flex align-center">
            <v-avatar class="insurance-avatar mr-3" color="gradient">
              <v-icon color="white" size="24">mdi-shield-check</v-icon>
            </v-avatar>
            <div>
              <div class="insurance-dialog-title">專業保險風險評估</div>
              <div class="insurance-dialog-subtitle">基於健康數據的精準風險分析與保費定價</div>
            </div>
          </div>
          <v-spacer></v-spacer>
          <v-btn
            icon
            variant="text"
            @click="showRisk = false"
            class="close-btn"
          >
            <v-icon>mdi-close</v-icon>
          </v-btn>
        </v-card-title>

        <v-card-text class="insurance-dialog-content">
          <!-- 風險評分總覽 -->
          <div class="risk-overview-section">
            <div class="risk-score-center">
              <div class="risk-score-card">
                <v-progress-circular
                  :model-value="insuranceAnalysis.riskScore"
                  :color="insuranceAnalysis.riskScore > 80 ? '#F44336' : insuranceAnalysis.riskScore > 60 ? '#FF9800' : '#4CAF50'"
                  size="120"
                  width="12"
                >
                  <div class="risk-score-content">
                    <div class="risk-score-number">{{ insuranceAnalysis.riskScore }}</div>
                    <div class="risk-score-label">風險評分</div>
                  </div>
                </v-progress-circular>
                <div class="risk-level-indicator">
                  <v-chip
                    :color="insuranceAnalysis.overallRiskLevel === 'low' ? 'success' : insuranceAnalysis.overallRiskLevel === 'medium' ? 'warning' : 'error'"
                    size="large"
                  >
                    {{ insuranceAnalysis.overallRiskLevel === 'low' ? '低風險' : insuranceAnalysis.overallRiskLevel === 'medium' ? '中等風險' : '高風險' }}
                  </v-chip>
                </div>
              </div>
            </div>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 風險分類詳細分析 -->
          <div class="risk-categories-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-chart-donut</v-icon>
              疾病風險分類評估
            </h3>
            <v-row>
              <v-col
                v-for="(category, key) in insuranceAnalysis.riskCategories"
                :key="key"
                cols="12" md="6"
              >
                <v-card class="risk-category-card" elevation="2">
                  <div class="category-header">
                    <div class="category-name">
                      {{ key === 'cardiovascular' ? '心血管疾病' : 
                         key === 'diabetes' ? '糖尿病' : 
                         key === 'liver' ? '肝臟疾病' : 
                         key === 'cancer' ? '癌症' : key }}
                    </div>
                    <v-chip
                      :color="category.level === 'low' ? 'success' : category.level === 'medium' ? 'warning' : 'error'"
                      size="small"
                    >
                      {{ category.score }}分
                    </v-chip>
                  </div>
                  
                  <v-progress-linear
                    :model-value="category.score"
                    :color="category.level === 'low' ? 'success' : category.level === 'medium' ? 'warning' : 'error'"
                    height="12"
                    rounded
                    class="category-progress"
                  ></v-progress-linear>

                  <div class="category-details">
                    <div class="category-description">{{ category.description }}</div>
                    <div class="category-factors">
                      <div class="factors-label">關鍵因素：</div>
                      <div class="factors-chips">
                        <v-chip
                          v-for="factor in category.factors"
                          :key="factor"
                          size="x-small"
                          color="grey-lighten-3"
                          class="mr-1 mb-1"
                        >
                          {{ factor }}
                        </v-chip>
                      </div>
                    </div>
                  </div>
                </v-card>
              </v-col>
            </v-row>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 健康指標影響分析 -->
          <div class="health-metrics-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-heart-pulse</v-icon>
              健康指標風險權重分析
            </h3>
            <v-row>
              <v-col
                v-for="metric in insuranceAnalysis.healthMetrics"
                :key="metric.name"
                cols="12" sm="6" md="3"
              >
                <div class="metric-weight-card">
                  <div class="metric-header">
                    <div class="metric-name">{{ metric.name }}</div>
                    <v-icon
                      :color="metric.status === 'normal' ? 'success' : 'warning'"
                      size="20"
                    >
                      {{ metric.status === 'normal' ? 'mdi-check-circle' : 'mdi-alert-circle' }}
                    </v-icon>
                  </div>
                  <div class="metric-value">{{ metric.value }} {{ metric.unit }}</div>
                  <div class="metric-weight">
                    <span class="weight-label">風險權重:</span>
                    <span class="weight-value">{{ (metric.weight * 100).toFixed(0) }}%</span>
                  </div>
                  <v-progress-linear
                    :model-value="metric.weight * 100"
                    color="#00B8D9"
                    height="6"
                    rounded
                  ></v-progress-linear>
                </div>
              </v-col>
            </v-row>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 年齡風險趨勢 -->
          <div class="age-risk-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-trending-up</v-icon>
              年齡風險趨勢預測
            </h3>
            <div class="age-risk-chart">
              <div class="current-age">
                <div class="age-label">目前年齡</div>
                <div class="age-value">{{ insuranceAnalysis.ageRiskFactors.currentAge }}歲</div>
              </div>
              <div class="age-projections">
                <div
                  v-for="(increase, period) in insuranceAnalysis.ageRiskFactors.riskIncrease"
                  :key="period"
                  class="age-projection-item"
                >
                  <div class="projection-period">{{ period.replace('years', '年後') }}</div>
                  <div class="projection-increase">風險增加 {{ increase }}%</div>
                  <v-progress-linear
                    :model-value="increase"
                    color="orange"
                    height="8"
                    rounded
                  ></v-progress-linear>
                </div>
              </div>
            </div>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 風險緩解策略 -->
          <div class="risk-mitigation-section">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-shield-plus</v-icon>
              風險緩解策略建議
            </h3>
            <v-row>
              <v-col
                v-for="strategy in insuranceAnalysis.riskMitigation"
                :key="strategy.category"
                cols="12" md="4"
              >
                <v-card class="mitigation-card" elevation="2">
                  <div class="mitigation-header">
                    <div class="mitigation-category">{{ strategy.category }}</div>
                    <v-chip color="success" size="small">
                      -{{ strategy.riskReduction }}% 風險
                    </v-chip>
                  </div>
                  <div class="mitigation-actions">
                    <div
                      v-for="action in strategy.actions"
                      :key="action"
                      class="mitigation-action"
                    >
                      <v-icon color="success" size="16" class="mr-2">mdi-check</v-icon>
                      {{ action }}
                    </div>
                  </div>
                  <div class="mitigation-timeframe">
                    <v-icon size="16" class="mr-1">mdi-clock-outline</v-icon>
                    執行期程: {{ strategy.timeframe }}
                  </div>
                </v-card>
              </v-col>
            </v-row>
          </div>

          <v-divider class="my-6"></v-divider>

          <!-- 專業建議 -->
          <div class="professional-recommendations">
            <h3 class="section-header">
              <v-icon class="section-icon">mdi-lightbulb-on</v-icon>
              專業核保建議
            </h3>
            <v-tabs v-model="insuranceActiveTab" color="#00B8D9" align-tabs="center">
              <v-tab
                v-for="rec in insuranceAnalysis.recommendations"
                :key="rec.type"
                :value="rec.type"
              >
                <v-icon start>
                  {{ rec.type === 'immediate' ? 'mdi-lightning-bolt' : rec.type === 'monitoring' ? 'mdi-monitor-eye' : 'mdi-clipboard-check' }}
                </v-icon>
                {{ rec.title }}
              </v-tab>
            </v-tabs>
            
            <v-window v-model="insuranceActiveTab" class="recommendations-content">
              <v-window-item
                v-for="rec in insuranceAnalysis.recommendations"
                :key="rec.type"
                :value="rec.type"
              >
                <div class="recommendation-items">
                  <div
                    v-for="(item, index) in rec.items"
                    :key="index"
                    class="recommendation-item"
                  >
                    <v-icon color="#00B8D9" size="16" class="mr-2">mdi-arrow-right-circle</v-icon>
                    {{ item }}
                  </div>
                </div>
              </v-window-item>
            </v-window>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>
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
  background-color: #00B8D9 !important;
  color: white !important;
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
  font-size: 2.5rem;
  font-weight: 900;
  color: #111827;
  letter-spacing: -0.5px;
  margin: 0;
}

.report-subtitle {
  font-size: 1.2rem;
  color: #888;
  margin: 0.5rem 0 0;
}

.meta-label {
  font-size: 1rem;
  color: #888;
  margin-bottom: 0.25rem;
}

.meta-value {
  font-size: 1.3rem;
  color: #111827;
  font-weight: 500;
}

/* 分隔線 */
:deep(.v-divider) {
  border-color: rgba(0, 0, 0, 0.05) !important;
}

/* 區塊標題 */
.section-title {
  font-size: 1.8rem;
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
  font-size: 16px;
  color: #888;
  font-weight: 500;
  margin-bottom: 1rem;  /* 與下方標題保持間距 */
  line-height: 1;
}

.metric-name {
  font-size: 1.3rem;
  color: #222;
  margin-bottom: 0.75rem;
  font-weight: 600;
}

.metric-range {
  font-size: 1rem;
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
  font-size: 1rem;
  color: #888;
}

.text-metric-value {
  font-size: 1.2rem;
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
  background-color: #00B8D9 !important;
  color: white !important;
}

.risk-btn {
  background-color: #00B8D9 !important;
  color: white !important;
  border: 1px solid #00B8D9 !important;
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
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
  margin-bottom: 1rem;
  text-align: center;
}

.analysis-content {
  color: #666;
  line-height: 1.6;
  margin: 0;
  font-size: 1.1rem;
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
    font-size: 2rem;
  }
  
  .section-title {
    font-size: 1.5rem;
  }
  
  .metric-ring {
    width: 120px;
    height: 120px;
    margin-bottom: 0.375rem;
  }
  
  .ring-value {
    font-size: 28px;
  }
  
  .metric-unit {
    font-size: 14px;
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
    font-size: 24px;
  }
  
  .metric-unit {
    font-size: 13px;
    margin-bottom: 0.5rem;
  }
}

/* 保險分析相關樣式 */
.risk-score-section {
  padding: 1rem 0;
}

.risk-subtitle {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
}

.risk-list {
  background: transparent !important;
}

.risk-item {
  padding: 12px !important;
  margin-bottom: 8px;
  background: #f8f9fa;
  border-radius: 8px;
}

.risk-type {
  font-weight: 500 !important;
  font-size: 1.1rem !important;
}

.risk-description {
  font-size: 1rem !important;
  color: #666 !important;
  margin-top: 4px !important;
}

.suggestions-list {
  list-style-type: none;
  padding: 0;
  margin: 0;
}

.suggestions-list li {
  padding: 8px 0;
  color: #444;
  position: relative;
  padding-left: 24px;
  font-size: 1.1rem;
}

.suggestions-list li::before {
  content: "•";
  color: #666;
  position: absolute;
  left: 8px;
}

/* 保險推薦相關樣式 */
.recommendations-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 1rem;
}

.recommendations-list {
  list-style-type: none;
  padding: 0;
  margin: 0 0 1rem;
}

.recommendations-list li {
  padding: 8px 0;
  color: #444;
  position: relative;
  padding-left: 24px;
  font-size: 1.1rem;
}

.recommendations-list li::before {
  content: "•";
  color: #666;
  position: absolute;
  left: 8px;
}

.insurance-recommendations {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.insurance-plan-card {
  padding: 1rem;
  border-radius: 12px !important;
  background: #f8f9fa !important;
}

.plan-name {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
  margin: 0;
}

.plan-coverage {
  font-size: 1rem;
  color: #666;
  margin: 0.5rem 0;
}

.plan-features {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

@media (max-width: 600px) {
  .insurance-plan-card {
    padding: 12px;
  }

  .plan-features {
    margin-top: 8px;
  }
}

/* AI 彈窗樣式 */
.ai-dialog-card {
  border-radius: 24px !important;
  overflow: hidden !important;
}

.ai-dialog-header {
  background: linear-gradient(135deg, #00B8D9 0%, #0093A6 100%) !important;
  color: white !important;
  padding: 1.5rem 2rem !important;
}

.ai-avatar {
  background: rgba(255, 255, 255, 0.2) !important;
}

.ai-dialog-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 0.25rem;
}

.ai-dialog-subtitle {
  font-size: 0.9rem;
  opacity: 0.9;
}

.close-btn {
  color: white !important;
}

.ai-dialog-content {
  padding: 2rem !important;
  max-height: 70vh;
}

/* 健康評分區域 */
.health-score-section {
  text-align: center;
  margin-bottom: 2rem;
}

.health-score-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.health-score-circle {
  margin-bottom: 1rem;
}

.health-score-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.health-score-number {
  font-size: 2rem;
  font-weight: 700;
  line-height: 1;
}

.health-score-label {
  font-size: 0.8rem;
  opacity: 0.8;
  margin-top: 0.25rem;
}

.health-level-chip {
  font-weight: 600 !important;
}

/* 區塊標題 */
.section-header {
  display: flex;
  align-items: center;
  font-size: 1.5rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 1rem;
}

.section-icon {
  margin-right: 0.5rem;
  color: #00B8D9;
}

/* AI 分析摘要 */
.ai-summary-section {
  margin-bottom: 2rem;
}

.ai-summary-content {
  font-size: 1.2rem;
  line-height: 1.7;
  color: #555;
  background: #f8f9fa;
  padding: 1.8rem;
  border-radius: 12px;
  border-left: 4px solid #00B8D9;
}

/* 疾病風險分析 */
.disease-risk-section {
  margin-bottom: 2rem;
}

.disease-risk-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
}

.disease-risk-card {
  padding: 1.5rem;
  border-radius: 16px !important;
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

.disease-risk-card:hover {
  transform: translateY(-2px);
  border-color: #00B8D9;
}

.disease-risk-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.disease-name {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
}

.risk-progress {
  margin: 1rem 0;
}

.disease-details {
  margin-top: 1rem;
}

.detail-label {
  font-size: 1.1rem;
  font-weight: 600;
  color: #666;
  margin-bottom: 0.5rem;
}

.risk-factor-tags {
  margin-bottom: 1rem;
}

.prevention-text {
  font-size: 1.1rem;
  color: #555;
  line-height: 1.5;
}

/* 健康趨勢 */
.health-trends-section {
  margin-bottom: 2rem;
}

.trends-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.trend-item {
  background: #f8f9fa;
  padding: 1.2rem;
  border-radius: 12px;
  text-align: center;
}

.trend-metric {
  font-size: 1.2rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 0.5rem;
}

.trend-indicator {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
}

.trend-text {
  font-size: 1.1rem;
  font-weight: 500;
}

/* 健康建議 */
.recommendations-section {
  margin-bottom: 2rem;
}

.recommendations-tabs {
  margin-top: 1rem;
}

.recommendations-content {
  padding: 1.5rem 0;
}

.recommendation-items {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.recommendation-item {
  display: flex;
  align-items: flex-start;
  font-size: 1.2rem;
  line-height: 1.6;
  color: #555;
}

/* 保險推薦 */
.insurance-section {
  margin-bottom: 1rem;
}

.insurance-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 1.5rem;
}

.insurance-card {
  padding: 1.5rem;
  border-radius: 16px !important;
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

.insurance-card:hover {
  transform: translateY(-2px);
  border-color: #00B8D9;
}

.insurance-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1rem;
}

.insurance-name {
  font-size: 1.4rem;
  font-weight: 600;
  color: #333;
  flex: 1;
}

.insurance-premium {
  text-align: right;
}

.premium-amount {
  font-size: 1.4rem;
  font-weight: 700;
  color: #00B8D9;
}

.premium-period {
  font-size: 1rem;
  color: #666;
}

.insurance-coverage {
  font-size: 1.2rem;
  color: #555;
  margin-bottom: 1rem;
  line-height: 1.5;
}

.suitability-section {
  margin: 1rem 0;
}

.suitability-label {
  font-size: 1.1rem;
  font-weight: 600;
  color: #666;
  margin-bottom: 0.5rem;
}

.suitability-percentage {
  font-size: 1rem;
  color: #00B8D9;
  font-weight: 600;
  text-align: right;
  margin-top: 0.25rem;
}

.insurance-features {
  margin-top: 1rem;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .ai-dialog-content {
    padding: 1rem !important;
  }
  
  .disease-risk-grid,
  .insurance-grid {
    grid-template-columns: 1fr;
  }
  
  .trends-grid {
    grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  }
  
  .health-score-number {
    font-size: 1.5rem;
  }
  
  .insurance-header {
    flex-direction: column;
    gap: 0.5rem;
    align-items: flex-start;
  }
  
  .insurance-premium {
    text-align: left;
  }
  
  /* 調整小螢幕字體大小 */
  .ai-summary-content {
    font-size: 1.1rem;
    padding: 1.5rem;
  }
  
  .disease-name {
    font-size: 1.2rem;
  }
  
  .detail-label {
    font-size: 1rem;
  }
  
  .prevention-text {
    font-size: 1rem;
  }
  
  .trend-metric {
    font-size: 1.1rem;
  }
  
  .recommendation-item {
    font-size: 1.1rem;
  }
  
  .insurance-name {
    font-size: 1.3rem;
  }
  
  .insurance-coverage {
    font-size: 1.1rem;
  }
  
  .section-header {
    font-size: 1.3rem;
  }
}

/* 保險風險分析彈窗樣式 */
.insurance-dialog-card {
  border-radius: 24px !important;
  overflow: hidden !important;
}

.insurance-dialog-header {
  background: linear-gradient(135deg, #1976D2 0%, #1565C0 100%) !important;
  color: white !important;
  padding: 1.5rem 2rem !important;
}

.insurance-avatar {
  background: rgba(255, 255, 255, 0.2) !important;
}

.insurance-dialog-title {
  font-size: 1.5rem;
  font-weight: 700;
  margin-bottom: 0.25rem;
}

.insurance-dialog-subtitle {
  font-size: 0.9rem;
  opacity: 0.9;
}

.insurance-dialog-content {
  padding: 2rem !important;
  max-height: 75vh;
}

/* 風險評分總覽 */
.risk-overview-section {
  margin-bottom: 2rem;
}

.risk-score-center {
  display: flex;
  justify-content: center;
  align-items: center;
}

.risk-score-card {
  text-align: center;
  padding: 2rem;
  background: #f8f9fa;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.risk-score-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.risk-score-number {
  font-size: 1.8rem;
  font-weight: 700;
  line-height: 1;
}

.risk-score-label {
  font-size: 0.8rem;
  opacity: 0.8;
  margin-top: 0.25rem;
}

.risk-level-indicator {
  margin-top: 1rem;
}

.premium-recommendation,
.claims-prediction {
  text-align: center;
  padding: 1rem;
  background: #f8f9fa;
  border-radius: 12px;
  height: 100%;
}

.premium-title,
.claims-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 1rem;
}

.premium-amount {
  font-size: 2rem;
  font-weight: 700;
  color: #1976D2;
  margin-bottom: 0.5rem;
}

.premium-comparison {
  font-size: 0.9rem;
  color: #666;
}

.standard-premium {
  color: #888;
}

.claims-stats {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.claims-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.5rem;
  background: white;
  border-radius: 8px;
}

.claims-period {
  font-weight: 600;
  color: #333;
}

.claims-probability {
  font-weight: 700;
  color: #FF5722;
}

/* 風險分類卡片 */
.risk-categories-section {
  margin-bottom: 2rem;
}

.risk-category-card {
  padding: 1.5rem;
  border-radius: 16px !important;
  height: 100%;
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

.risk-category-card:hover {
  transform: translateY(-2px);
  border-color: #1976D2;
}

.category-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.category-name {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
}

.category-progress {
  margin: 1rem 0;
}

.category-details {
  margin-top: 1rem;
}

.category-description {
  font-size: 1.1rem;
  color: #555;
  margin-bottom: 1rem;
  line-height: 1.5;
}

.category-factors {
  margin-top: 1rem;
}

.factors-label {
  font-size: 1rem;
  font-weight: 600;
  color: #666;
  margin-bottom: 0.5rem;
}

.factors-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.25rem;
}

/* 健康指標權重 */
.health-metrics-section {
  margin-bottom: 2rem;
}

.metric-weight-card {
  background: #f8f9fa;
  padding: 1.2rem;
  border-radius: 12px;
  text-align: center;
  height: 100%;
}

.metric-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.metric-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: #333;
}

.metric-value {
  font-size: 1.4rem;
  font-weight: 700;
  color: #1976D2;
  margin: 0.5rem 0;
}

.metric-weight {
  display: flex;
  justify-content: space-between;
  font-size: 1rem;
  margin-bottom: 0.5rem;
}

.weight-label {
  color: #666;
}

.weight-value {
  font-weight: 600;
  color: #333;
}

/* 年齡風險趨勢 */
.age-risk-section {
  margin-bottom: 2rem;
}

.age-risk-chart {
  display: flex;
  gap: 2rem;
  align-items: center;
  background: #f8f9fa;
  padding: 1.5rem;
  border-radius: 12px;
}

.current-age {
  text-align: center;
  min-width: 120px;
}

.age-label {
  font-size: 1.1rem;
  color: #666;
  margin-bottom: 0.5rem;
}

.age-value {
  font-size: 2rem;
  font-weight: 700;
  color: #1976D2;
}

.age-projections {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.age-projection-item {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.projection-period {
  min-width: 80px;
  font-weight: 600;
  color: #333;
  font-size: 1.1rem;
}

.projection-increase {
  min-width: 120px;
  font-size: 1.1rem;
  color: #555;
}

/* 風險緩解策略 */
.risk-mitigation-section {
  margin-bottom: 2rem;
}

.mitigation-card {
  padding: 1.5rem;
  border-radius: 16px !important;
  height: 100%;
  border: 2px solid transparent;
  transition: all 0.3s ease;
}

.mitigation-card:hover {
  transform: translateY(-2px);
  border-color: #4CAF50;
}

.mitigation-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.mitigation-category {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
}

.mitigation-actions {
  margin: 1rem 0;
}

.mitigation-action {
  display: flex;
  align-items: center;
  margin-bottom: 0.5rem;
  font-size: 1.1rem;
  color: #555;
}

.mitigation-timeframe {
  display: flex;
  align-items: center;
  font-size: 1rem;
  color: #666;
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #eee;
}

/* 競爭對手分析表格 */
.competitor-section {
  margin-bottom: 2rem;
}

.competitor-table {
  border-radius: 12px;
  overflow: hidden;
}

:deep(.competitor-table .v-data-table__wrapper) {
  border-radius: 12px;
}

:deep(.competitor-table th) {
  background: #f5f5f5 !important;
  font-weight: 600 !important;
  color: #333 !important;
}

.premium-cell {
  font-weight: 600;
  color: #1976D2;
}

.coverage-cell {
  font-weight: 600;
  color: #333;
}

.score-cell {
  min-width: 80px;
}

/* 專業建議 */
.professional-recommendations {
  margin-bottom: 1rem;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .insurance-dialog-content {
    padding: 1rem !important;
  }
  
  .age-risk-chart {
    flex-direction: column;
    gap: 1rem;
  }
  
  .age-projections {
    width: 100%;
  }
  
  .age-projection-item {
    flex-direction: column;
    text-align: center;
  }
  
  .risk-score-card {
    padding: 1.5rem;
  }
  
  .risk-score-number {
    font-size: 1.5rem;
  }
  
  /* 調整保險風險分析字體大小 */
  .category-name {
    font-size: 1.2rem;
  }
  
  .category-description {
    font-size: 1rem;
  }
  
  .metric-name {
    font-size: 1rem;
  }
  
  .metric-value {
    font-size: 1.3rem;
  }
  
  .mitigation-category {
    font-size: 1.2rem;
  }
  
  .mitigation-action {
    font-size: 1rem;
  }
  
  .projection-period,
  .projection-increase {
    font-size: 1rem;
  }
}

@media (max-width: 600px) {
  .risk-overview-section .v-col {
    margin-bottom: 1rem;
  }
}
</style> 