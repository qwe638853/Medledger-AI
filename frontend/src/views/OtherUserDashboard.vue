<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import healthCheckService from '../services/healthCheckService';
import { useRouter } from 'vue-router';

const authStore = useAuthStore();
const currentUser = computed(() => authStore.currentUser);
const token = ref(authStore.token);
const patientReports = ref([]);
const loading = ref(false);
const router = useRouter();

// 分類報告
const authorizedReports = computed(() => patientReports.value.filter(report => report.is_authorized));
const unauthorizedReports = computed(() => patientReports.value.filter(report => !report.is_authorized));

// 所有已授權報告
const allAuthorizedReports = ref([]);
const authorizedReportsDialog = ref(false);
const loadingAuthorizedReports = ref(false);

// 視圖模式
const viewMode = ref('search'); // 'search' 或 'authorized' 或 'pending' 或 'history'

// 搜尋相關
const searchLoading = ref(false);
const patientId = ref('');
const showSearchResults = ref(false);

// 標籤頁
const tab = ref('unauthorized');

// 授權請求對話框
const authRequestDialog = ref(false);
const selectedReport = ref(null);
const authReason = ref('');
const authExpiry = ref('');
const requestLoading = ref(false);

// Snackbar
const snackbar = ref(false);
const snackbarMessage = ref('');
const snackbarColor = ref('success');

// Dashboard 數據
const dashboardStats = ref({
  totalAuthorized: 0,
  pendingRequests: 0,
});

// 待授權請求相關
const myAccessRequests = ref([]);
const loadingMyRequests = ref(false);

// 移除 pendingRequestsTab，改用兩個獨立的計算屬性
const pendingRequests = computed(() => 
  myAccessRequests.value.filter(req => req.status === 'PENDING')
);

const historyRequests = computed(() => 
  myAccessRequests.value.filter(req => req.status !== 'PENDING')
);

// 分頁設置
const itemsPerPage = ref(5);

// 新增游標分頁相關函數
const cursor = ref(null);
const hasNextPage = ref(true);

const fetchDashboardStats = async () => {
  try {
    //dashboardStats.value = await healthCheckService.fetchDashboardStats();
  } catch (error) {
    console.error('獲取 Dashboard 數據時出錯:', error);
    // 錯誤處理已經在服務層完成
  }
};

// 計算剩餘天數
const getRemainingDays = (expiry) => {
  if (!expiry) return '';
  const now = Math.floor(Date.now() / 1000);
  const days = Math.ceil((expiry - now) / (24 * 60 * 60));
  if (days < 0) return '已過期';
  if (days === 0) return '今日到期';
  return `剩餘 ${days} 天`;
};

// 顯示授權詳情
const showAuthDetails = (item) => {
  // 切換展開狀態
  item.expanded = !item.expanded;
};

// 獲取所有已授權報告
const fetchAllAuthorizedReports = async () => {
  loadingAuthorizedReports.value = true;
  try {
    const response = await healthCheckService.fetchAuthorizedReports({
      cursor: cursor.value,
      limit: 5
    });
    console.log('已授權報告回應:', response);
    
    if (response && Array.isArray(response)) {
      // 直接處理回應陣列
      allAuthorizedReports.value = response.map(report => ({
        ...report,
        expanded: false,
        id: report.id ,
        requester_name: report.patient_id,
        granted_at: report.date ,
        expiry: report.expiry || null
      }));
      
      dashboardStats.value.totalAuthorized = allAuthorizedReports.value.length;
    } else if (response && response.reports) {
      // 處理包含 reports 欄位的回應
      cursor.value = response.next_cursor;
      hasNextPage.value = response.has_next_page;
      
      allAuthorizedReports.value = response.reports.map(report => ({
        ...report,
        expanded: false,
        id: report.report_id || report.id,
        company_name: report.company_name || '未知公司',
        requester_name: report.requester_name || '未知請求者',
        granted_at: report.granted_at || report.created_at || new Date().toISOString(),
        expiry: report.expiry || null
      }));
      
      dashboardStats.value.totalAuthorized = response.total || allAuthorizedReports.value.length;
    } else {
      // 如果沒有資料，設為空陣列
      allAuthorizedReports.value = [];
      dashboardStats.value.totalAuthorized = 0;
    }
    
    console.log('處理後的已授權報告:', allAuthorizedReports.value);
  } catch (error) {
    console.error('獲取已授權報告時出錯:', error);
    snackbarMessage.value = '獲取已授權報告時出錯';
    snackbarColor.value = 'error';
    snackbar.value = true;
    allAuthorizedReports.value = [];
  } finally {
    loadingAuthorizedReports.value = false;
  }
};

// 獲取我發出的授權請求
const fetchMyAccessRequests = async () => {
  loadingMyRequests.value = true;
  try {
    myAccessRequests.value = await healthCheckService.listMyAccessRequests();
    // 更新儀表板數據
    dashboardStats.value.pendingRequests = pendingRequests.value.length;
  } catch (error) {
    console.error('獲取授權請求失敗:', error);
    snackbarMessage.value = '獲取授權請求失敗';
    snackbarColor.value = 'error';
    snackbar.value = true;
  } finally {
    loadingMyRequests.value = false;
  }
};

// 切換到已授權報告視圖
const switchToAuthorizedView = async () => {
  if (allAuthorizedReports.value.length === 0) {
    await fetchAllAuthorizedReports();
  }
  viewMode.value = 'authorized';
  showSearchResults.value = false;
};

// 切換到搜尋視圖
const switchToSearchView = () => {
  viewMode.value = 'search';
};

// 切換到待授權請求視圖
const switchToPendingView = async () => {
  if (myAccessRequests.value.length === 0) {
    await fetchMyAccessRequests();
  }
  viewMode.value = 'pending';
  showSearchResults.value = false;
};

// 切換到歷史紀錄視圖
const switchToHistoryView = async () => {
  if (myAccessRequests.value.length === 0) {
    await fetchMyAccessRequests();
  }
  viewMode.value = 'history';
  showSearchResults.value = false;
};

// 打開授權報告對話框（查看詳細）
const openAuthorizedReportsDialog = () => {
  authorizedReportsDialog.value = true;
};

// 搜尋並獲取病患報告
const searchPatientReports = async () => {
  if (!patientId.value || patientId.value.trim() === '') {
    snackbarMessage.value = '請輸入病患身分證字號';
    snackbarColor.value = 'error';
    snackbar.value = true;
    return;
  }

  switchToSearchView();
  searchLoading.value = true;
  try {
    await fetchReportsByPatientId(patientId.value);
    showSearchResults.value = true;
  } catch (error) {
    console.error('搜尋病患報告時出錯:', error);
    // 錯誤訊息已在服務層處理
    patientReports.value = [];
    showSearchResults.value = false;
  } finally {
    searchLoading.value = false;
  }
};

// 獲取病患報告元數據
const fetchReportsByPatientId = async (id) => {
  try {
    // 使用新的 API 獲取報告元數據
    patientReports.value = await healthCheckService.fetchReportMetaByPatientID(id);
    
    // 根據授權狀態自動選擇預設Tab
    if (unauthorizedReports.value.length > 0) {
      tab.value = 'unauthorized';
    } else if (authorizedReports.value.length > 0) {
      tab.value = 'authorized';
    }
  } catch (error) {
    console.error('獲取報告失敗:', error);
    throw error;
  }
};

// 開啟授權請求對話框
const openAuthRequestDialog = (report) => {
  selectedReport.value = report;
  authReason.value = '';
  // 設定預設日期為30天後
  const today = new Date();
  const thirtyDaysLater = new Date(today);
  thirtyDaysLater.setDate(today.getDate() + 30);
  authExpiry.value = thirtyDaysLater.toISOString().substr(0, 10);
  authRequestDialog.value = true;
};

// 發送授權請求
const sendAuthRequest = async () => {
  if (!authReason.value || !authExpiry.value) {
    snackbarMessage.value = '請填寫完整授權資訊';
    snackbarColor.value = 'error';
    snackbar.value = true;
    return;
  }

  requestLoading.value = true;
  try {
    // 使用 requestReportAccess 函數發送授權請求
    await healthCheckService.requestReportAccess(
      selectedReport.value.id,
      patientId.value,
      authReason.value,
      authExpiry.value
    );

    authRequestDialog.value = false;
    snackbarMessage.value = '已送出授權請求';
    snackbarColor.value = 'success';
    snackbar.value = true;
    
    // 授權請求成功後重新獲取報告
    await fetchReportsByPatientId(patientId.value);
  } catch (error) {
    console.error('發送授權請求時出錯:', error);
    // 錯誤訊息已在服務層處理
  } finally {
    requestLoading.value = false;
  }
};

// 初次載入頁面時獲取全局數據
onMounted(async () => {
  // 先取得後端資料
  await fetchAllAuthorizedReports();
  await fetchMyAccessRequests();
  await fetchDashboardStats();

  /*
  // [前端測試用] 塞一筆假授權報告，繞過後端
  allAuthorizedReports.value.push({
    id: 'RPT123456',
    patient_id: 'A123456789',
    date: new Date().toISOString(),
    expiry: new Date(Date.now() + 7 * 24 * 60 * 60 * 1000).toISOString(), // 7天後到期
    content: '血壓: 120/80, 血糖: 90',
    is_authorized: true
  });
  // [前端測試用] 同步更新儀表板數量
  dashboardStats.value.totalAuthorized += 1;
  dashboardStats.value.totalPatients += 1;
  */
});

const handleLogout = () => {
  authStore.logout();
};

// 查看報告詳情
const goToReportDetail = (item) => {
  console.log('查看報告詳情:', item);
  const report_id = item.report_id || item.id;
  const patient_id = item.patient_id || item.patientId;
  
  if (!report_id || !patient_id) {
    snackbarMessage.value = '無法查看報告：缺少必要資訊';
    snackbarColor.value = 'error';
    snackbar.value = true;
    return;
  }
  
  router.push({ 
    name: 'ReportDetail', 
    params: { 
      report_id,
      patient_id
    }
  });
};

// 計算指標值（用於圓形進度條）
const calculateMetricValue = (value) => {
  if (!value) return 0;
  
  // 提取數值部分
  const numStr = value.toString().match(/[\d.]+/);
  if (!numStr) return 0;
  
  const numValue = parseFloat(numStr[0]);
  if (isNaN(numValue)) return 0;
  
  // 根據不同指標類型設定最大值
  const maxValues = {
    'Glu-AC': 200,    // 血糖
    'HbA1c': 10,      // 糖化血色素
    'LDL-C': 200,     // 低密度脂蛋白膽固醇
    'HDL-C': 100,     // 高密度脂蛋白膽固醇
    'TG': 500,        // 三酸甘油脂
    'T-CHO': 300,     // 總膽固醇
    'BP': 200,        // 血壓
    'WBC': 20,        // 白血球
    'RBC': 8,         // 紅血球
    'Hb': 20,         // 血紅素
    'Hct': 60,        // 血球容積比
    'PLT': 500,       // 血小板
    'AST（GOT）': 100, // 天門冬胺酸轉胺酶
    'ALT（GPT）': 100, // 丙胺酸轉胺酶
    'ALP': 200,       // 鹼性磷酸酶
    'UN': 50,         // 尿素氮
    'CRE': 5,         // 肌酸酐
    'U.A': 10,        // 尿酸
  };

  // 獲取對應的最大值，如果沒有定義則使用 100 作為預設值
  const maxValue = maxValues[Object.keys(maxValues).find(key => value.includes(key))] || 100;
  
  // 計算百分比值
  return Math.min(Math.max((numValue / maxValue) * 100, 0), 100);
};

// 獲取指標顏色
const getMetricColor = (value) => {
  if (!value) return 'grey';
  
  // 提取數值部分
  const numStr = value.toString().match(/[\d.]+/);
  if (!numStr) return 'grey';
  
  const numValue = parseFloat(numStr[0]);
  if (isNaN(numValue)) return 'grey';
  
  // 根據不同指標類型設定正常值範圍
  const ranges = {
    'Glu-AC': { min: 70, max: 100 },    // 血糖
    'HbA1c': { min: 4, max: 5.7 },      // 糖化血色素
    'LDL-C': { min: 0, max: 130 },      // 低密度脂蛋白膽固醇
    'HDL-C': { min: 40, max: 60 },      // 高密度脂蛋白膽固醇
    'TG': { min: 0, max: 150 },         // 三酸甘油脂
    'T-CHO': { min: 0, max: 200 },      // 總膽固醇
    'BP': { min: 90, max: 140 },        // 血壓（收縮壓）
    'WBC': { min: 4, max: 10 },         // 白血球
    'RBC': { min: 4.5, max: 5.5 },      // 紅血球
    'Hb': { min: 12, max: 16 },         // 血紅素
    'Hct': { min: 37, max: 50 },        // 血球容積比
    'PLT': { min: 150, max: 450 },      // 血小板
    'AST（GOT）': { min: 0, max: 40 },   // 天門冬胺酸轉胺酶
    'ALT（GPT）': { min: 0, max: 40 },   // 丙胺酸轉胺酶
    'ALP': { min: 30, max: 100 },       // 鹼性磷酸酶
    'UN': { min: 7, max: 20 },          // 尿素氮
    'CRE': { min: 0.6, max: 1.2 },      // 肌酸酐
    'U.A': { min: 2.5, max: 7.2 },      // 尿酸
  };

  // 獲取對應的範圍，如果沒有定義則使用預設範圍
  const range = ranges[Object.keys(ranges).find(key => value.includes(key))] || { min: 0, max: 100 };
  
  if (numValue < range.min) return 'red';
  if (numValue > range.max) return 'orange';
  return 'green';
};

// 格式化日期
const formatDate = (date) => {
  if (!date) return '-';
  return new Date(date).toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });
};

// 獲取授權到期 Chip 的顏色
const getExpiryChipColor = (expiry) => {
  if (!expiry) return 'green-lighten-4';
  
  const expiryDate = new Date(expiry);
  const now = new Date();
  const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24));
  
  if (daysUntilExpiry < 0) return 'red-lighten-4';
  if (daysUntilExpiry <= 7) return 'orange-lighten-4';
  return 'green-lighten-4';
};

// 獲取授權到期 Chip 的文字顏色
const getExpiryTextColor = (expiry) => {
  if (!expiry) return 'green-darken-2';
  
  const expiryDate = new Date(expiry);
  const now = new Date();
  const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24));
  
  if (daysUntilExpiry < 0) return 'red-darken-2';
  if (daysUntilExpiry <= 7) return 'orange-darken-2';
  return 'green-darken-2';
};

// 格式化授權到期日期
const formatExpiryDate = (expiry) => {
  if (!expiry) return '永久';
  
  const expiryDate = new Date(expiry);
  const now = new Date();
  const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24));
  
  if (daysUntilExpiry < 0) return '已過期';
  if (daysUntilExpiry <= 7) return `${daysUntilExpiry} 天後到期`;
  
  return expiryDate.toLocaleDateString('zh-TW', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  });
};

// 獲取請求狀態的顯示文字和顏色
const getRequestStatusInfo = (status) => {
  const statusMap = {
    'PENDING': { text: '待審核', color: 'warning', icon: 'mdi-clock-outline' },
    'APPROVED': { text: '已通過', color: 'success', icon: 'mdi-check-circle' },
    'REJECTED': { text: '已拒絕', color: 'error', icon: 'mdi-close-circle' }
  };
  return statusMap[status] || { text: '未知', color: 'grey', icon: 'mdi-help-circle' };
};
</script>

<template>
  <div class="dashboard-bg">
    <v-container class="dashboard-container pa-6 pa-sm-8">
      <!-- 頂部統計卡片區 -->
      <v-row class="mb-8">
        <!-- 用戶資訊卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card" elevation="2">
            <v-card-text class="pa-8">
              <div class="d-flex align-center mb-4">
                <div class="rounded-circle bg-primary-lighten-5 p-3">
                  <v-icon size="28" color="primary">mdi-briefcase-account</v-icon>
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
                <div class="text-subtitle-1 text-grey-darken-1">保險業者</div>
                <div class="text-h6 font-weight-bold">{{ currentUser }}</div>
                <v-chip
                  size="small"
                  color="success-lighten-5"
                  class="font-weight-medium px-2 mt-2"
                  variant="flat"
                >
                  <v-icon start size="16" color="success">mdi-shield-check</v-icon>
                  <span class="text-success">已驗證業者</span>
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

        <!-- 已授權報告統計卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card" elevation="2">
            <v-card-text class="pa-8">
              <div class="d-flex align-center mb-4">
                <div class="rounded-circle bg-success-lighten-5 p-3">
                  <v-icon size="28" color="success">mdi-file-document</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-grey-darken-1">已授權報告</div>
                <div class="d-flex align-center justify-space-between mt-2">
                  <div class="text-h4 font-weight-bold">{{ dashboardStats.totalAuthorized }}</div>
                  <v-btn
                    color="success"
                    variant="tonal"
                    @click="switchToAuthorizedView"
                    class="view-btn"
                  >
                    <v-icon>mdi-eye</v-icon>
                  </v-btn>
                </div>
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

        <!-- 待處理請求卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card" elevation="2">
            <v-card-text class="pa-8">
              <div class="d-flex align-center mb-4">
                <div class="rounded-circle bg-warning-lighten-5 p-3">
                  <v-icon size="28" color="warning">mdi-clock-outline</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-grey-darken-1">待處理請求</div>
                <div class="d-flex align-center justify-space-between mt-2">
                  <div class="text-h4 font-weight-bold">{{ dashboardStats.pendingRequests }}</div>
                  <v-btn
                    color="warning"
                    variant="tonal"
                    @click="switchToPendingView"
                    class="view-btn"
                  >
                    <v-icon>mdi-eye</v-icon>
                  </v-btn>
                </div>
                <div class="stat-progress mt-4">
                  <v-progress-linear
                    :model-value="(dashboardStats.pendingRequests / 10) * 100"
                    color="warning"
                    height="4"
                    rounded
                  ></v-progress-linear>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 歷史紀錄卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card" elevation="2">
            <v-card-text class="pa-8">
              <div class="d-flex align-center mb-4">
                <div class="rounded-circle bg-info-lighten-5 p-3">
                  <v-icon size="28" color="info">mdi-history</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-grey-darken-1">歷史紀錄</div>
                <div class="d-flex align-center justify-space-between mt-2">
                  <div class="text-h4 font-weight-bold">{{ historyRequests.length }}</div>
                  <v-btn
                    color="info"
                    variant="tonal"
                    @click="switchToHistoryView"
                    class="view-btn"
                  >
                    <v-icon>mdi-eye</v-icon>
                  </v-btn>
                </div>
                <div class="stat-progress mt-4">
                  <v-progress-linear
                    :model-value="(historyRequests.length / 10) * 100"
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

      <!-- 搜尋區塊 -->
      <v-row class="mb-8">
        <v-col cols="12">
          <v-card class="search-card pa-6 elevation-2 rounded-xl">
            <v-card-text class="pa-6">
              <v-row align="center">
                <v-col cols="12" md="4">
                  <div class="d-flex align-center">
                    <v-icon size="24" color="primary" class="me-2">mdi-magnify</v-icon>
                    <span class="text-h6 font-weight-bold">病患健康報告查詢</span>
                  </div>
                </v-col>
                <v-col cols="12" md="5">
                  <v-text-field
                    v-model="patientId"
                    label="請輸入病患身分證字號"
                    placeholder="例如：A123456789"
                    variant="outlined"
                    density="comfortable"
                    hide-details
                    prepend-inner-icon="mdi-card-account-details-outline"
                    class="search-input"
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="3" class="d-flex gap-3">
                  <v-btn
                    color="primary"
                    :loading="searchLoading"
                    :disabled="searchLoading"
                    @click="searchPatientReports"
                    prepend-icon="mdi-magnify"
                    elevation="0"
                    class="flex-grow-1 search-btn"
                  >
                    搜尋
                  </v-btn>
                </v-col>
              </v-row>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 表格區域 -->
      <v-row v-if="showSearchResults && viewMode === 'search'" justify="center">
        <v-col cols="12">
          <v-card class="result-card mb-6 elevation-2 rounded-xl">
            <v-card-title class="d-flex align-center py-4 px-6 bg-grey-lighten-4">
              <div class="d-flex align-center">
                <v-icon size="24" color="primary" class="me-3">mdi-file-document-outline</v-icon>
                <span class="text-h6 font-weight-bold">病患「{{ patientId }}」的健康報告</span>
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

            <!-- 標籤頁 -->
            <v-tabs
              v-model="tab"
              color="grey-darken-3"
              align-tabs="start"
              class="px-4 pt-2"
              slider-color="yellow-accent-3"
            >
              <v-tab value="unauthorized" class="font-weight-medium">
                <v-icon start class="me-1" color="orange-darken-2" size="small">mdi-lock-open-alert</v-icon>
                尚未授權報告
                <v-badge
                  :content="unauthorizedReports.length.toString()"
                  :color="unauthorizedReports.length > 0 ? 'orange-darken-1' : 'grey'"
                  offset-x="3"
                  offset-y="-3"
                  size="small"
                  class="mb-4 ms-2"
                ></v-badge>
              </v-tab>
              <v-tab value="authorized" class="font-weight-medium">
                <v-icon start class="me-1" color="green-darken-1" size="small">mdi-lock-check</v-icon>
                已授權報告
                <v-badge
                  :content="authorizedReports.length.toString()"
                  :color="authorizedReports.length > 0 ? 'green-darken-1' : 'grey'"
                  offset-x="3"
                  offset-y="-3"
                  size="small"
                  class="mb-4 ms-2"
                ></v-badge>
              </v-tab>
            </v-tabs>
            
            <v-divider></v-divider>
            
            <v-card-text class="pa-6">
              <v-window v-model="tab">
                <!-- 尚未授權報告標籤內容 -->
                <v-window-item value="unauthorized">
                  <v-data-table
                    :headers="[
                      { title: '報告編號', key: 'id', align: 'start', width: '120px' },
                      { title: '健康數據', key: 'content', align: 'start' },
                      { title: '日期', key: 'date', align: 'center', width: '120px' },
                      { title: '狀態', key: 'status', align: 'center', width: '100px' },
                      { title: '操作', key: 'actions', align: 'center', width: '100px' }
                    ]"
                    :items="unauthorizedReports"
                    :loading="loading"
                    loading-text="資料載入中..."
                    class="elevation-0"
                    hover
                    density="comfortable"
                  >
                    <template v-slot:item.id="{ item }">
                      <div class="d-flex align-center">
                        <v-icon size="18" color="primary" class="me-2">mdi-file-document</v-icon>
                        <span class="font-weight-medium">{{ item.id }}</span>
                      </div>
                    </template>

                    <template v-slot:item.content="{ item }">
                      <div class="text-truncate">{{ item.content }}</div>
                    </template>

                    <template v-slot:item.date="{ item }">
                      <div class="d-flex align-center justify-center">
                        <v-icon size="16" color="grey" class="me-1">mdi-calendar</v-icon>
                        {{ formatDate(item.date) }}
                      </div>
                    </template>

                    <template v-slot:item.status="{ item }">
                      <v-chip
                        size="small"
                        color="warning"
                        variant="tonal"
                        class="font-weight-medium"
                      >
                        未授權
                      </v-chip>
                    </template>

                    <template v-slot:item.actions="{ item }">
                      <div class="d-flex gap-2 justify-center">
                        <v-btn
                          icon
                          variant="text"
                          size="small"
                          color="primary"
                          @click="openAuthRequestDialog(item)"
                          class="view-report-btn"
                        >
                          <v-icon>mdi-key-chain</v-icon>
                          <v-tooltip
                            activator="parent"
                            location="top"
                            open-delay="200"
                          >
                            請求授權
                          </v-tooltip>
                        </v-btn>
                      </div>
                    </template>

                    <template v-slot:no-data>
                      <div class="text-center pa-5">
                        <v-icon size="40" color="grey-lighten-1" class="mb-3">
                          mdi-file-document-outline
                        </v-icon>
                        <div class="text-subtitle-1 font-weight-medium text-grey-darken-1">
                          無未授權報告
                        </div>
                        <div class="text-body-2 text-grey mt-2">
                          該病患沒有未授權的健康報告
                        </div>
                      </div>
                    </template>
                  </v-data-table>
                </v-window-item>
                
                <!-- 已授權報告標籤內容 -->
                <v-window-item value="authorized">
                  <v-data-table
                    :headers="[
                      { title: '報告編號', key: 'id', align: 'start', width: '120px' },
                      { title: '授權對象', key: 'requester', align: 'start', width: '200px' },
                      { title: '授權日期', key: 'granted_at', align: 'center', width: '120px' },
                      { title: '到期日期', key: 'expiry', align: 'center', width: '120px' },
                      { title: '操作', key: 'actions', align: 'center', width: '100px' }
                    ]"
                    :items="authorizedReports"
                    :loading="loading"
                    loading-text="資料載入中..."
                    class="elevation-0"
                    hover
                    density="comfortable"
                  >
                    <template v-slot:item.id="{ item }">
                      <div class="d-flex align-center">
                        <v-icon size="18" color="primary" class="me-2">mdi-file-document</v-icon>
                        <span class="font-weight-medium">{{ item.id }}</span>
                      </div>
                    </template>

                    <template v-slot:item.requester="{ item }">
                      <div class="d-flex align-center">
                        <v-avatar size="28" color="primary" class="me-2">
                          <v-icon color="white" size="16">mdi-account</v-icon>
                        </v-avatar>
                        <div class="d-flex flex-column">
                          <span class="font-weight-medium">{{ item.requester_name }}</span>
                          <span class="text-caption text-grey">{{ item.company_name }}</span>
                        </div>
                      </div>
                    </template>

                    <template v-slot:item.granted_at="{ item }">
                      <div class="d-flex align-center justify-center">
                        <v-icon size="16" color="grey" class="me-1">mdi-calendar</v-icon>
                        {{ formatDate(item.granted_at) }}
                      </div>
                    </template>

                    <template v-slot:item.expiry="{ item }">
                      <div class="d-flex flex-column align-center">
                        <v-chip
                          size="small"
                          :color="getExpiryChipColor(item.expiry)"
                          variant="tonal"
                          class="font-weight-medium"
                        >
                          {{ formatExpiryDate(item.expiry) }}
                        </v-chip>
                      </div>
                    </template>

                    <template v-slot:item.actions="{ item }">
                      <div class="d-flex gap-2 justify-center">
                        <v-btn
                          icon
                          variant="text"
                          size="small"
                          color="primary"
                          @click="goToReportDetail(item)"
                          class="view-report-btn"
                        >
                          <v-icon>mdi-eye</v-icon>
                          <v-tooltip
                            activator="parent"
                            location="top"
                            open-delay="200"
                          >
                            查看報告
                          </v-tooltip>
                        </v-btn>
                      </div>
                    </template>

                    <template v-slot:no-data>
                      <div class="text-center pa-5">
                        <v-icon size="40" color="grey-lighten-1" class="mb-3">
                          mdi-file-document-outline
                        </v-icon>
                        <div class="text-subtitle-1 font-weight-medium text-grey-darken-1">
                          尚無已授權報告
                        </div>
                        <div class="text-body-2 text-grey mt-2">
                          目前您沒有任何已授權的健康報告
                        </div>
                      </div>
                    </template>
                  </v-data-table>
                </v-window-item>
              </v-window>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 全部已授權報告區塊 -->
      <v-row v-if="viewMode === 'authorized'" justify="center">
        <v-col cols="12">
          <!-- 返回按鈕 -->
          <v-btn
            color="grey-darken-1"
            variant="outlined"
            size="small"
            @click="switchToSearchView"
            class="mb-6"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="mb-6 elevation-2 rounded-xl">
            <v-card-title class="d-flex align-center py-4 px-6 bg-grey-lighten-4">
              <div class="d-flex align-center">
                <v-icon size="24" color="primary" class="me-3">mdi-folder-account</v-icon>
                <span class="text-h6 font-weight-bold">所有已授權健康報告</span>
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
                  { 
                    title: '報告編號',
                    key: 'id',
                    align: 'start',
                    width: '150px'
                  },
                  { 
                    title: '授權對象',
                    key: 'requester_name',
                    align: 'start',
                    width: '200px'
                  },
                  { 
                    title: '授權到期日',
                    key: 'expiry',
                    align: 'center',
                    width: '150px'
                  },
                  { 
                    title: '操作',
                    key: 'actions',
                    align: 'center',
                    width: '100px',
                    sortable: false
                  }
                ]"
                :items="allAuthorizedReports"
                :loading="loadingAuthorizedReports"
                loading-text="正在載入已授權報告..."
                class="elevation-0"
                hover
                density="comfortable"
                v-model:items-per-page="itemsPerPage"
                :items-per-page-options="[5, 10, 15]"
              >
                <!-- 報告編號欄位 -->
                <template v-slot:item.id="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="primary" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.id }}</span>
                  </div>
                </template>

                <!-- 授權對象欄位 -->
                <template v-slot:item.requester_name="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="grey" class="me-2">mdi-account</v-icon>
                    <span class="font-weight-medium">{{ item.requester_name }}</span>
                  </div>
                </template>

                <!-- 授權到期日欄位 -->
                <template v-slot:item.expiry="{ item }">
                  <v-chip
                    size="small"
                    :color="getExpiryChipColor(item.expiry)"
                    variant="tonal"
                    class="font-weight-medium"
                  >
                    {{ formatExpiryDate(item.expiry) }}
                  </v-chip>
                </template>

                <!-- 操作按鈕欄位 -->
                <template v-slot:item.actions="{ item }">
                  <div class="d-flex gap-2 justify-center">
                    <v-btn
                      icon
                      variant="text"
                      size="small"
                      color="primary"
                      @click="goToReportDetail(item)"
                      class="view-report-btn"
                    >
                      <v-icon>mdi-eye</v-icon>
                      <v-tooltip
                        activator="parent"
                        location="top"
                        open-delay="200"
                      >
                        查看報告
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
                      尚無已授權報告
                    </div>
                    <div class="text-body-2 text-grey mt-2">
                      目前您沒有任何已授權的健康報告
                    </div>
                  </div>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 未搜尋或搜尋前的提示 -->
      <v-row v-if="!showSearchResults && viewMode === 'search'" justify="center">
        <v-col cols="12">
          <v-card class="text-center py-6 px-5" elevation="1" color="grey-lighten-5">
            <v-icon size="50" color="blue-lighten-3" class="mb-3">mdi-file-document-search</v-icon>
            <div class="text-h6 font-weight-medium text-blue-grey-darken-1 mb-1">請輸入病患身分證字號查詢健康報告</div>
            <div class="text-body-2 text-grey px-6 mb-0">
              您可以輸入病患身分證字號查詢特定病患，或點擊「授權報告」查看所有已授權報告
            </div>
          </v-card>
        </v-col>
      </v-row>

      <!-- 待授權請求列表區塊 -->
      <v-row v-if="viewMode === 'pending'" justify="center">
        <v-col cols="12">
          <!-- 返回按鈕 -->
          <v-btn
            color="grey-darken-1"
            variant="outlined"
            size="small"
            @click="switchToSearchView"
            class="mb-6"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="mb-6 elevation-2 rounded-xl">
            <v-card-title class="d-flex align-center py-4 px-6 bg-grey-lighten-4">
              <div class="d-flex align-center">
                <v-icon size="24" color="warning" class="me-3">mdi-clock-alert</v-icon>
                <span class="text-h6 font-weight-bold">待處理授權請求</span>
                <v-chip
                  class="ms-3"
                  color="warning"
                  size="small"
                  variant="tonal"
                >
                  {{ pendingRequests.length }} 筆請求
                </v-chip>
              </div>
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
                    title: '授權對象',
                    key: 'patientHash',
                    align: 'start',
                    width: '120px'
                  },
                  { 
                    title: '申請日期',
                    key: 'requestedAt',
                    align: 'center',
                    width: '120px'
                  }
                ]"
                :items="pendingRequests"
                :loading="loadingMyRequests"
                loading-text="正在載入授權請求..."
                class="elevation-0"
                hover
                density="comfortable"
                v-model:items-per-page="itemsPerPage"
                :items-per-page-options="[5, 10]"
              >
                <!-- 報告編號欄位 -->
                <template v-slot:item.reportId="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="primary" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.reportId }}</span>
                  </div>
                </template>

                <!-- 病患雜湊欄位 -->
                <template v-slot:item.patientHash="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="grey" class="me-2">mdi-identifier</v-icon>
                    <span class="font-weight-medium">{{ item.patientHash }}</span>
                  </div>
                </template>

                <!-- 申請日期欄位 -->
                <template v-slot:item.requestedAt="{ item }">
                  <div class="d-flex align-center justify-center">
                    <v-icon size="16" color="grey" class="me-1">mdi-calendar</v-icon>
                    {{ item.requestedAt }}
                  </div>
                </template>

                <!-- 無資料顯示 -->
                <template v-slot:no-data>
                  <div class="text-center pa-5">
                    <v-icon size="40" color="grey-lighten-1" class="mb-3">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium text-grey-darken-1">
                      無待處理請求
                    </div>
                    <div class="text-body-2 text-grey mt-2">
                      您目前沒有待處理的授權請求
                    </div>
                  </div>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 歷史紀錄列表區塊 -->
      <v-row v-if="viewMode === 'history'" justify="center">
        <v-col cols="12">
          <!-- 返回按鈕 -->
          <v-btn
            color="grey-darken-1"
            variant="outlined"
            size="small"
            @click="switchToSearchView"
            class="mb-6"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="mb-6 elevation-2 rounded-xl">
            <v-card-title class="d-flex align-center py-4 px-6 bg-grey-lighten-4">
              <div class="d-flex align-center">
                <v-icon size="24" color="primary" class="me-3">mdi-history</v-icon>
                <span class="text-h6 font-weight-bold">授權請求歷史紀錄</span>
              </div>
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
                    title: '授權對象',
                    key: 'patientHash',
                    align: 'start',
                    width: '120px'
                  },
                  { 
                    title: '申請日期',
                    key: 'requestedAt',
                     align: 'center',
                    width: '120px'
                  },
                  { 
                    title: '狀態',
                    key: 'status',
                    align: 'center',
                    width: '100px'
                  }
                ]"
                :items="historyRequests"
                :loading="loadingMyRequests"
                loading-text="正在載入歷史紀錄..."
                class="elevation-0"
                hover
                density="comfortable"
                v-model:items-per-page="itemsPerPage"
                :items-per-page-options="[5, 10]"
              >
                <!-- 報告編號欄位 -->
                <template v-slot:item.reportId="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="primary" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.reportId }}</span>
                  </div>
                </template>

                <!-- 病患雜湊欄位 -->
                <template v-slot:item.patientHash="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="grey" class="me-2">mdi-identifier</v-icon>
                    <span class="font-weight-medium">{{ item.patientHash }}</span>
                  </div>
                </template>

                <!-- 申請日期欄位 -->
                <template v-slot:item.requestedAt="{ item }">
                  <div class="d-flex align-center justify-center">
                    <v-icon size="16" color="grey" class="me-1">mdi-calendar</v-icon>
                    {{ item.requestedAt }}
                  </div>
                </template>

                <!-- 狀態欄位 -->
                <template v-slot:item.status="{ item }">
                  <v-chip
                    size="small"
                    :color="getRequestStatusInfo(item.status).color"
                    variant="tonal"
                    class="font-weight-medium"
                  >
                    {{ getRequestStatusInfo(item.status).text }}
                  </v-chip>
                </template>

                <!-- 無資料顯示 -->
                <template v-slot:no-data>
                  <div class="text-center pa-5">
                    <v-icon size="40" color="grey-lighten-1" class="mb-3">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium text-grey-darken-1">
                      無歷史紀錄
                    </div>
                    <div class="text-body-2 text-grey mt-2">
                      您目前沒有已處理的授權請求紀錄
                    </div>
                  </div>
                </template>
              </v-data-table>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- Snackbar 訊息 -->
      <v-snackbar
        v-model="snackbar"
        :color="snackbarColor"
        :timeout="3000"
        location="top"
      >
        {{ snackbarMessage }}
        <template v-slot:actions>
          <v-btn
            variant="text"
            icon="mdi-close"
            @click="snackbar = false"
          ></v-btn>
        </template>
      </v-snackbar>
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
:deep(.v-data-table) {
  border-radius: 24px !important;
  overflow: hidden !important;
}

:deep(.v-data-table-header) {
  background-color: #F8FAFC !important;
}

:deep(.v-data-table-header th) {
  font-size: 1rem !important;
  color: #64748B !important;
  font-weight: 600 !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  padding: 14px 16px !important;
}

:deep(.v-data-table-row) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.05) !important;
  transition: background-color 0.2s ease;
}

:deep(.v-data-table-row:hover) {
  background-color: #F8FAFC !important;
}

:deep(.v-data-table td) {
  padding: 14px 16px !important;
  color: #334155 !important;
  font-size: 1rem !important;
}

/* 操作按鈕樣式 */
.view-report-btn {
  transition: all 0.3s ease !important;
  width: 36px !important;
  height: 36px !important;
  border-radius: 8px !important;
}

.view-report-btn:hover {
  transform: translateY(-2px) !important;
  background-color: rgba(var(--v-theme-primary), 0.1) !important;
}

.view-report-btn:active {
  transform: scale(0.95) !important;
}

/* 返回按鈕樣式 */
.back-btn {
  transition: all 0.3s ease !important;
}

.back-btn:hover {
  transform: translateX(-2px) !important;
}

/* 狀態標籤樣式 */
.status-chip {
  font-weight: 500 !important;
}

/* RWD 適配 */
@media (max-width: 960px) {
  .dashboard-container {
    padding: 1rem;
  }
  
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
}

@media (max-width: 600px) {
  .info-card {
    min-height: 180px;
  }
  
  .info-card .v-card-text {
    padding: 1rem !important;
  }
  
  :deep(.v-data-table-row) {
    display: flex !important;
    flex-direction: column !important;
    padding: 16px !important;
    margin-bottom: 8px !important;
    background: white !important;
    border-radius: 8px !important;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05) !important;
  }
  
  :deep(.v-data-table td) {
    display: flex !important;
    justify-content: space-between !important;
    align-items: center !important;
    padding: 8px 0 !important;
    border: none !important;
  }
}

/* 動畫效果 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

/* 搜尋區塊樣式 */
.search-card {
  background: white;
}

.search-input {
  border-radius: 12px;
}

:deep(.search-input .v-field) {
  border-radius: 12px !important;
  background-color: white !important;
  border: 1px solid #e5e7eb !important;
  transition: all 0.3s ease !important;
}

:deep(.search-input .v-field:hover) {
  border-color: rgba(var(--v-theme-primary), 0.5) !important;
}

:deep(.search-input .v-field:focus-within) {
  border-color: rgb(var(--v-theme-primary)) !important;
  box-shadow: 0 2px 12px rgba(var(--v-theme-primary), 0.15) !important;
}

.search-btn {
  border-radius: 12px !important;
  height: 44px;
}

.search-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(var(--v-theme-primary), 0.15) !important;
}
</style>