<template>
  <div class="inventory-form-page">
    <div class="page-header">
      <h1 class="page-title">
        {{ isEdit ? t('inventory.editItem') : t('inventory.addItem') }}
      </h1>
    </div>

    <BaseCard>
      <LoadingSpinner v-if="loading" />
      <form v-else @submit.prevent="handleSubmit">
        <!-- Basic Information -->
        <div class="form-section">
          <h3 class="section-title">{{ t('inventory.basicInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('inventory.itemName')" :error="errors.name" required>
              <input
                v-model="form.name"
                type="text"
                class="form-control"
                :class="{ 'error': errors.name }"
                :placeholder="t('inventory.itemNamePlaceholder')"
              />
            </FormGroup>

            <FormGroup :label="t('inventory.category')" :error="errors.category" required>
              <select v-model="form.category" class="form-control" :class="{ 'error': errors.category }">
                <option value="">{{ t('inventory.selectCategory') }}</option>
                <option value="food">{{ t('inventory.categoryFood') }}</option>
                <option value="medicine">{{ t('inventory.categoryMedicine') }}</option>
                <option value="supplies">{{ t('inventory.categorySupplies') }}</option>
                <option value="equipment">{{ t('inventory.categoryEquipment') }}</option>
              </select>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('common.description')" :error="errors.description">
              <textarea
                v-model="form.description"
                class="form-control"
                :class="{ 'error': errors.description }"
                :placeholder="t('inventory.descriptionPlaceholder')"
                rows="3"
              ></textarea>
            </FormGroup>
          </div>
        </div>

        <!-- Stock Information -->
        <div class="form-section">
          <h3 class="section-title">{{ t('inventory.stockInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('inventory.stockLevel')" :error="errors.stock_level" required>
              <input
                v-model.number="form.stock_level"
                type="number"
                min="0"
                class="form-control"
                :class="{ 'error': errors.stock_level }"
                :placeholder="t('inventory.stockLevelPlaceholder')"
              />
            </FormGroup>

            <FormGroup :label="t('inventory.unit')" :error="errors.unit" required>
              <select v-model="form.unit" class="form-control" :class="{ 'error': errors.unit }">
                <option value="">{{ t('inventory.selectUnit') }}</option>
                <option value="kg">kg</option>
                <option value="g">g</option>
                <option value="l">l</option>
                <option value="ml">ml</option>
                <option value="szt">szt</option>
                <option value="op">op</option>
              </select>
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('inventory.minStock')" :error="errors.min_stock" required>
              <input
                v-model.number="form.min_stock"
                type="number"
                min="0"
                class="form-control"
                :class="{ 'error': errors.min_stock }"
                :placeholder="t('inventory.minStockPlaceholder')"
              />
            </FormGroup>

            <FormGroup :label="t('inventory.maxStock')" :error="errors.max_stock">
              <input
                v-model.number="form.max_stock"
                type="number"
                min="0"
                class="form-control"
                :class="{ 'error': errors.max_stock }"
                :placeholder="t('inventory.maxStockPlaceholder')"
              />
            </FormGroup>
          </div>
        </div>

        <!-- Pricing & Supplier -->
        <div class="form-section">
          <h3 class="section-title">{{ t('inventory.pricingInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('inventory.unitPrice')" :error="errors.unit_price">
              <input
                v-model.number="form.unit_price"
                type="number"
                step="0.01"
                min="0"
                class="form-control"
                :class="{ 'error': errors.unit_price }"
                :placeholder="t('inventory.unitPricePlaceholder')"
              />
            </FormGroup>

            <FormGroup :label="t('inventory.supplier')" :error="errors.supplier">
              <input
                v-model="form.supplier"
                type="text"
                class="form-control"
                :class="{ 'error': errors.supplier }"
                :placeholder="t('inventory.supplierPlaceholder')"
              />
            </FormGroup>
          </div>

          <div class="form-row">
            <FormGroup :label="t('inventory.expiryDate')" :error="errors.expiry_date">
              <input
                v-model="form.expiry_date"
                type="date"
                class="form-control"
                :class="{ 'error': errors.expiry_date }"
              />
            </FormGroup>

            <FormGroup :label="t('inventory.location')" :error="errors.location">
              <input
                v-model="form.location"
                type="text"
                class="form-control"
                :class="{ 'error': errors.location }"
                :placeholder="t('inventory.locationPlaceholder')"
              />
            </FormGroup>
          </div>
        </div>

        <!-- Notes -->
        <div class="form-section">
          <h3 class="section-title">{{ t('inventory.additionalInfo') }}</h3>

          <div class="form-row">
            <FormGroup :label="t('common.notes')" :error="errors.notes">
              <textarea
                v-model="form.notes"
                class="form-control"
                :class="{ 'error': errors.notes }"
                :placeholder="t('inventory.notesPlaceholder')"
                rows="4"
              ></textarea>
            </FormGroup>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="form-actions">
          <BaseButton type="button" variant="secondary" @click="goBack">
            {{ t('common.cancel') }}
          </BaseButton>
          <BaseButton type="submit" variant="primary" :disabled="submitting">
            {{ submitting ? t('common.saving') : t('common.save') }}
          </BaseButton>
        </div>
      </form>
    </BaseCard>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useNotificationStore } from '../../stores/notifications'
import { API } from '../../api'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import FormGroup from '../../components/base/FormGroup.vue'
import LoadingSpinner from '../../components/base/LoadingSpinner.vue'

const router = useRouter()
const route = useRoute()
const { t } = useI18n()
const notificationStore = useNotificationStore()

const isEdit = ref(false)
const loading = ref(false)
const submitting = ref(false)

const form = reactive({
  name: '',
  category: '',
  description: '',
  stock_level: 0,
  unit: '',
  min_stock: 0,
  max_stock: null,
  unit_price: null,
  supplier: '',
  expiry_date: '',
  location: '',
  notes: '',
})

const errors = reactive({})

async function fetchItem() {
  if (!route.params.id) return

  try {
    loading.value = true
    const response = await API.inventory.getById(route.params.id)
    const item = response.data

    Object.keys(form).forEach(key => {
      if (item[key] !== undefined) {
        form[key] = item[key]
      }
    })

    // Format date for input
    if (form.expiry_date) {
      form.expiry_date = new Date(form.expiry_date).toISOString().split('T')[0]
    }
  } catch (error) {
    console.error('Failed to fetch inventory item:', error)
    notificationStore.error(t('inventory.fetchError'))
    goBack()
  } finally {
    loading.value = false
  }
}

function validateForm() {
  Object.keys(errors).forEach(key => delete errors[key])
  let isValid = true

  if (!form.name || form.name.trim().length === 0) {
    errors.name = t('common.required')
    isValid = false
  }

  if (!form.category) {
    errors.category = t('common.required')
    isValid = false
  }

  if (form.stock_level === null || form.stock_level === undefined) {
    errors.stock_level = t('common.required')
    isValid = false
  }

  if (!form.unit) {
    errors.unit = t('common.required')
    isValid = false
  }

  if (form.min_stock === null || form.min_stock === undefined) {
    errors.min_stock = t('common.required')
    isValid = false
  }

  if (form.max_stock !== null && form.max_stock < form.min_stock) {
    errors.max_stock = t('inventory.maxStockError')
    isValid = false
  }

  return isValid
}

async function handleSubmit() {
  if (!validateForm()) {
    notificationStore.error(t('common.fixErrors'))
    return
  }

  try {
    submitting.value = true

    if (isEdit.value) {
      await API.inventory.update(route.params.id, form)
      notificationStore.success(t('inventory.updateSuccess'))
    } else {
      await API.inventory.create(form)
      notificationStore.success(t('inventory.createSuccess'))
    }

    goBack()
  } catch (error) {
    console.error('Failed to save inventory item:', error)
    notificationStore.error(
      isEdit.value ? t('inventory.updateError') : t('inventory.createError')
    )
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.push({ name: 'Inventory' })
}

onMounted(() => {
  if (route.params.id) {
    isEdit.value = true
    fetchItem()
  }
})
</script>

<style scoped>
.inventory-form-page {
  max-width: 900px;
  padding: 2rem;
}

.page-header {
  margin-bottom: 2rem;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin: 0;
}

.form-section {
  margin-bottom: 2rem;
}

.section-title {
  font-size: 1.25rem;
  font-weight: 600;
  margin-bottom: 1rem;
  padding-bottom: 0.5rem;
  border-bottom: 2px solid var(--border-color);
}

.form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
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

.form-control:disabled {
  background: var(--bg-secondary);
  cursor: not-allowed;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding-top: 1.5rem;
  border-top: 1px solid var(--border-color);
}
</style>
