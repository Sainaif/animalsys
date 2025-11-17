import axios from 'axios'
import { useAuthStore } from '@/stores/auth'
import router from '@/router'

const resolveBaseURL = () => {
  if (import.meta.env.VITE_API_URL) {
    return import.meta.env.VITE_API_URL
  }

  if (typeof window !== 'undefined') {
    const { protocol, hostname, port } = window.location
    if (port === '23001') {
      return `${protocol}//${hostname}:23000/api/v1`
    }
    return `${protocol}//${hostname}${port ? `:${port}` : ''}/api/v1`
  }

  return '/api/v1'
}

// Create axios instance
const api = axios.create({
  baseURL: resolveBaseURL(),
  timeout: 30000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// Request interceptor
api.interceptors.request.use(
  (config) => {
    const authStore = useAuthStore()

    // Add auth token to requests
    if (authStore.accessToken) {
      config.headers.Authorization = `Bearer ${authStore.accessToken}`
    }

    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Response interceptor
api.interceptors.response.use(
  (response) => {
    return response
  },
  async (error) => {
    const originalRequest = error.config

    // If error is 401 and we haven't retried yet
    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true

      const authStore = useAuthStore()

      try {
        // Try to refresh the token
        await authStore.refreshSession()

        // Retry the original request with new token
        originalRequest.headers.Authorization = `Bearer ${authStore.accessToken}`
        return api(originalRequest)
      } catch (refreshError) {
        // Refresh failed, logout user
        authStore.logout()
        router.push({ name: 'login' })
        return Promise.reject(refreshError)
      }
    }

    return Promise.reject(error)
  }
)

export default api
