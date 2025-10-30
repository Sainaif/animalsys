<template>
  <form @submit.prevent="handleSubmit" class="document-form">
    <!-- File Upload -->
    <div class="form-section">
      <h3 class="section-title">{{ t('documents.fileUpload') }}</h3>

      <div class="upload-area" @click="triggerFileInput" @dragover.prevent @drop.prevent="handleDrop">
        <input
          ref="fileInput"
          type="file"
          @change="handleFileSelect"
          style="display: none"
        />

        <div v-if="!selectedFile" class="upload-placeholder">
          <div class="upload-icon">📁</div>
          <p class="upload-text">{{ t('documents.dragDropOrClick') }}</p>
          <p class="upload-hint">{{ t('documents.maxFileSize') }}</p>
        </div>

        <div v-else class="selected-file">
          <div class="file-icon">{{ getFileIcon(selectedFile.type) }}</div>
          <div class="file-info">
            <div class="file-name">{{ selectedFile.name }}</div>
            <div class="file-size">{{ formatFileSize(selectedFile.size) }}</div>
          </div>
          <button type="button" @click.stop="removeFile" class="remove-btn">
            ✕
          </button>
        </div>
      </div>
    </div>

    <!-- Document Information -->
    <div class="form-section">
      <h3 class="section-title">{{ t('documents.documentInfo') }}</h3>

      <FormGroup :label="t('documents.name')" required>
        <input
          v-model="formData.name"
          type="text"
          :placeholder="t('documents.namePlaceholder')"
          required
          class="form-input"
        />
      </FormGroup>

      <div class="form-row">
        <FormGroup :label="t('documents.category')" required>
          <select v-model="formData.category" required class="form-select">
            <option value="">{{ t('documents.selectCategory') }}</option>
            <option value="medical">{{ t('documents.categoryMedical') }}</option>
            <option value="legal">{{ t('documents.categoryLegal') }}</option>
            <option value="financial">{{ t('documents.categoryFinancial') }}</option>
            <option value="administrative">{{ t('documents.categoryAdministrative') }}</option>
            <option value="other">{{ t('documents.categoryOther') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('documents.type')" required>
          <select v-model="formData.type" required class="form-select">
            <option value="">{{ t('documents.selectType') }}</option>
            <option value="contract">{{ t('documents.typeContract') }}</option>
            <option value="invoice">{{ t('documents.typeInvoice') }}</option>
            <option value="report">{{ t('documents.typeReport') }}</option>
            <option value="certificate">{{ t('documents.typeCertificate') }}</option>
            <option value="photo">{{ t('documents.typePhoto') }}</option>
            <option value="other">{{ t('documents.typeOther') }}</option>
          </select>
        </FormGroup>
      </div>

      <FormGroup :label="t('documents.description')">
        <textarea
          v-model="formData.description"
          :placeholder="t('documents.descriptionPlaceholder')"
          rows="4"
          class="form-textarea"
        />
      </FormGroup>
    </div>

    <!-- Entity Association (Optional) -->
    <div class="form-section">
      <h3 class="section-title">{{ t('documents.entityAssociation') }}</h3>

      <div class="form-row">
        <FormGroup :label="t('documents.entityType')">
          <select v-model="formData.entity_type" class="form-select">
            <option value="">{{ t('documents.noAssociation') }}</option>
            <option value="animal">{{ t('documents.entityAnimal') }}</option>
            <option value="adoption">{{ t('documents.entityAdoption') }}</option>
            <option value="volunteer">{{ t('documents.entityVolunteer') }}</option>
            <option value="donor">{{ t('documents.entityDonor') }}</option>
            <option value="partner">{{ t('documents.entityPartner') }}</option>
          </select>
        </FormGroup>

        <FormGroup v-if="formData.entity_type" :label="t('documents.entityId')">
          <input
            v-model="formData.entity_id"
            type="text"
            :placeholder="t('documents.entityIdPlaceholder')"
            class="form-input"
          />
        </FormGroup>
      </div>
    </div>

    <!-- Additional Information -->
    <div class="form-section">
      <h3 class="section-title">{{ t('documents.additionalInfo') }}</h3>

      <div class="form-row">
        <FormGroup :label="t('documents.issueDate')">
          <input
            v-model="formData.issue_date"
            type="date"
            class="form-input"
          />
        </FormGroup>

        <FormGroup :label="t('documents.expiryDate')">
          <input
            v-model="formData.expiry_date"
            type="date"
            class="form-input"
          />
        </FormGroup>
      </div>

      <FormGroup :label="t('documents.tags')">
        <input
          v-model="formData.tags"
          type="text"
          :placeholder="t('documents.tagsPlaceholder')"
          class="form-input"
        />
      </FormGroup>

      <FormGroup :label="t('documents.notes')">
        <textarea
          v-model="formData.notes"
          :placeholder="t('documents.notesPlaceholder')"
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
      <BaseButton type="submit" variant="primary" :loading="uploading" :disabled="!selectedFile">
        {{ t('documents.upload') }}
      </BaseButton>
    </div>
  </form>
</template>

<script setup>
import { ref } from 'vue'
import { useI18n } from 'vue-i18n'
import FormGroup from '../../components/base/FormGroup.vue'
import BaseButton from '../../components/base/BaseButton.vue'

const { t } = useI18n()

const emit = defineEmits(['submit', 'cancel'])

// Form data
const formData = ref({
  name: '',
  category: '',
  type: '',
  description: '',
  entity_type: '',
  entity_id: '',
  issue_date: '',
  expiry_date: '',
  tags: '',
  notes: ''
})

const fileInput = ref(null)
const selectedFile = ref(null)
const uploading = ref(false)

// Methods
const triggerFileInput = () => {
  fileInput.value?.click()
}

const handleFileSelect = (event) => {
  const file = event.target.files[0]
  if (file) {
    validateAndSetFile(file)
  }
}

const handleDrop = (event) => {
  const file = event.dataTransfer.files[0]
  if (file) {
    validateAndSetFile(file)
  }
}

const validateAndSetFile = (file) => {
  // Check file size (max 10MB)
  const maxSize = 10 * 1024 * 1024 // 10MB
  if (file.size > maxSize) {
    alert(t('documents.fileTooLarge'))
    return
  }

  selectedFile.value = file

  // Auto-fill name if empty
  if (!formData.value.name) {
    formData.value.name = file.name
  }
}

const removeFile = () => {
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

const formatFileSize = (bytes) => {
  if (!bytes) return '0 B'
  const k = 1024
  const sizes = ['B', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return Math.round(bytes / Math.pow(k, i) * 100) / 100 + ' ' + sizes[i]
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

const handleSubmit = () => {
  if (!selectedFile.value) {
    alert(t('documents.noFileSelected'))
    return
  }

  // Create FormData
  const formDataToSubmit = new FormData()
  formDataToSubmit.append('file', selectedFile.value)

  // Append other fields
  Object.keys(formData.value).forEach(key => {
    if (formData.value[key]) {
      formDataToSubmit.append(key, formData.value[key])
    }
  })

  emit('submit', formDataToSubmit)
}
</script>

<style scoped>
.document-form {
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

/* Upload Area */
.upload-area {
  border: 2px dashed var(--border-color);
  border-radius: 0.5rem;
  padding: 2rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  background: var(--background-secondary);
}

.upload-area:hover {
  border-color: var(--primary-color);
  background: var(--input-background);
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1rem;
}

.upload-icon {
  font-size: 4rem;
}

.upload-text {
  font-size: 1.1rem;
  font-weight: 500;
  color: var(--text-primary);
  margin: 0;
}

.upload-hint {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin: 0;
}

.selected-file {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: var(--background-primary);
  border-radius: 0.375rem;
}

.file-icon {
  font-size: 3rem;
}

.file-info {
  flex: 1;
  text-align: left;
}

.file-name {
  font-size: 1rem;
  font-weight: 500;
  color: var(--text-primary);
  word-break: break-all;
}

.file-size {
  font-size: 0.9rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.remove-btn {
  padding: 0.5rem;
  background: var(--danger-color);
  color: white;
  border: none;
  border-radius: 50%;
  width: 2rem;
  height: 2rem;
  cursor: pointer;
  font-size: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: opacity 0.2s;
}

.remove-btn:hover {
  opacity: 0.8;
}

/* Form Elements */
.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
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
