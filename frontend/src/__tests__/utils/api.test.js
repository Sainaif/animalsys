import { describe, it, expect, beforeEach, vi } from 'vitest'
import api from '../../utils/api'

describe('API Client', () => {
  beforeEach(() => {
    localStorage.clear()
  })

  it('should have baseURL configured', () => {
    expect(api.defaults.baseURL).toBeDefined()
  })

  it('should have JSON content-type header', () => {
    expect(api.defaults.headers['Content-Type']).toBe('application/json')
  })

  it('should add Authorization header when token exists', () => {
    const token = 'test-token-123'
    localStorage.setItem('token', token)

    const config = { headers: {} }
    const requestInterceptor = api.interceptors.request.handlers[0].fulfilled

    const result = requestInterceptor(config)

    expect(result.headers.Authorization).toBe(`Bearer ${token}`)
  })

  it('should not add Authorization header when token does not exist', () => {
    localStorage.removeItem('token')

    const config = { headers: {} }
    const requestInterceptor = api.interceptors.request.handlers[0].fulfilled

    const result = requestInterceptor(config)

    expect(result.headers.Authorization).toBeUndefined()
  })
})
