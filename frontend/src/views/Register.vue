<template>
  <div class="register">
    <h1>{{ $t('register.title') }}</h1>
    <form @submit.prevent="handleRegister">
      <div class="form-group">
        <label>{{ $t('register.username') }}</label>
        <input v-model="username" type="text" required />
      </div>
      <div class="form-group">
        <label>{{ $t('register.email') }}</label>
        <input v-model="email" type="email" required />
      </div>
      <div class="form-group">
        <label>{{ $t('register.password') }}</label>
        <input v-model="password" type="password" required />
      </div>
      <button type="submit">{{ $t('register.submit') }}</button>
      <p v-if="error" class="error">{{ error }}</p>
      <p v-if="success" class="success">{{ success }}</p>
    </form>
    <p>
      {{ $t('register.hasAccount') }}
      <router-link to="/login">{{ $t('register.login') }}</router-link>
    </p>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'Register',
  setup() {
    const store = useStore()
    const router = useRouter()
    const username = ref('')
    const email = ref('')
    const password = ref('')
    const error = ref('')
    const success = ref('')

    const handleRegister = async () => {
      try {
        error.value = ''
        success.value = ''
        await store.dispatch('auth/register', {
          username: username.value,
          email: email.value,
          password: password.value
        })
        success.value = 'Registration successful! Redirecting to login...'
        setTimeout(() => {
          router.push('/login')
        }, 2000)
      } catch (err) {
        error.value = err.response?.data?.error || 'Registration failed'
      }
    }

    return {
      username,
      email,
      password,
      error,
      success,
      handleRegister
    }
  }
}
</script>

<style scoped>
.register {
  max-width: 400px;
  margin: 50px auto;
  padding: 20px;
}

.form-group {
  margin-bottom: 15px;
}

.form-group label {
  display: block;
  margin-bottom: 5px;
}

.form-group input {
  width: 100%;
  padding: 8px;
  box-sizing: border-box;
}

button {
  width: 100%;
  padding: 10px;
  background-color: #42b983;
  color: white;
  border: none;
  cursor: pointer;
}

.error {
  color: red;
  margin-top: 10px;
}

.success {
  color: green;
  margin-top: 10px;
}
</style>
