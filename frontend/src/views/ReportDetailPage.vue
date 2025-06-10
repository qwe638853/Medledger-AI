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
  // 基本生化檢查
  'Glu-AC': { min: 70, max: 100, unit: 'mg/dL' },
  'HbA1c': { min: 4.0, max: 6.0, unit: '%' },
  'Glu-PC': { max: 140, unit: 'mg/dL' },
  
  // 脂質檢查
  'LDL-C': { max: 130, unit: 'mg/dL' },
  'HDL-C': { min: 40, unit: 'mg/dL' },
  'TG': { max: 150, unit: 'mg/dL' },
  'T-CHO': { max: 200, unit: 'mg/dL' },
  
  // 蛋白質檢查
  'Alb': { min: 3.5, max: 5.0, unit: 'g/dL' },
  'TP': { min: 6.0, max: 8.3, unit: 'g/dL' },
  
  // 肝功能檢查
  'AST（GOT）': { min: 5, max: 40, unit: 'U/L' },
  'ALT（GPT）': { min: 5, max: 35, unit: 'U/L' },
  'ALP': { min: 30, max: 120, unit: 'U/L' },
  'T-Bil': { min: 0.2, max: 1.2, unit: 'mg/dL' },
  'D-Bil': { min: 0.0, max: 0.3, unit: 'mg/dL' },
  
  // 腎功能檢查
  'UN': { min: 7, max: 20, unit: 'mg/dL' },
  'CRE': { min: 0.6, max: 1.2, unit: 'mg/dL' },
  'U.A': { min: 2.5, max: 7.0, unit: 'mg/dL' },
  
  // 血液常規檢查
  'Hb': { min: 11.0, max: 17.2, unit: 'g/dL' }, // 男13.1~17.2 女11.0~15.2
  'RBC': { min: 3.78, max: 5.9, unit: '10^6/uL' }, // 男4.21~5.9 女3.78~5.25
  'WBC': { min: 3.25, max: 9.16, unit: '10^3/uL' },
  'Hct': { min: 34.8, max: 51.5, unit: '%' }, // 男39.6~51.5 女34.8~46.3
  'PLT': { min: 150, max: 378, unit: '10^3/uL' },
  'Platelet': { min: 150, max: 378, unit: '10^3/uL' },
  'MCV': { min: 80.9, max: 99.3, unit: 'fL' },
  'MCH': { min: 25.5, max: 33.2, unit: 'pg' },
  'MCHC': { min: 31.0, max: 34.9, unit: 'g/dL' },
  'RDW-CV': { min: 11.6, max: 15.0, unit: '%' },
  
  // 白血球分類
  'Reticulocyte': { min: 0.87, max: 2.50, unit: '%' }, // 男1.05~2.50 女0.87~2.48
  'Neutrophils (seg)': { min: 41.6, max: 74.4, unit: '%' },
  'Lymphocytes (Lym)': { min: 18.0, max: 48.8, unit: '%' },
  'Monocytes (Mono)': { min: 3.3, max: 8.9, unit: '%' },
  'Eosinophils (Eso)': { min: 0.3, max: 7.9, unit: '%' },
  'Basophils (Baso)': { min: 0.2, max: 1.6, unit: '%' },
  
  // 凝血功能檢查
  'PT': { min: 9.7, max: 11.8, unit: 'sec' },
  'aPTT': { min: 25.6, max: 32.6, unit: 'sec' },
  
  // 發炎指標
  'E.S.R.': { min: 2, max: 15, unit: 'mm/hr' }, // 男2~10 女2~15
  'hsCRP': { min: 0.0, max: 3.0, unit: 'mg/L' },
  
  // 腫瘤標記
  'AFP': { min: 0.0, max: 20.0, unit: 'ng/mL' },
  'CEA': { min: 0.0, max: 5.0, unit: 'ng/mL' },
  'CA-125': { min: 0.0, max: 35.0, unit: 'U/mL' },
  'CA19-9': { min: 0.0, max: 37.0, unit: 'U/mL' },
  
  // 尿液檢查
  'Specific Gravity': { min: 1.003, max: 1.035 },
  'Specific Gravity (Dipstick)': { min: 1.003, max: 1.035 },
  'Color-Appearance': { normal: '淡黃色清澈' },
  'PH': { min: 5.0, max: 8.0 },
  'PH (Dipstick)': { min: 5.0, max: 8.0 },
  'Protein (Dipstick)': { normal: '陰性 (-)' },
  'Glucose (Dipstick)': { normal: '陰性 (-)' },
  'Bilirubin (Dipstick)': { normal: '陰性 (-)' },
  'Urobilinogen (Dipstick)': { max: 1.5, unit: 'mg/dL' },
  'Ketone': { normal: '陰性 (-)' },
  'Ketone(Dipstick)': { normal: '陰性 (-)' },
  'Nitrite': { normal: '陰性 (-)' },
  'Nitrite(Dipstick)': { normal: '陰性 (-)' },
  'Occult Blood': { normal: '陰性 (-)' },
  'O.B.(Dipstick)': { normal: '陰性 (-)' },
  'WBC Esterase': { normal: '陰性 (-)' },
  'WBC esterase(Dipstick)': { normal: '陰性 (-)' },
  
  // 尿沉渣檢查
  'RBC (Urine)': { min: 0, max: 2, unit: '/HPF' },
  'RBC (Sediment)': { min: 0, max: 2, unit: '/HPF' },
  'WBC (Urine)': { min: 0, max: 5, unit: '/HPF' },
  'WBC (Sediment)': { min: 0, max: 5, unit: '/HPF' },
  'Epith Cell (Sediment)': { min: 0, max: 5, unit: '/HPF' },
  'Casts (Sediment)': { min: 0, max: 2, unit: '/LPF' },
  'Crystal (Sediment)': { normal: '少量' },
  'Bacteria (Sediment)': { normal: '少量' },
  'Albumin (Dipstick)': { max: 30, unit: 'mg/g' },
  'Creatinine (Dipstick)': { min: 30, max: 300, unit: 'mg/dL' },
  'Albumin / Creatinine Ratio (Dipstick)': { max: 30, unit: 'mg/g' },
};
function getMetricColor(key, value) {
  const ref = METRIC_REF_RANGE[key];
  if (!ref) return 'grey';
  
  // 處理只有 normal 值的項目（如尿液檢查）
  if (ref.normal !== undefined) {
    const normalizedValue = value?.toString().toLowerCase().trim();
    const normalizedNormal = ref.normal?.toString().toLowerCase().trim();
    
    // 檢查是否為正常值
    if (normalizedValue === normalizedNormal || 
        normalizedValue === '陰性' || 
        normalizedValue === 'negative' || 
        normalizedValue === '-' ||
        normalizedValue === '(-)') {
      return 'green';
    } else if (normalizedValue === '陽性' || 
               normalizedValue === 'positive' || 
               normalizedValue === '+' ||
               normalizedValue?.includes('+')) {
      return 'red';
    } else {
      return 'orange'; // 其他異常值
    }
  }
  
  // 處理數值型項目
  const match = (value || '').toString().match(/-?\d+(\.\d+)?/);
  const num = match ? parseFloat(match[0]) : NaN;
  if (isNaN(num)) return 'grey';
  
  // 檢查是否超出正常範圍
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
  riskScore: 52,
  overallRiskLevel: 'low', // low, medium, high
  riskCategories: {
    cardiovascular: {
      score: 45,
      level: 'low',
      factors: ['血壓正常(127/61)', '血脂正常', 'HDL-C充足(54)'],
      impact: 'low',
      description: '心血管風險極低，各項指標表現優秀'
    },
    diabetes: {
      score: 25,
      level: 'low',
      factors: ['空腹血糖89正常', '糖化血色素4.1%優秀', '飯後血糖124正常'],
      impact: 'low',
      description: '糖尿病風險極低，血糖控制優良'
    },
    kidney: {
      score: 68,
      level: 'medium',
      factors: ['尿素氮略高(23)', '肌酸酐正常上限(1.2)', '尿液檢查正常'],
      impact: 'low',
      description: '腎功能需要關注，建議定期追蹤'
    },
    liver: {
      score: 20,
      level: 'low',
      factors: ['ALT正常(10)', 'AST正常(27)', '膽紅素正常'],
      impact: 'low',
      description: '肝功能正常，風險極低'
    },
    cancer: {
      score: 15,
      level: 'low',
      factors: ['AFP正常(14)', 'CEA正常(2.8)', 'CA-125正常(28)', 'CA19-9正常(29)'],
      impact: 'low',
      description: '腫瘤標記全部正常，癌症風險極低'
    }
  },
  healthMetrics: [
    { name: '血壓', value: 127, unit: 'mmHg', status: 'normal', weight: 0.20 },
    { name: '血糖', value: 89, unit: 'mg/dL', status: 'normal', weight: 0.25 },
    { name: '膽固醇', value: 164, unit: 'mg/dL', status: 'normal', weight: 0.20 },
    { name: '腎功能', value: 23, unit: 'mg/dL', status: 'elevated', weight: 0.25 },
    { name: '肝功能', value: 27, unit: 'U/L', status: 'normal', weight: 0.10 }
  ],
  ageRiskFactors: {
    currentAge: 35,
    riskIncrease: {
      '5years': 8,
      '10years': 18,
      '15years': 32
    }
  },
  riskMitigation: [
    {
      category: '腎功能保養',
      actions: ['每日充足飲水2500ml', '定期腎功能檢查', '控制蛋白質攝取'],
      riskReduction: 15,
      timeframe: '3個月'
    },
    {
      category: '生活習慣維持',
      actions: ['保持規律運動', '維持健康體重', '充足睡眠'],
      riskReduction: 10,
      timeframe: '持續進行'
    },
    {
      category: '定期健康監測',
      actions: ['每季度腎功能追蹤', '年度完整健檢', '血壓日常監測'],
      riskReduction: 20,
      timeframe: '持續進行'
    }
  ],
  recommendations: [
    {
      type: 'immediate',
      title: '立即建議',
      items: ['標準費率承保', '無需額外體檢', '建議腎功能定期追蹤', '整體健康狀況優良']
    },
    {
      type: 'monitoring',
      title: '持續監控',
      items: ['每3個月追蹤腎功能指標', '年度健康檢查', '維持健康生活習慣', '血壓日常監測']
    },
    {
      type: 'assessment',
      title: '風險評估',
      items: ['整體風險評級：低風險', '可提供健康優惠費率', '適合標準保險產品', '長期風險展望良好']
    }
  ]
});

// AI 分析和保單推薦
const aiAnalysis = ref({
  summary: '根據您的健康檢查結果，整體健康狀況優良。血糖、肝功能、血脂等主要指標均在正常範圍內，僅尿素氮(UN)略微偏高，建議注意腎功能保養和適量飲水。血壓控制良好，腫瘤標記正常，顯示低癌症風險。',
  healthScore: 92,
  riskLevel: 'low', // low, medium, high
  diseaseRisks: [
    {
      name: '心血管疾病',
      risk: 18,
      level: 'low',
      factors: ['血壓正常', '血脂控制良好', 'HDL-C達標'],
      prevention: '維持現有生活習慣，定期監測血壓'
    },
    {
      name: '糖尿病',
      risk: 12,
      level: 'low',
      factors: ['空腹血糖正常', '糖化血色素4.1%', '飯後血糖正常'],
      prevention: '繼續保持健康飲食和規律運動'
    },
    {
      name: '腎功能異常',
      risk: 28,
      level: 'low',
      factors: ['尿素氮略高(23)', '肌酸酐正常', '尿液檢查大致正常'],
      prevention: '增加水分攝取，定期追蹤腎功能指標'
    },
    {
      name: '肝功能異常',
      risk: 8,
      level: 'low',
      factors: ['ALT正常(10)', 'AST正常(27)', '膽紅素正常'],
      prevention: '維持健康生活習慣即可'
    }
  ],
  healthTrends: [
    { metric: '血糖控制', trend: 'stable', change: 0 },
    { metric: '血脂狀況', trend: 'stable', change: 0 },
    { metric: '腎功能', trend: 'monitoring', change: 5 },
    { metric: '血壓控制', trend: 'stable', change: 0 }
  ],
  recommendations: [
    {
      type: 'diet',
      title: '飲食建議',
      items: ['每日飲水量2000-2500毫升', '適量減少蛋白質攝取', '多食用新鮮蔬果', '控制鈉鹽攝入量']
    },
    {
      type: 'exercise',
      title: '運動建議',
      items: ['維持每週150分鐘中等強度運動', '適度有氧運動如快走、游泳', '避免過度劇烈運動', '運動後適當補充水分']
    },
    {
      type: 'lifestyle',
      title: '生活習慣',
      items: ['保持充足睡眠7-8小時', '定期監測血壓', '每3個月追蹤腎功能', '避免熬夜和過度疲勞']
    }
  ],
  insuranceRecommendations: [
    {
      name: '標準健康保障計劃',
      coverage: '住院醫療、手術、門診',
      features: ['標準費率承保', '無額外體檢要求', '全面醫療保障'],
      monthlyPremium: 2200,
      suitability: 95
    },
    {
      name: '優質健康防護計劃',
      coverage: '重大疾病、住院醫療、預防保健',
      features: ['健康優惠費率', '預防保健給付', '家族保障方案'],
      monthlyPremium: 2800,
      suitability: 88
    }
  ]
});

// 根據不同角色使用不同的 API endpoint
const fetchReportData = async () => {
  loading.value = true;
  errorMsg.value = '';
  
  try {
    let response;
    console.log("當前角色:", userRole.value);
    console.log("報告ID:", reportId);
    
    if (userRole.value === 'insurer') {
      // 保險公司調用授權報告 API
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
      console.log("保險公司獲取的報告數據:", report.value);
    } else {
      // 一般用戶調用自己的報告 API
      console.log('一般用戶正在調用 HandleReadMyReport API...');
      
      // 先檢查 store 中是否有基本數據
      const storeReport = userStore.currentReport;
      console.log('Store 中的報告數據:', storeReport);
      
      if (storeReport && storeReport.rawData && Object.keys(storeReport.rawData).length > 0) {
        // 如果 store 中已有完整數據，直接使用
        report.value = storeReport;
        console.log('使用 Store 中的完整數據:', report.value);
      } else {
        // 否則調用 API 獲取完整報告內容
        try {
          response = await healthCheckService.fetchReportDetail(reportId);
          
          if (response && response.success && response.resultJson) {
            const parsedData = JSON.parse(response.resultJson);
            
            report.value = {
              id: reportId,
              patient_id: patientId || userStore.user?.id,
              date: storeReport?.date || new Date().toISOString(),
              clinic_id: storeReport?.clinic_id || '未知診所',
              content: storeReport?.content || '健康檢查報告',
              rawData: parsedData
            };
            
            console.log('一般用戶獲取的完整報告數據:', report.value);
            
            // 更新 store 中的數據
            userStore.setCurrentReport(report.value);
          } else {
            throw new Error('API 回應格式異常或無數據');
          }
        } catch (apiError) {
          console.error('調用 HandleReadMyReport 失敗:', apiError);
          
          // 如果 API 調用失敗，嘗試使用 store 中的基本數據
          if (storeReport) {
            report.value = storeReport;
            console.log('API 失敗，使用 Store 中的基本數據:', report.value);
            errorMsg.value = '無法獲取最新報告數據，顯示基本資訊';
          } else {
            throw new Error('無法獲取報告數據：' + apiError.message);
          }
        }
      }
    }
    
    // 如果有數據，進行風險評估
    if (report.value?.rawData && Object.keys(report.value.rawData).length > 0) {
      console.log('開始進行風險評估:', report.value.rawData);
      evaluateRisk(report.value.rawData);
    } else {
      console.warn('沒有足夠的數據進行風險評估');
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
  
  // 處理只有 normal 值的項目（如尿液檢查）
  if (ref.normal !== undefined) {
    const normalizedValue = value?.toString().toLowerCase().trim();
    const normalizedNormal = ref.normal?.toString().toLowerCase().trim();
    
    // 如果是正常值，返回 100%（滿圓）
    if (normalizedValue === normalizedNormal || 
        normalizedValue === '陰性' || 
        normalizedValue === 'negative' || 
        normalizedValue === '-' ||
        normalizedValue === '(-)') {
      return 100;
    } else {
      return 75; // 異常值顯示 75%
    }
  }
  
  // 處理數值型項目
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
                      <template v-if="METRIC_REF_RANGE[key].normal !== undefined">
                        {{ METRIC_REF_RANGE[key].normal }}
                      </template>
                      <template v-else>
                        <template v-if="METRIC_REF_RANGE[key].min !== undefined">
                          {{ METRIC_REF_RANGE[key].min }}
                        </template>
                        <template v-else>0</template>
                        -
                        <template v-if="METRIC_REF_RANGE[key].max !== undefined">
                          {{ METRIC_REF_RANGE[key].max }}
                        </template>
                        <template v-else>∞</template>
                        <template v-if="METRIC_REF_RANGE[key].unit">
                          {{ METRIC_REF_RANGE[key].unit }}
                        </template>
                      </template>
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
    <v-dialog v-model="showAISummary" max-width="1000" scrollable>
      <v-card class="ai-dialog-card">
        <!-- 對話框標題 -->
        <v-card-title class="ai-dialog-header">
          <div class="ai-header-content">
            <div class="ai-header-left">
              <v-avatar class="ai-avatar mr-4" size="56" color="gradient">
                <v-icon color="white" size="28">mdi-robot-outline</v-icon>
              </v-avatar>
              <div class="ai-header-text">
                <div class="ai-dialog-title">AI 智能健康分析</div>
                <div class="ai-dialog-subtitle">基於您的健康數據進行專業分析與建議</div>
              </div>
            </div>
            <v-btn
              icon
              variant="text"
              @click="showAISummary = false"
              class="ai-close-btn"
              size="large"
            >
              <v-icon size="24">mdi-close</v-icon>
            </v-btn>
          </div>
        </v-card-title>

        <v-card-text class="ai-dialog-content">
          <!-- 健康評分區域 -->
          <div class="health-score-wrapper">
            <div class="health-score-main">
              <v-progress-circular
                :model-value="aiAnalysis.healthScore"
                :color="aiAnalysis.healthScore > 80 ? '#4CAF50' : aiAnalysis.healthScore > 60 ? '#FF9800' : '#F44336'"
                size="140"
                width="14"
                class="health-score-circle"
              >
                <div class="health-score-content">
                  <div class="health-score-number">{{ aiAnalysis.healthScore }}</div>
                </div>
              </v-progress-circular>
            </div>
            <div class="health-score-info">
              <div class="health-score-title">整體健康評分</div>
              <v-chip
                :color="aiAnalysis.riskLevel === 'low' ? 'success' : aiAnalysis.riskLevel === 'medium' ? 'warning' : 'error'"
                size="large"
                class="health-level-chip"
                variant="flat"
              >
                <v-icon start size="18">
                  {{ aiAnalysis.riskLevel === 'low' ? 'mdi-check-circle' : aiAnalysis.riskLevel === 'medium' ? 'mdi-alert-circle' : 'mdi-close-circle' }}
                </v-icon>
                {{ aiAnalysis.riskLevel === 'low' ? '健康狀況良好' : aiAnalysis.riskLevel === 'medium' ? '需要注意' : '需要關注' }}
              </v-chip>
              <div class="health-score-description">
                您的健康指標整體表現
                {{ aiAnalysis.healthScore > 80 ? '優秀' : aiAnalysis.healthScore > 60 ? '良好' : '需要改善' }}，
                建議持續關注並改善生活習慣。
              </div>
            </div>
          </div>

          <!-- AI 分析摘要 -->
          <div class="ai-summary-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large">
                <v-icon size="32" color="white">mdi-brain</v-icon>
              </div>
              <div class="section-title-large">AI 專業分析</div>
            </div>
            <div class="ai-summary-card">
              <div class="ai-summary-content">
                {{ aiAnalysis.summary }}
              </div>
            </div>
          </div>

          <!-- 疾病風險分析 -->
          <div class="disease-risk-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large risk-icon">
                <v-icon size="32" color="white">mdi-shield-alert</v-icon>
              </div>
              <div class="section-title-large">疾病風險評估</div>
            </div>
            <div class="disease-risk-cards">
              <div
                v-for="(disease, index) in aiAnalysis.diseaseRisks"
                :key="index"
                class="disease-card"
              >
                <div class="disease-card-header">
                  <div class="disease-info">
                    <div class="disease-name-large">{{ disease.name }}</div>
                    <div class="disease-risk-level">
                      <span class="risk-percentage">{{ disease.risk }}%</span>
                      <span class="risk-label">風險機率</span>
                    </div>
                  </div>
                  <div class="risk-indicator">
                    <v-chip
                      :color="disease.level === 'low' ? 'success' : disease.level === 'medium' ? 'warning' : 'error'"
                      size="large"
                      variant="flat"
                    >
                      {{ disease.level === 'low' ? '低風險' : disease.level === 'medium' ? '中風險' : '高風險' }}
                    </v-chip>
                  </div>
                </div>
                
                <div class="risk-progress-wrapper">
                  <v-progress-linear
                    :model-value="disease.risk"
                    :color="disease.level === 'low' ? 'success' : disease.level === 'medium' ? 'warning' : 'error'"
                    height="12"
                    rounded
                    class="risk-progress-bar"
                  ></v-progress-linear>
                </div>

                <div class="disease-details-large">
                  <div class="factors-section">
                    <div class="detail-label-large">主要影響因子</div>
                    <div class="factor-chips">
                      <v-chip
                        v-for="factor in disease.factors"
                        :key="factor"
                        size="default"
                        color="blue"
                        variant="tonal"
                        class="factor-chip"
                      >
                        {{ factor }}
                      </v-chip>
                    </div>
                  </div>
                  <div class="prevention-section">
                    <div class="detail-label-large">預防建議</div>
                    <div class="prevention-text-large">{{ disease.prevention }}</div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 健康趨勢 -->
          <div class="health-trends-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large trends-icon">
                <v-icon size="32" color="white">mdi-trending-up</v-icon>
              </div>
              <div class="section-title-large">健康趨勢分析</div>
            </div>
            <div class="trends-cards">
              <div
                v-for="trend in aiAnalysis.healthTrends"
                :key="trend.metric"
                class="trend-card"
              >
                <div class="trend-indicator-icon">
                  <v-icon
                    :color="trend.trend === 'improving' ? '#4CAF50' : trend.trend === 'concern' ? '#F44336' : '#9E9E9E'"
                    size="32"
                  >
                    {{ trend.trend === 'improving' ? 'mdi-trending-up' : trend.trend === 'concern' ? 'mdi-trending-down' : 'mdi-trending-neutral' }}
                  </v-icon>
                </div>
                <div class="trend-content">
                  <div class="trend-metric-large">{{ trend.metric }}</div>
                  <div class="trend-status">
                    <span class="trend-text-large">
                      {{ trend.trend === 'improving' ? '持續改善' : trend.trend === 'concern' ? '需要關注' : '保持穩定' }}
                    </span>
                    <span class="trend-change" v-if="trend.change !== 0">
                      {{ trend.change > 0 ? '+' : '' }}{{ trend.change }}%
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 健康建議 -->
          <div class="recommendations-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large recommendations-icon">
                <v-icon size="32" color="white">mdi-lightbulb</v-icon>
              </div>
              <div class="section-title-large">個人化健康建議</div>
            </div>
            
            <div class="recommendations-tabs-wrapper">
              <v-tabs v-model="activeTab" color="#00B8D9" align-tabs="center" class="custom-tabs">
                <v-tab
                  v-for="rec in aiAnalysis.recommendations"
                  :key="rec.type"
                  :value="rec.type"
                  class="recommendation-tab"
                >
                  <v-icon start size="24">
                    {{ rec.type === 'diet' ? 'mdi-food-apple' : rec.type === 'exercise' ? 'mdi-run' : 'mdi-heart' }}
                  </v-icon>
                  {{ rec.title }}
                </v-tab>
              </v-tabs>
              
              <v-window v-model="activeTab" class="recommendations-window">
                <v-window-item
                  v-for="rec in aiAnalysis.recommendations"
                  :key="rec.type"
                  :value="rec.type"
                >
                  <div class="recommendation-content">
                    <div
                      v-for="(item, index) in rec.items"
                      :key="index"
                      class="recommendation-item-large"
                    >
                      <div class="recommendation-icon">
                        <v-icon color="#00B8D9" size="20">mdi-check-circle</v-icon>
                      </div>
                      <div class="recommendation-text">{{ item }}</div>
                    </div>
                  </div>
                </v-window-item>
              </v-window>
            </div>
          </div>

          <!-- 保險推薦 -->
          <div class="insurance-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large insurance-icon">
                <v-icon size="32" color="white">mdi-shield-account</v-icon>
              </div>
              <div class="section-title-large">智能保險推薦</div>
            </div>
            <div class="insurance-cards">
              <div
                v-for="(plan, index) in aiAnalysis.insuranceRecommendations"
                :key="index"
                class="insurance-plan-card"
              >
                <div class="insurance-card-header">
                  <div class="plan-info">
                    <div class="plan-name-large">{{ plan.name }}</div>
                    <div class="plan-coverage-large">{{ plan.coverage }}</div>
                  </div>
                  <div class="plan-pricing">
                    <div class="premium-amount-large">NT$ {{ plan.monthlyPremium.toLocaleString() }}</div>
                    <div class="premium-period-large">/月</div>
                  </div>
                </div>
                
                <div class="suitability-wrapper">
                  <div class="suitability-header">
                    <span class="suitability-label-large">適合度評分</span>
                    <span class="suitability-score">{{ plan.suitability }}%</span>
                  </div>
                  <v-progress-linear
                    :model-value="plan.suitability"
                    color="#00B8D9"
                    height="8"
                    rounded
                    class="suitability-progress"
                  ></v-progress-linear>
                </div>

                <div class="plan-features-large">
                  <div class="feature-tags-title">方案特色</div>
                  <div class="feature-tags-container">
                    <v-chip
                      v-for="feature in plan.features"
                      :key="feature"
                      size="large"
                      color="primary"
                      variant="flat"
                      class="feature-chip-prominent"
                    >
                      <v-icon start size="18" color="white">
                        {{ feature.includes('免等待期') ? 'mdi-clock-fast' : 
                           feature.includes('一年一約') ? 'mdi-calendar-check' : 
                           feature.includes('可選擇醫院') ? 'mdi-hospital-building' : 
                           feature.includes('保額') ? 'mdi-shield-check' : 
                           feature.includes('理賠') ? 'mdi-lightning-bolt' : 
                           feature.includes('豁免') ? 'mdi-hand-heart' : 'mdi-check-circle' }}
                      </v-icon>
                      {{ feature }}
                    </v-chip>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- 保險風險分析彈窗 -->
    <v-dialog v-model="showRisk" max-width="1200" scrollable>
      <v-card class="insurance-dialog-card">
        <!-- 對話框標題 -->
        <v-card-title class="insurance-dialog-header">
          <div class="insurance-header-content">
            <div class="insurance-header-left">
              <v-avatar class="insurance-avatar mr-4" size="56" color="gradient">
                <v-icon color="white" size="28">mdi-shield-check</v-icon>
              </v-avatar>
              <div class="insurance-header-text">
                <div class="insurance-dialog-title">專業保險風險評估</div>
                <div class="insurance-dialog-subtitle">基於健康數據的精準風險分析與保費定價</div>
              </div>
            </div>
            <v-btn
              icon
              variant="text"
              @click="showRisk = false"
              class="insurance-close-btn"
              size="large"
            >
              <v-icon size="24">mdi-close</v-icon>
            </v-btn>
          </div>
        </v-card-title>

        <v-card-text class="insurance-dialog-content">
          <!-- 風險評分總覽 -->
          <div class="risk-overview-wrapper">
            <div class="risk-score-main">
              <v-progress-circular
                :model-value="insuranceAnalysis.riskScore"
                :color="insuranceAnalysis.riskScore > 80 ? '#F44336' : insuranceAnalysis.riskScore > 60 ? '#FF9800' : '#4CAF50'"
                size="160"
                width="16"
                class="risk-score-circle"
              >
                <div class="risk-score-content">
                  <div class="risk-score-number">{{ insuranceAnalysis.riskScore }}</div>
                </div>
              </v-progress-circular>
            </div>
            <div class="risk-score-info">
              <div class="risk-score-title">整體風險評分</div>
              <v-chip
                :color="insuranceAnalysis.overallRiskLevel === 'low' ? 'success' : insuranceAnalysis.overallRiskLevel === 'medium' ? 'warning' : 'error'"
                size="large"
                class="risk-level-chip"
                variant="flat"
              >
                <v-icon start size="18">
                  {{ insuranceAnalysis.overallRiskLevel === 'low' ? 'mdi-shield-check' : insuranceAnalysis.overallRiskLevel === 'medium' ? 'mdi-shield-alert' : 'mdi-shield-remove' }}
                </v-icon>
                {{ insuranceAnalysis.overallRiskLevel === 'low' ? '低風險客戶' : insuranceAnalysis.overallRiskLevel === 'medium' ? '中等風險客戶' : '高風險客戶' }}
              </v-chip>
              <div class="risk-score-description">
                根據健康指標分析，該客戶屬於
                {{ insuranceAnalysis.overallRiskLevel === 'low' ? '標準承保' : insuranceAnalysis.overallRiskLevel === 'medium' ? '加費承保' : '特殊核保' }}
                類別，建議相應調整保費策略。
              </div>
              <div class="risk-metrics-summary">
                <div class="summary-item">
                  <div class="summary-label">建議保費調整</div>
                  <div class="summary-value premium-adjustment">
                    {{ insuranceAnalysis.riskScore > 80 ? '+35%' : insuranceAnalysis.riskScore > 60 ? '+15%' : '標準費率' }}
                  </div>
                </div>
                <div class="summary-item">
                  <div class="summary-label">核保建議</div>
                  <div class="summary-value underwriting-advice">
                    {{ insuranceAnalysis.riskScore > 80 ? '需體檢' : insuranceAnalysis.riskScore > 60 ? '加強審核' : '標準承保' }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 疾病風險分類評估 -->
          <div class="risk-categories-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large risk-category-icon">
                <v-icon size="32" color="white">mdi-chart-donut</v-icon>
              </div>
              <div class="section-title-large">疾病風險分類評估</div>
            </div>
            <div class="risk-category-cards">
              <div
                v-for="(category, key) in insuranceAnalysis.riskCategories"
                :key="key"
                class="category-card"
              >
                <div class="category-header-large">
                  <div class="category-info">
                    <div class="category-name-large">
                      {{ key === 'cardiovascular' ? '心血管疾病' : 
                         key === 'diabetes' ? '糖尿病' : 
                         key === 'kidney' ? '腎臟疾病' :
                         key === 'liver' ? '肝臟疾病' : 
                         key === 'cancer' ? '癌症' : key }}
                    </div>
                    <div class="category-risk-score">
                      <span class="score-number">{{ category.score }}</span>
                      <span class="score-label">風險分</span>
                    </div>
                  </div>
                  <div class="category-level">
                    <v-chip
                      :color="category.level === 'low' ? 'success' : category.level === 'medium' ? 'warning' : 'error'"
                      size="large"
                      variant="flat"
                    >
                      {{ category.level === 'low' ? '低風險' : category.level === 'medium' ? '中風險' : '高風險' }}
                    </v-chip>
                  </div>
                </div>
                
                <div class="category-progress-wrapper">
                  <v-progress-linear
                    :model-value="category.score"
                    :color="category.level === 'low' ? 'success' : category.level === 'medium' ? 'warning' : 'error'"
                    height="14"
                    rounded
                    class="category-progress-bar"
                  ></v-progress-linear>
                </div>

                <div class="category-details-large">
                  <div class="category-description-large">{{ category.description }}</div>
                  <div class="category-factors-section">
                    <div class="factors-label-large">關鍵風險因素</div>
                    <div class="factors-chips-large">
                      <v-chip
                        v-for="factor in category.factors"
                        :key="factor"
                        size="default"
                        color="orange"
                        variant="tonal"
                        class="factor-chip-large"
                      >
                        {{ factor }}
                      </v-chip>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 健康指標影響分析 -->
          <div class="health-metrics-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large metrics-icon">
                <v-icon size="32" color="white">mdi-heart-pulse</v-icon>
              </div>
              <div class="section-title-large">健康指標風險權重</div>
            </div>
            <div class="metrics-cards">
              <div
                v-for="metric in insuranceAnalysis.healthMetrics"
                :key="metric.name"
                class="metric-card-large"
              >
                <div class="metric-header-large">
                  <div class="metric-info">
                    <div class="metric-name-large">{{ metric.name }}</div>
                    <div class="metric-value-display">
                      <span class="metric-number">{{ metric.value }}</span>
                      <span class="metric-unit-text">{{ metric.unit }}</span>
                    </div>
                  </div>
                  <div class="metric-status-icon">
                    <v-icon
                      :color="metric.status === 'normal' ? '#4CAF50' : '#FF9800'"
                      size="32"
                    >
                      {{ metric.status === 'normal' ? 'mdi-check-circle' : 'mdi-alert-circle' }}
                    </v-icon>
                  </div>
                </div>
                
                <div class="metric-weight-section">
                  <div class="weight-info">
                    <span class="weight-label-large">風險權重</span>
                    <span class="weight-value-large">{{ (metric.weight * 100).toFixed(0) }}%</span>
                  </div>
                  <v-progress-linear
                    :model-value="metric.weight * 100"
                    color="#1976D2"
                    height="10"
                    rounded
                    class="weight-progress"
                  ></v-progress-linear>
                </div>
              </div>
            </div>
          </div>

          

          <!-- 風險緩解策略 -->
          <div class="mitigation-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large mitigation-icon">
                <v-icon size="32" color="white">mdi-shield-plus</v-icon>
              </div>
              <div class="section-title-large">風險緩解策略建議</div>
            </div>
            <div class="mitigation-cards">
              <div
                v-for="strategy in insuranceAnalysis.riskMitigation"
                :key="strategy.category"
                class="mitigation-card-large"
              >
                <div class="mitigation-header-large">
                  <div class="mitigation-info">
                    <div class="mitigation-category-large">{{ strategy.category }}</div>
                    <div class="mitigation-timeframe-large">
                      <v-icon size="18" class="mr-1">mdi-clock-outline</v-icon>
                      {{ strategy.timeframe }}
                    </div>
                  </div>
                  <div class="mitigation-reduction">
                    <v-chip color="success" size="large" variant="flat">
                      <v-icon start size="16">mdi-trending-down</v-icon>
                      -{{ strategy.riskReduction }}% 風險
                    </v-chip>
                  </div>
                </div>
                <div class="mitigation-actions-large">
                  <div
                    v-for="action in strategy.actions"
                    :key="action"
                    class="mitigation-action-large"
                  >
                    <v-icon color="success" size="20" class="mr-3">mdi-check-bold</v-icon>
                    <span class="action-text">{{ action }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- 專業核保建議 -->
          <div class="professional-recommendations-wrapper">
            <div class="section-header-large">
              <div class="section-icon-large recommendations-icon">
                <v-icon size="32" color="white">mdi-lightbulb-on</v-icon>
              </div>
              <div class="section-title-large">專業核保建議</div>
            </div>
            
            <div class="recommendations-tabs-wrapper">
              <v-tabs v-model="insuranceActiveTab" color="#1976D2" align-tabs="center" class="insurance-custom-tabs">
                <v-tab
                  v-for="rec in insuranceAnalysis.recommendations"
                  :key="rec.type"
                  :value="rec.type"
                  class="insurance-recommendation-tab"
                >
                  <v-icon start size="24">
                    {{ rec.type === 'immediate' ? 'mdi-lightning-bolt' : rec.type === 'monitoring' ? 'mdi-monitor-eye' : 'mdi-clipboard-check' }}
                  </v-icon>
                  {{ rec.title }}
                </v-tab>
              </v-tabs>
              
              <v-window v-model="insuranceActiveTab" class="insurance-recommendations-window">
                <v-window-item
                  v-for="rec in insuranceAnalysis.recommendations"
                  :key="rec.type"
                  :value="rec.type"
                >
                  <div class="insurance-recommendation-content">
                    <div
                      v-for="(item, index) in rec.items"
                      :key="index"
                      class="insurance-recommendation-item-large"
                    >
                      <div class="insurance-recommendation-icon">
                        <v-icon color="#1976D2" size="20">mdi-arrow-right-circle</v-icon>
                      </div>
                      <div class="insurance-recommendation-text">{{ item }}</div>
                    </div>
                  </div>
                </v-window-item>
              </v-window>
            </div>
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
  border-radius: 20px !important;
  font-weight: 700 !important;
  padding: 0 32px !important;
  height: 56px !important;
  font-size: 1.2rem !important;
  min-width: 140px !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  box-shadow: 0 6px 20px rgba(0, 184, 217, 0.3) !important;
}

.back-btn:hover {
  transform: translateY(-3px) scale(1.02);
  box-shadow: 0 12px 32px rgba(0, 184, 217, 0.4) !important;
  background-color: #0093A6 !important;
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
  gap: 2rem;
  justify-content: center;
  margin-bottom: 3rem;
}

.action-btn {
  border-radius: 20px !important;
  font-weight: 700 !important;
  padding: 0 48px !important;
  height: 64px !important;
  font-size: 1.3rem !important;
  min-width: 220px !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15) !important;
  letter-spacing: 0.5px !important;
}

.ai-btn {
  background: linear-gradient(135deg, #00B8D9 0%, #0093A6 100%) !important;
  color: white !important;
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.4) !important;
}

.risk-btn {
  background: linear-gradient(135deg, #00B8D9 0%, #0093A6 100%) !important;
  color: white !important;
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.4) !important;
}

.action-btn:hover {
  transform: translateY(-4px) scale(1.05);
  box-shadow: 0 16px 40px rgba(0, 184, 217, 0.5) !important;
}

.ai-btn:hover {
  background: linear-gradient(135deg, #0093A6 0%, #007A8F 100%) !important;
}

.risk-btn:hover {
  background: linear-gradient(135deg, #0093A6 0%, #007A8F 100%) !important;
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

/* AI 彈窗樣式重構 - 統一主色調 */
.ai-dialog-card {
  border-radius: 28px !important;
  overflow: hidden !important;
  box-shadow: 0 20px 80px rgba(0, 0, 0, 0.12) !important;
}

.ai-dialog-header {
  background: linear-gradient(135deg, #00B8D9 0%, #0093A6 100%) !important;
  color: white !important;
  padding: 2rem !important;
  border: none !important;
}

.ai-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.ai-header-left {
  display: flex;
  align-items: center;
}

.ai-avatar {
  background: rgba(255, 255, 255, 0.2) !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15) !important;
}

.ai-header-text {
  margin-left: 0;
}

.ai-dialog-title {
  font-size: 1.8rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  letter-spacing: -0.5px;
}

.ai-dialog-subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
  font-weight: 400;
}

.ai-close-btn {
  color: rgba(255, 255, 255, 0.9) !important;
  background: rgba(255, 255, 255, 0.15) !important;
  border-radius: 16px !important;
  transition: all 0.3s ease !important;
  width: 56px !important;
  height: 56px !important;
  min-width: 56px !important;
}

.ai-close-btn:hover {
  background: rgba(255, 255, 255, 0.25) !important;
  transform: scale(1.1) rotate(90deg);
}

.ai-dialog-content {
  padding: 3rem !important;
  max-height: 75vh;
  background: #fafafa !important;
}

/* 健康評分區域重構 */
.health-score-wrapper {
  display: flex;
  align-items: center;
  gap: 3rem;
  background: white;
  padding: 3rem;
  border-radius: 24px;
  margin-bottom: 3rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.health-score-main {
  flex-shrink: 0;
}

.health-score-circle {
  filter: drop-shadow(0 8px 16px rgba(0, 0, 0, 0.1));
}

.health-score-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.health-score-number {
  font-size: 3rem;
  font-weight: 800;
  line-height: 1;
  color: inherit;
}

.health-score-unit {
  font-size: 1.2rem;
  font-weight: 500;
  opacity: 0.8;
  margin-top: 0.25rem;
}

.health-score-info {
  flex: 1;
}

.health-score-title {
  font-size: 2rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 1.5rem;
  letter-spacing: -0.5px;
}

.health-level-chip {
  font-size: 1.2rem !important;
  font-weight: 600 !important;
  padding: 1rem 1.5rem !important;
  height: auto !important;
  margin-bottom: 1.5rem;
}

.health-score-description {
  font-size: 1.3rem;
  line-height: 1.6;
  color: #666;
  font-weight: 400;
}

/* 大型區塊標題樣式 - 統一主色調 */
.section-header-large {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.section-icon-large {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #00B8D9 0%, #0093A6 100%);
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.3);
}

.section-icon-large.risk-icon {
  background: linear-gradient(135deg, #00B8D9 10%, #FF6B6B 90%);
  box-shadow: 0 8px 24px rgba(255, 107, 107, 0.3);
}

.section-icon-large.trends-icon {
  background: linear-gradient(135deg, #00B8D9 10%, #4ECDC4 90%);
  box-shadow: 0 8px 24px rgba(78, 205, 196, 0.3);
}

.section-icon-large.recommendations-icon {
  background: linear-gradient(135deg, #00B8D9 10%, #45B7D1 90%);
  box-shadow: 0 8px 24px rgba(69, 183, 209, 0.3);
}

.section-icon-large.insurance-icon {
  background: linear-gradient(135deg, #00B8D9 10%, #96CEB4 90%);
  box-shadow: 0 8px 24px rgba(150, 206, 180, 0.3);
}

.section-title-large {
  font-size: 2.2rem;
  font-weight: 700;
  color: #333;
  letter-spacing: -0.5px;
}

/* AI 分析摘要重構 */
.ai-summary-wrapper {
  margin-bottom: 3rem;
}

.ai-summary-card {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.ai-summary-content {
  font-size: 1.4rem;
  line-height: 1.8;
  color: #555;
  padding: 2.5rem;
  font-weight: 400;
  border-left: 6px solid #00B8D9;
  background: linear-gradient(135deg, #f0faff 0%, #e6f7ff 100%);
}

/* 疾病風險分析重構 */
.disease-risk-wrapper {
  margin-bottom: 3rem;
}

.disease-risk-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(420px, 1fr));
  gap: 2rem;
}

.disease-card {
  background: white;
  padding: 2.5rem;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 2px solid transparent;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.disease-card:hover {
  transform: translateY(-8px);
  border-color: #00B8D9;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
}

.disease-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  gap: 1.5rem;
}

.disease-info {
  flex: 1;
}

.disease-name-large {
  font-size: 1.6rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 1rem;
  letter-spacing: -0.5px;
}

.disease-risk-level {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.risk-percentage {
  font-size: 2.2rem;
  font-weight: 800;
  color: #00B8D9;
  line-height: 1;
}

.risk-label {
  font-size: 1rem;
  color: #888;
  font-weight: 500;
}

.risk-indicator {
  flex-shrink: 0;
}

.risk-progress-wrapper {
  margin: 1.5rem 0;
}

.risk-progress-bar {
  border-radius: 6px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.disease-details-large {
  margin-top: 2rem;
}

.factors-section,
.prevention-section {
  margin-bottom: 1.5rem;
}

.detail-label-large {
  font-size: 1.3rem;
  font-weight: 600;
  color: #444;
  margin-bottom: 1rem;
}

.factor-chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.factor-chip {
  font-size: 1rem !important;
  font-weight: 600 !important;
  padding: 0.5rem 1rem !important;
  height: auto !important;
  background: #e3f2fd !important;
  color: #1565c0 !important;
  border: 1px solid #bbdefb !important;
}

.prevention-text-large {
  font-size: 1.2rem;
  color: #555;
  line-height: 1.6;
  font-weight: 400;
}

/* 健康趨勢重構 */
.health-trends-wrapper {
  margin-bottom: 3rem;
}

.trends-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
}

.trend-card {
  background: white;
  padding: 2rem;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  gap: 1.5rem;
  transition: all 0.3s ease;
}

.trend-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
}

.trend-indicator-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  background: #f8f9fa;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.trend-content {
  flex: 1;
}

.trend-metric-large {
  font-size: 1.4rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 0.5rem;
}

.trend-status {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.trend-text-large {
  font-size: 1.1rem;
  font-weight: 500;
  color: #666;
}

.trend-change {
  font-size: 1rem;
  font-weight: 600;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  background: #f0faff;
  color: #00B8D9;
}

/* 健康建議重構 */
.recommendations-wrapper {
  margin-bottom: 3rem;
}

.recommendations-tabs-wrapper {
  background: white;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.custom-tabs {
  background: #f8f9fa !important;
}

.recommendation-tab {
  font-size: 1.3rem !important;
  font-weight: 700 !important;
  padding: 2rem 2.5rem !important;
  min-height: 72px !important;
  border-radius: 16px 16px 0 0 !important;
  transition: all 0.3s ease !important;
}

.recommendation-tab:hover {
  background: rgba(0, 184, 217, 0.1) !important;
  transform: translateY(-2px);
}

.recommendations-window {
  padding: 2.5rem;
}

.recommendation-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.recommendation-item-large {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.5rem;
  background: #f0faff;
  border-radius: 16px;
  border-left: 4px solid #00B8D9;
  transition: all 0.3s ease;
}

.recommendation-item-large:hover {
  background: #e6f7ff;
  transform: translateX(8px);
}

.recommendation-icon {
  flex-shrink: 0;
  margin-top: 0.25rem;
}

.recommendation-text {
  font-size: 1.2rem;
  line-height: 1.6;
  color: #555;
  font-weight: 400;
}

/* 保險推薦重構 */
.insurance-wrapper {
  margin-bottom: 2rem;
}

.insurance-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(450px, 1fr));
  gap: 2rem;
}

.insurance-plan-card {
  background: white;
  padding: 2.5rem;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 2px solid transparent;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.insurance-plan-card:hover {
  transform: translateY(-8px);
  border-color: #00B8D9;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
}

.insurance-card-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  gap: 2rem;
}

.plan-info {
  flex: 1;
}

.plan-name-large {
  font-size: 1.6rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 0.75rem;
  letter-spacing: -0.5px;
}

.plan-coverage-large {
  font-size: 1.2rem;
  color: #666;
  line-height: 1.5;
}

.plan-pricing {
  text-align: right;
  flex-shrink: 0;
}

.premium-amount-large {
  font-size: 2rem;
  font-weight: 800;
  color: #00B8D9;
  line-height: 1;
}

.premium-period-large {
  font-size: 1.1rem;
  color: #888;
  font-weight: 500;
}

.suitability-wrapper {
  margin: 2rem 0;
  padding: 1.5rem;
  background: #f0faff;
  border-radius: 16px;
}

.suitability-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.suitability-label-large {
  font-size: 1.2rem;
  font-weight: 600;
  color: #444;
}

.suitability-score {
  font-size: 1.4rem;
  font-weight: 700;
  color: #00B8D9;
}

.suitability-progress {
  border-radius: 4px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.plan-features-large {
  margin-top: 1.5rem;
}

.feature-tags-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid #f0faff;
}

.feature-tags-container {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.feature-chip-prominent {
  font-size: 1.1rem !important;
  font-weight: 600 !important;
  padding: 0.75rem 1.25rem !important;
  height: auto !important;
  min-height: 48px !important;
  background: linear-gradient(135deg, #00B8D9 0%, #0093A6 100%) !important;
  color: white !important;
  border-radius: 12px !important;
  box-shadow: 0 4px 16px rgba(0, 184, 217, 0.25) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  border: none !important;
}

.feature-chip-prominent:hover {
  transform: translateY(-2px) scale(1.02) !important;
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.35) !important;
}

.feature-chip-prominent .v-icon {
  margin-right: 0.5rem !important;
}

/* 舊的樣式移除 */
.feature-chip {
  display: none;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .ai-dialog-content {
    padding: 1.5rem !important;
  }
  
  .health-score-wrapper {
    flex-direction: column;
    text-align: center;
    gap: 2rem;
  }
  
  .section-header-large {
    flex-direction: column;
    text-align: center;
    gap: 1rem;
  }
  
  .section-title-large {
    font-size: 1.8rem;
  }
  
  .disease-risk-cards,
  .insurance-cards {
    grid-template-columns: 1fr;
  }
  
  .trends-cards {
    grid-template-columns: 1fr;
  }
  
  .disease-card-header,
  .insurance-card-header {
    flex-direction: column;
    gap: 1rem;
  }
  
  .plan-pricing {
    text-align: left;
  }
  
  .ai-summary-content {
    font-size: 1.2rem;
    padding: 2rem;
  }
  
  .disease-name-large {
    font-size: 1.4rem;
  }
  
  .recommendation-text {
    font-size: 1.1rem;
  }
  
  .health-score-title {
    font-size: 1.6rem;
  }
  
  .health-score-number {
    font-size: 2.5rem;
  }
  
  /* 保險特色標籤響應式 */
  .feature-tags-title {
    font-size: 1.2rem;
  }
  
  .feature-tags-container {
    gap: 0.75rem;
  }
  
  .feature-chip-prominent {
    font-size: 1rem !important;
    padding: 0.5rem 1rem !important;
    min-height: 44px !important;
  }
}

/* 保險風險分析彈窗樣式重構 - 統一主色調 */
.insurance-dialog-card {
  border-radius: 28px !important;
  overflow: hidden !important;
  box-shadow: 0 20px 80px rgba(0, 0, 0, 0.12) !important;
}

.insurance-dialog-header {
  background: linear-gradient(135deg, #00B8D9 0%, #0093A6 100%) !important;
  color: white !important;
  padding: 2rem !important;
  border: none !important;
}

.insurance-header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
}

.insurance-header-left {
  display: flex;
  align-items: center;
}

.insurance-avatar {
  background: rgba(255, 255, 255, 0.2) !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.15) !important;
}

.insurance-header-text {
  margin-left: 0;
}

.insurance-dialog-title {
  font-size: 1.8rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  letter-spacing: -0.5px;
}

.insurance-dialog-subtitle {
  font-size: 1.1rem;
  opacity: 0.9;
  font-weight: 400;
}

.insurance-close-btn {
  color: rgba(255, 255, 255, 0.9) !important;
  background: rgba(255, 255, 255, 0.15) !important;
  border-radius: 16px !important;
  transition: all 0.3s ease !important;
  width: 56px !important;
  height: 56px !important;
  min-width: 56px !important;
}

.insurance-close-btn:hover {
  background: rgba(255, 255, 255, 0.25) !important;
  transform: scale(1.1) rotate(90deg);
}

.insurance-dialog-content {
  padding: 3rem !important;
  max-height: 75vh;
  background: #fafafa !important;
}

/* 風險評分總覽重構 */
.risk-overview-wrapper {
  display: flex;
  align-items: center;
  gap: 3rem;
  background: white;
  padding: 3rem;
  border-radius: 24px;
  margin-bottom: 3rem;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

.risk-score-main {
  flex-shrink: 0;
}

.risk-score-circle {
  filter: drop-shadow(0 8px 16px rgba(0, 0, 0, 0.1));
}

.risk-score-content {
  display: flex;
  flex-direction: column;
  align-items: center;
}

.risk-score-number {
  font-size: 3.2rem;
  font-weight: 800;
  line-height: 1;
  color: inherit;
}

.risk-score-unit {
  font-size: 1.2rem;
  font-weight: 500;
  opacity: 0.8;
  margin-top: 0.25rem;
}

.risk-score-info {
  flex: 1;
}

.risk-score-title {
  font-size: 2rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 1.5rem;
  letter-spacing: -0.5px;
}

.risk-level-chip {
  font-size: 1.2rem !important;
  font-weight: 600 !important;
  padding: 1rem 1.5rem !important;
  height: auto !important;
  margin-bottom: 1.5rem;
}

.risk-score-description {
  font-size: 1.3rem;
  line-height: 1.6;
  color: #666;
  font-weight: 400;
  margin-bottom: 1.5rem;
}

.risk-metrics-summary {
  display: flex;
  gap: 2rem;
}

.summary-item {
  flex: 1;
  padding: 1.5rem;
  background: #f0faff;
  border-radius: 16px;
  border-left: 4px solid #00B8D9;
}

.summary-label {
  font-size: 1.1rem;
  color: #666;
  margin-bottom: 0.5rem;
  font-weight: 500;
}

.summary-value {
  font-size: 1.4rem;
  font-weight: 700;
}

.premium-adjustment {
  color: #00B8D9;
}

.underwriting-advice {
  color: #333;
}

/* 疾病風險分類重構 */
.risk-categories-wrapper {
  margin-bottom: 3rem;
}

.risk-category-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(480px, 1fr));
  gap: 2rem;
}

.category-card {
  background: white;
  padding: 2.5rem;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 2px solid transparent;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.category-card:hover {
  transform: translateY(-8px);
  border-color: #00B8D9;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
}

.category-header-large {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  gap: 1.5rem;
}

.category-info {
  flex: 1;
}

.category-name-large {
  font-size: 1.6rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 1rem;
  letter-spacing: -0.5px;
}

.category-risk-score {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.score-number {
  font-size: 2.2rem;
  font-weight: 800;
  color: #00B8D9;
  line-height: 1;
}

.score-label {
  font-size: 1rem;
  color: #888;
  font-weight: 500;
}

.category-level {
  flex-shrink: 0;
}

.category-progress-wrapper {
  margin: 1.5rem 0;
}

.category-progress-bar {
  border-radius: 7px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.category-details-large {
  margin-top: 2rem;
}

.category-description-large {
  font-size: 1.2rem;
  color: #555;
  line-height: 1.6;
  margin-bottom: 1.5rem;
  font-weight: 400;
}

.category-factors-section {
  margin-top: 1.5rem;
}

.factors-label-large {
  font-size: 1.3rem;
  font-weight: 600;
  color: #444;
  margin-bottom: 1rem;
}

.factors-chips-large {
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
}

.factor-chip-large {
  font-size: 1rem !important;
  font-weight: 600 !important;
  padding: 0.5rem 1rem !important;
  height: auto !important;
}

/* 健康指標權重重構 */
.health-metrics-wrapper {
  margin-bottom: 3rem;
}

.metrics-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 2rem;
}

.metric-card-large {
  background: white;
  padding: 2rem;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  transition: all 0.3s ease;
}

.metric-card-large:hover {
  transform: translateY(-4px);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
}

.metric-header-large {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 1.5rem;
  gap: 1rem;
}

.metric-info {
  flex: 1;
}

.metric-name-large {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
  margin-bottom: 0.5rem;
}

.metric-value-display {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
}

.metric-number {
  font-size: 1.8rem;
  font-weight: 700;
  color: #00B8D9;
}

.metric-unit-text {
  font-size: 1.1rem;
  color: #666;
  font-weight: 500;
}

.metric-status-icon {
  flex-shrink: 0;
}

.metric-weight-section {
  background: #f0faff;
  padding: 1.5rem;
  border-radius: 16px;
}

.weight-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.weight-label-large {
  font-size: 1.2rem;
  font-weight: 600;
  color: #444;
}

.weight-value-large {
  font-size: 1.4rem;
  font-weight: 700;
  color: #00B8D9;
}

.weight-progress {
  border-radius: 5px !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

/* 年齡風險趨勢重構 */
.age-risk-wrapper {
  margin-bottom: 3rem;
}

.age-risk-content {
  background: white;
  padding: 3rem;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  display: flex;
  gap: 3rem;
  align-items: center;
}

.current-age-display {
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  min-width: 160px;
  flex-shrink: 0;
}

.age-icon-wrapper {
  margin-bottom: 1rem;
}

.age-info {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.age-label-large {
  font-size: 1.2rem;
  color: #666;
  font-weight: 500;
}

.age-value-large {
  font-size: 2.5rem;
  font-weight: 800;
  color: #00B8D9;
  line-height: 1;
}

.age-projections-large {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.projection-card {
  background: #fff8f0;
  padding: 1.5rem;
  border-radius: 16px;
  border-left: 4px solid #FF9800;
}

.projection-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.projection-period-large {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
}

.projection-increase-large {
  font-size: 1.6rem;
  font-weight: 700;
  color: #FF9800;
}

.projection-progress {
  margin-bottom: 0.5rem;
  border-radius: 6px !important;
}

.projection-description {
  font-size: 1rem;
  color: #666;
  font-weight: 500;
}

/* 風險緩解策略重構 */
.mitigation-wrapper {
  margin-bottom: 3rem;
}

.mitigation-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
  gap: 2rem;
}

.mitigation-card-large {
  background: white;
  padding: 2.5rem;
  border-radius: 20px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border: 2px solid transparent;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
}

.mitigation-card-large:hover {
  transform: translateY(-8px);
  border-color: #4CAF50;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
}

.mitigation-header-large {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 2rem;
  gap: 1.5rem;
}

.mitigation-info {
  flex: 1;
}

.mitigation-category-large {
  font-size: 1.6rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 0.75rem;
  letter-spacing: -0.5px;
}

.mitigation-timeframe-large {
  display: flex;
  align-items: center;
  font-size: 1.1rem;
  color: #666;
  font-weight: 500;
}

.mitigation-reduction {
  flex-shrink: 0;
}

.mitigation-actions-large {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.mitigation-action-large {
  display: flex;
  align-items: center;
  padding: 1rem;
  background: #f0fff4;
  border-radius: 12px;
  border-left: 3px solid #4CAF50;
}

.action-text {
  font-size: 1.2rem;
  color: #333;
  font-weight: 500;
}

/* 專業核保建議重構 */
.professional-recommendations-wrapper {
  margin-bottom: 2rem;
}

.insurance-custom-tabs {
  background: #f8f9fa !important;
}

.insurance-recommendation-tab {
  font-size: 1.3rem !important;
  font-weight: 700 !important;
  padding: 2rem 2.5rem !important;
  min-height: 72px !important;
  border-radius: 16px 16px 0 0 !important;
  transition: all 0.3s ease !important;
}

.insurance-recommendation-tab:hover {
  background: rgba(0, 184, 217, 0.1) !important;
  transform: translateY(-2px);
}

.insurance-recommendations-window {
  padding: 2.5rem;
}

.insurance-recommendation-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.insurance-recommendation-item-large {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1.5rem;
  background: #f0faff;
  border-radius: 16px;
  border-left: 4px solid #00B8D9;
  transition: all 0.3s ease;
}

.insurance-recommendation-item-large:hover {
  background: #e6f7ff;
  transform: translateX(8px);
}

.insurance-recommendation-icon {
  flex-shrink: 0;
  margin-top: 0.25rem;
}

.insurance-recommendation-text {
  font-size: 1.2rem;
  line-height: 1.6;
  color: #555;
  font-weight: 400;
}

/* 區塊圖標顏色重構 - 保險風險分析專用 */
.section-icon-large.risk-category-icon {
  background: linear-gradient(135deg, #00B8D9 0%, #0093A6 100%);
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.3);
}

.section-icon-large.metrics-icon {
  background: linear-gradient(135deg, #00B8D9 10%, #FF6B6B 90%);
  box-shadow: 0 8px 24px rgba(255, 107, 107, 0.3);
}

.section-icon-large.age-icon {
  background: linear-gradient(135deg, #00B8D9 10%, #4ECDC4 90%);
  box-shadow: 0 8px 24px rgba(78, 205, 196, 0.3);
}

.section-icon-large.mitigation-icon {
  background: linear-gradient(135deg, #00B8D9 10%, #96CEB4 90%);
  box-shadow: 0 8px 24px rgba(150, 206, 180, 0.3);
}
</style> 