import axios from 'axios';

// API 基礎 URL
const API_BASE_URL = 'http://localhost:8080'; // 使用本地 gRPC-Gateway 伺服器端口

// 創建一個 axios 實例
const apiClient = axios.create({
  baseURL: API_BASE_URL,
  timeout: 30000,
  headers: {
    'Accept': 'application/json',
    'Content-Type': 'application/json'
  }
});

// 設置請求攔截器添加授權標頭
export const setupAuthInterceptor = (token) => {
  console.log('設置認證攔截器，令牌前綴:', token ? token.substring(0, 10) + '...' : 'N/A');
  
  // 清除之前的所有攔截器
  apiClient.interceptors.request.handlers = [];
  
  // 添加新的攔截器
  apiClient.interceptors.request.use(
    (config) => {
      // 優先使用傳入的令牌
      const authToken = token || localStorage.getItem('token');
      
      if (authToken) {
        // 注意: 從後端代碼看，後端期望的是純 token，不包含 "Bearer " 前綴
        // 檢查是否需要清除 "Bearer " 前綴
        let tokenToUse = authToken;
        if (authToken.startsWith('Bearer ')) {
          tokenToUse = authToken.substring(7); // 移除 "Bearer " 前綴
        }
        
        config.headers.Authorization = tokenToUse;
        console.log(`請求 ${config.url} 已添加授權標頭，令牌前綴: ${tokenToUse.substring(0, 10)}...`);
      } else {
        console.warn(`請求 ${config.url} 無授權令牌`);
      }
      return config;
    },
    (error) => {
      console.error('請求攔截器錯誤:', error);
      return Promise.reject(error);
    }
  );
};

// 錯誤處理
export const handleApiError = (error, action) => {
  console.group(`API 錯誤 - ${action}`);
  
  if (error.response) {
    // 後端回應錯誤
    console.log('錯誤狀態碼:', error.response.status);
    console.log('後端回應標頭:', error.response.headers);
    console.log('後端回應資料:', error.response.data);
    return error.response.data?.message || `${action}失敗：伺服器錯誤 (${error.response.status})`;
  } else if (error.request) {
    // 請求已發送但未收到回應
    console.log('未收到回應:', error.request);
    return `${action}失敗：未收到伺服器回應，請檢查網路連接`;
  } else {
    // 請求設置出錯
    console.log('請求錯誤:', error.message);
    return `${action}失敗：${error.message}`;
  }
  
  console.log('完整錯誤物件:', error);
  console.groupEnd();
};

// 通知函數
export const notifySuccess = (message) => {
  document.dispatchEvent(new CustomEvent('show-snackbar', { detail: { message, color: 'success' } }));
};

export const notifyError = (message) => {
  document.dispatchEvent(new CustomEvent('show-snackbar', { detail: { message, color: 'error' } }));
};

export default apiClient; 