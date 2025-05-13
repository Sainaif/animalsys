import axios from '../utils/api';

const state = {
  finances: []
};

const mutations = {
  SET_FINANCES(state, finances) {
    state.finances = finances;
  },
  ADD_FINANCE(state, finance) {
    state.finances.push(finance);
  }
};

const actions = {
  async fetchFinances({ commit }) {
    const res = await axios.get('/finances');
    commit('SET_FINANCES', res.data.data);
  },
  async addFinance({ commit }, finance) {
    const res = await axios.post('/finances', finance);
    commit('ADD_FINANCE', res.data.data);
  },
  async exportCsv() {
    const res = await axios.get('/finances/report/csv', { responseType: 'blob' });
    const url = window.URL.createObjectURL(new Blob([res.data]));
    const link = document.createElement('a');
    link.href = url;
    link.setAttribute('download', 'finances.csv');
    document.body.appendChild(link);
    link.click();
    link.remove();
  }
};

const getters = {
  finances: state => state.finances
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}; 