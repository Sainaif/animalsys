<template>
  <div class="donation-form-container">
    <div class="form-header">
      <Button
        icon="pi pi-arrow-left"
        class="p-button-text"
        @click="router.back()"
      />
      <h1>{{ isEdit ? 'Edit Donation' : $t('finance.addDonation') }}</h1>
    </div>

    <form @submit.prevent="handleSubmit">
      <Card>
        <template #title>
          {{ $t('finance.donation') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-field full-width">
              <label for="donor_id">Donor ID *</label>
              <InputText
                id="donor_id"
                v-model="formData.donor_id"
                required
              />
            </div>

            <div class="form-field">
              <label for="amount">{{ $t('finance.amount') }} *</label>
              <InputNumber
                id="amount"
                v-model="formData.amount"
                mode="currency"
                currency="USD"
                :min="0"
                required
              />
            </div>

            <div class="form-field">
              <label for="donation_date">{{ $t('finance.donationDate') }} *</label>
              <Calendar
                id="donation_date"
                v-model="formData.donation_date"
                date-format="yy-mm-dd"
                required
              />
            </div>

            <div class="form-field">
              <label for="donation_type">{{ $t('finance.donationType') }} *</label>
              <Dropdown
                id="donation_type"
                v-model="formData.donation_type"
                :options="donationTypeOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="form-field">
              <label for="payment_method">{{ $t('finance.paymentMethod') }} *</label>
              <Dropdown
                id="payment_method"
                v-model="formData.payment_method"
                :options="paymentMethodOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="form-field">
              <label for="payment_status">{{ $t('finance.paymentStatus') }}</label>
              <Dropdown
                id="payment_status"
                v-model="formData.payment_status"
                :options="paymentStatusOptions"
                option-label="label"
                option-value="value"
              />
            </div>

            <div class="form-field">
              <label for="transaction_id">Transaction ID</label>
              <InputText
                id="transaction_id"
                v-model="formData.transaction_id"
              />
            </div>

            <div class="form-field">
              <label
                for="tax_deductible"
                class="checkbox-label"
              >
                <Checkbox
                  id="tax_deductible"
                  v-model="formData.tax_deductible"
                  :binary="true"
                />
                {{ $t('finance.taxDeductible') }}
              </label>
            </div>

            <div class="form-field">
              <label
                for="anonymous"
                class="checkbox-label"
              >
                <Checkbox
                  id="anonymous"
                  v-model="formData.anonymous"
                  :binary="true"
                />
                {{ $t('finance.anonymous') }}
              </label>
            </div>

            <div
              v-if="formData.donation_type === 'recurring'"
              class="form-field"
            >
              <label for="recurring_frequency">{{ $t('finance.recurringFrequency') }}</label>
              <Dropdown
                id="recurring_frequency"
                v-model="formData.recurring_frequency"
                :options="frequencyOptions"
                option-label="label"
                option-value="value"
              />
            </div>

            <div class="form-field full-width">
              <label for="purpose">{{ $t('finance.purpose') }}</label>
              <Textarea
                id="purpose"
                v-model="formData.purpose"
                rows="2"
              />
            </div>

            <div class="form-field full-width">
              <label for="notes">Notes</label>
              <Textarea
                id="notes"
                v-model="formData.notes"
                rows="2"
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
import Checkbox from 'primevue/checkbox'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const toast = useToast()

const isEdit = computed(() => !!route.params.id)
const saving = ref(false)

const formData = reactive({
  donor_id: '',
  amount: 0,
  currency: 'USD',
  donation_date: new Date(),
  donation_type: 'one_time',
  payment_method: 'credit_card',
  payment_status: 'completed',
  transaction_id: '',
  receipt_sent: false,
  tax_deductible: true,
  purpose: '',
  anonymous: false,
  recurring_frequency: null,
  notes: ''
})

const donationTypeOptions = [
  { label: t('finance.oneTime'), value: 'one_time' },
  { label: t('finance.recurring'), value: 'recurring' },
  { label: t('finance.pledge'), value: 'pledge' }
]

const paymentMethodOptions = [
  { label: t('finance.cash'), value: 'cash' },
  { label: t('finance.creditCard'), value: 'credit_card' },
  { label: t('finance.debitCard'), value: 'debit_card' },
  { label: t('finance.bankTransfer'), value: 'bank_transfer' },
  { label: t('finance.check'), value: 'check' },
  { label: 'PayPal', value: 'paypal' }
]

const paymentStatusOptions = [
  { label: t('finance.pending'), value: 'pending' },
  { label: t('finance.completed'), value: 'completed' },
  { label: t('finance.failed'), value: 'failed' },
  { label: t('finance.refunded'), value: 'refunded' }
]

const frequencyOptions = [
  { label: t('finance.weekly'), value: 'weekly' },
  { label: t('finance.monthly'), value: 'monthly' },
  { label: t('finance.quarterly'), value: 'quarterly' },
  { label: t('finance.annually'), value: 'annually' }
]

const loadDonation = async () => {
  if (!isEdit.value) return
  try {
    const donation = await financeService.getDonation(route.params.id)
    Object.assign(formData, {
      ...donation,
      donation_date: donation.donation_date ? new Date(donation.donation_date) : new Date()
    })
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to load donation', life: 3000 })
    router.push('/staff/finance/donations')
  }
}

const handleSubmit = async () => {
  try {
    saving.value = true
    const dataToSend = {
      ...formData,
      donation_date: formData.donation_date ? formData.donation_date.toISOString().split('T')[0] : null
    }

    if (isEdit.value) {
      await financeService.updateDonation(route.params.id, dataToSend)
      toast.add({ severity: 'success', summary: t('common.success'), detail: t('finance.donationUpdated'), life: 3000 })
    } else {
      await financeService.createDonation(dataToSend)
      toast.add({ severity: 'success', summary: t('common.success'), detail: t('finance.donationCreated'), life: 3000 })
    }
    router.push('/staff/finance/donations')
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to save donation', life: 3000 })
  } finally {
    saving.value = false
  }
}

onMounted(loadDonation)
</script>

<style scoped>
.donation-form-container { max-width: 1000px; margin: 0 auto; }
.form-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 2rem; }
.form-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
.form-grid { display: grid; grid-template-columns: repeat(2, 1fr); gap: 1.5rem; }
.form-field { display: flex; flex-direction: column; gap: 0.5rem; }
.form-field label { font-weight: 600; color: var(--text-color); }
.checkbox-label { flex-direction: row !important; align-items: center; gap: 0.75rem !important; }
.full-width { grid-column: 1 / -1; }
.form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 2rem; }
@media (max-width: 768px) { .form-grid { grid-template-columns: 1fr; } }
</style>
