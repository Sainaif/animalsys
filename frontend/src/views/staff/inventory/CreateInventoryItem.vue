<template>
  <div class="card">
    <h2>Create Inventory Item</h2>
    <form @submit.prevent="handleSubmit">
      <div class="p-fluid formgrid grid">
        <div class="field col-12 md:col-6">
          <label for="name">Name</label>
          <InputText
            id="name"
            v-model="item.name"
            required
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="category">Category</label>
          <Dropdown
            id="category"
            v-model="item.category"
            :options="categoryOptions"
            required
          />
        </div>
        <div class="field col-12">
          <label for="description">Description</label>
          <Textarea
            id="description"
            v-model="item.description"
            rows="3"
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="sku">SKU</label>
          <InputText
            id="sku"
            v-model="item.sku"
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="unit">Unit</label>
          <InputText
            id="unit"
            v-model="item.unit"
            required
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="quantity_in_stock">Quantity in Stock</label>
          <InputNumber
            id="quantity_in_stock"
            v-model="item.quantity_in_stock"
            required
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="minimum_quantity">Minimum Quantity</label>
          <InputNumber
            id="minimum_quantity"
            v-model="item.minimum_quantity"
            required
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="maximum_quantity">Maximum Quantity</label>
          <InputNumber
            id="maximum_quantity"
            v-model="item.maximum_quantity"
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="unit_cost">Unit Cost</label>
          <InputNumber
            id="unit_cost"
            v-model="item.unit_cost"
            mode="currency"
            currency="USD"
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="supplier">Supplier</label>
          <InputText
            id="supplier"
            v-model="item.supplier"
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="location">Location</label>
          <InputText
            id="location"
            v-model="item.location"
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="expiration_date">Expiration Date</label>
          <Calendar
            id="expiration_date"
            v-model="item.expiration_date"
          />
        </div>
        <div class="field col-12 md:col-6">
          <label for="status">Status</label>
          <Dropdown
            id="status"
            v-model="item.status"
            :options="statusOptions"
            required
          />
        </div>
      </div>
      <Button
        type="submit"
        label="Create Item"
        class="mt-4"
      />
    </form>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
import { inventoryService } from '@/services/inventoryService'

import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Textarea from 'primevue/textarea'
import InputNumber from 'primevue/inputnumber'
import Calendar from 'primevue/calendar'
import Button from 'primevue/button'

const router = useRouter()
const toast = useToast()
const { t } = useI18n()

const item = ref({
  name: '',
  category: 'other',
  description: '',
  sku: '',
  unit: '',
  quantity_in_stock: 0,
  minimum_quantity: 0,
  maximum_quantity: null,
  unit_cost: null,
  supplier: '',
  location: '',
  expiration_date: '',
  status: 'in_stock'
})

const categoryOptions = ref(['food', 'medicine', 'supplies', 'equipment', 'other'])
const statusOptions = ref(['in_stock', 'low_stock', 'out_of_stock', 'expired'])

const handleSubmit = async () => {
  try {
    await inventoryService.createInventoryItem(item.value)
    toast.add({ severity: 'success', summary: t('common.success'), detail: t('inventory.itemCreated'), life: 3000 })
    router.push({ name: 'inventory' })
  } catch (error) {
    console.error('Error creating inventory item:', error)
    toast.add({ severity: 'error', summary: t('common.error'), detail: t('inventory.itemCreateFailed'), life: 3000 })
  }
}
</script>
