<template>
  <div class="communication-view">
    <div class="view-header">
      <h2 class="view-title">
        {{ communication.is_template ? t('communications.templateDetails') : t('communications.communicationDetails') }}
      </h2>
      <div class="view-actions">
        <button
          v-if="!communication.is_template && communication.status === 'draft'"
          @click="$emit('send')"
          class="btn btn-success"
        >
          {{ t('communications.send') }}
        </button>
        <button
          v-if="!communication.is_template && communication.status === 'scheduled'"
          @click="$emit('cancel-schedule')"
          class="btn btn-warning"
        >
          {{ t('communications.cancelSchedule') }}
        </button>
        <button
          v-if="communication.is_template"
          @click="$emit('use-template')"
          class="btn btn-primary"
        >
          {{ t('communications.useTemplate') }}
        </button>
        <button @click="$emit('edit')" class="btn btn-secondary">
          {{ t('common.edit') }}
        </button>
        <button @click="$emit('delete')" class="btn btn-danger">
          {{ t('common.delete') }}
        </button>
        <button @click="$emit('close')" class="btn btn-outline">
          {{ t('common.close') }}
        </button>
      </div>
    </div>

    <!-- Message Information -->
    <div class="info-section">
      <h3 class="section-title">{{ t('communications.messageInformation') }}</h3>
      <div class="info-grid">
        <div class="info-item">
          <label>{{ t('communications.type') }}</label>
          <span class="badge" :class="`badge-${communication.type}`">
            {{ t(`communications.types.${communication.type}`) }}
          </span>
        </div>

        <div class="info-item">
          <label>{{ t('communications.status') }}</label>
          <span class="badge" :class="`badge-${communication.status}`">
            {{ t(`communications.statuses.${communication.status}`) }}
          </span>
        </div>

        <div v-if="communication.subject" class="info-item full-width">
          <label>{{ t('communications.subject') }}</label>
          <span>{{ communication.subject }}</span>
        </div>

        <div class="info-item full-width">
          <label>{{ t('communications.message') }}</label>
          <div class="message-content">
            {{ communication.message }}
          </div>
        </div>

        <div v-if="communication.scheduled_time" class="info-item">
          <label>{{ t('communications.scheduledTime') }}</label>
          <span>{{ formatDateTime(communication.scheduled_time) }}</span>
        </div>

        <div v-if="communication.sent_at" class="info-item">
          <label>{{ t('communications.sentAt') }}</label>
          <span>{{ formatDateTime(communication.sent_at) }}</span>
        </div>
      </div>
    </div>

    <!-- Template Information -->
    <div v-if="communication.is_template" class="info-section">
      <h3 class="section-title">{{ t('communications.templateInformation') }}</h3>
      <div class="info-grid">
        <div class="info-item full-width">
          <label>{{ t('communications.templateName') }}</label>
          <span>{{ communication.template_name }}</span>
        </div>

        <div v-if="communication.template_description" class="info-item full-width">
          <label>{{ t('communications.templateDescription') }}</label>
          <span>{{ communication.template_description }}</span>
        </div>

        <div class="info-item">
          <label>{{ t('communications.timesUsed') }}</label>
          <span>{{ communication.times_used || 0 }}</span>
        </div>

        <div class="info-item">
          <label>{{ t('communications.lastUsed') }}</label>
          <span>{{ communication.last_used_at ? formatDateTime(communication.last_used_at) : t('communications.neverUsed') }}</span>
        </div>
      </div>
    </div>

    <!-- Recipients Information -->
    <div v-if="!communication.is_template" class="info-section">
      <h3 class="section-title">
        {{ t('communications.recipients') }}
        <span v-if="isBulk" class="bulk-indicator">📢 {{ t('communications.bulk') }}</span>
      </h3>

      <div class="info-grid">
        <div v-if="communication.recipient_type" class="info-item">
          <label>{{ t('communications.recipientType') }}</label>
          <span>{{ t(`communications.recipientTypes.${communication.recipient_type}`) }}</span>
        </div>

        <div v-if="communication.total_recipients" class="info-item">
          <label>{{ t('communications.totalRecipients') }}</label>
          <span>{{ communication.total_recipients }}</span>
        </div>

        <div v-if="communication.delivered_count !== undefined" class="info-item">
          <label>{{ t('communications.delivered') }}</label>
          <span>{{ communication.delivered_count }} / {{ communication.total_recipients }}</span>
        </div>

        <div v-if="communication.failed_count !== undefined && communication.failed_count > 0" class="info-item">
          <label>{{ t('communications.failed') }}</label>
          <span class="text-danger">{{ communication.failed_count }}</span>
        </div>
      </div>

      <!-- Recipients List -->
      <div v-if="recipients.length > 0" class="recipients-section">
        <h4 class="subsection-title">{{ t('communications.recipientsList') }}</h4>
        <div class="recipients-table">
          <table>
            <thead>
              <tr>
                <th>{{ t('common.name') }}</th>
                <th>{{ t('common.email') }}</th>
                <th>{{ t('communications.deliveryStatus') }}</th>
                <th>{{ t('communications.deliveredAt') }}</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="recipient in recipients" :key="recipient.id">
                <td>{{ recipient.name }}</td>
                <td>{{ recipient.email || recipient.phone }}</td>
                <td>
                  <span class="badge" :class="`badge-${recipient.delivery_status}`">
                    {{ t(`communications.deliveryStatuses.${recipient.delivery_status}`) }}
                  </span>
                </td>
                <td>{{ recipient.delivered_at ? formatDateTime(recipient.delivered_at) : '-' }}</td>
              </tr>
            </tbody>
          </table>
        </div>

        <div v-if="loading.recipients" class="loading-state">
          {{ t('common.loading') }}...
        </div>
      </div>

      <!-- Custom Recipients -->
      <div v-if="communication.custom_recipients" class="info-item full-width">
        <label>{{ t('communications.customRecipients') }}</label>
        <div class="custom-recipients">
          {{ communication.custom_recipients }}
        </div>
      </div>
    </div>

    <!-- Delivery Statistics -->
    <div v-if="!communication.is_template && communication.status === 'sent'" class="info-section">
      <h3 class="section-title">{{ t('communications.deliveryStatistics') }}</h3>
      <div class="statistics-grid">
        <div class="stat-card">
          <div class="stat-label">{{ t('communications.deliveryRate') }}</div>
          <div class="stat-value">{{ deliveryRate }}%</div>
          <div class="stat-progress">
            <div class="progress-bar" :style="{ width: `${deliveryRate}%` }"></div>
          </div>
        </div>

        <div v-if="communication.opened_count !== undefined" class="stat-card">
          <div class="stat-label">{{ t('communications.openRate') }}</div>
          <div class="stat-value">{{ openRate }}%</div>
          <div class="stat-progress">
            <div class="progress-bar progress-bar-success" :style="{ width: `${openRate}%` }"></div>
          </div>
        </div>

        <div v-if="communication.clicked_count !== undefined" class="stat-card">
          <div class="stat-label">{{ t('communications.clickRate') }}</div>
          <div class="stat-value">{{ clickRate }}%</div>
          <div class="stat-progress">
            <div class="progress-bar progress-bar-info" :style="{ width: `${clickRate}%` }"></div>
          </div>
        </div>
      </div>
    </div>

    <!-- Metadata -->
    <div class="info-section">
      <h3 class="section-title">{{ t('common.metadata') }}</h3>
      <div class="info-grid">
        <div class="info-item">
          <label>{{ t('common.createdBy') }}</label>
          <span>{{ communication.created_by || t('common.unknown') }}</span>
        </div>

        <div class="info-item">
          <label>{{ t('common.createdAt') }}</label>
          <span>{{ formatDateTime(communication.created_at) }}</span>
        </div>

        <div v-if="communication.updated_at" class="info-item">
          <label>{{ t('common.updatedAt') }}</label>
          <span>{{ formatDateTime(communication.updated_at) }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { API } from '@/api'

const { t } = useI18n()

const props = defineProps({
  communication: {
    type: Object,
    required: true
  }
})

defineEmits(['close', 'edit', 'delete', 'send', 'cancel-schedule', 'use-template'])

const recipients = ref([])
const loading = ref({
  recipients: false
})

const isBulk = computed(() => {
  return props.communication.total_recipients > 1 || props.communication.recipient_type
})

const deliveryRate = computed(() => {
  if (!props.communication.total_recipients) return 0
  const delivered = props.communication.delivered_count || 0
  return Math.round((delivered / props.communication.total_recipients) * 100)
})

const openRate = computed(() => {
  if (!props.communication.delivered_count) return 0
  const opened = props.communication.opened_count || 0
  return Math.round((opened / props.communication.delivered_count) * 100)
})

const clickRate = computed(() => {
  if (!props.communication.delivered_count) return 0
  const clicked = props.communication.clicked_count || 0
  return Math.round((clicked / props.communication.delivered_count) * 100)
})

const formatDateTime = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleString('pl-PL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadRecipients = async () => {
  if (props.communication.is_template) return

  loading.value.recipients = true
  try {
    // This endpoint would return the list of recipients for this communication
    // For now, we'll use mock data structure
    // const response = await API.communications.getRecipients(props.communication.id)
    // recipients.value = response.data.data || []
  } catch (error) {
    console.error('Error loading recipients:', error)
  } finally {
    loading.value.recipients = false
  }
}

onMounted(() => {
  loadRecipients()
})
</script>

<style scoped>
.communication-view {
  max-width: 1200px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid var(--border-color);
}

.view-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
}

.view-actions {
  display: flex;
  gap: 0.75rem;
}

.info-section {
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
  display: flex;
  align-items: center;
  gap: 1rem;
}

.subsection-title {
  font-size: 1rem;
  font-weight: 600;
  margin: 1.5rem 0 1rem;
  color: var(--text-primary);
}

.bulk-indicator {
  font-size: 0.875rem;
  font-weight: 500;
  padding: 0.25rem 0.75rem;
  background: var(--warning-light);
  color: var(--warning-dark);
  border-radius: 12px;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
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

.info-item label {
  font-size: 0.875rem;
  font-weight: 600;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.info-item span {
  font-size: 0.9375rem;
  color: var(--text-primary);
}

.message-content {
  padding: 1rem;
  background: var(--background);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  white-space: pre-wrap;
  line-height: 1.6;
  color: var(--text-primary);
}

.custom-recipients {
  padding: 1rem;
  background: var(--background);
  border: 1px solid var(--border-color);
  border-radius: 6px;
  white-space: pre-wrap;
  font-family: monospace;
  font-size: 0.875rem;
  color: var(--text-primary);
}

.badge {
  display: inline-block;
  padding: 0.375rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8125rem;
  font-weight: 500;
  text-align: center;
}

.badge-email {
  background: #dbeafe;
  color: #1e40af;
}

.badge-sms {
  background: #dcfce7;
  color: #15803d;
}

.badge-newsletter {
  background: #fef3c7;
  color: #92400e;
}

.badge-notification {
  background: #e0e7ff;
  color: #3730a3;
}

.badge-draft {
  background: #f3f4f6;
  color: #374151;
}

.badge-scheduled {
  background: #dbeafe;
  color: #1e40af;
}

.badge-sent {
  background: #dcfce7;
  color: #15803d;
}

.badge-failed {
  background: #fee2e2;
  color: #991b1b;
}

.badge-pending {
  background: #fef3c7;
  color: #92400e;
}

.badge-delivered {
  background: #dcfce7;
  color: #15803d;
}

.text-danger {
  color: #dc2626;
  font-weight: 600;
}

.recipients-section {
  margin-top: 1.5rem;
}

.recipients-table {
  overflow-x: auto;
}

.recipients-table table {
  width: 100%;
  border-collapse: collapse;
}

.recipients-table th,
.recipients-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.recipients-table th {
  background: var(--background);
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.recipients-table td {
  font-size: 0.9375rem;
  color: var(--text-primary);
}

.recipients-table tbody tr:hover {
  background: var(--hover-bg);
}

.statistics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
}

.stat-card {
  padding: 1rem;
  background: var(--background);
  border: 1px solid var(--border-color);
  border-radius: 6px;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--primary);
  margin-bottom: 0.75rem;
}

.stat-progress {
  height: 8px;
  background: var(--border-color);
  border-radius: 4px;
  overflow: hidden;
}

.progress-bar {
  height: 100%;
  background: var(--primary);
  transition: width 0.3s ease;
}

.progress-bar-success {
  background: #10b981;
}

.progress-bar-info {
  background: #3b82f6;
}

.loading-state {
  padding: 2rem;
  text-align: center;
  color: var(--text-secondary);
  font-size: 0.9375rem;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-primary {
  background: var(--primary);
  color: white;
}

.btn-primary:hover {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.btn-success {
  background: #10b981;
  color: white;
}

.btn-success:hover {
  background: #059669;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.btn-warning {
  background: #f59e0b;
  color: white;
}

.btn-warning:hover {
  background: #d97706;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
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

.btn-outline {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-outline:hover {
  background: var(--hover-bg);
}
</style>
