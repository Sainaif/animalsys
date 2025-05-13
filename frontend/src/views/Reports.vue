<template>
  <div>
    <h1>Raporty finansowe</h1>
    <form v-if="canAdd" @submit.prevent="addFinance">
      <input v-model="form.date" type="date" required />
      <select v-model="form.type" required>
        <option value="income">Przychód</option>
        <option value="expense">Wydatek</option>
      </select>
      <input v-model.number="form.amount" type="number" step="0.01" required placeholder="Kwota" />
      <input v-model="form.description" placeholder="Opis" />
      <input v-model="form.category" placeholder="Kategoria" />
      <button type="submit">Dodaj</button>
    </form>
    <button @click="exportCsv">Eksportuj CSV</button>
    <table>
      <thead>
        <tr>
          <th>Data</th><th>Typ</th><th>Kwota</th><th>Opis</th><th>Kategoria</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="f in finances" :key="f.id">
          <td>{{f.date}}</td>
          <td>{{f.type}}</td>
          <td>{{f.amount}}</td>
          <td>{{f.description}}</td>
          <td>{{f.category}}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';
const store = useStore();
const finances = computed(() => store.getters['finances/finances']);
const role = computed(() => store.getters['auth/userRole']);
const canAdd = computed(() => role.value === 'admin' || role.value === 'employee');
const form = ref({ date: '', type: 'income', amount: 0, description: '', category: '' });

onMounted(() => {
  store.dispatch('finances/fetchFinances');
});

const addFinance = async () => {
  await store.dispatch('finances/addFinance', form.value);
  store.dispatch('finances/fetchFinances');
};

const exportCsv = () => {
  store.dispatch('finances/exportCsv');
};
</script>

<style>
table { width: 100%; border-collapse: collapse; margin-top: 1em; }
th, td { border: 1px solid #ccc; padding: 0.5em; }
</style> 