<template>
  <div class="volunteer-detail">
    <LoadingSpinner v-if="loading" />

    <div v-else-if="volunteer" class="detail-container">
      <div class="detail-header">
        <Button icon="pi pi-arrow-left" class="p-button-text" @click="router.back()" />
        <h1>{{ volunteer.first_name }} {{ volunteer.last_name }}</h1>
        <div class="header-actions">
          <Button :label="$t('common.edit')" icon="pi pi-pencil" @click="router.push(`/staff/volunteers/${volunteer.id}/edit`)" />
          <Button :label="$t('common.delete')" icon="pi pi-trash" class="p-button-danger" @click="confirmDelete" />
        </div>
      </div>

      <Card class="status-card">
        <template #content>
          <div class="status-info">
            <div class="status-item">
              <label>{{ $t('event.volunteerStatus') }}</label>
              <Badge :variant="getStatusVariant(volunteer.volunteer_status)">{{ $t(`finance.${volunteer.volunteer_status}`) }}</Badge>
            </div>
            <div class="status-item">
              <label>{{ $t('event.totalHours') }}</label>
              <p>{{ volunteer.total_hours || 0 }} hours</p>
            </div>
            <div class="status-item">
              <label>{{ $t('event.orientationCompleted') }}</label>
              <Badge :variant="volunteer.orientation_completed ? 'success' : 'warning'">
                {{ volunteer.orientation_completed ? $t('common.yes') : $t('common.no') }}
              </Badge>
            </div>
          </div>
        </template>
      </Card>

      <TabView>
        <TabPanel header="Personal Info">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('event.firstName') }}</label>
                  <p>{{ volunteer.first_name }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.lastName') }}</label>
                  <p>{{ volunteer.last_name }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('finance.email') }}</label>
                  <p>{{ volunteer.email }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('finance.phone') }}</label>
                  <p>{{ volunteer.phone || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>Date of Birth</label>
                  <p>{{ formatDate(volunteer.date_of_birth) }}</p>
                </div>
                <div class="info-item">
                  <label>Start Date</label>
                  <p>{{ formatDate(volunteer.start_date) }}</p>
                </div>
                <div class="info-item full-width" v-if="volunteer.address">
                  <label>{{ $t('finance.address') }}</label>
                  <p>
                    {{ volunteer.address.street }}<br>
                    {{ volunteer.address.city }}, {{ volunteer.address.state }} {{ volunteer.address.postal_code }}<br>
                    {{ volunteer.address.country }}
                  </p>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel header="Emergency Contact">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('event.emergencyContact') }} Name</label>
                  <p>{{ volunteer.emergency_contact_name || 'N/A' }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.emergencyContact') }} Phone</label>
                  <p>{{ volunteer.emergency_contact_phone || 'N/A' }}</p>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel header="Skills & Availability">
          <Card>
            <template #content>
              <div class="info-item">
                <label>{{ $t('event.skills') }}</label>
                <div class="skills-list" v-if="volunteer.skills && volunteer.skills.length">
                  <Badge v-for="skill in volunteer.skills" :key="skill" variant="info">{{ skill }}</Badge>
                </div>
                <p v-else>N/A</p>
              </div>

              <div class="info-item" style="margin-top: 1.5rem;">
                <label>{{ $t('event.preferredRoles') }}</label>
                <div class="skills-list" v-if="volunteer.preferred_roles && volunteer.preferred_roles.length">
                  <Badge v-for="role in volunteer.preferred_roles" :key="role" variant="success">{{ role }}</Badge>
                </div>
                <p v-else>N/A</p>
              </div>

              <div class="info-item" style="margin-top: 1.5rem;">
                <label>{{ $t('event.availability') }}</label>
                <div class="availability-grid" v-if="volunteer.availability">
                  <div v-for="(available, day) in volunteer.availability" :key="day" class="availability-item">
                    <Checkbox :model-value="available" :binary="true" disabled />
                    <label>{{ day }}</label>
                  </div>
                </div>
              </div>
            </template>
          </Card>
        </TabPanel>

        <TabPanel header="Background Check">
          <Card>
            <template #content>
              <div class="info-grid">
                <div class="info-item">
                  <label>{{ $t('event.backgroundCheck') }} Status</label>
                  <Badge :variant="getBackgroundCheckVariant(volunteer.background_check_status)">
                    {{ volunteer.background_check_status ? $t(`event.${volunteer.background_check_status}`) : 'N/A' }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>Check Date</label>
                  <p>{{ formatDate(volunteer.background_check_date) }}</p>
                </div>
                <div class="info-item">
                  <label>{{ $t('event.orientationCompleted') }}</label>
                  <Badge :variant="volunteer.orientation_completed ? 'success' : 'warning'">
                    {{ volunteer.orientation_completed ? $t('common.yes') : $t('common.no') }}
                  </Badge>
                </div>
                <div class="info-item">
                  <label>Orientation Date</label>
                  <p>{{ formatDate(volunteer.orientation_date) }}</p>
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
import Checkbox from 'primevue/checkbox'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const volunteer = ref(null)
const loading = ref(true)

const loadVolunteer = async () => {
  try {
    loading.value = true
    volunteer.value = await eventService.getVolunteer(route.params.id)
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load volunteer', life: 3000 })
    router.push('/staff/volunteers')
  } finally {
    loading.value = false
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const getStatusVariant = (status) => ({ active: 'success', inactive: 'neutral', pending: 'warning' }[status] || 'neutral')
const getBackgroundCheckVariant = (status) => ({ pending: 'warning', approved: 'success', rejected: 'danger' }[status] || 'neutral')

const confirmDelete = () => {
  confirm.require({
    message: 'Are you sure you want to delete this volunteer?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await eventService.deleteVolunteer(volunteer.value.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('event.volunteerDeleted'), life: 3000 })
        router.push('/staff/volunteers')
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete volunteer', life: 3000 })
      }
    }
  })
}

onMounted(loadVolunteer)
</script>

<style scoped>
.volunteer-detail { max-width: 1200px; margin: 0 auto; }
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
.skills-list { display: flex; flex-wrap: wrap; gap: 0.5rem; }
.availability-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(150px, 1fr)); gap: 1rem; }
.availability-item { display: flex; align-items: center; gap: 0.5rem; }
.availability-item label { text-transform: capitalize; }
</style>
