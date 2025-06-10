<template>
  <v-card flat class="mt-6 pa-4 rounded-lg">
    <div class="section-title">
      <h3 class="text-h5 font-weight-bold">
        <v-icon color="primary" class="me-2">mdi-account-multiple</v-icon>
        角色選擇
      </h3>
      <div class="text-body-2 text-grey">請選擇您的帳號類型</div>
      <v-divider class="mt-2"></v-divider>
    </div>
    
    <div class="role-cards-container">
      <div
        v-for="role in roles"
        :key="role.value"
        class="role-card"
        :class="{ 'role-card--selected': modelValue === role.value }"
        @click="$emit('update:modelValue', role.value)"
      >
        <div class="role-card__icon-wrapper">
          <v-icon 
            :color="modelValue === role.value ? 'primary' : 'grey'" 
            size="32"
          >{{ role.icon }}</v-icon>
        </div>
        <div class="role-card__content">
          <h3 class="role-card__title">{{ role.text }}</h3>
        </div>
        <v-scale-transition>
          <div v-if="modelValue === role.value" class="role-card__check">
            <v-icon color="primary" size="24">mdi-check-circle</v-icon>
          </div>
        </v-scale-transition>
      </div>
    </div>
    
    <div class="mt-6 text-center">
      <v-btn
        class="primary-btn"
        size="large"
        @click="$emit('next')"
        :disabled="!modelValue"
      >
        下一步
        <v-icon class="ms-2">mdi-arrow-right-circle</v-icon>
      </v-btn>
    </div>
  </v-card>
</template>

<script setup>
import { defineProps, defineEmits } from 'vue';

const props = defineProps({
  modelValue: {
    type: String,
    required: true
  }
});

defineEmits(['update:modelValue', 'next']);

const roles = [
  { 
    text: '一般用戶',
    value: 'user',
    icon: 'mdi-account'
  },
  { 
    text: '醫療機構',
    value: 'medical',
    icon: 'mdi-hospital-building'
  },
  { 
    text: '保險業者',
    value: 'insurer',
    icon: 'mdi-shield-account'
  }
];
</script>

<style scoped>
/* 角色選擇卡片 */
.role-cards-container {
  display: flex;
  flex-direction: row;
  gap: 1.5rem;
  width: 100%;
  margin: 1.5rem 0;
  justify-content: center;
  flex-wrap: nowrap;
  overflow-x: unset;
}

.role-card {
  position: relative;
  min-width: 260px;
  max-width: 300px;
  flex: 1 1 0;
  padding: 1.5rem 1.2rem;
  font-size: 1.18rem;
  border: 3px solid #ececec;
  border-radius: 20px;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  cursor: pointer;
}

.role-card:hover {
  border-color: #00B8D9;
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.15);
  transform: translateY(-2px);
}

.role-card:hover .role-card__icon-wrapper {
  background: rgba(0, 184, 217, 0.15);
}

.role-card:hover .v-icon {
  color: #00B8D9 !important;
}

.role-card:hover .role-card__title {
  color: #111827;
  font-weight: 700;
}

.role-card--selected {
  border-color: #00B8D9;
  background-color: rgba(0, 184, 217, 0.1);
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.25);
}

.role-card--selected .role-card__icon-wrapper {
  background: rgba(0, 184, 217, 0.2);
}

.role-card--selected .v-icon {
  color: #00B8D9 !important;
}

.role-card--selected .role-card__title {
  color: #111827;
  font-weight: 700;
}

.role-card__icon-wrapper {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  border-radius: 50%;
  background: #f3f4f6;
  margin-right: 1rem;
  margin-bottom: 0.7rem;
  font-size: 2.5rem;
  transition: all 0.3s ease;
}

.role-card__content {
  flex: 1;
}

.role-card__title {
  font-size: 1.7rem;
  font-weight: 600;
  color: #4B5563;
  transition: all 0.3s ease;
}

.role-card__check {
  position: absolute;
  top: 1rem;
  right: 1rem;
}

/* 按鈕樣式 */
.primary-btn {
  background-color: #00B8D9 !important;
  color: white !important;
  border: none !important;
  box-shadow: 0 4px 16px rgba(0, 184, 217, 0.25) !important;
  border-radius: 20px !important;
  text-transform: none !important;
  font-size: 1.5rem !important;
  font-weight: 600 !important;
  letter-spacing: 0 !important;
  height: 56px !important;
  min-width: 140px !important;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1) !important;
}

.primary-btn:hover {
  background-color: #0095B0 !important;
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(0, 184, 217, 0.35) !important;
}

/* RWD 適配 */
@media (max-width: 900px) {
  .section-title h3 {
    font-size: 1.7rem !important;
  }
  
  .section-title .text-caption {
    font-size: 1.18rem !important;
  }
  
  .role-card {
    padding: 1.2rem;
  }
  
  .role-card__title {
    font-size: 1.4rem;
  }
  
  .role-card__icon-wrapper {
    width: 48px;
    height: 48px;
    font-size: 2rem;
  }
}

@media (max-width: 600px) {
  .role-cards-container {
    flex-direction: column;
    align-items: stretch;
  }
  
  .role-card {
    max-width: none;
  }
}
</style> 