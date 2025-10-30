<template>
  <div class="volunteer-view-page">
    <div class="page-header">
      <BaseButton variant="outline" size="small" @click="goBack">
        {{ t('common.back') }}
      </BaseButton>
      <div class="page-actions">
        <BaseButton
          v-if="canEdit"
          variant="primary"
          @click="editVolunteer"
        >
          {{ t('common.edit') }}
        </BaseButton>
        <BaseButton
          v-if="canDelete"
          variant="danger"
          @click="confirmDelete"
        >
          {{ t('common.delete') }}
        </BaseButton>
      </div>
    </div>

    <LoadingSpinner v-if="loading" />

    <div v-else-if="volunteer" class="volunteer-content">
      <!-- Header with Status -->
      <BaseCard>
        <div class="volunteer-header">
          <div>
            <h1 class="volunteer-name">{{ volunteer.first_name }} {{ volunteer.last_name }}</h1>
            <span :class="['status-badge-large', `status-badge--${volunteer.status}`]">
              {{ t(`volunteers.status${capitalize(volunteer.status)}`) }}
            </span>
          </div>
          <div class="volunteer-meta">
            <div class="meta-item">
              <span class="meta-label">{{ t('volunteers.registrationDate') }}:</span>
              <span class="meta-value">{{ formatDate(volunteer.registration_date) }}</span>
            </div>
          </div>
        </div>
      </BaseCard>

      <!-- Contact Info -->
      <BaseCard>
        <template #header>{{ t('volunteers.contactInfo') }}</template>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('common.email') }}:</span>
            <a :href="`mailto:${volunteer.email}`" class="info-link">{{ volunteer.email }}</a>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('common.phone') }}:</span>
            <a :href="`tel:${volunteer.phone}`" class="info-link">{{ volunteer.phone }}</a>
          </div>
          <div v-if="volunteer.address" class="info-item info-item-full">
            <span class="info-label">{{ t('common.address') }}:</span>
            <span class="info-value">{{ volunteer.address }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Hours Statistics -->
      <BaseCard>
        <template #header>{{ t('volunteers.volunteerHours') }}</template>
        <div class="hours-stats">
          <div class="stat-card">
            <div class="stat-value">{{ volunteer.total_hours || 0 }}</div>
            <div class="stat-label">{{ t('volunteers.totalHours') }}</div>
          </div>
          <div class="stat-card">
            <div class="stat-value">{{ volunteer.hours_this_month || 0 }}</div>
            <div class="stat-label">{{ t('volunteers.thisMonth') }}</div>
          </div>
        </div>
        <BaseButton
          variant="primary"
          size="small"
          @click="showLogHoursModal"
          class="log-hours-button"
        >
          {{ t('volunteers.logHours') }}
        </BaseButton>
      </BaseCard>

      <!-- Skills & Availability -->
      <div class="two-column-grid">
        <BaseCard v-if="volunteer.skills">
          <template #header>{{ t('volunteers.skills') }}</template>
          <p class="text-content">{{ volunteer.skills }}</p>
        </BaseCard>

        <BaseCard v-if="volunteer.availability">
          <template #header>{{ t('volunteers.availability') }}</template>
          <p class="text-content">{{ volunteer.availability }}</p>
        </BaseCard>
      </div>

      <!-- Training Records -->
      <BaseCard v-if="volunteer.trainings && volunteer.trainings.length > 0">
        <template #header>{{ t('volunteers.training') }}</template>
        <div class="trainings-list">
          <div
            v-for="training in volunteer.trainings"
            :key="training.id"
            class="training-item"
          >
            <div class="training-info">
              <span class="training-name">{{ training.training_type }}</span>
              <span class="training-date">{{ formatDate(training.completion_date) }}</span>
            </div>
            <span v-if="training.certificate_issued" class="certificate-badge">
              ✓ {{ t('volunteers.certified') }}
            </span>
          </div>
        </div>
        <BaseButton
          variant="outline"
          size="small"
          @click="showAddTrainingModal"
          class="add-training-button"
        >
          {{ t('volunteers.addTraining') }}
        </BaseButton>
      </BaseCard>

      <!-- Notes -->
      <BaseCard v-if="volunteer.notes">
        <template #header>{{ t('common.notes') }}</template>
        <p class="text-content">{{ volunteer.notes }}</p>
      </BaseCard>
    </div>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-model="deleteModal.show"
      :title="t('volunteers.deleteVolunteer')"
      size="small"
    >
      <p>{{ t('volunteers.deleteConfirm') }}</p>
      <p v-if="volunteer" class="confirm-name">
        <strong>{{ volunteer.first_name }} {{ volunteer.last_name }}</strong>
      </p>

      <template #footer>
        <BaseButton variant="outline" @click="deleteModal.show = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton
          variant="danger"
          :loading="deleteModal.loading"
          @click="deleteVolunteer"
        >
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Log Hours Modal -->
    <BaseModal
      v-model="logHoursModal.show"
      :title="t('volunteers.logHours')"
      size="medium"
    >
      <FormGroup :label="t('volunteers.hoursWorked')" required>
        <input
          v-model.number="logHoursModal.hours"
          type="number"
          min="0"
          step="0.5"
          :placeholder="t('volunteers.hoursWorked')"
        />
      </FormGroup>

      <FormGroup :label="t('common.date')" required>
        <input
          v-model="logHoursModal.date"
          type="date"
        />
      </FormGroup>

      <FormGroup :label="t('common.notes')">
        <textarea
          v-model="logHoursModal.notes"
          :placeholder="t('common.notes')"
          rows="3"
        ></textarea>
      </FormGroup>

      <template #footer>
        <BaseButton variant="outline" @click="logHoursModal.show = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton
          variant="primary"
          :loading="logHoursModal.loading"
          @click="logHours"
        >
          {{ t('volunteers.logHours') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notification'
import { API } from '../../api'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import FormGroup from '../../components/base/FormGroup.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const volunteer = ref(null)
const loading = ref(false)

const canEdit = computed(() => authStore.hasRole('employee'))
const canDelete = computed(() => authStore.hasRole('admin'))

const deleteModal = reactive({
  show: false,
  loading: false
})

const logHoursModal = reactive({
  show: false,
  hours: 0,
  date: new Date().toISOString().split('T')[0],
  notes: '',
  loading: false
})

onMounted(() => {
  fetchVolunteer()
})

async function fetchVolunteer() {
  loading.value = true
  try {
    const response = await API.volunteers.getById(route.params.id)
    volunteer.value = response.data
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
    router.push({ name: 'volunteers-list' })
  } finally {
    loading.value = false
  }
}

function formatDate(dateString) {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString()
}

function capitalize(str) {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1).replace('_', '')
}

function goBack() {
  router.back()
}

function editVolunteer() {
  router.push({ name: 'volunteer-edit', params: { id: volunteer.value.id } })
}

function confirmDelete() {
  deleteModal.show = true
}

async function deleteVolunteer() {
  deleteModal.loading = true
  try {
    await API.volunteers.delete(volunteer.value.id)
    notificationStore.success(t('volunteers.deleteSuccess'))
    router.push({ name: 'volunteers-list' })
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    deleteModal.loading = false
    deleteModal.show = false
  }
}

function showLogHoursModal() {
  logHoursModal.hours = 0
  logHoursModal.date = new Date().toISOString().split('T')[0]
  logHoursModal.notes = ''
  logHoursModal.show = true
}

async function logHours() {
  if (!logHoursModal.hours || logHoursModal.hours <= 0) {
    notificationStore.error(t('common.error'), t('volunteers.hoursRequired'))
    return
  }

  logHoursModal.loading = true
  try {
    await API.volunteers.logHours(volunteer.value.id, {
      hours: logHoursModal.hours,
      date: logHoursModal.date,
      notes: logHoursModal.notes
    })
    notificationStore.success(t('volunteers.hoursLoggedSuccess'))
    logHoursModal.show = false
    fetchVolunteer() // Reload to get updated stats
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    logHoursModal.loading = false
  }
}

function showAddTrainingModal() {
  // This would open a modal to add training records
  notificationStore.info(t('volunteers.comingSoon'), t('volunteers.trainingModalComingSoon'))
}
</script>

<style scoped>
.volunteer-view-page {
  padding: 2rem;
  max-width: 1200px;
  margin: 0 auto;
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

.volunteer-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.volunteer-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
}

.volunteer-name {
  font-size: 2rem;
  font-weight: bold;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.status-badge-large {
  display: inline-block;
  padding: 0.5rem 1rem;
  border-radius: 0.5rem;
  font-size: 0.875rem;
  font-weight: 600;
}

.status-badge--active {
  background-color: #d1fae5;
  color: #065f46;
}

.status-badge--inactive {
  background-color: #f3f4f6;
  color: #374151;
}

.status-badge--on_leave {
  background-color: #fef3c7;
  color: #92400e;
}

.volunteer-meta {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.5rem;
}

.meta-item {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 0.25rem;
}

.meta-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.meta-value {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.info-item-full {
  grid-column: 1 / -1;
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

.info-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 600;
}

.info-link:hover {
  text-decoration: underline;
}

.hours-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.stat-card {
  padding: 1.5rem;
  background-color: var(--bg-secondary);
  border-radius: 0.75rem;
  text-align: center;
}

.stat-value {
  font-size: 2.5rem;
  font-weight: bold;
  color: var(--primary-color);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: 0.5rem;
}

.log-hours-button,
.add-training-button {
  margin-top: 1rem;
}

.two-column-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
}

.text-content {
  line-height: 1.6;
  color: var(--text-primary);
  margin: 0;
  white-space: pre-wrap;
}

.trainings-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
  margin-bottom: 1rem;
}

.training-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 0.75rem;
  background-color: var(--bg-secondary);
  border-radius: 0.5rem;
}

.training-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.training-name {
  font-weight: 600;
  color: var(--text-primary);
}

.training-date {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.certificate-badge {
  padding: 0.25rem 0.75rem;
  background-color: #d1fae5;
  color: #065f46;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
}

.confirm-name {
  margin: 1rem 0;
  text-align: center;
}

input,
textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
  font-family: inherit;
}

input:focus,
textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

textarea {
  resize: vertical;
}

@media (max-width: 768px) {
  .volunteer-view-page {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .page-actions {
    width: 100%;
    flex-direction: column;
  }

  .volunteer-header {
    flex-direction: column;
    gap: 1rem;
  }

  .volunteer-meta {
    align-items: flex-start;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .two-column-grid {
    grid-template-columns: 1fr;
  }
}
</style>
