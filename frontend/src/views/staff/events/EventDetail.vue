<template>
  <div class="event-detail">
    <LoadingSpinner v-if="loading" />

    <div v-else-if="event" class="detail-container">
      <div class="detail-header">
        <Button icon="pi pi-arrow-left" class="p-button-text" @click="router.back()" />
        <h1>{{ event.name }}</h1>
        <div class="header-actions">
          <Button :label="$t('common.edit')" icon="pi pi-pencil" @click="router.push(`/staff/events/${event.id}/edit`)" />
          <Button :label="$t('common.delete')" icon="pi pi-trash" class="p-button-danger" @click="confirmDelete" />
        </div>
      </div>

      <Card class="status-card">
        <template #content>
          <div class="status-info">
            <div class="status-item">
              <label>{{ $t('common.status') }}</label>
              <Badge :variant="getStatusVariant(event.status)">{{ $t(`event.${event.status}`) }}</Badge>
            </div>
            <div class="status-item">
              <label>{{ $t('event.eventType') }}</label>
              <p>{{ $t(`event.${event.event_type}`) }}</p>
            </div>
            <div class="status-item">
              <label>{{ $t('event.startDate') }}</label>
              <p>{{ formatDate(event.start_date) }}</p>
            </div>
          </div>
        </template>
      </Card>

      <TabView>
        <TabPanel header="Details">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('event.eventName') }}</label>
                  <p>{{ event.name }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.location') }}</label>
                  <p>{{ event.location || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.startDate') }}</label>
                  <p>{{ formatDate(event.start_date) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.endDate') }}</label>
                  <p>{{ formatDate(event.end_date) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.organizer') }}</label>
                  <p>{{ event.organizer_name || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>Organizer Email</label>
                  <p>{{ event.organizer_email || 'N/A' }}</p>
                </div>
                <div class="info-item full-width">
                  <label>{{ $t('event.description') }}</label>
                  <p>{{ event.description || 'N/A' }}</p>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel header="Participants">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('event.maxParticipants') }}</label>
                  <p>{{ event.max_participants || 'Unlimited' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.registeredParticipants') }}</label>
                  <p>{{ event.registered_participants || 0 }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.volunteersNeeded') }}</label>
                  <p>{{ event.volunteers_needed || 0 }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.volunteersAssigned') }}</label>
                  <p>{{ event.volunteers_assigned || 0 }}</p>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel header="Budget" v-if="event.budget || event.raised_amount">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('event.budget') }}</label>
                  <p>{{ formatCurrency(event.budget) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.raisedAmount') }}</label>
                  <p>{{ formatCurrency(event.raised_amount) }}</p>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>
      </TabView>
    </div>

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { eventService } from '@/services/eventService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import TabView from 'primevue/tabview'
import TabPanel from 'primevue/tabpanel'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const event = ref(null)
const loading = ref(true)

const loadEvent = async () => {
  try {
    loading.value = true
    event.value = await eventService.getEvent(route.params.id)
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load event', life: 3000 })
    router.push('/staff/events')
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const formatCurrency = (amount) => amount ? `$${amount.toFixed(2)}` : '$0.00'
const getStatusVariant = (status) => ({ planned: 'neutral', active: 'success', completed: 'info', cancelled: 'danger' }[status] || 'neutral')

const confirmDelete = () => {
  confirm.require({
    message: 'Are you sure you want to delete this event?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await eventService.deleteEvent(event.value.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('event.eventDeleted'), life: 3000 })
        router.push('/staff/events')
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete event', life: 3000 })
      }
    }
  })
}

onMounted(loadEvent)
</script>

<style scoped>
.event-detail { max-width: 1200px; margin: 0 auto; }
.detail-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
.detail-header h1 { flex: 1; font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.header-actions { display: flex; gap: 0.5rem; }
.status-card { margin-bottom: 1.5rem; }
.status-info { display: flex; gap: 2rem; align-items: center; }
.status-item { display: flex; flex-direction: column; gap: 0.5rem; }
.status-item label { font-weight: 600; color: #6b7280; font-size: 0.875rem; text-transform: uppercase; }
.status-item p { color: #2c3e50; font-size: 1rem; margin: 0; }
.info-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1.5rem; }
.info-item { display: flex; flex-direction: column; gap: 0.5rem; }
.info-item label { font-weight: 600; color: #6b7280; font-size: 0.875rem; text-transform: uppercase; }
.info-item p { color: #2c3e50; font-size: 1rem; margin: 0; line-height: 1.6; }
.full-width { grid-column: 1 / -1; }
</style>
