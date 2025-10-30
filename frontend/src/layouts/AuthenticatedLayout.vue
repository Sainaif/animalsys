<template>
  <div class="authenticated-layout">
    <!-- Sidebar -->
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <RouterLink to="/app" class="logo">
          <img src="/logo.svg" alt="AnimalSys" class="logo-image" />
          <span v-if="!sidebarCollapsed" class="logo-text">AnimalSys</span>
        </RouterLink>
      </div>

      <nav class="sidebar-nav">
        <RouterLink to="/app" class="nav-item" exact>
          <span class="nav-icon">📊</span>
          <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.dashboard') }}</span>
        </RouterLink>

        <!-- Admin Section -->
        <template v-if="authStore.hasRole('admin')">
          <div class="nav-section">
            <span v-if="!sidebarCollapsed" class="nav-section-title">{{ t('nav.administration') }}</span>
          </div>
          <RouterLink to="/app/users" class="nav-item">
            <span class="nav-icon">👥</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.users') }}</span>
          </RouterLink>
        </template>

        <!-- Operations Section -->
        <template v-if="authStore.hasRole('employee')">
          <div class="nav-section">
            <span v-if="!sidebarCollapsed" class="nav-section-title">{{ t('nav.operations') }}</span>
          </div>
          <RouterLink to="/app/animals" class="nav-item">
            <span class="nav-icon">🐾</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.animals') }}</span>
          </RouterLink>
          <RouterLink to="/app/veterinary" class="nav-item">
            <span class="nav-icon">🏥</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.veterinary') }}</span>
          </RouterLink>
          <RouterLink to="/app/volunteers" class="nav-item">
            <span class="nav-icon">🤝</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.volunteers') }}</span>
          </RouterLink>
          <RouterLink to="/app/inventory" class="nav-item">
            <span class="nav-icon">📦</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.inventory') }}</span>
          </RouterLink>
        </template>

        <!-- Adoptions -->
        <div class="nav-section">
          <span v-if="!sidebarCollapsed" class="nav-section-title">{{ t('nav.adoptions') }}</span>
        </div>
        <RouterLink to="/app/adoptions" class="nav-item">
          <span class="nav-icon">❤️</span>
          <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.adoptions') }}</span>
        </RouterLink>

        <!-- Schedule (Volunteer+) -->
        <template v-if="authStore.hasRole('volunteer')">
          <RouterLink to="/app/schedules" class="nav-item">
            <span class="nav-icon">📅</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.schedules') }}</span>
          </RouterLink>
        </template>

        <!-- Management Section -->
        <template v-if="authStore.hasRole('employee')">
          <div class="nav-section">
            <span v-if="!sidebarCollapsed" class="nav-section-title">{{ t('nav.management') }}</span>
          </div>
          <RouterLink to="/app/finance" class="nav-item">
            <span class="nav-icon">💰</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.finance') }}</span>
          </RouterLink>
          <RouterLink to="/app/donors" class="nav-item">
            <span class="nav-icon">🎁</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.donors') }}</span>
          </RouterLink>
          <RouterLink to="/app/campaigns" class="nav-item">
            <span class="nav-icon">📢</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.campaigns') }}</span>
          </RouterLink>
          <RouterLink to="/app/partners" class="nav-item">
            <span class="nav-icon">🤝</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.partners') }}</span>
          </RouterLink>
          <RouterLink to="/app/communications" class="nav-item">
            <span class="nav-icon">✉️</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.communications') }}</span>
          </RouterLink>
          <RouterLink to="/app/documents" class="nav-item">
            <span class="nav-icon">📄</span>
            <span v-if="!sidebarCollapsed" class="nav-text">{{ t('nav.documents') }}</span>
          </RouterLink>
        </template>
      </nav>

      <!-- Collapse Toggle -->
      <button @click="toggleSidebar" class="sidebar-toggle">
        <span>{{ sidebarCollapsed ? '→' : '←' }}</span>
      </button>
    </aside>

    <!-- Main Area -->
    <div class="main-area">
      <!-- Header -->
      <header class="header">
        <div class="header-content">
          <h1 class="page-title">{{ pageTitle }}</h1>

          <div class="header-actions">
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

            <!-- User Menu -->
            <div class="user-menu" @click="toggleUserMenu">
              <div class="user-avatar">
                {{ userInitials }}
              </div>
              <span class="user-name">{{ userName }}</span>

              <!-- Dropdown -->
              <div v-if="userMenuOpen" class="user-dropdown">
                <RouterLink to="/app/profile" class="dropdown-item" @click="userMenuOpen = false">
                  {{ t('nav.profile') }}
                </RouterLink>
                <button @click="handleLogout" class="dropdown-item">
                  {{ t('auth.logout') }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </header>

      <!-- Content -->
      <main class="content">
        <RouterView />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { RouterLink, RouterView, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useThemeStore } from '../stores/theme'
import { useAuthStore } from '../stores/auth'
import { availableLanguages, languageInfo } from '../locales'

const { t, locale } = useI18n()
const route = useRoute()
const themeStore = useThemeStore()
const authStore = useAuthStore()

const sidebarCollapsed = ref(false)
const userMenuOpen = ref(false)

const currentLocale = computed({
  get: () => locale.value,
  set: (value) => {
    locale.value = value
    localStorage.setItem('locale', value)
  }
})

const pageTitle = computed(() => {
  // Extract title from route meta or name
  return route.meta.title || t(`nav.${route.name}`) || 'AnimalSys'
})

const userName = computed(() => {
  const user = authStore.user
  if (user?.first_name && user?.last_name) {
    return `${user.first_name} ${user.last_name}`
  }
  return user?.username || user?.email || 'User'
})

const userInitials = computed(() => {
  const user = authStore.user
  if (user?.first_name && user?.last_name) {
    return `${user.first_name[0]}${user.last_name[0]}`.toUpperCase()
  }
  return (user?.username?.[0] || 'U').toUpperCase()
})

function toggleSidebar() {
  sidebarCollapsed.value = !sidebarCollapsed.value
}

function toggleUserMenu() {
  userMenuOpen.value = !userMenuOpen.value
}

async function handleLogout() {
  userMenuOpen.value = false
  await authStore.logout()
}
</script>

<style scoped>
.authenticated-layout {
  display: flex;
  min-height: 100vh;
}

.sidebar {
  width: 250px;
  background-color: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  position: sticky;
  top: 0;
  height: 100vh;
  overflow-y: auto;
}

.sidebar.collapsed {
  width: 70px;
}

.sidebar-header {
  padding: 1rem;
  border-bottom: 1px solid var(--border-color);
}

.logo {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  text-decoration: none;
  color: var(--text-primary);
  font-weight: bold;
  font-size: 1.25rem;
}

.logo-image {
  width: 32px;
  height: 32px;
}

.sidebar-nav {
  flex: 1;
  padding: 1rem 0;
}

.nav-section {
  padding: 0.5rem 1rem;
  margin-top: 1rem;
}

.nav-section-title {
  font-size: 0.75rem;
  text-transform: uppercase;
  color: var(--text-tertiary);
  font-weight: 600;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.75rem 1rem;
  color: var(--text-secondary);
  text-decoration: none;
  transition: all 0.2s;
  border-left: 3px solid transparent;
}

.nav-item:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.nav-item.router-link-active {
  background-color: var(--bg-active);
  color: var(--primary-color);
  border-left-color: var(--primary-color);
}

.nav-icon {
  font-size: 1.25rem;
  min-width: 1.5rem;
}

.sidebar-toggle {
  padding: 0.75rem;
  border: none;
  border-top: 1px solid var(--border-color);
  background-color: var(--bg-secondary);
  cursor: pointer;
  color: var(--text-secondary);
  transition: all 0.2s;
}

.sidebar-toggle:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.main-area {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-width: 0;
}

.header {
  background-color: var(--bg-secondary);
  border-bottom: 1px solid var(--border-color);
  padding: 1rem 2rem;
  position: sticky;
  top: 0;
  z-index: 10;
}

.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.header-actions {
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

.user-menu {
  position: relative;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 0.25rem;
  transition: background-color 0.2s;
}

.user-menu:hover {
  background-color: var(--bg-hover);
}

.user-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  background-color: var(--primary-color);
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-weight: 600;
}

.user-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.5rem;
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 0.25rem;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  min-width: 200px;
  z-index: 100;
}

.dropdown-item {
  display: block;
  width: 100%;
  padding: 0.75rem 1rem;
  text-align: left;
  text-decoration: none;
  color: var(--text-secondary);
  background: none;
  border: none;
  cursor: pointer;
  transition: all 0.2s;
}

.dropdown-item:hover {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

.content {
  flex: 1;
  padding: 2rem;
  overflow-y: auto;
}

@media (max-width: 768px) {
  .sidebar {
    position: fixed;
    z-index: 100;
    transform: translateX(-100%);
  }

  .sidebar:not(.collapsed) {
    transform: translateX(0);
  }

  .content {
    padding: 1rem;
  }
}
</style>
