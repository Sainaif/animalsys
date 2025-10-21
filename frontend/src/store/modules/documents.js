import api from '../../utils/api'

const state = {
  documents: []
}

const mutations = {
  SET_DOCUMENTS(state, documents) {
    state.documents = documents
  },
  ADD_DOCUMENT(state, document) {
    state.documents.push(document)
  },
  DELETE_DOCUMENT(state, id) {
    state.documents = state.documents.filter(d => d.id !== id)
  }
}

const actions = {
  async fetchDocuments({ commit }) {
    const response = await api.get('/documents')
    commit('SET_DOCUMENTS', response.data.data)
  },
  async createDocument({ commit }, document) {
    const response = await api.post('/documents', document)
    commit('ADD_DOCUMENT', response.data.data)
  },
  async deleteDocument({ commit }, id) {
    await api.delete(`/documents/${id}`)
    commit('DELETE_DOCUMENT', id)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
