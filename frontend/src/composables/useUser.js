import axios from 'axios';
//處理用戶相關的操作
//包含註冊、忘記密碼、獲取健康檢查數據等功能
export function useUser(token, currentUser) {
    // 忘記密碼
    const forgotPassword = async (data) => {
        try {
            const response = await axios.post('/default/forget-password', {
                id_number: data.id_number,
                role: data.role
            });
            alert(response.data.message || '已發送重設密碼郵件，請檢查您的電子郵件');
        } catch (error) {
            const errorMessage = error.response?.data?.detail || error.message;
            alert(`重設密碼失敗：${errorMessage}`);
        }
    };

    // 註冊
    const register = async (data) => {
        try {
            const response = await axios.post('/default/register', {
                id_number: data.id_number,
                password: data.password,
                full_name: data.full_name,
                gender: data.gender,
                birth_date: data.birth_date,
                phone_number: data.phone_number,
                email: data.email,
                role: data.role
            }, {
                headers: {
                    'Content-Type': 'application/json'
                }
            });
            alert(response.data || '註冊成功，請登入');
        } catch (error) {
            const errorMessage = error.response?.data?.detail || error.message;
            alert(`註冊失敗：${errorMessage}`);
        }
    };

    // 獲取健康檢查數據
    const fetchData = async () => {
        try {
            const response = await axios.get(`/default/health-check/other/${currentUser.value}`, {
                headers: {
                    'Authorization': `Bearer ${token.value}`,
                    'Accept': 'application/json'
                }
            });
            return response.data || [];
        } catch (error) {
            if (error.response?.status === 401) {
                alert('認證過期，請重新登入');
                return null;
            } else if (error.response?.status === 404) {
                alert('找不到數據端點，請確認後端服務是否正常運行');
            } else {
                console.error('Fetch data failed:', error);
                alert('獲取數據失敗，請重試');
            }
            return [];
        }
    };

    return {
        forgotPassword,
        register,
        fetchData
    };
} 