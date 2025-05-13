import axios from '../utils/api';

const state = {
  token: localStorage.getItem('token') || '',
  user: null,
  role: localStorage.getItem('role') || '',
  error: null
};

const mutations = {
  SET_AUTH(state, { token, role }) {
    state.token = token;
    state.role = role;
    localStorage.setItem('token', token);
    localStorage.setItem('role', role);
  },
  SET_USER(state, user) {
    state.user = user;
  },
  SET_ERROR(state, error) {
    state.error = error;
  },
  LOGOUT(state) {
    state.token = '';
    state.role = '';
    state.user = null;
    localStorage.removeItem('token');
    localStorage.removeItem('role');
  }
};

const actions = {
  async login({ commit }, { username, password }) {
    try {
      const res = await axios.post('/auth/login', { username, password });
      commit('SET_AUTH', res.data.data);
      commit('SET_ERROR', null);
      return true;
    } catch (e) {
      commit('SET_ERROR', e.response?.data?.error || 'Błąd logowania');
      return false;
    }
  },
  async register({ commit }, { username, email, password }) {
    try {
      const res = await axios.post('/auth/register', { username, email, password });
      commit('SET_AUTH', res.data.data);
      commit('SET_ERROR', null);
      return true;
    } catch (e) {
      commit('SET_ERROR', e.response?.data?.error || 'Błąd rejestracji');
      return false;
    }
  },
  logout({ commit }) {
    commit('LOGOUT');
  }
};

const getters = {
  isAuthenticated: state => !!state.token,
  userRole: state => state.role,
  authError: state => state.error
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}; 