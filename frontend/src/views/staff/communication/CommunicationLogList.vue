<template>
  <div class="communication-log-list">
    <div class="page-header">
      <h1>{{ $t('communication.communicationLogs') }}</h1>
      <Button
        :label="$t('communication.addCommunicationLog')"
        icon="pi pi-plus"
        @click="router.push('/staff/communication/logs/new')"
      />
    </div>

    <Card class="filters-card">
      <template #content>
        <div class="filters">
          <Dropdown
            v-model="filters.communication_type"
            :options="typeOptions"
            :placeholder="$t('communication.communicationType')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadLogs"
          />
          <Dropdown
            v-model="filters.status"
            :options="statusOptions"
            :placeholder="$t('common.status')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadLogs"
          />
        </div>
      </template>
    </Card>

    <Card v-if="!loading && logs.length > 0">
      <template #content>
        <DataTable
          :value="logs"
          paginator
          :rows="20"
        >
          <Column
            field="communication_type"
            :header="$t('communication.communicationType')"
          >
            <template #body="slotProps">
              <Badge variant="info">
                {{ $t(`communication.${slotProps.data.communication_type}`) }}
              </Badge>
            </template>
          </Column>
          <Column
            field="subject"
            :header="$t('communication.subject')"
          />
          <Column
            field="recipient_name"
            :header="$t('communication.recipient')"
          />
          <Column
            field="communication_date"
            :header="$t('communication.communicationDate')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.communication_date) }}
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
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  @click="router.push(`/staff/communication/logs/${slotProps.data.id}/edit`)"
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
      v-if="!loading && logs.length === 0"
      :message="$t('communication.noCommunicationLogsFound')"
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
import Dropdown from 'primevue/dropdown'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const logs = ref([])
const loading = ref(true)
const filters = ref({
  communication_type: null,
  status: null
})

const typeOptions = [
  { label: t('communication.email'), value: 'email' },
  { label: t('communication.phone'), value: 'phone' },
  { label: t('communication.sms'), value: 'sms' },
  { label: t('communication.inPerson'), value: 'in_person' },
  { label: t('communication.other'), value: 'other' }
]

const statusOptions = [
  { label: t('communication.sent'), value: 'sent' },
  { label: t('communication.delivered'), value: 'delivered' },
  { label: t('communication.failed'), value: 'failed' },
  { label: t('communication.pending'), value: 'pending' }
]

const loadLogs = async () => {
  try {
    loading.value = true
    const params = Object.fromEntries(
      Object.entries(filters.value).filter(([, value]) => value != null)
    )
    const response = await communicationService.getCommunicationLogs(params)
    logs.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to load communication logs', life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const getStatusVariant = (status) => ({
  sent: 'success',
  delivered: 'success',
  failed: 'danger',
  pending: 'warning'
}[status] || 'neutral')

const confirmDelete = (log) => {
  confirm.require({
    message: t('common.deleteConfirmation'),
    header: t('common.confirm'),
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await communicationService.deleteCommunicationLog(log.id)
        toast.add({ severity: 'success', summary: t('common.success'), detail: t('communication.communicationLogDeleted'), life: 3000 })
        loadLogs()
      } catch (error) {
        toast.add({ severity: 'error', summary: t('common.error'), detail: t('common.deleteError'), life: 3000 })
      }
    }
  })
}

onMounted(loadLogs)
</script>

<style scoped>
.communication-log-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; flex-wrap: wrap; }
.filters > * { min-width: 200px; }
.action-buttons { display: flex; gap: 0.5rem; }
</style>
