import axios from '../utils/api';

const state = {
  documents: []
};

const mutations = {
  SET_DOCUMENTS(state, documents) {
    state.documents = documents;
  },
  ADD_DOCUMENT(state, document) {
    state.documents.push(document);
  }
};

const actions = {
  async fetchDocuments({ commit }) {
    const res = await axios.get('/documents');
    commit('SET_DOCUMENTS', res.data.data);
  },
  async uploadDocument({ commit }, { file, type }) {
    const formData = new FormData();
    formData.append('file', file);
    formData.append('type', type);
    const res = await axios.post('/documents', formData, {
      headers: { 'Content-Type': 'multipart/form-data' }
    });
    commit('ADD_DOCUMENT', res.data.data);
  }
};

const getters = {
  documents: state => state.documents
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}; 