import api from '../../utils/api'

const state = {
  users: []
}

const mutations = {
  SET_USERS(state, users) {
    state.users = users
  },
  UPDATE_USER(state, updatedUser) {
    const index = state.users.findIndex(u => u.id === updatedUser.id)
    if (index !== -1) {
      state.users.splice(index, 1, updatedUser)
    }
  },
  DELETE_USER(state, id) {
    state.users = state.users.filter(u => u.id !== id)
  }
}

const actions = {
  async fetchUsers({ commit }) {
    const response = await api.get('/users')
    commit('SET_USERS', response.data.data)
  },
  async updateUser({ commit }, { id, data }) {
    await api.put(`/users/${id}`, data)
    commit('UPDATE_USER', { id, ...data })
  },
  async deleteUser({ commit }, id) {
    await api.delete(`/users/${id}`)
    commit('DELETE_USER', id)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
