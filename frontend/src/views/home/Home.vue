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
                v-if="animal.photo_url"
                :src="animal.photo_url"
                :alt="animal.name"
              />
              <div v-else class="animal-placeholder">
                <i class="pi pi-image"></i>
              </div>
              <span class="animal-badge">{{ $t('home.adoptions.available') }}</span>
            </div>
            <div class="animal-info">
              <h3>{{ animal.name }}</h3>
              <div class="animal-details">
                <span><i class="pi pi-calendar"></i> {{ formatAge(animal) }}</span>
                <span><i class="pi pi-tag"></i> {{ animal.breed || animal.species }}</span>
                <span>
                  <i class="pi" :class="animal.gender === 'male' ? 'pi-mars' : 'pi-venus'"></i>
                  {{ $t(`home.adoptions.${animal.gender}`) }}
                </span>
              </div>
              <p v-if="animal.description" class="animal-description">
                {{ truncateText(animal.description, 100) }}
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
import { ref, computed, onMounted } from 'vue'
import { useI18n } from 'vue-i18n'
import { useToast } from 'primevue/usetoast'
import Button from 'primevue/button'
import InputText from 'primevue/inputtext'
import InputNumber from 'primevue/inputnumber'
import Textarea from 'primevue/textarea'
import Divider from 'primevue/divider'
import ProgressSpinner from 'primevue/progressspinner'
import api from '@/services/api'

const { t } = useI18n()
const toast = useToast()

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
      params: { limit: 6 }
    })
    animals.value = response.data.data || response.data
  } catch (error) {
    console.error('Error loading animals:', error)
    // Use mock data if API fails
    animals.value = []
  } finally {
    loading.value = false
  }
}

const loadStatistics = async () => {
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

const formatAge = (animal) => {
  if (!animal.age_years && !animal.age_months) return 'Unknown'
  const years = animal.age_years || 0
  const months = animal.age_months || 0

  if (years > 0) {
    return `${years} ${t('home.adoptions.years')}`
  }
  return `${months} ${t('home.adoptions.months')}`
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

const viewAnimal = (id) => {
  // In future, navigate to animal detail page
  console.log('View animal:', id)
}

const viewAllAnimals = () => {
  scrollToAnimals()
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
})
</script>

<style scoped>
.hero {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 6rem 2rem;
  text-align: center;
}

.hero-content {
  max-width: 800px;
  margin: 0 auto;
}

.hero-title {
  font-size: 3.5rem;
  font-weight: 800;
  margin-bottom: 1rem;
  line-height: 1.2;
}

.hero-subtitle {
  font-size: 1.5rem;
  margin-bottom: 2rem;
  opacity: 0.95;
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
  font-size: 2.5rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.subtitle {
  font-size: 1.25rem;
  color: #7f8c8d;
}

.about-section {
  padding: 5rem 0;
  background: #f8f9fa;
}

.about-content {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
  gap: 2rem;
}

.about-card {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
  transition: transform 0.3s;
}

.about-card:hover {
  transform: translateY(-5px);
}

.about-card i {
  font-size: 3rem;
  color: #e74c3c;
  margin-bottom: 1rem;
}

.about-card h3 {
  font-size: 1.1rem;
  color: #2c3e50;
  line-height: 1.6;
}

.statistics-section {
  padding: 5rem 0;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
}

.statistics-section .section-header h2 {
  color: white;
}

.statistics-section .subtitle {
  color: rgba(255, 255, 255, 0.9);
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.stat-card {
  background: rgba(255, 255, 255, 0.1);
  backdrop-filter: blur(10px);
  padding: 2rem;
  border-radius: 12px;
  text-align: center;
  transition: transform 0.3s;
}

.stat-card:hover {
  transform: translateY(-5px);
}

.stat-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 1rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
}

.stat-value {
  font-size: 3rem;
  font-weight: 800;
  margin-bottom: 0.5rem;
}

.stat-label {
  font-size: 1.1rem;
  opacity: 0.9;
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
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.animal-card {
  background: white;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s, box-shadow 0.3s;
}

.animal-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
}

.animal-image {
  position: relative;
  height: 250px;
  overflow: hidden;
  background: #f0f0f0;
}

.animal-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.animal-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 4rem;
  color: #bdc3c7;
}

.animal-badge {
  position: absolute;
  top: 1rem;
  right: 1rem;
  background: #e74c3c;
  color: white;
  padding: 0.5rem 1rem;
  border-radius: 20px;
  font-size: 0.875rem;
  font-weight: 600;
}

.animal-info {
  padding: 1.5rem;
}

.animal-info h3 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
  color: #2c3e50;
}

.animal-details {
  display: flex;
  flex-wrap: wrap;
  gap: 1rem;
  margin-bottom: 1rem;
  font-size: 0.9rem;
  color: #7f8c8d;
}

.animal-details span {
  display: flex;
  align-items: center;
  gap: 0.25rem;
}

.animal-description {
  color: #7f8c8d;
  margin-bottom: 1rem;
  line-height: 1.6;
}

.section-footer {
  text-align: center;
  margin-top: 3rem;
}

.help-section {
  padding: 5rem 0;
  background: #f8f9fa;
}

.help-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.help-card {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
  transition: transform 0.3s;
}

.help-card:hover {
  transform: translateY(-5px);
}

.help-icon {
  width: 80px;
  height: 80px;
  margin: 0 auto 1rem;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 2rem;
  color: white;
}

.help-card h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: #2c3e50;
}

.help-card p {
  color: #7f8c8d;
  margin-bottom: 1rem;
  line-height: 1.6;
}

.donation-section {
  padding: 5rem 0;
}

.donation-card {
  max-width: 600px;
  margin: 0 auto;
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.donation-type {
  display: flex;
  gap: 1rem;
  margin-bottom: 2rem;
}

.donation-type button {
  flex: 1;
}

.donation-amounts {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 1rem;
  margin-bottom: 2rem;
}

.donation-amounts button {
  width: 100%;
}

.donor-info {
  margin-bottom: 2rem;
}

.donor-info h3 {
  margin-bottom: 1rem;
  color: #2c3e50;
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
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
  color: #2c3e50;
}

.w-full {
  width: 100%;
}

.contact-section {
  padding: 5rem 0;
  background: #f8f9fa;
}

.contact-info {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 2rem;
}

.contact-card {
  background: white;
  padding: 2rem;
  border-radius: 12px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  text-align: center;
}

.contact-card i {
  font-size: 3rem;
  color: #e74c3c;
  margin-bottom: 1rem;
}

.contact-card h3 {
  font-size: 1.5rem;
  margin-bottom: 1rem;
  color: #2c3e50;
}

.contact-card a {
  color: #3498db;
  text-decoration: none;
}

.contact-card a:hover {
  text-decoration: underline;
}

@media (max-width: 768px) {
  .hero-title {
    font-size: 2rem;
  }

  .hero-subtitle {
    font-size: 1.1rem;
  }

  .section-header h2 {
    font-size: 2rem;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .donation-amounts {
    grid-template-columns: repeat(2, 1fr);
  }
}
</style>
