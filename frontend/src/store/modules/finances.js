import api from '../../utils/api'

const state = {
  finances: []
}

const mutations = {
  SET_FINANCES(state, finances) {
    state.finances = finances
  },
  ADD_FINANCE(state, finance) {
    state.finances.push(finance)
  },
  UPDATE_FINANCE(state, updatedFinance) {
    const index = state.finances.findIndex(f => f.id === updatedFinance.id)
    if (index !== -1) {
      state.finances.splice(index, 1, updatedFinance)
    }
  },
  DELETE_FINANCE(state, id) {
    state.finances = state.finances.filter(f => f.id !== id)
  }
}

const actions = {
  async fetchFinances({ commit }) {
    const response = await api.get('/finances')
    commit('SET_FINANCES', response.data.data)
  },
  async createFinance({ commit }, finance) {
    const response = await api.post('/finances', finance)
    commit('ADD_FINANCE', response.data.data)
  },
  async updateFinance({ commit }, { id, data }) {
    await api.put(`/finances/${id}`, data)
    commit('UPDATE_FINANCE', { id, ...data })
  },
  async deleteFinance({ commit }, id) {
    await api.delete(`/finances/${id}`)
    commit('DELETE_FINANCE', id)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
