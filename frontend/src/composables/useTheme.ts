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

  function applyTheme() {
    if (isDark.value) {
      document.documentElement.classList.add('dark')
      document.documentElement.classList.remove('light')
    } else {
      document.documentElement.classList.add('light')
      document.documentElement.classList.remove('dark')
    }
  }

  function toggleTheme() {
    const next = isDark.value ? 'light' : 'dark'
    settingsStore.updateSettings({ theme: next })
    applyTheme()
  }

  function setTheme(theme: 'dark' | 'light' | 'system') {
    settingsStore.updateSettings({ theme })
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

