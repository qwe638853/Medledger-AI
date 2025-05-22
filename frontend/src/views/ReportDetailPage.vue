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
import healthCheckService from '../services/healthCheckService';

const route = useRoute();
const router = useRouter();
const reportId = route.params.report_id;
const patientId = route.params.patient_id;
const loading = ref(true);
const report = ref({});
const metrics = ref({});
const errorMsg = ref('');

onMounted(async () => {
  loading.value = true;
  try {
    const response = await healthCheckService.fetchReportContent(reportId, patientId);
    if (response && response.resultJson) {
      try {
        metrics.value = JSON.parse(response.resultJson);
      } catch (e) {
        metrics.value = {};
      }
    } else if (typeof response === 'object') {
      metrics.value = response;
    }
    report.value = response;
  } catch (e) {
    errorMsg.value = '獲取報告內容失敗';
  } finally {
    loading.value = false;
  }
});

const numericMetrics = computed(() => {
  const result = {};
  for (const [key, value] of Object.entries(metrics.value)) {
    const match = (value || '').toString().match(/-?\d+(\.\d+)?/);
    if (match) result[key] = value;
  }
  return result;
});

const textMetrics = computed(() => {
  const result = {};
  for (const [key, value] of Object.entries(metrics.value)) {
    const match = (value || '').toString().match(/-?\d+(\.\d+)?/);
    if (!match) result[key] = value;
  }
  return result;
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
</script>

<template>
  <v-container class="py-8">
    <v-btn @click="router.back()" prepend-icon="mdi-arrow-left" class="mb-4">返回</v-btn>
    <v-card v-if="loading" class="pa-8 text-center"><v-progress-circular indeterminate color="primary" /></v-card>
    <v-alert v-else-if="errorMsg" type="error">{{ errorMsg }}</v-alert>
    <template v-else>
      <v-card class="mb-6 pa-6">
        <div class="text-h5 font-weight-bold mb-2">健康檢查報告詳情</div>
        <div class="mb-2">報告編號：{{ reportId }}</div>
        <div class="mb-2">病患 ID：{{ patientId }}</div>
        <div class="mb-2">日期：{{ report.value?.date || '-' }}</div>
      </v-card>

      <!-- 指標視覺化區塊 -->
      <div class="text-h6 font-weight-bold mb-4 mt-8">主要健康指標視覺化</div>
      <v-row>
        <v-col
          v-for="(value, key) in numericMetrics"
          :key="key"
          cols="12"
          sm="6"
          md="4"
        >
          <v-card outlined class="mb-4 metric-visual-card">
            <v-card-title class="pb-0">{{ METRIC_NAME_MAP[key] || key }}</v-card-title>
            <v-card-text class="d-flex flex-column align-center justify-center">
              <v-progress-circular
                :value="getMetricPercent(key, value)"
                :color="getMetricColor(key, value)"
                :size="100"
                :width="12"
                class="mb-2"
              >
                <span class="text-h6 font-weight-bold">{{ getMetricNumber(value) }}</span>
              </v-progress-circular>
              <div class="mt-2 text-caption grey--text">參考值：
                <span v-if="METRIC_REF_RANGE[key]">
                  <template v-if="METRIC_REF_RANGE[key].min !== undefined">{{ METRIC_REF_RANGE[key].min }}</template>
                  <template v-if="METRIC_REF_RANGE[key].max !== undefined"> - {{ METRIC_REF_RANGE[key].max }}</template>
                  <template v-if="METRIC_REF_RANGE[key].unit"> {{ METRIC_REF_RANGE[key].unit }}</template>
                </span>
                <span v-else>--</span>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 其他指標（非數值型） -->
      <div class="text-h6 font-weight-bold mb-4 mt-8">其他檢查項目</div>
      <v-row>
        <v-col v-for="(value, key) in textMetrics" :key="key + '-text'" cols="12" sm="6" md="4">
          <v-card outlined class="mb-2">
            <v-card-text>
              <span class="font-weight-bold">{{ METRIC_NAME_MAP[key] || key }}：</span>
              <span>{{ value }}</span>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </template>
  </v-container>
</template>

<style scoped>
.metric-visual-card {
  min-height: 220px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
}
</style> 