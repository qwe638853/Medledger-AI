<script setup>
import { ref, onMounted, computed } from 'vue';
import { useAuthStore } from '../stores';
import { healthCheckService, notifyError, notifySuccess } from '../services';

// 假設有這些服務
// import { authorizeService, llmSummaryService } from '../services';

const authStore = useAuthStore();
const currentUser = computed(() => authStore.currentUser);
const healthData = ref([]);
const loading = ref(false);

// 授權相關
const authorizeDialog = ref(false);
const authorizeTarget = ref('');
const authorizeTargets = ref(['醫師A', '醫院B', '家人C']); // 暫時使用假資料
const authorizing = ref(false);

// LLM 分析相關
const llmLoading = ref(false);
const llmSummary = ref('');

onMounted(async () => {
  loading.value = true;
  try {
    // 注意：fetchUserHealthData不需要傳userId參數，後端從JWT提取
    const healthResponse = await healthCheckService.fetchUserHealthData();
    healthData.value = healthResponse.map(report => ({
      id: report.reportId || report.id,
      content: report.content || JSON.stringify(report.testResults),
      date: report.timestamp || report.date || new Date().toISOString()
    }));
    
    // 暫時使用靜態授權對象列表，待授權API完成後替換
    // const targetsResponse = await healthCheckService.fetchAuthorizeTargets();
    // authorizeTargets.value = targetsResponse;
  } catch (error) {
    notifyError(`獲取健康數據失敗：${error.message}`);
    healthData.value = [];
  } finally {
    loading.value = false;
  }
});

const handleLogout = () => {
  authStore.logout();
};

// 處理授權行為
const handleAuthorize = async () => {
  if (!authorizeTarget.value || !healthData.value.length) {
    notifyError('請選擇授權對象且確保有健康數據');
    return;
  }
  
  authorizing.value = true;
  try {
    // 待後端 API 完成後實現
    // await healthCheckService.authorizeHealthData(
    //   authorizeTarget.value,
    //   healthData.value
    // );
    
    // 暫時使用模擬授權
    await new Promise(resolve => setTimeout(resolve, 800));
    notifySuccess(`已成功授權給 ${authorizeTarget.value}！`);
    authorizeDialog.value = false;
    authorizeTarget.value = '';
  } catch (error) {
    notifyError(`授權失敗：${error.message}`);
  } finally {
    authorizing.value = false;
  }
};

// 處理 LLM 分析
const handleLLMSummary = async () => {
  if (!healthData.value.length) {
    notifyError('無可分析的健康數據');
    return;
  }
  
  llmLoading.value = true;
  try {
    // 待後端 API 完成後實現
    // const summary = await healthCheckService.analyzeLLMSummary(healthData.value);
    
    // 暫時使用模擬 LLM 分析
    await new Promise(resolve => setTimeout(resolve, 1000));
    llmSummary.value = "【AI 健康摘要】根據您的健檢數據，血糖、血脂與血壓均在正常範圍內。建議維持均衡飲食和適度運動，每半年進行一次健康檢查。";
  } catch (error) {
    notifyError(`LLM 分析失敗：${error.message}`);
  } finally {
    llmLoading.value = false;
  }
};
</script>

<template>
  <v-container class="fill-height" fluid>
    <v-row align="center" justify="center">
      <v-col cols="12" sm="10" md="8" lg="7">
        <v-card class="pa-6 mb-6" elevation="8">
          <v-row align="center" justify="space-between" class="mb-4">
            <v-col cols="8">
              <h2 class="mb-1">👤 使用者儀表板</h2>
              <div class="subtitle-1">歡迎，{{ currentUser }}</div>
            </v-col>
            <v-col cols="4" class="d-flex justify-end align-center">
              <v-btn color="primary" @click="handleLogout" elevation="2">登出</v-btn>
            </v-col>
          </v-row>
        </v-card>

        <!-- 健檢報告列表 -->
        <v-card class="pa-4 mb-6" elevation="4">
          <h3 class="mb-3">健康檢查報告</h3>
          <v-data-table
            :headers="[
              { text: '報告編號', value: 'id', width: 120 },
              { text: '內容', value: 'content', width: 300 },
              { text: '日期', value: 'date', width: 150 }
            ]"
            :items="healthData"
            :loading="loading"
            loading-text="資料載入中..."
            class="elevation-0"
            dense
            hide-default-footer
            :no-data-text="'暫無資料'"
          >
            <template #item.content="{ item }">
              {{ item.content || item }}
            </template>
            <template #item.date="{ item }">
              {{ item.date || '-' }}
            </template>
          </v-data-table>
        </v-card>

        <!-- 授權區塊 -->
        <v-card class="pa-4 mb-6" elevation="4">
          <h3 class="mb-3">資料授權</h3>
          <v-row>
            <v-col cols="12" sm="8" md="6">
              <v-select
                v-model="authorizeTarget"
                :items="authorizeTargets"
                label="選擇授權對象"
                dense
                outlined
                clearable
              ></v-select>
            </v-col>
            <v-col cols="12" sm="4" md="6" class="d-flex align-end">
              <v-btn
                color="success"
                :disabled="!authorizeTarget || authorizing"
                @click="authorizeDialog = true"
                elevation="2"
              >
                <v-icon left>mdi-account-key</v-icon> 授權
              </v-btn>
            </v-col>
          </v-row>
        </v-card>

        <!-- LLM 分析區塊 -->
        <v-card class="pa-4 mb-6" elevation="4">
          <h3 class="mb-3">AI 健檢摘要分析</h3>
          <v-btn
            color="deep-purple accent-4"
            :loading="llmLoading"
            :disabled="llmLoading || !healthData.length"
            @click="handleLLMSummary"
            elevation="2"
            class="mb-3"
          >
            <v-icon left>mdi-robot</v-icon> 產生 AI 摘要
          </v-btn>
          <v-alert
            v-if="llmSummary"
            type="info"
            class="mt-3"
            border="left"
            colored-border
            elevation="1"
          >
            {{ llmSummary }}
          </v-alert>
        </v-card>
      </v-col>
    </v-row>

    <!-- 授權確認對話框 -->
    <v-dialog v-model="authorizeDialog" max-width="400">
      <v-card>
        <v-card-title class="headline">確認授權</v-card-title>
        <v-card-text>
          確定要將健康資料授權給 <b>{{ authorizeTarget }}</b> 嗎？
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn text @click="authorizeDialog = false">取消</v-btn>
          <v-btn color="success" :loading="authorizing" @click="handleAuthorize">確認</v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-container>
</template>

<style scoped>
.fill-height {
  min-height: 100vh;
  background: #f5f6fa;
}
h2 {
  font-size: 28px;
  font-weight: bold;
}
h3 {
  font-size: 20px;
  margin-top: 0;
  font-weight: 500;
}
.subtitle-1 {
  color: #666;
}
</style>