<template>
  <div class="empty-state">
    <i
      v-if="icon"
      class="pi"
      :class="icon"
    />
    <h3 v-if="title">
      {{ title }}
    </h3>
    <p>{{ message }}</p>
    <Button
      v-if="actionText"
      :label="actionText"
      :icon="actionIcon"
      @click="handleAction"
    />
  </div>
</template>

<script setup>
import Button from 'primevue/button'

defineProps({
  icon: {
    type: String,
    default: 'pi-inbox'
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    required: true
  },
  actionText: {
    type: String,
    default: ''
  },
  actionIcon: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['action'])

const handleAction = () => {
  emit('action')
}
</script>

<style scoped>
.empty-state {
  text-align: center;
  padding: 4rem 2rem;
  color: #6b7280;
}

.empty-state i {
  font-size: 4rem;
  color: #d1d5db;
  margin-bottom: 1.5rem;
}

.empty-state h3 {
  font-size: 1.25rem;
  font-weight: 600;
  color: #374151;
  margin-bottom: 0.5rem;
}

.empty-state p {
  font-size: 0.875rem;
  margin-bottom: 1.5rem;
}
</style>
