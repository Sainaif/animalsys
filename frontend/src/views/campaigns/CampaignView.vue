<template>
  <div class="campaign-view-page">
    <LoadingSpinner v-if="loading" />
    <div v-else>
      <div class="page-header">
        <div>
          <h1 class="page-title">{{ campaign.name }}</h1>
          <span class="badge" :class="`badge-${campaign.type}`">
            {{ t(`campaigns.type${campaign.type?.charAt(0).toUpperCase() + campaign.type?.slice(1)}`) }}
          </span>
          <span class="badge" :class="`badge-status-${campaign.status}`">
            {{ t(`campaigns.status${campaign.status?.charAt(0).toUpperCase() + campaign.status?.slice(1)}`) }}
          </span>
        </div>
        <div class="header-actions">
          <BaseButton
            v-if="authStore.hasRole('staff')"
            variant="secondary"
            @click="editCampaign"
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

      <!-- Progress Section -->
      <BaseCard class="progress-card">
        <div class="progress-header">
          <h2>{{ t('campaigns.progress') }}</h2>
          <span class="progress-percentage">{{ calculateProgress() }}%</span>
        </div>
        <div class="progress-bar-large">
          <div
            class="progress-fill-large"
            :style="{ width: `${calculateProgress()}%` }"
            :class="getProgressClass()"
          ></div>
        </div>
        <div v-if="campaign.type === 'fundraising'" class="progress-details">
          <div class="progress-stat">
            <span class="stat-label">{{ t('campaigns.raised') }}:</span>
            <span class="stat-value">{{ formatCurrency(campaign.current_amount || 0) }}</span>
          </div>
          <div class="progress-stat">
            <span class="stat-label">{{ t('campaigns.goal') }}:</span>
            <span class="stat-value">{{ formatCurrency(campaign.goal_amount) }}</span>
          </div>
        </div>
        <div v-else-if="campaign.type === 'adoption'" class="progress-details">
          <div class="progress-stat">
            <span class="stat-label">{{ t('campaigns.adoptionsCompleted') }}:</span>
            <span class="stat-value">{{ campaign.current_adoptions || 0 }}</span>
          </div>
          <div class="progress-stat">
            <span class="stat-label">{{ t('campaigns.goal') }}:</span>
            <span class="stat-value">{{ campaign.goal_adoptions }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Campaign Info -->
      <div class="info-grid">
        <BaseCard>
          <template #header>{{ t('campaigns.campaignInfo') }}</template>
          <div class="info-item">
            <span class="info-label">{{ t('campaigns.startDate') }}:</span>
            <span class="info-value">{{ formatDate(campaign.start_date) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('campaigns.endDate') }}:</span>
            <span class="info-value">{{ formatDate(campaign.end_date) }}</span>
          </div>
          <div v-if="campaign.description" class="info-item full-width">
            <span class="info-label">{{ t('common.description') }}:</span>
            <p class="description-text">{{ campaign.description }}</p>
          </div>
        </BaseCard>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      :title="t('campaigns.deleteCampaign')"
      @close="showDeleteModal = false"
    >
      <p>{{ t('campaigns.deleteCampaignConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteCampaign">
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

const campaign = ref({})
const loading = ref(false)
const showDeleteModal = ref(false)

async function fetchCampaign() {
  try {
    loading.value = true
    const response = await API.campaigns.getById(route.params.id)
    campaign.value = response.data
  } catch (error) {
    console.error('Failed to fetch campaign:', error)
    notificationStore.error(t('campaigns.fetchError'))
    router.push({ name: 'Campaigns' })
  } finally {
    loading.value = false
  }
}

function editCampaign() {
  router.push({ name: 'CampaignForm', params: { id: route.params.id } })
}

async function deleteCampaign() {
  try {
    await API.campaigns.delete(route.params.id)
    notificationStore.success(t('campaigns.deleteSuccess'))
    router.push({ name: 'Campaigns' })
  } catch (error) {
    console.error('Failed to delete campaign:', error)
    notificationStore.error(t('campaigns.deleteError'))
  }
}

function calculateProgress() {
  if (campaign.value.type === 'fundraising' && campaign.value.goal_amount) {
    const progress = ((campaign.value.current_amount || 0) / campaign.value.goal_amount) * 100
    return Math.min(Math.round(progress), 100)
  } else if (campaign.value.type === 'adoption' && campaign.value.goal_adoptions) {
    const progress = ((campaign.value.current_adoptions || 0) / campaign.value.goal_adoptions) * 100
    return Math.min(Math.round(progress), 100)
  }
  return 0
}

function getProgressClass() {
  const progress = calculateProgress()
  if (progress >= 100) return 'complete'
  if (progress >= 75) return 'high'
  if (progress >= 50) return 'medium'
  return 'low'
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

onMounted(() => {
  fetchCampaign()
})
</script>

<style scoped>
.campaign-view-page {
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

.badge-fundraising {
  background: #e8f5e9;
  color: #388e3c;
}

.badge-adoption {
  background: #e3f2fd;
  color: #1976d2;
}

.badge-event {
  background: #f3e5f5;
  color: #7b1fa2;
}

.badge-awareness {
  background: #fff3e0;
  color: #f57c00;
}

.badge-status-active {
  background: #e8f5e9;
  color: #388e3c;
}

.badge-status-completed {
  background: #e3f2fd;
  color: #1976d2;
}

.badge-status-upcoming {
  background: #fff3e0;
  color: #f57c00;
}

.badge-status-cancelled {
  background: #f5f5f5;
  color: #757575;
}

.progress-card {
  margin-bottom: 2rem;
}

.progress-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1rem;
}

.progress-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.progress-percentage {
  font-size: 2rem;
  font-weight: bold;
  color: var(--primary-color);
}

.progress-bar-large {
  height: 24px;
  background: var(--bg-secondary);
  border-radius: 12px;
  overflow: hidden;
  margin-bottom: 1.5rem;
}

.progress-fill-large {
  height: 100%;
  transition: width 0.3s ease;
}

.progress-fill-large.low {
  background: linear-gradient(90deg, #f44336, #e91e63);
}

.progress-fill-large.medium {
  background: linear-gradient(90deg, #ff9800, #ff5722);
}

.progress-fill-large.high {
  background: linear-gradient(90deg, #2196f3, #00bcd4);
}

.progress-fill-large.complete {
  background: linear-gradient(90deg, #4caf50, #8bc34a);
}

.progress-details {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.progress-stat {
  display: flex;
  flex-direction: column;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--text-primary);
}

.info-grid {
  display: grid;
  gap: 1.5rem;
}

.info-item {
  display: flex;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-item.full-width {
  flex-direction: column;
}

.info-label {
  font-weight: 600;
  min-width: 150px;
  color: var(--text-secondary);
}

.info-value {
  color: var(--text-primary);
}

.description-text {
  margin-top: 0.5rem;
  white-space: pre-wrap;
  line-height: 1.6;
}
</style>
