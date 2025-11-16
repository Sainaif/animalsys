<template>
  <div class="profile-page page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">
          My Profile
        </h1>
        <p class="page-subtitle">
          Review your account details, update preferences, and keep your login secure.
        </p>
      </div>
      <Button
        label="Sign out"
        icon="pi pi-sign-out"
        class="p-button-outlined"
        @click="handleSignOut"
      />
    </div>

    <div class="profile-grid">
      <Card class="profile-card">
        <template #title>
          <div class="card-title">
            <div
              class="avatar"
              :data-letter="userInitials"
            >
              <img
                v-if="authStore.user?.avatar"
                :src="authStore.user?.avatar"
                alt="Avatar"
              >
            </div>
            <div>
              <h2>{{ authStore.user?.first_name }} {{ authStore.user?.last_name }}</h2>
              <p>{{ authStore.user?.email }}</p>
              <Badge :variant="statusVariant">
                {{ formatStatus(authStore.user?.status) }}
              </Badge>
            </div>
          </div>
        </template>
        <template #content>
          <form
            class="form-grid"
            @submit.prevent="saveProfile"
          >
            <div class="form-group">
              <label for="firstName">First name</label>
              <InputText
                id="firstName"
                v-model="profileForm.first_name"
                required
              />
            </div>
            <div class="form-group">
              <label for="lastName">Last name</label>
              <InputText
                id="lastName"
                v-model="profileForm.last_name"
                required
              />
            </div>
            <div class="form-group full-width">
              <label for="email">Email</label>
              <InputText
                id="email"
                :value="profileForm.email"
                type="email"
                disabled
              />
              <small>This address is used for login and notifications.</small>
            </div>
            <div class="form-group">
              <label for="phone">Phone number</label>
              <InputText
                id="phone"
                v-model="profileForm.phone"
                placeholder="+48 123 456 789"
              />
            </div>
            <div class="form-group">
              <label for="language">Language</label>
              <Dropdown
                id="language"
                v-model="profileForm.language"
                :options="languageOptions"
                option-label="label"
                option-value="value"
              />
            </div>
            <div class="form-group">
              <label for="theme">Theme</label>
              <Dropdown
                id="theme"
                v-model="profileForm.theme"
                :options="themeOptions"
                option-label="label"
                option-value="value"
              />
            </div>
            <div class="form-group full-width">
              <label for="timezone">Time zone</label>
              <InputText
                id="timezone"
                :value="timezone"
                disabled
              />
            </div>
            <div class="form-actions">
              <Button
                type="submit"
                label="Save changes"
                icon="pi pi-check"
                :loading="savingProfile"
              />
            </div>
          </form>
        </template>
      </Card>

      <Card class="profile-card">
        <template #title>
          <h2>Account activity</h2>
        </template>
        <template #content>
          <ul class="detail-list">
            <li>
              <span>Role</span>
              <strong>{{ formatRole(authStore.user?.role) }}</strong>
            </li>
            <li>
              <span>Status</span>
              <strong>{{ formatStatus(authStore.user?.status) }}</strong>
            </li>
            <li>
              <span>Last login</span>
              <strong>{{ formatDate(authStore.user?.last_login) }}</strong>
            </li>
            <li>
              <span>Member since</span>
              <strong>{{ formatDate(authStore.user?.created_at) }}</strong>
            </li>
            <li>
              <span>Last profile update</span>
              <strong>{{ formatDate(authStore.user?.updated_at) }}</strong>
            </li>
            <li>
              <span>Preferred channel</span>
              <strong>{{ profileForm.language === 'pl' ? 'Polski' : 'English' }}</strong>
            </li>
          </ul>
          <div class="tip-box">
            <i class="pi pi-info-circle" />
            <p>Need to change your email or role? Contact an administrator.</p>
          </div>
        </template>
      </Card>
    </div>

    <div class="profile-grid">
      <Card class="profile-card">
        <template #title>
          <h2>Security</h2>
        </template>
        <template #content>
          <form
            class="form-grid"
            @submit.prevent="changePassword"
          >
            <div class="form-group full-width">
              <label for="currentPassword">Current password</label>
              <InputText
                id="currentPassword"
                v-model="passwordForm.current"
                type="password"
                autocomplete="current-password"
              />
            </div>
            <div class="form-group">
              <label for="newPassword">New password</label>
              <InputText
                id="newPassword"
                v-model="passwordForm.new"
                type="password"
                autocomplete="new-password"
                minlength="8"
                required
              />
            </div>
            <div class="form-group">
              <label for="confirmPassword">Confirm new password</label>
              <InputText
                id="confirmPassword"
                v-model="passwordForm.confirm"
                type="password"
                autocomplete="new-password"
                minlength="8"
                required
              />
            </div>
            <div class="form-actions">
              <Button
                type="submit"
                label="Update password"
                icon="pi pi-lock"
                :loading="changingPassword"
              />
            </div>
          </form>
          <Message
            severity="info"
            :closable="false"
            class="security-hint"
          >
            Use at least 8 characters with a mix of letters, numbers, and symbols.
          </Message>
        </template>
      </Card>

      <Card class="profile-card">
        <template #title>
          <h2>Sessions</h2>
        </template>
        <template #content>
          <p>You're currently logged in on this device.</p>
          <ul class="detail-list">
            <li>
              <span>Browser</span>
              <strong>{{ browserName }}</strong>
            </li>
            <li>
              <span>Platform</span>
              <strong>{{ platform }}</strong>
            </li>
            <li>
              <span>IP address</span>
              <strong>{{ ipAddress }}</strong>
            </li>
          </ul>
          <div class="session-actions">
            <Button
              label="Sign out everywhere"
              icon="pi pi-shield"
              class="p-button-text p-button-danger"
              @click="signOutAll"
            />
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
import Card from 'primevue/card'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Dropdown from 'primevue/dropdown'
import Message from 'primevue/message'
import Badge from '@/components/shared/Badge.vue'
import { useAuthStore } from '@/stores/auth'

const authStore = useAuthStore()
const toast = useToast()
const { locale } = useI18n()

const profileForm = ref({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  language: 'pl',
  theme: 'light'
})

const passwordForm = ref({
  current: '',
  new: '',
  confirm: ''
})

const savingProfile = ref(false)
const changingPassword = ref(false)

const languageOptions = [
  { label: 'Polski', value: 'pl' },
  { label: 'English', value: 'en' }
]

const themeOptions = [
  { label: 'Light', value: 'light' },
  { label: 'Dark', value: 'dark' }
]

const timezone = new Intl.DateTimeFormat().resolvedOptions().timeZone
const browserName = navigator.userAgent
const platform = navigator.platform
const ipAddress = 'Current session'

const userInitials = computed(() => {
  const first = authStore.user?.first_name?.[0] || ''
  const last = authStore.user?.last_name?.[0] || ''
  return `${first}${last}`.toUpperCase() || 'ME'
})

const statusVariant = computed(() => {
  switch (authStore.user?.status) {
    case 'active':
      return 'success'
    case 'inactive':
      return 'neutral'
    case 'suspended':
      return 'danger'
    default:
      return 'neutral'
  }
})

watch(
  () => authStore.user,
  (user) => {
    if (!user) return
    profileForm.value = {
      first_name: user.first_name,
      last_name: user.last_name,
      email: user.email,
      phone: user.phone || '',
      language: user.language || 'pl',
      theme: user.theme || 'light'
    }
  },
  { immediate: true }
)

const saveProfile = async () => {
  try {
    savingProfile.value = true
    await authStore.updateProfile({
      first_name: profileForm.value.first_name,
      last_name: profileForm.value.last_name,
      phone: profileForm.value.phone,
      language: profileForm.value.language,
      theme: profileForm.value.theme
    })
    document.documentElement.setAttribute('data-theme', profileForm.value.theme)
    localStorage.setItem('theme', profileForm.value.theme)
    locale.value = profileForm.value.language
    localStorage.setItem('locale', profileForm.value.language)
    toast.add({ severity: 'success', summary: 'Profile updated', detail: 'Changes saved successfully.', life: 3000 })
  } catch (error) {
    showError('Unable to update profile', error)
  } finally {
    savingProfile.value = false
  }
}

const changePassword = async () => {
  if (passwordForm.value.new !== passwordForm.value.confirm) {
    toast.add({ severity: 'warn', summary: 'Passwords do not match', detail: 'Please confirm the same password.', life: 3000 })
    return
  }

  try {
    changingPassword.value = true
    await authStore.changePassword(passwordForm.value.current, passwordForm.value.new)
    toast.add({ severity: 'success', summary: 'Password updated', detail: 'Use the new password next time you log in.', life: 3000 })
    passwordForm.value = { current: '', new: '', confirm: '' }
  } catch (error) {
    showError('Unable to change password', error)
  } finally {
    changingPassword.value = false
  }
}

const signOutAll = async () => {
  await authStore.logout()
  toast.add({ severity: 'info', summary: 'Session closed', detail: 'Please sign in again to continue.', life: 3000 })
}

const handleSignOut = async () => {
  await authStore.logout()
}

const formatStatus = (value) => {
  if (value === 'active') return 'Active'
  if (value === 'inactive') return 'Inactive'
  if (value === 'suspended') return 'Suspended'
  return value || 'Unknown'
}

const formatRole = (value) => {
  switch (value) {
    case 'super_admin':
      return 'Super admin'
    case 'admin':
      return 'Admin'
    case 'employee':
      return 'Employee'
    case 'volunteer':
      return 'Volunteer'
    default:
      return 'User'
  }
}

const formatDate = (value) => {
  if (!value) return 'Not yet'
  return new Date(value).toLocaleString()
}

const showError = (summary, error) => {
  const detail = error?.response?.data?.error || error?.message || 'Unexpected error occurred'
  toast.add({ severity: 'error', summary, detail, life: 4000 })
}
</script>

<style scoped>
.profile-page {
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
  margin: 0.4rem 0 0;
  color: #6b7280;
}

.profile-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 1.5rem;
}

.profile-card :deep(.p-card-title) {
  font-size: 1.25rem;
}

.card-title {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.avatar {
  width: 72px;
  height: 72px;
  border-radius: 999px;
  background: #e5e7eb;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
  font-size: 1.5rem;
  color: #374151;
  overflow: hidden;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.avatar::after {
  content: attr(data-letter);
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

.form-group small {
  color: #9ca3af;
}

.full-width {
  grid-column: 1 / -1;
}

.form-actions {
  grid-column: 1 / -1;
  display: flex;
  justify-content: flex-end;
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
  color: #374151;
  font-size: 0.95rem;
}

.detail-list span {
  color: #9ca3af;
}

.tip-box {
  margin-top: 1rem;
  padding: 0.75rem 1rem;
  border-radius: 0.75rem;
  background: #f3f4f6;
  display: flex;
  gap: 0.5rem;
  align-items: center;
  color: #4b5563;
}

.security-hint {
  margin-top: 1rem;
}

.session-actions {
  margin-top: 1rem;
  display: flex;
  justify-content: flex-end;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .card-title {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
