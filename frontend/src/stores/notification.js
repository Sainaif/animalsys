import { defineStore } from 'pinia'
import { ref } from 'vue'

let notificationIdCounter = 0

export const useNotificationStore = defineStore('notification', () => {
  // State
  const notifications = ref([])

  // Actions
  function addNotification({ type = 'info', title, message, duration = 5000 }) {
    const id = ++notificationIdCounter

    const notification = {
      id,
      type, // 'success', 'error', 'warning', 'info'
      title,
      message,
      duration,
      timestamp: Date.now(),
    }

    notifications.value.push(notification)

    // Auto remove after duration
    if (duration > 0) {
      setTimeout(() => {
        removeNotification(id)
      }, duration)
    }

    return id
  }

  function removeNotification(id) {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index !== -1) {
      notifications.value.splice(index, 1)
    }
  }

  function clearAll() {
    notifications.value = []
  }

  // Convenience methods
  function success(message, title = '') {
    return addNotification({ type: 'success', title, message })
  }

  function error(message, title = '') {
    return addNotification({ type: 'error', title, message, duration: 7000 })
  }

  function warning(message, title = '') {
    return addNotification({ type: 'warning', title, message })
  }

  function info(message, title = '') {
    return addNotification({ type: 'info', title, message })
  }

  return {
    // State
    notifications,

    // Actions
    addNotification,
    removeNotification,
    clearAll,
    success,
    error,
    warning,
    info,
  }
})
