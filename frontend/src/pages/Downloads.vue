<template>
  <div class="max-w-6xl mx-auto space-y-6 animate-slide-up">
    <!-- Top Action & Destination Bar -->
    <div class="glass-panel p-5 rounded-2xl border border-white/10 flex flex-col md:flex-row items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold">Download Manager</h1>
        <p class="text-xs text-slate-400 mt-0.5">
          Real-time chunked streaming with FairPlay signature replication and concurrent download workers.
        </p>
      </div>

      <div class="flex items-center space-x-3 w-full md:w-auto justify-end">
        <button
          type="button"
          class="btn-secondary text-xs px-3.5 py-2 flex items-center space-x-1.5"
          @click="browseFolder"
        >
          <span>📁 Destination: {{ settingsStore.settings.defaultDownloadFolder || 'Default Downloads' }}</span>
        </button>

        <button
          v-if="downloadsStore.completedDownloads.length > 0"
          type="button"
          class="btn-secondary text-xs px-3 py-2 text-slate-400 hover:text-white"
          @click="downloadsStore.clearCompleted()"
        >
          Clear Finished
        </button>
      </div>
    </div>

    <!-- Active Transfers Section -->
    <div class="space-y-3">
      <div class="flex items-center justify-between px-1">
        <div class="flex items-center space-x-2">
          <h2 class="text-sm font-bold uppercase tracking-wider text-slate-300">Active Queue</h2>
          <span v-if="downloadsStore.activeCount > 0" class="px-2 py-0.5 text-xs font-bold rounded-full bg-blue-500/20 text-blue-400">
            {{ downloadsStore.activeCount }}
          </span>
        </div>
        <span v-if="downloadsStore.activeCount > 0" class="text-xs font-mono text-emerald-400">
          Total Speed: {{ downloadsStore.totalSpeedFormatted }}
        </span>
      </div>

      <!-- Active Transfers List -->
      <div v-if="downloadsStore.activeDownloads.length > 0" class="space-y-3">
        <div
          v-for="task in downloadsStore.activeDownloads"
          :key="task.id"
          class="glass-panel p-5 rounded-2xl border border-white/10 shadow-lg space-y-3"
        >
          <!-- Task Header -->
          <div class="flex items-center justify-between">
            <div class="flex items-center space-x-3 min-w-0">
              <img
                :src="task.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
                :alt="task.appName"
                class="w-12 h-12 rounded-xl object-cover bg-slate-800 border border-white/10 shrink-0"
              />
              <div class="min-w-0">
                <h3 class="text-base font-bold truncate">{{ task.appName }}</h3>
                <div class="flex items-center space-x-2 text-xs text-slate-400 font-mono">
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
                class="p-2 rounded-lg bg-white/5 hover:bg-white/10 text-slate-300 hover:text-white transition"
                title="Pause Download"
                @click="downloadsStore.pauseDownload(task.id)"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </button>

              <button
                v-if="task.status === 'paused'"
                type="button"
                class="p-2 rounded-lg bg-emerald-500/20 hover:bg-emerald-500/30 text-emerald-400 transition"
                title="Resume Download"
                @click="downloadsStore.resumeDownload(task.id)"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </button>

              <button
                type="button"
                class="p-2 rounded-lg bg-rose-500/10 hover:bg-rose-500/20 text-rose-400 transition"
                title="Cancel Download"
                @click="downloadsStore.cancelDownload(task.id)"
              >
                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>

          <!-- Signing Phase Notification Banner -->
          <div v-if="task.status === 'signing'" class="flex items-center space-x-2.5 p-2.5 rounded-xl bg-indigo-500/15 border border-indigo-500/30 text-indigo-300 text-xs font-semibold animate-pulse">
            <svg class="w-4 h-4 animate-spin shrink-0 text-indigo-400" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
            </svg>
            <span>Download complete! Injecting Apple FairPlay DRM SINF signature certificates into final .ipa package...</span>
          </div>

          <!-- Progress Bar -->
          <div class="w-full bg-slate-900/60 rounded-full h-2.5 overflow-hidden border border-white/5 relative">
            <div
              class="h-full rounded-full transition-all duration-300 relative overflow-hidden"
              :class="task.status === 'signing' ? 'bg-gradient-to-r from-purple-500 via-indigo-500 to-cyan-400 animate-pulse' : (task.status === 'paused' ? 'bg-amber-500' : 'bg-gradient-to-r from-blue-600 to-cyan-400')"
              :style="{ width: `${Math.max(task.progress, 2)}%` }"
            >
              <div class="absolute inset-0 bg-white/20 animate-pulse-subtle"></div>
            </div>
          </div>

          <!-- Progress Stats Footer -->
          <div class="flex items-center justify-between text-xs text-slate-400 font-mono">
            <div class="flex items-center space-x-3">
              <span class="font-bold text-slate-200">{{ task.progress.toFixed(1) }}%</span>
              <span>{{ formatBytes(task.downloadedBytes) }} / {{ formatBytes(task.totalBytes) }}</span>
            </div>
            <div class="flex items-center space-x-3">
              <span v-if="task.status === 'signing'" class="text-indigo-400 font-bold flex items-center space-x-1.5 animate-pulse">
                <span class="inline-block w-2 h-2 rounded-full bg-indigo-400 animate-ping"></span>
                <span>Signing FairPlay DRM...</span>
              </span>
              <span v-else-if="task.status === 'downloading'" class="text-emerald-400 font-semibold">{{ task.formattedSpeed || 'Streaming...' }}</span>
              <span v-else-if="task.status === 'paused'" class="text-amber-400 font-semibold">Paused</span>
              <span v-if="task.formattedETA && task.status === 'downloading'">ETA: {{ task.formattedETA }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- No Active Downloads Empty Card -->
      <div v-else class="glass-panel p-8 rounded-2xl border border-white/5 text-center space-y-2">
        <div class="text-2xl">⬇</div>
        <p class="text-sm font-semibold">No active downloads</p>
        <p class="text-xs text-slate-400">Search for any iOS or tvOS application to queue downloads.</p>
      </div>
    </div>

    <!-- Completed & Past Downloads Section -->
    <div v-if="downloadsStore.completedDownloads.length > 0" class="space-y-3 pt-4">
      <div class="flex items-center justify-between px-1">
        <h2 class="text-sm font-bold uppercase tracking-wider text-slate-300">Completed & Past Transfers</h2>
      </div>

      <div class="glass-panel rounded-2xl border border-white/10 divide-y divide-white/5 overflow-hidden">
        <div
          v-for="task in downloadsStore.completedDownloads"
          :key="task.id"
          class="p-4 flex items-center justify-between hover:bg-white/5 transition"
        >
          <div class="flex items-center space-x-3 min-w-0">
            <img
              :src="task.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
              class="w-10 h-10 rounded-xl object-cover bg-slate-800 border border-white/10 shrink-0"
            />
            <div class="min-w-0">
              <div class="flex items-center space-x-2">
                <h4 class="text-sm font-bold truncate">{{ task.appName }}</h4>
                <span
                  class="px-2 py-0.2 text-[10px] font-semibold rounded uppercase"
                  :class="statusBadgeClass(task.status)"
                >
                  {{ task.status }}
                </span>
              </div>
              <p class="text-xs text-slate-400 truncate font-mono">{{ task.destinationPath }}</p>
              <div v-if="task.error" class="mt-1 p-2 rounded-lg bg-rose-500/15 border border-rose-500/30 text-xs text-rose-300 font-mono break-all flex items-start justify-between gap-2">
                <span>⚠️ {{ task.error }}</span>
                <button
                  type="button"
                  class="text-[10px] uppercase font-bold text-rose-400 hover:text-white underline shrink-0"
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
              class="btn-primary text-xs px-3 py-1.5"
              @click="downloadsStore.retryDownload(task.id)"
            >
              Retry
            </button>
            <button
              type="button"
              class="btn-secondary text-xs px-3 py-1.5"
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
      return 'bg-indigo-500/20 text-indigo-400 border border-indigo-500/30 animate-pulse'
    case 'completed':
      return 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
    case 'failed':
      return 'bg-rose-500/20 text-rose-400 border border-rose-500/30'
    case 'cancelled':
      return 'bg-slate-500/20 text-slate-400'
    default:
      return 'bg-blue-500/20 text-blue-400'
  }
}
</script>
