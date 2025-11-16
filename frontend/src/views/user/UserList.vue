<template>
  <div class="user-management page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">
          User Management
        </h1>
        <p class="page-subtitle">
          Invite new teammates, adjust permissions, and keep your workspace secure.
        </p>
      </div>
      <div class="header-actions">
        <Button
          label="Refresh"
          icon="pi pi-refresh"
          class="p-button-text"
          @click="loadUsers"
        />
        <Button
          label="Invite User"
          icon="pi pi-user-plus"
          :disabled="!authStore.isAdmin"
          @click="openCreateDialog"
        />
      </div>
    </div>

    <Card class="stats-card">
      <template #content>
        <div class="stats-grid">
          <div
            v-for="stat in userStats"
            :key="stat.label"
            class="stat-card"
          >
            <span class="stat-label">{{ stat.label }}</span>
            <strong class="stat-value">{{ stat.value }}</strong>
            <small
              v-if="stat.description"
              class="stat-description"
            >{{ stat.description }}</small>
          </div>
        </div>
      </template>
    </Card>

    <AdvancedFilter
      v-model="filters"
      :show-date-range="false"
      :show-export="false"
      search-placeholder="Search by name or email"
      @filter="applyFilters"
    >
      <div class="filter-field">
        <label>Role</label>
        <Dropdown
          v-model="filters.role"
          :options="roleOptions"
          option-label="label"
          option-value="value"
          placeholder="All roles"
          show-clear
          @change="applyFilters"
        />
      </div>
      <div class="filter-field">
        <label>Status</label>
        <Dropdown
          v-model="filters.status"
          :options="statusOptions"
          option-label="label"
          option-value="value"
          placeholder="All statuses"
          show-clear
          @change="applyFilters"
        />
      </div>
    </AdvancedFilter>

    <Card class="table-card">
      <template #content>
        <DataTable
          :value="users"
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
            header="Name"
          >
            <template #body="slotProps">
              <div class="user-name-cell">
                <div
                  class="avatar"
                  :data-letter="getInitials(slotProps.data)"
                >
                  <img
                    v-if="slotProps.data.avatar"
                    :src="slotProps.data.avatar"
                    :alt="slotProps.data.first_name"
                  >
                </div>
                <div>
                  <p class="user-name">
                    {{ slotProps.data.first_name }} {{ slotProps.data.last_name }}
                  </p>
                  <span class="user-email">{{ slotProps.data.email }}</span>
                </div>
              </div>
            </template>
          </Column>
          <Column
            field="role"
            header="Role"
          >
            <template #body="slotProps">
              <Badge :variant="getRoleVariant(slotProps.data.role)">
                {{ formatRole(slotProps.data.role) }}
              </Badge>
            </template>
          </Column>
          <Column
            field="status"
            header="Status"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ formatStatus(slotProps.data.status) }}
              </Badge>
            </template>
          </Column>
          <Column
            field="phone"
            header="Phone"
          >
            <template #body="slotProps">
              {{ slotProps.data.phone || '—' }}
            </template>
          </Column>
          <Column
            field="last_login"
            header="Last Login"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.last_login) }}
            </template>
          </Column>
          <Column
            header="Actions"
            :style="{ width: '150px' }"
          >
            <template #body="slotProps">
              <div
                class="action-buttons"
                @click.stop
              >
                <Button
                  v-tooltip.top="'View details'"
                  icon="pi pi-eye"
                  class="p-button-rounded p-button-text"
                  @click="openDetails({ data: slotProps.data })"
                />
                <Button
                  v-tooltip.top="'Edit user'"
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  :disabled="!authStore.isAdmin"
                  @click="openEditDialog(slotProps.data)"
                />
                <Button
                  v-tooltip.top="'More actions'"
                  icon="pi pi-ellipsis-v"
                  class="p-button-rounded p-button-text"
                  @click="openActionMenu($event, slotProps.data)"
                />
              </div>
            </template>
          </Column>
        </DataTable>
      </template>
    </Card>

    <LoadingSpinner v-if="loading && users.length === 0" />
    <EmptyState
      v-if="!loading && users.length === 0"
      message="No users found for the selected filters."
    />

    <!-- Create/Edit dialog -->
    <Dialog
      v-model:visible="userDialog.visible"
      modal
      :header="userDialog.mode === 'create' ? 'Invite User' : 'Update User'"
      :style="{ width: '520px' }"
    >
      <form
        class="form-grid"
        @submit.prevent="saveUser"
      >
        <div class="form-group">
          <label for="firstName">First name</label>
          <InputText
            id="firstName"
            v-model="userDialog.form.first_name"
            required
          />
        </div>
        <div class="form-group">
          <label for="lastName">Last name</label>
          <InputText
            id="lastName"
            v-model="userDialog.form.last_name"
            required
          />
        </div>
        <div class="form-group">
          <label for="email">Email</label>
          <InputText
            id="email"
            v-model="userDialog.form.email"
            :disabled="userDialog.mode === 'edit'"
            type="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="phone">Phone</label>
          <InputText
            id="phone"
            v-model="userDialog.form.phone"
          />
        </div>
        <div class="form-group">
          <label for="role">Role</label>
          <Dropdown
            id="role"
            v-model="userDialog.form.role"
            :options="roleOptions"
            option-label="label"
            option-value="value"
            required
          />
        </div>
        <div class="form-group">
          <label for="status">Status</label>
          <Dropdown
            id="status"
            v-model="userDialog.form.status"
            :options="statusOptions"
            option-label="label"
            option-value="value"
          />
        </div>
        <div class="form-group">
          <label for="language">Language</label>
          <Dropdown
            id="language"
            v-model="userDialog.form.language"
            :options="languageOptions"
            option-label="label"
            option-value="value"
          />
        </div>
        <div class="form-group">
          <label for="theme">Theme</label>
          <Dropdown
            id="theme"
            v-model="userDialog.form.theme"
            :options="themeOptions"
            option-label="label"
            option-value="value"
          />
        </div>
        <div
          v-if="userDialog.mode === 'create'"
          class="form-group"
        >
          <label for="password">Temporary password</label>
          <InputText
            id="password"
            v-model="userDialog.form.password"
            type="password"
            required
            minlength="8"
          />
        </div>

        <div class="dialog-actions">
          <Button
            label="Cancel"
            text
            type="button"
            @click="userDialog.visible = false"
          />
          <Button
            label="Save"
            icon="pi pi-check"
            type="submit"
            :loading="userDialog.saving"
            :disabled="userDialog.saving"
          />
        </div>
      </form>
    </Dialog>

    <!-- Detail dialog -->
    <Dialog
      v-model:visible="detailDialog.visible"
      modal
      header="User details"
      :style="{ width: '480px' }"
    >
      <div
        v-if="detailDialog.user"
        class="detail-panel"
      >
        <div class="detail-header">
          <div
            class="avatar large"
            :data-letter="getInitials(detailDialog.user)"
          >
            <img
              v-if="detailDialog.user.avatar"
              :src="detailDialog.user.avatar"
              :alt="detailDialog.user.first_name"
            >
          </div>
          <div>
            <h3>{{ detailDialog.user.first_name }} {{ detailDialog.user.last_name }}</h3>
            <p>{{ detailDialog.user.email }}</p>
            <Badge :variant="getStatusVariant(detailDialog.user.status)">
              {{ formatStatus(detailDialog.user.status) }}
            </Badge>
          </div>
        </div>
        <ul class="detail-list">
          <li>
            <span>Role</span>
            <strong>{{ formatRole(detailDialog.user.role) }}</strong>
          </li>
          <li>
            <span>Phone</span>
            <strong>{{ detailDialog.user.phone || '—' }}</strong>
          </li>
          <li>
            <span>Language</span>
            <strong>{{ detailDialog.user.language?.toUpperCase() }}</strong>
          </li>
          <li>
            <span>Theme</span>
            <strong>{{ detailDialog.user.theme }}</strong>
          </li>
          <li>
            <span>Last login</span>
            <strong>{{ formatDate(detailDialog.user.last_login) }}</strong>
          </li>
          <li>
            <span>Created</span>
            <strong>{{ formatDate(detailDialog.user.created_at) }}</strong>
          </li>
          <li>
            <span>Updated</span>
            <strong>{{ formatDate(detailDialog.user.updated_at) }}</strong>
          </li>
        </ul>
        <div class="detail-actions">
          <Button
            label="Reset password"
            icon="pi pi-lock"
            class="p-button-text"
            :disabled="!authStore.isAdmin"
            @click="openPasswordDialog(detailDialog.user)"
          />
          <Button
            label="Deactivate"
            icon="pi pi-user-minus"
            class="p-button-danger p-button-text"
            :disabled="!authStore.isAdmin || detailDialog.user.status !== 'active'"
            @click="setUserStatus(detailDialog.user, 'inactive')"
          />
        </div>
      </div>
    </Dialog>

    <!-- Reset password dialog -->
    <Dialog
      v-model:visible="passwordDialog.visible"
      modal
      header="Reset password"
      :style="{ width: '420px' }"
    >
      <form
        class="form-grid"
        @submit.prevent="resetPassword"
      >
        <p class="dialog-hint">
          Generate a secure password and share it with the user using a trusted channel. They will be asked to change it
          on first login.
        </p>
        <div class="form-group">
          <label for="newPassword">New password</label>
          <InputText
            id="newPassword"
            v-model="passwordDialog.password"
            type="password"
            minlength="8"
            required
          />
        </div>
        <div class="dialog-actions">
          <Button
            label="Cancel"
            text
            type="button"
            @click="passwordDialog.visible = false"
          />
          <Button
            label="Reset password"
            icon="pi pi-check"
            type="submit"
            :loading="passwordDialog.saving"
            :disabled="passwordDialog.saving"
          />
        </div>
      </form>
    </Dialog>

    <ConfirmDialog />
  </div>
</template>

<script setup>
import { onMounted, reactive, ref, computed } from 'vue'
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
import AdvancedFilter from '@/components/shared/AdvancedFilter.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import Badge from '@/components/shared/Badge.vue'
import { userService } from '@/services/userService'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const toast = useToast()
const confirm = useConfirm()

const users = ref([])
const loading = ref(false)
const pagination = reactive({
  limit: 20,
  offset: 0,
  total: 0
})

const filters = reactive({
  search: '',
  role: '',
  status: ''
})

const userDialog = reactive({
  visible: false,
  mode: 'create',
  saving: false,
  form: getDefaultForm()
})

const detailDialog = reactive({
  visible: false,
  user: null
})

const passwordDialog = reactive({
  visible: false,
  saving: false,
  userId: null,
  password: ''
})

const roleOptions = [
  { label: 'Super Admin', value: 'super_admin' },
  { label: 'Admin', value: 'admin' },
  { label: 'Employee', value: 'employee' },
  { label: 'Volunteer', value: 'volunteer' },
  { label: 'User', value: 'user' }
]

const statusOptions = [
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'Suspended', value: 'suspended' }
]

const languageOptions = [
  { label: 'Polish', value: 'pl' },
  { label: 'English', value: 'en' }
]

const themeOptions = [
  { label: 'Light', value: 'light' },
  { label: 'Dark', value: 'dark' }
]

const userStats = computed(() => {
  const counts = users.value.reduce(
    (acc, user) => {
      acc[user.status] = (acc[user.status] || 0) + 1
      acc[user.role] = (acc[user.role] || 0) + 1
      return acc
    },
    {}
  )

  return [
    { label: 'Total users', value: pagination.total },
    { label: 'Active', value: counts.active || 0, description: 'Authorized today' },
    { label: 'Administrators', value: (counts.super_admin || 0) + (counts.admin || 0) },
    { label: 'Inactive & suspended', value: (counts.inactive || 0) + (counts.suspended || 0) }
  ]
})

function getDefaultForm() {
  return {
    id: null,
    first_name: '',
    last_name: '',
    email: '',
    phone: '',
    role: 'employee',
    status: 'active',
    language: 'pl',
    theme: 'light',
    password: ''
  }
}

const applyFilters = () => {
  pagination.offset = 0
  loadUsers()
}

const onPage = (event) => {
  pagination.limit = event.rows
  pagination.offset = event.first
  loadUsers()
}

const loadUsers = async () => {
  try {
    loading.value = true
    const response = await userService.getUsers({
      limit: pagination.limit,
      offset: pagination.offset,
      search: filters.search || undefined,
      role: filters.role || undefined,
      status: filters.status || undefined
    })

    users.value = response.data
    pagination.total = response.total
  } catch (error) {
    showError('Unable to load users', error)
  } finally {
    loading.value = false
  }
}

const openCreateDialog = () => {
  userDialog.mode = 'create'
  userDialog.form = getDefaultForm()
  userDialog.visible = true
}

const openEditDialog = (user) => {
  userDialog.mode = 'edit'
  userDialog.form = {
    id: user.id,
    email: user.email,
    first_name: user.first_name,
    last_name: user.last_name,
    phone: user.phone || '',
    role: user.role,
    status: user.status,
    language: user.language || 'pl',
    theme: user.theme || 'light'
  }
  userDialog.visible = true
}

const saveUser = async () => {
  try {
    userDialog.saving = true
    if (userDialog.mode === 'create') {
      await userService.createUser({
        email: userDialog.form.email,
        password: userDialog.form.password,
        first_name: userDialog.form.first_name,
        last_name: userDialog.form.last_name,
        phone: userDialog.form.phone || undefined,
        role: userDialog.form.role,
        status: userDialog.form.status,
        language: userDialog.form.language,
        theme: userDialog.form.theme
      })
      toast.add({ severity: 'success', summary: 'User invited', detail: 'A new account has been created.', life: 3000 })
    } else {
      const payload = {
        first_name: userDialog.form.first_name,
        last_name: userDialog.form.last_name,
        phone: userDialog.form.phone || undefined,
        role: userDialog.form.role,
        status: userDialog.form.status,
        language: userDialog.form.language,
        theme: userDialog.form.theme
      }
      await userService.updateUser(userDialog.form.id, payload)
      toast.add({ severity: 'success', summary: 'User updated', detail: 'Changes were saved successfully.', life: 3000 })
    }
    userDialog.visible = false
    await loadUsers()
  } catch (error) {
    showError('Unable to save user', error)
  } finally {
    userDialog.saving = false
  }
}

const openDetails = (event) => {
  detailDialog.user = event.data
  detailDialog.visible = true
}

const openPasswordDialog = (user) => {
  passwordDialog.userId = user.id
  passwordDialog.password = ''
  passwordDialog.visible = true
}

const resetPassword = async () => {
  if (!passwordDialog.userId) return

  try {
    passwordDialog.saving = true
    await userService.resetPassword(passwordDialog.userId, passwordDialog.password)
    toast.add({ severity: 'success', summary: 'Password reset', detail: 'The user will need to sign in again.', life: 3000 })
    passwordDialog.visible = false
  } catch (error) {
    showError('Unable to reset password', error)
  } finally {
    passwordDialog.saving = false
  }
}

const setUserStatus = (user, status) => {
  confirm.require({
    message: `Are you sure you want to mark ${user.first_name} as ${status}?`,
    header: 'Confirm status change',
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: 'Yes, update',
    rejectLabel: 'Cancel',
    accept: async () => {
      try {
        await userService.updateUser(user.id, { status })
        toast.add({ severity: 'success', summary: 'Status updated', detail: 'The account state has changed.', life: 3000 })
        await loadUsers()
        detailDialog.visible = false
      } catch (error) {
        showError('Unable to update status', error)
      }
    }
  })
}

const openActionMenu = (event, user) => {
  event.preventDefault()
  event.stopPropagation()
  confirm.require({
    message: 'Select action for this user',
    header: `${user.first_name} ${user.last_name}`,
    acceptLabel: 'Reset password',
    rejectLabel: user.status === 'active' ? 'Deactivate' : 'Activate',
    accept: () => openPasswordDialog(user),
    reject: () => setUserStatus(user, user.status === 'active' ? 'inactive' : 'active')
  })
}

const formatDate = (value) => {
  if (!value) return 'Never'
  return new Date(value).toLocaleString()
}

const formatRole = (role) => {
  const option = roleOptions.find((item) => item.value === role)
  return option ? option.label : role
}

const formatStatus = (status) => {
  const option = statusOptions.find((item) => item.value === status)
  return option ? option.label : status
}

const getStatusVariant = (status) => {
  switch (status) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'neutral'
    case 'suspended':
      return 'danger'
    default:
      return 'neutral'
  }
}

const getRoleVariant = (role) => {
  switch (role) {
    case 'super_admin':
      return 'danger'
    case 'admin':
      return 'info'
    case 'employee':
      return 'success'
    case 'volunteer':
      return 'warning'
    default:
      return 'neutral'
  }
}

const getInitials = (user) => {
  const first = user.first_name?.charAt(0) || ''
  const last = user.last_name?.charAt(0) || ''
  return `${first}${last}`.toUpperCase()
}

const showError = (message, error) => {
  const detail = error?.response?.data?.error || error.message || 'Unexpected error'
  toast.add({ severity: 'error', summary: message, detail, life: 4000 })
}

onMounted(() => {
  loadUsers()
})
</script>

<style scoped>
.user-management {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
}

.page-title {
  margin: 0;
  font-size: 2rem;
  font-weight: 700;
}

.page-subtitle {
  margin: 0.35rem 0 0;
  color: #6b7280;
}

.header-actions {
  display: flex;
  gap: 0.5rem;
}

.stats-card {
  margin-top: -0.5rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.stat-card {
  padding: 1rem;
  border: 1px solid #e5e7eb;
  border-radius: 0.75rem;
  background: #f9fafb;
}

.stat-label {
  display: block;
  font-size: 0.85rem;
  color: #6b7280;
}

.stat-value {
  display: block;
  font-size: 1.7rem;
  color: #111827;
}

.stat-description {
  color: #6b7280;
}

.table-card {
  overflow: hidden;
}

.user-name-cell {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 999px;
  background: #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: #374151;
  text-transform: uppercase;
  overflow: hidden;
}

.avatar.large {
  width: 64px;
  height: 64px;
  font-size: 1.3rem;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar::after {
  content: attr(data-letter);
}

.user-name {
  margin: 0;
  font-weight: 600;
  color: #111827;
}

.user-email {
  font-size: 0.85rem;
  color: #6b7280;
}

.action-buttons {
  display: flex;
  align-items: center;
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
  color: #374151;
}

.dialog-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.5rem;
  margin-top: 1rem;
}

.detail-panel {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.detail-header {
  display: flex;
  gap: 1rem;
  align-items: center;
}

.detail-panel h3 {
  margin: 0;
  font-size: 1.25rem;
}

.detail-panel p {
  margin: 0.15rem 0 0.5rem;
  color: #6b7280;
}

.detail-list {
  list-style: none;
  padding: 0;
  margin: 0;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.detail-list li {
  display: flex;
  justify-content: space-between;
  font-size: 0.95rem;
  color: #374151;
}

.detail-list span {
  color: #9ca3af;
}

.detail-actions {
  display: flex;
  gap: 0.5rem;
}

.dialog-hint {
  color: #6b7280;
  margin-top: 0;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .header-actions {
    width: 100%;
    justify-content: flex-start;
    flex-wrap: wrap;
  }

  .detail-panel {
    text-align: left;
  }
}
</style>
