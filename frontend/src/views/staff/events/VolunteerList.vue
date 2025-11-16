<template>
  <div class="volunteer-list">
    <div class="page-header">
      <h1>{{ $t('event.volunteers') }}</h1>
      <Button
        :label="$t('event.addVolunteer')"
        icon="pi pi-plus"
        @click="router.push('/staff/volunteers/new')"
      />
    </div>

    <Card v-if="!loading && volunteers.length > 0">
      <template #content>
        <DataTable
          :value="volunteers"
          paginator
          :rows="20"
        >
          <Column
            field="first_name"
            :header="$t('event.firstName')"
          >
            <template #body="slotProps">
              {{ slotProps.data.first_name }} {{ slotProps.data.last_name }}
            </template>
          </Column>
          <Column
            field="email"
            :header="$t('finance.email')"
          />
          <Column
            field="phone"
            :header="$t('finance.phone')"
          />
          <Column
            field="total_hours"
            :header="$t('event.totalHours')"
          />
          <Column
            field="volunteer_status"
            :header="$t('event.volunteerStatus')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.volunteer_status)">
                {{ $t(`finance.${slotProps.data.volunteer_status}`) }}
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
      v-if="!loading && volunteers.length === 0"
      :message="$t('event.noVolunteersFound')"
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

const volunteers = ref([])
const loading = ref(true)

const loadVolunteers = async () => {
  try {
    loading.value = true
    const response = await eventService.getVolunteers()
    volunteers.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load volunteers', life: 3000 })
  } finally {
    loading.value = false
  }
}

const getStatusVariant = (status) => ({ active: 'success', inactive: 'neutral', pending: 'warning' }[status] || 'neutral')

const confirmDelete = (volunteer) => {
  confirm.require({
    message: 'Are you sure you want to delete this volunteer?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await eventService.deleteVolunteer(volunteer.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('event.volunteerDeleted'), life: 3000 })
        loadVolunteers()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete volunteer', life: 3000 })
      }
    }
  })
}

onMounted(loadVolunteers)
</script>

<style scoped>
.volunteer-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
</style>
