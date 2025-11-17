<template>
  <div class="application-list">
    <div class="page-header">
      <h1>{{ $t('adoption.applications') }}</h1>
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
            @change="loadApplications"
          />
        </div>
      </template>
    </Card>

    <Card v-if="!loading && applications.length > 0">
      <template #content>
        <DataTable
          :value="applications"
          paginator
          :rows="20"
        >
          <Column
            field="animal.name"
            header="Animal"
          >
            <template #body="slotProps">
              {{ formatAnimalName(slotProps.data) }}
            </template>
          </Column>
          <Column
            field="applicant_first_name"
            :header="$t('adoption.firstName')"
          />
          <Column
            field="applicant_last_name"
            :header="$t('adoption.lastName')"
          />
          <Column
            field="email"
            :header="$t('adoption.email')"
          />
          <Column
            field="application_date"
            :header="$t('adoption.applicationDate')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.application_date) }}
            </template>
          </Column>
          <Column
            field="status"
            :header="$t('adoption.status')"
          >
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
                @click="router.push(`/staff/adoptions/applications/${slotProps.data.id}`)"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState
      v-if="!loading && applications.length === 0"
      :message="$t('adoption.noApplicationsFound')"
    />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { adoptionService } from '@/services/adoptionService'
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

const applications = ref([])
const loading = ref(true)
const filters = reactive({ status: null })
const animalsCache = reactive({})

const statusOptions = computed(() => ([
  { label: t('adoption.pending'), value: 'pending' },
  { label: t('adoption.underReview'), value: 'under_review' },
  { label: t('adoption.approved'), value: 'approved' },
  { label: t('adoption.rejected'), value: 'rejected' }
]))

const loadApplications = async () => {
  try {
    loading.value = true
    const response = await adoptionService.getApplications(filters)
    applications.value = response.data
    await preloadAnimals(response.data)
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load applications', life: 3000 })
  } finally {
    loading.value = false
  }
}

const dateFormatter = computed(() => new Intl.DateTimeFormat(locale.value === 'pl' ? 'pl-PL' : 'en-US'))

const formatDate = (date) => {
  if (!date) return '—'
  const parsed = new Date(date)
  if (Number.isNaN(parsed.getTime())) return '—'
  return dateFormatter.value.format(parsed)
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
  pending: 'warning',
  under_review: 'info',
  approved: 'success',
  rejected: 'danger'
}[status?.toLowerCase?.()] || 'neutral')

onMounted(loadApplications)
</script>

<style scoped>
.application-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; }
</style>
