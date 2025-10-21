import api from '../../utils/api'

const state = {
  adoptions: []
}

const mutations = {
  SET_ADOPTIONS(state, adoptions) {
    state.adoptions = adoptions
  },
  ADD_ADOPTION(state, adoption) {
    state.adoptions.push(adoption)
  },
  UPDATE_ADOPTION(state, updatedAdoption) {
    const index = state.adoptions.findIndex(a => a.id === updatedAdoption.id)
    if (index !== -1) {
      state.adoptions.splice(index, 1, updatedAdoption)
    }
  },
  DELETE_ADOPTION(state, id) {
    state.adoptions = state.adoptions.filter(a => a.id !== id)
  }
}

const actions = {
  async fetchAdoptions({ commit }) {
    const response = await api.get('/adoptions')
    commit('SET_ADOPTIONS', response.data.data)
  },
  async createAdoption({ commit }, adoption) {
    const response = await api.post('/adoptions', adoption)
    commit('ADD_ADOPTION', response.data.data)
  },
  async updateAdoption({ commit }, { id, data }) {
    await api.put(`/adoptions/${id}`, data)
    commit('UPDATE_ADOPTION', { id, ...data })
  },
  async deleteAdoption({ commit }, id) {
    await api.delete(`/adoptions/${id}`)
    commit('DELETE_ADOPTION', id)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
