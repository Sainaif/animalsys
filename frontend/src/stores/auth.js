import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/services/api'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const accessToken = ref(localStorage.getItem('access_token') || null)
  const refreshToken = ref(localStorage.getItem('refresh_token') || null)

  // Getters
  const isAuthenticated = computed(() => !!accessToken.value)
  const userRole = computed(() => user.value?.role || null)
  const isAdmin = computed(() => {
    const role = userRole.value
    return role === 'super_admin' || role === 'admin'
  })
  const isSuperAdmin = computed(() => userRole.value === 'super_admin')

  // Actions
  async function login(email, password) {
    const response = await api.post('/auth/login', { email, password })
    const { access_token, refresh_token, user: userData } = response.data

    setTokens(access_token, refresh_token)
    user.value = userData

    return response.data
  }

  async function logout() {
    try {
      await api.post('/auth/logout')
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      clearAuth()
    }
  }

  async function refreshSession() {
    try {
      const response = await api.post('/auth/refresh', {
        refresh_token: refreshToken.value
      })

      const { access_token, refresh_token: newRefreshToken } = response.data

      setTokens(access_token, newRefreshToken || refreshToken.value)

      return response.data
    } catch (error) {
      clearAuth()
      throw error
    }
  }

  async function getCurrentUser() {
    try {
      const response = await api.get('/auth/me')
      user.value = response.data
      return response.data
    } catch (error) {
      clearAuth()
      throw error
    }
  }

  async function updateProfile(data) {
    const response = await api.put(`/users/${user.value.id}`, data)
    user.value = response.data
    return response.data
  }

  async function changePassword(oldPassword, newPassword) {
    await api.put('/auth/change-password', {
      old_password: oldPassword,
      new_password: newPassword
    })
  }

  function setTokens(access, refresh) {
    accessToken.value = access
    refreshToken.value = refresh
    localStorage.setItem('access_token', access)
    localStorage.setItem('refresh_token', refresh)
  }

  function clearAuth() {
    user.value = null
    accessToken.value = null
    refreshToken.value = null
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
  }

  // Initialize user if token exists
  if (accessToken.value) {
    getCurrentUser().catch(() => clearAuth())
  }

  return {
    // State
    user,
    accessToken,
    refreshToken,
    // Getters
    isAuthenticated,
    userRole,
    isAdmin,
    isSuperAdmin,
    // Actions
    login,
    logout,
    refreshSession,
    getCurrentUser,
    updateProfile,
    changePassword
  }
})
