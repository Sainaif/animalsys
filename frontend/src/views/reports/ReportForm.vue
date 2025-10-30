<template>
  <div class="report-form">
    <form @submit.prevent="handleSubmit">
      <!-- Report Type Selection -->
      <div class="form-section">
        <h3 class="section-title">{{ t('reports.reportType') }}</h3>

        <FormGroup :label="t('reports.selectReportType')" required>
          <select
            v-model="formData.type"
            required
            class="form-control"
            @change="handleTypeChange"
          >
            <option value="">{{ t('common.select') }}</option>
            <option value="financial">{{ t('reports.types.financial') }}</option>
            <option value="adoption">{{ t('reports.types.adoption') }}</option>
            <option value="volunteer">{{ t('reports.types.volunteer') }}</option>
            <option value="inventory">{{ t('reports.types.inventory') }}</option>
            <option value="veterinary">{{ t('reports.types.veterinary') }}</option>
            <option value="campaign">{{ t('reports.types.campaign') }}</option>
            <option value="donor">{{ t('reports.types.donor') }}</option>
            <option value="animal">{{ t('reports.types.animal') }}</option>
            <option value="statutory">{{ t('reports.types.statutory') }}</option>
            <option value="custom">{{ t('reports.types.custom') }}</option>
          </select>
        </FormGroup>

        <FormGroup :label="t('reports.name')" required>
          <input
            v-model="formData.name"
            type="text"
            class="form-control"
            required
            :placeholder="t('reports.namePlaceholder')"
          />
        </FormGroup>

        <FormGroup :label="t('reports.description')">
          <textarea
            v-model="formData.description"
            class="form-control"
            rows="3"
            :placeholder="t('reports.descriptionPlaceholder')"
          ></textarea>
        </FormGroup>
      </div>

      <!-- Date Range Parameters -->
      <div class="form-section">
        <h3 class="section-title">{{ t('reports.dateRange') }}</h3>

        <div class="date-range-row">
          <FormGroup :label="t('common.from')" required>
            <input
              v-model="formData.parameters.start_date"
              type="date"
              class="form-control"
              required
            />
          </FormGroup>

          <FormGroup :label="t('common.to')" required>
            <input
              v-model="formData.parameters.end_date"
              type="date"
              class="form-control"
              required
            />
          </FormGroup>
        </div>

        <div class="quick-date-buttons">
          <button type="button" @click="setThisMonth" class="btn btn-sm btn-outline">
            {{ t('reports.thisMonth') }}
          </button>
          <button type="button" @click="setLastMonth" class="btn btn-sm btn-outline">
            {{ t('reports.lastMonth') }}
          </button>
          <button type="button" @click="setThisQuarter" class="btn btn-sm btn-outline">
            {{ t('reports.thisQuarter') }}
          </button>
          <button type="button" @click="setThisYear" class="btn btn-sm btn-outline">
            {{ t('reports.thisYear') }}
          </button>
        </div>
      </div>

      <!-- Type-Specific Parameters -->
      <div v-if="formData.type" class="form-section">
        <h3 class="section-title">{{ t('reports.reportParameters') }}</h3>

        <!-- Financial Report Parameters -->
        <div v-if="formData.type === 'financial'">
          <FormGroup :label="t('reports.includeDetails')">
            <label class="checkbox-label">
              <input type="checkbox" v-model="formData.parameters.include_transactions" />
              {{ t('reports.includeTransactions') }}
            </label>
          </FormGroup>

          <FormGroup :label="t('reports.groupBy')">
            <select v-model="formData.parameters.group_by" class="form-control">
              <option value="category">{{ t('reports.byCategory') }}</option>
              <option value="month">{{ t('reports.byMonth') }}</option>
              <option value="week">{{ t('reports.byWeek') }}</option>
            </select>
          </FormGroup>
        </div>

        <!-- Adoption Report Parameters -->
        <div v-if="formData.type === 'adoption'">
          <FormGroup :label="t('reports.adoptionStatus')">
            <select v-model="formData.parameters.status" class="form-control">
              <option value="">{{ t('common.all') }}</option>
              <option value="submitted">{{ t('adoptions.submitted') }}</option>
              <option value="approved">{{ t('adoptions.approved') }}</option>
              <option value="rejected">{{ t('adoptions.rejected') }}</option>
              <option value="completed">{{ t('adoptions.completed') }}</option>
            </select>
          </FormGroup>

          <FormGroup :label="t('reports.includeStatistics')">
            <label class="checkbox-label">
              <input type="checkbox" v-model="formData.parameters.include_statistics" />
              {{ t('reports.includeStatisticsLabel') }}
            </label>
          </FormGroup>
        </div>

        <!-- Volunteer Report Parameters -->
        <div v-if="formData.type === 'volunteer'">
          <FormGroup :label="t('reports.volunteerStatus')">
            <select v-model="formData.parameters.status" class="form-control">
              <option value="">{{ t('common.all') }}</option>
              <option value="active">{{ t('volunteers.statusActive') }}</option>
              <option value="inactive">{{ t('volunteers.statusInactive') }}</option>
              <option value="on_leave">{{ t('volunteers.statusOnLeave') }}</option>
            </select>
          </FormGroup>

          <FormGroup :label="t('reports.includeHours')">
            <label class="checkbox-label">
              <input type="checkbox" v-model="formData.parameters.include_hours" />
              {{ t('reports.includeHoursLabel') }}
            </label>
          </FormGroup>
        </div>

        <!-- Inventory Report Parameters -->
        <div v-if="formData.type === 'inventory'">
          <FormGroup :label="t('reports.inventoryCategory')">
            <select v-model="formData.parameters.category" class="form-control">
              <option value="">{{ t('common.all') }}</option>
              <option value="food">{{ t('inventory.categoryFood') }}</option>
              <option value="medicine">{{ t('inventory.categoryMedicine') }}</option>
              <option value="supplies">{{ t('inventory.categorySupplies') }}</option>
              <option value="equipment">{{ t('inventory.categoryEquipment') }}</option>
            </select>
          </FormGroup>

          <FormGroup :label="t('reports.stockStatus')">
            <select v-model="formData.parameters.stock_status" class="form-control">
              <option value="">{{ t('common.all') }}</option>
              <option value="in_stock">{{ t('inventory.statusInStock') }}</option>
              <option value="low_stock">{{ t('inventory.statusLowStock') }}</option>
              <option value="out_of_stock">{{ t('inventory.statusOutOfStock') }}</option>
            </select>
          </FormGroup>
        </div>

        <!-- Animal Report Parameters -->
        <div v-if="formData.type === 'animal'">
          <FormGroup :label="t('reports.animalSpecies')">
            <select v-model="formData.parameters.species" class="form-control">
              <option value="">{{ t('common.all') }}</option>
              <option value="dog">{{ t('animals.dog') }}</option>
              <option value="cat">{{ t('animals.cat') }}</option>
              <option value="other">{{ t('animals.other') }}</option>
            </select>
          </FormGroup>

          <FormGroup :label="t('reports.animalStatus')">
            <select v-model="formData.parameters.status" class="form-control">
              <option value="">{{ t('common.all') }}</option>
              <option value="available">{{ t('animals.available') }}</option>
              <option value="adopted">{{ t('animals.adopted') }}</option>
              <option value="medical_care">{{ t('animals.medical_care') }}</option>
            </select>
          </FormGroup>
        </div>
      </div>

      <!-- Export Format -->
      <div class="form-section">
        <h3 class="section-title">{{ t('reports.exportFormat') }}</h3>

        <FormGroup :label="t('reports.selectFormat')">
          <div class="format-options">
            <label class="radio-label">
              <input type="radio" v-model="formData.format" value="pdf" />
              <span class="format-icon">📄</span>
              {{ t('reports.pdf') }}
            </label>
            <label class="radio-label">
              <input type="radio" v-model="formData.format" value="excel" />
              <span class="format-icon">📊</span>
              {{ t('reports.excel') }}
            </label>
            <label class="radio-label">
              <input type="radio" v-model="formData.format" value="csv" />
              <span class="format-icon">📋</span>
              {{ t('reports.csv') }}
            </label>
          </div>
        </FormGroup>
      </div>

      <!-- Schedule Options -->
      <div class="form-section">
        <h3 class="section-title">{{ t('reports.scheduleOptions') }}</h3>

        <FormGroup>
          <label class="checkbox-label">
            <input type="checkbox" v-model="scheduleReport" />
            {{ t('reports.scheduleForAutoGeneration') }}
          </label>
        </FormGroup>

        <div v-if="scheduleReport">
          <FormGroup :label="t('reports.frequency')">
            <select v-model="formData.schedule.frequency" class="form-control">
              <option value="daily">{{ t('reports.daily') }}</option>
              <option value="weekly">{{ t('reports.weekly') }}</option>
              <option value="monthly">{{ t('reports.monthly') }}</option>
              <option value="quarterly">{{ t('reports.quarterly') }}</option>
              <option value="yearly">{{ t('reports.yearly') }}</option>
            </select>
          </FormGroup>

          <FormGroup :label="t('reports.nextRunDate')">
            <input
              v-model="formData.schedule.next_run"
              type="datetime-local"
              class="form-control"
            />
          </FormGroup>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="form-actions">
        <button type="button" @click="$emit('cancel')" class="btn btn-secondary">
          {{ t('common.cancel') }}
        </button>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? t('reports.generating') : t('reports.generate') }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useI18n } from 'vue-i18n'
import FormGroup from '@/components/common/FormGroup.vue'

const { t } = useI18n()

const props = defineProps({
  report: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['submit', 'cancel'])

const loading = ref(false)
const scheduleReport = ref(false)

const formData = reactive({
  type: '',
  name: '',
  description: '',
  parameters: {
    start_date: '',
    end_date: '',
    include_transactions: false,
    include_statistics: false,
    include_hours: false,
    group_by: 'category',
    status: '',
    category: '',
    stock_status: '',
    species: ''
  },
  format: 'pdf',
  schedule: {
    frequency: 'monthly',
    next_run: ''
  }
})

const handleTypeChange = () => {
  // Auto-generate name based on type
  if (!formData.name || formData.name.includes(t('reports.types'))) {
    const date = new Date().toLocaleDateString('pl-PL')
    formData.name = `${t(`reports.types.${formData.type}`)} - ${date}`
  }
}

const setThisMonth = () => {
  const now = new Date()
  const firstDay = new Date(now.getFullYear(), now.getMonth(), 1)
  const lastDay = new Date(now.getFullYear(), now.getMonth() + 1, 0)

  formData.parameters.start_date = firstDay.toISOString().split('T')[0]
  formData.parameters.end_date = lastDay.toISOString().split('T')[0]
}

const setLastMonth = () => {
  const now = new Date()
  const firstDay = new Date(now.getFullYear(), now.getMonth() - 1, 1)
  const lastDay = new Date(now.getFullYear(), now.getMonth(), 0)

  formData.parameters.start_date = firstDay.toISOString().split('T')[0]
  formData.parameters.end_date = lastDay.toISOString().split('T')[0]
}

const setThisQuarter = () => {
  const now = new Date()
  const quarter = Math.floor(now.getMonth() / 3)
  const firstDay = new Date(now.getFullYear(), quarter * 3, 1)
  const lastDay = new Date(now.getFullYear(), (quarter + 1) * 3, 0)

  formData.parameters.start_date = firstDay.toISOString().split('T')[0]
  formData.parameters.end_date = lastDay.toISOString().split('T')[0]
}

const setThisYear = () => {
  const now = new Date()
  const firstDay = new Date(now.getFullYear(), 0, 1)
  const lastDay = new Date(now.getFullYear(), 11, 31)

  formData.parameters.start_date = firstDay.toISOString().split('T')[0]
  formData.parameters.end_date = lastDay.toISOString().split('T')[0]
}

const handleSubmit = () => {
  const submitData = {
    ...formData,
    schedule: scheduleReport.value ? formData.schedule : undefined
  }

  emit('submit', submitData)
}

// Initialize with current month
setThisMonth()
</script>

<style scoped>
.report-form {
  max-width: 900px;
}

.form-section {
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  padding: 1.5rem;
  margin-bottom: 1.5rem;
}

.section-title {
  font-size: 1.125rem;
  font-weight: 600;
  margin-bottom: 1.25rem;
  color: var(--text-primary);
  border-bottom: 2px solid var(--primary);
  padding-bottom: 0.5rem;
}

.form-control {
  width: 100%;
  padding: 0.625rem 0.875rem;
  border: 1px solid var(--border-color);
  border-radius: 6px;
  font-size: 0.9375rem;
  background: var(--input-bg);
  color: var(--text-primary);
  transition: border-color 0.2s;
}

.form-control:focus {
  outline: none;
  border-color: var(--primary);
  box-shadow: 0 0 0 3px var(--primary-alpha);
}

textarea.form-control {
  resize: vertical;
  font-family: inherit;
}

.date-range-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.quick-date-buttons {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
  margin-top: 1rem;
}

.checkbox-label,
.radio-label {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  padding: 0.5rem 0;
  font-size: 0.9375rem;
}

.checkbox-label input[type="checkbox"],
.radio-label input[type="radio"] {
  width: 18px;
  height: 18px;
  cursor: pointer;
}

.format-options {
  display: flex;
  gap: 1.5rem;
}

.format-icon {
  font-size: 1.25rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  justify-content: flex-end;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
  margin-top: 2rem;
}

.btn {
  padding: 0.625rem 1.5rem;
  border-radius: 6px;
  font-size: 0.9375rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
}

.btn-sm {
  padding: 0.375rem 0.875rem;
  font-size: 0.875rem;
}

.btn-primary {
  background: var(--primary);
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: var(--primary-dark);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
}

.btn-primary:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: var(--background);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-secondary:hover {
  background: var(--hover-bg);
  transform: translateY(-1px);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
}

.btn-outline {
  background: transparent;
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.btn-outline:hover {
  background: var(--hover-bg);
  border-color: var(--primary);
}
</style>
