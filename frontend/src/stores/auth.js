import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { api } from '../api/client'
import router from '../router'

const ROLE_HIERARCHY = {
  'super_admin': 6,
  'admin': 5,
  'employee': 4,
  'volunteer': 3,
  'user': 2,
  'guest': 1,
}

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const accessToken = ref(localStorage.getItem('access_token') || null)
  const refreshToken = ref(localStorage.getItem('refresh_token') || null)
  const loading = ref(false)
  const error = ref(null)

  // Getters
  const isAuthenticated = computed(() => !!accessToken.value && !!user.value)
  const userRole = computed(() => user.value?.role || 'guest')
  const userRoleLevel = computed(() => ROLE_HIERARCHY[userRole.value] || 0)

  // Actions
  async function login(credentials) {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/login', credentials)

      accessToken.value = response.data.access_token
      refreshToken.value = response.data.refresh_token
      user.value = response.data.user

      // Save to localStorage
      localStorage.setItem('access_token', accessToken.value)
      localStorage.setItem('refresh_token', refreshToken.value)
      localStorage.setItem('user', JSON.stringify(user.value))

      // Set default auth header
      api.defaults.headers.common['Authorization'] = `Bearer ${accessToken.value}`

      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Login failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function register(userData) {
    loading.value = true
    error.value = null

    try {
      const response = await api.post('/auth/register', userData)
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Registration failed'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    // Clear state
    user.value = null
    accessToken.value = null
    refreshToken.value = null

    // Clear localStorage
    localStorage.removeItem('access_token')
    localStorage.removeItem('refresh_token')
    localStorage.removeItem('user')

    // Clear auth header
    delete api.defaults.headers.common['Authorization']

    // Redirect to login
    router.push({ name: 'login' })
  }

  async function refreshAccessToken() {
    if (!refreshToken.value) {
      throw new Error('No refresh token available')
    }

    try {
      const response = await api.post('/auth/refresh', {
        refresh_token: refreshToken.value
      })

      accessToken.value = response.data.access_token
      localStorage.setItem('access_token', accessToken.value)
      api.defaults.headers.common['Authorization'] = `Bearer ${accessToken.value}`

      return accessToken.value
    } catch (err) {
      // If refresh fails, logout
      await logout()
      throw err
    }
  }

  async function fetchProfile() {
    try {
      const response = await api.get('/auth/profile')
      user.value = response.data
      localStorage.setItem('user', JSON.stringify(user.value))
      return user.value
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to fetch profile'
      throw err
    }
  }

  async function updateProfile(userData) {
    loading.value = true
    error.value = null

    try {
      const response = await api.put(`/users/${user.value.user_id}`, userData)
      user.value = { ...user.value, ...response.data }
      localStorage.setItem('user', JSON.stringify(user.value))
      return user.value
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to update profile'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function changePassword(passwordData) {
    loading.value = true
    error.value = null

    try {
      await api.post('/auth/change-password', passwordData)
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to change password'
      throw err
    } finally {
      loading.value = false
    }
  }

  async function initAuth() {
    // Try to restore session from localStorage
    const storedToken = localStorage.getItem('access_token')
    const storedUser = localStorage.getItem('user')

    if (storedToken && storedUser) {
      accessToken.value = storedToken
      user.value = JSON.parse(storedUser)
      api.defaults.headers.common['Authorization'] = `Bearer ${storedToken}`

      // Fetch fresh profile to verify token
      try {
        await fetchProfile()
      } catch (err) {
        // Token might be expired, try refresh
        try {
          await refreshAccessToken()
          await fetchProfile()
        } catch {
          // Both failed, logout
          await logout()
        }
      }
    }
  }

  function hasRole(requiredRole) {
    const requiredLevel = ROLE_HIERARCHY[requiredRole] || 0
    return userRoleLevel.value >= requiredLevel
  }

  function can(permission) {
    // Permission checking logic
    // For now, simple role-based
    return hasRole(permission)
  }

  return {
    // State
    user,
    accessToken,
    refreshToken,
    loading,
    error,

    // Getters
    isAuthenticated,
    userRole,
    userRoleLevel,

    // Actions
    login,
    register,
    logout,
    refreshAccessToken,
    fetchProfile,
    updateProfile,
    changePassword,
    initAuth,
    hasRole,
    can,
  }
})
