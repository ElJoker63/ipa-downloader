import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AppMetadata, AppDetailsOutput, SearchHistoryItem, Platform } from '../types'
import { WailsService } from '../services/wails'

export const useSearchStore = defineStore('search', () => {
  const query = ref('')
  const platform = ref<Platform>('ios')
  const results = ref<AppMetadata[]>([])
  const isLoading = ref(false)
  const errorMessage = ref<string | null>(null)

  const selectedApp = ref<AppDetailsOutput | null>(null)
  const isDetailsLoading = ref(false)
  const isDetailsModalOpen = ref(false)

  const searchHistory = ref<SearchHistoryItem[]>([])
  const limit = ref(15)

  async function search(searchQuery: string = query.value, targetPlatform: Platform = platform.value) {
    if (!searchQuery.trim()) {
      results.value = []
      return
    }

    query.value = searchQuery
    platform.value = targetPlatform
    isLoading.value = true
    errorMessage.value = null

    try {
      const data = await WailsService.searchApps(searchQuery.trim(), targetPlatform, limit.value)
      results.value = data || []
      fetchHistory()
    } catch (err: any) {
      errorMessage.value = err?.message || 'Failed to execute search'
      results.value = []
    } finally {
      isLoading.value = false
    }
  }

  async function openAppDetails(app: AppMetadata) {
    isDetailsModalOpen.value = true
    isDetailsLoading.value = true
    selectedApp.value = {
      metadata: app,
      versionHistory: [],
      isFavorite: app.isFavorite || false,
    }

    try {
      const details = await WailsService.getAppDetails(app.id, app.bundleId, platform.value)
      selectedApp.value = details
    } catch (err) {
      // Keep initial metadata if full lookup fails
    } finally {
      isDetailsLoading.value = false
    }
  }

  function closeAppDetails() {
    isDetailsModalOpen.value = false
    selectedApp.value = null
  }

  async function fetchHistory() {
    try {
      const items = await WailsService.getSearchHistory(15)
      searchHistory.value = items || []
    } catch (err) {
      // ignore
    }
  }

  async function clearHistory() {
    await WailsService.clearSearchHistory()
    searchHistory.value = []
  }

  return {
    query,
    platform,
    results,
    isLoading,
    errorMessage,
    selectedApp,
    isDetailsLoading,
    isDetailsModalOpen,
    searchHistory,
    limit,
    search,
    openAppDetails,
    closeAppDetails,
    fetchHistory,
    clearHistory,
  }
})
