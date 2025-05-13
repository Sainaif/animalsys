<template>
  <div>
    <h1>{{$t('register')}}</h1>
    <form @submit.prevent="onRegister">
      <div>
        <label for="username">Nazwa użytkownika</label>
        <input v-model="username" id="username" required />
      </div>
      <div>
        <label for="email">Email</label>
        <input v-model="email" id="email" type="email" required />
      </div>
      <div>
        <label for="password">Hasło</label>
        <input v-model="password" id="password" type="password" required />
      </div>
      <button type="submit">Zarejestruj</button>
      <div v-if="error" class="error">{{ error }}</div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

const store = useStore();
const router = useRouter();
const username = ref('');
const email = ref('');
const password = ref('');
const error = computed(() => store.getters['auth/authError']);

const onRegister = async () => {
  const ok = await store.dispatch('auth/register', { username: username.value, email: email.value, password: password.value });
  if (ok) {
    router.push('/animals');
  }
};
</script>

<style>
.error { color: red; margin-top: 1em; }
form { max-width: 300px; margin: 2em auto; display: flex; flex-direction: column; gap: 1em; }
</style> 