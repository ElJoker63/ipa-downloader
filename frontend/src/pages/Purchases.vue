<template>
  <div class="max-w-6xl mx-auto space-y-6 flex flex-col h-full animate-slide-up font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4 shrink-0">
      <div>
        <div class="flex items-center space-x-3">
          <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
            {{ t.purchases?.title || 'Apps Compradas' }}
          </h1>
          <span
            v-if="purchasesStore.totalCount > 0"
            class="px-2.5 py-0.5 text-xs font-semibold rounded-full bg-[#0A84FF]/20 text-[#0A84FF] border border-[#0A84FF]/30"
          >
            {{ purchasesStore.totalCount }} {{ t.purchases?.totalCountLabel || 'disponibles' }}
          </span>
        </div>
        <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
          {{ t.purchases?.subtitle || 'Aplicaciones vinculadas a tu cuenta de Apple ID listas para descargar' }}
        </p>
      </div>

      <!-- Controls Header: Filter & Refresh -->
      <div class="flex items-center space-x-2.5 w-full md:w-auto">
        <!-- Search filter -->
        <div class="relative w-full md:w-64">
          <svg class="w-3.5 h-3.5 text-[#7D8592] absolute left-3.5 top-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="purchasesStore.searchQuery"
            type="text"
            :placeholder="t.purchases?.filterPlaceholder || 'Filtrar en esta página...'"
            class="glass-input w-full pl-9 pr-4 py-2 text-xs"
          />
        </div>

        <!-- Reload Button -->
        <button
          type="button"
          class="btn-secondary px-3 py-2 text-xs flex items-center space-x-1.5 shrink-0"
          :disabled="purchasesStore.isLoading"
          @click="refreshPurchases"
        >
          <svg
            class="w-3.5 h-3.5"
            :class="{ 'animate-spin': purchasesStore.isLoading }"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
            stroke-width="2"
          >
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span class="hidden sm:inline">{{ t.common.refresh || 'Actualizar' }}</span>
        </button>
      </div>
    </div>

    <!-- Content Area -->
    <div class="flex-1 min-h-0 overflow-y-auto">
      <!-- Loading State -->
      <div v-if="purchasesStore.isLoading" class="flex flex-col items-center justify-center h-64 space-y-3">
        <svg class="w-8 h-8 text-[#0A84FF] animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-xs text-[#B8C0CC]">{{ t.purchases?.loading || 'Cargando aplicaciones de tu cuenta...' }}</p>
      </div>

      <!-- Error State -->
      <div v-else-if="purchasesStore.error" class="glass-card !rounded-[24px] p-8 text-center space-y-3 max-w-md mx-auto mt-12 !border-red-500/30">
        <svg class="w-10 h-10 text-[#FF453A] mx-auto opacity-80" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <h3 class="text-sm font-semibold text-[#FFFFFF]">{{ t.common.failed }}</h3>
        <p class="text-xs text-[#FF453A]">{{ purchasesStore.error }}</p>
        <button type="button" class="btn-primary liquid-shine text-xs px-4 py-1.5 mt-2" @click="refreshPurchases">
          {{ t.common.retry }}
        </button>
      </div>

      <!-- Grid of Purchased Apps (Matching Search Cards) -->
      <div v-else-if="filteredApps.length > 0" class="space-y-4 pb-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div
            v-for="app in filteredApps"
            :key="app.id"
            class="glass-card glass-interactive !rounded-[20px] p-5 flex flex-col justify-between space-y-4"
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
                    :title="favoritesStore.isFavorite(app.id) ? t.common.delete : t.nav.favorites"
                    @click="toggleFav(app)"
                  >
                    <svg v-if="favoritesStore.isFavorite(app.id)" class="w-4 h-4 text-[#FF453A] fill-current" viewBox="0 0 24 24">
                      <path d="M12 21.35l-1.45-1.32C5.4 15.36 2 12.28 2 8.5 2 5.42 4.42 3 7.5 3c1.74 0 3.41.81 4.5 2.09C13.09 3.81 14.76 3 16.5 3 19.58 3 22 5.42 22 8.5c0 3.78-3.4 6.86-8.55 11.54L12 21.35z"/>
                    </svg>
                    <svg v-else class="w-4 h-4 text-[#7D8592] hover:text-[#FF453A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z"/>
                    </svg>
                  </button>
                </div>
                <p class="text-xs text-[#B8C0CC] truncate mt-0.5">{{ app.developer || app.bundleId }}</p>
                <div class="flex flex-wrap items-center gap-1.5 mt-2">
                  <span class="px-2 py-0.5 text-[10px] font-semibold rounded-full bg-[#30D158]/15 text-[#30D158] border border-[#30D158]/30">
                    ✓ {{ t.purchases?.licensed || 'Comprada' }}
                  </span>
                  <span class="px-2 py-0.5 text-[10px] font-mono rounded-md bg-white/[0.06] text-[#B8C0CC] border border-white/[0.08]">
                    v{{ app.version || 'Latest' }}
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
                class="px-4 py-1.5 rounded-xl bg-gradient-to-r from-[#30D158] to-[#28CD41] hover:from-[#28CD41] hover:to-[#30D158] backdrop-blur-xl text-white text-xs font-semibold shadow-[inset_0_1px_0_0_rgba(255,255,255,0.35),0_4px_16px_rgba(48,209,88,0.3)] liquid-shine flex items-center space-x-1.5 transition-all duration-300 ease-liquid active:scale-95"
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
                disabled
                class="px-3.5 py-1.5 rounded-xl bg-white/[0.08] backdrop-blur-xl text-[#30D158] border border-[#30D158]/30 text-xs font-medium flex items-center space-x-1.5 cursor-default opacity-80 select-none"
                :title="`Aplicación ya descargada en tu biblioteca (v${downloadedAppsStore.getDownloadedByBundleId(app.bundleId)?.version})`"
              >
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                <span>{{ t.downloadedApps?.downloaded || 'Descargado' }}</span>
              </button>

              <button
                v-else
                type="button"
                class="btn-primary liquid-shine text-xs px-4 py-1.5 flex items-center space-x-1.5"
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

        <!-- Pagination Controls -->
        <div class="flex items-center justify-between px-2 py-4 border-t border-white/[0.08] text-xs text-[#B8C0CC]">
          <div class="text-[11px]">
            {{ t.purchases?.showing || 'Mostrando página' }} {{ purchasesStore.currentPage }} {{ t.purchases?.of || 'de' }} {{ totalPages }}
          </div>
          <div class="flex items-center space-x-2">
            <button
              type="button"
              class="btn-secondary px-3 py-1.5 text-xs flex items-center space-x-1"
              :disabled="purchasesStore.currentPage <= 1 || purchasesStore.isLoading"
              @click="purchasesStore.setPage(purchasesStore.currentPage - 1)"
            >
              <span>← {{ t.purchases?.prev || 'Anterior' }}</span>
            </button>
            <span class="px-2 font-mono text-[11px]">{{ purchasesStore.currentPage }} / {{ totalPages }}</span>
            <button
              type="button"
              class="btn-secondary px-3 py-1.5 text-xs flex items-center space-x-1"
              :disabled="purchasesStore.currentPage >= totalPages || purchasesStore.isLoading"
              @click="purchasesStore.setPage(purchasesStore.currentPage + 1)"
            >
              <span>{{ t.purchases?.next || 'Siguiente' }} →</span>
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="glass-card !rounded-[26px] p-12 text-center space-y-3 max-w-lg mx-auto mt-12">
        <svg class="w-12 h-12 text-[#7D8592] mx-auto opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" />
        </svg>
        <h3 class="text-base font-semibold text-[#FFFFFF]">{{ t.purchases?.emptyTitle || 'No se encontraron compras' }}</h3>
        <p class="text-xs text-[#B8C0CC]">
          {{ t.purchases?.emptyDesc || 'Inicia sesión con tu Apple ID para sincronizar tu catálogo de aplicaciones adquiridas.' }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { usePurchasesStore } from '../stores/purchases'
import { useSearchStore } from '../stores/search'
import { useFavoritesStore } from '../stores/favorites'
import { useDownloadsStore } from '../stores/downloads'
import { useDownloadedAppsStore } from '../stores/downloadedApps'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'
import type { AppMetadata } from '../types'

const purchasesStore = usePurchasesStore()
const searchStore = useSearchStore()
const favoritesStore = useFavoritesStore()
const downloadsStore = useDownloadsStore()
const downloadedAppsStore = useDownloadedAppsStore()
const { t } = useI18n()
const { showToast } = useNotifications()

onMounted(async () => {
  if (purchasesStore.purchasedApps.length === 0) {
    await purchasesStore.fetchPurchases(1)
  }
})

const totalPages = computed(() => {
  if (purchasesStore.totalCount === 0) return 1
  return Math.ceil(purchasesStore.totalCount / purchasesStore.pageSize)
})

const filteredApps = computed(() => {
  if (!purchasesStore.searchQuery.trim()) return purchasesStore.purchasedApps
  const q = purchasesStore.searchQuery.toLowerCase()
  return purchasesStore.purchasedApps.filter(
    (a) => a.name.toLowerCase().includes(q) || a.bundleId.toLowerCase().includes(q)
  )
})

async function refreshPurchases() {
  await purchasesStore.fetchPurchases(purchasesStore.currentPage)
  showToast(t.value.purchases?.refreshed || 'Compras actualizadas', 'success')
}

function openDetails(app: AppMetadata) {
  searchStore.openAppDetails(app)
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
  showToast(favoritesStore.isFavorite(app.id) ? t.value.search.addedFav : t.value.search.removedFav, app.name, 'info')
}

async function downloadApp(app: AppMetadata) {
  try {
    await downloadsStore.queueDownload(app)
    showToast(t.value.search.downloadQueued || `${t.value.common.downloading} ${app.name}`, app.name, 'info')
  } catch (err: any) {
    showToast(t.value.search.downloadError || 'Error al iniciar descarga', err?.message || 'Error al iniciar descarga', 'error')
  }
}
</script>
