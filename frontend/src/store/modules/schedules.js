import api from '../../utils/api'

const state = {
  schedules: []
}

const mutations = {
  SET_SCHEDULES(state, schedules) {
    state.schedules = schedules
  },
  ADD_SCHEDULE(state, schedule) {
    state.schedules.push(schedule)
  },
  UPDATE_SCHEDULE(state, updatedSchedule) {
    const index = state.schedules.findIndex(s => s.id === updatedSchedule.id)
    if (index !== -1) {
      state.schedules.splice(index, 1, updatedSchedule)
    }
  },
  DELETE_SCHEDULE(state, id) {
    state.schedules = state.schedules.filter(s => s.id !== id)
  }
}

const actions = {
  async fetchSchedules({ commit }) {
    const response = await api.get('/schedules')
    commit('SET_SCHEDULES', response.data.data)
  },
  async createSchedule({ commit }, schedule) {
    const response = await api.post('/schedules', schedule)
    commit('ADD_SCHEDULE', response.data.data)
  },
  async updateSchedule({ commit }, { id, data }) {
    await api.put(`/schedules/${id}`, data)
    commit('UPDATE_SCHEDULE', { id, ...data })
  },
  async deleteSchedule({ commit }, id) {
    await api.delete(`/schedules/${id}`)
    commit('DELETE_SCHEDULE', id)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
