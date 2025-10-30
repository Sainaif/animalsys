<template>
  <form @submit.prevent="handleSubmit" class="partner-form">
    <!-- Basic Information -->
    <div class="form-section">
      <h3 class="section-title">{{ t('partners.basicInfo') }}</h3>

      <FormGroup :label="t('partners.name')" required>
        <input
          v-model="formData.name"
          type="text"
          :placeholder="t('partners.namePlaceholder')"
          required
          class="form-input"
        />
      </FormGroup>

      <FormGroup :label="t('partners.partnerType')" required>
        <select v-model="formData.type" required class="form-select">
          <option value="">{{ t('partners.selectType') }}</option>
          <option value="veterinary">{{ t('partners.typeVeterinary') }}</option>
          <option value="shelter">{{ t('partners.typeShelter') }}</option>
          <option value="pet_store">{{ t('partners.typePetStore') }}</option>
          <option value="corporate">{{ t('partners.typeCorporate') }}</option>
          <option value="foundation">{{ t('partners.typeFoundation') }}</option>
          <option value="individual">{{ t('partners.typeIndividual') }}</option>
          <option value="other">{{ t('partners.typeOther') }}</option>
        </select>
      </FormGroup>

      <FormGroup :label="t('common.status')" required>
        <select v-model="formData.status" required class="form-select">
          <option value="active">{{ t('partners.statusActive') }}</option>
          <option value="inactive">{{ t('partners.statusInactive') }}</option>
          <option value="pending">{{ t('partners.statusPending') }}</option>
          <option value="suspended">{{ t('partners.statusSuspended') }}</option>
        </select>
      </FormGroup>

      <FormGroup :label="t('partners.description')">
        <textarea
          v-model="formData.description"
          :placeholder="t('partners.descriptionPlaceholder')"
          rows="4"
          class="form-textarea"
        />
      </FormGroup>
    </div>

    <!-- Contact Information -->
    <div class="form-section">
      <h3 class="section-title">{{ t('partners.contactInfo') }}</h3>

      <FormGroup :label="t('partners.contactPerson')">
        <input
          v-model="formData.contact_person"
          type="text"
          :placeholder="t('partners.contactPersonPlaceholder')"
          class="form-input"
        />
      </FormGroup>

      <div class="form-row">
        <FormGroup :label="t('partners.email')">
          <input
            v-model="formData.email"
            type="email"
            :placeholder="t('common.emailPlaceholder')"
            pattern="[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$"
            class="form-input"
          />
          <span v-if="formData.email && !isValidEmail(formData.email)" class="error-message">
            {{ t('common.invalidEmail') }}
          </span>
        </FormGroup>

        <FormGroup :label="t('partners.phone')">
          <input
            v-model="formData.phone"
            type="tel"
            :placeholder="t('common.phonePlaceholder')"
            class="form-input"
          />
        </FormGroup>
      </div>

      <FormGroup :label="t('partners.address')">
        <textarea
          v-model="formData.address"
          :placeholder="t('common.addressPlaceholder')"
          rows="3"
          class="form-textarea"
        />
      </FormGroup>

      <FormGroup :label="t('partners.website')">
        <input
          v-model="formData.website"
          type="url"
          :placeholder="t('partners.websitePlaceholder')"
          class="form-input"
        />
      </FormGroup>
    </div>

    <!-- Partnership Details -->
    <div class="form-section">
      <h3 class="section-title">{{ t('partners.partnershipDetails') }}</h3>

      <div class="form-row">
        <FormGroup :label="t('partners.partnershipStartDate')">
          <input
            v-model="formData.partnership_start_date"
            type="date"
            class="form-input"
          />
        </FormGroup>

        <FormGroup :label="t('partners.partnershipEndDate')">
          <input
            v-model="formData.partnership_end_date"
            type="date"
            class="form-input"
          />
        </FormGroup>
      </div>

      <FormGroup :label="t('partners.servicesProvided')">
        <textarea
          v-model="formData.services_provided"
          :placeholder="t('partners.servicesProvidedPlaceholder')"
          rows="3"
          class="form-textarea"
        />
      </FormGroup>

      <FormGroup :label="t('partners.notes')">
        <textarea
          v-model="formData.notes"
          :placeholder="t('partners.notesPlaceholder')"
          rows="4"
          class="form-textarea"
        />
      </FormGroup>
    </div>

    <!-- Form Actions -->
    <div class="form-actions">
      <BaseButton type="button" variant="secondary" @click="$emit('cancel')">
        {{ t('common.cancel') }}
      </BaseButton>
      <BaseButton type="submit" variant="primary" :loading="saving">
        {{ partner ? t('common.save') : t('partners.addPartner') }}
      </BaseButton>
    </div>
  </form>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import FormGroup from '../../components/base/FormGroup.vue'
import BaseButton from '../../components/base/BaseButton.vue'

const { t } = useI18n()

const props = defineProps({
  partner: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['submit', 'cancel'])

// Form data
const formData = ref({
  name: '',
  type: '',
  status: 'active',
  description: '',
  contact_person: '',
  email: '',
  phone: '',
  address: '',
  website: '',
  partnership_start_date: '',
  partnership_end_date: '',
  services_provided: '',
  notes: ''
})

const saving = ref(false)

// Watch for partner prop changes
watch(() => props.partner, (newPartner) => {
  if (newPartner) {
    formData.value = {
      name: newPartner.name || '',
      type: newPartner.type || '',
      status: newPartner.status || 'active',
      description: newPartner.description || '',
      contact_person: newPartner.contact_person || '',
      email: newPartner.email || '',
      phone: newPartner.phone || '',
      address: newPartner.address || '',
      website: newPartner.website || '',
      partnership_start_date: newPartner.partnership_start_date ? newPartner.partnership_start_date.split('T')[0] : '',
      partnership_end_date: newPartner.partnership_end_date ? newPartner.partnership_end_date.split('T')[0] : '',
      services_provided: newPartner.services_provided || '',
      notes: newPartner.notes || ''
    }
  }
}, { immediate: true })

// Methods
const isValidEmail = (email) => {
  const re = /^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$/i
  return re.test(email)
}

const handleSubmit = () => {
  // Validate email if provided
  if (formData.value.email && !isValidEmail(formData.value.email)) {
    return
  }

  emit('submit', formData.value)
}
</script>

<style scoped>
.partner-form {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border-color);
  padding-bottom: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.form-input,
.form-select,
.form-textarea {
  width: 100%;
  padding: 0.625rem;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  font-size: 0.95rem;
  background: var(--input-background);
  color: var(--text-primary);
  font-family: inherit;
}

.form-input:focus,
.form-select:focus,
.form-textarea:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 80px;
}

.error-message {
  color: var(--danger-color);
  font-size: 0.875rem;
  margin-top: 0.25rem;
  display: block;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
}

/* Responsive */
@media (max-width: 768px) {
  .form-row {
    grid-template-columns: 1fr;
  }

  .form-actions {
    flex-direction: column-reverse;
  }

  .form-actions button {
    width: 100%;
  }
}
</style>
