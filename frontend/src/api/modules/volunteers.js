import { api } from '../client'

export const volunteersApi = {
  // List volunteers
  list(params = {}) {
    return api.get('/volunteers', { params })
  },

  // Get volunteer by ID
  getById(id) {
    return api.get(`/volunteers/${id}`)
  },

  // Create volunteer
  create(data) {
    return api.post('/volunteers', data)
  },

  // Update volunteer
  update(id, data) {
    return api.put(`/volunteers/${id}`, data)
  },

  // Delete volunteer
  delete(id) {
    return api.delete(`/volunteers/${id}`)
  },

  // Add training
  addTraining(id, training) {
    return api.post(`/volunteers/${id}/trainings`, training)
  },

  // Log hours
  logHours(id, hours) {
    return api.post(`/volunteers/${id}/hours`, hours)
  },

  // Get volunteer statistics
  getStatistics(id) {
    return api.get(`/volunteers/${id}/statistics`)
  },

  // Update status
  updateStatus(id, status) {
    return api.put(`/volunteers/${id}/status`, { status })
  },
}

export default volunteersApi
