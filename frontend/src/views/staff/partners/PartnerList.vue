<template>
  <div class="partner-list">
    <div class="page-header">
      <h1>{{ $t('partner.partners') }}</h1>
      <Button
        :label="$t('partner.addPartner')"
        icon="pi pi-plus"
        @click="router.push('/staff/partners/new')"
      />
    </div>

    <Card class="filters-card">
      <template #content>
        <div class="filters">
          <Dropdown
            v-model="filters.partner_type"
            :options="partnerTypeOptions"
            :placeholder="$t('partner.partnerType')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadPartners"
          />
          <Dropdown
            v-model="filters.status"
            :options="statusOptions"
            :placeholder="$t('common.status')"
            option-label="label"
            option-value="value"
            show-clear
            @change="loadPartners"
          />
        </div>
      </template>
    </Card>

    <Card v-if="!loading && partners.length > 0">
      <template #content>
        <DataTable
          :value="partners"
          paginator
          :rows="20"
        >
          <Column
            field="organization_name"
            :header="$t('partner.organizationName')"
          />
          <Column
            field="partner_type"
            :header="$t('partner.partnerType')"
          >
            <template #body="slotProps">
              {{ $t(`partner.${slotProps.data.partner_type}`) }}
            </template>
          </Column>
          <Column
            field="contact_person"
            :header="$t('partner.contactPerson')"
          />
          <Column
            field="email"
            header="Email"
          />
          <Column
            field="phone"
            :header="$t('finance.phone')"
          />
          <Column
            field="status"
            :header="$t('common.status')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ $t(`finance.${slotProps.data.status}`) }}
              </Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  @click="router.push(`/staff/partners/${slotProps.data.id}/edit`)"
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
      v-if="!loading && partners.length === 0"
      :message="$t('partner.noPartnersFound')"
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
import { partnerService } from '@/services/partnerService'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import ConfirmDialog from 'primevue/confirmdialog'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()
const toast = useToast()
const confirm = useConfirm()

const partners = ref([])
const loading = ref(true)
const filters = ref({
  partner_type: null,
  status: null
})

const partnerTypeOptions = [
  { label: t('partner.shelter'), value: 'shelter' },
  { label: t('partner.rescue'), value: 'rescue' },
  { label: t('partner.veterinary'), value: 'veterinary' },
  { label: t('partner.foster'), value: 'foster' },
  { label: t('communication.other'), value: 'other' }
]

const statusOptions = [
  { label: t('finance.active'), value: 'active' },
  { label: t('finance.inactive'), value: 'inactive' },
  { label: t('communication.pending'), value: 'pending' }
]

const loadPartners = async () => {
  try {
    loading.value = true
    const params = Object.fromEntries(
      Object.entries(filters.value).filter(([, value]) => value != null)
    )
    const response = await partnerService.getPartners(params)
    partners.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: t('common.error'), detail: 'Failed to load partners', life: 3000 })
  } finally {
    loading.value = false
  }
}

const getStatusVariant = (status) => ({
  active: 'success',
  inactive: 'neutral',
  pending: 'warning'
}[status] || 'neutral')

const confirmDelete = (partner) => {
  confirm.require({
    message: t('common.deleteConfirmation'),
    header: t('common.confirm'),
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await partnerService.deletePartner(partner.id)
        toast.add({ severity: 'success', summary: t('common.success'), detail: t('partner.partnerDeleted'), life: 3000 })
        loadPartners()
      } catch (error) {
        toast.add({ severity: 'error', summary: t('common.error'), detail: t('common.deleteError'), life: 3000 })
      }
    }
  })
}

onMounted(loadPartners)
</script>

<style scoped>
.partner-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
.filters-card { margin-bottom: 1.5rem; }
.filters { display: flex; gap: 1rem; flex-wrap: wrap; }
.filters > * { min-width: 200px; }
.action-buttons { display: flex; gap: 0.5rem; }
</style>
