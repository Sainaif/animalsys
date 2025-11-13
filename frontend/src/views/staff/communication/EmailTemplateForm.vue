<template>
  <div class="email-template-form">
    <div class="form-header">
      <Button icon="pi pi-arrow-left" class="p-button-text" @click="router.back()" />
      <h1>{{ isEdit ? $t('common.edit') : $t('communication.addEmailTemplate') }}</h1>
    </div>

    <Card>
      <template #content>
        <form @submit.prevent="handleSubmit" class="form-grid">
          <div class="form-field full-width">
            <label for="name">{{ $t('communication.templateName') }} *</label>
            <InputText id="name" v-model="formData.name" required />
          </div>

          <div class="form-field">
            <label for="template_type">{{ $t('communication.templateType') }} *</label>
            <Dropdown id="template_type" v-model="formData.template_type" :options="templateTypeOptions" option-label="label" option-value="value" required />
          </div>

          <div class="form-field">
            <label for="is_active">{{ $t('communication.isActive') }}</label>
            <Checkbox id="is_active" v-model="formData.is_active" :binary="true" />
          </div>

          <div class="form-field full-width">
            <label for="subject">{{ $t('communication.subject') }} *</label>
            <InputText id="subject" v-model="formData.subject" required />
          </div>

          <div class="form-field full-width">
            <label for="body">{{ $t('communication.body') }} *</label>
            <Textarea id="body" v-model="formData.body" rows="10" required />
          </div>

          <div class="form-field full-width">
            <label for="variables">{{ $t('communication.variables') }}</label>
            <Chips id="variables" v-model="formData.variables" separator="," />
            <small>Available variables: {name}, {email}, {phone}, {date}</small>
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
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { communicationService } from '@/services/communicationService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Checkbox from 'primevue/checkbox'
import Chips from 'primevue/chips'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)

const formData = reactive({
  name: '',
  subject: '',
  body: '',
  template_type: 'general',
  variables: [],
  is_active: true
})

const templateTypeOptions = [
  { label: t('communication.adoption'), value: 'adoption' },
  { label: t('communication.donation'), value: 'donation' },
  { label: t('communication.event'), value: 'event' },
  { label: t('communication.general'), value: 'general' },
  { label: t('communication.newsletter'), value: 'newsletter' }
]

const loadTemplate = async () => {
  if (!isEdit.value) return

  try {
    const response = await communicationService.getEmailTemplate(route.params.id)
    Object.assign(formData, response.data)
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load template', life: 3000 })
    router.push('/staff/communication/templates')
  }
}

const handleSubmit = async () => {
  try {
    if (isEdit.value) {
      await communicationService.updateEmailTemplate(route.params.id, formData)
      toast.add({ severity: 'success', summary: 'Success', detail: t('communication.emailTemplateUpdated'), life: 3000 })
    } else {
      await communicationService.createEmailTemplate(formData)
      toast.add({ severity: 'success', summary: 'Success', detail: t('communication.emailTemplateCreated'), life: 3000 })
    }
    router.push('/staff/communication/templates')
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save template', life: 3000 })
  }
}

onMounted(loadTemplate)
</script>

<style scoped>
.email-template-form { max-width: 900px; margin: 0 auto; }
.form-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
.form-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.5rem; }
.form-field label { font-weight: 600; color: #374151; }
.full-width { grid-column: 1 / -1; }
.form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
</style>
