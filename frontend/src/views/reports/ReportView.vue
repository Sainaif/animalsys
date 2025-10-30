<template>
  <div class="report-view">
    <div class="view-header">
      <h2 class="view-title">{{ report.name }}</h2>
      <div class="view-actions">
        <button
          v-if="report.status === 'completed'"
          @click="$emit('export', report, 'pdf')"
          class="btn btn-primary"
        >
          📄 {{ t('reports.exportPdf') }}
        </button>
        <button
          v-if="report.status === 'completed'"
          @click="$emit('export', report, 'excel')"
          class="btn btn-success"
        >
          📊 {{ t('reports.exportExcel') }}
        </button>
        <button
          v-if="report.status === 'completed'"
          @click="$emit('export', report, 'csv')"
          class="btn btn-secondary"
        >
          📋 {{ t('reports.exportCsv') }}
        </button>
        <button @click="$emit('delete', report)" class="btn btn-danger">
          {{ t('common.delete') }}
        </button>
        <button @click="$emit('close')" class="btn btn-outline">
          {{ t('common.close') }}
        </button>
      </div>
    </div>

    <!-- Report Information -->
    <div class="info-section">
      <h3 class="section-title">{{ t('reports.reportInformation') }}</h3>
      <div class="info-grid">
        <div class="info-item">
          <label>{{ t('reports.type') }}</label>
          <span class="badge" :class="`badge-${report.type}`">
            {{ t(`reports.types.${report.type}`) }}
          </span>
        </div>

        <div class="info-item">
          <label>{{ t('reports.status') }}</label>
          <span class="badge" :class="`badge-${report.status}`">
            {{ t(`reports.statuses.${report.status}`) }}
          </span>
        </div>

        <div v-if="report.description" class="info-item full-width">
          <label>{{ t('reports.description') }}</label>
          <span>{{ report.description }}</span>
        </div>

        <div class="info-item">
          <label>{{ t('reports.dateRange') }}</label>
          <span>{{ formatDate(report.parameters?.start_date) }} - {{ formatDate(report.parameters?.end_date) }}</span>
        </div>

        <div class="info-item">
          <label>{{ t('reports.format') }}</label>
          <span>{{ report.format?.toUpperCase() }}</span>
        </div>

        <div class="info-item">
          <label>{{ t('reports.generatedAt') }}</label>
          <span>{{ formatDateTime(report.created_at) }}</span>
        </div>

        <div v-if="report.generated_by" class="info-item">
          <label>{{ t('reports.generatedBy') }}</label>
          <span>{{ report.generated_by }}</span>
        </div>
      </div>
    </div>

    <!-- Report Summary -->
    <div v-if="report.summary" class="info-section">
      <h3 class="section-title">{{ t('reports.summary') }}</h3>
      <div class="summary-cards">
        <div
          v-for="(value, key) in report.summary"
          :key="key"
          class="summary-card"
        >
          <div class="summary-label">{{ formatSummaryLabel(key) }}</div>
          <div class="summary-value">{{ formatSummaryValue(value) }}</div>
        </div>
      </div>
    </div>

    <!-- Report Data Preview -->
    <div v-if="report.status === 'completed' && report.data" class="info-section">
      <h3 class="section-title">{{ t('reports.dataPreview') }}</h3>

      <!-- Financial Report Data -->
      <div v-if="report.type === 'financial' && report.data.transactions" class="data-table">
        <table>
          <thead>
            <tr>
              <th>{{ t('common.date') }}</th>
              <th>{{ t('finances.category') }}</th>
              <th>{{ t('finances.type') }}</th>
              <th>{{ t('finances.amount') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in report.data.transactions.slice(0, 10)" :key="index">
              <td>{{ formatDate(item.date) }}</td>
              <td>{{ item.category }}</td>
              <td>
                <span class="badge" :class="`badge-${item.type}`">
                  {{ item.type }}
                </span>
              </td>
              <td class="amount-cell" :class="item.type === 'income' ? 'positive' : 'negative'">
                {{ formatCurrency(item.amount) }}
              </td>
            </tr>
          </tbody>
        </table>
        <p v-if="report.data.transactions.length > 10" class="preview-note">
          {{ t('reports.showingFirstN', { n: 10, total: report.data.transactions.length }) }}
        </p>
      </div>

      <!-- Adoption Report Data -->
      <div v-if="report.type === 'adoption' && report.data.adoptions" class="data-table">
        <table>
          <thead>
            <tr>
              <th>{{ t('common.date') }}</th>
              <th>{{ t('adoptions.applicantName') }}</th>
              <th>{{ t('animals.name') }}</th>
              <th>{{ t('common.status') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in report.data.adoptions.slice(0, 10)" :key="index">
              <td>{{ formatDate(item.date) }}</td>
              <td>{{ item.applicant_name }}</td>
              <td>{{ item.animal_name }}</td>
              <td>
                <span class="badge" :class="`badge-${item.status}`">
                  {{ item.status }}
                </span>
              </td>
            </tr>
          </tbody>
        </table>
        <p v-if="report.data.adoptions.length > 10" class="preview-note">
          {{ t('reports.showingFirstN', { n: 10, total: report.data.adoptions.length }) }}
        </p>
      </div>

      <!-- Volunteer Report Data -->
      <div v-if="report.type === 'volunteer' && report.data.volunteers" class="data-table">
        <table>
          <thead>
            <tr>
              <th>{{ t('common.name') }}</th>
              <th>{{ t('volunteers.totalHours') }}</th>
              <th>{{ t('common.status') }}</th>
              <th>{{ t('volunteers.skills') }}</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in report.data.volunteers.slice(0, 10)" :key="index">
              <td>{{ item.name }}</td>
              <td>{{ item.total_hours }} {{ t('volunteers.hoursShort') }}</td>
              <td>
                <span class="badge" :class="`badge-${item.status}`">
                  {{ item.status }}
                </span>
              </td>
              <td>{{ item.skills }}</td>
            </tr>
          </tbody>
        </table>
        <p v-if="report.data.volunteers.length > 10" class="preview-note">
          {{ t('reports.showingFirstN', { n: 10, total: report.data.volunteers.length }) }}
        </p>
      </div>

      <!-- Generic Data Table for other report types -->
      <div v-if="!['financial', 'adoption', 'volunteer'].includes(report.type) && Array.isArray(report.data)" class="data-table">
        <table>
          <thead>
            <tr>
              <th v-for="(key, index) in Object.keys(report.data[0] || {})" :key="index">
                {{ formatTableHeader(key) }}
              </th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="(item, index) in report.data.slice(0, 10)" :key="index">
              <td v-for="(value, key) in item" :key="key">
                {{ formatCellValue(value) }}
              </td>
            </tr>
          </tbody>
        </table>
        <p v-if="report.data.length > 10" class="preview-note">
          {{ t('reports.showingFirstN', { n: 10, total: report.data.length }) }}
        </p>
      </div>
    </div>

    <!-- Report Status Messages -->
    <div v-if="report.status === 'generating'" class="status-message generating">
      <div class="status-icon">⏳</div>
      <div class="status-text">
        <h4>{{ t('reports.generatingReport') }}</h4>
        <p>{{ t('reports.generatingMessage') }}</p>
      </div>
    </div>

    <div v-if="report.status === 'failed'" class="status-message failed">
      <div class="status-icon">❌</div>
      <div class="status-text">
        <h4>{{ t('reports.generationFailed') }}</h4>
        <p>{{ report.error_message || t('reports.generationFailedMessage') }}</p>
      </div>
    </div>

    <!-- Metadata -->
    <div class="info-section">
      <h3 class="section-title">{{ t('common.metadata') }}</h3>
      <div class="info-grid">
        <div class="info-item">
          <label>{{ t('common.created') }}</label>
          <span>{{ formatDateTime(report.created_at) }}</span>
        </div>

        <div v-if="report.updated_at" class="info-item">
          <label>{{ t('common.updated') }}</label>
          <span>{{ formatDateTime(report.updated_at) }}</span>
        </div>

        <div v-if="report.file_size" class="info-item">
          <label>{{ t('reports.fileSize') }}</label>
          <span>{{ formatFileSize(report.file_size) }}</span>
        </div>

        <div v-if="report.record_count" class="info-item">
          <label>{{ t('reports.recordCount') }}</label>
          <span>{{ report.record_count }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const props = defineProps({
  report: {
    type: Object,
    required: true
  }
})

defineEmits(['close', 'export', 'delete'])

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit'
  })
}

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

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount)
}

const formatFileSize = (bytes) => {
  if (!bytes) return '-'
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(1024))
  return Math.round(bytes / Math.pow(1024, i) * 100) / 100 + ' ' + sizes[i]
}

const formatSummaryLabel = (key) => {
  // Convert snake_case to Title Case
  return key.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
}

const formatSummaryValue = (value) => {
  if (typeof value === 'number') {
    return value.toLocaleString('pl-PL')
  }
  return value
}

const formatTableHeader = (key) => {
  return key.split('_').map(word => word.charAt(0).toUpperCase() + word.slice(1)).join(' ')
}

const formatCellValue = (value) => {
  if (value === null || value === undefined) return '-'
  if (typeof value === 'boolean') return value ? t('common.yes') : t('common.no')
  if (typeof value === 'number') return value.toLocaleString('pl-PL')
  return value
}
</script>

<style scoped>
.report-view {
  max-width: 1200px;
}

.view-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
  padding-bottom: 1rem;
  border-bottom: 2px solid var(--border-color);
  flex-wrap: wrap;
  gap: 1rem;
}

.view-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--text-primary);
}

.view-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
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

.badge {
  display: inline-block;
  padding: 0.375rem 0.75rem;
  border-radius: 12px;
  font-size: 0.8125rem;
  font-weight: 500;
}

.badge-financial {
  background: #dcfce7;
  color: #15803d;
}

.badge-adoption {
  background: #dbeafe;
  color: #1e40af;
}

.badge-volunteer {
  background: #e0e7ff;
  color: #3730a3;
}

.badge-inventory {
  background: #fef3c7;
  color: #92400e;
}

.badge-completed {
  background: #dcfce7;
  color: #15803d;
}

.badge-generating {
  background: #fef3c7;
  color: #92400e;
}

.badge-failed {
  background: #fee2e2;
  color: #991b1b;
}

.summary-cards {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.summary-card {
  background: var(--background);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1rem;
  text-align: center;
}

.summary-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
}

.summary-value {
  font-size: 1.75rem;
  font-weight: 700;
  color: var(--primary);
}

.data-table {
  overflow-x: auto;
}

.data-table table {
  width: 100%;
  border-collapse: collapse;
}

.data-table th,
.data-table td {
  padding: 0.75rem;
  text-align: left;
  border-bottom: 1px solid var(--border-color);
}

.data-table th {
  background: var(--background);
  font-weight: 600;
  font-size: 0.875rem;
  color: var(--text-secondary);
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.data-table td {
  font-size: 0.9375rem;
  color: var(--text-primary);
}

.data-table tbody tr:hover {
  background: var(--hover-bg);
}

.amount-cell {
  font-weight: 600;
}

.amount-cell.positive {
  color: #15803d;
}

.amount-cell.negative {
  color: #dc2626;
}

.preview-note {
  margin-top: 1rem;
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-style: italic;
}

.status-message {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1.5rem;
  border-radius: 8px;
  margin-bottom: 1.5rem;
}

.status-message.generating {
  background: #fef3c7;
  border: 1px solid #f59e0b;
}

.status-message.failed {
  background: #fee2e2;
  border: 1px solid #ef4444;
}

.status-icon {
  font-size: 3rem;
}

.status-text h4 {
  margin: 0 0 0.5rem 0;
  font-size: 1.125rem;
  font-weight: 600;
}

.status-text p {
  margin: 0;
  font-size: 0.9375rem;
  color: var(--text-secondary);
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 6px;
  font-size: 0.875rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
  white-space: nowrap;
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

.btn-danger {
  background: #ef4444;
  color: white;
}

.btn-danger:hover {
  background: #dc2626;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
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
