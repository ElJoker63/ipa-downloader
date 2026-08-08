import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { DownloadTask } from '../types'
import { WailsService } from '../services/wails'

export const useHistoryStore = defineStore('history', () => {
  const history = ref<DownloadTask[]>([])
  const filterStatus = ref<'all' | 'completed' | 'failed' | 'cancelled'>('all')
  const searchQuery = ref('')
  const isLoading = ref(false)

  const filteredHistory = computed(() => {
    return history.value.filter((item) => {
      const matchesStatus =
        filterStatus.value === 'all' || item.status === filterStatus.value

      const query = searchQuery.value.toLowerCase().trim()
      const matchesSearch =
        !query ||
        item.appName.toLowerCase().includes(query) ||
        item.bundleId.toLowerCase().includes(query) ||
        item.developer.toLowerCase().includes(query) ||
        item.destinationPath.toLowerCase().includes(query)

      return matchesStatus && matchesSearch
    })
  })

  async function fetchHistory() {
    isLoading.value = true
    try {
      history.value = (await WailsService.getDownloadHistory()) || []
    } finally {
      isLoading.value = false
    }
  }

  async function deleteItem(id: string) {
    await WailsService.deleteHistoryItem(id)
    history.value = history.value.filter((h) => h.id !== id)
  }

  async function clearHistory() {
    await WailsService.clearDownloadHistory()
    history.value = []
  }

  async function openFolder(path: string) {
    await WailsService.openFolder(path)
  }

  async function openFile(path: string) {
    await WailsService.openFile(path)
  }

  async function revealInExplorer(path: string) {
    await WailsService.revealInExplorer(path)
  }

  async function copyPath(path: string) {
    try {
      await navigator.clipboard.writeText(path)
    } catch {
      // fallback
    }
  }

  return {
    history,
    filterStatus,
    searchQuery,
    filteredHistory,
    isLoading,
    fetchHistory,
    deleteItem,
    clearHistory,
    openFolder,
    openFile,
    revealInExplorer,
    copyPath,
  }
})
