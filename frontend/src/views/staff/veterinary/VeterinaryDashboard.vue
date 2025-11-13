<template>
  <div class="veterinary-dashboard">
    <div class="page-header">
      <h1>{{ $t('veterinary.title') }}</h1>
    </div>

    <div class="dashboard-grid">
      <Card class="module-card" @click="router.push('/staff/veterinary/visits')">
        <template #content>
          <div class="module-content">
            <i class="pi pi-calendar module-icon"></i>
            <h2>{{ $t('veterinary.visits') }}</h2>
            <p>Manage veterinary visits and checkups</p>
          </div>
        </template>
      </Card>

      <Card class="module-card" @click="router.push('/staff/veterinary/vaccinations')">
        <template #content>
          <div class="module-content">
            <i class="pi pi-shield module-icon"></i>
            <h2>{{ $t('veterinary.vaccinations') }}</h2>
            <p>Track vaccination records and schedules</p>
          </div>
        </template>
      </Card>

      <Card class="module-card" @click="router.push('/staff/veterinary/medications')">
        <template #content>
          <div class="module-content">
            <i class="pi pi-box module-icon"></i>
            <h2>{{ $t('veterinary.medications') }}</h2>
            <p>Manage medications and prescriptions</p>
          </div>
        </template>
      </Card>

      <Card class="module-card" @click="router.push('/staff/veterinary/treatment-plans')">
        <template #content>
          <div class="module-content">
            <i class="pi pi-list module-icon"></i>
            <h2>{{ $t('veterinary.treatmentPlans') }}</h2>
            <p>Manage treatment plans and progress</p>
          </div>
        </template>
      </Card>

      <Card class="module-card" @click="router.push('/staff/veterinary/medical-conditions')">
        <template #content>
          <div class="module-content">
            <i class="pi pi-heart module-icon"></i>
            <h2>{{ $t('veterinary.medicalConditions') }}</h2>
            <p>Track medical conditions and diagnoses</p>
          </div>
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'
import Card from 'primevue/card'

const router = useRouter()
</script>

<style scoped>
.veterinary-dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

.page-header {
  margin-bottom: 2rem;
}

.page-header h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin: 0;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

.module-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
}

.module-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.module-content {
  text-align: center;
  padding: 2rem 1rem;
}

.module-icon {
  font-size: 3rem;
  color: #3b82f6;
  margin-bottom: 1rem;
}

.module-content h2 {
  font-size: 1.5rem;
  font-weight: 600;
  color: #2c3e50;
  margin: 0 0 0.5rem 0;
}

.module-content p {
  color: #6b7280;
  margin: 0;
}
</style>
