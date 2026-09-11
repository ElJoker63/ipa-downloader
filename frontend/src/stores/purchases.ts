import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { AppMetadata } from '../types'
import { WailsService } from '../services/wails'

export const usePurchasesStore = defineStore('purchases', () => {
  const purchasedApps = ref<AppMetadata[]>([])
  const currentPage = ref(1)
  const pageSize = ref(20)
  const totalCount = ref(0)
  const isLoading = ref(false)
  const searchQuery = ref('')
  const error = ref<string | null>(null)

  async function fetchPurchases(page: number = currentPage.value) {
    isLoading.value = true
    error.value = null
    try {
      const output = await WailsService.getPurchasedApps(page, pageSize.value)
      purchasedApps.value = output.results || []
      totalCount.value = output.totalCount
      currentPage.value = output.page
    } catch (err: any) {
      error.value = err?.message || 'Failed to load purchased apps'
      console.error('Error loading purchases:', err)
    } finally {
      isLoading.value = false
    }
  }

  function setPage(page: number) {
    if (page >= 1 && (totalCount.value === 0 || page <= Math.ceil(totalCount.value / pageSize.value))) {
      currentPage.value = page
      fetchPurchases(page)
    }
  }

  return {
    purchasedApps,
    currentPage,
    pageSize,
    totalCount,
    isLoading,
    searchQuery,
    error,
    fetchPurchases,
    setPage,
  }
})
