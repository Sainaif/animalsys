import { api } from '../client'

export const animalsApi = {
  // Get available animals for adoption (public)
  getAvailable() {
    return api.get('/animals/available')
  },

  // Get animal by ID
  getById(id) {
    return api.get(`/animals/${id}`)
  },

  // List animals with filters
  list(params = {}) {
    return api.get('/animals', { params })
  },

  // Create animal
  create(data) {
    return api.post('/animals', data)
  },

  // Update animal
  update(id, data) {
    return api.put(`/animals/${id}`, data)
  },

  // Delete animal
  delete(id) {
    return api.delete(`/animals/${id}`)
  },

  // Add medical record
  addMedicalRecord(id, record) {
    return api.post(`/animals/${id}/medical-records`, record)
  },

  // Add photo
  addPhoto(id, photoUrl) {
    return api.post(`/animals/${id}/photos`, { photo_url: photoUrl })
  },
}

export default animalsApi
