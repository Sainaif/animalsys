<template>
  <div class="veterinary-form-page">
    <div class="page-header">
      <h1 class="page-title">
        {{ isEdit ? t('veterinary.editVisit') : t('veterinary.addVisit') }}
      </h1>
    </div>

    <BaseCard>
      <LoadingSpinner v-if="loading" />
      <form v-else @submit.prevent="handleSubmit">
        <!-- Visit Information -->
        <div class="form-section">
          <h3 class="section-title">{{ t('veterinary.visitInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('veterinary.animal')" :error="errors.animal_id" required>
              <select v-model="form.animal_id" class="form-control" :class="{ 'error': errors.animal_id }">
                <option value="">{{ t('veterinary.selectAnimal') }}</option>
                <option v-for="animal in animals" :key="animal.id" :value="animal.id">
                  {{ animal.name }} ({{ animal.species }})
                </option>
              </select>
            </FormGroup>

            <FormGroup :label="t('veterinary.visitType')" :error="errors.type" required>
              <select v-model="form.type" class="form-control" :class="{ 'error': errors.type }">
                <option value="">{{ t('veterinary.selectVisitType') }}</option>
                <option value="checkup">{{ t('veterinary.typeCheckup') }}</option>
                <option value="vaccination">{{ t('veterinary.typeVaccination') }}</option>
                <option value="treatment">{{ t('veterinary.typeTreatment') }}</option>
                <option value="surgery">{{ t('veterinary.typeSurgery') }}</option>
                <option value="emergency">{{ t('veterinary.typeEmergency') }}</option>
              </select>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('veterinary.visitDate')" :error="errors.visit_date" required>
              <input
                v-model="form.visit_date"
                type="datetime-local"
                class="form-control"
                :class="{ 'error': errors.visit_date }"
              />
            </FormGroup>

            <FormGroup :label="t('common.status')" :error="errors.status">
              <select v-model="form.status" class="form-control">
                <option value="scheduled">{{ t('veterinary.statusScheduled') }}</option>
                <option value="completed">{{ t('veterinary.statusCompleted') }}</option>
                <option value="cancelled">{{ t('veterinary.statusCancelled') }}</option>
              </select>
            </FormGroup>
          </div>
        </div>

        <!-- Veterinarian & Clinic -->
        <div class="form-section">
          <h3 class="section-title">{{ t('veterinary.veterinarianInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('veterinary.veterinarian')" :error="errors.veterinarian_name">
              <input
                v-model="form.veterinarian_name"
                type="text"
                class="form-control"
                :placeholder="t('veterinary.veterinarianPlaceholder')"
              />
            </FormGroup>

            <FormGroup :label="t('veterinary.clinic')" :error="errors.clinic_name">
              <input
                v-model="form.clinic_name"
                type="text"
                class="form-control"
                :placeholder="t('veterinary.clinicPlaceholder')"
              />
            </FormGroup>
          </div>
        </div>

        <!-- Medical Details -->
        <div class="form-section">
          <h3 class="section-title">{{ t('veterinary.medicalDetails') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('veterinary.diagnosis')" :error="errors.diagnosis">
              <textarea
                v-model="form.diagnosis"
                class="form-control"
                :placeholder="t('veterinary.diagnosisPlaceholder')"
                rows="3"
              ></textarea>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('veterinary.treatment')" :error="errors.treatment">
              <textarea
                v-model="form.treatment"
                class="form-control"
                :placeholder="t('veterinary.treatmentPlaceholder')"
                rows="3"
              ></textarea>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('veterinary.prescription')" :error="errors.prescription">
              <textarea
                v-model="form.prescription"
                class="form-control"
                :placeholder="t('veterinary.prescriptionPlaceholder')"
                rows="2"
              ></textarea>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('veterinary.cost')" :error="errors.cost">
              <input
                v-model.number="form.cost"
                type="number"
                step="0.01"
                min="0"
                class="form-control"
                :placeholder="t('veterinary.costPlaceholder')"
              />
            </FormGroup>

            <FormGroup :label="t('veterinary.nextVisit')" :error="errors.next_visit_date">
              <input
                v-model="form.next_visit_date"
                type="date"
                class="form-control"
              />
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('common.notes')" :error="errors.notes">
              <textarea
                v-model="form.notes"
                class="form-control"
                :placeholder="t('veterinary.notesPlaceholder')"
                rows="4"
              ></textarea>
            </FormGroup>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="form-actions">
          <BaseButton type="button" variant="secondary" @click="goBack">
            {{ t('common.cancel') }}
          </BaseButton>
          <BaseButton type="submit" variant="primary" :disabled="submitting">
            {{ submitting ? t('common.saving') : t('common.save') }}
          </BaseButton>
        </div>
      </form>
    </BaseCard>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '../../stores/notifications'
import { API } from '../../api'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import FormGroup from '../../components/base/FormGroup.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const notificationStore = useNotificationStore()

const isEdit = ref(false)
const loading = ref(false)
const submitting = ref(false)
const animals = ref([])

const form = reactive({
  animal_id: '',
  type: '',
  visit_date: '',
  status: 'scheduled',
  veterinarian_name: '',
  clinic_name: '',
  diagnosis: '',
  treatment: '',
  prescription: '',
  cost: null,
  next_visit_date: '',
  notes: '',
})

const errors = reactive({})

async function fetchAnimals() {
  try {
    const response = await API.animals.list({ limit: 1000 })
    animals.value = response.data.data || []
  } catch (error) {
    console.error('Failed to fetch animals:', error)
  }
}

async function fetchVisit() {
  if (!route.params.id) return

  try {
    loading.value = true
    const response = await API.veterinary.getById(route.params.id)
    const visit = response.data

    Object.keys(form).forEach(key => {
      if (visit[key] !== undefined) {
        form[key] = visit[key]
      }
    })

    // Format dates for input
    if (form.visit_date) {
      form.visit_date = new Date(form.visit_date).toISOString().slice(0, 16)
    }
    if (form.next_visit_date) {
      form.next_visit_date = new Date(form.next_visit_date).toISOString().split('T')[0]
    }
  } catch (error) {
    console.error('Failed to fetch visit:', error)
    notificationStore.error(t('veterinary.fetchError'))
    goBack()
  } finally {
    loading.value = false
  }
}

function validateForm() {
  Object.keys(errors).forEach(key => delete errors[key])
  let isValid = true

  if (!form.animal_id) {
    errors.animal_id = t('common.required')
    isValid = false
  }

  if (!form.type) {
    errors.type = t('common.required')
    isValid = false
  }

  if (!form.visit_date) {
    errors.visit_date = t('common.required')
    isValid = false
  }

  return isValid
}

async function handleSubmit() {
  if (!validateForm()) {
    notificationStore.error(t('common.fixErrors'))
    return
  }

  try {
    submitting.value = true

    if (isEdit.value) {
      await API.veterinary.update(route.params.id, form)
      notificationStore.success(t('veterinary.updateSuccess'))
    } else {
      await API.veterinary.create(form)
      notificationStore.success(t('veterinary.createSuccess'))
    }

    goBack()
  } catch (error) {
    console.error('Failed to save visit:', error)
    notificationStore.error(
      isEdit.value ? t('veterinary.updateError') : t('veterinary.createError')
    )
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.push({ name: 'Veterinary' })
}

onMounted(() => {
  fetchAnimals()
  if (route.params.id) {
    isEdit.value = true
    fetchVisit()
  }
})
</script>

<style scoped>
.veterinary-form-page {
  max-width: 900px;
  padding: 2rem;
}

.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin: 0;
}

.form-section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--border-color);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-control {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-control.error {
  border-color: var(--danger-color);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}
</style>
