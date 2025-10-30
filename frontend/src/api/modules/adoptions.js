import { api } from '../client'

export const adoptionsApi = {
  // List adoption applications
  list(params = {}) {
    return api.get('/adoptions', { params })
  },

  // Get adoption by ID
  getById(id) {
    return api.get(`/adoptions/${id}`)
  },

  // Create adoption application
  create(data) {
    return api.post('/adoptions', data)
  },

  // Update adoption
  update(id, data) {
    return api.put(`/adoptions/${id}`, data)
  },

  // Update adoption status
  updateStatus(id, status) {
    return api.put(`/adoptions/${id}/status`, { status })
  },

  // Schedule interview
  scheduleInterview(id, data) {
    return api.put(`/adoptions/${id}/interview`, data)
  },

  // Approve adoption
  approve(id) {
    return api.put(`/adoptions/${id}/approve`)
  },

  // Reject adoption
  reject(id, reason) {
    return api.put(`/adoptions/${id}/reject`, { reason })
  },

  // Complete adoption
  complete(id, data) {
    return api.put(`/adoptions/${id}/complete`, data)
  },

  // Upload contract
  uploadContract(id, contractUrl) {
    return api.post(`/adoptions/${id}/contract`, { contract_url: contractUrl })
  },

  // Add follow-up
  addFollowUp(id, followUp) {
    return api.post(`/adoptions/${id}/follow-ups`, followUp)
  },
}

export default adoptionsApi
