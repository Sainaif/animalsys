<template>
  <div class="volunteers-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('volunteers.title') }}</h1>
      <RouterLink to="/app/volunteers/create">
        <BaseButton variant="primary">
          ➕ {{ t('volunteers.addVolunteer') }}
        </BaseButton>
      </RouterLink>
    </div>

    <!-- Filters -->
    <BaseCard class="filters-card">
      <div class="filters">
        <FormGroup :label="t('common.search')">
          <input
            v-model="filters.search"
            type="text"
            :placeholder="t('volunteers.searchPlaceholder')"
            @input="handleFilterChange"
          />
        </FormGroup>

        <FormGroup :label="t('common.status')">
          <select v-model="filters.status" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="active">{{ t('volunteers.statusActive') }}</option>
            <option value="inactive">{{ t('volunteers.statusInactive') }}</option>
            <option value="on_leave">{{ t('volunteers.statusOnLeave') }}</option>
          </select>
        </FormGroup>
      </div>
    </BaseCard>

    <!-- Table -->
    <BaseCard padding="none">
      <DataTable
        :columns="columns"
        :data="volunteers"
        :loading="loading"
        :total="total"
        :current-page="pagination.page"
        :per-page="pagination.limit"
        :sort-by="sort.sortBy"
        :sort-order="sort.sortOrder"
        has-actions
        @sort="handleSort"
        @page-change="handlePageChange"
      >
        <template #cell-name="{ row }">
          <RouterLink :to="`/app/volunteers/${row.id}`" class="link">
            {{ row.first_name }} {{ row.last_name }}
          </RouterLink>
        </template>

        <template #cell-status="{ value }">
          <span :class="['status-badge', `status-badge--${value}`]">
            {{ t(`volunteers.status${capitalize(value)}`) }}
          </span>
        </template>

        <template #cell-registration_date="{ value }">
          {{ formatDate(value) }}
        </template>

        <template #cell-total_hours="{ value }">
          {{ value || 0 }} {{ t('volunteers.hoursShort') }}
        </template>

        <template #actions="{ row }">
          <BaseButton
            size="small"
            variant="ghost"
            @click="viewVolunteer(row.id)"
          >
            {{ t('common.view') }}
          </BaseButton>
          <BaseButton
            size="small"
            variant="ghost"
            @click="editVolunteer(row.id)"
          >
            {{ t('common.edit') }}
          </BaseButton>
          <BaseButton
            v-if="canDelete"
            size="small"
            variant="danger"
            @click="confirmDelete(row)"
          >
            {{ t('common.delete') }}
          </BaseButton>
        </template>
      </DataTable>
    </BaseCard>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-model="deleteModal.show"
      :title="t('volunteers.deleteVolunteer')"
      size="small"
    >
      <p>{{ t('volunteers.deleteConfirm') }}</p>
      <p v-if="deleteModal.volunteer" class="volunteer-name">
        <strong>{{ deleteModal.volunteer.first_name }} {{ deleteModal.volunteer.last_name }}</strong>
      </p>

      <template #footer>
        <BaseButton variant="outline" @click="deleteModal.show = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" :loading="deleteModal.loading" @click="deleteVolunteer">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { API } from '../../api'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notification'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import DataTable from '../../components/base/DataTable.vue'
import FormGroup from '../../components/base/FormGroup.vue'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const volunteers = ref([])
const loading = ref(false)
const total = ref(0)

const canDelete = computed(() => authStore.hasRole('admin'))

const filters = reactive({
  search: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  limit: 10
})

const sort = reactive({
  sortBy: 'registration_date',
  sortOrder: 'desc'
})

const deleteModal = reactive({
  show: false,
  volunteer: null,
  loading: false
})

const columns = [
  { key: 'name', label: t('common.name'), sortable: true },
  { key: 'email', label: t('common.email'), sortable: false },
  { key: 'phone', label: t('common.phone'), sortable: false },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'registration_date', label: t('volunteers.registrationDate'), sortable: true },
  { key: 'total_hours', label: t('volunteers.totalHours'), sortable: true }
]

onMounted(() => {
  fetchVolunteers()
})

async function fetchVolunteers() {
  loading.value = true

  try {
    const params = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      search: filters.search || undefined,
      status: filters.status || undefined
    }

    const response = await API.volunteers.list(params)
    volunteers.value = response.data.data || []
    total.value = response.data.total || 0
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    loading.value = false
  }
}

function handleFilterChange() {
  pagination.page = 1
  fetchVolunteers()
}

function handleSort({ sortBy, sortOrder }) {
  sort.sortBy = sortBy
  sort.sortOrder = sortOrder
  fetchVolunteers()
}

function handlePageChange(page) {
  pagination.page = page
  fetchVolunteers()
}

function viewVolunteer(id) {
  router.push({ name: 'volunteer-view', params: { id } })
}

function editVolunteer(id) {
  router.push({ name: 'volunteer-edit', params: { id } })
}

function confirmDelete(volunteer) {
  deleteModal.volunteer = volunteer
  deleteModal.show = true
}

async function deleteVolunteer() {
  if (!deleteModal.volunteer) return

  deleteModal.loading = true

  try {
    await API.volunteers.delete(deleteModal.volunteer.id)
    notificationStore.success(t('volunteers.deleteSuccess'))
    deleteModal.show = false
    deleteModal.volunteer = null
    fetchVolunteers()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    deleteModal.loading = false
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
</script>

<style scoped>
.volunteers-page {
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
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.link:hover {
  text-decoration: underline;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
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

.volunteer-name {
  margin: 1rem 0;
  text-align: center;
}

@media (max-width: 768px) {
  .volunteers-page {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .filters {
    grid-template-columns: 1fr;
  }
}
</style>
