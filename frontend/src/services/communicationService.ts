import axios from 'axios'
import type { EmailTemplate, EmailCampaign, CommunicationLog, SMSTemplate } from '@/types/communication'

const API_URL = import.meta.env.VITE_API_URL || '/api/v1'

export const communicationService = {
  // Email Templates (uses backend /v1/templates endpoint)
  getEmailTemplates: (params?: any) => axios.get(`${API_URL}/templates`, { params }),
  getEmailTemplate: (id: number) => axios.get(`${API_URL}/templates/${id}`),
  createEmailTemplate: (data: EmailTemplate) => axios.post(`${API_URL}/templates`, data),
  updateEmailTemplate: (id: number, data: EmailTemplate) => axios.put(`${API_URL}/templates/${id}`, data),
  deleteEmailTemplate: (id: number) => axios.delete(`${API_URL}/templates/${id}`),

  // Email Campaigns (Note: May need backend implementation for bulk email campaigns)
  // Currently mapped to communications endpoint - bulk sending may need additional backend work
  getEmailCampaigns: (params?: any) => axios.get(`${API_URL}/communications`, { params }),
  getEmailCampaign: (id: number) => axios.get(`${API_URL}/communications/${id}`),
  createEmailCampaign: (data: EmailCampaign) => axios.post(`${API_URL}/communications`, data),
  updateEmailCampaign: (id: number, data: EmailCampaign) => axios.put(`${API_URL}/communications/${id}`, data),
  deleteEmailCampaign: (id: number) => axios.delete(`${API_URL}/communications/${id}`),
  sendEmailCampaign: (id: number) => axios.post(`${API_URL}/communications/${id}/send`),

  // Communication Logs
  getCommunicationLogs: (params?: any) => axios.get(`${API_URL}/communications`, { params }),
  getCommunicationLog: (id: number) => axios.get(`${API_URL}/communications/${id}`),
  createCommunicationLog: (data: CommunicationLog) => axios.post(`${API_URL}/communications`, data),
  updateCommunicationLog: (id: number, data: CommunicationLog) => axios.put(`${API_URL}/communications/${id}`, data),
  deleteCommunicationLog: (id: number) => axios.delete(`${API_URL}/communications/${id}`),

  // SMS Templates (uses same templates endpoint, filter by type on frontend)
  getSMSTemplates: (params?: any) => axios.get(`${API_URL}/templates`, { params: { ...params, template_type: 'sms' } }),
  getSMSTemplate: (id: number) => axios.get(`${API_URL}/templates/${id}`),
  createSMSTemplate: (data: SMSTemplate) => axios.post(`${API_URL}/templates`, data),
  updateSMSTemplate: (id: number, data: SMSTemplate) => axios.put(`${API_URL}/templates/${id}`, data),
  deleteSMSTemplate: (id: number) => axios.delete(`${API_URL}/templates/${id}`)
}
