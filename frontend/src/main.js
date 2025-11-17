import { createApp } from 'vue'
import { createPinia } from 'pinia'
import PrimeVue from 'primevue/config'
import ToastService from 'primevue/toastservice'
import ConfirmationService from 'primevue/confirmationservice'
import BadgeDirective from 'primevue/badgedirective'
import Tooltip from 'primevue/tooltip'
import { createI18n } from 'vue-i18n'

import App from './App.vue'
import router from './router'

// PrimeVue styles
import 'primevue/resources/themes/lara-light-blue/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'

// Custom styles
import './assets/styles/main.css'

// i18n messages
import en from './i18n/en.json'
import pl from './i18n/pl.json'

const savedLocale = localStorage.getItem('locale') || 'pl'
const savedTheme = localStorage.getItem('theme') || 'light'
document.documentElement.setAttribute('data-theme', savedTheme)

// Create i18n instance
const i18n = createI18n({
  legacy: false,
  locale: savedLocale,
  fallbackLocale: 'en',
  linking: {
    enabled: false
  },
  messages: {
    en,
    pl
  }
})

// Create app
const app = createApp(App)

app.use(createPinia())
app.use(router)
app.use(PrimeVue, { ripple: true })
app.use(ToastService)
app.use(ConfirmationService)
app.use(i18n)

app.directive('badge', BadgeDirective)
app.directive('tooltip', Tooltip)

app.mount('#app')
