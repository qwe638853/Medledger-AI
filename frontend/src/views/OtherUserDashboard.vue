<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import healthCheckService from '../services/healthCheckService';

const authStore = useAuthStore();
const currentUser = computed(() => authStore.currentUser);
const token = ref(authStore.token);
const patientReports = ref([]);
const loading = ref(false);

// 分類報告
const authorizedReports = computed(() => patientReports.value.filter(report => report.is_authorized));
const unauthorizedReports = computed(() => patientReports.value.filter(report => !report.is_authorized));

// 所有已授權報告
const allAuthorizedReports = ref([]);
const authorizedReportsDialog = ref(false);
const loadingAuthorizedReports = ref(false);

// 視圖模式
const viewMode = ref('search'); // 'search' 或 'authorized'

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
  totalPatients: 0
});

// 獲取 Dashboard 數據
const fetchDashboardStats = async () => {
  try {
    dashboardStats.value = await healthCheckService.fetchDashboardStats();
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
    fetchDashboardStats(); // 更新儀表板數據
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
onMounted(() => {
  fetchAllAuthorizedReports();
  fetchDashboardStats();
});

const handleLogout = () => {
  authStore.logout();
};
</script>

<template>
  <div class="dashboard-bg">
    <v-container class="dashboard-container py-6 mx-auto" max-width="1280">
      <!-- 四張統計卡片置於同一行 -->
      <v-row class="mb-6" justify="center">
        <!-- 歡迎卡片 -->
        <v-col cols="12" sm="6" md="3" class="px-2">
          <v-card class="rounded-lg welcome-card" elevation="2" height="100">
            <v-card-text class="d-flex align-center pa-4">
              <div class="rounded-circle bg-blue-lighten-5 p-3 me-3">
                <v-icon size="32" color="blue-darken-2">mdi-shield-account</v-icon>
              </div>
              <div class="flex-grow-1">
                <div class="text-subtitle-1 text-grey-darken-1">歡迎，保險顧問</div>
                <div class="text-h6 font-weight-bold text-blue-darken-3">{{ currentUser }}</div>
              </div>
              <v-btn 
                color="blue-darken-2" 
                @click="handleLogout" 
                elevation="1" 
                icon="mdi-logout"
                size="small"
              ></v-btn>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 已授權報告數量卡片 -->
        <v-col cols="12" sm="6" md="3" class="px-2">
          <v-card class="stat-card" elevation="2" height="100" @click="switchToAuthorizedView">
            <v-card-text class="pa-4">
              <div class="d-flex align-center">
                <div class="rounded-circle bg-green-lighten-5 p-2 me-3">
                  <v-icon size="24" color="green-darken-2">mdi-file-document-check</v-icon>
                </div>
                <div>
                  <div class="text-overline text-green-darken-1">已授權報告</div>
                  <div class="text-h4 font-weight-bold text-green-darken-3">{{ dashboardStats.totalAuthorized }}</div>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 待處理請求卡片 -->
        <v-col cols="12" sm="6" md="3" class="px-2">
          <v-card class="stat-card" elevation="2" height="100">
            <v-card-text class="pa-4">
              <div class="d-flex align-center">
                <div class="rounded-circle bg-orange-lighten-5 p-2 me-3">
                  <v-icon size="24" color="orange-darken-2">mdi-clock-outline</v-icon>
                </div>
                <div>
                  <div class="text-overline text-orange-darken-1">待處理請求</div>
                  <div class="text-h4 font-weight-bold text-orange-darken-3">{{ dashboardStats.pendingRequests }}</div>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 已授權客戶卡片 -->
        <v-col cols="12" sm="6" md="3" class="px-2">
          <v-card class="stat-card" elevation="2" height="100">
            <v-card-text class="pa-4">
              <div class="d-flex align-center">
                <div class="rounded-circle bg-purple-lighten-5 p-2 me-3">
                  <v-icon size="24" color="purple-darken-2">mdi-account-group</v-icon>
                </div>
                <div>
                  <div class="text-overline text-purple-darken-1">授權病患數</div>
                  <div class="text-h4 font-weight-bold text-purple-darken-3">{{ dashboardStats.totalPatients }}</div>
                </div>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 搜尋區塊水平置中 -->
      <v-row class="mb-6" justify="center">
        <v-col cols="12">
          <v-card class="rounded-lg" elevation="2">
            <v-card-text class="pa-5">
              <v-row align="center">
                <v-col cols="12" md="4" class="pb-md-0">
                  <div class="d-flex align-center mb-md-0 mb-3">
                    <v-icon size="24" color="blue-darken-2" class="me-2">mdi-account-search</v-icon>
                    <span class="text-h6 font-weight-bold text-blue-darken-2">病患健康報告查詢</span>
                  </div>
                </v-col>
                <v-col cols="12" md="5" class="pb-md-0 pt-md-0">
                  <v-text-field
                    v-model="patientId"
                    label="請輸入病患身分證字號"
                    placeholder="例如：A123456789"
                    variant="outlined"
                    density="compact"
                    hide-details
                    prepend-inner-icon="mdi-account-card"
                    bg-color="white"
                  ></v-text-field>
            </v-col>
                <v-col cols="12" md="3" class="pt-md-0 d-flex gap-3">
                  <v-btn
                    color="blue-darken-2"
                    :loading="searchLoading"
                    :disabled="searchLoading"
                    @click="searchPatientReports"
                    prepend-icon="mdi-magnify"
                    elevation="1"
                    class="flex-grow-1"
                    :class="{'v-btn--active': viewMode === 'search'}"
                  >
                    搜尋
                  </v-btn>
                  <v-btn
                    color="green-darken-1"
                    :loading="loadingAuthorizedReports"
                    :disabled="loadingAuthorizedReports"
                    @click="switchToAuthorizedView"
                    prepend-icon="mdi-folder-account"
                    elevation="1"
                    class="flex-grow-1"
                    :class="{'v-btn--active': viewMode === 'authorized'}"
                  >
                    授權報告
                  </v-btn>
            </v-col>
          </v-row>
            </v-card-text>
        </v-card>
        </v-col>
      </v-row>

      <!-- 搜尋結果區塊 -->
      <v-row v-if="showSearchResults && viewMode === 'search'" justify="center">
        <v-col cols="12">
          <v-card class="rounded-lg" elevation="2">
            <v-card-title class="py-3 px-5 bg-blue-lighten-5 d-flex align-center">
              <v-icon size="24" color="blue-darken-3" class="me-2">mdi-clipboard-text</v-icon>
              <span class="text-h6 font-weight-bold text-blue-darken-3">病患「{{ patientId }}」的健康報告</span>
            </v-card-title>
            
            <!-- 標籤頁 -->
            <v-tabs
              v-model="tab"
              color="blue-darken-1"
              align-tabs="start"
              class="px-3 pt-2 bg-blue-lighten-5"
              slider-color="blue-darken-3"
              height="80"
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
          <v-card class="rounded-lg" elevation="2">
            <v-card-title class="py-3 px-5 bg-green-lighten-5 d-flex align-center">
              <v-icon size="24" color="green-darken-2" class="me-2">mdi-folder-account</v-icon>
              <span class="text-h6 font-weight-bold text-green-darken-3">所有已授權健康報告</span>
              <v-spacer></v-spacer>
          <v-btn
                color="blue-darken-2"
                variant="text"
                size="small"
                @click="switchToSearchView"
                prepend-icon="mdi-keyboard-return"
              >
                返回搜尋
          </v-btn>
            </v-card-title>
            
            <v-data-table
              :headers="[
                { title: '報告編號', key: 'id', align: 'start', width: '100px' },
                { title: '病患 ID', key: 'patient_id', align: 'start', width: '130px' },
                { title: '健康數據', key: 'content', align: 'start' },
                { title: '報告日期', key: 'date', align: 'center', width: '120px' },
                { title: '授權到期', key: 'expiry', align: 'center', width: '120px' }
              ]"
              :items="allAuthorizedReports"
              :loading="loadingAuthorizedReports"
              loading-text="正在載入已授權報告..."
              class="elevation-0"
              hover
              item-value="id"
              density="compact"
            >
              <template v-slot:item.patient_id="{ item }">
                <v-chip size="small" color="blue-lighten-4" class="font-weight-medium">
                  {{ item.patient_id }}
                </v-chip>
              </template>
              <template v-slot:item.content="{ item }">
                <div style="max-width: 400px; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">
                  {{ item.content }}
                </div>
              </template>
              <template v-slot:item.date="{ item }">
                {{ item.date || '-' }}
              </template>
              <template v-slot:item.expiry="{ item }">
                <v-chip 
                  size="small" 
                  :color="new Date(item.expiry) < new Date() ? 'error-lighten-4' : 'green-lighten-4'"
                  :text-color="new Date(item.expiry) < new Date() ? 'error-darken-2' : 'green-darken-2'"
                  variant="outlined"
                >
                  {{ item.expiry || '永久' }}
                </v-chip>
              </template>
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

      <!-- 授權請求對話框 -->
      <v-dialog v-model="authRequestDialog" max-width="600" persistent>
        <v-card class="auth-request-dialog">
          <v-card-title class="py-4 px-6 bg-deep-purple-lighten-5 rounded-t-lg">
            <v-icon size="28" color="deep-purple-darken-1" class="me-3">mdi-key-chain</v-icon>
            <span class="text-h6 font-weight-medium text-deep-purple-darken-1">請求健康報告授權</span>
          </v-card-title>
          
          <!-- 報告資訊區塊 -->
          <v-card-text class="pa-6 pt-5">
            <v-sheet
              v-if="selectedReport"
              class="info-card mb-5 pa-4 rounded-lg"
              color="blue-lighten-5"
              elevation="0"
            >
              <div class="d-flex align-center mb-3">
                <v-icon color="blue-darken-1" class="me-2">mdi-clipboard-text</v-icon>
                <div class="text-subtitle-1 font-weight-medium text-blue-darken-1">報告資訊</div>
              </div>
              <v-divider class="mb-3"></v-divider>
              <div class="d-flex flex-column gap-1">
                <div class="d-flex align-center">
                  <div class="text-body-2 font-weight-medium text-grey-darken-2 me-2 info-label">報告編號：</div>
                  <div class="text-body-1">{{ selectedReport.id }}</div>
                </div>
                <div class="d-flex align-center">
                  <div class="text-body-2 font-weight-medium text-grey-darken-2 me-2 info-label">病患 ID：</div>
                  <div class="text-body-1">{{ patientId }}</div>
                </div>
                <div class="d-flex align-center">
                  <div class="text-body-2 font-weight-medium text-grey-darken-2 me-2 info-label">報告日期：</div>
                  <div class="text-body-1">{{ selectedReport.date || '未指定' }}</div>
                </div>
              </div>
            </v-sheet>

            <!-- 授權理由輸入 -->
            <v-textarea
              v-model="authReason"
              label="授權理由"
              placeholder="請詳細說明請求授權此報告的原因..."
              variant="outlined"
              rows="3"
              auto-grow
              class="mb-5"
              hide-details="auto"
              counter="200"
              :rules="[v => !!v || '請輸入授權理由']"
              bg-color="white"
            ></v-textarea>

            <!-- 授權到期日選擇 -->
            <div class="mb-1 font-weight-medium text-grey-darken-2">授權到期日</div>
            <v-date-picker
              v-model="authExpiry"
              class="mb-3 elevation-1 rounded-lg mx-auto"
              :min="new Date().toISOString().substr(0, 10)"
              color="deep-purple"
              width="100%"
            ></v-date-picker>
            
            <!-- 錯誤提示 -->
            <v-alert
              v-if="!authReason || !authExpiry"
              type="warning"
              variant="tonal"
              border="start"
              class="mt-4 rounded-lg"
              density="compact"
            >
              請填寫完整授權資訊才能繼續
            </v-alert>
          </v-card-text>
          
          <v-divider></v-divider>
          
          <v-card-actions class="pa-5">
            <v-spacer></v-spacer>
            <v-btn
              color="grey"
              variant="flat"
              @click="authRequestDialog = false"
              class="me-3"
            >
              取消
            </v-btn>
            <v-btn
              color="deep-purple"
              @click="sendAuthRequest"
              :loading="requestLoading"
              :disabled="requestLoading || !authReason || !authExpiry"
              prepend-icon="mdi-send"
              elevation="1"
            >
              送出授權請求
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-dialog>

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
  background-color: #f5f8fb;
  min-height: 100vh;
  display: flex;
  align-items: flex-start;
  justify-content: center;
}

.dashboard-container {
  width: 100%;
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
}

:deep(.v-data-table-header th) {
  font-weight: bold !important;
  color: rgba(0, 0, 0, 0.7) !important;
  background-color: #f0f7ff !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
}

:deep(.v-data-table-row:hover) {
  background-color: #f0f7ff !important;
}

:deep(.v-data-table__wrapper) {
  table-layout: fixed !important;
}

:deep(.v-data-table table) {
  table-layout: fixed;
  width: 100%;
}

:deep(.v-data-table td),
:deep(.v-data-table th) {
  padding: 0 8px !important;
  box-sizing: border-box;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

:deep(.v-tabs .v-tab) {
  font-size: 0.9rem;
  min-height: 42px;
  padding-bottom: 10px;
  position: relative;
  z-index: 5;
}

:deep(.v-tab--selected) {
  font-weight: bold !important;
}

:deep(.v-tabs .v-slide-group__content) {
  z-index: 5;
}

:deep(.v-tabs .v-tabs-slider-wrapper) {
  z-index: 4;
  height: 3px !important;
  margin-bottom: 5px;
}

.v-card {
  transition: transform 0.2s, box-shadow 0.2s;
  border-radius: 8px !important;
}

.v-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1) !important;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(0, 0, 0, 0.1) !important;
  cursor: pointer;
}

.welcome-card {
  background: linear-gradient(135deg, #e3f2fd 0%, #bbdefb 100%);
}

.w-100 {
  width: 100%;
}

.h-100 {
  height: 100%;
}

.rounded-circle {
  border-radius: 50% !important;
}

.p-2 {
  padding: 8px;
}

.p-3 {
  padding: 12px;
}

.gap-3 {
  gap: 12px;
}

.v-btn--active {
  font-weight: bold;
  transform: scale(1.05);
}

.text-truncate {
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 授權請求對話框樣式 */
.auth-request-dialog {
  border-radius: 12px;
  overflow: hidden;
}

.info-card {
  border: 1px solid rgba(25, 118, 210, 0.1);
}

.info-label {
  min-width: 80px;
}

:deep(.v-date-picker) {
  border-radius: 12px;
}
</style>