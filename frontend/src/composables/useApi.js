import { ref } from 'vue'
import { useNotificationStore } from '../stores/notification'

export function useApi(apiFunc, options = {}) {
  const { showErrorNotification = true, showSuccessNotification = false } = options

  const data = ref(null)
  const loading = ref(false)
  const error = ref(null)

  const notificationStore = useNotificationStore()

  async function execute(...args) {
    loading.value = true
    error.value = null

    try {
      const response = await apiFunc(...args)
      data.value = response.data

      if (showSuccessNotification && options.successMessage) {
        notificationStore.success(options.successMessage)
      }

      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || err.message || 'An error occurred'

      if (showErrorNotification) {
        notificationStore.error(error.value)
      }

      throw err
    } finally {
      loading.value = false
    }
  }

  function reset() {
    data.value = null
    error.value = null
    loading.value = false
  }

  return {
    data,
    loading,
    error,
    execute,
    reset,
  }
}

export default useApi
