<template>
  <div class="animals-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('animals.title') }}</h1>
      <RouterLink to="/app/animals/create">
        <BaseButton variant="primary">
          ➕ {{ t('animals.addAnimal') }}
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
            :placeholder="t('animals.searchPlaceholder')"
            @input="handleFilterChange"
          />
        </FormGroup>

        <FormGroup :label="t('animals.species')">
          <select v-model="filters.species" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="dog">{{ t('animals.dog') }}</option>
            <option value="cat">{{ t('animals.cat') }}</option>
            <option value="other">{{ t('animals.other') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('common.status')">
          <select v-model="filters.status" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="available">{{ t('animals.available') }}</option>
            <option value="adopted">{{ t('animals.adopted') }}</option>
            <option value="medical_care">{{ t('animals.medicalCare') }}</option>
            <option value="reserved">{{ t('animals.reserved') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('animals.gender')">
          <select v-model="filters.gender" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="male">{{ t('animals.male') }}</option>
            <option value="female">{{ t('animals.female') }}</option>
          </select>
        </FormGroup>
      </div>
    </BaseCard>

    <!-- Table -->
    <BaseCard padding="none">
      <DataTable
        :columns="columns"
        :data="animals"
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
          <RouterLink :to="`/app/animals/${row.id}`" class="animal-name-link">
            {{ row.name }}
          </RouterLink>
        </template>

        <template #cell-status="{ value }">
          <span :class="['status-badge', `status-badge--${value}`]">
            {{ t(`animals.${value}`) }}
          </span>
        </template>

        <template #cell-species="{ value }">
          {{ t(`animals.${value}`) }}
        </template>

        <template #cell-gender="{ value }">
          {{ t(`animals.${value}`) }}
        </template>

        <template #actions="{ row }">
          <BaseButton
            size="small"
            variant="ghost"
            @click="viewAnimal(row.id)"
          >
            {{ t('common.view') }}
          </BaseButton>
          <BaseButton
            size="small"
            variant="ghost"
            @click="editAnimal(row.id)"
          >
            {{ t('common.edit') }}
          </BaseButton>
          <BaseButton
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
      :title="t('animals.deleteAnimal')"
      size="small"
    >
      <p>{{ t('animals.deleteConfirm') }}</p>
      <p v-if="deleteModal.animal" class="animal-name">
        <strong>{{ deleteModal.animal.name }}</strong>
      </p>

      <template #footer>
        <BaseButton variant="secondary" @click="deleteModal.show = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" :loading="deleteModal.loading" @click="deleteAnimal">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { API } from '../../api'
import { useNotificationStore } from '../../stores/notification'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import DataTable from '../../components/base/DataTable.vue'
import FormGroup from '../../components/base/FormGroup.vue'

const { t } = useI18n()
const router = useRouter()
const notificationStore = useNotificationStore()

const animals = ref([])
const loading = ref(false)
const total = ref(0)

const filters = reactive({
  search: '',
  species: '',
  status: '',
  gender: ''
})

const pagination = reactive({
  page: 1,
  limit: 10
})

const sort = reactive({
  sortBy: 'name',
  sortOrder: 'asc'
})

const deleteModal = reactive({
  show: false,
  animal: null,
  loading: false
})

const columns = [
  { key: 'name', label: t('common.name'), sortable: true },
  { key: 'species', label: t('animals.species'), sortable: true },
  { key: 'breed', label: t('animals.breed'), sortable: false },
  { key: 'gender', label: t('animals.gender'), sortable: false },
  { key: 'age', label: t('animals.age'), sortable: true },
  { key: 'status', label: t('common.status'), sortable: true }
]

onMounted(() => {
  fetchAnimals()
})

async function fetchAnimals() {
  loading.value = true

  try {
    const params = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      search: filters.search || undefined,
      species: filters.species || undefined,
      status: filters.status || undefined,
      gender: filters.gender || undefined
    }

    const response = await API.animals.list(params)
    animals.value = response.data.data || []
    total.value = response.data.total || 0
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    loading.value = false
  }
}

function handleFilterChange() {
  pagination.page = 1
  fetchAnimals()
}

function handleSort({ sortBy, sortOrder }) {
  sort.sortBy = sortBy
  sort.sortOrder = sortOrder
  fetchAnimals()
}

function handlePageChange(page) {
  pagination.page = page
  fetchAnimals()
}

function viewAnimal(id) {
  router.push(`/app/animals/${id}`)
}

function editAnimal(id) {
  router.push(`/app/animals/${id}/edit`)
}

function confirmDelete(animal) {
  deleteModal.animal = animal
  deleteModal.show = true
}

async function deleteAnimal() {
  if (!deleteModal.animal) return

  deleteModal.loading = true

  try {
    await API.animals.delete(deleteModal.animal.id)
    notificationStore.success(t('animals.deleteSuccess'))
    deleteModal.show = false
    deleteModal.animal = null
    fetchAnimals()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    deleteModal.loading = false
  }
}
</script>

<style scoped>
.animals-page {
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

.animal-name-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.animal-name-link:hover {
  text-decoration: underline;
}

.status-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
}

.status-badge--available {
  background-color: #d1fae5;
  color: #065f46;
}

.status-badge--adopted {
  background-color: #dbeafe;
  color: #1e40af;
}

.status-badge--medical_care {
  background-color: #fef3c7;
  color: #92400e;
}

.status-badge--reserved {
  background-color: #e0e7ff;
  color: #3730a3;
}

.status-badge--deceased {
  background-color: #f3f4f6;
  color: #374151;
}

.animal-name {
  margin: 1rem 0;
  text-align: center;
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
