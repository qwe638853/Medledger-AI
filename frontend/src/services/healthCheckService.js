import apiClient, { handleApiError, notifyError, notifySuccess } from './apiService';
import { v4 as uuidv4 } from 'uuid';
import axios from 'axios';

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
      
      console.log('上傳報告響應:', response);
      
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

/**
 * 上傳 JSON 格式的健康檢查報告數據
 * @param {string} patientId - 病人身分證
 * @param {Object} jsonData - 從 JSON 文件解析出的健康數據
 * @param {string} fileName - 原始文件名
 * @param {Function} progressCallback - 上傳進度回調函數
 * @returns {Promise} - 上傳結果的Promise
 */
export const uploadJsonHealthData = async (patientId, jsonData, fileName, progressCallback) => {
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
      
      // 生成唯一報告ID
      const reportId = `report_${uuidv4().substring(0, 8)}`;
      
      // 將原始 JSON 數據轉為字符串
      const testResultsJson = JSON.stringify(jsonData);
      
      // 準備上傳到區塊鏈的請求數據
      const uploadData = {
        report_id: reportId,
        user_id: patientId, // 使用病人身分證
        test_results_json: testResultsJson,
      };
      
      // 發送請求
      const response = await apiClient.post('/v1/upload', uploadData, {
        headers: {
          'Authorization': token,
          'Content-Type': 'application/json'
        }
      });
      
      // 清除進度模擬
      clearProgressSimulation();
      progressCallback(100);
      
      console.log('上傳 JSON 報告響應:', response.data);
      
      // 結合API響應和解析的數據
      return {
        reportId: reportId,
        patientId: patientId,
        fileName: fileName,
        uploadTime: new Date().toISOString(),
        success: response.data.success || true,
        message: response.data.message || '上傳成功',
        testResults: jsonData  // 保存原始解析的 JSON 數據
      };
    }
  } catch (error) {
    const errorMsg = handleApiError(error, '上傳健康報告');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 獲取特定病患的健康報告元數據（僅報告基本信息，不含健檢數據）
 * @param {string} patientId - 病患身分證字號
 * @returns {Promise<Array>} - 包含報告元數據的Promise
 */
export const fetchReportMetaByPatientID = async (patientId) => {
  try {
    if (!patientId || patientId.trim() === '') {
      throw new Error('病患身分證字號不能為空');
    }

    const response = await apiClient.get(`/v1/reports/meta/${patientId}`);
    console.log(response.data);
    if (response.data && response.data.reports) {
      // 格式化回傳的數據以符合前端顯示需求
      return response.data.reports.map(report => ({
        id: report.reportId,
        clinic_id: report.clinicId,
        content: `健康檢查報告 - ${report.reportId}`, // 不含詳細健檢數據
        date: new Date(Number(report.createdAt) * 1000).toISOString().split('T')[0], // 先轉換為數字再轉換為日期
        is_authorized: false // 預設為未授權
      }));
    }
    return [];
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取病患報告元數據');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 獲取已授權的健康報告列表
 * @returns {Promise<Array>} - 包含已授權報告的Promise
 */
export const fetchAuthorizedReports = async () => {
  try {
    const response = await apiClient.get('/v1/reports/authorized');
    
    if (response.data && response.data.reports) {
      return response.data.reports.map(report => ({
        id: report.reportId,
        patient_id: report.patientId,
        date: report.date,
        expiry: report.expiry
      }));
    }
    return [];
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取已授權報告');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 獲取儀表板統計數據
 * @returns {Promise<Object>} - 包含統計數據的Promise
 */
export const fetchDashboardStats = async () => {
  try {
    const response = await apiClient.get('/v1/dashboard/summary');
    
    if (response.data) {
      return {
        totalAuthorized: response.data.total_authorized,
        pendingRequests: response.data.pending_requests,
        totalPatients: response.data.total_patients
      };
    }
    return {
      totalAuthorized: 0,
      pendingRequests: 0,
      totalPatients: 0
    };
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取儀表板統計數據');
    notifyError(errorMsg);
    // 返回默認值，避免UI顯示錯誤
    return {
      totalAuthorized: 0,
      pendingRequests: 0,
      totalPatients: 0
    };
  }
};

/**
 * 請求病患健康報告授權
 * @param {string} reportId - 報告ID
 * @param {string} patientId - 病患身分證字號
 * @param {string} reason - 請求授權理由
 * @param {string} expiry - 授權到期日期
 * @returns {Promise<Object>} - 請求結果的Promise
 */
export const requestReportAccess = async (reportId, patientId, reason, expiry) => {
  try {
    const response = await apiClient.post('/v1/access/request', {
      report_id: reportId,
      patient_id: patientId,
      reason: reason,
      expiry: new Date(expiry).getTime() / 1000  // 轉換為 Unix 時間戳
    });
    
    if (response.data && response.data.success) {
      notifySuccess('已成功送出授權請求');
      return {
        success: true,
        requestId: response.data.request_id,
        message: '授權請求已送出'
      };
    } else {
      throw new Error(response.data ? response.data.message : '請求授權失敗');
    }
  } catch (error) {
    const errorMsg = handleApiError(error, '請求授權');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 獲取病患(我)收到的訪問請求
 * @returns {Promise<Array>} - 包含訪問請求的Promise
 */
export const fetchAccessRequests = async () => {
  try {
    const response = await apiClient.get('/v1/access/requests');
    console.log('獲取授權請求回應:', response);
    
    // 處理可能的空回應
    if (!response.data) {
      console.warn('授權請求回應為空');
      return [];
    }
    
    if (response.data && response.data.requests) {
      // 適應實際後端回應格式，正確映射欄位
      return response.data.requests.map(request => ({
        id: request.requestId , // 使用requesterId作為id
        reportId: request.reportId , // 使用reportId作為reportId
        requesterId: request.requesterId , // 使用requesterId作為requesterId
        requesterName: request.targetHash, // 提供fallback
        reason: request.reason || '',
        requestTime: request.requestedAt || request.request_time || 0, // 使用requestedAt作為requestTime
        status: request.status || 'UNKNOWN',
        expiry: request.expiry || 0
      }));
    }
    return [];
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取授權請求');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 同意授權請求
 * @param {string} requestId - 請求ID
 * @returns {Promise<Object>} - 同意授權結果的Promise
 */
export const approveAccessRequest = async (requestId) => {
  try {
    console.log('同意授權請求1:', requestId);
    const response = await apiClient.post('/v1/access/approve', {
      request_id: requestId
    });
    
    // 改進對回應的處理，處理可能的空回應
    console.log('授權回應:', response);
    
    // 即使後端返回空數據，也視為成功
    if (!response.data) {
      notifySuccess('已成功授權報告');
      return {
        success: true,
        message: '授權成功'
      };
    }
    
    if (response.data && response.data.success) {
      notifySuccess('已成功授權報告');
      return {
        success: true,
        message: '授權成功'
      };
    } else {
      throw new Error(response.data ? response.data.message : '授權失敗');
    }
  } catch (error) {
    const errorMsg = handleApiError(error, '同意授權');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 拒絕授權請求
 * @param {string} requestId - 請求ID
 * @returns {Promise<Object>} - 拒絕授權結果的Promise
 */
export const rejectAccessRequest = async (requestId) => {
  try {
    const response = await apiClient.post('/v1/access/reject', {
      request_id: requestId
    });
    
    // 改進對回應的處理，處理可能的空回應
    console.log('拒絕授權回應:', response);
    
    // 即使後端返回空數據，也視為成功
    if (!response.data) {
      notifySuccess('已拒絕授權請求');
      return {
        success: true,
        message: '已拒絕授權請求'
      };
    }
    
    if (response.data && response.data.success) {
      notifySuccess('已拒絕授權請求');
      return {
        success: true,
        message: '已拒絕授權請求'
      };
    } else {
      throw new Error(response.data ? response.data.message : '拒絕授權失敗');
    }
  } catch (error) {
    const errorMsg = handleApiError(error, '拒絕授權');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 獲取已授予的授權票據列表
 * @returns {Promise<Array>} - 包含授權票據的Promise
 */
export const fetchGrantedTickets = async () => {
  try {
    const response = await apiClient.get('/v1/access/granted');
    console.log('獲取已授權票據回應:', response);
    
    // 處理可能的空回應
    if (!response.data) {
      console.warn('已授權票據回應為空');
      return [];
    }
    
    if (response.data && response.data.tickets) {
      // 適應實際後端回應格式，正確映射欄位
      return response.data.tickets.map(ticket => ({
        reportId: ticket.reportId || ticket.report_id || '',
        targetId: ticket.targetId || ticket.target_id || ticket.requesterId || '', 
        targetName: ticket.targetName || ticket.target_name || ticket.targetId || ticket.requesterId || '未知對象',
        grantTime: ticket.grantedAt || ticket.grant_time || ticket.requestedAt || 0,
        expiry: ticket.expiry || 0
      }));
    }
    return [];
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取已授權票據');
    notifyError(errorMsg);
    throw error;
  }
};

// 獲取報告詳細內容
export const fetchReportContent = async (reportId, patientId) => {
  try {
    const response = await apiClient.get(`/v1/reports/authorized/${patientId}/${reportId}`, {
      headers: {
        'Authorization': `Bearer ${localStorage.getItem('token')}`
      }
    });
    return response.data;
  } catch (error) {
    console.error('獲取報告內容失敗:', error);
    throw error;
  }
}

// 導出健康檢查服務對象
export default {
  fetchUserHealthData,
  fetchAuthorizeTargets,
  authorizeHealthData,
  analyzeLLMSummary,
  fetchOtherHealthData,
  uploadHealthReport,
  batchUploadHealthReports,
  uploadJsonHealthData,
  fetchReportMetaByPatientID,
  fetchAuthorizedReports,
  fetchDashboardStats,
  requestReportAccess,
  fetchAccessRequests,
  approveAccessRequest,
  rejectAccessRequest,
  fetchGrantedTickets,
  fetchReportContent
}; 