<template>
  <div class="medication-list">
    <div class="page-header">
      <h1>{{ $t('veterinary.medications') }}</h1>
      <Button
        :label="$t('veterinary.addMedication')"
        icon="pi pi-plus"
        @click="router.push('/staff/veterinary/medications/new')"
      />
    </div>

    <Card class="filters-card">
      <template #content>
        <div class="filters">
          <Dropdown
            v-model="filters.status"
            :options="statusOptions"
            :placeholder="$t('veterinary.status')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadMedications"
          />
        </div>
      </template>
    </Card>

    <Card v-if="!loading && medications.length > 0">
      <template #content>
        <DataTable
          :value="medications"
          paginator
          :rows="20"
        >
          <Column
            field="animal.name"
            header="Animal"
          >
            <template #body="slotProps">
              {{ formatAnimalName(slotProps.data.animal) }}
            </template>
          </Column>
          <Column
            field="medication_name"
            :header="$t('veterinary.medicationName')"
          />
          <Column
            field="dosage"
            :header="$t('veterinary.dosage')"
          />
          <Column
            field="frequency"
            :header="$t('veterinary.frequency')"
          />
          <Column
            field="start_date"
            :header="$t('veterinary.startDate')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.start_date) }}
            </template>
          </Column>
          <Column
            field="status"
            :header="$t('common.status')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ $t(`veterinary.${slotProps.data.status}`) }}
              </Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <Button
                icon="pi pi-trash"
                class="p-button-rounded p-button-text p-button-danger"
                @click="confirmDelete(slotProps.data)"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState
      v-if="!loading && medications.length === 0"
      :message="$t('veterinary.noMedicationsFound')"
    />
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

const medications = ref([])
const loading = ref(true)
const filters = reactive({ status: null })

const statusOptions = computed(() => ([
  { label: t('veterinary.active'), value: 'active' },
  { label: t('veterinary.completed'), value: 'completed' },
  { label: t('veterinary.discontinued'), value: 'discontinued' }
]))

const loadMedications = async () => {
  try {
    loading.value = true
    const response = await veterinaryService.getMedications(filters)
    medications.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load medications', life: 3000 })
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

const getStatusVariant = (status) => ({
  active: 'success',
  completed: 'info',
  discontinued: 'neutral'
}[status?.toLowerCase?.()] || 'neutral')

const confirmDelete = (medication) => {
  confirm.require({
    message: 'Are you sure you want to delete this medication?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await veterinaryService.deleteMedication(medication.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('veterinary.medicationDeleted'), life: 3000 })
        loadMedications()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete medication', life: 3000 })
      }
    }
  })
}

onMounted(loadMedications)
</script>

<style scoped>
.medication-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; }
</style>
