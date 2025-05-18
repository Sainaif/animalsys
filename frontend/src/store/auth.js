import api from '../utils/api';

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
  SET_USER(state, { token, role, user }) {
    state.token = token;
    state.role = role;
    state.user = user;
    localStorage.setItem('token', token);
    localStorage.setItem('role', role);
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
      const res = await api.post('/auth/login', { username, password });
      const { token, role } = res.data.data;
      localStorage.setItem('token', token);
      localStorage.setItem('role', role);
      commit('SET_AUTH', { token, role });
      commit('SET_ERROR', null);
    } catch (err) {
      commit('SET_ERROR', err.response?.data?.message || 'Login failed');
    }
  },
  async register({ commit }, { username, email, password }) {
    try {
      const res = await api.post('/auth/register', { username, email, password });
      commit('SET_USER', { token: res.data.data.token, role: res.data.data.role, user: res.data.data.user });
      commit('SET_ERROR', null);
    } catch (err) {
      commit('SET_ERROR', err.response?.data?.message || 'Registration failed');
    }
  },
  logout({ commit }) {
    localStorage.removeItem('token');
    localStorage.removeItem('role');
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