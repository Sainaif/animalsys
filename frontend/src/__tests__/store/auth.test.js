import { describe, it, expect, beforeEach, vi } from 'vitest'
import { createStore } from 'vuex'
import auth from '../../store/modules/auth'

describe('Auth Store Module', () => {
  let store

  beforeEach(() => {
    localStorage.clear()
    store = createStore({
      modules: {
        auth
      }
    })
  })

  describe('getters', () => {
    it('isAuthenticated should return false when no token', () => {
      expect(store.getters['auth/isAuthenticated']).toBe(false)
    })

    it('isAuthenticated should return true when token exists', () => {
      store.commit('auth/SET_TOKEN', 'test-token')
      expect(store.getters['auth/isAuthenticated']).toBe(true)
    })

    it('user should return null initially', () => {
      expect(store.getters['auth/user']).toBe(null)
    })

    it('user should return user object when set', () => {
      const user = { id: '1', username: 'test', role: 'user' }
      store.commit('auth/SET_USER', user)
      expect(store.getters['auth/user']).toEqual(user)
    })
  })

  describe('mutations', () => {
    it('SET_TOKEN should set token in state and localStorage', () => {
      const token = 'test-token-123'
      store.commit('auth/SET_TOKEN', token)

      expect(store.state.auth.token).toBe(token)
      expect(localStorage.getItem('token')).toBe(token)
    })

    it('SET_TOKEN with null should remove token from state and localStorage', () => {
      localStorage.setItem('token', 'old-token')
      store.commit('auth/SET_TOKEN', null)

      expect(store.state.auth.token).toBe(null)
      expect(localStorage.getItem('token')).toBe(null)
    })

    it('SET_USER should set user in state and localStorage', () => {
      const user = { id: '1', username: 'test', role: 'admin' }
      store.commit('auth/SET_USER', user)

      expect(store.state.auth.user).toEqual(user)
      expect(JSON.parse(localStorage.getItem('user'))).toEqual(user)
    })

    it('SET_USER with null should remove user from state and localStorage', () => {
      localStorage.setItem('user', JSON.stringify({ id: '1' }))
      store.commit('auth/SET_USER', null)

      expect(store.state.auth.user).toBe(null)
      expect(localStorage.getItem('user')).toBe(null)
    })
  })

  describe('actions', () => {
    it('logout should clear token and user', () => {
      store.commit('auth/SET_TOKEN', 'token')
      store.commit('auth/SET_USER', { id: '1' })

      store.dispatch('auth/logout')

      expect(store.state.auth.token).toBe(null)
      expect(store.state.auth.user).toBe(null)
      expect(localStorage.getItem('token')).toBe(null)
      expect(localStorage.getItem('user')).toBe(null)
    })
  })
})
