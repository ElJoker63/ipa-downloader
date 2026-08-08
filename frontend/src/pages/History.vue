<template>
  <div class="max-w-6xl mx-auto space-y-6 animate-slide-up">
    <!-- Header with Filters & Clear Actions -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl font-bold">Download History</h1>
        <p class="text-xs text-slate-400 mt-0.5">Complete record of downloaded iOS & tvOS packages with instant directory reveal and retry.</p>
      </div>

      <div class="flex items-center space-x-3 w-full md:w-auto justify-end">
        <!-- Status Filter -->
        <select
          v-model="historyStore.filterStatus"
          class="glass-input px-3 py-2 rounded-xl text-xs font-semibold outline-none"
        >
          <option value="all">All Statuses</option>
          <option value="completed">Completed</option>
          <option value="downloading">Downloading</option>
          <option value="failed">Failed</option>
          <option value="cancelled">Cancelled</option>
        </select>

        <button
          v-if="historyStore.history.length > 0"
          type="button"
          class="btn-secondary text-xs px-3.5 py-2 text-rose-400 hover:text-rose-300"
          @click="clearHistory"
        >
          Clear History
        </button>
      </div>
    </div>

    <!-- History Table / Cards -->
    <div v-if="historyStore.filteredHistory.length > 0" class="glass-panel rounded-2xl border border-white/10 overflow-hidden divide-y divide-white/5">
      <div
        v-for="item in historyStore.filteredHistory"
        :key="item.id"
        class="p-5 flex flex-col md:flex-row items-start md:items-center justify-between gap-4 hover:bg-white/5 transition"
      >
        <!-- App Details -->
        <div class="flex items-center space-x-4 min-w-0">
          <img
            :src="item.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
            :alt="item.appName"
            class="w-12 h-12 rounded-xl object-cover bg-slate-800 border border-white/10 shrink-0"
          />
          <div class="min-w-0">
            <div class="flex items-center space-x-2">
              <h3 class="text-sm font-bold truncate">{{ item.appName }}</h3>
              <span
                class="px-2 py-0.5 text-[10px] font-semibold rounded uppercase"
                :class="statusBadgeClass(item.status)"
              >
                {{ item.status }}
              </span>
            </div>
            <p class="text-xs text-slate-400 font-mono truncate mt-0.5" :title="item.destinationPath">{{ item.destinationPath }}</p>
            <div class="flex items-center space-x-3 text-[11px] text-slate-500 mt-1 font-mono">
              <span>v{{ item.version }}</span>
              <span>•</span>
              <span>{{ formatDate(item.createdAt) }}</span>
              <span v-if="item.error" class="text-rose-400 truncate">• {{ item.error }}</span>
            </div>
          </div>
        </div>

        <!-- Action Buttons: Retry, Open Folder, Delete, Copy Path -->
        <div class="flex flex-wrap items-center gap-2 shrink-0">
          <button
            v-if="item.status === 'failed' || item.status === 'cancelled'"
            type="button"
            class="btn-primary text-xs px-3 py-1.5"
            @click="retry(item.id)"
          >
            Retry
          </button>

          <button
            type="button"
            class="btn-secondary text-xs px-3 py-1.5"
            title="Reveal in File Explorer"
            @click="revealInExplorer(item.destinationPath)"
          >
            Open Folder
          </button>

          <button
            type="button"
            class="btn-secondary text-xs px-2.5 py-1.5"
            title="Copy File Path to Clipboard"
            @click="copyPath(item.destinationPath)"
          >
            📋
          </button>

          <button
            type="button"
            class="btn-secondary text-xs px-2.5 py-1.5 text-rose-400 hover:text-rose-300"
            title="Delete Record"
            @click="deleteItem(item.id)"
          >
            ✕
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="glass-panel p-12 rounded-2xl border border-white/10 text-center space-y-3">
      <div class="w-12 h-12 rounded-2xl bg-white/5 flex items-center justify-center mx-auto text-slate-400 text-2xl">
        📜
      </div>
      <h3 class="text-base font-bold">No Download History Found</h3>
      <p class="text-xs text-slate-400 max-w-sm mx-auto">
        Your downloaded packages will appear here with instant file management shortcuts.
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useHistoryStore } from '../stores/history'
import { useDownloadsStore } from '../stores/downloads'
import { useNotifications } from '../composables/useNotifications'

const router = useRouter()
const historyStore = useHistoryStore()
const downloadsStore = useDownloadsStore()
const { showToast } = useNotifications()

onMounted(() => {
  historyStore.fetchHistory()
})

async function retry(id: string) {
  await downloadsStore.retryDownload(id)
  showToast('Download Retried', 'Package re-queued', 'info')
  router.push('/downloads')
}

function revealInExplorer(path: string) {
  historyStore.revealInExplorer(path)
}

async function copyPath(path: string) {
  await historyStore.copyPath(path)
  showToast('Copied', 'Path copied to clipboard', 'info')
}

async function deleteItem(id: string) {
  await historyStore.deleteItem(id)
  showToast('Deleted', 'History item removed', 'info')
}

async function clearHistory() {
  await historyStore.clearHistory()
  showToast('Cleared', 'Download history cleared', 'info')
}

function formatDate(d: string) {
  if (!d) return ''
  try {
    return new Date(d).toLocaleString()
  } catch {
    return d
  }
}

function statusBadgeClass(status: string) {
  switch (status) {
    case 'completed':
      return 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
    case 'failed':
      return 'bg-rose-500/20 text-rose-400 border border-rose-500/30'
    case 'downloading':
      return 'bg-blue-500/20 text-blue-400 animate-pulse'
    default:
      return 'bg-slate-500/20 text-slate-400'
  }
}
</script>
