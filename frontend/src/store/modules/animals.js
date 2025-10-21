import api from '../../utils/api'

const state = () => ({
  animals: []
})

const mutations = {
  SET_ANIMALS(state, animals) {
    state.animals = animals
  },
  ADD_ANIMAL(state, animal) {
    state.animals.push(animal)
  },
  UPDATE_ANIMAL(state, updatedAnimal) {
    const index = state.animals.findIndex(a => a.id === updatedAnimal.id)
    if (index !== -1) {
      state.animals.splice(index, 1, updatedAnimal)
    }
  },
  DELETE_ANIMAL(state, id) {
    state.animals = state.animals.filter(a => a.id !== id)
  }
}

const actions = {
  async fetchAnimals({ commit }) {
    const response = await api.get('/animals')
    commit('SET_ANIMALS', response.data.data)
  },
  async createAnimal({ commit }, animal) {
    const response = await api.post('/animals', animal)
    commit('ADD_ANIMAL', response.data.data)
  },
  async updateAnimal({ commit }, { id, data }) {
    await api.put(`/animals/${id}`, data)
    commit('UPDATE_ANIMAL', { id, ...data })
  },
  async deleteAnimal({ commit }, id) {
    await api.delete(`/animals/${id}`)
    commit('DELETE_ANIMAL', id)
  }
}

export default {
  namespaced: true,
  state,
  mutations,
  actions
}
