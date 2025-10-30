<template>
  <div class="partner-view-page">
    <div v-if="loading" class="loading-container">
      <LoadingSpinner />
    </div>

    <div v-else-if="partner" class="partner-content">
      <!-- Page Header -->
      <div class="page-header">
        <div class="header-left">
          <router-link to="/partners" class="back-link">
            ← {{ t('common.back') }}
          </router-link>
          <h1 class="page-title">{{ partner.name }}</h1>
          <span :class="['badge', `status-${partner.status}`]">
            {{ t(`partners.status${capitalize(partner.status)}`) }}
          </span>
        </div>
        <BaseButton variant="primary" @click="$router.push('/partners')">
          {{ t('common.backToList') }}
        </BaseButton>
      </div>

      <!-- Partner Information Card -->
      <BaseCard>
        <template #header>
          <h2 class="card-title">{{ t('partners.partnerInfo') }}</h2>
        </template>

        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('partners.partnerType') }}</span>
            <span :class="['badge', 'type-badge', `type-${partner.type}`]">
              {{ t(`partners.type${capitalize(partner.type)}`) }}
            </span>
          </div>

          <div class="info-item" v-if="partner.description">
            <span class="info-label">{{ t('partners.description') }}</span>
            <span class="info-value">{{ partner.description }}</span>
          </div>

          <div class="info-item" v-if="partner.services_provided">
            <span class="info-label">{{ t('partners.servicesProvided') }}</span>
            <span class="info-value">{{ partner.services_provided }}</span>
          </div>

          <div class="info-item" v-if="partner.partnership_start_date">
            <span class="info-label">{{ t('partners.partnershipStartDate') }}</span>
            <span class="info-value">{{ formatDate(partner.partnership_start_date) }}</span>
          </div>

          <div class="info-item" v-if="partner.partnership_end_date">
            <span class="info-label">{{ t('partners.partnershipEndDate') }}</span>
            <span class="info-value">{{ formatDate(partner.partnership_end_date) }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Contact Information Card -->
      <BaseCard>
        <template #header>
          <h2 class="card-title">{{ t('partners.contactInfo') }}</h2>
        </template>

        <div class="info-grid">
          <div class="info-item" v-if="partner.contact_person">
            <span class="info-label">{{ t('partners.contactPerson') }}</span>
            <span class="info-value">👤 {{ partner.contact_person }}</span>
          </div>

          <div class="info-item" v-if="partner.email">
            <span class="info-label">{{ t('partners.email') }}</span>
            <span class="info-value">
              <a :href="`mailto:${partner.email}`" class="link">✉️ {{ partner.email }}</a>
            </span>
          </div>

          <div class="info-item" v-if="partner.phone">
            <span class="info-label">{{ t('partners.phone') }}</span>
            <span class="info-value">
              <a :href="`tel:${partner.phone}`" class="link">📞 {{ partner.phone }}</a>
            </span>
          </div>

          <div class="info-item" v-if="partner.website">
            <span class="info-label">{{ t('partners.website') }}</span>
            <span class="info-value">
              <a :href="partner.website" target="_blank" rel="noopener noreferrer" class="link">
                🌐 {{ partner.website }}
              </a>
            </span>
          </div>

          <div class="info-item full-width" v-if="partner.address">
            <span class="info-label">{{ t('partners.address') }}</span>
            <span class="info-value">📍 {{ partner.address }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Agreements Section -->
      <BaseCard>
        <template #header>
          <div class="section-header">
            <h2 class="card-title">{{ t('partners.agreements') }}</h2>
            <BaseButton size="small" variant="primary" @click="openAddAgreementModal">
              ➕ {{ t('partners.addAgreement') }}
            </BaseButton>
          </div>
        </template>

        <div v-if="agreements.length === 0" class="empty-state">
          <EmptyState
            icon="📄"
            :title="t('partners.noAgreements')"
            :description="t('partners.noAgreementsMessage')"
          />
        </div>

        <div v-else class="agreements-list">
          <div v-for="agreement in agreements" :key="agreement.id" class="agreement-item">
            <div class="agreement-header">
              <div class="agreement-title">
                <span class="agreement-name">{{ agreement.title }}</span>
                <span :class="['badge', 'agreement-status', `status-${agreement.status}`]">
                  {{ t(`partners.agreement${capitalize(agreement.status)}`) }}
                </span>
              </div>
              <div class="agreement-actions">
                <BaseButton size="small" variant="secondary" @click="editAgreement(agreement)">
                  {{ t('common.edit') }}
                </BaseButton>
                <BaseButton size="small" variant="danger" @click="confirmDeleteAgreement(agreement)">
                  {{ t('common.delete') }}
                </BaseButton>
              </div>
            </div>

            <div class="agreement-details">
              <div v-if="agreement.description" class="agreement-description">
                {{ agreement.description }}
              </div>

              <div class="agreement-meta">
                <span v-if="agreement.start_date">
                  📅 {{ t('partners.startDate') }}: {{ formatDate(agreement.start_date) }}
                </span>
                <span v-if="agreement.end_date">
                  📅 {{ t('partners.endDate') }}: {{ formatDate(agreement.end_date) }}
                </span>
                <span v-if="agreement.value">
                  💰 {{ t('partners.value') }}: {{ formatCurrency(agreement.value) }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </BaseCard>

      <!-- Notes Section -->
      <BaseCard v-if="partner.notes">
        <template #header>
          <h2 class="card-title">{{ t('partners.notes') }}</h2>
        </template>
        <div class="notes-content">
          {{ partner.notes }}
        </div>
      </BaseCard>

      <!-- Add/Edit Agreement Modal -->
      <BaseModal
        v-if="showAgreementModal"
        :title="editingAgreement ? t('partners.editAgreement') : t('partners.addAgreement')"
        size="medium"
        @close="closeAgreementModal"
      >
        <form @submit.prevent="handleAgreementSubmit" class="agreement-form">
          <FormGroup :label="t('partners.agreementTitle')" required>
            <input
              v-model="agreementFormData.title"
              type="text"
              :placeholder="t('partners.agreementTitlePlaceholder')"
              required
              class="form-input"
            />
          </FormGroup>

          <FormGroup :label="t('partners.description')">
            <textarea
              v-model="agreementFormData.description"
              :placeholder="t('partners.descriptionPlaceholder')"
              rows="3"
              class="form-textarea"
            />
          </FormGroup>

          <div class="form-row">
            <FormGroup :label="t('partners.startDate')">
              <input
                v-model="agreementFormData.start_date"
                type="date"
                class="form-input"
              />
            </FormGroup>

            <FormGroup :label="t('partners.endDate')">
              <input
                v-model="agreementFormData.end_date"
                type="date"
                class="form-input"
              />
            </FormGroup>
          </div>

          <FormGroup :label="t('partners.value')">
            <input
              v-model.number="agreementFormData.value"
              type="number"
              step="0.01"
              :placeholder="t('partners.valuePlaceholder')"
              class="form-input"
            />
          </FormGroup>

          <FormGroup :label="t('common.status')" required>
            <select v-model="agreementFormData.status" required class="form-select">
              <option value="active">{{ t('partners.agreementActive') }}</option>
              <option value="expired">{{ t('partners.agreementExpired') }}</option>
              <option value="terminated">{{ t('partners.agreementTerminated') }}</option>
            </select>
          </FormGroup>

          <div class="form-actions">
            <BaseButton type="button" variant="secondary" @click="closeAgreementModal">
              {{ t('common.cancel') }}
            </BaseButton>
            <BaseButton type="submit" variant="primary" :loading="savingAgreement">
              {{ t('common.save') }}
            </BaseButton>
          </div>
        </form>
      </BaseModal>

      <!-- Delete Agreement Confirmation Modal -->
      <BaseModal
        v-if="showDeleteAgreementModal"
        :title="t('partners.deleteAgreement')"
        size="small"
        @close="showDeleteAgreementModal = false"
      >
        <p>{{ t('partners.deleteAgreementConfirm') }}</p>
        <template #footer>
          <BaseButton variant="secondary" @click="showDeleteAgreementModal = false">
            {{ t('common.cancel') }}
          </BaseButton>
          <BaseButton variant="danger" @click="deleteAgreement" :loading="deletingAgreement">
            {{ t('common.delete') }}
          </BaseButton>
        </template>
      </BaseModal>
    </div>

    <div v-else class="error-container">
      <EmptyState
        icon="❌"
        :title="t('partners.notFound')"
        :description="t('partners.notFoundMessage')"
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRoute, useRouter } from 'vue-router'
import { API } from '../../api'
import { useNotificationStore } from '../../stores/notification'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'
import EmptyState from '../../components/base/EmptyState.vue'
import FormGroup from '../../components/base/FormGroup.vue'

const { t } = useI18n()
const route = useRoute()
const router = useRouter()
const notificationStore = useNotificationStore()

// State
const partner = ref(null)
const agreements = ref([])
const loading = ref(false)

// Agreement modal state
const showAgreementModal = ref(false)
const showDeleteAgreementModal = ref(false)
const editingAgreement = ref(null)
const agreementToDelete = ref(null)
const savingAgreement = ref(false)
const deletingAgreement = ref(false)

const agreementFormData = ref({
  title: '',
  description: '',
  start_date: '',
  end_date: '',
  value: null,
  status: 'active'
})

// Methods
const fetchPartner = async () => {
  loading.value = true
  try {
    const response = await API.partners.getById(route.params.id)
    partner.value = response.data.data
  } catch (error) {
    notificationStore.error(t('partners.fetchError'))
    console.error('Error fetching partner:', error)
  } finally {
    loading.value = false
  }
}

const fetchAgreements = async () => {
  try {
    const response = await API.partners.getAgreements(route.params.id)
    agreements.value = response.data.data || []
  } catch (error) {
    console.error('Error fetching agreements:', error)
  }
}

const openAddAgreementModal = () => {
  editingAgreement.value = null
  agreementFormData.value = {
    title: '',
    description: '',
    start_date: '',
    end_date: '',
    value: null,
    status: 'active'
  }
  showAgreementModal.value = true
}

const editAgreement = (agreement) => {
  editingAgreement.value = agreement
  agreementFormData.value = {
    title: agreement.title || '',
    description: agreement.description || '',
    start_date: agreement.start_date ? agreement.start_date.split('T')[0] : '',
    end_date: agreement.end_date ? agreement.end_date.split('T')[0] : '',
    value: agreement.value || null,
    status: agreement.status || 'active'
  }
  showAgreementModal.value = true
}

const closeAgreementModal = () => {
  showAgreementModal.value = false
  editingAgreement.value = null
}

const handleAgreementSubmit = async () => {
  savingAgreement.value = true
  try {
    if (editingAgreement.value) {
      await API.partners.updateAgreement(
        route.params.id,
        editingAgreement.value.id,
        agreementFormData.value
      )
      notificationStore.success(t('partners.agreementUpdateSuccess'))
    } else {
      await API.partners.addAgreement(route.params.id, agreementFormData.value)
      notificationStore.success(t('partners.agreementCreateSuccess'))
    }
    closeAgreementModal()
    fetchAgreements()
  } catch (error) {
    notificationStore.error(
      editingAgreement.value
        ? t('partners.agreementUpdateError')
        : t('partners.agreementCreateError')
    )
    console.error('Error saving agreement:', error)
  } finally {
    savingAgreement.value = false
  }
}

const confirmDeleteAgreement = (agreement) => {
  agreementToDelete.value = agreement
  showDeleteAgreementModal.value = true
}

const deleteAgreement = async () => {
  deletingAgreement.value = true
  try {
    await API.partners.deleteAgreement(route.params.id, agreementToDelete.value.id)
    notificationStore.success(t('partners.agreementDeleteSuccess'))
    showDeleteAgreementModal.value = false
    agreementToDelete.value = null
    fetchAgreements()
  } catch (error) {
    notificationStore.error(t('partners.agreementDeleteError'))
    console.error('Error deleting agreement:', error)
  } finally {
    deletingAgreement.value = false
  }
}

const formatDate = (date) => {
  if (!date) return ''
  return new Date(date).toLocaleDateString('pl-PL')
}

const formatCurrency = (amount) => {
  if (amount === null || amount === undefined) return ''
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount)
}

const capitalize = (str) => {
  if (!str) return ''
  return str.charAt(0).toUpperCase() + str.slice(1)
}

// Lifecycle
onMounted(() => {
  fetchPartner()
  fetchAgreements()
})
</script>

<style scoped>
.partner-view-page {
  max-width: 1200px;
  padding: 2rem;
}

.loading-container,
.error-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: 400px;
}

.partner-content {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

/* Page Header */
.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1rem;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.back-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.back-link:hover {
  text-decoration: underline;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin: 0;
}

/* Card */
.card-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin: 0;
}

.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
}

/* Info Grid */
.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
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

.link {
  color: var(--primary-color);
  text-decoration: none;
}

.link:hover {
  text-decoration: underline;
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

.type-badge {
  background: var(--background-secondary);
  color: var(--text-primary);
}

.type-badge.type-veterinary {
  background: #e3f2fd;
  color: #1976d2;
}

.type-badge.type-shelter {
  background: #f3e5f5;
  color: #7b1fa2;
}

.type-badge.type-pet_store {
  background: #fff3e0;
  color: #e65100;
}

.type-badge.type-corporate {
  background: #e8f5e9;
  color: #2e7d32;
}

.type-badge.type-foundation {
  background: #fce4ec;
  color: #c2185b;
}

.type-badge.type-individual {
  background: #e0f2f1;
  color: #00695c;
}

.type-badge.type-other {
  background: #f5f5f5;
  color: #616161;
}

.status-active {
  background: #d4edda;
  color: #155724;
}

.status-inactive {
  background: #f8d7da;
  color: #721c24;
}

.status-pending {
  background: #fff3cd;
  color: #856404;
}

.status-suspended {
  background: #f5f5f5;
  color: #6c757d;
}

/* Agreements List */
.agreements-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.agreement-item {
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background: var(--background-secondary);
}

.agreement-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.75rem;
}

.agreement-title {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.agreement-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: var(--text-primary);
}

.agreement-status {
  font-size: 0.75rem;
}

.agreement-actions {
  display: flex;
  gap: 0.5rem;
}

.agreement-details {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.agreement-description {
  font-size: 0.95rem;
  color: var(--text-secondary);
  margin-bottom: 0.5rem;
}

.agreement-meta {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  font-size: 0.9rem;
  color: var(--text-secondary);
}

/* Notes */
.notes-content {
  white-space: pre-wrap;
  color: var(--text-secondary);
  line-height: 1.6;
}

/* Agreement Form */
.agreement-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
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
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .header-left {
    flex-wrap: wrap;
  }

  .info-grid {
    grid-template-columns: 1fr;
  }

  .agreement-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }

  .agreement-actions {
    width: 100%;
  }

  .agreement-actions button {
    flex: 1;
  }
}
</style>
