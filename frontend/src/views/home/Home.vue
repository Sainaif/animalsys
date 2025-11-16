<template>
  <div class="home-page">
    <!-- Hero Section -->
    <section class="hero">
      <div class="hero-content">
        <h1 class="hero-title">{{ $t('home.hero.title') }}</h1>
        <p class="hero-subtitle">{{ $t('home.hero.subtitle') }}</p>
        <div class="hero-actions">
          <Button
            :label="$t('home.hero.cta')"
            icon="pi pi-heart"
            class="p-button-lg p-button-danger"
            @click="scrollToAnimals"
          />
          <Button
            :label="$t('home.hero.donateCta')"
            icon="pi pi-dollar"
            class="p-button-lg p-button-outlined p-button-danger"
            @click="scrollToDonation"
          />
        </div>
      </div>
    </section>

    <!-- About Section -->
    <section id="about" class="about-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ $t('home.about.title') }}</h2>
          <p class="subtitle">{{ $t('home.about.subtitle') }}</p>
        </div>
        <div class="about-content">
          <div class="about-card">
            <i class="pi pi-heart"></i>
            <h3>{{ $t('home.about.mission') }}</h3>
          </div>
          <div class="about-card">
            <i class="pi pi-eye"></i>
            <h3>{{ $t('home.about.vision') }}</h3>
          </div>
          <div class="about-card">
            <i class="pi pi-star"></i>
            <h3>{{ $t('home.about.values') }}</h3>
          </div>
        </div>
      </div>
    </section>

    <!-- Statistics Section -->
    <section class="statistics-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ $t('home.statistics.title') }}</h2>
          <p class="subtitle">{{ $t('home.statistics.subtitle') }}</p>
        </div>
        <div class="stats-grid">
          <div v-for="stat in statistics" :key="stat.label" class="stat-card">
            <div class="stat-icon" :style="{ backgroundColor: stat.color }">
              <i class="pi" :class="stat.icon"></i>
            </div>
            <div class="stat-value">{{ formatNumber(stat.value) }}</div>
            <div class="stat-label">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Animals Section -->
    <section id="animals" class="animals-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ $t('home.adoptions.title') }}</h2>
          <p class="subtitle">{{ $t('home.adoptions.subtitle') }}</p>
        </div>

        <div v-if="loading" class="loading-container">
          <ProgressSpinner />
        </div>

        <div v-else-if="animals.length > 0" class="animals-grid">
          <div v-for="animal in animals" :key="animal.id" class="animal-card">
            <div class="animal-image">
              <img
                v-if="getAnimalImageSrc(animal)"
                :src="getAnimalImageSrc(animal)"
                :alt="getAnimalName(animal)"
              />
              <div v-else class="animal-placeholder">
                <i class="pi pi-image"></i>
              </div>
              <span class="animal-badge">{{ getStatusLabel(animal) }}</span>
            </div>
            <div class="animal-info">
              <h3>{{ getAnimalName(animal) }}</h3>
              <div class="animal-details">
                <span><i class="pi pi-calendar"></i> {{ formatAge(animal) }}</span>
                <span><i class="pi pi-tag"></i> {{ getAnimalSpeciesLabel(animal) }}</span>
                <span v-if="getAnimalColorLabel(animal)"><i class="pi pi-palette"></i> {{ getAnimalColorLabel(animal) }}</span>
                <span>
                  <i class="pi" :class="getGenderIcon(animal)"></i>
                  {{ getGenderLabel(animal) }}
                </span>
              </div>
              <p v-if="getAnimalDescription(animal)" class="animal-description">
                {{ truncateText(getAnimalDescription(animal), 100) }}
              </p>
              <Button
                :label="$t('home.adoptions.viewProfile')"
                class="p-button-outlined p-button-sm"
                @click="viewAnimal(animal.id)"
              />
            </div>
          </div>
        </div>

        <div class="section-footer">
          <Button
            :label="$t('home.adoptions.allAnimals')"
            icon="pi pi-arrow-right"
            class="p-button-lg p-button-outlined"
            @click="viewAllAnimals"
          />
        </div>
      </div>
    </section>

    <!-- How to Help Section -->
    <section id="help" class="help-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ $t('home.howToHelp.title') }}</h2>
          <p class="subtitle">{{ $t('home.howToHelp.subtitle') }}</p>
        </div>
        <div class="help-grid">
          <div class="help-card">
            <div class="help-icon" style="background-color: #3498db;">
              <i class="pi pi-dollar"></i>
            </div>
            <h3>{{ $t('home.howToHelp.donate.title') }}</h3>
            <p>{{ $t('home.howToHelp.donate.description') }}</p>
            <Button
              :label="$t('common.learnMore')"
              class="p-button-text"
              @click="scrollToDonation"
            />
          </div>
          <div class="help-card">
            <div class="help-icon" style="background-color: #2ecc71;">
              <i class="pi pi-users"></i>
            </div>
            <h3>{{ $t('home.howToHelp.volunteer.title') }}</h3>
            <p>{{ $t('home.howToHelp.volunteer.description') }}</p>
            <Button
              :label="$t('common.learnMore')"
              class="p-button-text"
              @click="scrollToContact"
            />
          </div>
          <div class="help-card">
            <div class="help-icon" style="background-color: #e74c3c;">
              <i class="pi pi-heart"></i>
            </div>
            <h3>{{ $t('home.howToHelp.adopt.title') }}</h3>
            <p>{{ $t('home.howToHelp.adopt.description') }}</p>
            <Button
              :label="$t('common.learnMore')"
              class="p-button-text"
              @click="scrollToAnimals"
            />
          </div>
          <div class="help-card">
            <div class="help-icon" style="background-color: #9b59b6;">
              <i class="pi pi-share-alt"></i>
            </div>
            <h3>{{ $t('home.howToHelp.spread.title') }}</h3>
            <p>{{ $t('home.howToHelp.spread.description') }}</p>
            <Button
              :label="$t('common.learnMore')"
              class="p-button-text"
            />
          </div>
        </div>
      </div>
    </section>

    <!-- Donation Section -->
    <section id="donation" class="donation-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ $t('home.donation.title') }}</h2>
          <p class="subtitle">{{ $t('home.donation.subtitle') }}</p>
        </div>
        <div class="donation-card">
          <div class="donation-form">
            <div class="donation-type">
              <Button
                :label="$t('home.donation.oneTime')"
                :class="{ 'p-button-outlined': donationType !== 'one-time' }"
                @click="donationType = 'one-time'"
              />
              <Button
                :label="$t('home.donation.monthly')"
                :class="{ 'p-button-outlined': donationType !== 'monthly' }"
                @click="donationType = 'monthly'"
              />
            </div>

            <div class="donation-amounts">
              <Button
                v-for="amount in donationAmounts"
                :key="amount"
                :label="`$${amount}`"
                :class="{ 'p-button-outlined': donationAmount !== amount }"
                @click="donationAmount = amount"
              />
              <InputNumber
                v-model="customAmount"
                :placeholder="$t('home.donation.customAmount')"
                mode="currency"
                currency="USD"
                @input="donationAmount = null"
              />
            </div>

            <Divider />

            <div class="donor-info">
              <h3>{{ $t('home.donation.yourInfo') }}</h3>
              <div class="form-grid">
                <div class="form-field">
                  <label>{{ $t('home.donation.firstName') }}</label>
                  <InputText v-model="donorInfo.firstName" />
                </div>
                <div class="form-field">
                  <label>{{ $t('home.donation.lastName') }}</label>
                  <InputText v-model="donorInfo.lastName" />
                </div>
                <div class="form-field">
                  <label>{{ $t('home.donation.email') }}</label>
                  <InputText v-model="donorInfo.email" type="email" />
                </div>
                <div class="form-field">
                  <label>{{ $t('home.donation.phone') }}</label>
                  <InputText v-model="donorInfo.phone" />
                </div>
              </div>
              <div class="form-field">
                <label>{{ $t('home.donation.message') }}</label>
                <Textarea v-model="donorInfo.message" rows="3" />
              </div>
            </div>

            <Button
              :label="$t('home.donation.submit')"
              icon="pi pi-check"
              class="p-button-lg p-button-success w-full"
              :disabled="!isValidDonation"
              @click="submitDonation"
            />
          </div>
        </div>
      </div>
    </section>

    <!-- Contact Section -->
    <section id="contact" class="contact-section">
      <div class="container">
        <div class="section-header">
          <h2>{{ $t('nav.contact') }}</h2>
          <p class="subtitle">{{ $t('home.hero.subtitle') }}</p>
        </div>
        <div class="contact-info">
          <div class="contact-card">
            <i class="pi pi-envelope"></i>
            <h3>{{ $t('home.footer.contact') }}</h3>
            <a href="mailto:info@animalfoundation.org">info@animalfoundation.org</a>
          </div>
          <div class="contact-card">
            <i class="pi pi-phone"></i>
            <h3>Phone</h3>
            <a href="tel:+1234567890">+1 (234) 567-890</a>
          </div>
          <div class="contact-card">
            <i class="pi pi-map-marker"></i>
            <h3>Address</h3>
            <p>123 Animal Street<br>City, State 12345</p>
          </div>
        </div>
      </div>
    </section>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import { useAuthStore } from '@/stores/auth'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Divider from 'primevue/divider'
import ProgressSpinner from 'primevue/progressspinner'
import api from '@/services/api'
import { getLocalizedValue, getAnimalImage, getAnimalGender, translateValue } from '@/utils/animalHelpers'

const { t, locale } = useI18n()
const router = useRouter()
const route = useRoute()
const toast = useToast()
const authStore = useAuthStore()

// Animals
const animals = ref([])
const loading = ref(true)

// Statistics
const statistics = ref([
  {
    label: t('home.statistics.animalsRescued'),
    value: 0,
    icon: 'pi-heart',
    color: '#3498db'
  },
  {
    label: t('home.statistics.successfulAdoptions'),
    value: 0,
    icon: 'pi-check-circle',
    color: '#2ecc71'
  },
  {
    label: t('home.statistics.activeVolunteers'),
    value: 0,
    icon: 'pi-users',
    color: '#9b59b6'
  },
  {
    label: t('home.statistics.donationsReceived'),
    value: 0,
    icon: 'pi-dollar',
    color: '#e74c3c'
  }
])

// Donation
const donationType = ref('one-time')
const donationAmount = ref(50)
const customAmount = ref(null)
const donationAmounts = [25, 50, 100, 250, 500]
const donorInfo = ref({
  firstName: '',
  lastName: '',
  email: '',
  phone: '',
  message: ''
})

const isValidDonation = computed(() => {
  const amount = donationAmount.value || customAmount.value
  return amount > 0 && donorInfo.value.firstName && donorInfo.value.email
})

// Methods
const loadAnimals = async () => {
  try {
    loading.value = true
    const response = await api.get('/public/animals', {
      params: { limit: 6, available_only: true, sort_by: 'created_at', sort_order: 'desc' }
    })
    const list = response.data?.animals || response.data?.data || []
    animals.value = Array.isArray(list) ? list : []
  } catch (error) {
    console.error('Error loading animals:', error)
    // Use mock data if API fails
    animals.value = []
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
  if (!authStore.isAuthenticated) {
    statistics.value[0].value = 150
    statistics.value[1].value = 89
    statistics.value[2].value = 42
    statistics.value[3].value = 50000
    return
  }

  try {
    // Try to load real statistics (this endpoint might require auth)
    const [animalStats] = await Promise.allSettled([
      api.get('/animals/statistics')
    ])

    if (animalStats.status === 'fulfilled') {
      const data = animalStats.value.data
      statistics.value[0].value = data.total_animals || 150
      statistics.value[1].value = data.by_status?.adopted || 89
    } else {
      // Mock data
      statistics.value[0].value = 150
      statistics.value[1].value = 89
      statistics.value[2].value = 42
      statistics.value[3].value = 50000
    }
  } catch (error) {
    console.error('Error loading statistics:', error)
    // Use mock data
    statistics.value[0].value = 150
    statistics.value[1].value = 89
    statistics.value[2].value = 42
    statistics.value[3].value = 50000
  }
}

const formatNumber = (num) => {
  if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

const getAnimalName = (animal) => {
  const name = getLocalizedValue(animal?.name, locale.value)
  return name || t('animal.unknown')
}

const getAnimalDescription = (animal) => getLocalizedValue(animal?.description, locale.value)

const getAnimalImageSrc = (animal) => getAnimalImage(animal)

const getAnimalSpeciesLabel = (animal) => {
  if (!animal) return ''
  if (animal.breed) {
    return animal.breed
  }
  const translated = translateValue(animal.species || '', t, 'animal.speciesNames')
  return translated || animal.species || ''
}

const getAnimalColorLabel = (animal) => {
  if (!animal?.color) return ''
  return translateValue(animal.color, t, 'animal.colorNames')
}

const getGenderLabel = (animal) => {
  const gender = getAnimalGender(animal)
  if (gender === 'male' || gender === 'female') {
    return t(`home.adoptions.${gender}`)
  }
  return t('animal.unknown')
}

const getGenderIcon = (animal) => getAnimalGender(animal) === 'male' ? 'pi-mars' : 'pi-venus'

const getStatusLabel = (animal) => {
  const status = (animal?.status || '').toLowerCase()
  switch (status) {
    case 'available':
      return t('animal.available')
    case 'adopted':
      return t('animal.adopted')
    case 'fostered':
      return t('animal.fostered')
    case 'under_treatment':
      return t('animal.underTreatment')
    default:
      return t('animal.status')
  }
}

const formatAge = (animal) => {
  if (!animal) return t('animal.unknown')

  if (animal.age_years || animal.age_months) {
    const years = animal.age_years || 0
    if (years > 0) {
      return `${years} ${t('home.adoptions.years')}`
    }
    return `${animal.age_months || 0} ${t('home.adoptions.months')}`
  }

  if (typeof animal.age === 'number') {
    return `${animal.age} ${t('home.adoptions.years')}`
  }

  if (animal.date_of_birth) {
    const dob = new Date(animal.date_of_birth)
    if (!Number.isNaN(dob.getTime())) {
      const diff = Date.now() - dob.getTime()
      const years = Math.floor(diff / (1000 * 60 * 60 * 24 * 365))
      if (years > 0) {
        return `${years} ${t('home.adoptions.years')}`
      }
      const months = Math.max(1, Math.floor(diff / (1000 * 60 * 60 * 24 * 30)))
      return `${months} ${t('home.adoptions.months')}`
    }
  }

  return t('animal.unknown')
}

const truncateText = (text, length) => {
  if (!text) return ''
  if (text.length <= length) return text
  return text.substring(0, length) + '...'
}

const scrollToAnimals = () => {
  document.getElementById('animals')?.scrollIntoView({ behavior: 'smooth' })
}

const scrollToDonation = () => {
  document.getElementById('donation')?.scrollIntoView({ behavior: 'smooth' })
}

const scrollToContact = () => {
  document.getElementById('contact')?.scrollIntoView({ behavior: 'smooth' })
}

const scrollToHash = (hash) => {
  if (!hash) return
  const target = hash.replace('#', '')
  requestAnimationFrame(() => {
    document.getElementById(target)?.scrollIntoView({ behavior: 'smooth' })
  })
}

const viewAnimal = (id) => {
  router.push({
    name: 'public-animals',
    query: id ? { animal: id } : undefined
  })
}

const viewAllAnimals = () => {
  router.push({ name: 'public-animals' })
}

const submitDonation = async () => {
  try {
    const amount = donationAmount.value || customAmount.value

    // Here you would integrate with a payment processor
    // For now, just show a success message
    toast.add({
      severity: 'success',
      summary: t('home.donation.thankYou'),
      detail: `Thank you for your ${amount} USD donation!`,
      life: 5000
    })

    // Reset form
    donationAmount.value = 50
    customAmount.value = null
    donorInfo.value = {
      firstName: '',
      lastName: '',
      email: '',
      phone: '',
      message: ''
    }
  } catch (error) {
    console.error('Error submitting donation:', error)
    toast.add({
      severity: 'error',
      summary: 'Error',
      detail: 'Failed to process donation. Please try again.',
      life: 5000
    })
  }
}

onMounted(() => {
  loadAnimals()
  loadStatistics()
  scrollToHash(route.hash)
})

watch(
  () => route.hash,
  (hash, prev) => {
    if (route.path === '/') {
      scrollToHash(hash)
    }
  }
)
</script>

<style scoped>
:global(body) {
  background-color: var(--surface-ground);
  font-family: 'Inter', system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
  color: var(--text-color);
}

.home-page {
  background: var(--surface-ground);
  color: var(--text-color);
}

.hero {
  position: relative;
  overflow: hidden;
  background: radial-gradient(circle at top, #ffe5ec 0%, #ffd6e0 40%, #f5f7fb 100%);
  color: #2a1f4f;
  padding: 6rem 2rem;
  text-align: center;
}

.hero::after {
  content: '';
  position: absolute;
  inset: 0;
  background: url("data:image/svg+xml,%3Csvg width='400' height='400' viewBox='0 0 400 400' xmlns='http://www.w3.org/2000/svg'%3E%3Cg fill='none' fill-opacity='0.2'%3E%3Cpath d='M0 200C120 120 280 280 400 200V400H0z' fill='%23ffcbd7'/%3E%3Cpath d='M0 0C130 60 270 60 400 0V400H0z' fill='%23ffe9f0'/%3E%3C/g%3E%3C/svg%3E");
  opacity: 0.5;
  pointer-events: none;
}

.hero-content {
  position: relative;
  max-width: 900px;
  margin: 0 auto;
}

.hero-title {
  font-size: clamp(2.5rem, 4vw, 4rem);
  font-weight: 800;
  margin-bottom: 1rem;
  line-height: 1.1;
}

.hero-subtitle {
  font-size: 1.3rem;
  margin-bottom: 2.5rem;
  opacity: 0.9;
}

.hero-actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 2rem;
}

.section-header {
  text-align: center;
  margin-bottom: 3rem;
}

.section-header h2 {
  font-size: clamp(2rem, 3vw, 2.8rem);
  font-weight: 700;
  margin-bottom: 0.6rem;
  color: var(--text-color);
  letter-spacing: -0.02em;
}

.subtitle {
  font-size: 1.15rem;
  color: var(--text-muted);
}

.about-section {
  padding: 5rem 0;
}

.about-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
}

.about-card {
  background: var(--card-bg);
  padding: 2rem;
  border-radius: 1.5rem;
  border: 1px solid var(--border-color);
  box-shadow: 0 20px 40px rgba(15, 23, 42, 0.08);
  text-align: center;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.about-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 30px 60px rgba(15, 23, 42, 0.15);
}

.about-card i {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 56px;
  height: 56px;
  border-radius: 18px;
  background: #fff1f5;
  color: #f43f5e;
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.about-card h3 {
  font-size: 1.1rem;
  color: var(--text-color);
  line-height: 1.6;
}

.statistics-section {
  padding: 5rem 0;
  background: radial-gradient(circle at top, rgba(255, 255, 255, 0.9) 0%, var(--surface-ground) 100%);
}

.statistics-section .section-header h2 {
  color: var(--text-color);
}

.statistics-section .subtitle {
  color: var(--text-muted);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.5rem;
}

.stat-card {
  background: var(--card-bg);
  padding: 2rem;
  border-radius: 1.5rem;
  border: 1px solid var(--border-color);
  text-align: center;
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.08);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 26px 60px rgba(79, 70, 229, 0.15);
}

.stat-icon {
  width: 72px;
  height: 72px;
  margin: 0 auto 1rem;
  border-radius: 1.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.75rem;
  color: white;
  background: linear-gradient(135deg, #f43f5e, #f97316);
  box-shadow: 0 14px 30px rgba(244, 63, 94, 0.35);
}

.stat-value {
  font-size: 2.5rem;
  font-weight: 800;
  margin-bottom: 0.2rem;
  color: var(--text-color);
}

.stat-label {
  font-size: 1rem;
  color: var(--text-muted);
}

.animals-section {
  padding: 5rem 0;
}

.loading-container {
  text-align: center;
  padding: 3rem;
}

.animals-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.75rem;
  margin-bottom: 2rem;
}

.animal-card {
  background: var(--card-bg);
  border-radius: 1.5rem;
  overflow: hidden;
  border: 1px solid var(--border-color);
  box-shadow: 0 18px 45px rgba(79, 70, 229, 0.12);
  transition: transform 0.35s ease, box-shadow 0.35s ease;
}

.animal-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 28px 60px rgba(79, 70, 229, 0.2);
}

.animal-image {
  position: relative;
  height: 230px;
  overflow: hidden;
  background: #f8fafc;
}

.animal-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.4s ease;
}

.animal-card:hover .animal-image img {
  transform: scale(1.05);
}

.animal-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 4rem;
  color: #cbd5f5;
}

.animal-badge {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: linear-gradient(120deg, #f43f5e, #f97316);
  color: white;
  padding: 0.4rem 1rem;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 600;
  box-shadow: 0 10px 30px rgba(244, 63, 94, 0.35);
}

.animal-info {
  padding: 1.5rem;
}

.animal-info h3 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
  color: var(--text-color);
}

.animal-details {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  color: var(--text-muted);
}

.animal-details span {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.animal-description {
  color: var(--text-muted);
  margin-bottom: 1rem;
  line-height: 1.6;
}

.section-footer {
  text-align: center;
  margin-top: 3rem;
}

.help-section {
  padding: 5rem 0;
  background: linear-gradient(180deg, #fdf2f8 0%, #ffffff 35%);
}

.help-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
}

.help-card {
  border-radius: 1.25rem;
  padding: 2rem;
  border: 1px solid rgba(244, 114, 182, 0.2);
  background: rgba(255, 255, 255, 0.75);
  box-shadow: 0 25px 50px rgba(244, 114, 182, 0.15);
  backdrop-filter: blur(6px);
  text-align: center;
  transition: transform 0.3s ease, box-shadow 0.3s ease;
}

.help-card:hover {
  transform: translateY(-8px);
  box-shadow: 0 35px 70px rgba(244, 63, 94, 0.25);
}

.help-icon {
  width: 70px;
  height: 70px;
  margin: 0 auto 1rem;
  border-radius: 1rem;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.8rem;
  color: white;
  background: #f472b6;
  box-shadow: 0 15px 30px rgba(244, 114, 182, 0.4);
}

.help-card h3 {
  font-size: 1.35rem;
  margin-bottom: 0.75rem;
  color: var(--text-color);
}

.help-card p {
  color: var(--text-muted);
  margin-bottom: 1rem;
  line-height: 1.6;
}

.donation-section {
  padding: 5rem 0;
}

.donation-card {
  max-width: 640px;
  margin: 0 auto;
  background: var(--card-bg);
  padding: 2.5rem;
  border-radius: 1.5rem;
  border: 1px solid var(--border-color);
  box-shadow: 0 25px 60px rgba(14, 116, 144, 0.15);
}

.donation-type {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(160px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.donation-amounts {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.donor-info {
  margin-bottom: 2rem;
}

.donor-info h3 {
  margin-bottom: 1rem;
  color: var(--text-color);
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1rem;
  margin-bottom: 1rem;
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.form-field label {
  font-weight: 600;
  color: var(--text-muted);
}

.w-full {
  width: 100%;
}

.contact-section {
  padding: 5rem 0;
  background: radial-gradient(circle at top, #ffffff 0%, #f0f4ff 100%);
}

.contact-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 1.5rem;
}

.contact-card {
  background: var(--card-bg);
  padding: 2rem;
  border-radius: 1.25rem;
  border: 1px solid var(--border-color);
  text-align: center;
  box-shadow: 0 15px 35px rgba(148, 163, 184, 0.2);
}

.contact-card i {
  font-size: 2.5rem;
  color: #6366f1;
  margin-bottom: 0.75rem;
}

.contact-card h3 {
  font-size: 1.3rem;
  margin-bottom: 0.75rem;
  color: var(--text-color);
}

.contact-card a {
  color: #2563eb;
  text-decoration: none;
  font-weight: 600;
}

.contact-card a:hover {
  text-decoration: underline;
}

:global([data-theme='dark'] .hero) {
  color: #f8fafc;
  background: radial-gradient(circle at top, rgba(36, 16, 64, 0.9) 0%, #0f172a 65%, #0b1120 100%);
}

:global([data-theme='dark'] .statistics-section) {
  background: radial-gradient(circle at top, rgba(30, 27, 75, 0.9) 0%, #0f172a 100%);
}

:global([data-theme='dark'] .help-section) {
  background: linear-gradient(180deg, rgba(36, 16, 64, 0.95) 0%, #0f172a 60%);
}

:global([data-theme='dark'] .contact-section) {
  background: radial-gradient(circle at top, rgba(30, 27, 75, 0.75) 0%, #0b1120 100%);
}

@media (max-width: 768px) {
  .hero {
    padding: 4rem 1.5rem;
  }

  .hero-actions {
    flex-direction: column;
  }

  .container {
    padding: 0 1.5rem;
  }
}
</style>
