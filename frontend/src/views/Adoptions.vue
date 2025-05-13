<template>
  <div>
    <h1>{{$t('adoptions')}}</h1>
    <h2>Dostępne zwierzęta do adopcji</h2>
    <ul>
      <li v-for="animal in availableAnimals" :key="animal.id">
        {{animal.name}} ({{animal.species}}, {{animal.breed}})
        <button @click="selectAnimal(animal)">Złóż wniosek</button>
      </li>
    </ul>
    <AdoptionForm v-if="selectedAnimal" @apply="applyForAdoption" />
    <h2>Wnioski adopcyjne</h2>
    <table>
      <thead>
        <tr>
          <th>Zwierzę</th><th>Status</th><th>Dane wniosku</th><th v-if="canApprove">Akcje</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="adoption in adoptions" :key="adoption.id">
          <td>{{animalName(adoption.animal_id)}}</td>
          <td>{{adoption.status}}</td>
          <td>{{adoption.application_data | json}}</td>
          <td v-if="canApprove">
            <button @click="updateStatus(adoption.id, 'approved')">Akceptuj</button>
            <button @click="updateStatus(adoption.id, 'rejected')">Odrzuć</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import AdoptionForm from '../components/AdoptionForm.vue';

const store = useStore();
const adoptions = computed(() => store.getters['adoptions/adoptions']);
const animals = computed(() => store.getters['animals/animals']);
const role = computed(() => store.getters['auth/userRole']);
const canApprove = computed(() => role.value === 'admin' || role.value === 'employee');
const selectedAnimal = ref(null);

onMounted(() => {
  store.dispatch('adoptions/fetchAdoptions');
  store.dispatch('animals/fetchAnimals');
});

const availableAnimals = computed(() => animals.value.filter(a => a.status === 'available'));

const selectAnimal = animal => {
  selectedAnimal.value = animal;
};

const applyForAdoption = async (form) => {
  await store.dispatch('adoptions/applyForAdoption', {
    animal_id: selectedAnimal.value.id,
    application_data: form
  });
  selectedAnimal.value = null;
  store.dispatch('adoptions/fetchAdoptions');
};

const updateStatus = async (id, status) => {
  await store.dispatch('adoptions/updateAdoptionStatus', { id, status });
  store.dispatch('adoptions/fetchAdoptions');
};

const animalName = (id) => {
  const a = animals.value.find(a => a.id === id);
  return a ? a.name : id;
};
</script>

<script>
export default {
  filters: {
    json(val) {
      return JSON.stringify(val);
    }
  }
};
</script>

<style>
table { width: 100%; border-collapse: collapse; margin-top: 1em; }
th, td { border: 1px solid #ccc; padding: 0.5em; }
</style> 