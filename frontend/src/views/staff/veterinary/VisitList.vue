<template>
  <div class="visit-list">
    <div class="page-header">
      <h1>{{ $t('veterinary.visits') }}</h1>
      <Button
        :label="$t('veterinary.addVisit')"
        icon="pi pi-plus"
        @click="router.push('/staff/veterinary/visits/new')"
      />
    </div>

    <Card class="filters-card">
      <template #content>
        <div class="filters">
          <Dropdown
            v-model="filters.visit_type"
            :options="visitTypeOptions"
            :placeholder="$t('veterinary.visitType')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadVisits"
          />
        </div>
      </template>
    </Card>

    <Card v-if="!loading && visits.length > 0">
      <template #content>
        <DataTable :value="visits" paginator :rows="20">
          <Column field="animal.name" header="Animal">
            <template #body="slotProps">
              {{ formatAnimalName(slotProps.data.animal) }}
            </template>
          </Column>
          <Column field="visit_date" :header="$t('veterinary.visitDate')">
            <template #body="slotProps">
              {{ formatDate(slotProps.data.visit_date) }}
            </template>
          </Column>
          <Column field="visit_type" :header="$t('veterinary.visitType')">
            <template #body="slotProps">
              {{ formatVisitType(slotProps.data.visit_type) }}
            </template>
          </Column>
          <Column field="veterinarian_name" :header="$t('veterinary.veterinarianName')" />
          <Column field="reason" :header="$t('veterinary.reason')">
            <template #body="slotProps">
              {{ truncate(slotProps.data.reason, 50) }}
            </template>
          </Column>
          <Column field="follow_up_required" :header="$t('veterinary.followUpRequired')">
            <template #body="slotProps">
              <Badge :variant="slotProps.data.follow_up_required ? 'warning' : 'success'">
                {{ slotProps.data.follow_up_required ? $t('common.yes') : $t('common.no') }}
              </Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  @click="router.push(`/staff/veterinary/visits/${slotProps.data.id}/edit`)"
                />
                <Button
                  icon="pi pi-trash"
                  class="p-button-rounded p-button-text p-button-danger"
                  @click="confirmDelete(slotProps.data)"
                />
              </div>
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState v-if="!loading && visits.length === 0" :message="$t('veterinary.noVisitsFound')" />
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { veterinaryService } from '@/services/veterinaryService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import { getLocalizedValue } from '@/utils/animalHelpers'

const router = useRouter()
const { t, locale } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const visits = ref([])
const loading = ref(true)
const filters = reactive({ visit_type: null })

const visitTypeOptions = computed(() => ([
  { label: t('veterinary.checkup'), value: 'checkup' },
  { label: t('veterinary.emergency'), value: 'emergency' },
  { label: t('veterinary.surgery'), value: 'surgery' },
  { label: t('veterinary.vaccination'), value: 'vaccination' },
  { label: t('veterinary.followUp'), value: 'follow_up' },
  { label: t('veterinary.other'), value: 'other' }
]))

const loadVisits = async () => {
  try {
    loading.value = true
    const response = await veterinaryService.getVisits(filters)
    visits.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load visits', life: 3000 })
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

const formatAnimalName = (animal) => {
  if (!animal) return t('animal.unknown')
  if (typeof animal === 'string') return animal
  const localized = getLocalizedValue(animal.name || animal, locale.value)
  return localized || animal.name?.en || animal.name?.pl || t('animal.unknown')
}

const formatVisitType = (type) => {
  if (!type) return '—'
  const key = `veterinary.${type}`
  const translation = t(key)
  return translation !== key ? translation : type
}

const truncate = (text, length) => {
  if (!text) return '—'
  return text.length > length ? `${text.substring(0, length)}…` : text
}

const confirmDelete = (visit) => {
  confirm.require({
    message: 'Are you sure you want to delete this visit?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await veterinaryService.deleteVisit(visit.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('veterinary.visitDeleted'), life: 3000 })
        loadVisits()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete visit', life: 3000 })
      }
    }
  })
}

onMounted(loadVisits)
</script>

<style scoped>
.visit-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; }
.action-buttons { display: flex; gap: 0.25rem; }
</style>
