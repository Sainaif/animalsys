<template>
  <div class="veterinary-view-page">
    <LoadingSpinner v-if="loading" />
    <div v-else>
      <div class="page-header">
        <div>
          <h1 class="page-title">{{ t('veterinary.visitDetails') }}</h1>
          <span class="badge" :class="`badge-${visit.type}`">
            {{ t(`veterinary.type${visit.type?.charAt(0).toUpperCase() + visit.type?.slice(1)}`) }}
          </span>
          <span class="badge" :class="`badge-status-${visit.status}`">
            {{ t(`veterinary.status${visit.status?.charAt(0).toUpperCase() + visit.status?.slice(1)}`) }}
          </span>
        </div>
        <div class="header-actions">
          <BaseButton
            v-if="authStore.hasRole('staff')"
            variant="secondary"
            @click="editVisit"
          >
            {{ t('common.edit') }}
          </BaseButton>
          <BaseButton
            v-if="authStore.hasRole('admin')"
            variant="danger"
            @click="showDeleteModal = true"
          >
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </div>

      <!-- Visit Information -->
      <div class="info-grid">
        <!-- Basic Info -->
        <BaseCard>
          <template #header>{{ t('veterinary.visitInfo') }}</template>
          <div class="info-item">
            <span class="info-label">{{ t('veterinary.animal') }}:</span>
            <router-link
              v-if="visit.animal"
              :to="{ name: 'AnimalView', params: { id: visit.animal_id } }"
              class="info-value animal-link"
            >
              {{ visit.animal.name }}
            </router-link>
            <span v-else class="info-value">-</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('veterinary.visitDate') }}:</span>
            <span class="info-value">{{ formatDateTime(visit.visit_date) }}</span>
          </div>
          <div v-if="visit.next_visit_date" class="info-item">
            <span class="info-label">{{ t('veterinary.nextVisit') }}:</span>
            <span class="info-value">{{ formatDate(visit.next_visit_date) }}</span>
          </div>
          <div v-if="visit.cost" class="info-item">
            <span class="info-label">{{ t('veterinary.cost') }}:</span>
            <span class="info-value">{{ formatCurrency(visit.cost) }}</span>
          </div>
        </BaseCard>

        <!-- Veterinarian Info -->
        <BaseCard>
          <template #header>{{ t('veterinary.veterinarianInfo') }}</template>
          <div v-if="visit.veterinarian_name" class="info-item">
            <span class="info-label">{{ t('veterinary.veterinarian') }}:</span>
            <span class="info-value">{{ visit.veterinarian_name }}</span>
          </div>
          <div v-if="visit.clinic_name" class="info-item">
            <span class="info-label">{{ t('veterinary.clinic') }}:</span>
            <span class="info-value">{{ visit.clinic_name }}</span>
          </div>
          <div v-if="!visit.veterinarian_name && !visit.clinic_name" class="info-item">
            <span class="info-value">{{ t('veterinary.noVeterinarianInfo') }}</span>
          </div>
        </BaseCard>
      </div>

      <!-- Medical Details -->
      <BaseCard v-if="visit.diagnosis || visit.treatment || visit.prescription">
        <template #header>{{ t('veterinary.medicalDetails') }}</template>

        <div v-if="visit.diagnosis" class="medical-section">
          <h4 class="medical-title">{{ t('veterinary.diagnosis') }}</h4>
          <p class="medical-content">{{ visit.diagnosis }}</p>
        </div>

        <div v-if="visit.treatment" class="medical-section">
          <h4 class="medical-title">{{ t('veterinary.treatment') }}</h4>
          <p class="medical-content">{{ visit.treatment }}</p>
        </div>

        <div v-if="visit.prescription" class="medical-section">
          <h4 class="medical-title">{{ t('veterinary.prescription') }}</h4>
          <p class="medical-content">{{ visit.prescription }}</p>
        </div>
      </BaseCard>

      <!-- Notes -->
      <BaseCard v-if="visit.notes">
        <template #header>{{ t('common.notes') }}</template>
        <p class="notes-content">{{ visit.notes }}</p>
      </BaseCard>
    </div>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      :title="t('veterinary.deleteVisit')"
      @close="showDeleteModal = false"
    >
      <p>{{ t('veterinary.deleteVisitConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteVisit">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notifications'
import { API } from '../../api'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const visit = ref({})
const loading = ref(false)
const showDeleteModal = ref(false)

async function fetchVisit() {
  try {
    loading.value = true
    const response = await API.veterinary.getById(route.params.id)
    visit.value = response.data
  } catch (error) {
    console.error('Failed to fetch visit:', error)
    notificationStore.error(t('veterinary.fetchError'))
    router.push({ name: 'Veterinary' })
  } finally {
    loading.value = false
  }
}

function editVisit() {
  router.push({ name: 'VeterinaryForm', params: { id: route.params.id } })
}

async function deleteVisit() {
  try {
    await API.veterinary.delete(route.params.id)
    notificationStore.success(t('veterinary.deleteSuccess'))
    router.push({ name: 'Veterinary' })
  } catch (error) {
    console.error('Failed to delete visit:', error)
    notificationStore.error(t('veterinary.deleteError'))
  }
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
}

function formatDate(date) {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('pl-PL')
}

function formatDateTime(date) {
  if (!date) return '-'
  return new Date(date).toLocaleString('pl-PL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(() => {
  fetchVisit()
})
</script>

<style scoped>
.veterinary-view-page {
  max-width: 1200px;
  padding: 2rem;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin: 0 0 0.5rem 0;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  margin-right: 0.5rem;
}

.badge-checkup {
  background: #e3f2fd;
  color: #1976d2;
}

.badge-vaccination {
  background: #e8f5e9;
  color: #388e3c;
}

.badge-treatment {
  background: #f3e5f5;
  color: #7b1fa2;
}

.badge-surgery {
  background: #ffebee;
  color: #c62828;
}

.badge-emergency {
  background: #ff5252;
  color: #ffffff;
}

.badge-status-scheduled {
  background: #fff3e0;
  color: #f57c00;
}

.badge-status-completed {
  background: #e8f5e9;
  color: #388e3c;
}

.badge-status-cancelled {
  background: #f5f5f5;
  color: #757575;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.info-item {
  display: flex;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-weight: 600;
  min-width: 180px;
  color: var(--text-secondary);
}

.info-value {
  color: var(--text-primary);
}

.animal-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.animal-link:hover {
  text-decoration: underline;
}

.medical-section {
  margin-bottom: 1.5rem;
}

.medical-section:last-child {
  margin-bottom: 0;
}

.medical-title {
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.75rem;
}

.medical-content {
  white-space: pre-wrap;
  line-height: 1.6;
  color: var(--text-primary);
  margin: 0;
}

.notes-content {
  white-space: pre-wrap;
  line-height: 1.6;
  color: var(--text-primary);
}
</style>
