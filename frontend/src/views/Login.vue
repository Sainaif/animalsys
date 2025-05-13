<template>
  <div>
    <h1>{{$t('login')}}</h1>
    <form @submit.prevent="onLogin">
      <div>
        <label for="username">Nazwa użytkownika</label>
        <input v-model="username" id="username" required />
      </div>
      <div>
        <label for="password">Hasło</label>
        <input v-model="password" id="password" type="password" required />
      </div>
      <button type="submit">Zaloguj</button>
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
const password = ref('');
const error = computed(() => store.getters['auth/authError']);

const onLogin = async () => {
  const ok = await store.dispatch('auth/login', { username: username.value, password: password.value });
  if (ok) {
    router.push('/animals');
  }
};
</script>

<style>
.error { color: red; margin-top: 1em; }
form { max-width: 300px; margin: 2em auto; display: flex; flex-direction: column; gap: 1em; }
</style> 