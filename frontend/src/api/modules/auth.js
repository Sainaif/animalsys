import { api } from '../client'

export const authApi = {
  // Login
  login(credentials) {
    return api.post('/auth/login', credentials)
  },

  // Register
  register(userData) {
    return api.post('/auth/register', userData)
  },

  // Refresh token
  refresh(refreshToken) {
    return api.post('/auth/refresh', { refresh_token: refreshToken })
  },

  // Get profile
  getProfile() {
    return api.get('/auth/profile')
  },

  // Update profile
  updateProfile(profileData) {
    return api.put('/auth/profile', profileData)
  },

  // Change password
  changePassword(passwordData) {
    return api.post('/auth/change-password', passwordData)
  },
}

export default authApi
