<template>
  <div class="dashboard">
    <h1 class="page-title">{{ t('nav.dashboard') }}</h1>

    <!-- Welcome Section -->
    <div class="welcome-card">
      <h2>{{ t('dashboard.welcome') }}, {{ userName }}!</h2>
      <p>{{ t('dashboard.welcomeMessage') }}</p>
    </div>

    <!-- Quick Stats -->
    <div class="stats-grid">
      <div class="stat-card">
        <div class="stat-icon">🐾</div>
        <div class="stat-content">
          <div class="stat-value">0</div>
          <div class="stat-label">{{ t('dashboard.totalAnimals') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">❤️</div>
        <div class="stat-content">
          <div class="stat-value">0</div>
          <div class="stat-label">{{ t('dashboard.pendingAdoptions') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">🤝</div>
        <div class="stat-content">
          <div class="stat-value">0</div>
          <div class="stat-label">{{ t('dashboard.activeVolunteers') }}</div>
        </div>
      </div>

      <div class="stat-card">
        <div class="stat-icon">📅</div>
        <div class="stat-content">
          <div class="stat-value">0</div>
          <div class="stat-label">{{ t('dashboard.upcomingEvents') }}</div>
        </div>
      </div>
    </div>

    <!-- Quick Actions -->
    <div class="section">
      <h2 class="section-title">{{ t('dashboard.quickActions') }}</h2>
      <div class="actions-grid">
        <RouterLink v-if="authStore.hasRole('employee')" to="/app/animals/create" class="action-card">
          <span class="action-icon">➕</span>
          <span class="action-label">{{ t('dashboard.addAnimal') }}</span>
        </RouterLink>

        <RouterLink to="/app/adoptions/apply" class="action-card">
          <span class="action-icon">📝</span>
          <span class="action-label">{{ t('dashboard.newAdoption') }}</span>
        </RouterLink>

        <RouterLink v-if="authStore.hasRole('volunteer')" to="/app/schedules" class="action-card">
          <span class="action-icon">📅</span>
          <span class="action-label">{{ t('dashboard.viewSchedule') }}</span>
        </RouterLink>

        <RouterLink to="/app/profile" class="action-card">
          <span class="action-icon">👤</span>
          <span class="action-label">{{ t('dashboard.myProfile') }}</span>
        </RouterLink>
      </div>
    </div>

    <!-- Recent Activity (placeholder) -->
    <div class="section">
      <h2 class="section-title">{{ t('dashboard.recentActivity') }}</h2>
      <div class="activity-list">
        <p class="no-data">{{ t('common.noData') }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { RouterLink } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { useAuthStore } from '../stores/auth'

const { t } = useI18n()
const authStore = useAuthStore()

const userName = computed(() => {
  const user = authStore.user
  if (user?.first_name) {
    return user.first_name
  }
  return user?.username || 'User'
})
</script>

<style scoped>
.dashboard {
  max-width: 1400px;
}

.page-title {
  font-size: 2rem;
  font-weight: bold;
  margin-bottom: 2rem;
  color: var(--text-primary);
}

.welcome-card {
  background: linear-gradient(135deg, var(--primary-color) 0%, var(--secondary-color) 100%);
  color: white;
  padding: 2rem;
  border-radius: 1rem;
  margin-bottom: 2rem;
}

.welcome-card h2 {
  font-size: 1.5rem;
  margin-bottom: 0.5rem;
}

.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 3rem;
}

.stat-card {
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  padding: 1.5rem;
  display: flex;
  align-items: center;
  gap: 1rem;
  transition: all 0.2s;
}

.stat-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.stat-icon {
  font-size: 2.5rem;
}

.stat-content {
  flex: 1;
}

.stat-value {
  font-size: 2rem;
  font-weight: bold;
  color: var(--primary-color);
}

.stat-label {
  font-size: 0.875rem;
  color: var(--text-secondary);
  margin-top: 0.25rem;
}

.section {
  margin-bottom: 3rem;
}

.section-title {
  font-size: 1.5rem;
  font-weight: 600;
  margin-bottom: 1.5rem;
  color: var(--text-primary);
}

.actions-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 1rem;
}

.action-card {
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  padding: 1.5rem;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.75rem;
  text-decoration: none;
  color: var(--text-primary);
  transition: all 0.2s;
}

.action-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  border-color: var(--primary-color);
}

.action-icon {
  font-size: 2rem;
}

.action-label {
  font-weight: 500;
  text-align: center;
}

.activity-list {
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 0.75rem;
  padding: 2rem;
}

.no-data {
  text-align: center;
  color: var(--text-secondary);
}
</style>
