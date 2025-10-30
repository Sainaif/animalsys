<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="buttonClasses"
    @click="handleClick"
  >
    <span v-if="loading" class="loading-spinner"></span>
    <slot v-else></slot>
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  variant: {
    type: String,
    default: 'primary', // primary, secondary, outline, danger, success, ghost
    validator: (value) => ['primary', 'secondary', 'outline', 'danger', 'success', 'ghost'].includes(value)
  },
  size: {
    type: String,
    default: 'medium', // small, medium, large
    validator: (value) => ['small', 'medium', 'large'].includes(value)
  },
  type: {
    type: String,
    default: 'button',
    validator: (value) => ['button', 'submit', 'reset'].includes(value)
  },
  disabled: {
    type: Boolean,
    default: false
  },
  loading: {
    type: Boolean,
    default: false
  },
  fullWidth: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['click'])

const buttonClasses = computed(() => {
  return [
    'base-button',
    `base-button--${props.variant}`,
    `base-button--${props.size}`,
    {
      'base-button--disabled': props.disabled,
      'base-button--loading': props.loading,
      'base-button--full-width': props.fullWidth
    }
  ]
})

function handleClick(event) {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>

<style scoped>
.base-button {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  font-weight: 500;
  border-radius: 0.375rem;
  transition: all 0.2s;
  cursor: pointer;
  border: none;
  font-family: inherit;
}

.base-button:focus {
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.5);
}

/* Sizes */
.base-button--small {
  padding: 0.5rem 0.75rem;
  font-size: 0.875rem;
}

.base-button--medium {
  padding: 0.625rem 1rem;
  font-size: 1rem;
}

.base-button--large {
  padding: 0.75rem 1.5rem;
  font-size: 1.125rem;
}

/* Variants */
.base-button--primary {
  background-color: var(--primary-color);
  color: white;
}

.base-button--primary:hover:not(.base-button--disabled) {
  opacity: 0.9;
  transform: translateY(-1px);
}

.base-button--secondary {
  background-color: var(--bg-secondary);
  color: var(--text-primary);
  border: 1px solid var(--border-color);
}

.base-button--secondary:hover:not(.base-button--disabled) {
  background-color: var(--bg-hover);
}

.base-button--outline {
  background-color: transparent;
  color: var(--primary-color);
  border: 1px solid var(--primary-color);
}

.base-button--outline:hover:not(.base-button--disabled) {
  background-color: var(--primary-color);
  color: white;
}

.base-button--danger {
  background-color: #ef4444;
  color: white;
}

.base-button--danger:hover:not(.base-button--disabled) {
  background-color: #dc2626;
}

.base-button--success {
  background-color: #10b981;
  color: white;
}

.base-button--success:hover:not(.base-button--disabled) {
  background-color: #059669;
}

.base-button--ghost {
  background-color: transparent;
  color: var(--text-secondary);
}

.base-button--ghost:hover:not(.base-button--disabled) {
  background-color: var(--bg-hover);
  color: var(--text-primary);
}

/* States */
.base-button--disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.base-button--loading {
  cursor: wait;
}

.base-button--full-width {
  width: 100%;
}

.loading-spinner {
  width: 1rem;
  height: 1rem;
  border: 2px solid transparent;
  border-top-color: currentColor;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
