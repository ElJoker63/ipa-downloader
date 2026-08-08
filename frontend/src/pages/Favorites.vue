<template>
  <div class="max-w-6xl mx-auto space-y-6 flex flex-col h-full animate-slide-up font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4 shrink-0">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
          Saved Favorites
        </h1>
        <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
          Bookmarked applications stored in local SQLite database for instant one-click downloads.
        </p>
      </div>

      <!-- Search in Favorites Input -->
      <div class="relative w-full md:w-64">
        <svg class="w-3.5 h-3.5 text-[#7D8592] absolute left-3.5 top-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Filter favorites..."
          class="glass-input w-full pl-9 pr-4 py-2 text-xs"
        />
      </div>
    </div>

    <!-- Favorites Grid (30px Blur Glass Cards) -->
    <div class="flex-1 min-h-0 overflow-y-auto">
      <div v-if="filteredFavorites.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pb-6">
        <div
          v-for="app in filteredFavorites"
          :key="app.appId"
          class="glass-card p-5 rounded-[18px] flex flex-col justify-between space-y-4"
        >
          <!-- Card Header -->
          <div class="flex items-start space-x-3.5">
            <img
              :src="app.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
              :alt="app.name"
              class="w-14 h-14 rounded-[14px] object-cover bg-[#171A21] border border-white/[0.18] shadow-md shrink-0"
            />
            <div class="min-w-0 flex-1">
              <div class="flex items-center justify-between">
                <h3 class="text-sm font-semibold truncate text-[#FFFFFF]" :title="app.name">{{ app.name }}</h3>
                <button
                  type="button"
                  class="text-[#FF453A] hover:opacity-80 transition duration-150"
                  title="Remove from favorites"
                  @click="removeFavorite(app.appId, app.name)"
                >
                  <svg class="w-4 h-4 fill-current" viewBox="0 0 24 24">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                  </svg>
                </button>
              </div>
              <p class="text-xs text-[#B8C0CC] truncate mt-0.5">{{ app.developer }}</p>
              <div class="flex items-center space-x-2 mt-2">
                <span class="px-2 py-0.5 text-[10px] font-semibold rounded-full bg-[#0A84FF]/15 text-[#0A84FF] border border-[#0A84FF]/30">
                  {{ app.formattedPrice || 'Free' }}
                </span>
                <span class="px-2 py-0.5 text-[10px] font-mono rounded-md bg-white/[0.06] text-[#B8C0CC] border border-white/[0.08]">
                  v{{ app.version }}
                </span>
              </div>
            </div>
          </div>

          <!-- Card Actions Footer -->
          <div class="flex items-center justify-end space-x-2 pt-3 border-t border-white/[0.08]">
            <button
              type="button"
              class="btn-primary text-xs px-4 py-1.5 flex items-center space-x-1.5 w-full shadow-sm"
              @click="downloadFavorite(app)"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              <span>Download Package</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="glass-card p-12 rounded-[22px] text-center space-y-3 max-w-lg mx-auto mt-12">
        <svg class="w-12 h-12 text-[#FFD60A] mx-auto opacity-80" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
        </svg>
        <h3 class="text-base font-semibold text-[#FFFFFF]">No favorites saved yet</h3>
        <p class="text-xs text-[#B8C0CC]">
          Click the heart icon on any search result to bookmark apps here for quick access.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useFavoritesStore } from '../stores/favorites'
import { useDownloadsStore } from '../stores/downloads'
import { useNotifications } from '../composables/useNotifications'
import type { FavoriteApp, AppMetadata } from '../types'

const favoritesStore = useFavoritesStore()
const downloadsStore = useDownloadsStore()
const { showToast } = useNotifications()

const searchQuery = ref('')

onMounted(async () => {
  await favoritesStore.fetchFavorites()
})

const filteredFavorites = computed(() => {
  if (!searchQuery.value.trim()) return favoritesStore.favorites
  const q = searchQuery.value.toLowerCase()
  return favoritesStore.favorites.filter((f) => f.name.toLowerCase().includes(q) || f.bundleId.toLowerCase().includes(q))
})

async function removeFavorite(appId: number, name: string) {
  await favoritesStore.removeFavorite(appId)
  showToast('Removed from Favorites', name, 'info')
}

async function downloadFavorite(fav: FavoriteApp) {
  const meta: AppMetadata = {
    id: fav.appId,
    bundleId: fav.bundleId,
    name: fav.name,
    developer: fav.developer,
    version: fav.version,
    price: fav.price,
    formattedPrice: fav.formattedPrice,
    artworkUrl: fav.artworkUrl,
    primaryGenre: fav.primaryGenre || '',
    genres: [fav.primaryGenre || ''],
    minimumOsVersion: '12.0',
    fileSizeBytes: 0,
    formattedSize: '0 MB',
    averageUserRating: 0,
    userRatingCount: 0,
    contentAdvisoryRating: '',
    releaseDate: fav.createdAt ? String(fav.createdAt) : '',
    description: '',
    screenshots: [],
    supportedPlatforms: ['ios'],
    isFavorite: true,
  }

  try {
    await downloadsStore.queueDownload(meta, 'ios')
    showToast('Download Queued', `Starting download for ${fav.name}`, 'info')
  } catch (err: any) {
    showToast('Download Failed', err?.message || 'Failed to queue download', 'error')
  }
}
</script>
