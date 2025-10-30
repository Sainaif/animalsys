<template>
  <div class="veterinary-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.veterinary') }}</h1>
      <BaseButton
        v-if="authStore.hasRole('staff')"
        variant="primary"
        @click="createVisit"
      >
        ➕ {{ t('veterinary.addVisit') }}
      </BaseButton>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">🏥</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.total_visits || 0 }}</div>
          <div class="stat-label">{{ t('veterinary.totalVisits') }}</div>
        </div>
      </div>
      <div class="stat-card warning">
        <div class="stat-icon">📅</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.upcoming_count || 0 }}</div>
          <div class="stat-label">{{ t('veterinary.upcomingVisits') }}</div>
        </div>
      </div>
      <div class="stat-card success">
        <div class="stat-icon">💉</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.vaccinations_this_month || 0 }}</div>
          <div class="stat-label">{{ t('veterinary.vaccinationsThisMonth') }}</div>
        </div>
      </div>
      <div class="stat-card">
        <div class="stat-icon">🩺</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.checkups_this_month || 0 }}</div>
          <div class="stat-label">{{ t('veterinary.checkupsThisMonth') }}</div>
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
            :placeholder="t('veterinary.searchPlaceholder')"
            @input="handleFilterChange"
          />
        </FormGroup>

        <FormGroup :label="t('veterinary.visitType')">
          <select v-model="filters.type" class="form-control" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="checkup">{{ t('veterinary.typeCheckup') }}</option>
            <option value="vaccination">{{ t('veterinary.typeVaccination') }}</option>
            <option value="treatment">{{ t('veterinary.typeTreatment') }}</option>
            <option value="surgery">{{ t('veterinary.typeSurgery') }}</option>
            <option value="emergency">{{ t('veterinary.typeEmergency') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('common.status')">
          <select v-model="filters.status" class="form-control" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="scheduled">{{ t('veterinary.statusScheduled') }}</option>
            <option value="completed">{{ t('veterinary.statusCompleted') }}</option>
            <option value="cancelled">{{ t('veterinary.statusCancelled') }}</option>
          </select>
        </FormGroup>
      </div>
    </BaseCard>

    <!-- Visits Table -->
    <BaseCard>
      <LoadingSpinner v-if="loading" />
      <EmptyState
        v-else-if="!visits || visits.length === 0"
        icon="🏥"
        :title="t('veterinary.noVisits')"
        :description="t('veterinary.noVisitsMessage')"
      />
      <DataTable
        v-else
        :columns="columns"
        :data="visits"
        :current-page="pagination.page"
        :total-pages="pagination.total_pages"
        @sort="handleSort"
        @page-change="handlePageChange"
      >
        <template #cell-animal="{ row }">
          <router-link :to="{ name: 'AnimalView', params: { id: row.animal_id } }" class="animal-link">
            {{ row.animal?.name || '-' }}
          </router-link>
        </template>

        <template #cell-type="{ row }">
          <span class="badge" :class="`badge-${row.type}`">
            {{ t(`veterinary.type${row.type?.charAt(0).toUpperCase() + row.type?.slice(1)}`) }}
          </span>
        </template>

        <template #cell-visit_date="{ row }">
          {{ formatDateTime(row.visit_date) }}
        </template>

        <template #cell-veterinarian="{ row }">
          {{ row.veterinarian?.name || row.veterinarian_name || '-' }}
        </template>

        <template #cell-status="{ row }">
          <span class="badge" :class="`badge-status-${row.status}`">
            {{ t(`veterinary.status${row.status?.charAt(0).toUpperCase() + row.status?.slice(1)}`) }}
          </span>
        </template>

        <template #cell-cost="{ row }">
          {{ row.cost ? formatCurrency(row.cost) : '-' }}
        </template>

        <template #cell-actions="{ row }">
          <div class="actions">
            <BaseButton
              variant="secondary"
              size="small"
              @click="viewVisit(row.id)"
            >
              {{ t('common.view') }}
            </BaseButton>
            <BaseButton
              v-if="authStore.hasRole('staff')"
              variant="secondary"
              size="small"
              @click="editVisit(row.id)"
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

const visits = ref([])
const stats = ref({})
const loading = ref(false)
const showDeleteModal = ref(false)
const visitToDelete = ref(null)

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
  field: 'visit_date',
  order: 'desc',
})

const columns = [
  { key: 'animal', label: t('veterinary.animal'), sortable: false },
  { key: 'type', label: t('veterinary.visitType'), sortable: true },
  { key: 'visit_date', label: t('veterinary.visitDate'), sortable: true },
  { key: 'veterinarian', label: t('veterinary.veterinarian'), sortable: false },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'cost', label: t('veterinary.cost'), sortable: true },
  { key: 'actions', label: t('common.actions'), sortable: false },
]

async function fetchVisits() {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      sort_by: sort.field,
      sort_order: sort.order,
      ...filters,
    }

    const response = await API.veterinary.list(params)
    visits.value = response.data.data || []
    pagination.total = response.data.total || 0
    pagination.total_pages = response.data.total_pages || 0
  } catch (error) {
    console.error('Failed to fetch veterinary visits:', error)
    notificationStore.error(t('veterinary.fetchError'))
  } finally {
    loading.value = false
  }
}

async function fetchStatistics() {
  try {
    const response = await API.veterinary.getStatistics()
    stats.value = response.data || {}
  } catch (error) {
    console.error('Failed to fetch statistics:', error)
  }
}

function handleFilterChange() {
  pagination.page = 1
  fetchVisits()
}

function handleSort(field) {
  if (sort.field === field) {
    sort.order = sort.order === 'asc' ? 'desc' : 'asc'
  } else {
    sort.field = field
    sort.order = 'asc'
  }
  fetchVisits()
}

function handlePageChange(page) {
  pagination.page = page
  fetchVisits()
}

function createVisit() {
  router.push({ name: 'VeterinaryForm' })
}

function viewVisit(id) {
  router.push({ name: 'VeterinaryView', params: { id } })
}

function editVisit(id) {
  router.push({ name: 'VeterinaryForm', params: { id } })
}

function confirmDelete(id) {
  visitToDelete.value = id
  showDeleteModal.value = true
}

async function deleteVisit() {
  try {
    await API.veterinary.delete(visitToDelete.value)
    notificationStore.success(t('veterinary.deleteSuccess'))
    showDeleteModal.value = false
    visitToDelete.value = null
    fetchVisits()
    fetchStatistics()
  } catch (error) {
    console.error('Failed to delete veterinary visit:', error)
    notificationStore.error(t('veterinary.deleteError'))
  }
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
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
  fetchVisits()
  fetchStatistics()
})
</script>

<style scoped>
.veterinary-page {
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

.stat-card.warning {
  border-left-color: #ff9800;
}

.stat-card.success {
  border-left-color: #4caf50;
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

.animal-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.animal-link:hover {
  text-decoration: underline;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
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

.actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}
</style>
