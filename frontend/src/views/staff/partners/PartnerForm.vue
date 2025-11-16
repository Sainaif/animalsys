<template>
  <div class="partner-form">
    <div class="form-header">
      <Button
        icon="pi pi-arrow-left"
        class="p-button-text"
        @click="router.back()"
      />
      <h1>{{ isEdit ? $t('common.edit') : $t('partner.addPartner') }}</h1>
    </div>

    <Card>
      <template #content>
        <form
          class="form-grid"
          @submit.prevent="handleSubmit"
        >
          <div class="form-field full-width">
            <label for="organization_name">{{ $t('partner.organizationName') }} *</label>
            <InputText
              id="organization_name"
              v-model="formData.organization_name"
              required
            />
          </div>

          <div class="form-field">
            <label for="partner_type">{{ $t('partner.partnerType') }} *</label>
            <Dropdown
              id="partner_type"
              v-model="formData.partner_type"
              :options="partnerTypeOptions"
              option-label="label"
              option-value="value"
              required
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

          <div class="form-field">
            <label for="contact_person">{{ $t('partner.contactPerson') }} *</label>
            <InputText
              id="contact_person"
              v-model="formData.contact_person"
              required
            />
          </div>

          <div class="form-field">
            <label for="email">Email *</label>
            <InputText
              id="email"
              v-model="formData.email"
              type="email"
              required
            />
          </div>

          <div class="form-field">
            <label for="phone">{{ $t('finance.phone') }}</label>
            <InputText
              id="phone"
              v-model="formData.phone"
            />
          </div>

          <div class="form-field">
            <label for="website">{{ $t('partner.website') }}</label>
            <InputText
              id="website"
              v-model="formData.website"
            />
          </div>

          <div class="form-field">
            <label for="license_number">{{ $t('partner.licenseNumber') }}</label>
            <InputText
              id="license_number"
              v-model="formData.license_number"
            />
          </div>

          <div class="form-field">
            <label for="capacity">{{ $t('partner.capacity') }}</label>
            <InputNumber
              id="capacity"
              v-model="formData.capacity"
            />
          </div>

          <div class="form-field full-width">
            <label for="address">{{ $t('finance.address') }}</label>
            <Textarea
              id="address"
              v-model="addressString"
              rows="3"
              placeholder="Street, City, State, Postal Code, Country"
            />
          </div>

          <div class="form-field full-width">
            <label for="notes">{{ $t('communication.notes') }}</label>
            <Textarea
              id="notes"
              v-model="formData.notes"
              rows="3"
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
import { partnerService } from '@/services/partnerService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)

const formData = reactive({
  organization_name: '',
  partner_type: 'shelter',
  contact_person: '',
  email: '',
  phone: '',
  address: null,
  website: '',
  license_number: '',
  capacity: null,
  status: 'active',
  notes: ''
})

const addressString = ref('')

const partnerTypeOptions = [
  { label: t('partner.shelter'), value: 'shelter' },
  { label: t('partner.rescue'), value: 'rescue' },
  { label: t('partner.veterinary'), value: 'veterinary' },
  { label: t('partner.foster'), value: 'foster' },
  { label: t('communication.other'), value: 'other' }
]

const statusOptions = [
  { label: t('finance.active'), value: 'active' },
  { label: t('finance.inactive'), value: 'inactive' },
  { label: t('communication.pending'), value: 'pending' }
]

const loadPartner = async () => {
  if (!isEdit.value) return

  try {
    const response = await partnerService.getPartner(route.params.id)
    Object.assign(formData, response.data)
    if (response.data.address) {
      addressString.value = `${response.data.address.street}, ${response.data.address.city}, ${response.data.address.state}, ${response.data.address.postal_code}, ${response.data.address.country}`
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load partner', life: 3000 })
    router.push('/staff/partners')
  }
}

const handleSubmit = async () => {
  try {
    const dataToSend = { ...formData }
    if (addressString.value) {
      const parts = addressString.value.split(',').map(s => s.trim())
      dataToSend.address = {
        street: parts[0] || '',
        city: parts[1] || '',
        state: parts[2] || '',
        postal_code: parts[3] || '',
        country: parts[4] || ''
      }
    }

    if (isEdit.value) {
      await partnerService.updatePartner(route.params.id, dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('partner.partnerUpdated'), life: 3000 })
    } else {
      await partnerService.createPartner(dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('partner.partnerCreated'), life: 3000 })
    }
    router.push('/staff/partners')
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save partner', life: 3000 })
  }
}

onMounted(loadPartner)
</script>

<style scoped>
.partner-form { max-width: 900px; margin: 0 auto; }
.form-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
.form-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.5rem; }
.form-field label { font-weight: 600; color: #374151; }
.full-width { grid-column: 1 / -1; }
.form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
</style>
