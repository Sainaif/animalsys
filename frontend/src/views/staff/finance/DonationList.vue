<template>
  <div class="donation-list">
    <div class="page-header">
      <h1>{{ $t('finance.donations') }}</h1>
      <Button
        :label="$t('finance.addDonation')"
        icon="pi pi-plus"
        @click="router.push('/staff/finance/donations/new')"
      />
    </div>

    <Card v-if="!loading && donations.length > 0">
      <template #content>
        <DataTable
          :value="donations"
          paginator
          :rows="20"
        >
          <Column
            field="donor.first_name"
            header="Donor"
          >
            <template #body="slotProps">
              {{ slotProps.data.donor?.organization_name || `${slotProps.data.donor?.first_name || ''} ${slotProps.data.donor?.last_name || ''}` }}
            </template>
          </Column>
          <Column
            field="amount"
            :header="$t('finance.amount')"
          >
            <template #body="slotProps">
              {{ formatCurrency(slotProps.data.amount) }}
            </template>
          </Column>
          <Column
            field="donation_date"
            :header="$t('finance.donationDate')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.donation_date) }}
            </template>
          </Column>
          <Column
            field="donation_type"
            :header="$t('finance.donationType')"
          >
            <template #body="slotProps">
              {{ translateFinanceKey(slotProps.data.donation_type) }}
            </template>
          </Column>
          <Column
            field="payment_status"
            :header="$t('finance.paymentStatus')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.payment_status)">
                {{ translateFinanceKey(slotProps.data.payment_status) }}
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
      v-if="!loading && donations.length === 0"
      :message="$t('finance.noDonationsFound')"
    />
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { financeService } from '@/services/financeService'
import Card from 'primevue/card'
import Button from 'primevue/button'
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

const donations = ref([])
const loading = ref(true)

const loadDonations = async () => {
  try {
    loading.value = true
    const response = await financeService.getDonations()
    donations.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: t('finance.loadDonationsError'), life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const formatCurrency = (amount) => (typeof amount === 'number' ? `$${amount.toFixed(2)}` : '$0.00')
const getStatusVariant = (status) => ({ pending: 'warning', completed: 'success', failed: 'danger', refunded: 'neutral' }[status] || 'neutral')
const translateFinanceKey = (value) => {
  if (!value) return t('common.notAvailable')
  const key = `finance.${value}`
  const translated = t(key)
  return translated === key ? value : translated
}

const confirmDelete = (donation) => {
  confirm.require({
    message: t('common.deleteConfirmation'),
    header: t('common.confirm'),
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await financeService.deleteDonation(donation.id)
        toast.add({ severity: 'success', summary: t('common.success'), detail: t('finance.donationDeleted'), life: 3000 })
        loadDonations()
      } catch (error) {
        toast.add({ severity: 'error', summary: t('common.error'), detail: t('common.deleteError'), life: 3000 })
      }
    }
  })
}

onMounted(loadDonations)
</script>

<style scoped>
.donation-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
</style>
