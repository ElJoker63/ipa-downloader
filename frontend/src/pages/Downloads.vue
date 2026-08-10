<template>
  <div class="max-w-6xl mx-auto space-y-6 flex flex-col h-full animate-slide-up font-sans">
    <!-- Header Area -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4 shrink-0">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
          {{ t.downloads.title }}
        </h1>
        <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
          {{ t.downloads.subtitle }}
        </p>
      </div>

      <div class="flex items-center space-x-3">
        <!-- Destination Folder Display -->
        <div class="hidden lg:flex items-center space-x-2 px-3 py-1.5 rounded-xl bg-white/[0.04] border border-white/[0.08]">
          <svg class="w-3.5 h-3.5 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          <div class="flex flex-col">
            <span class="text-[9px] uppercase font-bold text-[#7D8592] leading-none">{{ t.downloads.destinationFolder }}</span>
            <span class="text-[11px] text-[#B8C0CC] font-mono truncate max-w-[180px]">{{ settingsStore.settings.defaultDownloadFolder || 'Not Set' }}</span>
          </div>
        </div>

        <button
          @click="downloadsStore.clearCompleted"
          class="btn-secondary text-xs px-3.5 py-2 text-[#7D8592] hover:text-[#FFFFFF]"
        >
          {{ t.downloads.clearCompleted }}
        </button>
      </div>
    </div>

    <div class="flex-1 min-h-0 overflow-y-auto space-y-8 pr-1">
      <!-- Active Transfers Section -->
      <section v-if="downloadsStore.activeDownloads.length > 0" class="space-y-4">
        <div class="flex items-center justify-between px-1">
          <div class="flex items-center space-x-3">
            <span class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC]">{{ t.downloads.activeTransfers }}</span>
            <span
              class="px-2 py-0.5 text-[10px] font-medium rounded-full bg-[#0A84FF]/20 text-[#0A84FF] border border-[#0A84FF]/30"
            >
              {{ downloadsStore.activeCount }}
            </span>
          </div>
          <div v-if="downloadsStore.activeCount > 0" class="text-xs font-mono text-[#30D158] font-medium">
            {{ t.downloads.totalSpeed }} {{ downloadsStore.totalSpeedFormatted }}
          </div>
        </div>

        <div class="grid grid-cols-1 gap-4">
          <div
            v-for="task in downloadsStore.activeDownloads"
            :key="task.id"
            class="glass-card p-5 rounded-[22px] border border-white/[0.08] hover:border-white/20 transition-all flex flex-col space-y-4 shadow-xl"
          >
            <div class="flex items-start justify-between gap-4">
              <div class="flex items-center space-x-4 min-w-0">
                <img
                  :src="task.artworkUrl || 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'"
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

              <div class="flex items-center space-x-1 shrink-0">
                <button
                  v-if="task.status === 'downloading'"
                  @click="downloadsStore.pauseDownload(task.id)"
                  class="p-2 rounded-xl bg-white/[0.08] hover:bg-white/[0.14] text-[#B8C0CC] hover:text-white transition duration-150"
                  title="Pause"
                >
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </button>
                <button
                  v-else-if="task.status === 'paused'"
                  @click="downloadsStore.resumeDownload(task.id)"
                  class="p-2 rounded-xl bg-[#30D158]/15 hover:bg-[#30D158]/25 text-[#30D158] border border-[#30D158]/30 transition duration-150"
                  title="Resume"
                >
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z" /><path d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>
                </button>
                <button
                  @click="downloadsStore.cancelDownload(task.id)"
                  class="p-2 rounded-xl bg-[#FF453A]/15 hover:bg-[#FF453A]/25 text-[#FF453A] border border-[#FF453A]/30 transition duration-150"
                  title="Cancel"
                >
                  <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2"><path d="M6 18L18 6M6 6l12 12" /></svg>
                </button>
              </div>
            </div>

            <!-- Signing Stage Overlay -->
            <div v-if="task.status === 'signing'" class="px-4 py-3 rounded-[12px] bg-[#0A84FF]/15 border border-[#0A84FF]/30 flex items-center space-x-3 animate-pulse">
              <svg class="animate-spin h-4 w-4 text-[#64D2FF] shrink-0" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span class="text-xs font-semibold text-white truncate">{{ t.downloads.signingNotice }}</span>
            </div>

            <div class="space-y-2">
              <div class="w-full bg-white/5 rounded-full h-2 overflow-hidden border border-white/5">
                <div
                  class="h-full rounded-full transition-all duration-300"
                  :class="task.status === 'paused' ? 'bg-[#FFD60A]' : 'bg-[#0A84FF] shadow-[0_0_12px_rgba(10,132,255,0.5)]'"
                  :style="{ width: `${task.progress}%` }"
                ></div>
              </div>
              <div class="flex items-center justify-between text-xs text-[#B8C0CC] font-mono">
                <div class="flex items-center space-x-3">
                  <span class="font-bold text-[#FFFFFF]">{{ task.status === 'signing' ? '100.0%' : `${task.progress.toFixed(1)}%` }}</span>
                  <span>{{ formatBytes(task.downloadedBytes) }} / {{ formatBytes(task.totalBytes) }}</span>
                </div>
                <div class="flex items-center space-x-3">
                  <span v-if="task.status === 'signing'" class="text-[#64D2FF] font-semibold flex items-center space-x-1.5 animate-pulse">
                    <svg class="w-3.5 h-3.5 text-[#64D2FF] shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                      <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                    </svg>
                    <span>{{ t.downloads.signingFairPlay }}</span>
                  </span>
                  <span v-else-if="task.status === 'downloading'" class="text-[#30D158] font-semibold">{{ task.formattedSpeed || t.common.streaming }}</span>
                  <span v-else-if="task.status === 'paused'" class="text-[#FFD60A] font-semibold">{{ t.downloads.paused }}</span>
                  <span v-if="task.formattedETA && task.status === 'downloading'" class="text-[#B8C0CC]">{{ t.downloads.eta }} {{ task.formattedETA }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
      </section>

      <!-- No Active Downloads Empty Card -->
      <div v-else class="glass-card p-10 rounded-[18px] text-center space-y-2">
        <div class="w-14 h-14 rounded-full bg-white/[0.05] border border-white/10 flex items-center justify-center mx-auto text-[#7D8592]">
          <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
        </div>
        <p class="text-sm font-semibold text-[#FFFFFF]">{{ t.downloads.noActive }}</p>
        <p class="text-xs text-[#B8C0CC]">{{ t.downloads.noActiveDesc }}</p>
      </div>

      <!-- Completed & History Section -->
      <section v-if="downloadsStore.completedDownloads.length > 0" class="space-y-4">
        <h2 class="text-xs font-semibold uppercase tracking-wider text-[#B8C0CC] px-1">{{ t.downloads.completedTitle }}</h2>
        <div class="space-y-3">
          <div
            v-for="task in downloadsStore.completedDownloads"
            :key="task.id"
            class="p-4 rounded-[18px] bg-[#171A21]/40 border border-white/[0.06] hover:border-white/10 transition-colors flex items-center justify-between group"
          >
            <div class="flex items-center space-x-4 min-w-0">
              <img :src="task.artworkUrl" class="w-10 h-10 rounded-xl object-cover border border-white/10 shrink-0" />
              <div class="min-w-0">
                <div class="flex items-center space-x-2">
                  <h4 class="text-sm font-semibold truncate text-[#FFFFFF]">{{ task.appName }}</h4>
                  <span
                    class="px-2 py-0.5 text-[10px] font-medium rounded-full uppercase"
                    :class="statusBadgeClass(task.status)"
                  >
                    {{ statusText(task.status) }}
                  </span>
                  <span
                    v-if="task.status === 'completed' && task.checksum"
                    class="px-2 py-0.5 text-[10px] font-bold rounded-full bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30 flex items-center space-x-1"
                  >
                    <svg class="w-2.5 h-2.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="3" d="M5 13l4 4L19 7" />
                    </svg>
                    <span>{{ t.common.verified }}</span>
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
                    {{ t.downloads.copyError }}
                  </button>
                </div>
              </div>
            </div>

            <div class="flex items-center space-x-2 shrink-0">
              <button
                v-if="deviceStore.devices.length > 0 && task.status === 'completed' && task.type === 'app'"
                type="button"
                class="btn-primary text-xs px-3.5 py-1.5 flex items-center space-x-1.5"
                @click="handleInstallClick(task.destinationPath)"
              >

                <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
                <span>{{ t.common.install }}</span>
              </button>
              <button
                v-if="task.status === 'failed'"

                type="button"
                class="btn-primary text-xs px-3.5 py-1.5"
                @click="downloadsStore.retryDownload(task.id)"
              >
                {{ t.common.retry }}
              </button>
              <button
                type="button"
                class="btn-secondary text-xs px-3.5 py-1.5"
                @click="revealInExplorer(task.destinationPath)"
              >
                {{ t.downloads.showInFolder }}
              </button>
              <button
                type="button"
                class="p-1.5 rounded-lg text-[#7D8592] hover:text-[#FF453A] hover:bg-[#FF453A]/15 transition duration-150"
                :title="t.downloads.deleteFile"
                @click="handleDeleteFile(task.destinationPath, task.id)"
              >

                <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>

          </div>
        </div>
      </section>
    </div>

    <!-- Device Picker Modal -->
    <div
      v-if="ipaToInstall && deviceStore.devices.length > 1"
      class="fixed inset-0 bg-black/60 backdrop-blur-md z-[80] flex items-center justify-center p-6"
      @click.self="ipaToInstall = null"
    >
      <div class="w-full max-w-sm bg-[#1C1C1E] border border-white/10 rounded-3xl shadow-2xl overflow-hidden p-6 space-y-6">
        <div class="text-center space-y-2">
          <h3 class="text-lg font-bold text-white">{{ t.common.selectTargetDevice }}</h3>
          <p class="text-xs text-[#8E8E93]">{{ t.common.multipleDevicesDesc }}</p>
        </div>

        <div class="space-y-2 max-h-60 overflow-y-auto pr-1">
          <button
            v-for="dev in deviceStore.devices"
            :key="dev.udid"
            @click="installToDevice(ipaToInstall!, dev.udid)"
            class="w-full p-4 rounded-2xl bg-white/[0.04] border border-white/10 hover:bg-white/[0.08] hover:border-white/20 transition text-left flex items-center justify-between group"
          >
            <div class="flex items-center space-x-3">
              <div class="w-10 h-10 rounded-xl bg-[#0A84FF]/10 flex items-center justify-center text-[#0A84FF]">
                <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
              </div>
              <div>
                <div class="text-sm font-bold text-white group-hover:text-[#0A84FF] transition-colors">{{ dev.name }}</div>
                <div class="text-[10px] text-[#8E8E93] font-mono">{{ dev.iosVersion }} • {{ truncateUDID(dev.udid) }}</div>
              </div>
            </div>
            <svg class="w-4 h-4 text-white/20 group-hover:text-[#0A84FF] transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" /></svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { useDownloadsStore } from '../stores/downloads'
import { useSettingsStore } from '../stores/settings'
import { useDeviceStore } from '../stores/device'
import { useModalStore } from '../stores/modal'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'
import { WailsService } from '../services/wails'

const downloadsStore = useDownloadsStore()
const settingsStore = useSettingsStore()
const deviceStore = useDeviceStore()
const modalStore = useModalStore()
const { t } = useI18n()
const { showToast } = useNotifications()

const ipaToInstall = ref<string | null>(null)

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

function statusText(status: string) {
  switch (status) {
    case 'queued': return t.value.common.queued
    case 'downloading': return t.value.common.downloading
    case 'signing': return t.value.downloads.signingFairPlay
    case 'completed': return t.value.common.verified
    case 'failed': return t.value.common.failed
    case 'cancelled': return t.value.common.cancelled
    case 'paused': return t.value.downloads.paused
    default: return status
  }
}

function formatBytes(bytes: number, decimals = 2) {
  if (bytes === 0) return '0 Bytes'
  const k = 1024
  const dm = decimals < 0 ? 0 : decimals
  const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i]
}

function revealInExplorer(path: string) {
  WailsService.revealInExplorer(path)
}

function handleInstallClick(path: string) {
  if (deviceStore.devices.length === 0) {
    showToast(t.value.downloads.noDeviceToast, t.value.downloads.noDeviceToastDesc, 'error')
    return
  }
  if (deviceStore.devices.length === 1) {
    installToDevice(path, deviceStore.devices[0].udid)
  } else {
    ipaToInstall.value = path
  }
}

async function installToDevice(path: string, udid: string) {
  ipaToInstall.value = null
  try {
    showToast(t.value.downloads.installStarted, t.value.downloads.installStartedDesc, 'info')
    await deviceStore.installIPA(path, udid)
  } catch (err: any) {
    showToast(t.value.downloads.installFailed, err.message || String(err), 'error')
  }
}

function handleDeleteFile(path: string, id: string) {
  modalStore.confirm(
    t.value.downloads.deleteFile,
    t.value.downloads.deleteFilePrompt,
    async () => {
      try {
        await WailsService.deleteFile(path)
        downloadsStore.removeTask(id)
        showToast(t.value.downloads.deletedToast, t.value.downloads.deletedDesc, 'info')
      } catch (err: any) {
        showToast(t.value.downloads.deleteFailed, err.message || String(err), 'error')
      }
    },
    t.value.common.delete
  )
}

async function copyErrorText(err: string) {
  try {
    await navigator.clipboard.writeText(err)
    showToast(t.value.downloads.copyErrorTitle, t.value.downloads.copied, 'success')
  } catch {
    showToast('Copy Failed', '', 'error')
  }
}

function truncateUDID(udid: string) {
  if (!udid || udid.length <= 12) return udid
  return `${udid.slice(0, 6)}...${udid.slice(-6)}`
}
</script>
