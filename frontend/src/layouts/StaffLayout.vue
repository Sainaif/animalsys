<template>
  <div class="staff-layout">
    <aside class="sidebar" :class="{ collapsed: sidebarCollapsed }">
      <div class="sidebar-header">
        <router-link to="/dashboard" class="brand">
          <i class="pi pi-heart"></i>
          <span v-if="!sidebarCollapsed">Animal Foundation</span>
        </router-link>
        <button class="collapse-btn" @click="toggleSidebar">
          <i class="pi" :class="sidebarCollapsed ? 'pi-angle-right' : 'pi-angle-left'"></i>
        </button>
      </div>

      <nav class="sidebar-nav">
        <router-link to="/dashboard" class="nav-item">
          <i class="pi pi-th-large"></i>
          <span v-if="!sidebarCollapsed">{{ $t('nav.dashboard') }}</span>
        </router-link>
        <router-link to="/staff/animals" class="nav-item">
          <i class="pi pi-heart"></i>
          <span v-if="!sidebarCollapsed">{{ $t('nav.animals') }}</span>
        </router-link>
        <router-link to="/staff/veterinary" class="nav-item">
          <i class="pi pi-plus"></i>
          <span v-if="!sidebarCollapsed">{{ $t('nav.veterinary') }}</span>
        </router-link>
        <router-link to="/staff/adoptions/applications" class="nav-item">
          <i class="pi pi-users"></i>
          <span v-if="!sidebarCollapsed">{{ $t('nav.adoptions') }}</span>
        </router-link>
        <router-link to="/contacts" class="nav-item">
          <i class="pi pi-phone"></i>
          <span v-if="!sidebarCollapsed">{{ $t('nav.contacts') }}</span>
        </router-link>

        <Divider v-if="!sidebarCollapsed" />

        <router-link v-if="isAdmin" to="/users" class="nav-item">
          <i class="pi pi-user"></i>
          <span v-if="!sidebarCollapsed">{{ $t('nav.users') }}</span>
        </router-link>
        <router-link v-if="isAdmin" to="/settings" class="nav-item">
          <i class="pi pi-cog"></i>
          <span v-if="!sidebarCollapsed">{{ $t('nav.settings') }}</span>
        </router-link>
      </nav>

      <div class="sidebar-footer">
        <router-link to="/profile" class="nav-item">
          <i class="pi pi-user"></i>
          <span v-if="!sidebarCollapsed">{{ $t('nav.profile') }}</span>
        </router-link>
      </div>
    </aside>

    <div class="main-wrapper">
      <header class="topbar">
        <div class="topbar-left">
          <h2 class="page-title">{{ pageTitle }}</h2>
        </div>
        <div class="topbar-right">
          <Dropdown
            v-model="currentLocale"
            :options="locales"
            optionLabel="label"
            optionValue="value"
            @change="changeLocale"
            class="locale-dropdown"
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
            icon="pi pi-home"
            class="p-button-text p-button-rounded"
            v-tooltip.bottom="'Home'"
            @click="goToHome"
          />

          <Button
            icon="pi pi-bell"
            class="p-button-text p-button-rounded"
            v-tooltip.bottom="'Notifications'"
            v-badge="3"
          />

          <div class="user-menu">
            <Button
              :label="userDisplayName"
              icon="pi pi-user"
              class="p-button-text"
              @click="toggleUserMenu"
            />
            <Menu ref="userMenuRef" :model="userMenuItems" :popup="true" />
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

const router = useRouter()
const route = useRoute()
const { locale, t } = useI18n()
const authStore = useAuthStore()

const sidebarCollapsed = ref(false)
const userMenuRef = ref()

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
  const name = route.name
  if (name) {
    const key = `nav.${name}`
    const translation = t(key)
    return translation !== key ? translation : name
  }
  return ''
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
}

.sidebar {
  width: 280px;
  background: #2c3e50;
  color: white;
  display: flex;
  flex-direction: column;
  transition: width 0.3s;
  position: fixed;
  left: 0;
  top: 0;
  bottom: 0;
  z-index: 1000;
}

.sidebar.collapsed {
  width: 80px;
}

.sidebar-header {
  padding: 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid rgba(255, 255, 255, 0.1);
}

.brand {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  color: white;
  text-decoration: none;
  font-size: 1.25rem;
  font-weight: 700;
  white-space: nowrap;
  overflow: hidden;
}

.brand i {
  font-size: 1.75rem;
  color: #e74c3c;
  flex-shrink: 0;
}

.collapse-btn {
  background: none;
  border: none;
  color: white;
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
  color: rgba(255, 255, 255, 0.8);
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
  background: rgba(255, 255, 255, 0.1);
  color: white;
}

.nav-item.router-link-active {
  background: rgba(231, 76, 60, 0.2);
  color: white;
  border-left: 3px solid #e74c3c;
}

.sidebar-footer {
  border-top: 1px solid rgba(255, 255, 255, 0.1);
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
  background: white;
  padding: 1rem 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 999;
}

.page-title {
  font-size: 1.5rem;
  font-weight: 600;
  color: #2c3e50;
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
