import { ref } from 'vue';
//處理頁面導航相關的邏輯
//包含頁腳顯示、菜單項目等
export function useNavigation() {
    const showFooter = ref(false);
    const menuItems = ref([
        { title: '首頁', path: '/' },
        { title: '健康紀錄', path: '/records' },
        { title: '聯絡我們', path: '/contact' }
    ]);

    const goToHome = (showLoginForm) => {
        showLoginForm.value = false;
    };

    return {
        showFooter,
        menuItems,
        goToHome
    };
} 