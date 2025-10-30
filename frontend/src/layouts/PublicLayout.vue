<template>
  <div class="public-layout">
    <!-- Header -->
    <header class="header">
      <div class="container">
        <nav class="nav">
          <RouterLink to="/" class="logo">
            <img src="/logo.svg" alt="AnimalSys" class="logo-image" />
            <span class="logo-text">AnimalSys</span>
          </RouterLink>

          <div class="nav-links">
            <RouterLink to="/" class="nav-link">{{ t('nav.home') }}</RouterLink>
            <RouterLink to="/animals" class="nav-link">{{ t('nav.animals') }}</RouterLink>
            <RouterLink to="/campaigns" class="nav-link">{{ t('nav.campaigns') }}</RouterLink>
          </div>

          <div class="nav-actions">
            <!-- Language Selector -->
            <select v-model="currentLocale" class="locale-selector">
              <option v-for="lang in availableLanguages" :key="lang" :value="lang">
                {{ languageInfo[lang].nativeName }}
              </option>
            </select>

            <!-- Theme Toggle -->
            <button @click="themeStore.toggleTheme()" class="theme-toggle" :title="t('common.toggleTheme')">
              <span v-if="themeStore.isDark">☀️</span>
              <span v-else>🌙</span>
            </button>

            <!-- Auth Links -->
            <template v-if="!authStore.isAuthenticated">
              <RouterLink to="/login" class="btn btn-outline">{{ t('auth.login') }}</RouterLink>
              <RouterLink to="/register" class="btn btn-primary">{{ t('auth.register') }}</RouterLink>
            </template>
            <template v-else>
              <RouterLink to="/app" class="btn btn-primary">{{ t('nav.dashboard') }}</RouterLink>
            </template>
          </div>
        </nav>
      </div>
    </header>

    <!-- Main Content -->
    <main class="main-content">
      <RouterView />
    </main>

    <!-- Footer -->
    <footer class="footer">
      <div class="container">
        <p>&copy; {{ currentYear }} AnimalSys. {{ t('common.allRightsReserved') }}</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink, RouterView } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useThemeStore } from '../stores/theme'
import { useAuthStore } from '../stores/auth'
import { availableLanguages, languageInfo } from '../locales'

const { t, locale } = useI18n()
const themeStore = useThemeStore()
const authStore = useAuthStore()

const currentLocale = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value
    localStorage.setItem('locale', value)
  }
})

const currentYear = new Date().getFullYear()
</script>

<style scoped>
.public-layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.header {
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  padding: 1rem 0;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 1rem;
}

.nav {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 2rem;
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
  font-size: 1.5rem;
  font-weight: bold;
  color: var(--text-primary);
}

.logo-image {
  width: 40px;
  height: 40px;
}

.nav-links {
  display: flex;
  gap: 2rem;
  flex: 1;
}

.nav-link {
  text-decoration: none;
  color: var(--text-secondary);
  transition: color 0.2s;
}

.nav-link:hover,
.nav-link.router-link-active {
  color: var(--primary-color);
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.locale-selector {
  padding: 0.5rem;
  border: 1px solid var(--border-color);
  border-radius: 0.25rem;
  background-color: var(--bg-primary);
  color: var(--text-primary);
}

.theme-toggle {
  padding: 0.5rem;
  border: none;
  background: transparent;
  font-size: 1.25rem;
  cursor: pointer;
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 0.25rem;
  text-decoration: none;
  font-weight: 500;
  transition: all 0.2s;
  display: inline-block;
}

.btn-outline {
  border: 1px solid var(--primary-color);
  color: var(--primary-color);
  background: transparent;
}

.btn-outline:hover {
  background-color: var(--primary-color);
  color: white;
}

.btn-primary {
  background-color: var(--primary-color);
  color: white;
  border: none;
}

.btn-primary:hover {
  opacity: 0.9;
}

.main-content {
  flex: 1;
}

.footer {
  background-color: var(--bg-secondary);
  border-top: 1px solid var(--border-color);
  padding: 2rem 0;
  text-align: center;
  color: var(--text-secondary);
}

@media (max-width: 768px) {
  .nav {
    flex-direction: column;
    gap: 1rem;
  }

  .nav-links {
    flex-direction: column;
    gap: 0.5rem;
    width: 100%;
  }

  .nav-actions {
    width: 100%;
    justify-content: center;
  }
}
</style>
