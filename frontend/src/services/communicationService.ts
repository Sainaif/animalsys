import axios from 'axios'
import type { EmailTemplate, EmailCampaign, CommunicationLog, SMSTemplate } from '@/types/communication'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api'

export const communicationService = {
  // Email Templates
  getEmailTemplates: (params?: any) => axios.get(`${API_URL}/email-templates`, { params }),
  getEmailTemplate: (id: number) => axios.get(`${API_URL}/email-templates/${id}`),
  createEmailTemplate: (data: EmailTemplate) => axios.post(`${API_URL}/email-templates`, data),
  updateEmailTemplate: (id: number, data: EmailTemplate) => axios.put(`${API_URL}/email-templates/${id}`, data),
  deleteEmailTemplate: (id: number) => axios.delete(`${API_URL}/email-templates/${id}`),

  // Email Campaigns
  getEmailCampaigns: (params?: any) => axios.get(`${API_URL}/email-campaigns`, { params }),
  getEmailCampaign: (id: number) => axios.get(`${API_URL}/email-campaigns/${id}`),
  createEmailCampaign: (data: EmailCampaign) => axios.post(`${API_URL}/email-campaigns`, data),
  updateEmailCampaign: (id: number, data: EmailCampaign) => axios.put(`${API_URL}/email-campaigns/${id}`, data),
  deleteEmailCampaign: (id: number) => axios.delete(`${API_URL}/email-campaigns/${id}`),
  sendEmailCampaign: (id: number) => axios.post(`${API_URL}/email-campaigns/${id}/send`),

  // Communication Logs
  getCommunicationLogs: (params?: any) => axios.get(`${API_URL}/communication-logs`, { params }),
  getCommunicationLog: (id: number) => axios.get(`${API_URL}/communication-logs/${id}`),
  createCommunicationLog: (data: CommunicationLog) => axios.post(`${API_URL}/communication-logs`, data),
  updateCommunicationLog: (id: number, data: CommunicationLog) => axios.put(`${API_URL}/communication-logs/${id}`, data),
  deleteCommunicationLog: (id: number) => axios.delete(`${API_URL}/communication-logs/${id}`),

  // SMS Templates
  getSMSTemplates: (params?: any) => axios.get(`${API_URL}/sms-templates`, { params }),
  getSMSTemplate: (id: number) => axios.get(`${API_URL}/sms-templates/${id}`),
  createSMSTemplate: (data: SMSTemplate) => axios.post(`${API_URL}/sms-templates`, data),
  updateSMSTemplate: (id: number, data: SMSTemplate) => axios.put(`${API_URL}/sms-templates/${id}`, data),
  deleteSMSTemplate: (id: number) => axios.delete(`${API_URL}/sms-templates/${id}`)
}
