import axios from 'axios'
import type { Partner, AnimalTransfer, PartnerAgreement } from '@/types/partner'

const API_URL = import.meta.env.VITE_API_URL || '/api/v1'

export const partnerService = {
  // Partners
  getPartners: (params?: any) => axios.get(`${API_URL}/partners`, { params }),
  getPartner: (id: number) => axios.get(`${API_URL}/partners/${id}`),
  createPartner: (data: Partner) => axios.post(`${API_URL}/partners`, data),
  updatePartner: (id: number, data: Partner) => axios.put(`${API_URL}/partners/${id}`, data),
  deletePartner: (id: number) => axios.delete(`${API_URL}/partners/${id}`),

  // Animal Transfers
  getAnimalTransfers: (params?: any) => axios.get(`${API_URL}/transfers`, { params }),
  getAnimalTransfer: (id: number) => axios.get(`${API_URL}/transfers/${id}`),
  createAnimalTransfer: (data: AnimalTransfer) => axios.post(`${API_URL}/transfers`, data),
  updateAnimalTransfer: (id: number, data: AnimalTransfer) => axios.put(`${API_URL}/transfers/${id}`, data),
  deleteAnimalTransfer: (id: number) => axios.delete(`${API_URL}/transfers/${id}`),

  // Partner Agreements (Note: May need backend implementation)
  getPartnerAgreements: (params?: any) => axios.get(`${API_URL}/partner-agreements`, { params }),
  getPartnerAgreement: (id: number) => axios.get(`${API_URL}/partner-agreements/${id}`),
  createPartnerAgreement: (data: PartnerAgreement) => axios.post(`${API_URL}/partner-agreements`, data),
  updatePartnerAgreement: (id: number, data: PartnerAgreement) => axios.put(`${API_URL}/partner-agreements/${id}`, data),
  deletePartnerAgreement: (id: number) => axios.delete(`${API_URL}/partner-agreements/${id}`)
}
