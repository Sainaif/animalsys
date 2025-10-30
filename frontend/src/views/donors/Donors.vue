<template>
  <div class="donors-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.donors') }}</h1>
      <BaseButton
        v-if="authStore.hasRole('staff')"
        variant="primary"
        @click="createDonor"
      >
        ➕ {{ t('donors.addDonor') }}
      </BaseButton>
    </div>

    <!-- Filters -->
    <BaseCard class="filters-card">
      <div class="filters">
        <FormGroup :label="t('common.search')">
          <input
            v-model="filters.search"
            type="text"
            class="form-control"
            :placeholder="t('donors.searchPlaceholder')"
            @input="handleFilterChange"
          />
        </FormGroup>

        <FormGroup :label="t('donors.type')">
          <select v-model="filters.type" class="form-control" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="individual">{{ t('donors.typeIndividual') }}</option>
            <option value="company">{{ t('donors.typeCompany') }}</option>
            <option value="foundation">{{ t('donors.typeFoundation') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('common.status')">
          <select v-model="filters.status" class="form-control" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="active">{{ t('donors.statusActive') }}</option>
            <option value="inactive">{{ t('donors.statusInactive') }}</option>
          </select>
        </FormGroup>
      </div>
    </BaseCard>

    <!-- Donors Table -->
    <BaseCard>
      <LoadingSpinner v-if="loading" />
      <EmptyState
        v-else-if="!donors || donors.length === 0"
        icon="🎁"
        :title="t('donors.noDonors')"
        :description="t('donors.noDonorsMessage')"
      />
      <DataTable
        v-else
        :columns="columns"
        :data="donors"
        :current-page="pagination.page"
        :total-pages="pagination.total_pages"
        @sort="handleSort"
        @page-change="handlePageChange"
      >
        <template #cell-type="{ row }">
          <span class="badge" :class="`badge-${row.type}`">
            {{ t(`donors.type${row.type.charAt(0).toUpperCase() + row.type.slice(1)}`) }}
          </span>
        </template>

        <template #cell-status="{ row }">
          <span class="badge" :class="`badge-${row.status}`">
            {{ t(`donors.status${row.status.charAt(0).toUpperCase() + row.status.slice(1)}`) }}
          </span>
        </template>

        <template #cell-total_donated="{ row }">
          <span class="amount-donated">{{ formatCurrency(row.total_donated) }}</span>
        </template>

        <template #cell-last_donation_date="{ row }">
          {{ row.last_donation_date ? formatDate(row.last_donation_date) : '-' }}
        </template>

        <template #cell-actions="{ row }">
          <div class="actions">
            <BaseButton
              variant="secondary"
              size="small"
              @click="viewDonor(row.id)"
            >
              {{ t('common.view') }}
            </BaseButton>
            <BaseButton
              v-if="authStore.hasRole('staff')"
              variant="secondary"
              size="small"
              @click="editDonor(row.id)"
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
      :title="t('donors.deleteDonor')"
      @close="showDeleteModal = false"
    >
      <p>{{ t('donors.deleteDonorConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteDonor">
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

const donors = ref([])
const loading = ref(false)
const showDeleteModal = ref(false)
const donorToDelete = ref(null)

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
  field: 'created_at',
  order: 'desc',
})

const columns = [
  { key: 'name', label: t('donors.name'), sortable: true },
  { key: 'type', label: t('donors.type'), sortable: true },
  { key: 'email', label: t('common.email'), sortable: false },
  { key: 'phone', label: t('common.phone'), sortable: false },
  { key: 'total_donated', label: t('donors.totalDonated'), sortable: true },
  { key: 'last_donation_date', label: t('donors.lastDonation'), sortable: true },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'actions', label: t('common.actions'), sortable: false },
]

async function fetchDonors() {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      sort_by: sort.field,
      sort_order: sort.order,
      ...filters,
    }

    const response = await API.donors.list(params)
    donors.value = response.data.data || []
    pagination.total = response.data.total || 0
    pagination.total_pages = response.data.total_pages || 0
  } catch (error) {
    console.error('Failed to fetch donors:', error)
    notificationStore.error(t('donors.fetchError'))
  } finally {
    loading.value = false
  }
}

function handleFilterChange() {
  pagination.page = 1
  fetchDonors()
}

function handleSort(field) {
  if (sort.field === field) {
    sort.order = sort.order === 'asc' ? 'desc' : 'asc'
  } else {
    sort.field = field
    sort.order = 'asc'
  }
  fetchDonors()
}

function handlePageChange(page) {
  pagination.page = page
  fetchDonors()
}

function createDonor() {
  router.push({ name: 'DonorForm' })
}

function viewDonor(id) {
  router.push({ name: 'DonorView', params: { id } })
}

function editDonor(id) {
  router.push({ name: 'DonorForm', params: { id } })
}

function confirmDelete(id) {
  donorToDelete.value = id
  showDeleteModal.value = true
}

async function deleteDonor() {
  try {
    await API.donors.delete(donorToDelete.value)
    notificationStore.success(t('donors.deleteSuccess'))
    showDeleteModal.value = false
    donorToDelete.value = null
    fetchDonors()
  } catch (error) {
    console.error('Failed to delete donor:', error)
    notificationStore.error(t('donors.deleteError'))
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

onMounted(() => {
  fetchDonors()
})
</script>

<style scoped>
.donors-page {
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

.badge-individual {
  background: #e3f2fd;
  color: #1976d2;
}

.badge-company {
  background: #f3e5f5;
  color: #7b1fa2;
}

.badge-foundation {
  background: #fff3e0;
  color: #f57c00;
}

.badge-active {
  background: #d4edda;
  color: #155724;
}

.badge-inactive {
  background: #f8d7da;
  color: #721c24;
}

.amount-donated {
  font-weight: 600;
  color: var(--success-color);
}

.actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}
</style>
