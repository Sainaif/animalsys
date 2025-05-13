<template>
  <form @submit.prevent="onSubmit">
    <div>
      <label>Imię i nazwisko</label>
      <input v-model="form.applicant_name" required />
    </div>
    <div>
      <label>Email</label>
      <input v-model="form.applicant_email" type="email" required />
    </div>
    <div>
      <label>Uwagi</label>
      <textarea v-model="form.notes"></textarea>
    </div>
    <button type="submit">Wyślij wniosek</button>
  </form>
</template>

<script setup>
import { ref } from 'vue';
const emit = defineEmits(['apply']);
const form = ref({
  applicant_name: '',
  applicant_email: '',
  notes: ''
});
const onSubmit = () => {
  emit('apply', { ...form.value });
};
</script> 