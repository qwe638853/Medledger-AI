import apiClient, { setupAuthInterceptor, handleApiError, notifySuccess, notifyError } from './apiService';
import authService from './authService';
import healthCheckService from './healthCheckService';
import userService from './userService';

export {
  apiClient,
  setupAuthInterceptor,
  handleApiError,
  notifySuccess,
  notifyError,
  authService,
  healthCheckService,
  userService
};

export default {
  api: apiClient,
  auth: authService,
  healthCheck: healthCheckService,
  user: userService,
  setupAuth: setupAuthInterceptor
}; 