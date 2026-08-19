<template>
  <div class="max-w-6xl mx-auto space-y-6 flex flex-col h-full animate-slide-up font-sans">
    <!-- Search Bar & Platform Selector Hero (8px Grid Spacing) -->
    <div class="space-y-4 shrink-0">
      <div class="flex flex-col md:flex-row items-center justify-between gap-4">
        <div>
          <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
            {{ t.search.title }}
          </h1>
          <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
            {{ t.search.subtitle }}
          </p>
        </div>

        <!-- Platform Segmented Pills (macOS / visionOS Style) -->
        <div class="flex items-center p-1 rounded-[14px] bg-white/[0.06] border border-white/[0.12] shrink-0 backdrop-blur-md">
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-[10px] text-xs font-medium transition-all duration-200 flex items-center space-x-1.5"
            :class="searchStore.platform === 'ios' ? 'bg-[#0A84FF] text-white shadow-sm shadow-[#0A84FF]/40' : 'text-[#B8C0CC] hover:text-white'"
            @click="changePlatform('ios')"
          >
            <span>{{ t.search.iphone }}</span>
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-[10px] text-xs font-medium transition-all duration-200 flex items-center space-x-1.5"
            :class="searchStore.platform === 'ipados' ? 'bg-[#0A84FF] text-white shadow-sm shadow-[#0A84FF]/40' : 'text-[#B8C0CC] hover:text-white'"
            @click="changePlatform('ipados')"
          >
            <span>{{ t.search.ipad }}</span>
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-[10px] text-xs font-medium transition-all duration-200 flex items-center space-x-1.5"
            :class="searchStore.platform === 'tvos' ? 'bg-[#0A84FF] text-white shadow-sm shadow-[#0A84FF]/40' : 'text-[#B8C0CC] hover:text-white'"
            @click="changePlatform('tvos')"
          >
            <span>{{ t.search.appleTv }}</span>
          </button>
        </div>
      </div>

      <!-- Hero Glass Search Input -->
      <div class="relative">
        <svg class="w-4 h-4 text-[#7D8592] absolute left-4 top-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="searchTerm"
          type="text"
          :placeholder="t.search.placeholder"
          class="glass-input w-full pl-11 pr-28 py-3 text-sm"
          @input="onSearchInput"
          @keydown.enter="searchStore.search(searchTerm)"
        />
        <div class="absolute right-3.5 top-3 flex items-center space-x-2">
          <span v-if="searchStore.isLoading" class="text-xs text-[#0A84FF] font-medium animate-pulse flex items-center space-x-1">
            <svg class="animate-spin h-3.5 w-3.5 text-[#0A84FF]" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
            </svg>
            <span>{{ t.search.searching }}</span>
          </span>
          <span v-else class="px-2 py-0.5 rounded-md bg-white/[0.08] text-[10px] font-mono text-[#7D8592] border border-white/[0.12]">Ctrl+K</span>
        </div>
      </div>

      <!-- Professional Filter Bar -->
      <div class="flex flex-wrap items-center gap-3 py-1">
        <!-- Region Selector -->
        <GlassDropdown
          v-model="searchStore.country"
          :options="countries"
          @change="searchStore.search(searchTerm)"
          class="min-w-[160px]"
        >
          <template #icon="{ selected }">
            <span v-if="selected" class="text-sm leading-none">{{ selected.flag }}</span>
          </template>
        </GlassDropdown>

        <!-- Category Dropdown -->
        <GlassDropdown
          v-model="searchStore.category"
          :options="categories"
          @change="searchStore.search(searchTerm)"
          class="min-w-[150px]"
        >
          <template #icon>
            <svg class="w-3.5 h-3.5 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
            </svg>
          </template>
        </GlassDropdown>

        <!-- Sort Dropdown -->
        <GlassDropdown
          v-model="searchStore.sortBy"
          :options="sortOptions"
          @change="searchStore.search(searchTerm)"
          class="min-w-[130px]"
        >
          <template #icon>
            <svg class="w-3.5 h-3.5 text-[#5E5CE6]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 4h13M3 8h9m-9 4h6m4 0l4-4m0 0l4 4m-4-4v12" />
            </svg>
          </template>
        </GlassDropdown>

        <div class="h-4 w-px bg-white/10 mx-1"></div>

        <!-- Result Limit -->
        <div class="flex items-center space-x-2 text-[10px] font-bold text-[#7D8592] uppercase tracking-wider">
          <span>{{ t.common.show }}:</span>
          <GlassDropdown
            v-model="searchStore.limit"
            :options="limitOptions"
            @change="searchStore.search(searchTerm)"
            class="min-w-[110px]"
          />
        </div>
      </div>



      <!-- Recent Searches Chips -->

      <div v-if="searchStore.searchHistory.length > 0 && searchStore.results.length === 0 && !searchStore.isLoading" class="flex flex-wrap items-center gap-2 pt-1">
        <span class="text-xs text-[#7D8592] font-medium">{{ t.search.recentSearches }}</span>
        <button
          v-for="item in searchStore.searchHistory.slice(0, 6)"
          :key="item.id"
          type="button"
          class="px-2.5 py-1 rounded-lg bg-white/[0.06] hover:bg-white/[0.12] border border-white/[0.08] text-xs text-[#B8C0CC] hover:text-white transition duration-150"
          @click="selectRecent(item.term)"
        >
          {{ item.term }}
        </button>
      </div>
    </div>

    <!-- Search Results Grid (30px Blur Glass Cards) -->
    <div class="flex-1 min-h-0 overflow-y-auto">
      <div v-if="searchStore.results.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pb-6">
        <div
          v-for="app in searchStore.results"
          :key="app.id"
          class="glass-card p-5 rounded-[18px] flex flex-col justify-between space-y-4"
        >
          <!-- App Header -->
          <div class="flex items-start space-x-3.5">
            <img
              :src="app.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
              :alt="app.name"
              class="w-16 h-16 rounded-[16px] object-cover bg-[#171A21] border border-white/[0.18] shadow-md shrink-0"
              loading="lazy"
            />
            <div class="min-w-0 flex-1">
              <div class="flex items-center justify-between gap-1">
                <h3 class="text-sm font-semibold truncate text-[#FFFFFF]" :title="app.name">{{ app.name }}</h3>
                <button
                  type="button"
                  class="text-sm text-[#B8C0CC] hover:text-[#FF453A] transition-transform duration-150 hover:scale-110"
                  :title="app.isFavorite ? 'Remove Favorite' : 'Add to Favorites'"
                  @click="toggleFav(app)"
                >
                  <svg v-if="app.isFavorite" class="w-4 h-4 text-[#FF453A] fill-current" viewBox="0 0 24 24">
                    <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                  </svg>
                  <svg v-else class="w-4 h-4 text-[#7D8592] hover:text-[#FF453A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                  </svg>
                </button>
              </div>
              <p class="text-xs text-[#B8C0CC] truncate mt-0.5">{{ app.developer }}</p>
              <div class="flex flex-wrap items-center gap-1.5 mt-2">
                <span class="px-2 py-0.5 text-[10px] font-semibold rounded-full bg-[#0A84FF]/15 text-[#0A84FF] border border-[#0A84FF]/30">
                  {{ app.formattedPrice || t.common.free }}
                </span>
                <span class="px-2 py-0.5 text-[10px] font-mono rounded-md bg-white/[0.06] text-[#B8C0CC] border border-white/[0.08]">
                  v{{ app.version }}
                </span>
                <span
                  v-if="downloadedAppsStore.isUpdateAvailable(app.bundleId, app.version)"
                  class="px-2 py-0.5 text-[10px] font-bold rounded-full bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30 animate-pulse"
                >
                  ↑ {{ t.downloadedApps?.updateAvailable || 'Nueva versión' }}
                </span>
                <span v-if="app.primaryGenre" class="text-[10px] text-[#7D8592] truncate max-w-[90px]">
                  {{ app.primaryGenre }}
                </span>
              </div>
            </div>
          </div>

          <!-- App Actions Footer -->
          <div class="flex items-center justify-between pt-3 border-t border-white/[0.08] gap-2">
            <button
              type="button"
              class="btn-secondary text-xs px-3 py-1.5 flex-1"
              @click="openDetails(app)"
            >
              {{ t.common.viewDetails }}
            </button>

            <!-- Action Buttons based on local downloaded state -->
            <button
              v-if="downloadedAppsStore.isUpdateAvailable(app.bundleId, app.version)"
              type="button"
              class="px-4 py-1.5 rounded-xl bg-gradient-to-r from-[#30D158] to-[#28CD41] hover:from-[#28CD41] hover:to-[#30D158] text-white text-xs font-semibold shadow-md shadow-[#30D158]/20 flex items-center space-x-1.5 transition-all duration-200"
              @click="downloadApp(app)"
              :title="`Actualizar a v${app.version}`"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
              <span>{{ t.downloadedApps?.update || 'Actualizar' }}</span>
            </button>

            <button
              v-else-if="downloadedAppsStore.getDownloadedByBundleId(app.bundleId)"
              type="button"
              class="px-3.5 py-1.5 rounded-xl bg-white/[0.08] hover:bg-white/[0.14] text-[#30D158] border border-[#30D158]/30 text-xs font-medium flex items-center space-x-1.5 transition"
              @click="downloadApp(app)"
              :title="`Ya descargado (v${downloadedAppsStore.getDownloadedByBundleId(app.bundleId)?.version}). Clic para volver a descargar.`"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
              </svg>
              <span>{{ t.downloadedApps?.downloaded || 'Descargado' }}</span>
            </button>

            <button
              v-else
              type="button"
              class="btn-primary text-xs px-4 py-1.5 flex items-center space-x-1.5 shadow-sm"
              @click="downloadApp(app)"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              <span>{{ t.common.download }}</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State Glass Panel -->
      <div v-else-if="!searchStore.isLoading" class="glass-card p-12 rounded-[22px] text-center space-y-3 max-w-lg mx-auto mt-12">
        <svg class="w-12 h-12 text-[#0A84FF] mx-auto opacity-80" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <h3 class="text-base font-semibold text-[#FFFFFF]">{{ t.search.emptyTitle }}</h3>
        <p class="text-xs text-[#B8C0CC]">
          {{ t.search.emptyDesc }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useSearchStore } from '../stores/search'
import { useDownloadsStore } from '../stores/downloads'
import { useFavoritesStore } from '../stores/favorites'
import { useDownloadedAppsStore } from '../stores/downloadedApps'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'
import GlassDropdown from '../components/GlassDropdown.vue'
import type { AppMetadata } from '../types'


const searchStore = useSearchStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const downloadedAppsStore = useDownloadedAppsStore()
const { t } = useI18n()
const { showToast } = useNotifications()

const searchTerm = ref('')

const countries = computed(() => [
  { id: 'US', name: 'United States', flag: '🇺🇸' },
  { id: 'ES', name: 'Spain', flag: '🇪🇸' },
  { id: 'GB', name: 'United Kingdom', flag: '🇬🇧' },
  { id: 'MX', name: 'Mexico', flag: '🇲🇽' },
  { id: 'AR', name: 'Argentina', flag: '🇦🇷' },
  { id: 'CO', name: 'Colombia', flag: '🇨🇴' },
  { id: 'CL', name: 'Chile', flag: '🇨🇱' },
  { id: 'BR', name: 'Brazil', flag: '🇧🇷' },
  { id: 'JP', name: 'Japan', flag: '🇯🇵' },
  { id: 'CN', name: 'China', flag: '🇨🇳' },
])

const categories = computed(() => [
  { id: '0', name: t.value.search.allCategories },
  { id: '6014', name: t.value.search.games },
  { id: '6002', name: t.value.search.utilities },
  { id: '6007', name: t.value.search.productivity },
  { id: '6017', name: t.value.search.education },
  { id: '6016', name: t.value.search.entertainment },
  { id: '6005', name: t.value.search.social },
  { id: '6015', name: t.value.search.finance },
  { id: '6013', name: t.value.search.health },
  { id: '6012', name: t.value.search.lifestyle },
  { id: '6011', name: t.value.search.music },
  { id: '6008', name: t.value.search.photoVideo },
  { id: '6000', name: t.value.search.business },
])

const sortOptions = computed(() => [
  { id: 'relevance', name: t.value.search.relevance },
  { id: 'popular', name: t.value.search.popularity },
  { id: 'rating', name: t.value.search.rating },
  { id: 'recent', name: t.value.search.releaseDate },
])

const limitOptions = computed(() => [
  { id: 15, name: `15 ${t.value.common.results}` },
  { id: 30, name: `30 ${t.value.common.results}` },
  { id: 50, name: `50 ${t.value.common.results}` },
])

let searchTimeout: any = null

onMounted(async () => {

  await searchStore.fetchHistory()
})

function onSearchInput() {
  if (searchTimeout) clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    searchStore.search(searchTerm.value)
  }, 400)
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
  showToast(app.isFavorite ? t.value.search.addedFav : t.value.search.removedFav, app.name, 'info')
}

function openDetails(app: AppMetadata) {
  searchStore.openAppDetails(app)
}

async function downloadApp(app: AppMetadata) {
  try {
    await downloadsStore.queueDownload(app, searchStore.platform)
    showToast(t.value.search.downloadQueued, app.name, 'info')
  } catch (err: any) {
    showToast(t.value.search.downloadError, err?.message || 'Failed to queue download', 'error')
  }
}
</script>
