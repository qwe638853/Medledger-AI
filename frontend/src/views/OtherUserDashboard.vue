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
        id: report.id,
        patientName: report.patientName || '未知病患', // 病患真實姓名
        granted_at: report.date,
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
        patientName: report.patientName || '未知病患', // 病患真實姓名
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
    const result = await healthCheckService.requestReportAccess(
      selectedReport.value.id,
      patientId.value,
      authReason.value,
      authExpiry.value
    );

    if (result && result.success) {
    authRequestDialog.value = false;
      snackbarMessage.value = '已成功送出授權請求';
    snackbarColor.value = 'success';
    snackbar.value = true;
    
      // 授權請求成功後重新獲取報告和請求列表
      await Promise.all([
        fetchReportsByPatientId(patientId.value),
        fetchMyAccessRequests()
      ]);
      
      // 更新儀表板數據
      dashboardStats.value.pendingRequests = pendingRequests.value.length;
    } else {
      throw new Error(result.message || '授權請求失敗');
    }
  } catch (error) {
    console.error('發送授權請求時出錯:', error);
    snackbarMessage.value = error.message || '發送授權請求失敗';
    snackbarColor.value = 'error';
    snackbar.value = true;
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
  if (!expiry) return 'success';
  
  const expiryDate = new Date(expiry);
  const now = new Date();
  const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24));
  
  if (daysUntilExpiry < 0) return 'error';
  if (daysUntilExpiry <= 7) return 'warning';
  return 'success';
};

// 獲取授權到期 Chip 的文字顏色
const getExpiryTextColor = (expiry) => {
  if (!expiry) return 'white';
  
  const expiryDate = new Date(expiry);
  const now = new Date();
  const daysUntilExpiry = Math.ceil((expiryDate - now) / (1000 * 60 * 60 * 24));
  
  if (daysUntilExpiry < 0) return 'white';
  if (daysUntilExpiry <= 7) return 'white';
  return 'white';
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
          <v-card class="info-card user-info-card" elevation="0">
            <v-card-text class="pa-0">
              <div class="user-card-wrapper">
                <!-- 背景裝飾 -->
                <div class="user-card-bg-decoration"></div>
                
                <!-- 用戶資訊主體 -->
                <div class="user-card-main">
                  <div class="user-avatar-section">
                    <div class="user-avatar-container">
                      <v-icon class="user-avatar-icon">mdi-briefcase-account</v-icon>
                    </div>
                    <div class="user-status-indicator"></div>
                  </div>
                  
                  <div class="user-info-content">
                    <div class="user-name">{{ currentUser }}</div>
                    <div class="user-role">
                      <v-icon size="14" class="me-1">mdi-shield-account</v-icon>
                      保險業者
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

        <!-- 已授權報告統計卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card stat-card" elevation="0">
            <v-card-text class="pa-0">
              <div class="card-content">
                <div class="card-header-section">
                  <div class="icon-container">
                    <v-icon>mdi-file-document</v-icon>
                  </div>
                  
                  <div class="card-info">
                    <div class="stat-number">{{ dashboardStats.totalAuthorized }}</div>
                    <div class="stat-label">已授權報告</div>
                  </div>
                </div>
                
                <div class="card-action-section">
                  <v-btn
                    class="view-btn-simple"
                    @click="switchToAuthorizedView"
                    variant="text"
                    color="primary"
                    size="large"
                    block
                  >
                    <v-icon start size="30">mdi-eye</v-icon>
                    查看詳情
                  </v-btn>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 待處理請求卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card stat-card" elevation="0">
            <v-card-text class="pa-0">
              <div class="card-content">
                <div class="card-header-section">
                  <div class="icon-container">
                    <v-icon>mdi-clock-alert</v-icon>
                  </div>
                  
                  <div class="card-info">
                    <div class="stat-number">{{ dashboardStats.pendingRequests }}</div>
                    <div class="stat-label">待處理請求</div>
                  </div>
                </div>
                
                <div class="card-action-section">
                  <v-btn
                    class="view-btn-simple"
                    @click="switchToPendingView"
                    variant="text"
                    color="warning"
                    size="large"
                    block
                  >
                    <v-icon start size="30">mdi-eye</v-icon>
                    查看詳情
                  </v-btn>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 歷史紀錄卡片 -->
        <v-col cols="12" sm="6" md="3">
          <v-card class="info-card stat-card" elevation="0">
            <v-card-text class="pa-0">
              <div class="card-content">
                <div class="card-header-section">
                  <div class="icon-container">
                    <v-icon>mdi-history</v-icon>
                  </div>
                  
                  <div class="card-info">
                    <div class="stat-number">{{ historyRequests.length }}</div>
                    <div class="stat-label">歷史紀錄</div>
                  </div>
                </div>
                
                <div class="card-action-section">
                  <v-btn
                    class="view-btn-simple"
                    @click="switchToHistoryView"
                    variant="text"
                    color="info"
                    size="large"
                    block
                  >
                    <v-icon start size="30">mdi-eye</v-icon>
                    查看詳情
                  </v-btn>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 搜尋區塊 -->
      <v-row class="mb-8">
        <v-col cols="12">
          <v-card class="main-card search-card" elevation="0">
            <div class="card-header">
              <div class="card-title">
                <div class="card-title-icon">
                  <v-icon>mdi-magnify</v-icon>
                </div>
                <div>
                  <div class="card-title-text">病患健康報告查詢</div>
                </div>
              </div>
            </div>
            
            <div class="search-content">
              <v-row align="center" class="pa-6">
                <v-col cols="12" md="8">
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
                <v-col cols="12" md="4">
                  <v-btn
                    color="primary"
                    :loading="searchLoading"
                    :disabled="searchLoading"
                    @click="searchPatientReports"
                    prepend-icon="mdi-magnify"
                    elevation="0"
                    class="search-btn"
                    block
                  >
                    搜尋報告
                  </v-btn>
                </v-col>
              </v-row>
            </div>
          </v-card>
        </v-col>
      </v-row>

      <!-- 表格區域 -->
      <v-row v-if="showSearchResults && viewMode === 'search'" justify="center">
        <v-col cols="12">
          <v-card class="main-card result-card" elevation="0">
            <div class="card-header">
              <div class="card-title-row">
                <div class="card-title">
                  <div class="card-title-icon">
                    <v-icon>mdi-file-document-outline</v-icon>
                  </div>
                  <div>
                    <div class="card-title-text">病患「{{ patientId }}」的健康報告</div>
                  </div>
                </div>
                
                <!-- 標籤頁移到右上角 -->
                <v-tabs
                  v-model="tab"
                  color="primary"
                  align-tabs="end"
                  class="header-tabs"
                >
                  <v-tab value="unauthorized" class="font-weight-medium">
                    <v-icon start class="me-2" size="18">mdi-lock-open-alert</v-icon>
                    尚未授權報告
                    <v-badge
                      :content="unauthorizedReports.length.toString()"
                      :color="unauthorizedReports.length > 0 ? 'warning' : 'grey'"
                      offset-x="8"
                      offset-y="-8"
                      size="small"
                      class="ms-2"
                    ></v-badge>
                  </v-tab>
                  <v-tab value="authorized" class="font-weight-medium">
                    <v-icon start class="me-2" size="18">mdi-lock-check</v-icon>
                    已授權報告
                    <v-badge
                      :content="authorizedReports.length.toString()"
                      :color="authorizedReports.length > 0 ? 'success' : 'grey'"
                      offset-x="8"
                      offset-y="-8"
                      size="small"
                      class="ms-2"
                    ></v-badge>
                  </v-tab>
                </v-tabs>
              </div>
            </div>

            <!-- 移除原來的標籤頁位置 -->
            <div class="table-container">
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
                        <v-icon size="18" class="me-2">mdi-file-document</v-icon>
                        <span class="font-weight-medium">{{ item.id }}</span>
                    </div>
                  </template>

                    <template v-slot:item.content="{ item }">
                      <div class="text-truncate">{{ item.content }}</div>
                    </template>

                  <template v-slot:item.date="{ item }">
                      <div class="d-flex align-center justify-center">
                        <v-icon size="16" class="me-1">mdi-calendar</v-icon>
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
                      <div class="text-center pa-8">
                        <v-icon size="64" color="grey-lighten-2" class="mb-4">
                          mdi-file-document-outline
                        </v-icon>
                        <div class="text-subtitle-1 font-weight-medium mb-2">
                          無未授權報告
                        </div>
                        <div class="text-body-2">
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
                      { title: '授權對象', key: 'patientName', align: 'start', width: '120px' },
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
                        <v-icon size="18" class="me-2">mdi-file-document</v-icon>
                        <span class="font-weight-medium">{{ item.id }}</span>
                    </div>
            </template>

                    <template v-slot:item.patientName="{ item }">
                      <div class="d-flex align-center">
                        <v-icon size="18" class="me-2">mdi-account</v-icon>
                        <span class="font-weight-medium">{{ item.patientName || item.patientHash || '未知用戶' }}</span>
                      </div>
            </template>

                    <template v-slot:item.granted_at="{ item }">
                      <div class="d-flex align-center justify-center">
                        <v-icon size="16" class="me-1">mdi-calendar</v-icon>
                        {{ formatDate(item.granted_at) }}
                      </div>
                    </template>

                    <template v-slot:item.expiry="{ item }">
                      <div class="d-flex flex-column align-center">
                    <v-chip
                          size="small"
                          :color="getExpiryChipColor(item.expiry)"
                          variant="flat"
                          class="font-weight-bold text-white"
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
                      <div class="text-center pa-8">
                        <v-icon size="64" color="grey-lighten-2" class="mb-4">
                          mdi-file-document-outline
                        </v-icon>
                        <div class="text-subtitle-1 font-weight-medium mb-2">
                          尚無已授權報告
                        </div>
                        <div class="text-body-2">
                          目前您沒有任何已授權的健康報告
                        </div>
                    </div>
            </template>
          </v-data-table>
              </v-window-item>
            </v-window>
            </div>
        </v-card>
        </v-col>
      </v-row>

      <!-- 全部已授權報告區塊 -->
      <v-row v-if="viewMode === 'authorized'" justify="center">
        <v-col cols="12">
          <!-- 返回按鈕 -->
          <v-btn
            variant="outlined"
            size="small"
            @click="switchToSearchView"
            class="mb-6"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="main-card" elevation="0">
            <div class="card-header">
              <div class="card-title">
                <div class="card-title-icon">
                  <v-icon>mdi-folder-account</v-icon>
                </div>
                <div>
                  <div class="card-title-text">所有已授權健康報告</div>
                </div>
              </div>
            </div>
            
            <div class="table-container">
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
                    key: 'patientName',
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
                    <v-icon size="18" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.id }}</span>
                  </div>
              </template>
              
                                <!-- 授權對象欄位 -->
                <template v-slot:item.patientName="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" class="me-2">mdi-account</v-icon>
                    <span class="font-weight-medium">{{ item.patientName || item.patient_id || '未知病患' }}</span>
                  </div>
                </template>
              
                <!-- 授權到期日欄位 -->
              <template v-slot:item.expiry="{ item }">
                <v-chip
                  size="small"
                  :color="getExpiryChipColor(item.expiry)"
                    variant="flat"
                    class="font-weight-bold text-white"
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
                <div class="text-center pa-8">
                    <v-icon size="64" color="grey-lighten-2" class="mb-4">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium mb-2">
                      尚無已授權報告
                    </div>
                    <div class="text-body-2">
                      目前您沒有任何已授權的健康報告
                    </div>
                </div>
              </template>
            </v-data-table>
            </div>
          </v-card>
        </v-col>
      </v-row>

      <!-- 未搜尋或搜尋前的提示 -->
      <v-row v-if="!showSearchResults && viewMode === 'search'" justify="center">
        <v-col cols="12">
          <v-card class="main-card text-center" elevation="0">
            <div class="pa-12">
              <v-icon size="80" color="grey-lighten-2" class="mb-6">mdi-file-document-search</v-icon>
              <div class="text-h6 font-weight-medium mb-3">請輸入病患身分證字號查詢健康報告</div>
              <div class="text-body-1 text-grey">
                您可以輸入病患身分證字號查詢特定病患，或點擊上方統計卡片查看相關報告
              </div>
            </div>
          </v-card>
        </v-col>
      </v-row>

      <!-- 待授權請求列表區塊 -->
      <v-row v-if="viewMode === 'pending'" justify="center">
        <v-col cols="12">
          <!-- 返回按鈕 -->
          <v-btn
            variant="outlined"
            size="small"
            @click="switchToSearchView"
            class="mb-6"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="main-card" elevation="0">
            <div class="card-header">
              <div class="card-title">
                <div class="card-title-icon">
                  <v-icon>mdi-clock-alert</v-icon>
                </div>
                <div>
                  <div class="card-title-text">待處理授權請求</div>
                </div>
                <v-chip
                  v-if="pendingRequests.length > 0"
                  class="ms-4"
                  color="warning"
                  size="small"
                  variant="tonal"
                >
                  {{ pendingRequests.length }} 筆請求
                </v-chip>
              </div>
            </div>
          
            <div class="table-container">
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
                    key: 'patientName',
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
                    <v-icon size="18" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.reportId }}</span>
                  </div>
                </template>

                <!-- 病患姓名欄位 -->
                <template v-slot:item.patientName="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" class="me-2">mdi-account</v-icon>
                    <span class="font-weight-medium">{{ item.patientName || item.patientHash || '未知用戶' }}</span>
                  </div>
                </template>

                <!-- 申請日期欄位 -->
                <template v-slot:item.requestedAt="{ item }">
                  <div class="d-flex align-center justify-center">
                    <v-icon size="16" class="me-1">mdi-calendar</v-icon>
                    {{ item.requestedAt }}
                  </div>
                </template>

                <!-- 無資料顯示 -->
                <template v-slot:no-data>
                  <div class="text-center pa-8">
                    <v-icon size="64" color="grey-lighten-2" class="mb-4">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium mb-2">
                      無待處理請求
                    </div>
                    <div class="text-body-2">
                      您目前沒有待處理的授權請求
                    </div>
                  </div>
                </template>
              </v-data-table>
            </div>
          </v-card>
        </v-col>
      </v-row>

      <!-- 歷史紀錄列表區塊 -->
      <v-row v-if="viewMode === 'history'" justify="center">
        <v-col cols="12">
          <!-- 返回按鈕 -->
          <v-btn
            variant="outlined"
            size="small"
            @click="switchToSearchView"
            class="mb-6"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="main-card" elevation="0">
            <div class="card-header">
              <div class="card-title">
                <div class="card-title-icon">
                  <v-icon>mdi-history</v-icon>
                </div>
                <div>
                  <div class="card-title-text">授權請求歷史紀錄</div>
                </div>
              </div>
            </div>

            <div class="table-container">
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
                    key: 'patientName',
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
                    <v-icon size="18" class="me-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.reportId }}</span>
                  </div>
                </template>

                <!-- 病患雜湊欄位 -->
                <template v-slot:item.patientHash="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" class="me-2">mdi-identifier</v-icon>
                    <span class="font-weight-medium">{{ item.patientHash }}</span>
                  </div>
                </template>

                <!-- 申請日期欄位 -->
                <template v-slot:item.requestedAt="{ item }">
                  <div class="d-flex align-center justify-center">
                    <v-icon size="16" class="me-1">mdi-calendar</v-icon>
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
                  <div class="text-center pa-8">
                    <v-icon size="64" color="grey-lighten-2" class="mb-4">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-subtitle-1 font-weight-medium mb-2">
                      無歷史紀錄
                    </div>
                    <div class="text-body-2">
                      您目前沒有已處理的授權請求紀錄
                    </div>
                  </div>
                </template>
              </v-data-table>
            </div>
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

      <!-- 授權請求對話框 -->
      <v-dialog v-model="authRequestDialog" max-width="500">
        <v-card class="rounded-lg">
          <div class="card-header">
            <div class="card-title">
              <div class="card-title-icon">
                <v-icon>mdi-key-chain</v-icon>
              </div>
              <div>
                <div class="card-title-text">請求授權</div>
              </div>
            </div>
          </div>

          <v-card-text class="pa-6">
            <v-form @submit.prevent="sendAuthRequest">
              <div class="mb-4">
                <div class="text-subtitle-1 font-weight-medium mb-2">報告編號</div>
                <div class="text-body-1">{{ selectedReport?.id || '未知' }}</div>
              </div>

            <v-textarea
              v-model="authReason"
              label="授權理由"
                placeholder="請說明需要授權的原因..."
                :rules="[v => !!v || '請填寫授權理由']"
              rows="3"
                class="mb-4"
              hide-details="auto"
                variant="outlined"
                density="comfortable"
            ></v-textarea>

              <v-text-field
              v-model="authExpiry"
                label="授權到期日"
                type="date"
                :rules="[v => !!v || '請選擇到期日']"
                class="mb-4"
                hide-details="auto"
                variant="outlined"
                density="comfortable"
              ></v-text-field>
            </v-form>
          </v-card-text>
          
          <v-card-actions class="pa-6 pt-0">
            <v-spacer></v-spacer>
            <v-btn
              variant="outlined"
              @click="authRequestDialog = false"
              :disabled="requestLoading"
            >
              取消
            </v-btn>
            <v-btn
              color="primary"
              @click="sendAuthRequest"
              :loading="requestLoading"
              :disabled="requestLoading"
            >
              送出請求
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>
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
  min-height: 180px;
  padding: 0 !important;
  display: flex;
  flex-direction: column;
  position: relative;
}

.card-content.clickable {
  cursor: default !important;
}

.card-content.clickable:hover {
  background: transparent !important;
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
  width: 10px !important;
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
  font-size: 20px !important;
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

/* 搜尋內容區域 */
.search-content {
  background: white !important;
  padding: 0 !important;
}

/* 搜尋輸入框樣式 */
.search-input {
  border-radius: 16px !important;
}

:deep(.search-input .v-field) {
  border-radius: 16px !important;
  background-color: #f8fafc !important;
  border: 2px solid #e2e8f0 !important;
  transition: all 0.3s ease !important;
}

:deep(.search-input .v-field:hover) {
  border-color: rgba(0, 184, 217, 0.3) !important;
  background-color: white !important;
}

:deep(.search-input .v-field:focus-within) {
  border-color: #00B8D9 !important;
  background-color: white !important;
  box-shadow: 0 4px 16px rgba(0, 184, 217, 0.15) !important;
}

/* 搜尋按鈕樣式 */
.search-btn {
  border-radius: 16px !important;
  height: 56px !important;
  font-weight: 600 !important;
  font-size: 1.1rem !important;
  background: #00B8D9 !important;
  color: white !important;
  transition: all 0.3s ease !important;
}

.search-btn:hover {
  background: #0099B8 !important;
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.3) !important;
}

:deep(.v-table){
  font-size: 20px !important;
}
/* 表格樣式 */
:deep(.v-data-table) {
  border-radius: 0 !important;
  overflow: visible !important;
  box-shadow: none !important;
  background: white;
  border: none !important;
}

:deep(.v-data-table-header) {
  background: #f8fafc !important;
  border-bottom: 2px solid #e2e8f0 !important;
}

:deep(.v-data-table-header th) {
  font-size: 1.1rem !important;
  font-weight: 700 !important;
  color: #22292F !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  padding: 1.5rem !important;
  line-height: 1.4 !important;
  white-space: nowrap;
  background: transparent !important;
  vertical-align: middle;
  border-bottom: none !important;
  border-right: none !important;
  text-align: center !important;
}

:deep(.v-data-table tbody td) {
  font-size: 1rem !important;
  color: #22292F !important;
  padding: 1.5rem !important;
  line-height: 1.6 !important;
  vertical-align: middle;
  border-bottom: 1px solid #f1f5f9 !important;
  border-right: none !important;
  background: transparent !important;
}

:deep(.v-data-table-row) {
  height: auto !important;
  min-height: 80px !important;
  transition: all 0.3s ease;
  border: none !important;
}

:deep(.v-data-table-row:hover) {
  background: #f8fafc !important;
  transform: none !important;
  box-shadow: none !important;
}

/* 操作按鈕樣式 */
:deep(.view-report-btn) {
  border-radius: 12px !important;
  background: #f8fafc !important;
  border: 1px solid #e2e8f0 !important;
  transition: all 0.3s ease !important;
  width: 48px !important;
  height: 48px !important;
}

:deep(.view-report-btn:hover) {
  background: #00B8D9 !important;
  color: white !important;
  border-color: transparent !important;
  transform: scale(1.05) !important;
}

:deep(.view-report-btn .v-icon) {
  font-size: 22px !important;
  color: #8898AA !important;
}

:deep(.view-report-btn:hover .v-icon) {
  color: white !important;
}

/* 狀態標籤樣式 */
:deep(.v-chip) {
  border-radius: 16px !important;
  font-weight: 600 !important;
  padding: 0.6rem 1.2rem !important;
  font-size: 0.95rem !important;
  border: none !important;
}

/* 返回按鈕樣式 */
:deep(.v-btn[variant="outlined"]) {
  border-radius: 12px !important;
  border: 2px solid #e2e8f0 !important;
  color: #8898AA !important;
  font-weight: 600 !important;
  transition: all 0.3s ease !important;
}

:deep(.v-btn[variant="outlined"]:hover) {
  border-color: #00B8D9 !important;
  color: #00B8D9 !important;
  background: rgba(0, 184, 217, 0.05) !important;
  transform: translateX(-2px) !important;
}

/* 無資料狀態 */
:deep(.v-data-table__empty-wrapper) {
  padding: 4rem 2rem !important;
  text-align: center;
}

:deep(.v-data-table__empty-wrapper .v-icon) {
  font-size: 64px !important;
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

/* 表格內圖標顏色 */
:deep(.v-data-table .d-flex.align-center .v-icon) {
  color: #00B8D9 !important;
  margin-right: 0.5rem !important;
}

:deep(.v-data-table .v-avatar) {
  background: rgba(0, 184, 217, 0.1) !important;
  border: 1px solid rgba(0, 184, 217, 0.2) !important;
}

:deep(.v-data-table .v-avatar .v-icon) {
  color: #00B8D9 !important;
}

/* 標籤頁樣式 */
:deep(.v-tabs) {
  background: white !important;
}

:deep(.v-tab) {
  color: #8898AA !important;
  font-weight: 600 !important;
  text-transform: none !important;
  font-size: 1rem !important;
}

:deep(.v-tab--selected) {
  color: #00B8D9 !important;
}

:deep(.v-tabs-slider) {
  background: #00B8D9 !important;
  height: 3px !important;
}

/* 對話框樣式 */
:deep(.v-dialog .v-card) {
  border-radius: 24px !important;
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.15) !important;
}

:deep(.v-dialog .v-card-title) {
  background: #f8fafc !important;
  color: #22292F !important;
  font-weight: 700 !important;
}

/* Snackbar 樣式 */
:deep(.v-snackbar) {
  border-radius: 16px !important;
  font-weight: 600 !important;
}

/* 表格容器樣式 */
.table-container {
  background: white !important;
  padding: 0 !important;
}

/* 響應式設計 */
@media (max-width: 1200px) {
  .dashboard-container {
    padding: 1.5rem !important;
  }

  .card-content {
    min-height: 160px;
  }

  .card-header-section {
    padding: 1.75rem !important;
  }

  .card-action-section {
    padding: 0 1.25rem 1.25rem 1.25rem !important;
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

  :deep(.v-data-table-header th) {
    font-size: 1rem !important;
    padding: 1.25rem !important;
  }

  :deep(.v-data-table tbody td) {
    font-size: 0.95rem !important;
    padding: 1.25rem !important;
  }

  /* 響應式卡片標題列 */
  .card-title-row {
    gap: 1.5rem;
  }

  .header-tabs :deep(.v-tab) {
    font-size: 0.9rem !important;
    padding: 0 1rem !important;
  }
}

@media (max-width: 960px) {
  .dashboard-container {
    padding: 1rem !important;
  }

  .info-card {
    min-height: 160px;
  }

  .card-content {
    min-height: 160px;
  }

  .card-header-section {
    padding: 1.5rem !important;
    gap: 1.25rem;
  }

  .card-action-section {
    padding: 0 1rem 1rem 1rem !important;
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

  .card-header {
    padding: 1.5rem !important;
  }

  :deep(.v-data-table-header th) {
    font-size: 0.95rem !important;
    padding: 1rem !important;
  }

  :deep(.v-data-table tbody td) {
    font-size: 0.9rem !important;
    padding: 1rem !important;
  }

  /* 響應式卡片標題列 */
  .card-title-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .header-tabs {
    align-self: stretch;
  }

  .header-tabs :deep(.v-tabs) {
    width: 100%;
  }

  .header-tabs :deep(.v-tab) {
    font-size: 0.85rem !important;
    padding: 0 0.75rem !important;
    flex: 1;
  }
}

@media (max-width: 600px) {
  .dashboard-container {
    padding: 0.75rem !important;
  }

  .info-card {
    min-height: 140px;
  }

  .card-content {
    min-height: 140px;
  }

  .card-header-section {
    padding: 1.25rem !important;
    gap: 1rem;
  }

  .card-action-section {
    padding: 0 0.75rem 0.75rem 0.75rem !important;
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

  .card-header {
    padding: 1.5rem !important;
  }

  .view-btn-new {
    height: 32px !important;
    font-size: 0.8rem !important;
    margin-top: 0.75rem !important;
  }

  :deep(.v-data-table td),
  :deep(.v-data-table-header th) {
    padding: 0.75rem 0.5rem !important;
    font-size: 0.85rem !important;
  }

  /* 響應式卡片標題列 */
  .card-title-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .card-title-text {
    font-size: 1.2rem !important;
  }

  .header-tabs :deep(.v-tab) {
    font-size: 0.8rem !important;
    padding: 0 0.5rem !important;
    height: 40px !important;
  }

  .header-tabs :deep(.v-badge) {
    margin-left: 0.25rem !important;
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

/* 統計卡片查看按鈕 */
.view-btn {
  position: absolute !important;
  top: 1.5rem !important;
  right: 1.5rem !important;
  width: 36px !important;
  height: 36px !important;
  border-radius: 12px !important;
  transition: all 0.3s ease !important;
  opacity: 0.8 !important;
}

.view-btn:hover {
  opacity: 1 !important;
  transform: scale(1.1) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.view-btn .v-icon {
  font-size: 18px !important;
}

/* 卡片頭部區域 */
.card-header-section {
  flex: 1;
  padding: 2rem !important;
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

/* 卡片操作區域 */
.card-action-section {
  padding: 0 1.5rem 1.5rem 1.5rem !important;
  border-top: none !important;
  background: transparent !important;
  display: flex;
  justify-content: center;
  align-items: center;
}

/* 新的查看按鈕樣式 */
.view-btn-new {
  border-radius: 12px !important;
  font-weight: 600 !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
  height: 48px !important;
  margin-top: 1rem !important;
  font-size: 1.1rem !important;
  background: transparent !important;
  border: none !important;
  text-align: center !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
}

.view-btn-new:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15) !important;
}

.view-btn-new .v-icon {
  margin-right: 0.5rem !important;
}

/* 移除舊的查看按鈕樣式 */
.view-btn {
  display: none !important;
}

/* 用戶卡片獨立樣式 */
.user-card-custom {
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

.user-card-custom:hover {
  transform: translateY(-4px);
  box-shadow: 0 16px 48px rgba(0, 184, 217, 0.15) !important;
  border-color: rgba(0, 184, 217, 0.2) !important;
}

.user-card-custom .v-card-text {
  height: 100%;
  padding: 0 !important;
}

/* 用戶卡片內容布局 */
.user-card-content {
  height: 100%;
  min-height: 160px;
  padding: 2rem !important;
  display: flex;
  flex-direction: column;
  justify-content: center;
  position: relative;
}

/* 用戶資訊區域 - 水平布局 */
.user-info-section {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  flex: 1;
}

/* 用戶圖標容器 */
.user-icon-container {
  width: 64px;
  height: 64px;
  border-radius: 18px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(0, 184, 217, 0.1) !important;
  border: 2px solid rgba(0, 184, 217, 0.2) !important;
  flex-shrink: 0;
}

.user-icon-container .v-icon {
  font-size: 32px !important;
  color: #00B8D9 !important;
}

/* 用戶文字資訊 */
.user-text-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  text-align: center;
}

.user-card-title {
  font-size: 1.5rem !important;
  font-weight: 700 !important;
  color: #22292F !important;
  line-height: 1.2;
  margin: 0;
}

.user-card-subtitle {
  font-size: 1.1rem !important;
  font-weight: 500 !important;
  color: #8898AA !important;
  margin: 0;
}

/* 響應式設計 */
@media (max-width: 768px) {
  .user-card-content {
    padding: 1.5rem !important;
    min-height: 140px;
  }
  
  .user-info-section {
    gap: 1rem;
  }
  
  .user-icon-container {
    width: 56px;
    height: 56px;
    border-radius: 16px;
  }
  
  .user-icon-container .v-icon {
    font-size: 28px !important;
  }
  
  .user-card-title {
    font-size: 1.3rem !important;
  }
  
  .user-card-subtitle {
    font-size: 1rem !important;
  }
}

@media (max-width: 600px) {
  .user-card-content {
    padding: 1.25rem !important;
    min-height: 120px;
  }
  
  .user-info-section {
    gap: 0.75rem;
  }
  
  .user-icon-container {
    width: 48px;
    height: 48px;
    border-radius: 14px;
  }
  
  .user-icon-container .v-icon {
    font-size: 24px !important;
  }
  
  .user-card-title {
    font-size: 1.2rem !important;
  }
  
  .user-card-subtitle {
    font-size: 0.95rem !important;
  }
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

/* 卡片標題列樣式 */
.card-title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  width: 100%;
  gap: 2rem;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  margin: 0 !important;
  flex-shrink: 0;
}

/* 頭部標籤頁樣式 */
.header-tabs {
  flex-shrink: 0;
  min-width: auto;
}

.header-tabs :deep(.v-tabs-slider) {
  background: #00B8D9 !important;
  height: 3px !important;
}

.header-tabs :deep(.v-tab) {
  color: #8898AA !important;
  font-weight: 600 !important;
  text-transform: none !important;
  font-size: 1rem !important;
  min-width: auto !important;
  padding: 0 1.5rem !important;
  height: 48px !important;
}

.header-tabs :deep(.v-tab--selected) {
  color: #00B8D9 !important;
}

.header-tabs :deep(.v-badge) {
  margin-left: 0.5rem !important;
}

/* 簡潔的查看按鈕樣式 */
.view-btn-simple {
  border-radius: 8px !important;
  font-weight: 600 !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  transition: all 0.3s ease !important;
  height: 40px !important;
  margin-top: 1rem !important;
  font-size: 1.3rem !important;
  background: transparent !important;
  border: none !important;
  text-align: center !important;
  display: flex !important;
  align-items: center !important;
  justify-content: center !important;
}

.view-btn-simple:hover {
  background: rgba(0, 184, 217, 0.08) !important;
  transform: translateY(-1px) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

/* 卡片操作區域調整 */
.card-action-section {
  padding: 0 1.5rem 1.5rem 1.5rem !important;
  border-top: none !important;
  background: transparent !important;
  display: flex;
  justify-content: center;
  align-items: center;
}
</style>