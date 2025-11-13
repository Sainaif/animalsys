<template>
  <div class="public-layout">
    <nav class="navbar">
      <div class="navbar-container">
        <router-link to="/" class="navbar-brand">
          <i class="pi pi-heart"></i>
          <span>Animal Foundation</span>
        </router-link>

        <div class="navbar-menu" :class="{ active: menuActive }">
          <router-link to="/" class="nav-link">{{ $t('nav.home') }}</router-link>
          <a href="#about" class="nav-link">{{ $t('nav.aboutUs') }}</a>
          <a href="#animals" class="nav-link">{{ $t('nav.adoptAnimal') }}</a>
          <a href="#help" class="nav-link">{{ $t('nav.howToHelp') }}</a>
          <a href="#contact" class="nav-link">{{ $t('nav.contact') }}</a>

          <div class="nav-actions">
            <Button
              :label="$t('common.donate')"
              icon="pi pi-heart"
              class="p-button-rounded p-button-danger"
              @click="scrollToDonation"
            />
            <router-link to="/login" class="staff-login-btn">
              <Button
                :label="$t('auth.staffLogin')"
                icon="pi pi-sign-in"
                class="p-button-rounded p-button-outlined"
              />
            </router-link>
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
          </div>
        </div>

        <button class="mobile-menu-toggle" @click="menuActive = !menuActive">
          <i class="pi" :class="menuActive ? 'pi-times' : 'pi-bars'"></i>
        </button>
      </div>
    </nav>

    <main class="main-content">
      <router-view />
    </main>

    <footer class="footer">
      <div class="footer-container">
        <div class="footer-section">
          <h3><i class="pi pi-heart"></i> Animal Foundation</h3>
          <p>{{ $t('home.hero.subtitle') }}</p>
        </div>
        <div class="footer-section">
          <h4>{{ $t('home.footer.about') }}</h4>
          <a href="#about">{{ $t('nav.aboutUs') }}</a>
          <a href="#animals">{{ $t('nav.adoptAnimal') }}</a>
          <a href="#help">{{ $t('nav.howToHelp') }}</a>
        </div>
        <div class="footer-section">
          <h4>{{ $t('home.footer.contact') }}</h4>
          <a href="mailto:info@animalfoundation.org">info@animalfoundation.org</a>
          <a href="tel:+1234567890">+1 (234) 567-890</a>
        </div>
        <div class="footer-section">
          <h4>{{ $t('home.footer.followUs') }}</h4>
          <div class="social-links">
            <a href="#" aria-label="Facebook"><i class="pi pi-facebook"></i></a>
            <a href="#" aria-label="Twitter"><i class="pi pi-twitter"></i></a>
            <a href="#" aria-label="Instagram"><i class="pi pi-instagram"></i></a>
          </div>
        </div>
      </div>
      <div class="footer-bottom">
        <p>&copy; 2025 Animal Foundation. {{ $t('home.footer.rights') }}</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import Button from 'primevue/button'
import Dropdown from 'primevue/dropdown'

const { locale } = useI18n()
const menuActive = ref(false)

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

const scrollToDonation = () => {
  const element = document.getElementById('donation')
  if (element) {
    element.scrollIntoView({ behavior: 'smooth' })
  }
}
</script>

<style scoped>
.public-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.navbar {
  background: white;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: sticky;
  top: 0;
  z-index: 1000;
  padding: 1rem 0;
}

.navbar-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
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
  gap: 2rem;
}

.nav-link {
  color: #333;
  text-decoration: none;
  font-weight: 500;
  transition: color 0.3s;
}

.nav-link:hover {
  color: #e74c3c;
}

.nav-actions {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.staff-login-btn {
  text-decoration: none;
}

.locale-dropdown {
  width: 120px;
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
  color: #333;
}

.main-content {
  flex: 1;
}

.footer {
  background: #2c3e50;
  color: white;
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
  color: #ecf0f1;
  text-decoration: none;
  margin-bottom: 0.5rem;
  transition: color 0.3s;
}

.footer-section a:hover {
  color: #e74c3c;
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
    background: white;
    flex-direction: column;
    padding: 2rem;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
    transform: translateY(-150%);
    transition: transform 0.3s;
    gap: 1rem;
  }

  .navbar-menu.active {
    transform: translateY(0);
  }

  .nav-actions {
    flex-direction: column;
    width: 100%;
  }

  .nav-actions button {
    width: 100%;
  }
}
</style>
