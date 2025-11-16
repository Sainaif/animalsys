<template>
  <div class="animal-detail">
    <LoadingSpinner v-if="loading" />

    <div
      v-else-if="animal"
      class="detail-container"
    >
      <div
        class="hero-card"
        :class="{ 'no-photo': !animalPhoto }"
      >
        <img
          v-if="animalPhoto"
          :src="animalPhoto"
          :alt="animalName"
          class="hero-image"
        >
        <div
          v-else
          class="hero-placeholder"
        >
          <i class="pi pi-heart" />
        </div>
        <div class="hero-overlay">
          <div class="hero-header">
            <Button
              icon="pi pi-arrow-left"
              class="p-button-rounded p-button-text"
              @click="router.back()"
            />
            <div class="hero-actions">
              <Button
                icon="pi pi-pencil"
                class="p-button-rounded p-button-text"
                :aria-label="$t('animal.editAnimal')"
                @click="router.push(`/staff/animals/${animal.id}/edit`)"
              />
              <Button
                icon="pi pi-trash"
                class="p-button-rounded p-button-text p-button-danger"
                :aria-label="$t('animal.deleteAnimal')"
                @click="confirmDelete"
              />
            </div>
          </div>
          <div class="hero-content">
            <Badge :variant="getStatusVariant(animal.status)">
              {{ statusLabel }}
            </Badge>
            <h1>{{ animalName }}</h1>
            <p class="hero-meta">
              {{ speciesLabel }} • {{ formatAge(animal) }} • {{ genderLabel }}
            </p>
            <div class="hero-tags">
              <span class="tag">{{ locationLabel }}</span>
              <span class="tag">{{ formatDate(animal.intake_date) }}</span>
              <span
                v-if="animal.breed"
                class="tag"
              >{{ animal.breed }}</span>
            </div>
          </div>
        </div>
      </div>

      <div class="summary-grid">
        <Card class="summary-card">
          <template #title>
            {{ $t('animal.medicalInfo') }}
          </template>
          <template #content>
            <div class="summary-content">
              <div class="summary-item">
                <span>{{ $t('animal.vaccinated') }}</span>
                <Badge :variant="animal.vaccinated ? 'success' : 'neutral'">
                  {{ animal.vaccinated ? $t('common.yes') : $t('common.no') }}
                </Badge>
              </div>
              <div class="summary-item">
                <span>{{ $t('animal.spayedNeutered') }}</span>
                <Badge :variant="animal.spayed_neutered ? 'success' : 'neutral'">
                  {{ animal.spayed_neutered ? $t('common.yes') : $t('common.no') }}
                </Badge>
              </div>
              <div class="summary-item">
                <span>{{ $t('animal.weight') }}</span>
                <strong>{{ animal.weight ? `${animal.weight} kg` : '—' }}</strong>
              </div>
              <div class="summary-item">
                <span>{{ $t('animal.microchipId') }}</span>
                <strong>{{ animal.microchip_id || '—' }}</strong>
              </div>
            </div>
          </template>
        </Card>

        <Card class="summary-card">
          <template #title>
            {{ $t('animal.behaviorInfo') }}
          </template>
          <template #content>
            <div class="summary-content">
              <div class="summary-item">
                <span>{{ $t('animal.goodWithKids') }}</span>
                <Badge :variant="animal.good_with_kids ? 'success' : 'neutral'">
                  {{ animal.good_with_kids ? $t('common.yes') : $t('common.no') }}
                </Badge>
              </div>
              <div class="summary-item">
                <span>{{ $t('animal.goodWithDogs') }}</span>
                <Badge :variant="animal.good_with_dogs ? 'success' : 'neutral'">
                  {{ animal.good_with_dogs ? $t('common.yes') : $t('common.no') }}
                </Badge>
              </div>
              <div class="summary-item">
                <span>{{ $t('animal.goodWithCats') }}</span>
                <Badge :variant="animal.good_with_cats ? 'success' : 'neutral'">
                  {{ animal.good_with_cats ? $t('common.yes') : $t('common.no') }}
                </Badge>
              </div>
              <div class="summary-item">
                <span>{{ $t('animal.houseTrained') }}</span>
                <Badge :variant="animal.house_trained ? 'success' : 'neutral'">
                  {{ animal.house_trained ? $t('common.yes') : $t('common.no') }}
                </Badge>
              </div>
            </div>
          </template>
        </Card>
      </div>

      <TabView>
        <TabPanel :header="$t('animal.basicInfo')">
          <Card class="info-card">
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('animal.name') }}</label>
                  <p>{{ animalName }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.species') }}</label>
                  <p>{{ speciesLabel }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.breed') }}</label>
                  <p>{{ animal.breed || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.gender') }}</label>
                  <p>{{ $t(`animal.${animal.sex}`) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.age') }}</label>
                  <p>{{ formatAge(animal) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.status') }}</label>
                  <Badge :variant="getStatusVariant(animal.status)">
                    {{ $t(`animal.${animal.status}`) }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.intakeDate') }}</label>
                  <p>{{ formatDate(animal.intake_date) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.color') }}</label>
                  <p>{{ animal.color || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.weight') }}</label>
                  <p>{{ animal.weight ? `${animal.weight} kg` : 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.microchipId') }}</label>
                  <p>{{ animal.microchip_id || 'N/A' }}</p>
                </div>
              </div>
              <div class="info-item full-width">
                <label>{{ $t('animal.description') }}</label>
                <p>{{ description || '—' }}</p>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel :header="$t('animal.medicalInfo')">
          <Card class="info-card">
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('animal.spayedNeutered') }}</label>
                  <Badge :variant="animal.spayed_neutered ? 'success' : 'neutral'">
                    {{ animal.spayed_neutered ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.vaccinated') }}</label>
                  <Badge :variant="animal.vaccinated ? 'success' : 'neutral'">
                    {{ animal.vaccinated ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
              </div>
              <div class="info-item full-width">
                <label>{{ $t('animal.medicalHistory') }}</label>
                <p>{{ animal.medical_history || 'N/A' }}</p>
              </div>
              <div class="info-item full-width">
                <label>{{ $t('animal.specialNeeds') }}</label>
                <p>{{ animal.special_needs || 'N/A' }}</p>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel :header="$t('animal.behaviorInfo')">
          <Card class="info-card">
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('animal.goodWithKids') }}</label>
                  <Badge :variant="animal.good_with_kids ? 'success' : 'neutral'">
                    {{ animal.good_with_kids ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.goodWithDogs') }}</label>
                  <Badge :variant="animal.good_with_dogs ? 'success' : 'neutral'">
                    {{ animal.good_with_dogs ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.goodWithCats') }}</label>
                  <Badge :variant="animal.good_with_cats ? 'success' : 'neutral'">
                    {{ animal.good_with_cats ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>{{ $t('animal.houseTrained') }}</label>
                  <Badge :variant="animal.house_trained ? 'success' : 'neutral'">
                    {{ animal.house_trained ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
              </div>
              <div class="info-item full-width">
                <label>{{ $t('animal.temperament') }}</label>
                <div class="temperament-badges">
                  <Badge
                    v-for="temp in animal.temperament || []"
                    :key="temp"
                    variant="primary"
                  >
                    {{ temp }}
                  </Badge>
                  <span v-if="!animal.temperament?.length">—</span>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>
      </TabView>
    </div>

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { animalService } from '@/services/animalService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import { getLocalizedValue, translateValue, getAnimalImage } from '@/utils/animalHelpers'

const router = useRouter()
const route = useRoute()
const { t, locale } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const animal = ref(null)
const loading = ref(true)

const loadAnimal = async () => {
  try {
    loading.value = true
    animal.value = await animalService.getAnimal(route.params.id)
  } catch (error) {
    console.error('Error loading animal:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load animal',
      life: 3000
    })
    router.push('/staff/animals')
  } finally {
    loading.value = false
  }
}

const formatAge = (animal) => {
  if (animal.age_years) return `${animal.age_years} ${t('animal.ageYears')}`
  if (animal.age_months) return `${animal.age_months} ${t('animal.ageMonths')}`
  return 'N/A'
}

const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString()
}

const getStatusVariant = (status) => {
  const variants = {
    available: 'success',
    adopted: 'info',
    under_treatment: 'warning',
    fostered: 'primary',
    transferred: 'neutral',
    deceased: 'danger'
  }
  return variants[status] || 'neutral'
}

const confirmDelete = () => {
  confirm.require({
    message: t('animal.confirmDelete'),
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await animalService.deleteAnimal(animal.value.id)
        toast.add({
          severity: 'success',
          summary: 'Success',
          detail: t('animal.animalDeleted'),
          life: 3000
        })
        router.push('/staff/animals')
      } catch (error) {
        console.error('Error deleting animal:', error)
        toast.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to delete animal',
          life: 3000
        })
      }
    }
  })
}

onMounted(() => {
  loadAnimal()
})

const animalName = computed(() => {
  if (!animal.value) return ''
  const source = animal.value.name || animal.value.Name
  if (typeof source === 'string') return source
  return getLocalizedValue(source, locale.value) || source?.en || source?.pl || ''
})

const speciesLabel = computed(() => {
  if (!animal.value?.species) return '—'
  const translated = translateValue(animal.value.species, t, 'animal.speciesNames')
  return translated || animal.value.species
})

const genderLabel = computed(() => {
  if (!animal.value?.sex) return t('animal.gender')
  const key = `animal.${animal.value.sex}`
  const translation = t(key)
  return translation !== key ? translation : animal.value.sex
})

const statusLabel = computed(() => {
  if (!animal.value?.status) return t('animal.status')
  const key = `animal.${animal.value.status}`
  const translation = t(key)
  return translation !== key ? translation : animal.value.status
})

const description = computed(() => {
  const source = animal.value?.description
  if (!source) return ''
  if (typeof source === 'string') return source
  return getLocalizedValue(source, locale.value) || source?.en || source?.pl || ''
})

const animalPhoto = computed(() => {
  const photo = getAnimalImage(animal.value)
  return photo?.includes('placehold') ? '' : photo
})

const locationLabel = computed(() => {
  return animal.value?.shelter?.location || '—'
})
</script>

<style scoped>
.animal-detail {
  max-width: 1200px;
  margin: 0 auto;
}

.hero-card {
  position: relative;
  border-radius: 24px;
  overflow: hidden;
  margin-bottom: 2rem;
  background: var(--card-bg);
  box-shadow: 0 30px 60px rgba(15, 23, 42, 0.18);
  border: 1px solid var(--border-color);
}

.hero-image {
  width: 100%;
  height: 320px;
  object-fit: cover;
  display: block;
  filter: brightness(0.85);
}

.hero-card.no-photo {
  background: linear-gradient(135deg, rgba(79, 70, 229, 0.6), rgba(16, 185, 129, 0.6));
}

.hero-placeholder {
  width: 100%;
  height: 320px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 4rem;
  color: rgba(255, 255, 255, 0.8);
}

.hero-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 1.5rem;
  background: linear-gradient(180deg, rgba(15, 23, 42, 0.1) 0%, rgba(15, 23, 42, 0.65) 100%);
  color: #fff;
}

:global([data-theme='dark']) .hero-overlay {
  background: linear-gradient(180deg, rgba(2, 6, 23, 0.4) 0%, rgba(2, 6, 23, 0.85) 100%);
}

.hero-header {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  justify-content: space-between;
}

.hero-actions {
  display: flex;
  gap: 0.5rem;
}

.hero-content {
  margin-top: auto;
}

.hero-content h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0.25rem 0;
  color: #fff;
}

.hero-meta {
  opacity: 0.9;
  margin: 0;
  color: rgba(255, 255, 255, 0.9);
}

.hero-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 0.75rem;
}

.tag {
  background: rgba(255, 255, 255, 0.15);
  border-radius: 999px;
  padding: 0.35rem 0.85rem;
  font-size: 0.85rem;
  color: #fff;
}

.summary-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(260px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.summary-card :deep(.p-card-content) {
  padding: 1rem 0;
}

.summary-card {
  border-radius: 16px;
  border: 1px solid var(--border-color);
  background: var(--card-bg);
  color: var(--text-color);
  box-shadow: 0 12px 26px rgba(15, 23, 42, 0.15);
}

.summary-content {
  display: grid;
  grid-template-columns: repeat(2, minmax(0, 1fr));
  gap: 1rem;
}

.summary-item {
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}

.summary-item span {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.summary-item strong {
  color: var(--text-color);
}

.detail-container :deep(.p-tabview) {
  border-radius: 16px;
  background: var(--card-bg);
  box-shadow: 0 20px 45px rgba(15, 23, 42, 0.12);
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  color: var(--text-color);
}

.detail-container :deep(.p-tabview-nav) {
  background: transparent;
  border: none;
  margin-bottom: 1rem;
}

.detail-container :deep(.p-tabview-nav li) {
  margin-right: 0.5rem;
}

.detail-container :deep(.p-tabview-nav-link) {
  border-radius: 999px;
  padding: 0.5rem 1.25rem;
  color: var(--text-muted);
}

.detail-container :deep(.p-tabview-nav li.p-highlight .p-tabview-nav-link) {
  background: var(--primary-color);
  color: #fff;
}

.detail-container :deep(.p-tabview-panels) {
  background: transparent;
  border: none;
  padding: 0;
}

.info-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 16px;
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.08);
}

.info-card :deep(.p-card-content) {
  padding: 0;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item label {
  font-weight: 600;
  color: var(--text-muted);
  font-size: 0.875rem;
  text-transform: uppercase;
}

:global([data-theme='dark']) .info-item label {
  color: var(--text-muted);
}

.info-item p {
  color: var(--text-color);
  font-size: 1rem;
  margin: 0;
}

:global([data-theme='dark']) .info-item p {
  color: var(--text-color);
}

.full-width {
  grid-column: 1 / -1;
}

.temperament-badges {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

@media (max-width: 768px) {
  .hero-overlay {
    padding: 1rem;
  }

  .hero-content h1 {
    font-size: 1.75rem;
  }

  .summary-content {
    grid-template-columns: 1fr;
  }
}
</style>
