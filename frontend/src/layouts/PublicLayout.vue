<template>
  <div class="public-layout">
    <nav class="navbar">
      <div class="navbar-container">
        <router-link
          to="/"
          class="navbar-brand"
        >
          <i class="pi pi-heart" />
          <span>{{ organization.shortName }}</span>
        </router-link>

        <div
          class="navbar-menu"
          :class="{ active: menuActive }"
        >
          <div class="nav-links">
            <router-link
              to="/"
              class="nav-link"
            >
              {{ $t('nav.home') }}
            </router-link>
            <button
              type="button"
              class="nav-link nav-link-button"
              @click="navigateToSection('about')"
            >
              {{ $t('nav.aboutUs') }}
            </button>
            <router-link
              to="/animals"
              class="nav-link"
            >
              {{ $t('nav.adoptAnimal') }}
            </router-link>
            <button
              type="button"
              class="nav-link nav-link-button"
              @click="navigateToSection('help')"
            >
              {{ $t('nav.howToHelp') }}
            </button>
            <button
              type="button"
              class="nav-link nav-link-button"
              @click="navigateToSection('contact')"
            >
              {{ $t('nav.contact') }}
            </button>
          </div>

          <div class="nav-actions">
            <Button
              :label="$t('common.donate')"
              icon="pi pi-heart"
              class="p-button-rounded p-button-danger"
              @click="scrollToDonation"
            />
            <router-link
              to="/login"
              class="staff-login-btn"
            >
              <Button
                :label="$t('auth.staffLogin')"
                icon="pi pi-sign-in"
                class="p-button-rounded p-button-outlined"
              />
            </router-link>
            <Button
              :icon="theme === 'light' ? 'pi pi-moon' : 'pi pi-sun'"
              class="p-button-rounded p-button-text theme-toggle"
              :aria-label="theme === 'light' ? $t('common.enableDarkMode') : $t('common.enableLightMode')"
              @click="toggleTheme"
            />
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
          </div>
        </div>

        <button
          class="mobile-menu-toggle"
          type="button"
          :aria-expanded="menuActive"
          :aria-label="menuActive ? $t('common.closeMenu') : $t('common.openMenu')"
          @click="menuActive = !menuActive"
        >
          <i
            class="pi"
            :class="menuActive ? 'pi-times' : 'pi-bars'"
          />
        </button>
      </div>
    </nav>

    <main class="main-content">
      <router-view />
    </main>

    <footer class="footer">
      <div class="footer-container">
        <div class="footer-section">
          <h3><i class="pi pi-heart" /> {{ organization.shortName }}</h3>
          <p>{{ organization.tagline }}</p>
        </div>
        <div class="footer-section">
          <h4>{{ $t('home.footer.about') }}</h4>
          <a href="#about">{{ $t('nav.aboutUs') }}</a>
          <router-link to="/animals">
            {{ $t('nav.adoptAnimal') }}
          </router-link>
          <a href="#help">{{ $t('nav.howToHelp') }}</a>
        </div>
        <div class="footer-section">
          <h4>{{ $t('home.footer.contact') }}</h4>
          <a :href="`mailto:${organization.contact.email}`">{{ organization.contact.email }}</a>
          <a :href="`tel:${organization.contact.phone}`">{{ organization.contact.phone }}</a>
          <p>{{ organization.contact.address }}</p>
        </div>
        <div class="footer-section">
          <h4>{{ $t('home.footer.followUs') }}</h4>
          <div class="social-links">
            <a
              href="#"
              aria-label="Facebook"
            ><i class="pi pi-facebook" /></a>
            <a
              href="#"
              aria-label="Twitter"
            ><i class="pi pi-twitter" /></a>
            <a
              href="#"
              aria-label="Instagram"
            ><i class="pi pi-instagram" /></a>
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        <p>&copy; {{ organization.currentYear }} {{ organization.legalName }}. {{ $t('home.footer.rights') }}</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed, watch, nextTick } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'
import useTheme from '@/composables/useTheme'
import organization from '@/config/organization'

const { locale } = useI18n()
const menuActive = ref(false)
const router = useRouter()
const route = useRoute()
const { theme, toggleTheme } = useTheme()

const locales = [
  { label: 'English', value: 'en' },
  { label: 'Polski', value: 'pl' }
]

const currentLocale = computed({
  get: () => locale.value,
  set: (val) => locale.value = val
})

const changeLocale = () => {
  localStorage.setItem('locale', currentLocale.value)
}

const getLocaleFlag = (loc) => {
  return loc === 'en' ? '🇬🇧' : '🇵🇱'
}

const closeMenu = () => {
  menuActive.value = false
}

watch(
  () => route.fullPath,
  () => closeMenu()
)

const scrollToSection = (sectionId) => {
  if (!sectionId) return
  const element = document.getElementById(sectionId)
  if (element) {
    const headerHeight = document.querySelector('.navbar')?.offsetHeight || 80
    const offset = element.getBoundingClientRect().top + window.scrollY - (headerHeight + 16)
    window.scrollTo({ top: Math.max(offset, 0), behavior: 'smooth' })
  }
}

const navigateToSection = async (sectionId) => {
  if (!sectionId) return
  closeMenu()
  const hash = `#${sectionId}`
  if (route.path !== '/') {
    await router.push({ path: '/', hash })
  } else if (route.hash !== hash) {
    router.replace({ hash })
  }
  nextTick(() => scrollToSection(sectionId))
}

const scrollToDonation = () => {
  navigateToSection('donation')
}

watch(
  () => route.hash,
  (hash) => {
    if (route.path === '/' && hash) {
      const target = hash.replace('#', '')
      nextTick(() => scrollToSection(target))
    }
  },
  { immediate: true }
)
</script>

<style scoped>
.public-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background: var(--surface-ground);
  color: var(--text-color);
}

.navbar {
  background: var(--nav-bg);
  box-shadow: 0 6px 30px rgba(15, 23, 42, 0.08);
  position: sticky;
  top: 0;
  z-index: 1000;
  padding: 1rem 0;
  backdrop-filter: blur(16px);
  border-bottom: 1px solid var(--nav-border);
}

.navbar-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  gap: 1.5rem;
}

.navbar-brand {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 1.5rem;
  font-weight: 700;
  color: #e74c3c;
  text-decoration: none;
  transition: transform 0.3s;
}

.navbar-brand:hover {
  transform: scale(1.05);
}

.navbar-brand i {
  font-size: 2rem;
}

.navbar-menu {
  display: flex;
  align-items: center;
  flex-wrap: nowrap;
  gap: 1.5rem;
  flex: 1;
  justify-content: space-between;
}

.nav-links {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  flex: 1;
  min-width: 0;
}

.nav-link {
  color: var(--nav-text, var(--text-color));
  text-decoration: none;
  font-weight: 600;
  font-size: 0.95rem;
  transition: color 0.3s, transform 0.3s;
  display: inline-flex;
  align-items: center;
  white-space: nowrap;
  padding: 0.35rem 0;
  letter-spacing: 0.01em;
}

.nav-link:hover,
.nav-link-button:focus-visible {
  color: var(--primary-color);
  transform: translateY(-1px);
}

.nav-link-button {
  border: none;
  background: transparent;
  padding: 0;
  font: inherit;
  cursor: pointer;
  color: inherit;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 0.85rem;
  flex-shrink: 0;
}

.staff-login-btn {
  text-decoration: none;
}

.theme-toggle {
  color: var(--nav-text);
}

.locale-dropdown {
  width: 120px;
  background: var(--card-bg);
  border-radius: 999px;
  border: 1px solid var(--border-color);
  transition: border-color 0.3s, box-shadow 0.3s;
}

.locale-dropdown :deep(.p-dropdown-label),
.locale-dropdown :deep(.p-dropdown-trigger-icon) {
  color: var(--text-color);
}

.locale-dropdown :deep(.p-dropdown-label) {
  padding-left: 0.75rem;
}

.locale-dropdown :deep(.p-dropdown-panel) {
  background: var(--card-bg);
  color: var(--text-color);
  border: 1px solid var(--border-color);
  box-shadow: 0 12px 24px rgba(15, 23, 42, 0.12);
}

.locale-flag {
  font-size: 1.2rem;
  margin-right: 0.5rem;
}

.mobile-menu-toggle {
  display: none;
  background: none;
  border: none;
  font-size: 1.5rem;
  cursor: pointer;
  color: var(--nav-text, #333);
}

.main-content {
  flex: 1;
}

.footer {
  background: #121a2c;
  color: #f8fafc;
  padding: 3rem 0 1rem;
  margin-top: 4rem;
}

.footer-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.footer-section h3 {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 1rem;
  color: #e74c3c;
}

.footer-section h4 {
  margin-bottom: 1rem;
}

.footer-section a {
  display: block;
  color: rgba(255, 255, 255, 0.8);
  text-decoration: none;
  margin-bottom: 0.5rem;
  transition: color 0.3s;
}

.footer-section a:hover {
  color: #f87171;
}

.social-links {
  display: flex;
  gap: 1rem;
  font-size: 1.5rem;
}

.social-links a {
  display: inline-flex;
  width: 40px;
  height: 40px;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.1);
  border-radius: 50%;
  transition: all 0.3s;
}

.social-links a:hover {
  background: #e74c3c;
  transform: translateY(-3px);
}

.footer-bottom {
  text-align: center;
  padding-top: 2rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
  max-width: 1200px;
  margin: 0 auto;
}

@media (max-width: 968px) {
  .mobile-menu-toggle {
    display: block;
  }

  .navbar-menu {
    position: fixed;
  top: 80px;
  left: 0;
  right: 0;
  background: var(--nav-bg);
  flex-direction: column;
  padding: 2rem;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  transform: translateY(-150%);
  transition: transform 0.3s;
  gap: 1.5rem;
  align-items: flex-start;
  }

  .navbar-menu.active {
    transform: translateY(0);
  }

  .nav-links {
    flex-direction: column;
    width: 100%;
    gap: 0.75rem;
  }

  .nav-link {
    width: 100%;
    justify-content: flex-start;
    white-space: normal;
  }

  .nav-actions {
    flex-direction: column;
    width: 100%;
    gap: 0.75rem;
  }

  .nav-actions button,
  .nav-actions .p-button {
    width: 100%;
  }
}
</style>
