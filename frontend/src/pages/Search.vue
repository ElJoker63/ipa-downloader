<template>
  <div class="max-w-6xl mx-auto space-y-6 flex flex-col h-full animate-slide-up">
    <!-- Search Bar & Platform Selector Hero -->
    <div class="space-y-4 shrink-0">
      <div class="flex flex-col md:flex-row items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl font-extrabold tracking-tight text-white dark:text-white light:text-slate-900">
            App Store Search
          </h1>
          <p class="text-xs text-slate-400 dark:text-slate-400 light:text-slate-600 mt-0.5">
            Query the App Store in real time with high-resolution artwork, screenshots, and version builds.
          </p>
        </div>

        <!-- Platform Segmented Pills -->
        <div class="flex items-center p-1 rounded-2xl bg-white/5 dark:bg-white/5 light:bg-black/5 border border-white/10 shrink-0">
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all duration-200 flex items-center space-x-1.5"
            :class="searchStore.platform === 'ios' ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-md shadow-blue-600/30' : 'text-slate-400 hover:text-white'"
            @click="changePlatform('ios')"
          >
            <span>📱</span>
            <span>iPhone</span>
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all duration-200 flex items-center space-x-1.5"
            :class="searchStore.platform === 'ipados' ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-md shadow-blue-600/30' : 'text-slate-400 hover:text-white'"
            @click="changePlatform('ipados')"
          >
            <span>📱</span>
            <span>iPad</span>
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-xl text-xs font-bold transition-all duration-200 flex items-center space-x-1.5"
            :class="searchStore.platform === 'tvos' ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-md shadow-blue-600/30' : 'text-slate-400 hover:text-white'"
            @click="changePlatform('tvos')"
          >
            <span>📺</span>
            <span>Apple TV</span>
          </button>
        </div>
      </div>

      <!-- Hero Search Input -->
      <div class="relative">
        <span class="absolute left-4 top-3.5 text-slate-400 text-lg">🔍</span>
        <input
          v-model="searchTerm"
          type="text"
          placeholder="Search by app name, developer, or keyword... (e.g. Spotify, Telegram, Minecraft)"
          class="glass-input w-full pl-12 pr-28 py-3 rounded-2xl text-sm font-medium outline-none shadow-lg"
          @input="onSearchInput"
          @keydown.enter="searchStore.search(searchTerm)"
        />
        <div class="absolute right-3.5 top-3 flex items-center space-x-2">
          <span v-if="searchStore.isLoading" class="text-xs text-blue-400 font-semibold animate-pulse flex items-center space-x-1">
            <svg class="animate-spin h-3.5 w-3.5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
            </svg>
            <span>Searching...</span>
          </span>
          <span v-else class="px-2 py-0.5 rounded-md bg-white/10 text-[10px] font-mono text-slate-400">Ctrl+K</span>
        </div>
      </div>

      <!-- Recent Searches Chips -->
      <div v-if="searchStore.searchHistory.length > 0 && searchStore.results.length === 0 && !searchStore.isLoading" class="flex flex-wrap items-center gap-2 pt-1">
        <span class="text-xs text-slate-500 font-semibold">Recent searches:</span>
        <button
          v-for="item in searchStore.searchHistory.slice(0, 6)"
          :key="item.id"
          type="button"
          class="px-2.5 py-1 rounded-lg bg-white/5 hover:bg-white/10 border border-white/5 text-xs text-slate-300 transition"
          @click="selectRecent(item.term)"
        >
          {{ item.term }}
        </button>
      </div>
    </div>

    <!-- Search Results Grid -->
    <div class="flex-1 min-h-0 overflow-y-auto">
      <div v-if="searchStore.results.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pb-6">
        <div
          v-for="app in searchStore.results"
          :key="app.id"
          class="glass-card-interactive p-4 rounded-2xl border border-white/10 flex flex-col justify-between space-y-4"
        >
          <!-- App Header -->
          <div class="flex items-start space-x-3.5">
            <img
              :src="app.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
              :alt="app.name"
              class="w-16 h-16 rounded-2xl object-cover bg-slate-800 border border-white/10 shadow-lg shadow-black/40 shrink-0"
              loading="lazy"
            />
            <div class="min-w-0 flex-1">
              <div class="flex items-center justify-between gap-1">
                <h3 class="text-sm font-bold truncate text-white dark:text-white light:text-slate-900" :title="app.name">{{ app.name }}</h3>
                <button
                  type="button"
                  class="text-sm transition-transform hover:scale-125"
                  :title="app.isFavorite ? 'Remove Favorite' : 'Add to Favorites'"
                  @click="toggleFav(app)"
                >
                  <span>{{ app.isFavorite ? '❤️' : '🤍' }}</span>
                </button>
              </div>
              <p class="text-xs text-slate-400 truncate mt-0.5">{{ app.developer }}</p>
              <div class="flex flex-wrap items-center gap-1.5 mt-2">
                <span class="px-2 py-0.5 text-[10px] font-bold rounded-full bg-blue-500/20 text-blue-400 border border-blue-500/30">
                  {{ app.formattedPrice || 'Free' }}
                </span>
                <span class="px-2 py-0.5 text-[10px] font-mono rounded bg-white/5 text-slate-300 border border-white/5">
                  v{{ app.version }}
                </span>
                <span v-if="app.primaryGenre" class="text-[10px] text-slate-400 truncate max-w-[90px]">
                  {{ app.primaryGenre }}
                </span>
              </div>
            </div>
          </div>

          <!-- App Actions Footer -->
          <div class="flex items-center justify-between pt-3 border-t border-white/5 gap-2">
            <button
              type="button"
              class="btn-secondary text-xs px-3 py-1.5 flex-1"
              @click="openDetails(app)"
            >
              View Details
            </button>
            <button
              type="button"
              class="btn-primary text-xs px-4 py-1.5 flex items-center space-x-1.5 shadow-md"
              @click="downloadApp(app)"
            >
              <span>⬇</span>
              <span>Download</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="!searchStore.isLoading" class="glass-panel p-12 rounded-3xl border border-white/5 text-center space-y-3 max-w-lg mx-auto mt-12">
        <div class="text-4xl">🔍</div>
        <h3 class="text-base font-bold text-white">Search Apple App Store</h3>
        <p class="text-xs text-slate-400">
          Type any application name to search across the entire iOS, iPadOS, or tvOS catalog in real time.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useSearchStore } from '../stores/search'
import { useDownloadsStore } from '../stores/downloads'
import { useFavoritesStore } from '../stores/favorites'
import { useNotifications } from '../composables/useNotifications'
import type { AppMetadata } from '../types'

const searchStore = useSearchStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const { showToast } = useNotifications()

const searchTerm = ref('')

onMounted(async () => {
  await searchStore.fetchHistory()
})

function onSearchInput() {
  searchStore.search(searchTerm.value)
}

function selectRecent(term: string) {
  searchTerm.value = term
  searchStore.search(term)
}

function changePlatform(p: 'ios' | 'ipados' | 'tvos') {
  searchStore.platform = p
  if (searchTerm.value) {
    searchStore.search(searchTerm.value)
  }
}

async function toggleFav(app: AppMetadata) {
  await favoritesStore.toggleFavorite({
    appId: app.id,
    bundleId: app.bundleId,
    name: app.name,
    developer: app.developer,
    version: app.version,
    price: app.price,
    formattedPrice: app.formattedPrice,
    artworkUrl: app.artworkUrl,
    primaryGenre: app.primaryGenre,
    createdAt: new Date().toISOString(),
  })
  app.isFavorite = !app.isFavorite
  showToast(app.isFavorite ? 'Added to Favorites' : 'Removed from Favorites', app.name, 'info')
}

function openDetails(app: AppMetadata) {
  searchStore.openAppDetails(app)
}

async function downloadApp(app: AppMetadata) {
  try {
    await downloadsStore.queueDownload(app, searchStore.platform)
    showToast('Download Queued', `Starting download for ${app.name}`, 'info')
  } catch (err: any) {
    showToast('Download Error', err?.message || 'Failed to queue download', 'error')
  }
}
</script>
