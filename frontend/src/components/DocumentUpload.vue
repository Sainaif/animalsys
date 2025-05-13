<template>
  <form @submit.prevent="onUpload">
    <input type="file" @change="onFile" required />
    <select v-model="type" required>
      <option value="medical">Raport medyczny</option>
      <option value="contract">Umowa</option>
      <option value="other">Inny</option>
    </select>
    <button type="submit">Wyślij</button>
  </form>
</template>

<script setup>
import { ref } from 'vue';
import { useStore } from 'vuex';
const emit = defineEmits(['uploaded']);
const store = useStore();
const file = ref(null);
const type = ref('medical');
const onFile = e => { file.value = e.target.files[0]; };
const onUpload = async () => {
  if (!file.value) return;
  await store.dispatch('documents/uploadDocument', { file: file.value, type: type.value });
  emit('uploaded');
};
</script> 