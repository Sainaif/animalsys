<template>
  <div class="inventory-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.inventory') }}</h1>
      <BaseButton
        v-if="authStore.hasRole('staff')"
        variant="primary"
        @click="createItem"
      >
        ➕ {{ t('inventory.addItem') }}
      </BaseButton>
    </div>

    <!-- Statistics Cards -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">📦</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.total_items || 0 }}</div>
          <div class="stat-label">{{ t('inventory.totalItems') }}</div>
        </div>
      </div>
      <div class="stat-card warning">
        <div class="stat-icon">⚠️</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.low_stock_count || 0 }}</div>
          <div class="stat-label">{{ t('inventory.lowStockItems') }}</div>
        </div>
      </div>
      <div class="stat-card danger">
        <div class="stat-icon">⏰</div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.expiring_soon_count || 0 }}</div>
          <div class="stat-label">{{ t('inventory.expiringSoon') }}</div>
        </div>
      </div>
      <div class="stat-card success">
        <div class="stat-icon">💰</div>
        <div class="stat-info">
          <div class="stat-value">{{ formatCurrency(stats.total_value || 0) }}</div>
          <div class="stat-label">{{ t('inventory.totalValue') }}</div>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <BaseCard class="filters-card">
      <div class="filters">
        <FormGroup :label="t('common.search')">
          <input
            v-model="filters.search"
            type="text"
            class="form-control"
            :placeholder="t('inventory.searchPlaceholder')"
            @input="handleFilterChange"
          />
        </FormGroup>

        <FormGroup :label="t('inventory.category')">
          <select v-model="filters.category" class="form-control" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="food">{{ t('inventory.categoryFood') }}</option>
            <option value="medicine">{{ t('inventory.categoryMedicine') }}</option>
            <option value="supplies">{{ t('inventory.categorySupplies') }}</option>
            <option value="equipment">{{ t('inventory.categoryEquipment') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('common.status')">
          <select v-model="filters.status" class="form-control" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="in_stock">{{ t('inventory.statusInStock') }}</option>
            <option value="low_stock">{{ t('inventory.statusLowStock') }}</option>
            <option value="out_of_stock">{{ t('inventory.statusOutOfStock') }}</option>
          </select>
        </FormGroup>
      </div>
    </BaseCard>

    <!-- Inventory Table -->
    <BaseCard>
      <LoadingSpinner v-if="loading" />
      <EmptyState
        v-else-if="!items || items.length === 0"
        icon="📦"
        :title="t('inventory.noItems')"
        :description="t('inventory.noItemsMessage')"
      />
      <DataTable
        v-else
        :columns="columns"
        :data="items"
        :current-page="pagination.page"
        :total-pages="pagination.total_pages"
        @sort="handleSort"
        @page-change="handlePageChange"
      >
        <template #cell-category="{ row }">
          <span class="badge" :class="`badge-${row.category}`">
            {{ t(`inventory.category${row.category.charAt(0).toUpperCase() + row.category.slice(1)}`) }}
          </span>
        </template>

        <template #cell-stock_level="{ row }">
          <div class="stock-level">
            <span
              class="stock-badge"
              :class="{
                'low': row.stock_level <= row.min_stock && row.stock_level > 0,
                'out': row.stock_level === 0
              }"
            >
              {{ row.stock_level }}
            </span>
            <span class="stock-unit">{{ row.unit }}</span>
          </div>
        </template>

        <template #cell-expiry_date="{ row }">
          <span
            v-if="row.expiry_date"
            :class="{
              'expiring-soon': isExpiringSoon(row.expiry_date),
              'expired': isExpired(row.expiry_date)
            }"
          >
            {{ formatDate(row.expiry_date) }}
          </span>
          <span v-else>-</span>
        </template>

        <template #cell-value="{ row }">
          {{ formatCurrency((row.stock_level || 0) * (row.unit_price || 0)) }}
        </template>

        <template #cell-actions="{ row }">
          <div class="actions">
            <BaseButton
              variant="secondary"
              size="small"
              @click="viewItem(row.id)"
            >
              {{ t('common.view') }}
            </BaseButton>
            <BaseButton
              v-if="authStore.hasRole('staff')"
              variant="secondary"
              size="small"
              @click="editItem(row.id)"
            >
              {{ t('common.edit') }}
            </BaseButton>
            <BaseButton
              v-if="authStore.hasRole('admin')"
              variant="danger"
              size="small"
              @click="confirmDelete(row.id)"
            >
              {{ t('common.delete') }}
            </BaseButton>
          </div>
        </template>
      </DataTable>
    </BaseCard>

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
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notifications'
import { API } from '../../api'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import DataTable from '../../components/base/DataTable.vue'
import FormGroup from '../../components/base/FormGroup.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'
import EmptyState from '../../components/base/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const items = ref([])
const stats = ref({})
const loading = ref(false)
const showDeleteModal = ref(false)
const itemToDelete = ref(null)

const filters = reactive({
  search: '',
  category: '',
  status: '',
})

const pagination = reactive({
  page: 1,
  limit: 10,
  total: 0,
  total_pages: 0,
})

const sort = reactive({
  field: 'name',
  order: 'asc',
})

const columns = [
  { key: 'name', label: t('inventory.itemName'), sortable: true },
  { key: 'category', label: t('inventory.category'), sortable: true },
  { key: 'stock_level', label: t('inventory.stockLevel'), sortable: true },
  { key: 'min_stock', label: t('inventory.minStock'), sortable: false },
  { key: 'expiry_date', label: t('inventory.expiryDate'), sortable: true },
  { key: 'value', label: t('inventory.value'), sortable: false },
  { key: 'actions', label: t('common.actions'), sortable: false },
]

async function fetchItems() {
  try {
    loading.value = true
    const params = {
      page: pagination.page,
      limit: pagination.limit,
      sort_by: sort.field,
      sort_order: sort.order,
      ...filters,
    }

    const response = await API.inventory.list(params)
    items.value = response.data.data || []
    pagination.total = response.data.total || 0
    pagination.total_pages = response.data.total_pages || 0
  } catch (error) {
    console.error('Failed to fetch inventory items:', error)
    notificationStore.error(t('inventory.fetchError'))
  } finally {
    loading.value = false
  }
}

async function fetchStatistics() {
  try {
    const response = await API.inventory.getStatistics()
    stats.value = response.data || {}
  } catch (error) {
    console.error('Failed to fetch statistics:', error)
  }
}

function handleFilterChange() {
  pagination.page = 1
  fetchItems()
}

function handleSort(field) {
  if (sort.field === field) {
    sort.order = sort.order === 'asc' ? 'desc' : 'asc'
  } else {
    sort.field = field
    sort.order = 'asc'
  }
  fetchItems()
}

function handlePageChange(page) {
  pagination.page = page
  fetchItems()
}

function createItem() {
  router.push({ name: 'InventoryForm' })
}

function viewItem(id) {
  router.push({ name: 'InventoryView', params: { id } })
}

function editItem(id) {
  router.push({ name: 'InventoryForm', params: { id } })
}

function confirmDelete(id) {
  itemToDelete.value = id
  showDeleteModal.value = true
}

async function deleteItem() {
  try {
    await API.inventory.delete(itemToDelete.value)
    notificationStore.success(t('inventory.deleteSuccess'))
    showDeleteModal.value = false
    itemToDelete.value = null
    fetchItems()
    fetchStatistics()
  } catch (error) {
    console.error('Failed to delete inventory item:', error)
    notificationStore.error(t('inventory.deleteError'))
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

function isExpiringSoon(date) {
  if (!date) return false
  const expiryDate = new Date(date)
  const today = new Date()
  const daysUntilExpiry = Math.ceil((expiryDate - today) / (1000 * 60 * 60 * 24))
  return daysUntilExpiry > 0 && daysUntilExpiry <= 30
}

function isExpired(date) {
  if (!date) return false
  return new Date(date) < new Date()
}

onMounted(() => {
  fetchItems()
  fetchStatistics()
})
</script>

<style scoped>
.inventory-page {
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

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
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
  font-size: 2.5rem;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: bold;
  color: var(--text-primary);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.filters-card {
  margin-bottom: 1.5rem;
}

.filters {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
}

.form-control {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid var(--border-color);
  border-radius: 4px;
  background: var(--bg-primary);
  color: var(--text-primary);
  font-size: 0.875rem;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary-color);
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

.stock-level {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.stock-badge {
  font-weight: 600;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  background: #e8f5e9;
  color: #388e3c;
}

.stock-badge.low {
  background: #fff3e0;
  color: #f57c00;
}

.stock-badge.out {
  background: #ffebee;
  color: #c62828;
}

.stock-unit {
  font-size: 0.875rem;
  color: var(--text-secondary);
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

.actions {
  display: flex;
  gap: 0.5rem;
  flex-wrap: wrap;
}
</style>
