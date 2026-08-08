import { computed } from 'vue'
import { useSettingsStore } from '../stores/settings'

export function useTheme() {
  const settingsStore = useSettingsStore()

  const currentTheme = computed(() => 'dark')
  const isDark = computed(() => true)

  function applyTheme() {
    document.documentElement.classList.add('dark')
    document.documentElement.classList.remove('light')
  }

  function toggleTheme() {
    // Theme toggling is disabled, keeping dark mode only
  }

  function setTheme(theme: 'dark') {
    settingsStore.updateSettings({ theme: 'dark' })
    applyTheme()
  }

  return {
    currentTheme,
    isDark,
    toggleTheme,
    setTheme,
    applyTheme,
  }
}
