<template>
  <div class="schedules">
    <h1>{{ $t('schedules.title') }}</h1>
    <table>
      <thead>
        <tr>
          <th>{{ $t('schedules.employeeId') }}</th>
          <th>{{ $t('schedules.date') }}</th>
          <th>{{ $t('schedules.time') }}</th>
          <th>{{ $t('schedules.tasks') }}</th>
          <th>{{ $t('schedules.status') }}</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="schedule in schedules" :key="schedule.id">
          <td>{{ schedule.employee_id }}</td>
          <td>{{ schedule.shift_date }}</td>
          <td>{{ schedule.shift_time }}</td>
          <td>{{ schedule.tasks.join(', ') }}</td>
          <td>{{ schedule.status }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { computed, onMounted } from 'vue'
import { useStore } from 'vuex'

export default {
  name: 'Schedules',
  setup() {
    const store = useStore()
    const schedules = computed(() => store.state.schedules.schedules)

    onMounted(() => {
      store.dispatch('schedules/fetchSchedules')
    })

    return {
      schedules
    }
  }
}
</script>

<style scoped>
.schedules {
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
