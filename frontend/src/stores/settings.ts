import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AppSettings } from '../types'
import { WailsService } from '../services/wails'
import { useI18n } from '../i18n'

export const useSettingsStore = defineStore('settings', () => {
  const settings = ref<AppSettings>({
    theme: 'dark',
    language: 'en',
    defaultDownloadFolder: '',
    maxConcurrentDownloads: 3,
    autoCheckUpdates: true,
    autoAcquireLicense: true,
    rememberCredentials: true,
    searchLimit: 15,
  })

  const cacheSize = ref('0 KB')
  const isLoading = ref(false)

  async function fetchSettings() {
    isLoading.value = true
    try {
      const data = await WailsService.getSettings()
      if (data) {
        settings.value = data
      }
      cacheSize.value = await WailsService.getCacheSize()
      applyTheme(settings.value.theme)

      // Update i18n if language is set in settings
      if (settings.value.language) {
        const { setLanguage } = useI18n()
        setLanguage(settings.value.language as any)
      }
    } finally {
      isLoading.value = false
    }
  }

  async function updateSettings(newSettings: Partial<AppSettings>) {
    settings.value = { ...settings.value, ...newSettings }
    await WailsService.saveSettings(settings.value)
    if (newSettings.theme) {
      applyTheme(newSettings.theme)
    }
  }

  async function browseFolder() {
    const selected = await WailsService.selectDownloadDirectory(settings.value.defaultDownloadFolder)
    if (selected) {
      settings.value.defaultDownloadFolder = selected
    }
  }

  async function clearCache() {
    await WailsService.clearAppCache()
    cacheSize.value = '0 KB'
  }

  async function exportLogs(path: string = '') {
    return WailsService.exportLogs(path)
  }

  function applyTheme(theme: string = 'dark') {
    const root = document.documentElement
    root.classList.add('dark')
    root.classList.remove('light')
  }

  return {
    settings,
    cacheSize,
    isLoading,
    fetchSettings,
    updateSettings,
    browseFolder,
    clearCache,
    exportLogs,
    applyTheme,
  }
})
