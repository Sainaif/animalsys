<template>
  <div class="donor-form-page">
    <div class="page-header">
      <h1 class="page-title">
        {{ isEdit ? t('donors.editDonor') : t('donors.addDonor') }}
      </h1>
    </div>

    <BaseCard>
      <LoadingSpinner v-if="loading" />
      <form v-else @submit.prevent="handleSubmit">
        <!-- Basic Information -->
        <div class="form-section">
          <h3 class="section-title">{{ t('donors.basicInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('donors.type')" :error="errors.type" required>
              <select v-model="form.type" class="form-control" :class="{ 'error': errors.type }">
                <option value="">{{ t('donors.selectType') }}</option>
                <option value="individual">{{ t('donors.typeIndividual') }}</option>
                <option value="company">{{ t('donors.typeCompany') }}</option>
                <option value="foundation">{{ t('donors.typeFoundation') }}</option>
              </select>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('donors.name')" :error="errors.name" required>
              <input
                v-model="form.name"
                type="text"
                class="form-control"
                :class="{ 'error': errors.name }"
                :placeholder="t('donors.namePlaceholder')"
              />
            </FormGroup>
          </div>
        </div>

        <!-- Contact Information -->
        <div class="form-section">
          <h3 class="section-title">{{ t('common.contactInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('common.email')" :error="errors.email">
              <input
                v-model="form.email"
                type="email"
                class="form-control"
                :class="{ 'error': errors.email }"
                :placeholder="t('common.emailPlaceholder')"
              />
            </FormGroup>

            <FormGroup :label="t('common.phone')" :error="errors.phone">
              <input
                v-model="form.phone"
                type="tel"
                class="form-control"
                :class="{ 'error': errors.phone }"
                :placeholder="t('common.phonePlaceholder')"
              />
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('common.address')" :error="errors.address">
              <textarea
                v-model="form.address"
                class="form-control"
                :class="{ 'error': errors.address }"
                :placeholder="t('common.addressPlaceholder')"
                rows="3"
              ></textarea>
            </FormGroup>
          </div>
        </div>

        <!-- Additional Information -->
        <div class="form-section">
          <h3 class="section-title">{{ t('donors.additionalInfo') }}</h3>

          <div class="form-row">
            <FormGroup
              v-if="form.type === 'company' || form.type === 'foundation'"
              :label="t('donors.taxId')"
              :error="errors.tax_id"
            >
              <input
                v-model="form.tax_id"
                type="text"
                class="form-control"
                :class="{ 'error': errors.tax_id }"
                :placeholder="t('donors.taxIdPlaceholder')"
              />
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('common.status')" :error="errors.status">
              <select v-model="form.status" class="form-control" :class="{ 'error': errors.status }">
                <option value="active">{{ t('donors.statusActive') }}</option>
                <option value="inactive">{{ t('donors.statusInactive') }}</option>
              </select>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('donors.preferredContactMethod')" :error="errors.preferred_contact_method">
              <select v-model="form.preferred_contact_method" class="form-control">
                <option value="">{{ t('donors.selectContactMethod') }}</option>
                <option value="email">{{ t('common.email') }}</option>
                <option value="phone">{{ t('common.phone') }}</option>
                <option value="mail">{{ t('donors.mail') }}</option>
              </select>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('donors.notes')" :error="errors.notes">
              <textarea
                v-model="form.notes"
                class="form-control"
                :class="{ 'error': errors.notes }"
                :placeholder="t('donors.notesPlaceholder')"
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

const form = reactive({
  type: '',
  name: '',
  email: '',
  phone: '',
  address: '',
  tax_id: '',
  status: 'active',
  preferred_contact_method: '',
  notes: '',
})

const errors = reactive({})

async function fetchDonor() {
  if (!route.params.id) return

  try {
    loading.value = true
    const response = await API.donors.getById(route.params.id)
    const donor = response.data

    Object.keys(form).forEach(key => {
      if (donor[key] !== undefined) {
        form[key] = donor[key]
      }
    })
  } catch (error) {
    console.error('Failed to fetch donor:', error)
    notificationStore.error(t('donors.fetchError'))
    goBack()
  } finally {
    loading.value = false
  }
}

function validateForm() {
  Object.keys(errors).forEach(key => delete errors[key])
  let isValid = true

  if (!form.type || form.type.trim().length === 0) {
    errors.type = t('common.required')
    isValid = false
  }

  if (!form.name || form.name.trim().length === 0) {
    errors.name = t('common.required')
    isValid = false
  }

  if (form.email && !isValidEmail(form.email)) {
    errors.email = t('common.invalidEmail')
    isValid = false
  }

  return isValid
}

function isValidEmail(email) {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

async function handleSubmit() {
  if (!validateForm()) {
    notificationStore.error(t('common.fixErrors'))
    return
  }

  try {
    submitting.value = true

    if (isEdit.value) {
      await API.donors.update(route.params.id, form)
      notificationStore.success(t('donors.updateSuccess'))
    } else {
      await API.donors.create(form)
      notificationStore.success(t('donors.createSuccess'))
    }

    goBack()
  } catch (error) {
    console.error('Failed to save donor:', error)
    notificationStore.error(
      isEdit.value ? t('donors.updateError') : t('donors.createError')
    )
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.push({ name: 'Donors' })
}

onMounted(() => {
  if (route.params.id) {
    isEdit.value = true
    fetchDonor()
  }
})
</script>

<style scoped>
.donor-form-page {
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

.form-control:disabled {
  background: var(--bg-secondary);
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}
</style>
