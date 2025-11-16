<template>
  <div class="inventory-list">
    <div class="page-header">
      <h1>{{ $t('inventory.title') }}</h1>
      <Button
        :label="$t('inventory.addItem')"
        icon="pi pi-plus"
        @click="router.push('/staff/inventory/new')"
      />
    </div>

    <Card v-if="!loading && items.length > 0">
      <template #content>
        <DataTable
          :value="items"
          paginator
          :rows="20"
        >
          <Column
            field="name"
            :header="$t('inventory.itemName')"
          />
          <Column
            field="category"
            :header="$t('inventory.category')"
          >
            <template #body="slotProps">
              {{ $t(`inventory.${slotProps.data.category}`) }}
            </template>
          </Column>
          <Column
            field="quantity_in_stock"
            :header="$t('inventory.quantityInStock')"
          />
          <Column
            field="minimum_quantity"
            :header="$t('inventory.minimumQuantity')"
          />
          <Column
            field="unit"
            :header="$t('inventory.unit')"
          />
          <Column
            field="status"
            :header="$t('common.status')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ $t(`inventory.${slotProps.data.status}`) }}
              </Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  @click="router.push(`/staff/inventory/${slotProps.data.id}/edit`)"
                />
                <Button
                  icon="pi pi-trash"
                  class="p-button-rounded p-button-text p-button-danger"
                  @click="confirmDelete(slotProps.data)"
                />
              </div>
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading" />
    <EmptyState
      v-if="!loading && items.length === 0"
      :message="$t('inventory.noItemsFound')"
    />
    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import { inventoryService } from '@/services/inventoryService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const items = ref([])
const loading = ref(true)

const loadItems = async () => {
  try {
    loading.value = true
    const response = await inventoryService.getInventoryItems()
    items.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load inventory items', life: 3000 })
  } finally {
    loading.value = false
  }
}

const getStatusVariant = (status) => ({
  in_stock: 'success',
  low_stock: 'warning',
  out_of_stock: 'danger',
  expired: 'danger'
}[status] || 'neutral')

const confirmDelete = (item) => {
  confirm.require({
    message: 'Are you sure you want to delete this item?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await inventoryService.deleteInventoryItem(item.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('inventory.itemDeleted'), life: 3000 })
        loadItems()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete item', life: 3000 })
      }
    }
  })
}

onMounted(loadItems)
</script>

<style scoped>
.inventory-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.action-buttons { display: flex; gap: 0.5rem; }
</style>
