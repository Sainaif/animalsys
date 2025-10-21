<template>
  <div class="adoptions">
    <h1>{{ $t('adoptions.title') }}</h1>
    <button @click="showForm = true">{{ $t('adoptions.add') }}</button>

    <div v-if="showForm" class="form-container">
      <h2>{{ $t('adoptions.addNew') }}</h2>
      <form @submit.prevent="handleCreate">
        <input v-model="newAdoption.animal_id" placeholder="Animal ID" required />
        <input v-model="newAdoption.user_id" placeholder="User ID" required />
        <select v-model="newAdoption.status">
          <option value="pending">Pending</option>
          <option value="approved">Approved</option>
          <option value="rejected">Rejected</option>
        </select>
        <button type="submit">{{ $t('adoptions.save') }}</button>
        <button type="button" @click="showForm = false">{{ $t('adoptions.cancel') }}</button>
      </form>
    </div>

    <table>
      <thead>
        <tr>
          <th>{{ $t('adoptions.animalId') }}</th>
          <th>{{ $t('adoptions.userId') }}</th>
          <th>{{ $t('adoptions.status') }}</th>
          <th>{{ $t('adoptions.actions') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="adoption in adoptions" :key="adoption.id">
          <td>{{ adoption.animal_id }}</td>
          <td>{{ adoption.user_id }}</td>
          <td>{{ adoption.status }}</td>
          <td>
            <button @click="updateStatus(adoption.id, 'approved')">Approve</button>
            <button @click="updateStatus(adoption.id, 'rejected')">Reject</button>
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'Adoptions',
  setup() {
    const store = useStore()
    const showForm = ref(false)
    const newAdoption = ref({
      animal_id: '',
      user_id: '',
      status: 'pending',
      application_data: {}
    })

    const adoptions = computed(() => store.state.adoptions.adoptions)

    const handleCreate = async () => {
      await store.dispatch('adoptions/createAdoption', newAdoption.value)
      newAdoption.value = { animal_id: '', user_id: '', status: 'pending', application_data: {} }
      showForm.value = false
    }

    const updateStatus = async (id, status) => {
      await store.dispatch('adoptions/updateAdoption', { id, data: { status } })
    }

    onMounted(() => {
      store.dispatch('adoptions/fetchAdoptions')
    })

    return {
      adoptions,
      showForm,
      newAdoption,
      handleCreate,
      updateStatus
    }
  }
}
</script>

<style scoped>
.adoptions {
  padding: 20px;
}

.form-container {
  margin: 20px 0;
  padding: 20px;
  border: 1px solid #ddd;
  border-radius: 8px;
}

form input, form select {
  display: block;
  width: 100%;
  margin-bottom: 10px;
  padding: 8px;
  box-sizing: border-box;
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
