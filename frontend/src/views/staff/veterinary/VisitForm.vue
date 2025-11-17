<template>
  <div class="visit-form-container">
    <div class="form-header">
      <Button
        icon="pi pi-arrow-left"
        class="p-button-text"
        @click="router.back()"
      />
      <h1>{{ isEdit ? 'Edit Visit' : $t('veterinary.addVisit') }}</h1>
    </div>

    <form @submit.prevent="handleSubmit">
      <Card>
        <template #title>
          {{ $t('veterinary.visit') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-field">
              <label for="animal_id">Animal ID *</label>
              <InputText
                id="animal_id"
                v-model="formData.animal_id"
                required
              />
            </div>

            <div class="form-field">
              <label for="visit_date">{{ $t('veterinary.visitDate') }} *</label>
              <Calendar
                id="visit_date"
                v-model="formData.visit_date"
                date-format="yy-mm-dd"
                required
              />
            </div>

            <div class="form-field">
              <label for="visit_type">{{ $t('veterinary.visitType') }} *</label>
              <Dropdown
                id="visit_type"
                v-model="formData.visit_type"
                :options="visitTypeOptions"
                option-label="label"
                option-value="value"
                required
              />
            </div>

            <div class="form-field">
              <label for="veterinarian_name">{{ $t('veterinary.veterinarianName') }} *</label>
              <InputText
                id="veterinarian_name"
                v-model="formData.veterinarian_name"
                required
              />
            </div>

            <div class="form-field">
              <label for="clinic_name">{{ $t('veterinary.clinicName') }}</label>
              <InputText
                id="clinic_name"
                v-model="formData.clinic_name"
              />
            </div>

            <div class="form-field">
              <label for="weight">{{ $t('veterinary.weight') }} (kg)</label>
              <InputNumber
                id="weight"
                v-model="formData.weight"
                :min="0"
                :max-fraction-digits="2"
              />
            </div>

            <div class="form-field">
              <label for="temperature">{{ $t('veterinary.temperature') }} (°C)</label>
              <InputNumber
                id="temperature"
                v-model="formData.temperature"
                :min="0"
                :max-fraction-digits="1"
              />
            </div>

            <div class="form-field">
              <label for="cost">{{ $t('veterinary.cost') }}</label>
              <InputNumber
                id="cost"
                v-model="formData.cost"
                mode="currency"
                currency="USD"
                :min="0"
              />
            </div>

            <div class="form-field full-width">
              <label for="reason">{{ $t('veterinary.reason') }} *</label>
              <Textarea
                id="reason"
                v-model="formData.reason"
                rows="3"
                required
              />
            </div>

            <div class="form-field full-width">
              <label for="diagnosis">{{ $t('veterinary.diagnosis') }}</label>
              <Textarea
                id="diagnosis"
                v-model="formData.diagnosis"
                rows="3"
              />
            </div>

            <div class="form-field full-width">
              <label for="treatment_provided">{{ $t('veterinary.treatmentProvided') }}</label>
              <Textarea
                id="treatment_provided"
                v-model="formData.treatment_provided"
                rows="3"
              />
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
                {{ $t('veterinary.followUpRequired') }}
              </label>
            </div>

            <div
              v-if="formData.follow_up_required"
              class="form-field"
            >
              <label for="follow_up_date">{{ $t('veterinary.followUpDate') }}</label>
              <Calendar
                id="follow_up_date"
                v-model="formData.follow_up_date"
                date-format="yy-mm-dd"
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
import { veterinaryService } from '@/services/veterinaryService'
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
  animal_id: '',
  visit_date: new Date(),
  visit_type: 'checkup',
  veterinarian_name: '',
  clinic_name: '',
  reason: '',
  diagnosis: '',
  treatment_provided: '',
  medications_prescribed: '',
  follow_up_required: false,
  follow_up_date: null,
  weight: null,
  temperature: null,
  heart_rate: null,
  cost: null,
  notes: ''
})

const visitTypeOptions = [
  { label: t('veterinary.checkup'), value: 'checkup' },
  { label: t('veterinary.emergency'), value: 'emergency' },
  { label: t('veterinary.surgery'), value: 'surgery' },
  { label: 'Vaccination', value: 'vaccination' },
  { label: t('veterinary.followUp'), value: 'follow_up' },
  { label: t('veterinary.other'), value: 'other' }
]

const loadVisit = async () => {
  if (!isEdit.value) return
  try {
    const visit = await veterinaryService.getVisit(route.params.id)
    Object.assign(formData, {
      ...visit,
      visit_date: visit.visit_date ? new Date(visit.visit_date) : null,
      follow_up_date: visit.follow_up_date ? new Date(visit.follow_up_date) : null
    })
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load visit', life: 3000 })
    router.push('/staff/veterinary/visits')
  }
}

const handleSubmit = async () => {
  try {
    saving.value = true
    const dataToSend = {
      ...formData,
      visit_date: formData.visit_date ? formData.visit_date.toISOString().split('T')[0] : null,
      follow_up_date: formData.follow_up_required && formData.follow_up_date ? formData.follow_up_date.toISOString().split('T')[0] : null
    }

    if (isEdit.value) {
      await veterinaryService.updateVisit(route.params.id, dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('veterinary.visitUpdated'), life: 3000 })
    } else {
      await veterinaryService.createVisit(dataToSend)
      toast.add({ severity: 'success', summary: 'Success', detail: t('veterinary.visitCreated'), life: 3000 })
    }
    router.push('/staff/veterinary/visits')
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to save visit', life: 3000 })
  } finally {
    saving.value = false
  }
}

onMounted(loadVisit)
</script>

<style scoped>
.visit-form-container { max-width: 1000px; margin: 0 auto; }
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
