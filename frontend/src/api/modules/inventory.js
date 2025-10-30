import { api } from '../client'

/**
 * Inventory API module
 * Handles all inventory-related API calls
 */
export const inventoryApi = {
  /**
   * Get list of inventory items with optional filters
   * @param {Object} params - Query parameters (search, category, status, page, limit, sort)
   * @returns {Promise} API response with inventory items list
   */
  list(params = {}) {
    return api.get('/inventory', { params })
  },

  /**
   * Get inventory item by ID
   * @param {string} id - Inventory item ID
   * @returns {Promise} API response with inventory item details
   */
  getById(id) {
    return api.get(`/inventory/${id}`)
  },

  /**
   * Create new inventory item
   * @param {Object} data - Inventory item data
   * @returns {Promise} API response with created inventory item
   */
  create(data) {
    return api.post('/inventory', data)
  },

  /**
   * Update existing inventory item
   * @param {string} id - Inventory item ID
   * @param {Object} data - Updated inventory item data
   * @returns {Promise} API response with updated inventory item
   */
  update(id, data) {
    return api.put(`/inventory/${id}`, data)
  },

  /**
   * Delete inventory item
   * @param {string} id - Inventory item ID
   * @returns {Promise} API response
   */
  delete(id) {
    return api.delete(`/inventory/${id}`)
  },

  /**
   * Get inventory item movements/history
   * @param {string} id - Inventory item ID
   * @param {Object} params - Query parameters (start_date, end_date, page, limit)
   * @returns {Promise} API response with movements history
   */
  getMovements(id, params = {}) {
    return api.get(`/inventory/${id}/movements`, { params })
  },

  /**
   * Add stock movement (in/out/adjustment)
   * @param {string} id - Inventory item ID
   * @param {Object} movement - Movement data (type, quantity, reason, notes)
   * @returns {Promise} API response
   */
  addMovement(id, movement) {
    return api.post(`/inventory/${id}/movements`, movement)
  },

  /**
   * Get inventory statistics
   * @returns {Promise} API response with inventory statistics
   */
  getStatistics() {
    return api.get('/inventory/statistics')
  },

  /**
   * Get low stock items
   * @returns {Promise} API response with low stock items
   */
  getLowStock() {
    return api.get('/inventory/low-stock')
  },

  /**
   * Get expired or expiring soon items
   * @param {number} days - Number of days to check (default 30)
   * @returns {Promise} API response with expiring items
   */
  getExpiring(days = 30) {
    return api.get('/inventory/expiring', { params: { days } })
  },
}

export default inventoryApi
