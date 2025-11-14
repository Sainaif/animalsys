<template>
  <span class="badge" :class="variantClass">
    <slot>{{ text }}</slot>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  text: String,
  variant: {
    type: String,
    default: 'neutral',
    validator: (value) => ['success', 'warning', 'danger', 'info', 'neutral', 'primary'].includes(value)
  }
})

const variantClass = computed(() => `badge-${props.variant}`)
</script>

<style scoped>
.badge {
  display: inline-flex;
  align-items: center;
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  white-space: nowrap;
}

.badge-success {
  background-color: #d1fae5;
  color: #065f46;
}

.badge-warning {
  background-color: #fef3c7;
  color: #92400e;
}

.badge-danger {
  background-color: #fee2e2;
  color: #991b1b;
}

.badge-info {
  background-color: #dbeafe;
  color: #1e40af;
}

.badge-primary {
  background-color: #e0e7ff;
  color: #3730a3;
}

.badge-neutral {
  background-color: #f3f4f6;
  color: #374151;
}
</style>
