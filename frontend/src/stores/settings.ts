import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AppSettings } from '../types'
import { WailsService } from '../services/wails'

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

  function applyTheme(theme: string) {
    const root = document.documentElement
    if (theme === 'dark') {
      root.classList.add('dark')
      root.classList.remove('light')
    } else if (theme === 'light') {
      root.classList.add('light')
      root.classList.remove('dark')
    } else {
      // System
      const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches
      if (prefersDark) {
        root.classList.add('dark')
        root.classList.remove('light')
      } else {
        root.classList.add('light')
        root.classList.remove('dark')
      }
    }
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
