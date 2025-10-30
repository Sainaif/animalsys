import api from '../client'

export const documentsApi = {
  // Basic CRUD operations
  list(params = {}) {
    return api.get('/documents', { params })
  },

  getById(id) {
    return api.get(`/documents/${id}`)
  },

  create(document) {
    return api.post('/documents', document)
  },

  update(id, document) {
    return api.put(`/documents/${id}`, document)
  },

  delete(id) {
    return api.delete(`/documents/${id}`)
  },

  // File upload
  upload(formData) {
    return api.post('/documents/upload', formData, {
      headers: {
        'Content-Type': 'multipart/form-data'
      }
    })
  },

  // File download
  download(id) {
    return api.get(`/documents/${id}/download`, {
      responseType: 'blob'
    })
  },

  // Get documents by entity
  getByEntity(entityType, entityId) {
    return api.get('/documents/by-entity', {
      params: { entity_type: entityType, entity_id: entityId }
    })
  },

  // Get documents by type
  getByType(type) {
    return api.get('/documents/by-type', { params: { type } })
  },

  // Get documents by category
  getByCategory(category) {
    return api.get('/documents/by-category', { params: { category } })
  },

  // Search documents
  search(query) {
    return api.get('/documents/search', { params: { q: query } })
  },

  // Statistics
  getStatistics() {
    return api.get('/documents/statistics')
  },

  // Get recent documents
  getRecent(limit = 10) {
    return api.get('/documents/recent', { params: { limit } })
  },

  // Get expiring documents (for contracts, certificates, etc.)
  getExpiring(days = 30) {
    return api.get('/documents/expiring', { params: { days } })
  }
}
