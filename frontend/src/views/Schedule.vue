<template>
  <div>
    <h1>{{$t('schedule')}}</h1>
    <button v-if="canAssign" @click="showForm = true">Przydziel dyżur</button>
    <form v-if="showForm" @submit.prevent="addSchedule">
      <input v-model="form.employee_id" placeholder="ID pracownika" required />
      <input v-model="form.shift_date" type="date" required />
      <input v-model="form.shift_time" placeholder="Godzina" required />
      <input v-model="form.tasks" placeholder="Zadania (przecinki)" />
      <button type="submit">Dodaj</button>
    </form>
    <table>
      <thead>
        <tr>
          <th>Pracownik</th><th>Data</th><th>Godzina</th><th>Zadania</th><th>Status</th><th>Akcje</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="s in schedules" :key="s.id">
          <td>{{s.employee_id}}</td>
          <td>{{s.shift_date}}</td>
          <td>{{s.shift_time}}</td>
          <td>{{s.tasks?.join(', ')}}</td>
          <td>{{s.status}}</td>
          <td>
            <button v-if="canSwap" @click="requestSwap(s.id)">Zaproponuj zamianę</button>
            <button v-if="canAbsence" @click="requestAbsence(s.id)">Zgłoś nieobecność</button>
            <button v-if="canUpdateStatus" @click="updateStatus(s.id, 'swapped')">Zatwierdź zamianę</button>
            <button v-if="canUpdateStatus" @click="updateStatus(s.id, 'absent_approved')">Zatwierdź nieobecność</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useStore } from 'vuex';

const store = useStore();
const schedules = computed(() => store.getters['schedule/schedules']);
const role = computed(() => store.getters['auth/userRole']);
const canAssign = computed(() => role.value === 'admin' || role.value === 'employee');
const canSwap = computed(() => role.value === 'volunteer');
const canAbsence = computed(() => role.value === 'volunteer');
const canUpdateStatus = computed(() => role.value === 'admin' || role.value === 'employee');
const showForm = ref(false);
const form = ref({ employee_id: '', shift_date: '', shift_time: '', tasks: '' });

onMounted(() => {
  store.dispatch('schedule/fetchSchedules');
});

const addSchedule = async () => {
  await store.dispatch('schedule/addSchedule', {
    ...form.value,
    tasks: form.value.tasks.split(',').map(s => s.trim()).filter(Boolean)
  });
  showForm.value = false;
  store.dispatch('schedule/fetchSchedules');
};

const requestSwap = async (id) => {
  const target_employee_id = prompt('Podaj ID pracownika do zamiany:');
  if (target_employee_id) {
    await store.dispatch('schedule/requestSwap', { id, target_employee_id });
    store.dispatch('schedule/fetchSchedules');
  }
};

const requestAbsence = async (id) => {
  const reason = prompt('Podaj powód nieobecności:');
  if (reason) {
    await store.dispatch('schedule/requestAbsence', { id, reason });
    store.dispatch('schedule/fetchSchedules');
  }
};

const updateStatus = async (id, status) => {
  await store.dispatch('schedule/updateScheduleStatus', { id, status });
  store.dispatch('schedule/fetchSchedules');
};
</script>

<style>
table { width: 100%; border-collapse: collapse; margin-top: 1em; }
th, td { border: 1px solid #ccc; padding: 0.5em; }
</style> 