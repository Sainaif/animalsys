<template>
  <div class="finances">
    <h1>{{ $t('finances.title') }}</h1>
    <table>
      <thead>
        <tr>
          <th>{{ $t('finances.date') }}</th>
          <th>{{ $t('finances.type') }}</th>
          <th>{{ $t('finances.amount') }}</th>
          <th>{{ $t('finances.description') }}</th>
          <th>{{ $t('finances.category') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="finance in finances" :key="finance.id">
          <td>{{ finance.date }}</td>
          <td>{{ finance.type }}</td>
          <td>${{ finance.amount.toFixed(2) }}</td>
          <td>{{ finance.description }}</td>
          <td>{{ finance.category }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'Finances',
  setup() {
    const store = useStore()
    const finances = computed(() => store.state.finances.finances)

    onMounted(() => {
      store.dispatch('finances/fetchFinances')
    })

    return {
      finances
    }
  }
}
</script>

<style scoped>
.finances {
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
