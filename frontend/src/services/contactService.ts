import api from './api'
import type { PaginatedResponse } from '@/types/common'
import type { Contact, ContactFilters, ContactPayload } from '@/types/contact'

const toPaginatedResponse = (payload: any): PaginatedResponse<Contact> => {
  const list = Array.isArray(payload?.contacts)
    ? payload.contacts
    : Array.isArray(payload?.data)
      ? payload.data
      : Array.isArray(payload)
        ? payload
        : []

  return {
    data: list,
    total: typeof payload?.total === 'number' ? payload.total : list.length,
    limit: typeof payload?.limit === 'number' ? payload.limit : list.length,
    offset: typeof payload?.offset === 'number' ? payload.offset : 0
  }
}

export const contactService = {
  async getContacts(params?: ContactFilters): Promise<PaginatedResponse<Contact>> {
    const response = await api.get('/contacts', { params })
    return toPaginatedResponse(response.data)
  },

  async getContact(id: string): Promise<Contact> {
    const response = await api.get(`/contacts/${id}`)
    return response.data
  },

  async createContact(payload: ContactPayload): Promise<Contact> {
    const response = await api.post('/contacts', payload)
    return response.data
  },

  async updateContact(id: string, payload: Partial<ContactPayload>): Promise<Contact> {
    const response = await api.put(`/contacts/${id}`, payload)
    return response.data
  },

  async deleteContact(id: string): Promise<void> {
    await api.delete(`/contacts/${id}`)
  }
}
