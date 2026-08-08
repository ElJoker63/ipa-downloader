<template>
  <div class="max-w-6xl mx-auto space-y-5 flex flex-col h-full animate-slide-up font-sans">
    <!-- Header with Filters & Controls -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4 shrink-0">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
          Real-Time Diagnostic Logs
        </h1>
        <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
          Stream Apple Storefront requests, authentication cycles, chunked transfer bytes, and FairPlay SINF DRM signing events.
        </p>
      </div>

      <div class="flex flex-wrap items-center gap-2">
        <!-- Severity Filter -->
        <select
          v-model="logsStore.filterLevel"
          class="glass-input px-3 py-1.5 rounded-[12px] text-xs font-semibold outline-none cursor-pointer"
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
          class="btn-secondary text-xs px-3 py-1.5 flex items-center space-x-1.5"
          :class="logsStore.autoScroll ? 'border-[#0A84FF]/40 text-[#0A84FF]' : 'text-[#7D8592]'"
          @click="logsStore.autoScroll = !logsStore.autoScroll"
        >
          <span>Auto-Scroll: {{ logsStore.autoScroll ? 'ON' : 'OFF' }}</span>
        </button>

        <button
          type="button"
          class="btn-secondary text-xs px-3 py-1.5 text-[#B8C0CC] hover:text-white"
          @click="copyLogs"
        >
          Copy
        </button>

        <button
          type="button"
          class="btn-secondary text-xs px-3 py-1.5 text-[#B8C0CC] hover:text-white"
          @click="exportLogs"
        >
          Export
        </button>

        <button
          type="button"
          class="btn-secondary text-xs px-3 py-1.5 text-[#FF453A] hover:bg-[#FF453A]/15 hover:border-[#FF453A]/30"
          @click="clearLogs"
        >
          Clear
        </button>
      </div>
    </div>

    <!-- Live Terminal Glass Box -->
    <div
      ref="logContainerRef"
      class="flex-1 rounded-[18px] border border-white/[0.12] bg-[#0A0D14]/90 backdrop-blur-[30px] p-4 font-mono text-xs overflow-y-auto space-y-1.5 shadow-[0_12px_40px_rgba(0,0,0,0.35)] min-h-[420px]"
    >
      <div v-if="logsStore.filteredLogs.length === 0" class="text-[#7D8592] text-center py-24 font-sans text-xs">
        No log entries recorded yet. Activities and streaming events will appear here in real time.
      </div>

      <div
        v-for="entry in logsStore.filteredLogs"
        :key="entry.id"
        class="flex items-start space-x-2.5 leading-relaxed hover:bg-white/[0.04] p-1 rounded-lg transition duration-100"
      >
        <span class="text-[#7D8592] shrink-0 font-mono text-[11px]">{{ formatTime(entry.timestamp) }}</span>
        <span
          class="px-1.5 py-0.2 rounded-[6px] text-[10px] font-bold shrink-0 uppercase"
          :class="levelBadgeClass(entry.level)"
        >
          {{ entry.level }}
        </span>
        <span v-if="entry.context" class="text-[#64D2FF] font-semibold shrink-0">[{{ entry.context }}]</span>
        <span class="text-[#FFFFFF] break-all flex-1">{{ entry.message }}</span>
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

function formatTime(isoString: string): string {
  if (!isoString) return ''
  const date = new Date(isoString)
  return date.toLocaleTimeString()
}

function levelBadgeClass(level: string): string {
  switch (level?.toUpperCase()) {
    case 'SUCCESS':
      return 'bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30'
    case 'WARN':
      return 'bg-[#FFD60A]/20 text-[#FFD60A] border border-[#FFD60A]/30'
    case 'ERROR':
      return 'bg-[#FF453A]/20 text-[#FF453A] border border-[#FF453A]/30'
    case 'DEBUG':
      return 'bg-white/[0.08] text-[#7D8592]'
    default:
      return 'bg-[#0A84FF]/20 text-[#64D2FF] border border-[#0A84FF]/30'
  }
}

async function copyLogs() {
  try {
    const text = logsStore.logs
      .map((l) => `[${formatTime(l.timestamp)}] [${l.level}] [${l.context}] ${l.message}`)
      .join('\n')
    await navigator.clipboard.writeText(text)
    showToast('Logs Copied', 'All log entries copied to clipboard', 'info')
  } catch {
    showToast('Copy Failed', 'Could not copy logs to clipboard', 'error')
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
  showToast('Logs Cleared', 'All log entries removed', 'info')
}
</script>
