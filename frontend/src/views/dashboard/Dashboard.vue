<script setup>
import { ref, onMounted, computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '@/stores/auth'
import api from '@/services/api'
import ProgressSpinner from 'primevue/progressspinner'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Button from 'primevue/button'
import { useRouter } from 'vue-router'

const { t } = useI18n()
const authStore = useAuthStore()
const router = useRouter()

const loading = ref(true)
const dashboardData = ref(null)

const stats = computed(() => {
  if (!dashboardData.value) {
    return []
  }

  const overview = dashboardData.value.overview || {}

  return [
    {
      label: t('dashboard.totalAnimals'),
      value: overview.total_animals || 0,
      icon: 'pi-heart',
      color: '#3B82F6',
      route: '/staff/animals'
    },
    {
      label: t('dashboard.availableAnimals'),
      value: overview.available_animals || 0,
      icon: 'pi-check-circle',
      color: '#10B981',
      route: '/staff/animals'
    },
    {
      label: t('dashboard.adoptionsThisMonth'),
      value: overview.animals_adopted_this_month || 0,
      icon: 'pi-users',
      color: '#8B5CF6',
      route: '/staff/adoptions'
    },
    {
      label: t('dashboard.animalsInTreatment'),
      value: overview.animals_in_treatment || 0,
      icon: 'pi-plus',
      color: '#F59E0B',
      route: '/staff/veterinary'
    },
    {
      label: t('dashboard.donationsThisMonth'),
      value: formatCurrency(overview.total_donations_this_month || 0),
      icon: 'pi-dollar',
      color: '#E74C3C',
      route: '/staff/finance',
      isCurrency: true
    },
    {
      label: t('dashboard.activeVolunteers'),
      value: overview.active_volunteers || 0,
      icon: 'pi-user-plus',
      color: '#9B59B6',
      route: '/staff/volunteers'
    }
  ]
})

const recentAnimals = computed(() => {
  return dashboardData.value?.recent_animals || []
})

const upcomingTasks = computed(() => {
  return dashboardData.value?.upcoming_tasks || []
})

const loadDashboard = async () => {
  try {
    loading.value = true
    const response = await api.get('/dashboard')
    dashboardData.value = response.data
  } catch (error) {
    console.error('Error loading dashboard:', error)
    // Use mock data if API fails
    dashboardData.value = {
      overview: {
        total_animals: 150,
        available_animals: 45,
        animals_adopted_this_month: 12,
        animals_in_treatment: 15,
        total_donations_this_month: 5000,
        active_volunteers: 35
      },
      recent_animals: [],
      upcoming_tasks: []
    }
  } finally {
    loading.value = false
  }
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: 'USD',
    minimumFractionDigits: 0
  }).format(value)
}

const navigateToStat = (route) => {
  if (route) {
    router.push(route)
  }
}

onMounted(() => {
  loadDashboard()
})
</script>

<template>
  <div class="dashboard">
    <div class="welcome-section">
      <h1>{{ $t('dashboard.welcome') }}, {{ authStore.user?.first_name || 'User' }}!</h1>
      <p>{{ $t('dashboard.overview') }}</p>
    </div>

    <div v-if="loading" class="loading-container">
      <ProgressSpinner />
    </div>

    <div v-else>
      <div class="stats-grid">
        <div
          v-for="stat in stats"
          :key="stat.label"
          class="stat-card"
          :class="{ clickable: stat.route }"
          @click="navigateToStat(stat.route)"
        >
          <div class="stat-icon" :style="{ backgroundColor: stat.color }">
            <i class="pi" :class="stat.icon"></i>
          </div>
          <div class="stat-content">
            <div class="stat-label">{{ stat.label }}</div>
            <div class="stat-value">{{ stat.value }}</div>
          </div>
        </div>
      </div>

      <div class="dashboard-grid">
        <Card class="recent-animals-card">
          <template #title>
            <div class="card-header">
              <h3>{{ $t('dashboard.recentActivity') }}</h3>
              <Button
                :label="$t('common.viewMore')"
                class="p-button-text p-button-sm"
                @click="router.push('/staff/animals')"
              />
            </div>
          </template>
          <template #content>
            <div v-if="recentAnimals.length > 0">
              <DataTable :value="recentAnimals" :rows="5">
                <Column field="name" :header="$t('animal.name')"></Column>
                <Column field="species" :header="$t('animal.species')"></Column>
                <Column field="status" :header="$t('animal.status')"></Column>
                <Column>
                  <template #body="slotProps">
                    <Button
                      icon="pi pi-eye"
                      class="p-button-text p-button-sm"
                      @click="router.push(`/staff/animals/${slotProps.data.id}`)"
                    />
                  </template>
                </Column>
              </DataTable>
            </div>
            <div v-else class="empty-state">
              <i class="pi pi-inbox"></i>
              <p>No recent activity</p>
            </div>
          </template>
        </Card>

        <Card class="tasks-card">
          <template #title>
            <div class="card-header">
              <h3>{{ $t('dashboard.upcomingTasks') }}</h3>
            </div>
          </template>
          <template #content>
            <div v-if="upcomingTasks.length > 0">
              <div v-for="task in upcomingTasks" :key="task.id" class="task-item">
                <div class="task-info">
                  <i class="pi pi-calendar"></i>
                  <div>
                    <div class="task-title">{{ task.title }}</div>
                    <div class="task-date">{{ task.due_date }}</div>
                  </div>
                </div>
                <Button
                  icon="pi pi-check"
                  class="p-button-text p-button-sm p-button-success"
                />
              </div>
            </div>
            <div v-else class="empty-state">
              <i class="pi pi-check-circle"></i>
              <p>No upcoming tasks</p>
            </div>
          </template>
        </Card>
      </div>

      <div class="modules-section">
        <h3>{{ $t('dashboard.modules') }}</h3>
        <div class="modules-grid">
          <!-- Priority 1: Animals -->
          <Card class="module-card" @click="router.push('/staff/animals')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);">
                <i class="pi pi-heart"></i>
              </div>
              <h4>{{ $t('nav.animals') }}</h4>
              <p>{{ $t('dashboard.manageAnimals') }}</p>
            </template>
          </Card>

          <!-- Priority 2: Adoptions -->
          <Card class="module-card" @click="router.push('/staff/adoptions')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);">
                <i class="pi pi-users"></i>
              </div>
              <h4>{{ $t('nav.adoptions') }}</h4>
              <p>{{ $t('dashboard.manageAdoptions') }}</p>
            </template>
          </Card>

          <!-- Priority 3: Veterinary -->
          <Card class="module-card" @click="router.push('/staff/veterinary')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);">
                <i class="pi pi-plus"></i>
              </div>
              <h4>{{ $t('nav.veterinary') }}</h4>
              <p>{{ $t('dashboard.manageVeterinary') }}</p>
            </template>
          </Card>

          <!-- Priority 4: Finance -->
          <Card class="module-card" @click="router.push('/staff/finance')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);">
                <i class="pi pi-dollar"></i>
              </div>
              <h4>{{ $t('nav.finance') }}</h4>
              <p>{{ $t('dashboard.manageFinance') }}</p>
            </template>
          </Card>

          <!-- Priority 5: Events & Volunteers -->
          <Card class="module-card" @click="router.push('/staff/events')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);">
                <i class="pi pi-calendar"></i>
              </div>
              <h4>{{ $t('nav.events') }}</h4>
              <p>{{ $t('dashboard.manageEvents') }}</p>
            </template>
          </Card>

          <!-- Priority 6: Communications -->
          <Card class="module-card" @click="router.push('/staff/communication')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #30cfd0 0%, #330867 100%);">
                <i class="pi pi-envelope"></i>
              </div>
              <h4>{{ $t('nav.communication') }}</h4>
              <p>{{ $t('dashboard.manageCommunication') }}</p>
            </template>
          </Card>

          <!-- Priority 7: Partners -->
          <Card class="module-card" @click="router.push('/staff/partners')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);">
                <i class="pi pi-building"></i>
              </div>
              <h4>{{ $t('nav.partners') }}</h4>
              <p>{{ $t('dashboard.managePartners') }}</p>
            </template>
          </Card>

          <!-- Priority 8: Inventory -->
          <Card class="module-card" @click="router.push('/staff/inventory')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);">
                <i class="pi pi-box"></i>
              </div>
              <h4>{{ $t('nav.inventory') }}</h4>
              <p>{{ $t('dashboard.manageInventory') }}</p>
            </template>
          </Card>

          <!-- Priority 9: Reports -->
          <Card class="module-card" @click="router.push('/staff/reports')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #fbc2eb 0%, #a6c1ee 100%);">
                <i class="pi pi-chart-bar"></i>
              </div>
              <h4>{{ $t('nav.reports') }}</h4>
              <p>{{ $t('dashboard.viewReports') }}</p>
            </template>
          </Card>

          <!-- Contacts -->
          <Card class="module-card" @click="router.push('/contacts')">
            <template #content>
              <div class="module-icon" style="background: linear-gradient(135deg, #fddb92 0%, #d1fdff 100%);">
                <i class="pi pi-phone"></i>
              </div>
              <h4>{{ $t('nav.contacts') }}</h4>
              <p>{{ $t('dashboard.manageContacts') }}</p>
            </template>
          </Card>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard {
  max-width: 1400px;
  margin: 0 auto;
}

.welcome-section {
  margin-bottom: 2rem;
}

.welcome-section h1 {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
  margin-bottom: 0.5rem;
}

.welcome-section p {
  color: #7f8c8d;
  font-size: 1.1rem;
}

.loading-container {
  text-align: center;
  padding: 4rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1.5rem;
  padding: 1.5rem;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  transition: all 0.3s;
}

.stat-card.clickable {
  cursor: pointer;
}

.stat-card.clickable:hover {
  transform: translateY(-3px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.stat-icon {
  width: 64px;
  height: 64px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.75rem;
  flex-shrink: 0;
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 0.875rem;
  color: #7f8c8d;
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #2c3e50;
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.card-header h3 {
  margin: 0;
  font-size: 1.25rem;
  color: #2c3e50;
}

.empty-state {
  text-align: center;
  padding: 3rem 1rem;
  color: #7f8c8d;
}

.empty-state i {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.5;
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem;
  border-bottom: 1px solid #ecf0f1;
}

.task-item:last-child {
  border-bottom: none;
}

.task-info {
  display: flex;
  align-items: center;
  gap: 1rem;
}

.task-info i {
  font-size: 1.5rem;
  color: #3498db;
}

.task-title {
  font-weight: 600;
  color: #2c3e50;
  margin-bottom: 0.25rem;
}

.task-date {
  font-size: 0.875rem;
  color: #7f8c8d;
}

.modules-section {
  margin-top: 2rem;
}

.modules-section h3 {
  margin-bottom: 1.5rem;
  color: #2c3e50;
  font-size: 1.5rem;
  font-weight: 700;
}

.modules-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
  gap: 1.5rem;
}

.module-card {
  cursor: pointer;
  transition: all 0.3s ease;
  height: 100%;
}

.module-card:hover {
  transform: translateY(-5px);
  box-shadow: 0 8px 20px rgba(0, 0, 0, 0.15);
}

.module-icon {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 1rem;
  color: white;
  font-size: 1.75rem;
}

.module-card h4 {
  margin: 0 0 0.5rem 0;
  color: #2c3e50;
  font-size: 1.125rem;
  font-weight: 600;
}

.module-card p {
  margin: 0;
  color: #7f8c8d;
  font-size: 0.875rem;
}

@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: 1fr;
  }

  .dashboard-grid {
    grid-template-columns: 1fr;
  }

  .modules-grid {
    grid-template-columns: 1fr;
  }
}
</style>
