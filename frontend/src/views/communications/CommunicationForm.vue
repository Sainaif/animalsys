<template>
  <div class="communication-form">
    <h2 class="form-title">
      {{ communication?.id ? t('communications.editCommunication') : t('communications.newCommunication') }}
    </h2>

    <form @submit.prevent="handleSubmit">
      <!-- Message Information Section -->
      <div class="form-section">
        <h3 class="section-title">{{ t('communications.messageInformation') }}</h3>

        <FormGroup :label="t('communications.type')" required>
          <select
            v-model="formData.type"
            required
            class="form-control"
          >
            <option value="">{{ t('common.select') }}</option>
            <option value="email">{{ t('communications.types.email') }}</option>
            <option value="sms">{{ t('communications.types.sms') }}</option>
            <option value="newsletter">{{ t('communications.types.newsletter') }}</option>
            <option value="notification">{{ t('communications.types.notification') }}</option>
          </select>
        </FormGroup>

        <FormGroup
          v-if="formData.type !== 'sms'"
          :label="t('communications.subject')"
          :required="formData.type === 'email'"
        >
          <input
            v-model="formData.subject"
            type="text"
            class="form-control"
            :required="formData.type === 'email'"
            :placeholder="t('communications.subjectPlaceholder')"
          />
        </FormGroup>

        <FormGroup :label="t('communications.message')" required>
          <textarea
            v-model="formData.message"
            required
            class="form-control"
            rows="8"
            :placeholder="t('communications.messagePlaceholder')"
          ></textarea>
        </FormGroup>
      </div>

      <!-- Recipients Section -->
      <div v-if="!formData.is_template" class="form-section">
        <h3 class="section-title">{{ t('communications.recipients') }}</h3>

        <FormGroup :label="t('communications.recipientType')" required>
          <select
            v-model="formData.recipient_type"
            required
            class="form-control"
            @change="loadRecipients"
          >
            <option value="">{{ t('common.select') }}</option>
            <option value="volunteers">{{ t('communications.recipientTypes.volunteers') }}</option>
            <option value="donors">{{ t('communications.recipientTypes.donors') }}</option>
            <option value="adopters">{{ t('communications.recipientTypes.adopters') }}</option>
            <option value="partners">{{ t('communications.recipientTypes.partners') }}</option>
            <option value="custom">{{ t('communications.recipientTypes.custom') }}</option>
          </select>
        </FormGroup>

        <FormGroup
          v-if="formData.recipient_type && formData.recipient_type !== 'custom'"
          :label="t('communications.selectRecipients')"
          required
        >
          <div class="recipient-selection">
            <label class="checkbox-label">
              <input
                type="checkbox"
                v-model="selectAllRecipients"
                @change="toggleAllRecipients"
              />
              {{ t('communications.selectAll') }}
            </label>

            <div v-if="!selectAllRecipients" class="recipients-list">
              <label
                v-for="recipient in availableRecipients"
                :key="recipient.id"
                class="checkbox-label recipient-item"
              >
                <input
                  type="checkbox"
                  :value="recipient.id"
                  v-model="formData.recipient_ids"
                />
                <span>{{ recipient.name || recipient.email }}</span>
                <span class="recipient-contact">{{ recipient.email || recipient.phone }}</span>
              </label>
            </div>

            <div v-if="loading.recipients" class="loading-state">
              {{ t('common.loading') }}...
            </div>

            <div v-if="!loading.recipients && availableRecipients.length === 0" class="empty-state">
              {{ t('communications.noRecipientsAvailable') }}
            </div>
          </div>
        </FormGroup>

        <FormGroup
          v-if="formData.recipient_type === 'custom'"
          :label="t('communications.customRecipients')"
          required
        >
          <textarea
            v-model="formData.custom_recipients"
            required
            class="form-control"
            rows="4"
            :placeholder="t('communications.customRecipientsPlaceholder')"
          ></textarea>
          <small class="form-text">{{ t('communications.customRecipientsHelp') }}</small>
        </FormGroup>
      </div>

      <!-- Scheduling Section -->
      <div v-if="!formData.is_template" class="form-section">
        <h3 class="section-title">{{ t('communications.scheduling') }}</h3>

        <FormGroup>
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="scheduleForLater"
            />
            {{ t('communications.scheduleForLater') }}
          </label>
        </FormGroup>

        <FormGroup
          v-if="scheduleForLater"
          :label="t('communications.scheduledTime')"
          required
        >
          <input
            v-model="formData.scheduled_time"
            type="datetime-local"
            class="form-control"
            :required="scheduleForLater"
            :min="minScheduleTime"
          />
        </FormGroup>
      </div>

      <!-- Template Options Section -->
      <div class="form-section">
        <h3 class="section-title">{{ t('communications.templateOptions') }}</h3>

        <FormGroup>
          <label class="checkbox-label">
            <input
              type="checkbox"
              v-model="formData.is_template"
            />
            {{ t('communications.saveAsTemplate') }}
          </label>
        </FormGroup>

        <FormGroup
          v-if="formData.is_template"
          :label="t('communications.templateName')"
          required
        >
          <input
            v-model="formData.template_name"
            type="text"
            class="form-control"
            :required="formData.is_template"
            :placeholder="t('communications.templateNamePlaceholder')"
          />
        </FormGroup>

        <FormGroup
          v-if="formData.is_template"
          :label="t('communications.templateDescription')"
        >
          <textarea
            v-model="formData.template_description"
            class="form-control"
            rows="3"
            :placeholder="t('communications.templateDescriptionPlaceholder')"
          ></textarea>
        </FormGroup>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button type="button" @click="$emit('cancel')" class="btn btn-secondary">
          {{ t('common.cancel') }}
        </button>
        <button type="submit" class="btn btn-primary" :disabled="loading.submit">
          {{ loading.submit ? t('common.saving') : getSaveButtonText() }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useI18n } from 'vue-i18n'
import { API } from '@/api'
import FormGroup from '@/components/common/FormGroup.vue'

const { t } = useI18n()

const props = defineProps({
  communication: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['submit', 'cancel'])

const formData = reactive({
  type: '',
  subject: '',
  message: '',
  recipient_type: '',
  recipient_ids: [],
  custom_recipients: '',
  scheduled_time: '',
  is_template: false,
  template_name: '',
  template_description: '',
  status: 'draft'
})

const loading = reactive({
  submit: false,
  recipients: false
})

const scheduleForLater = ref(false)
const selectAllRecipients = ref(false)
const availableRecipients = ref([])

const minScheduleTime = computed(() => {
  const now = new Date()
  now.setMinutes(now.getMinutes() + 5) // Minimum 5 minutes from now
  return now.toISOString().slice(0, 16)
})

const getSaveButtonText = () => {
  if (formData.is_template) {
    return props.communication?.id ? t('common.save') : t('communications.saveTemplate')
  }
  if (scheduleForLater.value) {
    return t('communications.schedule')
  }
  return props.communication?.id ? t('common.save') : t('communications.sendNow')
}

const loadRecipients = async () => {
  if (!formData.recipient_type || formData.recipient_type === 'custom') {
    availableRecipients.value = []
    return
  }

  loading.recipients = true
  try {
    const response = await API.communications.getRecipients(formData.recipient_type)
    availableRecipients.value = response.data.data || []
  } catch (error) {
    console.error('Error loading recipients:', error)
    availableRecipients.value = []
  } finally {
    loading.recipients = false
  }
}

const toggleAllRecipients = () => {
  if (selectAllRecipients.value) {
    formData.recipient_ids = []
  } else {
    formData.recipient_ids = []
  }
}

watch(() => formData.is_template, (isTemplate) => {
  if (isTemplate) {
    // Clear recipient information when saving as template
    scheduleForLater.value = false
    formData.scheduled_time = ''
  }
})

const handleSubmit = async () => {
  loading.submit = true
  try {
    const submitData = {
      type: formData.type,
      subject: formData.subject,
      message: formData.message,
      is_template: formData.is_template
    }

    if (formData.is_template) {
      // Template submission
      submitData.template_name = formData.template_name
      submitData.template_description = formData.template_description
    } else {
      // Regular communication submission
      if (formData.recipient_type === 'custom') {
        submitData.custom_recipients = formData.custom_recipients
      } else {
        if (selectAllRecipients.value) {
          submitData.recipient_type = formData.recipient_type
          submitData.send_to_all = true
        } else {
          submitData.recipient_ids = formData.recipient_ids
        }
      }

      if (scheduleForLater.value) {
        submitData.scheduled_time = formData.scheduled_time
        submitData.status = 'scheduled'
      } else {
        submitData.status = 'draft'
      }
    }

    emit('submit', submitData)
  } catch (error) {
    console.error('Error submitting form:', error)
  } finally {
    loading.submit = false
  }
}

onMounted(() => {
  if (props.communication) {
    Object.assign(formData, {
      type: props.communication.type || '',
      subject: props.communication.subject || '',
      message: props.communication.message || '',
      recipient_type: props.communication.recipient_type || '',
      recipient_ids: props.communication.recipient_ids || [],
      custom_recipients: props.communication.custom_recipients || '',
      scheduled_time: props.communication.scheduled_time || '',
      is_template: props.communication.is_template || false,
      template_name: props.communication.template_name || '',
      template_description: props.communication.template_description || '',
      status: props.communication.status || 'draft'
    })

    if (props.communication.scheduled_time) {
      scheduleForLater.value = true
    }

    if (props.communication.recipient_type) {
      loadRecipients()
    }
  }
})
</script>

<style scoped>
.communication-form {
  max-width: 900px;
}

.form-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 2rem;
  color: var(--text-primary);
}

.form-section {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 1.25rem;
  color: var(--text-primary);
  border-bottom: 2px solid var(--primary);
  padding-bottom: 0.5rem;
}

.form-control {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 0.9375rem;
  background: var(--input-bg);
  color: var(--text-primary);
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-alpha);
}

textarea.form-control {
  resize: vertical;
  font-family: inherit;
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  padding: 0.5rem 0;
  font-size: 0.9375rem;
}

.checkbox-label input[type="checkbox"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.recipient-selection {
  border: 1px solid var(--border-color);
  border-radius: 6px;
  padding: 1rem;
  background: var(--input-bg);
}

.recipients-list {
  max-height: 300px;
  overflow-y: auto;
  margin-top: 1rem;
  border-top: 1px solid var(--border-color);
  padding-top: 1rem;
}

.recipient-item {
  padding: 0.75rem;
  border-bottom: 1px solid var(--border-color);
  display: flex;
  align-items: center;
  gap: 0.75rem;
  transition: background-color 0.2s;
}

.recipient-item:last-child {
  border-bottom: none;
}

.recipient-item:hover {
  background: var(--hover-bg);
}

.recipient-contact {
  margin-left: auto;
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.form-text {
  display: block;
  margin-top: 0.5rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.loading-state,
.empty-state {
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.9375rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
  margin-top: 2rem;
}

.btn {
  padding: 0.625rem 1.5rem;
  border-radius: 6px;
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-primary {
  background: var(--primary);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--background);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--hover-bg);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}
</style>
