<template>
  <div class="campaign-form-container">
    <div class="form-header">
      <Button
        icon="pi pi-arrow-left"
        class="p-button-text"
        @click="router.back()"
      />
      <h1>{{ isEdit ? 'Edit Campaign' : $t('finance.addCampaign') }}</h1>
    </div>

    <form @submit.prevent="handleSubmit">
      <Card>
        <template #title>
          {{ $t('finance.campaign') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-field full-width">
              <label for="name">{{ $t('finance.campaignName') }} *</label>
              <InputText
                id="name"
                v-model="formData.name"
                required
              />
            </div>

            <div class="form-field full-width">
              <label for="description">{{ $t('event.description') }}</label>
              <Textarea
                id="description"
                v-model="formData.description"
                rows="3"
              />
            </div>

            <div class="form-field">
              <label for="campaign_type">{{ $t('finance.campaignType') }} *</label>
              <Dropdown
                id="campaign_type"
                v-model="formData.campaign_type"
                :options="campaignTypeOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="form-field">
              <label for="status">Status</label>
              <Dropdown
                id="status"
                v-model="formData.status"
                :options="statusOptions"
                option-label="label"
                option-value="value"
              />
            </div>

            <div class="form-field">
              <label for="start_date">{{ $t('finance.startDate') }} *</label>
              <Calendar
                id="start_date"
                v-model="formData.start_date"
                date-format="yy-mm-dd"
                required
              />
            </div>

            <div class="form-field">
              <label for="end_date">{{ $t('finance.endDate') }}</label>
              <Calendar
                id="end_date"
                v-model="formData.end_date"
                date-format="yy-mm-dd"
              />
            </div>

            <div class="form-field">
              <label for="goal_amount">{{ $t('finance.goalAmount') }}</label>
              <InputNumber
                id="goal_amount"
                v-model="formData.goal_amount"
                mode="currency"
                currency="USD"
                :min="0"
              />
            </div>

            <div class="form-field">
              <label for="coordinator_name">{{ $t('finance.coordinator') }}</label>
              <InputText
                id="coordinator_name"
                v-model="formData.coordinator_name"
              />
            </div>

            <div class="form-field">
              <label for="coordinator_email">Coordinator Email</label>
              <InputText
                id="coordinator_email"
                v-model="formData.coordinator_email"
                type="email"
              />
            </div>

            <div class="form-field">
              <label for="target_audience">{{ $t('finance.targetAudience') }}</label>
              <InputText
                id="target_audience"
                v-model="formData.target_audience"
              />
            </div>

            <div class="form-field full-width">
              <label for="notes">Notes</label>
              <Textarea
                id="notes"
                v-model="formData.notes"
                rows="3"
              />
            </div>
          </div>
        </template>
      </Card>

      <div class="form-actions">
        <Button
          type="button"
          :label="$t('common.cancel')"
          class="p-button-secondary"
          @click="router.back()"
        />
        <Button
          type="submit"
          :label="$t('common.save')"
          icon="pi pi-check"
          :loading="saving"
        />
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { financeService } from '@/services/financeService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Calendar from 'primevue/calendar'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)
const saving = ref(false)

const formData = reactive({
  name: '',
  description: '',
  campaign_type: 'fundraising',
  start_date: new Date(),
  end_date: null,
  goal_amount: 0,
  currency: 'USD',
  status: 'planning',
  target_audience: '',
  coordinator_name: '',
  coordinator_email: '',
  notes: ''
})

const campaignTypeOptions = [
  { label: t('finance.fundraising'), value: 'fundraising' },
  { label: t('finance.awareness'), value: 'awareness' },
  { label: t('adoption.title'), value: 'adoption' },
  { label: t('finance.event'), value: 'event' },
  { label: 'Other', value: 'other' }
]

const statusOptions = [
  { label: t('finance.planning'), value: 'planning' },
  { label: t('finance.active'), value: 'active' },
  { label: t('finance.completed'), value: 'completed' },
  { label: t('finance.cancelled'), value: 'cancelled' }
]

const loadCampaign = async () => {
  if (!isEdit.value) return
  try {
    const campaign = await financeService.getCampaign(route.params.id)
    Object.assign(formData, {
      ...campaign,
      start_date: campaign.start_date ? new Date(campaign.start_date) : new Date(),
      end_date: campaign.end_date ? new Date(campaign.end_date) : null
    })
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load campaign', life: 3000 })
    router.push('/staff/finance/campaigns')
  }
}

const handleSubmit = async () => {
  try {
    saving.value = true
    const dataToSend = {
      ...formData,
      start_date: formData.start_date ? formData.start_date.toISOString().split('T')[0] : null,
      end_date: formData.end_date ? formData.end_date.toISOString().split('T')[0] : null
    }

    if (isEdit.value) {
      await financeService.updateCampaign(route.params.id, dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('finance.campaignUpdated'), life: 3000 })
    } else {
      await financeService.createCampaign(dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('finance.campaignCreated'), life: 3000 })
    }
    router.push('/staff/finance/campaigns')
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save campaign', life: 3000 })
  } finally {
    saving.value = false
  }
}

onMounted(loadCampaign)
</script>

<style scoped>
.campaign-form-container { max-width: 1000px; margin: 0 auto; }
.form-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
.form-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.5rem; }
.form-field label { font-weight: 600; color: #374151; }
.full-width { grid-column: 1 / -1; }
.form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 2rem; }
@media (max-width: 768px) { .form-grid { grid-template-columns: 1fr; } }
</style>
