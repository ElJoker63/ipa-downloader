import { computed } from 'vue'
import { useSettingsStore } from '../stores/settings'

export function useTheme() {
  const settingsStore = useSettingsStore()

  const currentTheme = computed(() => settingsStore.settings.theme)
  const isDark = computed(() => {
    if (settingsStore.settings.theme === 'dark') return true
    if (settingsStore.settings.theme === 'light') return false
    return window.matchMedia('(prefers-color-scheme: dark)').matches
  })

  function toggleTheme() {
    const next = isDark.value ? 'light' : 'dark'
    settingsStore.updateSettings({ theme: next })
  }

  function setTheme(theme: 'dark' | 'light' | 'system') {
    settingsStore.updateSettings({ theme })
  }

  return {
    currentTheme,
    isDark,
    toggleTheme,
    setTheme,
  }
}
