import api from '../client'

/**
 * Reports API Module
 * Handles report generation, export, and scheduling
 */
export const reportsApi = {
  /**
   * Get list of available report types
   */
  getReportTypes(params = {}) {
    return api.get('/reports/types', { params })
  },

  /**
   * Get list of generated reports
   */
  list(params = {}) {
    return api.get('/reports', { params })
  },

  /**
   * Get report by ID
   */
  getById(id) {
    return api.get(`/reports/${id}`)
  },

  /**
   * Generate a new report
   */
  generate(reportData) {
    return api.post('/reports/generate', reportData)
  },

  /**
   * Delete a report
   */
  delete(id) {
    return api.delete(`/reports/${id}`)
  },

  /**
   * Export report in specified format
   */
  export(id, format) {
    return api.get(`/reports/${id}/export`, {
      params: { format },
      responseType: 'blob'
    })
  },

  /**
   * Get financial report
   */
  getFinancialReport(params = {}) {
    return api.get('/reports/financial', { params })
  },

  /**
   * Get adoption report
   */
  getAdoptionReport(params = {}) {
    return api.get('/reports/adoption', { params })
  },

  /**
   * Get volunteer report
   */
  getVolunteerReport(params = {}) {
    return api.get('/reports/volunteer', { params })
  },

  /**
   * Get inventory report
   */
  getInventoryReport(params = {}) {
    return api.get('/reports/inventory', { params })
  },

  /**
   * Get veterinary report
   */
  getVeterinaryReport(params = {}) {
    return api.get('/reports/veterinary', { params })
  },

  /**
   * Get campaign report
   */
  getCampaignReport(params = {}) {
    return api.get('/reports/campaign', { params })
  },

  /**
   * Get statutory report (for legal compliance)
   */
  getStatutoryReport(params = {}) {
    return api.get('/reports/statutory', { params })
  },

  /**
   * Get donor report
   */
  getDonorReport(params = {}) {
    return api.get('/reports/donor', { params })
  },

  /**
   * Get animal report
   */
  getAnimalReport(params = {}) {
    return api.get('/reports/animal', { params })
  },

  /**
   * Get custom report with custom parameters
   */
  getCustomReport(params = {}) {
    return api.post('/reports/custom', params)
  },

  /**
   * Schedule a report for automatic generation
   */
  schedule(scheduleData) {
    return api.post('/reports/schedule', scheduleData)
  },

  /**
   * Get scheduled reports
   */
  getScheduled(params = {}) {
    return api.get('/reports/scheduled', { params })
  },

  /**
   * Cancel a scheduled report
   */
  cancelSchedule(id) {
    return api.delete(`/reports/scheduled/${id}`)
  },

  /**
   * Get report statistics
   */
  getStatistics() {
    return api.get('/reports/statistics')
  },

  /**
   * Get recent reports
   */
  getRecent(params = {}) {
    return api.get('/reports/recent', { params })
  },

  /**
   * Get reports by type
   */
  getByType(type, params = {}) {
    return api.get('/reports/by-type', {
      params: { type, ...params }
    })
  },

  /**
   * Get reports by date range
   */
  getByDateRange(startDate, endDate, params = {}) {
    return api.get('/reports/date-range', {
      params: { start_date: startDate, end_date: endDate, ...params }
    })
  },

  /**
   * Validate report parameters
   */
  validateParameters(reportType, parameters) {
    return api.post('/reports/validate', { report_type: reportType, parameters })
  }
}
