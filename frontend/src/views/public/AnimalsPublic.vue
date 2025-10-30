<template>
  <div class="animals-public-page">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">{{ t('animals.availableForAdoption') }}</h1>
        <p class="page-subtitle">{{ t('animals.findYourCompanion') }}</p>
      </div>

      <LoadingSpinner v-if="loading" />

      <div v-else-if="animals.length === 0" class="empty-state">
        <EmptyState
          icon="🐾"
          :title="t('animals.noAnimalsAvailable')"
          :description="t('animals.checkBackSoon')"
        />
      </div>

      <div v-else class="animals-grid">
        <div v-for="animal in animals" :key="animal.id" class="animal-card" @click="viewAnimal(animal.id)">
          <div class="animal-image">
            <img v-if="animal.photos && animal.photos[0]" :src="animal.photos[0]" :alt="animal.name" />
            <div v-else class="animal-image-placeholder">🐾</div>
          </div>
          <div class="animal-info">
            <h3 class="animal-name">{{ animal.name }}</h3>
            <p class="animal-details">
              {{ t(`animals.${animal.species}`) }} • {{ animal.breed }} • {{ animal.age }} {{ t('animals.years') }}
            </p>
            <p class="animal-description">{{ animal.description }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { API } from '../../api'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'
import EmptyState from '../../components/base/EmptyState.vue'

const { t } = useI18n()
const router = useRouter()

const animals = ref([])
const loading = ref(false)

onMounted(() => {
  fetchAnimals()
})

async function fetchAnimals() {
  loading.value = true
  try {
    const response = await API.animals.getAvailable()
    animals.value = response.data || []
  } catch (error) {
    console.error('Failed to fetch animals:', error)
  } finally {
    loading.value = false
  }
}

function viewAnimal(id) {
  router.push(`/animals/${id}`)
}
</script>

<style scoped>
.animals-public-page {
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

.animals-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
}

.animal-card {
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  overflow: hidden;
  cursor: pointer;
  transition: all 0.2s;
}

.animal-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
}

.animal-image {
  width: 100%;
  height: 250px;
  overflow: hidden;
  background-color: var(--bg-primary);
}

.animal-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.animal-image-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 4rem;
  opacity: 0.3;
}

.animal-info {
  padding: 1.5rem;
}

.animal-name {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.animal-details {
  color: var(--text-secondary);
  margin: 0 0 0.75rem 0;
  font-size: 0.875rem;
}

.animal-description {
  color: var(--text-primary);
  line-height: 1.6;
  margin: 0;
  display: -webkit-box;
  -webkit-line-clamp: 3;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

@media (max-width: 768px) {
  .page-title {
    font-size: 2rem;
  }

  .animals-grid {
    grid-template-columns: 1fr;
    gap: 1.5rem;
  }
}
</style>
