<template>
  <div class="contacts-page page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">
          {{ $t('contacts.title') }}
        </h1>
        <p class="page-subtitle">
          {{ $t('contacts.subtitle') }}
        </p>
      </div>
      <div class="header-actions">
        <Button
          :label="$t('contacts.actions.export')"
          icon="pi pi-download"
          class="p-button-text"
          @click="exportContacts"
        />
        <Button
          :label="$t('contacts.actions.add')"
          icon="pi pi-plus"
          @click="openCreateDialog"
        />
      </div>
    </div>

    <Card class="stats-card">
      <template #content>
        <div class="stats-grid">
          <div
            v-for="stat in contactStats"
            :key="stat.label"
            class="stat-card"
          >
            <span class="stat-label">{{ stat.label }}</span>
            <strong class="stat-value">{{ stat.value }}</strong>
            <small v-if="stat.description">{{ stat.description }}</small>
          </div>
        </div>
      </template>
    </Card>

    <AdvancedFilter
      v-model="filters"
      :show-date-range="false"
      :show-export="false"
      :search-placeholder="$t('contacts.searchPlaceholder')"
      @filter="applyFilters"
    >
      <div class="filter-field">
        <label>{{ $t('contacts.filters.type') }}</label>
        <Dropdown
          v-model="filters.type"
          :options="typeOptions"
          option-label="label"
          option-value="value"
          :placeholder="$t('contacts.filters.allTypes')"
          show-clear
          @change="applyFilters"
        />
      </div>
      <div class="filter-field">
        <label>{{ $t('contacts.filters.status') }}</label>
        <Dropdown
          v-model="filters.status"
          :options="statusOptions"
          option-label="label"
          option-value="value"
          :placeholder="$t('contacts.filters.allStatuses')"
          show-clear
          @change="applyFilters"
        />
      </div>
      <div class="filter-field">
        <label>{{ $t('contacts.filters.owner') }}</label>
        <Dropdown
          v-model="filters.owner_id"
          :options="ownerOptions"
          option-label="label"
          option-value="value"
          :placeholder="$t('contacts.filters.allOwners')"
          show-clear
          @change="applyFilters"
        />
      </div>
    </AdvancedFilter>

    <Card class="table-card">
      <template #content>
        <DataTable
          :value="contacts"
          :rows="pagination.limit"
          :total-records="pagination.total"
          :lazy="true"
          paginator
          data-key="id"
          :loading="loading"
          row-hover
          @page="onPage"
          @row-click="openDetails"
        >
          <Column
            field="name"
            :header="$t('contacts.table.name')"
          >
            <template #body="slotProps">
              <div class="contact-name">
                <div
                  class="avatar"
                  :data-letter="getInitials(slotProps.data)"
                >
                  <i class="pi pi-user" />
                </div>
                <div>
                  <p class="contact-title">
                    {{ slotProps.data.first_name }} {{ slotProps.data.last_name }}
                    <Badge :variant="getTypeVariant(slotProps.data.type)">
                      {{ formatType(slotProps.data.type) }}
                    </Badge>
                  </p>
                  <span class="contact-email">{{ slotProps.data.email }}</span>
                </div>
              </div>
            </template>
          </Column>
          <Column
            field="organization"
            :header="$t('contacts.table.organisation')"
          >
            <template #body="slotProps">
              {{ slotProps.data.organization || $t('common.notAvailable') }}
            </template>
          </Column>
          <Column
            field="tags"
            :header="$t('contacts.table.tags')"
          >
            <template #body="slotProps">
              <Tag
                v-for="tag in slotProps.data.tags"
                :key="tag"
                :value="tag"
                class="mr-2"
                severity="info"
              />
              <span v-if="!slotProps.data.tags?.length">{{ $t('common.notAvailable') }}</span>
            </template>
          </Column>
          <Column
            field="owner"
            :header="$t('contacts.table.owner')"
          >
            <template #body="slotProps">
              {{ slotProps.data.owner_name || $t('common.unassigned') }}
            </template>
          </Column>
          <Column
            field="status"
            :header="$t('contacts.form.status')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ formatStatus(slotProps.data.status) }}
              </Badge>
            </template>
          </Column>
          <Column
            field="last_contacted_at"
            :header="$t('contacts.table.lastContact')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.last_contacted_at) }}
            </template>
          </Column>
          <Column
            field="next_follow_up_at"
            :header="$t('contacts.table.nextFollowUp')"
          >
            <template #body="slotProps">
              <span :class="{ overdue: isOverdue(slotProps.data.next_follow_up_at) }">
                {{ formatDate(slotProps.data.next_follow_up_at) }}
              </span>
            </template>
          </Column>
          <Column
            :header="$t('common.actions')"
            :style="{ width: '140px' }"
          >
            <template #body="slotProps">
              <div
                class="action-buttons"
                @click.stop
              >
                <Button
                  v-tooltip.top="$t('common.edit')"
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  @click="openEditDialog(slotProps.data)"
                />
                <Button
                  v-tooltip.top="$t('contacts.actions.copyEmail')"
                  icon="pi pi-copy"
                  class="p-button-rounded p-button-text"
                  @click="copyEmail(slotProps.data.email)"
                />
                <Button
                  v-tooltip.top="$t('common.delete')"
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

    <LoadingSpinner v-if="loading && contacts.length === 0" />
    <EmptyState
      v-if="!loading && contacts.length === 0"
      :message="$t('contacts.messages.empty')"
    />

    <!-- Create/Edit dialog -->
    <Dialog
      v-model:visible="contactDialog.visible"
      modal
      :header="contactDialog.mode === 'create' ? $t('contacts.dialog.createTitle') : $t('contacts.dialog.editTitle')"
      :style="{ width: '640px' }"
    >
      <form
        class="form-grid"
        @submit.prevent="saveContact"
      >
        <div class="form-group">
          <label for="firstName">{{ $t('contacts.form.firstName') }}</label>
          <InputText
            id="firstName"
            v-model="contactDialog.form.first_name"
            required
          />
        </div>
        <div class="form-group">
          <label for="lastName">{{ $t('contacts.form.lastName') }}</label>
          <InputText
            id="lastName"
            v-model="contactDialog.form.last_name"
            required
          />
        </div>
        <div class="form-group">
          <label for="email">{{ $t('contacts.form.email') }}</label>
          <InputText
            id="email"
            v-model="contactDialog.form.email"
            type="email"
          />
        </div>
        <div class="form-group">
          <label for="phone">{{ $t('contacts.form.phone') }}</label>
          <InputText
            id="phone"
            v-model="contactDialog.form.phone"
          />
        </div>
        <div class="form-group">
          <label for="organization">{{ $t('contacts.form.organisation') }}</label>
          <InputText
            id="organization"
            v-model="contactDialog.form.organization"
          />
        </div>
        <div class="form-group">
          <label for="type">{{ $t('contacts.form.type') }}</label>
          <Dropdown
            id="type"
            v-model="contactDialog.form.type"
            :options="typeOptions"
            option-label="label"
            option-value="value"
            required
          />
        </div>
        <div class="form-group">
          <label for="status">{{ $t('contacts.form.status') }}</label>
          <Dropdown
            id="status"
            v-model="contactDialog.form.status"
            :options="statusOptions"
            option-label="label"
            option-value="value"
            required
          />
        </div>
        <div class="form-group">
          <label for="owner">{{ $t('contacts.form.owner') }}</label>
          <Dropdown
            id="owner"
            v-model="contactDialog.form.owner_id"
            :options="ownerOptions"
            option-label="label"
            option-value="value"
            :placeholder="$t('common.unassigned')"
            show-clear
          />
        </div>
        <div class="form-group full-width">
          <label for="tags">{{ $t('contacts.form.tags') }}</label>
          <Chips
            id="tags"
            v-model="contactDialog.form.tags"
            separator=","
          />
        </div>
        <div class="form-group">
          <label for="nextFollowUp">{{ $t('contacts.form.nextFollowUp') }}</label>
          <Calendar
            id="nextFollowUp"
            v-model="contactDialog.form.next_follow_up_at"
            date-format="yy-mm-dd"
            show-icon
          />
        </div>
        <div class="form-group full-width">
          <label for="notes">{{ $t('contacts.form.notes') }}</label>
          <Textarea
            id="notes"
            v-model="contactDialog.form.notes"
            rows="3"
          />
        </div>
        <div class="dialog-actions">
          <Button
            :label="$t('common.cancel')"
            text
            type="button"
            @click="contactDialog.visible = false"
          />
          <Button
            :label="$t('contacts.actions.save')"
            icon="pi pi-check"
            type="submit"
            :loading="contactDialog.saving"
          />
        </div>
      </form>
    </Dialog>

    <!-- Detail dialog -->
    <Dialog
      v-model:visible="detailDialog.visible"
      modal
      :header="detailDialog.contact ? getContactName(detailDialog.contact) : $t('contacts.detail.title')"
      :style="{ width: '720px' }"
    >
      <div
        v-if="detailDialog.contact"
        class="detail-layout"
      >
        <section>
          <h3>{{ $t('contacts.detail.info') }}</h3>
          <ul class="detail-list">
            <li>
              <span>{{ $t('contacts.detail.email') }}</span>
              <strong>{{ detailDialog.contact.email || $t('common.notAvailable') }}</strong>
            </li>
            <li>
              <span>{{ $t('contacts.detail.phone') }}</span>
              <strong>{{ detailDialog.contact.phone || $t('common.notAvailable') }}</strong>
            </li>
            <li>
              <span>{{ $t('contacts.detail.organisation') }}</span>
              <strong>{{ detailDialog.contact.organization || $t('common.notAvailable') }}</strong>
            </li>
            <li>
              <span>{{ $t('contacts.detail.owner') }}</span>
              <strong>{{ detailDialog.contact.owner_name || $t('common.unassigned') }}</strong>
            </li>
            <li>
              <span>{{ $t('contacts.detail.status') }}</span>
              <strong>{{ formatStatus(detailDialog.contact.status) }}</strong>
            </li>
            <li>
              <span>{{ $t('contacts.detail.nextFollowUp') }}</span>
              <strong>{{ formatDate(detailDialog.contact.next_follow_up_at) }}</strong>
            </li>
          </ul>
          <div class="tag-list">
            <Tag
              v-for="tag in detailDialog.contact.tags"
              :key="tag"
              :value="tag"
              severity="info"
            />
            <span v-if="!detailDialog.contact.tags?.length">{{ $t('contacts.detail.tagsEmpty') }}</span>
          </div>
        </section>
        <section>
          <h3>{{ $t('contacts.detail.recent') }}</h3>
          <div
            v-if="detailDialog.contact.activities?.length"
            class="timeline"
          >
            <div
              v-for="activity in detailDialog.contact.activities"
              :key="activity.id"
              class="timeline-item"
            >
              <div class="timeline-time">
                {{ formatDate(activity.created_at) }}
              </div>
              <div class="timeline-card">
                <p class="timeline-title">
                  <i :class="activityIcon(activity.type)" />
                  {{ activity.subject }}
                </p>
                <p class="timeline-body">
                  {{ activity.description || $t('contacts.detail.noNotes') }}
                </p>
                <small>{{ $t('contacts.detail.author', { name: activity.created_by?.name || $t('contacts.detail.authorFallback') }) }}</small>
              </div>
            </div>
          </div>
          <EmptyState
            v-else
            :message="$t('contacts.detail.noActivities')"
          />
        </section>
      </div>
    </Dialog>

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useConfirm } from 'primevue/useconfirm'
import Card from 'primevue/card'
import Button from 'primevue/button'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Dropdown from 'primevue/dropdown'
import InputText from 'primevue/inputtext'
import Dialog from 'primevue/dialog'
import ConfirmDialog from 'primevue/confirmdialog'
import Tag from 'primevue/tag'
import Chips from 'primevue/chips'
import Calendar from 'primevue/calendar'
import Textarea from 'primevue/textarea'
import AdvancedFilter from '@/components/shared/AdvancedFilter.vue'
import Badge from '@/components/shared/Badge.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import { contactService } from '@/services/contactService'

const toast = useToast()
const confirm = useConfirm()
const { t } = useI18n()

const contacts = ref([])
const loading = ref(false)
const pagination = reactive({
  limit: 20,
  offset: 0,
  total: 0
})

const filters = reactive({
  search: '',
  type: '',
  status: '',
  owner_id: ''
})

const contactDialog = reactive({
  visible: false,
  mode: 'create',
  saving: false,
  form: getDefaultForm()
})

const detailDialog = reactive({
  visible: false,
  contact: null
})

const typeOptions = computed(() => [
  { label: t('contacts.types.adopter'), value: 'adopter' },
  { label: t('contacts.types.donor'), value: 'donor' },
  { label: t('contacts.types.volunteer'), value: 'volunteer' },
  { label: t('contacts.types.partner'), value: 'partner' },
  { label: t('contacts.types.vendor'), value: 'vendor' },
  { label: t('contacts.types.other'), value: 'other' }
])

const statusOptions = computed(() => [
  { label: t('contacts.statuses.active'), value: 'active' },
  { label: t('contacts.statuses.prospect'), value: 'prospect' },
  { label: t('contacts.statuses.inactive'), value: 'inactive' },
  { label: t('contacts.statuses.archived'), value: 'archived' }
])

const ownerOptions = computed(() => {
  const map = new Map()
  contacts.value.forEach((contact) => {
    if (contact.owner_id && contact.owner_name) {
      map.set(contact.owner_id, contact.owner_name)
    }
  })
  return Array.from(map, ([value, label]) => ({ value, label }))
})

const contactStats = computed(() => {
  const now = new Date()
  const upcoming = contacts.value.filter((contact) => {
    if (!contact.next_follow_up_at) return false
    const date = new Date(contact.next_follow_up_at)
    return date > now && date.getTime() - now.getTime() < 1000 * 60 * 60 * 24 * 7
  }).length

  return [
    { label: t('contacts.stats.total'), value: pagination.total },
    {
      label: t('contacts.stats.active'),
      value: contacts.value.filter((contact) => contact.status === 'active').length,
      description: t('contacts.stats.activeDescription')
    },
    {
      label: t('contacts.stats.prospects'),
      value: contacts.value.filter((contact) => contact.status === 'prospect').length,
      description: t('contacts.stats.prospectsDescription')
    },
    {
      label: t('contacts.stats.upcoming'),
      value: upcoming,
      description: t('contacts.stats.upcomingDescription')
    }
  ]
})

function getDefaultForm() {
  return {
    id: null,
    first_name: '',
    last_name: '',
    organization: '',
    email: '',
    phone: '',
    type: 'adopter',
    status: 'prospect',
    owner_id: '',
    tags: [],
    next_follow_up_at: null,
    notes: ''
  }
}

const loadContacts = async () => {
  try {
    loading.value = true
    const response = await contactService.getContacts({
      limit: pagination.limit,
      offset: pagination.offset,
      search: filters.search || undefined,
      type: filters.type || undefined,
      status: filters.status || undefined,
      owner_id: filters.owner_id || undefined
    })
    contacts.value = response.data
    pagination.total = response.total
  } catch (error) {
    contacts.value = []
    pagination.total = 0
    showError(t('contacts.notifications.loadError'), error)
  } finally {
    loading.value = false
  }
}

const applyFilters = () => {
  pagination.offset = 0
  loadContacts()
}

const onPage = (event) => {
  pagination.limit = event.rows
  pagination.offset = event.first
  loadContacts()
}

const openCreateDialog = () => {
  contactDialog.mode = 'create'
  contactDialog.form = getDefaultForm()
  contactDialog.visible = true
}

const openEditDialog = (contact) => {
  contactDialog.mode = 'edit'
  contactDialog.form = {
    id: contact.id,
    first_name: contact.first_name,
    last_name: contact.last_name,
    organization: contact.organization,
    email: contact.email,
    phone: contact.phone,
    type: contact.type,
    status: contact.status,
    owner_id: contact.owner_id || '',
    tags: [...(contact.tags || [])],
    next_follow_up_at: contact.next_follow_up_at ? new Date(contact.next_follow_up_at) : null,
    notes: contact.notes || ''
  }
  contactDialog.visible = true
}

const saveContact = async () => {
  const payload = {
    first_name: contactDialog.form.first_name,
    last_name: contactDialog.form.last_name,
    organization: contactDialog.form.organization || undefined,
    email: contactDialog.form.email || undefined,
    phone: contactDialog.form.phone || undefined,
    type: contactDialog.form.type,
    status: contactDialog.form.status,
    owner_id: contactDialog.form.owner_id || undefined,
    tags: contactDialog.form.tags || [],
    next_follow_up_at: contactDialog.form.next_follow_up_at
      ? new Date(contactDialog.form.next_follow_up_at).toISOString()
      : null,
    notes: contactDialog.form.notes || undefined
  }

  try {
    contactDialog.saving = true
    if (contactDialog.mode === 'create') {
      await contactService.createContact(payload)
      toast.add({ severity: 'success', summary: t('contacts.notifications.createSuccess'), detail: t('contacts.notifications.createDetail'), life: 3000 })
    } else {
      await contactService.updateContact(contactDialog.form.id, payload)
      toast.add({ severity: 'success', summary: t('contacts.notifications.updateSuccess'), detail: t('contacts.notifications.updateDetail'), life: 3000 })
    }
    contactDialog.visible = false
    await loadContacts()
  } catch (error) {
    showError(t('contacts.notifications.saveError'), error)
  } finally {
    contactDialog.saving = false
  }
}

const openDetails = (event) => {
  detailDialog.contact = event.data
  detailDialog.visible = true
}

const confirmDelete = (contact) => {
  const name = getContactName(contact)
  confirm.require({
    message: t('contacts.confirmDeleteMessage', { name }),
    header: t('contacts.confirmDeleteTitle'),
    acceptLabel: t('common.delete'),
    rejectLabel: t('common.cancel'),
    icon: 'pi pi-exclamation-triangle',
    accept: async () => {
      try {
        await contactService.deleteContact(contact.id)
        toast.add({ severity: 'success', summary: t('contacts.notifications.deleteSuccess'), detail: t('contacts.notifications.deleteDetail'), life: 3000 })
        loadContacts()
      } catch (error) {
        showError(t('contacts.notifications.deleteError'), error)
      }
    }
  })
}

const exportContacts = () => {
  if (!contacts.value.length) {
    toast.add({ severity: 'info', summary: t('contacts.notifications.nothingToExport'), detail: t('contacts.notifications.nothingToExportDetail'), life: 3000 })
    return
  }
  const headers = [
    t('contacts.form.firstName'),
    t('contacts.form.lastName'),
    t('contacts.form.email'),
    t('contacts.form.phone'),
    t('contacts.form.type'),
    t('contacts.form.status'),
    t('contacts.form.owner')
  ]
  const rows = contacts.value.map((contact) => [
    contact.first_name,
    contact.last_name,
    contact.email || '',
    contact.phone || '',
    formatType(contact.type),
    formatStatus(contact.status),
    contact.owner_name || ''
  ])
  const csv = [headers, ...rows].map((row) => row.map((cell) => `"${cell ?? ''}"`).join(',')).join('\n')
  const blob = new Blob([csv], { type: 'text/csv;charset=utf-8;' })
  const url = window.URL.createObjectURL(blob)
  const link = document.createElement('a')
  link.href = url
  link.download = `contacts-${new Date().toISOString().slice(0, 10)}.csv`
  link.click()
  window.URL.revokeObjectURL(url)
  toast.add({ severity: 'success', summary: t('contacts.notifications.exportReady'), detail: t('contacts.notifications.exportReadyDetail'), life: 3000 })
}

const copyEmail = async (email) => {
  if (!email) return
  try {
    await navigator.clipboard.writeText(email)
    toast.add({ severity: 'info', summary: t('contacts.notifications.emailCopied'), detail: email, life: 2000 })
  } catch (error) {
    showError(t('contacts.notifications.copyError'), error)
  }
}

const getInitials = (contact) => {
  const first = contact.first_name?.[0] || ''
  const last = contact.last_name?.[0] || ''
  return `${first}${last}`.toUpperCase() || 'C'
}

const getContactName = (contact) => {
  if (!contact) return t('contacts.detail.title')
  const name = [contact.first_name, contact.last_name].filter(Boolean).join(' ').trim()
  return name || t('contacts.detail.title')
}

const formatDate = (value) => {
  if (!value) return t('common.notAvailable')
  return new Date(value).toLocaleDateString()
}

const formatType = (type) => {
  const option = typeOptions.value.find((item) => item.value === type)
  return option ? option.label : type
}

const formatStatus = (status) => {
  const option = statusOptions.value.find((item) => item.value === status)
  return option ? option.label : status
}

const getTypeVariant = (type) => {
  switch (type) {
    case 'donor':
      return 'success'
    case 'adopter':
      return 'info'
    case 'volunteer':
      return 'warning'
    case 'partner':
      return 'neutral'
    default:
      return 'secondary'
  }
}

const getStatusVariant = (status) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'prospect':
      return 'warning'
    case 'inactive':
      return 'neutral'
    case 'archived':
      return 'secondary'
    default:
      return 'neutral'
  }
}

const isOverdue = (value) => {
  if (!value) return false
  return new Date(value) < new Date()
}

const activityIcon = (type) => {
  switch (type) {
    case 'call':
      return 'pi pi-phone'
    case 'email':
      return 'pi pi-envelope'
    case 'meeting':
      return 'pi pi-calendar'
    default:
      return 'pi pi-comment'
  }
}

const showError = (summary, error) => {
  const detail = error?.response?.data?.error || error?.message || t('common.genericError')
  toast.add({ severity: 'error', summary, detail, life: 4000 })
}

onMounted(() => {
  loadContacts()
})
</script>

<style scoped>
.contacts-page {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1rem;
}

.page-title {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
}

.page-subtitle {
  margin: 0.35rem 0 0;
  color: var(--text-muted);
}

.header-actions {
  display: flex;
  gap: 0.75rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.stat-card {
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  background: var(--card-bg);
}

.stat-label {
  display: block;
  font-size: 0.85rem;
  color: var(--text-muted);
}

.stat-value {
  display: block;
  font-size: 1.65rem;
  color: var(--heading-color);
}

.table-card {
  overflow: hidden;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 1rem;
}

.contact-name {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.contact-title {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin: 0;
  font-weight: 600;
}

.contact-email {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.avatar {
  width: 42px;
  height: 42px;
  border-radius: 999px;
  background: var(--card-muted-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-color);
  font-weight: 600;
}

.avatar::after {
  content: attr(data-letter);
}

.action-buttons {
  display: flex;
  gap: 0.25rem;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.35rem;
}

.form-group label {
  font-weight: 600;
  color: var(--text-color);
}

.full-width {
  grid-column: 1 / -1;
}

.dialog-actions {
  grid-column: 1 / -1;
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
}

.detail-layout {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
}

.detail-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.detail-list li {
  display: flex;
  justify-content: space-between;
  color: var(--text-color);
}

.detail-list span {
  color: var(--text-muted);
}

.tag-list {
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
  margin-top: 1rem;
}

.timeline {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.timeline-item {
  display: grid;
  grid-template-columns: 120px auto;
  gap: 1rem;
}

.timeline-time {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.timeline-card {
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  padding: 0.75rem 1rem;
  background: var(--card-bg);
}

.timeline-title {
  margin: 0;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-weight: 600;
}

.timeline-body {
  margin: 0.35rem 0;
  color: var(--text-muted);
}

.overdue {
  color: #b91c1c;
  font-weight: 600;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    flex-wrap: wrap;
  }

  .timeline-item {
    grid-template-columns: 1fr;
  }
}
</style>
