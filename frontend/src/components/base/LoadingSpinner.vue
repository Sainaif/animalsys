<template>
  <div :class="spinnerClasses">
    <div class="spinner"></div>
    <p v-if="text" class="spinner-text">{{ text }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  size: {
    type: String,
    default: 'medium', // small, medium, large
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  },
  text: {
    type: String,
    default: ''
  },
  fullscreen: {
    type: Boolean,
    default: false
  }
})

const spinnerClasses = computed(() => {
  return [
    'loading-spinner',
    `loading-spinner--${props.size}`,
    {
      'loading-spinner--fullscreen': props.fullscreen
    }
  ]
})
</script>

<style scoped>
.loading-spinner {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 1rem;
  padding: 2rem;
}

.loading-spinner--fullscreen {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  z-index: 9999;
}

.spinner {
  border-radius: 50%;
  border: 3px solid var(--border-color);
  border-top-color: var(--primary-color);
  animation: spin 0.8s linear infinite;
}

.loading-spinner--small .spinner {
  width: 1.5rem;
  height: 1.5rem;
  border-width: 2px;
}

.loading-spinner--medium .spinner {
  width: 2.5rem;
  height: 2.5rem;
  border-width: 3px;
}

.loading-spinner--large .spinner {
  width: 4rem;
  height: 4rem;
  border-width: 4px;
}

.spinner-text {
  color: var(--text-secondary);
  font-size: 0.875rem;
  margin: 0;
}

.loading-spinner--fullscreen .spinner-text {
  color: white;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
