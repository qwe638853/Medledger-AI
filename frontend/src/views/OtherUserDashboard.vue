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
const pendingItemsPerPage = ref(5);
const historyItemsPerPage = ref(5);

const fetchDashboardStats = async () => {
  try {
    //dashboardStats.value = await healthCheckService.fetchDashboardStats();
  } catch (error) {
    console.error('獲取 Dashboard 數據時出錯:', error);
    // 錯誤處理已經在服務層完成
  }
};


// 獲取所有已授權報告
const fetchAllAuthorizedReports = async () => {
  loadingAuthorizedReports.value = true;
  try {
    allAuthorizedReports.value = await healthCheckService.fetchAuthorizedReports();
    console.log("tag",allAuthorizedReports.value);
    
    dashboardStats.value.totalAuthorized = allAuthorizedReports.value.length;
    console.log("tag",dashboardStats.value);
  } catch (error) {
    console.error('獲取已授權報告時出錯:', error);
    snackbarMessage.value = '獲取已授權報告時出錯';
    snackbarColor.value = 'error';
    snackbar.value = true;
    // 錯誤處理已經在服務層完成
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

// 查看內容按鈕的事件改為：
const goToReportDetail = (item) => {
  // 確保 item.report_id 和 item.patient_id 有值，若來源 key 為 id/patientId 也一併處理
  console.log(item);
  const report_id = item.id;
  const patient_id = item.patient_id ;
  router.push({ name: 'ReportDetail', params: { report_id, patient_id } });
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
    <v-container class="dashboard-container py-8 mx-auto" max-width="1800">
      <!-- 四張統計卡片置於同一行 -->
      <v-row class="mb-8" justify="center">
        <!-- 歡迎卡片 -->
        <v-col cols="12" sm="6" md="3" class="px-3">
          <v-card class="rounded-xl welcome-card" elevation="2" height="100">
            <v-card-text class="d-flex align-center pa-4">
              <div class="rounded-circle bg-yellow-lighten-4 p-3 me-3">
                <v-icon size="32" color="grey-darken-3">mdi-account-outline</v-icon>
              </div>
              <div class="flex-grow-1">
                <div class="text-subtitle-1 text-grey-darken-1">歡迎，保險顧問</div>
                <div class="text-h6 font-weight-bold text-grey-darken-3">{{ currentUser }}</div>
              </div>
              <v-btn 
                color="grey-darken-3" 
                @click="handleLogout" 
                elevation="1" 
                icon
                size="large"
                class="modern-btn d-flex align-center logout-btn"
              >
                <v-icon size="28" class="me-2">mdi-logout-variant</v-icon>
                <span class="logout-text">登出</span>
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 已授權報告數量卡片 -->
        <v-col cols="12" sm="6" md="3" class="px-3">
          <v-card class="stat-card rounded-xl" elevation="2" height="100">
            <v-card-text class="pa-4 d-flex align-center justify-space-between">
              <div class="d-flex align-center">
                <div class="rounded-circle bg-yellow-lighten-4 p-2 me-3">
                  <v-icon size="24" color="grey-darken-3">mdi-file-document-outline</v-icon>
                </div>
                <div>
                  <div class="text-overline text-grey-darken-1">已授權報告</div>
                  <div class="text-h4 font-weight-bold text-grey-darken-3">{{ dashboardStats.totalAuthorized }}</div>
                </div>
              </div>
              <v-btn
                color="grey-darken-3"
                class="modern-btn d-flex align-center view-all-btn"
                size="large"
                elevation="1"
                style="min-width: 90px; height: 48px; font-size: 1.15rem; font-weight: 600;"
                @click="switchToAuthorizedView"
              >
                <v-icon size="24" class="me-2">mdi-eye</v-icon>
                <span class="logout-text">查看</span>
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 待處理請求卡片 -->
        <v-col cols="12" sm="6" md="3" class="px-3">
          <v-card class="stat-card rounded-xl" elevation="2" height="100" @click="switchToPendingView">
            <v-card-text class="pa-4 d-flex align-center justify-space-between">
              <div class="d-flex align-center">
                <div class="rounded-circle bg-orange-lighten-5 p-2 me-3">
                  <v-icon size="24" color="orange-darken-2">mdi-clock-outline</v-icon>
                </div>
                <div>
                  <div class="text-overline text-orange-darken-1">待處理請求</div>
                  <div class="text-h4 font-weight-bold text-orange-darken-3">{{ dashboardStats.pendingRequests }}</div>
                </div>
              </div>
              <v-btn
                color="orange-darken-2"
                class="modern-btn d-flex align-center view-all-btn"
                size="large"
                elevation="1"
                style="min-width: 90px; height: 48px; font-size: 1.15rem; font-weight: 600;"
                @click.stop="switchToPendingView"
              >
                <v-icon size="24" class="me-2">mdi-eye</v-icon>
                <span class="logout-text">查看</span>
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 歷史紀錄卡片 -->
        <v-col cols="12" sm="6" md="3" class="px-3">
          <v-card class="stat-card rounded-xl" elevation="2" height="100" @click="switchToHistoryView">
            <v-card-text class="pa-4 d-flex align-center justify-space-between">
              <div class="d-flex align-center">
                <div class="rounded-circle bg-blue-grey-lighten-5 p-2 me-3">
                  <v-icon size="24" color="blue-grey-darken-1">mdi-history</v-icon>
                </div>
                <div>
                  <div class="text-overline text-blue-grey-darken-1">歷史紀錄</div>
                  <div class="text-h6 font-weight-medium text-blue-grey-darken-2">查看已處理請求</div>
                </div>
              </div>
              <v-btn
                color="blue-grey-darken-1"
                class="modern-btn d-flex align-center view-all-btn"
                size="large"
                elevation="1"
                style="min-width: 90px; height: 48px; font-size: 1.15rem; font-weight: 600;"
                @click.stop="switchToHistoryView"
              >
                <v-icon size="24" class="me-2">mdi-eye</v-icon>
                <span class="logout-text">查看</span>
              </v-btn>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 搜尋區塊 -->
      <v-row class="mb-8" justify="center">
        <v-col cols="12">
          <v-card class="rounded-xl search-card" elevation="3">
            <v-card-text class="pa-6">
              <v-row align="center">
                <v-col cols="12" md="4" class="pb-md-0">
                  <div class="d-flex align-center mb-md-0 mb-3">
                    <v-icon size="24" color="grey-darken-3" class="me-2">mdi-magnify</v-icon>
                    <span class="text-h6 font-weight-bold text-grey-darken-3">病患健康報告查詢</span>
                  </div>
                </v-col>
                <v-col cols="12" md="5" class="pb-md-0 pt-md-0">
                  <v-text-field
                    v-model="patientId"
                    label="請輸入病患身分證字號"
                    placeholder="例如：A123456789"
                    variant="outlined"
                    density="comfortable"
                    hide-details
                    prepend-inner-icon="mdi-card-account-details-outline"
                    class="modern-input"
                  ></v-text-field>
                </v-col>
                <v-col cols="12" md="3" class="pt-md-0 d-flex gap-3">
                  <v-btn
                    color="yellow-accent-3"
                    :loading="searchLoading"
                    :disabled="searchLoading"
                    @click="searchPatientReports"
                    prepend-icon="mdi-magnify"
                    elevation="2"
                    class="flex-grow-1 modern-btn primary-btn"
                    :class="{'v-btn--active': viewMode === 'search'}"
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
          <v-card class="rounded-xl result-card" elevation="3">
            <v-card-title class="py-4 px-6 bg-grey-lighten-4">
              <div class="d-flex align-center">
                <v-icon size="24" color="grey-darken-3" class="me-3">mdi-file-document-outline</v-icon>
                <span class="text-h6 font-weight-bold text-grey-darken-3">病患「{{ patientId }}」的健康報告</span>
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
            
            <!-- 間隔條 -->
            <div class="tab-content-divider"></div>
            
            <v-window v-model="tab" class="result-window">
              <!-- 尚未授權報告標籤內容 -->
              <v-window-item value="unauthorized" class="pb-4 pt-2">
                <v-data-table
                  :headers="[
                    { title: '報告編號', key: 'id', align: 'start', sortable: true, width: '90px' },
                    { title: '健康數據', key: 'content', align: 'start', width: '45%' },
                    { title: '日期', key: 'date', align: 'center', sortable: true, width: '110px' },
                    { title: '狀態', key: 'status', align: 'center', width: '90px' },
                    { title: '操作', key: 'actions', align: 'center', width: '110px', fixed: true }
                  ]"
                  :items="unauthorizedReports"
                  :loading="loading"
                  loading-text="資料載入中..."
                  class="elevation-0 report-table"
                  hover
                  item-value="id"
                  density="compact"
                  fixed-header
                >
                  <template v-slot:item.content="{ item }">
                    <div class="text-truncate content-cell">
                      {{ item.content }}
                    </div>
                  </template>
                  <template v-slot:item.date="{ item }">
                    {{ item.date || '-' }}
                  </template>
                  <template v-slot:item.status="{ item }">
                    <v-chip
                      size="x-small"
                      color="orange-lighten-4"
                      text-color="orange-darken-3"
                      variant="outlined"
                      prepend-icon="mdi-lock-open-alert"
                      class="font-weight-medium status-chip"
                    >
                      未授權
                    </v-chip>
                  </template>
                  <template v-slot:item.actions="{ item }">
                    <div class="actions-cell">
                      <v-btn 
                        color="deep-purple" 
                        variant="flat"
                        size="small"
                        @click="openAuthRequestDialog(item)"
                        prepend-icon="mdi-key-chain"
                        rounded
                        class="action-btn"
                      >
                        <span class="d-none d-sm-inline">請求授權</span>
                        <span class="d-sm-none">授權</span>
                      </v-btn>
                    </div>
                  </template>
                  <template v-slot:no-data>
                    <div class="text-center pa-4">
                      <v-icon size="40" color="blue-grey-lighten-2" class="mb-2">mdi-lock-open-outline</v-icon>
                      <div class="text-subtitle-1 font-weight-medium text-blue-grey">無未授權報告</div>
                      <div class="text-body-2 text-grey">該病患沒有未授權的健康報告</div>
                    </div>
                  </template>
                </v-data-table>
              </v-window-item>
              
              <!-- 已授權報告標籤內容 -->
              <v-window-item value="authorized" class="pb-4 pt-2">
          <v-data-table
            :headers="[
                    { title: '報告編號', key: 'id', align: 'start', sortable: true, width: '90px' },
                    { title: '健康數據', key: 'content', align: 'start', width: '45%' },
                    { title: '日期', key: 'date', align: 'center', sortable: true, width: '110px' },
                    { title: '狀態', key: 'status', align: 'center', width: '90px' },
                    { title: '授權到期', key: 'expiry', align: 'center', width: '110px' }
            ]"
                  :items="authorizedReports"
            :loading="loading"
            loading-text="資料載入中..."
                  class="elevation-0 report-table"
                  hover
                  item-value="id"
                  density="compact"
                  fixed-header
          >
                  <template v-slot:item.content="{ item }">
                    <div class="text-truncate content-cell">
                      {{ item.content }}
                    </div>
            </template>
                  <template v-slot:item.date="{ item }">
              {{ item.date || '-' }}
            </template>
                  <template v-slot:item.status="{ item }">
                    <v-chip
                      size="x-small"
                      color="green-lighten-4"
                      text-color="green-darken-3"
                      variant="outlined"
                      prepend-icon="mdi-lock-check"
                      class="font-weight-medium status-chip"
                    >
                      已授權
                    </v-chip>
                  </template>
                  <template v-slot:item.expiry="{ item }">
                    {{ item.expiry || '永久' }}
                  </template>
                  <template v-slot:no-data>
                    <div class="text-center pa-4">
                      <v-icon size="40" color="blue-grey-lighten-2" class="mb-2">mdi-lock-check-outline</v-icon>
                      <div class="text-subtitle-1 font-weight-medium text-blue-grey">無已授權報告</div>
                      <div class="text-body-2 text-grey">該病患沒有已授權的健康報告</div>
                    </div>
            </template>
          </v-data-table>
              </v-window-item>
            </v-window>
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
            class="mb-6 back-btn"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="authorized-reports-card" elevation="0">
            <!-- 標題區塊 -->
            <div class="auth-reports-header px-6 py-4 d-flex align-center">
              <v-icon
                size="28"
                color="#24B47E"
                class="me-3"
              >mdi-folder-account</v-icon>
              <span class="text-h5 font-weight-bold">所有已授權健康報告</span>
            </div>
            
            <!-- 表格區塊 -->
            <v-data-table
              :headers="[
                { 
                  title: '報告編號',
                  key: 'id',
                  align: 'start',
                  width: '120px'
                },
                { 
                  title: '病患 ID',
                  key: 'patient_id',
                  align: 'start',
                  width: '150px'
                },
                { 
                  title: '報告日期',
                  key: 'date',
                  align: 'center',
                  width: '140px'
                },
                { 
                  title: '授權到期',
                  key: 'expiry',
                  align: 'center',
                  width: '140px'
                },
                { 
                  title: '',
                  key: 'actions',
                  align: 'center',
                  width: '80px',
                  sortable: false
                }
              ]"
              :items="allAuthorizedReports"
              :loading="loadingAuthorizedReports"
              loading-text="正在載入已授權報告..."
              class="authorized-reports-table"
              hover
              v-model:items-per-page="itemsPerPage"
              :items-per-page-options="[10, 20, 50]"
            >
              <!-- 報告編號欄位 -->
              <template v-slot:item.id="{ item }">
                <div class="id-cell">
                  {{ item.id.substring(0, 4) }}...{{ item.id.slice(-4) }}
                  <div class="id-tooltip">{{ item.id }}</div>
                </div>
              </template>

              <!-- 病患 ID 欄位 -->
              <template v-slot:item.patient_id="{ item }">
                <div class="id-cell">
                  {{ item.patient_id.substring(0, 4) }}...{{ item.patient_id.slice(-4) }}
                  <div class="id-tooltip">{{ item.patient_id }}</div>
                </div>
              </template>
              
              <!-- 報告日期欄位 -->
              <template v-slot:item.date="{ item }">
                <div class="date-cell">
                  <v-icon
                    size="16"
                    color="grey-darken-1"
                    class="me-1"
                  >mdi-calendar-outline</v-icon>
                  {{ formatDate(item.date) }}
                </div>
              </template>
              
              <!-- 授權到期欄位 -->
              <template v-slot:item.expiry="{ item }">
                <v-chip
                  size="small"
                  :color="getExpiryChipColor(item.expiry)"
                  :text-color="getExpiryTextColor(item.expiry)"
                  variant="outlined"
                  class="expiry-chip"
                >
                  {{ formatExpiryDate(item.expiry) }}
                </v-chip>
              </template>
              
              <!-- 操作按鈕欄位 -->
              <template v-slot:item.actions="{ item }">
                <v-btn
                  icon
                  variant="text"
                  size="small"
                  @click="goToReportDetail(item)"
                  class="view-report-btn"
                >
                  <v-icon>mdi-eye-outline</v-icon>
                </v-btn>
              </template>

              <!-- 無資料顯示 -->
              <template v-slot:no-data>
                <div class="text-center pa-5">
                  <v-icon size="48" color="grey-lighten-1" class="mb-3">mdi-folder-open-outline</v-icon>
                  <div class="text-h6 font-weight-medium text-grey-darken-1 mb-1">尚無已授權報告</div>
                  <div class="text-body-2 text-grey">目前您沒有任何已授權的健康報告</div>
                </div>
              </template>
            </v-data-table>
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
            class="mb-6 back-btn"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="pending-requests-card" elevation="0">
            <!-- 標題區塊 -->
            <div class="pending-header px-6 py-4 d-flex align-center">
              <v-icon
                size="28"
                color="orange-darken-2"
                class="me-3"
              >mdi-clock-alert</v-icon>
              <span class="text-h6 font-weight-bold">待處理授權請求</span>
              <v-chip
                class="ms-3"
                color="orange-lighten-4"
                text-color="orange-darken-2"
                size="small"
              >
                {{ pendingRequests.length }} 筆請求
              </v-chip>
            </div>
            
            <!-- 表格區塊 -->
            <v-data-table
              :headers="[
                { 
                  title: '報告編號',
                  key: 'reportId',
                  align: 'start',
                  width: '120px'
                },
                { 
                  title: '病患雜湊',
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
              class="pending-requests-table"
              hover
              v-model:items-per-page="pendingItemsPerPage"
              :items-per-page-options="[5, 10]"
            >
              <!-- 報告編號欄位 -->
              <template v-slot:item.reportId="{ item }">
                <div class="id-cell">
                  {{ item.reportId.substring(0, 4) }}...{{ item.reportId.slice(-4) }}
                  <div class="id-tooltip">{{ item.reportId }}</div>
                </div>
              </template>

              <!-- 病患雜湊欄位 -->
              <template v-slot:item.patientHash="{ item }">
                <div class="id-cell">
                  {{ item.patientHash.substring(0, 4) }}...{{ item.patientHash.slice(-4) }}
                  <div class="id-tooltip">{{ item.patientHash }}</div>
                </div>
              </template>

              <!-- 申請日期欄位 -->
              <template v-slot:item.requestedAt="{ item }">
                <div class="date-cell">
                  <v-icon size="16" color="grey-darken-1" class="me-1">mdi-calendar-outline</v-icon>
                  {{ item.requestedAt }}
                </div>
              </template>

              <!-- 無資料顯示 -->
              <template v-slot:no-data>
                <div class="text-center pa-5">
                  <v-icon size="40" color="grey-lighten-1" class="mb-3">mdi-clipboard-text-clock-outline</v-icon>
                  <div class="text-subtitle-1 font-weight-medium text-grey-darken-1">無待處理請求</div>
                  <div class="text-body-2 text-grey">您目前沒有待處理的授權請求</div>
                </div>
              </template>
            </v-data-table>
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
            class="mb-6 back-btn"
            elevation="0"
          >
            <v-icon start size="18">mdi-arrow-left</v-icon>
            返回搜尋
          </v-btn>

          <v-card class="history-requests-card" elevation="0">
            <!-- 標題區塊 -->
            <div class="history-header px-6 py-4 d-flex align-center">
              <v-icon
                size="28"
                color="blue-grey-darken-1"
                class="me-3"
              >mdi-history</v-icon>
              <span class="text-h6 font-weight-bold">授權請求歷史紀錄</span>
            </div>
            
            <!-- 表格區塊 -->
            <v-data-table
              :headers="[
                { 
                  title: '報告編號',
                  key: 'reportId',
                  align: 'start',
                  width: '120px'
                },
                { 
                  title: '病患雜湊',
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
              class="history-requests-table"
              hover
              v-model:items-per-page="historyItemsPerPage"
              :items-per-page-options="[5, 10]"
            >
              <!-- 報告編號欄位 -->
              <template v-slot:item.reportId="{ item }">
                <div class="id-cell">
                  {{ item.reportId.substring(0, 4) }}...{{ item.reportId.slice(-4) }}
                  <div class="id-tooltip">{{ item.reportId }}</div>
                </div>
              </template>

              <!-- 病患雜湊欄位 -->
              <template v-slot:item.patientHash="{ item }">
                <div class="id-cell">
                  {{ item.patientHash.substring(0, 4) }}...{{ item.patientHash.slice(-4) }}
                  <div class="id-tooltip">{{ item.patientHash }}</div>
                </div>
              </template>

              <!-- 申請日期欄位 -->
              <template v-slot:item.requestedAt="{ item }">
                <div class="date-cell">
                  <v-icon size="16" color="grey-darken-1" class="me-1">mdi-calendar-outline</v-icon>
                  {{ item.requestedAt }}
                </div>
              </template>

              <!-- 狀態欄位 -->
              <template v-slot:item.status="{ item }">
                <v-chip
                  size="small"
                  :color="getRequestStatusInfo(item.status).color + '-lighten-4'"
                  :text-color="getRequestStatusInfo(item.status).color + '-darken-2'"
                  variant="outlined"
                  :prepend-icon="getRequestStatusInfo(item.status).icon"
                  class="status-chip"
                >
                  {{ getRequestStatusInfo(item.status).text }}
                </v-chip>
              </template>

              <!-- 無資料顯示 -->
              <template v-slot:no-data>
                <div class="text-center pa-5">
                  <v-icon size="40" color="grey-lighten-1" class="mb-3">mdi-history</v-icon>
                  <div class="text-subtitle-1 font-weight-medium text-grey-darken-1">無歷史紀錄</div>
                  <div class="text-body-2 text-grey">您目前沒有已處理的授權請求紀錄</div>
                </div>
              </template>
            </v-data-table>
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
.dashboard-bg {
  background-color: #F9F7F4;
  min-height: 100vh;
  display: flex;
  align-items: flex-start;
  justify-content: center;
}

.dashboard-container {
  padding: 2rem;
  background: var(--background-color);
  min-height: 100vh;
  max-width: 1800px;
  margin: 0 auto;
}

.dashboard-header {
  margin-bottom: 2rem;
}

.dashboard-title {
  font-family: 'Inter', sans-serif;
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-color);
  margin-bottom: 1rem;
  letter-spacing: -0.5px;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.modern-btn {
  border-radius: 24px !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  font-weight: 600 !important;
  transition: all 0.3s ease !important;
}

.modern-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
  opacity: 0.95;
}

.primary-btn {
  background-color: #F8F441 !important;
  color: #333 !important;
}

.modern-input {
  border-radius: 12px !important;
  background-color: white !important;
}

.modern-input :deep(.v-field) {
  border-radius: 12px !important;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05) !important;
  border: 1px solid #e5e7eb !important;
}

.modern-input :deep(.v-field:focus-within) {
  border-color: #F8F441 !important;
  box-shadow: 0 2px 8px rgba(248, 244, 65, 0.2) !important;
}

.search-card, .result-card {
  background: white;
  border-radius: 24px !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.05) !important;
  overflow: hidden;
}

.welcome-card {
  background: linear-gradient(135deg, #fff 0%, #F9F7F4 100%);
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.stat-card {
  background: white;
  transition: all 0.3s ease;
  border: 1px solid rgba(0, 0, 0, 0.05);
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1) !important;
}

.search-section {
  background: var(--white);
  border-radius: var(--border-radius-lg);
  padding: 2rem;
  box-shadow: var(--shadow-md);
  margin-bottom: 2rem;
  border: 1px solid var(--border-color);
}

.search-form {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.search-input {
  flex: 1;
  padding: 0.75rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  font-size: 1rem;
  transition: all 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.view-toggle {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.view-button {
  padding: 0.75rem 1.5rem;
  border-radius: var(--border-radius-md);
  font-weight: 500;
  color: var(--muted-color);
  background: var(--background-color);
  border: 1px solid var(--border-color);
  transition: all 0.2s ease;
}

.view-button.active {
  background: var(--primary-color);
  color: var(--white);
  border-color: var(--primary-color);
}

.reports-section {
  background: var(--white);
  border-radius: var(--border-radius-lg);
  padding: 2rem;
  box-shadow: var(--shadow-md);
  border: 1px solid var(--border-color);
}

.reports-tabs {
  display: flex;
  gap: 2rem;
  margin-bottom: 2rem;
  border-bottom: 1px solid var(--border-color);
}

.tab-button {
  padding: 0.75rem 0;
  font-weight: 500;
  color: var(--muted-color);
  border-bottom: 2px solid transparent;
  transition: all 0.2s ease;
}

.tab-button.active {
  color: var(--primary-color);
  border-bottom-color: var(--primary-color);
}

.reports-grid {
  display: grid;
  gap: 1rem;
}

.report-card {
  background: var(--background-color);
  border-radius: var(--border-radius-md);
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  transition: all 0.2s ease;
}

.report-card:hover {
  transform: translateX(4px);
  box-shadow: var(--shadow-sm);
}

.report-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.report-title {
  font-weight: 600;
  color: var(--text-color);
}

.report-date {
  font-size: 0.875rem;
  color: var(--muted-color);
}

.report-content {
  color: var(--muted-color);
  margin-bottom: 1rem;
}

.report-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.btn {
  font-family: 'Inter', sans-serif;
  padding: 0.75rem 1.5rem;
  border-radius: var(--border-radius-lg);
  font-weight: 600;
  transition: all 0.2s ease;
  display: inline-flex;
  align-items: center;
  gap: 0.5rem;
}

.btn-primary {
  background: var(--primary-color);
  color: var(--white);
  border: none;
}

.btn-secondary {
  background: var(--white);
  color: var(--text-color);
  border: 1px solid var(--border-color);
}

.btn:hover {
  transform: translateY(-2px);
  box-shadow: var(--shadow-md);
}

.dialog {
  border-radius: var(--border-radius-lg);
  overflow: hidden;
}

.dialog-header {
  padding: 1.5rem;
  background: var(--background-color);
  border-bottom: 1px solid var(--border-color);
}

.dialog-content {
  padding: 2rem;
}

.form-group {
  margin-bottom: 1.5rem;
}

.form-label {
  display: block;
  margin-bottom: 0.5rem;
  font-weight: 500;
  color: var(--text-color);
}

.form-input {
  width: 100%;
  padding: 0.75rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: var(--border-radius-md);
  font-size: 1rem;
  transition: all 0.2s ease;
}

.form-input:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

@media (max-width: 768px) {
  .dashboard-container {
    padding: 1rem;
  }
  
  .dashboard-title {
    font-size: 1.75rem;
  }
  
  .stats-grid {
    grid-template-columns: 1fr;
  }
  
  .search-form {
    flex-direction: column;
  }
  
  .view-toggle {
    flex-direction: column;
  }
  
  .view-button {
    width: 100%;
  }
  
  .reports-tabs {
    overflow-x: auto;
    padding-bottom: 0.5rem;
  }
  
  .report-footer {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
  }
}

/* 動畫效果 */
@keyframes slideIn {
  from {
    opacity: 0;
    transform: translateY(20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.fade-enter-active {
  animation: slideIn 0.3s ease-out;
}

/* 操作按鈕樣式 */
.action-btn {
  min-width: auto !important;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  height: 32px;
  padding: 0 8px;
}

/* 健康數據內容單元格樣式 */
.content-cell {
  max-width: 100%;
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 狀態標籤樣式 */
.status-chip {
  width: 100%;
  justify-content: center;
}

/* 操作單元格樣式 */
.actions-cell {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
}

/* 搜尋結果窗口 */
.result-window {
  width: 100%;
  overflow-x: auto;
  padding-top: 6px;
}

/* 報告表格樣式 */
.report-table {
  width: 100%;
}

/* 標籤與內容間距 */
.tab-content-divider {
  height: 6px;
  background-color: #f0f7ff;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
  margin-top: 0;
  position: relative;
  z-index: 1;
}

@media (max-width: 600px) {
  .action-btn {
    font-size: 0.8rem;
    padding: 0 8px !important;
    max-width: 100%;
}
  
  :deep(.v-data-table__wrapper) {
    overflow-x: auto;
  }
  
  :deep(.v-data-table-header th) {
    white-space: nowrap;
  }
}

:deep(.v-data-table) {
  background-color: transparent !important;
  border-radius: 16px !important;
  overflow: hidden;
}

:deep(.v-data-table-header th) {
  font-weight: 600 !important;
  color: #333 !important;
  background-color: #f9f9f9 !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  padding: 16px !important;
  border-bottom: 1px solid #eee !important;
}

:deep(.v-data-table-row) {
  transition: all 0.2s ease !important;
}

:deep(.v-data-table-row:hover) {
  background-color: #fafafa !important;
  transform: translateY(-1px);
}

:deep(.v-data-table td) {
  padding: 16px !important;
  color: #444 !important;
  font-weight: 400 !important;
}

:deep(.v-tabs) {
  border-bottom: 1px solid #eee;
}

:deep(.v-tabs .v-tab) {
  text-transform: none !important;
  letter-spacing: 0 !important;
  font-weight: 500 !important;
  color: #666 !important;
  min-height: 48px !important;
  padding: 0 24px !important;
}

:deep(.v-tab--selected) {
  color: #333 !important;
  font-weight: 600 !important;
}

:deep(.v-tabs-slider) {
  height: 2px !important;
}

.unauthorized-chip {
  background-color: #FFF8E1 !important;
  color: #F9A825 !important;
}

.authorized-chip {
  background-color: #F1F8E9 !important;
  color: #558B2F !important;
}

.auth-request-dialog {
  border-radius: 24px !important;
  overflow: hidden !important;
}

:deep(.v-dialog > .v-card) {
  border-radius: 24px !important;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1) !important;
}

:deep(.v-text-field) {
  border-radius: 12px !important;
}

:deep(.v-text-field .v-field) {
  border-radius: 12px !important;
  background-color: #fff !important;
  border: 1px solid #e5e7eb !important;
}

:deep(.v-text-field .v-field:focus-within) {
  border-color: #F8F441 !important;
  box-shadow: 0 2px 8px rgba(248, 244, 65, 0.2) !important;
}

@media (max-width: 600px) {
  .dashboard-container {
    padding: 16px !important;
  }

  .search-card, .result-card {
    border-radius: 16px !important;
  }

  .v-col {
    padding: 8px !important;
  }

  :deep(.v-data-table-header) {
    display: none !important;
  }

  :deep(.v-data-table-row) {
    display: flex !important;
    flex-direction: column !important;
    padding: 16px !important;
    border-bottom: 1px solid #eee !important;
  }

  :deep(.v-data-table td) {
    padding: 4px 0 !important;
    border: none !important;
  }

  :deep(.v-data-table td::before) {
    content: attr(data-label);
    font-weight: 600 !important;
    margin-right: 8px !important;
    color: #666 !important;
  }

  .modern-btn {
    width: 100% !important;
    margin-bottom: 8px !important;
  }
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease, transform 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
  transform: translateY(10px);
}

.v-card {
  transition: all 0.3s ease !important;
}

.v-card:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.1) !important;
}

.v-btn {
  transition: all 0.3s ease !important;
}

.v-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

.logout-btn {
  min-width: 110px !important;
  height: 48px !important;
  font-size: 1.15rem !important;
  padding: 0 20px !important;
}

.logout-text {
  font-size: 1.15rem;
  font-weight: 600;
  letter-spacing: 1px;
  display: inline-block;
}

.view-all-btn {
  align-self: center;
  margin-left: 12px;
  border-radius: 24px !important;
  letter-spacing: 1px;
}

/* 已授權報告卡片樣式 */
.authorized-reports-card {
  background: white;
  border-radius: 24px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  overflow: hidden;
}

.auth-reports-header {
  background: #F9F7F4;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

/* 返回按鈕樣式 */
.back-btn {
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
  border-radius: 12px !important;
  font-weight: 500 !important;
  letter-spacing: 0 !important;
}

.back-btn:hover {
  background-color: #f5f5f5 !important;
  transform: translateY(-1px);
}

/* 表格樣式 */
.authorized-reports-table {
  background: white !important;
}

:deep(.v-data-table) {
  border-radius: 24px !important;
  overflow: hidden !important;
}

:deep(.v-data-table__wrapper) {
  border-radius: 24px !important;
}

:deep(.v-data-table-header th) {
  font-weight: 600 !important;
  color: #333 !important;
  background: white !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  padding: 16px !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05) !important;
}

:deep(.v-data-table-row) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.05) !important;
  transition: all 0.2s ease !important;
}

:deep(.v-data-table-row:nth-child(even)) {
  background-color: rgba(0, 0, 0, 0.01) !important;
}

:deep(.v-data-table-row:hover) {
  background-color: rgba(36, 180, 126, 0.05) !important;
}

:deep(.v-data-table td) {
  color: #666 !important;
  font-weight: 400 !important;
  padding: 12px 16px !important;
}

/* ID 欄位樣式 */
.id-cell {
  position: relative;
  cursor: pointer;
}

.id-tooltip {
  display: none;
  position: absolute;
  top: -30px;
  left: 50%;
  transform: translateX(-50%);
  background: #333;
  color: white;
  padding: 4px 8px;
  border-radius: 4px;
  font-size: 12px;
  white-space: nowrap;
  z-index: 1000;
}

.id-cell:hover .id-tooltip {
  display: block;
}

/* 日期欄位樣式 */
.date-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
}

/* 查看報告按鈕樣式 */
.view-report-btn {
  color: #666 !important;
  background: rgba(0, 0, 0, 0.03) !important;
  border-radius: 8px !important;
  transition: all 0.2s ease !important;
}

.view-report-btn:hover {
  background: #F8F441 !important;
  color: #333 !important;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(248, 244, 65, 0.2) !important;
}

/* 分頁元件樣式 */
:deep(.v-data-table-footer) {
  background: white !important;
  border-top: 1px solid rgba(0, 0, 0, 0.05) !important;
  padding: 12px 16px !important;
}

:deep(.v-data-table-footer select) {
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
  border-radius: 8px !important;
  padding: 4px 8px !important;
}

:deep(.v-data-table-footer .v-btn) {
  color: #666 !important;
  border: 1px solid rgba(0, 0, 0, 0.1) !important;
  border-radius: 8px !important;
}

:deep(.v-data-table-footer .v-btn:hover) {
  background: rgba(36, 180, 126, 0.1) !important;
  color: #24B47E !important;
}

/* RWD 適配 */
@media (max-width: 600px) {
  .authorized-reports-card {
    border-radius: 16px;
  }

  :deep(.v-data-table-row) {
    display: flex !important;
    flex-direction: column !important;
    padding: 16px !important;
    margin-bottom: 8px !important;
    background: white !important;
    border-radius: 12px !important;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05) !important;
  }

  :deep(.v-data-table td) {
    display: flex !important;
    justify-content: space-between !important;
    align-items: center !important;
    padding: 8px 0 !important;
    border: none !important;
  }

  :deep(.v-data-table td::before) {
    content: attr(data-label);
    font-weight: 500 !important;
    color: #333 !important;
  }

  .date-cell {
    justify-content: flex-end;
  }

  .view-report-btn {
    margin-left: auto;
  }
}

/* 待授權請求卡片樣式 */
.pending-requests-card, .history-requests-card {
  background: white;
  border-radius: 24px;
  border: 1px solid rgba(0, 0, 0, 0.05);
  overflow: hidden;
  height: 100%;
}

.pending-header {
  background: #FFF8E1;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.history-header {
  background: #ECEFF1;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.pending-requests-table, .history-requests-table {
  background: white !important;
}
</style>