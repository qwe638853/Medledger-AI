import apiClient, { handleApiError, notifyError, notifySuccess } from './apiService';
import { v4 as uuidv4 } from 'uuid';

/**
 * 獲取用戶已上傳的健康檢查數據
 * @returns {Promise} - 包含用戶健康檢查數據的Promise
 */
export const fetchUserHealthData = async () => {
  try {
    // 正確串接後端 HandleListMyReports API
    const response = await apiClient.get('/v1/reports');
    
    // 根據 HandleListMyReports 的回傳結構處理數據
    if (response.data && response.data.reports) {
      return response.data.reports; // 回傳報告陣列
    }
    return [];
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取健康數據');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 獲取可授權對象列表
 * @returns {Promise<Array>} - 包含可授權對象的Promise
 */
export const fetchAuthorizeTargets = async () => {
  try {
    // 待後端實現授權目標 API
    const response = await apiClient.get('/v1/auth/targets');
    return response.data.targets || [];
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取授權對象列表');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 授權健康數據給指定對象
 * @param {string} targetId - 授權對象ID
 * @param {Array} healthData - 要授權的健康數據
 * @returns {Promise} - 授權結果的Promise
 */
export const authorizeHealthData = async (targetId, healthData) => {
  try {
    // 待後端實現授權 API
    const response = await apiClient.post('/v1/auth/authorize', {
      targetId,
      reportIds: healthData.map(data => data.id)
    });
    
    if (response.data.success) {
      notifySuccess('授權成功！');
      return response.data;
    } else {
      throw new Error(response.data.message || '授權失敗');
    }
  } catch (error) {
    const errorMsg = handleApiError(error, '授權健康數據');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 使用 LLM 分析健康數據
 * @param {Array} healthData - 要分析的健康數據
 * @returns {Promise<string>} - 分析結果的Promise
 */
export const analyzeLLMSummary = async (healthData) => {
  try {
    // 待後端實現 LLM 分析 API
    const response = await apiClient.post('/v1/llm/analyze', {
      reportIds: healthData.map(data => data.id)
    });
    
    if (response.data.success) {
      return response.data.summary;
    } else {
      throw new Error(response.data.message || 'LLM 分析失敗');
    }
  } catch (error) {
    const errorMsg = handleApiError(error, 'LLM 分析');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 獲取他人的健康檢查數據（適用於醫療機構等有權限查看的角色）
 * @param {string} userId - 用戶ID
 * @returns {Promise} - 包含健康檢查數據的Promise
 */
export const fetchOtherHealthData = async (userId) => {
  try {
    const response = await apiClient.get(`/default/health-check/other/${userId}`);
    return response.data;
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取健康數據');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 解析健康檢查報告檔案內容
 * @param {File} file - 健康檢查報告檔案
 * @returns {Promise<Object>} - 解析後的健康數據
 */
const parseHealthReportFile = async (file) => {
  // 這裡應該實現實際的檔案解析邏輯
  // 根據檔案類型（Excel、CSV、PDF等）進行不同的處理
  // 為了示範，我們返回模擬數據
  
  return new Promise((resolve) => {
    setTimeout(() => {
      // 模擬解析結果，實際應用中應該從檔案內容提取
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

/**
 * 上傳健康檢查報告
 * @param {File} file - 要上傳的文件
 * @param {string} userId - 用戶ID
 * @param {Function} progressCallback - 上傳進度回調函數
 * @returns {Promise} - 上傳結果的Promise
 */
export const uploadHealthReport = async (file, userId, progressCallback) => {
  try {
    // 檢查 localStorage 中是否存在令牌
    const token = localStorage.getItem('token');
    if (!token) {
      throw new Error('未登入或會話已過期，請重新登入');
    }
    
    // 模擬上傳進度
    if (progressCallback) {
      const simulateProgress = () => {
        let progress = 0;
        const interval = setInterval(() => {
          progress += 5;
          if (progress >= 90) {
            clearInterval(interval);
          }
          progressCallback(progress);
        }, 100);
        return () => clearInterval(interval);
      };
      const clearProgressSimulation = simulateProgress();
      
      // 解析報告檔案
      const testResults = await parseHealthReportFile(file);
      
      // 生成唯一報告ID
      const reportId = `report_${uuidv4().substring(0, 8)}`;
      
      // 將數據轉為JSON字符串
      const testResultsJson = JSON.stringify(testResults);
      
      // 準備上傳到區塊鏈的請求數據
      const uploadData = {
        report_id: reportId,
        patient_hash: userId, // 使用用戶ID作為患者哈希
        test_results_json: testResultsJson
      };
      
      // 手動設置授權標頭
      const response = await apiClient.post('/v1/upload', uploadData, {
        headers: {
          'Authorization': token,
          'Content-Type': 'application/json'
        }
      });
      
      // 清除進度模擬
      clearProgressSimulation();
      progressCallback(100);
      
      console.log('上傳報告響應:', response.data);
      
      // 結合API響應和解析的數據
      return {
        reportId: reportId,
        userId: userId,
        fileName: file.name,
        fileSize: file.size,
        uploadTime: new Date().toISOString(),
        success: response.data.success,
        message: response.data.message,
        testResults: testResults
      };
    }
  } catch (error) {
    const errorMsg = handleApiError(error, '上傳健康報告');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 批量上傳健康檢查報告
 * @param {Array<File>} files - 要上傳的文件數組
 * @param {string} userId - 用戶ID
 * @param {Function} progressCallback - 上傳進度回調函數
 * @param {Function} fileCompletedCallback - 單個文件上傳完成回調
 * @returns {Promise} - 所有上傳結果的Promise
 */
export const batchUploadHealthReports = async (
  files, 
  userId, 
  progressCallback,
  fileCompletedCallback
) => {
  const results = [];
  
  for (let i = 0; i < files.length; i++) {
    try {
      const result = await uploadHealthReport(files[i], userId, progressCallback);
      results.push(result);
      
      if (fileCompletedCallback) {
        fileCompletedCallback(result, i, files.length);
      }
    } catch (error) {
      console.error(`上傳文件 ${files[i].name} 失敗:`, error);
      // 將錯誤添加到結果中，但繼續上傳其他文件
      results.push({
        error: true,
        fileName: files[i].name,
        message: error.message || '上傳失敗'
      });
    }
  }
  
  return results;
};

// 導出健康檢查服務對象
export default {
  fetchUserHealthData,
  fetchAuthorizeTargets,
  authorizeHealthData,
  analyzeLLMSummary,
  fetchOtherHealthData,
  uploadHealthReport,
  batchUploadHealthReports
}; 