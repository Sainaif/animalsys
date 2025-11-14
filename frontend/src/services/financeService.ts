import api from './api'
import type { PaginatedResponse, QueryParams } from '@/types/common'
import type { Donor, Donation, Campaign, FinanceStatistics } from '@/types/finance'

export const financeService = {
  // Donors
  async getDonors(params?: QueryParams): Promise<PaginatedResponse<Donor>> {
    const response = await api.get('/donors', { params })
    return response.data
  },

  async getDonor(id: string): Promise<Donor> {
    const response = await api.get(`/donors/${id}`)
    return response.data
  },

  async createDonor(data: Partial<Donor>): Promise<Donor> {
    const response = await api.post('/donors', data)
    return response.data
  },

  async updateDonor(id: string, data: Partial<Donor>): Promise<Donor> {
    const response = await api.put(`/donors/${id}`, data)
    return response.data
  },

  async deleteDonor(id: string): Promise<void> {
    await api.delete(`/donors/${id}`)
  },

  async getDonorDonations(donorId: string): Promise<Donation[]> {
    const response = await api.get(`/donors/${donorId}/donations`)
    return response.data
  },

  // Donations
  async getDonations(params?: QueryParams): Promise<PaginatedResponse<Donation>> {
    const response = await api.get('/donations', { params })
    return response.data
  },

  async getDonation(id: string): Promise<Donation> {
    const response = await api.get(`/donations/${id}`)
    return response.data
  },

  async createDonation(data: Partial<Donation>): Promise<Donation> {
    const response = await api.post('/donations', data)
    return response.data
  },

  async updateDonation(id: string, data: Partial<Donation>): Promise<Donation> {
    const response = await api.put(`/donations/${id}`, data)
    return response.data
  },

  async deleteDonation(id: string): Promise<void> {
    await api.delete(`/donations/${id}`)
  },

  async sendReceipt(id: string): Promise<void> {
    await api.post(`/donations/${id}/send-receipt`)
  },

  // Campaigns
  async getCampaigns(params?: QueryParams): Promise<PaginatedResponse<Campaign>> {
    const response = await api.get('/campaigns', { params })
    return response.data
  },

  async getCampaign(id: string): Promise<Campaign> {
    const response = await api.get(`/campaigns/${id}`)
    return response.data
  },

  async createCampaign(data: Partial<Campaign>): Promise<Campaign> {
    const response = await api.post('/campaigns', data)
    return response.data
  },

  async updateCampaign(id: string, data: Partial<Campaign>): Promise<Campaign> {
    const response = await api.put(`/campaigns/${id}`, data)
    return response.data
  },

  async deleteCampaign(id: string): Promise<void> {
    await api.delete(`/campaigns/${id}`)
  },

  async getCampaignDonations(campaignId: string): Promise<Donation[]> {
    const response = await api.get(`/campaigns/${campaignId}/donations`)
    return response.data
  },

  // Statistics
  async getStatistics(): Promise<FinanceStatistics> {
    const response = await api.get('/finance/statistics')
    return response.data
  }
}
