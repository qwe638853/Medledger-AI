<script setup>
import { ref, onMounted } from 'vue';
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
  // 過濾出支持的文件類型（例如Excel、CSV、PDF等）
  const supportedFiles = droppedFiles.filter(file => {
    const fileType = file.type;
    return fileType.includes('excel') || 
           fileType.includes('spreadsheet') || 
           fileType.includes('csv') || 
           fileType.includes('pdf');
  });
  
  if (supportedFiles.length === 0) {
    showSnackbar('請上傳支持的檔案格式（Excel、CSV或PDF）', 'error');
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
};

const handleFileUpload = async () => {
  if (!files.value.length) {
    showSnackbar('請選擇檔案', 'error');
    return;
  }

  isUploading.value = true;
  uploadProgress.value = 0;
  
  try {
    const updateProgress = (progress) => {
      uploadProgress.value = progress;
    };
    
    const fileCompletedCallback = (result) => {
      // 將新上傳的數據添加到列表
      uploadedData.value.push(result);
    };
    
    await healthCheckService.batchUploadHealthReports(
      files.value,
      currentUser.value,
      updateProgress,
      fileCompletedCallback
    );
    
    showSnackbar(`成功上傳 ${files.value.length} 個檔案！`, 'success');
    files.value = []; // 清空選擇的文件
  } catch (error) {
    showSnackbar(`上傳失敗：${error.message}`, 'error');
  } finally {
    isUploading.value = false;
    uploadProgress.value = 0;
  }
};

const getFileIcon = (file) => {
  const fileType = file.type || '';
  if (fileType.includes('excel') || fileType.includes('spreadsheet')) {
    return 'mdi-file-excel';
  } else if (fileType.includes('pdf')) {
    return 'mdi-file-pdf';
  } else if (fileType.includes('csv')) {
    return 'mdi-file-delimited';
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
                    <p class="text-caption text-grey">支援 Excel、CSV、PDF 格式</p>
                    <input
                      ref="fileInput"
                      type="file"
                      multiple
                      class="d-none"
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
                        <v-icon :color="getFileIcon(file) === 'mdi-file-pdf' ? 'red' : 'green'">
                          {{ getFileIcon(file) }}
                        </v-icon>
                      </template>
                      <v-list-item-title>{{ file.name }}</v-list-item-title>
                      <v-list-item-subtitle>{{ formatFileSize(file.size) }}</v-list-item-subtitle>
                    </v-list-item>
                  </v-list>
                  <v-card-actions>
                    <v-spacer></v-spacer>
                    <v-btn
                      color="primary"
                      :loading="isUploading"
                      :disabled="isUploading"
                      @click="handleFileUpload"
                    >
                      <v-icon left>mdi-cloud-upload</v-icon>
                      上傳
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
                          用戶：{{ item.user_id || currentUser }}
                        </v-list-item-title>
                        <v-list-item-subtitle>
                          數據：{{ item.content || JSON.stringify(item) }}
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