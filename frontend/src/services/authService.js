import apiClient, { handleApiError, notifySuccess, notifyError } from './apiService';

/**
 * 用戶登入
 * @param {Object} data - 登入資料
 * @param {string} data.username - 使用者名稱
 * @param {string} data.password - 密碼
 * @returns {Promise} - 包含登入結果的 Promise
 */
export const login = async (data) => {
  try {
    console.log('發送登入請求:', { 
      data: { user_id: data.username, password: data.password }
    });
    
    const response = await apiClient.post('/v1/login', {
      user_id: data.username,
      password: data.password
    });
    
    console.log('登入響應:', {
      success: response.data.success,
      message: response.data.message,
      hasToken: !!response.data.token,
      tokenPrefix: response.data.token ? response.data.token.substring(0, 10) + '...' : 'N/A'
    });
    
    return response.data;
  } catch (error) {
    const errorMsg = handleApiError(error, '登入');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 用戶註冊
 * @param {Object} data - 註冊資料
 * @param {string} data.username - 使用者名稱
 * @param {string} data.password - 密碼
 * @param {string} data.name - 姓名
 * @param {string} data.date - 出生日期
 * @param {string} data.email - 電子郵件
 * @param {string} data.phone - 電話號碼
 * @returns {Promise} - 包含註冊結果的 Promise
 */
export const register = async (data) => {
  try {
    const response = await apiClient.post('/v1/register/user', {
      user_id: data.username,
      password: data.password,
      name: data.name,
      date: data.date,
      email: data.email,
      phone: data.phone
    });
    
    return response.data;
  } catch (error) {
    const errorMsg = handleApiError(error, '註冊');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 保險業者註冊
 * @param {Object} data - 註冊資料
 * @param {string} data.insurerId - 保險業者ID
 * @param {string} data.password - 密碼
 * @param {string} data.companyName - 公司名稱
 * @param {string} data.contactPerson - 聯絡人
 * @param {string} data.email - 電子郵件
 * @param {string} data.phone - 電話號碼
 * @returns {Promise} - 包含註冊結果的 Promise
 */
export const registerInsurer = async (data) => {
  try {
    console.log('發送保險業者註冊請求:', { 
      insurerId: data.insurerId,
      hasPassword: !!data.password,
      companyName: data.companyName
    });
    
    const response = await apiClient.post('/v1/register/insurer', {
      insurer_id: data.insurerId,
      password: data.password,
      company_name: data.companyName,
      contact_person: data.contactPerson,
      email: data.email,
      phone: data.phone
    });
    
    return response.data;
  } catch (error) {
    const errorMsg = handleApiError(error, '保險業者註冊');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 忘記密碼
 * @param {Object} data - 忘記密碼資料
 * @param {string} data.id_number - 使用者 ID
 * @param {string} data.role - 用戶角色
 * @returns {Promise} - 包含忘記密碼結果的 Promise
 */
export const forgotPassword = async (data) => {
  try {
    const response = await apiClient.post('/default/forget-password', {
      id_number: data.id_number,
      role: data.role
    });
    
    notifySuccess(response.data.message || '已發送重設密碼郵件，請檢查您的電子郵件');
    return response.data;
  } catch (error) {
    const errorMsg = handleApiError(error, '忘記密碼');
    notifyError(errorMsg);
    throw error;
  }
};

export default {
  login,
  register,
  registerInsurer,
  forgotPassword
}; 