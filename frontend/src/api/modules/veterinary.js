import { api } from '../client'

/**
 * Veterinary API module
 * Handles all veterinary-related API calls
 */
export const veterinaryApi = {
  /**
   * Get list of veterinary visits with optional filters
   * @param {Object} params - Query parameters (search, animal_id, type, status, page, limit, sort)
   * @returns {Promise} API response with visits list
   */
  list(params = {}) {
    return api.get('/veterinary/visits', { params })
  },

  /**
   * Get veterinary visit by ID
   * @param {string} id - Visit ID
   * @returns {Promise} API response with visit details
   */
  getById(id) {
    return api.get(`/veterinary/visits/${id}`)
  },

  /**
   * Create new veterinary visit
   * @param {Object} data - Visit data
   * @returns {Promise} API response with created visit
   */
  create(data) {
    return api.post('/veterinary/visits', data)
  },

  /**
   * Update existing veterinary visit
   * @param {string} id - Visit ID
   * @param {Object} data - Updated visit data
   * @returns {Promise} API response with updated visit
   */
  update(id, data) {
    return api.put(`/veterinary/visits/${id}`, data)
  },

  /**
   * Delete veterinary visit
   * @param {string} id - Visit ID
   * @returns {Promise} API response
   */
  delete(id) {
    return api.delete(`/veterinary/visits/${id}`)
  },

  /**
   * Get visits for specific animal
   * @param {string} animalId - Animal ID
   * @param {Object} params - Query parameters
   * @returns {Promise} API response with animal's visits
   */
  getAnimalVisits(animalId, params = {}) {
    return api.get(`/animals/${animalId}/veterinary/visits`, { params })
  },

  /**
   * Get vaccinations for specific animal
   * @param {string} animalId - Animal ID
   * @returns {Promise} API response with vaccinations
   */
  getVaccinations(animalId) {
    return api.get(`/animals/${animalId}/vaccinations`)
  },

  /**
   * Add vaccination record
   * @param {string} animalId - Animal ID
   * @param {Object} vaccination - Vaccination data
   * @returns {Promise} API response
   */
  addVaccination(animalId, vaccination) {
    return api.post(`/animals/${animalId}/vaccinations`, vaccination)
  },

  /**
   * Get medications for specific animal
   * @param {string} animalId - Animal ID
   * @returns {Promise} API response with medications
   */
  getMedications(animalId) {
    return api.get(`/animals/${animalId}/medications`)
  },

  /**
   * Add medication record
   * @param {string} animalId - Animal ID
   * @param {Object} medication - Medication data
   * @returns {Promise} API response
   */
  addMedication(animalId, medication) {
    return api.post(`/animals/${animalId}/medications`, medication)
  },

  /**
   * Get upcoming vaccinations/checkups
   * @param {number} days - Number of days to look ahead (default 30)
   * @returns {Promise} API response with upcoming items
   */
  getUpcoming(days = 30) {
    return api.get('/veterinary/upcoming', { params: { days } })
  },

  /**
   * Get veterinary statistics
   * @returns {Promise} API response with veterinary statistics
   */
  getStatistics() {
    return api.get('/veterinary/statistics')
  },

  /**
   * Get list of veterinarians
   * @returns {Promise} API response with veterinarians list
   */
  getVeterinarians() {
    return api.get('/veterinary/veterinarians')
  },

  /**
   * Get list of clinics
   * @returns {Promise} API response with clinics list
   */
  getClinics() {
    return api.get('/veterinary/clinics')
  },
}

export default veterinaryApi
