import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { LogEntry } from '../types'
import { WailsService } from '../services/wails'

export const useLogsStore = defineStore('logs', () => {
  const logs = ref<LogEntry[]>([])
  const filterLevel = ref<string>('ALL')
  const searchQuery = ref('')
  const autoScroll = ref(true)
  const isLoading = ref(false)

  const filteredLogs = computed(() => {
    return logs.value.filter((item) => {
      const matchesLevel =
        filterLevel.value === 'ALL' || item.level.toUpperCase() === filterLevel.value.toUpperCase()

      const query = searchQuery.value.toLowerCase().trim()
      const matchesSearch =
        !query ||
        item.message.toLowerCase().includes(query) ||
        (item.context && item.context.toLowerCase().includes(query))

      return matchesLevel && matchesSearch
    })
  })

  async function fetchLogs() {
    isLoading.value = true
    try {
      logs.value = (await WailsService.getLogs(500)) || []
    } finally {
      isLoading.value = false
    }
  }

  async function clearLogs() {
    await WailsService.clearLogs()
    logs.value = []
  }

  async function exportLogs(path: string = '') {
    return WailsService.exportLogs(path)
  }

  function initListeners() {
    WailsService.onEvent('log:entry', (entry: LogEntry) => {
      if (entry) {
        logs.value.push(entry)
        if (logs.value.length > 2000) {
          logs.value.shift()
        }
      }
    })
  }

  return {
    logs,
    filterLevel,
    searchQuery,
    autoScroll,
    filteredLogs,
    isLoading,
    fetchLogs,
    clearLogs,
    exportLogs,
    initListeners,
  }
})
