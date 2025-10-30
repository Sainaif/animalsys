<template>
  <div class="schedule-view">
    <div v-if="schedule" class="schedule-content">
      <!-- Shift Information -->
      <div class="info-section">
        <h3 class="section-title">{{ t('schedules.shiftInfo') }}</h3>

        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('schedules.date') }}</span>
            <span class="info-value">{{ formatFullDate(schedule.date) }}</span>
          </div>

          <div class="info-item">
            <span class="info-label">{{ t('schedules.time') }}</span>
            <span class="info-value">
              {{ formatTime(schedule.start_time) }} - {{ formatTime(schedule.end_time) }}
            </span>
          </div>

          <div class="info-item">
            <span class="info-label">{{ t('schedules.shiftType') }}</span>
            <span :class="['badge', 'type-badge', `type-${schedule.shift_type}`]">
              {{ t(`schedules.type${capitalize(schedule.shift_type)}`) }}
            </span>
          </div>

          <div class="info-item">
            <span class="info-label">{{ t('common.status') }}</span>
            <span :class="['badge', `status-${schedule.status}`]">
              {{ t(`schedules.status${capitalize(schedule.status)}`) }}
            </span>
          </div>

          <div class="info-item">
            <span class="info-label">{{ t('schedules.role') }}</span>
            <span class="info-value">{{ schedule.role }}</span>
          </div>

          <div v-if="schedule.location" class="info-item">
            <span class="info-label">{{ t('schedules.location') }}</span>
            <span class="info-value">📍 {{ schedule.location }}</span>
          </div>

          <div v-if="schedule.required_skills" class="info-item full-width">
            <span class="info-label">{{ t('schedules.requiredSkills') }}</span>
            <span class="info-value">{{ schedule.required_skills }}</span>
          </div>

          <div v-if="schedule.description" class="info-item full-width">
            <span class="info-label">{{ t('schedules.description') }}</span>
            <span class="info-value">{{ schedule.description }}</span>
          </div>
        </div>
      </div>

      <!-- Volunteer Assignment -->
      <div class="info-section">
        <h3 class="section-title">{{ t('schedules.volunteerAssignment') }}</h3>

        <div v-if="schedule.assigned_volunteer" class="volunteer-card">
          <div class="volunteer-info">
            <div class="volunteer-avatar">👤</div>
            <div class="volunteer-details">
              <div class="volunteer-name">
                {{ schedule.assigned_volunteer.first_name }} {{ schedule.assigned_volunteer.last_name }}
              </div>
              <div v-if="schedule.assigned_volunteer.email" class="volunteer-contact">
                ✉️ {{ schedule.assigned_volunteer.email }}
              </div>
              <div v-if="schedule.assigned_volunteer.phone" class="volunteer-contact">
                📞 {{ schedule.assigned_volunteer.phone }}
              </div>
            </div>
          </div>

          <BaseButton
            v-if="schedule.status !== 'completed' && schedule.status !== 'cancelled'"
            size="small"
            variant="secondary"
            @click="showUnassignModal = true"
          >
            {{ t('schedules.unassignVolunteer') }}
          </BaseButton>
        </div>

        <div v-else class="empty-volunteer">
          <p>{{ t('schedules.noVolunteerAssigned') }}</p>
          <BaseButton
            v-if="schedule.status !== 'completed' && schedule.status !== 'cancelled'"
            size="small"
            variant="primary"
            @click="showAssignModal = true"
          >
            {{ t('schedules.assignVolunteer') }}
          </BaseButton>
        </div>
      </div>

      <!-- Notes -->
      <div v-if="schedule.notes" class="info-section">
        <h3 class="section-title">{{ t('schedules.notes') }}</h3>
        <div class="notes-content">
          {{ schedule.notes }}
        </div>
      </div>

      <!-- Actions -->
      <div class="actions-section">
        <BaseButton variant="secondary" @click="$emit('close')">
          {{ t('common.close') }}
        </BaseButton>
      </div>
    </div>

    <!-- Assign Volunteer Modal -->
    <BaseModal
      v-if="showAssignModal"
      :title="t('schedules.assignVolunteer')"
      size="small"
      @close="showAssignModal = false"
    >
      <div class="assign-form">
        <FormGroup :label="t('schedules.selectVolunteer')" required>
          <select v-model="selectedVolunteerId" class="form-select">
            <option value="">{{ t('schedules.chooseVolunteer') }}</option>
            <option v-for="volunteer in volunteers" :key="volunteer.id" :value="volunteer.id">
              {{ volunteer.first_name }} {{ volunteer.last_name }}
            </option>
          </select>
        </FormGroup>
      </div>

      <template #footer>
        <BaseButton variant="secondary" @click="showAssignModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton
          variant="primary"
          @click="assignVolunteer"
          :loading="assigning"
          :disabled="!selectedVolunteerId"
        >
          {{ t('schedules.assign') }}
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Unassign Confirmation Modal -->
    <BaseModal
      v-if="showUnassignModal"
      :title="t('schedules.unassignVolunteer')"
      size="small"
      @close="showUnassignModal = false"
    >
      <p>{{ t('schedules.unassignVolunteerConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showUnassignModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="unassignVolunteer" :loading="unassigning">
          {{ t('schedules.unassign') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { API } from '../../api'
import { useNotificationStore } from '../../stores/notification'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import FormGroup from '../../components/base/FormGroup.vue'

const { t } = useI18n()
const notificationStore = useNotificationStore()

const props = defineProps({
  schedule: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['close', 'refresh'])

// State
const volunteers = ref([])
const selectedVolunteerId = ref('')
const showAssignModal = ref(false)
const showUnassignModal = ref(false)
const assigning = ref(false)
const unassigning = ref(false)

// Methods
const fetchVolunteers = async () => {
  try {
    const response = await API.volunteers.list({ status: 'active' })
    volunteers.value = response.data.data || []
  } catch (error) {
    console.error('Error fetching volunteers:', error)
  }
}

const assignVolunteer = async () => {
  if (!selectedVolunteerId.value) return

  assigning.value = true
  try {
    await API.schedules.assignVolunteer(props.schedule.id, selectedVolunteerId.value)
    notificationStore.success(t('schedules.assignSuccess'))
    showAssignModal.value = false
    selectedVolunteerId.value = ''
    emit('refresh')
    emit('close')
  } catch (error) {
    notificationStore.error(t('schedules.assignError'))
    console.error('Error assigning volunteer:', error)
  } finally {
    assigning.value = false
  }
}

const unassignVolunteer = async () => {
  unassigning.value = true
  try {
    await API.schedules.unassignVolunteer(
      props.schedule.id,
      props.schedule.assigned_volunteer_id
    )
    notificationStore.success(t('schedules.unassignSuccess'))
    showUnassignModal.value = false
    emit('refresh')
    emit('close')
  } catch (error) {
    notificationStore.error(t('schedules.unassignError'))
    console.error('Error unassigning volunteer:', error)
  } finally {
    unassigning.value = false
  }
}

const formatFullDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('pl-PL', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const formatTime = (timeString) => {
  if (!timeString) return ''
  const date = timeString.includes('T') ? new Date(timeString) : new Date(`2000-01-01T${timeString}`)
  return date.toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' })
}

const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1).replace('_', '')
}

// Lifecycle
onMounted(() => {
  fetchVolunteers()
})
</script>

<style scoped>
.schedule-view {
  padding: 1rem;
}

.schedule-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Info Section */
.info-section {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.section-title {
  font-size: 1.1rem;
  font-weight: 600;
  margin: 0;
  color: var(--text-primary);
  border-bottom: 2px solid var(--border-color);
  padding-bottom: 0.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.info-item.full-width {
  grid-column: 1 / -1;
}

.info-label {
  font-size: 0.875rem;
  font-weight: 500;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.025em;
}

.info-value {
  font-size: 1rem;
  color: var(--text-primary);
}

/* Badges */
.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.85rem;
  font-weight: 500;
  display: inline-block;
  width: fit-content;
}

.type-badge.type-morning {
  background: #fff3cd;
  color: #856404;
}

.type-badge.type-afternoon {
  background: #cfe2ff;
  color: #084298;
}

.type-badge.type-evening {
  background: #e7d4f8;
  color: #6f42c1;
}

.type-badge.type-night {
  background: #d3d3d3;
  color: #495057;
}

.type-badge.type-full_day,
.type-badge.type-fullday {
  background: #d1ecf1;
  color: #0c5460;
}

.status-open {
  background: #fff3cd;
  color: #856404;
}

.status-filled {
  background: #d4edda;
  color: #155724;
}

.status-completed {
  background: #d1ecf1;
  color: #0c5460;
}

.status-cancelled {
  background: #f8d7da;
  color: #721c24;
}

/* Volunteer Card */
.volunteer-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background: var(--background-secondary);
}

.volunteer-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.volunteer-avatar {
  width: 3rem;
  height: 3rem;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--primary-color);
  color: white;
  border-radius: 50%;
  font-size: 1.5rem;
}

.volunteer-details {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.volunteer-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.volunteer-contact {
  font-size: 0.9rem;
  color: var(--text-secondary);
}

.empty-volunteer {
  padding: 2rem;
  text-align: center;
  border: 2px dashed var(--border-color);
  border-radius: 0.5rem;
  background: var(--background-secondary);
}

.empty-volunteer p {
  margin: 0 0 1rem 0;
  color: var(--text-secondary);
}

/* Notes */
.notes-content {
  white-space: pre-wrap;
  color: var(--text-secondary);
  line-height: 1.6;
  padding: 1rem;
  background: var(--background-secondary);
  border-radius: 0.375rem;
}

/* Actions */
.actions-section {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid var(--border-color);
}

/* Assign Form */
.assign-form {
  padding: 1rem 0;
}

.form-select {
  width: 100%;
  padding: 0.625rem;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  font-size: 0.95rem;
  background: var(--input-background);
  color: var(--text-primary);
}

.form-select:focus {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(66, 153, 225, 0.1);
}

/* Responsive */
@media (max-width: 768px) {
  .info-grid {
    grid-template-columns: 1fr;
  }

  .volunteer-card {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }
}
</style>
