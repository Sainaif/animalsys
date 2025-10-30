<template>
  <div class="documents-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.documents') }}</h1>
      <BaseButton variant="primary" @click="openUploadModal">
        ➕ {{ t('documents.upload') }}
      </BaseButton>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">📄</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.totalDocuments }}</div>
            <div class="stat-label">{{ t('documents.totalDocuments') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">📁</div>
          <div class="stat-details">
            <div class="stat-value">{{ formatFileSize(statistics.totalSize) }}</div>
            <div class="stat-label">{{ t('documents.totalSize') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">📅</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.recentUploads }}</div>
            <div class="stat-label">{{ t('documents.recentUploads') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">⚠️</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.expiringDocuments }}</div>
            <div class="stat-label">{{ t('documents.expiringDocuments') }}</div>
          </div>
        </div>
      </BaseCard>
    </div>

    <!-- Filters and Search -->
    <BaseCard>
      <div class="filters">
        <div class="filter-group">
          <label>{{ t('documents.category') }}</label>
          <select v-model="filters.category" class="filter-select">
            <option value="">{{ t('common.all') }}</option>
            <option value="medical">{{ t('documents.categoryMedical') }}</option>
            <option value="legal">{{ t('documents.categoryLegal') }}</option>
            <option value="financial">{{ t('documents.categoryFinancial') }}</option>
            <option value="administrative">{{ t('documents.categoryAdministrative') }}</option>
            <option value="other">{{ t('documents.categoryOther') }}</option>
          </select>
        </div>

        <div class="filter-group">
          <label>{{ t('documents.type') }}</label>
          <select v-model="filters.type" class="filter-select">
            <option value="">{{ t('common.all') }}</option>
            <option value="contract">{{ t('documents.typeContract') }}</option>
            <option value="invoice">{{ t('documents.typeInvoice') }}</option>
            <option value="report">{{ t('documents.typeReport') }}</option>
            <option value="certificate">{{ t('documents.typeCertificate') }}</option>
            <option value="photo">{{ t('documents.typePhoto') }}</option>
            <option value="other">{{ t('documents.typeOther') }}</option>
          </select>
        </div>

        <div class="filter-group">
          <label>{{ t('common.search') }}</label>
          <input
            v-model="filters.search"
            type="text"
            :placeholder="t('documents.searchPlaceholder')"
            class="filter-input"
          />
        </div>
      </div>
    </BaseCard>

    <!-- Documents Table -->
    <BaseCard>
      <DataTable
        :columns="columns"
        :data="filteredDocuments"
        :loading="loading"
        @sort="handleSort"
      >
        <template #cell-name="{ row }">
          <div class="document-name-cell">
            <span class="document-icon">{{ getFileIcon(row.file_type) }}</span>
            <span class="document-name">{{ row.name }}</span>
          </div>
        </template>

        <template #cell-category="{ row }">
          <span :class="['badge', 'category-badge', `category-${row.category}`]">
            {{ t(`documents.category${capitalize(row.category)}`) }}
          </span>
        </template>

        <template #cell-type="{ row }">
          <span :class="['badge', 'type-badge']">
            {{ t(`documents.type${capitalize(row.type)}`) }}
          </span>
        </template>

        <template #cell-file_size="{ row }">
          {{ formatFileSize(row.file_size) }}
        </template>

        <template #cell-uploaded_at="{ row }">
          {{ formatDate(row.uploaded_at) }}
        </template>

        <template #cell-uploaded_by="{ row }">
          <span v-if="row.uploaded_by">
            👤 {{ row.uploaded_by.first_name }} {{ row.uploaded_by.last_name }}
          </span>
          <span v-else class="text-muted">{{ t('documents.unknown') }}</span>
        </template>

        <template #cell-actions="{ row }">
          <div class="actions">
            <BaseButton size="small" variant="secondary" @click="openViewModal(row)">
              {{ t('common.view') }}
            </BaseButton>
            <BaseButton size="small" variant="primary" @click="downloadDocument(row)">
              {{ t('documents.download') }}
            </BaseButton>
            <BaseButton size="small" variant="danger" @click="confirmDelete(row)">
              {{ t('common.delete') }}
            </BaseButton>
          </div>
        </template>
      </DataTable>
    </BaseCard>

    <!-- Upload Modal -->
    <BaseModal
      v-if="showUploadModal"
      :title="t('documents.uploadDocument')"
      size="large"
      @close="closeUploadModal"
    >
      <DocumentForm
        @submit="handleUpload"
        @cancel="closeUploadModal"
      />
    </BaseModal>

    <!-- View Modal -->
    <BaseModal
      v-if="showViewModal"
      :title="t('documents.documentDetails')"
      size="medium"
      @close="showViewModal = false"
    >
      <DocumentView
        :document="viewingDocument"
        @close="showViewModal = false"
        @download="downloadDocument"
        @delete="confirmDelete"
      />
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      :title="t('documents.deleteDocument')"
      size="small"
      @close="showDeleteModal = false"
    >
      <p>{{ t('documents.deleteDocumentConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteDocument" :loading="deleting">
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
import DocumentForm from './DocumentForm.vue'
import DocumentView from './DocumentView.vue'

const { t } = useI18n()
const notificationStore = useNotificationStore()

// State
const documents = ref([])
const statistics = ref({
  totalDocuments: 0,
  totalSize: 0,
  recentUploads: 0,
  expiringDocuments: 0
})
const loading = ref(false)
const filters = ref({
  category: '',
  type: '',
  search: ''
})

// Modal state
const showUploadModal = ref(false)
const showViewModal = ref(false)
const showDeleteModal = ref(false)
const viewingDocument = ref(null)
const documentToDelete = ref(null)
const deleting = ref(false)

// Table columns
const columns = [
  { key: 'name', label: t('documents.name'), sortable: true },
  { key: 'category', label: t('documents.category'), sortable: true },
  { key: 'type', label: t('documents.type'), sortable: true },
  { key: 'file_size', label: t('documents.fileSize'), sortable: true },
  { key: 'uploaded_at', label: t('documents.uploadedAt'), sortable: true },
  { key: 'uploaded_by', label: t('documents.uploadedBy'), sortable: false },
  { key: 'actions', label: t('common.actions'), sortable: false }
]

// Computed
const filteredDocuments = computed(() => {
  let result = [...documents.value]

  // Filter by category
  if (filters.value.category) {
    result = result.filter(d => d.category === filters.value.category)
  }

  // Filter by type
  if (filters.value.type) {
    result = result.filter(d => d.type === filters.value.type)
  }

  // Filter by search
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter(d =>
      d.name?.toLowerCase().includes(search) ||
      d.description?.toLowerCase().includes(search)
    )
  }

  return result
})

// Methods
const fetchDocuments = async () => {
  loading.value = true
  try {
    const response = await API.documents.list()
    documents.value = response.data.data || []
  } catch (error) {
    notificationStore.error(t('documents.fetchError'))
    console.error('Error fetching documents:', error)
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    const response = await API.documents.getStatistics()
    statistics.value = response.data.data || statistics.value
  } catch (error) {
    console.error('Error fetching statistics:', error)
  }
}

const openUploadModal = () => {
  showUploadModal.value = true
}

const closeUploadModal = () => {
  showUploadModal.value = false
}

const handleUpload = async (formData) => {
  try {
    await API.documents.upload(formData)
    notificationStore.success(t('documents.uploadSuccess'))
    closeUploadModal()
    fetchDocuments()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(t('documents.uploadError'))
    console.error('Error uploading document:', error)
  }
}

const openViewModal = (document) => {
  viewingDocument.value = document
  showViewModal.value = true
}

const downloadDocument = async (document) => {
  try {
    const response = await API.documents.download(document.id)

    // Create a blob from the response
    const blob = new Blob([response.data])
    const url = window.URL.createObjectURL(blob)

    // Create a temporary link and trigger download
    const link = document.createElement('a')
    link.href = url
    link.download = document.name || 'document'
    document.body.appendChild(link)
    link.click()

    // Cleanup
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)

    notificationStore.success(t('documents.downloadSuccess'))
  } catch (error) {
    notificationStore.error(t('documents.downloadError'))
    console.error('Error downloading document:', error)
  }
}

const confirmDelete = (document) => {
  documentToDelete.value = document
  showDeleteModal.value = true
  showViewModal.value = false
}

const deleteDocument = async () => {
  deleting.value = true
  try {
    await API.documents.delete(documentToDelete.value.id)
    notificationStore.success(t('documents.deleteSuccess'))
    showDeleteModal.value = false
    documentToDelete.value = null
    fetchDocuments()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(t('documents.deleteError'))
    console.error('Error deleting document:', error)
  } finally {
    deleting.value = false
  }
}

const handleSort = (column, direction) => {
  documents.value.sort((a, b) => {
    const aVal = a[column.key]
    const bVal = b[column.key]
    if (direction === 'asc') {
      return aVal > bVal ? 1 : -1
    } else {
      return aVal < bVal ? 1 : -1
    }
  })
}

const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: '2-digit',
    day: '2-digit',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const getFileIcon = (fileType) => {
  if (!fileType) return '📄'

  const type = fileType.toLowerCase()
  if (type.includes('pdf')) return '📕'
  if (type.includes('word') || type.includes('doc')) return '📘'
  if (type.includes('excel') || type.includes('sheet')) return '📗'
  if (type.includes('image') || type.includes('png') || type.includes('jpg') || type.includes('jpeg')) return '🖼️'
  if (type.includes('video')) return '🎥'
  if (type.includes('audio')) return '🎵'
  if (type.includes('zip') || type.includes('rar')) return '📦'

  return '📄'
}

const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1)
}

// Lifecycle
onMounted(() => {
  fetchDocuments()
  fetchStatistics()
})
</script>

<style scoped>
.documents-page {
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

/* Filters */
.filters {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
  padding: 1rem;
}

.filter-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.filter-group label {
  font-size: 0.9rem;
  font-weight: 500;
  color: var(--text-secondary);
}

.filter-select,
.filter-input {
  padding: 0.5rem;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  font-size: 0.9rem;
  background: var(--input-background);
  color: var(--text-primary);
}

.filter-select:focus,
.filter-input:focus {
  outline: none;
  border-color: var(--primary-color);
}

/* Document Name Cell */
.document-name-cell {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.document-icon {
  font-size: 1.5rem;
}

.document-name {
  font-weight: 500;
}

/* Badges */
.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.85rem;
  font-weight: 500;
  display: inline-block;
}

.category-badge.category-medical {
  background: #e3f2fd;
  color: #1976d2;
}

.category-badge.category-legal {
  background: #f3e5f5;
  color: #7b1fa2;
}

.category-badge.category-financial {
  background: #e8f5e9;
  color: #2e7d32;
}

.category-badge.category-administrative {
  background: #fff3e0;
  color: #e65100;
}

.category-badge.category-other {
  background: #f5f5f5;
  color: #616161;
}

.type-badge {
  background: var(--background-secondary);
  color: var(--text-primary);
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
@media (max-width: 768px) {
  .actions {
    flex-direction: column;
  }

  .actions button {
    width: 100%;
  }
}
</style>
