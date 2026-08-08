<template>
  <div class="max-w-6xl mx-auto space-y-6 animate-slide-up">
    <!-- Search Bar & Platform Selector -->
    <div class="glass-panel p-5 rounded-2xl border border-white/10 space-y-4">
      <div class="flex flex-col md:flex-row gap-4 items-center justify-between">
        <!-- Search Input -->
        <div class="relative flex-1 w-full">
          <span class="absolute inset-y-0 left-0 flex items-center pl-4 pointer-events-none text-slate-400">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
            </svg>
          </span>
          <input
            id="main-search-input"
            v-model="searchInput"
            type="text"
            placeholder="Search apps on the App Store (e.g., Spotify, WhatsApp, Slack, Fortnite)..."
            class="glass-input w-full pl-12 pr-10 py-3 rounded-xl text-sm font-medium"
            @input="onSearchInput"
            @keydown.enter="triggerImmediateSearch"
          />
          <button
            v-if="searchInput"
            type="button"
            class="absolute inset-y-0 right-0 flex items-center pr-3 text-slate-400 hover:text-white"
            @click="clearInput"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Platform Segmented Buttons -->
        <div class="flex p-1 rounded-xl bg-slate-900/60 border border-white/5 shrink-0">
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-lg text-xs font-semibold transition"
            :class="searchStore.platform === 'ios' ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-white'"
            @click="setPlatform('ios')"
          >
            📱 iPhone
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-lg text-xs font-semibold transition"
            :class="searchStore.platform === 'ipados' ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-white'"
            @click="setPlatform('ipados')"
          >
            💻 iPad
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-lg text-xs font-semibold transition"
            :class="searchStore.platform === 'tvos' ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-white'"
            @click="setPlatform('tvos')"
          >
            📺 Apple TV
          </button>
        </div>
      </div>

      <!-- Recent Search History Chips -->
      <div v-if="searchStore.searchHistory.length > 0 && !searchStore.query" class="flex flex-wrap items-center gap-2 pt-1">
        <span class="text-xs font-semibold text-slate-500 uppercase tracking-wider">Recent:</span>
        <button
          v-for="item in searchStore.searchHistory.slice(0, 8)"
          :key="item.id"
          type="button"
          class="px-2.5 py-1 rounded-lg text-xs bg-white/5 hover:bg-white/10 text-slate-300 transition flex items-center space-x-1.5 border border-white/5"
          @click="selectHistory(item.term, item.platform)"
        >
          <span>{{ item.term }}</span>
          <span class="text-[10px] text-slate-500">({{ item.platform }})</span>
        </button>
        <button
          type="button"
          class="text-xs text-slate-500 hover:text-slate-300 underline pl-1"
          @click="searchStore.clearHistory()"
        >
          Clear
        </button>
      </div>
    </div>

    <!-- Error Banner -->
    <div v-if="searchStore.errorMessage" class="p-4 rounded-xl bg-rose-500/10 border border-rose-500/20 text-sm text-rose-400 flex items-center justify-between">
      <span>{{ searchStore.errorMessage }}</span>
      <button type="button" class="btn-secondary text-xs px-3 py-1" @click="triggerImmediateSearch">Retry</button>
    </div>

    <!-- Loading Skeleton Grid -->
    <div v-if="searchStore.isLoading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div v-for="i in 6" :key="i" class="glass-panel p-5 rounded-2xl border border-white/5 skeleton h-48 space-y-4"></div>
    </div>

    <!-- Results Grid -->
    <div v-else-if="searchStore.results.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      <div
        v-for="app in searchStore.results"
        :key="app.id"
        class="glass-panel p-5 rounded-2xl border border-white/10 hover:border-white/20 transition-all duration-200 flex flex-col justify-between group hover:-translate-y-0.5 shadow-lg"
      >
        <div class="flex items-start space-x-4">
          <!-- App Artwork Icon -->
          <img
            :src="app.artworkUrl || app.artworkUrl512 || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
            :alt="app.name"
            class="w-16 h-16 rounded-2xl shadow-md object-cover bg-slate-800 border border-white/10 shrink-0 group-hover:scale-105 transition"
            loading="lazy"
          />

          <!-- App Metadata -->
          <div class="flex-1 min-w-0">
            <div class="flex items-start justify-between gap-1">
              <h3 class="text-base font-bold truncate text-slate-100 dark:text-slate-100 light:text-slate-900 group-hover:text-blue-400 transition" :title="app.name">
                {{ app.name }}
              </h3>
              <span class="px-2 py-0.5 text-[10px] font-semibold rounded bg-emerald-500/20 text-emerald-400 shrink-0">
                {{ app.formattedPrice || 'Free' }}
              </span>
            </div>

            <p class="text-xs text-slate-400 truncate mt-0.5" :title="app.developer">{{ app.developer }}</p>
            <p class="text-[11px] font-mono text-slate-500 truncate mt-0.5">{{ app.bundleId }}</p>

            <div class="flex items-center space-x-2 mt-2">
              <span class="px-2 py-0.5 text-[10px] font-mono rounded bg-white/5 text-slate-300">
                v{{ app.version }}
              </span>
              <span v-if="app.formattedSize" class="text-[10px] text-slate-400">
                {{ app.formattedSize }}
              </span>
            </div>
          </div>
        </div>

        <!-- Card Bottom Actions -->
        <div class="flex items-center justify-between pt-4 mt-3 border-t border-white/5">
          <button
            type="button"
            class="text-xs font-semibold text-blue-400 hover:text-blue-300 flex items-center space-x-1"
            @click="searchStore.openAppDetails(app)"
          >
            <span>View Details</span>
            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
            </svg>
          </button>

          <div class="flex items-center space-x-2">
            <!-- Favorite Toggle Button -->
            <button
              type="button"
              class="p-2 rounded-lg text-slate-400 hover:text-rose-400 hover:bg-rose-500/10 transition"
              :title="favoritesStore.isFavorite(app.id) ? 'Remove Favorite' : 'Save Favorite'"
              @click="toggleFavorite(app)"
            >
              <svg class="w-4 h-4" :class="favoritesStore.isFavorite(app.id) ? 'text-rose-500 fill-rose-500' : ''" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
              </svg>
            </button>

            <!-- Download Button -->
            <button
              type="button"
              class="btn-primary text-xs px-3.5 py-1.5 flex items-center space-x-1.5"
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
    </div>

    <!-- Empty State -->
    <div v-else-if="!searchStore.isLoading && searchInput" class="glass-panel p-12 rounded-2xl border border-white/10 text-center space-y-3">
      <div class="w-12 h-12 rounded-2xl bg-white/5 flex items-center justify-center mx-auto text-slate-400">
        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.172 16.172a4 4 0 015.656 0M9 10h.01M15 10h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>
      <h3 class="text-base font-bold">No Applications Found</h3>
      <p class="text-xs text-slate-400 max-w-sm mx-auto">
        No apps matched "{{ searchInput }}" for {{ searchStore.platform }}. Try another keyword or switch platforms.
      </p>
    </div>

    <!-- Initial Prompt -->
    <div v-else-if="!searchStore.isLoading && !searchInput" class="glass-panel p-12 rounded-2xl border border-white/10 text-center space-y-4">
      <div class="w-14 h-14 rounded-2xl bg-blue-500/10 text-blue-400 flex items-center justify-center mx-auto shadow-inner">
        <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
      </div>
      <div>
        <h3 class="text-lg font-bold">Search App Store Catalog</h3>
        <p class="text-xs text-slate-400 max-w-md mx-auto mt-1">
          Type an application name, developer, or bundle identifier to search live across millions of apps.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSearchStore } from '../stores/search'
import { useDownloadsStore } from '../stores/downloads'
import { useFavoritesStore } from '../stores/favorites'
import { useNotifications } from '../composables/useNotifications'
import type { AppMetadata, Platform } from '../types'

const router = useRouter()
const searchStore = useSearchStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const { showToast } = useNotifications()

const searchInput = ref(searchStore.query)
let debounceTimer: any = null

onMounted(() => {
  searchStore.fetchHistory()
  favoritesStore.fetchFavorites()
})

function onSearchInput() {
  clearTimeout(debounceTimer)
  debounceTimer = setTimeout(() => {
    searchStore.search(searchInput.value, searchStore.platform)
  }, 350)
}

function triggerImmediateSearch() {
  clearTimeout(debounceTimer)
  searchStore.search(searchInput.value, searchStore.platform)
}

function setPlatform(p: Platform) {
  searchStore.platform = p
  if (searchInput.value) {
    searchStore.search(searchInput.value, p)
  }
}

function selectHistory(term: string, p: string) {
  searchInput.value = term
  searchStore.platform = (p as Platform) || 'ios'
  searchStore.search(term, searchStore.platform)
}

function clearInput() {
  searchInput.value = ''
  searchStore.results = []
  searchStore.query = ''
}

async function toggleFavorite(app: AppMetadata) {
  const result = await favoritesStore.toggleFavorite(app)
  app.isFavorite = result
  showToast(result ? 'Added to Favorites' : 'Removed from Favorites', app.name, 'info')
}

async function downloadApp(app: AppMetadata) {
  try {
    await downloadsStore.queueDownload(app, searchStore.platform)
    showToast('Download Queued', `Starting download for ${app.name}`, 'success')
    router.push('/downloads')
  } catch (err: any) {
    showToast('Download Failed', err?.message || 'Could not queue download', 'error')
  }
}
</script>
