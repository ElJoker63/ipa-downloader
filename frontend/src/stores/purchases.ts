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
  // True while a background refresh (silent, cache already shown) is in flight.
  // Kept separate from isLoading so the UI can show a subtle indicator instead
  // of the full-page spinner while this runs.
  const isRefreshing = ref(false)
  const hasLoadedOnce = ref(false)
  const searchQuery = ref('')
  const error = ref<string | null>(null)

  function applyOutput(output: { results?: AppMetadata[]; totalCount: number; page: number }) {
    purchasedApps.value = output.results || []
    totalCount.value = output.totalCount
    currentPage.value = output.page
  }

  // Live fetch from Apple. When silent is true (background refresh behind an
  // already-shown cache) failures are only logged, never surfaced as an error
  // banner or a loading spinner — the cache stays on screen untouched.
  async function fetchPurchases(page: number = currentPage.value, opts: { silent?: boolean } = {}) {
    const silent = opts.silent === true

    if (silent) {
      isRefreshing.value = true
    } else {
      isLoading.value = true
      error.value = null
    }

    try {
      const output = await WailsService.getPurchasedApps(page, pageSize.value)
      applyOutput(output)
      hasLoadedOnce.value = true
    } catch (err: any) {
      console.error('Error loading purchases:', err)
      if (!silent) {
        error.value = err?.message || 'Failed to load purchased apps'
      }
    } finally {
      if (silent) {
        isRefreshing.value = false
      } else {
        isLoading.value = false
      }
    }
  }

  // Cache-first entry point: shows whatever is cached locally immediately (no
  // spinner, no network wait), then quietly refreshes from Apple in the
  // background in case something changed. Falls back to a normal foreground
  // fetch only when nothing is cached yet (e.g. first launch ever).
  async function loadPurchases(page: number = currentPage.value) {
    const cached = await WailsService.getCachedPurchasedApps(page, pageSize.value)

    if (cached) {
      applyOutput(cached)
      hasLoadedOnce.value = true
      isLoading.value = false
      error.value = null
      fetchPurchases(page, { silent: true })
      return
    }

    await fetchPurchases(page)
  }

  function setPage(page: number) {
    if (page >= 1 && (totalCount.value === 0 || page <= Math.ceil(totalCount.value / pageSize.value))) {
      loadPurchases(page)
    }
  }

  return {
    purchasedApps,
    currentPage,
    pageSize,
    totalCount,
    isLoading,
    isRefreshing,
    hasLoadedOnce,
    searchQuery,
    error,
    fetchPurchases,
    loadPurchases,
    setPage,
  }
})
