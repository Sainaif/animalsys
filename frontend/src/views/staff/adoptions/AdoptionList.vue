<template>
  <div class="adoption-list">
    <div class="page-header">
      <h1>{{ $t('adoption.title') }}</h1>
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
              {{ slotProps.data.animal?.name || 'N/A' }}
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
                {{ $t(`adoption.${slotProps.data.status}`) }}
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
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { adoptionService } from '@/services/adoptionService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()

const adoptions = ref([])
const loading = ref(true)
const filters = reactive({ status: null })

const statusOptions = ref([
  { label: t('adoption.active'), value: 'active' },
  { label: t('adoption.returned'), value: 'returned' },
  { label: t('adoption.completed'), value: 'completed' }
])

const loadAdoptions = async () => {
  try {
    loading.value = true
    const response = await adoptionService.getAdoptions(filters)
    adoptions.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load adoptions', life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const formatCurrency = (amount) => amount ? `$${amount.toFixed(2)}` : 'N/A'
const getStatusVariant = (status) => ({
  active: 'success', returned: 'warning', completed: 'info'
}[status] || 'neutral')

onMounted(loadAdoptions)
</script>

<style scoped>
.adoption-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; }
</style>
