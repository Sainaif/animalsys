<template>
  <div class="document-view">
    <div v-if="document" class="document-content">
      <!-- File Preview -->
      <div class="preview-section">
        <div class="file-preview">
          <div class="file-icon-large">{{ getFileIcon(document.file_type) }}</div>
          <div class="file-details">
            <h2 class="file-name">{{ document.name }}</h2>
            <div class="file-meta">
              <span>{{ document.file_type || t('documents.unknown') }}</span>
              <span>•</span>
              <span>{{ formatFileSize(document.file_size) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Document Information -->
      <div class="info-section">
        <h3 class="section-title">{{ t('documents.documentInfo') }}</h3>

        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('documents.category') }}</span>
            <span :class="['badge', 'category-badge', `category-${document.category}`]">
              {{ t(`documents.category${capitalize(document.category)}`) }}
            </span>
          </div>

          <div class="info-item">
            <span class="info-label">{{ t('documents.type') }}</span>
            <span class="badge type-badge">
              {{ t(`documents.type${capitalize(document.type)}`) }}
            </span>
          </div>

          <div v-if="document.description" class="info-item full-width">
            <span class="info-label">{{ t('documents.description') }}</span>
            <span class="info-value">{{ document.description }}</span>
          </div>

          <div v-if="document.issue_date" class="info-item">
            <span class="info-label">{{ t('documents.issueDate') }}</span>
            <span class="info-value">{{ formatDate(document.issue_date) }}</span>
          </div>

          <div v-if="document.expiry_date" class="info-item">
            <span class="info-label">{{ t('documents.expiryDate') }}</span>
            <span class="info-value" :class="{ 'text-warning': isExpiringSoon(document.expiry_date) }">
              {{ formatDate(document.expiry_date) }}
              <span v-if="isExpiringSoon(document.expiry_date)" class="expiry-warning">
                ⚠️ {{ t('documents.expiringSoon') }}
              </span>
            </span>
          </div>

          <div v-if="document.tags" class="info-item full-width">
            <span class="info-label">{{ t('documents.tags') }}</span>
            <div class="tags">
              <span v-for="tag in getTags(document.tags)" :key="tag" class="tag">
                {{ tag }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Entity Association -->
      <div v-if="document.entity_type" class="info-section">
        <h3 class="section-title">{{ t('documents.entityAssociation') }}</h3>

        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('documents.entityType') }}</span>
            <span class="info-value">
              {{ t(`documents.entity${capitalize(document.entity_type)}`) }}
            </span>
          </div>

          <div class="info-item">
            <span class="info-label">{{ t('documents.entityId') }}</span>
            <span class="info-value">{{ document.entity_id }}</span>
          </div>
        </div>
      </div>

      <!-- Upload Information -->
      <div class="info-section">
        <h3 class="section-title">{{ t('documents.uploadInfo') }}</h3>

        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('documents.uploadedBy') }}</span>
            <span class="info-value" v-if="document.uploaded_by">
              👤 {{ document.uploaded_by.first_name }} {{ document.uploaded_by.last_name }}
            </span>
            <span class="info-value text-muted" v-else>{{ t('documents.unknown') }}</span>
          </div>

          <div class="info-item">
            <span class="info-label">{{ t('documents.uploadedAt') }}</span>
            <span class="info-value">{{ formatDateTime(document.uploaded_at) }}</span>
          </div>
        </div>
      </div>

      <!-- Notes -->
      <div v-if="document.notes" class="info-section">
        <h3 class="section-title">{{ t('documents.notes') }}</h3>
        <div class="notes-content">
          {{ document.notes }}
        </div>
      </div>

      <!-- Actions -->
      <div class="actions-section">
        <BaseButton variant="secondary" @click="$emit('close')">
          {{ t('common.close') }}
        </BaseButton>
        <BaseButton variant="primary" @click="$emit('download', document)">
          {{ t('documents.download') }}
        </BaseButton>
        <BaseButton variant="danger" @click="$emit('delete', document)">
          {{ t('common.delete') }}
        </BaseButton>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useI18n } from 'vue-i18n'
import BaseButton from '../../components/base/BaseButton.vue'

const { t } = useI18n()

const props = defineProps({
  document: {
    type: Object,
    required: true
  }
})

defineEmits(['close', 'download', 'delete'])

// Methods
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
    month: 'long',
    day: 'numeric'
  })
}

const formatDateTime = (dateString) => {
  if (!dateString) return ''
  return new Date(dateString).toLocaleDateString('pl-PL', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
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

const getTags = (tagsString) => {
  if (!tagsString) return []
  return tagsString.split(',').map(tag => tag.trim()).filter(tag => tag)
}

const isExpiringSoon = (expiryDate) => {
  if (!expiryDate) return false
  const expiry = new Date(expiryDate)
  const now = new Date()
  const daysUntilExpiry = Math.floor((expiry - now) / (1000 * 60 * 60 * 24))
  return daysUntilExpiry <= 30 && daysUntilExpiry >= 0
}

const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1)
}
</script>

<style scoped>
.document-view {
  padding: 1rem;
}

.document-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Preview Section */
.preview-section {
  background: var(--background-secondary);
  border-radius: 0.5rem;
  padding: 2rem;
}

.file-preview {
  display: flex;
  align-items: center;
  gap: 1.5rem;
}

.file-icon-large {
  font-size: 5rem;
}

.file-details {
  flex: 1;
}

.file-name {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0 0 0.5rem 0;
  color: var(--text-primary);
  word-break: break-word;
}

.file-meta {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.9rem;
  color: var(--text-secondary);
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

.text-muted {
  color: var(--text-secondary);
  font-style: italic;
}

.text-warning {
  color: var(--warning-color);
}

.expiry-warning {
  display: block;
  font-size: 0.85rem;
  margin-top: 0.25rem;
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
  border: 1px solid var(--border-color);
}

/* Tags */
.tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.tag {
  padding: 0.25rem 0.75rem;
  background: var(--primary-color);
  color: white;
  border-radius: 1rem;
  font-size: 0.85rem;
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

/* Responsive */
@media (max-width: 768px) {
  .file-preview {
    flex-direction: column;
    text-align: center;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .actions-section {
    flex-direction: column;
  }

  .actions-section button {
    width: 100%;
  }
}
</style>
