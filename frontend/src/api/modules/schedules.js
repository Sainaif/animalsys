import api from '../client'

export const schedulesApi = {
  // Basic CRUD operations
  list(params = {}) {
    return api.get('/schedules', { params })
  },

  getById(id) {
    return api.get(`/schedules/${id}`)
  },

  create(schedule) {
    return api.post('/schedules', schedule)
  },

  update(id, schedule) {
    return api.put(`/schedules/${id}`, schedule)
  },

  delete(id) {
    return api.delete(`/schedules/${id}`)
  },

  // Volunteer schedules
  getVolunteerSchedules(volunteerId, params = {}) {
    return api.get(`/volunteers/${volunteerId}/schedules`, { params })
  },

  // Assign volunteer to schedule
  assignVolunteer(id, volunteerId) {
    return api.post(`/schedules/${id}/assign`, { volunteer_id: volunteerId })
  },

  // Unassign volunteer from schedule
  unassignVolunteer(id, volunteerId) {
    return api.post(`/schedules/${id}/unassign`, { volunteer_id: volunteerId })
  },

  // Swap requests
  getSwapRequests(id) {
    return api.get(`/schedules/${id}/swap-requests`)
  },

  createSwapRequest(id, swapRequest) {
    return api.post(`/schedules/${id}/swap-requests`, swapRequest)
  },

  approveSwapRequest(scheduleId, requestId) {
    return api.post(`/schedules/${scheduleId}/swap-requests/${requestId}/approve`)
  },

  rejectSwapRequest(scheduleId, requestId) {
    return api.post(`/schedules/${scheduleId}/swap-requests/${requestId}/reject`)
  },

  // Get schedules by date range
  getByDateRange(startDate, endDate) {
    return api.get('/schedules/date-range', {
      params: { start_date: startDate, end_date: endDate }
    })
  },

  // Get upcoming schedules
  getUpcoming(days = 7) {
    return api.get('/schedules/upcoming', { params: { days } })
  },

  // Get schedules by status
  getByStatus(status) {
    return api.get('/schedules/by-status', { params: { status } })
  },

  // Statistics
  getStatistics() {
    return api.get('/schedules/statistics')
  },

  // Check availability
  checkAvailability(volunteerId, startTime, endTime) {
    return api.post('/schedules/check-availability', {
      volunteer_id: volunteerId,
      start_time: startTime,
      end_time: endTime
    })
  }
}
