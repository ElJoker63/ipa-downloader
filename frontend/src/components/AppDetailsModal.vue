<template>
  <div
    v-if="searchStore.isDetailsModalOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md transition-all duration-200"
    @click.self="closeModal"
  >
    <div
      v-if="app"
      class="glass-card w-full max-w-4xl max-h-[88vh] rounded-[22px] border border-white/[0.18] shadow-[0_16px_48px_rgba(0,0,0,0.45)] overflow-hidden flex flex-col animate-modal font-sans"
    >
      <!-- Modal Header -->
      <div class="p-6 border-b border-white/[0.08] flex items-start justify-between gap-4 shrink-0 bg-white/[0.02]">
        <div class="flex items-start space-x-4 min-w-0">
          <img
            :src="app.artworkUrl512 || app.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
            :alt="app.name"
            class="w-20 h-20 rounded-[18px] object-cover bg-[#171A21] border border-white/[0.18] shadow-lg shrink-0"
          />
          <div class="min-w-0">
            <h2 class="text-xl font-bold truncate text-[#FFFFFF]">{{ app.name }}</h2>
            <p class="text-xs text-[#B8C0CC] truncate mt-0.5">{{ app.developer }}</p>
            <div class="flex flex-wrap items-center gap-2 mt-2">
              <span class="px-2.5 py-0.5 text-xs font-semibold rounded-full bg-[#0A84FF]/20 text-[#0A84FF] border border-[#0A84FF]/30">
                {{ app.formattedPrice || t.common.free }}
              </span>
              <span class="px-2.5 py-0.5 text-xs font-mono rounded-md bg-white/[0.08] text-[#B8C0CC] border border-white/[0.12]">
                v{{ app.version }}
              </span>
              <span v-if="app.formattedSize" class="text-xs text-[#7D8592] font-mono">
                {{ app.formattedSize }}
              </span>
              <span v-if="app.primaryGenre" class="text-xs text-[#7D8592]">
                {{ app.primaryGenre }}
              </span>
            </div>
          </div>
        </div>

        <!-- Header Actions (Close & Download/Update) -->
        <div class="flex items-center space-x-2 shrink-0">
          <button
            v-if="app && downloadedAppsStore.isUpdateAvailable(app.bundleId, app.version)"
            type="button"
            class="px-4 py-2 rounded-xl bg-gradient-to-r from-[#30D158] to-[#28CD41] hover:from-[#28CD41] hover:to-[#30D158] text-white text-xs font-semibold shadow-md shadow-[#30D158]/20 flex items-center space-x-1.5 transition-all duration-200"
            @click="downloadCurrent"
          >
            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            <span>{{ t.downloadedApps?.update || 'Actualizar' }}</span>
          </button>

          <button
            v-else-if="app && downloadedAppsStore.getDownloadedByBundleId(app.bundleId)"
            type="button"
            disabled
            class="px-4 py-2 rounded-xl bg-white/[0.08] text-[#30D158] border border-[#30D158]/30 text-xs font-medium flex items-center space-x-1.5 cursor-default opacity-80 select-none"
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
            class="btn-primary text-xs px-4 py-2 flex items-center space-x-1.5 shadow-sm"
            @click="downloadCurrent"
          >
            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            <span>{{ t.common.download }}</span>
          </button>

          <button
            type="button"
            class="p-2 rounded-xl text-[#7D8592] hover:text-white hover:bg-white/[0.08] transition duration-150"
            @click="closeModal"
          >
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>

      <!-- Modal Body (Scrollable) -->
      <div class="p-6 overflow-y-auto space-y-6 flex-1">
        <!-- Screenshots Lightbox Gallery -->
        <div v-if="screenshots.length > 0" class="space-y-2">
          <div class="flex items-center justify-between">
            <h3 class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">{{ t.details.screenshots }}</h3>
            <span class="text-[11px] text-[#7D8592] font-normal">{{ t.common.clickToEnlarge }}</span>
          </div>
          <div class="flex space-x-3 overflow-x-auto pb-3 pt-1">
            <div
              v-for="(img, idx) in screenshots"
              :key="idx"
              class="group relative rounded-[14px] overflow-hidden border border-white/[0.12] shadow-md hover:border-[#0A84FF]/60 transition-all duration-200 cursor-pointer shrink-0"
              @click="openLightbox(idx)"
            >
              <img
                :src="img"
                class="h-64 object-cover bg-[#171A21] group-hover:scale-[1.03] transition-transform duration-200"
                loading="lazy"
              />
              <div class="absolute inset-0 bg-black/30 opacity-0 group-hover:opacity-100 flex items-center justify-center transition-opacity duration-150 backdrop-blur-[2px]">
                <div class="p-2 rounded-full bg-white/20 text-white border border-white/30 shadow-lg">
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0zM10 7v6m3-3H7" />
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Description -->
        <div v-if="app.description" class="space-y-2">
          <h3 class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">{{ t.details.description }}</h3>
          <div class="p-4 rounded-[14px] bg-white/[0.04] border border-white/[0.08] text-xs text-[#B8C0CC] leading-relaxed whitespace-pre-line max-h-48 overflow-y-auto font-normal">
            {{ app.description }}
          </div>
        </div>

        <!-- Metadata Properties Grid -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs">
          <div class="p-3 rounded-[12px] bg-white/[0.04] border border-white/[0.08]">
            <span class="text-[#7D8592] block">{{ t.details.bundleId }}</span>
            <span class="font-mono text-[#FFFFFF] truncate block mt-0.5">{{ app.bundleId }}</span>
          </div>

          <div class="p-3 rounded-[12px] bg-white/[0.04] border border-white/[0.08]">
            <span class="text-[#7D8592] block">{{ t.details.minOs }}</span>
            <span class="font-mono text-[#FFFFFF] block mt-0.5">iOS {{ app.minimumOsVersion || '12.0' }}</span>
          </div>

          <div class="p-3 rounded-[12px] bg-white/[0.04] border border-white/[0.08]">
            <span class="text-[#7D8592] block">{{ t.details.releaseDate }}</span>
            <span class="text-[#FFFFFF] block mt-0.5">{{ formatDate(app.releaseDate) }}</span>
          </div>

          <div class="p-3 rounded-[12px] bg-white/[0.04] border border-white/[0.08]">
            <span class="text-[#7D8592] block">{{ t.details.adamId }}</span>
            <span class="font-mono text-[#FFFFFF] block mt-0.5">{{ app.id }}</span>
          </div>
        </div>

        <!-- Version History List with All Previous Builds Available -->
        <div class="space-y-3">
          <div class="flex items-center justify-between gap-2">
            <div class="flex items-center space-x-2">
              <h3 class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">{{ t.details.versionBuilds }}</h3>
              <span v-if="filteredVersions.length > 0" class="px-2 py-0.5 rounded-full text-[10px] font-mono font-bold bg-white/[0.08] text-[#B8C0CC] border border-white/[0.12]">
                {{ filteredVersions.length }}
              </span>
            </div>

            <!-- Version Filter Search Input -->
            <input
              v-if="allVersions.length > 5"
              v-model="versionFilter"
              type="text"
              :placeholder="t.details.versionFilterPlaceholder"
              class="glass-input px-3 py-1 text-xs w-48 font-mono"
            />
          </div>

          <!-- Loading Indicator when fetching version builds -->
          <div v-if="searchStore.isDetailsLoading" class="p-4 rounded-[14px] bg-white/[0.03] border border-white/[0.08] text-center text-xs text-[#0A84FF] font-medium flex items-center justify-center space-x-2 animate-pulse">
            <svg class="animate-spin h-4 w-4 text-[#0A84FF]" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
            </svg>
            <span>{{ t.details.loadingVersions }}</span>
          </div>

          <!-- Version History Scrollable List -->
          <div v-else-if="filteredVersions.length > 0" class="rounded-[14px] border border-white/[0.08] divide-y divide-white/[0.08] max-h-60 overflow-y-auto bg-white/[0.02]">
            <div
              v-for="(v, index) in filteredVersions"
              :key="v.externalVersionId"
              class="p-3.5 flex items-center justify-between text-xs hover:bg-white/[0.04] transition duration-150 gap-3"
            >
              <div class="flex items-center space-x-3 min-w-0">
                <div class="flex items-center space-x-2">
                  <span class="font-bold text-[#FFFFFF]">{{ v.displayVersion || `Build ${v.externalVersionId}` }}</span>
                  <span
                    v-if="index === 0 || v.displayVersion?.includes('Latest')"
                    class="px-2 py-0.5 rounded-full text-[9px] font-bold uppercase bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30"
                  >
                    {{ t.details.latestBadge }}
                  </span>
                </div>
                <span class="text-[11px] font-mono text-[#7D8592] truncate">{{ t.common.buildId }}: {{ v.externalVersionId }}</span>
              </div>

              <button
                v-if="isVersionDownloaded(v)"
                type="button"
                disabled
                class="px-3 py-1.5 rounded-xl bg-white/[0.08] text-[#30D158] border border-[#30D158]/30 text-xs font-medium flex items-center space-x-1.5 cursor-default opacity-80 select-none shrink-0"
              >
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
                </svg>
                <span>{{ t.downloadedApps?.downloaded || 'Descargado' }}</span>
              </button>
              <button
                v-else
                type="button"
                class="btn-secondary text-xs px-3.5 py-1.5 shrink-0 flex items-center space-x-1.5 hover:border-[#0A84FF]/60 hover:text-[#0A84FF]"
                @click="downloadVersion(v.externalVersionId)"
              >
                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                <span>{{ t.details.downloadBuild }}</span>
              </button>
            </div>
          </div>

          <!-- Informative Banner when not signed in -->
          <div v-if="!authStore.isLoggedIn && allVersions.length <= 1" class="p-3.5 rounded-[12px] bg-[#0A84FF]/10 border border-[#0A84FF]/25 text-xs text-[#64D2FF] flex items-center justify-between gap-3">
            <div class="flex items-center space-x-2.5">
              <svg class="w-4 h-4 text-[#64D2FF] shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>{{ t.details.signInForFullHistory }}</span>
            </div>
            <router-link to="/" class="btn-primary text-xs px-3 py-1 shrink-0" @click="closeModal">
              {{ t.auth.signInButton }}
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Fullscreen HD Screenshot Lightbox Viewer -->
    <div
      v-if="isLightboxOpen"
      class="fixed inset-0 z-60 flex items-center justify-center p-6 bg-black/85 backdrop-blur-[40px] animate-modal"
      @click.self="closeLightbox"
    >
      <!-- Close Button -->
      <button
        type="button"
        class="absolute top-6 right-6 p-2.5 rounded-full bg-white/10 hover:bg-white/20 text-white border border-white/20 shadow-xl transition duration-150 z-70"
        title="Close Preview (Esc)"
        @click="closeLightbox"
      >
        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>

      <!-- Previous Button -->
      <button
        v-if="screenshots.length > 1"
        type="button"
        class="absolute left-6 top-1/2 -translate-y-1/2 p-3 rounded-full bg-white/10 hover:bg-white/25 text-white border border-white/20 shadow-xl transition duration-150 z-70"
        title="Previous Screenshot (←)"
        @click.stop="prevScreenshot"
      >
        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M15 19l-7-7 7-7" />
        </svg>
      </button>

      <!-- Active Screenshot Display -->
      <div class="relative max-w-5xl max-h-[85vh] flex flex-col items-center justify-center">
        <img
          :src="screenshots[activeScreenshotIndex]"
          :alt="`Screenshot ${activeScreenshotIndex + 1}`"
          class="max-h-[80vh] max-w-full rounded-[20px] object-contain shadow-[0_24px_64px_rgba(0,0,0,0.6)] border border-white/[0.18]"
        />

        <!-- Image Counter Capsule -->
        <div class="mt-4 px-4 py-1.5 rounded-full bg-white/10 border border-white/20 backdrop-blur-md text-xs font-mono text-white shadow-md">
          {{ activeScreenshotIndex + 1 }} / {{ screenshots.length }}
        </div>
      </div>

      <!-- Next Button -->
      <button
        v-if="screenshots.length > 1"
        type="button"
        class="absolute right-6 top-1/2 -translate-y-1/2 p-3 rounded-full bg-white/10 hover:bg-white/25 text-white border border-white/20 shadow-xl transition duration-150 z-70"
        title="Next Screenshot (→)"
        @click.stop="nextScreenshot"
      >
        <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M9 5l7 7-7 7" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useSearchStore } from '../stores/search'
import { useDownloadsStore } from '../stores/downloads'
import { useAuthStore } from '../stores/auth'
import { useDownloadedAppsStore } from '../stores/downloadedApps'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'

const searchStore = useSearchStore()
const downloadsStore = useDownloadsStore()
const authStore = useAuthStore()
const downloadedAppsStore = useDownloadedAppsStore()
const { t } = useI18n()
const { showToast } = useNotifications()

const isLightboxOpen = ref(false)
const activeScreenshotIndex = ref(0)
const versionFilter = ref('')

const app = computed(() => searchStore.selectedApp?.metadata)

const screenshots = computed(() => {
  if (!app.value) return []
  return app.value.screenshots && app.value.screenshots.length > 0
    ? app.value.screenshots
    : app.value.ipadScreenshots || []
})

const allVersions = computed(() => {
  return searchStore.selectedApp?.versionHistory || []
})

const filteredVersions = computed(() => {
  if (!versionFilter.value.trim()) return allVersions.value
  const q = versionFilter.value.toLowerCase()
  return allVersions.value.filter(
    (v) => v.externalVersionId.toLowerCase().includes(q) || (v.displayVersion && v.displayVersion.toLowerCase().includes(q))
  )
})

function openLightbox(idx: number) {
  activeScreenshotIndex.value = idx
  isLightboxOpen.value = true
}

function closeLightbox() {
  isLightboxOpen.value = false
}

function prevScreenshot() {
  if (screenshots.value.length === 0) return
  activeScreenshotIndex.value =
    (activeScreenshotIndex.value - 1 + screenshots.value.length) % screenshots.value.length
}

function nextScreenshot() {
  if (screenshots.value.length === 0) return
  activeScreenshotIndex.value =
    (activeScreenshotIndex.value + 1) % screenshots.value.length
}

function handleKeydown(e: KeyboardEvent) {
  if (!isLightboxOpen.value) return
  if (e.key === 'ArrowLeft') {
    prevScreenshot()
  } else if (e.key === 'ArrowRight') {
    nextScreenshot()
  } else if (e.key === 'Escape') {
    closeLightbox()
  }
}

onMounted(() => {
  window.addEventListener('keydown', handleKeydown)
})

onUnmounted(() => {
  window.removeEventListener('keydown', handleKeydown)
})

function closeModal() {
  if (isLightboxOpen.value) {
    closeLightbox()
    return
  }
  searchStore.isDetailsModalOpen = false
}

function formatDate(d: string): string {
  if (!d) return '--'
  try {
    return new Date(d).toLocaleDateString()
  } catch {
    return d
  }
}

function isVersionDownloaded(v: any): boolean {
  if (!app.value) return false
  return downloadedAppsStore.downloadedIPAs.some(
    d => d.bundleId === app.value?.bundleId && (d.fileName.includes(v.externalVersionId) || d.version === v.displayVersion)
  )
}

async function downloadCurrent() {
  if (!app.value) return
  try {
    await downloadsStore.queueDownload(app.value, searchStore.platform)
    searchStore.isDetailsModalOpen = false
    showToast(t.value.search.downloadQueued, app.value.name, 'info')
  } catch (err: any) {
    showToast(t.value.search.downloadError, err?.message || 'Failed to queue download', 'error')
  }
}

async function downloadVersion(versionId: string) {
  if (!app.value) return
  try {
    await downloadsStore.queueDownload(app.value, searchStore.platform, versionId)
    searchStore.isDetailsModalOpen = false
    showToast(t.value.search.downloadQueued, `Build ${versionId} for ${app.value.name}`, 'info')
  } catch (err: any) {
    showToast(t.value.search.downloadError, err?.message || 'Failed to queue build', 'error')
  }
}
</script>
