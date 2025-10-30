<template>
  <div class="animal-form-page">
    <div class="page-header">
      <h1 class="page-title">{{ isEdit ? t('animals.editAnimal') : t('animals.addAnimal') }}</h1>
      <BaseButton variant="outline" size="small" @click="goBack">
        {{ t('common.cancel') }}
      </BaseButton>
    </div>

    <LoadingSpinner v-if="loading" />

    <form v-else @submit.prevent="handleSubmit" class="animal-form">
      <BaseCard>
        <template #header>{{ t('animals.basicInfo') }}</template>

        <div class="form-row">
          <FormGroup
            :label="t('animals.name')"
            :error="errors.name"
            required
          >
            <input
              v-model="form.name"
              type="text"
              :placeholder="t('animals.name')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('animals.species')"
            :error="errors.species"
            required
          >
            <select v-model="form.species" required>
              <option value="">{{ t('common.select') }}</option>
              <option value="dog">{{ t('animals.dog') }}</option>
              <option value="cat">{{ t('animals.cat') }}</option>
              <option value="other">{{ t('animals.other') }}</option>
            </select>
          </FormGroup>
        </div>

        <div class="form-row">
          <FormGroup
            :label="t('animals.breed')"
            :error="errors.breed"
            required
          >
            <input
              v-model="form.breed"
              type="text"
              :placeholder="t('animals.breed')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('animals.gender')"
            :error="errors.gender"
            required
          >
            <select v-model="form.gender" required>
              <option value="">{{ t('common.select') }}</option>
              <option value="male">{{ t('animals.male') }}</option>
              <option value="female">{{ t('animals.female') }}</option>
            </select>
          </FormGroup>
        </div>

        <div class="form-row">
          <FormGroup
            :label="t('animals.age')"
            :error="errors.age"
            required
          >
            <input
              v-model.number="form.age"
              type="number"
              min="0"
              step="0.1"
              :placeholder="t('animals.age')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('animals.weight')"
            :error="errors.weight"
            required
          >
            <input
              v-model.number="form.weight"
              type="number"
              min="0"
              step="0.1"
              :placeholder="t('animals.weight')"
              required
            />
          </FormGroup>
        </div>

        <div class="form-row">
          <FormGroup
            :label="t('animals.color')"
            :error="errors.color"
            required
          >
            <input
              v-model="form.color"
              type="text"
              :placeholder="t('animals.color')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('animals.size')"
            :error="errors.size"
            required
          >
            <select v-model="form.size" required>
              <option value="">{{ t('common.select') }}</option>
              <option value="small">{{ t('animals.size_small') }}</option>
              <option value="medium">{{ t('animals.size_medium') }}</option>
              <option value="large">{{ t('animals.size_large') }}</option>
            </select>
          </FormGroup>
        </div>

        <FormGroup
          :label="t('animals.description')"
          :error="errors.description"
          required
        >
          <textarea
            v-model="form.description"
            :placeholder="t('animals.description')"
            rows="4"
            required
          ></textarea>
        </FormGroup>
      </BaseCard>

      <BaseCard>
        <template #header>{{ t('animals.healthInfo') }}</template>

        <div class="checkbox-group">
          <label class="checkbox-label">
            <input v-model="form.sterilized" type="checkbox" />
            <span>{{ t('animals.sterilized') }}</span>
          </label>
          <label class="checkbox-label">
            <input v-model="form.vaccinated" type="checkbox" />
            <span>{{ t('animals.vaccinated') }}</span>
          </label>
          <label class="checkbox-label">
            <input v-model="form.chipped" type="checkbox" />
            <span>{{ t('animals.chipped') }}</span>
          </label>
        </div>

        <FormGroup
          v-if="form.chipped"
          :label="t('animals.chipNumber')"
          :error="errors.chip_number"
        >
          <input
            v-model="form.chip_number"
            type="text"
            :placeholder="t('animals.chipNumber')"
          />
        </FormGroup>

        <FormGroup
          :label="t('animals.medicalNotes')"
          :error="errors.medical_notes"
        >
          <textarea
            v-model="form.medical_notes"
            :placeholder="t('animals.medicalNotes')"
            rows="3"
          ></textarea>
        </FormGroup>
      </BaseCard>

      <BaseCard>
        <template #header>{{ t('animals.additionalInfo') }}</template>

        <FormGroup
          :label="t('animals.status')"
          :error="errors.status"
          required
        >
          <select v-model="form.status" required>
            <option value="">{{ t('common.select') }}</option>
            <option value="available">{{ t('animals.available') }}</option>
            <option value="adopted">{{ t('animals.adopted') }}</option>
            <option value="reserved">{{ t('animals.reserved') }}</option>
            <option value="medical_care">{{ t('animals.medical_care') }}</option>
            <option value="quarantine">{{ t('animals.quarantine') }}</option>
          </select>
        </FormGroup>

        <FormGroup
          :label="t('animals.location')"
          :error="errors.location"
        >
          <input
            v-model="form.location"
            type="text"
            :placeholder="t('animals.location')"
          />
        </FormGroup>

        <FormGroup
          :label="t('animals.admissionDate')"
          :error="errors.admission_date"
          required
        >
          <input
            v-model="form.admission_date"
            type="date"
            required
          />
        </FormGroup>

        <FormGroup
          :label="t('animals.behavioralNotes')"
          :error="errors.behavioral_notes"
        >
          <textarea
            v-model="form.behavioral_notes"
            :placeholder="t('animals.behavioralNotes')"
            rows="3"
          ></textarea>
        </FormGroup>

        <FormGroup
          :label="t('animals.specialNeeds')"
          :error="errors.special_needs"
        >
          <textarea
            v-model="form.special_needs"
            :placeholder="t('animals.specialNeeds')"
            rows="3"
          ></textarea>
        </FormGroup>
      </BaseCard>

      <div class="form-actions">
        <BaseButton variant="outline" type="button" @click="goBack">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="primary" type="submit" :loading="submitting">
          {{ isEdit ? t('common.update') : t('common.create') }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '../../stores/notification'
import { API } from '../../api'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseCard from '../../components/base/BaseCard.vue'
import FormGroup from '../../components/base/FormGroup.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const notificationStore = useNotificationStore()

const form = ref({
  name: '',
  species: '',
  breed: '',
  gender: '',
  age: null,
  weight: null,
  color: '',
  size: '',
  description: '',
  status: 'available',
  sterilized: false,
  vaccinated: false,
  chipped: false,
  chip_number: '',
  medical_notes: '',
  behavioral_notes: '',
  special_needs: '',
  location: '',
  admission_date: new Date().toISOString().split('T')[0],
})

const errors = ref({})
const loading = ref(false)
const submitting = ref(false)

const isEdit = computed(() => !!route.params.id)

onMounted(() => {
  if (isEdit.value) {
    fetchAnimal()
  }
})

async function fetchAnimal() {
  loading.value = true
  try {
    const response = await API.animals.getById(route.params.id)
    const animal = response.data

    // Populate form with animal data
    form.value = {
      name: animal.name || '',
      species: animal.species || '',
      breed: animal.breed || '',
      gender: animal.gender || '',
      age: animal.age || null,
      weight: animal.weight || null,
      color: animal.color || '',
      size: animal.size || '',
      description: animal.description || '',
      status: animal.status || 'available',
      sterilized: animal.sterilized || false,
      vaccinated: animal.vaccinated || false,
      chipped: animal.chipped || false,
      chip_number: animal.chip_number || '',
      medical_notes: animal.medical_notes || '',
      behavioral_notes: animal.behavioral_notes || '',
      special_needs: animal.special_needs || '',
      location: animal.location || '',
      admission_date: animal.admission_date ? animal.admission_date.split('T')[0] : new Date().toISOString().split('T')[0],
    }
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
    router.push({ name: 'animals-list' })
  } finally {
    loading.value = false
  }
}

function validateForm() {
  errors.value = {}
  let isValid = true

  if (!form.value.name || form.value.name.trim().length === 0) {
    errors.value.name = t('animals.nameRequired')
    isValid = false
  }

  if (!form.value.species) {
    errors.value.species = t('animals.speciesRequired')
    isValid = false
  }

  if (!form.value.breed || form.value.breed.trim().length === 0) {
    errors.value.breed = t('animals.breedRequired')
    isValid = false
  }

  if (!form.value.gender) {
    errors.value.gender = t('animals.genderRequired')
    isValid = false
  }

  if (!form.value.age || form.value.age <= 0) {
    errors.value.age = t('animals.ageRequired')
    isValid = false
  }

  if (!form.value.weight || form.value.weight <= 0) {
    errors.value.weight = t('animals.weightRequired')
    isValid = false
  }

  if (!form.value.description || form.value.description.trim().length === 0) {
    errors.value.description = t('animals.descriptionRequired')
    isValid = false
  }

  return isValid
}

async function handleSubmit() {
  if (!validateForm()) {
    return
  }

  submitting.value = true
  try {
    if (isEdit.value) {
      await API.animals.update(route.params.id, form.value)
      notificationStore.success(t('animals.updateSuccess'))
    } else {
      await API.animals.create(form.value)
      notificationStore.success(t('animals.createSuccess'))
    }
    router.push({ name: 'animals-list' })
  } catch (error) {
    notificationStore.error(t('common.error'), error.response?.data?.error || error.message)
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.animal-form-page {
  padding: 2rem;
  max-width: 1000px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  color: var(--text-primary);
  margin: 0;
}

.animal-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

input,
select,
textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
  font-family: inherit;
}

input:focus,
select:focus,
textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

textarea {
  resize: vertical;
}

.checkbox-group {
  display: flex;
  gap: 2rem;
  margin-bottom: 1rem;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  color: var(--text-primary);
}

.checkbox-label input[type="checkbox"] {
  width: auto;
  cursor: pointer;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1rem;
}

@media (max-width: 768px) {
  .animal-form-page {
    padding: 1rem;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
}
</style>
