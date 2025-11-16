import { ref } from 'vue'

const THEME_STORAGE_KEY = 'theme'
const theme = ref(localStorage.getItem(THEME_STORAGE_KEY) || 'light')

const applyTheme = (value) => {
  const normalized = value === 'dark' ? 'dark' : 'light'
  document.documentElement.setAttribute('data-theme', normalized)
  localStorage.setItem(THEME_STORAGE_KEY, normalized)
}

applyTheme(theme.value)

const useTheme = () => {
  const setTheme = (value) => {
    theme.value = value
    applyTheme(theme.value)
  }

  const toggleTheme = () => {
    setTheme(theme.value === 'light' ? 'dark' : 'light')
  }

  return {
    theme,
    toggleTheme,
    setTheme
  }
}

export { theme, applyTheme }
export default useTheme
