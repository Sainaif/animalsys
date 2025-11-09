<script setup>
import { ref } from 'vue'

const stats = ref([
  { label: 'Total Animals', value: 0, icon: 'pi-heart', color: '#3B82F6' },
  { label: 'Adoptions This Month', value: 0, icon: 'pi-users', color: '#10B981' },
  { label: 'Upcoming Vaccinations', value: 0, icon: 'pi-calendar', color: '#F59E0B' },
  { label: 'Active Volunteers', value: 0, icon: 'pi-user-plus', color: '#8B5CF6' }
])
</script>

<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">Dashboard</h1>
    </div>

    <div class="stats-grid">
      <div v-for="stat in stats" :key="stat.label" class="stat-card card">
        <div class="stat-icon" :style="{ backgroundColor: stat.color }">
          <i class="pi" :class="stat.icon"></i>
        </div>
        <div class="stat-content">
          <div class="stat-label">{{ stat.label }}</div>
          <div class="stat-value">{{ stat.value }}</div>
        </div>
      </div>
    </div>

    <div class="dashboard-content">
      <div class="card">
        <h3>Welcome to Animal Foundation CRM</h3>
        <p>The system is being built. Features will be added progressively.</p>
      </div>
    </div>
  </div>
</template>

<style scoped>
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
  margin-bottom: 2rem;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1.5rem;
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
}

.stat-content {
  flex: 1;
}

.stat-label {
  font-size: 0.875rem;
  color: #6c757d;
  margin-bottom: 0.25rem;
}

.stat-value {
  font-size: 2rem;
  font-weight: 700;
  color: #212529;
}

.dashboard-content {
  margin-top: 2rem;
}
</style>
