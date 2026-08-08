<template>
  <div class="max-w-6xl mx-auto space-y-6 flex flex-col h-full animate-slide-up font-sans">
    <!-- Header Section -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4 shrink-0">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
          {{ t.history.title }}
        </h1>
        <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
          {{ t.history.subtitle }}
        </p>
      </div>

      <!-- Action Controls -->
      <div class="flex items-center space-x-2.5">
        <button
          v-if="historyStore.history.length > 0"
          type="button"
          class="btn-secondary text-xs px-3.5 py-2 text-[#FF453A] hover:bg-[#FF453A]/15 hover:border-[#FF453A]/30"
          @click="clearHistory"
        >
          {{ t.history.clearHistory }}
        </button>
      </div>
    </div>

    <!-- History Table & List (Clean macOS Table) -->
    <div class="flex-1 min-h-0 overflow-y-auto">
      <div v-if="historyStore.history.length > 0" class="glass-card rounded-[18px] divide-y divide-white/[0.08] overflow-hidden">
        <div
          v-for="item in historyStore.history"
          :key="item.id"
          class="p-4 flex items-center justify-between hover:bg-white/[0.04] transition duration-150 gap-4"
        >
          <!-- App Info -->
          <div class="flex items-center space-x-3.5 min-w-0 flex-1">
            <img
              :src="item.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
              :alt="item.appName"
              class="w-11 h-11 rounded-[12px] object-cover bg-[#171A21] border border-white/[0.18] shadow-sm shrink-0"
            />
            <div class="min-w-0 flex-1">
              <div class="flex items-center space-x-2.5">
                <h3 class="text-sm font-semibold truncate text-[#FFFFFF]">{{ item.appName }}</h3>
                <span
                  class="px-2 py-0.5 text-[10px] font-semibold rounded-full uppercase"
                  :class="statusBadgeClass(item.status)"
                >
                  {{ item.status }}
                </span>
              </div>
              <p class="text-xs text-[#B8C0CC] truncate font-mono mt-0.5">{{ item.destinationPath }}</p>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex items-center space-x-2 shrink-0">
            <button
              v-if="deviceStore.isConnected && item.status === 'completed'"
              type="button"
              class="btn-primary text-xs px-3 py-1.5 flex items-center space-x-1.5"
              @click="installToDevice(item.destinationPath)"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
              <span>Install</span>
            </button>
            <button
              type="button"
              class="btn-secondary text-xs px-3 py-1.5 flex items-center space-x-1"
              :title="t.history.copyPath"
              @click="copyPath(item.destinationPath)"
            >
              <svg class="w-3.5 h-3.5 text-[#B8C0CC]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z" />
              </svg>
              <span>{{ t.history.copyPath }}</span>
            </button>

            <button
              type="button"
              class="btn-secondary text-xs px-3 py-1.5 flex items-center space-x-1"
              @click="revealInExplorer(item.destinationPath)"
            >
              <svg class="w-3.5 h-3.5 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              <span>{{ t.history.showInFolder }}</span>
            </button>

            <button
              type="button"
              class="p-1.5 rounded-lg text-[#7D8592] hover:text-[#FF453A] hover:bg-[#FF453A]/15 transition duration-150"
              :title="t.history.deleteRecord"
              @click="deleteItem(item.id)"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="glass-card p-12 rounded-[22px] text-center space-y-3 max-w-lg mx-auto mt-12">
        <svg class="w-10 h-10 text-[#7D8592] mx-auto opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <h3 class="text-base font-semibold text-[#FFFFFF]">{{ t.history.emptyTitle }}</h3>
        <p class="text-xs text-[#B8C0CC]">
          {{ t.history.emptyDesc }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useHistoryStore } from '../stores/history'
import { useDeviceStore } from '../stores/device'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'
import { useRouter } from 'vue-router'

const historyStore = useHistoryStore()
const deviceStore = useDeviceStore()
const { t } = useI18n()
const { showToast } = useNotifications()
const router = useRouter()

onMounted(async () => {
  await historyStore.fetchHistory()
})

async function installToDevice(path: string) {
  try {
    router.push('/apps')
    setTimeout(() => {
      deviceStore.installIPA(path)
    }, 300)
  } catch (err: any) {
    showToast('Installation Failed', err.message || err, 'error')
  }
}

function revealInExplorer(path: string) {
  historyStore.revealInExplorer(path)
}

async function copyPath(path: string) {
  try {
    await navigator.clipboard.writeText(path)
    showToast(t.value.history.copyPath, path, 'info')
  } catch {
    showToast('Copy Failed', path, 'error')
  }
}

async function deleteItem(id: string) {
  await historyStore.deleteItem(id)
  showToast(t.value.history.deleteRecord, id, 'info')
}

async function clearHistory() {
  await historyStore.clearHistory()
  showToast(t.value.history.clearedTitle, t.value.history.clearedDesc, 'info')
}

function statusBadgeClass(status: string) {
  switch (status) {
    case 'completed':
      return 'bg-[#30D158]/15 text-[#30D158] border border-[#30D158]/30'
    case 'failed':
      return 'bg-[#FF453A]/15 text-[#FF453A] border border-[#FF453A]/30'
    case 'cancelled':
      return 'bg-white/[0.08] text-[#7D8592]'
    default:
      return 'bg-[#0A84FF]/15 text-[#0A84FF] border border-[#0A84FF]/30'
  }
}
</script>
