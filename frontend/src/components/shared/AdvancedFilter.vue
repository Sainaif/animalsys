<template>
  <Card class="advanced-filter">
    <template #content>
      <div class="filter-grid">
        <!-- Search -->
        <div
          v-if="showSearch"
          class="filter-field"
        >
          <label>{{ $t('common.search') }}</label>
          <InputText
            v-model="localFilters.search"
            :placeholder="searchPlaceholder"
            @input="emitFilters"
          />
        </div>

        <!-- Date Range -->
        <div
          v-if="showDateRange"
          class="filter-field"
        >
          <label>{{ $t('reports.from') }}</label>
          <Calendar
            v-model="localFilters.dateFrom"
            date-format="yy-mm-dd"
            show-icon
            @date-select="emitFilters"
          />
        </div>

        <div
          v-if="showDateRange"
          class="filter-field"
        >
          <label>{{ $t('reports.to') }}</label>
          <Calendar
            v-model="localFilters.dateTo"
            date-format="yy-mm-dd"
            show-icon
            @date-select="emitFilters"
          />
        </div>

        <!-- Custom filter slots -->
        <slot />

        <!-- Action buttons -->
        <div class="filter-actions">
          <Button
            :label="$t('common.filter')"
            icon="pi pi-filter"
            @click="emitFilters"
          />
          <Button
            v-if="showExport"
            :label="$t('common.export')"
            icon="pi pi-download"
            class="p-button-secondary"
            @click="$emit('export')"
          />
          <Button
            label="Clear"
            icon="pi pi-times"
            class="p-button-text"
            @click="clearFilters"
          />
        </div>
      </div>
    </template>
  </Card>
</template>

<script setup>
import { ref, watch } from 'vue'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import Calendar from 'primevue/calendar'
import Button from 'primevue/button'

const props = defineProps({
  modelValue: {
    type: Object,
    default: () => ({})
  },
  showSearch: {
    type: Boolean,
    default: true
  },
  showDateRange: {
    type: Boolean,
    default: false
  },
  showExport: {
    type: Boolean,
    default: true
  },
  searchPlaceholder: {
    type: String,
    default: 'Search...'
  }
})

const emit = defineEmits(['update:modelValue', 'filter', 'export'])

const localFilters = ref({
  search: props.modelValue.search || '',
  dateFrom: props.modelValue.dateFrom || null,
  dateTo: props.modelValue.dateTo || null,
  ...props.modelValue
})

watch(() => props.modelValue, (newVal) => {
  localFilters.value = { ...localFilters.value, ...newVal }
}, { deep: true })

const emitFilters = () => {
  const filters = { ...localFilters.value }
  // Convert dates to ISO strings if present
  if (filters.dateFrom) {
    filters.dateFrom = filters.dateFrom.toISOString().split('T')[0]
  }
  if (filters.dateTo) {
    filters.dateTo = filters.dateTo.toISOString().split('T')[0]
  }
  emit('update:modelValue', filters)
  emit('filter', filters)
}

const clearFilters = () => {
  localFilters.value = {
    search: '',
    dateFrom: null,
    dateTo: null
  }
  emitFilters()
}
</script>

<style scoped>
.advanced-filter { margin-bottom: 1.5rem; }
.filter-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(200px, 1fr)); gap: 1rem; align-items: end; }
.filter-field { display: flex; flex-direction: column; gap: 0.5rem; }
.filter-field label { font-weight: 600; color: var(--text-color); font-size: 0.875rem; }
.filter-actions { display: flex; gap: 0.5rem; align-items: center; }
</style>
