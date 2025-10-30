import { api } from '../client'

/**
 * Campaigns API module
 * Handles all campaign-related API calls
 */
export const campaignsApi = {
  /**
   * Get list of campaigns with optional filters
   * @param {Object} params - Query parameters (search, type, status, page, limit, sort)
   * @returns {Promise} API response with campaigns list
   */
  list(params = {}) {
    return api.get('/campaigns', { params })
  },

  /**
   * Get campaign by ID
   * @param {string} id - Campaign ID
   * @returns {Promise} API response with campaign details
   */
  getById(id) {
    return api.get(`/campaigns/${id}`)
  },

  /**
   * Create new campaign
   * @param {Object} data - Campaign data
   * @returns {Promise} API response with created campaign
   */
  create(data) {
    return api.post('/campaigns', data)
  },

  /**
   * Update existing campaign
   * @param {string} id - Campaign ID
   * @param {Object} data - Updated campaign data
   * @returns {Promise} API response with updated campaign
   */
  update(id, data) {
    return api.put(`/campaigns/${id}`, data)
  },

  /**
   * Delete campaign
   * @param {string} id - Campaign ID
   * @returns {Promise} API response
   */
  delete(id) {
    return api.delete(`/campaigns/${id}`)
  },

  /**
   * Get active campaigns
   * @returns {Promise} API response with active campaigns
   */
  getActive() {
    return api.get('/campaigns/active')
  },

  /**
   * Get campaign statistics
   * @param {string} id - Campaign ID
   * @returns {Promise} API response with campaign statistics
   */
  getStatistics(id) {
    return api.get(`/campaigns/${id}/statistics`)
  },

  /**
   * Update campaign progress
   * @param {string} id - Campaign ID
   * @param {Object} progress - Progress data
   * @returns {Promise} API response
   */
  updateProgress(id, progress) {
    return api.post(`/campaigns/${id}/progress`, progress)
  },

  /**
   * Get campaign milestones
   * @param {string} id - Campaign ID
   * @returns {Promise} API response with milestones
   */
  getMilestones(id) {
    return api.get(`/campaigns/${id}/milestones`)
  },

  /**
   * Add campaign milestone
   * @param {string} id - Campaign ID
   * @param {Object} milestone - Milestone data
   * @returns {Promise} API response
   */
  addMilestone(id, milestone) {
    return api.post(`/campaigns/${id}/milestones`, milestone)
  },

  /**
   * Get campaign donations
   * @param {string} id - Campaign ID
   * @param {Object} params - Query parameters
   * @returns {Promise} API response with donations
   */
  getDonations(id, params = {}) {
    return api.get(`/campaigns/${id}/donations`, { params })
  },

  /**
   * Get overall campaign statistics
   * @returns {Promise} API response with overall statistics
   */
  getOverallStatistics() {
    return api.get('/campaigns/statistics')
  },
}

export default campaignsApi
