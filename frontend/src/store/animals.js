import api from '../utils/api';

const state = {
  animals: []
};

const mutations = {
  SET_ANIMALS(state, animals) {
    state.animals = animals;
  },
  ADD_ANIMAL(state, animal) {
    state.animals.push(animal);
  },
  UPDATE_ANIMAL(state, animal) {
    const idx = state.animals.findIndex(a => a.id === animal.id);
    if (idx !== -1) state.animals[idx] = animal;
  },
  DELETE_ANIMAL(state, id) {
    state.animals = state.animals.filter(a => a.id !== id);
  }
};

const actions = {
  async fetchAnimals({ commit }) {
    const res = await api.get('/animals');
    commit('SET_ANIMALS', res.data.data);
  },
  async addAnimal({ commit }, animal) {
    const res = await api.post('/animals', animal);
    commit('ADD_ANIMAL', res.data.data);
  },
  async updateAnimal({ commit }, animal) {
    await api.put(`/animals/${animal.id}`, animal);
    commit('UPDATE_ANIMAL', animal);
  },
  async deleteAnimal({ commit }, id) {
    await api.delete(`/animals/${id}`);
    commit('DELETE_ANIMAL', id);
  }
};

const getters = {
  animals: state => state.animals
};

export default {
  namespaced: true,
  state,
  mutations,
  actions,
  getters
};