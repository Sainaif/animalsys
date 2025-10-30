<template>
  <div class="donor-view-page">
    <LoadingSpinner v-if="loading" />
    <div v-else>
      <div class="page-header">
        <div>
          <h1 class="page-title">{{ donor.name }}</h1>
          <span class="badge" :class="`badge-${donor.type}`">
            {{ t(`donors.type${donor.type.charAt(0).toUpperCase() + donor.type.slice(1)}`) }}
          </span>
          <span class="badge" :class="`badge-${donor.status}`">
            {{ t(`donors.status${donor.status.charAt(0).toUpperCase() + donor.status.slice(1)}`) }}
          </span>
        </div>
        <div class="header-actions">
          <BaseButton
            v-if="authStore.hasRole('staff')"
            variant="primary"
            @click="showAddDonationModal = true"
          >
            ➕ {{ t('donors.addDonation') }}
          </BaseButton>
          <BaseButton
            v-if="authStore.hasRole('staff')"
            variant="secondary"
            @click="editDonor"
          >
            {{ t('common.edit') }}
          </BaseButton>
          <BaseButton
            v-if="authStore.hasRole('admin')"
            variant="danger"
            @click="showDeleteModal = true"
          >
            {{ t('common.delete') }}
          </BaseButton>
        </div>
      </div>

      <!-- Donor Information -->
      <div class="info-grid">
        <!-- Contact Information -->
        <BaseCard>
          <template #header>{{ t('common.contactInfo') }}</template>
          <div class="info-item">
            <span class="info-label">{{ t('common.email') }}:</span>
            <span class="info-value">{{ donor.email || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('common.phone') }}:</span>
            <span class="info-value">{{ donor.phone || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('common.address') }}:</span>
            <span class="info-value">{{ donor.address || '-' }}</span>
          </div>
          <div v-if="donor.tax_id" class="info-item">
            <span class="info-label">{{ t('donors.taxId') }}:</span>
            <span class="info-value">{{ donor.tax_id }}</span>
          </div>
          <div v-if="donor.preferred_contact_method" class="info-item">
            <span class="info-label">{{ t('donors.preferredContactMethod') }}:</span>
            <span class="info-value">{{ donor.preferred_contact_method }}</span>
          </div>
        </BaseCard>

        <!-- Donation Statistics -->
        <BaseCard>
          <template #header>{{ t('donors.donationStats') }}</template>
          <div class="stats-grid">
            <div class="stat-card">
              <div class="stat-value">{{ formatCurrency(donor.total_donated || 0) }}</div>
              <div class="stat-label">{{ t('donors.totalDonated') }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ donor.donation_count || 0 }}</div>
              <div class="stat-label">{{ t('donors.donationCount') }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">
                {{ donor.last_donation_date ? formatDate(donor.last_donation_date) : '-' }}
              </div>
              <div class="stat-label">{{ t('donors.lastDonation') }}</div>
            </div>
            <div class="stat-card">
              <div class="stat-value">{{ formatCurrency(donor.average_donation || 0) }}</div>
              <div class="stat-label">{{ t('donors.averageDonation') }}</div>
            </div>
          </div>
        </BaseCard>
      </div>

      <!-- Notes -->
      <BaseCard v-if="donor.notes">
        <template #header>{{ t('donors.notes') }}</template>
        <p class="notes-content">{{ donor.notes }}</p>
      </BaseCard>

      <!-- Donation History -->
      <BaseCard>
        <template #header>{{ t('donors.donationHistory') }}</template>
        <LoadingSpinner v-if="loadingDonations" />
        <EmptyState
          v-else-if="!donations || donations.length === 0"
          icon="🎁"
          :title="t('donors.noDonations')"
          :description="t('donors.noDonationsMessage')"
        />
        <div v-else class="donations-list">
          <div
            v-for="donation in donations"
            :key="donation.id"
            class="donation-item"
          >
            <div class="donation-info">
              <div class="donation-amount">{{ formatCurrency(donation.amount) }}</div>
              <div class="donation-date">{{ formatDate(donation.date) }}</div>
            </div>
            <div class="donation-details">
              <span v-if="donation.method" class="donation-method">
                {{ t(`donors.method${donation.method.charAt(0).toUpperCase() + donation.method.slice(1)}`) }}
              </span>
              <span v-if="donation.purpose" class="donation-purpose">
                {{ donation.purpose }}
              </span>
            </div>
            <div v-if="donation.notes" class="donation-notes">
              {{ donation.notes }}
            </div>
          </div>
        </div>
      </BaseCard>
    </div>

    <!-- Add Donation Modal -->
    <BaseModal
      v-if="showAddDonationModal"
      :title="t('donors.addDonation')"
      size="medium"
      @close="showAddDonationModal = false"
    >
      <form @submit.prevent="addDonation">
        <FormGroup :label="t('donors.amount')" :error="donationErrors.amount" required>
          <input
            v-model.number="donationForm.amount"
            type="number"
            step="0.01"
            min="0"
            class="form-control"
            :class="{ 'error': donationErrors.amount }"
            :placeholder="t('donors.amountPlaceholder')"
          />
        </FormGroup>

        <FormGroup :label="t('donors.date')" :error="donationErrors.date" required>
          <input
            v-model="donationForm.date"
            type="date"
            class="form-control"
            :class="{ 'error': donationErrors.date }"
          />
        </FormGroup>

        <FormGroup :label="t('donors.method')" :error="donationErrors.method">
          <select v-model="donationForm.method" class="form-control">
            <option value="">{{ t('donors.selectMethod') }}</option>
            <option value="cash">{{ t('donors.methodCash') }}</option>
            <option value="transfer">{{ t('donors.methodTransfer') }}</option>
            <option value="card">{{ t('donors.methodCard') }}</option>
            <option value="check">{{ t('donors.methodCheck') }}</option>
            <option value="other">{{ t('donors.methodOther') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('donors.purpose')" :error="donationErrors.purpose">
          <input
            v-model="donationForm.purpose"
            type="text"
            class="form-control"
            :placeholder="t('donors.purposePlaceholder')"
          />
        </FormGroup>

        <FormGroup :label="t('donors.notes')" :error="donationErrors.notes">
          <textarea
            v-model="donationForm.notes"
            class="form-control"
            :placeholder="t('donors.notesPlaceholder')"
            rows="3"
          ></textarea>
        </FormGroup>
      </form>

      <template #footer>
        <BaseButton variant="secondary" @click="showAddDonationModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="primary" @click="addDonation" :disabled="submittingDonation">
          {{ submittingDonation ? t('common.saving') : t('common.save') }}
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      :title="t('donors.deleteDonor')"
      @close="showDeleteModal = false"
    >
      <p>{{ t('donors.deleteDonorConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteDonor">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notifications'
import { API } from '../../api'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import FormGroup from '../../components/base/FormGroup.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'
import EmptyState from '../../components/base/EmptyState.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const donor = ref({})
const donations = ref([])
const loading = ref(false)
const loadingDonations = ref(false)
const showDeleteModal = ref(false)
const showAddDonationModal = ref(false)
const submittingDonation = ref(false)

const donationForm = reactive({
  amount: null,
  date: new Date().toISOString().split('T')[0],
  method: '',
  purpose: '',
  notes: '',
})

const donationErrors = reactive({})

async function fetchDonor() {
  try {
    loading.value = true
    const response = await API.donors.getById(route.params.id)
    donor.value = response.data
  } catch (error) {
    console.error('Failed to fetch donor:', error)
    notificationStore.error(t('donors.fetchError'))
    router.push({ name: 'Donors' })
  } finally {
    loading.value = false
  }
}

async function fetchDonations() {
  try {
    loadingDonations.value = true
    const response = await API.donors.getDonations(route.params.id)
    donations.value = response.data.data || response.data || []
  } catch (error) {
    console.error('Failed to fetch donations:', error)
    notificationStore.error(t('donors.fetchDonationsError'))
  } finally {
    loadingDonations.value = false
  }
}

function validateDonationForm() {
  Object.keys(donationErrors).forEach(key => delete donationErrors[key])
  let isValid = true

  if (!donationForm.amount || donationForm.amount <= 0) {
    donationErrors.amount = t('common.required')
    isValid = false
  }

  if (!donationForm.date) {
    donationErrors.date = t('common.required')
    isValid = false
  }

  return isValid
}

async function addDonation() {
  if (!validateDonationForm()) {
    notificationStore.error(t('common.fixErrors'))
    return
  }

  try {
    submittingDonation.value = true
    await API.donors.addDonation(route.params.id, donationForm)
    notificationStore.success(t('donors.donationAddSuccess'))
    showAddDonationModal.value = false

    // Reset form
    donationForm.amount = null
    donationForm.date = new Date().toISOString().split('T')[0]
    donationForm.method = ''
    donationForm.purpose = ''
    donationForm.notes = ''

    // Refresh data
    fetchDonor()
    fetchDonations()
  } catch (error) {
    console.error('Failed to add donation:', error)
    notificationStore.error(t('donors.donationAddError'))
  } finally {
    submittingDonation.value = false
  }
}

function editDonor() {
  router.push({ name: 'DonorForm', params: { id: route.params.id } })
}

async function deleteDonor() {
  try {
    await API.donors.delete(route.params.id)
    notificationStore.success(t('donors.deleteSuccess'))
    router.push({ name: 'Donors' })
  } catch (error) {
    console.error('Failed to delete donor:', error)
    notificationStore.error(t('donors.deleteError'))
  }
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
}

function formatDate(date) {
  if (!date) return '-'
  return new Date(date).toLocaleDateString('pl-PL')
}

onMounted(() => {
  fetchDonor()
  fetchDonations()
})
</script>

<style scoped>
.donor-view-page {
  max-width: 1200px;
  padding: 2rem;
}

.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin: 0 0 0.5rem 0;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
  margin-right: 0.5rem;
}

.badge-individual {
  background: #e3f2fd;
  color: #1976d2;
}

.badge-company {
  background: #f3e5f5;
  color: #7b1fa2;
}

.badge-foundation {
  background: #fff3e0;
  color: #f57c00;
}

.badge-active {
  background: #d4edda;
  color: #155724;
}

.badge-inactive {
  background: #f8d7da;
  color: #721c24;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.info-item {
  display: flex;
  padding: 0.75rem 0;
  border-bottom: 1px solid var(--border-color);
}

.info-item:last-child {
  border-bottom: none;
}

.info-label {
  font-weight: 600;
  min-width: 180px;
  color: var(--text-secondary);
}

.info-value {
  color: var(--text-primary);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
  gap: 1rem;
}

.stat-card {
  padding: 1rem;
  background: var(--bg-secondary);
  border-radius: 8px;
  text-align: center;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--primary-color);
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.notes-content {
  white-space: pre-wrap;
  line-height: 1.6;
  color: var(--text-primary);
}

.donations-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.donation-item {
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-secondary);
}

.donation-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.donation-amount {
  font-size: 1.25rem;
  font-weight: bold;
  color: var(--success-color);
}

.donation-date {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.donation-details {
  display: flex;
  gap: 1rem;
  margin-bottom: 0.5rem;
}

.donation-method,
.donation-purpose {
  font-size: 0.875rem;
  padding: 0.25rem 0.5rem;
  background: var(--bg-primary);
  border-radius: 4px;
}

.donation-notes {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-style: italic;
  margin-top: 0.5rem;
}

.form-control {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-control.error {
  border-color: var(--danger-color);
}
</style>
