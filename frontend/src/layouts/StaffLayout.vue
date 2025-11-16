<template>
  <div class="staff-layout">
    <aside
      class="sidebar"
      :class="{ collapsed: sidebarCollapsed }"
    >
      <div class="sidebar-header">
        <router-link
          to="/dashboard"
          class="brand"
        >
          <i class="pi pi-heart" />
          <span v-if="!sidebarCollapsed">Animal Foundation</span>
        </router-link>
        <button
          class="collapse-btn"
          @click="toggleSidebar"
        >
          <i
            class="pi"
            :class="sidebarCollapsed ? 'pi-angle-right' : 'pi-angle-left'"
          />
        </button>
      </div>

      <nav class="sidebar-nav">
        <router-link
          to="/dashboard"
          class="nav-item"
        >
          <i class="pi pi-th-large" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.dashboard') }}</span>
        </router-link>
        <router-link
          to="/staff/animals"
          class="nav-item"
        >
          <i class="pi pi-heart" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.animals') }}</span>
        </router-link>
        <router-link
          to="/staff/veterinary"
          class="nav-item"
        >
          <i class="pi pi-plus" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.veterinary') }}</span>
        </router-link>
        <router-link
          to="/staff/adoptions/applications"
          class="nav-item"
        >
          <i class="pi pi-users" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.adoptions') }}</span>
        </router-link>
        <router-link
          to="/staff/finance"
          class="nav-item"
        >
          <i class="pi pi-money-bill" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.finance') }}</span>
        </router-link>
        <router-link
          to="/staff/events"
          class="nav-item"
        >
          <i class="pi pi-calendar" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.events') }}</span>
        </router-link>
        <router-link
          to="/contacts"
          class="nav-item"
        >
          <i class="pi pi-phone" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.contacts') }}</span>
        </router-link>

        <Divider v-if="!sidebarCollapsed" />

        <router-link
          v-if="isAdmin"
          to="/users"
          class="nav-item"
        >
          <i class="pi pi-user" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.users') }}</span>
        </router-link>
        <router-link
          v-if="isAdmin"
          to="/settings"
          class="nav-item"
        >
          <i class="pi pi-cog" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.settings') }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <router-link
          to="/profile"
          class="nav-item"
        >
          <i class="pi pi-user" />
          <span v-if="!sidebarCollapsed">{{ $t('nav.profile') }}</span>
        </router-link>
      </div>
    </aside>

    <div class="main-wrapper">
      <header class="topbar">
        <div class="topbar-left">
          <h2 class="page-title">
            {{ pageTitle }}
          </h2>
        </div>
        <div class="topbar-right">
          <Dropdown
            v-model="currentLocale"
            :options="locales"
            option-label="label"
            option-value="value"
            class="locale-dropdown"
            @change="changeLocale"
          >
            <template #value="slotProps">
              <span class="locale-flag">{{ getLocaleFlag(slotProps.value) }}</span>
            </template>
            <template #option="slotProps">
              <span class="locale-flag">{{ getLocaleFlag(slotProps.option.value) }}</span>
              <span>{{ slotProps.option.label }}</span>
            </template>
          </Dropdown>

          <Button
            :icon="theme === 'light' ? 'pi pi-moon' : 'pi pi-sun'"
            class="p-button-text p-button-rounded theme-toggle"
            :aria-label="theme === 'light' ? $t('common.enableDarkMode') : $t('common.enableLightMode')"
            @click="toggleTheme"
          />

          <Button
            v-tooltip.bottom="$t('nav.home')"
            icon="pi pi-home"
            class="p-button-text p-button-rounded"
            :aria-label="$t('nav.home')"
            @click="goToHome"
          />

          <Button
            v-tooltip.bottom="$t('nav.notifications')"
            v-badge="3"
            icon="pi pi-bell"
            class="p-button-text p-button-rounded"
            :aria-label="$t('nav.notifications')"
          />

          <div class="user-menu">
            <Button
              :label="userDisplayName"
              icon="pi pi-user"
              class="p-button-text"
              @click="toggleUserMenu"
            />
            <Menu
              ref="userMenuRef"
              :model="userMenuItems"
              :popup="true"
            />
          </div>
        </div>
      </header>

      <main class="content">
        <Toast />
        <ConfirmDialog />
        <router-view />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import Menu from 'primevue/menu'
import Divider from 'primevue/divider'
import Toast from 'primevue/toast'
import ConfirmDialog from 'primevue/confirmdialog'
import useTheme from '@/composables/useTheme'

const router = useRouter()
const route = useRoute()
const { locale, t } = useI18n()
const authStore = useAuthStore()
const { theme, toggleTheme } = useTheme()

watch(
  () => authStore.isAuthenticated,
  (isAuthenticated) => {
    if (!isAuthenticated && route.meta?.requiresAuth) {
      router.push({
        name: 'login',
        query: { redirect: route.fullPath }
      })
    }
  },
  { immediate: true }
)

const sidebarCollapsed = ref(false)
const userMenuRef = ref()

const routeTitleMap = {
  dashboard: 'nav.dashboard',
  animals: 'nav.animals',
  'animal-create': 'nav.animals',
  'animal-detail': 'nav.animals',
  'animal-edit': 'nav.animals',
  'adoption-applications': 'adoption.applications',
  'adoption-application-detail': 'adoption.applicationDetail',
  adoptions: 'adoption.title',
  'adoption-create': 'adoption.createAdoption',
  'adoption-detail': 'adoption.applicationDetail',
  veterinary: 'nav.veterinary',
  'veterinary-visits': 'veterinary.visits',
  'veterinary-visit-create': 'veterinary.addVisit',
  'veterinary-visit-edit': 'veterinary.visit',
  'veterinary-vaccinations': 'veterinary.vaccinations',
  'veterinary-medications': 'veterinary.medications',
  finance: 'nav.finance',
  events: 'nav.events',
  'finance-donations': 'nav.finance',
  'finance-donors': 'nav.donors',
  volunteers: 'nav.volunteers',
  partners: 'nav.partners',
  inventory: 'nav.inventory',
  reports: 'nav.reports',
  contacts: 'nav.contacts'
}

const locales = [
  { label: 'English', value: 'en' },
  { label: 'Polski', value: 'pl' }
]

const currentLocale = computed({
  get: () => locale.value,
  set: (val) => locale.value = val
})

const isAdmin = computed(() => authStore.isAdmin)
const userDisplayName = computed(() => {
  const user = authStore.user
  if (user?.first_name && user?.last_name) {
    return `${user.first_name} ${user.last_name}`
  }
  return user?.email || 'User'
})

const pageTitle = computed(() => {
  const routeName = route.name ? String(route.name) : ''
  const key = routeTitleMap[routeName] || route.meta?.titleKey || (routeName ? `nav.${routeName}` : '')
  if (key) {
    const translation = t(key)
    if (translation !== key) {
      return translation
    }
  }
  return routeName ? routeName.replace(/-/g, ' ') : ''
})

const userMenuItems = computed(() => [
  {
    label: t('nav.profile'),
    icon: 'pi pi-user',
    command: () => router.push('/profile')
  },
  {
    separator: true
  },
  {
    label: t('auth.logout'),
    icon: 'pi pi-sign-out',
    command: () => handleLogout()
  }
])

const toggleSidebar = () => {
  sidebarCollapsed.value = !sidebarCollapsed.value
  localStorage.setItem('sidebarCollapsed', sidebarCollapsed.value)
}

const toggleUserMenu = (event) => {
  userMenuRef.value.toggle(event)
}

const changeLocale = () => {
  localStorage.setItem('locale', currentLocale.value)
}

const getLocaleFlag = (loc) => {
  return loc === 'en' ? '🇬🇧' : '🇵🇱'
}

const goToHome = () => {
  router.push('/')
}

const handleLogout = async () => {
  await authStore.logout()
  router.push('/login')
}

// Load saved sidebar state
const savedState = localStorage.getItem('sidebarCollapsed')
if (savedState !== null) {
  sidebarCollapsed.value = savedState === 'true'
}
</script>

<style scoped>
.staff-layout {
  display: flex;
  min-height: 100vh;
  background: var(--surface-ground);
  color: var(--text-color);
}

.sidebar {
  width: 280px;
  background: var(--sidebar-bg, #1f2937);
  color: var(--sidebar-text, #f8fafc);
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 1000;
  border-right: 1px solid var(--sidebar-border, rgba(255, 255, 255, 0.08));
  box-shadow: 8px 0 40px rgba(15, 23, 42, 0.45);
}

.sidebar.collapsed {
  width: 80px;
}

.sidebar-header {
  padding: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--sidebar-border, rgba(255, 255, 255, 0.1));
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: var(--sidebar-text, #fff);
  text-decoration: none;
  font-size: 1.25rem;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
}

.brand i {
  font-size: 1.75rem;
  color: var(--brand-accent, #f97316);
  flex-shrink: 0;
}

.collapse-btn {
  background: none;
  border: none;
  color: inherit;
  cursor: pointer;
  padding: 0.5rem;
  border-radius: 4px;
  transition: background 0.3s;
  flex-shrink: 0;
}

.collapse-btn:hover {
  background: rgba(255, 255, 255, 0.1);
}

.sidebar-nav {
  flex: 1;
  padding: 1rem 0;
  overflow-y: auto;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  padding: 0.875rem 1.5rem;
  color: rgba(248, 250, 252, 0.8);
  text-decoration: none;
  transition: all 0.3s;
  white-space: nowrap;
  overflow: hidden;
}

.nav-item i {
  font-size: 1.25rem;
  flex-shrink: 0;
  width: 24px;
  text-align: center;
}

.nav-item:hover {
  background: rgba(255, 255, 255, 0.08);
  color: var(--sidebar-text, #fff);
}

.nav-item.router-link-active {
  background: var(--sidebar-active-bg, rgba(59, 130, 246, 0.25));
  color: var(--sidebar-text, #fff);
  border-left: 3px solid var(--primary-color, #3b82f6);
}

.sidebar-footer {
  border-top: 1px solid var(--sidebar-border, rgba(255, 255, 255, 0.1));
  padding: 1rem 0;
}

.main-wrapper {
  flex: 1;
  margin-left: 280px;
  display: flex;
  flex-direction: column;
  transition: margin-left 0.3s;
}

.sidebar.collapsed + .main-wrapper {
  margin-left: 80px;
}

.topbar {
  background: var(--topbar-bg, rgba(255, 255, 255, 0.9));
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 6px 20px rgba(15, 23, 42, 0.12);
  position: sticky;
  top: 0;
  z-index: 999;
  backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--border-color);
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: var(--topbar-text, #0f172a);
  margin: 0;
}

.topbar-right {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.locale-dropdown {
  width: 140px;
}

.locale-flag {
  font-size: 1.2rem;
  margin-right: 0.5rem;
}

.user-menu {
  position: relative;
}

.content {
  flex: 1;
  padding: 2rem;
  background: var(--surface-ground);
  min-height: 100%;
}

.theme-toggle {
  color: var(--topbar-text, var(--text-color));
}

@media (max-width: 968px) {
  .sidebar {
    transform: translateX(-100%);
  }

  .sidebar.collapsed {
    transform: translateX(0);
  }

  .main-wrapper {
    margin-left: 0;
  }

  .topbar {
    padding: 1rem;
  }

  .page-title {
    font-size: 1.25rem;
  }

  .content {
    padding: 1rem;
  }
}
</style>
