import api from './api'
import type { FoundationSettings, OrganizationSettingsPayload, EmailSettingsPayload } from '@/types/settings'

export const settingsService = {
  async getSettings(): Promise<FoundationSettings> {
    const response = await api.get('/settings')
    return response.data
  },

  async updateOrganization(payload: OrganizationSettingsPayload) {
    const response = await api.put('/settings/organization', payload)
    return response.data
  },

  async updateEmailSettings(payload: EmailSettingsPayload) {
    const response = await api.put('/settings/email', payload)
    return response.data
  }
}
