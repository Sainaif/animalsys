<template>
  <div class="adoption-form-container">
    <div class="form-header">
      <Button
        icon="pi pi-arrow-left"
        class="p-button-text"
        @click="router.back()"
      />
      <h1>{{ $t('adoption.createAdoption') }}</h1>
    </div>

    <form @submit.prevent="handleSubmit">
      <Card>
        <template #title>
          {{ $t('adoption.applicantInfo') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div
              v-if="applicationId"
              class="form-field full-width"
            >
              <label>{{ $t('adoption.applicationDetail') }}</label>
              <InputText
                :value="`Application #${applicationId}`"
                disabled
              />
            </div>

            <div class="form-field">
              <label for="adopter_first_name">{{ $t('adoption.firstName') }} *</label>
              <InputText
                id="adopter_first_name"
                v-model="formData.adopter_first_name"
                required
              />
            </div>

            <div class="form-field">
              <label for="adopter_last_name">{{ $t('adoption.lastName') }} *</label>
              <InputText
                id="adopter_last_name"
                v-model="formData.adopter_last_name"
                required
              />
            </div>

            <div class="form-field">
              <label for="adopter_email">{{ $t('adoption.email') }} *</label>
              <InputText
                id="adopter_email"
                v-model="formData.adopter_email"
                type="email"
                required
              />
            </div>

            <div class="form-field">
              <label for="adopter_phone">{{ $t('adoption.phone') }}</label>
              <InputText
                id="adopter_phone"
                v-model="formData.adopter_phone"
              />
            </div>

            <div class="form-field full-width">
              <label for="adopter_address">{{ $t('adoption.address') }}</label>
              <Textarea
                id="adopter_address"
                v-model="formData.adopter_address"
                rows="3"
              />
            </div>
          </div>
        </template>
      </Card>

      <Card class="mt-3">
        <template #title>
          {{ $t('adoption.adoptionInfo') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label for="adoption_date">{{ $t('adoption.adoptionDate') }} *</label>
              <Calendar
                id="adoption_date"
                v-model="formData.adoption_date"
                date-format="yy-mm-dd"
                required
              />
            </div>

            <div class="form-field">
              <label for="adoption_fee">{{ $t('adoption.adoptionFee') }} *</label>
              <InputNumber
                id="adoption_fee"
                v-model="formData.adoption_fee"
                mode="currency"
                currency="USD"
                :min="0"
                required
              />
            </div>

            <div class="form-field">
              <label for="payment_status">{{ $t('adoption.paymentStatus') }}</label>
              <Dropdown
                id="payment_status"
                v-model="formData.payment_status"
                :options="paymentStatusOptions"
                option-label="label"
                option-value="value"
              />
            </div>

            <div class="form-field">
              <label for="payment_method">{{ $t('adoption.paymentMethod') }}</label>
              <Dropdown
                id="payment_method"
                v-model="formData.payment_method"
                :options="paymentMethodOptions"
                option-label="label"
                option-value="value"
              />
            </div>

            <div class="form-field">
              <label
                for="contract_signed"
                class="checkbox-label"
              >
                <Checkbox
                  id="contract_signed"
                  v-model="formData.contract_signed"
                  :binary="true"
                />
                {{ $t('adoption.contractSigned') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="microchip_transferred"
                class="checkbox-label"
              >
                <Checkbox
                  id="microchip_transferred"
                  v-model="formData.microchip_transferred"
                  :binary="true"
                />
                {{ $t('adoption.microchipTransferred') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="return_policy_explained"
                class="checkbox-label"
              >
                <Checkbox
                  id="return_policy_explained"
                  v-model="formData.return_policy_explained"
                  :binary="true"
                />
                {{ $t('adoption.returnPolicyExplained') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="follow_up_required"
                class="checkbox-label"
              >
                <Checkbox
                  id="follow_up_required"
                  v-model="formData.follow_up_required"
                  :binary="true"
                />
                {{ $t('adoption.followUpRequired') }}
              </label>
            </div>
          </div>

          <div
            v-if="formData.follow_up_required"
            class="form-field full-width"
          >
            <label for="follow_up_schedule">{{ $t('adoption.followUpSchedule') }}</label>
            <Textarea
              id="follow_up_schedule"
              v-model="formData.follow_up_schedule"
              rows="3"
              placeholder="e.g., 1 week, 1 month, 3 months"
            />
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
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { adoptionService } from '@/services/adoptionService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Calendar from 'primevue/calendar'
import Checkbox from 'primevue/checkbox'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const applicationId = route.query.applicationId
const saving = ref(false)

const formData = reactive({
  application_id: applicationId || null,
  animal_id: route.query.animalId || null,
  adopter_first_name: '',
  adopter_last_name: '',
  adopter_email: '',
  adopter_phone: '',
  adopter_address: '',
  adoption_date: new Date(),
  adoption_fee: 0,
  payment_status: 'pending',
  payment_method: null,
  contract_signed: false,
  microchip_transferred: false,
  return_policy_explained: false,
  follow_up_required: false,
  follow_up_schedule: '',
  status: 'active'
})

const paymentStatusOptions = [
  { label: 'Pending', value: 'pending' },
  { label: 'Paid', value: 'paid' },
  { label: 'Refunded', value: 'refunded' }
]

const paymentMethodOptions = [
  { label: 'Cash', value: 'cash' },
  { label: 'Credit Card', value: 'credit_card' },
  { label: 'Debit Card', value: 'debit_card' },
  { label: 'Bank Transfer', value: 'bank_transfer' },
  { label: 'Other', value: 'other' }
]

const loadApplicationData = async () => {
  if (!applicationId) return

  try {
    const application = await adoptionService.getApplication(applicationId)
    formData.animal_id = application.animal_id
    formData.adopter_first_name = application.applicant_first_name
    formData.adopter_last_name = application.applicant_last_name
    formData.adopter_email = application.email
    formData.adopter_phone = application.phone
    if (application.address) {
      formData.adopter_address = `${application.address.street}, ${application.address.city}, ${application.address.state} ${application.address.postal_code}, ${application.address.country}`
    }
  } catch (error) {
    console.error('Error loading application:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to load application data',
      life: 3000
    })
  }
}

const handleSubmit = async () => {
  try {
    saving.value = true

    const dataToSend = {
      ...formData,
      adoption_date: formData.adoption_date ? formData.adoption_date.toISOString().split('T')[0] : null,
      follow_up_schedule: formData.follow_up_required ? formData.follow_up_schedule : null
    }

    await adoptionService.createAdoption(dataToSend)
    toast.add({
      severity: 'success',
      summary: 'Success',
      detail: t('adoption.adoptionCreated'),
      life: 3000
    })
    router.push('/staff/adoptions')
  } catch (error) {
    console.error('Error creating adoption:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to create adoption',
      life: 3000
    })
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadApplicationData()
})
</script>

<style scoped>
.adoption-form-container {
  max-width: 1000px;
  margin: 0 auto;
}

.form-header {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 2rem;
}

.form-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-field label {
  font-weight: 600;
  color: #374151;
}

.checkbox-label {
  flex-direction: row !important;
  align-items: center;
  gap: 0.75rem !important;
}

.full-width {
  grid-column: 1 / -1;
}

.mt-3 {
  margin-top: 1.5rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  margin-top: 2rem;
}

@media (max-width: 768px) {
  .form-grid {
    grid-template-columns: 1fr;
  }
}
</style>
