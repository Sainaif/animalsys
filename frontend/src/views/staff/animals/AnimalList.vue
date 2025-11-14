<template>
  <div class="animal-list">
    <div class="page-header">
      <h1>{{ $t('animal.title') }}</h1>
      <div class="header-actions">
        <Button
          :label="$t('common.export')"
          icon="pi pi-download"
          class="p-button-secondary"
          @click="exportAnimals"
        />
        <Button
          :label="$t('animal.addAnimal')"
          icon="pi pi-plus"
          @click="router.push('/staff/animals/new')"
        />
      </div>
    </div>

    <!-- Filters -->
    <Card class="filters-card">
      <template #content>
        <div class="filters">
          <span class="p-input-icon-left search-box">
            <i class="pi pi-search" />
            <InputText
              v-model="filters.search"
              :placeholder="$t('animal.searchAnimals')"
              @input="handleSearch"
            />
          </span>

          <Dropdown
            v-model="filters.species"
            :options="speciesOptions"
            :placeholder="$t('animal.filterBySpecies')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadAnimals"
          />

          <Dropdown
            v-model="filters.status"
            :options="statusOptions"
            :placeholder="$t('animal.filterByStatus')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadAnimals"
          />

          <Dropdown
            v-model="filters.size"
            :options="sizeOptions"
            :placeholder="$t('animal.filterBySize')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadAnimals"
          />
        </div>
      </template>
    </Card>

    <!-- Animals Table/Grid -->
    <Card v-if="!loading && animals.length > 0" class="animals-card">
      <template #content>
        <DataTable
          :value="animals"
          :rows="pagination.limit"
          :total-records="pagination.total"
          :lazy="true"
          paginator
          @page="onPage"
        >
          <Column field="photo_url" header="" style="width: 80px">
            <template #body="slotProps">
              <img
                v-if="slotProps.data.photo_url"
                :src="slotProps.data.photo_url"
                :alt="slotProps.data.name"
                class="animal-thumbnail"
              />
              <div v-else class="animal-placeholder">
                <i class="pi pi-image"></i>
              </div>
            </template>
          </Column>

          <Column field="name" :header="$t('animal.name')" sortable>
            <template #body="slotProps">
              <router-link :to="`/staff/animals/${slotProps.data.id}`" class="animal-link">
                {{ slotProps.data.name }}
              </router-link>
            </template>
          </Column>

          <Column field="species" :header="$t('animal.species')" sortable />

          <Column field="breed" :header="$t('animal.breed')" sortable />

          <Column field="age" :header="$t('animal.age')">
            <template #body="slotProps">
              {{ formatAge(slotProps.data) }}
            </template>
          </Column>

          <Column field="sex" :header="$t('animal.gender')">
            <template #body="slotProps">
              {{ $t(`animal.${slotProps.data.sex}`) }}
            </template>
          </Column>

          <Column field="status" :header="$t('animal.status')">
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ $t(`animal.${slotProps.data.status}`) }}
              </Badge>
            </template>
          </Column>

          <Column :header="$t('common.actions')" style="width: 150px">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  icon="pi pi-eye"
                  class="p-button-rounded p-button-text"
                  v-tooltip.top="'View'"
                  @click="router.push(`/staff/animals/${slotProps.data.id}`)"
                />
                <Button
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  v-tooltip.top="'Edit'"
                  @click="router.push(`/staff/animals/${slotProps.data.id}/edit`)"
                />
                <Button
                  icon="pi pi-trash"
                  class="p-button-rounded p-button-text p-button-danger"
                  v-tooltip.top="'Delete'"
                  @click="confirmDelete(slotProps.data)"
                />
              </div>
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />

    <EmptyState
      v-if="!loading && animals.length === 0"
      :message="$t('animal.noAnimalsFound')"
      :action-text="$t('animal.addAnimal')"
      icon="pi-heart"
      @action="router.push('/staff/animals/new')"
    />

    <!-- Delete Confirmation -->
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { animalService } from '@/services/animalService'
import { exportService } from '@/services/exportService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const animals = ref([])
const loading = ref(true)
const pagination = reactive({
  limit: 20,
  offset: 0,
  total: 0
})

const filters = reactive({
  search: '',
  species: null,
  status: null,
  size: null
})

const speciesOptions = [
  { label: 'Dog', value: 'dog' },
  { label: 'Cat', value: 'cat' },
  { label: 'Rabbit', value: 'rabbit' },
  { label: 'Bird', value: 'bird' },
  { label: 'Other', value: 'other' }
]

const statusOptions = ref([
  { label: t('animal.available'), value: 'available' },
  { label: t('animal.adopted'), value: 'adopted' },
  { label: t('animal.underTreatment'), value: 'under_treatment' },
  { label: t('animal.fostered'), value: 'fostered' },
  { label: t('animal.transferred'), value: 'transferred' }
])

const sizeOptions = ref([
  { label: t('animal.small'), value: 'small' },
  { label: t('animal.medium'), value: 'medium' },
  { label: t('animal.large'), value: 'large' },
  { label: t('animal.extraLarge'), value: 'extra_large' }
])

const loadAnimals = async () => {
  try {
    loading.value = true
    const params = {
      limit: pagination.limit,
      offset: pagination.offset,
      ...filters
    }

    // Remove null/undefined values
    Object.keys(params).forEach(key => {
      if (params[key] === null || params[key] === undefined || params[key] === '') {
        delete params[key]
      }
    })

    const response = await animalService.getAnimals(params)
    animals.value = response.data
    pagination.total = response.total
  } catch (error) {
    console.error('Error loading animals:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load animals',
      life: 3000
    })
  } finally {
    loading.value = false
  }
}

let searchTimeout = null
const handleSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    pagination.offset = 0
    loadAnimals()
  }, 500)
}

const onPage = (event) => {
  pagination.offset = event.first
  pagination.limit = event.rows
  loadAnimals()
}

const formatAge = (animal) => {
  if (animal.age_years) {
    return `${animal.age_years} ${t('animal.ageYears')}`
  }
  if (animal.age_months) {
    return `${animal.age_months} ${t('animal.ageMonths')}`
  }
  return 'N/A'
}

const getStatusVariant = (status) => {
  const variants = {
    available: 'success',
    adopted: 'info',
    under_treatment: 'warning',
    fostered: 'primary',
    transferred: 'neutral',
    deceased: 'danger'
  }
  return variants[status] || 'neutral'
}

const confirmDelete = (animal) => {
  confirm.require({
    message: t('animal.confirmDelete'),
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await animalService.deleteAnimal(animal.id)
        toast.add({
          severity: 'success',
          summary: 'Success',
          detail: t('animal.animalDeleted'),
          life: 3000
        })
        loadAnimals()
      } catch (error) {
        console.error('Error deleting animal:', error)
        toast.add({
          severity: 'error',
          summary: 'Error',
          detail: 'Failed to delete animal',
          life: 3000
        })
      }
    }
  })
}

const exportAnimals = () => {
  try {
    if (animals.value.length === 0) {
      toast.add({
        severity: 'warn',
        summary: 'Warning',
        detail: 'No animals to export',
        life: 3000
      })
      return
    }

    const columns = [
      { field: 'name', header: 'Name' },
      { field: 'species', header: 'Species' },
      { field: 'breed', header: 'Breed' },
      { field: 'sex', header: 'Gender' },
      { field: 'status', header: 'Status' },
      { field: 'intake_date', header: 'Intake Date' }
    ]

    const timestamp = new Date().toISOString().split('T')[0]
    exportService.exportToCSV(animals.value, `animals_${timestamp}.csv`, columns)

    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: 'Animals exported successfully',
      life: 3000
    })
  } catch (error) {
    console.error('Error exporting animals:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to export animals',
      life: 3000
    })
  }
}

onMounted(() => {
  loadAnimals()
})
</script>

<style scoped>
.animal-list {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.filters-card {
  margin-bottom: 1.5rem;
}

.filters {
  display: grid;
  grid-template-columns: 2fr 1fr 1fr 1fr;
  gap: 1rem;
}

.search-box {
  width: 100%;
}

.search-box input {
  width: 100%;
}

.animals-card {
  margin-bottom: 2rem;
}

.animal-thumbnail {
  width: 60px;
  height: 60px;
  object-fit: cover;
  border-radius: 8px;
}

.animal-placeholder {
  width: 60px;
  height: 60px;
  background: #f3f4f6;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #9ca3af;
  font-size: 1.5rem;
}

.animal-link {
  color: #3b82f6;
  text-decoration: none;
  font-weight: 600;
}

.animal-link:hover {
  text-decoration: underline;
}

.action-buttons {
  display: flex;
  gap: 0.25rem;
}

@media (max-width: 968px) {
  .filters {
    grid-template-columns: 1fr;
  }
}
</style>
