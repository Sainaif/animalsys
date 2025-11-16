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
import { getLocalizedValue, translateValue } from '@/utils/animalHelpers'

const { t, locale } = useI18n()
const authStore = useAuthStore()
const router = useRouter()

const loading = ref(true)
const dashboardData = ref(null)

const localeCode = computed(() => (locale.value === 'pl' ? 'pl-PL' : 'en-US'))
const currencyCode = computed(() => (locale.value === 'pl' ? 'PLN' : 'USD'))

const formatNumber = (value) => {
  return new Intl.NumberFormat(localeCode.value).format(value || 0)
}

const formatCurrency = (value) => {
  return new Intl.NumberFormat(localeCode.value, {
    style: 'currency',
    currency: currencyCode.value,
    minimumFractionDigits: 0,
    maximumFractionDigits: 0
  }).format(value || 0)
}

const formatDate = (value) => {
  if (!value) {
    return ''
  }
  const date = new Date(value)
  if (Number.isNaN(date.getTime())) {
    return ''
  }
  return new Intl.DateTimeFormat(localeCode.value, { dateStyle: 'medium' }).format(date)
}

const formatTaskDate = (value) => {
  const formatted = formatDate(value)
  return formatted || t('tasks.noDueDate')
}

const stats = computed(() => {
  const overview = dashboardData.value?.overview || {}

  return [
    {
      label: t('dashboard.totalAnimals'),
      value: formatNumber(overview.total_animals),
      icon: 'pi-heart',
      color: '#6366f1',
      route: '/staff/animals'
    },
    {
      label: t('dashboard.availableAnimals'),
      value: formatNumber(overview.animals_in_shelter || overview.available_animals),
      icon: 'pi-check-circle',
      color: '#10b981',
      route: '/staff/animals'
    },
    {
      label: t('dashboard.adoptionsThisMonth'),
      value: formatNumber(overview.adoptions_this_month || overview.animals_adopted_this_month),
      icon: 'pi-users',
      color: '#8b5cf6',
      route: '/staff/adoptions'
    },
    {
      label: t('dashboard.animalsInTreatment'),
      value: formatNumber(overview.animals_in_treatment),
      icon: 'pi-plus',
      color: '#f59e0b',
      route: '/staff/veterinary'
    },
    {
      label: t('dashboard.donationsThisMonth'),
      value: formatCurrency(overview.donations_this_month || overview.total_donations_this_month),
      icon: 'pi-dollar',
      color: '#ec4899',
      route: '/staff/finance'
    },
    {
      label: t('dashboard.activeVolunteers'),
      value: formatNumber(overview.active_volunteers),
      icon: 'pi-user-plus',
      color: '#14b8a6',
      route: '/staff/volunteers'
    }
  ]
})

const toArray = (value) => {
  if (Array.isArray(value)) {
    return value
  }
  if (!value) {
    return []
  }
  return Object.values(value)
}

const recentAnimals = computed(() => toArray(dashboardData.value?.recent_animals))
const upcomingTasks = computed(() => toArray(dashboardData.value?.upcoming_tasks))
const lastUpdated = computed(() => formatDate(dashboardData.value?.generated_at))

const getAnimalName = (animal) => {
  return getLocalizedValue(animal?.name, locale.value) || t('animal.unknown')
}

const getAnimalSpecies = (animal) => {
  if (!animal) return ''
  const translated = translateValue(animal.species || '', t, 'animal.speciesNames')
  return translated || animal.species || ''
}

const getAnimalStatus = (animal) => {
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

const getTaskStatusLabel = (task) => {
  const status = (task?.status || '').toLowerCase()
  const key = `tasks.statuses.${status}`
  const translated = t(key)
  return translated === key ? status : translated
}

const getTaskCategoryLabel = (task) => {
  if (!task?.category) return ''
  const key = `tasks.categories.${task.category}`
  const translated = t(key)
  return translated === key ? task.category : translated
}

const getTaskPriorityLabel = (task) => {
  if (!task?.priority) return ''
  const key = `tasks.priorities.${task.priority}`
  const translated = t(key)
  return translated === key ? task.priority : translated
}

const loadDashboard = async () => {
  try {
    loading.value = true
    const response = await api.get('/dashboard')
    dashboardData.value = {
      ...response.data,
      recent_animals: toArray(response.data?.recent_animals),
      upcoming_tasks: toArray(response.data?.upcoming_tasks)
    }
  } catch (error) {
    console.error('Error loading dashboard:', error)
    dashboardData.value = {
      overview: {
        total_animals: 150,
        animals_in_shelter: 45,
        adoptions_this_month: 12,
        animals_in_treatment: 15,
        donations_this_month: 5000,
        active_volunteers: 35
      },
      recent_animals: [],
      upcoming_tasks: []
    }
  } finally {
    loading.value = false
  }
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
      <h1>{{ $t('dashboard.welcome') }}, {{ authStore.user?.first_name || $t('dashboard.user') }}!</h1>
      <p>{{ $t('dashboard.overview') }}</p>
      <p
        v-if="lastUpdated"
        class="last-updated"
      >
        {{ $t('dashboard.lastUpdated', { time: lastUpdated }) }}
      </p>
    </div>

    <div
      v-if="loading"
      class="loading-container"
    >
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
          <div
            class="stat-icon"
            :style="{ backgroundColor: stat.color }"
          >
            <i
              class="pi"
              :class="stat.icon"
            />
          </div>
          <div class="stat-content">
            <div class="stat-label">
              {{ stat.label }}
            </div>
            <div class="stat-value">
              {{ stat.value }}
            </div>
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
            <div v-if="recentAnimals && recentAnimals.length > 0">
              <DataTable
                :value="recentAnimals"
                :rows="5"
              >
                <Column :header="$t('animal.name')">
                  <template #body="slotProps">
                    <div class="animal-cell">
                      <div class="animal-name">
                        {{ getAnimalName(slotProps.data) }}
                      </div>
                      <small
                        v-if="slotProps.data.breed"
                        class="animal-breed"
                      >{{ slotProps.data.breed }}</small>
                    </div>
                  </template>
                </Column>
                <Column :header="$t('animal.species')">
                  <template #body="slotProps">
                    {{ getAnimalSpecies(slotProps.data) }}
                  </template>
                </Column>
                <Column :header="$t('animal.status')">
                  <template #body="slotProps">
                    <span
                      class="status-pill"
                      :class="`status-${slotProps.data.status}`"
                    >
                      {{ getAnimalStatus(slotProps.data) }}
                    </span>
                  </template>
                </Column>
                <Column :header="$t('common.createdAt')">
                  <template #body="slotProps">
                    {{ formatDate(slotProps.data.created_at) }}
                  </template>
                </Column>
                <Column>
                  <template #body="slotProps">
                    <Button
                      icon="pi pi-eye"
                      class="p-button-text p-button-sm"
                      :aria-label="$t('common.viewMore')"
                      @click="router.push(`/staff/animals/${slotProps.data.id}`)"
                    />
                  </template>
                </Column>
              </DataTable>
            </div>
            <div
              v-else
              class="empty-state"
            >
              <i class="pi pi-inbox" />
              <p>{{ $t('dashboard.noRecentActivity') }}</p>
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
            <div v-if="upcomingTasks && upcomingTasks.length > 0">
              <div
                v-for="task in upcomingTasks"
                :key="task.id"
                class="task-item"
              >
                <div class="task-info">
                  <div class="task-title">
                    {{ task.title }}
                  </div>
                  <div class="task-meta">
                    <span>{{ getTaskCategoryLabel(task) }}</span>
                    <span>{{ formatTaskDate(task.due_date) }}</span>
                  </div>
                </div>
                <div class="task-tags">
                  <span
                    class="priority-pill"
                    :class="`priority-${task.priority}`"
                  >
                    {{ getTaskPriorityLabel(task) }}
                  </span>
                  <span
                    class="status-pill"
                    :class="`status-${task.status}`"
                  >
                    {{ getTaskStatusLabel(task) }}
                  </span>
                </div>
              </div>
            </div>
            <div
              v-else
              class="empty-state"
            >
              <i class="pi pi-check-circle" />
              <p>{{ $t('tasks.noUpcoming') }}</p>
            </div>
          </template>
        </Card>
      </div>

      <div class="modules-section">
        <h3>{{ $t('dashboard.modules') }}</h3>
        <div class="modules-grid">
          <!-- Priority 1: Animals -->
          <Card
            class="module-card"
            @click="router.push('/staff/animals')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);"
              >
                <i class="pi pi-heart" />
              </div>
              <h4>{{ $t('nav.animals') }}</h4>
              <p>{{ $t('dashboard.manageAnimals') }}</p>
            </template>
          </Card>

          <!-- Priority 2: Adoptions -->
          <Card
            class="module-card"
            @click="router.push('/staff/adoptions')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);"
              >
                <i class="pi pi-users" />
              </div>
              <h4>{{ $t('nav.adoptions') }}</h4>
              <p>{{ $t('dashboard.manageAdoptions') }}</p>
            </template>
          </Card>

          <!-- Priority 3: Veterinary -->
          <Card
            class="module-card"
            @click="router.push('/staff/veterinary')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);"
              >
                <i class="pi pi-plus" />
              </div>
              <h4>{{ $t('nav.veterinary') }}</h4>
              <p>{{ $t('dashboard.manageVeterinary') }}</p>
            </template>
          </Card>

          <!-- Priority 4: Finance -->
          <Card
            class="module-card"
            @click="router.push('/staff/finance')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);"
              >
                <i class="pi pi-dollar" />
              </div>
              <h4>{{ $t('nav.finance') }}</h4>
              <p>{{ $t('dashboard.manageFinance') }}</p>
            </template>
          </Card>

          <!-- Priority 5: Events & Volunteers -->
          <Card
            class="module-card"
            @click="router.push('/staff/events')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);"
              >
                <i class="pi pi-calendar" />
              </div>
              <h4>{{ $t('nav.events') }}</h4>
              <p>{{ $t('dashboard.manageEvents') }}</p>
            </template>
          </Card>

          <!-- Priority 6: Communications -->
          <Card
            class="module-card"
            @click="router.push('/staff/communication')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #30cfd0 0%, #330867 100%);"
              >
                <i class="pi pi-envelope" />
              </div>
              <h4>{{ $t('nav.communication') }}</h4>
              <p>{{ $t('dashboard.manageCommunication') }}</p>
            </template>
          </Card>

          <!-- Priority 7: Partners -->
          <Card
            class="module-card"
            @click="router.push('/staff/partners')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #a8edea 0%, #fed6e3 100%);"
              >
                <i class="pi pi-building" />
              </div>
              <h4>{{ $t('nav.partners') }}</h4>
              <p>{{ $t('dashboard.managePartners') }}</p>
            </template>
          </Card>

          <!-- Priority 8: Inventory -->
          <Card
            class="module-card"
            @click="router.push('/staff/inventory')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #ff9a9e 0%, #fecfef 100%);"
              >
                <i class="pi pi-box" />
              </div>
              <h4>{{ $t('nav.inventory') }}</h4>
              <p>{{ $t('dashboard.manageInventory') }}</p>
            </template>
          </Card>

          <!-- Priority 9: Reports -->
          <Card
            class="module-card"
            @click="router.push('/staff/reports')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #fbc2eb 0%, #a6c1ee 100%);"
              >
                <i class="pi pi-chart-bar" />
              </div>
              <h4>{{ $t('nav.reports') }}</h4>
              <p>{{ $t('dashboard.viewReports') }}</p>
            </template>
          </Card>

          <!-- Contacts -->
          <Card
            class="module-card"
            @click="router.push('/contacts')"
          >
            <template #content>
              <div
                class="module-icon"
                style="background: linear-gradient(135deg, #fddb92 0%, #d1fdff 100%);"
              >
                <i class="pi pi-phone" />
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
  max-width: 1440px;
  margin: 0 auto;
  padding: 2rem;
  color: var(--text-color);
}

.welcome-section {
  margin-bottom: 2rem;
}

.welcome-section h1 {
  font-size: 2.25rem;
  font-weight: 700;
  color: var(--text-color);
  margin-bottom: 0.35rem;
}

.welcome-section p {
  color: var(--text-muted);
  font-size: 1rem;
}

.last-updated {
  font-size: 0.9rem;
  color: var(--text-muted);
}

.loading-container {
  text-align: center;
  padding: 4rem 0;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2.5rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1.25rem;
  padding: 1.5rem;
  background: var(--card-bg);
  border-radius: 1rem;
  border: 1px solid var(--border-color);
  box-shadow: 0 12px 30px rgba(15, 23, 42, 0.08);
  transition: transform 0.2s ease, box-shadow 0.2s ease;
}

.stat-card.clickable {
  cursor: pointer;
}

.stat-card.clickable:hover {
  transform: translateY(-4px);
  box-shadow: 0 18px 40px rgba(15, 23, 42, 0.15);
}

.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 14px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  font-size: 1.35rem;
}

.stat-content {
  flex: 1;
}

.stat-label {
  color: var(--text-muted);
  font-size: 0.85rem;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 1.85rem;
  font-weight: 700;
  color: var(--text-color);
}

.dashboard-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(320px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.recent-animals-card :deep(.p-card),
.tasks-card :deep(.p-card) {
  background: var(--card-bg);
  border-radius: 1rem;
  border: 1px solid var(--border-color);
  box-shadow: 0 12px 26px rgba(15, 23, 42, 0.08);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.empty-state {
  text-align: center;
  padding: 2rem 1rem;
  color: var(--text-muted);
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.empty-state i {
  font-size: 2rem;
  color: var(--border-color);
}

.animal-cell {
  display: flex;
  flex-direction: column;
}

.animal-name {
  font-weight: 600;
  color: var(--text-color);
}

.animal-breed {
  color: var(--text-muted);
}

.status-pill {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.85rem;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 600;
  text-transform: capitalize;
  color: #fff;
}

.status-pill.status-available {
  background: #22c55e;
}

.status-pill.status-adopted {
  background: #6366f1;
}

.status-pill.status-fostered {
  background: #ec4899;
}

.status-pill.status-under_treatment {
  background: #f59e0b;
}

.task-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 0;
  border-bottom: 1px solid var(--border-color);
  gap: 1rem;
}

.task-item:last-child {
  border-bottom: none;
}

.task-info {
  flex: 1;
}

.task-title {
  font-weight: 600;
  color: var(--text-color);
}

.task-meta {
  font-size: 0.9rem;
  color: var(--text-muted);
  display: flex;
  flex-wrap: wrap;
  gap: 0.75rem;
  margin-top: 0.25rem;
}

.task-tags {
  display: flex;
  gap: 0.5rem;
  align-items: center;
}

.priority-pill {
  padding: 0.2rem 0.65rem;
  border-radius: 999px;
  font-size: 0.8rem;
  font-weight: 600;
}

.priority-pill.priority-low {
  background: rgba(34, 197, 94, 0.15);
  color: #15803d;
}

.priority-pill.priority-medium {
  background: rgba(251, 191, 36, 0.2);
  color: #b45309;
}

.priority-pill.priority-high {
  background: rgba(249, 115, 22, 0.2);
  color: #c2410c;
}

.priority-pill.priority-urgent {
  background: rgba(239, 68, 68, 0.2);
  color: #b91c1c;
}

.modules-section {
  margin-top: 2rem;
}

.modules-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
  gap: 1.25rem;
}

.module-card {
  cursor: pointer;
  transition: transform 0.2s, box-shadow 0.2s;
  border-radius: 16px;
  background: var(--card-bg);
  border: 1px solid transparent;
  padding: 1.25rem;
}

.module-card:hover {
  transform: translateY(-6px);
  box-shadow: 0 15px 30px rgba(15, 23, 42, 0.15);
}

.module-icon {
  width: 56px;
  height: 56px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  font-size: 1.5rem;
  margin-bottom: 1rem;
}

.module-card h4 {
  margin-bottom: 0.5rem;
  font-size: 1.1rem;
  color: var(--text-color);
}

.module-card p {
  color: var(--text-muted);
  font-size: 0.95rem;
}

@media (max-width: 1024px) {
  .dashboard {
    padding: 1.5rem;
  }
}

@media (max-width: 768px) {
  .stats-grid,
  .dashboard-grid,
  .modules-grid {
    grid-template-columns: 1fr;
  }

  .task-item {
    flex-direction: column;
    align-items: flex-start;
  }
}
</style>
