<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import Password from 'primevue/password'
import { useToast } from 'primevue/usetoast'
import { useI18n } from 'vue-i18n'

const router = useRouter()
const authStore = useAuthStore()
const toast = useToast()
const { t } = useI18n()

const email = ref('')
const password = ref('')
const loading = ref(false)

const handleLogin = async () => {
  loading.value = true
  try {
    await authStore.login(email.value, password.value)
    toast.add({ severity: 'success', summary: t('common.success'), detail: t('auth.loginSuccess'), life: 3000 })
    router.push({ name: 'dashboard' })
  } catch (error) {
    toast.add({
      severity: 'error',
      summary: t('common.error'),
      detail: error.response?.data?.message || t('auth.loginFailed'),
      life: 5000
    })
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="login-container">
    <div class="login-card card">
      <div class="login-header">
        <h1 class="login-title">
          {{ $t('auth.loginTitle') }}
        </h1>
        <p class="login-subtitle">
          {{ $t('auth.loginSubtitle') }}
        </p>
      </div>

      <form
        class="login-form"
        @submit.prevent="handleLogin"
      >
        <div class="field">
          <label for="email">{{ $t('auth.email') }}</label>
          <InputText
            id="email"
            v-model="email"
            type="email"
            :placeholder="$t('auth.emailPlaceholder')"
            required
            class="w-full"
          />
        </div>

        <div class="field">
          <label for="password">{{ $t('auth.password') }}</label>
          <Password
            id="password"
            v-model="password"
            :placeholder="$t('auth.passwordPlaceholder')"
            :feedback="false"
            toggle-mask
            required
            class="w-full"
          />
        </div>

        <Button
          type="submit"
          :label="$t('auth.login')"
          icon="pi pi-sign-in"
          :loading="loading"
          class="w-full"
        />
      </form>
    </div>
  </div>
</template>

<style scoped>
.login-container {
  display: flex;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.9) 0%, rgba(118, 75, 162, 0.9) 100%);
  padding: 1rem;
}

.login-card {
  width: 100%;
  max-width: 420px;
  background: var(--card-bg, #ffffff);
  border-radius: 1.25rem;
  border: 1px solid var(--border-color, rgba(255, 255, 255, 0.15));
  box-shadow: 0 25px 60px rgba(15, 23, 42, 0.25);
}

.login-header {
  text-align: center;
  margin-bottom: 2rem;
}

.login-title {
  font-size: 1.8rem;
  font-weight: 700;
  color: var(--text-color, #0f172a);
  margin-bottom: 0.5rem;
}

.login-subtitle {
  color: var(--text-muted, #6c757d);
  font-size: 0.95rem;
}

.login-form {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.field label {
  font-weight: 600;
  color: var(--text-muted, #495057);
  font-size: 0.95rem;
}

:global([data-theme='dark'] .login-container) {
  background: radial-gradient(circle at top, rgba(8, 18, 43, 0.95) 0%, rgba(5, 8, 22, 0.98) 100%);
}

:global([data-theme='dark'] .login-card) {
  box-shadow: 0 30px 80px rgba(3, 7, 18, 0.8);
}
</style>
