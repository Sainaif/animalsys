<template>
  <div class="animal-gallery">
    <section class="gallery-hero">
      <div class="container">
        <h1>{{ $t('home.adoptions.title') }}</h1>
        <p>{{ $t('home.adoptions.subtitle') }}</p>
      </div>
    </section>

    <section class="gallery-filters">
      <div class="container filters-container">
        <Dropdown
          v-model="filters.species"
          :options="speciesOptions"
          optionLabel="label"
          optionValue="value"
          class="filter-dropdown"
          @change="applyFilters"
        />
        <span class="p-input-icon-left search-input">
          <i class="pi pi-search" />
          <InputText
            v-model="filters.search"
            :placeholder="$t('home.adoptions.searchPlaceholder')"
            @keyup.enter="applyFilters"
          />
        </span>
        <Button
          :label="$t('common.search')"
          icon="pi pi-search"
          class="p-button-outlined"
          @click="applyFilters"
        />
      </div>
    </section>

    <section class="gallery-grid">
      <div class="container">
        <div v-if="featuredAnimal" class="featured-card">
          <div class="featured-header">
            <h3>{{ $t('home.adoptions.viewProfile') }}</h3>
            <Button
              icon="pi pi-times"
              class="p-button-rounded p-button-text"
              @click="clearHighlight"
            />
          </div>
          <div class="animal-card featured">
            <div class="animal-image">
              <img :src="getAnimalImageSrc(featuredAnimal)" :alt="getAnimalName(featuredAnimal)" />
              <span class="animal-badge">{{ getStatusLabel(featuredAnimal) }}</span>
            </div>
            <div class="animal-info">
              <h3>{{ getAnimalName(featuredAnimal) }}</h3>
              <p class="animal-meta">
                <span><i class="pi pi-tag" /> {{ getAnimalSpeciesLabel(featuredAnimal) }}</span>
                <span v-if="getAnimalColorLabel(featuredAnimal)"><i class="pi pi-palette" /> {{ getAnimalColorLabel(featuredAnimal) }}</span>
                <span><i class="pi pi-calendar" /> {{ formatAge(featuredAnimal) }}</span>
              </p>
              <p class="animal-description">
                {{ truncateText(getAnimalDescription(featuredAnimal), 160) }}
              </p>
            </div>
          </div>
        </div>

        <div v-if="loading && animals.length === 0" class="loading-state">
          <ProgressSpinner />
        </div>

        <div v-else-if="animals.length === 0" class="empty-state">
          <i class="pi pi-paw" />
          <p>{{ $t('animal.noAnimalsFound') }}</p>
        </div>

        <div v-else class="animals-grid">
          <div
            v-for="animal in animals"
            :key="animal.id"
            class="animal-card"
          >
            <div class="animal-image">
              <img :src="getAnimalImageSrc(animal)" :alt="getAnimalName(animal)" />
              <span class="animal-badge">{{ getStatusLabel(animal) }}</span>
            </div>
            <div class="animal-info">
              <h3>{{ getAnimalName(animal) }}</h3>
              <p class="animal-meta">
                <span><i class="pi pi-tag" /> {{ getAnimalSpeciesLabel(animal) }}</span>
                <span v-if="getAnimalColorLabel(animal)"><i class="pi pi-palette" /> {{ getAnimalColorLabel(animal) }}</span>
                <span><i class="pi pi-calendar" /> {{ formatAge(animal) }}</span>
              </p>
              <p class="animal-description">
                {{ truncateText(getAnimalDescription(animal), 140) }}
              </p>
            </div>
          </div>
        </div>

        <div v-if="hasMore" class="load-more">
          <Button
            :label="$t('home.adoptions.loadMore')"
            icon="pi pi-arrow-down"
            :loading="loading"
            @click="loadMore"
          />
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import ProgressSpinner from 'primevue/progressspinner'
import api from '@/services/api'
import { getLocalizedValue, getAnimalImage, translateValue } from '@/utils/animalHelpers'

const { t, locale } = useI18n()
const route = useRoute()
const router = useRouter()

const animals = ref([])
const loading = ref(false)
const featuredAnimal = ref(null)
const filters = reactive({
  species: '',
  search: ''
})
const pagination = reactive({
  limit: 12,
  offset: 0,
  total: 0
})

const baseSpecies = ['Dog', 'Cat', 'Rabbit', 'Bird']
const speciesOptions = ref([])

const updateSpeciesOptions = () => {
  speciesOptions.value = [
    { label: t('home.adoptions.allSpecies'), value: '' },
    ...baseSpecies.map((name) => ({
      label: translateValue(name, t, 'animal.speciesNames'),
      value: name
    }))
  ]
}

updateSpeciesOptions()
watch(() => locale.value, updateSpeciesOptions)

const hasMore = computed(() => animals.value.length < pagination.total)

const loadAnimals = async (reset = false) => {
  if (loading.value) return

  loading.value = true
  if (reset) {
    animals.value = []
    pagination.offset = 0
  } else {
    pagination.offset = animals.value.length
  }

  try {
    const params = {
      limit: pagination.limit,
      offset: pagination.offset,
      species: filters.species || undefined,
      search: filters.search || undefined,
      available_only: true,
      sort_by: 'created_at',
      sort_order: 'desc'
    }
    const response = await api.get('/public/animals', { params })
    const list = response.data?.animals || []
    pagination.total = response.data?.total || list.length
    animals.value = reset ? list : [...animals.value, ...list]
    pagination.offset = animals.value.length
  } catch (error) {
    console.error('Error loading animals list:', error)
    if (reset) {
      animals.value = []
      pagination.total = 0
    }
  } finally {
    loading.value = false
  }
}

const applyFilters = () => {
  loadAnimals(true)
}

const loadMore = () => {
  if (hasMore.value) {
    loadAnimals()
  }
}

const loadHighlightedAnimal = async (id) => {
  try {
    const response = await api.get(`/public/animals/${id}`)
    featuredAnimal.value = response.data
  } catch (error) {
    console.warn('Unable to load highlighted animal', error)
    featuredAnimal.value = null
  }
}

const clearHighlight = () => {
  router.replace({ name: 'public-animals' })
  featuredAnimal.value = null
}

const getAnimalName = (animal) => {
  const localized = getLocalizedValue(animal?.name, locale.value)
  return localized || t('animal.unknown')
}

const getAnimalDescription = (animal) => {
  const localized = getLocalizedValue(animal?.description, locale.value)
  return localized || ''
}

const getAnimalImageSrc = (animal) => getAnimalImage(animal)

const getAnimalSpeciesLabel = (animal) => {
  if (!animal) return ''
  if (animal.breed) return animal.breed
  const translated = translateValue(animal.species || '', t, 'animal.speciesNames')
  return translated || animal.species || ''
}

const getAnimalColorLabel = (animal) => {
  if (!animal?.color) return ''
  return translateValue(animal.color, t, 'animal.colorNames')
}

const getStatusLabel = (animal) => {
  const status = (animal?.status || '').toLowerCase()
  switch (status) {
    case 'available':
      return t('animal.available')
    case 'adopted':
      return t('animal.adopted')
    case 'fostered':
      return t('animal.fostered')
    default:
      return t('animal.status')
  }
}

const formatAge = (animal) => {
  if (animal.date_of_birth) {
    const dob = new Date(animal.date_of_birth)
    if (!Number.isNaN(dob.getTime())) {
      const diff = Date.now() - dob.getTime()
      const years = Math.floor(diff / (1000 * 60 * 60 * 24 * 365))
      if (years > 0) {
        return `${years} ${t('home.adoptions.years')}`
      }
      const months = Math.max(1, Math.floor(diff / (1000 * 60 * 60 * 24 * 30)))
      return `${months} ${t('home.adoptions.months')}`
    }
  }
  return t('animal.unknown')
}

const truncateText = (text, length) => {
  if (!text) return ''
  if (text.length <= length) return text
  return `${text.substring(0, length)}...`
}

onMounted(() => {
  if (route.query.animal) {
    loadHighlightedAnimal(route.query.animal)
  }
  loadAnimals(true)
})

watch(
  () => route.query.animal,
  (id) => {
    if (id) {
      loadHighlightedAnimal(id)
    } else {
      featuredAnimal.value = null
    }
  }
)
</script>

<style scoped>
.animal-gallery {
  background-color: var(--surface-ground);
  min-height: 100vh;
  color: var(--text-color);
}

.gallery-hero {
  background: linear-gradient(135deg, rgba(249, 197, 209, 0.9) 0%, rgba(251, 194, 235, 0.85) 100%);
  padding: 3rem 0;
  text-align: center;
  color: #4a154b;
}

.gallery-hero h1 {
  font-size: 2.5rem;
  margin-bottom: 0.5rem;
}

.gallery-filters {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.08);
  padding: 1.5rem 0;
}

.filters-container {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  align-items: center;
}

.filter-dropdown {
  min-width: 220px;
}

.search-input {
  flex: 1;
  min-width: 220px;
}

.gallery-grid {
  padding: 2rem 0 4rem;
}

.animals-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1.5rem;
}

.featured-card {
  background: var(--card-bg);
  border-radius: 20px;
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  box-shadow: 0 20px 45px rgba(15, 23, 42, 0.12);
  margin-bottom: 2rem;
}

.featured-card .animal-card {
  box-shadow: none;
  margin: 0;
}

.featured-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.animal-card {
  background: var(--card-bg);
  border-radius: 16px;
  border: 1px solid var(--border-color);
  box-shadow: 0 18px 45px rgba(79, 70, 229, 0.12);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.animal-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 14px 35px rgba(0, 0, 0, 0.12);
}

.animal-image {
  position: relative;
  padding-top: 65%;
  overflow: hidden;
}

.animal-image img {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.animal-badge {
  position: absolute;
  top: 15px;
  left: 15px;
  background-color: rgba(255, 255, 255, 0.9);
  color: #e74c3c;
  padding: 0.35rem 0.9rem;
  border-radius: 999px;
  font-weight: 600;
  font-size: 0.9rem;
}

.animal-info {
  padding: 1.25rem;
  flex: 1;
}

.animal-info h3 {
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.animal-meta {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
  color: var(--text-muted);
  font-size: 0.95rem;
  margin-bottom: 0.75rem;
}

.animal-meta span {
  display: flex;
  align-items: center;
  gap: 0.35rem;
}

.animal-description {
  color: var(--text-muted);
  line-height: 1.5;
  margin: 0;
}

.load-more {
  display: flex;
  justify-content: center;
  margin-top: 2rem;
}

.loading-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 3rem 0;
  color: var(--text-muted);
  gap: 1rem;
}

.empty-state i {
  font-size: 3rem;
}

:global([data-theme='dark'] .gallery-hero) {
  color: #f8fafc;
  background: linear-gradient(135deg, rgba(78, 54, 127, 0.9) 0%, rgba(112, 63, 134, 0.85) 100%);
}

@media (max-width: 768px) {
  .filters-container {
    flex-direction: column;
    align-items: stretch;
  }

  .filter-dropdown,
  .search-input {
    width: 100%;
  }
}
</style>
