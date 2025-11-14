<template>
  <div class="donor-list">
    <div class="page-header">
      <h1>{{ $t('finance.donors') }}</h1>
      <div class="header-actions">
        <Button :label="$t('common.export')" icon="pi pi-download" class="p-button-secondary" @click="exportDonors" />
        <Button :label="$t('finance.addDonor')" icon="pi pi-plus" @click="router.push('/staff/finance/donors/new')" />
      </div>
    </div>

    <Card class="filters-card">
      <template #content>
        <div class="filters">
          <InputText v-model="filters.search" :placeholder="$t('common.search')" @input="handleSearch" />
          <Dropdown v-model="filters.donor_type" :options="donorTypeOptions" :placeholder="$t('finance.donorType')" option-label="label" option-value="value" show-clear @change="loadDonors" />
          <Dropdown v-model="filters.donor_status" :options="statusOptions" :placeholder="$t('finance.donorStatus')" option-label="label" option-value="value" show-clear @change="loadDonors" />
        </div>
      </template>
    </Card>

    <Card v-if="!loading && donors.length > 0">
      <template #content>
        <DataTable :value="donors" paginator :rows="20">
          <Column field="first_name" :header="$t('finance.firstName')">
            <template #body="slotProps">
              {{ slotProps.data.organization_name || `${slotProps.data.first_name || ''} ${slotProps.data.last_name || ''}` }}
            </template>
          </Column>
          <Column field="email" :header="$t('finance.email')" />
          <Column field="donor_type" :header="$t('finance.donorType')">
            <template #body="slotProps">{{ $t(`finance.${slotProps.data.donor_type}`) }}</template>
          </Column>
          <Column field="total_donated" :header="$t('finance.totalDonated')">
            <template #body="slotProps">{{ formatCurrency(slotProps.data.total_donated) }}</template>
          </Column>
          <Column field="donation_count" :header="$t('finance.donationCount')" />
          <Column field="donor_status" :header="$t('common.status')">
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.donor_status)">{{ $t(`finance.${slotProps.data.donor_status}`) }}</Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <Button icon="pi pi-trash" class="p-button-rounded p-button-text p-button-danger" @click="confirmDelete(slotProps.data)" />
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState v-if="!loading && donors.length === 0" :message="$t('finance.noDonorsFound')" />
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { financeService } from '@/services/financeService'
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

const donors = ref([])
const loading = ref(true)
const filters = reactive({ search: '', donor_type: null, donor_status: null })

const donorTypeOptions = [
  { label: t('finance.individual'), value: 'individual' },
  { label: t('finance.organization'), value: 'organization' },
  { label: t('finance.corporate'), value: 'corporate' },
  { label: t('finance.foundation'), value: 'foundation' }
]

const statusOptions = [
  { label: t('finance.active'), value: 'active' },
  { label: t('finance.inactive'), value: 'inactive' },
  { label: t('finance.lapsed'), value: 'lapsed' }
]

const loadDonors = async () => {
  try {
    loading.value = true
    const response = await financeService.getDonors(filters)
    donors.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load donors', life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatCurrency = (amount) => amount ? `$${amount.toFixed(2)}` : '$0.00'
const getStatusVariant = (status) => ({ active: 'success', inactive: 'neutral', lapsed: 'warning' }[status] || 'neutral')

const confirmDelete = (donor) => {
  confirm.require({
    message: 'Are you sure you want to delete this donor?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await financeService.deleteDonor(donor.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('finance.donorDeleted'), life: 3000 })
        loadDonors()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete donor', life: 3000 })
      }
    }
  })
}

let searchTimeout = null
const handleSearch = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    loadDonors()
  }, 500)
}

const exportDonors = () => {
  try {
    if (donors.value.length === 0) {
      toast.add({ severity: 'warn', summary: 'Warning', detail: 'No donors to export', life: 3000 })
      return
    }

    const columns = [
      { field: 'first_name', header: 'Name' },
      { field: 'email', header: 'Email' },
      { field: 'phone', header: 'Phone' },
      { field: 'donor_type', header: 'Donor Type' },
      { field: 'total_donated', header: 'Total Donated' },
      { field: 'donation_count', header: 'Donation Count' },
      { field: 'donor_status', header: 'Status' }
    ]

    const timestamp = new Date().toISOString().split('T')[0]
    exportService.exportToCSV(donors.value, `donors_${timestamp}.csv`, columns)

    toast.add({ severity: 'success', summary: 'Success', detail: 'Donors exported successfully', life: 3000 })
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to export donors', life: 3000 })
  }
}

onMounted(loadDonors)
</script>

<style scoped>
.donor-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.header-actions { display: flex; gap: 0.75rem; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; flex-wrap: wrap; }
.filters > * { min-width: 200px; flex: 1; }
</style>
