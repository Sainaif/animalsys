import { api } from '../client'

export const financeApi = {
  // List transactions
  list(params = {}) {
    return api.get('/finance/transactions', { params })
  },

  // Get transaction by ID
  getById(id) {
    return api.get(`/finance/transactions/${id}`)
  },

  // Create transaction
  create(data) {
    return api.post('/finance/transactions', data)
  },

  // Update transaction
  update(id, data) {
    return api.put(`/finance/transactions/${id}`, data)
  },

  // Delete transaction
  delete(id) {
    return api.delete(`/finance/transactions/${id}`)
  },

  // Get financial report
  getReport(startDate, endDate) {
    return api.get('/finance/report', {
      params: {
        start_date: startDate,
        end_date: endDate
      }
    })
  },

  // Get dashboard statistics
  getDashboardStats() {
    return api.get('/finance/dashboard')
  },
}

export default financeApi
