import api from '../client'

export const communicationsApi = {
  // Basic CRUD operations
  list(params = {}) {
    return api.get('/communications', { params })
  },

  getById(id) {
    return api.get(`/communications/${id}`)
  },

  create(communication) {
    return api.post('/communications', communication)
  },

  update(id, communication) {
    return api.put(`/communications/${id}`, communication)
  },

  delete(id) {
    return api.delete(`/communications/${id}`)
  },

  // Send communication
  send(id) {
    return api.post(`/communications/${id}/send`)
  },

  // Templates
  getTemplates(params = {}) {
    return api.get('/communications/templates', { params })
  },

  getTemplateById(id) {
    return api.get(`/communications/templates/${id}`)
  },

  createTemplate(template) {
    return api.post('/communications/templates', template)
  },

  updateTemplate(id, template) {
    return api.put(`/communications/templates/${id}`, template)
  },

  deleteTemplate(id) {
    return api.delete(`/communications/templates/${id}`)
  },

  // Recipients
  getRecipients(type) {
    return api.get('/communications/recipients', { params: { type } })
  },

  // Send bulk communication
  sendBulk(communication) {
    return api.post('/communications/send-bulk', communication)
  },

  // Get by status
  getByStatus(status) {
    return api.get('/communications/by-status', { params: { status } })
  },

  // Get by type
  getByType(type) {
    return api.get('/communications/by-type', { params: { type } })
  },

  // Statistics
  getStatistics() {
    return api.get('/communications/statistics')
  },

  // Schedule communication
  schedule(id, scheduledTime) {
    return api.post(`/communications/${id}/schedule`, { scheduled_time: scheduledTime })
  },

  // Cancel scheduled communication
  cancelSchedule(id) {
    return api.post(`/communications/${id}/cancel-schedule`)
  },

  // Get scheduled communications
  getScheduled() {
    return api.get('/communications/scheduled')
  },

  // Get communication history
  getHistory(params = {}) {
    return api.get('/communications/history', { params })
  }
}
