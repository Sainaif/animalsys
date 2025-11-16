<template>
  <div class="reports-dashboard">
    <h1>{{ $t('reports.title') }}</h1>

    <div class="reports-grid">
      <Card class="report-card">
        <template #header>
          <div class="card-header">
            <i class="pi pi-chart-bar" />
            <h2>{{ $t('reports.animalReports') }}</h2>
          </div>
        </template>
        <template #content>
          <p>Animal intake, outcomes, and population statistics</p>
          <Button
            :label="$t('reports.generateReport')"
            @click="generateReport('animals')"
          />
        </template>
      </Card>

      <Card class="report-card">
        <template #header>
          <div class="card-header">
            <i class="pi pi-heart" />
            <h2>{{ $t('reports.adoptionReports') }}</h2>
          </div>
        </template>
        <template #content>
          <p>Adoption applications, success rates, and trends</p>
          <Button
            :label="$t('reports.generateReport')"
            @click="generateReport('adoptions')"
          />
        </template>
      </Card>

      <Card class="report-card">
        <template #header>
          <div class="card-header">
            <i class="pi pi-dollar" />
            <h2>{{ $t('reports.financeReports') }}</h2>
          </div>
        </template>
        <template #content>
          <p>Donations, expenses, and financial summaries</p>
          <Button
            :label="$t('reports.generateReport')"
            @click="generateReport('finance')"
          />
        </template>
      </Card>

      <Card class="report-card">
        <template #header>
          <div class="card-header">
            <i class="pi pi-box" />
            <h2>{{ $t('reports.inventoryReports') }}</h2>
          </div>
        </template>
        <template #content>
          <p>Stock levels, usage patterns, and procurement needs</p>
          <Button
            :label="$t('reports.generateReport')"
            @click="generateReport('inventory')"
          />
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { useToast } from 'primevue/usetoast'
import Card from 'primevue/card'
import Button from 'primevue/button'

const toast = useToast()

const generateReport = (reportType) => {
  toast.add({
    severity: 'info',
    summary: 'Report Generation',
    detail: `Generating ${reportType} report...`,
    life: 3000
  })
}
</script>

<style scoped>
.reports-dashboard { max-width: 1400px; margin: 0 auto; }
.reports-dashboard h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin-bottom: 2rem; }
.reports-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 2rem; }
.report-card { height: 100%; }
.card-header { display: flex; align-items: center; gap: 1rem; padding: 1rem; background: #f3f4f6; }
.card-header i { font-size: 2rem; color: #3b82f6; }
.card-header h2 { font-size: 1.25rem; font-weight: 600; color: #2c3e50; margin: 0; }
.report-card p { color: #6b7280; margin-bottom: 1rem; }
</style>
