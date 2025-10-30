<template>
  <div class="animal-view-page">
    <div class="page-header">
      <BaseButton variant="outline" size="small" @click="goBack">
        {{ t('common.back') }}
      </BaseButton>
      <div class="page-actions">
        <BaseButton v-if="canEdit" variant="primary" @click="editAnimal">
          {{ t('common.edit') }}
        </BaseButton>
        <BaseButton v-if="canDelete" variant="danger" @click="confirmDelete">
          {{ t('common.delete') }}
        </BaseButton>
      </div>
    </div>

    <LoadingSpinner v-if="loading" />

    <div v-else-if="animal" class="animal-content">
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
            <span :class="['status-badge', `status-badge--${animal.status}`]">
              {{ t(`animals.${animal.status}`) }}
            </span>
            <span class="info-badge">{{ t(`animals.${animal.species}`) }}</span>
            <span class="info-badge">{{ animal.gender === 'male' ? t('animals.male') : t('animals.female') }}</span>
          </div>

          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">{{ t('animals.breed') }}:</span>
              <span class="info-value">{{ animal.breed }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ t('animals.age') }}:</span>
              <span class="info-value">{{ animal.age }} {{ t('animals.years') }}</span>
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
            <div class="info-item">
              <span class="info-label">{{ t('animals.admissionDate') }}:</span>
              <span class="info-value">{{ formatDate(animal.admission_date) }}</span>
            </div>
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
          <div class="info-grid">
            <div class="info-item">
              <span class="info-label">{{ t('animals.sterilized') }}:</span>
              <span class="info-value">{{ animal.sterilized ? t('common.yes') : t('common.no') }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ t('animals.vaccinated') }}:</span>
              <span class="info-value">{{ animal.vaccinated ? t('common.yes') : t('common.no') }}</span>
            </div>
            <div class="info-item">
              <span class="info-label">{{ t('animals.chipped') }}:</span>
              <span class="info-value">{{ animal.chipped ? t('common.yes') : t('common.no') }}</span>
            </div>
          </div>
          <div v-if="animal.chip_number" class="chip-info">
            <span class="info-label">{{ t('animals.chipNumber') }}:</span>
            <span class="info-value">{{ animal.chip_number }}</span>
          </div>
          <div v-if="animal.medical_notes" class="medical-notes">
            <h4>{{ t('animals.medicalNotes') }}</h4>
            <p>{{ animal.medical_notes }}</p>
          </div>
        </BaseCard>

        <BaseCard v-if="animal.behavioral_notes">
          <template #header>{{ t('animals.behavioralNotes') }}</template>
          <p>{{ animal.behavioral_notes }}</p>
        </BaseCard>

        <BaseCard v-if="animal.special_needs">
          <template #header>{{ t('animals.specialNeeds') }}</template>
          <p>{{ animal.special_needs }}</p>
        </BaseCard>

        <BaseCard v-if="animal.location">
          <template #header>{{ t('animals.location') }}</template>
          <p>{{ animal.location }}</p>
        </BaseCard>
      </div>
    </div>

    <BaseModal v-model:show="showDeleteModal" size="small">
      <template #header>{{ t('animals.deleteConfirmTitle') }}</template>
      <p>{{ t('animals.deleteConfirmMessage', { name: animal?.name }) }}</p>
      <template #footer>
        <BaseButton variant="outline" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteAnimal" :loading="deleting">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notification'
import { API } from '../../api'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const animal = ref(null)
const loading = ref(false)
const showDeleteModal = ref(false)
const deleting = ref(false)

const canEdit = computed(() => authStore.hasRole('employee'))
const canDelete = computed(() => authStore.hasRole('admin'))

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
    router.push({ name: 'animals-list' })
  } finally {
    loading.value = false
  }
}

function formatDate(dateString) {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString()
}

function goBack() {
  router.back()
}

function editAnimal() {
  router.push({ name: 'animal-edit', params: { id: animal.value.id } })
}

function confirmDelete() {
  showDeleteModal.value = true
}

async function deleteAnimal() {
  deleting.value = true
  try {
    await API.animals.delete(animal.value.id)
    notificationStore.success(t('animals.deleteSuccess'))
    router.push({ name: 'animals-list' })
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    deleting.value = false
    showDeleteModal.value = false
  }
}
</script>

<style scoped>
.animal-view-page {
  padding: 2rem;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-actions {
  display: flex;
  gap: 1rem;
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

.status-badge {
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
}

.status-badge--available {
  background-color: #d4edda;
  color: #155724;
}

.status-badge--adopted {
  background-color: #d1ecf1;
  color: #0c5460;
}

.status-badge--reserved {
  background-color: #fff3cd;
  color: #856404;
}

.status-badge--medical_care {
  background-color: #f8d7da;
  color: #721c24;
}

.status-badge--quarantine {
  background-color: #e2e3e5;
  color: #383d41;
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

.chip-info,
.medical-notes {
  margin-top: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
}

.medical-notes h4 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.medical-notes p {
  line-height: 1.6;
  color: var(--text-primary);
  margin: 0;
}

@media (max-width: 768px) {
  .animal-view-page {
    padding: 1rem;
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
