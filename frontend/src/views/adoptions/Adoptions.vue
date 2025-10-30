<template>
  <div class="adoptions-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('adoptions.title') }}</h1>
      <RouterLink to="/app/adoptions/apply">
        <BaseButton variant="primary">
          ➕ {{ t('adoptions.newApplication') }}
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
            :placeholder="t('adoptions.searchPlaceholder')"
            @input="handleFilterChange"
          />
        </FormGroup>

        <FormGroup :label="t('common.status')">
          <select v-model="filters.status" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="submitted">{{ t('adoptions.submitted') }}</option>
            <option value="under_review">{{ t('adoptions.underReview') }}</option>
            <option value="interview_scheduled">{{ t('adoptions.interviewScheduled') }}</option>
            <option value="approved">{{ t('adoptions.approved') }}</option>
            <option value="rejected">{{ t('adoptions.rejected') }}</option>
            <option value="completed">{{ t('adoptions.completed') }}</option>
          </select>
        </FormGroup>
      </div>
    </BaseCard>

    <!-- Table -->
    <BaseCard padding="none">
      <DataTable
        :columns="columns"
        :data="adoptions"
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
        <template #cell-animal_name="{ row }">
          <RouterLink
            v-if="row.animal_id"
            :to="`/app/animals/${row.animal_id}`"
            class="link"
          >
            {{ row.animal_name }}
          </RouterLink>
          <span v-else>-</span>
        </template>

        <template #cell-status="{ value }">
          <span :class="['status-badge', `status-badge--${value}`]">
            {{ t(`adoptions.${value}`) }}
          </span>
        </template>

        <template #cell-application_date="{ value }">
          {{ formatDate(value) }}
        </template>

        <template #actions="{ row }">
          <BaseButton
            size="small"
            variant="ghost"
            @click="viewAdoption(row.id)"
          >
            {{ t('common.view') }}
          </BaseButton>
          <BaseButton
            v-if="canApprove && row.status === 'under_review'"
            size="small"
            variant="success"
            @click="approveAdoption(row.id)"
          >
            {{ t('adoptions.approve') }}
          </BaseButton>
          <BaseButton
            v-if="canApprove && row.status === 'under_review'"
            size="small"
            variant="danger"
            @click="showRejectModal(row)"
          >
            {{ t('adoptions.reject') }}
          </BaseButton>
        </template>
      </DataTable>
    </BaseCard>

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

const adoptions = ref([])
const loading = ref(false)
const total = ref(0)

const canApprove = computed(() => authStore.hasRole('employee'))

const filters = reactive({
  search: '',
  status: ''
})

const pagination = reactive({
  page: 1,
  limit: 10
})

const sort = reactive({
  sortBy: 'application_date',
  sortOrder: 'desc'
})

const rejectModal = reactive({
  show: false,
  adoption: null,
  reason: '',
  loading: false
})

const columns = [
  { key: 'applicant_name', label: t('adoptions.applicantName'), sortable: true },
  { key: 'animal_name', label: t('animals.name'), sortable: false },
  { key: 'application_date', label: t('adoptions.applicationDate'), sortable: true },
  { key: 'status', label: t('common.status'), sortable: true }
]

onMounted(() => {
  fetchAdoptions()
})

async function fetchAdoptions() {
  loading.value = true

  try {
    const params = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      search: filters.search || undefined,
      status: filters.status || undefined
    }

    const response = await API.adoptions.list(params)
    adoptions.value = response.data.data || []
    total.value = response.data.total || 0
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    loading.value = false
  }
}

function handleFilterChange() {
  pagination.page = 1
  fetchAdoptions()
}

function handleSort({ sortBy, sortOrder }) {
  sort.sortBy = sortBy
  sort.sortOrder = sortOrder
  fetchAdoptions()
}

function handlePageChange(page) {
  pagination.page = page
  fetchAdoptions()
}

function viewAdoption(id) {
  router.push(`/app/adoptions/${id}`)
}

async function approveAdoption(id) {
  try {
    await API.adoptions.approve(id)
    notificationStore.success(t('adoptions.approveSuccess'))
    fetchAdoptions()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  }
}

function showRejectModal(adoption) {
  rejectModal.adoption = adoption
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
    await API.adoptions.reject(rejectModal.adoption.id, rejectModal.reason)
    notificationStore.success(t('adoptions.rejectSuccess'))
    rejectModal.show = false
    rejectModal.adoption = null
    rejectModal.reason = ''
    fetchAdoptions()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    rejectModal.loading = false
  }
}

function formatDate(dateString) {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString()
}
</script>

<style scoped>
.adoptions-page {
  max-width: 1400px;
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

textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
  font-family: inherit;
  resize: vertical;
}

textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

@media (max-width: 768px) {
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
