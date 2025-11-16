<template>
  <div class="donor-form-container">
    <div class="form-header">
      <Button
        icon="pi pi-arrow-left"
        class="p-button-text"
        @click="router.back()"
      />
      <h1>{{ isEdit ? 'Edit Donor' : $t('finance.addDonor') }}</h1>
    </div>

    <form @submit.prevent="handleSubmit">
      <Card>
        <template #title>
          {{ $t('finance.donor') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-field full-width">
              <label for="donor_type">{{ $t('finance.donorType') }} *</label>
              <Dropdown
                id="donor_type"
                v-model="formData.donor_type"
                :options="donorTypeOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div
              v-if="formData.donor_type === 'individual'"
              class="form-field"
            >
              <label for="first_name">{{ $t('finance.firstName') }} *</label>
              <InputText
                id="first_name"
                v-model="formData.first_name"
                :required="formData.donor_type === 'individual'"
              />
            </div>

            <div
              v-if="formData.donor_type === 'individual'"
              class="form-field"
            >
              <label for="last_name">{{ $t('finance.lastName') }} *</label>
              <InputText
                id="last_name"
                v-model="formData.last_name"
                :required="formData.donor_type === 'individual'"
              />
            </div>

            <div
              v-if="formData.donor_type !== 'individual'"
              class="form-field"
              :class="{ 'full-width': formData.donor_type !== 'individual' }"
            >
              <label for="organization_name">{{ $t('finance.organizationName') }} *</label>
              <InputText
                id="organization_name"
                v-model="formData.organization_name"
                :required="formData.donor_type !== 'individual'"
              />
            </div>

            <div class="form-field">
              <label for="email">{{ $t('finance.email') }} *</label>
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

            <div class="form-field full-width">
              <label for="address">{{ $t('finance.address') }}</label>
              <Textarea
                id="address"
                v-model="addressString"
                rows="3"
                placeholder="Street, City, State, Postal Code, Country"
              />
            </div>

            <div class="form-field">
              <label for="donor_status">{{ $t('finance.donorStatus') }}</label>
              <Dropdown
                id="donor_status"
                v-model="formData.donor_status"
                :options="statusOptions"
                option-label="label"
                option-value="value"
              />
            </div>

            <div class="form-field full-width">
              <label for="notes">{{ $t('common.notes') }}</label>
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
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)
const saving = ref(false)
const addressString = ref('')

const formData = reactive({
  donor_type: 'individual',
  first_name: '',
  last_name: '',
  organization_name: '',
  email: '',
  phone: '',
  address: null,
  donor_status: 'active',
  notes: ''
})

const donorTypeOptions = [
  { label: t('finance.individual'), value: 'individual' },
  { label: t('finance.organization'), value: 'organization' },
  { label: t('finance.corporate'), value: 'corporate' },
  { label: t('finance.foundation'), value: 'foundation' }
]

const statusOptions = [
  { label: t('finance.active'), value: 'active' },
  { label: t('finance.inactive'), value: 'inactive' },
  { label: t('finance.lapsed'), value: 'lapsed' }
]

const loadDonor = async () => {
  if (!isEdit.value) return
  try {
    const donor = await financeService.getDonor(route.params.id)
    Object.assign(formData, donor)
    if (donor.address) {
      addressString.value = `${donor.address.street}, ${donor.address.city}, ${donor.address.state}, ${donor.address.postal_code}, ${donor.address.country}`
    }
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load donor', life: 3000 })
    router.push('/staff/finance/donors')
  }
}

const handleSubmit = async () => {
  try {
    saving.value = true
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
      await financeService.updateDonor(route.params.id, dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('finance.donorUpdated'), life: 3000 })
    } else {
      await financeService.createDonor(dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('finance.donorCreated'), life: 3000 })
    }
    router.push('/staff/finance/donors')
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save donor', life: 3000 })
  } finally {
    saving.value = false
  }
}

onMounted(loadDonor)
</script>

<style scoped>
.donor-form-container { max-width: 1000px; margin: 0 auto; }
.form-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
.form-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.5rem; }
.form-field label { font-weight: 600; color: #374151; }
.full-width { grid-column: 1 / -1; }
.form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 2rem; }
@media (max-width: 768px) { .form-grid { grid-template-columns: 1fr; } }
</style>
