<template>
  <div class="vaccination-list">
    <div class="page-header">
      <h1>{{ $t('veterinary.vaccinations') }}</h1>
      <Button :label="$t('veterinary.addVaccination')" icon="pi pi-plus" @click="router.push('/staff/veterinary/vaccinations/new')" />
    </div>

    <Card v-if="!loading && vaccinations.length > 0">
      <template #content>
        <DataTable :value="vaccinations" paginator :rows="20">
          <Column field="animal.name" header="Animal">
            <template #body="slotProps">{{ formatAnimalName(slotProps.data.animal) }}</template>
          </Column>
          <Column field="vaccine_name" :header="$t('veterinary.vaccineName')" />
          <Column field="vaccination_date" :header="$t('veterinary.vaccinationDate')">
            <template #body="slotProps">{{ formatDate(slotProps.data.vaccination_date) }}</template>
          </Column>
          <Column field="next_due_date" :header="$t('veterinary.nextDueDate')">
            <template #body="slotProps">{{ formatDate(slotProps.data.next_due_date) }}</template>
          </Column>
          <Column field="veterinarian_name" :header="$t('veterinary.veterinarianName')" />
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <Button icon="pi pi-trash" class="p-button-rounded p-button-text p-button-danger" @click="confirmDelete(slotProps.data)" />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState v-if="!loading && vaccinations.length === 0" :message="$t('veterinary.noVaccinationsFound')" />
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { veterinaryService } from '@/services/veterinaryService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ConfirmDialog from 'primevue/confirmdialog'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import { getLocalizedValue } from '@/utils/animalHelpers'

const router = useRouter()
const { t, locale } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const vaccinations = ref([])
const loading = ref(true)

const loadVaccinations = async () => {
  try {
    loading.value = true
    const response = await veterinaryService.getVaccinations()
    vaccinations.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load vaccinations', life: 3000 })
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

const confirmDelete = (vaccination) => {
  confirm.require({
    message: 'Are you sure you want to delete this vaccination?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await veterinaryService.deleteVaccination(vaccination.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('veterinary.vaccinationDeleted'), life: 3000 })
        loadVaccinations()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete vaccination', life: 3000 })
      }
    }
  })
}

onMounted(loadVaccinations)
</script>

<style scoped>
.vaccination-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
</style>
