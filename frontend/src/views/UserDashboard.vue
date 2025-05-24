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
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="11" md="10" lg="9">
        <v-card class="pa-6 mb-6" elevation="8">
          <v-row align="center" justify="space-between" class="mb-4">
            <v-col cols="8">
              <h2 class="mb-1">👤 使用者儀表板</h2>
              <div class="subtitle-1">歡迎，{{ currentUser }}</div>
            </v-col>
            <v-col cols="4" class="d-flex justify-end align-center">
              <v-btn color="primary" @click="handleLogout" elevation="2">登出</v-btn>
            </v-col>
          </v-row>
        </v-card>

        <!-- 健檢報告列表 -->
        <v-card class="pa-4 mb-6" elevation="4">
          <h3 class="mb-3">健康檢查報告</h3>
          <v-data-table
            :headers="[
              { text: '報告編號', value: 'id', width: 120 },
              { text: '內容', value: 'content', width: 300 },
              { text: '日期', value: 'date', width: 150 },
              { text: '操作', value: 'actions', sortable: false, align: 'center', width: 100 }
            ]"
            :items="healthData"
            :loading="loading"
            loading-text="資料載入中..."
            class="elevation-0"
            dense
            hide-default-footer
            :no-data-text="'暫無資料'"
          >
            <template v-slot:item.content="{ item }">
              <span class="content-preview">{{ item.content }}</span>
            </template>
            <template v-slot:item.date="{ item }">
              {{ formatDate(item.date) }}
            </template>
            <template v-slot:item.actions="{ item }">
              <v-tooltip bottom>
                <template v-slot:activator="{ on, attrs }">
                  <v-btn
                    small
                    color="info"
                    icon
                    @click="viewReportDetail(item)"
                    elevation="1"
                    v-bind="attrs"
                    v-on="on"
                  >
                    <v-icon>mdi-magnify</v-icon>
                  </v-btn>
                </template>
                <span>查看詳細數據</span>
              </v-tooltip>
            </template>
          </v-data-table>
        </v-card>

        <!-- 資料授權區塊 -->
        <v-card class="pa-4 mb-6" elevation="4">
          <h3 class="mb-3">資料授權管理</h3>
          
          <v-tabs v-model="authTab" color="primary" slider-color="primary">
            <v-tab value="requests" class="text-none">
              <v-icon start>mdi-clipboard-alert</v-icon>
              授權請求
              <v-badge
                v-if="accessRequests.length > 0"
                :content="accessRequests.length.toString()"
                color="error"
                offset-x="5"
                offset-y="-5"
              ></v-badge>
            </v-tab>
            <v-tab value="authorized" class="text-none">
              <v-icon start>mdi-clipboard-check</v-icon>
              已授權清單
              <v-badge
                v-if="authorizedTickets.length > 0"
                :content="authorizedTickets.length.toString()"
                color="success"
                offset-x="5"
                offset-y="-5"
              ></v-badge>
            </v-tab>
          </v-tabs>
          
          <v-divider class="mb-3"></v-divider>
          
          <v-window v-model="authTab">
            <!-- 授權請求分頁 -->
            <v-window-item value="requests">
              <v-data-table
                :headers="[
                  { title: '報告編號', key: 'reportId', width: '120px' },
                  { title: '請求者', key: 'requesterName', width: '120px' },
                  { title: '授權理由', key: 'reason', width: '200px' },
                  { title: '申請時間', key: 'requestTime', width: '150px' },
                  { title: '到期時間', key: 'expiry', width: '150px' },
                  { title: '操作', key: 'actions', width: '180px', sortable: false }
                ]"
                :items="accessRequests"
                :loading="loadingRequests"
                loading-text="載入中..."
                no-data-text="目前沒有授權請求"
                hide-default-footer
                class="elevation-0"
                :class="{'opacity-50': authProcessing}"
              >
                <template v-slot:item.requestTime="{ item }">
                  {{ formatTimestamp(item.requestTime) }}
                </template>
                
                <template v-slot:item.expiry="{ item }">
                  {{ item.expiry ? formatTimestamp(item.expiry) : '永久' }}
                </template>
                
                <template v-slot:item.reason="{ item }">
                  <div class="reason-cell">{{ item.reason }}</div>
                </template>
                
                <template v-slot:item.actions="{ item }">
                  <div class="d-flex gap-2">
                    <v-btn
                      color="success"
                      size="small"
                      :loading="authProcessing"
                      :disabled="authProcessing"
                      @click="approveRequest(item.id)"
                      prepend-icon="mdi-check"
                    >
                      同意
                    </v-btn>
                    <v-btn
                      color="error"
                      size="small"
                      :loading="authProcessing"
                      :disabled="authProcessing"
                      @click="rejectRequest(item.id)"
                      prepend-icon="mdi-close"
                    >
                      拒絕
                    </v-btn>
                  </div>
                </template>
              </v-data-table>
              
              <div v-if="!loadingRequests && accessRequests.length === 0" class="text-center py-5">
                <v-icon size="64" color="grey-lighten-1">mdi-inbox-outline</v-icon>
                <div class="text-h6 mt-2 text-grey-darken-1">目前沒有待處理的授權請求</div>
                <div class="text-body-2 mt-1 text-grey">當保險業者請求訪問您的健康報告時，將顯示在這裡</div>
              </div>
            </v-window-item>
            
            <!-- 已授權清單分頁 -->
            <v-window-item value="authorized">
              <v-data-table
                :headers="[
                  { title: '報告編號', key: 'reportId', width: '120px' },
                  { title: '授權對象', key: 'targetName', width: '120px' },
                  { title: '授權時間', key: 'grantTime', width: '150px' },
                  { title: '到期時間', key: 'expiry', width: '150px' },
                  { title: '狀態', key: 'status', width: '100px' }
                ]"
                :items="authorizedTickets"
                :loading="loadingTickets"
                loading-text="載入中..."
                no-data-text="目前沒有已授權報告"
                hide-default-footer
                class="elevation-0"
              >
                <template v-slot:item.grantTime="{ item }">
                  {{ formatTimestamp(item.grantTime) }}
                </template>
                
                <template v-slot:item.expiry="{ item }">
                  <span v-if="item.expiry && item.expiry !== '0'">
                    {{ formatTimestamp(item.expiry) }}
                  </span>
                  <span v-else class="text-green">永久</span>
                </template>
                
                <template v-slot:item.status="{ item }">
                  <v-chip
                    :color="new Date().getTime() > item.expiry * 1000 && item.expiry !== 0 ? 'grey' : 'success'"
                    size="small"
                    variant="outlined"
                  >
                    {{ new Date().getTime() > item.expiry * 1000 && item.expiry !== 0 ? '已過期' : '有效' }}
                  </v-chip>
                </template>
              </v-data-table>
              
              <div v-if="!loadingTickets && authorizedTickets.length === 0" class="text-center py-5">
                <v-icon size="64" color="grey-lighten-1">mdi-shield-outline</v-icon>
                <div class="text-h6 mt-2 text-grey-darken-1">目前沒有已授權的健康報告</div>
                <div class="text-body-2 mt-1 text-grey">當您同意授權請求後，授權記錄將顯示在這裡</div>
              </div>
            </v-window-item>
          </v-window>
        </v-card>

        <!-- LLM 分析區塊 -->
        <v-card class="pa-4 mb-6" elevation="4">
          <h3 class="mb-3">AI 健檢摘要分析</h3>
          <v-btn
            color="deep-purple accent-4"
            :loading="llmLoading"
            :disabled="llmLoading || !healthData.length"
            @click="handleLLMSummary"
            elevation="2"
            class="mb-3"
          >
            <v-icon left>mdi-robot</v-icon> 產生 AI 摘要
          </v-btn>
          <v-alert
            v-if="llmSummary"
            type="info"
            class="mt-3"
            border="left"
            colored-border
            elevation="1"
          >
            {{ llmSummary }}
          </v-alert>
        </v-card>
      </v-col>
    </v-row>

    <!-- 健康報告詳細資料對話框 -->
    <v-dialog v-model="detailDialog" max-width="900" scrollable>
      <v-card v-if="selectedReport" class="report-detail-card">
        <v-card-title class="headline primary--text">
          <v-icon large color="primary" class="mr-2">mdi-clipboard-pulse</v-icon>
          健康檢查報告詳情
        </v-card-title>
        
        <v-card-subtitle>
          報告編號：{{ selectedReport.id }} | 日期：{{ formatDate(selectedReport.date) }}
        </v-card-subtitle>
        
        <v-divider></v-divider>
        
        <v-card-text>
          <v-container>
            <v-row v-if="reportMetrics.length === 0">
              <v-col cols="12" class="text-center">
                <v-alert type="info" outlined>
                  此報告無法解析為視覺化指標，請查看原始數據
                </v-alert>
              </v-col>
            </v-row>
            
            <template v-else>
              <v-row>
                <v-col cols="12">
                  <h3 class="font-weight-bold mb-3">健康指標視覺化</h3>
                </v-col>
              </v-row>
              
              <!-- 視覺化圓圈指標 -->
              <v-row>
                <v-col
                  v-for="metric in reportMetrics.filter(m => !m.isText)"
                  :key="metric.key"
                  cols="12"
                  sm="6"
                  md="4"
                  class="text-center mb-4"
                >
                  <div class="metric-container">
                    <v-progress-circular
                      :rotate="-90"
                      :size="120"
                      :width="15"
                      :value="metric.percentage"
                      :color="metric.color"
                      :class="{'pulse-animation': isOutsideNormalRange(metric)}"
                    >
                      {{ metric.value }}
                    </v-progress-circular>
                    <div class="metric-details mt-2">
                      <div class="metric-name">{{ metric.name }}</div>
                      <div class="metric-value">
                        {{ metric.value }} <span class="metric-unit">{{ metric.unit }}</span>
                      </div>
                      <v-chip
                        x-small
                        :color="metric.color"
                        text-color="white"
                        class="mt-1"
                      >
                        {{ metric.status }}
                      </v-chip>
                    </div>
                  </div>
                </v-col>
              </v-row>
              
              <!-- 文字指標 -->
              <v-row v-if="reportMetrics.some(m => m.isText)">
                <v-col cols="12">
                  <h3 class="font-weight-bold mt-3 mb-3">其他健康數據</h3>
                </v-col>
                <v-col
                  v-for="metric in reportMetrics.filter(m => m.isText)"
                  :key="metric.key"
                  cols="12"
                  sm="6"
                  class="mb-3"
                >
                  <v-card outlined class="text-metric-card pa-3">
                    <div class="text-metric-name">{{ metric.name }}</div>
                    <div class="text-metric-value">{{ metric.textValue }}</div>
                  </v-card>
                </v-col>
              </v-row>
            </template>
            
            <!-- 參考範圍說明 -->
            <v-row class="mt-3">
              <v-col cols="12">
                <v-alert
                  type="info"
                  text
                  dense
                  colored-border
                  border="left"
                >
                  <div class="text-caption">
                    <strong>指標說明：</strong>
                    <span class="mr-2"><v-icon x-small color="green">mdi-circle</v-icon> 正常</span>
                    <span class="mr-2"><v-icon x-small color="orange">mdi-circle</v-icon> 偏低</span>
                    <span><v-icon x-small color="red">mdi-circle</v-icon> 偏高</span>
                  </div>
                </v-alert>
              </v-col>
            </v-row>
            
            <!-- 原始數據 -->
            <v-row class="mt-3">
              <v-col cols="12">
                <v-expansion-panels flat>
                  <v-expansion-panel>
                    <v-expansion-panel-header class="pb-1">
                      <div class="text-subtitle-1 font-weight-medium">
                        <v-icon small class="mr-1">mdi-code-json</v-icon>
                        查看原始數據
                      </div>
                    </v-expansion-panel-header>
                    <v-expansion-panel-content>
                      <v-simple-table dense>
                        <thead>
                          <tr>
                            <th class="text-left">指標</th>
                            <th class="text-left">數值</th>
                          </tr>
                        </thead>
                        <tbody>
                          <tr v-for="(value, key) in selectedReport.rawData" :key="key">
                            <td>{{ key }}</td>
                            <td>{{ value }}</td>
                          </tr>
                        </tbody>
                      </v-simple-table>
                    </v-expansion-panel-content>
                  </v-expansion-panel>
                </v-expansion-panels>
              </v-col>
            </v-row>
          </v-container>
        </v-card-text>
        
        <v-divider></v-divider>
        
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn
            color="primary"
            text
            @click="detailDialog = false"
          >
            關閉
          </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<style scoped>
.fill-height {
  min-height: 100vh;
  background: #f5f6fa;
}
h2 {
  font-size: 28px;
  font-weight: bold;
}
h3 {
  font-size: 20px;
  margin-top: 0;
  font-weight: 500;
}
.subtitle-1 {
  color: #666;
}

/* 報告列表樣式 */
.content-preview {
  display: inline-block;
  max-width: 300px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* 視覺化指標樣式 */
.report-detail-card {
  max-height: 90vh;
}
.metric-container {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 16px;
  border-radius: 12px;
  background-color: #f9f9f9;
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
  height: 100%;
  transition: all 0.3s ease;
}
.metric-container:hover {
  transform: translateY(-5px);
  box-shadow: 0 5px 15px rgba(0,0,0,0.1);
  background-color: #f0f9ff;
}
.metric-details {
  text-align: center;
  padding-top: 12px;
}
.metric-name {
  font-weight: bold;
  font-size: 1.1rem;
  color: #333;
}
.metric-value {
  font-size: 1.2rem;
  margin-top: 5px;
  color: #424242;
  font-weight: 500;
}
.metric-unit {
  font-size: 0.8rem;
  color: #666;
}
.text-metric-card {
  background-color: #f5f5f5;
  height: 100%;
  transition: all 0.2s;
}
.text-metric-card:hover {
  background-color: #e8f5e9;
}
.text-metric-name {
  font-weight: bold;
  margin-bottom: 5px;
  color: #424242;
}
.text-metric-value {
  font-size: 0.9rem;
  color: #616161;
}

/* 異常值動畫效果 */
@keyframes pulse {
  0% {
    box-shadow: 0 0 0 0 rgba(255, 82, 82, 0.4);
  }
  70% {
    box-shadow: 0 0 0 10px rgba(255, 82, 82, 0);
  }
  100% {
    box-shadow: 0 0 0 0 rgba(255, 82, 82, 0);
  }
}
.pulse-animation {
  animation: pulse 2s infinite;
}

/* 授權管理相關樣式 */
.reason-cell {
  max-width: 200px;
  white-space: normal;
  word-break: break-word;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.opacity-50 {
  opacity: 0.5;
  pointer-events: none;
}

.gap-2 {
  gap: 8px;
}

:deep(.v-data-table .v-table__wrapper) {
  border-radius: 8px;
  overflow: hidden;
}

:deep(.v-data-table thead) {
  background-color: #f5f5f5;
}

:deep(.v-data-table tbody tr:hover) {
  background-color: rgba(0, 0, 0, 0.03);
}

.text-green {
  color: #2e7d32;
}
</style>