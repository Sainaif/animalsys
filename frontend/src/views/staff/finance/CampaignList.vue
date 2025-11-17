<template>
  <div class="campaign-list">
    <div class="page-header">
      <h1>{{ $t('finance.campaigns') }}</h1>
      <Button
        :label="$t('finance.addCampaign')"
        icon="pi pi-plus"
        @click="router.push('/staff/finance/campaigns/new')"
      />
    </div>

    <Card v-if="!loading && campaigns.length > 0">
      <template #content>
        <DataTable
          :value="campaigns"
          paginator
          :rows="20"
        >
          <Column :header="$t('finance.campaignName')">
            <template #body="slotProps">
              {{ getCampaignName(slotProps.data) }}
            </template>
          </Column>
          <Column :header="$t('finance.campaignType')">
            <template #body="slotProps">
              {{ translateFinanceKey(slotProps.data.campaign_type || slotProps.data.type) }}
            </template>
          </Column>
          <Column
            field="start_date"
            :header="$t('finance.startDate')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.start_date) }}
            </template>
          </Column>
          <Column :header="$t('finance.goalAmount')">
            <template #body="slotProps">
              {{ formatCurrency(getGoalAmount(slotProps.data)) }}
            </template>
          </Column>
          <Column :header="$t('finance.raisedAmount')">
            <template #body="slotProps">
              {{ formatCurrency(getRaisedAmount(slotProps.data)) }}
            </template>
          </Column>
          <Column
            field="status"
            :header="$t('common.status')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ translateFinanceKey(slotProps.data.status) }}
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
      v-if="!loading && campaigns.length === 0"
      :message="$t('finance.noCampaignsFound')"
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
import { getLocalizedValue } from '@/utils/animalHelpers'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const { t, locale } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const campaigns = ref([])
const loading = ref(true)

const loadCampaigns = async () => {
  try {
    loading.value = true
    const response = await financeService.getCampaigns()
    campaigns.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: t('finance.loadCampaignsError'), life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const formatCurrency = (amount) => (typeof amount === 'number' ? `$${amount.toFixed(2)}` : '$0.00')
const getStatusVariant = (status) => ({ planning: 'neutral', active: 'success', completed: 'info', cancelled: 'danger', draft: 'neutral', paused: 'warning' }[status] || 'neutral')
const translateFinanceKey = (value) => {
  if (!value) return t('common.notAvailable')
  const key = `finance.${value}`
  const translated = t(key)
  return translated === key ? value : translated
}
const getGoalAmount = (campaign) => {
  if (typeof campaign?.goal_amount === 'number') {
    return campaign.goal_amount
  }
  if (typeof campaign?.goal === 'number') {
    return campaign.goal
  }
  return 0
}
const getRaisedAmount = (campaign) => {
  if (typeof campaign?.raised_amount === 'number') {
    return campaign.raised_amount
  }
  if (typeof campaign?.current_amount === 'number') {
    return campaign.current_amount
  }
  if (typeof campaign?.currentAmount === 'number') {
    return campaign.currentAmount
  }
  return 0
}
const getCampaignName = (campaign) => getLocalizedValue(campaign?.name, locale.value) || t('common.notAvailable')

const confirmDelete = (campaign) => {
  confirm.require({
    message: t('common.deleteConfirmation'),
    header: t('common.confirm'),
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await financeService.deleteCampaign(campaign.id)
        toast.add({ severity: 'success', summary: t('common.success'), detail: t('finance.campaignDeleted'), life: 3000 })
        loadCampaigns()
      } catch (error) {
        toast.add({ severity: 'error', summary: t('common.error'), detail: t('common.deleteError'), life: 3000 })
      }
    }
  })
}

onMounted(loadCampaigns)
</script>

<style scoped>
.campaign-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
</style>
