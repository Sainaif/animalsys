<template>
  <form @submit.prevent="onSubmit">
    <div>
      <label>Nazwa</label>
      <input v-model="form.name" required />
    </div>
    <div>
      <label>Gatunek</label>
      <input v-model="form.species" required />
    </div>
    <div>
      <label>Rasa</label>
      <input v-model="form.breed" required />
    </div>
    <div>
      <label>Wiek</label>
      <input v-model.number="form.age" type="number" min="0" required />
    </div>
    <div>
      <label>Historia zdrowotna (oddziel przecinkami)</label>
      <input v-model="healthHistoryStr" />
    </div>
    <button type="submit">Zapisz</button>
  </form>
</template>

<script setup>
import { ref, watch, computed } from 'vue';
const props = defineProps({ animal: Object });
const emit = defineEmits(['save']);
const form = ref({
  name: '',
  species: '',
  breed: '',
  age: 0,
  health_history: []
});
const healthHistoryStr = ref('');

watch(() => props.animal, (val) => {
  if (val) {
    form.value = { ...val };
    healthHistoryStr.value = val.health_history?.join(', ') || '';
  }
}, { immediate: true });

watch(healthHistoryStr, val => {
  form.value.health_history = val.split(',').map(s => s.trim()).filter(Boolean);
});

const onSubmit = () => {
  emit('save', { ...form.value });
};
</script> 