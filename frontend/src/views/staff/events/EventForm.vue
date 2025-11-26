<template>
  <div class="event-form">
    <div class="page-header">
      <h1>{{ isEdit ? $t('event.editEvent') : $t('event.addEvent') }}</h1>
    </div>

    <Card>
      <template #content>
        <form @submit.prevent="saveEvent">
          <div class="p-fluid formgrid grid">
            <div class="field col-12 md:col-6">
              <label for="name">{{ $t('event.eventName') }}</label>
              <InputText
                id="name"
                v-model="event.name"
                required
              />
            </div>

            <div class="field col-12 md:col-6">
              <label for="event_type">{{ $t('event.eventType') }}</label>
              <Dropdown
                id="event_type"
                v-model="event.event_type"
                :options="eventTypes"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="field col-12 md:col-6">
              <label for="start_date">{{ $t('event.startDate') }}</label>
              <Calendar
                id="start_date"
                v-model="event.start_date"
                required
              />
            </div>

            <div class="field col-12 md:col-6">
              <label for="budget">{{ $t('event.budget') }}</label>
              <InputNumber
                id="budget"
                v-model="event.budget"
              />
            </div>
          </div>

          <div class="form-actions">
            <Button
              :label="$t('common.save')"
              type="submit"
              :loading="loading"
            />
            <Button
              :label="$t('common.cancel')"
              class="p-button-secondary"
              @click="router.back()"
            />
          </div>
        </form>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { eventService } from '@/services/eventService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Calendar from 'primevue/calendar'
import InputNumber from 'primevue/inputnumber'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const event = ref({})
const loading = ref(false)
const isEdit = ref(false)

const eventTypes = [
  { label: t('event.fundraiser'), value: 'fundraiser' },
  { label: t('event.adoptionEvent'), value: 'adoption_event' },
  { label: t('event.education'), value: 'education' },
  { label: t('event.community'), value: 'community' },
  { label: t('event.shopping'), value: 'shopping' },
  { label: t('event.contribution'), value: 'contribution' },
  { label: t('event.other'), value: 'other' },
]

const loadEvent = async () => {
  if (route.params.id) {
    isEdit.value = true
    try {
      loading.value = true
      event.value = await eventService.getEvent(route.params.id)
    } catch (error) {
      toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to load event', life: 3000 })
      router.push('/staff/events')
    } finally {
      loading.value = false
    }
  }
}

const saveEvent = async () => {
  try {
    loading.value = true
    if (isEdit.value) {
      await eventService.updateEvent(event.value.id, event.value)
      toast.add({ severity: 'success', summary: t('common.success'), detail: t('event.eventUpdated'), life: 3000 })
    } else {
      await eventService.createEvent(event.value)
      toast.add({ severity: 'success', summary: t('common.success'), detail: t('event.eventCreated'), life: 3000 })
    }
    router.push('/staff/events')
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to save event', life: 3000 })
  } finally {
    loading.value = false
  }
}

onMounted(loadEvent)
</script>

<style scoped>
.event-form {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.form-actions {
  margin-top: 2rem;
  display: flex;
  gap: 1rem;
}
</style>
