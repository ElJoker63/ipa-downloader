import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { FavoriteApp, AppMetadata } from '../types'
import { WailsService } from '../services/wails'

export const useFavoritesStore = defineStore('favorites', () => {
  const favorites = ref<FavoriteApp[]>([])
  const searchQuery = ref('')
  const isLoading = ref(false)

  async function fetchFavorites() {
    isLoading.value = true
    try {
      if (searchQuery.value.trim()) {
        favorites.value = await WailsService.searchFavorites(searchQuery.value.trim())
      } else {
        favorites.value = await WailsService.getFavorites()
      }
    } finally {
      isLoading.value = false
    }
  }

  async function toggleFavorite(app: AppMetadata | FavoriteApp): Promise<boolean> {
    const isFavApp = 'appId' in app
    const favItem: FavoriteApp = {
      appId: isFavApp ? app.appId : (app as AppMetadata).id,
      bundleId: isFavApp ? app.bundleId : (app as AppMetadata).bundleId,
      name: isFavApp ? app.name : (app as AppMetadata).name,
      developer: isFavApp ? app.developer : (app as AppMetadata).developer,
      version: isFavApp ? app.version : (app as AppMetadata).version,
      price: isFavApp ? app.price : (app as AppMetadata).price,
      formattedPrice: isFavApp ? app.formattedPrice : (app as AppMetadata).formattedPrice,
      artworkUrl: isFavApp ? app.artworkUrl : (app as AppMetadata).artworkUrl,
      primaryGenre: isFavApp ? app.primaryGenre : (app as AppMetadata).primaryGenre,
      createdAt: new Date().toISOString(),
    }

    const newState = await WailsService.toggleFavorite(favItem)
    await fetchFavorites()
    return newState
  }

  async function removeFavorite(appId: number) {
    await WailsService.removeFavorite(appId)
    favorites.value = favorites.value.filter((f) => f.appId !== appId)
  }

  function isFavorite(appId: number): boolean {
    return favorites.value.some((f) => f.appId === appId)
  }

  function initListeners() {
    WailsService.onEvent('favorites:updated', () => {
      fetchFavorites()
    })
  }

  return {
    favorites,
    searchQuery,
    isLoading,
    fetchFavorites,
    toggleFavorite,
    removeFavorite,
    isFavorite,
    initListeners,
  }
})
