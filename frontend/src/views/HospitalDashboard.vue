<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores/auth';
import { healthCheckService, setupAuthInterceptor } from '../services';

const authStore = useAuthStore();
const currentUser = ref(authStore.currentUser);
const token = ref(authStore.token);
const uploadedData = ref([]);
const loading = ref(false);
const files = ref([]);
const uploadMessage = ref('');
const isDragging = ref(false);
const uploadProgress = ref(0);
const isUploading = ref(false);
const snackbar = ref({
  show: false,
  message: '',
  color: 'success'
});

// 新增病人身分證欄位
const patientId = ref('');
const patientIdRules = [
  v => !!v || '請輸入病人身分證字號',
  v => /^[A-Z][12]\d{8}$/.test(v) || '身分證字號格式不正確'
];

// 新增預覽對話框
const previewDialog = ref(false);
const parsedData = ref(null);
const parseError = ref(null);

onMounted(async () => {
  // 確保在組件掛載時重新設置 API 客戶端的身份驗證
  console.log('HospitalDashboard 掛載，當前用戶:', currentUser.value);
  console.log('令牌存在:', !!token.value);
  
  // 從 localStorage 中獲取最新令牌
  const storedToken = localStorage.getItem('token');
  if (storedToken) {
    console.log('發現存儲的令牌，長度:', storedToken.length);
    token.value = storedToken;
    
    // 重新設置 API 客戶端的身份驗證
    setupAuthInterceptor(storedToken);
  } else {
    console.warn('找不到存儲的令牌，可能需要重新登入');
    showSnackbar('會話可能已過期，請重新登入', 'warning');
  }
  
  //await fetchUploadedData();
});

const fetchUploadedData = async () => {
  loading.value = true;
  try {
    const data = await healthCheckService.fetchUserHealthData(currentUser.value);
    uploadedData.value = data || [];
  } catch (error) {
    showSnackbar(`獲取已上傳數據失敗：${error.message}`, 'error');
    uploadedData.value = [];
  } finally {
    loading.value = false;
  }
};

const showSnackbar = (message, color = 'success') => {
  snackbar.value = {
    show: true,
    message,
    color
  };
  
  // 也發送事件以兼容現有代碼
  document.dispatchEvent(new CustomEvent('show-snackbar', {
    detail: { message, color }
  }));
};

const handleFileDrop = (e) => {
  e.preventDefault();
  isDragging.value = false;
  
  if (!e.dataTransfer.files.length) return;
  
  const droppedFiles = Array.from(e.dataTransfer.files);
  // 過濾出支持的文件類型（包括JSON）
  const supportedFiles = droppedFiles.filter(file => {
    const fileType = file.type;
    const fileName = file.name.toLowerCase();
    return fileType.includes('excel') || 
           fileType.includes('spreadsheet') || 
           fileType.includes('csv') || 
           fileType.includes('pdf') ||
           fileType.includes('json') ||
           fileName.endsWith('.json');
  });
  
  if (supportedFiles.length === 0) {
    showSnackbar('請上傳支持的檔案格式（Excel、CSV、PDF或JSON）', 'error');
    return;
  }
  
  files.value = supportedFiles;
};

const handleFileSelect = (e) => {
  const selectedFiles = Array.from(e.target.files || []);
  if (selectedFiles.length === 0) return;
  
  files.value = selectedFiles;
};

const clearFiles = () => {
  files.value = [];
  parseError.value = null;
};

// 讀取 JSON 文件內容
const readJSONFile = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    
    reader.onload = (event) => {
      try {
        const jsonData = JSON.parse(event.target.result);
        resolve(jsonData);
      } catch (error) {
        reject(new Error(`無法解析 JSON 文件: ${error.message}`));
      }
    };
    
    reader.onerror = () => {
      reject(new Error('讀取文件時發生錯誤'));
    };
    
    reader.readAsText(file);
  });
};

// 新增預覽功能
const previewFiles = async () => {
  if (!files.value.length) {
    showSnackbar('請選擇檔案', 'error');
    return;
  }
  
  if (!patientId.value || !patientIdRules[1](patientId.value)) {
    showSnackbar('請輸入有效的病人身分證字號', 'error');
    return;
  }
  
  parseError.value = null;
  
  try {
    // 解析文件內容
    parsedData.value = await Promise.all(
      files.value.map(async (file) => {
        try {
          const fileExt = file.name.split('.').pop().toLowerCase();
          let parsed;
          
          // 針對不同文件格式進行處理
          if (fileExt === 'json' || file.type.includes('json')) {
            // 實際讀取 JSON 文件內容
            parsed = await readJSONFile(file);
            console.log('已解析 JSON 檔案:', file.name, parsed);
          } else {
            // 其他格式暫時仍使用模擬數據
            parsed = await parseHealthReportFile(file);
          }
          
          return {
            fileName: file.name,
            fileSize: formatFileSize(file.size),
            fileType: fileExt,
            patientId: patientId.value,
            reportData: parsed,
            isJson: fileExt === 'json' || file.type.includes('json')
          };
        } catch (err) {
          console.error(`解析文件 ${file.name} 失敗:`, err);
          throw err;
        }
      })
    );
    
    previewDialog.value = true;
  } catch (error) {
    parseError.value = error.message;
    showSnackbar(`解析檔案失敗：${error.message}`, 'error');
  }
};

// 解析健康檢查報告文件
const parseHealthReportFile = async (file) => {
  // 這裡使用已有的函數，但實際上我們先預覽
  return new Promise((resolve) => {
    setTimeout(() => {
      // 模擬解析結果
      resolve({
        'Glu-AC': Math.floor(Math.random() * 50 + 70) + ' mg/dL', // 血糖值
        'HbA1c': (Math.random() * 2 + 4).toFixed(1) + ' %',       // 糖化血色素
        'LDL-C': Math.floor(Math.random() * 70 + 80) + ' mg/dL',  // 低密度脂蛋白膽固醇
        'HDL-C': Math.floor(Math.random() * 20 + 40) + ' mg/dL',  // 高密度脂蛋白膽固醇
        'BP': `${Math.floor(Math.random() * 40 + 100)}/${Math.floor(Math.random() * 20 + 60)} mmHg` // 血壓
      });
    }, 500);
  });
};

// 修改上傳邏輯，確保JSON資料傳送至後端
const handleFileUpload = async () => {
  if (!parsedData.value || !patientId.value) {
    showSnackbar('請先預覽檔案並確認病人身分證字號', 'error');
    return;
  }

  isUploading.value = true;
  uploadProgress.value = 0;
  
  try {
    const updateProgress = (progress) => {
      uploadProgress.value = progress;
    };
    
    const uploadResults = [];
    
    // 上傳每個檔案
    for (let i = 0; i < files.value.length; i++) {
      const file = files.value[i];
      const fileData = parsedData.value[i];
      updateProgress(Math.floor((i / files.value.length) * 90));
      
      // 針對 JSON 文件特殊處理
      if (fileData.isJson) {
        try {
          console.log(`上傳 JSON 數據給病人 ${patientId.value}:`, fileData.reportData);
          
          // 直接使用解析好的 JSON 數據
          const result = await healthCheckService.uploadJsonHealthData(
            patientId.value,
            fileData.reportData,
            file.name,
            (progress) => {
              const baseProgress = Math.floor((i / files.value.length) * 90);
              const fileProgress = Math.floor((progress / 100) * (90 / files.value.length));
              updateProgress(baseProgress + fileProgress);
            }
          );
          
          uploadResults.push({
            ...result,
            fileName: file.name,
            fileType: 'json',
            uploadTime: new Date().toISOString()
          });
        } catch (error) {
          console.error(`上傳 JSON 數據失敗:`, error);
          uploadResults.push({
            success: false,
            error: true,
            fileName: file.name,
            message: `JSON 上傳失敗: ${error.message}`
          });
        }
      } else {
        // 其他類型文件使用原來的方法上傳
        try {
          const result = await healthCheckService.uploadHealthReport(
            file,
            patientId.value,
            (progress) => {
              const baseProgress = Math.floor((i / files.value.length) * 90);
              const fileProgress = Math.floor((progress / 100) * (90 / files.value.length));
              updateProgress(baseProgress + fileProgress);
            }
          );
          
          uploadResults.push(result);
        } catch (error) {
          console.error(`上傳文件 ${file.name} 失敗:`, error);
          uploadResults.push({
            success: false,
            error: true,
            fileName: file.name,
            message: `上傳失敗: ${error.message}`
          });
        }
      }
    }
    
    // 將上傳結果添加到列表
    uploadedData.value = [...uploadedData.value, ...uploadResults];
    
    updateProgress(100);
    
    showSnackbar(`成功上傳 ${files.value.length} 個檔案至病人 ${patientId.value}！`, 'success');
    files.value = []; // 清空選擇的文件
    patientId.value = ''; // 清空病人身分證號
    previewDialog.value = false; // 關閉預覽
    parsedData.value = null; // 清空解析數據
  } catch (error) {
    showSnackbar(`上傳失敗：${error.message}`, 'error');
  } finally {
    isUploading.value = false;
    uploadProgress.value = 0;
  }
};

const getFileIcon = (file) => {
  if (!file || !file.type) return 'mdi-file';
  
  const fileType = file.type;
  const fileName = file.name || '';
  
  if (fileType.includes('excel') || fileType.includes('spreadsheet')) {
    return 'mdi-file-excel';
  } else if (fileType.includes('pdf')) {
    return 'mdi-file-pdf';
  } else if (fileType.includes('csv')) {
    return 'mdi-file-delimited';
  } else if (fileType.includes('json') || fileName.toLowerCase().endsWith('.json')) {
    return 'mdi-code-json';
  } else {
    return 'mdi-file';
  }
};

const formatFileSize = (size) => {
  if (size < 1024) {
    return size + ' B';
  } else if (size < 1024 * 1024) {
    return (size / 1024).toFixed(2) + ' KB';
  } else {
    return (size / (1024 * 1024)).toFixed(2) + ' MB';
  }
};

const handleLogout = () => {
  authStore.logout();
};

// 格式化 JSON 數據顯示
const formatJSONDisplay = (data) => {
  if (typeof data !== 'object' || data === null) {
    return data;
  }
  
  // 將最頂層的物件或陣列展平成鍵值對
  const flattenedData = {};
  
  if (Array.isArray(data)) {
    // 處理陣列
    data.forEach((item, index) => {
      flattenedData[`項目 ${index + 1}`] = typeof item === 'object' ? JSON.stringify(item) : item;
    });
  } else {
    // 處理物件
    Object.keys(data).forEach(key => {
      const value = data[key];
      if (typeof value === 'object' && value !== null) {
        flattenedData[key] = JSON.stringify(value);
      } else {
        flattenedData[key] = value;
      }
    });
  }
  
  return flattenedData;
};
</script>

<template>
  <div class="dashboard-page">
    <v-container class="py-2">
      <v-row justify="center">
        <v-col cols="12" sm="10" md="10" lg="10" xl="10">
          <!-- 主要卡片容器 -->
          <v-card class="main-card" elevation="0">
            <!-- 頂部標題區 -->
            <div class="header-section">
              <div class="d-flex align-center mb-6">
                <v-icon size="44" class="header-icon">mdi-hospital-building</v-icon>
                <div class="ml-6">
                  <h1 class="header-title">健檢中心</h1>
                </div>
              </div>
            </div>

            <!-- 病人身分證字號輸入區 -->
            <v-row class="mb-8">
              <v-col cols="12" md="8" lg="9" xl="10">
                <v-text-field
                  v-model="patientId"
                  label="病人身分證字號"
                  placeholder="例：A123456789"
                  :rules="patientIdRules"
                  variant="outlined"
                  density="comfortable"
                  bg-color="#f5f5f5"
                  class="id-input"
                  hide-details="auto"
                >
                  <template v-slot:prepend>
                    <v-icon size="36">mdi-account-outline</v-icon>
                  </template>
                </v-text-field>
              </v-col>
              <v-col cols="12" md="4" lg="3" xl="2" class="d-flex justify-end align-center">
                <v-btn
                  class="logout-btn"
                  elevation="0"
                  @click="handleLogout"
                >
                  <v-icon start size="24">mdi-logout-variant</v-icon>
                  登出
                </v-btn>
              </v-col>
            </v-row>

            <!-- 文件上傳區域 -->
            <v-card
              class="upload-zone mb-8"
              :class="{ 'dragging': isDragging }"
              elevation="0"
              @dragover.prevent="isDragging = true"
              @dragleave.prevent="isDragging = false"
              @drop="handleFileDrop"
            >
              <div class="upload-content">
                <v-icon size="64" class="upload-icon">mdi-cloud-upload-outline</v-icon>
                <h3 class="upload-title">拖曳健檢報告至此處上傳</h3>
                <p class="upload-text">
                  
                  <v-btn
                    class="select-btn"
                    elevation="0"
                    @click="$refs.fileInput.click()"
                  >
                    <v-icon start size="24">mdi-file-plus-outline</v-icon>
                    選擇檔案
                  </v-btn>
                </p>
                <p class="upload-hint">支援 Excel、CSV、PDF 與 JSON 格式</p>
                <input
                  ref="fileInput"
                  type="file"
                  multiple
                  class="d-none"
                  accept=".xlsx,.xls,.csv,.pdf,.json"
                  @change="handleFileSelect"
                />
              </div>
            </v-card>

            <!-- 已選擇文件列表 -->
            <v-expand-transition>
              <v-card
                v-if="files.length > 0"
                class="selected-files mb-8"
                elevation="0"
              >
                <div class="d-flex align-center justify-space-between pa-4">
                  <div class="d-flex align-center">
                    <v-icon size="28" class="mr-2">mdi-file-multiple-outline</v-icon>
                    <span class="selected-title">已選擇 {{ files.length }} 個檔案</span>
                  </div>
                  <v-btn
                    class="clear-btn"
                    elevation="0"
                    @click="clearFiles"
                  >
                    <v-icon size="24">mdi-close</v-icon>
                    清除
                  </v-btn>
                </div>

                <v-divider class="mx-4"></v-divider>

                <v-list class="file-list">
                  <v-list-item
                    v-for="(file, index) in files"
                    :key="index"
                    class="file-item"
                  >
                    <template v-slot:prepend>
                      <v-icon
                        :color="file.name.toLowerCase().endsWith('.json') ? '#6B7280' : '#374151'"
                        size="24"
                      >
                        {{ getFileIcon(file) }}
                      </v-icon>
                    </template>
                    <v-list-item-title class="file-name">{{ file.name }}</v-list-item-title>
                    <v-list-item-subtitle class="file-size">{{ formatFileSize(file.size) }}</v-list-item-subtitle>
                  </v-list-item>
                </v-list>

                <v-alert
                  v-if="parseError"
                  type="error"
                  class="mx-4 mt-4"
                  variant="tonal"
                  density="comfortable"
                >
                  {{ parseError }}
                </v-alert>

                <v-card-actions class="pa-4">
                  <v-spacer></v-spacer>
                  <v-btn
                    class="preview-btn"
                    :disabled="!files.length || !patientId"
                    @click="previewFiles"
                  >
                    <v-icon start size="24">mdi-eye-outline</v-icon>
                    預覽
                  </v-btn>
                </v-card-actions>
              </v-card>
            </v-expand-transition>

            <!-- 已上傳數據列表 -->
            <v-card class="uploaded-list" elevation="0">
              <div class="d-flex align-center mb-4">
                <v-icon size="28" class="mr-2">mdi-text-box-check-outline</v-icon>
                <h3 class="uploaded-title">已上傳健檢報告</h3>
              </div>

              <div v-if="loading" class="d-flex align-center justify-center pa-8">
                <v-progress-circular
                  indeterminate
                  color="#111827"
                  size="36"
                ></v-progress-circular>
                <span class="ml-4 text-body-1">載入中...</span>
              </div>

              <template v-else>
                <v-list v-if="uploadedData.length > 0" class="report-list">
                  <v-list-item
                    v-for="(item, index) in uploadedData"
                    :key="index"
                    class="report-item"
                  >
                    <template v-slot:prepend>
                      <v-icon size="24">mdi-file-check-outline</v-icon>
                    </template>
                    <v-list-item-title class="report-patient">
                      病人：{{ item.userId || item.patient_hash || '未知' }}
                    </v-list-item-title>
                    <v-list-item-subtitle class="report-info">
                      {{ item.fileName || '未命名' }} · {{ item.uploadTime || new Date().toLocaleString() }}
                    </v-list-item-subtitle>
                  </v-list-item>
                </v-list>

                <div v-else class="empty-state">
                  <v-icon size="64" class="empty-icon">mdi-text-box-outline</v-icon>
                  <p class="empty-text">目前沒有已上傳的報告</p>
                  <p class="empty-hint">拖曳檔案或點擊上方按鈕開始上傳</p>
                </div>
              </template>
            </v-card>
          </v-card>
        </v-col>
      </v-row>
    </v-container>

    <!-- 預覽對話框 -->
    <v-dialog
      v-model="previewDialog"
      max-width="800"
      persistent
      class="preview-dialog"
    >
      <v-card class="preview-card">
        <v-card-title class="preview-header">
          <v-icon start size="28">mdi-eye-outline</v-icon>
          預覽健檢報告
        </v-card-title>
        
        <v-card-subtitle class="preview-subtitle">
          病人身分證：{{ patientId }}
        </v-card-subtitle>

        <v-card-text class="preview-content">
          <v-expansion-panels v-if="parsedData">
            <v-expansion-panel
              v-for="(fileData, index) in parsedData"
              :key="index"
              class="preview-panel"
            >
              <v-expansion-panel-title class="preview-panel-title">
                <v-icon
                  class="mr-2"
                  size="24"
                  :color="fileData.fileType === 'json' ? '#6B7280' : '#374151'"
                >
                  {{ fileData.fileType === 'json' ? 'mdi-code-json' : 
                     fileData.fileType === 'pdf' ? 'mdi-file-pdf-outline' : 
                     fileData.fileType === 'csv' ? 'mdi-file-delimited-outline' : 
                     'mdi-file-excel-outline' }}
                </v-icon>
                {{ fileData.fileName }}
                <span class="preview-file-size">({{ fileData.fileSize }})</span>
              </v-expansion-panel-title>

              <v-expansion-panel-text>
                <v-table class="preview-table">
                  <thead>
                    <tr>
                      <th>項目</th>
                      <th>數值</th>
                    </tr>
                  </thead>
                  <tbody>
                    <tr v-for="(value, key) in formatJSONDisplay(fileData.reportData)" :key="key">
                      <td>{{ key }}</td>
                      <td>{{ value }}</td>
                    </tr>
                  </tbody>
                </v-table>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-card-text>

        <v-divider></v-divider>

        <v-card-actions class="preview-actions">
          <v-spacer></v-spacer>
          <v-btn
            class="cancel-btn"
            elevation="0"
            @click="previewDialog = false"
          >
            取消
          </v-btn>
          <v-btn
            class="confirm-btn"
            elevation="0"
            :loading="isUploading"
            :disabled="isUploading"
            @click="handleFileUpload"
          >
            <v-icon start size="24">mdi-cloud-upload-outline</v-icon>
            確認上傳
          </v-btn>
        </v-card-actions>

        <v-progress-linear
          v-if="isUploading"
          :value="uploadProgress"
          height="4"
          color="#111827"
          class="upload-progress"
        ></v-progress-linear>
      </v-card>
    </v-dialog>

    <!-- 通知提示 -->
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color === 'success' ? '#111827' : '#EF4444'"
      timeout="3000"
      location="top"
      class="notification"
    >
      {{ snackbar.message }}
      <template v-slot:actions>
        <v-btn
          variant="text"
          class="notification-close"
          @click="snackbar.show = false"
        >
          關閉
        </v-btn>
      </template>
    </v-snackbar>
  </div>
</template>

<style scoped>
/* 全局樣式 */
.dashboard-page {
  background-color: #F9F7F4;
  min-height: 100vh;
}

/* 主卡片容器 */
.main-card {
  border-radius: 28px !important;
  margin-top: 30px;
  background: white !important;
  padding: 1rem 5rem !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.03) !important;
}

/* 頂部標題區 */
.header-icon {
  color: #111827;
  margin: 20px 0 0 0;
}

.header-title {

  font-size: 2.5rem !important;
  font-weight: 900;
  color: #111827;
  margin: 20px 0 0 0;
  letter-spacing: -0.5px;
}

.header-subtitle {
  font-size: 1.2rem; 
  color: #6B7280;
  margin: 0.25rem 0 0;
}

/* 身分證輸入區 */
.id-input {
  border-radius: 16px !important;
}

.id-input :deep(.v-field__outline) {
  border-color: rgba(0, 0, 0, 0.05) !important;
}

.input-hint {
  font-size: 1.05rem;
  color: #888;
  margin-top: 0.5rem;
}

/* 上傳區域 */
.upload-zone {
  border: 2px dashed #E5E7EB !important;
  border-radius: 24px !important;
  background: white !important;
  transition: all 0.2s ease !important;
  min-height: 300px;
}

.upload-zone:hover {
  border-color: #111827 !important;
  background: #FAFAFA !important;
}

.upload-zone.dragging {
  border-color: #111827 !important;
  background: #F9FAFB !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05) !important;
}

.upload-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  text-align: center;
}

.upload-icon {
  color: #111827;
  margin-bottom: 1.5rem;
}

.upload-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #111827;
  margin-bottom: 1rem;
}

.upload-text {
  font-size: 1.2rem;
  color: #6B7280;
  margin-bottom: 1rem;
}

.upload-hint {
  font-size: 1.05rem;
  font-weight: 600;
  color: #888;
  margin: 0;
}

/* 選擇檔案按鈕 */
.select-btn {
  background-color: #00B8D9 !important;
  color: #111827 !important;
  border-radius: 16px !important;
  font-weight: 600 !important;
  padding: 0 28px !important;
  height: 52px !important;
  margin: 0 8px !important;
  transition: all 0.2s ease !important;
  font-size: 1.1rem !important;
}

.select-btn:hover {
  background-color: #00B8D9 !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 184, 217, 0.25) !important;
}

/* 已選擇文件列表 */
.selected-files {
  border-radius: 24px !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
}

.selected-title {
  font-size: 1.35rem;
  font-weight: 600;
  color: #111827;
}

.clear-btn {
  color: #6B7280 !important;
  height: 48px !important;
  font-size: 1.1rem !important;
}

.file-list {
  padding: 0.5rem 0 !important;
}

.file-item {
  padding: 0.75rem 1rem !important;
}

.file-name {
  font-size: 1.05rem !important;
  color: #111827 !important;
}

.file-size {
  font-size: 0.9rem !important;
  color: #6B7280 !important;
}

/* 預覽按鈕 */
.preview-btn {
  background-color: #00B8D9 !important;
  color: #111827 !important;
  border-radius: 16px !important;
  font-weight: 600 !important;
  padding: 0 28px !important;
  height: 52px !important;
  transition: all 0.2s ease !important;
  font-size: 1.1rem !important;
}

.preview-btn:hover {
  background-color: #00B8D9 !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 184, 217, 0.25) !important;
}

/* 已上傳列表 */
.uploaded-list {
  border-radius: 24px !important;
  padding: 2rem !important;
  border: 1px solid rgba(0, 0, 0, 0.05) !important;
}

.uploaded-title {
  font-size: 1.5rem;
  font-weight: 700;
  color: #111827;
  margin: 0;
}

.report-list {
  padding: 0 !important;
}

.report-item {
  padding: 1rem !important;
  border-bottom: 1px solid #E5E7EB !important;
}

.report-item:last-child {
  border-bottom: none !important;
}

.report-patient {
  font-size: 1.05rem !important;
  color: #111827 !important;
}

.report-info {
  font-size: 0.9rem !important;
  color: #6B7280 !important;
}

/* 空狀態 */
.empty-state {
  text-align: center;
  padding: 4rem 0;
}

.empty-icon {
  color: #9CA3AF;
  margin-bottom: 1rem;
}

.empty-text {
  font-size: 1.35rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 0.5rem;
}

.empty-hint {
  font-size: 1.05rem;
  font-weight: 600;
  color: #6B7280;
  margin: 0;
}

/* 登出按鈕 */
.logout-btn {
  background-color: #00B8D9 !important;
  color: #111827 !important;
  font-weight: 600 !important;
  padding: 0 28px !important;
  height: 60px !important;
  border-radius: 16px !important;
  margin-top: -36px;
  transition: all 0.2s ease !important;
  font-size: 1.1rem !important;
}

.logout-btn:hover {
  background-color: #00B8D9 !important;
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 184, 217, 0.25) !important;
}

/* 預覽對話框 */
.preview-dialog :deep(.v-overlay__content) {
  border-radius: 28px !important;
  overflow: hidden !important;
}

.preview-card {
  border-radius: 28px !important;
}

.preview-header {
  font-size: 1.5rem !important;
  font-weight: 700 !important;
  color: #111827 !important;
  padding: 1.5rem !important;
}

.preview-subtitle {
  font-size: 1.2rem !important;
  color: #6B7280 !important;
  padding: 0 1.5rem 1rem !important;
}

.preview-content {
  padding: 1.5rem !important;
}

.preview-panel {
  border-radius: 16px !important;
  overflow: hidden !important;
  border: 1px solid #E5E7EB !important;
  margin-bottom: 0.5rem !important;
}

.preview-panel-title {
  font-size: 1.05rem !important;
  color: #111827 !important;
}

.preview-file-size {
  color: #6B7280;
  margin-left: 0.5rem;
}

.preview-table {
  border: 1px solid #E5E7EB !important;
  border-radius: 8px !important;
  overflow: hidden !important;
}

.preview-table th {
  background: #F9FAFB !important;
  color: #374151 !important;
  font-weight: 600 !important;
  font-size: 1.05rem !important;
  padding: 0.75rem 1rem !important;
}

.preview-table td {
  padding: 0.75rem 1rem !important;
  color: #111827 !important;
  font-size: 1.05rem !important;
}

.preview-actions {
  padding: 1rem 1.5rem !important;
}

.cancel-btn {
  color: #6B7280 !important;
  background-color: #00B8D9 !important;
  height: 48px !important;
  font-size: 1.1rem !important;
}

.confirm-btn {
  background-color: #00B8D9 !important;
  color: #111827 !important;
  border-radius: 16px !important;
  font-weight: 600 !important;
  margin-left: 0.5rem !important;
  height: 48px !important;
  font-size: 1.1rem !important;
}

.upload-progress {
  border-radius: 0 0 28px 28px !important;
}

/* RWD 適配 */
@media (max-width: 960px) {
  .main-card {
    padding: 1.5rem !important;
  }
  
  .header-title {
    font-size: 1.8rem;
  }
  
  .upload-content {
    padding: 2rem;
  }
  
  .upload-title {
    font-size: 1.35rem;
  }
}

@media (max-width: 600px) {
  .upload-zone {
    min-height: 240px;
  }
  
  .upload-content {
    padding: 1.5rem;
  }
  
  .select-btn,
  .preview-btn,
  .logout-btn {
    width: 100%;
  }
  
  .preview-dialog :deep(.v-overlay__content) {
    width: 90vw !important;
  }
  
  .upload-title {
    font-size: 1.35rem;
  }
}
</style>