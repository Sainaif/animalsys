<template>
  <div class="error-page">
    <div class="error-content">
      <h1 class="error-code">403</h1>
      <h2 class="error-title">{{ t('errors.unauthorized') }}</h2>
      <p class="error-message">{{ t('errors.unauthorizedMessage') }}</p>
      <BaseButton variant="primary" @click="goBack">
        {{ t('common.back') }}
      </BaseButton>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import BaseButton from '../../components/base/BaseButton.vue'

const router = useRouter()
const { t } = useI18n()

function goBack() {
  router.back()
}
</script>

<style scoped>
.error-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background-color: var(--bg-primary);
}

.error-content {
  text-align: center;
  max-width: 500px;
}

.error-code {
  font-size: 6rem;
  font-weight: bold;
  color: #ef4444;
  margin: 0 0 1rem 0;
  line-height: 1;
}

.error-title {
  font-size: 2rem;
  margin: 0 0 1rem 0;
  color: var(--text-primary);
}

.error-message {
  font-size: 1.125rem;
  color: var(--text-secondary);
  margin: 0 0 2rem 0;
}
</style>
