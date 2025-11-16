<template>
  <div class="settings-page page-container">
    <div class="page-header">
      <div>
        <h1 class="page-title">
          {{ $t('settings.title') }}
        </h1>
        <p class="page-subtitle">
          {{ lastUpdatedLabel }}
        </p>
      </div>
      <Button
        label="Reload"
        icon="pi pi-refresh"
        class="p-button-text"
        :loading="loading"
        @click="loadSettings"
      />
    </div>

    <LoadingSpinner v-if="loading" />

    <form
      v-else
      class="settings-grid"
      @submit.prevent="saveSettings"
    >
      <Card class="settings-card">
        <template #title>
          {{ $t('settings.general.title') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-group">
              <label for="foundationName">{{ $t('settings.general.foundationName') }}</label>
              <InputText
                id="foundationName"
                v-model.trim="form.name"
                required
              />
            </div>
            <div class="form-group">
              <label for="legalName">{{ $t('settings.general.legalName') }}</label>
              <InputText
                id="legalName"
                v-model.trim="form.legal_name"
              />
            </div>
            <div class="form-group full-width">
              <label for="description">{{ $t('settings.general.description') }}</label>
              <Textarea
                id="description"
                v-model="form.description"
                rows="3"
                auto-resize
              />
            </div>
            <div class="form-group">
              <label for="contactEmail">{{ $t('settings.general.contactEmail') }}</label>
              <InputText
                id="contactEmail"
                v-model.trim="form.contact.email"
                type="email"
                required
              />
            </div>
            <div class="form-group">
              <label for="contactPhone">{{ $t('settings.general.contactPhone') }}</label>
              <InputText
                id="contactPhone"
                v-model.trim="form.contact.phone"
              />
            </div>
            <div class="form-group">
              <label for="website">{{ $t('settings.general.website') }}</label>
              <InputText
                id="website"
                v-model.trim="form.contact.website"
                placeholder="https://example.org"
              />
            </div>
          </div>
        </template>
      </Card>

      <Card class="settings-card">
        <template #title>
          {{ $t('settings.localization.title') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-group">
              <label for="language">{{ $t('settings.localization.defaultLanguage') }}</label>
              <Dropdown
                id="language"
                v-model="form.defaultLanguage"
                :options="languageOptions"
                option-label="label"
                option-value="value"
              />
            </div>
            <div class="form-group full-width">
              <label>{{ $t('settings.general.address') }}</label>
              <div class="address-grid">
                <InputText
                  v-model.trim="form.address.street"
                  placeholder="Street"
                />
                <InputText
                  v-model.trim="form.address.city"
                  placeholder="City"
                />
                <InputText
                  v-model.trim="form.address.state"
                  placeholder="State/Region"
                />
                <InputText
                  v-model.trim="form.address.zip_code"
                  placeholder="Postal code"
                />
                <InputText
                  v-model.trim="form.address.country"
                  placeholder="Country"
                />
              </div>
            </div>
          </div>
        </template>
      </Card>

      <Card class="settings-card">
        <template #title>
          {{ $t('settings.email.title') }}
        </template>
        <template #content>
          <div class="form-grid">
            <div class="form-group">
              <label for="fromName">{{ $t('settings.email.fromName') }}</label>
              <InputText
                id="fromName"
                v-model.trim="form.email.from_name"
                required
              />
            </div>
            <div class="form-group">
              <label for="fromEmail">{{ $t('settings.email.fromEmail') }}</label>
              <InputText
                id="fromEmail"
                v-model.trim="form.email.from_email"
                type="email"
                required
              />
            </div>
            <div class="form-group">
              <label for="smtpServer">{{ $t('settings.email.smtpServer') }}</label>
              <InputText
                id="smtpServer"
                v-model.trim="form.email.smtp_host"
              />
            </div>
            <div class="form-group">
              <label for="smtpPort">{{ $t('settings.email.smtpPort') }}</label>
              <InputText
                id="smtpPort"
                v-model.number="form.email.smtp_port"
                type="number"
                min="1"
              />
            </div>
            <div class="form-group">
              <label for="smtpUsername">{{ $t('settings.email.smtpUsername') }}</label>
              <InputText
                id="smtpUsername"
                v-model.trim="form.email.smtp_username"
              />
            </div>
            <div class="form-group">
              <label for="smtpPassword">{{ $t('settings.email.smtpPassword') }}</label>
              <InputText
                id="smtpPassword"
                v-model="form.email.smtp_password"
                type="password"
                autocomplete="new-password"
              />
              <small>{{ $t('settings.email.smtpPasswordHint') }}</small>
            </div>
            <div class="form-group toggle">
              <Checkbox
                v-model="form.email.enable_tls"
                input-id="tlsToggle"
                :binary="true"
              />
              <label for="tlsToggle">{{ $t('settings.email.enableTls') }}</label>
            </div>
          </div>
        </template>
      </Card>

      <div class="form-actions">
        <Button
          type="button"
          label="Cancel"
          class="p-button-text"
          :disabled="saving"
          @click="loadSettings"
        />
        <Button
          type="submit"
          :label="$t('common.save')"
          icon="pi pi-check"
          :loading="saving"
        />
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref, computed, onMounted } from 'vue'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'
import Card from 'primevue/card'
import InputText from 'primevue/inputtext'
import Textarea from 'primevue/textarea'
import Dropdown from 'primevue/dropdown'
import Checkbox from 'primevue/checkbox'
import Button from 'primevue/button'
import LoadingSpinner from '@/components/shared/LoadingSpinner.vue'
import { settingsService } from '@/services/settingsService'

const toast = useToast()
const { t, locale } = useI18n()

const loading = ref(true)
const saving = ref(false)
const metadata = ref({ updated_at: null })
let emptySettingsNoticeShown = false

const form = reactive({
  name: '',
  legal_name: '',
  description: '',
  contact: {
    email: '',
    phone: '',
    website: ''
  },
  address: {
    street: '',
    city: '',
    state: '',
    zip_code: '',
    country: ''
  },
  defaultLanguage: localStorage.getItem('locale') || locale.value || 'pl',
  email: {
    smtp_host: '',
    smtp_port: 587,
    smtp_username: '',
    smtp_password: '',
    from_email: '',
    from_name: '',
    enable_tls: true
  }
})

const languageOptions = [
  { label: 'English', value: 'en' },
  { label: 'Polski', value: 'pl' }
]

const lastUpdatedLabel = computed(() => {
  if (!metadata.value.updated_at) {
    return t('settings.general.neverUpdated')
  }
  return `${t('settings.general.lastUpdated')} ${new Date(metadata.value.updated_at).toLocaleString()}`
})

const optional = (value) => {
  if (value === null || value === undefined) return undefined
  const trimmed = value.toString().trim()
  return trimmed === '' ? undefined : trimmed
}

const applySettings = (settings) => {
  form.name = settings?.name || ''
  form.legal_name = settings?.legal_name || ''
  form.description = settings?.description || ''
  form.contact.email = settings?.contact_info?.email || ''
  form.contact.phone = settings?.contact_info?.phone || ''
  form.contact.website = settings?.contact_info?.website || ''
  form.address.street = settings?.address?.street || ''
  form.address.city = settings?.address?.city || ''
  form.address.state = settings?.address?.state || ''
  form.address.zip_code = settings?.address?.zip_code || ''
  form.address.country = settings?.address?.country || ''
  form.email.smtp_host = settings?.email_settings?.smtp_host || ''
  form.email.smtp_port = settings?.email_settings?.smtp_port ?? 587
  form.email.smtp_username = settings?.email_settings?.smtp_username || ''
  form.email.smtp_password = ''
  form.email.from_email = settings?.email_settings?.from_email || ''
  form.email.from_name = settings?.email_settings?.from_name || ''
  form.email.enable_tls = settings?.email_settings?.enable_tls ?? true
  metadata.value.updated_at = settings?.updated_at || null
}

const loadSettings = async () => {
  try {
    loading.value = true
    const data = await settingsService.getSettings()
    applySettings(data)
  } catch (error) {
    if (error?.response?.status === 404) {
      applySettings({})
      metadata.value.updated_at = null
      if (!emptySettingsNoticeShown) {
        toast.add({
          severity: 'info',
          summary: t('settings.general.noSettingsConfigured'),
          detail: t('settings.general.noSettingsConfiguredDescription'),
          life: 4000
        })
        emptySettingsNoticeShown = true
      }
    } else {
      showError(t('settings.general.loadError'), error)
    }
  } finally {
    loading.value = false
  }
}

const buildOrganizationPayload = () => {
  const hasAddress = Object.values(form.address).some((value) => optional(value))
  return {
    name: form.name,
    legal_name: optional(form.legal_name),
    description: optional(form.description),
    contact_info: {
      email: form.contact.email,
      phone: optional(form.contact.phone),
      website: optional(form.contact.website)
    },
    address: hasAddress
      ? {
          street: optional(form.address.street),
          city: optional(form.address.city),
          state: optional(form.address.state),
          zip_code: optional(form.address.zip_code),
          country: optional(form.address.country)
        }
      : null
  }
}

const buildEmailPayload = () => {
  const payload = {
    smtp_host: optional(form.email.smtp_host),
    smtp_port: form.email.smtp_port || 587,
    smtp_username: optional(form.email.smtp_username),
    smtp_password: optional(form.email.smtp_password),
    from_email: form.email.from_email,
    from_name: form.email.from_name,
    enable_tls: form.email.enable_tls
  }

  if (!payload.smtp_password) {
    delete payload.smtp_password
  }

  return payload
}

const saveSettings = async () => {
  if (!form.name.trim()) {
    toast.add({ severity: 'warn', summary: t('settings.general.validation'), detail: t('settings.general.foundationNameRequired'), life: 3000 })
    return
  }

  if (!form.contact.email.trim() || !form.email.from_email.trim()) {
    toast.add({ severity: 'warn', summary: t('settings.general.validation'), detail: t('settings.general.emailRequired'), life: 3000 })
    return
  }

  try {
    saving.value = true
    const [organization] = await Promise.all([
      settingsService.updateOrganization(buildOrganizationPayload()),
      settingsService.updateEmailSettings(buildEmailPayload())
    ])
    form.name = organization.name
    form.legal_name = organization.legal_name || ''
    form.description = organization.description || ''
    await loadSettings()
    locale.value = form.defaultLanguage
    localStorage.setItem('locale', form.defaultLanguage)
    toast.add({ severity: 'success', summary: t('settings.saveSuccess'), life: 3000 })
  } catch (error) {
    showError(t('settings.general.saveError'), error)
  } finally {
    saving.value = false
  }
}

const showError = (summary, error) => {
  const detail = error?.response?.data?.error || error?.message || t('common.genericError')
  toast.add({ severity: 'error', summary, detail, life: 4000 })
}

onMounted(() => {
  loadSettings()
})
</script>

<style scoped>
.settings-page {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.page-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  margin: 0;
  font-size: 2rem;
}

.page-subtitle {
  margin: 0.2rem 0 0;
  color: #6b7280;
}

.settings-grid {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.settings-card :deep(.p-card-title) {
  font-size: 1.2rem;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
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

.address-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 0.5rem;
}

.toggle {
  flex-direction: row;
  align-items: center;
  gap: 0.5rem;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
}

@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 0.75rem;
  }
}
</style>
