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
  <v-container class="dashboard-container fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="11" md="10" lg="9">
        <!-- 頂部用戶資訊卡片 -->
        <v-card class="user-card mb-8" elevation="2">
          <v-row align="center" justify="space-between">
            <v-col cols="12" sm="8">
              <div class="d-flex align-center">
                <div class="rounded-circle bg-yellow-lighten-4 p-3 me-3">
                  <v-icon size="32" color="grey-darken-3">mdi-account-outline</v-icon>
                </div>
                <div>
                  <h2 class="user-title mb-1">使用者儀表板</h2>
                  <div class="user-subtitle">{{ currentUser }}</div>
                </div>
              </div>
            </v-col>
            <v-col cols="12" sm="4" class="d-flex justify-end">
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
            </v-col>
          </v-row>
        </v-card>

        <!-- 健檢報告列表 -->
        <v-card class="report-card mb-8" elevation="0">
          <div class="d-flex align-center mb-6">
            <v-icon size="24" class="section-icon mr-3">mdi-file-document-outline</v-icon>
            <h3 class="section-title">健康檢查報告</h3>
          </div>
          
          <v-data-table
            :headers="[
              { title: '報告編號', key: 'id', width: '120px' },
              { title: '內容摘要', key: 'content', width: '300px' },
              { title: '檢查日期', key: 'date', width: '150px' },
              { title: '', key: 'actions', sortable: false, align: 'end', width: '80px' }
            ]"
            :items="healthData"
            :loading="loading"
            loading-text="資料載入中..."
            class="report-table"
            :custom-class="{ 'report-row': true }"
            hide-default-footer
          >
            <template v-slot:item.content="{ item }">
              <span class="content-preview">{{ item.content }}</span>
            </template>
            
            <template v-slot:item.date="{ item }">
              <span class="date-text">{{ formatDate(item.date) }}</span>
            </template>
            
            <template v-slot:item.actions="{ item }">
                  <v-btn
                    icon
                variant="text"
                class="view-btn"
                    @click="viewReportDetail(item)"
                  >
                <v-icon>mdi-eye-outline</v-icon>
                  </v-btn>
                </template>

            <template v-slot:no-data>
              <div class="empty-state">
                <v-icon size="40" class="empty-icon">mdi-file-outline</v-icon>
                <div class="empty-text">暫無健康檢查報告</div>
              </div>
            </template>
          </v-data-table>
        </v-card>

        <!-- 資料授權管理 -->
        <v-card class="auth-card" elevation="0">
          <div class="d-flex align-center mb-6">
            <v-icon size="24" class="section-icon mr-3">mdi-shield-outline</v-icon>
            <h3 class="section-title">資料授權管理</h3>
          </div>
          
          <v-tabs
            v-model="authTab"
            class="auth-tabs mb-6"
          >
            <v-tab value="requests" class="auth-tab">
              <v-icon start class="mr-2">mdi-clipboard-text-outline</v-icon>
              授權請求
              <v-badge
                v-if="accessRequests.length"
                :content="accessRequests.length.toString()"
                color="warning"
                offset-x="5"
                offset-y="-5"
              ></v-badge>
            </v-tab>
            <v-tab value="authorized" class="auth-tab">
              <v-icon start class="mr-2">mdi-check-circle-outline</v-icon>
              已授權清單
            </v-tab>
          </v-tabs>
          
          <v-window v-model="authTab">
            <v-window-item value="requests">
              <!-- 授權請求表格 -->
              <v-data-table
                :headers="[
                  { title: '報告編號', key: 'reportId' },
                  { title: '請求者', key: 'requesterName' },
                  { title: '授權理由', key: 'reason' },
                  { title: '申請時間', key: 'requestTime' },
                  { title: '操作', key: 'actions', align: 'end' }
                ]"
                :items="accessRequests"
                :loading="loadingRequests"
                class="auth-table"
                hide-default-footer
              >
                <!-- 表格內容模板 -->
                <template v-slot:item.actions="{ item }">
                  <div class="d-flex gap-2 justify-end">
                    <v-btn
                      class="approve-btn"
                      size="small"
                      :loading="authProcessing"
                      @click="approveRequest(item.id)"
                    >
                      <v-icon start size="18">mdi-check</v-icon>
                      同意
                    </v-btn>
                    <v-btn
                      class="reject-btn"
                      size="small"
                      :loading="authProcessing"
                      @click="rejectRequest(item.id)"
                    >
                      <v-icon start size="18">mdi-close</v-icon>
                      拒絕
                    </v-btn>
                  </div>
                </template>

                <template v-slot:no-data>
                  <div class="empty-state">
                    <v-icon size="40" class="empty-icon">mdi-inbox-outline</v-icon>
                    <div class="empty-title">暫無待處理的授權請求</div>
                    <div class="empty-subtitle">當有新的授權請求時會顯示在這裡</div>
              </div>
                </template>
              </v-data-table>
            </v-window-item>
          </v-window>
        </v-card>
      </v-col>
    </v-row>
  </v-container>
</template>

<style scoped>
/* 全局樣式 */
.dashboard-container {
  background-color: #F9F7F4;
  min-height: 100vh;
  padding: 2rem;
}

/* 卡片基礎樣式 */
:deep(.v-card) {
  border-radius: 28px !important;
  background: white !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03) !important;
  padding: 2rem !important;
}

/* 用戶資訊區域 */
.user-title {
  font-size: 1.75rem;
  font-weight: 900;
  color: #111827;
  letter-spacing: -0.5px;
}

.user-subtitle {
  font-size: 1rem;
  color: #888;
  font-weight: 400;
}

.user-icon {
  /* 移除原本顏色設定，統一用圓形背景 */
  color: unset;
}

.logout-btn {
  min-width: 110px !important;
  height: 48px !important;
  font-size: 1.15rem !important;
  padding: 0 20px !important;
  background-color: #333 !important;
  color: #fff !important;
  border-radius: 24px !important;
  font-weight: 600 !important;
  letter-spacing: 1px;
}

.logout-text {
  font-size: 1.15rem;
  font-weight: 600;
  letter-spacing: 1px;
  display: inline-block;
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

/* 區塊標題樣式 */
.section-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
  letter-spacing: -0.5px;
}

.section-icon {
  color: #111827;
}

/* 表格樣式 */
:deep(.v-data-table) {
  background: transparent !important;
  border-radius: 16px !important;
}

:deep(.v-data-table-header) {
  background: transparent !important;
}

:deep(.v-data-table-header th) {
  font-size: 0.875rem !important;
  color: #888 !important;
  font-weight: 500 !important;
  text-transform: none !important;
  letter-spacing: 0 !important;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05) !important;
}

:deep(.v-data-table-row) {
  transition: background-color 0.2s ease;
}

:deep(.v-data-table-row:nth-child(odd)) {
  background-color: rgba(0, 0, 0, 0.01);
}

:deep(.v-data-table-row:hover) {
  background-color: rgba(0, 0, 0, 0.02) !important;
}

/* 操作按鈕 */
.approve-btn {
  background-color: #43AA8B !important;
  color: white !important;
  border-radius: 12px !important;
}

.reject-btn {
  background-color: transparent !important;
  color: #666 !important;
  border: 1px solid #ddd !important;
  border-radius: 12px !important;
}

/* Tab 樣式 */
:deep(.v-tab) {
  text-transform: none !important;
  letter-spacing: 0 !important;
  font-weight: 500 !important;
  color: #666 !important;
  min-width: 120px !important;
}

:deep(.v-tab--selected) {
  color: #111827 !important;
  font-weight: 600 !important;
}

:deep(.v-tab--selected .v-tab__slider) {
  background-color: #F8F441 !important;
  height: 3px !important;
}

/* 空狀態樣式 */
.empty-state {
  text-align: center;
  padding: 3rem 0;
}

.empty-icon {
  color: #888;
  margin-bottom: 1rem;
}

.empty-title {
  font-size: 1.1rem;
  color: #666;
  font-weight: 500;
  margin-bottom: 0.5rem;
}

.empty-subtitle {
  font-size: 0.875rem;
  color: #888;
}

/* RWD 適配 */
@media (max-width: 960px) {
  .dashboard-container {
    padding: 1rem;
}
  
  :deep(.v-card) {
    padding: 1.5rem !important;
}
  
  .user-title {
    font-size: 1.5rem;
  }
  
  .section-title {
    font-size: 1.25rem;
  }
}

@media (max-width: 600px) {
  .user-card .v-row {
    flex-direction: column;
}

  .user-card .v-col {
    padding: 0.5rem 0;
}

  .logout-btn {
    width: 100%;
}

  :deep(.v-data-table) {
    font-size: 0.875rem;
  }
}
</style>