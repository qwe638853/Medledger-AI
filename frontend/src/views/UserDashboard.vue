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
        date: report.createdAt || report.timestamp || report.created_at || report.date || new Date().toISOString(),
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
    authorizedTickets.value = await healthCheckService.fetchGrantedTickets();
    console.log('載入已授權票據完成:', authorizedTickets.value);
    
    // 檢查獲取的數據是否完整
    if (authorizedTickets.value.length > 0) {
      authorizedTickets.value.forEach((ticket, index) => {
        if (!ticket.reportId || !ticket.targetId || !ticket.grantTime) {
          console.warn(`授權票據 #${index} 資料不完整:`, ticket);
        }
      });
    }
  } catch (error) {
    console.error('載入已授權票據失敗:', error);
    notifyError(`無法載入已授權票據：${error.message || '未知錯誤'}`);
    authorizedTickets.value = []; // 確保失敗時清空列表
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

// 格式化時間戳為日期
const formatTimestamp = (timestamp) => {
  if (!timestamp) return '未設定';
  
  try {
    console.log('格式化時間戳:', timestamp, typeof timestamp);
    
    // 如果是字符串數字，轉為數字
    if (typeof timestamp === 'string') {
      timestamp = parseInt(timestamp, 10);
    }
    
    // 確保是以秒為單位的時間戳
    if (timestamp < 10000000000) {
      // 如果時間戳是以秒為單位
      return new Date(timestamp * 1000).toLocaleDateString('zh-TW', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    } else {
      // 如果時間戳是以毫秒為單位
      return new Date(timestamp).toLocaleDateString('zh-TW', {
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit'
      });
    }
  } catch (e) {
    console.error('格式化時間戳失敗:', e, timestamp);
    return timestamp.toString();
  }
};

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
  const report_id = item.id;
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

// 格式化日期顯示
function formatDate(dateString) {
  if (!dateString) return '未知日期';
  
  try {
    const date = new Date(dateString);
    return date.toLocaleDateString('zh-TW', {
      year: 'numeric',
      month: '2-digit',
      day: '2-digit',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch (e) {
    return dateString;
  }
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
</script>

<template>
  <div class="dashboard-bg">
    <v-container class="dashboard-container py-8 mx-auto" max-width="1800">
      <!-- 頂部統計卡片區 -->
      <v-row class="mb-8">
        <!-- 用戶資訊卡片 -->
        <v-col cols="12" md="6" lg="3">
          <v-card class="user-info-card" elevation="0">
            <v-card-text class="pa-6">
              <div class="d-flex align-center mb-2">
                <div class="rounded-circle bg-primary-lighten-5 p-3">
                  <v-icon size="32" color="primary">mdi-account</v-icon>
                </div>
                <v-btn
                  color="grey-darken-3"
                  @click="handleLogout"
                  elevation="0"
                  icon
                  size="small"
                  class="ml-auto"
                >
                  <v-icon>mdi-logout</v-icon>
                </v-btn>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-medium-emphasis">歡迎回來</div>
                <div class="text-h6 font-weight-bold">{{ currentUser }}</div>
              </div>
              <div class="mt-2 text-caption text-medium-emphasis">
                上次登入：{{ new Date().toLocaleDateString('zh-TW') }}
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 健康報告統計卡片 -->
        <v-col cols="12" md="6" lg="3">
          <v-card class="stat-card" elevation="0">
            <v-card-text class="pa-6">
              <div class="d-flex align-center mb-2">
                <div class="rounded-circle bg-success-lighten-5 p-3">
                  <v-icon size="32" color="success">mdi-file-document</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-medium-emphasis">健康報告總數</div>
                <div class="text-h4 font-weight-bold">{{ healthData.length }}</div>
              </div>
              <div class="stat-progress mt-4">
                <v-progress-linear
                  model-value="70"
                  color="success"
                  height="4"
                  rounded
                ></v-progress-linear>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 待處理授權請求卡片 -->
        <v-col cols="12" md="6" lg="3">
          <v-card class="stat-card" elevation="0">
            <v-card-text class="pa-6">
              <div class="d-flex align-center mb-2">
                <div class="rounded-circle bg-warning-lighten-5 p-3">
                  <v-icon size="32" color="warning">mdi-clock-outline</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-medium-emphasis">待處理授權</div>
                <div class="text-h4 font-weight-bold">{{ accessRequests.length }}</div>
              </div>
              <div class="stat-progress mt-4">
                <v-progress-linear
                  :model-value="(accessRequests.length / 10) * 100"
                  color="warning"
                  height="4"
                  rounded
                ></v-progress-linear>
              </div>
            </v-card-text>
          </v-card>
        </v-col>

        <!-- 已授權報告卡片 -->
        <v-col cols="12" md="6" lg="3">
          <v-card class="stat-card" elevation="0">
            <v-card-text class="pa-6">
              <div class="d-flex align-center mb-2">
                <div class="rounded-circle bg-info-lighten-5 p-3">
                  <v-icon size="32" color="info">mdi-shield-check</v-icon>
                </div>
              </div>
              <div class="mt-4">
                <div class="text-subtitle-1 text-medium-emphasis">已授權報告</div>
                <div class="text-h4 font-weight-bold">{{ authorizedTickets.length }}</div>
              </div>
              <div class="stat-progress mt-4">
                <v-progress-linear
                  :model-value="(authorizedTickets.length / healthData.length) * 100"
                  color="info"
                  height="4"
                  rounded
                ></v-progress-linear>
              </div>
            </v-card-text>
          </v-card>
        </v-col>
      </v-row>

      <!-- 主要內容區 -->
      <v-row>
        <!-- 健康報告列表 -->
        <v-col cols="12" lg="8">
          <v-card class="report-card" elevation="0">
            <v-card-title class="d-flex align-center px-6 py-4 bg-surface">
              <v-icon size="24" color="primary" class="mr-3">mdi-file-document-multiple</v-icon>
              <span class="text-h6 font-weight-bold">健康檢查報告</span>
              <v-spacer></v-spacer>
              <v-btn
                color="primary"
                variant="tonal"
                size="small"
                prepend-icon="mdi-filter-variant"
                class="mr-2"
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
            </v-card-title>

            <v-divider></v-divider>

            <v-card-text class="pa-6">
              <v-data-table
                :headers="[
                  { title: '報告編號', key: 'id', align: 'start' },
                  { title: '檢查項目', key: 'content', align: 'start' },
                  { title: '檢查日期', key: 'date', align: 'center' },
                  { title: '狀態', key: 'status', align: 'center' },
                  { title: '操作', key: 'actions', align: 'end' }
                ]"
                :items="healthData"
                :loading="loading"
                hover
                class="report-table"
              >
                <!-- 報告編號 -->
                <template v-slot:item.id="{ item }">
                  <div class="d-flex align-center">
                    <v-icon size="18" color="primary" class="mr-2">mdi-file-document</v-icon>
                    <span class="font-weight-medium">{{ item.id }}</span>
                  </div>
                </template>

                <!-- 檢查項目 -->
                <template v-slot:item.content="{ item }">
                  <div class="content-cell">
                    <span class="text-truncate">{{ item.content }}</span>
                  </div>
                </template>

                <!-- 檢查日期 -->
                <template v-slot:item.date="{ item }">
                  <div class="date-cell">
                    <v-icon size="16" color="grey" class="mr-1">mdi-calendar</v-icon>
                    {{ formatDate(item.date) }}
                  </div>
                </template>

                <!-- 狀態 -->
                <template v-slot:item.status="{ item }">
                  <v-chip
                    size="small"
                    :color="item.status === 'normal' ? 'success' : 'warning'"
                    variant="tonal"
                    class="font-weight-medium"
                  >
                    {{ item.status === 'normal' ? '正常' : '需追蹤' }}
                  </v-chip>
                </template>

                <!-- 操作按鈕 -->
                <template v-slot:item.actions="{ item }">
                  <div class="d-flex gap-2 justify-end">
                    <v-btn
                      icon="mdi-eye"
                      variant="text"
                      size="small"
                      color="primary"
                      @click="viewReportDetail(item)"
                    ></v-btn>
                    <v-btn
                      icon="mdi-share-variant"
                      variant="text"
                      size="small"
                      color="grey"
                    ></v-btn>
                  </div>
                </template>

                <!-- 無資料顯示 -->
                <template v-slot:no-data>
                  <div class="d-flex flex-column align-center py-8">
                    <v-icon size="64" color="grey-lighten-1" class="mb-4">
                      mdi-file-document-outline
                    </v-icon>
                    <div class="text-h6 font-weight-medium text-grey-darken-1">
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
        </v-col>

        <!-- 授權管理區域 -->
        <v-col cols="12" lg="4">
          <v-card class="auth-card" elevation="0">
            <v-card-title class="d-flex align-center px-6 py-4 bg-surface">
              <v-icon size="24" color="primary" class="mr-3">mdi-shield-account</v-icon>
              <span class="text-h6 font-weight-bold">授權管理</span>
            </v-card-title>

            <v-divider></v-divider>

            <v-card-text class="pa-0">
              <v-tabs
                v-model="authTab"
                color="primary"
                align-tabs="center"
                class="auth-tabs"
              >
                <v-tab value="requests" class="text-body-2 font-weight-medium">
                  <v-icon start size="18">mdi-clock-outline</v-icon>
                  待處理請求
                  <v-badge
                    v-if="accessRequests.length"
                    :content="accessRequests.length.toString()"
                    color="error"
                    offset-x="3"
                    offset-y="-3"
                  ></v-badge>
                </v-tab>
                <v-tab value="authorized" class="text-body-2 font-weight-medium">
                  <v-icon start size="18">mdi-shield-check</v-icon>
                  已授權清單
                </v-tab>
              </v-tabs>

              <v-window v-model="authTab" class="auth-window">
                <v-window-item value="requests">
                  <div class="pa-4">
                    <v-list class="request-list pa-0" v-if="accessRequests.length">
                      <v-list-item
                        v-for="request in accessRequests"
                        :key="request.id"
                        class="request-item mb-3"
                        rounded="lg"
                      >
                        <div class="d-flex flex-column">
                          <div class="d-flex align-center mb-2">
                            <v-avatar size="32" color="primary" class="mr-3">
                              <v-icon color="white" size="16">mdi-account</v-icon>
                            </v-avatar>
                            <div>
                              <div class="text-subtitle-2 font-weight-medium">
                                {{ request.requesterName }}
                              </div>
                              <div class="text-caption text-grey">
                                請求報告 #{{ request.reportId }}
                              </div>
                            </div>
                          </div>
                          <div class="text-caption text-grey-darken-1 mb-3">
                            {{ request.reason }}
                          </div>
                          <div class="d-flex align-center justify-space-between">
                            <div class="text-caption text-grey">
                              <v-icon size="14" class="mr-1">mdi-clock-outline</v-icon>
                              {{ formatTimestamp(request.requestTime) }}
                            </div>
                            <div class="d-flex gap-2">
                              <v-btn
                                color="success"
                                size="small"
                                variant="tonal"
                                :loading="authProcessing"
                                @click="approveRequest(request.id)"
                              >
                                同意
                              </v-btn>
                              <v-btn
                                color="error"
                                size="small"
                                variant="tonal"
                                :loading="authProcessing"
                                @click="rejectRequest(request.id)"
                              >
                                拒絕
                              </v-btn>
                            </div>
                          </div>
                        </div>
                      </v-list-item>
                    </v-list>
                    <div v-else class="d-flex flex-column align-center py-8">
                      <v-icon size="64" color="grey-lighten-1" class="mb-4">
                        mdi-tray-full
                      </v-icon>
                      <div class="text-h6 font-weight-medium text-grey-darken-1">
                        無待處理請求
                      </div>
                      <div class="text-body-2 text-grey mt-2">
                        目前沒有需要處理的授權請求
                      </div>
                    </div>
                  </div>
                </v-window-item>

                <v-window-item value="authorized">
                  <div class="pa-4">
                    <v-list class="auth-list pa-0" v-if="authorizedTickets.length">
                      <v-list-item
                        v-for="ticket in authorizedTickets"
                        :key="ticket.id"
                        class="auth-item mb-3"
                        rounded="lg"
                      >
                        <div class="d-flex flex-column">
                          <div class="d-flex align-center mb-2">
                            <v-avatar size="32" color="success" class="mr-3">
                              <v-icon color="white" size="16">mdi-shield-check</v-icon>
                            </v-avatar>
                            <div>
                              <div class="text-subtitle-2 font-weight-medium">
                                報告 #{{ ticket.reportId }}
                              </div>
                              <div class="text-caption text-grey">
                                授權給：{{ ticket.targetId }}
                              </div>
                            </div>
                          </div>
                          <div class="d-flex align-center justify-space-between">
                            <div class="text-caption text-grey">
                              <v-icon size="14" class="mr-1">mdi-calendar</v-icon>
                              {{ formatTimestamp(ticket.grantTime) }}
                            </div>
                            <v-chip
                              size="x-small"
                              color="success"
                              variant="tonal"
                              class="font-weight-medium"
                            >
                              已授權
                            </v-chip>
                          </div>
                        </div>
                      </v-list-item>
                    </v-list>
                    <div v-else class="d-flex flex-column align-center py-8">
                      <v-icon size="64" color="grey-lighten-1" class="mb-4">
                        mdi-shield-outline
                      </v-icon>
                      <div class="text-h6 font-weight-medium text-grey-darken-1">
                        無已授權報告
                      </div>
                      <div class="text-body-2 text-grey mt-2">
                        您目前沒有授權給他人的報告
                      </div>
                    </div>
                  </div>
                </v-window-item>
              </v-window>
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
  padding: 2rem;
  min-height: 100vh;
  max-width: 1800px;
  margin: 0 auto;
}

/* 卡片基礎樣式 */
:deep(.v-card) {
  border-radius: 24px !important;
  background: white !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
  transition: all 0.3s ease !important;
  overflow: hidden !important;
}

:deep(.v-card:hover) {
  transform: translateY(-2px);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.08) !important;
}

/* 統計卡片樣式 */
.stat-card {
  height: 100%;
}

.stat-progress {
  opacity: 0.7;
}

/* 圓形圖標背景 */
.rounded-circle {
  width: 52px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  transition: all 0.3s ease;
}

.rounded-circle:hover {
  transform: scale(1.05);
}

/* 表格樣式 */
.report-table {
  border-radius: 20px !important;
  overflow: hidden !important;
}

:deep(.v-data-table) {
  border-radius: 20px !important;
  overflow: hidden !important;
}

:deep(.v-data-table-header th) {
  font-size: 0.875rem !important;
  color: #64748B !important;
  font-weight: 600 !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  padding: 12px 16px !important;
}

:deep(.v-data-table-row) {
  border-bottom: 1px solid rgba(0, 0, 0, 0.05) !important;
  transition: background-color 0.2s ease;
}

:deep(.v-data-table-row:hover) {
  background-color: #F8FAFC !important;
}

:deep(.v-data-table td) {
  padding: 12px 16px !important;
  color: #334155 !important;
  font-size: 0.875rem !important;
}

/* 內容單元格樣式 */
.content-cell {
  max-width: 300px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 日期單元格樣式 */
.date-cell {
  display: flex;
  align-items: center;
  justify-content: center;
  color: #64748B;
}

/* 授權管理樣式 */
.auth-tabs {
  background-color: #F8FAFC;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.auth-window {
  height: calc(100vh - 400px);
  overflow-y: auto;
}

.request-list, .auth-list {
  gap: 12px;
}

.request-item, .auth-item {
  background-color: #F8FAFC !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
  transition: all 0.2s ease !important;
}

.request-item:hover, .auth-item:hover {
  transform: translateX(4px);
  background-color: white !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05) !important;
}

/* 按鈕樣式 */
:deep(.v-btn) {
  text-transform: none !important;
  letter-spacing: 0 !important;
  font-weight: 600 !important;
  transition: all 0.2s ease !important;
}

:deep(.v-btn:hover) {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1) !important;
}

/* 空狀態樣式 */
.empty-state {
  padding: 48px 24px;
  text-align: center;
  background-color: #F8FAFC;
  border-radius: 12px;
}

/* RWD 適配 */
@media (max-width: 960px) {
  .dashboard-container {
    padding: 1rem;
  }

  .auth-window {
    height: 500px;
  }
}

@media (max-width: 600px) {
  :deep(.v-card) {
    border-radius: 12px !important;
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

  :deep(.v-data-table td::before) {
    content: attr(data-label);
    font-weight: 500 !important;
    color: #64748B !important;
  }

  .content-cell {
    max-width: 100%;
  }

  .date-cell {
    justify-content: flex-end;
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
</style>