<template>
  <div class="home">
    <h1>{{ $t('home.title') }}</h1>
    <p>{{ $t('home.welcome') }}</p>
    <div class="dashboard">
      <div class="card">
        <h3>{{ $t('home.animals') }}</h3>
        <p>{{ $t('home.manageAnimals') }}</p>
        <router-link to="/animals">{{ $t('home.view') }}</router-link>
      </div>
      <div class="card">
        <h3>{{ $t('home.adoptions') }}</h3>
        <p>{{ $t('home.manageAdoptions') }}</p>
        <router-link to="/adoptions">{{ $t('home.view') }}</router-link>
      </div>
      <div class="card">
        <h3>{{ $t('home.schedules') }}</h3>
        <p>{{ $t('home.manageSchedules') }}</p>
        <router-link to="/schedules">{{ $t('home.view') }}</router-link>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'Home'
}
</script>

<style scoped>
.home {
  padding: 20px;
}

.dashboard {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 20px;
  margin-top: 30px;
}

.card {
  border: 1px solid #ddd;
  padding: 20px;
  border-radius: 8px;
}

.card h3 {
  margin-top: 0;
}

.card a {
  display: inline-block;
  margin-top: 10px;
  color: #42b983;
  text-decoration: none;
}
</style>
