<template>
  <div :class="groupClasses">
    <label v-if="label" :for="inputId" class="form-label">
      {{ label }}
      <span v-if="required" class="required-mark">*</span>
    </label>

    <div class="form-input-wrapper">
      <slot></slot>
    </div>

    <p v-if="error" class="form-error">{{ error }}</p>
    <p v-else-if="hint" class="form-hint">{{ hint }}</p>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  label: {
    type: String,
    default: ''
  },
  inputId: {
    type: String,
    default: ''
  },
  error: {
    type: String,
    default: ''
  },
  hint: {
    type: String,
    default: ''
  },
  required: {
    type: Boolean,
    default: false
  }
})

const groupClasses = computed(() => {
  return [
    'form-group',
    {
      'form-group--error': props.error
    }
  ]
})
</script>

<style scoped>
.form-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
  margin-bottom: 1rem;
}

.form-label {
  font-weight: 500;
  color: var(--text-primary);
  font-size: 0.875rem;
}

.required-mark {
  color: #ef4444;
  margin-left: 0.25rem;
}

.form-input-wrapper {
  position: relative;
}

.form-error {
  color: #ef4444;
  font-size: 0.875rem;
  margin: 0;
}

.form-hint {
  color: var(--text-tertiary);
  font-size: 0.875rem;
  margin: 0;
}

.form-group--error .form-input-wrapper :deep(input),
.form-group--error .form-input-wrapper :deep(textarea),
.form-group--error .form-input-wrapper :deep(select) {
  border-color: #ef4444;
}

/* Global input styles */
:deep(input[type="text"]),
:deep(input[type="email"]),
:deep(input[type="password"]),
:deep(input[type="number"]),
:deep(input[type="date"]),
:deep(input[type="time"]),
:deep(textarea),
:deep(select) {
  width: 100%;
  padding: 0.625rem 0.75rem;
  border: 1px solid var(--border-color);
  border-radius: 0.375rem;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  font-size: 1rem;
  font-family: inherit;
  transition: border-color 0.2s;
}

:deep(input:focus),
:deep(textarea:focus),
:deep(select:focus) {
  outline: none;
  border-color: var(--primary-color);
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

:deep(input:disabled),
:deep(textarea:disabled),
:deep(select:disabled) {
  opacity: 0.6;
  cursor: not-allowed;
  background-color: var(--bg-disabled);
}

:deep(textarea) {
  min-height: 100px;
  resize: vertical;
}
</style>
