<template>
  <form @submit.prevent="handleSubmit" class="schedule-form">
    <!-- Shift Information -->
    <div class="form-section">
      <h3 class="section-title">{{ t('schedules.shiftInfo') }}</h3>

      <FormGroup :label="t('schedules.date')" required>
        <input
          v-model="formData.date"
          type="date"
          required
          class="form-input"
        />
      </FormGroup>

      <div class="form-row">
        <FormGroup :label="t('schedules.startTime')" required>
          <input
            v-model="formData.start_time"
            type="time"
            required
            class="form-input"
          />
        </FormGroup>

        <FormGroup :label="t('schedules.endTime')" required>
          <input
            v-model="formData.end_time"
            type="time"
            required
            class="form-input"
          />
        </FormGroup>
      </div>

      <FormGroup :label="t('schedules.shiftType')" required>
        <select v-model="formData.shift_type" required class="form-select">
          <option value="">{{ t('schedules.selectShiftType') }}</option>
          <option value="morning">{{ t('schedules.typeMorning') }}</option>
          <option value="afternoon">{{ t('schedules.typeAfternoon') }}</option>
          <option value="evening">{{ t('schedules.typeEvening') }}</option>
          <option value="night">{{ t('schedules.typeNight') }}</option>
          <option value="full_day">{{ t('schedules.typeFullDay') }}</option>
        </select>
      </FormGroup>

      <FormGroup :label="t('schedules.role')" required>
        <input
          v-model="formData.role"
          type="text"
          :placeholder="t('schedules.rolePlaceholder')"
          required
          class="form-input"
        />
      </FormGroup>

      <FormGroup :label="t('schedules.description')">
        <textarea
          v-model="formData.description"
          :placeholder="t('schedules.descriptionPlaceholder')"
          rows="3"
          class="form-textarea"
        />
      </FormGroup>
    </div>

    <!-- Volunteer Assignment -->
    <div class="form-section">
      <h3 class="section-title">{{ t('schedules.volunteerAssignment') }}</h3>

      <FormGroup :label="t('schedules.assignVolunteer')">
        <select v-model="formData.assigned_volunteer_id" class="form-select">
          <option value="">{{ t('schedules.noVolunteerAssigned') }}</option>
          <option v-for="volunteer in volunteers" :key="volunteer.id" :value="volunteer.id">
            {{ volunteer.first_name }} {{ volunteer.last_name }}
          </option>
        </select>
      </FormGroup>

      <FormGroup :label="t('common.status')" required>
        <select v-model="formData.status" required class="form-select">
          <option value="open">{{ t('schedules.statusOpen') }}</option>
          <option value="filled">{{ t('schedules.statusFilled') }}</option>
          <option value="completed">{{ t('schedules.statusCompleted') }}</option>
          <option value="cancelled">{{ t('schedules.statusCancelled') }}</option>
        </select>
      </FormGroup>
    </div>

    <!-- Additional Information -->
    <div class="form-section">
      <h3 class="section-title">{{ t('schedules.additionalInfo') }}</h3>

      <FormGroup :label="t('schedules.location')">
        <input
          v-model="formData.location"
          type="text"
          :placeholder="t('schedules.locationPlaceholder')"
          class="form-input"
        />
      </FormGroup>

      <FormGroup :label="t('schedules.requiredSkills')">
        <input
          v-model="formData.required_skills"
          type="text"
          :placeholder="t('schedules.requiredSkillsPlaceholder')"
          class="form-input"
        />
      </FormGroup>

      <FormGroup :label="t('schedules.notes')">
        <textarea
          v-model="formData.notes"
          :placeholder="t('schedules.notesPlaceholder')"
          rows="3"
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
        {{ schedule ? t('common.save') : t('schedules.addShift') }}
      </BaseButton>
    </div>
  </form>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { API } from '../../api'
import FormGroup from '../../components/base/FormGroup.vue'
import BaseButton from '../../components/base/BaseButton.vue'

const { t } = useI18n()

const props = defineProps({
  schedule: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['submit', 'cancel'])

// Form data
const formData = ref({
  date: '',
  start_time: '',
  end_time: '',
  shift_type: '',
  role: '',
  description: '',
  assigned_volunteer_id: '',
  status: 'open',
  location: '',
  required_skills: '',
  notes: ''
})

const volunteers = ref([])
const saving = ref(false)

// Watch for schedule prop changes
watch(() => props.schedule, (newSchedule) => {
  if (newSchedule) {
    formData.value = {
      date: newSchedule.date || '',
      start_time: newSchedule.start_time ? formatTimeForInput(newSchedule.start_time) : '',
      end_time: newSchedule.end_time ? formatTimeForInput(newSchedule.end_time) : '',
      shift_type: newSchedule.shift_type || '',
      role: newSchedule.role || '',
      description: newSchedule.description || '',
      assigned_volunteer_id: newSchedule.assigned_volunteer_id || '',
      status: newSchedule.status || 'open',
      location: newSchedule.location || '',
      required_skills: newSchedule.required_skills || '',
      notes: newSchedule.notes || ''
    }
  }
}, { immediate: true })

// Methods
const formatTimeForInput = (timeString) => {
  if (!timeString) return ''
  // If it's a full datetime, extract time part
  if (timeString.includes('T')) {
    return timeString.split('T')[1].substring(0, 5)
  }
  // If it's already in HH:MM format
  if (timeString.length === 5) {
    return timeString
  }
  // If it's in HH:MM:SS format
  return timeString.substring(0, 5)
}

const fetchVolunteers = async () => {
  try {
    const response = await API.volunteers.list({ status: 'active' })
    volunteers.value = response.data.data || []
  } catch (error) {
    console.error('Error fetching volunteers:', error)
  }
}

const handleSubmit = () => {
  emit('submit', formData.value)
}

// Lifecycle
onMounted(() => {
  fetchVolunteers()
})
</script>

<style scoped>
.schedule-form {
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
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
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
