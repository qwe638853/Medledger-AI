import apiClient, { handleApiError, notifyError, notifySuccess } from './apiService';

/**
 * 獲取用戶個人資料
 * @param {string} userId - 用戶ID
 * @returns {Promise} - 包含用戶個人資料的Promise
 */
export const getUserProfile = async (userId) => {
  try {
    const response = await apiClient.get(`/default/user/${userId}`);
    return response.data;
  } catch (error) {
    const errorMsg = handleApiError(error, '獲取用戶資料');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 更新用戶個人資料
 * @param {string} userId - 用戶ID
 * @param {Object} profileData - 更新的個人資料
 * @returns {Promise} - 包含更新結果的Promise
 */
export const updateUserProfile = async (userId, profileData) => {
  try {
    const response = await apiClient.put(`/default/user/${userId}`, profileData);
    notifySuccess('個人資料更新成功');
    return response.data;
  } catch (error) {
    const errorMsg = handleApiError(error, '更新用戶資料');
    notifyError(errorMsg);
    throw error;
  }
};

/**
 * 更改密碼
 * @param {string} userId - 用戶ID
 * @param {Object} passwordData - 密碼數據
 * @param {string} passwordData.oldPassword - 舊密碼
 * @param {string} passwordData.newPassword - 新密碼
 * @returns {Promise} - 包含更改結果的Promise
 */
export const changePassword = async (userId, passwordData) => {
  try {
    const response = await apiClient.post(`/default/user/${userId}/change-password`, passwordData);
    notifySuccess('密碼更改成功');
    return response.data;
  } catch (error) {
    const errorMsg = handleApiError(error, '更改密碼');
    notifyError(errorMsg);
    throw error;
  }
};

export default {
  getUserProfile,
  updateUserProfile,
  changePassword
}; 