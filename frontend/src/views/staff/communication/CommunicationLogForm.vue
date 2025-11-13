<template>
  <div class="communication-log-form">
    <div class="form-header">
      <Button icon="pi pi-arrow-left" class="p-button-text" @click="router.back()" />
      <h1>{{ isEdit ? $t('common.edit') : $t('communication.addCommunicationLog') }}</h1>
    </div>

    <Card>
      <template #content>
        <form @submit.prevent="handleSubmit" class="form-grid">
          <div class="form-field">
            <label for="communication_type">{{ $t('communication.communicationType') }} *</label>
            <Dropdown id="communication_type" v-model="formData.communication_type" :options="typeOptions" option-label="label" option-value="value" required />
          </div>

          <div class="form-field">
            <label for="communication_date">{{ $t('communication.communicationDate') }} *</label>
            <Calendar id="communication_date" v-model="formData.communication_date" date-format="yy-mm-dd" show-icon required />
          </div>

          <div class="form-field full-width">
            <label for="subject">{{ $t('communication.subject') }} *</label>
            <InputText id="subject" v-model="formData.subject" required />
          </div>

          <div class="form-field">
            <label for="recipient_type">{{ $t('communication.recipientType') }} *</label>
            <Dropdown id="recipient_type" v-model="formData.recipient_type" :options="recipientTypeOptions" option-label="label" option-value="value" required />
          </div>

          <div class="form-field">
            <label for="recipient_name">{{ $t('communication.recipient') }}</label>
            <InputText id="recipient_name" v-model="formData.recipient_name" />
          </div>

          <div class="form-field">
            <label for="sender_name">{{ $t('communication.sender') }}</label>
            <InputText id="sender_name" v-model="formData.sender_name" />
          </div>

          <div class="form-field">
            <label for="status">{{ $t('common.status') }} *</label>
            <Dropdown id="status" v-model="formData.status" :options="statusOptions" option-label="label" option-value="value" required />
          </div>

          <div class="form-field full-width">
            <label for="message">{{ $t('communication.message') }}</label>
            <Textarea id="message" v-model="formData.message" rows="5" />
          </div>

          <div class="form-field full-width">
            <label for="notes">{{ $t('communication.notes') }}</label>
            <Textarea id="notes" v-model="formData.notes" rows="3" />
          </div>

          <div class="form-actions full-width">
            <Button :label="$t('common.cancel')" class="p-button-secondary" @click="router.back()" type="button" />
            <Button :label="$t('common.save')" type="submit" />
          </div>
        </form>
      </template>
    </Card>
  </div>
</template>

<script setup>
import { reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { communicationService } from '@/services/communicationService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Calendar from 'primevue/calendar'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)

const formData = reactive({
  communication_type: 'email',
  subject: '',
  message: '',
  recipient_type: 'contact',
  recipient_name: '',
  sender_name: '',
  communication_date: new Date(),
  status: 'sent',
  notes: ''
})

const typeOptions = [
  { label: t('communication.email'), value: 'email' },
  { label: t('communication.phone'), value: 'phone' },
  { label: t('communication.sms'), value: 'sms' },
  { label: t('communication.inPerson'), value: 'in_person' },
  { label: t('communication.other'), value: 'other' }
]

const recipientTypeOptions = [
  { label: t('communication.donors'), value: 'donor' },
  { label: t('communication.adopters'), value: 'adopter' },
  { label: t('communication.volunteers'), value: 'volunteer' },
  { label: 'Contact', value: 'contact' },
  { label: 'Staff', value: 'staff' }
]

const statusOptions = [
  { label: t('communication.sent'), value: 'sent' },
  { label: t('communication.delivered'), value: 'delivered' },
  { label: t('communication.failed'), value: 'failed' },
  { label: t('communication.pending'), value: 'pending' }
]

const loadLog = async () => {
  if (!isEdit.value) return

  try {
    const response = await communicationService.getCommunicationLog(route.params.id)
    Object.assign(formData, {
      ...response.data,
      communication_date: response.data.communication_date ? new Date(response.data.communication_date) : new Date()
    })
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load log', life: 3000 })
    router.push('/staff/communication/logs')
  }
}

const handleSubmit = async () => {
  try {
    const dataToSend = {
      ...formData,
      communication_date: formData.communication_date ? formData.communication_date.toISOString().split('T')[0] : null
    }

    if (isEdit.value) {
      await communicationService.updateCommunicationLog(route.params.id, dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('communication.communicationLogUpdated'), life: 3000 })
    } else {
      await communicationService.createCommunicationLog(dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('communication.communicationLogCreated'), life: 3000 })
    }
    router.push('/staff/communication/logs')
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save log', life: 3000 })
  }
}

onMounted(loadLog)
</script>

<style scoped>
.communication-log-form { max-width: 900px; margin: 0 auto; }
.form-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
.form-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.5rem; }
.form-field label { font-weight: 600; color: #374151; }
.full-width { grid-column: 1 / -1; }
.form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
</style>
