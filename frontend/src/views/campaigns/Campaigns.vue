<template>
  <div class="campaigns-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.campaigns') }}</h1>
      <BaseButton
        v-if="authStore.hasRole('staff')"
        variant="primary"
        @click="createCampaign"
      >
        ➕ {{ t('campaigns.addCampaign') }}
      </BaseButton>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">🎯</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.total_campaigns || 0 }}</div>
          <div class="stat-label">{{ t('campaigns.totalCampaigns') }}</div>
        </div>
      </div>
      <div class="stat-card success">
        <div class="stat-icon">✅</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.active_campaigns || 0 }}</div>
          <div class="stat-label">{{ t('campaigns.activeCampaigns') }}</div>
        </div>
      </div>
      <div class="stat-card warning">
        <div class="stat-icon">💰</div>
        <div class="stat-info">
          <div class="stat-value">{{ formatCurrency(stats.total_raised || 0) }}</div>
          <div class="stat-label">{{ t('campaigns.totalRaised') }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">📊</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.average_progress || 0 }}%</div>
          <div class="stat-label">{{ t('campaigns.averageProgress') }}</div>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <BaseCard class="filters-card">
      <div class="filters">
        <FormGroup :label="t('common.search')">
          <input
            v-model="filters.search"
            type="text"
            class="form-control"
            :placeholder="t('campaigns.searchPlaceholder')"
            @input="handleFilterChange"
          />
        </FormGroup>

        <FormGroup :label="t('campaigns.campaignType')">
          <select v-model="filters.type" class="form-control" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="fundraising">{{ t('campaigns.typeFundraising') }}</option>
            <option value="adoption">{{ t('campaigns.typeAdoption') }}</option>
            <option value="event">{{ t('campaigns.typeEvent') }}</option>
            <option value="awareness">{{ t('campaigns.typeAwareness') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('common.status')">
          <select v-model="filters.status" class="form-control" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="active">{{ t('campaigns.statusActive') }}</option>
            <option value="completed">{{ t('campaigns.statusCompleted') }}</option>
            <option value="upcoming">{{ t('campaigns.statusUpcoming') }}</option>
            <option value="cancelled">{{ t('campaigns.statusCancelled') }}</option>
          </select>
        </FormGroup>
      </div>
    </BaseCard>

    <!-- Campaigns Table -->
    <BaseCard>
      <LoadingSpinner v-if="loading" />
      <EmptyState
        v-else-if="!campaigns || campaigns.length === 0"
        icon="🎯"
        :title="t('campaigns.noCampaigns')"
        :description="t('campaigns.noCampaignsMessage')"
      />
      <DataTable
        v-else
        :columns="columns"
        :data="campaigns"
        :current-page="pagination.page"
        :total-pages="pagination.total_pages"
        @sort="handleSort"
        @page-change="handlePageChange"
      >
        <template #cell-type="{ row }">
          <span class="badge" :class="`badge-${row.type}`">
            {{ t(`campaigns.type${row.type?.charAt(0).toUpperCase() + row.type?.slice(1)}`) }}
          </span>
        </template>

        <template #cell-dates="{ row }">
          <div class="dates-cell">
            <div>{{ formatDate(row.start_date) }}</div>
            <div class="date-separator">→</div>
            <div>{{ formatDate(row.end_date) }}</div>
          </div>
        </template>

        <template #cell-progress="{ row }">
          <div class="progress-cell">
            <div class="progress-bar">
              <div
                class="progress-fill"
                :style="{ width: `${calculateProgress(row)}%` }"
                :class="getProgressClass(row)"
              ></div>
            </div>
            <span class="progress-text">{{ calculateProgress(row) }}%</span>
          </div>
        </template>

        <template #cell-goal="{ row }">
          <div v-if="row.type === 'fundraising'" class="goal-cell">
            <div class="goal-current">{{ formatCurrency(row.current_amount || 0) }}</div>
            <div class="goal-target">{{ t('campaigns.of') }} {{ formatCurrency(row.goal_amount) }}</div>
          </div>
          <div v-else-if="row.type === 'adoption'" class="goal-cell">
            <div class="goal-current">{{ row.current_adoptions || 0 }}</div>
            <div class="goal-target">{{ t('campaigns.of') }} {{ row.goal_adoptions }}</div>
          </div>
          <span v-else>-</span>
        </template>

        <template #cell-status="{ row }">
          <span class="badge" :class="`badge-status-${row.status}`">
            {{ t(`campaigns.status${row.status?.charAt(0).toUpperCase() + row.status?.slice(1)}`) }}
          </span>
        </template>

        <template #cell-actions="{ row }">
          <div class="actions">
            <BaseButton
              variant="secondary"
              size="small"
              @click="viewCampaign(row.id)"
            >
              {{ t('common.view') }}
            </BaseButton>
            <BaseButton
              v-if="authStore.hasRole('staff')"
              variant="secondary"
              size="small"
              @click="editCampaign(row.id)"
            >
              {{ t('common.edit') }}
            </BaseButton>
            <BaseButton
              v-if="authStore.hasRole('admin')"
              variant="danger"
              size="small"
              @click="confirmDelete(row.id)"
            >
              {{ t('common.delete') }}
            </BaseButton>
          </div>
        </template>
      </DataTable>
    </BaseCard>

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
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notifications'
import { API } from '../../api'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import DataTable from '../../components/base/DataTable.vue'
import FormGroup from '../../components/base/FormGroup.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'
import EmptyState from '../../components/base/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const campaigns = ref([])
const stats = ref({})
const loading = ref(false)
const showDeleteModal = ref(false)
const campaignToDelete = ref(null)

const filters = reactive({
  search: '',
  type: '',
  status: '',
})

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  total_pages: 0,
})

const sort = reactive({
  field: 'start_date',
  order: 'desc',
})

const columns = [
  { key: 'name', label: t('campaigns.name'), sortable: true },
  { key: 'type', label: t('campaigns.campaignType'), sortable: true },
  { key: 'dates', label: t('campaigns.dates'), sortable: false },
  { key: 'progress', label: t('campaigns.progress'), sortable: false },
  { key: 'goal', label: t('campaigns.goal'), sortable: false },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'actions', label: t('common.actions'), sortable: false },
]

async function fetchCampaigns() {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      sort_by: sort.field,
      sort_order: sort.order,
      ...filters,
    }

    const response = await API.campaigns.list(params)
    campaigns.value = response.data.data || []
    pagination.total = response.data.total || 0
    pagination.total_pages = response.data.total_pages || 0
  } catch (error) {
    console.error('Failed to fetch campaigns:', error)
    notificationStore.error(t('campaigns.fetchError'))
  } finally {
    loading.value = false
  }
}

async function fetchStatistics() {
  try {
    const response = await API.campaigns.getOverallStatistics()
    stats.value = response.data || {}
  } catch (error) {
    console.error('Failed to fetch statistics:', error)
  }
}

function handleFilterChange() {
  pagination.page = 1
  fetchCampaigns()
}

function handleSort(field) {
  if (sort.field === field) {
    sort.order = sort.order === 'asc' ? 'desc' : 'asc'
  } else {
    sort.field = field
    sort.order = 'asc'
  }
  fetchCampaigns()
}

function handlePageChange(page) {
  pagination.page = page
  fetchCampaigns()
}

function createCampaign() {
  router.push({ name: 'CampaignForm' })
}

function viewCampaign(id) {
  router.push({ name: 'CampaignView', params: { id } })
}

function editCampaign(id) {
  router.push({ name: 'CampaignForm', params: { id } })
}

function confirmDelete(id) {
  campaignToDelete.value = id
  showDeleteModal.value = true
}

async function deleteCampaign() {
  try {
    await API.campaigns.delete(campaignToDelete.value)
    notificationStore.success(t('campaigns.deleteSuccess'))
    showDeleteModal.value = false
    campaignToDelete.value = null
    fetchCampaigns()
    fetchStatistics()
  } catch (error) {
    console.error('Failed to delete campaign:', error)
    notificationStore.error(t('campaigns.deleteError'))
  }
}

function calculateProgress(campaign) {
  if (campaign.type === 'fundraising' && campaign.goal_amount) {
    const progress = ((campaign.current_amount || 0) / campaign.goal_amount) * 100
    return Math.min(Math.round(progress), 100)
  } else if (campaign.type === 'adoption' && campaign.goal_adoptions) {
    const progress = ((campaign.current_adoptions || 0) / campaign.goal_adoptions) * 100
    return Math.min(Math.round(progress), 100)
  }
  return 0
}

function getProgressClass(campaign) {
  const progress = calculateProgress(campaign)
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
  fetchCampaigns()
  fetchStatistics()
})
</script>

<style scoped>
.campaigns-page {
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--bg-secondary);
  border-radius: 8px;
  border-left: 4px solid var(--primary-color);
}

.stat-card.success {
  border-left-color: #4caf50;
}

.stat-card.warning {
  border-left-color: #ff9800;
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: bold;
  color: var(--text-primary);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.filters-card {
  margin-bottom: 1.5rem;
}

.filters {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.form-control {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.875rem;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
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

.dates-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.875rem;
}

.date-separator {
  color: var(--text-secondary);
}

.progress-cell {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.progress-bar {
  flex: 1;
  height: 8px;
  background: var(--bg-secondary);
  border-radius: 4px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  transition: width 0.3s ease;
}

.progress-fill.low {
  background: #f44336;
}

.progress-fill.medium {
  background: #ff9800;
}

.progress-fill.high {
  background: #2196f3;
}

.progress-fill.complete {
  background: #4caf50;
}

.progress-text {
  font-size: 0.875rem;
  font-weight: 600;
  min-width: 45px;
}

.goal-cell {
  display: flex;
  flex-direction: column;
}

.goal-current {
  font-weight: 600;
}

.goal-target {
  font-size: 0.75rem;
  color: var(--text-secondary);
}

.actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}
</style>
