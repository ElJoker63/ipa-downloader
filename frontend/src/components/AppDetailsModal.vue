<template>
  <div
    v-if="searchStore.isDetailsModalOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md transition-all duration-200"
    @click.self="searchStore.isDetailsModalOpen = false"
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
                {{ app.formattedPrice || 'Free' }}
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

        <!-- Header Actions (Close & Download) -->
        <div class="flex items-center space-x-2 shrink-0">
          <button
            type="button"
            class="btn-primary text-xs px-4 py-2 flex items-center space-x-1.5 shadow-sm"
            @click="downloadCurrent"
          >
            <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
            <span>Download</span>
          </button>
          <button
            type="button"
            class="p-2 rounded-xl text-[#7D8592] hover:text-white hover:bg-white/[0.08] transition duration-150"
            @click="searchStore.isDetailsModalOpen = false"
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
          <h3 class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">Screenshots</h3>
          <div class="flex space-x-3 overflow-x-auto pb-2 pt-1">
            <img
              v-for="(img, idx) in screenshots"
              :key="idx"
              :src="img"
              class="h-64 rounded-[14px] object-cover bg-[#171A21] border border-white/[0.12] shadow-md hover:scale-[1.02] transition-transform duration-200 cursor-pointer"
              loading="lazy"
            />
          </div>
        </div>

        <!-- Description -->
        <div v-if="app.description" class="space-y-2">
          <h3 class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">Description</h3>
          <div class="p-4 rounded-[14px] bg-white/[0.04] border border-white/[0.08] text-xs text-[#B8C0CC] leading-relaxed whitespace-pre-line max-h-48 overflow-y-auto font-normal">
            {{ app.description }}
          </div>
        </div>

        <!-- Metadata Properties Grid -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-3 text-xs">
          <div class="p-3 rounded-[12px] bg-white/[0.04] border border-white/[0.08]">
            <span class="text-[#7D8592] block">Bundle Identifier</span>
            <span class="font-mono text-[#FFFFFF] truncate block mt-0.5">{{ app.bundleId }}</span>
          </div>

          <div class="p-3 rounded-[12px] bg-white/[0.04] border border-white/[0.08]">
            <span class="text-[#7D8592] block">Minimum OS Version</span>
            <span class="font-mono text-[#FFFFFF] block mt-0.5">iOS {{ app.minimumOsVersion || '12.0' }}</span>
          </div>

          <div class="p-3 rounded-[12px] bg-white/[0.04] border border-white/[0.08]">
            <span class="text-[#7D8592] block">Release Date</span>
            <span class="text-[#FFFFFF] block mt-0.5">{{ formatDate(app.releaseDate) }}</span>
          </div>

          <div class="p-3 rounded-[12px] bg-white/[0.04] border border-white/[0.08]">
            <span class="text-[#7D8592] block">Apple Adam ID</span>
            <span class="font-mono text-[#FFFFFF] block mt-0.5">{{ app.id }}</span>
          </div>
        </div>

        <!-- Version History List with Direct Version Download -->
        <div v-if="searchStore.selectedApp?.versionHistory && searchStore.selectedApp.versionHistory.length > 0" class="space-y-2">
          <h3 class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">Available Version Builds</h3>
          <div class="rounded-[14px] border border-white/[0.08] divide-y divide-white/[0.08] max-h-48 overflow-y-auto bg-white/[0.02]">
            <div
              v-for="v in searchStore.selectedApp.versionHistory"
              :key="v.externalVersionId"
              class="p-3 flex items-center justify-between text-xs hover:bg-white/[0.04] transition duration-150"
            >
              <div class="flex items-center space-x-2.5">
                <span class="font-bold text-[#FFFFFF]">Build {{ v.displayVersion || v.externalVersionId }}</span>
                <span class="text-[11px] font-mono text-[#7D8592]">ID: {{ v.externalVersionId }}</span>
              </div>
              <button
                type="button"
                class="btn-secondary text-xs px-3 py-1"
                @click="downloadVersion(v.externalVersionId)"
              >
                Download Build
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useSearchStore } from '../stores/search'
import { useDownloadsStore } from '../stores/downloads'
import { useNotifications } from '../composables/useNotifications'

const searchStore = useSearchStore()
const downloadsStore = useDownloadsStore()
const { showToast } = useNotifications()

const app = computed(() => searchStore.selectedApp?.metadata)

const screenshots = computed(() => {
  if (!app.value) return []
  return app.value.screenshots && app.value.screenshots.length > 0
    ? app.value.screenshots
    : app.value.ipadScreenshots || []
})

function formatDate(d: string): string {
  if (!d) return '--'
  try {
    return new Date(d).toLocaleDateString()
  } catch {
    return d
  }
}

async function downloadCurrent() {
  if (!app.value) return
  try {
    await downloadsStore.queueDownload(app.value, searchStore.platform)
    searchStore.isDetailsModalOpen = false
    showToast('Download Queued', `Starting download for ${app.value.name}`, 'info')
  } catch (err: any) {
    showToast('Download Error', err?.message || 'Failed to queue download', 'error')
  }
}

async function downloadVersion(versionId: string) {
  if (!app.value) return
  try {
    await downloadsStore.queueDownload(app.value, searchStore.platform, versionId)
    searchStore.isDetailsModalOpen = false
    showToast('Download Queued', `Starting build ${versionId} for ${app.value.name}`, 'info')
  } catch (err: any) {
    showToast('Download Error', err?.message || 'Failed to queue build', 'error')
  }
}
</script>
