import api from './api'
import type { PaginatedResponse, QueryParams } from '@/types/common'
import type { Donor, Donation, Campaign, FinanceStatistics } from '@/types/finance'

const mapPaginatedResponse = <T>(payload: any, key: string, params?: QueryParams): PaginatedResponse<T> => {
  const keyedCollection = Array.isArray(payload?.[key]) ? payload[key] : undefined
  const dataCollection = Array.isArray(payload?.data) ? payload.data : undefined
  const rawCollection = Array.isArray(payload) ? payload : undefined
  const collection: T[] = (keyedCollection ?? dataCollection ?? rawCollection ?? []) as T[]

  const total = typeof payload?.total === 'number' ? payload.total : collection.length
  const limit = typeof payload?.limit === 'number'
    ? payload.limit
    : params?.limit ?? collection.length
  const offset = typeof payload?.offset === 'number'
    ? payload.offset
    : params?.offset ?? 0

  return {
    data: collection,
    total,
    limit,
    offset
  }
}

const extractCollection = <T>(payload: any, key: string): T[] => {
  if (Array.isArray(payload?.[key])) {
    return payload[key]
  }
  if (Array.isArray(payload?.data)) {
    return payload.data
  }
  if (Array.isArray(payload)) {
    return payload
  }
  return []
}

export const financeService = {
  // Donors
  async getDonors(params?: QueryParams): Promise<PaginatedResponse<Donor>> {
    const response = await api.get('/donors', { params })
    return mapPaginatedResponse<Donor>(response.data, 'donors', params)
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
    return extractCollection<Donation>(response.data, 'donations')
  },

  // Donations
  async getDonations(params?: QueryParams): Promise<PaginatedResponse<Donation>> {
    const response = await api.get('/donations', { params })
    return mapPaginatedResponse<Donation>(response.data, 'donations', params)
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
    return mapPaginatedResponse<Campaign>(response.data, 'campaigns', params)
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
    return extractCollection<Donation>(response.data, 'donations')
  },

  // Statistics
  async getStatistics(): Promise<FinanceStatistics> {
    const response = await api.get('/finance/statistics')
    return response.data
  }
}
