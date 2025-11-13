import axios from 'axios'
import type { Partner, AnimalTransfer, PartnerAgreement } from '@/types/partner'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api'

export const partnerService = {
  // Partners
  getPartners: (params?: any) => axios.get(`${API_URL}/partners`, { params }),
  getPartner: (id: number) => axios.get(`${API_URL}/partners/${id}`),
  createPartner: (data: Partner) => axios.post(`${API_URL}/partners`, data),
  updatePartner: (id: number, data: Partner) => axios.put(`${API_URL}/partners/${id}`, data),
  deletePartner: (id: number) => axios.delete(`${API_URL}/partners/${id}`),

  // Animal Transfers
  getAnimalTransfers: (params?: any) => axios.get(`${API_URL}/animal-transfers`, { params }),
  getAnimalTransfer: (id: number) => axios.get(`${API_URL}/animal-transfers/${id}`),
  createAnimalTransfer: (data: AnimalTransfer) => axios.post(`${API_URL}/animal-transfers`, data),
  updateAnimalTransfer: (id: number, data: AnimalTransfer) => axios.put(`${API_URL}/animal-transfers/${id}`, data),
  deleteAnimalTransfer: (id: number) => axios.delete(`${API_URL}/animal-transfers/${id}`),

  // Partner Agreements
  getPartnerAgreements: (params?: any) => axios.get(`${API_URL}/partner-agreements`, { params }),
  getPartnerAgreement: (id: number) => axios.get(`${API_URL}/partner-agreements/${id}`),
  createPartnerAgreement: (data: PartnerAgreement) => axios.post(`${API_URL}/partner-agreements`, data),
  updatePartnerAgreement: (id: number, data: PartnerAgreement) => axios.put(`${API_URL}/partner-agreements/${id}`, data),
  deletePartnerAgreement: (id: number) => axios.delete(`${API_URL}/partner-agreements/${id}`)
}
