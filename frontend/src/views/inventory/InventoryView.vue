<template>
  <div class="inventory-view-page">
    <LoadingSpinner v-if="loading" />
    <div v-else>
      <div class="page-header">
        <div>
          <h1 class="page-title">{{ item.name }}</h1>
          <span class="badge" :class="`badge-${item.category}`">
            {{ t(`inventory.category${item.category?.charAt(0).toUpperCase() + item.category?.slice(1)}`) }}
          </span>
        </div>
        <div class="header-actions">
          <BaseButton
            v-if="authStore.hasRole('staff')"
            variant="primary"
            @click="showAddMovementModal = true"
          >
            ➕ {{ t('inventory.addMovement') }}
          </BaseButton>
          <BaseButton
            v-if="authStore.hasRole('staff')"
            variant="secondary"
            @click="editItem"
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

      <!-- Stock Status Cards -->
      <div class="stats-grid">
        <div class="stat-card" :class="getStockStatusClass()">
          <div class="stat-icon">📦</div>
          <div class="stat-info">
            <div class="stat-value">{{ item.stock_level }} {{ item.unit }}</div>
            <div class="stat-label">{{ t('inventory.currentStock') }}</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">⚠️</div>
          <div class="stat-info">
            <div class="stat-value">{{ item.min_stock }} {{ item.unit }}</div>
            <div class="stat-label">{{ t('inventory.minStock') }}</div>
          </div>
        </div>
        <div v-if="item.max_stock" class="stat-card">
          <div class="stat-icon">📊</div>
          <div class="stat-info">
            <div class="stat-value">{{ item.max_stock }} {{ item.unit }}</div>
            <div class="stat-label">{{ t('inventory.maxStock') }}</div>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon">💰</div>
          <div class="stat-info">
            <div class="stat-value">{{ formatCurrency((item.stock_level || 0) * (item.unit_price || 0)) }}</div>
            <div class="stat-label">{{ t('inventory.totalValue') }}</div>
          </div>
        </div>
      </div>

      <!-- Item Information -->
      <div class="info-grid">
        <!-- Basic Information -->
        <BaseCard>
          <template #header>{{ t('inventory.basicInfo') }}</template>
          <div v-if="item.description" class="info-item">
            <span class="info-label">{{ t('common.description') }}:</span>
            <span class="info-value">{{ item.description }}</span>
          </div>
          <div v-if="item.supplier" class="info-item">
            <span class="info-label">{{ t('inventory.supplier') }}:</span>
            <span class="info-value">{{ item.supplier }}</span>
          </div>
          <div v-if="item.location" class="info-item">
            <span class="info-label">{{ t('inventory.location') }}:</span>
            <span class="info-value">{{ item.location }}</span>
          </div>
          <div v-if="item.unit_price" class="info-item">
            <span class="info-label">{{ t('inventory.unitPrice') }}:</span>
            <span class="info-value">{{ formatCurrency(item.unit_price) }}</span>
          </div>
          <div v-if="item.expiry_date" class="info-item">
            <span class="info-label">{{ t('inventory.expiryDate') }}:</span>
            <span class="info-value" :class="getExpiryClass(item.expiry_date)">
              {{ formatDate(item.expiry_date) }}
            </span>
          </div>
        </BaseCard>

        <!-- Notes -->
        <BaseCard v-if="item.notes">
          <template #header>{{ t('common.notes') }}</template>
          <p class="notes-content">{{ item.notes }}</p>
        </BaseCard>
      </div>

      <!-- Stock Movements History -->
      <BaseCard>
        <template #header>{{ t('inventory.movementHistory') }}</template>
        <LoadingSpinner v-if="loadingMovements" />
        <EmptyState
          v-else-if="!movements || movements.length === 0"
          icon="📋"
          :title="t('inventory.noMovements')"
          :description="t('inventory.noMovementsMessage')"
        />
        <div v-else class="movements-list">
          <div
            v-for="movement in movements"
            :key="movement.id"
            class="movement-item"
            :class="`movement-${movement.type}`"
          >
            <div class="movement-header">
              <div class="movement-type">
                <span class="type-badge" :class="`badge-${movement.type}`">
                  {{ t(`inventory.movement${movement.type?.charAt(0).toUpperCase() + movement.type?.slice(1)}`) }}
                </span>
                <span class="movement-quantity">
                  {{ movement.type === 'out' ? '-' : '+' }}{{ movement.quantity }} {{ item.unit }}
                </span>
              </div>
              <div class="movement-date">{{ formatDate(movement.date) }}</div>
            </div>
            <div v-if="movement.reason" class="movement-reason">
              {{ movement.reason }}
            </div>
            <div v-if="movement.notes" class="movement-notes">
              {{ movement.notes }}
            </div>
            <div v-if="movement.created_by" class="movement-meta">
              {{ t('common.by') }}: {{ movement.created_by.full_name || movement.created_by.email }}
            </div>
          </div>
        </div>
      </BaseCard>
    </div>

    <!-- Add Movement Modal -->
    <BaseModal
      v-if="showAddMovementModal"
      :title="t('inventory.addMovement')"
      size="medium"
      @close="showAddMovementModal = false"
    >
      <form @submit.prevent="addMovement">
        <FormGroup :label="t('inventory.movementType')" :error="movementErrors.type" required>
          <select v-model="movementForm.type" class="form-control" :class="{ 'error': movementErrors.type }">
            <option value="">{{ t('inventory.selectMovementType') }}</option>
            <option value="in">{{ t('inventory.movementIn') }}</option>
            <option value="out">{{ t('inventory.movementOut') }}</option>
            <option value="adjustment">{{ t('inventory.movementAdjustment') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('inventory.quantity')" :error="movementErrors.quantity" required>
          <input
            v-model.number="movementForm.quantity"
            type="number"
            min="1"
            class="form-control"
            :class="{ 'error': movementErrors.quantity }"
            :placeholder="t('inventory.quantityPlaceholder')"
          />
        </FormGroup>

        <FormGroup :label="t('inventory.reason')" :error="movementErrors.reason">
          <input
            v-model="movementForm.reason"
            type="text"
            class="form-control"
            :placeholder="t('inventory.reasonPlaceholder')"
          />
        </FormGroup>

        <FormGroup :label="t('common.notes')" :error="movementErrors.notes">
          <textarea
            v-model="movementForm.notes"
            class="form-control"
            :placeholder="t('inventory.movementNotesPlaceholder')"
            rows="3"
          ></textarea>
        </FormGroup>
      </form>

      <template #footer>
        <BaseButton variant="secondary" @click="showAddMovementModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="primary" @click="addMovement" :disabled="submittingMovement">
          {{ submittingMovement ? t('common.saving') : t('common.save') }}
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-if="showDeleteModal"
      :title="t('inventory.deleteItem')"
      @close="showDeleteModal = false"
    >
      <p>{{ t('inventory.deleteItemConfirm') }}</p>
      <template #footer>
        <BaseButton variant="secondary" @click="showDeleteModal = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" @click="deleteItem">
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

const item = ref({})
const movements = ref([])
const loading = ref(false)
const loadingMovements = ref(false)
const showDeleteModal = ref(false)
const showAddMovementModal = ref(false)
const submittingMovement = ref(false)

const movementForm = reactive({
  type: '',
  quantity: null,
  reason: '',
  notes: '',
})

const movementErrors = reactive({})

async function fetchItem() {
  try {
    loading.value = true
    const response = await API.inventory.getById(route.params.id)
    item.value = response.data
  } catch (error) {
    console.error('Failed to fetch inventory item:', error)
    notificationStore.error(t('inventory.fetchError'))
    router.push({ name: 'Inventory' })
  } finally {
    loading.value = false
  }
}

async function fetchMovements() {
  try {
    loadingMovements.value = true
    const response = await API.inventory.getMovements(route.params.id)
    movements.value = response.data.data || response.data || []
  } catch (error) {
    console.error('Failed to fetch movements:', error)
    notificationStore.error(t('inventory.fetchMovementsError'))
  } finally {
    loadingMovements.value = false
  }
}

function validateMovementForm() {
  Object.keys(movementErrors).forEach(key => delete movementErrors[key])
  let isValid = true

  if (!movementForm.type) {
    movementErrors.type = t('common.required')
    isValid = false
  }

  if (!movementForm.quantity || movementForm.quantity <= 0) {
    movementErrors.quantity = t('common.required')
    isValid = false
  }

  return isValid
}

async function addMovement() {
  if (!validateMovementForm()) {
    notificationStore.error(t('common.fixErrors'))
    return
  }

  try {
    submittingMovement.value = true
    await API.inventory.addMovement(route.params.id, movementForm)
    notificationStore.success(t('inventory.movementAddSuccess'))
    showAddMovementModal.value = false

    // Reset form
    movementForm.type = ''
    movementForm.quantity = null
    movementForm.reason = ''
    movementForm.notes = ''

    // Refresh data
    fetchItem()
    fetchMovements()
  } catch (error) {
    console.error('Failed to add movement:', error)
    notificationStore.error(t('inventory.movementAddError'))
  } finally {
    submittingMovement.value = false
  }
}

function editItem() {
  router.push({ name: 'InventoryForm', params: { id: route.params.id } })
}

async function deleteItem() {
  try {
    await API.inventory.delete(route.params.id)
    notificationStore.success(t('inventory.deleteSuccess'))
    router.push({ name: 'Inventory' })
  } catch (error) {
    console.error('Failed to delete inventory item:', error)
    notificationStore.error(t('inventory.deleteError'))
  }
}

function getStockStatusClass() {
  if (item.value.stock_level === 0) return 'danger'
  if (item.value.stock_level <= item.value.min_stock) return 'warning'
  return 'success'
}

function getExpiryClass(date) {
  if (!date) return ''
  const expiryDate = new Date(date)
  const today = new Date()
  const daysUntilExpiry = Math.ceil((expiryDate - today) / (1000 * 60 * 60 * 24))

  if (daysUntilExpiry < 0) return 'expired'
  if (daysUntilExpiry <= 30) return 'expiring-soon'
  return ''
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
  fetchItem()
  fetchMovements()
})
</script>

<style scoped>
.inventory-view-page {
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
}

.badge-food {
  background: #fff3e0;
  color: #f57c00;
}

.badge-medicine {
  background: #f3e5f5;
  color: #7b1fa2;
}

.badge-supplies {
  background: #e3f2fd;
  color: #1976d2;
}

.badge-equipment {
  background: #e8f5e9;
  color: #388e3c;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
  background: var(--bg-secondary);
  border-radius: 8px;
  border-left: 4px solid var(--primary-color);
}

.stat-card.warning {
  border-left-color: #ff9800;
}

.stat-card.danger {
  border-left-color: #f44336;
}

.stat-card.success {
  border-left-color: #4caf50;
}

.stat-icon {
  font-size: 2rem;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--text-primary);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
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
  min-width: 150px;
  color: var(--text-secondary);
}

.info-value {
  color: var(--text-primary);
}

.notes-content {
  white-space: pre-wrap;
  line-height: 1.6;
  color: var(--text-primary);
}

.expiring-soon {
  color: #f57c00;
  font-weight: 600;
}

.expired {
  color: #c62828;
  font-weight: 600;
  text-decoration: line-through;
}

.movements-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.movement-item {
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-radius: 8px;
  background: var(--bg-secondary);
  border-left-width: 4px;
}

.movement-item.movement-in {
  border-left-color: #4caf50;
}

.movement-item.movement-out {
  border-left-color: #f44336;
}

.movement-item.movement-adjustment {
  border-left-color: #ff9800;
}

.movement-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.5rem;
}

.movement-type {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.type-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  font-weight: 600;
  text-transform: uppercase;
}

.type-badge.badge-in {
  background: #e8f5e9;
  color: #388e3c;
}

.type-badge.badge-out {
  background: #ffebee;
  color: #c62828;
}

.type-badge.badge-adjustment {
  background: #fff3e0;
  color: #f57c00;
}

.movement-quantity {
  font-weight: 600;
  font-size: 1.125rem;
}

.movement-date {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.movement-reason {
  font-weight: 500;
  margin-bottom: 0.25rem;
}

.movement-notes {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-style: italic;
  margin-bottom: 0.5rem;
}

.movement-meta {
  font-size: 0.75rem;
  color: var(--text-secondary);
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
