<template>
  <div class="users-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.users') }}</h1>
      <BaseButton variant="primary">
        ➕ {{ t('users.addUser') }}
      </BaseButton>
    </div>

    <BaseCard>
      <EmptyState
        icon="👥"
        :title="t('users.comingSoon')"
        :description="t('users.comingSoonMessage')"
      />
    </BaseCard>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import EmptyState from '../../components/base/EmptyState.vue'

const { t } = useI18n()
</script>

<style scoped>
.users-page {
  max-width: 1400px;
  padding: 2rem;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin: 0;
}

@media (max-width: 768px) {
  .users-page {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
}
</style>
