<template>
  <div class="email-template-list">
    <div class="page-header">
      <h1>{{ $t('communication.emailTemplates') }}</h1>
      <Button
        :label="$t('communication.addEmailTemplate')"
        icon="pi pi-plus"
        @click="router.push('/staff/communication/templates/new')"
      />
    </div>

    <Card v-if="!loading && templates.length > 0">
      <template #content>
        <DataTable
          :value="templates"
          paginator
          :rows="20"
        >
          <Column
            field="name"
            :header="$t('communication.templateName')"
          />
          <Column
            field="subject"
            :header="$t('communication.subject')"
          />
          <Column
            field="template_type"
            :header="$t('communication.templateType')"
          >
            <template #body="slotProps">
              {{ $t(`communication.${slotProps.data.template_type}`) }}
            </template>
          </Column>
          <Column
            field="is_active"
            :header="$t('communication.isActive')"
          >
            <template #body="slotProps">
              <Badge :variant="slotProps.data.is_active ? 'success' : 'neutral'">
                {{ slotProps.data.is_active ? $t('common.yes') : $t('common.no') }}
              </Badge>
            </template>
          </Column>
          <Column :header="$t('common.actions')">
            <template #body="slotProps">
              <div class="action-buttons">
                <Button
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  @click="router.push(`/staff/communication/templates/${slotProps.data.id}/edit`)"
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
      v-if="!loading && templates.length === 0"
      :message="$t('communication.noEmailTemplatesFound')"
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
import { communicationService } from '@/services/communicationService'
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

const templates = ref([])
const loading = ref(true)

const loadTemplates = async () => {
  try {
    loading.value = true
    const response = await communicationService.getEmailTemplates()
    templates.value = response.data
  } catch (error) {
    toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to load email templates', life: 3000 })
  } finally {
    loading.value = false
  }
}

const confirmDelete = (template) => {
  confirm.require({
    message: 'Are you sure you want to delete this template?',
    header: 'Confirmation',
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await communicationService.deleteEmailTemplate(template.id)
        toast.add({ severity: 'success', summary: 'Success', detail: t('communication.emailTemplateDeleted'), life: 3000 })
        loadTemplates()
      } catch (error) {
        toast.add({ severity: 'error', summary: 'Error', detail: 'Failed to delete template', life: 3000 })
      }
    }
  })
}

onMounted(loadTemplates)
</script>

<style scoped>
.email-template-list { max-width: 1400px; margin: 0 auto; }
.page-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.action-buttons { display: flex; gap: 0.5rem; }
</style>
