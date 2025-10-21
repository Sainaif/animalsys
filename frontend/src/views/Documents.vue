<template>
  <div class="documents">
    <h1>{{ $t('documents.title') }}</h1>
    <table>
      <thead>
        <tr>
          <th>{{ $t('documents.filename') }}</th>
          <th>{{ $t('documents.type') }}</th>
          <th>{{ $t('documents.uploadDate') }}</th>
          <th>{{ $t('documents.size') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="document in documents" :key="document.id">
          <td>{{ document.filename }}</td>
          <td>{{ document.type }}</td>
          <td>{{ document.upload_date }}</td>
          <td>{{ formatSize(document.size) }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'Documents',
  setup() {
    const store = useStore()
    const documents = computed(() => store.state.documents.documents)

    const formatSize = (bytes) => {
      return (bytes / 1024).toFixed(2) + ' KB'
    }

    onMounted(() => {
      store.dispatch('documents/fetchDocuments')
    })

    return {
      documents,
      formatSize
    }
  }
}
</script>

<style scoped>
.documents {
  padding: 20px;
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
