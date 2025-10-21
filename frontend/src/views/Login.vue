<template>
  <div class="login">
    <h1>{{ $t('login.title') }}</h1>
    <form @submit.prevent="handleLogin">
      <div class="form-group">
        <label>{{ $t('login.username') }}</label>
        <input v-model="username" type="text" required />
      </div>
      <div class="form-group">
        <label>{{ $t('login.password') }}</label>
        <input v-model="password" type="password" required />
      </div>
      <button type="submit">{{ $t('login.submit') }}</button>
      <p v-if="error" class="error">{{ error }}</p>
    </form>
    <p>
      {{ $t('login.noAccount') }}
      <router-link to="/register">{{ $t('login.register') }}</router-link>
    </p>
  </div>
</template>

<script>
import { ref } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'Login',
  setup() {
    const store = useStore()
    const router = useRouter()
    const username = ref('')
    const password = ref('')
    const error = ref('')

    const handleLogin = async () => {
      try {
        error.value = ''
        await store.dispatch('auth/login', {
          username: username.value,
          password: password.value
        })
        router.push('/')
      } catch (err) {
        error.value = err.response?.data?.error || 'Login failed'
      }
    }

    return {
      username,
      password,
      error,
      handleLogin
    }
  }
}
</script>

<style scoped>
.login {
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
</style>
