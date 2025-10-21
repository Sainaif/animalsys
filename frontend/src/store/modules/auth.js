import api from '../../utils/api'

const state = () => ({
  token: localStorage.getItem('token') || null,
  user: JSON.parse(localStorage.getItem('user')) || null
})

const getters = {
  isAuthenticated: state => !!state.token,
  user: state => state.user
}

const mutations = {
  SET_TOKEN(state, token) {
    state.token = token
    if (token) {
      localStorage.setItem('token', token)
    } else {
      localStorage.removeItem('token')
    }
  },
  SET_USER(state, user) {
    state.user = user
    if (user) {
      localStorage.setItem('user', JSON.stringify(user))
    } else {
      localStorage.removeItem('user')
    }
  }
}

const actions = {
  async login({ commit }, credentials) {
    try {
      const response = await api.post('/auth/login', credentials)
      const { token, user } = response.data.data
      commit('SET_TOKEN', token)
      commit('SET_USER', user)
      return response.data
    } catch (error) {
      throw error
    }
  },
  async register({ commit }, userData) {
    try {
      const response = await api.post('/auth/register', userData)
      return response.data
    } catch (error) {
      throw error
    }
  },
  logout({ commit }) {
    commit('SET_TOKEN', null)
    commit('SET_USER', null)
  }
}

export default {
  namespaced: true,
  state,
  getters,
  mutations,
  actions
}
