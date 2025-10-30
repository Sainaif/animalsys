<template>
  <div :class="['app', themeStore.isDark ? 'dark' : 'light']">
    <!-- Loading overlay -->
    <div v-if="loading" class="loading-overlay">
      <div class="loading-spinner"></div>
    </div>

    <!-- Main content -->
    <RouterView />

    <!-- Global notifications -->
    <NotificationContainer />
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue'
import { RouterView } from 'vue-router'
import { useThemeStore } from './stores/theme'
import { useAuthStore } from './stores/auth'
import NotificationContainer from './components/common/NotificationContainer.vue'

const themeStore = useThemeStore()
const authStore = useAuthStore()

const loading = computed(() => authStore.loading)

onMounted(async () => {
  // Initialize theme from localStorage
  themeStore.initTheme()

  // Try to restore session
  await authStore.initAuth()
})
</script>

<style scoped>
.app {
  min-height: 100vh;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  transition: background-color 0.3s, color 0.3s;
}

.loading-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
}

.loading-spinner {
  width: 50px;
  height: 50px;
  border: 4px solid var(--border-color);
  border-top-color: var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
