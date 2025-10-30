<template>
  <div class="adoption-form-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('adoptions.newApplication') }}</h1>
      <BaseButton variant="outline" size="small" @click="goBack">
        {{ t('common.cancel') }}
      </BaseButton>
    </div>

    <form @submit.prevent="handleSubmit" class="adoption-form">
      <BaseCard>
        <template #header>{{ t('adoptions.animalSelection') }}</template>

        <FormGroup
          :label="t('animals.name')"
          :error="errors.animal_id"
          required
        >
          <select v-model="form.animal_id" required>
            <option value="">{{ t('adoptions.selectAnimal') }}</option>
            <option
              v-for="animal in availableAnimals"
              :key="animal.id"
              :value="animal.id"
            >
              {{ animal.name }} ({{ t(`animals.${animal.species}`) }})
            </option>
          </select>
        </FormGroup>
      </BaseCard>

      <BaseCard>
        <template #header>{{ t('adoptions.applicantInfo') }}</template>

        <div class="form-row">
          <FormGroup
            :label="t('common.firstName')"
            :error="errors.first_name"
            required
          >
            <input
              v-model="form.first_name"
              type="text"
              :placeholder="t('common.firstName')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('common.lastName')"
            :error="errors.last_name"
            required
          >
            <input
              v-model="form.last_name"
              type="text"
              :placeholder="t('common.lastName')"
              required
            />
          </FormGroup>
        </div>

        <div class="form-row">
          <FormGroup
            :label="t('common.email')"
            :error="errors.email"
            required
          >
            <input
              v-model="form.email"
              type="email"
              :placeholder="t('common.email')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('common.phone')"
            :error="errors.phone"
            required
          >
            <input
              v-model="form.phone"
              type="tel"
              :placeholder="t('common.phone')"
              required
            />
          </FormGroup>
        </div>

        <FormGroup
          :label="t('common.address')"
          :error="errors.address"
          required
        >
          <textarea
            v-model="form.address"
            :placeholder="t('common.address')"
            rows="3"
            required
          ></textarea>
        </FormGroup>
      </BaseCard>

      <BaseCard>
        <template #header>{{ t('adoptions.housingInfo') }}</template>

        <FormGroup
          :label="t('adoptions.housingType')"
          :error="errors.housing_type"
          required
        >
          <select v-model="form.housing_type" required>
            <option value="">{{ t('common.select') }}</option>
            <option value="house">{{ t('adoptions.house') }}</option>
            <option value="apartment">{{ t('adoptions.apartment') }}</option>
            <option value="other">{{ t('common.other') }}</option>
          </select>
        </FormGroup>

        <div class="checkbox-group">
          <label class="checkbox-label">
            <input v-model="form.has_yard" type="checkbox" />
            <span>{{ t('adoptions.hasYard') }}</span>
          </label>
          <label class="checkbox-label">
            <input v-model="form.has_other_pets" type="checkbox" />
            <span>{{ t('adoptions.hasOtherPets') }}</span>
          </label>
        </div>

        <FormGroup
          v-if="form.has_other_pets"
          :label="t('adoptions.otherPetsDescription')"
          :error="errors.other_pets_description"
        >
          <textarea
            v-model="form.other_pets_description"
            :placeholder="t('adoptions.otherPetsDescription')"
            rows="3"
          ></textarea>
        </FormGroup>
      </BaseCard>

      <BaseCard>
        <template #header>{{ t('adoptions.experienceInfo') }}</template>

        <FormGroup
          :label="t('adoptions.previousPetExperience')"
          :error="errors.previous_pet_experience"
          required
        >
          <textarea
            v-model="form.previous_pet_experience"
            :placeholder="t('adoptions.previousPetExperience')"
            rows="4"
            required
          ></textarea>
        </FormGroup>

        <FormGroup
          :label="t('adoptions.reasonForAdoption')"
          :error="errors.reason_for_adoption"
          required
        >
          <textarea
            v-model="form.reason_for_adoption"
            :placeholder="t('adoptions.reasonForAdoption')"
            rows="4"
            required
          ></textarea>
        </FormGroup>
      </BaseCard>

      <div class="form-actions">
        <BaseButton variant="outline" type="button" @click="goBack">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="primary" type="submit" :loading="submitting">
          {{ t('adoptions.submitApplication') }}
        </BaseButton>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '../../stores/notification'
import { API } from '../../api'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseCard from '../../components/base/BaseCard.vue'
import FormGroup from '../../components/base/FormGroup.vue'

const { t } = useI18n()
const router = useRouter()
const notificationStore = useNotificationStore()

const availableAnimals = ref([])

const form = ref({
  animal_id: '',
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  address: '',
  housing_type: '',
  has_yard: false,
  has_other_pets: false,
  other_pets_description: '',
  previous_pet_experience: '',
  reason_for_adoption: ''
})

const errors = ref({})
const submitting = ref(false)

onMounted(() => {
  fetchAvailableAnimals()
})

async function fetchAvailableAnimals() {
  try {
    const response = await API.animals.getAvailable()
    availableAnimals.value = response.data || []
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  }
}

function validateForm() {
  errors.value = {}
  let isValid = true

  if (!form.value.animal_id) {
    errors.value.animal_id = t('adoptions.animalRequired')
    isValid = false
  }

  if (!form.value.first_name || form.value.first_name.trim().length === 0) {
    errors.value.first_name = t('common.required')
    isValid = false
  }

  if (!form.value.last_name || form.value.last_name.trim().length === 0) {
    errors.value.last_name = t('common.required')
    isValid = false
  }

  if (!form.value.email || !form.value.email.includes('@')) {
    errors.value.email = t('common.invalid')
    isValid = false
  }

  if (!form.value.phone || form.value.phone.trim().length === 0) {
    errors.value.phone = t('common.required')
    isValid = false
  }

  if (!form.value.address || form.value.address.trim().length === 0) {
    errors.value.address = t('common.required')
    isValid = false
  }

  if (!form.value.housing_type) {
    errors.value.housing_type = t('common.required')
    isValid = false
  }

  if (!form.value.previous_pet_experience || form.value.previous_pet_experience.trim().length === 0) {
    errors.value.previous_pet_experience = t('common.required')
    isValid = false
  }

  if (!form.value.reason_for_adoption || form.value.reason_for_adoption.trim().length === 0) {
    errors.value.reason_for_adoption = t('common.required')
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
    await API.adoptions.create(form.value)
    notificationStore.success(t('adoptions.submitSuccess'))
    router.push({ name: 'adoptions-list' })
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
.adoption-form-page {
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

.adoption-form {
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
  .adoption-form-page {
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
