<template>
  <div class="schedules-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.schedules') }}</h1>
      <BaseButton variant="primary" @click="openCreateModal">
        ➕ {{ t('schedules.addShift') }}
      </BaseButton>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">📅</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.totalSchedules }}</div>
            <div class="stat-label">{{ t('schedules.totalShifts') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">✅</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.filledSchedules }}</div>
            <div class="stat-label">{{ t('schedules.filledShifts') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">⚠️</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.openSchedules }}</div>
            <div class="stat-label">{{ t('schedules.openShifts') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">🔄</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.swapRequests }}</div>
            <div class="stat-label">{{ t('schedules.swapRequests') }}</div>
          </div>
        </div>
      </BaseCard>
    </div>

    <!-- View Toggle and Filters -->
    <BaseCard>
      <div class="controls-bar">
        <div class="view-toggle">
          <button
            :class="['toggle-btn', { active: currentView === 'calendar' }]"
            @click="currentView = 'calendar'"
          >
            📅 {{ t('schedules.calendarView') }}
          </button>
          <button
            :class="['toggle-btn', { active: currentView === 'list' }]"
            @click="currentView = 'list'"
          >
            📋 {{ t('schedules.listView') }}
          </button>
        </div>

        <div class="filters">
          <div class="filter-group">
            <select v-model="filters.shiftType" class="filter-select">
              <option value="">{{ t('schedules.allShiftTypes') }}</option>
              <option value="morning">{{ t('schedules.typeMorning') }}</option>
              <option value="afternoon">{{ t('schedules.typeAfternoon') }}</option>
              <option value="evening">{{ t('schedules.typeEvening') }}</option>
              <option value="night">{{ t('schedules.typeNight') }}</option>
              <option value="full_day">{{ t('schedules.typeFullDay') }}</option>
            </select>
          </div>

          <div class="filter-group">
            <select v-model="filters.status" class="filter-select">
              <option value="">{{ t('schedules.allStatuses') }}</option>
              <option value="open">{{ t('schedules.statusOpen') }}</option>
              <option value="filled">{{ t('schedules.statusFilled') }}</option>
              <option value="completed">{{ t('schedules.statusCompleted') }}</option>
              <option value="cancelled">{{ t('schedules.statusCancelled') }}</option>
            </select>
          </div>
        </div>
      </div>
    </BaseCard>

    <!-- Calendar View -->
    <BaseCard v-if="currentView === 'calendar'">
      <div class="calendar-view">
        <div class="calendar-header">
          <BaseButton size="small" variant="secondary" @click="previousWeek">
            ← {{ t('schedules.previousWeek') }}
          </BaseButton>
          <h3 class="calendar-title">
            {{ formatWeekRange(currentWeekStart) }}
          </h3>
          <BaseButton size="small" variant="secondary" @click="nextWeek">
            {{ t('schedules.nextWeek') }} →
          </BaseButton>
        </div>

        <div class="calendar-grid">
          <div
            v-for="day in weekDays"
            :key="day.date"
            class="calendar-day"
          >
            <div class="day-header">
              <span class="day-name">{{ day.name }}</span>
              <span class="day-date">{{ formatDate(day.date) }}</span>
            </div>

            <div class="day-schedules">
              <div
                v-for="schedule in getSchedulesForDay(day.date)"
                :key="schedule.id"
                :class="['schedule-item', `status-${schedule.status}`, `type-${schedule.shift_type}`]"
                @click="openViewModal(schedule)"
              >
                <div class="schedule-time">
                  {{ formatTime(schedule.start_time) }} - {{ formatTime(schedule.end_time) }}
                </div>
                <div class="schedule-info">
                  <span class="schedule-role">{{ schedule.role }}</span>
                  <span v-if="schedule.assigned_volunteer" class="schedule-volunteer">
                    👤 {{ schedule.assigned_volunteer.first_name }} {{ schedule.assigned_volunteer.last_name }}
                  </span>
                  <span v-else class="schedule-open">{{ t('schedules.noVolunteer') }}</span>
                </div>
              </div>

              <div v-if="getSchedulesForDay(day.date).length === 0" class="no-schedules">
                {{ t('schedules.noShifts') }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </BaseCard>

    <!-- List View -->
    <BaseCard v-if="currentView === 'list'">
      <DataTable
        :columns="columns"
        :data="filteredSchedules"
        :loading="loading"
        @sort="handleSort"
      >
        <template #cell-date="{ row }">
          {{ formatFullDate(row.date) }}
        </template>

        <template #cell-time="{ row }">
          {{ formatTime(row.start_time) }} - {{ formatTime(row.end_time) }}
        </template>

        <template #cell-shift_type="{ row }">
          <span :class="['badge', 'type-badge', `type-${row.shift_type}`]">
            {{ t(`schedules.type${capitalize(row.shift_type)}`) }}
          </span>
        </template>

        <template #cell-status="{ row }">
          <span :class="['badge', `status-${row.status}`]">
            {{ t(`schedules.status${capitalize(row.status)}`) }}
          </span>
        </template>

        <template #cell-volunteer="{ row }">
          <span v-if="row.assigned_volunteer">
            👤 {{ row.assigned_volunteer.first_name }} {{ row.assigned_volunteer.last_name }}
          </span>
          <span v-else class="text-muted">{{ t('schedules.noVolunteer') }}</span>
        </template>

        <template #cell-actions="{ row }">
          <div class="actions">
            <BaseButton size="small" variant="secondary" @click="openViewModal(row)">
              {{ t('common.view') }}
            </BaseButton>
            <BaseButton size="small" variant="secondary" @click="openEditModal(row)">
              {{ t('common.edit') }}
            </BaseButton>
            <BaseButton size="small" variant="danger" @click="confirmDelete(row)">
              {{ t('common.delete') }}
            </BaseButton>
          </div>
        </template>
      </DataTable>
    </BaseCard>

    <!-- Create/Edit Modal -->
    <BaseModal
      v-if="showModal"
      :title="editingSchedule ? t('schedules.editShift') : t('schedules.addShift')"
      size="large"
      @close="closeModal"
    >
      <ScheduleForm
        :schedule="editingSchedule"
        @submit="handleSubmit"
        @cancel="closeModal"
      />
    </BaseModal>

    <!-- View Modal -->
    <BaseModal
      v-if="showViewModal"
      :title="t('schedules.shiftDetails')"
      size="medium"
      @close="showViewModal = false"
    >
      <ScheduleView
        :schedule="viewingSchedule"
        @close="showViewModal = false"
        @refresh="fetchSchedules"
      />
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      :title="t('schedules.deleteShift')"
      size="small"
      @close="showDeleteModal = false"
    >
      <p>{{ t('schedules.deleteShiftConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteSchedule" :loading="deleting">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { API } from '../../api'
import { useNotificationStore } from '../../stores/notification'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import DataTable from '../../components/base/DataTable.vue'
import ScheduleForm from './ScheduleForm.vue'
import ScheduleView from './ScheduleView.vue'

const { t } = useI18n()
const notificationStore = useNotificationStore()

// State
const schedules = ref([])
const statistics = ref({
  totalSchedules: 0,
  filledSchedules: 0,
  openSchedules: 0,
  swapRequests: 0
})
const loading = ref(false)
const currentView = ref('calendar') // 'calendar' or 'list'
const currentWeekStart = ref(getStartOfWeek(new Date()))
const filters = ref({
  shiftType: '',
  status: ''
})

// Modal state
const showModal = ref(false)
const showViewModal = ref(false)
const showDeleteModal = ref(false)
const editingSchedule = ref(null)
const viewingSchedule = ref(null)
const scheduleToDelete = ref(null)
const deleting = ref(false)

// Table columns
const columns = [
  { key: 'date', label: t('schedules.date'), sortable: true },
  { key: 'time', label: t('schedules.time'), sortable: false },
  { key: 'shift_type', label: t('schedules.shiftType'), sortable: true },
  { key: 'role', label: t('schedules.role'), sortable: true },
  { key: 'volunteer', label: t('schedules.volunteer'), sortable: false },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'actions', label: t('common.actions'), sortable: false }
]

// Computed
const weekDays = computed(() => {
  const days = []
  const start = new Date(currentWeekStart.value)

  for (let i = 0; i < 7; i++) {
    const date = new Date(start)
    date.setDate(start.getDate() + i)
    days.push({
      date: date.toISOString().split('T')[0],
      name: date.toLocaleDateString('pl-PL', { weekday: 'long' })
    })
  }

  return days
})

const filteredSchedules = computed(() => {
  let result = [...schedules.value]

  // Filter by shift type
  if (filters.value.shiftType) {
    result = result.filter(s => s.shift_type === filters.value.shiftType)
  }

  // Filter by status
  if (filters.value.status) {
    result = result.filter(s => s.status === filters.value.status)
  }

  return result
})

// Methods
function getStartOfWeek(date) {
  const d = new Date(date)
  const day = d.getDay()
  const diff = d.getDate() - day + (day === 0 ? -6 : 1) // Adjust for Monday start
  return new Date(d.setDate(diff))
}

const fetchSchedules = async () => {
  loading.value = true
  try {
    const response = await API.schedules.list()
    schedules.value = response.data.data || []
  } catch (error) {
    notificationStore.error(t('schedules.fetchError'))
    console.error('Error fetching schedules:', error)
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    const response = await API.schedules.getStatistics()
    statistics.value = response.data.data || statistics.value
  } catch (error) {
    console.error('Error fetching statistics:', error)
  }
}

const getSchedulesForDay = (date) => {
  return filteredSchedules.value.filter(s => s.date === date)
}

const previousWeek = () => {
  const newStart = new Date(currentWeekStart.value)
  newStart.setDate(newStart.getDate() - 7)
  currentWeekStart.value = newStart
}

const nextWeek = () => {
  const newStart = new Date(currentWeekStart.value)
  newStart.setDate(newStart.getDate() + 7)
  currentWeekStart.value = newStart
}

const formatWeekRange = (startDate) => {
  const start = new Date(startDate)
  const end = new Date(startDate)
  end.setDate(end.getDate() + 6)

  return `${start.toLocaleDateString('pl-PL')} - ${end.toLocaleDateString('pl-PL')}`
}

const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('pl-PL', { day: '2-digit', month: '2-digit' })
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
  // Handle both full datetime and time-only formats
  const date = timeString.includes('T') ? new Date(timeString) : new Date(`2000-01-01T${timeString}`)
  return date.toLocaleTimeString('pl-PL', { hour: '2-digit', minute: '2-digit' })
}

const openCreateModal = () => {
  editingSchedule.value = null
  showModal.value = true
}

const openEditModal = (schedule) => {
  editingSchedule.value = schedule
  showModal.value = true
}

const openViewModal = (schedule) => {
  viewingSchedule.value = schedule
  showViewModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingSchedule.value = null
}

const handleSubmit = async (scheduleData) => {
  try {
    if (editingSchedule.value) {
      await API.schedules.update(editingSchedule.value.id, scheduleData)
      notificationStore.success(t('schedules.updateSuccess'))
    } else {
      await API.schedules.create(scheduleData)
      notificationStore.success(t('schedules.createSuccess'))
    }
    closeModal()
    fetchSchedules()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(
      editingSchedule.value ? t('schedules.updateError') : t('schedules.createError')
    )
    console.error('Error saving schedule:', error)
  }
}

const confirmDelete = (schedule) => {
  scheduleToDelete.value = schedule
  showDeleteModal.value = true
}

const deleteSchedule = async () => {
  deleting.value = true
  try {
    await API.schedules.delete(scheduleToDelete.value.id)
    notificationStore.success(t('schedules.deleteSuccess'))
    showDeleteModal.value = false
    scheduleToDelete.value = null
    fetchSchedules()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(t('schedules.deleteError'))
    console.error('Error deleting schedule:', error)
  } finally {
    deleting.value = false
  }
}

const handleSort = (column, direction) => {
  schedules.value.sort((a, b) => {
    const aVal = a[column.key]
    const bVal = b[column.key]
    if (direction === 'asc') {
      return aVal > bVal ? 1 : -1
    } else {
      return aVal < bVal ? 1 : -1
    }
  })
}

const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1).replace('_', '')
}

// Lifecycle
onMounted(() => {
  fetchSchedules()
  fetchStatistics()
})
</script>

<style scoped>
.schedules-page {
  max-width: 1600px;
  padding: 2rem;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin: 0;
}

/* Statistics Grid */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  padding: 1.5rem;
}

.stat-content {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-details {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
  color: var(--primary-color);
}

.stat-label {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

/* Controls Bar */
.controls-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
  padding: 1rem;
  flex-wrap: wrap;
}

.view-toggle {
  display: flex;
  gap: 0.5rem;
}

.toggle-btn {
  padding: 0.5rem 1rem;
  border: 1px solid var(--border-color);
  background: var(--background-primary);
  color: var(--text-primary);
  border-radius: 0.375rem;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.toggle-btn:hover {
  background: var(--background-secondary);
}

.toggle-btn.active {
  background: var(--primary-color);
  color: white;
  border-color: var(--primary-color);
}

.filters {
  display: flex;
  gap: 1rem;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.filter-select {
  padding: 0.5rem;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  font-size: 0.9rem;
  background: var(--input-background);
  color: var(--text-primary);
}

/* Calendar View */
.calendar-view {
  padding: 1rem;
}

.calendar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.calendar-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
}

.calendar-grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: 1rem;
}

.calendar-day {
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  min-height: 200px;
  display: flex;
  flex-direction: column;
}

.day-header {
  padding: 0.75rem;
  background: var(--background-secondary);
  border-bottom: 1px solid var(--border-color);
  border-radius: 0.5rem 0.5rem 0 0;
}

.day-name {
  display: block;
  font-weight: 600;
  font-size: 0.9rem;
  text-transform: capitalize;
}

.day-date {
  display: block;
  font-size: 0.85rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.day-schedules {
  padding: 0.5rem;
  flex: 1;
  overflow-y: auto;
}

.schedule-item {
  padding: 0.5rem;
  margin-bottom: 0.5rem;
  border-radius: 0.375rem;
  cursor: pointer;
  transition: transform 0.2s;
  border-left: 3px solid;
}

.schedule-item:hover {
  transform: translateX(2px);
}

.schedule-item.status-open {
  background: #fff3cd;
  border-left-color: #ffc107;
}

.schedule-item.status-filled {
  background: #d4edda;
  border-left-color: #28a745;
}

.schedule-item.status-completed {
  background: #d1ecf1;
  border-left-color: #17a2b8;
}

.schedule-item.status-cancelled {
  background: #f8d7da;
  border-left-color: #dc3545;
}

.schedule-time {
  font-size: 0.85rem;
  font-weight: 600;
  margin-bottom: 0.25rem;
}

.schedule-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.8rem;
}

.schedule-role {
  font-weight: 500;
}

.schedule-volunteer {
  color: var(--text-secondary);
}

.schedule-open {
  color: #856404;
  font-style: italic;
}

.no-schedules {
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.9rem;
}

/* Badges */
.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.85rem;
  font-weight: 500;
  display: inline-block;
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

/* Actions */
.actions {
  display: flex;
  gap: 0.5rem;
}

.text-muted {
  color: var(--text-secondary);
  font-style: italic;
}

/* Responsive */
@media (max-width: 1200px) {
  .calendar-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

@media (max-width: 768px) {
  .calendar-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .controls-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .view-toggle,
  .filters {
    width: 100%;
  }

  .toggle-btn {
    flex: 1;
  }
}
</style>
