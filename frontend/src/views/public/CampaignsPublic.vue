<template>
  <div class="campaigns-public-page">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">{{ t('nav.campaigns') }}</h1>
        <p class="page-subtitle">{{ t('campaigns.publicSubtitle') }}</p>
      </div>

      <BaseCard>
        <EmptyState
          icon="🎯"
          :title="t('campaigns.comingSoon')"
          :description="t('campaigns.comingSoonMessage')"
        />
      </BaseCard>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import BaseCard from '../../components/base/BaseCard.vue'
import EmptyState from '../../components/base/EmptyState.vue'

const { t } = useI18n()
</script>

<style scoped>
.campaigns-public-page {
  min-height: 100vh;
  padding: 2rem 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.page-header {
  text-align: center;
  margin-bottom: 3rem;
}

.page-title {
  font-size: 2.5rem;
  font-weight: bold;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.page-subtitle {
  font-size: 1.25rem;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }
}
</style>
