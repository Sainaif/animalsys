import api from '../client'

export const partnersApi = {
  // Basic CRUD operations
  list(params = {}) {
    return api.get('/partners', { params })
  },

  getById(id) {
    return api.get(`/partners/${id}`)
  },

  create(partner) {
    return api.post('/partners', partner)
  },

  update(id, partner) {
    return api.put(`/partners/${id}`, partner)
  },

  delete(id) {
    return api.delete(`/partners/${id}`)
  },

  // Agreement management
  getAgreements(id, params = {}) {
    return api.get(`/partners/${id}/agreements`, { params })
  },

  addAgreement(id, agreement) {
    return api.post(`/partners/${id}/agreements`, agreement)
  },

  updateAgreement(id, agreementId, agreement) {
    return api.put(`/partners/${id}/agreements/${agreementId}`, agreement)
  },

  deleteAgreement(id, agreementId) {
    return api.delete(`/partners/${id}/agreements/${agreementId}`)
  },

  // Statistics
  getStatistics(id) {
    return api.get(`/partners/${id}/statistics`)
  },

  // Active partners
  getActive() {
    return api.get('/partners/active')
  },

  // Partners by type
  getByType(type) {
    return api.get('/partners/by-type', { params: { type } })
  },

  // Overall statistics
  getOverallStatistics() {
    return api.get('/partners/statistics')
  }
}
