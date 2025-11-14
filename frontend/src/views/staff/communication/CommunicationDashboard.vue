<template>
  <div class="communication-dashboard">
    <h1>{{ $t('communication.title') }}</h1>

    <div class="dashboard-grid">
      <Card class="stat-card" @click="router.push('/staff/communication/templates')">
        <template #content>
          <div class="stat-content">
            <i class="pi pi-envelope" style="font-size: 2.5rem; color: #3b82f6;"></i>
            <div class="stat-info">
              <h3>{{ $t('communication.emailTemplates') }}</h3>
              <p class="stat-number">{{ stats.emailTemplates || 0 }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card" @click="router.push('/staff/communication/campaigns')">
        <template #content>
          <div class="stat-content">
            <i class="pi pi-send" style="font-size: 2.5rem; color: #10b981;"></i>
            <div class="stat-info">
              <h3>{{ $t('communication.emailCampaigns') }}</h3>
              <p class="stat-number">{{ stats.emailCampaigns || 0 }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card" @click="router.push('/staff/communication/logs')">
        <template #content>
          <div class="stat-content">
            <i class="pi pi-book" style="font-size: 2.5rem; color: #8b5cf6;"></i>
            <div class="stat-info">
              <h3>{{ $t('communication.communicationLogs') }}</h3>
              <p class="stat-number">{{ stats.communicationLogs || 0 }}</p>
            </div>
          </div>
        </template>
      </Card>

      <Card class="stat-card">
        <template #content>
          <div class="stat-content">
            <i class="pi pi-chart-line" style="font-size: 2.5rem; color: #f59e0b;"></i>
            <div class="stat-info">
              <h3>{{ $t('communication.sentCount') }}</h3>
              <p class="stat-number">{{ stats.totalSent || 0 }}</p>
            </div>
          </div>
        </template>
      </Card>
    </div>

    <div class="recent-section">
      <Card>
        <template #header>
          <h2>{{ $t('communication.emailCampaigns') }}</h2>
        </template>
        <template #content>
          <DataTable v-if="recentCampaigns.length > 0" :value="recentCampaigns" :rows="5">
            <Column field="name" :header="$t('communication.campaignName')" />
            <Column field="status" :header="$t('common.status')">
              <template #body="slotProps">
                <Badge :variant="getStatusVariant(slotProps.data.status)">
                  {{ $t(`communication.${slotProps.data.status}`) }}
                </Badge>
              </template>
            </Column>
            <Column field="sent_count" :header="$t('communication.sentCount')" />
            <Column field="scheduled_date" :header="$t('communication.scheduledDate')">
              <template #body="slotProps">
                {{ formatDate(slotProps.data.scheduled_date) }}
              </template>
            </Column>
          </DataTable>
          <EmptyState v-else :message="$t('communication.noEmailCampaignsFound')" />
        </template>
      </Card>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { communicationService } from '@/services/communicationService'
import Card from 'primevue/card'
import DataTable from 'primevue/datatable'
import Column from 'primevue/column'
import Badge from '@/components/shared/Badge.vue'
import EmptyState from '@/components/shared/EmptyState.vue'

const router = useRouter()
const { t } = useI18n()

const stats = ref({
  emailTemplates: 0,
  emailCampaigns: 0,
  communicationLogs: 0,
  totalSent: 0
})

const recentCampaigns = ref([])

const loadDashboard = async () => {
  try {
    const [templatesRes, campaignsRes, logsRes] = await Promise.all([
      communicationService.getEmailTemplates(),
      communicationService.getEmailCampaigns({ limit: 5 }),
      communicationService.getCommunicationLogs()
    ])

    stats.value.emailTemplates = templatesRes.data?.length || 0
    stats.value.emailCampaigns = campaignsRes.data?.length || 0
    stats.value.communicationLogs = logsRes.data?.length || 0
    stats.value.totalSent = campaignsRes.data?.reduce((sum, c) => sum + (c.sent_count || 0), 0) || 0

    recentCampaigns.value = campaignsRes.data || []
  } catch (error) {
    console.error('Failed to load communication dashboard:', error)
  }
}

const formatDate = (date) => date ? new Date(date).toLocaleDateString() : 'N/A'
const getStatusVariant = (status) => ({
  draft: 'neutral',
  scheduled: 'warning',
  sending: 'info',
  sent: 'success',
  failed: 'danger'
}[status] || 'neutral')

onMounted(loadDashboard)
</script>

<style scoped>
.communication-dashboard { max-width: 1400px; margin: 0 auto; }
.communication-dashboard h1 { font-size: 2rem; font-weight: 700; color: #2c3e50; margin-bottom: 2rem; }
.dashboard-grid { display: grid; grid-template-columns: repeat(auto-fit, minmax(250px, 1fr)); gap: 1.5rem; margin-bottom: 2rem; }
.stat-card { cursor: pointer; transition: transform 0.2s; }
.stat-card:hover { transform: translateY(-5px); }
.stat-content { display: flex; align-items: center; gap: 1.5rem; }
.stat-info { flex: 1; }
.stat-info h3 { font-size: 1rem; color: #6b7280; margin: 0 0 0.5rem 0; }
.stat-number { font-size: 2rem; font-weight: 700; color: #2c3e50; margin: 0; }
.recent-section { margin-top: 2rem; }
.recent-section h2 { font-size: 1.5rem; font-weight: 700; color: #2c3e50; margin: 0; padding: 1rem; }
</style>
