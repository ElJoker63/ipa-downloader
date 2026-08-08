<template>
  <div v-if="searchStore.isDetailsModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 md:p-6 bg-black/70 backdrop-blur-md animate-fade-in" @click.self="searchStore.closeAppDetails()">
    <div class="glass-panel w-full max-w-4xl max-h-[90vh] flex flex-col rounded-2xl shadow-2xl border border-white/10 dark:border-white/10 light:border-black/10 bg-[#131B2E]/95 dark:bg-[#131B2E]/95 light:bg-white/95 text-slate-100 dark:text-slate-100 light:text-slate-900 overflow-hidden animate-slide-up">
      <!-- Header Bar -->
      <div class="flex items-center justify-between px-6 py-4 border-b border-white/10 dark:border-white/10 light:border-black/10 bg-slate-900/40 dark:bg-slate-900/40 light:bg-slate-50">
        <div class="flex items-center space-x-3">
          <span class="px-2.5 py-1 text-xs font-semibold uppercase tracking-wider rounded-md bg-blue-500/20 text-blue-400 border border-blue-500/30">
            {{ app.primaryGenre || 'Application' }}
          </span>
          <span class="text-xs font-mono text-slate-400 dark:text-slate-400 light:text-slate-500">{{ app.bundleId }}</span>
        </div>
        <button
          type="button"
          class="p-1.5 rounded-lg text-slate-400 hover:text-white hover:bg-white/10 transition"
          @click="searchStore.closeAppDetails()"
        >
          <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>

      <!-- Scrollable Body -->
      <div class="overflow-y-auto p-6 space-y-6 flex-1">
        <!-- Hero Overview -->
        <div class="flex flex-col md:flex-row gap-6 items-start">
          <img
            :src="app.artworkUrl512 || app.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
            :alt="app.name"
            class="w-24 h-24 md:w-32 md:h-32 rounded-2xl shadow-xl object-cover bg-slate-800 border border-white/10 shrink-0"
          />
          <div class="flex-1 space-y-2">
            <div class="flex items-start justify-between">
              <div>
                <h2 class="text-2xl font-bold text-slate-100 dark:text-slate-100 light:text-slate-900">{{ app.name }}</h2>
                <p class="text-sm text-slate-400 dark:text-slate-400 light:text-slate-600 font-medium">{{ app.developer }}</p>
              </div>
              <div class="text-right">
                <span class="px-3 py-1 text-sm font-semibold rounded-full bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">
                  {{ app.formattedPrice || 'Free' }}
                </span>
              </div>
            </div>

            <!-- Quick Meta Tags -->
            <div class="flex flex-wrap gap-2 pt-2">
              <span class="px-2.5 py-0.5 text-xs rounded-md bg-white/5 dark:bg-white/5 light:bg-slate-200 border border-white/5 text-slate-300 dark:text-slate-300 light:text-slate-700">
                v{{ app.version }}
              </span>
              <span v-if="app.formattedSize" class="px-2.5 py-0.5 text-xs rounded-md bg-white/5 dark:bg-white/5 light:bg-slate-200 border border-white/5 text-slate-300 dark:text-slate-300 light:text-slate-700">
                💾 {{ app.formattedSize }}
              </span>
              <span v-if="app.minimumOsVersion" class="px-2.5 py-0.5 text-xs rounded-md bg-white/5 dark:bg-white/5 light:bg-slate-200 border border-white/5 text-slate-300 dark:text-slate-300 light:text-slate-700">
                📱 iOS {{ app.minimumOsVersion }}+
              </span>
              <span v-if="app.contentAdvisoryRating" class="px-2.5 py-0.5 text-xs rounded-md bg-white/5 dark:bg-white/5 light:bg-slate-200 border border-white/5 text-slate-300 dark:text-slate-300 light:text-slate-700">
                {{ app.contentAdvisoryRating }}
              </span>
              <span v-if="app.averageUserRating" class="px-2.5 py-0.5 text-xs rounded-md bg-amber-500/15 text-amber-300 border border-amber-500/20 flex items-center gap-1">
                ⭐ {{ app.averageUserRating.toFixed(1) }} ({{ app.userRatingCount || 0 }})
              </span>
            </div>

            <!-- Primary Action Buttons -->
            <div class="flex flex-wrap items-center gap-3 pt-4">
              <button
                type="button"
                class="btn-primary text-sm px-5 py-2.5 flex items-center space-x-2"
                :disabled="isDownloading"
                @click="downloadLatest"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
                <span>Download IPA</span>
              </button>

              <button
                type="button"
                class="btn-secondary text-sm px-4 py-2.5 flex items-center space-x-2"
                @click="toggleFav"
              >
                <svg class="w-4 h-4" :class="isFav ? 'text-rose-500 fill-rose-500' : 'text-slate-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
                <span>{{ isFav ? 'In Favorites' : 'Add to Favorites' }}</span>
              </button>
            </div>
          </div>
        </div>

        <!-- Screenshots Gallery Carousel -->
        <div v-if="screenshots.length > 0" class="space-y-2">
          <h3 class="text-sm font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-400 light:text-slate-600">Screenshots</h3>
          <div class="flex gap-4 overflow-x-auto pb-4 pt-1 snap-x scrollbar-thin">
            <img
              v-for="(src, idx) in screenshots"
              :key="idx"
              :src="src"
              :alt="`Screenshot ${idx + 1}`"
              class="h-64 md:h-80 rounded-xl shadow-md border border-white/10 object-contain bg-slate-900 shrink-0 snap-start hover:scale-[1.02] transition cursor-pointer"
              @click="activeLightbox = src"
            />
          </div>
        </div>

        <!-- Description -->
        <div v-if="app.description" class="space-y-2">
          <h3 class="text-sm font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-400 light:text-slate-600">Description</h3>
          <div
            class="text-sm leading-relaxed text-slate-300 dark:text-slate-300 light:text-slate-700 whitespace-pre-line rounded-xl p-4 bg-slate-900/40 dark:bg-slate-900/40 light:bg-slate-100 border border-white/5 dark:border-white/5 light:border-black/5"
            :class="isDescExpanded ? '' : 'max-h-36 overflow-hidden relative'"
          >
            {{ app.description }}
            <div v-if="!isDescExpanded" class="absolute inset-x-0 bottom-0 h-16 bg-gradient-to-t from-slate-900/90 dark:from-slate-900/90 light:from-slate-100 to-transparent flex items-end justify-center pb-2">
              <button type="button" class="text-xs font-semibold text-blue-400 hover:text-blue-300" @click="isDescExpanded = true">Show More</button>
            </div>
          </div>
          <button v-if="isDescExpanded" type="button" class="text-xs font-semibold text-blue-400 hover:text-blue-300 mt-1" @click="isDescExpanded = false">Show Less</button>
        </div>

        <!-- Version History -->
        <div v-if="versionHistory.length > 0" class="space-y-2">
          <h3 class="text-sm font-semibold uppercase tracking-wider text-slate-400 dark:text-slate-400 light:text-slate-600">Version History</h3>
          <div class="rounded-xl border border-white/10 dark:border-white/10 light:border-black/10 overflow-hidden divide-y divide-white/5 dark:divide-white/5 light:divide-black/5 bg-slate-900/40 dark:bg-slate-900/40 light:bg-slate-100 max-h-56 overflow-y-auto">
            <div
              v-for="v in versionHistory"
              :key="v.externalVersionId"
              class="flex items-center justify-between p-3 text-sm hover:bg-white/5 transition"
            >
              <div class="flex items-center space-x-3">
                <span class="font-mono text-xs text-slate-400">{{ v.externalVersionId }}</span>
                <span class="font-medium text-slate-200 dark:text-slate-200 light:text-slate-800">{{ v.displayVersion }}</span>
              </div>
              <div class="flex items-center space-x-3">
                <span v-if="v.formattedDate" class="text-xs text-slate-400">{{ v.formattedDate }}</span>
                <button
                  type="button"
                  class="btn-secondary text-xs px-2.5 py-1"
                  @click="downloadVersion(v.externalVersionId)"
                >
                  Download
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Lightbox Modal -->
    <div v-if="activeLightbox" class="fixed inset-0 z-60 flex items-center justify-center p-4 bg-black/90 backdrop-blur-xl" @click="activeLightbox = null">
      <img :src="activeLightbox" class="max-w-full max-h-[95vh] rounded-2xl shadow-2xl object-contain" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useSearchStore } from '../stores/search'
import { useDownloadsStore } from '../stores/downloads'
import { useFavoritesStore } from '../stores/favorites'
import { useNotifications } from '../composables/useNotifications'

const router = useRouter()
const searchStore = useSearchStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const { showToast } = useNotifications()

const isDescExpanded = ref(false)
const isDownloading = ref(false)
const activeLightbox = ref<string | null>(null)

const app = computed(() => searchStore.selectedApp?.metadata || ({} as any))
const versionHistory = computed(() => searchStore.selectedApp?.versionHistory || [])
const screenshots = computed(() => {
  const meta = app.value
  const list: string[] = []
  if (meta.screenshots) list.push(...meta.screenshots)
  if (meta.ipadScreenshots) list.push(...meta.ipadScreenshots)
  return list
})

const isFav = computed(() => {
  if (searchStore.selectedApp?.isFavorite) return true
  return favoritesStore.isFavorite(app.value.id)
})

async function downloadLatest() {
  isDownloading.value = true
  try {
    await downloadsStore.queueDownload(app.value, searchStore.platform)
    showToast('Download Started', `Queued ${app.value.name} for download`, 'success')
    searchStore.closeAppDetails()
    router.push('/downloads')
  } catch (err: any) {
    showToast('Download Failed', err?.message || 'Could not queue download', 'error')
  } finally {
    isDownloading.value = false
  }
}

async function downloadVersion(versionId: string) {
  try {
    await downloadsStore.queueDownload(app.value, searchStore.platform, versionId)
    showToast('Download Started', `Downloading version ${versionId}`, 'success')
    searchStore.closeAppDetails()
    router.push('/downloads')
  } catch (err: any) {
    showToast('Download Failed', err?.message || 'Could not queue download', 'error')
  }
}

async function toggleFav() {
  const result = await favoritesStore.toggleFavorite(app.value)
  if (searchStore.selectedApp) {
    searchStore.selectedApp.isFavorite = result
  }
  showToast(result ? 'Added to Favorites' : 'Removed from Favorites', app.value.name, 'info')
}
</script>
