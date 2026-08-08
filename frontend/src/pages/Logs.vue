<template>
  <div class="max-w-6xl mx-auto space-y-4 flex flex-col h-full animate-slide-up">
    <!-- Header with Filters & Controls -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4 shrink-0">
      <div>
        <h1 class="text-2xl font-bold">Real-Time Diagnostic Logs</h1>
        <p class="text-xs text-slate-400 mt-0.5">Stream backend Apple Storefront events, authentication cycles, and download progress in real time.</p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <!-- Severity Filter -->
        <select
          v-model="logsStore.filterLevel"
          class="glass-input px-3 py-1.5 rounded-xl text-xs font-semibold outline-none"
        >
          <option value="ALL">All Levels</option>
          <option value="INFO">INFO</option>
          <option value="SUCCESS">SUCCESS</option>
          <option value="WARN">WARN</option>
          <option value="ERROR">ERROR</option>
          <option value="DEBUG">DEBUG</option>
        </select>

        <!-- Auto-scroll Toggle -->
        <button
          type="button"
          class="btn-secondary text-xs px-3 py-1.5 flex items-center space-x-1"
          :class="logsStore.autoScroll ? 'text-blue-400 border-blue-500/30' : 'text-slate-400'"
          @click="logsStore.autoScroll = !logsStore.autoScroll"
        >
          <span>Auto-Scroll: {{ logsStore.autoScroll ? 'ON' : 'OFF' }}</span>
        </button>

        <button
          type="button"
          class="btn-secondary text-xs px-3 py-1.5 text-slate-300"
          @click="copyLogs"
        >
          Copy
        </button>

        <button
          type="button"
          class="btn-secondary text-xs px-3 py-1.5 text-slate-300"
          @click="exportLogs"
        >
          Export
        </button>

        <button
          type="button"
          class="btn-secondary text-xs px-3 py-1.5 text-rose-400 hover:text-rose-300"
          @click="clearLogs"
        >
          Clear
        </button>
      </div>
    </div>

    <!-- Live Terminal Log Box -->
    <div
      ref="logContainerRef"
      class="flex-1 rounded-2xl border border-white/10 bg-[#060913] p-4 font-mono text-xs overflow-y-auto space-y-1.5 shadow-2xl min-h-[420px]"
    >
      <div v-if="logsStore.filteredLogs.length === 0" class="text-slate-500 text-center py-24">
        No log entries recorded yet. Activities will stream here in real time.
      </div>

      <div
        v-for="entry in logsStore.filteredLogs"
        :key="entry.id"
        class="flex items-start space-x-2.5 leading-relaxed hover:bg-white/5 p-1 rounded transition"
      >
        <span class="text-slate-500 shrink-0 font-mono text-[11px]">{{ formatTime(entry.timestamp) }}</span>
        <span
          class="px-1.5 py-0.2 rounded text-[10px] font-bold shrink-0 uppercase"
          :class="levelBadgeClass(entry.level)"
        >
          {{ entry.level }}
        </span>
        <span v-if="entry.context" class="text-indigo-400 font-semibold shrink-0">[{{ entry.context }}]</span>
        <span class="text-slate-200 dark:text-slate-200 light:text-slate-300 break-all flex-1">{{ entry.message }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick, watch } from 'vue'
import { useLogsStore } from '../stores/logs'
import { useNotifications } from '../composables/useNotifications'

const logsStore = useLogsStore()
const { showToast } = useNotifications()
const logContainerRef = ref<HTMLElement | null>(null)

onMounted(() => {
  logsStore.fetchLogs()
  scrollToBottom()
})

watch(
  () => logsStore.logs.length,
  () => {
    if (logsStore.autoScroll) {
      nextTick(() => {
        scrollToBottom()
      })
    }
  }
)

function scrollToBottom() {
  if (logContainerRef.value) {
    logContainerRef.value.scrollTop = logContainerRef.value.scrollHeight
  }
}

function formatTime(timestamp: string) {
  if (!timestamp) return ''
  try {
    const d = new Date(timestamp)
    return d.toLocaleTimeString()
  } catch {
    return timestamp
  }
}

function levelBadgeClass(level: string) {
  switch (level?.toUpperCase()) {
    case 'SUCCESS':
      return 'bg-emerald-500/20 text-emerald-400 border border-emerald-500/30'
    case 'ERROR':
      return 'bg-rose-500/20 text-rose-400 border border-rose-500/30'
    case 'WARN':
      return 'bg-amber-500/20 text-amber-300 border border-amber-500/30'
    case 'DEBUG':
      return 'bg-purple-500/20 text-purple-300 border border-purple-500/30'
    default:
      return 'bg-blue-500/20 text-blue-400 border border-blue-500/30'
  }
}

async function copyLogs() {
  const text = logsStore.filteredLogs
    .map((l) => `[${l.timestamp}] [${l.level}] [${l.context || ''}] ${l.message}`)
    .join('\n')
  try {
    await navigator.clipboard.writeText(text)
    showToast('Copied', 'Logs copied to clipboard', 'success')
  } catch {
    showToast('Copy Failed', 'Clipboard access denied', 'error')
  }
}

async function exportLogs() {
  const path = await logsStore.exportLogs()
  if (path) {
    showToast('Logs Exported', `Saved to ${path}`, 'success')
  }
}

async function clearLogs() {
  await logsStore.clearLogs()
  showToast('Logs Cleared', 'Diagnostic log buffer cleared', 'info')
}
</script>
