<template>
  <div class="max-w-6xl mx-auto space-y-6 flex flex-col h-full animate-slide-up font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4 shrink-0">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
          Download Manager
        </h1>
        <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
          Concurrent streaming transfers with real-time speed, pause/resume, and FairPlay SINF DRM signing.
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center space-x-2.5">
        <button
          type="button"
          class="btn-secondary text-xs px-3.5 py-2 flex items-center space-x-1.5"
          @click="browseFolder"
        >
          <svg class="w-3.5 h-3.5 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          <span>Destination Folder</span>
        </button>

        <button
          v-if="downloadsStore.completedDownloads.length > 0"
          type="button"
          class="btn-secondary text-xs px-3.5 py-2 text-[#7D8592] hover:text-[#FFFFFF]"
          @click="downloadsStore.clearCompleted()"
        >
          Clear Completed
        </button>
      </div>
    </div>

    <!-- Active Transfers Section -->
    <div class="space-y-3 shrink-0">
      <div class="flex items-center justify-between px-1">
        <div class="flex items-center space-x-2">
          <span class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">Active Transfers</span>
          <span
            v-if="downloadsStore.activeCount > 0"
            class="px-2 py-0.5 text-[10px] font-medium rounded-full bg-[#0A84FF]/20 text-[#0A84FF] border border-[#0A84FF]/30"
          >
            {{ downloadsStore.activeCount }}
          </span>
        </div>
        <span v-if="downloadsStore.activeCount > 0" class="text-xs font-mono text-[#30D158] font-medium">
          Total Speed: {{ downloadsStore.totalSpeedFormatted }}
        </span>
      </div>

      <!-- Active Transfers List -->
      <div v-if="downloadsStore.activeDownloads.length > 0" class="space-y-3.5">
        <div
          v-for="task in downloadsStore.activeDownloads"
          :key="task.id"
          class="glass-card p-5 rounded-[18px] space-y-3.5"
        >
          <!-- Task Header -->
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3.5 min-w-0">
              <img
                :src="task.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
                :alt="task.appName"
                class="w-12 h-12 rounded-[14px] object-cover bg-[#171A21] border border-white/[0.18] shadow-md shrink-0"
              />
              <div class="min-w-0">
                <h3 class="text-sm font-semibold truncate text-[#FFFFFF]">{{ task.appName }}</h3>
                <div class="flex items-center space-x-2 text-xs text-[#B8C0CC] font-mono mt-0.5">
                  <span>{{ task.bundleId }}</span>
                  <span>•</span>
                  <span>v{{ task.version }}</span>
                </div>
              </div>
            </div>

            <!-- Action Controls (Pause / Resume / Cancel) -->
            <div class="flex items-center space-x-2">
              <button
                v-if="task.status === 'downloading'"
                type="button"
                class="p-2 rounded-xl bg-white/[0.08] hover:bg-white/[0.14] text-[#B8C0CC] hover:text-white transition duration-150"
                title="Pause Download"
                @click="downloadsStore.pauseDownload(task.id)"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </button>

              <button
                v-if="task.status === 'paused'"
                type="button"
                class="p-2 rounded-xl bg-[#30D158]/15 hover:bg-[#30D158]/25 text-[#30D158] border border-[#30D158]/30 transition duration-150"
                title="Resume Download"
                @click="downloadsStore.resumeDownload(task.id)"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </button>

              <button
                type="button"
                class="p-2 rounded-xl bg-[#FF453A]/15 hover:bg-[#FF453A]/25 text-[#FF453A] border border-[#FF453A]/30 transition duration-150"
                title="Cancel Download"
                @click="downloadsStore.cancelDownload(task.id)"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Signing Phase Notification Banner -->
          <div v-if="task.status === 'signing'" class="flex items-center space-x-2.5 p-3 rounded-[12px] bg-[#0A84FF]/15 border border-[#0A84FF]/30 text-[#64D2FF] text-xs font-medium animate-pulse">
            <svg class="w-4 h-4 animate-spin shrink-0 text-[#0A84FF]" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
            </svg>
            <span>Download complete! Injecting Apple FairPlay DRM SINF signature certificates into final .ipa package...</span>
          </div>

          <!-- Progress Bar -->
          <div class="w-full bg-black/40 rounded-full h-2.5 overflow-hidden border border-white/[0.08] relative">
            <div
              class="h-full rounded-full transition-all duration-300 relative overflow-hidden"
              :class="task.status === 'signing' ? 'bg-gradient-to-r from-[#0A84FF] via-[#64D2FF] to-[#30D158] animate-pulse' : (task.status === 'paused' ? 'bg-[#FFD60A]' : 'bg-[#0A84FF]')"
              :style="{ width: `${Math.max(task.progress, 2)}%` }"
            >
              <div class="absolute inset-0 bg-white/20 animate-pulse-subtle"></div>
            </div>
          </div>

          <!-- Progress Stats Footer -->
          <div class="flex items-center justify-between text-xs text-[#B8C0CC] font-mono">
            <div class="flex items-center space-x-3">
              <span class="font-bold text-[#FFFFFF]">{{ task.progress.toFixed(1) }}%</span>
              <span>{{ formatBytes(task.downloadedBytes) }} / {{ formatBytes(task.totalBytes) }}</span>
            </div>
            <div class="flex items-center space-x-3">
              <span v-if="task.status === 'signing'" class="text-[#64D2FF] font-semibold flex items-center space-x-1.5 animate-pulse">
                <span class="inline-block w-2 h-2 rounded-full bg-[#64D2FF] animate-ping"></span>
                <span>Signing FairPlay DRM...</span>
              </span>
              <span v-else-if="task.status === 'downloading'" class="text-[#30D158] font-semibold">{{ task.formattedSpeed || 'Streaming...' }}</span>
              <span v-else-if="task.status === 'paused'" class="text-[#FFD60A] font-semibold">Paused</span>
              <span v-if="task.formattedETA && task.status === 'downloading'" class="text-[#B8C0CC]">ETA: {{ task.formattedETA }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- No Active Downloads Empty Card -->
      <div v-else class="glass-card p-10 rounded-[18px] text-center space-y-2">
        <svg class="w-8 h-8 text-[#7D8592] mx-auto opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
        </svg>
        <p class="text-sm font-semibold text-[#FFFFFF]">No active downloads</p>
        <p class="text-xs text-[#B8C0CC]">Search for any iOS or tvOS application to queue downloads.</p>
      </div>
    </div>

    <!-- Completed & Past Downloads Section -->
    <div v-if="downloadsStore.completedDownloads.length > 0" class="space-y-3 pt-2">
      <div class="flex items-center justify-between px-1">
        <h2 class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">Completed & Past Transfers</h2>
      </div>

      <div class="glass-card rounded-[18px] divide-y divide-white/[0.08] overflow-hidden">
        <div
          v-for="task in downloadsStore.completedDownloads"
          :key="task.id"
          class="p-4 flex items-center justify-between hover:bg-white/[0.04] transition duration-150"
        >
          <div class="flex items-center space-x-3.5 min-w-0">
            <img
              :src="task.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
              class="w-10 h-10 rounded-[12px] object-cover bg-[#171A21] border border-white/[0.18] shrink-0 shadow-sm"
            />
            <div class="min-w-0">
              <div class="flex items-center space-x-2">
                <h4 class="text-sm font-semibold truncate text-[#FFFFFF]">{{ task.appName }}</h4>
                <span
                  class="px-2 py-0.5 text-[10px] font-medium rounded-full uppercase"
                  :class="statusBadgeClass(task.status)"
                >
                  {{ task.status }}
                </span>
              </div>
              <p class="text-xs text-[#B8C0CC] truncate font-mono mt-0.5">{{ task.destinationPath }}</p>
              <div v-if="task.error" class="mt-1.5 p-2 rounded-lg bg-[#FF453A]/15 border border-[#FF453A]/30 text-xs text-[#FF453A] font-mono break-all flex items-start justify-between gap-2">
                <span>{{ task.error }}</span>
                <button
                  type="button"
                  class="text-[10px] uppercase font-bold text-[#FF453A] hover:text-white underline shrink-0"
                  @click="copyErrorText(task.error)"
                >
                  Copy Error
                </button>
              </div>
            </div>
          </div>

          <div class="flex items-center space-x-2 shrink-0">
            <button
              v-if="task.status === 'failed'"
              type="button"
              class="btn-primary text-xs px-3.5 py-1.5"
              @click="downloadsStore.retryDownload(task.id)"
            >
              Retry
            </button>
            <button
              type="button"
              class="btn-secondary text-xs px-3.5 py-1.5"
              @click="revealInExplorer(task.destinationPath)"
            >
              Show in Folder
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useDownloadsStore } from '../stores/downloads'
import { useSettingsStore } from '../stores/settings'
import { useHistoryStore } from '../stores/history'
import { useNotifications } from '../composables/useNotifications'

const downloadsStore = useDownloadsStore()
const settingsStore = useSettingsStore()
const historyStore = useHistoryStore()
const { showToast } = useNotifications()

onMounted(() => {
  downloadsStore.fetchDownloads()
  settingsStore.fetchSettings()
})

async function browseFolder() {
  await settingsStore.browseFolder()
  showToast('Download Folder Updated', settingsStore.settings.defaultDownloadFolder, 'info')
}

function revealInExplorer(path: string) {
  historyStore.revealInExplorer(path)
}

async function copyErrorText(errText: string) {
  try {
    await navigator.clipboard.writeText(errText)
    showToast('Error Copied', 'Error details copied to clipboard', 'info')
  } catch {
    showToast('Copy Failed', errText, 'error')
  }
}

function formatBytes(bytes: number): string {
  if (!bytes || bytes <= 0) return '0 B'
  const unit = 1024
  if (bytes < unit) return `${bytes} B`
  const div = 1024
  const exp = Math.floor(Math.log(bytes) / Math.log(div))
  const letter = 'KMGTPE'[exp - 1] || 'K'
  return `${(bytes / Math.pow(div, exp)).toFixed(1)} ${letter}B`
}

function statusBadgeClass(status: string) {
  switch (status) {
    case 'signing':
      return 'bg-[#0A84FF]/20 text-[#64D2FF] border border-[#0A84FF]/30 animate-pulse'
    case 'completed':
      return 'bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30'
    case 'failed':
      return 'bg-[#FF453A]/20 text-[#FF453A] border border-[#FF453A]/30'
    case 'cancelled':
      return 'bg-white/[0.08] text-[#7D8592]'
    default:
      return 'bg-[#0A84FF]/20 text-[#0A84FF]'
  }
}
</script>
