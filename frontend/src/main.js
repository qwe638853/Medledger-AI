import { createApp } from 'vue';
import { createPinia } from 'pinia';
import router from './router'; // 導入路由配置
import App from './App.vue';
import 'vuetify/styles';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import '@mdi/font/css/materialdesignicons.css';
import { aliases, mdi } from 'vuetify/iconsets/mdi-svg';
import './style.css';
import { initializeApp as initAuth } from './stores';

const vuetify = createVuetify({
  components,
  directives,
  theme: {
    defaultTheme: 'light'
  }
});

const app = createApp(App);
const pinia = createPinia();
app.use(pinia);
app.use(router);
app.use(vuetify);

// 在應用程式啟動時初始化身份驗證
initAuth();

app.config.errorHandler = (err, vm, info) => {
  console.error('Vue error:', err, info);
  document.dispatchEvent(new CustomEvent('show-snackbar', {
    detail: { message: '應用程式發生錯誤，請稍後再試或聯繫管理員。', color: 'error' }
  }));
};

app.mount('#app');