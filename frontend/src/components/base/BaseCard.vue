<template>
  <div :class="cardClasses">
    <div v-if="$slots.header || title" class="card-header">
      <slot name="header">
        <h3 v-if="title" class="card-title">{{ title }}</h3>
      </slot>
    </div>

    <div class="card-body">
      <slot></slot>
    </div>

    <div v-if="$slots.footer" class="card-footer">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  padding: {
    type: String,
    default: 'normal', // none, small, normal, large
    validator: (value) => ['none', 'small', 'normal', 'large'].includes(value)
  },
  hoverable: {
    type: Boolean,
    default: false
  }
})

const cardClasses = computed(() => {
  return [
    'base-card',
    `base-card--padding-${props.padding}`,
    {
      'base-card--hoverable': props.hoverable
    }
  ]
})
</script>

<style scoped>
.base-card {
  background-color: var(--bg-secondary);
  border: 1px solid var(--border-color);
  border-radius: 0.5rem;
  overflow: hidden;
  transition: all 0.2s;
}

.base-card--hoverable:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.card-header {
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--bg-primary);
}

.card-title {
  margin: 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--text-primary);
}

.card-body {
  color: var(--text-primary);
}

/* Padding variants */
.base-card--padding-none .card-body {
  padding: 0;
}

.base-card--padding-small .card-body {
  padding: 0.75rem;
}

.base-card--padding-normal .card-body {
  padding: 1.5rem;
}

.base-card--padding-large .card-body {
  padding: 2rem;
}

.card-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--border-color);
  background-color: var(--bg-primary);
}
</style>
