<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores';
import { healthCheckService, notifyError, notifySuccess } from '../services';
import { useRouter } from 'vue-router';
import { useUserStore } from '../stores/user';

// 假設有這些服務
// import { authorizeService, llmSummaryService } from '../services';

const authStore = useAuthStore();
const currentUser = computed(() => authStore.currentUser);
const userRole = computed(() => authStore.userRole);
const healthData = ref([]);
const loading = ref(false);

// 授權相關
const authTab = ref('requests');
const accessRequests = ref([]);
const authorizedTickets = ref([]);
const loadingRequests = ref(false);
const loadingTickets = ref(false);
const authProcessing = ref(false);

// LLM 分析相關
const llmLoading = ref(false);
const llmSummary = ref('');

// 詳細資料對話框
const detailDialog = ref(false);
const selectedReport = ref(null);
const reportMetrics = ref([]);

// 定義健康指標的正常值範圍和對應顏色
const healthRanges = {
  'Glu-AC': { min: 70, max: 100, unit: 'mg/dL', name: '空腹血糖' },
  'HbA1c': { min: 4, max: 5.6, unit: '%', name: '糖化血色素' },
  'LDL-C': { min: 0, max: 100, unit: 'mg/dL', name: '低密度脂蛋白膽固醇' },
  'HDL-C': { min: 40, max: 60, unit: 'mg/dL', name: '高密度脂蛋白膽固醇' },
  'BP-sys': { min: 90, max: 120, unit: 'mmHg', name: '收縮壓' },
  'BP-dia': { min: 60, max: 80, unit: 'mmHg', name: '舒張壓' },
  'T-CHO': { min: 120, max: 200, unit: 'mg/dL', name: '總膽固醇' },
  'TG': { min: 0, max: 150, unit: 'mg/dL', name: '三酸甘油脂' },
  'U.A': { min: 3, max: 7, unit: 'mg/dL', name: '尿酸' },
  'AST（GOT）': { min: 5, max: 40, unit: 'U/L', name: '麩草轉氨酶' },
  'ALT（GPT）': { min: 5, max: 35, unit: 'U/L', name: '麩丙轉氨酶' },
  'CRE': { min: 0.7, max: 1.3, unit: 'mg/dL', name: '肌酐酸' },
  'Hb': { min: 12, max: 16, unit: 'g/dL', name: '血紅蛋白' },
  'PLT': { min: 150, max: 400, unit: 'x10^3/uL', name: '血小板' },
  'WBC': { min: 4, max: 10, unit: 'x10^3/uL', name: '白血球' }
};

/*
//健檢報告測試假資料區塊
 healthData.value = [
   {
     id: 'RPT-TEST-001',
     content: 'Glu-AC: 92, HbA1c: 5.2...',
     date: '2024-05-01T10:00:00Z',
     rawData: {
       'Glu-AC': 92,
       'HbA1c': 5.2,
       'LDL-C': 98,
       'HDL-C': 55,
       'BP': '118/76',
       'T-CHO': 180,
       'TG': 110,
       'U.A': 5.8,
       'AST（GOT）': 28,
       'ALT（GPT）': 22,
       'CRE': 1.0,
       'Hb': 14.2,
       'PLT': 250,
       'WBC': 6.5,
       '備註': '一切正常'
     },
     originalReport: {}
   },
   {
     id: 'RPT-TEST-002',
     content: 'Glu-AC: 110, HbA1c: 6.1...',
     date: '2024-04-15T09:30:00Z',
     rawData: {
       'Glu-AC': 110,
       'HbA1c': 6.1,
       'LDL-C': 130,
       'HDL-C': 38,
       'BP': '135/88',
       'T-CHO': 210,
       'TG': 180,
       'U.A': 7.2,
       'AST（GOT）': 45,
       'ALT（GPT）': 40,
       'CRE': 1.4,
       'Hb': 11.8,
       'PLT': 180,
       'WBC': 11.2,
       '備註': '需追蹤血壓與血糖'
     },
     originalReport: {}
   }
 ];
*/

//健檢報告測試假資料區塊 end


onMounted(async () => {
  loading.value = true;
  try {
    // 從後端獲取健康數據 - 注意這裡的 API 對應到 HandleListMyReports
    const healthResponse = await healthCheckService.fetchUserHealthData();
    console.log('從後端獲取的健康數據:', healthResponse);
    
    // 處理來自後端的報告數據
    healthData.value = healthResponse.map(report => {
      // 嘗試解析 resultJson 字段 (如果是 JSON 字符串)
      let parsedResults = {};
      
      try {
        if (report.resultJson) {
          if (typeof report.resultJson === 'string') {
            parsedResults = JSON.parse(report.resultJson);
          } else if (typeof report.resultJson === 'object') {
            parsedResults = report.resultJson;
          }
        } else if (report.testResults) {
          if (typeof report.testResults === 'string') {
            parsedResults = JSON.parse(report.testResults);
          } else if (typeof report.testResults === 'object') {
            parsedResults = report.testResults;
          }
        } else if (report.test_results_json) {
          if (typeof report.test_results_json === 'string') {
            parsedResults = JSON.parse(report.test_results_json);
          } else if (typeof report.test_results_json === 'object') {
            parsedResults = report.test_results_json;
          }
        }
      } catch (e) {
        console.error('解析測試結果失敗:', e);
      }
      
      // 生成預覽內容
      const previewContent = Object.keys(parsedResults).length > 0 
        ? Object.keys(parsedResults).slice(0, 2).map(k => `${k}: ${parsedResults[k]}`).join(', ') + '...'
        : (report.content || '無資料').substring(0, 50);
      
      return {
        id: report.reportId || report.report_id || report.id || '未知',
        content: previewContent,
        date: report.createdAt ? parseInt(report.createdAt) : null,
        rawData: parsedResults,
        originalReport: report // 保存原始報告數據
      };
    });
    
    console.log('處理後的健康數據:', healthData.value);
    
    // 載入授權請求和已授權票據
    await Promise.all([
      loadAccessRequests(),
      loadGrantedTickets()
    ]);
  } catch (error) {
    console.error('獲取健康數據失敗:', error);
    notifyError(`獲取健康數據失敗：${error.message}`);
    healthData.value = [];
  } finally {
    loading.value = false;
  }
});


// 載入授權請求
const loadAccessRequests = async () => {
  loadingRequests.value = true;
  try {
    console.log('開始載入授權請求...');
    accessRequests.value = await healthCheckService.fetchAccessRequests();
    console.log('載入授權請求完成:', accessRequests.value);
    
    // 檢查獲取的數據是否完整
    if (accessRequests.value.length > 0) {
      accessRequests.value.forEach((request, index) => {
        if (!request.reportId || !request.reason || !request.requestTime) {
          console.warn(`授權請求 #${index} 資料不完整:`, request);
        }
      });
    }
  } catch (error) {
    console.error('載入授權請求失敗:', error);
    notifyError(`無法載入授權請求：${error.message || '未知錯誤'}`);
    accessRequests.value = []; // 確保失敗時清空列表
  } finally {
    loadingRequests.value = false;
  }
};

// 載入已授權票據
const loadGrantedTickets = async () => {
  loadingTickets.value = true;
  try {
    console.log('開始載入已授權票據...');
    const response = await healthCheckService.fetchGrantedTickets();
    console.log('載入已授權票據完成:', response);
    
    if (response && response.tickets && Array.isArray(response.tickets)) {
      authorizedTickets.value = response.tickets.map(ticket => ({
        id: ticket.reportId,
        requesterName: ticket.requesterName || '未知請求者',
        reportId: ticket.reportId,
        requestTime: parseInt(ticket.grantTime),
        expiry: parseInt(ticket.expiryTime),
        status: 'APPROVED',
        companyName: ticket.companyName || '未知公司'
      }));
    } else {
      authorizedTickets.value = [];
    }
    
    console.log('處理後的已授權票據:', authorizedTickets.value);
  } catch (error) {
    console.error('載入已授權票據失敗:', error);
    notifyError(`無法載入已授權票據：${error.message || '未知錯誤'}`);
    authorizedTickets.value = [];
  } finally {
    loadingTickets.value = false;
  }
};

// 同意授權請求
const approveRequest = async (requestId) => {
  authProcessing.value = true;
  try {
    console.log('開始處理同意授權請求:', requestId);
    const result = await healthCheckService.approveAccessRequest(requestId);
    console.log('授權結果:', result);
    
    if (result && result.success) {
      notifySuccess('授權請求已成功處理');
      // 重新載入授權請求和授權票據
      await Promise.all([
        loadAccessRequests(),
        loadGrantedTickets()
      ]);
    } else {
      throw new Error('未能成功處理授權請求');
    }
  } catch (error) {
    console.error('同意授權請求失敗:', error);
    notifyError(`授權處理失敗：${error.message || '未知錯誤'}`);
  } finally {
    authProcessing.value = false;
  }
};

// 拒絕授權請求
const rejectRequest = async (requestId) => {
  authProcessing.value = true;
  try {
    console.log('開始處理拒絕授權請求:', requestId);
    const result = await healthCheckService.rejectAccessRequest(requestId);
    console.log('拒絕結果:', result);
    
    if (result && result.success) {
      notifySuccess('已拒絕授權請求');
      // 重新載入授權請求
      await loadAccessRequests();
    } else {
      throw new Error('未能成功處理拒絕請求');
    }
  } catch (error) {
    console.error('拒絕授權請求失敗:', error);
    notifyError(`拒絕處理失敗：${error.message || '未知錯誤'}`);
  } finally {
    authProcessing.value = false;
  }
};

// 格式化日期顯示
function formatDate(dateString) {
  if (!dateString) return '未知日期';
  
  try {
    let timestamp;
    
    // 處理數字或字符串形式的時間戳
    if (typeof dateString === 'number' || !isNaN(Number(dateString))) {
      // 將字符串轉換為數字
      const numericTimestamp = Number(dateString);
      
      // 檢查是否為秒級時間戳
      timestamp = numericTimestamp < 10000000000 
        ? numericTimestamp * 1000  // 轉換秒為毫秒
        : numericTimestamp;        // 已經是毫秒
      
      return new Date(timestamp).toLocaleDateString('zh-TW', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
      });
    }
    
    // 如果不是有效的時間戳，嘗試作為一般日期字符串處理
    return new Date(dateString).toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
    });
  } catch (e) {
    console.error('日期格式化錯誤:', e, '輸入值:', dateString, '類型:', typeof dateString);
    return '日期格式錯誤';
  }
}

// 嘗試將內容解析為JSON對象
function parseReportContent(content) {
  if (!content) return {};
  
  try {
    if (typeof content === 'string') {
      return JSON.parse(content);
    }
    return content;
  } catch (e) {
    console.error('解析報告內容失敗:', e);
    return {};
  }
}

const router = useRouter();

// 處理查看詳細資料
function viewReportDetail(item) {
  console.log('查看報告詳情:', item);
  const report_id = item.report_id || item.id;
  
  if (!report_id) {
    notifyError('無法查看報告：缺少報告編號');
    return;
  }
  
  const reportData = healthData.value.find(report => report.id === report_id);
  const userStore = useUserStore();
  userStore.setCurrentReport(reportData);
  
  router.push({ 
    name: 'ReportDetail', 
    params: { 
      report_id,
      patient_id: currentUser.value
    }
  });
}

// 添加健康指標到視覺化列表
function addMetric(key, value) {
  const range = healthRanges[key];
  const percentage = calculatePercentage(value, range.min, range.max);
  const { color, status } = getStatusInfo(value, range.min, range.max);
  
  reportMetrics.value.push({
    key,
    name: range.name,
    value,
    unit: range.unit,
    percentage,
    color,
    status
  });
}

// 計算百分比位置
function calculatePercentage(value, min, max) {
  // 確保值在0-100範圍內
  const percentage = Math.min(Math.max(((value - min) / (max - min)) * 100, 0), 100);
  return Math.round(percentage);
}

// 獲取狀態信息
function getStatusInfo(value, min, max) {
  if (value < min) {
    // 數值越低，顏色越深
    const severity = Math.min(100, Math.max(0, (min - value) / min * 100));
    if (severity > 30) {
      return { color: 'deep-orange', status: '偏低' };
    } else {
      return { color: 'orange', status: '偏低' };
    }
  } else if (value > max) {
    // 數值越高，顏色越深
    const severity = Math.min(100, Math.max(0, (value - max) / max * 100));
    if (severity > 30) {
      return { color: 'red darken-4', status: '嚴重偏高' };
    } else {
      return { color: 'red', status: '偏高' };
    }
  } else {
    // 正常值 - 在範圍靠近邊界時顯示淺綠色
    const distToMid = Math.abs((value - (min + max) / 2) / ((max - min) / 2));
    if (distToMid > 0.7) {
      return { color: 'light-green', status: '正常' };
    } else {
      return { color: 'green', status: '正常' };
    }
  }
}

// 判斷是否在正常範圍之外
function isOutsideNormalRange(metric) {
  return metric.status === '偏高' || metric.status === '嚴重偏高' || metric.status === '偏低';
}

const handleLogout = () => {
  authStore.logout();
};

// 處理 LLM 分析
const handleLLMSummary = async () => {
  if (!healthData.value.length) {
    notifyError('無可分析的健康數據');
    return;
  }
  
  llmLoading.value = true;
  try {
    // 待後端 API 完成後實現
    // const summary = await healthCheckService.analyzeLLMSummary(healthData.value);
    
    // 暫時使用模擬 LLM 分析
    await new Promise(resolve => setTimeout(resolve, 1000));
    llmSummary.value = "【AI 健康摘要】根據您的健檢數據，血糖、血脂與血壓均在正常範圍內。建議維持均衡飲食和適度運動，每半年進行一次健康檢查。";
  } catch (error) {
    notifyError(`LLM 分析失敗：${error.message}`);
  } finally {
    llmLoading.value = false;
  }
};

// 獲取到期時間顯示的顏色
const getExpiryChipColor = (expiryTime) => {
  const now = Math.floor(Date.now() / 1000);
  const timeLeft = expiryTime - now;
  
  if (timeLeft <= 0) {
    return 'error';
  } else if (timeLeft <= 7 * 24 * 60 * 60) { // 7天內
    return 'warning';
  } else {
    return 'success';
  }
};

// 獲取到期時間文字顏色
const getExpiryTextColor = (expiryTime) => {
  const now = Math.floor(Date.now() / 1000);
  return expiryTime - now <= 0 ? 'white' : '';
};

// 狀態相關函數
const getStatusColor = (status) => {
  const colors = {
    'PENDING': 'warning',
    'APPROVED': 'success',
    'REJECTED': 'error'
  };
  return colors[status] || 'grey';
};

const getStatusText = (status) => {
  const texts = {
    'PENDING': '待處理',
    'APPROVED': '已同意',
    'REJECTED': '已拒絕'
  };
  return texts[status] || '未知';
};

// 新增狀態管理
const showPendingOnly = ref(false);

// 計算待處理請求數量
const pendingRequestsCount = computed(() => {
  return accessRequests.value.filter(req => req.status === 'PENDING').length;
});

// 根據篩選條件顯示請求列表
const filteredAccessRequests = computed(() => {
  if (showPendingOnly.value) {
    return accessRequests.value.filter(req => req.status === 'PENDING');
  }
  // 當顯示已授權時，返回已授權票據
  return authorizedTickets.value;
});

// 新增計算剩餘天數的函數
const getRemainingDays = (expiry) => {
  if (!expiry) return '-';
  const now = Math.floor(Date.now() / 1000);
  const days = Math.ceil((expiry - now) / (24 * 60 * 60));
  if (days < 0) return '已過期';
  if (days === 0) return '今日到期';
  return `剩餘 ${days} 天`;
};
</script>

<template>
  <div class="dashboard-bg">
    <v-container class="dashboard-container py-8 mx-auto" max-width="1400">
      <!-- 頂部統計卡片區 -->
      <v-row class="mb-6">
        <!-- 用戶資訊卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card" elevation="2">
            <v-card-text class="pa-8">
              <div class="d-flex align-center mb-4">
                <div class="rounded-circle bg-primary-lighten-5 p-3">
                  <v-icon size="28" color="primary">mdi-account</v-icon>
                </div>
                <v-btn
                  color="grey-darken-1"
                  @click="handleLogout"
                  variant="text"
                  icon
                  size="small"
                  class="ml-auto"
                >
                  <v-icon size="20">mdi-logout</v-icon>
                </v-btn>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-grey-darken-1">歡迎回來</div>
                <div class="text-h6 font-weight-bold">{{ currentUser }}</div>
                <v-chip
                  size="small"
                  color="success-lighten-5"
                  class="font-weight-medium px-2 mt-2"
                  variant="flat"
                >
                  <v-icon start size="16" color="success">mdi-check-circle</v-icon>
                  <span class="text-success">一般使用者</span>
                </v-chip>
              </div>
              <v-divider class="my-3"></v-divider>
              <div class="d-flex align-center justify-space-between text-caption text-grey">
                <span>
                  <v-icon size="14" class="mr-1">mdi-clock-outline</v-icon>
                  {{ new Date().toLocaleDateString('zh-TW') }}
                </span>
                <span>ID: {{ currentUser }}</span>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 健康報告統計卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card" elevation="2">
            <v-card-text class="pa-8">
              <div class="d-flex align-center mb-4">
                <div class="rounded-circle bg-success-lighten-5 p-3">
                  <v-icon size="28" color="success">mdi-file-document</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-grey-darken-1">健康報告總數</div>
                <div class="text-h4 font-weight-bold">{{ healthData.length }}</div>
                <div class="stat-progress mt-4">
                  <v-progress-linear
                    model-value="70"
                    color="success"
                    height="4"
                    rounded
                  ></v-progress-linear>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 待處理授權請求卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card" elevation="2">
            <v-card-text class="pa-8">
              <div class="d-flex align-center mb-4">
                <div class="rounded-circle bg-warning-lighten-5 p-3">
                  <v-icon size="28" color="warning">mdi-clock-outline</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-grey-darken-1">待處理授權</div>
                <div class="text-h4 font-weight-bold">{{ accessRequests.length }}</div>
                <div class="stat-progress mt-4">
                  <v-progress-linear
                    :model-value="(accessRequests.length / 10) * 100"
                    color="warning"
                    height="4"
                    rounded
                  ></v-progress-linear>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 已授權報告卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card" elevation="2">
            <v-card-text class="pa-8">
              <div class="d-flex align-center mb-4">
                <div class="rounded-circle bg-info-lighten-5 p-3">
                  <v-icon size="28" color="info">mdi-shield-check</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-grey-darken-1">已授權報告</div>
                <div class="text-h4 font-weight-bold">{{ authorizedTickets.length }}</div>
                <div class="stat-progress mt-4">
                  <v-progress-linear
                    :model-value="(authorizedTickets.length / healthData.length) * 100"
                    color="info"
                    height="4"
                    rounded
                  ></v-progress-linear>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 主要內容區 -->
      <v-row>
        <v-col cols="12">
          <!-- 健康檢查報告卡片 -->
          <v-card class="mb-6 elevation-2 rounded-xl">
            <v-card-title class="d-flex align-center py-4 px-6 bg-grey-lighten-4">
              <div class="d-flex align-center">
                <v-icon size="24" color="primary" class="me-3">mdi-file-document-outline</v-icon>
                <span class="text-h6 font-weight-bold">我的健康檢查報告</span>
              </div>
              <v-spacer></v-spacer>
              <div class="d-flex gap-2">
                <v-btn
                  color="primary"
                  variant="tonal"
                  size="small"
                  prepend-icon="mdi-filter-variant"
                >
                  篩選
                </v-btn>
                <v-btn
                  color="primary"
                  variant="tonal"
                  size="small"
                  prepend-icon="mdi-sort-variant"
                >
                  排序
                </v-btn>
              </div>
            </v-card-title>

            <v-card-text class="pa-6">
              <v-data-table
                :headers="[
                  { title: '報告編號', key: 'id', align: 'start', width: '150px' },
                  { title: '檢查日期', key: 'date', align: 'center', width: '150px' },
                  { title: '檢查類型', key: 'type', align: 'center', width: '150px' },
                  { title: '操作', key: 'actions', align: 'center', width: '120px' }
                ]"
                :items="healthData"
                :loading="loading"
                class="elevation-0"
                hover
                density="comfortable"
              >
                <!-- 報告編號欄位 -->
                <template v-slot:item.id="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="primary" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.id }}</span>
                  </div>
                </template>

                <!-- 檢查日期欄位 -->
                <template v-slot:item.date="{ item }">
                  <div class="d-flex align-center justify-center">
                    <v-icon size="16" color="grey" class="me-1">mdi-calendar</v-icon>
                    {{ formatDate(item.date) }}
                  </div>
                </template>

                <!-- 檢查類型欄位 -->
                <template v-slot:item.type="{ item }">
                  <v-chip
                    size="small"
                    color="primary"
                    variant="tonal"
                    class="font-weight-medium"
                  >
                    {{ item.type || '一般體檢' }}
                  </v-chip>
                </template>

                <!-- 操作欄位 -->
                <template v-slot:item.actions="{ item }">
                  <div class="d-flex gap-2 justify-center">
                    <v-btn
                      icon
                      variant="text"
                      size="small"
                      color="primary"
                      @click="viewReportDetail(item)"
                      class="view-report-btn"
                    >
                      <v-icon>mdi-eye</v-icon>
                      <v-tooltip
                        activator="parent"
                        location="top"
                        open-delay="200"
                      >
                        查看詳情
                      </v-tooltip>
                    </v-btn>
                    <v-btn
                      icon
                      variant="text"
                      size="small"
                      color="grey"
                      class="share-btn"
                    >
                      <v-icon>mdi-share-variant</v-icon>
                      <v-tooltip
                        activator="parent"
                        location="top"
                        open-delay="200"
                      >
                        分享
                      </v-tooltip>
                    </v-btn>
                  </div>
                </template>

                <!-- 無資料顯示 -->
                <template v-slot:no-data>
                  <div class="text-center pa-5">
                    <v-icon size="40" color="grey-lighten-1" class="mb-3">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium text-grey-darken-1">
                      尚無健康檢查報告
                    </div>
                    <div class="text-body-2 text-grey mt-2">
                      您目前沒有任何健康檢查報告記錄
                    </div>
                  </div>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>

          <!-- 授權管理卡片 -->
          <v-card elevation="2" class="rounded-xl">
            <v-card-title class="d-flex align-center py-4 px-6 bg-grey-lighten-4">
              <div class="d-flex align-center">
                <v-icon size="24" color="primary" class="me-3">mdi-key-chain</v-icon>
                <span class="text-h6 font-weight-bold">授權管理</span>
                <v-chip
                  v-if="pendingRequestsCount > 0"
                  color="warning"
                  size="small"
                  class="ms-3"
                >
                  {{ pendingRequestsCount }} 個待處理
                </v-chip>
              </div>
              <v-spacer></v-spacer>
              <v-btn-group density="comfortable" variant="tonal">
                <v-btn
                  :color="showPendingOnly ? 'warning' : 'grey'"
                  @click="showPendingOnly = true"
                >
                  待處理
                </v-btn>
                <v-btn
                  :color="!showPendingOnly ? 'primary' : 'grey'"
                  @click="showPendingOnly = false"
                >
                  已授權
                </v-btn>
              </v-btn-group>
            </v-card-title>

            <v-card-text class="pa-6">
              <v-data-table
                :headers="[
                  { 
                    title: '報告編號',
                    key: 'reportId',
                    align: 'start',
                    width: '120px'
                  },
                  { 
                    title: '請求者',
                    key: 'requesterName',
                    align: 'start',
                    width: '150px'
                  },
                  { 
                    title: '授權理由',
                    key: 'reason',
                    align: 'start'
                  },
                  { 
                    title: '申請日期',
                    key: 'requestTime',
                    align: 'center',
                    width: '120px'
                  },
                  { 
                    title: '到期日期',
                    key: 'expiry',
                    align: 'center',
                    width: '120px'
                  },
                  { 
                    title: '操作',
                    key: 'actions',
                    align: 'center',
                    width: '180px',
                    sortable: false
                  }
                ]"
                :items="filteredAccessRequests"
                :loading="showPendingOnly ? loadingRequests : loadingTickets"
                :loading-text="showPendingOnly ? '正在載入授權請求...' : '正在載入已授權報告...'"
                class="elevation-0"
                hover
                density="comfortable"
              >
                <!-- 報告編號欄位 -->
                <template v-slot:item.reportId="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="primary" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.reportId }}</span>
                  </div>
                </template>

                <!-- 請求者欄位 -->
                <template v-slot:item.requesterName="{ item }">
                  <div class="d-flex align-center">
                    <v-avatar size="28" color="primary" class="me-2">
                      <v-icon color="white" size="16">mdi-account</v-icon>
                    </v-avatar>
                    <div class="d-flex flex-column">
                      <span class="font-weight-medium">{{ item.requesterName }}</span>
                      <span class="text-caption text-grey">{{ item.companyName }}</span>
                    </div>
                  </div>
                </template>

                <!-- 授權理由欄位 -->
                <template v-slot:item.reason="{ item }">
                  <div class="text-body-2">{{ item.reason || '無' }}</div>
                </template>

                <!-- 申請日期欄位 -->
                <template v-slot:item.requestTime="{ item }">
                  <div class="d-flex align-center justify-center">
                    <v-icon size="16" color="grey" class="me-1">mdi-calendar</v-icon>
                    {{ formatDate(item.requestTime || item.grantTime) }}
                  </div>
                </template>

                <!-- 到期日期欄位 -->
                <template v-slot:item.expiry="{ item }">
                  <div v-if="item.expiry" class="d-flex align-center justify-center">
                    <v-chip
                      size="small"
                      :color="getExpiryChipColor(item.expiry)"
                      variant="tonal"
                      class="font-weight-medium"
                    >
                      {{ getRemainingDays(item.expiry) }}
                    </v-chip>
                  </div>
                  <span v-else>-</span>
                </template>

                <!-- 操作按鈕欄位 -->
                <template v-slot:item.actions="{ item }">
                  <div class="d-flex gap-2 justify-center">
                    <template v-if="item.status === 'PENDING'">
                      <v-btn
                        :loading="authProcessing"
                        @click="approveRequest(item.id)"
                        color="success"
                        variant="flat"
                        size="small"
                        class="modern-btn approve-btn"
                        elevation="0"
                      >
                        <template v-slot:prepend>
                          <v-icon size="18">mdi-check-circle-outline</v-icon>
                        </template>
                        授權
                      </v-btn>
                      <v-btn
                        :loading="authProcessing"
                        @click="rejectRequest(item.id)"
                        color="error"
                        variant="flat"
                        size="small"
                        class="modern-btn reject-btn"
                        elevation="0"
                      >
                        <template v-slot:prepend>
                          <v-icon size="18">mdi-close-circle-outline</v-icon>
                        </template>
                        拒絕
                      </v-btn>
                    </template>
                    <v-chip
                      v-else
                      :color="getStatusColor(item.status || 'APPROVED')"
                      size="small"
                      class="status-chip"
                      variant="flat"
                    >
                      <template v-slot:prepend>
                        <v-icon size="16" :color="getStatusColor(item.status || 'APPROVED')">
                          {{ item.status === 'REJECTED' ? 'mdi-close-circle' : 'mdi-check-circle' }}
                        </v-icon>
                      </template>
                      {{ getStatusText(item.status || 'APPROVED') }}
                    </v-chip>
                  </div>
                </template>

                <!-- 無資料顯示 -->
                <template v-slot:no-data>
                  <div class="text-center pa-5">
                    <v-icon size="40" color="grey-lighten-1" class="mb-3">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium text-grey-darken-1">
                      {{ showPendingOnly ? '無待處理請求' : '無已授權報告' }}
                    </div>
                    <div class="text-body-2 text-grey mt-2">
                      {{ showPendingOnly ? '目前沒有待處理的授權請求' : '目前沒有已授權的報告' }}
                    </div>
                  </div>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<style scoped>
/* 全局樣式 */
.dashboard-bg {
  background-color: #F8FAFC;
  min-height: 100vh;
  display: flex;
  align-items: flex-start;
  justify-content: center;
}

.dashboard-container {
  max-width: 1400px !important;
  margin: 0 auto;
  padding: 2rem;
  width: 100%;
}

/* 資訊卡片統一樣式 */
.info-card {
  height: 100%;
  min-height: 220px;
  display: flex;
  flex-direction: column;
  border-radius: 24px !important;
  background: white !important;
  transition: all 0.3s ease !important;
}

.info-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12) !important;
}

.info-card .v-card-text {
  height: 100%;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

/* 圓形圖標背景統一樣式 */
.rounded-circle {
  width: 48px;
  height: 48px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.rounded-circle:hover {
  transform: scale(1.05);
}

/* 進度條統一樣式 */
.stat-progress {
  margin-top: auto;
  padding-top: 1rem;
}

:deep(.v-progress-linear) {
  border-radius: 8px;
  opacity: 0.8;
}

/* 卡片內文字統一樣式 */
.text-h4 {
  font-size: 2.25rem !important;
  line-height: 2.75rem !important;
  font-weight: 700 !important;
}

.text-h6 {
  font-size: 1.25rem !important;
  line-height: 1.75rem !important;
  font-weight: 600 !important;
}

.text-subtitle-1 {
  font-size: 1rem !important;
  line-height: 1.5rem !important;
  opacity: 0.85;
}

.text-body-2 {
  font-size: 0.95rem !important;
  line-height: 1.5rem !important;
}

/* 表格樣式 */
:deep(.v-data-table-header th) {
  font-size: 1rem !important;
  color: #64748B !important;
  font-weight: 600 !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  padding: 14px 16px !important;
}

:deep(.v-data-table td) {
  padding: 14px 16px !important;
  color: #334155 !important;
  font-size: 1rem !important;
}

/* 狀態標籤樣式 */
:deep(.v-chip) {
  font-size: 0.9rem !important;
}

/* RWD 適配 */
@media (max-width: 960px) {
  .info-card {
    min-height: 200px;
  }
  
  .info-card .v-card-text {
    padding: 1.5rem !important;
  }
  
  .rounded-circle {
    width: 40px;
    height: 40px;
  }
  
  .text-h4 {
    font-size: 2rem !important;
    line-height: 2.5rem !important;
  }
  
  .text-h6 {
    font-size: 1.15rem !important;
    line-height: 1.6rem !important;
  }
}

@media (max-width: 600px) {
  .info-card {
    min-height: 180px;
  }
  
  .info-card .v-card-text {
    padding: 1rem !important;
  }
  
  .text-h4 {
    font-size: 1.75rem !important;
    line-height: 2.25rem !important;
  }
  
  .text-h6 {
    font-size: 1.1rem !important;
    line-height: 1.5rem !important;
  }
  
  :deep(.v-data-table td) {
    font-size: 0.95rem !important;
  }
}

/* 操作按鈕樣式 */
.view-report-btn,
.share-btn {
  transition: all 0.3s ease !important;
  width: 36px !important;
  height: 36px !important;
  border-radius: 8px !important;
}

.view-report-btn:hover,
.share-btn:hover {
  transform: translateY(-2px) !important;
  background-color: rgba(var(--v-theme-primary), 0.1) !important;
}

.view-report-btn:active,
.share-btn:active {
  transform: scale(0.95) !important;
}

/* 現代化按鈕基礎樣式 */
.modern-btn {
  min-width: 86px !important;
  height: 32px !important;
  border-radius: 8px !important;
  font-weight: 600 !important;
  letter-spacing: 0.3px !important;
  text-transform: none !important;
  transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1) !important;
  position: relative;
  overflow: hidden !important;
}

/* 授權按鈕特殊樣式 */
.approve-btn {
  background-color: rgb(var(--v-theme-success), 0.12) !important;
  color: rgb(var(--v-theme-success)) !important;
  border: 1px solid rgb(var(--v-theme-success), 0.1) !important;
}

.approve-btn:hover {
  background-color: rgb(var(--v-theme-success), 0.18) !important;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(var(--v-theme-success), 0.15) !important;
}

/* 拒絕按鈕特殊樣式 */
.reject-btn {
  background-color: rgb(var(--v-theme-error), 0.12) !important;
  color: rgb(var(--v-theme-error)) !important;
  border: 1px solid rgb(var(--v-theme-error), 0.1) !important;
}

.reject-btn:hover {
  background-color: rgb(var(--v-theme-error), 0.18) !important;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(var(--v-theme-error), 0.15) !important;
}

/* 按鈕點擊效果 */
.modern-btn:active {
  transform: scale(0.96) !important;
  box-shadow: none !important;
}

/* 狀態標籤樣式 */
.status-chip {
  min-width: 86px !important;
  height: 32px !important;
  font-weight: 600 !important;
  letter-spacing: 0.3px !important;
  border: 1px solid rgba(var(--v-theme-primary), 0.1) !important;
}

/* 按鈕圖標樣式 */
.modern-btn .v-icon {
  margin-right: 4px !important;
  transition: transform 0.2s ease !important;
}

.modern-btn:hover .v-icon {
  transform: scale(1.1) !important;
}

/* 載入狀態樣式 */
.modern-btn.v-btn--loading {
  opacity: 0.8;
}

/* RWD 適配 */
@media (max-width: 600px) {
  .modern-btn {
    min-width: 72px !important;
    height: 28px !important;
    font-size: 0.8rem !important;
  }
  
  .modern-btn .v-icon {
    font-size: 16px !important;
  }
  
  .status-chip {
    min-width: 72px !important;
    height: 28px !important;
    font-size: 0.8rem !important;
  }
}
</style>