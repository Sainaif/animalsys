import { api } from '../client'

/**
 * Donors API module
 * Handles all donor-related API calls
 */
export const donorsApi = {
  /**
   * Get list of donors with optional filters
   * @param {Object} params - Query parameters (search, type, page, limit, sort)
   * @returns {Promise} API response with donors list
   */
  list(params = {}) {
    return api.get('/donors', { params })
  },

  /**
   * Get donor by ID
   * @param {string} id - Donor ID
   * @returns {Promise} API response with donor details
   */
  getById(id) {
    return api.get(`/donors/${id}`)
  },

  /**
   * Create new donor
   * @param {Object} data - Donor data
   * @returns {Promise} API response with created donor
   */
  create(data) {
    return api.post('/donors', data)
  },

  /**
   * Update existing donor
   * @param {string} id - Donor ID
   * @param {Object} data - Updated donor data
   * @returns {Promise} API response with updated donor
   */
  update(id, data) {
    return api.put(`/donors/${id}`, data)
  },

  /**
   * Delete donor
   * @param {string} id - Donor ID
   * @returns {Promise} API response
   */
  delete(id) {
    return api.delete(`/donors/${id}`)
  },

  /**
   * Get donor's donation history
   * @param {string} id - Donor ID
   * @param {Object} params - Query parameters (start_date, end_date, page, limit)
   * @returns {Promise} API response with donation history
   */
  getDonations(id, params = {}) {
    return api.get(`/donors/${id}/donations`, { params })
  },

  /**
   * Add donation to donor
   * @param {string} id - Donor ID
   * @param {Object} donation - Donation data
   * @returns {Promise} API response
   */
  addDonation(id, donation) {
    return api.post(`/donors/${id}/donations`, donation)
  },

  /**
   * Get donor statistics
   * @param {string} id - Donor ID
   * @returns {Promise} API response with donor statistics
   */
  getStatistics(id) {
    return api.get(`/donors/${id}/statistics`)
  },

  /**
   * Update donor status
   * @param {string} id - Donor ID
   * @param {string} status - New status (active, inactive)
   * @returns {Promise} API response
   */
  updateStatus(id, status) {
    return api.put(`/donors/${id}/status`, { status })
  },
}

export default donorsApi
