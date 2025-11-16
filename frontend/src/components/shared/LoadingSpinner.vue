<template>
  <div
    v-if="fullPage"
    class="loading-overlay"
  >
    <ProgressSpinner :style="{ width: spinnerSize, height: spinnerSize }" />
  </div>
  <div
    v-else
    class="loading-inline"
    :class="{ 'loading-center': center }"
  >
    <ProgressSpinner :style="{ width: spinnerSize, height: spinnerSize }" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import ProgressSpinner from 'primevue/progressspinner'

const props = defineProps({
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  },
  fullPage: {
    type: Boolean,
    default: false
  },
  center: {
    type: Boolean,
    default: true
  }
})

const spinnerSize = computed(() => {
  const sizes = {
    sm: '30px',
    md: '50px',
    lg: '80px'
  }
  return sizes[props.size]
})
</script>

<style scoped>
.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(255, 255, 255, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.loading-inline {
  padding: 2rem;
}

.loading-center {
  display: flex;
  align-items: center;
  justify-content: center;
}
</style>
