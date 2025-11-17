<template>
  <div class="user-management page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">
          {{ $t('users.title') }}
        </h1>
        <p class="page-subtitle">
          {{ $t('users.subtitle') }}
        </p>
      </div>
      <div class="header-actions">
        <Button
          :label="$t('users.actions.refresh')"
          icon="pi pi-refresh"
          class="p-button-text"
          @click="loadUsers"
        />
        <Button
          :label="$t('users.actions.inviteUser')"
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
      :search-placeholder="$t('users.searchPlaceholder')"
      @filter="applyFilters"
    >
      <div class="filter-field">
        <label>{{ $t('users.filters.role') }}</label>
        <Dropdown
          v-model="filters.role"
          :options="roleOptions"
          option-label="label"
          option-value="value"
          :placeholder="$t('users.filters.allRoles')"
          show-clear
          @change="applyFilters"
        />
      </div>
      <div class="filter-field">
        <label>{{ $t('users.filters.status') }}</label>
        <Dropdown
          v-model="filters.status"
          :options="statusOptions"
          option-label="label"
          option-value="value"
          :placeholder="$t('users.filters.allStatuses')"
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
            :header="$t('users.table.name')"
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
            :header="$t('users.table.role')"
          >
            <template #body="slotProps">
              <Badge :variant="getRoleVariant(slotProps.data.role)">
                {{ formatRole(slotProps.data.role) }}
              </Badge>
            </template>
          </Column>
          <Column
            field="status"
            :header="$t('users.table.status')"
          >
            <template #body="slotProps">
              <Badge :variant="getStatusVariant(slotProps.data.status)">
                {{ formatStatus(slotProps.data.status) }}
              </Badge>
            </template>
          </Column>
          <Column
            field="phone"
            :header="$t('users.table.phone')"
          >
            <template #body="slotProps">
              {{ slotProps.data.phone || $t('common.notAvailable') }}
            </template>
          </Column>
          <Column
            field="last_login"
            :header="$t('users.table.lastLogin')"
          >
            <template #body="slotProps">
              {{ formatDate(slotProps.data.last_login) }}
            </template>
          </Column>
          <Column
            :header="$t('common.actions')"
            :style="{ width: '150px' }"
          >
            <template #body="slotProps">
              <div
                class="action-buttons"
                @click.stop
              >
                <Button
                  v-tooltip.top="$t('users.tooltips.viewDetails')"
                  icon="pi pi-eye"
                  class="p-button-rounded p-button-text"
                  @click="openDetails({ data: slotProps.data })"
                />
                <Button
                  v-tooltip.top="$t('users.tooltips.editUser')"
                  icon="pi pi-pencil"
                  class="p-button-rounded p-button-text"
                  :disabled="!authStore.isAdmin"
                  @click="openEditDialog(slotProps.data)"
                />
                <Button
                  v-tooltip.top="$t('users.tooltips.moreActions')"
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
      :message="$t('users.messages.empty')"
    />

    <!-- Create/Edit dialog -->
    <Dialog
      v-model:visible="userDialog.visible"
      modal
      :header="userDialog.mode === 'create' ? $t('users.dialog.createTitle') : $t('users.dialog.editTitle')"
      :style="{ width: '520px' }"
    >
      <form
        class="form-grid"
        @submit.prevent="saveUser"
      >
        <div class="form-group">
          <label for="firstName">{{ $t('users.form.firstName') }}</label>
          <InputText
            id="firstName"
            v-model="userDialog.form.first_name"
            required
          />
        </div>
        <div class="form-group">
          <label for="lastName">{{ $t('users.form.lastName') }}</label>
          <InputText
            id="lastName"
            v-model="userDialog.form.last_name"
            required
          />
        </div>
        <div class="form-group">
          <label for="email">{{ $t('users.form.email') }}</label>
          <InputText
            id="email"
            v-model="userDialog.form.email"
            :disabled="userDialog.mode === 'edit'"
            type="email"
            required
          />
        </div>
        <div class="form-group">
          <label for="phone">{{ $t('users.form.phone') }}</label>
          <InputText
            id="phone"
            v-model="userDialog.form.phone"
          />
        </div>
        <div class="form-group">
          <label for="role">{{ $t('users.form.role') }}</label>
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
          <label for="status">{{ $t('users.form.status') }}</label>
          <Dropdown
            id="status"
            v-model="userDialog.form.status"
            :options="statusOptions"
            option-label="label"
            option-value="value"
          />
        </div>
        <div class="form-group">
          <label for="language">{{ $t('users.form.language') }}</label>
          <Dropdown
            id="language"
            v-model="userDialog.form.language"
            :options="languageOptions"
            option-label="label"
            option-value="value"
          />
        </div>
        <div class="form-group">
          <label for="theme">{{ $t('users.form.theme') }}</label>
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
          <label for="password">{{ $t('users.form.password') }}</label>
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
            :label="$t('common.cancel')"
            text
            type="button"
            @click="userDialog.visible = false"
          />
          <Button
            :label="$t('common.save')"
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
      :header="$t('users.dialog.detailTitle')"
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
            <span>{{ $t('users.detail.role') }}</span>
            <strong>{{ formatRole(detailDialog.user.role) }}</strong>
          </li>
          <li>
            <span>{{ $t('users.detail.phone') }}</span>
            <strong>{{ detailDialog.user.phone || $t('common.notAvailable') }}</strong>
          </li>
          <li>
            <span>{{ $t('users.detail.language') }}</span>
            <strong>{{ formatLanguage(detailDialog.user.language) }}</strong>
          </li>
          <li>
            <span>{{ $t('users.detail.theme') }}</span>
            <strong>{{ formatTheme(detailDialog.user.theme) }}</strong>
          </li>
          <li>
            <span>{{ $t('users.detail.lastLogin') }}</span>
            <strong>{{ formatDate(detailDialog.user.last_login) }}</strong>
          </li>
          <li>
            <span>{{ $t('users.detail.created') }}</span>
            <strong>{{ formatDate(detailDialog.user.created_at) }}</strong>
          </li>
          <li>
            <span>{{ $t('users.detail.updated') }}</span>
            <strong>{{ formatDate(detailDialog.user.updated_at) }}</strong>
          </li>
        </ul>
        <div class="detail-actions">
          <Button
            :label="$t('users.actions.resetPassword')"
            icon="pi pi-lock"
            class="p-button-text"
            :disabled="!authStore.isAdmin"
            @click="openPasswordDialog(detailDialog.user)"
          />
          <Button
            :label="$t('users.actions.deactivate')"
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
      :header="$t('users.dialog.passwordTitle')"
      :style="{ width: '420px' }"
    >
      <form
        class="form-grid"
        @submit.prevent="resetPassword"
      >
        <p class="dialog-hint">
          {{ $t('users.password.hint') }}
        </p>
        <div class="form-group">
          <label for="newPassword">{{ $t('users.password.newPassword') }}</label>
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
            :label="$t('common.cancel')"
            text
            type="button"
            @click="passwordDialog.visible = false"
          />
          <Button
            :label="$t('users.actions.resetPassword')"
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
import AdvancedFilter from '@/components/shared/AdvancedFilter.vue'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import EmptyState from '@/components/shared/EmptyState.vue'
import Badge from '@/components/shared/Badge.vue'
import { userService } from '@/services/userService'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const toast = useToast()
const confirm = useConfirm()
const { t } = useI18n()

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

const roleOptions = computed(() => [
  { label: t('users.roles.super_admin'), value: 'super_admin' },
  { label: t('users.roles.admin'), value: 'admin' },
  { label: t('users.roles.employee'), value: 'employee' },
  { label: t('users.roles.volunteer'), value: 'volunteer' },
  { label: t('users.roles.user'), value: 'user' }
])

const statusOptions = computed(() => [
  { label: t('users.statuses.active'), value: 'active' },
  { label: t('users.statuses.inactive'), value: 'inactive' },
  { label: t('users.statuses.suspended'), value: 'suspended' }
])

const languageOptions = computed(() => [
  { label: t('users.languages.pl'), value: 'pl' },
  { label: t('users.languages.en'), value: 'en' }
])

const themeOptions = computed(() => [
  { label: t('users.themes.light'), value: 'light' },
  { label: t('users.themes.dark'), value: 'dark' }
])

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
    { label: t('users.stats.total'), value: pagination.total },
    { label: t('users.stats.active'), value: counts.active || 0, description: t('users.stats.activeDescription') },
    { label: t('users.stats.admins'), value: (counts.super_admin || 0) + (counts.admin || 0) },
    { label: t('users.stats.inactiveSuspended'), value: (counts.inactive || 0) + (counts.suspended || 0) }
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
    showError(t('users.notifications.loadError'), error)
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
      toast.add({ severity: 'success', summary: t('users.notifications.inviteSuccess'), detail: t('users.notifications.inviteDetail'), life: 3000 })
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
      toast.add({ severity: 'success', summary: t('users.notifications.updateSuccess'), detail: t('users.notifications.updateDetail'), life: 3000 })
    }
    userDialog.visible = false
    await loadUsers()
  } catch (error) {
    showError(t('users.notifications.saveError'), error)
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
    toast.add({ severity: 'success', summary: t('users.notifications.passwordReset'), detail: t('users.notifications.passwordResetDetail'), life: 3000 })
    passwordDialog.visible = false
  } catch (error) {
    showError(t('users.notifications.passwordError'), error)
  } finally {
    passwordDialog.saving = false
  }
}

const setUserStatus = (user, status) => {
  const name = [user.first_name, user.last_name].filter(Boolean).join(' ').trim() || user.email
  confirm.require({
    message: t('users.confirmStatusMessage', { name, status: formatStatus(status) }),
    header: t('users.confirmStatusTitle'),
    icon: 'pi pi-exclamation-triangle',
    acceptLabel: t('users.actions.confirmUpdate'),
    rejectLabel: t('common.cancel'),
    accept: async () => {
      try {
        await userService.updateUser(user.id, { status })
        toast.add({ severity: 'success', summary: t('users.notifications.updateStatusSuccess'), detail: t('users.notifications.updateStatusDetail'), life: 3000 })
        await loadUsers()
        detailDialog.visible = false
      } catch (error) {
        showError(t('users.notifications.updateStatusError'), error)
      }
    }
  })
}

const openActionMenu = (event, user) => {
  event.preventDefault()
  event.stopPropagation()
  confirm.require({
    message: t('users.actions.selectAction'),
    header: `${user.first_name} ${user.last_name}`,
    acceptLabel: t('users.actions.resetPassword'),
    rejectLabel: user.status === 'active' ? t('users.actions.deactivate') : t('users.actions.activate'),
    accept: () => openPasswordDialog(user),
    reject: () => setUserStatus(user, user.status === 'active' ? 'inactive' : 'active')
  })
}

const formatDate = (value) => {
  if (!value) return t('users.messages.neverLoggedIn')
  return new Date(value).toLocaleString()
}

const formatRole = (role) => {
  const option = roleOptions.value.find((item) => item.value === role)
  return option ? option.label : role
}

const formatStatus = (status) => {
  const option = statusOptions.value.find((item) => item.value === status)
  return option ? option.label : status
}

const formatLanguage = (language) => {
  if (!language) return t('common.notAvailable')
  const option = languageOptions.value.find((item) => item.value === language)
  return option ? option.label : language.toUpperCase()
}

const formatTheme = (theme) => {
  if (!theme) return t('common.notAvailable')
  const option = themeOptions.value.find((item) => item.value === theme)
  return option ? option.label : theme
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
  const detail = error?.response?.data?.error || error.message || t('common.genericError')
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
  color: var(--text-muted);
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
  font-size: 1.7rem;
  color: var(--heading-color);
}

.stat-description {
  color: var(--text-muted);
}

.table-card {
  overflow: hidden;
  background: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 1rem;
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
  background: var(--card-muted-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  color: var(--text-color);
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
  color: var(--heading-color);
}

.user-email {
  font-size: 0.85rem;
  color: var(--text-muted);
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
  color: var(--text-color);
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
  color: var(--text-muted);
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
  color: var(--text-color);
}

.detail-list span {
  color: var(--text-muted);
}

.detail-actions {
  display: flex;
  gap: 0.5rem;
}

.dialog-hint {
  color: var(--text-muted);
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
