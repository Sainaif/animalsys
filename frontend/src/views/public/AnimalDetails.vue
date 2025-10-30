<template>
  <div class="animal-details-page">
    <div class="container">
      <LoadingSpinner v-if="loading" />

      <div v-else-if="animal" class="animal-content">
        <BaseButton variant="outline" size="small" @click="goBack" class="back-button">
          {{ t('common.back') }}
        </BaseButton>

        <div class="animal-header">
          <div class="animal-photos">
            <div v-if="animal.photos && animal.photos.length > 0" class="photo-gallery">
              <img
                v-for="(photo, index) in animal.photos"
                :key="index"
                :src="photo"
                :alt="`${animal.name} - ${index + 1}`"
                class="photo-item"
                :class="{ 'photo-main': index === 0 }"
              />
            </div>
            <div v-else class="photo-placeholder">🐾</div>
          </div>

          <div class="animal-main-info">
            <h1 class="animal-name">{{ animal.name }}</h1>
            <div class="animal-badges">
              <span class="info-badge">{{ t(`animals.${animal.species}`) }}</span>
              <span class="info-badge">{{ animal.gender === 'male' ? t('animals.male') : t('animals.female') }}</span>
              <span class="info-badge">{{ animal.age }} {{ t('animals.years') }}</span>
            </div>

            <div class="info-grid">
              <div class="info-item">
                <span class="info-label">{{ t('animals.breed') }}:</span>
                <span class="info-value">{{ animal.breed }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">{{ t('animals.color') }}:</span>
                <span class="info-value">{{ animal.color }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">{{ t('animals.size') }}:</span>
                <span class="info-value">{{ t(`animals.size_${animal.size}`) }}</span>
              </div>
              <div class="info-item">
                <span class="info-label">{{ t('animals.weight') }}:</span>
                <span class="info-value">{{ animal.weight }} kg</span>
              </div>
            </div>

            <div class="adopt-action">
              <RouterLink to="/register">
                <BaseButton variant="primary" size="large">
                  ❤️ {{ t('animals.adoptMe') }}
                </BaseButton>
              </RouterLink>
            </div>
          </div>
        </div>

        <div class="animal-sections">
          <BaseCard>
            <template #header>{{ t('animals.description') }}</template>
            <p class="description-text">{{ animal.description }}</p>
          </BaseCard>

          <BaseCard>
            <template #header>{{ t('animals.healthInfo') }}</template>
            <div class="health-badges">
              <span v-if="animal.sterilized" class="health-badge">✓ {{ t('animals.sterilized') }}</span>
              <span v-if="animal.vaccinated" class="health-badge">✓ {{ t('animals.vaccinated') }}</span>
              <span v-if="animal.chipped" class="health-badge">✓ {{ t('animals.chipped') }}</span>
            </div>
          </BaseCard>

          <BaseCard v-if="animal.behavioral_notes">
            <template #header>{{ t('animals.behavioralNotes') }}</template>
            <p>{{ animal.behavioral_notes }}</p>
          </BaseCard>
        </div>
      </div>

      <EmptyState
        v-else
        icon="🔍"
        :title="t('animals.notFound')"
        :description="t('animals.notFoundMessage')"
      >
        <RouterLink to="/animals">
          <BaseButton variant="primary">
            {{ t('animals.viewAll') }}
          </BaseButton>
        </RouterLink>
      </EmptyState>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { API } from '../../api'
import { useNotificationStore } from '../../stores/notification'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseCard from '../../components/base/BaseCard.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'
import EmptyState from '../../components/base/EmptyState.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const notificationStore = useNotificationStore()

const animal = ref(null)
const loading = ref(false)

onMounted(() => {
  fetchAnimal()
})

async function fetchAnimal() {
  loading.value = true
  try {
    const response = await API.animals.getById(route.params.id)
    animal.value = response.data
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    loading.value = false
  }
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.animal-details-page {
  min-height: 100vh;
  padding: 2rem 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.back-button {
  margin-bottom: 1.5rem;
}

.animal-content {
  max-width: 1200px;
}

.animal-header {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 2rem;
  margin-bottom: 2rem;
}

.animal-photos {
  background-color: var(--bg-secondary);
  border-radius: 0.75rem;
  overflow: hidden;
}

.photo-gallery {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 0.5rem;
  padding: 0.5rem;
}

.photo-item {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 0.5rem;
}

.photo-main {
  grid-column: 1 / -1;
  height: 400px;
}

.photo-placeholder {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 6rem;
  opacity: 0.3;
}

.animal-main-info {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.animal-name {
  font-size: 2.5rem;
  font-weight: bold;
  color: var(--text-primary);
  margin: 0;
}

.animal-badges {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

.info-badge {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
  background-color: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.info-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.info-value {
  font-size: 1rem;
  color: var(--text-primary);
  font-weight: 600;
}

.adopt-action {
  padding-top: 1rem;
}

.animal-sections {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.description-text {
  line-height: 1.6;
  color: var(--text-primary);
  margin: 0;
}

.health-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
}

.health-badge {
  padding: 0.5rem 1rem;
  background-color: #d1fae5;
  color: #065f46;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
}

@media (max-width: 768px) {
  .animal-details-page {
    padding: 1rem 0;
  }

  .animal-header {
    grid-template-columns: 1fr;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .animal-name {
    font-size: 2rem;
  }
}
</style>
