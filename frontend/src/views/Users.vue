<template>
  <div class="users">
    <h1>{{ $t('users.title') }}</h1>
    <table>
      <thead>
        <tr>
          <th>{{ $t('users.username') }}</th>
          <th>{{ $t('users.email') }}</th>
          <th>{{ $t('users.role') }}</th>
          <th>{{ $t('users.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="user in users" :key="user.id">
          <td>{{ user.username }}</td>
          <td>{{ user.email }}</td>
          <td>{{ user.role }}</td>
          <td>
            <button @click="deleteUser(user.id)">{{ $t('users.delete') }}</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'Users',
  setup() {
    const store = useStore()
    const users = computed(() => store.state.users.users)

    const deleteUser = async (id) => {
      if (confirm('Are you sure?')) {
        await store.dispatch('users/deleteUser', id)
      }
    }

    onMounted(() => {
      store.dispatch('users/fetchUsers')
    })

    return {
      users,
      deleteUser
    }
  }
}
</script>

<style scoped>
.users {
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
