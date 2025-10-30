<template>
  <div class="data-table">
    <div class="table-wrapper">
      <table>
        <thead>
          <tr>
            <th
              v-for="column in columns"
              :key="column.key"
              :class="{ 'sortable': column.sortable }"
              @click="column.sortable && handleSort(column.key)"
            >
              <div class="th-content">
                <span>{{ column.label }}</span>
                <span v-if="column.sortable" class="sort-icon">
                  <span v-if="sortBy === column.key">
                    {{ sortOrder === 'asc' ? '↑' : '↓' }}
                  </span>
                  <span v-else class="sort-placeholder">⇅</span>
                </span>
              </div>
            </th>
            <th v-if="hasActions" class="actions-column">{{ t('common.actions') }}</th>
          </tr>
        </thead>

        <tbody>
          <tr v-if="loading">
            <td :colspan="columns.length + (hasActions ? 1 : 0)" class="loading-cell">
              <div class="loading-spinner"></div>
            </td>
          </tr>

          <tr v-else-if="!data || data.length === 0">
            <td :colspan="columns.length + (hasActions ? 1 : 0)" class="empty-cell">
              {{ t('common.noData') }}
            </td>
          </tr>

          <tr v-else v-for="(row, index) in data" :key="row.id || index">
            <td v-for="column in columns" :key="column.key">
              <slot :name="`cell-${column.key}`" :row="row" :value="getCellValue(row, column.key)">
                {{ getCellValue(row, column.key) }}
              </slot>
            </td>

            <td v-if="hasActions" class="actions-cell">
              <slot name="actions" :row="row"></slot>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <div v-if="pagination && total > 0" class="table-pagination">
      <div class="pagination-info">
        {{ t('common.page') }} {{ currentPage }} {{ t('common.of') }} {{ totalPages }}
        ({{ t('common.total') }}: {{ total }})
      </div>

      <div class="pagination-controls">
        <BaseButton
          size="small"
          variant="secondary"
          :disabled="currentPage === 1"
          @click="handlePageChange(currentPage - 1)"
        >
          ← {{ t('common.previous') }}
        </BaseButton>

        <BaseButton
          size="small"
          variant="secondary"
          :disabled="currentPage === totalPages"
          @click="handlePageChange(currentPage + 1)"
        >
          {{ t('common.next') }} →
        </BaseButton>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import BaseButton from './BaseButton.vue'

const { t } = useI18n()

const props = defineProps({
  columns: {
    type: Array,
    required: true,
    // [{ key: 'name', label: 'Name', sortable: true }]
  },
  data: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  },
  pagination: {
    type: Boolean,
    default: true
  },
  total: {
    type: Number,
    default: 0
  },
  currentPage: {
    type: Number,
    default: 1
  },
  perPage: {
    type: Number,
    default: 10
  },
  sortBy: {
    type: String,
    default: ''
  },
  sortOrder: {
    type: String,
    default: 'asc', // asc, desc
    validator: (value) => ['asc', 'desc'].includes(value)
  },
  hasActions: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['sort', 'page-change'])

const totalPages = computed(() => {
  return Math.ceil(props.total / props.perPage)
})

function getCellValue(row, key) {
  // Support nested keys like 'user.name'
  const keys = key.split('.')
  let value = row
  for (const k of keys) {
    value = value?.[k]
    if (value === undefined) return '-'
  }
  return value ?? '-'
}

function handleSort(key) {
  const newOrder = props.sortBy === key && props.sortOrder === 'asc' ? 'desc' : 'asc'
  emit('sort', { sortBy: key, sortOrder: newOrder })
}

function handlePageChange(page) {
  if (page >= 1 && page <= totalPages.value) {
    emit('page-change', page)
  }
}
</script>

<style scoped>
.data-table {
  width: 100%;
}

.table-wrapper {
  overflow-x: auto;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
}

table {
  width: 100%;
  border-collapse: collapse;
  background-color: var(--bg-secondary);
}

thead {
  background-color: var(--bg-primary);
  border-bottom: 2px solid var(--border-color);
}

th {
  padding: 0.75rem 1rem;
  text-align: left;
  font-weight: 600;
  color: var(--text-primary);
  font-size: 0.875rem;
  white-space: nowrap;
}

th.sortable {
  cursor: pointer;
  user-select: none;
}

th.sortable:hover {
  background-color: var(--bg-hover);
}

.th-content {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.sort-icon {
  color: var(--text-secondary);
  font-size: 0.875rem;
}

.sort-placeholder {
  opacity: 0.3;
}

td {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

tr:last-child td {
  border-bottom: none;
}

tbody tr:hover {
  background-color: var(--bg-hover);
}

.actions-column,
.actions-cell {
  text-align: right;
  width: 150px;
}

.actions-cell {
  display: flex;
  gap: 0.5rem;
  justify-content: flex-end;
}

.loading-cell,
.empty-cell {
  text-align: center;
  padding: 3rem 1rem;
  color: var(--text-secondary);
}

.loading-spinner {
  width: 2rem;
  height: 2rem;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
  margin: 0 auto;
}

.table-pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-top: none;
  border-radius: 0 0 0.5rem 0.5rem;
  background-color: var(--bg-secondary);
}

.pagination-info {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.pagination-controls {
  display: flex;
  gap: 0.5rem;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  th, td {
    padding: 0.5rem;
    font-size: 0.875rem;
  }

  .table-pagination {
    flex-direction: column;
    gap: 1rem;
  }
}
</style>
