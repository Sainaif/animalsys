<template>
  <div class="event-dashboard">
    <div class="page-header">
      <h1>{{ $t('event.title') }} & {{ $t('event.volunteers') }}</h1>
    </div>

    <div class="dashboard-grid">
      <Card
        class="module-card"
        @click="router.push('/staff/events')"
      >
        <template #content>
          <div class="module-content">
            <i class="pi pi-calendar module-icon" />
            <h2>{{ $t('event.title') }}</h2>
            <p>Manage events and activities</p>
          </div>
        </template>
      </Card>

      <Card
        class="module-card"
        @click="router.push('/staff/volunteers')"
      >
        <template #content>
          <div class="module-content">
            <i class="pi pi-users module-icon" />
            <h2>{{ $t('event.volunteers') }}</h2>
            <p>Manage volunteer database</p>
          </div>
        </template>
      </Card>

      <Card
        class="module-card"
        @click="router.push('/staff/shifts')"
      >
        <template #content>
          <div class="module-content">
            <i class="pi pi-clock module-icon" />
            <h2>{{ $t('event.shifts') }}</h2>
            <p>Schedule and track volunteer shifts</p>
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
.event-dashboard { max-width: 1400px; margin: 0 auto; }
.page-header { margin-bottom: 2rem; }
.page-header h1 { font-size: 2rem; font-weight: 700; color: var(--heading-color); margin: 0; }
.dashboard-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1.5rem; }
.module-card { cursor: pointer; transition: transform 0.2s, box-shadow 0.2s; }
.module-card:hover { transform: translateY(-4px); box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15); }
.module-content { text-align: center; padding: 2rem 1rem; }
.module-icon { font-size: 3rem; color: #f59e0b; margin-bottom: 1rem; }
.module-content h2 { font-size: 1.5rem; font-weight: 600; color: var(--heading-color); margin: 0 0 0.5rem 0; }
.module-content p { color: var(--text-muted); margin: 0; }
</style>
