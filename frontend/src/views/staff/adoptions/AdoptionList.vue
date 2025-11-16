<template>
  <div class="adoption-list">
    <div class="page-header">
      <h1>{{ $t('adoption.title') }}</h1>
      <Button :label="$t('common.export')" icon="pi pi-download" class="p-button-secondary" @click="exportAdoptions" />
    </div>

    <Card class="filters-card">
      <template #content>
        <div class="filters">
          <Dropdown
            v-model="filters.status"
            :options="statusOptions"
            :placeholder="$t('adoption.status')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadAdoptions"
          />
        </div>
      </template>
    </Card>

    <Card v-if="!loading && adoptions.length > 0">
      <template #content>
        <DataTable :value="adoptions" paginator :rows="20">
          <Column field="animal.name" header="Animal">
            <template #body="slotProps">
                {{ formatAnimalName(slotProps.data) }}
            </template>
          </Column>
          <Column field="adopter_first_name" :header="$t('adoption.firstName')" />
          <Column field="adopter_last_name" :header="$t('adoption.lastName')" />
          <Column field="adoption_date" :header="$t('adoption.adoptionDate')">
            <template #body="slotProps">
              {{ formatDate(slotProps.data.adoption_date) }}
            </template>
          </Column>
          <Column field="adoption_fee" :header="$t('adoption.adoptionFee')">
            <template #body="slotProps">
              {{ formatCurrency(slotProps.data.adoption_fee) }}
            </template>
          </Column>
          <Column field="status" :header="$t('adoption.status')">
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ formatStatusLabel(slotProps.data.status) }}
              </Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <Button
                icon="pi pi-eye"
                class="p-button-rounded p-button-text"
                @click="router.push(`/staff/adoptions/${slotProps.data.id}`)"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState v-if="!loading && adoptions.length === 0" :message="$t('adoption.noAdoptionsFound')" />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { adoptionService } from '@/services/adoptionService'
import { exportService } from '@/services/exportService'
import { animalService } from '@/services/animalService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import { getLocalizedValue } from '@/utils/animalHelpers'

const router = useRouter()
const { t, locale } = useI18n()
const toast = useToast()

const adoptions = ref([])
const loading = ref(true)
const filters = reactive({ status: null })
const animalsCache = reactive({})

const statusOptions = computed(() => ([
  { label: t('adoption.active'), value: 'active' },
  { label: t('adoption.returned'), value: 'returned' },
  { label: t('adoption.completed'), value: 'completed' }
]))

const loadAdoptions = async () => {
  try {
    loading.value = true
    const response = await adoptionService.getAdoptions(filters)
    adoptions.value = response.data
    await preloadAnimals(response.data)
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load adoptions', life: 3000 })
  } finally {
    loading.value = false
  }
}

const currencyFormatter = computed(() => new Intl.NumberFormat(
  locale.value === 'pl' ? 'pl-PL' : 'en-US',
  {
    style: 'currency',
    currency: locale.value === 'pl' ? 'PLN' : 'USD',
    minimumFractionDigits: 0
  }
))

const dateFormatter = computed(() => new Intl.DateTimeFormat(locale.value === 'pl' ? 'pl-PL' : 'en-US'))

const formatDate = (date) => {
  if (!date) return '—'
  const parsed = new Date(date)
  if (Number.isNaN(parsed.getTime())) return '—'
  return dateFormatter.value.format(parsed)
}

const formatCurrency = (amount) => {
  if (typeof amount !== 'number') return '—'
  return currencyFormatter.value.format(amount)
}

const preloadAnimals = async (records) => {
  const missingIds = Array.from(new Set(
    records
      .map((item) => item.animal_id || item.animal?.id || item.animal)
      .filter((id) => typeof id === 'string' && id && !animalsCache[id])
  ))

  if (missingIds.length === 0) {
    return
  }

  await Promise.allSettled(
    missingIds.map(async (id) => {
      try {
        const animal = await animalService.getAnimal(id)
        animalsCache[id] = animal
      } catch (error) {
        console.warn('Failed to load animal', id, error)
      }
    })
  )
}

const resolveAnimal = (record) => {
  if (!record) {
    return null
  }
  if (record.animal && typeof record.animal === 'object') {
    return record.animal
  }
  const id = record.animal_id || (typeof record.animal === 'string' ? record.animal : null)
  if (id && animalsCache[id]) {
    return animalsCache[id]
  }
  return null
}

const formatAnimalName = (record) => {
  const animal = resolveAnimal(record)
  if (!animal) {
    return t('animal.unknown')
  }
  const localized = getLocalizedValue(animal.name || animal.Name || animal, locale.value)
  return localized || animal.name?.en || animal.name?.pl || animal.Name?.en || t('animal.unknown')
}

const formatStatusLabel = (status) => {
  if (!status) return t('common.status')
  const normalized = status.toLowerCase()
  const key = `adoption.${normalized}`
  const translation = t(key)
  return translation !== key ? translation : status
}

const getStatusVariant = (status) => ({
  active: 'success',
  returned: 'warning',
  completed: 'info'
}[status?.toLowerCase?.()] || 'neutral')

const exportAdoptions = () => {
  try {
    if (adoptions.value.length === 0) {
      toast.add({ severity: 'warn', summary: 'Warning', detail: 'No adoptions to export', life: 3000 })
      return
    }

    const columns = [
      { field: 'animal.name', header: 'Animal' },
      { field: 'adopter_first_name', header: 'First Name' },
      { field: 'adopter_last_name', header: 'Last Name' },
      { field: 'adopter_email', header: 'Email' },
      { field: 'adoption_date', header: 'Adoption Date' },
      { field: 'adoption_fee', header: 'Fee' },
      { field: 'status', header: 'Status' }
    ]

    const timestamp = new Date().toISOString().split('T')[0]
    exportService.exportToCSV(adoptions.value, `adoptions_${timestamp}.csv`, columns)

    toast.add({ severity: 'success', summary: 'Success', detail: 'Adoptions exported successfully', life: 3000 })
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to export adoptions', life: 3000 })
  }
}

onMounted(loadAdoptions)
</script>

<style scoped>
.adoption-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; }
</style>
