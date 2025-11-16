<template>
  <div class="email-campaign-form">
    <div class="form-header">
      <Button
        icon="pi pi-arrow-left"
        class="p-button-text"
        @click="router.back()"
      />
      <h1>{{ isEdit ? $t('common.edit') : $t('communication.addEmailCampaign') }}</h1>
    </div>

    <Card>
      <template #content>
        <form
          class="form-grid"
          @submit.prevent="handleSubmit"
        >
          <div class="form-field full-width">
            <label for="name">{{ $t('communication.campaignName') }} *</label>
            <InputText
              id="name"
              v-model="formData.name"
              required
            />
          </div>

          <div class="form-field">
            <label for="template_id">{{ $t('communication.template') }} *</label>
            <Dropdown
              id="template_id"
              v-model="formData.template_id"
              :options="templates"
              option-label="name"
              option-value="id"
              :placeholder="$t('communication.template')"
              required
            />
          </div>

          <div class="form-field">
            <label for="recipient_type">{{ $t('communication.recipientType') }} *</label>
            <Dropdown
              id="recipient_type"
              v-model="formData.recipient_type"
              :options="recipientTypeOptions"
              option-label="label"
              option-value="value"
              required
            />
          </div>

          <div class="form-field">
            <label for="scheduled_date">{{ $t('communication.scheduledDate') }}</label>
            <Calendar
              id="scheduled_date"
              v-model="formData.scheduled_date"
              date-format="yy-mm-dd"
              show-icon
            />
          </div>

          <div class="form-field">
            <label for="status">{{ $t('common.status') }} *</label>
            <Dropdown
              id="status"
              v-model="formData.status"
              :options="statusOptions"
              option-label="label"
              option-value="value"
              required
            />
          </div>

          <div class="form-actions full-width">
            <Button
              :label="$t('common.cancel')"
              class="p-button-secondary"
              type="button"
              @click="router.back()"
            />
            <Button
              :label="$t('common.save')"
              type="submit"
            />
          </div>
        </form>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { communicationService } from '@/services/communicationService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Calendar from 'primevue/calendar'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)

const formData = reactive({
  name: '',
  template_id: null,
  recipient_type: 'all',
  scheduled_date: null,
  status: 'draft'
})

const templates = ref([])

const recipientTypeOptions = [
  { label: t('communication.donors'), value: 'donors' },
  { label: t('communication.adopters'), value: 'adopters' },
  { label: t('communication.volunteers'), value: 'volunteers' },
  { label: t('communication.all'), value: 'all' },
  { label: t('communication.custom'), value: 'custom' }
]

const statusOptions = [
  { label: t('communication.draft'), value: 'draft' },
  { label: t('communication.scheduled'), value: 'scheduled' },
  { label: t('communication.sent'), value: 'sent' }
]

const loadTemplates = async () => {
  try {
    const response = await communicationService.getEmailTemplates()
    templates.value = response.data
  } catch (error) {
    console.error('Failed to load templates:', error)
  }
}

const loadCampaign = async () => {
  if (!isEdit.value) return

  try {
    const response = await communicationService.getEmailCampaign(route.params.id)
    Object.assign(formData, {
      ...response.data,
      scheduled_date: response.data.scheduled_date ? new Date(response.data.scheduled_date) : null
    })
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load campaign', life: 3000 })
    router.push('/staff/communication/campaigns')
  }
}

const handleSubmit = async () => {
  try {
    const dataToSend = {
      ...formData,
      scheduled_date: formData.scheduled_date ? formData.scheduled_date.toISOString().split('T')[0] : null
    }

    if (isEdit.value) {
      await communicationService.updateEmailCampaign(route.params.id, dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('communication.emailCampaignUpdated'), life: 3000 })
    } else {
      await communicationService.createEmailCampaign(dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('communication.emailCampaignCreated'), life: 3000 })
    }
    router.push('/staff/communication/campaigns')
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save campaign', life: 3000 })
  }
}

onMounted(async () => {
  await loadTemplates()
  await loadCampaign()
})
</script>

<style scoped>
.email-campaign-form { max-width: 900px; margin: 0 auto; }
.form-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
.form-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.5rem; }
.form-field label { font-weight: 600; color: #374151; }
.full-width { grid-column: 1 / -1; }
.form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
</style>
