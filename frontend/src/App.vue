<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import PublicLayout from '@/layouts/PublicLayout.vue'
import StaffLayout from '@/layouts/StaffLayout.vue'

const route = useRoute()

const layout = computed(() => {
  return route.meta.layout || 'staff'
})
</script>

<template>
  <div id="app">
    <component :is="layout === 'public' ? PublicLayout : layout === 'staff' ? StaffLayout : 'div'">
      <RouterView />
    </component>
  </div>
</template>

<style>
#app {
  min-height: 100vh;
  background-color: var(--surface-ground);
}
</style>
