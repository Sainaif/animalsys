import api from '../utils/api';

const state = {
  adoptions: []
};

const mutations = {
  SET_ADOPTIONS(state, adoptions) {
    state.adoptions = adoptions;
  },
  ADD_ADOPTION(state, adoption) {
    state.adoptions.push(adoption);
  },
  UPDATE_ADOPTION(state, adoption) {
    const idx = state.adoptions.findIndex(a => a.id === adoption.id);
    if (idx !== -1) state.adoptions[idx] = adoption;
  }
};

const actions = {
  async fetchAdoptions({ commit }) {
    const res = await api.get('/adoptions');
    commit('SET_ADOPTIONS', res.data.data);
  },
  async applyForAdoption({ commit }, { animal_id, application_data }) {
    const res = await api.post('/adoptions', { animal_id, application_data });
    commit('ADD_ADOPTION', res.data.data);
  },
  async updateAdoptionStatus({ commit }, { id, status }) {
    await api.put(`/adoptions/${id}/status`, { status });
  }
};

const getters = {
  adoptions: state => state.adoptions
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
};