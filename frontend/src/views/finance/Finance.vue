<template>
  <div class="finance-page">
    <div class="page-header">
      <h1 class="page-title">{{ t('nav.finance') }}</h1>
      <BaseButton variant="primary" @click="showTransactionModal">
        ➕ {{ t('finance.addTransaction') }}
      </BaseButton>
    </div>

    <!-- Financial Dashboard -->
    <div class="dashboard-stats">
      <div class="stat-card stat-card--income">
        <div class="stat-icon">💰</div>
        <div class="stat-content">
          <div class="stat-label">{{ t('finances.income') }}</div>
          <div class="stat-value">{{ formatCurrency(stats.totalIncome) }}</div>
          <div class="stat-change">{{ t('finances.thisMonth') }}</div>
        </div>
      </div>
      <div class="stat-card stat-card--expense">
        <div class="stat-icon">💸</div>
        <div class="stat-content">
          <div class="stat-label">{{ t('finances.expense') }}</div>
          <div class="stat-value">{{ formatCurrency(stats.totalExpense) }}</div>
          <div class="stat-change">{{ t('finances.thisMonth') }}</div>
        </div>
      </div>
      <div class="stat-card stat-card--balance">
        <div class="stat-icon">📊</div>
        <div class="stat-content">
          <div class="stat-label">{{ t('finances.balance') }}</div>
          <div class="stat-value" :class="{ 'text-negative': stats.balance < 0 }">
            {{ formatCurrency(stats.balance) }}
          </div>
          <div class="stat-change">{{ t('finances.net') }}</div>
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
            :placeholder="t('finances.searchPlaceholder')"
            @input="handleFilterChange"
          />
        </FormGroup>

        <FormGroup :label="t('finances.type')">
          <select v-model="filters.type" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="income">{{ t('finances.income') }}</option>
            <option value="expense">{{ t('finances.expense') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('finances.category')">
          <select v-model="filters.category" @change="handleFilterChange">
            <option value="">{{ t('common.all') }}</option>
            <option value="donations">{{ t('finances.donations') }}</option>
            <option value="grants">{{ t('finances.grants') }}</option>
            <option value="veterinary">{{ t('finances.veterinary') }}</option>
            <option value="food">{{ t('finances.food') }}</option>
            <option value="utilities">{{ t('finances.utilities') }}</option>
            <option value="salaries">{{ t('finances.salaries') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('finances.dateRange')">
          <div class="date-range">
            <input
              v-model="filters.startDate"
              type="date"
              @change="handleFilterChange"
            />
            <span>-</span>
            <input
              v-model="filters.endDate"
              type="date"
              @change="handleFilterChange"
            />
          </div>
        </FormGroup>
      </div>
    </BaseCard>

    <!-- Transactions Table -->
    <BaseCard padding="none">
      <DataTable
        :columns="columns"
        :data="transactions"
        :loading="loading"
        :total="total"
        :current-page="pagination.page"
        :per-page="pagination.limit"
        :sort-by="sort.sortBy"
        :sort-order="sort.sortOrder"
        has-actions
        @sort="handleSort"
        @page-change="handlePageChange"
      >
        <template #cell-date="{ value }">
          {{ formatDate(value) }}
        </template>

        <template #cell-type="{ value }">
          <span :class="['type-badge', `type-badge--${value}`]">
            {{ t(`finances.${value}`) }}
          </span>
        </template>

        <template #cell-category="{ value }">
          {{ t(`finances.${value}`) }}
        </template>

        <template #cell-amount="{ value, row }">
          <span :class="{ 'amount-income': row.type === 'income', 'amount-expense': row.type === 'expense' }">
            {{ row.type === 'income' ? '+' : '-' }}{{ formatCurrency(value) }}
          </span>
        </template>

        <template #actions="{ row }">
          <BaseButton
            size="small"
            variant="ghost"
            @click="viewTransaction(row)"
          >
            {{ t('common.view') }}
          </BaseButton>
          <BaseButton
            v-if="canDelete"
            size="small"
            variant="danger"
            @click="confirmDelete(row)"
          >
            {{ t('common.delete') }}
          </BaseButton>
        </template>
      </DataTable>
    </BaseCard>

    <!-- Delete Confirmation Modal -->
    <BaseModal
      v-model="deleteModal.show"
      :title="t('finances.deleteTransaction')"
      size="small"
    >
      <p>{{ t('finances.deleteConfirm') }}</p>
      <p v-if="deleteModal.transaction" class="transaction-details">
        <strong>{{ formatCurrency(deleteModal.transaction.amount) }}</strong>
        - {{ deleteModal.transaction.description }}
      </p>

      <template #footer>
        <BaseButton variant="outline" @click="deleteModal.show = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton variant="danger" :loading="deleteModal.loading" @click="deleteTransaction">
          {{ t('common.delete') }}
        </BaseButton>
      </template>
    </BaseModal>

    <!-- Transaction Modal -->
    <BaseModal
      v-model="transactionModal.show"
      :title="transactionModal.isEdit ? t('finances.editTransaction') : t('finance.addTransaction')"
      size="medium"
    >
      <form @submit.prevent="saveTransaction" class="transaction-form">
        <FormGroup :label="t('finances.type')" required>
          <select v-model="transactionModal.form.type" required>
            <option value="">{{ t('common.select') }}</option>
            <option value="income">{{ t('finances.income') }}</option>
            <option value="expense">{{ t('finances.expense') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('finances.category')" required>
          <select v-model="transactionModal.form.category" required>
            <option value="">{{ t('common.select') }}</option>
            <option value="donations">{{ t('finances.donations') }}</option>
            <option value="grants">{{ t('finances.grants') }}</option>
            <option value="veterinary">{{ t('finances.veterinary') }}</option>
            <option value="food">{{ t('finances.food') }}</option>
            <option value="utilities">{{ t('finances.utilities') }}</option>
            <option value="salaries">{{ t('finances.salaries') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('finances.amount')" required>
          <input
            v-model.number="transactionModal.form.amount"
            type="number"
            step="0.01"
            min="0"
            :placeholder="t('finances.amount')"
            required
          />
        </FormGroup>

        <FormGroup :label="t('common.description')" required>
          <textarea
            v-model="transactionModal.form.description"
            :placeholder="t('common.description')"
            rows="3"
            required
          ></textarea>
        </FormGroup>

        <FormGroup :label="t('common.date')" required>
          <input
            v-model="transactionModal.form.date"
            type="date"
            required
          />
        </FormGroup>
      </form>

      <template #footer>
        <BaseButton variant="outline" @click="transactionModal.show = false">
          {{ t('common.cancel') }}
        </BaseButton>
        <BaseButton
          variant="primary"
          :loading="transactionModal.loading"
          @click="saveTransaction"
        >
          {{ transactionModal.isEdit ? t('common.update') : t('common.create') }}
        </BaseButton>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { API } from '../../api'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notification'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import BaseModal from '../../components/base/BaseModal.vue'
import DataTable from '../../components/base/DataTable.vue'
import FormGroup from '../../components/base/FormGroup.vue'

const { t } = useI18n()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const transactions = ref([])
const loading = ref(false)
const total = ref(0)

const stats = reactive({
  totalIncome: 0,
  totalExpense: 0,
  balance: 0
})

const canDelete = computed(() => authStore.hasRole('admin'))

const filters = reactive({
  search: '',
  type: '',
  category: '',
  startDate: '',
  endDate: ''
})

const pagination = reactive({
  page: 1,
  limit: 10
})

const sort = reactive({
  sortBy: 'date',
  sortOrder: 'desc'
})

const deleteModal = reactive({
  show: false,
  transaction: null,
  loading: false
})

const transactionModal = reactive({
  show: false,
  isEdit: false,
  loading: false,
  form: {
    type: '',
    category: '',
    amount: 0,
    description: '',
    date: new Date().toISOString().split('T')[0]
  }
})

const columns = [
  { key: 'date', label: t('common.date'), sortable: true },
  { key: 'type', label: t('finances.type'), sortable: true },
  { key: 'category', label: t('finances.category'), sortable: false },
  { key: 'description', label: t('common.description'), sortable: false },
  { key: 'amount', label: t('finances.amount'), sortable: true }
]

onMounted(() => {
  fetchTransactions()
  fetchStats()
})

async function fetchTransactions() {
  loading.value = true

  try {
    const params = {
      limit: pagination.limit,
      offset: (pagination.page - 1) * pagination.limit,
      search: filters.search || undefined,
      type: filters.type || undefined,
      category: filters.category || undefined,
      start_date: filters.startDate || undefined,
      end_date: filters.endDate || undefined
    }

    const response = await API.finance.list(params)
    transactions.value = response.data.data || []
    total.value = response.data.total || 0
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    loading.value = false
  }
}

async function fetchStats() {
  try {
    const response = await API.finance.getDashboardStats()
    Object.assign(stats, response.data)
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

function handleFilterChange() {
  pagination.page = 1
  fetchTransactions()
}

function handleSort({ sortBy, sortOrder }) {
  sort.sortBy = sortBy
  sort.sortOrder = sortOrder
  fetchTransactions()
}

function handlePageChange(page) {
  pagination.page = page
  fetchTransactions()
}

function formatDate(dateString) {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleDateString()
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('pl-PL', {
    style: 'currency',
    currency: 'PLN'
  }).format(amount || 0)
}

function showTransactionModal() {
  transactionModal.isEdit = false
  transactionModal.form = {
    type: '',
    category: '',
    amount: 0,
    description: '',
    date: new Date().toISOString().split('T')[0]
  }
  transactionModal.show = true
}

function viewTransaction(transaction) {
  transactionModal.isEdit = true
  transactionModal.form = {
    type: transaction.type,
    category: transaction.category,
    amount: transaction.amount,
    description: transaction.description,
    date: transaction.date ? transaction.date.split('T')[0] : new Date().toISOString().split('T')[0]
  }
  transactionModal.transactionId = transaction.id
  transactionModal.show = true
}

async function saveTransaction() {
  if (!transactionModal.form.type || !transactionModal.form.category || !transactionModal.form.amount || !transactionModal.form.description) {
    notificationStore.error(t('common.error'), t('common.required'))
    return
  }

  transactionModal.loading = true
  try {
    if (transactionModal.isEdit) {
      await API.finance.update(transactionModal.transactionId, transactionModal.form)
      notificationStore.success(t('finances.updateSuccess'))
    } else {
      await API.finance.create(transactionModal.form)
      notificationStore.success(t('finances.createSuccess'))
    }
    transactionModal.show = false
    fetchTransactions()
    fetchStats()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    transactionModal.loading = false
  }
}

function confirmDelete(transaction) {
  deleteModal.transaction = transaction
  deleteModal.show = true
}

async function deleteTransaction() {
  if (!deleteModal.transaction) return

  deleteModal.loading = true

  try {
    await API.finance.delete(deleteModal.transaction.id)
    notificationStore.success(t('finances.deleteSuccess'))
    deleteModal.show = false
    deleteModal.transaction = null
    fetchTransactions()
    fetchStats()
  } catch (error) {
    notificationStore.error(t('common.error'), error.message)
  } finally {
    deleteModal.loading = false
  }
}
</script>

<style scoped>
.finance-page {
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

.dashboard-stats {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1.5rem;
  border-radius: 0.75rem;
  background-color: var(--bg-secondary);
}

.stat-card--income {
  border-left: 4px solid #10b981;
}

.stat-card--expense {
  border-left: 4px solid #ef4444;
}

.stat-card--balance {
  border-left: 4px solid #3b82f6;
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 1.75rem;
  font-weight: bold;
  color: var(--text-primary);
  margin-bottom: 0.25rem;
}

.stat-value.text-negative {
  color: #ef4444;
}

.stat-change {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.filters-card {
  margin-bottom: 1.5rem;
}

.filters {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.date-range {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.type-badge {
  display: inline-block;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 500;
}

.type-badge--income {
  background-color: #d1fae5;
  color: #065f46;
}

.type-badge--expense {
  background-color: #fee2e2;
  color: #991b1b;
}

.amount-income {
  color: #10b981;
  font-weight: 600;
}

.amount-expense {
  color: #ef4444;
  font-weight: 600;
}

.transaction-details {
  margin: 1rem 0;
  text-align: center;
}

.transaction-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

input,
select,
textarea {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
  font-family: inherit;
}

input:focus,
select:focus,
textarea:focus {
  outline: none;
  border-color: var(--primary-color);
}

textarea {
  resize: vertical;
}

@media (max-width: 768px) {
  .finance-page {
    padding: 1rem;
  }

  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .dashboard-stats {
    grid-template-columns: 1fr;
  }

  .filters {
    grid-template-columns: 1fr;
  }

  .date-range {
    flex-direction: column;
  }
}
</style>
