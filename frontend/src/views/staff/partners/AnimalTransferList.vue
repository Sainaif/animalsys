<template>
  <div class="animal-transfer-list">
    <div class="page-header">
      <h1>{{ $t('partner.animalTransfers') }}</h1>
      <Button
        :label="$t('partner.addTransfer')"
        icon="pi pi-plus"
        @click="router.push('/staff/partners/transfers/new')"
      />
    </div>

    <Card class="filters-card">
      <template #content>
        <div class="filters">
          <Dropdown
            v-model="filters.transfer_type"
            :options="transferTypeOptions"
            :placeholder="$t('partner.transferType')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadTransfers"
          />
          <Dropdown
            v-model="filters.status"
            :options="statusOptions"
            :placeholder="$t('common.status')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadTransfers"
          />
        </div>
      </template>
    </Card>

    <Card v-if="!loading && transfers.length > 0">
      <template #content>
        <DataTable
          :value="transfers"
          paginator
          :rows="20"
        >
          <Column
            field="animal.name"
            :header="$t('animal.name')"
          >
            <template #body="slotProps">
              {{ slotProps.data.animal?.name || 'N/A' }}
            </template>
          </Column>
          <Column
            field="transfer_type"
            :header="$t('partner.transferType')"
          >
            <template #body="slotProps">
              <Badge variant="info">
                {{ $t(`partner.${slotProps.data.transfer_type}`) }}
              </Badge>
            </template>
          </Column>
          <Column
            field="transfer_date"
            :header="$t('partner.transferDate')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.transfer_date) }}
            </template>
          </Column>
          <Column
            field="from_organization.organization_name"
            :header="$t('partner.fromOrganization')"
          >
            <template #body="slotProps">
              {{ slotProps.data.from_organization?.organization_name || 'N/A' }}
            </template>
          </Column>
          <Column
            field="to_organization.organization_name"
            :header="$t('partner.toOrganization')"
          >
            <template #body="slotProps">
              {{ slotProps.data.to_organization?.organization_name || 'N/A' }}
            </template>
          </Column>
          <Column
            field="status"
            :header="$t('common.status')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ $t(`partner.${slotProps.data.status}`) }}
              </Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  @click="router.push(`/staff/partners/transfers/${slotProps.data.id}/edit`)"
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
      v-if="!loading && transfers.length === 0"
      :message="$t('partner.noTransfersFound')"
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
import { partnerService } from '@/services/partnerService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const transfers = ref([])
const loading = ref(true)
const filters = ref({
  transfer_type: null,
  status: null
})

const transferTypeOptions = [
  { label: t('partner.incoming'), value: 'incoming' },
  { label: t('partner.outgoing'), value: 'outgoing' },
  { label: t('partner.temporary'), value: 'temporary' }
]

const statusOptions = [
  { label: t('communication.pending'), value: 'pending' },
  { label: t('partner.inTransit'), value: 'in_transit' },
  { label: t('finance.completed'), value: 'completed' },
  { label: t('finance.cancelled'), value: 'cancelled' }
]

const loadTransfers = async () => {
  try {
    loading.value = true
    const params = Object.fromEntries(
      Object.entries(filters.value).filter(([, value]) => value != null)
    )
    const response = await partnerService.getAnimalTransfers(params)
    transfers.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to load transfers', life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const getStatusVariant = (status) => ({
  pending: 'warning',
  in_transit: 'info',
  completed: 'success',
  cancelled: 'danger'
}[status] || 'neutral')

const confirmDelete = (transfer) => {
  confirm.require({
    message: t('common.deleteConfirmation'),
    header: t('common.confirm'),
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await partnerService.deleteAnimalTransfer(transfer.id)
        toast.add({ severity: 'success', summary: t('common.success'), detail: t('partner.transferDeleted'), life: 3000 })
        loadTransfers()
      } catch (error) {
        toast.add({ severity: 'error', summary: t('common.error'), detail: t('common.deleteError'), life: 3000 })
      }
    }
  })
}

onMounted(loadTransfers)
</script>

<style scoped>
.animal-transfer-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; flex-wrap: wrap; }
.filters > * { min-width: 200px; }
.action-buttons { display: flex; gap: 0.5rem; }
</style>
