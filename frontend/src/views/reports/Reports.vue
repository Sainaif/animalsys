<template>
  <div class="reports-container">
    <div class="reports-header">
      <h1 class="page-title">{{ t('reports.title') }}</h1>
      <button @click="showReportForm = true" class="btn btn-primary">
        {{ t('reports.generate') }}
      </button>
    </div>

    <!-- Statistics Dashboard -->
    <div class="statistics-grid">
      <div class="stat-card">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <div class="stat-value">{{ statistics.totalReports || 0 }}</div>
          <div class="stat-label">{{ t('reports.totalReports') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📅</div>
        <div class="stat-content">
          <div class="stat-value">{{ statistics.thisMonth || 0 }}</div>
          <div class="stat-label">{{ t('reports.generatedThisMonth') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">⏰</div>
        <div class="stat-content">
          <div class="stat-value">{{ statistics.scheduled || 0 }}</div>
          <div class="stat-label">{{ t('reports.scheduledReports') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">⭐</div>
        <div class="stat-content">
          <div class="stat-value">{{ statistics.favorites || 0 }}</div>
          <div class="stat-label">{{ t('reports.favoriteReports') }}</div>
        </div>
      </div>
    </div>

    <!-- Quick Report Actions -->
    <div class="quick-reports-section">
      <h2 class="section-title">{{ t('reports.quickReports') }}</h2>
      <div class="quick-reports-grid">
        <div
          v-for="reportType in quickReportTypes"
          :key="reportType.type"
          class="quick-report-card"
          @click="generateQuickReport(reportType.type)"
        >
          <div class="quick-report-icon">{{ reportType.icon }}</div>
          <div class="quick-report-title">{{ t(`reports.types.${reportType.type}`) }}</div>
          <div class="quick-report-description">{{ t(`reports.descriptions.${reportType.type}`) }}</div>
        </div>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="filters-section">
      <div class="filters-row">
        <div class="search-box">
          <input
            v-model="filters.search"
            type="text"
            :placeholder="t('reports.searchPlaceholder')"
            class="search-input"
          />
        </div>

        <select v-model="filters.type" class="filter-select">
          <option value="">{{ t('reports.allTypes') }}</option>
          <option value="financial">{{ t('reports.types.financial') }}</option>
          <option value="adoption">{{ t('reports.types.adoption') }}</option>
          <option value="volunteer">{{ t('reports.types.volunteer') }}</option>
          <option value="inventory">{{ t('reports.types.inventory') }}</option>
          <option value="veterinary">{{ t('reports.types.veterinary') }}</option>
          <option value="campaign">{{ t('reports.types.campaign') }}</option>
          <option value="donor">{{ t('reports.types.donor') }}</option>
          <option value="animal">{{ t('reports.types.animal') }}</option>
          <option value="statutory">{{ t('reports.types.statutory') }}</option>
          <option value="custom">{{ t('reports.types.custom') }}</option>
        </select>

        <select v-model="filters.status" class="filter-select">
          <option value="">{{ t('reports.allStatuses') }}</option>
          <option value="completed">{{ t('reports.statuses.completed') }}</option>
          <option value="generating">{{ t('reports.statuses.generating') }}</option>
          <option value="scheduled">{{ t('reports.statuses.scheduled') }}</option>
          <option value="failed">{{ t('reports.statuses.failed') }}</option>
        </select>

        <div class="date-range-filter">
          <input
            v-model="filters.startDate"
            type="date"
            class="date-input"
            :placeholder="t('common.from')"
          />
          <span class="date-separator">-</span>
          <input
            v-model="filters.endDate"
            type="date"
            class="date-input"
            :placeholder="t('common.to')"
          />
        </div>

        <button @click="loadReports" class="btn btn-secondary">
          {{ t('common.filter') }}
        </button>
      </div>
    </div>

    <!-- Reports List -->
    <div class="reports-list">
      <h2 class="section-title">{{ t('reports.generatedReports') }}</h2>

      <DataTable
        v-if="!loading && reports.length > 0"
        :columns="tableColumns"
        :data="reports"
        :sortable="true"
        @sort="handleSort"
      >
        <template #cell-name="{ row }">
          <div class="report-name-cell">
            <span class="report-icon">{{ getReportIcon(row.type) }}</span>
            <span class="report-name">{{ row.name }}</span>
          </div>
        </template>

        <template #cell-type="{ row }">
          <span class="badge" :class="`badge-${row.type}`">
            {{ t(`reports.types.${row.type}`) }}
          </span>
        </template>

        <template #cell-status="{ row }">
          <span class="badge" :class="`badge-${row.status}`">
            {{ t(`reports.statuses.${row.status}`) }}
          </span>
        </template>

        <template #cell-createdAt="{ row }">
          {{ formatDate(row.created_at) }}
        </template>

        <template #cell-actions="{ row }">
          <div class="action-buttons">
            <button
              @click="viewReport(row)"
              class="btn-icon"
              :title="t('common.view')"
            >
              👁️
            </button>
            <button
              v-if="row.status === 'completed'"
              @click="exportReport(row, 'pdf')"
              class="btn-icon"
              :title="t('reports.exportPdf')"
            >
              📄
            </button>
            <button
              v-if="row.status === 'completed'"
              @click="exportReport(row, 'excel')"
              class="btn-icon"
              :title="t('reports.exportExcel')"
            >
              📊
            </button>
            <button
              @click="confirmDelete(row)"
              class="btn-icon btn-danger"
              :title="t('common.delete')"
            >
              🗑️
            </button>
          </div>
        </template>
      </DataTable>

      <EmptyState
        v-if="!loading && reports.length === 0"
        :title="t('reports.noReports')"
        :message="t('reports.noReportsMessage')"
      />

      <LoadingSpinner v-if="loading" />
    </div>

    <!-- Report Form Modal -->
    <BaseModal
      v-if="showReportForm"
      @close="closeReportForm"
      :title="t('reports.generateReport')"
      size="large"
    >
      <ReportForm
        :report="editingReport"
        @submit="handleReportGenerate"
        @cancel="closeReportForm"
      />
    </BaseModal>

    <!-- Report View Modal -->
    <BaseModal
      v-if="showReportView"
      @close="closeReportView"
      :title="t('reports.reportDetails')"
      size="large"
    >
      <ReportView
        :report="viewingReport"
        @close="closeReportView"
        @export="exportReport"
        @delete="confirmDelete"
      />
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      @close="showDeleteModal = false"
      :title="t('reports.deleteReport')"
      size="small"
    >
      <div class="modal-content">
        <p>{{ t('reports.deleteReportConfirm') }}</p>
        <div class="modal-actions">
          <button @click="showDeleteModal = false" class="btn btn-secondary">
            {{ t('common.cancel') }}
          </button>
          <button @click="deleteReport" class="btn btn-danger">
            {{ t('common.delete') }}
          </button>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '@/stores/notification'
import { API } from '@/api'
import DataTable from '@/components/base/DataTable.vue'
import BaseModal from '@/components/base/BaseModal.vue'
import LoadingSpinner from '@/components/base/LoadingSpinner.vue'
import EmptyState from '@/components/base/EmptyState.vue'
import ReportForm from './ReportForm.vue'
import ReportView from './ReportView.vue'

const { t } = useI18n()
const notificationStore = useNotificationStore()

const loading = ref(false)
const reports = ref([])
const statistics = ref({})
const showReportForm = ref(false)
const showReportView = ref(false)
const showDeleteModal = ref(false)
const editingReport = ref(null)
const viewingReport = ref(null)
const reportToDelete = ref(null)

const filters = reactive({
  search: '',
  type: '',
  status: '',
  startDate: '',
  endDate: ''
})

const quickReportTypes = [
  { type: 'financial', icon: '💰' },
  { type: 'adoption', icon: '🏠' },
  { type: 'volunteer', icon: '👥' },
  { type: 'inventory', icon: '📦' },
  { type: 'veterinary', icon: '⚕️' },
  { type: 'campaign', icon: '📢' },
  { type: 'donor', icon: '💝' },
  { type: 'animal', icon: '🐾' },
  { type: 'statutory', icon: '📋' }
]

const tableColumns = computed(() => [
  { key: 'name', label: t('reports.name'), sortable: true },
  { key: 'type', label: t('reports.type'), sortable: true },
  { key: 'status', label: t('reports.status'), sortable: true },
  { key: 'createdAt', label: t('common.created'), sortable: true },
  { key: 'actions', label: t('common.actions'), sortable: false }
])

const getReportIcon = (type) => {
  const iconMap = {
    financial: '💰',
    adoption: '🏠',
    volunteer: '👥',
    inventory: '📦',
    veterinary: '⚕️',
    campaign: '📢',
    donor: '💝',
    animal: '🐾',
    statutory: '📋',
    custom: '⚙️'
  }
  return iconMap[type] || '📊'
}

const formatDate = (dateString) => {
  if (!dateString) return '-'
  const date = new Date(dateString)
  return date.toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const loadReports = async () => {
  loading.value = true
  try {
    const params = {
      search: filters.search || undefined,
      type: filters.type || undefined,
      status: filters.status || undefined,
      start_date: filters.startDate || undefined,
      end_date: filters.endDate || undefined
    }

    const response = await API.reports.list(params)
    reports.value = response.data.data || []
  } catch (error) {
    console.error('Error loading reports:', error)
    notificationStore.error(t('reports.fetchError'))
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
  try {
    const response = await API.reports.getStatistics()
    statistics.value = response.data.data || {}
  } catch (error) {
    console.error('Error loading statistics:', error)
  }
}

const generateQuickReport = async (reportType) => {
  try {
    loading.value = true
    const reportData = {
      type: reportType,
      name: `${t(`reports.types.${reportType}`)} - ${new Date().toLocaleDateString('pl-PL')}`,
      parameters: getDefaultParameters(reportType)
    }

    const response = await API.reports.generate(reportData)

    notificationStore.success(t('reports.generateSuccess'))
    loadReports()
    loadStatistics()

    // View the generated report
    viewingReport.value = response.data.data
    showReportView.value = true
  } catch (error) {
    console.error('Error generating report:', error)
    notificationStore.error(t('reports.generateError'))
  } finally {
    loading.value = false
  }
}

const getDefaultParameters = (reportType) => {
  const now = new Date()
  const firstDayOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
  const lastDayOfMonth = new Date(now.getFullYear(), now.getMonth() + 1, 0)

  return {
    start_date: firstDayOfMonth.toISOString().split('T')[0],
    end_date: lastDayOfMonth.toISOString().split('T')[0]
  }
}

const handleReportGenerate = async (reportData) => {
  try {
    loading.value = true
    const response = await API.reports.generate(reportData)

    notificationStore.success(t('reports.generateSuccess'))
    closeReportForm()
    loadReports()
    loadStatistics()

    // View the generated report
    viewingReport.value = response.data.data
    showReportView.value = true
  } catch (error) {
    console.error('Error generating report:', error)
    notificationStore.error(t('reports.generateError'))
  } finally {
    loading.value = false
  }
}

const viewReport = (report) => {
  viewingReport.value = report
  showReportView.value = true
}

const exportReport = async (report, format) => {
  try {
    const response = await API.reports.export(report.id, format)

    // Create blob and download
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `${report.name}.${format === 'excel' ? 'xlsx' : format}`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    notificationStore.success(t('reports.exportSuccess'))
  } catch (error) {
    console.error('Error exporting report:', error)
    notificationStore.error(t('reports.exportError'))
  }
}

const confirmDelete = (report) => {
  reportToDelete.value = report
  showDeleteModal.value = true
  showReportView.value = false
}

const deleteReport = async () => {
  if (!reportToDelete.value) return

  try {
    await API.reports.delete(reportToDelete.value.id)

    notificationStore.success(t('reports.deleteSuccess'))
    showDeleteModal.value = false
    reportToDelete.value = null
    loadReports()
    loadStatistics()
  } catch (error) {
    console.error('Error deleting report:', error)
    notificationStore.error(t('reports.deleteError'))
  }
}

const handleSort = ({ column, direction }) => {
  // Implement sorting logic
  const sorted = [...reports.value].sort((a, b) => {
    const aVal = a[column]
    const bVal = b[column]

    if (direction === 'asc') {
      return aVal > bVal ? 1 : -1
    } else {
      return aVal < bVal ? 1 : -1
    }
  })

  reports.value = sorted
}

const closeReportForm = () => {
  showReportForm.value = false
  editingReport.value = null
}

const closeReportView = () => {
  showReportView.value = false
  viewingReport.value = null
}

onMounted(() => {
  loadReports()
  loadStatistics()
})
</script>

<style scoped>
.reports-container {
  padding: 2rem;
  max-width: 1400px;
  margin: 0 auto;
}

.reports-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: 700;
  color: var(--text-primary);
  margin: 0;
}

.statistics-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 12px;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: transform 0.2s, box-shadow 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: var(--primary);
  line-height: 1;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.quick-reports-section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 1rem;
}

.quick-reports-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 1rem;
}

.quick-report-card {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s;
}

.quick-report-card:hover {
  background: var(--hover-bg);
  transform: translateY(-2px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  border-color: var(--primary);
}

.quick-report-icon {
  font-size: 2.5rem;
  margin-bottom: 0.75rem;
}

.quick-report-title {
  font-size: 1rem;
  font-weight: 600;
  color: var(--text-primary);
  margin-bottom: 0.5rem;
}

.quick-report-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
  line-height: 1.4;
}

.filters-section {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 2rem;
}

.filters-row {
  display: flex;
  gap: 1rem;
  flex-wrap: wrap;
  align-items: center;
}

.search-box {
  flex: 1;
  min-width: 250px;
}

.search-input {
  width: 100%;
  padding: 0.625rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 0.9375rem;
  background: var(--input-bg);
  color: var(--text-primary);
}

.filter-select {
  padding: 0.625rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 0.9375rem;
  background: var(--input-bg);
  color: var(--text-primary);
  min-width: 150px;
}

.date-range-filter {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.date-input {
  padding: 0.625rem 1rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 0.9375rem;
  background: var(--input-bg);
  color: var(--text-primary);
}

.date-separator {
  color: var(--text-secondary);
}

.reports-list {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
}

.report-name-cell {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.report-icon {
  font-size: 1.5rem;
}

.report-name {
  font-weight: 500;
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

.badge-veterinary {
  background: #fce7f3;
  color: #9f1239;
}

.badge-campaign {
  background: #fee2e2;
  color: #991b1b;
}

.badge-donor {
  background: #f3e8ff;
  color: #6b21a8;
}

.badge-animal {
  background: #dbeafe;
  color: #1e40af;
}

.badge-statutory {
  background: #f3f4f6;
  color: #374151;
}

.badge-custom {
  background: #e0f2fe;
  color: #0c4a6e;
}

.badge-completed {
  background: #dcfce7;
  color: #15803d;
}

.badge-generating {
  background: #fef3c7;
  color: #92400e;
}

.badge-scheduled {
  background: #dbeafe;
  color: #1e40af;
}

.badge-failed {
  background: #fee2e2;
  color: #991b1b;
}

.action-buttons {
  display: flex;
  gap: 0.5rem;
}

.btn-icon {
  background: none;
  border: none;
  cursor: pointer;
  font-size: 1.25rem;
  padding: 0.25rem;
  transition: transform 0.2s;
}

.btn-icon:hover {
  transform: scale(1.2);
}

.btn-icon.btn-danger:hover {
  filter: brightness(1.2);
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

.btn-primary:hover {
  background: var(--primary-dark);
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

.modal-content {
  padding: 1rem;
}

.modal-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  margin-top: 1.5rem;
}
</style>
