import api from './api'
import type { PaginatedResponse, QueryParams } from '@/types/common'
import type { Event, Volunteer, Shift, EventStatistics } from '@/types/event'

export const eventService = {
  // Events
  async getEvents(params?: QueryParams): Promise<PaginatedResponse<Event>> {
    const response = await api.get('/events', { params })
    return response.data
  },

  async getEvent(id: string): Promise<Event> {
    const response = await api.get(`/events/${id}`)
    return response.data
  },

  async createEvent(data: Partial<Event>): Promise<Event> {
    const response = await api.post('/events', data)
    return response.data
  },

  async updateEvent(id: string, data: Partial<Event>): Promise<Event> {
    const response = await api.put(`/events/${id}`, data)
    return response.data
  },

  async deleteEvent(id: string): Promise<void> {
    await api.delete(`/events/${id}`)
  },

  // Volunteers
  async getVolunteers(params?: QueryParams): Promise<PaginatedResponse<Volunteer>> {
    const response = await api.get('/volunteers', { params })
    return response.data
  },

  async getVolunteer(id: string): Promise<Volunteer> {
    const response = await api.get(`/volunteers/${id}`)
    return response.data
  },

  async createVolunteer(data: Partial<Volunteer>): Promise<Volunteer> {
    const response = await api.post('/volunteers', data)
    return response.data
  },

  async updateVolunteer(id: string, data: Partial<Volunteer>): Promise<Volunteer> {
    const response = await api.put(`/volunteers/${id}`, data)
    return response.data
  },

  async deleteVolunteer(id: string): Promise<void> {
    await api.delete(`/volunteers/${id}`)
  },

  async getVolunteerShifts(volunteerId: string): Promise<Shift[]> {
    const response = await api.get(`/volunteers/${volunteerId}/shifts`)
    return response.data
  },

  // Shifts
  async getShifts(params?: QueryParams): Promise<PaginatedResponse<Shift>> {
    const response = await api.get('/shifts', { params })
    return response.data
  },

  async getShift(id: string): Promise<Shift> {
    const response = await api.get(`/shifts/${id}`)
    return response.data
  },

  async createShift(data: Partial<Shift>): Promise<Shift> {
    const response = await api.post('/shifts', data)
    return response.data
  },

  async updateShift(id: string, data: Partial<Shift>): Promise<Shift> {
    const response = await api.put(`/shifts/${id}`, data)
    return response.data
  },

  async deleteShift(id: string): Promise<void> {
    await api.delete(`/shifts/${id}`)
  },

  async getEventShifts(eventId: string): Promise<Shift[]> {
    const response = await api.get(`/events/${eventId}/shifts`)
    return response.data
  },

  // Statistics
  async getStatistics(): Promise<EventStatistics> {
    const response = await api.get('/events/statistics')
    return response.data
  }
}
