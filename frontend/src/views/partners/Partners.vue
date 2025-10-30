<template>
  <div class="partners-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.partners') }}</h1>
      <BaseButton variant="primary" @click="openCreateModal">
        ➕ {{ t('partners.addPartner') }}
      </BaseButton>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">🤝</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.totalPartners }}</div>
            <div class="stat-label">{{ t('partners.totalPartners') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">✅</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.activePartners }}</div>
            <div class="stat-label">{{ t('partners.activePartners') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">📄</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.activeAgreements }}</div>
            <div class="stat-label">{{ t('partners.activeAgreements') }}</div>
          </div>
        </div>
      </BaseCard>

      <BaseCard class="stat-card">
        <div class="stat-content">
          <div class="stat-icon">⏰</div>
          <div class="stat-details">
            <div class="stat-value">{{ statistics.expiringAgreements }}</div>
            <div class="stat-label">{{ t('partners.expiringAgreements') }}</div>
          </div>
        </div>
      </BaseCard>
    </div>

    <!-- Filters -->
    <BaseCard>
      <div class="filters">
        <div class="filter-group">
          <label>{{ t('partners.partnerType') }}</label>
          <select v-model="filters.type" class="filter-select">
            <option value="">{{ t('common.all') }}</option>
            <option value="veterinary">{{ t('partners.typeVeterinary') }}</option>
            <option value="shelter">{{ t('partners.typeShelter') }}</option>
            <option value="pet_store">{{ t('partners.typePetStore') }}</option>
            <option value="corporate">{{ t('partners.typeCorporate') }}</option>
            <option value="foundation">{{ t('partners.typeFoundation') }}</option>
            <option value="individual">{{ t('partners.typeIndividual') }}</option>
            <option value="other">{{ t('partners.typeOther') }}</option>
          </select>
        </div>

        <div class="filter-group">
          <label>{{ t('common.status') }}</label>
          <select v-model="filters.status" class="filter-select">
            <option value="">{{ t('common.all') }}</option>
            <option value="active">{{ t('partners.statusActive') }}</option>
            <option value="inactive">{{ t('partners.statusInactive') }}</option>
            <option value="pending">{{ t('partners.statusPending') }}</option>
            <option value="suspended">{{ t('partners.statusSuspended') }}</option>
          </select>
        </div>

        <div class="filter-group">
          <label>{{ t('common.search') }}</label>
          <input
            v-model="filters.search"
            type="text"
            :placeholder="t('partners.searchPlaceholder')"
            class="filter-input"
          />
        </div>
      </div>
    </BaseCard>

    <!-- Partners Table -->
    <BaseCard>
      <DataTable
        :columns="columns"
        :data="filteredPartners"
        :loading="loading"
        @sort="handleSort"
      >
        <template #cell-name="{ row }">
          <router-link :to="`/partners/${row.id}`" class="partner-link">
            {{ row.name }}
          </router-link>
        </template>

        <template #cell-type="{ row }">
          <span :class="['badge', 'type-badge', `type-${row.type}`]">
            {{ t(`partners.type${capitalize(row.type)}`) }}
          </span>
        </template>

        <template #cell-status="{ row }">
          <span :class="['badge', `status-${row.status}`]">
            {{ t(`partners.status${capitalize(row.status)}`) }}
          </span>
        </template>

        <template #cell-agreements="{ row }">
          <span class="agreements-count">
            {{ row.active_agreements || 0 }} {{ t('partners.active') }}
          </span>
        </template>

        <template #cell-contact="{ row }">
          <div class="contact-info">
            <div v-if="row.contact_person">👤 {{ row.contact_person }}</div>
            <div v-if="row.email">✉️ {{ row.email }}</div>
            <div v-if="row.phone">📞 {{ row.phone }}</div>
          </div>
        </template>

        <template #cell-actions="{ row }">
          <div class="actions">
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
      :title="editingPartner ? t('partners.editPartner') : t('partners.addPartner')"
      size="large"
      @close="closeModal"
    >
      <PartnerForm
        :partner="editingPartner"
        @submit="handleSubmit"
        @cancel="closeModal"
      />
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      :title="t('partners.deletePartner')"
      size="small"
      @close="showDeleteModal = false"
    >
      <p>{{ t('partners.deletePartnerConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deletePartner" :loading="deleting">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useRouter } from 'vue-router'
import { API } from '../../api'
import { useNotificationStore } from '../../stores/notification'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import DataTable from '../../components/base/DataTable.vue'
import PartnerForm from './PartnerForm.vue'

const { t } = useI18n()
const router = useRouter()
const notificationStore = useNotificationStore()

// State
const partners = ref([])
const statistics = ref({
  totalPartners: 0,
  activePartners: 0,
  activeAgreements: 0,
  expiringAgreements: 0
})
const loading = ref(false)
const filters = ref({
  type: '',
  status: '',
  search: ''
})

// Modal state
const showModal = ref(false)
const showDeleteModal = ref(false)
const editingPartner = ref(null)
const partnerToDelete = ref(null)
const deleting = ref(false)

// Table columns
const columns = [
  { key: 'name', label: t('partners.name'), sortable: true },
  { key: 'type', label: t('partners.partnerType'), sortable: true },
  { key: 'status', label: t('common.status'), sortable: true },
  { key: 'agreements', label: t('partners.agreements'), sortable: false },
  { key: 'contact', label: t('partners.contactInfo'), sortable: false },
  { key: 'actions', label: t('common.actions'), sortable: false }
]

// Computed
const filteredPartners = computed(() => {
  let result = [...partners.value]

  // Filter by type
  if (filters.value.type) {
    result = result.filter(p => p.type === filters.value.type)
  }

  // Filter by status
  if (filters.value.status) {
    result = result.filter(p => p.status === filters.value.status)
  }

  // Filter by search
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter(p =>
      p.name?.toLowerCase().includes(search) ||
      p.contact_person?.toLowerCase().includes(search) ||
      p.email?.toLowerCase().includes(search)
    )
  }

  return result
})

// Methods
const fetchPartners = async () => {
  loading.value = true
  try {
    const response = await API.partners.list()
    partners.value = response.data.data || []
  } catch (error) {
    notificationStore.error(t('partners.fetchError'))
    console.error('Error fetching partners:', error)
  } finally {
    loading.value = false
  }
}

const fetchStatistics = async () => {
  try {
    const response = await API.partners.getOverallStatistics()
    statistics.value = response.data.data || statistics.value
  } catch (error) {
    console.error('Error fetching statistics:', error)
  }
}

const openCreateModal = () => {
  editingPartner.value = null
  showModal.value = true
}

const openEditModal = (partner) => {
  editingPartner.value = partner
  showModal.value = true
}

const closeModal = () => {
  showModal.value = false
  editingPartner.value = null
}

const handleSubmit = async (partnerData) => {
  try {
    if (editingPartner.value) {
      await API.partners.update(editingPartner.value.id, partnerData)
      notificationStore.success(t('partners.updateSuccess'))
    } else {
      await API.partners.create(partnerData)
      notificationStore.success(t('partners.createSuccess'))
    }
    closeModal()
    fetchPartners()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(
      editingPartner.value ? t('partners.updateError') : t('partners.createError')
    )
    console.error('Error saving partner:', error)
  }
}

const confirmDelete = (partner) => {
  partnerToDelete.value = partner
  showDeleteModal.value = true
}

const deletePartner = async () => {
  deleting.value = true
  try {
    await API.partners.delete(partnerToDelete.value.id)
    notificationStore.success(t('partners.deleteSuccess'))
    showDeleteModal.value = false
    partnerToDelete.value = null
    fetchPartners()
    fetchStatistics()
  } catch (error) {
    notificationStore.error(t('partners.deleteError'))
    console.error('Error deleting partner:', error)
  } finally {
    deleting.value = false
  }
}

const handleSort = (column, direction) => {
  // Implement sorting logic
  partners.value.sort((a, b) => {
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
  return str.charAt(0).toUpperCase() + str.slice(1)
}

// Lifecycle
onMounted(() => {
  fetchPartners()
  fetchStatistics()
})
</script>

<style scoped>
.partners-page {
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

/* Partner Link */
.partner-link {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.partner-link:hover {
  text-decoration: underline;
}

/* Badges */
.badge {
  padding: 0.25rem 0.75rem;
  border-radius: 1rem;
  font-size: 0.85rem;
  font-weight: 500;
  display: inline-block;
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

/* Contact Info */
.contact-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  font-size: 0.85rem;
}

.agreements-count {
  font-weight: 500;
}

/* Actions */
.actions {
  display: flex;
  gap: 0.5rem;
}

/* Dark mode support */
@media (prefers-color-scheme: dark) {
  .type-badge.type-veterinary {
    background: #1a237e;
    color: #bbdefb;
  }

  .type-badge.type-shelter {
    background: #4a148c;
    color: #e1bee7;
  }

  .type-badge.type-pet_store {
    background: #e65100;
    color: #ffe0b2;
  }

  .type-badge.type-corporate {
    background: #1b5e20;
    color: #c8e6c9;
  }

  .type-badge.type-foundation {
    background: #880e4f;
    color: #f8bbd0;
  }

  .type-badge.type-individual {
    background: #004d40;
    color: #b2dfdb;
  }
}
</style>
