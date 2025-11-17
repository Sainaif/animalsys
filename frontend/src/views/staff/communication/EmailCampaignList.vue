<template>
  <div class="email-campaign-list">
    <div class="page-header">
      <h1>{{ $t('communication.emailCampaigns') }}</h1>
      <Button
        :label="$t('communication.addEmailCampaign')"
        icon="pi pi-plus"
        @click="router.push('/staff/communication/campaigns/new')"
      />
    </div>

    <Card v-if="!loading && campaigns.length > 0">
      <template #content>
        <DataTable
          :value="campaigns"
          paginator
          :rows="20"
        >
          <Column
            field="name"
            :header="$t('communication.campaignName')"
          />
          <Column
            field="recipient_type"
            :header="$t('communication.recipientType')"
          >
            <template #body="slotProps">
              {{ $t(`communication.${slotProps.data.recipient_type}`) }}
            </template>
          </Column>
          <Column
            field="scheduled_date"
            :header="$t('communication.scheduledDate')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.scheduled_date) }}
            </template>
          </Column>
          <Column
            field="status"
            :header="$t('common.status')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ $t(`communication.${slotProps.data.status}`) }}
              </Badge>
            </template>
          </Column>
          <Column
            field="sent_count"
            :header="$t('communication.sentCount')"
          />
          <Column
            field="opened_count"
            :header="$t('communication.openedCount')"
          />
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  v-if="slotProps.data.status === 'draft'"
                  icon="pi pi-send"
                  class="p-button-rounded p-button-text p-button-success"
                  @click="confirmSend(slotProps.data)"
                />
                <Button
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  @click="router.push(`/staff/communication/campaigns/${slotProps.data.id}/edit`)"
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
    <EmptyState
      v-if="!loading && campaigns.length === 0"
      :message="$t('communication.noEmailCampaignsFound')"
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
import { communicationService } from '@/services/communicationService'
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

const campaigns = ref([])
const loading = ref(true)

const loadCampaigns = async () => {
  try {
    loading.value = true
    const response = await communicationService.getEmailCampaigns()
    campaigns.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to load email campaigns', life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const getStatusVariant = (status) => ({
  draft: 'neutral',
  scheduled: 'warning',
  sending: 'info',
  sent: 'success',
  failed: 'danger'
}[status] || 'neutral')

const confirmSend = (campaign) => {
  confirm.require({
    message: `Are you sure you want to send campaign "${campaign.name}"?`,
    header: t('common.confirm'),
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await communicationService.sendEmailCampaign(campaign.id)
        toast.add({ severity: 'success', summary: t('common.success'), detail: t('communication.emailCampaignSent'), life: 3000 })
        loadCampaigns()
      } catch (error) {
        toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to send campaign', life: 3000 })
      }
    }
  })
}

const confirmDelete = (campaign) => {
  confirm.require({
    message: t('common.deleteConfirmation'),
    header: t('common.confirm'),
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await communicationService.deleteEmailCampaign(campaign.id)
        toast.add({ severity: 'success', summary: t('common.success'), detail: t('communication.emailCampaignDeleted'), life: 3000 })
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
.email-campaign-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
.action-buttons { display: flex; gap: 0.5rem; }
</style>
