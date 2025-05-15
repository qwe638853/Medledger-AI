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
  <v-container>
    <v-row>
      <v-col cols="12">
        <v-card class="mb-6">
          <v-card-title class="text-h5 primary--text">
            <v-icon large color="primary" class="mr-2">mdi-hospital-building</v-icon>
            健檢中心儀表板
          </v-card-title>
          <v-card-subtitle>
            歡迎，{{ currentUser }}
          </v-card-subtitle>
          <v-card-text>
            <!-- 病人身分證字號輸入 -->
            <v-row>
              <v-col cols="12" md="6">
                <v-text-field
                  v-model="patientId"
                  label="病人身分證字號"
                  placeholder="例：A123456789"
                  :rules="patientIdRules"
                  outlined
                  dense
                  hint="請輸入病人身分證字號，格式為一個大寫英文字母後跟九個數字"
                  persistent-hint
                >
                  <template v-slot:prepend>
                    <v-icon>mdi-account-badge</v-icon>
                  </template>
                </v-text-field>
              </v-col>
            </v-row>
            
            <v-row>
              <v-col cols="12">
                <!-- 文件上傳區域 -->
                <v-card
                  class="file-drop-zone"
                  :class="{ 'dragging': isDragging }"
                  @dragover.prevent="isDragging = true"
                  @dragleave.prevent="isDragging = false"
                  @drop="handleFileDrop"
                >
                  <div class="d-flex flex-column align-center justify-center pa-6">
                    <v-icon size="64" color="primary">mdi-cloud-upload</v-icon>
                    <h3 class="text-h5 mt-4 primary--text">拖曳健檢報告至此處上傳</h3>
                    <p class="text-body-1 text-center my-4">
                      或者
                      <v-btn
                        color="primary"
                        rounded
                        prepend-icon="mdi-file-plus"
                        @click="$refs.fileInput.click()"
                      >
                        選擇檔案
                      </v-btn>
                    </p>
                    <p class="text-caption text-grey">支援 Excel、CSV、PDF 與 JSON 格式</p>
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
              </v-col>
            </v-row>
            
            <!-- 已選擇文件列表 -->
            <v-row v-if="files.length > 0" class="mt-4">
              <v-col cols="12">
                <v-card outlined>
                  <v-card-title class="text-h6">
                    <v-icon left>mdi-file-multiple</v-icon>
                    已選擇 {{ files.length }} 個檔案
                    <v-spacer></v-spacer>
                    <v-btn color="error" variant="text" @click="clearFiles">
                      <v-icon>mdi-delete</v-icon>
                      清除
                    </v-btn>
                  </v-card-title>
                  <v-list>
                    <v-list-item v-for="(file, index) in files" :key="index">
                      <template v-slot:prepend>
                        <v-icon :color="
                          file.name.toLowerCase().endsWith('.json') ? 'deep-purple' : 
                          getFileIcon(file) === 'mdi-file-pdf' ? 'red' : 'green'
                        ">
                          {{ getFileIcon(file) }}
                        </v-icon>
                      </template>
                      <v-list-item-title>{{ file.name }}</v-list-item-title>
                      <v-list-item-subtitle>{{ formatFileSize(file.size) }}</v-list-item-subtitle>
                    </v-list-item>
                  </v-list>
                  
                  <!-- 解析錯誤提示 -->
                  <v-alert
                    v-if="parseError"
                    type="error"
                    class="ma-3"
                    dense
                  >
                    {{ parseError }}
                  </v-alert>
                  
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <!-- 修改為預覽按鈕 -->
                    <v-btn
                      color="info"
                      :disabled="!files.length || !patientId"
                      @click="previewFiles"
                      class="mr-2"
                    >
                      <v-icon left>mdi-eye</v-icon>
                      預覽
                    </v-btn>
                  </v-card-actions>
                </v-card>
              </v-col>
            </v-row>
            
            <!-- 已上傳數據列表 -->
            <v-row class="mt-4">
              <v-col cols="12">
                <v-card outlined>
                  <v-card-title class="text-h6">
                    <v-icon left>mdi-clipboard-text</v-icon>
                    已上傳健檢報告
                  </v-card-title>
                  <v-card-text>
                    <div v-if="loading">
                      <v-progress-circular
                        indeterminate
                        color="primary"
                      ></v-progress-circular>
                      加載中...
                    </div>
                    <v-list v-else-if="uploadedData.length > 0">
                      <v-list-item v-for="(item, index) in uploadedData" :key="index">
                        <template v-slot:prepend>
                          <v-icon color="primary">mdi-file-check</v-icon>
                        </template>
                        <v-list-item-title>
                          病人：{{ item.userId || item.patient_hash || '未知' }}
                        </v-list-item-title>
                        <v-list-item-subtitle>
                          檔案：{{ item.fileName || '未命名' }}
                        </v-list-item-subtitle>
                        <v-list-item-subtitle>
                          時間：{{ item.uploadTime || new Date().toLocaleString() }}
                        </v-list-item-subtitle>
                      </v-list-item>
                    </v-list>
                    <div v-else class="text-center pa-4">
                      <v-icon large color="grey">mdi-folder-open</v-icon>
                      <p class="text-body-1 mt-2">暫無已上傳資料</p>
                    </div>
                  </v-card-text>
                </v-card>
              </v-col>
            </v-row>
          </v-card-text>
          <v-card-actions>
            <v-spacer></v-spacer>
            <v-btn color="primary" variant="outlined" @click="handleLogout">
              <v-icon left>mdi-logout</v-icon>
              登出
            </v-btn>
          </v-card-actions>
        </v-card>
      </v-col>
    </v-row>
    
    <!-- 預覽對話框 -->
    <v-dialog v-model="previewDialog" max-width="800" persistent>
      <v-card>
        <v-card-title class="text-h5">
          <v-icon left>mdi-eye</v-icon>
          預覽健檢報告
        </v-card-title>
        <v-card-subtitle>
          病人身分證：{{ patientId }}
        </v-card-subtitle>
        <v-card-text>
          <v-expansion-panels v-if="parsedData">
            <v-expansion-panel v-for="(fileData, index) in parsedData" :key="index">
              <v-expansion-panel-title>
                <v-icon class="mr-2" :color="
                  fileData.fileType === 'json' ? 'deep-purple' : 
                  fileData.fileType === 'pdf' ? 'red' : 'green'
                ">
                  {{ fileData.fileType === 'json' ? 'mdi-code-json' : 
                     fileData.fileType === 'pdf' ? 'mdi-file-pdf' : 
                     fileData.fileType === 'csv' ? 'mdi-file-delimited' : 
                     'mdi-file-excel' }}
                </v-icon>
                {{ fileData.fileName }} ({{ fileData.fileSize }})
              </v-expansion-panel-title>
              <v-expansion-panel-text>
                <v-simple-table>
                  <template v-slot:default>
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
                  </template>
                </v-simple-table>
              </v-expansion-panel-text>
            </v-expansion-panel>
          </v-expansion-panels>
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="grey" text @click="previewDialog = false">
            取消
          </v-btn>
          <v-btn
            color="primary"
            :loading="isUploading"
            :disabled="isUploading"
            @click="handleFileUpload"
          >
            <v-icon left>mdi-cloud-upload</v-icon>
            確認上傳
          </v-btn>
        </v-card-actions>
        <v-progress-linear
          v-if="isUploading"
          :value="uploadProgress"
          height="10"
          color="primary"
          striped
        ></v-progress-linear>
      </v-card>
    </v-dialog>
    
    <!-- 通知提示 -->
    <v-snackbar
      v-model="snackbar.show"
      :color="snackbar.color"
      timeout="5000"
    >
      {{ snackbar.message }}
      <template v-slot:actions>
        <v-btn
          variant="text"
          @click="snackbar.show = false"
        >
          關閉
        </v-btn>
      </template>
    </v-snackbar>
  </v-container>
</template>

<style scoped>
.file-drop-zone {
  border: 2px dashed #1976d2;
  border-radius: 8px;
  background-color: #f5f9ff;
  transition: all 0.3s ease;
  min-height: 250px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.file-drop-zone:hover {
  background-color: #e3f2fd;
  border-color: #0d47a1;
}

.file-drop-zone.dragging {
  background-color: #bbdefb;
  border-color: #1565c0;
  box-shadow: 0 0 10px rgba(25, 118, 210, 0.2);
}
</style>