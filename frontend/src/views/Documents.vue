<template>
  <div>
    <h1>{{$t('documents')}}</h1>
    <DocumentUpload v-if="isAuthenticated" @uploaded="refresh" />
    <table>
      <thead>
        <tr>
          <th>Nazwa</th><th>Typ</th><th>Rozmiar</th><th>Data</th><th>Pobierz</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="doc in documents" :key="doc.id">
          <td>{{doc.filename}}</td>
          <td>{{doc.type}}</td>
          <td>{{(doc.size/1024).toFixed(1)}} KB</td>
          <td>{{doc.upload_date}}</td>
          <td><a :href="downloadUrl(doc.id)" target="_blank">Pobierz</a></td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script setup>
import { computed, onMounted } from 'vue';
import { useStore } from 'vuex';
import DocumentUpload from '../components/DocumentUpload.vue';

const store = useStore();
const documents = computed(() => store.getters['documents/documents']);
const isAuthenticated = computed(() => store.getters['auth/isAuthenticated']);
const refresh = () => store.dispatch('documents/fetchDocuments');
const downloadUrl = id => `${import.meta.env.VITE_API_URL}/api/documents/${id}`;

onMounted(refresh);
</script>

<style>
table { width: 100%; border-collapse: collapse; margin-top: 1em; }
th, td { border: 1px solid #ccc; padding: 0.5em; }
</style> 