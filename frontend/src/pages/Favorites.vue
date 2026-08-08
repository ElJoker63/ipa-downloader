<template>
  <div class="max-w-6xl mx-auto space-y-6 animate-slide-up">
    <!-- Header with Search within Favorites -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold">Saved Favorites</h1>
        <p class="text-xs text-slate-400 mt-0.5">Quickly access and download your favorite iOS and tvOS apps bookmarked locally.</p>
      </div>

      <!-- Search in Favorites -->
      <div class="relative w-full md:w-72">
        <span class="absolute inset-y-0 left-0 flex items-center pl-3.5 pointer-events-none text-slate-400">
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </span>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Filter favorites..."
          class="glass-input w-full pl-10 pr-4 py-2 rounded-xl text-xs font-medium"
          @input="onSearch"
        />
      </div>
    </div>

    <!-- Favorites Grid -->
    <div v-if="favoritesStore.favorites.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="app in favoritesStore.favorites"
        :key="app.appId"
        class="glass-panel p-5 rounded-2xl border border-white/10 hover:border-white/20 transition duration-200 flex flex-col justify-between group shadow-lg"
      >
        <div class="flex items-start space-x-4">
          <img
            :src="app.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
            :alt="app.name"
            class="w-14 h-14 rounded-2xl shadow-md object-cover bg-slate-800 border border-white/10 shrink-0"
          />

          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-1">
              <h3 class="text-sm font-bold truncate text-slate-100 dark:text-slate-100 light:text-slate-900" :title="app.name">
                {{ app.name }}
              </h3>
              <span class="px-2 py-0.5 text-[10px] font-semibold rounded bg-emerald-500/20 text-emerald-400 shrink-0">
                {{ app.formattedPrice || 'Free' }}
              </span>
            </div>

            <p class="text-xs text-slate-400 truncate mt-0.5">{{ app.developer }}</p>
            <p class="text-[11px] font-mono text-slate-500 truncate mt-0.5">{{ app.bundleId }}</p>
          </div>
        </div>

        <div class="flex items-center justify-between pt-4 mt-3 border-t border-white/5">
          <button
            type="button"
            class="text-xs text-rose-400 hover:text-rose-300 transition flex items-center space-x-1"
            @click="removeFavorite(app.appId)"
          >
            <span>Remove</span>
          </button>

          <button
            type="button"
            class="btn-primary text-xs px-3 py-1.5 flex items-center space-x-1.5"
            @click="downloadApp(app)"
          >
            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            <span>Download</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="glass-panel p-12 rounded-2xl border border-white/10 text-center space-y-3">
      <div class="w-12 h-12 rounded-2xl bg-white/5 flex items-center justify-center mx-auto text-amber-400 text-2xl">
        ⭐
      </div>
      <h3 class="text-base font-bold">No Saved Favorites</h3>
      <p class="text-xs text-slate-400 max-w-sm mx-auto">
        Search for apps on the Search tab and click the heart icon to save them for quick access.
      </p>
      <div class="pt-2">
        <router-link to="/search" class="btn-primary text-xs px-4 py-2 inline-block">
          Explore Apps
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useFavoritesStore } from '../stores/favorites'
import { useDownloadsStore } from '../stores/downloads'
import { useNotifications } from '../composables/useNotifications'
import type { FavoriteApp, AppMetadata } from '../types'

const router = useRouter()
const favoritesStore = useFavoritesStore()
const downloadsStore = useDownloadsStore()
const { showToast } = useNotifications()

const searchQuery = ref('')

onMounted(() => {
  favoritesStore.fetchFavorites()
})

function onSearch() {
  favoritesStore.searchQuery = searchQuery.value
  favoritesStore.fetchFavorites()
}

async function removeFavorite(appId: number) {
  await favoritesStore.removeFavorite(appId)
  showToast('Removed', 'App removed from favorites', 'info')
}

async function downloadApp(fav: FavoriteApp) {
  const meta: AppMetadata = {
    id: fav.appId,
    bundleId: fav.bundleId,
    name: fav.name,
    developer: fav.developer,
    version: fav.version,
    price: fav.price,
    formattedPrice: fav.formattedPrice,
    artworkUrl: fav.artworkUrl,
    screenshots: [],
    description: '',
    releaseDate: '',
    minimumOsVersion: '12.0',
    fileSizeBytes: 0,
    formattedSize: '',
    averageUserRating: 0,
    userRatingCount: 0,
    contentAdvisoryRating: '',
    genres: [],
    primaryGenre: fav.primaryGenre || 'Application',
    supportedPlatforms: ['iOS'],
    isFavorite: true,
  }

  try {
    await downloadsStore.queueDownload(meta, 'ios')
    showToast('Download Queued', `Starting download for ${meta.name}`, 'success')
    router.push('/downloads')
  } catch (err: any) {
    showToast('Download Failed', err?.message || 'Could not queue download', 'error')
  }
}
</script>
