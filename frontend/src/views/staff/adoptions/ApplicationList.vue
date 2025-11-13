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
        <DataTable :value="applications" paginator :rows="20">
          <Column field="animal.name" header="Animal">
            <template #body="slotProps">
              {{ slotProps.data.animal?.name || 'N/A' }}
            </template>
          </Column>
          <Column field="applicant_first_name" :header="$t('adoption.firstName')" />
          <Column field="applicant_last_name" :header="$t('adoption.lastName')" />
          <Column field="email" :header="$t('adoption.email')" />
          <Column field="application_date" :header="$t('adoption.applicationDate')">
            <template #body="slotProps">
              {{ formatDate(slotProps.data.application_date) }}
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
                @click="router.push(`/staff/adoptions/applications/${slotProps.data.id}`)"
              />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState v-if="!loading && applications.length === 0" :message="$t('adoption.noApplicationsFound')" />
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

const applications = ref([])
const loading = ref(true)
const filters = reactive({ status: null })

const statusOptions = ref([
  { label: t('adoption.pending'), value: 'pending' },
  { label: t('adoption.underReview'), value: 'under_review' },
  { label: t('adoption.approved'), value: 'approved' },
  { label: t('adoption.rejected'), value: 'rejected' }
])

const loadApplications = async () => {
  try {
    loading.value = true
    const response = await adoptionService.getApplications(filters)
    applications.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load applications', life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const getStatusVariant = (status) => ({
  pending: 'warning', under_review: 'info', approved: 'success', rejected: 'danger'
}[status] || 'neutral')

onMounted(loadApplications)
</script>

<style scoped>
.application-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; }
</style>
