<template>
  <div class="adoption-view-page">
    <div class="page-header">
      <BaseButton variant="outline" size="small" @click="goBack">
        {{ t('common.back') }}
      </BaseButton>
      <div class="page-actions">
        <BaseButton
          v-if="canManage && adoption?.status === 'under_review'"
          variant="success"
          @click="approveAdoption"
          :loading="approving"
        >
          {{ t('adoptions.approve') }}
        </BaseButton>
        <BaseButton
          v-if="canManage && adoption?.status === 'under_review'"
          variant="danger"
          @click="showRejectModal"
        >
          {{ t('adoptions.reject') }}
        </BaseButton>
        <BaseButton
          v-if="canManage && adoption?.status === 'approved'"
          variant="primary"
          @click="showCompleteModal"
        >
          {{ t('adoptions.complete') }}
        </BaseButton>
      </div>
    </div>

    <LoadingSpinner v-if="loading" />

    <div v-else-if="adoption" class="adoption-content">
      <!-- Status -->
      <BaseCard>
        <template #header>{{ t('common.status') }}</template>
        <div class="status-section">
          <span :class="['status-badge-large', `status-badge--${adoption.status}`]">
            {{ t(`adoptions.${adoption.status}`) }}
          </span>
          <div class="status-dates">
            <div class="date-item">
              <span class="date-label">{{ t('adoptions.applicationDate') }}:</span>
              <span class="date-value">{{ formatDate(adoption.application_date) }}</span>
            </div>
            <div v-if="adoption.approval_date" class="date-item">
              <span class="date-label">{{ t('adoptions.approvalDate') }}:</span>
              <span class="date-value">{{ formatDate(adoption.approval_date) }}</span>
            </div>
          </div>
        </div>
      </BaseCard>

      <!-- Animal Info -->
      <BaseCard v-if="adoption.animal_id">
        <template #header>{{ t('adoptions.animalInfo') }}</template>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('animals.name') }}:</span>
            <RouterLink
              :to="`/app/animals/${adoption.animal_id}`"
              class="info-link"
            >
              {{ adoption.animal_name }}
            </RouterLink>
          </div>
        </div>
      </BaseCard>

      <!-- Applicant Info -->
      <BaseCard>
        <template #header>{{ t('adoptions.applicantInfo') }}</template>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('adoptions.applicantName') }}:</span>
            <span class="info-value">{{ adoption.applicant_name }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('common.email') }}:</span>
            <span class="info-value">{{ adoption.applicant_email }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('common.phone') }}:</span>
            <span class="info-value">{{ adoption.applicant_phone }}</span>
          </div>
          <div class="info-item info-item-full">
            <span class="info-label">{{ t('common.address') }}:</span>
            <span class="info-value">{{ adoption.applicant_address }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Housing Info -->
      <BaseCard>
        <template #header>{{ t('adoptions.housingInfo') }}</template>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('adoptions.housingType') }}:</span>
            <span class="info-value">{{ t(`adoptions.${adoption.housing_type}`) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('adoptions.hasYard') }}:</span>
            <span class="info-value">{{ adoption.has_yard ? t('common.yes') : t('common.no') }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('adoptions.hasOtherPets') }}:</span>
            <span class="info-value">{{ adoption.has_other_pets ? t('common.yes') : t('common.no') }}</span>
          </div>
          <div v-if="adoption.other_pets_description" class="info-item info-item-full">
            <span class="info-label">{{ t('adoptions.otherPetsDescription') }}:</span>
            <span class="info-value">{{ adoption.other_pets_description }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Experience & Reason -->
      <BaseCard>
        <template #header>{{ t('adoptions.experienceInfo') }}</template>
        <div class="text-section">
          <h4>{{ t('adoptions.previousPetExperience') }}</h4>
          <p>{{ adoption.previous_pet_experience }}</p>
        </div>
        <div class="text-section">
          <h4>{{ t('adoptions.reasonForAdoption') }}</h4>
          <p>{{ adoption.reason_for_adoption }}</p>
        </div>
      </BaseCard>

      <!-- Interview Info -->
      <BaseCard v-if="adoption.interview_date">
        <template #header>{{ t('adoptions.interviewInfo') }}</template>
        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('adoptions.interviewDate') }}:</span>
            <span class="info-value">{{ formatDate(adoption.interview_date) }}</span>
          </div>
          <div v-if="adoption.interview_notes" class="info-item info-item-full">
            <span class="info-label">{{ t('common.notes') }}:</span>
            <span class="info-value">{{ adoption.interview_notes }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Rejection Reason -->
      <BaseCard v-if="adoption.status === 'rejected' && adoption.rejection_reason">
        <template #header>{{ t('adoptions.rejectionReason') }}</template>
        <p>{{ adoption.rejection_reason }}</p>
      </BaseCard>
    </div>

    <!-- Reject Modal -->
    <BaseModal
      v-model="rejectModal.show"
      :title="t('adoptions.rejectApplication')"
      size="medium"
    >
      <FormGroup :label="t('adoptions.rejectionReason')" required>
        <textarea
          v-model="rejectModal.reason"
          :placeholder="t('adoptions.rejectionReason')"
          rows="4"
          required
        ></textarea>
      </FormGroup>

      <template #footer>
        <BaseButton variant="outline" @click="rejectModal.show = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton
          variant="danger"
          :loading="rejectModal.loading"
          @click="rejectAdoption"
        >
          {{ t('adoptions.reject') }}
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Complete Modal -->
    <BaseModal
      v-model="completeModal.show"
      :title="t('adoptions.completeAdoption')"
      size="medium"
    >
      <FormGroup :label="t('adoptions.adoptionDate')" required>
        <input
          v-model="completeModal.adoptionDate"
          type="date"
          required
        />
      </FormGroup>

      <FormGroup :label="t('common.notes')">
        <textarea
          v-model="completeModal.notes"
          :placeholder="t('common.notes')"
          rows="3"
        ></textarea>
      </FormGroup>

      <template #footer>
        <BaseButton variant="outline" @click="completeModal.show = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton
          variant="primary"
          :loading="completeModal.loading"
          @click="completeAdoption"
        >
          {{ t('adoptions.complete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute, RouterLink } from 'vue-router'
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

const adoption = ref(null)
const loading = ref(false)
const approving = ref(false)

const canManage = computed(() => authStore.hasRole('employee'))

const rejectModal = reactive({
  show: false,
  reason: '',
  loading: false
})

const completeModal = reactive({
  show: false,
  adoptionDate: new Date().toISOString().split('T')[0],
  notes: '',
  loading: false
})

onMounted(() => {
  fetchAdoption()
})

async function fetchAdoption() {
  loading.value = true
  try {
    const response = await API.adoptions.getById(route.params.id)
    adoption.value = response.data
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
    router.push({ name: 'adoptions-list' })
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

async function approveAdoption() {
  approving.value = true
  try {
    await API.adoptions.approve(adoption.value.id)
    notificationStore.success(t('adoptions.approveSuccess'))
    fetchAdoption()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    approving.value = false
  }
}

function showRejectModal() {
  rejectModal.reason = ''
  rejectModal.show = true
}

async function rejectAdoption() {
  if (!rejectModal.reason.trim()) {
    notificationStore.error(t('common.error'), t('adoptions.rejectionReasonRequired'))
    return
  }

  rejectModal.loading = true

  try {
    await API.adoptions.reject(adoption.value.id, rejectModal.reason)
    notificationStore.success(t('adoptions.rejectSuccess'))
    rejectModal.show = false
    fetchAdoption()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    rejectModal.loading = false
  }
}

function showCompleteModal() {
  completeModal.adoptionDate = new Date().toISOString().split('T')[0]
  completeModal.notes = ''
  completeModal.show = true
}

async function completeAdoption() {
  if (!completeModal.adoptionDate) {
    notificationStore.error(t('common.error'), t('common.required'))
    return
  }

  completeModal.loading = true

  try {
    await API.adoptions.complete(adoption.value.id, {
      adoption_date: completeModal.adoptionDate,
      notes: completeModal.notes
    })
    notificationStore.success(t('adoptions.completeSuccess'))
    completeModal.show = false
    fetchAdoption()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    completeModal.loading = false
  }
}
</script>

<style scoped>
.adoption-view-page {
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

.adoption-content {
  max-width: 1200px;
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.status-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.status-badge-large {
  display: inline-block;
  padding: 0.75rem 1.5rem;
  border-radius: 0.5rem;
  font-size: 1rem;
  font-weight: 600;
  width: fit-content;
}

.status-badge--submitted {
  background-color: #e0e7ff;
  color: #3730a3;
}

.status-badge--under_review {
  background-color: #fef3c7;
  color: #92400e;
}

.status-badge--interview_scheduled {
  background-color: #dbeafe;
  color: #1e40af;
}

.status-badge--approved {
  background-color: #d1fae5;
  color: #065f46;
}

.status-badge--rejected {
  background-color: #fee2e2;
  color: #991b1b;
}

.status-badge--completed {
  background-color: #d1fae5;
  color: #065f46;
}

.status-dates {
  display: flex;
  gap: 2rem;
}

.date-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.date-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.date-value {
  font-size: 1rem;
  color: var(--text-primary);
  font-weight: 600;
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
  font-size: 1rem;
}

.info-link:hover {
  text-decoration: underline;
}

.text-section {
  margin-bottom: 1.5rem;
}

.text-section:last-child {
  margin-bottom: 0;
}

.text-section h4 {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin: 0 0 0.5rem 0;
}

.text-section p {
  line-height: 1.6;
  color: var(--text-primary);
  margin: 0;
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
  .adoption-view-page {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .page-actions {
    flex-direction: column;
    width: 100%;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .status-dates {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>
