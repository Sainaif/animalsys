<template>
  <div class="volunteer-form-page">
    <div class="page-header">
      <h1 class="page-title">{{ isEdit ? t('volunteers.editVolunteer') : t('volunteers.addVolunteer') }}</h1>
      <BaseButton variant="outline" size="small" @click="goBack">
        {{ t('common.cancel') }}
      </BaseButton>
    </div>

    <LoadingSpinner v-if="loading" />

    <form v-else @submit.prevent="handleSubmit" class="volunteer-form">
      <BaseCard>
        <template #header>{{ t('volunteers.personalInfo') }}</template>

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
        >
          <textarea
            v-model="form.address"
            :placeholder="t('common.address')"
            rows="2"
          ></textarea>
        </FormGroup>
      </BaseCard>

      <BaseCard>
        <template #header>{{ t('volunteers.volunteering') }}</template>

        <div class="form-row">
          <FormGroup
            :label="t('common.status')"
            :error="errors.status"
            required
          >
            <select v-model="form.status" required>
              <option value="">{{ t('common.select') }}</option>
              <option value="active">{{ t('volunteers.statusActive') }}</option>
              <option value="inactive">{{ t('volunteers.statusInactive') }}</option>
              <option value="on_leave">{{ t('volunteers.statusOnLeave') }}</option>
            </select>
          </FormGroup>

          <FormGroup
            :label="t('volunteers.registrationDate')"
            :error="errors.registration_date"
            required
          >
            <input
              v-model="form.registration_date"
              type="date"
              required
            />
          </FormGroup>
        </div>

        <FormGroup
          :label="t('volunteers.skills')"
          :error="errors.skills"
          :hint="t('volunteers.skillsHint')"
        >
          <textarea
            v-model="form.skills"
            :placeholder="t('volunteers.skillsPlaceholder')"
            rows="3"
          ></textarea>
        </FormGroup>

        <FormGroup
          :label="t('volunteers.availability')"
          :error="errors.availability"
          :hint="t('volunteers.availabilityHint')"
        >
          <textarea
            v-model="form.availability"
            :placeholder="t('volunteers.availabilityPlaceholder')"
            rows="3"
          ></textarea>
        </FormGroup>

        <FormGroup
          :label="t('common.notes')"
          :error="errors.notes"
        >
          <textarea
            v-model="form.notes"
            :placeholder="t('common.notes')"
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
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  address: '',
  status: 'active',
  registration_date: new Date().toISOString().split('T')[0],
  skills: '',
  availability: '',
  notes: ''
})

const errors = ref({})
const loading = ref(false)
const submitting = ref(false)

const isEdit = computed(() => !!route.params.id)

onMounted(() => {
  if (isEdit.value) {
    fetchVolunteer()
  }
})

async function fetchVolunteer() {
  loading.value = true
  try {
    const response = await API.volunteers.getById(route.params.id)
    const volunteer = response.data

    form.value = {
      first_name: volunteer.first_name || '',
      last_name: volunteer.last_name || '',
      email: volunteer.email || '',
      phone: volunteer.phone || '',
      address: volunteer.address || '',
      status: volunteer.status || 'active',
      registration_date: volunteer.registration_date ? volunteer.registration_date.split('T')[0] : new Date().toISOString().split('T')[0],
      skills: volunteer.skills || '',
      availability: volunteer.availability || '',
      notes: volunteer.notes || ''
    }
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
    router.push({ name: 'volunteers-list' })
  } finally {
    loading.value = false
  }
}

function validateForm() {
  errors.value = {}
  let isValid = true

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

  if (!form.value.status) {
    errors.value.status = t('common.required')
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
      await API.volunteers.update(route.params.id, form.value)
      notificationStore.success(t('volunteers.updateSuccess'))
    } else {
      await API.volunteers.create(form.value)
      notificationStore.success(t('volunteers.createSuccess'))
    }
    router.push({ name: 'volunteers-list' })
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
.volunteer-form-page {
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

.volunteer-form {
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

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1rem;
}

@media (max-width: 768px) {
  .volunteer-form-page {
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
