<template>
  <div class="event-list">
    <div class="page-header">
      <h1>{{ $t('event.title') }}</h1>
      <Button :label="$t('event.addEvent')" icon="pi pi-plus" @click="router.push('/staff/events/new')" />
    </div>

    <Card v-if="!loading && events.length > 0">
      <template #content>
        <DataTable :value="events" paginator :rows="20">
          <Column field="name" :header="$t('event.eventName')" />
          <Column field="event_type" :header="$t('event.eventType')">
            <template #body="slotProps">{{ $t(`event.${slotProps.data.event_type}`) }}</template>
          </Column>
          <Column field="start_date" :header="$t('event.startDate')">
            <template #body="slotProps">{{ formatDate(slotProps.data.start_date) }}</template>
          </Column>
          <Column field="location" :header="$t('event.location')" />
          <Column field="registered_participants" :header="$t('event.registeredParticipants')" />
          <Column field="status" :header="$t('common.status')">
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">{{ $t(`event.${slotProps.data.status}`) }}</Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button icon="pi pi-eye" class="p-button-rounded p-button-text" @click="router.push(`/staff/events/${slotProps.data.id}`)" />
                <Button icon="pi pi-trash" class="p-button-rounded p-button-text p-button-danger" @click="confirmDelete(slotProps.data)" />
              </div>
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState v-if="!loading && events.length === 0" :message="$t('event.noEventsFound')" />
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { eventService } from '@/services/eventService'
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

const events = ref([])
const loading = ref(true)

const loadEvents = async () => {
  try {
    loading.value = true
    const response = await eventService.getEvents()
    events.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load events', life: 3000 })
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const getStatusVariant = (status) => ({ planned: 'neutral', active: 'success', completed: 'info', cancelled: 'danger' }[status] || 'neutral')

const confirmDelete = (event) => {
  confirm.require({
    message: 'Are you sure you want to delete this event?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await eventService.deleteEvent(event.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('event.eventDeleted'), life: 3000 })
        loadEvents()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete event', life: 3000 })
      }
    }
  })
}

onMounted(loadEvents)
</script>

<style scoped>
.event-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
</style>
