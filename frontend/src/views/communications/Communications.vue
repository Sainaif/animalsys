<template>
  <div class="communications-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.communications') }}</h1>
      <BaseButton variant="primary" @click="openCreateModal">
        ➕ {{ t('communications.newMessage') }}
      </BaseButton>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">✉️</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.totalSent }}</div>
            <div class="stat-label">{{ t('communications.totalSent') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">📅</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.scheduled }}</div>
            <div class="stat-label">{{ t('communications.scheduled') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">✅</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.delivered }}</div>
            <div class="stat-label">{{ t('communications.delivered') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">❌</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.failed }}</div>
            <div class="stat-label">{{ t('communications.failed') }}</div>
          </div>
        </div>
      </BaseCard>
    </div>

    <!-- Filters and Tabs -->
    <BaseCard>
      <div class="controls-bar">
        <div class="view-tabs">
          <button
            :class="['tab-btn', { active: currentTab === 'all' }]"
            @click="currentTab = 'all'"
          >
            📋 {{ t('communications.all') }}
          </button>
          <button
            :class="['tab-btn', { active: currentTab === 'scheduled' }]"
            @click="currentTab = 'scheduled'"
          >
            📅 {{ t('communications.scheduled') }}
          </button>
          <button
            :class="['tab-btn', { active: currentTab === 'templates' }]"
            @click="currentTab = 'templates'"
          >
            📝 {{ t('communications.templates') }}
          </button>
        </div>

        <div class="filters">
          <div class="filter-group">
            <select v-model="filters.type" class="filter-select">
              <option value="">{{ t('communications.allTypes') }}</option>
              <option value="email">{{ t('communications.typeEmail') }}</option>
              <option value="sms">{{ t('communications.typeSms') }}</option>
              <option value="newsletter">{{ t('communications.typeNewsletter') }}</option>
              <option value="notification">{{ t('communications.typeNotification') }}</option>
            </select>
          </div>

          <div class="filter-group" v-if="currentTab !== 'templates'">
            <select v-model="filters.status" class="filter-select">
              <option value="">{{ t('communications.allStatuses') }}</option>
              <option value="draft">{{ t('communications.statusDraft') }}</option>
              <option value="scheduled">{{ t('communications.statusScheduled') }}</option>
              <option value="sent">{{ t('communications.statusSent') }}</option>
              <option value="failed">{{ t('communications.statusFailed') }}</option>
            </select>
          </div>
        </div>
      </div>
    </BaseCard>

    <!-- Communications List -->
    <BaseCard v-if="currentTab !== 'templates'">
      <DataTable
        :columns="communicationColumns"
        :data="filteredCommunications"
        :loading="loading"
        @sort="handleSort"
      >
        <template #cell-subject="{ row }">
          <div class="subject-cell">
            <span class="subject-text">{{ row.subject }}</span>
            <span v-if="row.is_bulk" class="bulk-badge">📢 {{ t('communications.bulk') }}</span>
          </div>
        </template>

        <template #cell-type="{ row }">
          <span :class="['badge', 'type-badge', `type-${row.type}`]">
            {{ t(`communications.type${capitalize(row.type)}`) }}
          </span>
        </template>

        <template #cell-status="{ row }">
          <span :class="['badge', `status-${row.status}`]">
            {{ t(`communications.status${capitalize(row.status)}`) }}
          </span>
        </template>

        <template #cell-recipients="{ row }">
          <span class="recipients-count">
            {{ row.recipient_count || 0 }} {{ t('communications.recipients') }}
          </span>
        </template>

        <template #cell-scheduled_time="{ row }">
          {{ row.scheduled_time ? formatDateTime(row.scheduled_time) : '-' }}
        </template>

        <template #cell-sent_at="{ row }">
          {{ row.sent_at ? formatDateTime(row.sent_at) : '-' }}
        </template>

        <template #cell-actions="{ row }">
          <div class="actions">
            <BaseButton size="small" variant="secondary" @click="openViewModal(row)">
              {{ t('common.view') }}
            </BaseButton>
            <BaseButton
              v-if="row.status === 'draft'"
              size="small"
              variant="primary"
              @click="confirmSend(row)"
            >
              {{ t('communications.send') }}
            </BaseButton>
            <BaseButton
              v-if="row.status === 'scheduled'"
              size="small"
              variant="secondary"
              @click="cancelSchedule(row)"
            >
              {{ t('communications.cancel') }}
            </BaseButton>
            <BaseButton size="small" variant="danger" @click="confirmDelete(row)">
              {{ t('common.delete') }}
            </BaseButton>
          </div>
        </template>
      </DataTable>
    </BaseCard>

    <!-- Templates List -->
    <BaseCard v-if="currentTab === 'templates'">
      <div class="templates-header">
        <h3>{{ t('communications.templates') }}</h3>
        <BaseButton size="small" variant="primary" @click="openTemplateModal">
          ➕ {{ t('communications.newTemplate') }}
        </BaseButton>
      </div>

      <div class="templates-grid">
        <div
          v-for="template in filteredTemplates"
          :key="template.id"
          class="template-card"
        >
          <div class="template-header">
            <h4 class="template-name">{{ template.name }}</h4>
            <span :class="['badge', 'type-badge', `type-${template.type}`]">
              {{ t(`communications.type${capitalize(template.type)}`) }}
            </span>
          </div>

          <p class="template-description">{{ template.description || t('communications.noDescription') }}</p>

          <div class="template-actions">
            <BaseButton size="small" variant="secondary" @click="useTemplate(template)">
              {{ t('communications.useTemplate') }}
            </BaseButton>
            <BaseButton size="small" variant="secondary" @click="editTemplate(template)">
              {{ t('common.edit') }}
            </BaseButton>
            <BaseButton size="small" variant="danger" @click="confirmDeleteTemplate(template)">
              {{ t('common.delete') }}
            </BaseButton>
          </div>
        </div>
      </div>
    </BaseCard>

    <!-- Create/Edit Modal -->
    <BaseModal
      v-if="showModal"
      :title="editingCommunication ? t('communications.editMessage') : t('communications.newMessage')"
      size="large"
      @close="closeModal"
    >
      <CommunicationForm
        :communication="editingCommunication"
        @submit="handleSubmit"
        @cancel="closeModal"
      />
    </BaseModal>

    <!-- View Modal -->
    <BaseModal
      v-if="showViewModal"
      :title="t('communications.messageDetails')"
      size="medium"
      @close="showViewModal = false"
    >
      <CommunicationView
        :communication="viewingCommunication"
        @close="showViewModal = false"
      />
    </BaseModal>

    <!-- Template Modal -->
    <BaseModal
      v-if="showTemplateModal"
      :title="editingTemplate ? t('communications.editTemplate') : t('communications.newTemplate')"
      size="large"
      @close="closeTemplateModal"
    >
      <CommunicationForm
        :communication="editingTemplate"
        :is-template="true"
        @submit="handleTemplateSubmit"
        @cancel="closeTemplateModal"
      />
    </BaseModal>

    <!-- Send Confirmation Modal -->
    <BaseModal
      v-if="showSendModal"
      :title="t('communications.sendMessage')"
      size="small"
      @close="showSendModal = false"
    >
      <p>{{ t('communications.sendMessageConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showSendModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="primary" @click="sendCommunication" :loading="sending">
          {{ t('communications.send') }}
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      :title="t('communications.deleteMessage')"
      size="small"
      @close="showDeleteModal = false"
    >
      <p>{{ t('communications.deleteMessageConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteCommunication" :loading="deleting">
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
import CommunicationForm from './CommunicationForm.vue'
import CommunicationView from './CommunicationView.vue'

const { t } = useI18n()
const notificationStore = useNotificationStore()

// State
const communications = ref([])
const templates = ref([])
const statistics = ref({
  totalSent: 0,
  scheduled: 0,
  delivered: 0,
  failed: 0
})
const loading = ref(false)
const currentTab = ref('all')
const filters = ref({
  type: '',
  status: ''
})

// Modal state
const showModal = ref(false)
const showViewModal = ref(false)
const showTemplateModal = ref(false)
const showSendModal = ref(false)
const showDeleteModal = ref(false)
const editingCommunication = ref(null)
const editingTemplate = ref(null)
const viewingCommunication = ref(null)
const communicationToSend = ref(null)
const communicationToDelete = ref(null)
const sending = ref(false)
const deleting = ref(false)

// Table columns
const communicationColumns = [
  { key: 'subject', label: t('communications.subject'), sortable: true },
  { key: 'type', label: t('communications.type'), sortable: true },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'recipients', label: t('communications.recipients'), sortable: false },
  { key: 'scheduled_time', label: t('communications.scheduledTime'), sortable: true },
  { key: 'sent_at', label: t('communications.sentAt'), sortable: true },
  { key: 'actions', label: t('common.actions'), sortable: false }
]

// Computed
const filteredCommunications = computed(() => {
  let result = [...communications.value]

  // Filter by tab
  if (currentTab.value === 'scheduled') {
    result = result.filter(c => c.status === 'scheduled')
  }

  // Filter by type
  if (filters.value.type) {
    result = result.filter(c => c.type === filters.value.type)
  }

  // Filter by status
  if (filters.value.status) {
    result = result.filter(c => c.status === filters.value.status)
  }

  return result
})

const filteredTemplates = computed(() => {
  let result = [...templates.value]

  // Filter by type
  if (filters.value.type) {
    result = result.filter(t => t.type === filters.value.type)
  }

  return result
})

// Methods
const fetchCommunications = async () => {
  loading.value = true
  try {
    const response = await API.communications.list()
    communications.value = response.data.data || []
  } catch (error) {
    notificationStore.error(t('communications.fetchError'))
    console.error('Error fetching communications:', error)
  } finally {
    loading.value = false
  }
}

const fetchTemplates = async () => {
  try {
    const response = await API.communications.getTemplates()
    templates.value = response.data.data || []
  } catch (error) {
    console.error('Error fetching templates:', error)
  }
}

const fetchStatistics = async () => {
  try {
    const response = await API.communications.getStatistics()
    statistics.value = response.data.data || statistics.value
  } catch (error) {
    console.error('Error fetching statistics:', error)
  }
}

const openCreateModal = () => {
  editingCommunication.value = null
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingCommunication.value = null
}

const handleSubmit = async (communicationData) => {
  try {
    if (editingCommunication.value) {
      await API.communications.update(editingCommunication.value.id, communicationData)
      notificationStore.success(t('communications.updateSuccess'))
    } else {
      await API.communications.create(communicationData)
      notificationStore.success(t('communications.createSuccess'))
    }
    closeModal()
    fetchCommunications()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(
      editingCommunication.value ? t('communications.updateError') : t('communications.createError')
    )
    console.error('Error saving communication:', error)
  }
}

const openViewModal = (communication) => {
  viewingCommunication.value = communication
  showViewModal.value = true
}

const confirmSend = (communication) => {
  communicationToSend.value = communication
  showSendModal.value = true
}

const sendCommunication = async () => {
  sending.value = true
  try {
    await API.communications.send(communicationToSend.value.id)
    notificationStore.success(t('communications.sendSuccess'))
    showSendModal.value = false
    communicationToSend.value = null
    fetchCommunications()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(t('communications.sendError'))
    console.error('Error sending communication:', error)
  } finally {
    sending.value = false
  }
}

const cancelSchedule = async (communication) => {
  try {
    await API.communications.cancelSchedule(communication.id)
    notificationStore.success(t('communications.cancelSuccess'))
    fetchCommunications()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(t('communications.cancelError'))
    console.error('Error canceling schedule:', error)
  }
}

const confirmDelete = (communication) => {
  communicationToDelete.value = communication
  showDeleteModal.value = true
}

const deleteCommunication = async () => {
  deleting.value = true
  try {
    await API.communications.delete(communicationToDelete.value.id)
    notificationStore.success(t('communications.deleteSuccess'))
    showDeleteModal.value = false
    communicationToDelete.value = null
    fetchCommunications()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(t('communications.deleteError'))
    console.error('Error deleting communication:', error)
  } finally {
    deleting.value = false
  }
}

// Template methods
const openTemplateModal = () => {
  editingTemplate.value = null
  showTemplateModal.value = true
}

const closeTemplateModal = () => {
  showTemplateModal.value = false
  editingTemplate.value = null
}

const editTemplate = (template) => {
  editingTemplate.value = template
  showTemplateModal.value = true
}

const useTemplate = (template) => {
  editingCommunication.value = {
    ...template,
    id: null, // Clear ID so it creates a new communication
    status: 'draft'
  }
  showModal.value = true
}

const handleTemplateSubmit = async (templateData) => {
  try {
    if (editingTemplate.value) {
      await API.communications.updateTemplate(editingTemplate.value.id, templateData)
      notificationStore.success(t('communications.templateUpdateSuccess'))
    } else {
      await API.communications.createTemplate(templateData)
      notificationStore.success(t('communications.templateCreateSuccess'))
    }
    closeTemplateModal()
    fetchTemplates()
  } catch (error) {
    notificationStore.error(
      editingTemplate.value ? t('communications.templateUpdateError') : t('communications.templateCreateError')
    )
    console.error('Error saving template:', error)
  }
}

const confirmDeleteTemplate = (template) => {
  communicationToDelete.value = template
  communicationToDelete.value.isTemplate = true
  showDeleteModal.value = true
}

const handleSort = (column, direction) => {
  communications.value.sort((a, b) => {
    const aVal = a[column.key]
    const bVal = b[column.key]
    if (direction === 'asc') {
      return aVal > bVal ? 1 : -1
    } else {
      return aVal < bVal ? 1 : -1
    }
  })
}

const formatDateTime = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1)
}

// Lifecycle
onMounted(() => {
  fetchCommunications()
  fetchTemplates()
  fetchStatistics()
})
</script>

<style scoped>
.communications-page {
  max-width: 1400px;
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

.view-tabs {
  display: flex;
  gap: 0.5rem;
}

.tab-btn {
  padding: 0.5rem 1rem;
  border: 1px solid var(--border-color);
  background: var(--background-primary);
  color: var(--text-primary);
  border-radius: 0.375rem;
  cursor: pointer;
  font-size: 0.9rem;
  transition: all 0.2s;
}

.tab-btn:hover {
  background: var(--background-secondary);
}

.tab-btn.active {
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

/* Subject Cell */
.subject-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.subject-text {
  font-weight: 500;
}

.bulk-badge {
  font-size: 0.75rem;
  padding: 0.125rem 0.5rem;
  background: var(--primary-color);
  color: white;
  border-radius: 0.25rem;
}

/* Badges */
.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.85rem;
  font-weight: 500;
  display: inline-block;
}

.type-badge.type-email {
  background: #e3f2fd;
  color: #1976d2;
}

.type-badge.type-sms {
  background: #f3e5f5;
  color: #7b1fa2;
}

.type-badge.type-newsletter {
  background: #e8f5e9;
  color: #2e7d32;
}

.type-badge.type-notification {
  background: #fff3e0;
  color: #e65100;
}

.status-draft {
  background: #f5f5f5;
  color: #616161;
}

.status-scheduled {
  background: #fff3cd;
  color: #856404;
}

.status-sent {
  background: #d4edda;
  color: #155724;
}

.status-failed {
  background: #f8d7da;
  color: #721c24;
}

.recipients-count {
  font-weight: 500;
}

/* Actions */
.actions {
  display: flex;
  gap: 0.5rem;
}

/* Templates */
.templates-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
}

.templates-header h3 {
  margin: 0;
  font-size: 1.25rem;
}

.templates-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.template-card {
  padding: 1.5rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background: var(--background-secondary);
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.template-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 0.5rem;
}

.template-name {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.template-description {
  margin: 0;
  font-size: 0.9rem;
  color: var(--text-secondary);
  flex: 1;
}

.template-actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}

/* Responsive */
@media (max-width: 768px) {
  .controls-bar {
    flex-direction: column;
    align-items: stretch;
  }

  .view-tabs {
    width: 100%;
  }

  .tab-btn {
    flex: 1;
    font-size: 0.8rem;
  }

  .filters {
    width: 100%;
    flex-direction: column;
  }

  .actions {
    flex-direction: column;
  }

  .actions button {
    width: 100%;
  }

  .templates-grid {
    grid-template-columns: 1fr;
  }
}
</style>
