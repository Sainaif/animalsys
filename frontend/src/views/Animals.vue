<template>
  <div>
    <h1>{{$t('animals')}}</h1>
    <button v-if="canEdit" @click="showForm = true">Dodaj zwierzę</button>
    <AnimalForm v-if="showForm" @save="saveAnimal" :animal="editAnimal" />
    <table>
      <thead>
        <tr>
          <th>Nazwa</th><th>Gatunek</th><th>Rasa</th><th>Wiek</th><th>Historia zdrowotna</th><th>Status</th><th v-if="canEdit">Akcje</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="animal in animals" :key="animal.id">
          <td>{{animal.name}}</td>
          <td>{{animal.species}}</td>
          <td>{{animal.breed}}</td>
          <td>{{animal.age}}</td>
          <td>{{animal.health_history?.join(', ')}}</td>
          <td>{{animal.status}}</td>
          <td v-if="canEdit">
            <button @click="edit(animal)">Edytuj</button>
            <button @click="remove(animal.id)">Usuń</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import AnimalForm from '../components/AnimalForm.vue';

const store = useStore();
const animals = computed(() => store.getters['animals/animals']);
const role = computed(() => store.getters['auth/userRole']);
const canEdit = computed(() => role.value === 'admin' || role.value === 'employee');
const showForm = ref(false);
const editAnimal = ref(null);

onMounted(() => {
  store.dispatch('animals/fetchAnimals');
});

const saveAnimal = async animal => {
  if (editAnimal.value) {
    await store.dispatch('animals/updateAnimal', animal);
  } else {
    await store.dispatch('animals/addAnimal', animal);
  }
  showForm.value = false;
  editAnimal.value = null;
  store.dispatch('animals/fetchAnimals');
};

const edit = animal => {
  editAnimal.value = { ...animal };
  showForm.value = true;
};

const remove = async id => {
  await store.dispatch('animals/deleteAnimal', id);
  store.dispatch('animals/fetchAnimals');
};
</script>

<style>
table { width: 100%; border-collapse: collapse; margin-top: 1em; }
th, td { border: 1px solid #ccc; padding: 0.5em; }
</style> 