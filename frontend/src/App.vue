<template>
  <div id="app">
    <nav v-if="isAuthenticated">
      <router-link to="/">{{ $t('nav.home') }}</router-link>
      <router-link to="/animals">{{ $t('nav.animals') }}</router-link>
      <router-link to="/adoptions">{{ $t('nav.adoptions') }}</router-link>
      <router-link to="/schedules">{{ $t('nav.schedules') }}</router-link>
      <router-link to="/documents">{{ $t('nav.documents') }}</router-link>
      <router-link to="/finances" v-if="['admin', 'employee'].includes(user?.role)">{{ $t('nav.finances') }}</router-link>
      <router-link to="/users" v-if="user?.role === 'admin'">{{ $t('nav.users') }}</router-link>
      <button @click="logout">{{ $t('nav.logout') }}</button>
    </nav>
    <router-view />
  </div>
</template>

<script>
import { computed } from 'vue'
import { useStore } from 'vuex'
import { useRouter } from 'vue-router'

export default {
  name: 'App',
  setup() {
    const store = useStore()
    const router = useRouter()

    const isAuthenticated = computed(() => store.getters['auth/isAuthenticated'])
    const user = computed(() => store.getters['auth/user'])

    const logout = () => {
      store.dispatch('auth/logout')
      router.push('/login')
    }

    return {
      isAuthenticated,
      user,
      logout
    }
  }
}
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
  color: #2c3e50;
}

nav {
  padding: 30px;
  background-color: #f3f3f3;
}

nav a {
  font-weight: bold;
  color: #2c3e50;
  margin-right: 20px;
  text-decoration: none;
}

nav a.router-link-exact-active {
  color: #42b983;
}

nav button {
  margin-left: 20px;
  padding: 5px 15px;
  cursor: pointer;
}
</style>
