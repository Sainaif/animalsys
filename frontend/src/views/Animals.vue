<template>
  <div class="animals">
    <h1>{{ $t('animals.title') }}</h1>
    <button @click="showForm = true" v-if="canCreate">{{ $t('animals.add') }}</button>

    <div v-if="showForm" class="form-container">
      <h2>{{ $t('animals.addNew') }}</h2>
      <form @submit.prevent="handleCreate">
        <input v-model="newAnimal.name" placeholder="Name" required />
        <input v-model="newAnimal.species" placeholder="Species" required />
        <input v-model="newAnimal.breed" placeholder="Breed" required />
        <input v-model.number="newAnimal.age" type="number" placeholder="Age" required />
        <select v-model="newAnimal.status" required>
          <option value="available">Available</option>
          <option value="adopted">Adopted</option>
          <option value="deceased">Deceased</option>
        </select>
        <button type="submit">{{ $t('animals.save') }}</button>
        <button type="button" @click="showForm = false">{{ $t('animals.cancel') }}</button>
      </form>
    </div>

    <table>
      <thead>
        <tr>
          <th>{{ $t('animals.name') }}</th>
          <th>{{ $t('animals.species') }}</th>
          <th>{{ $t('animals.breed') }}</th>
          <th>{{ $t('animals.age') }}</th>
          <th>{{ $t('animals.status') }}</th>
          <th v-if="canCreate">{{ $t('animals.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="animal in animals" :key="animal.id">
          <td>{{ animal.name }}</td>
          <td>{{ animal.species }}</td>
          <td>{{ animal.breed }}</td>
          <td>{{ animal.age }}</td>
          <td>{{ animal.status }}</td>
          <td v-if="canCreate">
            <button @click="deleteAnimal(animal.id)">{{ $t('animals.delete') }}</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'Animals',
  setup() {
    const store = useStore()
    const showForm = ref(false)
    const newAnimal = ref({
      name: '',
      species: '',
      breed: '',
      age: 0,
      status: 'available',
      health_history: []
    })

    const animals = computed(() => store.state.animals.animals)
    const user = computed(() => store.getters['auth/user'])
    const canCreate = computed(() => ['admin', 'employee'].includes(user.value?.role))

    const handleCreate = async () => {
      await store.dispatch('animals/createAnimal', newAnimal.value)
      newAnimal.value = { name: '', species: '', breed: '', age: 0, status: 'available', health_history: [] }
      showForm.value = false
    }

    const deleteAnimal = async (id) => {
      if (confirm('Are you sure?')) {
        await store.dispatch('animals/deleteAnimal', id)
      }
    }

    onMounted(() => {
      store.dispatch('animals/fetchAnimals')
    })

    return {
      animals,
      showForm,
      newAnimal,
      canCreate,
      handleCreate,
      deleteAnimal
    }
  }
}
</script>

<style scoped>
.animals {
  padding: 20px;
}

.form-container {
  margin: 20px 0;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

form input, form select {
  display: block;
  width: 100%;
  margin-bottom: 10px;
  padding: 8px;
  box-sizing: border-box;
}

table {
  width: 100%;
  border-collapse: collapse;
  margin-top: 20px;
}

th, td {
  border: 1px solid #ddd;
  padding: 8px;
  text-align: left;
}

th {
  background-color: #f3f3f3;
}
</style>
