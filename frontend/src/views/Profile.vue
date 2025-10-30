<template>
  <div class="profile-page">
    <h1 class="page-title">{{ t('nav.profile') }}</h1>

    <div class="profile-content">
      <!-- Profile Info Card -->
      <BaseCard>
        <template #header>{{ t('profile.personalInfo') }}</template>

        <form @submit.prevent="updateProfile" class="profile-form">
          <div class="form-row">
            <FormGroup
              :label="t('common.firstName')"
              :error="errors.first_name"
              required
            >
              <input
                v-model="profileForm.first_name"
                type="text"
                :placeholder="t('common.firstName')"
                required
              />
            </FormGroup>

            <FormGroup
              :label="t('common.lastName')"
              :error="errors.last_name"
              required
            >
              <input
                v-model="profileForm.last_name"
                type="text"
                :placeholder="t('common.lastName')"
                required
              />
            </FormGroup>
          </div>

          <FormGroup
            :label="t('common.email')"
            :error="errors.email"
            required
          >
            <input
              v-model="profileForm.email"
              type="email"
              :placeholder="t('common.email')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('common.phone')"
            :error="errors.phone"
          >
            <input
              v-model="profileForm.phone"
              type="tel"
              :placeholder="t('common.phone')"
            />
          </FormGroup>

          <div class="form-actions">
            <BaseButton variant="primary" type="submit" :loading="updatingProfile">
              {{ t('profile.updateProfile') }}
            </BaseButton>
          </div>
        </form>
      </BaseCard>

      <!-- Account Info Card -->
      <BaseCard>
        <template #header>{{ t('profile.accountInfo') }}</template>

        <div class="info-grid">
          <div class="info-item">
            <span class="info-label">{{ t('profile.username') }}:</span>
            <span class="info-value">{{ user?.username || '-' }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('profile.role') }}:</span>
            <span class="info-value">{{ t(`profile.roles.${user?.role}`) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('profile.createdAt') }}:</span>
            <span class="info-value">{{ formatDate(user?.created_at) }}</span>
          </div>
          <div class="info-item">
            <span class="info-label">{{ t('profile.lastLogin') }}:</span>
            <span class="info-value">{{ formatDate(user?.last_login) }}</span>
          </div>
        </div>
      </BaseCard>

      <!-- Change Password Card -->
      <BaseCard>
        <template #header>{{ t('profile.changePassword') }}</template>

        <form @submit.prevent="changePassword" class="password-form">
          <FormGroup
            :label="t('profile.currentPassword')"
            :error="passwordErrors.current_password"
            required
          >
            <input
              v-model="passwordForm.current_password"
              type="password"
              :placeholder="t('profile.currentPassword')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('profile.newPassword')"
            :error="passwordErrors.new_password"
            required
          >
            <input
              v-model="passwordForm.new_password"
              type="password"
              :placeholder="t('profile.newPassword')"
              required
            />
          </FormGroup>

          <FormGroup
            :label="t('profile.confirmNewPassword')"
            :error="passwordErrors.confirm_password"
            required
          >
            <input
              v-model="passwordForm.confirm_password"
              type="password"
              :placeholder="t('profile.confirmNewPassword')"
              required
            />
          </FormGroup>

          <div class="form-actions">
            <BaseButton variant="primary" type="submit" :loading="changingPassword">
              {{ t('profile.changePassword') }}
            </BaseButton>
          </div>
        </form>
      </BaseCard>

      <!-- Preferences Card -->
      <BaseCard>
        <template #header>{{ t('profile.preferences') }}</template>

        <div class="preferences-section">
          <div class="preference-item">
            <div class="preference-info">
              <span class="preference-label">{{ t('profile.language') }}</span>
              <span class="preference-description">{{ t('profile.selectLanguage') }}</span>
            </div>
            <select v-model="selectedLocale" @change="changeLanguage" class="preference-select">
              <option value="pl">Polski</option>
              <option value="en">English</option>
            </select>
          </div>

          <div class="preference-item">
            <div class="preference-info">
              <span class="preference-label">{{ t('profile.theme') }}</span>
              <span class="preference-description">{{ t('profile.selectTheme') }}</span>
            </div>
            <select v-model="themeStore.theme" @change="themeStore.setTheme(themeStore.theme)" class="preference-select">
              <option value="light">{{ t('profile.lightTheme') }}</option>
              <option value="dark">{{ t('profile.darkTheme') }}</option>
            </select>
          </div>
        </div>
      </BaseCard>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'
import { useThemeStore } from '../stores/theme'
import { useNotificationStore } from '../stores/notification'
import { API } from '../api'
import BaseCard from '../components/base/BaseCard.vue'
import BaseButton from '../components/base/BaseButton.vue'
import FormGroup from '../components/base/FormGroup.vue'

const { t, locale } = useI18n()
const authStore = useAuthStore()
const themeStore = useThemeStore()
const notificationStore = useNotificationStore()

const user = computed(() => authStore.user)
const selectedLocale = ref(locale.value)

const profileForm = reactive({
  first_name: '',
  last_name: '',
  email: '',
  phone: ''
})

const passwordForm = reactive({
  current_password: '',
  new_password: '',
  confirm_password: ''
})

const errors = ref({})
const passwordErrors = ref({})
const updatingProfile = ref(false)
const changingPassword = ref(false)

onMounted(() => {
  loadProfile()
})

function loadProfile() {
  if (user.value) {
    profileForm.first_name = user.value.first_name || ''
    profileForm.last_name = user.value.last_name || ''
    profileForm.email = user.value.email || ''
    profileForm.phone = user.value.phone || ''
  }
}

async function updateProfile() {
  errors.value = {}

  if (!profileForm.first_name.trim()) {
    errors.value.first_name = t('common.required')
    return
  }

  if (!profileForm.last_name.trim()) {
    errors.value.last_name = t('common.required')
    return
  }

  if (!profileForm.email.trim() || !profileForm.email.includes('@')) {
    errors.value.email = t('common.invalid')
    return
  }

  updatingProfile.value = true

  try {
    await API.auth.updateProfile(profileForm)
    await authStore.fetchProfile()
    loadProfile() // Reload form with updated data
    notificationStore.success(t('profile.updateSuccess'))
  } catch (error) {
    notificationStore.error(t('common.error'), error.response?.data?.error || error.message)
  } finally {
    updatingProfile.value = false
  }
}

async function changePassword() {
  passwordErrors.value = {}

  if (!passwordForm.current_password) {
    passwordErrors.value.current_password = t('common.required')
    return
  }

  if (!passwordForm.new_password || passwordForm.new_password.length < 8) {
    passwordErrors.value.new_password = t('auth.passwordTooShort')
    return
  }

  if (passwordForm.new_password !== passwordForm.confirm_password) {
    passwordErrors.value.confirm_password = t('auth.passwordsDoNotMatch')
    return
  }

  changingPassword.value = true

  try {
    await API.auth.changePassword({
      old_password: passwordForm.current_password,
      new_password: passwordForm.new_password
    })

    // Clear form
    passwordForm.current_password = ''
    passwordForm.new_password = ''
    passwordForm.confirm_password = ''

    notificationStore.success(t('profile.passwordChangeSuccess'))
  } catch (error) {
    notificationStore.error(t('common.error'), error.response?.data?.error || error.message)
  } finally {
    changingPassword.value = false
  }
}

function changeLanguage() {
  locale.value = selectedLocale.value
  localStorage.setItem('locale', selectedLocale.value)
  notificationStore.success(t('profile.languageChanged'))
}

function formatDate(dateString) {
  if (!dateString) return '-'
  return new Date(dateString).toLocaleString()
}
</script>

<style scoped>
.profile-page {
  padding: 2rem;
  max-width: 1000px;
  margin: 0 auto;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  color: var(--text-primary);
  margin: 0 0 2rem 0;
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.profile-form,
.password-form {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

input,
select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
}

input:focus,
select:focus {
  outline: none;
  border-color: var(--primary-color);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  margin-top: 0.5rem;
}

.info-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 1.5rem;
}

.info-item {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.info-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.info-value {
  font-size: 1rem;
  color: var(--text-primary);
  font-weight: 600;
}

.preferences-section {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.preference-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
}

.preference-info {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.preference-label {
  font-size: 1rem;
  color: var(--text-primary);
  font-weight: 600;
}

.preference-description {
  font-size: 0.875rem;
  color: var(--text-secondary);
}

.preference-select {
  width: auto;
  min-width: 150px;
}

@media (max-width: 768px) {
  .profile-page {
    padding: 1rem;
  }

  .form-row,
  .info-grid {
    grid-template-columns: 1fr;
  }

  .preference-item {
    flex-direction: column;
    align-items: flex-start;
    gap: 1rem;
  }

  .preference-select {
    width: 100%;
  }
}
</style>
