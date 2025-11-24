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

import { ref, reactive, onMounted, computed } from 'vue';
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
const aiAnalysisLoading = ref(false);
const insurerAnalysisLoading = ref(false);

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

// AI 分析和保單推薦（適配新的資料結構）
const aiAnalysis = ref({
  summary: '',
  personalAnalysis: '',
  riskSummary: '',
  riskLevel: 'low', // low, medium, monitor, high
  diseaseRisks: [],
  protectionPlan: [],
  insuranceRecommendations: []
});
const aiAnalysisCache = reactive({
  user: null,
  insurer: null
});
const aiAnalysisUpdatedAt = reactive({
  user: null,
  insurer: null
});
const currentAnalysisType = computed(() => (userRole.value === 'insurer' ? 'insurer' : 'user'));

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
      console.log('保險業者 API 回應:', response);
      
      if (!response) {
        throw new Error('無法獲取報告數據');
      }
      
      // 安全地解析 JSON
      let parsedData = {};
      if (response.resultJson) {
        try {
          parsedData = JSON.parse(response.resultJson);
          console.log('保險業者解析後的數據:', parsedData);
        } catch (parseError) {
          console.error('JSON 解析失敗:', parseError);
          console.error('原始 resultJson:', response.resultJson);
          throw new Error('報告數據格式錯誤');
        }
      } else {
        console.warn('resultJson 為空或 undefined:', response.resultJson);
        throw new Error('未獲取到報告內容');
      }
      
      report.value = {
        id: reportId,
        patient_id: patientId,
        date: new Date().toISOString(),
        rawData: parsedData
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
          
          console.log('API 回應詳情:', response); // 添加更詳細的調試信息
          console.log('response.success:', response?.success);
          console.log('response.resultJson:', response?.resultJson);
          console.log('response.resultJson type:', typeof response?.resultJson);
          
          if (response && response.success && response.resultJson) {
            console.log('開始解析 resultJson:', response.resultJson);
            const parsedData = JSON.parse(response.resultJson);
            console.log('解析後的數據:', parsedData);
            
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
            console.error('API 回應驗證失敗:');
            console.error('response:', response);
            console.error('response.success:', response?.success);
            console.error('response.resultJson:', response?.resultJson);
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
const insuranceActiveTab = ref('immediate');

// 競爭對手分析表格標題
const competitorHeaders = [
  { title: '公司', key: 'company', align: 'start' },
  { title: '保費', key: 'premium', align: 'center' },
  { title: '保額', key: 'coverage', align: 'center' },
  { title: '風險調整', key: 'riskAdjustment', align: 'center' },
  { title: '競爭力評分', key: 'score', align: 'center' }
];

// 載入健康分析數據（含快取）
const loadHealthAnalysis = async (forceRefresh = false) => {
  if (!reportId) {
    console.error('報告 ID 不存在，無法載入分析');
    return;
  }

  const analysisType = currentAnalysisType.value;

  // 若已有快取且未指定強制更新，直接載入快取內容
  if (!forceRefresh && aiAnalysisCache[analysisType]) {
    aiAnalysis.value = JSON.parse(JSON.stringify(aiAnalysisCache[analysisType]));
    showAISummary.value = true;
    aiAnalysisLoading.value = false;
    return;
  }

  aiAnalysisLoading.value = true;
  showAISummary.value = true; // 顯示彈窗並呈現載入狀態

  try {
    console.log(`[loadHealthAnalysis] 開始載入分析: reportId=${reportId}`);
    
    const targetPatientId = patientId || userStore.user?.id || authStore.currentUser;
    if (analysisType === 'insurer' && !targetPatientId) {
      throw new Error('缺少病患身份，無法進行保險分析');
    }
    const response = await healthCheckService.getHealthAnalysis(reportId, analysisType, targetPatientId);
    
    console.log('[loadHealthAnalysis] API 回應:', response);
    
    if (response.success) {
      const analysisData = analysisType === 'user' 
        ? (response.user_analysis || response.userAnalysis || {})
        : (response.insurer_analysis || response.insurerAnalysis || {});
      
      if (analysisData) {
        if (analysisType === 'user') {
          const protectionPlanRaw = analysisData.protection_plan || analysisData.protectionPlan || [];
          const protectionPlan = Array.isArray(protectionPlanRaw)
            ? protectionPlanRaw
            : (protectionPlanRaw ? [protectionPlanRaw] : []);
          const insuranceRecommendationRaw = analysisData.insurance_recommendation || analysisData.insuranceRecommendation || [];
          const insuranceRecommendation = Array.isArray(insuranceRecommendationRaw)
            ? insuranceRecommendationRaw
            : (insuranceRecommendationRaw ? [insuranceRecommendationRaw] : []);
          const riskSummary = analysisData.risk_level_summary || analysisData.riskLevelSummary || analysisData.risk_level || analysisData.riskLevel || '';
          const riskLevelRaw = analysisData.risk_level || analysisData.riskLevel || riskSummary;
          aiAnalysis.value = {
            summary: analysisData.summary || '',
            personalAnalysis: analysisData.personal_analysis || analysisData.personalAnalysis || '',
            riskSummary,
            riskLevel: mapRiskLevel(riskLevelRaw),
            diseaseRisks: (analysisData.disease_risks || analysisData.diseaseRisks || []).map(risk => {
              const rawLevelText = (risk.risk_level || risk.riskLevel || risk.risk_level_summary || risk.riskLevelSummary || risk.risk_level_label || risk.riskLevelLabel || '').trim();
              const normalizedLevel = mapRiskLevel(rawLevelText);
              return {
              name: risk.disease || '',
                levelLabel: rawLevelText || '未提供風險等級',
                level: normalizedLevel,
                factors: Array.isArray(risk.main_factors || risk.mainFactors)
                  ? (risk.main_factors || risk.mainFactors)
                  : ((risk.main_factors || risk.mainFactors) ? [risk.main_factors || risk.mainFactors] : []),
              prevention: risk.advice || ''
              };
            }),
            protectionPlan,
            insuranceRecommendations: Array.isArray(insuranceRecommendation) ? insuranceRecommendation : []
          };
        } else {
          const riskSummary = analysisData.risk_level_label || analysisData.riskLevelLabel || '';
          const coreRecommendations = analysisData.core_recommendation || analysisData.coreRecommendation || [];
          aiAnalysis.value = {
            summary: analysisData.summary || '',
            personalAnalysis: '',
            riskSummary,
            riskLevel: mapRiskLevel(riskSummary),
            diseaseRisks: (analysisData.disease_risk_evaluation || analysisData.diseaseRiskEvaluation || []).map(risk => {
              const rawLevelText = (risk.risk_level || risk.riskLevel || risk.risk_level_summary || risk.riskLevelSummary || risk.risk_level_label || risk.riskLevelLabel || '').trim();
              const normalizedLevel = mapRiskLevel(rawLevelText);
              return {
              name: risk.disease || '',
                levelLabel: rawLevelText || '未提供風險等級',
                level: normalizedLevel,
                factors: Array.isArray(risk.main_factors || risk.mainFactors)
                  ? (risk.main_factors || risk.mainFactors)
                  : ((risk.main_factors || risk.mainFactors) ? [risk.main_factors || risk.mainFactors] : []),
              prevention: risk.advice || ''
              };
            }),
            protectionPlan: Array.isArray(coreRecommendations) ? coreRecommendations : [],
            insuranceRecommendations: []
          };
        }

        // 更新快取與最後更新時間
        aiAnalysisCache[analysisType] = JSON.parse(JSON.stringify(aiAnalysis.value));
        aiAnalysisUpdatedAt[analysisType] = new Date();
      }
    } else {
      throw new Error(response.message || '分析失敗');
    }
  } catch (error) {
    console.error('[loadHealthAnalysis] 載入分析失敗:', error);
    // 保持原有的假數據，不更新
  } finally {
    aiAnalysisLoading.value = false;
  }
};

// 載入保險業者專業風險評估
const loadInsurerRiskAnalysis = async (forceRefresh = false) => {
  if (!reportId) {
    console.error('報告 ID 不存在，無法載入風險評估');
    return;
  }

  if (!patientId) {
    console.error('病患 ID 不存在，無法載入風險評估');
    errorMsg.value = '缺少病患身份，無法進行風險評估';
    return;
  }

  // 若已有快取且未指定強制更新，直接載入快取內容
  if (!forceRefresh && aiAnalysisCache.insurer) {
    // 從快取的 aiAnalysis 轉換為 insuranceAnalysis
    updateInsuranceAnalysisFromCache();
    showRisk.value = true;
    insurerAnalysisLoading.value = false;
    return;
  }

  insurerAnalysisLoading.value = true;
  // 先顯示彈窗，但不顯示內容（會顯示載入狀態）
  if (!showRisk.value) {
    showRisk.value = true;
  }
  
  // 清空舊數據，避免顯示假資料
  insuranceAnalysis.value = {
    riskScore: 0,
    overallRiskLevel: 'low',
    riskCategories: {},
    healthMetrics: [],
    ageRiskFactors: {},
    riskMitigation: [],
    recommendations: []
  };

  try {
    console.log(`[loadInsurerRiskAnalysis] 開始載入風險評估: reportId=${reportId}, patientId=${patientId}`);
    
    const response = await healthCheckService.getHealthAnalysis(reportId, 'insurer', patientId);
    
    console.log('[loadInsurerRiskAnalysis] API 回應:', response);
    
    if (response.success) {
      const analysisData = response.insurer_analysis || response.insurerAnalysis || {};
      
      if (analysisData) {
        // 將後端數據映射到 insuranceAnalysis 結構（先更新，確保顯示真實數據）
        updateInsuranceAnalysisFromBackend(analysisData);
        
        // 更新 aiAnalysis（用於快取）
        const riskSummary = analysisData.risk_level_label || analysisData.riskLevelLabel || '';
        const coreRecommendations = analysisData.core_recommendation || analysisData.coreRecommendation || [];
        aiAnalysis.value = {
          summary: analysisData.summary || '',
          personalAnalysis: '',
          riskSummary,
          riskLevel: mapRiskLevel(riskSummary),
          diseaseRisks: (analysisData.disease_risk_evaluation || analysisData.diseaseRiskEvaluation || []).map(risk => {
            const rawLevelText = (risk.risk_level || risk.riskLevel || '').trim();
            const normalizedLevel = mapRiskLevel(rawLevelText);
            return {
              name: risk.disease || '',
              levelLabel: rawLevelText || '未提供風險等級',
              level: normalizedLevel,
              factors: Array.isArray(risk.main_factors || risk.mainFactors)
                ? (risk.main_factors || risk.mainFactors)
                : ((risk.main_factors || risk.mainFactors) ? [risk.main_factors || risk.mainFactors] : []),
              prevention: risk.advice || ''
            };
          }),
          protectionPlan: Array.isArray(coreRecommendations) ? coreRecommendations : [],
          insuranceRecommendations: []
        };

        // 更新快取
        aiAnalysisCache.insurer = JSON.parse(JSON.stringify(aiAnalysis.value));
        aiAnalysisUpdatedAt.insurer = new Date();
      } else {
        throw new Error('未獲取到分析數據');
      }
    } else {
      throw new Error(response.message || '風險評估失敗');
    }
  } catch (error) {
    console.error('[loadInsurerRiskAnalysis] 載入風險評估失敗:', error);
    errorMsg.value = error.message || '載入風險評估失敗，請稍後再試';
    // 保持原有的假數據，不更新
  } finally {
    insurerAnalysisLoading.value = false;
  }
};

// 從後端數據更新 insuranceAnalysis
const updateInsuranceAnalysisFromBackend = (analysisData) => {
  const riskLevelLabel = analysisData.risk_level_label || analysisData.riskLevelLabel || '';
  const overallRiskScore = analysisData.overall_risk_score || analysisData.overallRiskScore || 0;
  const diseaseRisks = analysisData.disease_risk_evaluation || analysisData.diseaseRiskEvaluation || [];
  const coreRecommendations = analysisData.core_recommendation || analysisData.coreRecommendation || [];
  
  // 根據風險等級標籤判斷風險等級
  let overallRiskLevel = 'low';
  let riskScore = overallRiskScore;
  
  if (riskLevelLabel.includes('高') || riskLevelLabel.toLowerCase().includes('high')) {
    overallRiskLevel = 'high';
    if (riskScore === 0) riskScore = 75; // 如果後端沒有提供分數，根據等級設定預設值
  } else if (riskLevelLabel.includes('中') || riskLevelLabel.toLowerCase().includes('medium')) {
    overallRiskLevel = 'medium';
    if (riskScore === 0) riskScore = 55;
  } else {
    overallRiskLevel = 'low';
    if (riskScore === 0) riskScore = 35;
  }

  // 將疾病風險評估轉換為 riskCategories
  const riskCategories = {};
  const categoryMapping = {
    '心血管': 'cardiovascular',
    '心臟': 'cardiovascular',
    '糖尿病': 'diabetes',
    '血糖': 'diabetes',
    '腎臟': 'kidney',
    '腎': 'kidney',
    '肝臟': 'liver',
    '肝': 'liver',
    '癌症': 'cancer',
    '腫瘤': 'cancer'
  };

  diseaseRisks.forEach(risk => {
    const diseaseName = risk.disease || '';
    let categoryKey = null;
    
    // 根據疾病名稱匹配分類
    for (const [chineseName, key] of Object.entries(categoryMapping)) {
      if (diseaseName.includes(chineseName)) {
        categoryKey = key;
        break;
      }
    }
    
    // 如果沒有匹配到，使用第一個單詞作為 key
    if (!categoryKey) {
      categoryKey = diseaseName.toLowerCase().replace(/\s+/g, '_');
    }

    const riskLevelText = risk.risk_level || risk.riskLevel || '';
    let level = 'low';
    let score = 30;
    
    if (riskLevelText.includes('高') || riskLevelText.toLowerCase().includes('high')) {
      level = 'high';
      score = 75;
    } else if (riskLevelText.includes('中') || riskLevelText.toLowerCase().includes('medium')) {
      level = 'medium';
      score = 55;
    }

    riskCategories[categoryKey] = {
      score,
      level,
      factors: Array.isArray(risk.main_factors || risk.mainFactors)
        ? (risk.main_factors || risk.mainFactors)
        : ((risk.main_factors || risk.mainFactors) ? [risk.main_factors || risk.mainFactors] : []),
      impact: level === 'high' ? 'high' : level === 'medium' ? 'medium' : 'low',
      description: risk.advice || '無詳細說明'
    };
  });

  // 更新 insuranceAnalysis（完全替換假資料）
  insuranceAnalysis.value = {
    riskScore,
    overallRiskLevel,
    riskCategories: riskCategories, // 只使用後端返回的數據
    healthMetrics: insuranceAnalysis.value.healthMetrics || [], // 保留健康指標（如果有的話）
    ageRiskFactors: insuranceAnalysis.value.ageRiskFactors || {}, // 保留年齡風險因素（如果有的話）
    riskMitigation: insuranceAnalysis.value.riskMitigation || [], // 保留風險緩解措施（如果有的話）
    recommendations: [
      {
        type: 'immediate',
        title: '立即建議',
        items: coreRecommendations.length > 0 
          ? coreRecommendations.slice(0, 4)
          : ['標準費率承保', '無需額外體檢', '建議定期追蹤', '整體健康狀況評估中']
      },
      {
        type: 'monitoring',
        title: '持續監控',
        items: coreRecommendations.length > 4
          ? coreRecommendations.slice(4)
          : ['每3個月追蹤健康指標', '年度健康檢查', '維持健康生活習慣', '定期風險評估']
      },
      {
        type: 'assessment',
        title: '風險評估',
        items: [
          `整體風險評級：${riskLevelLabel}`,
          overallRiskLevel === 'low' ? '可提供健康優惠費率' : overallRiskLevel === 'medium' ? '建議加費承保' : '需特殊核保',
          '適合標準保險產品',
          '長期風險展望需持續觀察'
        ]
      }
    ]
  };
};

// 從快取更新 insuranceAnalysis
const updateInsuranceAnalysisFromCache = () => {
  if (!aiAnalysisCache.insurer) return;
  
  const cached = aiAnalysisCache.insurer;
  const riskSummary = cached.riskSummary || '';
  
  let overallRiskLevel = 'low';
  let riskScore = 35;
  
  if (riskSummary.includes('高') || riskSummary.toLowerCase().includes('high')) {
    overallRiskLevel = 'high';
    riskScore = 75;
  } else if (riskSummary.includes('中') || riskSummary.toLowerCase().includes('medium')) {
    overallRiskLevel = 'medium';
    riskScore = 55;
  }

  // 將 diseaseRisks 轉換為 riskCategories
  const riskCategories = {};
  const categoryMapping = {
    '心血管': 'cardiovascular',
    '心臟': 'cardiovascular',
    '糖尿病': 'diabetes',
    '血糖': 'diabetes',
    '腎臟': 'kidney',
    '腎': 'kidney',
    '肝臟': 'liver',
    '肝': 'liver',
    '癌症': 'cancer',
    '腫瘤': 'cancer'
  };

  (cached.diseaseRisks || []).forEach(risk => {
    const diseaseName = risk.name || '';
    let categoryKey = null;
    
    for (const [chineseName, key] of Object.entries(categoryMapping)) {
      if (diseaseName.includes(chineseName)) {
        categoryKey = key;
        break;
      }
    }
    
    if (!categoryKey) {
      categoryKey = diseaseName.toLowerCase().replace(/\s+/g, '_');
    }

    riskCategories[categoryKey] = {
      score: risk.level === 'high' ? 75 : risk.level === 'medium' ? 55 : 30,
      level: risk.level || 'low',
      factors: risk.factors || [],
      impact: risk.level === 'high' ? 'high' : risk.level === 'medium' ? 'medium' : 'low',
      description: risk.prevention || '無詳細說明'
    };
  });

  insuranceAnalysis.value = {
    ...insuranceAnalysis.value,
    riskScore,
    overallRiskLevel,
    riskCategories: {
      ...insuranceAnalysis.value.riskCategories,
      ...riskCategories
    }
  };
};

// 輔助函數：根據健康分數獲取風險等級
const mapRiskLevel = (level) => {
  if (!level) return 'medium';
  const levelLower = (level || '').toLowerCase();
  // 處理 "需注意"、"邊緣偏高" 等特殊情況
  if (levelLower.includes('需注意') || levelLower.includes('monitor') || levelLower.includes('邊緣')) return 'monitor';
  if (levelLower.includes('低') || levelLower.includes('low')) return 'low';
  if (levelLower.includes('中') || levelLower.includes('medium') || levelLower.includes('中風險') || levelLower.includes('中低')) return 'medium';
  if (levelLower.includes('高') || levelLower.includes('high')) return 'high';
  return 'medium'; // 默認
};

const formatDateTime = (value) => {
  if (!value) return '';
  const date = value instanceof Date ? value : new Date(value);
  if (Number.isNaN(date.getTime())) return '';
  return date.toLocaleString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  });
};

const getDiseaseCardClass = (level) => {
  switch (level) {
    case 'low':
      return 'disease-card--low';
    case 'monitor':
      return 'disease-card--monitor';
    case 'high':
      return 'disease-card--high';
    default:
      return 'disease-card--medium';
  }
};

const getRiskLabelClass = (level) => {
  switch (level) {
    case 'low':
      return 'risk-label--low';
    case 'monitor':
      return 'risk-label--monitor';
    case 'high':
      return 'risk-label--high';
    default:
      return 'risk-label--medium';
  }
};

const getRiskIcon = (level) => {
  switch (level) {
    case 'low':
      return 'mdi-shield-check';
    case 'monitor':
      return 'mdi-eye-alert-outline';
    case 'high':
      return 'mdi-shield-alert';
    default:
      return 'mdi-alert-circle';
  }
};

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
            :loading="aiAnalysisLoading"
            @click="loadHealthAnalysis"
          >
            <v-icon start size="20">mdi-robot-outline</v-icon>
            AI 智能分析
          </v-btn>
          
          <v-btn
            v-if="userRole === 'insurer'"
            class="action-btn risk-btn"
            elevation="0"
            color="#00B8D9"
            :loading="insurerAnalysisLoading"
            @click="loadInsurerRiskAnalysis"
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
            <div class="ai-header-right">
              <div
                v-if="aiAnalysisUpdatedAt[currentAnalysisType] && !aiAnalysisLoading"
                class="ai-last-updated"
              >
                上次更新：{{ formatDateTime(aiAnalysisUpdatedAt[currentAnalysisType]) }}
            </div>
              <v-btn
                class="ai-refresh-btn"
                variant="flat"
                color="white"
                :disabled="aiAnalysisLoading"
                :loading="aiAnalysisLoading"
                @click.stop="loadHealthAnalysis(true)"
              >
                <v-icon start size="18">mdi-refresh</v-icon>
                重新分析
              </v-btn>
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
          </div>
        </v-card-title>

        <v-card-text class="ai-dialog-content">
          <!-- 載入狀態 -->
          <div v-if="aiAnalysisLoading" class="loading-wrapper">
            <v-progress-circular
              indeterminate
              color="#00B8D9"
              size="64"
              width="6"
            />
            <div class="loading-text">正在分析您的健康數據，這可能需要幾分鐘時間，請耐心等待...</div>
            <div class="loading-subtext">AI 正在仔細分析您的健康指標，請勿關閉此視窗</div>
          </div>

          <!-- 分析內容 -->
          <template v-else>
            <section
              v-if="aiAnalysis.summary || aiAnalysis.riskSummary"
              class="analysis-summary-card"
            >
              <div class="analysis-summary-card__icon">
                <v-icon size="24" color="#0f766e">mdi-heart-pulse</v-icon>
                </div>
              <div class="analysis-summary-card__content">
                <div class="analysis-summary-card__title">健康總結</div>
                <p v-if="aiAnalysis.summary" class="analysis-summary-card__text">
                  {{ aiAnalysis.summary }}
                </p>
                <div
                  v-if="aiAnalysis.riskSummary"
                  class="analysis-summary-card__badge"
                >
                  <v-icon size="16" color="#006c7d">mdi-alert-decagram-outline</v-icon>
                  <span>{{ aiAnalysis.riskSummary }}</span>
              </div>
            </div>
            </section>

            <div class="analysis-grid">
              <div class="analysis-column analysis-column--main">
                <section class="analysis-card">
                  <header class="analysis-card__header">
                    <div class="analysis-card__icon analysis-card__icon--brain">
                      <v-icon size="22" color="#4338ca">mdi-brain</v-icon>
              </div>
                    <div class="analysis-card__titles">
                      <div class="analysis-card__title">AI 專業分析</div>
                      <div class="analysis-card__subtitle">
                        深度解析您的檢查結果，提供可理解的專業洞察。
            </div>
                </div>
                  </header>
                  <div class="analysis-card__body">
                    <p v-if="aiAnalysis.personalAnalysis" class="analysis-text">
                      {{ aiAnalysis.personalAnalysis }}
                    </p>
                    <div v-else class="analysis-empty">
                      暫無詳細分析內容，請稍後再試。
                </div>
              </div>
                </section>

                <section
                  v-if="aiAnalysis.diseaseRisks && aiAnalysis.diseaseRisks.length"
                  class="analysis-card"
                >
                  <header class="analysis-card__header">
                    <div class="analysis-card__icon analysis-card__icon--risk">
                      <v-icon size="22" color="#b91c1c">mdi-shield-alert</v-icon>
              </div>
                    <div class="analysis-card__titles">
                      <div class="analysis-card__title">疾病風險評估</div>
                      <div class="analysis-card__subtitle">
                        彙整重要疾病的風險程度與主要觀察指標。
            </div>
                    </div>
                  </header>
                  <div class="risk-list">
              <div
                v-for="(disease, index) in aiAnalysis.diseaseRisks"
                :key="index"
                      class="risk-item"
                    >
                      <div class="risk-item__header">
                        <div class="risk-item__title">{{ disease.name }}</div>
                        <span
                          class="risk-item__badge"
                          :class="'risk-item__badge--' + (disease.level || 'medium')"
                        >
                          {{ disease.levelLabel }}
                        </span>
                  </div>
                      <div
                        v-if="disease.factors && disease.factors.length"
                        class="risk-item__chips"
                      >
                        <span
                        v-for="factor in disease.factors"
                        :key="factor"
                          class="analysis-chip"
                      >
                        {{ factor }}
                        </span>
                    </div>
                      <p v-if="disease.prevention" class="risk-item__text">
                        {{ disease.prevention }}
                      </p>
                  </div>
                  </div>
                </section>

                <section
                  v-if="aiAnalysis.healthTrends && aiAnalysis.healthTrends.length"
                  class="analysis-card"
                >
                  <header class="analysis-card__header">
                    <div class="analysis-card__icon analysis-card__icon--trend">
                      <v-icon size="22" color="#0ea5e9">mdi-trending-up</v-icon>
              </div>
                    <div class="analysis-card__titles">
                      <div class="analysis-card__title">健康趨勢分析</div>
                      <div class="analysis-card__subtitle">
                        檢視近期指標變化，掌握改善與需關注的方向。
            </div>
                    </div>
                  </header>
                  <div class="trend-list">
              <div
                v-for="trend in aiAnalysis.healthTrends"
                :key="trend.metric"
                      class="trend-item"
                    >
                      <div
                        class="trend-item__icon"
                        :class="{
                          'trend-item__icon--good': trend.trend === 'improving',
                          'trend-item__icon--warn': trend.trend === 'concern'
                        }"
                      >
                        <v-icon size="20" color="white">
                          {{ trend.trend === 'improving'
                            ? 'mdi-trending-up'
                            : trend.trend === 'concern'
                              ? 'mdi-trending-down'
                              : 'mdi-trending-neutral' }}
                  </v-icon>
                </div>
                      <div class="trend-item__content">
                        <div class="trend-item__title">{{ trend.metric }}</div>
                        <div class="trend-item__meta">
                          <span>
                            {{ trend.trend === 'improving'
                              ? '持續改善'
                              : trend.trend === 'concern'
                                ? '需要關注'
                                : '保持穩定' }}
                    </span>
                          <span
                            v-if="typeof trend.change === 'number' && trend.change !== 0"
                            class="trend-item__change"
                            :class="{ 'trend-item__change--up': trend.change > 0 }"
                          >
                      {{ trend.change > 0 ? '+' : '' }}{{ trend.change }}%
                    </span>
                  </div>
                </div>
              </div>
            </div>
                </section>
          </div>

              <div class="analysis-column analysis-column--side">
                <section
                  v-if="aiAnalysis.protectionPlan && aiAnalysis.protectionPlan.length"
                  class="analysis-card"
                >
                  <header class="analysis-card__header">
                    <div class="analysis-card__icon analysis-card__icon--plan">
                      <v-icon size="22" color="#0f766e">mdi-clipboard-check-outline</v-icon>
                      </div>
                    <div class="analysis-card__titles">
                      <div class="analysis-card__title">健康行動計畫</div>
                      <div class="analysis-card__subtitle">
                        依據分析提供可立即執行的步驟。
                    </div>
                  </div>
                  </header>
                  <ul class="action-list">
                    <li
                      v-for="(item, index) in aiAnalysis.protectionPlan"
                :key="index"
                      class="action-item"
                    >
                      <v-icon size="18" color="#0f766e">mdi-check-circle</v-icon>
                      <span>{{ item }}</span>
                    </li>
                  </ul>
                </section>

                <section
                  v-if="aiAnalysis.insuranceRecommendations && aiAnalysis.insuranceRecommendations.length"
                  class="analysis-card"
                >
                  <header class="analysis-card__header">
                    <div class="analysis-card__icon analysis-card__icon--insurance">
                      <v-icon size="22" color="#ea580c">mdi-shield-check</v-icon>
                  </div>
                    <div class="analysis-card__titles">
                      <div class="analysis-card__title">保險建議</div>
                      <div class="analysis-card__subtitle">
                        對應風險的保障建議與挑選重點。
                </div>
              </div>
                  </header>
                  <ul class="action-list action-list--highlight">
                    <li
                      v-for="(item, index) in aiAnalysis.insuranceRecommendations"
                      :key="index"
                      class="action-item"
                    >
                      <v-icon size="18" color="#f59e0b">mdi-lightbulb-on</v-icon>
                      <span>{{ item }}</span>
                    </li>
                  </ul>
                </section>
            </div>
          </div>
          </template>
        </v-card-text>
      </v-card>
    </v-dialog>

    <!-- 保險風險分析彈窗 -->
    <v-dialog v-model="showRisk" max-width="1200" scrollable>
      <v-card class="ai-dialog-card">
        <!-- 對話框標題 -->
        <v-card-title class="ai-dialog-header">
          <div class="ai-header-content">
            <div class="ai-header-left">
              <v-avatar class="ai-avatar mr-4" size="56" color="#00B8D9">
                <v-icon color="white" size="28">mdi-shield-check</v-icon>
              </v-avatar>
              <div class="ai-header-text">
                <div class="ai-dialog-title">專業保險風險評估</div>
                <div class="ai-dialog-subtitle">基於健康數據的精準風險分析與核保建議</div>
              </div>
            </div>
            <div class="ai-header-right">
              <div
                v-if="aiAnalysisUpdatedAt.insurer && !insurerAnalysisLoading"
                class="ai-last-updated"
              >
                上次更新：{{ formatDateTime(aiAnalysisUpdatedAt.insurer) }}
              </div>
              <v-btn
                class="ai-refresh-btn"
                variant="flat"
                color="white"
                :disabled="insurerAnalysisLoading"
                :loading="insurerAnalysisLoading"
                @click.stop="loadInsurerRiskAnalysis(true)"
              >
                <v-icon start size="18">mdi-refresh</v-icon>
                重新分析
              </v-btn>
              <v-btn
                icon
                variant="text"
                @click="showRisk = false"
                class="ai-close-btn"
                size="large"
              >
                <v-icon size="24">mdi-close</v-icon>
              </v-btn>
            </div>
          </div>
        </v-card-title>

        <v-card-text class="ai-dialog-content">
          <!-- 載入狀態 -->
          <div v-if="insurerAnalysisLoading" class="loading-wrapper">
            <v-progress-circular
              indeterminate
              color="#00B8D9"
              size="64"
              width="6"
            />
            <div class="loading-text">正在進行專業風險評估分析，這可能需要幾分鐘時間，請耐心等待...</div>
            <div class="loading-subtext">AI 正在仔細分析健康指標與風險因素，請勿關閉此視窗</div>
          </div>

          <!-- 分析內容 -->
          <template v-else>
            <section
              v-if="aiAnalysis.summary || aiAnalysis.riskSummary"
              class="analysis-summary-card"
            >
              <div class="analysis-summary-card__icon">
                <v-icon size="24" color="#0f766e">mdi-shield-check</v-icon>
              </div>
              <div class="analysis-summary-card__content">
                <div class="analysis-summary-card__title">風險評估總結</div>
                <p v-if="aiAnalysis.summary" class="analysis-summary-card__text">
                  {{ aiAnalysis.summary }}
                </p>
                <div
                  v-if="aiAnalysis.riskSummary"
                  class="analysis-summary-card__badge"
                >
                  <v-icon size="16" color="#006c7d">mdi-alert-decagram-outline</v-icon>
                  <span>{{ aiAnalysis.riskSummary }}</span>
                </div>
              </div>
            </section>

            <div class="analysis-grid">
              <div class="analysis-column analysis-column--main">
                <section
                  v-if="aiAnalysis.diseaseRisks && aiAnalysis.diseaseRisks.length"
                  class="analysis-card"
                >
                  <header class="analysis-card__header">
                    <div class="analysis-card__icon analysis-card__icon--risk">
                      <v-icon size="22" color="#b91c1c">mdi-shield-alert</v-icon>
                    </div>
                    <div class="analysis-card__titles">
                      <div class="analysis-card__title">疾病風險評估</div>
                      <div class="analysis-card__subtitle">
                        彙整重要疾病的風險程度與主要觀察指標。
                      </div>
                    </div>
                  </header>
                  <div class="risk-list">
                    <div
                      v-for="(disease, index) in aiAnalysis.diseaseRisks"
                      :key="index"
                      class="risk-item"
                    >
                      <div class="risk-item__header">
                        <div class="risk-item__title">{{ disease.name }}</div>
                        <span
                          class="risk-item__badge"
                          :class="'risk-item__badge--' + (disease.level || 'medium')"
                        >
                          {{ disease.levelLabel }}
                        </span>
                      </div>
                      <div
                        v-if="disease.factors && disease.factors.length"
                        class="risk-item__chips"
                      >
                        <span
                          v-for="factor in disease.factors"
                          :key="factor"
                          class="analysis-chip"
                        >
                          {{ factor }}
                        </span>
                      </div>
                      <p v-if="disease.prevention" class="risk-item__text">
                        {{ disease.prevention }}
                      </p>
                    </div>
                  </div>
                </section>
                <div v-else class="analysis-empty">
                  暫無疾病風險評估數據，請稍後再試。
                </div>
              </div>

              <div class="analysis-column analysis-column--side">
                <section
                  v-if="aiAnalysis.protectionPlan && aiAnalysis.protectionPlan.length"
                  class="analysis-card"
                >
                  <header class="analysis-card__header">
                    <div class="analysis-card__icon analysis-card__icon--plan">
                      <v-icon size="22" color="#0f766e">mdi-clipboard-check-outline</v-icon>
                    </div>
                    <div class="analysis-card__titles">
                      <div class="analysis-card__title">專業核保建議</div>
                      <div class="analysis-card__subtitle">
                        依據風險分析提供核保建議與監控重點。
                      </div>
                    </div>
                  </header>
                  <ul class="action-list">
                    <li
                      v-for="(item, index) in aiAnalysis.protectionPlan"
                      :key="index"
                      class="action-item"
                    >
                      <v-icon size="18" color="#0f766e">mdi-check-circle</v-icon>
                      <span>{{ item }}</span>
                    </li>
                  </ul>
                </section>
              </div>
            </div>
          </template>
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

.ai-header-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.ai-last-updated {
  color: rgba(255, 255, 255, 0.85);
  font-size: 0.95rem;
  font-weight: 500;
}

.ai-refresh-btn {
  border-radius: 16px !important;
  background: rgba(255, 255, 255, 0.2) !important;
  color: white !important;
  font-weight: 600 !important;
  letter-spacing: 0.5px;
  text-transform: none !important;
}

.ai-refresh-btn:hover {
  background: rgba(255, 255, 255, 0.3) !important;
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

.loading-wrapper {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 4rem 2rem;
  gap: 2rem;
}

.loading-text {
  font-size: 1.3rem;
  font-weight: 600;
  color: #333;
  text-align: center;
}

.loading-subtext {
  font-size: 1rem;
  color: #666;
  text-align: center;
  margin-top: 0.5rem;
}

/* 健康總結區域 */
.health-overview-wrapper {
  display: flex;
  align-items: flex-start;
  gap: 1.5rem;
  background: linear-gradient(135deg, rgba(0, 184, 217, 0.18) 0%, rgba(0, 147, 166, 0.08) 100%);
  border-radius: 24px;
  padding: 2.25rem 2.8rem;
  margin-bottom: 2.5rem;
  border: 1px solid rgba(0, 184, 217, 0.2);
  box-shadow: 0 12px 36px rgba(0, 147, 166, 0.18);
}

.overview-icon {
  width: 64px;
  height: 64px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #0093A6 0%, #00B8D9 100%);
  box-shadow: 0 10px 28px rgba(0, 147, 166, 0.35);
}

.overview-content {
  display: flex;
  flex-direction: column;
  gap: 1rem;
  color: #0f172a;
}

.overview-title {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.6px;
}

.overview-text {
  font-size: 1.25rem;
  line-height: 1.8;
  color: #1f2937;
  font-weight: 500;
}

.overview-subtext {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.05rem;
  line-height: 1.6;
  padding: 0.75rem 1.25rem;
  border-radius: 16px;
  background: rgba(0, 147, 166, 0.12);
  color: #006c7d;
  font-weight: 600;
  width: fit-content;
}

/* 大型區塊標題樣式 - 統一主色調 */
.ai-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.section-header-large {
  display: flex;
  align-items: flex-start;
  gap: 1.25rem;
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

.section-title-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  flex: 1;
}

.section-title-large {
  font-size: 2.05rem;
  font-weight: 700;
  color: #333;
  letter-spacing: -0.5px;
}

.section-subtitle {
  font-size: 1rem;
  color: #5f6b7a;
  line-height: 1.55;
}

/* AI 分析新版卡片樣式 */
.analysis-summary-card {
  display: flex;
  gap: 1.5rem;
  background: linear-gradient(135deg, rgba(0, 184, 217, 0.16) 0%, rgba(0, 147, 166, 0.05) 100%);
  border-radius: 24px;
  padding: 2.4rem 2.8rem;
  border: 1px solid rgba(0, 184, 217, 0.18);
  box-shadow: 0 12px 34px rgba(0, 147, 166, 0.18);
  margin-bottom: 2.5rem;
}

.analysis-summary-card__icon {
  min-width: 56px;
  width: 56px;
  height: 56px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #d1fae5;
  box-shadow: inset 0 0 0 1px rgba(15, 118, 110, 0.12);
}

.analysis-summary-card__content {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
  color: #0f172a;
}

.analysis-summary-card__title {
  font-size: 1.75rem;
  font-weight: 700;
  letter-spacing: -0.5px;
}

.analysis-summary-card__text {
  font-size: 1.15rem;
  line-height: 1.8;
  color: #1f2937;
  font-weight: 500;
}

.analysis-summary-card__badge {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 1.1rem;
  background: rgba(0, 147, 166, 0.12);
  border-radius: 999px;
  font-size: 0.95rem;
  font-weight: 600;
  color: #006c7d;
  width: fit-content;
}

.analysis-grid {
  display: grid;
  grid-template-columns: minmax(0, 1.75fr) minmax(0, 1fr);
  gap: 1.75rem;
}

.analysis-column {
  display: flex;
  flex-direction: column;
  gap: 1.75rem;
}

.analysis-card {
  background: #ffffff;
  border-radius: 20px;
  padding: 2rem 2.3rem;
  box-shadow: 0 10px 32px rgba(15, 23, 42, 0.08);
  border: 1px solid rgba(15, 23, 42, 0.05);
  display: flex;
  flex-direction: column;
  gap: 1.4rem;
}

.analysis-card__header {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
}

.analysis-card__icon {
  width: 48px;
  height: 48px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f1f5f9;
  box-shadow: inset 0 0 0 1px rgba(148, 163, 184, 0.2);
}

.analysis-card__icon--brain {
  background: #e0e7ff;
  box-shadow: inset 0 0 0 1px rgba(99, 102, 241, 0.25);
}

.analysis-card__icon--risk {
  background: #fee2e2;
  box-shadow: inset 0 0 0 1px rgba(239, 68, 68, 0.25);
}

.analysis-card__icon--trend {
  background: #e0f2fe;
  box-shadow: inset 0 0 0 1px rgba(14, 165, 233, 0.25);
}

.analysis-card__icon--plan {
  background: #dcfce7;
  box-shadow: inset 0 0 0 1px rgba(34, 197, 94, 0.25);
}

.analysis-card__icon--insurance {
  background: #ffedd5;
  box-shadow: inset 0 0 0 1px rgba(249, 115, 22, 0.25);
}

.analysis-card__titles {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.analysis-card__title {
  font-size: 1.45rem;
  font-weight: 700;
  color: #111827;
  letter-spacing: -0.3px;
}

.analysis-card__subtitle {
  font-size: 0.95rem;
  color: #6b7280;
  line-height: 1.5;
}

.analysis-card__body {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.analysis-text {
  font-size: 1.1rem;
  line-height: 1.8;
  color: #374151;
}

.analysis-empty {
  padding: 1.2rem;
  border: 1px dashed rgba(148, 163, 184, 0.7);
  border-radius: 16px;
  text-align: center;
  color: #94a3b8;
  font-size: 0.95rem;
}

.analysis-chip {
  display: inline-flex;
  align-items: center;
  padding: 0.35rem 0.75rem;
  background: #e8f4ff;
  color: #0c4a6e;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 600;
}

.risk-list {
  display: flex;
  flex-direction: column;
  gap: 1.2rem;
}

.risk-item {
  border: 1px solid rgba(15, 23, 42, 0.06);
  border-radius: 16px;
  background: #f8fafc;
  padding: 1.25rem 1.5rem;
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.risk-item__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.risk-item__title {
  font-weight: 600;
  font-size: 1.2rem;
  color: #111827;
}

.risk-item__badge {
  padding: 0.25rem 0.9rem;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 600;
  color: #fff;
}

.risk-item__badge--low {
  background: #22c55e;
}

.risk-item__badge--medium {
  background: #f59e0b;
}

.risk-item__badge--monitor {
  background: #3b82f6;
}

.risk-item__badge--high {
  background: #f97316;
}

.risk-item__chips {
  display: flex;
  flex-wrap: wrap;
  gap: 0.6rem;
}

.risk-item__text {
  color: #4b5563;
  line-height: 1.6;
  font-size: 0.98rem;
}

.trend-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.trend-item {
  display: flex;
  gap: 1rem;
  align-items: center;
  padding: 1rem 1.25rem;
  border-radius: 16px;
  background: #f8fafc;
  border: 1px solid rgba(15, 23, 42, 0.05);
}

.trend-item__icon {
  width: 40px;
  height: 40px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #94a3b8;
}

.trend-item__icon--good {
  background: #34d399;
}

.trend-item__icon--warn {
  background: #f97316;
}

.trend-item__content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.trend-item__title {
  font-weight: 600;
  color: #111827;
  font-size: 1.05rem;
}

.trend-item__meta {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  font-size: 0.9rem;
  color: #64748b;
}

.trend-item__change {
  font-weight: 600;
  color: #0ea5e9;
}

.trend-item__change--up {
  color: #f97316;
}

.action-list {
  display: flex;
  flex-direction: column;
  gap: 0.9rem;
}

.action-item {
  display: flex;
  align-items: flex-start;
  gap: 0.8rem;
  padding: 0.85rem 1.1rem;
  background: #f0fdf4;
  border: 1px solid rgba(74, 222, 128, 0.25);
  border-radius: 14px;
  color: #0f172a;
  font-size: 0.98rem;
  line-height: 1.5;
}

.action-list--highlight .action-item {
  background: #fff7ed;
  border-color: rgba(251, 191, 36, 0.35);
}



/* AI 分析摘要重構 */
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

.summary-section {
  margin-bottom: 1.5rem;
}

.summary-label {
  font-size: 1.2rem;
  font-weight: 600;
  color: #00B8D9;
  margin-bottom: 0.75rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.summary-label::before {
  content: '';
  width: 4px;
  height: 20px;
  background: #00B8D9;
  border-radius: 2px;
}

.summary-text {
  font-size: 1.3rem;
  line-height: 1.8;
  color: #555;
  font-weight: 400;
}

.summary-empty {
  font-size: 1.15rem;
  color: #9ca3af;
  padding: 1.5rem;
  background: #f4f6f8;
  border-radius: 16px;
  text-align: center;
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

.disease-card--low {
  border-color: rgba(76, 175, 80, 0.25);
  background: linear-gradient(135deg, rgba(76, 175, 80, 0.08) 0%, rgba(76, 175, 80, 0.02) 100%);
}

.disease-card--medium {
  border-color: rgba(255, 152, 0, 0.25);
  background: linear-gradient(135deg, rgba(255, 152, 0, 0.08) 0%, rgba(255, 152, 0, 0.02) 100%);
}

.disease-card--monitor {
  border-color: rgba(33, 150, 243, 0.25);
  background: linear-gradient(135deg, rgba(33, 150, 243, 0.1) 0%, rgba(33, 150, 243, 0.03) 100%);
}

.disease-card--high {
  border-color: rgba(244, 67, 54, 0.25);
  background: linear-gradient(135deg, rgba(244, 67, 54, 0.1) 0%, rgba(244, 67, 54, 0.03) 100%);
}

.disease-card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  gap: 1.5rem;
}

.disease-name-group {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.disease-name-large {
  font-size: 1.6rem;
  font-weight: 700;
  color: #333;
  margin-bottom: 1rem;
  letter-spacing: -0.5px;
}

.risk-label {
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
  border-radius: 999px;
  padding: 0.45rem 1.35rem;
  font-weight: 600;
  font-size: 1rem;
  letter-spacing: 0.02em;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
  transition: all 0.3s ease;
}

.risk-label .v-icon {
  color: inherit !important;
}

.risk-label--low {
  background: rgba(76, 175, 80, 0.18);
  color: #2e7d32;
}

.risk-label--medium {
  background: rgba(255, 152, 0, 0.2);
  color: #e65100;
}

.risk-label--monitor {
  background: rgba(33, 150, 243, 0.2);
  color: #1565c0;
}

.risk-label--high {
  background: rgba(244, 67, 54, 0.22);
  color: #b71c1c;
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


.plan-icon {
  background: linear-gradient(135deg, #34d399 0%, #059669 100%);
}

.plan-list {
  margin-top: 1.5rem;
  padding: 2rem 2.5rem;
  background: white;
  border-radius: 20px;
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.plan-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.25rem;
  border-radius: 16px;
  background: #f0fdf4;
  border: 1px solid rgba(52, 211, 153, 0.25);
  transition: all 0.3s ease;
  color: #14532d;
  font-size: 1.05rem;
  line-height: 1.5;
}

.plan-item:hover {
  transform: translateX(6px);
  box-shadow: 0 12px 24px rgba(34, 197, 94, 0.15);
}

.plan-text {
  flex: 1;
}

.insurance-icon {
  background: linear-gradient(135deg, #fbbf24 0%, #f97316 100%);
}

.insurance-list {
  margin-top: 1.5rem;
  padding: 2rem 2.5rem;
  background: white;
  border-radius: 20px;
  box-shadow: 0 12px 36px rgba(0, 0, 0, 0.08);
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.insurance-item {
  display: flex;
  align-items: flex-start;
  gap: 1rem;
  padding: 1rem 1.25rem;
  border-radius: 16px;
  background: #fff7eb;
  border: 1px solid rgba(253, 186, 116, 0.35);
  color: #92400e;
  font-size: 1.05rem;
  line-height: 1.6;
  transition: all 0.3s ease;
}

.insurance-item:hover {
  transform: translateX(6px);
  box-shadow: 0 10px 22px rgba(249, 115, 22, 0.15);
}

.insurance-text {
  flex: 1;
}

.insurance-item .v-icon {
  margin-top: 0.2rem;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .ai-dialog-content {
    padding: 1.5rem !important;
  }

  .ai-header-content {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }

  .ai-header-right {
    width: 100%;
    justify-content: flex-end;
    flex-wrap: wrap;
    gap: 0.75rem;
  }
  
  .health-overview-wrapper {
    flex-direction: column;
    align-items: flex-start;
    text-align: left;
    padding: 1.75rem;
    gap: 1rem;
  }

  .overview-icon {
    width: 56px;
    height: 56px;
  }

  .overview-text {
    font-size: 1.1rem;
  }
  
  .analysis-summary-card {
    flex-direction: column;
    padding: 1.8rem;
  }

  .analysis-grid {
    grid-template-columns: 1fr;
  }

  .analysis-column--side {
    position: static;
  }

  .analysis-card {
    padding: 1.6rem 1.8rem;
  }

  .analysis-card__icon {
    width: 44px;
    height: 44px;
  }

  .analysis-card__title {
    font-size: 1.35rem;
  }

  .analysis-summary-card__icon {
    width: 48px;
    height: 48px;
  }

  .analysis-summary-card__title {
    font-size: 1.5rem;
  }

  .analysis-summary-card__text {
    font-size: 1.05rem;
  }

  .trend-item {
    align-items: flex-start;
  }

  .trend-item__meta {
    flex-wrap: wrap;
  }
  
  .section-header-large {
    flex-direction: column;
    gap: 1rem;
  }

  .section-title-group {
    width: 100%;
    text-align: left;
  }
  
  .section-title-large {
    font-size: 1.8rem;
  }
  
  .disease-risk-cards {
    grid-template-columns: 1fr;
  }
  
  .plan-list,
  .insurance-list {
    padding: 1.5rem 1.75rem;
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
  margin-bottom: 0;
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