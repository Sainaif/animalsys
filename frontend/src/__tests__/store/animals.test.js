import { describe, it, expect, beforeEach } from 'vitest'
import { createStore } from 'vuex'
import animals from '../../store/modules/animals'

describe('Animals Store Module', () => {
  let store

  beforeEach(() => {
    store = createStore({
      modules: {
        animals
      }
    })
  })

  describe('mutations', () => {
    it('SET_ANIMALS should set animals array', () => {
      const animalList = [
        { id: '1', name: 'Max', species: 'Dog' },
        { id: '2', name: 'Bella', species: 'Cat' }
      ]

      store.commit('animals/SET_ANIMALS', animalList)

      expect(store.state.animals.animals).toEqual(animalList)
      expect(store.state.animals.animals).toHaveLength(2)
    })

    it('ADD_ANIMAL should add animal to the list', () => {
      const animal = { id: '1', name: 'Max', species: 'Dog' }

      store.commit('animals/ADD_ANIMAL', animal)

      expect(store.state.animals.animals).toContain(animal)
      expect(store.state.animals.animals).toHaveLength(1)
    })

    it('UPDATE_ANIMAL should update existing animal', () => {
      const animal = { id: '1', name: 'Max', species: 'Dog' }
      store.commit('animals/SET_ANIMALS', [animal])

      const updatedAnimal = { id: '1', name: 'Max Updated', species: 'Dog' }
      store.commit('animals/UPDATE_ANIMAL', updatedAnimal)

      expect(store.state.animals.animals[0].name).toBe('Max Updated')
    })

    it('DELETE_ANIMAL should remove animal from the list', () => {
      const animals = [
        { id: '1', name: 'Max', species: 'Dog' },
        { id: '2', name: 'Bella', species: 'Cat' }
      ]
      store.commit('animals/SET_ANIMALS', animals)

      store.commit('animals/DELETE_ANIMAL', '1')

      expect(store.state.animals.animals).toHaveLength(1)
      expect(store.state.animals.animals[0].id).toBe('2')
    })
  })

  describe('state', () => {
    it('should initialize with empty animals array', () => {
      expect(store.state.animals.animals).toEqual([])
    })
  })
})
