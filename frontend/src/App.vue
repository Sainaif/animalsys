<template>
  <div>
    <nav class="navbar">
      <router-link to="/animals">{{$t('animals')}}</router-link>
      <router-link to="/adoptions">{{$t('adoptions')}}</router-link>
      <router-link to="/schedule">{{$t('schedule')}}</router-link>
      <router-link to="/documents">{{$t('documents')}}</router-link>
      <router-link v-if="isAdmin" to="/users">{{$t('users')}}</router-link>
      <router-link v-if="!isAuthenticated" to="/login">{{$t('login')}}</router-link>
      <router-link v-if="!isAuthenticated" to="/register">{{$t('register')}}</router-link>
      <button v-if="isAuthenticated" @click="logout">Wyloguj</button>
      <select v-model="$i18n.locale">
        <option value="pl">PL</option>
        <option value="en">EN</option>
      </select>
    </nav>
    <main>
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useStore } from 'vuex';
import { useRouter } from 'vue-router';

const store = useStore();
const router = useRouter();
const isAuthenticated = computed(() => store.getters['auth/isAuthenticated']);
const isAdmin = computed(() => store.getters['auth/userRole'] === 'admin');
const logout = () => {
  store.dispatch('auth/logout');
  router.push('/login');
};
</script>

<style>
.navbar {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  background: #f8f8f8;
  padding: 1rem;
  align-items: center;
}
main {
  padding: 1rem;
}
</style> 