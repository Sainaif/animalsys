import api from './api'
import type { AdoptionApplication, Adoption, AdoptionStatistics } from '@/types/adoption'
import type { PaginatedResponse, QueryParams } from '@/types/common'

export const adoptionService = {
  // Applications
  async getApplications(params?: QueryParams): Promise<PaginatedResponse<AdoptionApplication>> {
    const response = await api.get('/adoptions/applications', { params })
    return response.data
  },

  async getApplication(id: string): Promise<AdoptionApplication> {
    const response = await api.get(`/adoptions/applications/${id}`)
    return response.data
  },

  async createApplication(data: Partial<AdoptionApplication>): Promise<AdoptionApplication> {
    const response = await api.post('/adoptions/applications', data)
    return response.data
  },

  async updateApplication(id: string, data: Partial<AdoptionApplication>): Promise<AdoptionApplication> {
    const response = await api.put(`/adoptions/applications/${id}`, data)
    return response.data
  },

  async deleteApplication(id: string): Promise<void> {
    await api.delete(`/adoptions/applications/${id}`)
  },

  async approveApplication(id: string, notes?: string): Promise<AdoptionApplication> {
    const response = await api.post(`/adoptions/applications/${id}/approve`, { notes })
    return response.data
  },

  async rejectApplication(id: string, reason: string): Promise<AdoptionApplication> {
    const response = await api.post(`/adoptions/applications/${id}/reject`, { reason })
    return response.data
  },

  async getPendingApplications(): Promise<AdoptionApplication[]> {
    const response = await api.get('/adoptions/applications/pending')
    return response.data
  },

  // Adoptions
  async getAdoptions(params?: QueryParams): Promise<PaginatedResponse<Adoption>> {
    const response = await api.get('/adoptions', { params })
    return response.data
  },

  async getAdoption(id: string): Promise<Adoption> {
    const response = await api.get(`/adoptions/${id}`)
    return response.data
  },

  async createAdoption(data: Partial<Adoption>): Promise<Adoption> {
    const response = await api.post('/adoptions', data)
    return response.data
  },

  async updateAdoption(id: string, data: Partial<Adoption>): Promise<Adoption> {
    const response = await api.put(`/adoptions/${id}`, data)
    return response.data
  },

  async deleteAdoption(id: string): Promise<void> {
    await api.delete(`/adoptions/${id}`)
  },

  async returnAdoption(id: string, reason: string): Promise<Adoption> {
    const response = await api.post(`/adoptions/${id}/return`, { reason })
    return response.data
  },

  async completeAdoption(id: string): Promise<Adoption> {
    const response = await api.post(`/adoptions/${id}/complete`)
    return response.data
  },

  async getStatistics(): Promise<AdoptionStatistics> {
    const response = await api.get('/adoptions/statistics')
    return response.data
  }
}
