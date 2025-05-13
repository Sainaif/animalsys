import axios from '../utils/api';

const state = {
  schedules: []
};

const mutations = {
  SET_SCHEDULES(state, schedules) {
    state.schedules = schedules;
  },
  ADD_SCHEDULE(state, schedule) {
    state.schedules.push(schedule);
  },
  UPDATE_SCHEDULE(state, schedule) {
    const idx = state.schedules.findIndex(s => s.id === schedule.id);
    if (idx !== -1) state.schedules[idx] = schedule;
  }
};

const actions = {
  async fetchSchedules({ commit }) {
    const res = await axios.get('/schedules');
    commit('SET_SCHEDULES', res.data.data);
  },
  async addSchedule({ commit }, schedule) {
    const res = await axios.post('/schedules', schedule);
    commit('ADD_SCHEDULE', res.data.data);
  },
  async requestSwap({ commit }, { id, target_employee_id }) {
    await axios.put(`/schedules/${id}/swap`, { target_employee_id });
  },
  async requestAbsence({ commit }, { id, reason }) {
    await axios.put(`/schedules/${id}/absence`, { reason });
  },
  async updateScheduleStatus({ commit }, { id, status }) {
    await axios.put(`/schedules/${id}/status`, { status });
  }
};

const getters = {
  schedules: state => state.schedules
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
}; 