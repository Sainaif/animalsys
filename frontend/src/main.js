import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { createI18n } from 'vue-i18n'
import App from './App.vue'
import router from './router'
import { messages, availableLanguages } from './locales'
import './styles/main.css'

// Create i18n instance
const i18n = createI18n({
  legacy: false,
  locale: localStorage.getItem('locale') || 'pl',
  fallbackLocale: 'en',
  messages,
  availableLocales: availableLanguages,
})

// Create Pinia instance
const pinia = createPinia()

// Create Vue app
const app = createApp(App)

// Use plugins
app.use(pinia)
app.use(router)
app.use(i18n)

// Mount app
app.mount('#app')
