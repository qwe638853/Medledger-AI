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

// 動態生成表格頭部，已授權清單不顯示授權理由
const tableHeaders = computed(() => {
  if (showPendingOnly.value) {
    // 待處理清單，包含授權理由
    return [
      { title: '報告編號', key: 'reportId', align: 'start', width: '180px' },
      { title: '請求者', key: 'requesterName', align: 'start', width: '220px' },
      { title: '授權理由', key: 'reason', align: 'start', width: '250px' },
      { title: '申請日期', key: 'requestTime', align: 'center', width: '160px' },
      { title: '到期日期', key: 'expiry', align: 'center', width: '160px' },
      { title: '操作', key: 'actions', align: 'center', width: '240px', sortable: false }
    ];
  } else {
    // 已授權清單，不包含授權理由
    return [
      { title: '報告編號', key: 'reportId', align: 'start', width: '180px' },
      { title: '請求者', key: 'requesterName', align: 'start', width: '220px' },
      { title: '申請日期', key: 'requestTime', align: 'center', width: '160px' },
      { title: '到期日期', key: 'expiry', align: 'center', width: '160px' },
      { title: '操作', key: 'actions', align: 'center', width: '240px', sortable: false }
    ];
  }
});
</script>

<template>
  <div class="dashboard-bg">
    <v-container class="dashboard-container py-8 mx-auto" max-width="1400">
      <!-- 頂部統計卡片區 -->
      <v-row class="mb-8">
        <!-- 用戶資訊卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card user-info-card" elevation="0">
            <v-card-text class="pa-0">
              <div class="user-card-wrapper">
                <!-- 背景裝飾 -->
                <div class="user-card-bg-decoration"></div>
                
                <!-- 用戶資訊主體 -->
                <div class="user-card-main">
                  <div class="user-avatar-section">
                    <div class="user-avatar-container">
                      <v-icon class="user-avatar-icon">mdi-account-circle</v-icon>
                    </div>
                    <div class="user-status-indicator"></div>
                  </div>
                  
                  <div class="user-info-content">
                    <div class="user-name">{{ currentUser }}</div>
                    <div class="user-role">
                      <v-icon size="14" class="me-1">mdi-account</v-icon>
                      一般使用者
                    </div>
                  </div>
                </div>
                
                <!-- 登出按鈕 -->
                <v-btn
                  class="user-logout-btn"
                  @click="handleLogout"
                  icon
                  size="small"
                  variant="text"
                >
                  <v-icon>mdi-logout</v-icon>
                  <v-tooltip activator="parent" location="bottom">
                    登出
                  </v-tooltip>
                </v-btn>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 健康報告統計卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card stat-card" elevation="0">
            <v-card-text class="pa-0">
              <div class="card-content">
                <div class="icon-container">
                  <v-icon>mdi-file-document-multiple</v-icon>
                </div>
                
                <div class="card-info">
                  <div class="stat-number">{{ healthData.length }}</div>
                  <div class="stat-label">健康報告</div>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 待處理授權請求卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card stat-card" elevation="0">
            <v-card-text class="pa-0">
              <div class="card-content">
                <div class="icon-container">
                  <v-icon>mdi-clock-alert</v-icon>
                </div>
                
                <div class="card-info">
                  <div class="stat-number">{{ pendingRequestsCount }}</div>
                  <div class="stat-label">待處理授權</div>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 已授權報告卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card stat-card" elevation="0">
            <v-card-text class="pa-0">
              <div class="card-content">
                <div class="icon-container">
                  <v-icon>mdi-shield-check</v-icon>
                </div>
                
                <div class="card-info">
                  <div class="stat-number">{{ authorizedTickets.length }}</div>
                  <div class="stat-label">已授權報告</div>
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
          <v-card class="main-card health-report-card mb-8" elevation="0">
            <div class="card-header">
              <div class="card-title">
                <div class="card-title-icon">
                  <v-icon>mdi-file-document-outline</v-icon>
                </div>
                <div>
                  <div class="card-title-text">我的健康檢查報告</div>
                </div>
              </div>
            </div>

            <div class="table-container">
              <v-data-table
                :headers="[
                  { title: '報告編號', key: 'id', align: 'start', width: '200px' },
                  { title: '檢查日期', key: 'date', align: 'center', width: '180px' },
                  { title: '檢查類型', key: 'type', align: 'center', width: '180px' },
                  { title: '操作', key: 'actions', align: 'center', width: '150px', sortable: false }
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
                    <v-icon size="18" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.id }}</span>
                  </div>
                </template>
                
                <!-- 檢查日期欄位 -->
                <template v-slot:item.date="{ item }">
                  <div class="d-flex align-center justify-center">
                    <v-icon size="16" class="me-1">mdi-calendar</v-icon>
                    {{ formatDate(item.date) }}
                  </div>
                </template>
                
                <!-- 檢查類型欄位 -->
                <template v-slot:item.type="{ item }">
                  <v-chip
                    class="type-chip"
                    size="large"
                  >
                    <v-icon start size="16">mdi-medical-bag</v-icon>
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
                      @click="viewReportDetail(item)"
                      class="view-report-btn"
                    >
                      <v-icon>mdi-eye</v-icon>
                      <v-tooltip activator="parent" location="top">
                        查看詳情
                      </v-tooltip>
                    </v-btn>
                    <v-btn
                      icon
                      variant="text"
                      size="small"
                      class="share-btn"
                    >
                      <v-icon>mdi-share-variant</v-icon>
                      <v-tooltip activator="parent" location="top">
                        分享
                      </v-tooltip>
                    </v-btn>
                  </div>
                </template>

                <!-- 無資料顯示 -->
                <template v-slot:no-data>
                  <div class="text-center pa-8">
                    <v-icon size="64" color="grey-lighten-2" class="mb-4">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium mb-2">
                      尚無健康檢查報告
                    </div>
                    <div class="text-body-2">
                      您目前沒有任何健康檢查報告記錄
                    </div>
                  </div>
                </template>
              </v-data-table>
            </div>
          </v-card>

          <!-- 授權管理卡片 -->
          <v-card class="main-card auth-management-card" elevation="0">
            <div class="card-header">
              <div class="card-title">
                <div class="card-title-icon">
                  <v-icon>mdi-key-chain</v-icon>
                </div>
                <div>
                  <div class="card-title-text">授權管理</div>
                </div>
                <v-chip
                  v-if="pendingRequestsCount > 0"
                  class="pending-badge"
                  size="small"
                >
                  {{ pendingRequestsCount }} 個待處理
                </v-chip>
                
                <div class="auth-tabs">
                  <v-btn
                    :class="['auth-tab-btn', { active: showPendingOnly }]"
                    @click="showPendingOnly = true"
                    variant="text"
                    size="large"
                  >
                    <v-icon start size="20">mdi-clock-outline</v-icon>
                    待處理
                  </v-btn>
                  <v-btn
                    :class="['auth-tab-btn', { active: !showPendingOnly }]"
                    @click="showPendingOnly = false"
                    variant="text"
                    size="large"
                  >
                    <v-icon start size="20">mdi-check-circle</v-icon>
                    已授權
                  </v-btn>
                </div>
              </div>
            </div>

            <div class="table-container">
              <v-data-table
                :headers="tableHeaders"
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
                    <v-icon size="18" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.reportId }}</span>
                  </div>
                </template>

                <!-- 請求者欄位 -->
                <template v-slot:item.requesterName="{ item }">
                  <div class="d-flex align-center">
                    <v-avatar size="32" class="me-3">
                      <v-icon color="white" size="16">mdi-account</v-icon>
                    </v-avatar>
                    <div class="d-flex flex-column">
                      <span class="font-weight-medium">{{ item.requesterName }}</span>
                      <span class="text-caption text-grey">{{ item.companyName }}</span>
                    </div>
                  </div>
                </template>

                <!-- 授權理由欄位 -->
                <template v-slot:item.reason="{ item }" v-if="showPendingOnly">
                  <div class="text-body-2">{{ item.reason || '無' }}</div>
                </template>

                <!-- 申請日期欄位 -->
                <template v-slot:item.requestTime="{ item }">
                  <div class="d-flex align-center justify-center">
                    <v-icon size="16" class="me-1">mdi-calendar</v-icon>
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
                      class="expiry-chip"
                    >
                      <v-icon start size="14">mdi-clock-outline</v-icon>
                      {{ getRemainingDays(item.expiry) }}
                    </v-chip>
                  </div>
                  <span v-else class="text-grey">-</span>
                </template>

                <!-- 操作按鈕欄位 -->
                <template v-slot:item.actions="{ item }">
                  <div class="d-flex gap-2 justify-center">
                    <template v-if="item.status === 'PENDING'">
                      <v-btn
                        :loading="authProcessing"
                        @click="approveRequest(item.id)"
                        class="modern-btn approve-btn"
                        size="small"
                        elevation="0"
                      >
                        <template v-slot:prepend>
                          <v-icon size="16">mdi-check-circle-outline</v-icon>
                        </template>
                        授權
                      </v-btn>
                      <v-btn
                        :loading="authProcessing"
                        @click="rejectRequest(item.id)"
                        class="modern-btn reject-btn"
                        size="small"
                        elevation="0"
                      >
                        <template v-slot:prepend>
                          <v-icon size="16">mdi-close-circle-outline</v-icon>
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
                        <v-icon size="14">
                          {{ item.status === 'REJECTED' ? 'mdi-close-circle' : 'mdi-check-circle' }}
                        </v-icon>
                      </template>
                      {{ getStatusText(item.status || 'APPROVED') }}
                    </v-chip>
                  </div>
                </template>

                <!-- 無資料顯示 -->
                <template v-slot:no-data>
                  <div class="text-center pa-8">
                    <v-icon size="64" color="grey-lighten-2" class="mb-4">
                      {{ showPendingOnly ? 'mdi-clock-outline' : 'mdi-shield-check' }}
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium mb-2">
                      {{ showPendingOnly ? '無待處理請求' : '無已授權報告' }}
                    </div>
                    <div class="text-body-2">
                      {{ showPendingOnly ? '目前沒有待處理的授權請求' : '目前沒有已授權的報告' }}
                    </div>
                  </div>
                </template>
              </v-data-table>
            </div>
          </v-card>
        </v-col>
      </v-row>
    </v-container>
  </div>
</template>

<style scoped>
/* 全局樣式 */
.dashboard-bg {
  background: #F8FAFC;
  min-height: 100vh;
  padding: 2rem 0;
}

.dashboard-container {
  max-width: 1400px !important;
  margin: 0 auto;
  padding: 2rem !important;
  width: 100%;
}

/* 統計卡片區域樣式 */
.info-card {
  height: 100%;
  min-height: 160px;
  border-radius: 24px !important;
  background: white !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  border: 1px solid rgba(0, 184, 217, 0.1) !important;
  box-shadow: 0 8px 32px rgba(0, 184, 217, 0.08) !important;
  margin-bottom: 2rem;
  overflow: hidden;
  position: relative;
}

.info-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 16px 48px rgba(0, 184, 217, 0.15) !important;
  border-color: rgba(0, 184, 217, 0.2) !important;
}

.info-card .v-card-text {
  height: 100%;
  padding: 0 !important;
}

/* 卡片內容布局 */
.card-content {
  height: 100%;
  min-height: 160px;
  padding: 2rem !important;
  display: flex;
  align-items: center;
  gap: 1.5rem;
  position: relative;
}

/* 圖標容器 */
.icon-container {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 184, 217, 0.1) !important;
  border: 2px solid rgba(0, 184, 217, 0.2) !important;
  flex-shrink: 0;
}

.icon-container .v-icon {
  font-size: 36px !important;
  color: #00B8D9 !important;
}

/* 卡片資訊區域 */
.card-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  text-align: center;
}

/* 用戶卡片特殊樣式 */
.user-card .card-title {
  font-size: 1.5rem !important;
  font-weight: 700 !important;
  color: #22292F !important;
  line-height: 1.2;
}

.user-card .card-subtitle {
  font-size: 1.1rem !important;
  font-weight: 500 !important;
  color: #8898AA !important;
}

/* 統計卡片樣式 */
.stat-card .stat-number {
  font-size: 3rem !important;
  font-weight: 800 !important;
  line-height: 1 !important;
  color: #22292F !important;
  margin-bottom: 0.25rem;
}

.stat-card .stat-label {
  font-size: 1.2rem !important;
  font-weight: 600 !important;
  color: #8898AA !important;
  line-height: 1.2;
}

/* 登出按鈕 */
.logout-btn {
  position: absolute !important;
  top: 1.25rem !important;
  right: 1.25rem !important;
  width: 32px !important;
  height: 32px !important;
  border-radius: 12px !important;
  background: rgba(239, 68, 68, 0.1) !important;
  color: #ef4444 !important;
  transition: all 0.3s ease !important;
}

.logout-btn:hover {
  background: #ef4444 !important;
  color: white !important;
  transform: scale(1.05) !important;
}

.logout-btn .v-icon {
  font-size: 16px !important;
}

/* 主要內容卡片樣式 */
.main-card {
  border-radius: 24px !important;
  background: white !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08) !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
  overflow: hidden !important;
  margin-bottom: 2rem !important;
}

.card-header {
  background: white !important;
  border-bottom: 1px solid #e2e8f0 !important;
  padding: 2.5rem !important;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin: 0 !important;
}

.card-title-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 184, 217, 0.1) !important;
  border: 1px solid rgba(0, 184, 217, 0.2) !important;
}

.card-title-icon .v-icon {
  font-size: 28px !important;
  color: #00B8D9 !important;
}

.card-title-text {
  font-size: 1.75rem !important;
  font-weight: 700 !important;
  color: #22292F !important;
  margin: 0 !important;
}

.card-subtitle {
  font-size: 1.1rem !important;
  color: #8898AA !important;
  margin-top: 0.5rem !important;
}

/* 授權管理標籤切換 - 重新設計 */
.auth-tabs {
  margin-left: auto;
  margin-top: 0;
  display: flex;
  gap: 0.75rem;
}

.auth-tab-btn {
  border-radius: 16px !important;
  padding: 1rem 2rem !important;
  font-weight: 600 !important;
  font-size: 1.1rem !important;
  transition: all 0.3s ease !important;
  text-transform: none !important;
  min-width: 140px !important;
  min-height: 56px !important;
  border: 2px solid transparent !important;
  background: transparent !important;
  color: #8898AA !important;
}

.auth-tab-btn.active {
  background: #00B8D9 !important;
  color: white !important;
  border-color: #00B8D9 !important;
  box-shadow: 0 4px 12px rgba(0, 184, 217, 0.3) !important;
}

.auth-tab-btn:not(.active) {
  background: #f8fafc !important;
  border-color: #e2e8f0 !important;
  color: #8898AA !important;
}

.auth-tab-btn:not(.active):hover {
  background: #f1f5f9 !important;
  border-color: #cbd5e1 !important;
  color: #22292F !important;
  transform: translateY(-2px) !important;
}

/* 待處理數量標籤 */
.pending-badge {
  margin-left: 1.5rem !important;
  background: #00B8D9 !important;
  color: white !important;
  border-radius: 20px !important;
  padding: 0.4rem 1rem !important;
  font-size: 0.9rem !important;
  font-weight: 600 !important;
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

/* 表格容器樣式 */
.table-container {
  padding: 0 !important;
  background: white !important;
}

/* 保持原有的表格樣式但優化 */
:deep(.v-table) {
  font-size: 20px !important;
}

:deep(.v-data-table) {
  font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
  border-radius: 0 !important;
  overflow: visible !important;
  box-shadow: none !important;
  background: white;
  border: none !important;
}

:deep(.v-data-table__wrapper) {
  overflow: visible !important;
  border-radius: 0 !important;
}

:deep(.v-data-table-header) {
  background: #f8fafc !important;
  border-bottom: 2px solid #e2e8f0 !important;
}

:deep(.v-data-table-header th) {
  font-size: 1.2rem !important;
  font-weight: 700 !important;
  color: #22292F !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  padding: 1.75rem 1.5rem !important;
  line-height: 1.4 !important;
  white-space: nowrap;
  background: transparent !important;
  vertical-align: middle;
  border-bottom: none !important;
  border-right: none !important;
}

:deep(.v-data-table tbody td) {
  font-size: 1.15rem !important;
  color: #22292F !important;
  padding: 1.75rem 1.5rem !important;
  line-height: 1.6 !important;
  vertical-align: middle;
  border-bottom: 1px solid #f1f5f9 !important;
  border-right: none !important;
  background: transparent !important;
}

:deep(.v-data-table-row) {
  height: auto !important;
  min-height: 90px !important;
  transition: all 0.3s ease;
  border: none !important;
}

:deep(.v-data-table-row:hover) {
  background: #f8fafc !important;
  transform: none !important;
  box-shadow: none !important;
}

:deep(.v-data-table-row:last-child td) {
  border-bottom: none !important;
}

/* 表格分頁樣式 */
:deep(.v-data-table-footer) {
  font-size: 1rem !important;
  padding: 1.75rem !important;
  border-top: 2px solid #e2e8f0 !important;
  background: #f8fafc !important;
  display: flex !important;
  align-items: center !important;
  justify-content: space-between !important;
  flex-wrap: wrap !important;
  gap: 1rem !important;
}

:deep(.v-data-table-footer .v-pagination__item) {
  font-size: 1rem !important;
  height: 44px !important;
  min-width: 44px !important;
  border-radius: 8px !important;
  border: 1px solid #e2e8f0 !important;
  margin: 0 2px !important;
  transition: all 0.3s ease !important;
}

:deep(.v-data-table-footer .v-pagination__item--is-active) {
  background: #00B8D9 !important;
  color: white !important;
  border-color: transparent !important;
}

/* 按鈕樣式統一 */
.modern-btn {
  margin: 3px;
  border-radius: 12px !important;
  font-weight: 600 !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
  padding: 1rem 2rem !important;
  font-size: 1rem !important;
  min-height: 48px !important;
}

.modern-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15) !important;
}

.modern-btn.approve-btn {
  background: #00B8D9 !important;
  color: white !important;
  border: none !important;
}

.modern-btn.approve-btn:hover {
  background: #0099B8 !important;
}

.modern-btn.reject-btn {
  background: #ef4444 !important;
  color: white !important;
  border: none !important;
}

/* 狀態標籤樣式 */
.status-chip {
  border-radius: 20px !important;
  font-weight: 600 !important;
  padding: 0.6rem 1.2rem !important;
  font-size: 0.95rem !important;
  border: none !important;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1) !important;
}

:deep(.v-chip--size-small) {
  border-radius: 16px !important;
  font-weight: 600 !important;
  padding: 0.5rem 1rem !important;
  font-size: 0.9rem !important;
}

:deep(.v-chip--size-large) {
  border-radius: 20px !important;
  font-weight: 600 !important;
  padding: 0.7rem 1.4rem !important;
  font-size: 1rem !important;
}

/* 檢查類型標籤 */
:deep(.type-chip) {
  background: #00B8D9 !important;
  color: white !important;
  border-radius: 16px !important;
  font-weight: 600 !important;
  padding: 0.7rem 1.4rem !important;
  font-size: 1rem !important;
  border: none !important;
}

/* 操作按鈕 */
:deep(.view-report-btn),
:deep(.share-btn) {
  border-radius: 12px !important;
  margin-left: 5px;
  background: #f8fafc !important;
  border: 1px solid #e2e8f0 !important;
  transition: all 0.3s ease !important;
  width: 52px !important;
  height: 52px !important;
}

:deep(.view-report-btn:hover) {
  background: #00B8D9 !important;
  color: white !important;
  border-color: transparent !important;
  transform: scale(1.05) !important;
}

:deep(.share-btn:hover) {
  background: #8898AA !important;
  color: white !important;
  border-color: transparent !important;
  transform: scale(1.05) !important;
}

:deep(.view-report-btn .v-icon),
:deep(.share-btn .v-icon) {
  font-size: 26px !important;
}

/* 表格內圖標 */
:deep(.v-data-table .d-flex.align-center .v-icon) {
  margin-right: 0.75rem !important;
  font-size: 24px !important;
}

/* 報告編號圖標 */
:deep(.v-data-table .d-flex.align-center .v-icon[class*="mdi-file-document"]) {
  color: #00B8D9 !important;
}

/* 日期圖標 */
:deep(.v-data-table .d-flex.align-center .v-icon[class*="mdi-calendar"]) {
  color: #8898AA !important;
}

/* 用戶頭像圖標 */
:deep(.v-data-table .v-avatar) {
  padding-left: 10px;
  width: 44px !important;
  height: 44px !important;
  border-radius: 12px !important;
  background: rgba(0, 184, 217, 0.1) !important;
  border: 10px solid rgba(0, 184, 217, 0.2) !important;
}

:deep(.v-data-table .v-avatar .v-icon) {
  font-size: 24px !important;
  color: #00B8D9 !important;
}

/* 操作按鈕圖標 */
:deep(.view-report-btn .v-icon) {
  font-size: 26px !important;
  color: #8898AA !important;
}

:deep(.view-report-btn:hover .v-icon) {
  color: white !important;
}

:deep(.share-btn .v-icon) {
  font-size: 26px !important;
  color: #8898AA !important;
}

:deep(.share-btn:hover .v-icon) {
  color: white !important;
}

/* 檢查類型標籤圖標 */
:deep(.type-chip .v-icon) {
  font-size: 20px !important;
  color: white !important;
}

/* 到期時間標籤圖標 */
.expiry-chip {
  border-radius: 16px !important;
  font-weight: 600 !important;
  padding: 0.6rem 1.2rem !important;
  font-size: 1rem !important;
  border: none !important;
}

.expiry-chip .v-icon {
  font-size: 18px !important;
  color: #8898AA !important;
}

/* 狀態圖標顏色 - 只保留必要的狀態區分 */
:deep(.v-chip .v-icon[class*="mdi-check-circle"]) {
  color: white !important;
}

:deep(.v-chip--variant-flat.v-chip--color-success) {
  background: #10b981 !important;
  color: white !important;
}

:deep(.v-chip .v-icon[class*="mdi-close-circle"]) {
  color: white !important;
}

:deep(.v-chip--variant-flat.v-chip--color-error) {
  background: #ef4444 !important;
  color: white !important;
}

:deep(.v-chip .v-icon[class*="mdi-clock"]) {
  color: #f59e0b !important;
}

/* 無資料狀態 */
:deep(.v-data-table__empty-wrapper) {
  padding: 4rem 2rem !important;
  text-align: center;
}

:deep(.v-data-table__empty-wrapper .v-icon) {
  font-size: 72px !important;
  color: #cbd5e1 !important;
  margin-bottom: 1rem !important;
}

:deep(.v-data-table__empty-wrapper .text-subtitle-1) {
  color: #22292F !important;
  font-weight: 600 !important;
  margin-bottom: 0.5rem !important;
  font-size: 1.2rem !important;
}

:deep(.v-data-table__empty-wrapper .text-body-2) {
  color: #8898AA !important;
  font-size: 1rem !important;
}

/* 表格內特定元素的樣式調整 */
:deep(.v-data-table .font-weight-medium) {
  font-size: 1.15rem !important;
  font-weight: 600 !important;
  color: #22292F !important;
}

:deep(.v-data-table .text-caption) {
  font-size: 1rem !important;
  color: #8898AA !important;
}

:deep(.v-data-table .text-body-2) {
  font-size: 1.1rem !important;
  line-height: 1.4 !important;
  color: #22292F !important;
}

/* 響應式設計 */
@media (max-width: 1200px) {
  .dashboard-container {
    padding: 1.5rem !important;
  }

  .card-content {
    min-height: 140px;
    padding: 1.75rem !important;
  }

  .icon-container {
    width: 64px;
    height: 64px;
  }

  .icon-container .v-icon {
    font-size: 32px !important;
  }

  .stat-card .stat-number {
    font-size: 2.5rem !important;
  }

  .user-card .card-title {
    font-size: 1.3rem !important;
  }

  .card-title-text {
    font-size: 1.5rem !important;
  }

  .card-subtitle {
    font-size: 1rem !important;
  }

  :deep(.v-data-table-header th) {
    font-size: 1rem !important;
    padding: 1.5rem 1rem !important;
  }

  :deep(.v-data-table tbody td) {
    font-size: 1rem !important;
    padding: 1.5rem 1rem !important;
  }
}

@media (max-width: 960px) {
  .dashboard-container {
    padding: 1rem !important;
  }

  .info-card {
    min-height: 140px;
  }

  .card-content {
    min-height: 140px;
    padding: 1.5rem !important;
    gap: 1.25rem;
  }

  .icon-container {
    width: 56px;
    height: 56px;
  }

  .icon-container .v-icon {
    font-size: 28px !important;
  }

  .stat-card .stat-number {
    font-size: 2.25rem !important;
  }

  .stat-card .stat-label {
    font-size: 1.1rem !important;
  }

  .user-card .card-title {
    font-size: 1.2rem !important;
  }

  .user-card .card-subtitle {
    font-size: 1rem !important;
  }

  .logout-btn {
    top: 1rem !important;
    right: 1rem !important;
    width: 20px !important;
    height: 32px !important;
  }

  .logout-btn .v-icon {
    font-size: 16px !important;
  }

  .card-header {
    padding: 2rem !important;
  }

  .card-title {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .auth-tabs {
    margin-left: 0;
    margin-top: 1rem;
  }

  .auth-tab-btn {
    padding: 0.6rem 1.2rem !important;
    font-size: 0.9rem !important;
    min-width: 100px !important;
  }

  :deep(.v-data-table-header th) {
    font-size: 0.95rem !important;
    padding: 1.25rem 0.75rem !important;
  }

  :deep(.v-data-table tbody td) {
    font-size: 0.9rem !important;
    padding: 1.25rem 0.75rem !important;
  }

  :deep(.v-data-table-footer) {
    flex-direction: column !important;
    align-items: stretch !important;
    gap: 1rem !important;
  }
}

@media (max-width: 600px) {
  .dashboard-container {
    padding: 0.75rem !important;
  }

  .info-card {
    min-height: 120px;
  }

  .card-content {
    min-height: 120px;
    padding: 1.25rem !important;
    gap: 1rem;
  }

  .icon-container {
    width: 48px;
    height: 48px;
    border-radius: 16px;
  }

  .icon-container .v-icon {
    font-size: 24px !important;
  }

  .stat-card .stat-number {
    font-size: 2rem !important;
  }

  .stat-card .stat-label {
    font-size: 1rem !important;
  }

  .user-card .card-title {
    font-size: 1.1rem !important;
  }

  .user-card .card-subtitle {
    font-size: 0.95rem !important;
  }

  .logout-btn {
    top: 1.25rem !important;
    right: 1.25rem !important;
    width: 32px !important;
    height: 32px !important;
  }

  .logout-btn .v-icon {
    font-size: 16px !important;
  }

  .card-header {
    padding: 1.5rem !important;
  }

  .modern-btn {
    padding: 0.75rem 1.5rem !important;
    font-size: 0.9rem !important;
    min-height: 42px !important;
  }

  .status-chip {
    padding: 0.4rem 0.8rem !important;
    font-size: 0.8rem !important;
  }

  :deep(.view-report-btn),
  :deep(.share-btn) {
    width: 40px !important;
    height: 40px !important;
  }

  :deep(.view-report-btn .v-icon),
  :deep(.share-btn .v-icon) {
    font-size: 18px !important;
  }

  :deep(.v-data-table td),
  :deep(.v-data-table-header th) {
    padding: 1rem 0.5rem !important;
    font-size: 0.85rem !important;
  }

  .auth-tabs {
    flex-direction: column;
    gap: 0.5rem;
  }

  .auth-tab-btn {
    width: 100%;
    justify-content: center;
  }
}

/* 動畫效果 */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.info-card {
  animation: fadeInUp 0.6s ease-out;
}

.main-card {
  animation: fadeInUp 0.8s ease-out;
}

/* 載入狀態 */
:deep(.v-data-table--loading) {
  position: relative;
}

:deep(.v-data-table--loading::after) {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  backdrop-filter: blur(2px);
  z-index: 10;
}

/* 健康檢查報告卡片的圖標顏色 */
.health-report-card .card-title-icon {
  background: rgba(0, 184, 217, 0.1) !important;
  border: 1px solid rgba(0, 184, 217, 0.2) !important;
}

/* 授權管理卡片的圖標顏色 */
.auth-management-card .card-title-icon {
  background: rgba(0, 184, 217, 0.1) !important;
  border: 1px solid rgba(0, 184, 217, 0.2) !important;
}

/* 報告編號圖標 */
:deep(.v-data-table .d-flex.align-center .v-icon[class*="mdi-file-document"]) {
  color: #64748b !important;
}

/* 日期圖標 */
:deep(.v-data-table .d-flex.align-center .v-icon[class*="mdi-calendar"]) {
  color: #64748b !important;
}

/* 用戶頭像圖標 */
:deep(.v-data-table .v-avatar) {
  width: 44px !important;
  height: 44px !important;
  border-radius: 12px !important;
  background: #ffffff00 !important;
}

:deep(.v-data-table .v-avatar .v-icon) {
  font-size: 24px !important;
  color: #64748b !important;
}

/* 操作按鈕圖標 */
:deep(.view-report-btn .v-icon) {
  font-size: 26px !important;
  color: #64748b !important;
}

:deep(.view-report-btn:hover .v-icon) {
  color: white !important;
}

:deep(.share-btn .v-icon) {
  font-size: 26px !important;
  color: #64748b !important;
}

:deep(.share-btn:hover .v-icon) {
  color: white !important;
}

/* 檢查類型標籤圖標 */
:deep(.type-chip .v-icon) {
  font-size: 20px !important;
  color: white !important;
}

/* 到期時間標籤圖標 */
.expiry-chip .v-icon {
  font-size: 18px !important;
  color: #64748b !important;
}

/* 狀態圖標顏色 - 只保留必要的狀態區分 */
:deep(.v-chip .v-icon[class*="mdi-check-circle"]) {
  color: white !important;
}

:deep(.v-chip--variant-flat.v-chip--color-success) {
  background: #10b981 !important;
  color: white !important;
}

:deep(.v-chip .v-icon[class*="mdi-close-circle"]) {
  color: white !important;
}

:deep(.v-chip--variant-flat.v-chip--color-error) {
  background: #ef4444 !important;
  color: white !important;
}

:deep(.v-chip .v-icon[class*="mdi-clock"]) {
  color: #f59e0b !important;
}

/* 用戶資訊卡片統一樣式 */
.user-info-card {
  height: 100%;
  min-height: 180px;
  border-radius: 24px !important;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%) !important;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1) !important;
  border: none !important;
  box-shadow: 0 12px 40px rgba(102, 126, 234, 0.15) !important;
  margin-bottom: 2rem;
  overflow: hidden;
  position: relative;
}

.user-info-card:hover {
  transform: translateY(-6px) scale(1.02);
  box-shadow: 0 20px 60px rgba(102, 126, 234, 0.25) !important;
}

.user-info-card .v-card-text {
  height: 100%;
  padding: 0 !important;
}

/* 用戶卡片包裝器 */
.user-card-wrapper {
  height: 100%;
  min-height: 180px;
  padding: 2rem !important;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  position: relative;
  z-index: 2;
}

/* 背景裝飾 */
.user-card-bg-decoration {
  position: absolute;
  top: -20px;
  right: -20px;
  width: 120px;
  height: 120px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  border-radius: 50%;
  z-index: 1;
}

.user-card-bg-decoration::before {
  content: '';
  position: absolute;
  bottom: -40px;
  left: -40px;
  width: 80px;
  height: 80px;
  background: radial-gradient(circle, rgba(255, 255, 255, 0.08) 0%, transparent 70%);
  border-radius: 50%;
}

/* 用戶卡片主體 */
.user-card-main {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  flex: 1;
}

/* 用戶頭像區域 */
.user-avatar-section {
  position: relative;
  flex-shrink: 0;
}

.user-avatar-container {
  width: 72px;
  height: 72px;
  border-radius: 20px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.2) !important;
  border: 3px solid rgba(255, 255, 255, 0.3) !important;
  backdrop-filter: blur(10px);
  transition: all 0.3s ease;
}

.user-avatar-container:hover {
  transform: scale(1.05);
  background: rgba(255, 255, 255, 0.25) !important;
}

.user-avatar-icon {
  font-size: 36px !important;
  color: rgba(255, 255, 255, 0.9) !important;
}

/* 狀態指示器 */
.user-status-indicator {
  position: absolute;
  bottom: 4px;
  right: 4px;
  width: 16px;
  height: 16px;
  background: #4ade80;
  border: 3px solid rgba(255, 255, 255, 0.8);
  border-radius: 50%;
  animation: pulse-green 2s infinite;
}

@keyframes pulse-green {
  0%, 100% {
    opacity: 1;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.1);
  }
}

/* 用戶資訊內容 */
.user-info-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.user-name {
  font-size: 1.5rem !important;
  font-weight: 700 !important;
  color: rgba(255, 255, 255, 0.95) !important;
  line-height: 1.2;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  margin: 0;
}

.user-role {
  font-size: 1rem !important;
  font-weight: 500 !important;
  color: rgba(255, 255, 255, 0.8) !important;
  display: flex;
  align-items: center;
  gap: 0.25rem;
  margin: 0;
}

.user-role .v-icon {
  color: rgba(255, 255, 255, 0.7) !important;
}

/* 登出按鈕 */
.user-logout-btn {
  position: absolute !important;
  top: 1.5rem !important;
  right: 1.5rem !important;
  width: 40px !important;
  height: 40px !important;
  border-radius: 14px !important;
  background: rgba(255, 255, 255, 0.15) !important;
  color: rgba(255, 255, 255, 0.9) !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  backdrop-filter: blur(10px);
  border: 1px solid rgba(255, 255, 255, 0.2) !important;
  z-index: 3;
}

.user-logout-btn:hover {
  background: rgba(239, 68, 68, 0.8) !important;
  color: white !important;
  transform: scale(1.05) rotate(5deg) !important;
  border-color: rgba(239, 68, 68, 0.3) !important;
  box-shadow: 0 6px 20px rgba(239, 68, 68, 0.3) !important;
}

.user-logout-btn .v-icon {
  font-size: 20px !important;
  transition: all 0.3s ease;
}

.user-logout-btn:hover .v-icon {
  transform: rotate(-5deg);
}

/* 用戶資訊卡片響應式設計 */
@media (max-width: 1200px) {
  .user-card-wrapper {
    min-height: 160px;
    padding: 1.75rem !important;
  }
  
  .user-avatar-container {
    width: 64px;
    height: 64px;
  }
  
  .user-avatar-icon {
    font-size: 32px !important;
  }
  
  .user-name {
    font-size: 1.3rem !important;
  }
}

@media (max-width: 960px) {
  .user-info-card {
    min-height: 160px;
  }
  
  .user-card-wrapper {
    min-height: 160px;
    padding: 1.5rem !important;
    gap: 1.25rem;
  }
  
  .user-card-main {
    gap: 1.25rem;
  }
  
  .user-avatar-container {
    width: 56px;
    height: 56px;
  }
  
  .user-avatar-icon {
    font-size: 28px !important;
  }
  
  .user-name {
    font-size: 1.2rem !important;
  }
  
  .user-role {
    font-size: 0.95rem !important;
  }
  
  .user-logout-btn {
    top: 1rem !important;
    right: 1rem !important;
    width: 36px !important;
    height: 36px !important;
  }
  
  .user-logout-btn .v-icon {
    font-size: 18px !important;
  }
}

@media (max-width: 600px) {
  .user-info-card {
    min-height: 140px;
  }
  
  .user-card-wrapper {
    min-height: 140px;
    padding: 1.25rem !important;
    gap: 1rem;
  }
  
  .user-card-main {
    gap: 1rem;
  }
  
  .user-avatar-container {
    width: 48px;
    height: 48px;
    border-radius: 16px;
  }
  
  .user-avatar-icon {
    font-size: 24px !important;
  }
  
  .user-name {
    font-size: 1.1rem !important;
  }
  
  .user-role {
    font-size: 0.9rem !important;
  }
  
  .user-logout-btn {
    top: 1rem !important;
    right: 1rem !important;
    width: 32px !important;
    height: 32px !important;
  }
  
  .user-logout-btn .v-icon {
    font-size: 16px !important;
  }
}
</style>