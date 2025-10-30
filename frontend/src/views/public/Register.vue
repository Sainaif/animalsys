<template>
  <div class="register-page">
    <div class="register-container">
      <BaseCard class="register-card">
        <h1 class="register-title">{{ t('auth.register') }}</h1>
        <p class="register-subtitle">{{ t('auth.registerSubtitle') }}</p>

        <form @submit.prevent="handleRegister" class="register-form">
          <FormGroup :label="t('auth.email')" input-id="email" :error="errors.email" required>
            <input id="email" v-model="form.email" type="email" required :disabled="loading" />
          </FormGroup>

          <FormGroup :label="t('auth.username')" input-id="username" :error="errors.username" required>
            <input id="username" v-model="form.username" type="text" required :disabled="loading" />
          </FormGroup>

          <div class="form-row">
            <FormGroup :label="t('common.firstName')" input-id="firstName" required>
              <input id="firstName" v-model="form.firstName" type="text" required :disabled="loading" />
            </FormGroup>

            <FormGroup :label="t('common.lastName')" input-id="lastName" required>
              <input id="lastName" v-model="form.lastName" type="text" required :disabled="loading" />
            </FormGroup>
          </div>

          <FormGroup :label="t('auth.password')" input-id="password" :error="errors.password" required>
            <input id="password" v-model="form.password" type="password" required :disabled="loading" />
          </FormGroup>

          <FormGroup :label="t('auth.confirmPassword')" input-id="confirmPassword" :error="errors.confirmPassword" required>
            <input id="confirmPassword" v-model="form.confirmPassword" type="password" required :disabled="loading" />
          </FormGroup>

          <BaseButton type="submit" variant="primary" full-width :loading="loading">
            {{ t('auth.registerButton') }}
          </BaseButton>
        </form>

        <div class="register-footer">
          <p>
            {{ t('auth.alreadyHaveAccount') }}
            <RouterLink to="/login">{{ t('auth.loginNow') }}</RouterLink>
          </p>
        </div>
      </BaseCard>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useRouter, RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../../stores/auth'
import { useNotificationStore } from '../../stores/notification'
import BaseCard from '../../components/base/BaseCard.vue'
import BaseButton from '../../components/base/BaseButton.vue'
import FormGroup from '../../components/base/FormGroup.vue'

const { t } = useI18n()
const router = useRouter()
const authStore = useAuthStore()
const notificationStore = useNotificationStore()

const form = reactive({
  email: '',
  username: '',
  firstName: '',
  lastName: '',
  password: '',
  confirmPassword: ''
})

const errors = reactive({
  email: '',
  username: '',
  password: '',
  confirmPassword: ''
})

const loading = ref(false)

function validateForm() {
  errors.email = ''
  errors.username = ''
  errors.password = ''
  errors.confirmPassword = ''

  let isValid = true

  if (form.password.length < 8) {
    errors.password = t('auth.passwordTooShort')
    isValid = false
  }

  if (form.password !== form.confirmPassword) {
    errors.confirmPassword = t('auth.passwordsDoNotMatch')
    isValid = false
  }

  return isValid
}

async function handleRegister() {
  if (!validateForm()) {
    return
  }

  loading.value = true

  try {
    await authStore.register({
      email: form.email,
      username: form.username,
      first_name: form.firstName,
      last_name: form.lastName,
      password: form.password,
      role: 'user'
    })

    notificationStore.success(t('auth.registerSuccess'))
    router.push({ name: 'login' })
  } catch (err) {
    const errorMsg = err.response?.data?.error || t('auth.registerError')
    notificationStore.error(errorMsg)
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.register-page {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 2rem;
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
}

.register-container {
  width: 100%;
  max-width: 500px;
}

.register-card {
  padding: 2rem;
}

.register-title {
  font-size: 2rem;
  font-weight: bold;
  text-align: center;
  margin-bottom: 0.5rem;
}

.register-subtitle {
  text-align: center;
  color: var(--text-secondary);
  margin-bottom: 2rem;
}

.register-form {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
}

.register-footer {
  margin-top: 1.5rem;
  text-align: center;
  color: var(--text-secondary);
}

.register-footer a {
  color: var(--primary-color);
  text-decoration: none;
  font-weight: 500;
}

.register-footer a:hover {
  text-decoration: underline;
}

@media (max-width: 640px) {
  .form-row {
    grid-template-columns: 1fr;
  }
}
</style>
