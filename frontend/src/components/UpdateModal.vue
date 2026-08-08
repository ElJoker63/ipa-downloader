<template>
  <div v-if="visible" class="fixed inset-0 z-[100] flex items-center justify-center bg-black/80 backdrop-blur-xl p-6">
    <div class="w-full max-w-md bg-[#1C1C1E] border border-white/10 rounded-3xl shadow-2xl overflow-hidden flex flex-col">
      <!-- Header -->
      <div class="p-6 border-b border-white/10 flex items-center space-x-4 bg-white/[0.02]">
        <div class="w-12 h-12 rounded-2xl bg-[#0A84FF] flex items-center justify-center text-white shrink-0 shadow-lg shadow-[#0A84FF]/20">
          <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
          </svg>
        </div>
        <div>
          <h3 class="text-xl font-bold text-white">Update Required</h3>
          <p class="text-xs text-[#8E8E93] mt-0.5 font-medium">A new version of IPA Downloader is available.</p>
        </div>
      </div>

      <!-- Content -->
      <div class="p-6 space-y-6">
        <div class="flex items-center justify-between">
          <div class="space-y-1">
            <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">Latest Version</div>
            <div class="text-lg font-bold text-[#30D158]">v{{ updateInfo?.latestVersion }}</div>
          </div>
          <div class="h-8 w-px bg-white/10"></div>
          <div class="space-y-1 text-right">
            <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">Current Version</div>
            <div class="text-lg font-bold text-[#8E8E93]">v{{ updateInfo?.currentVersion }}</div>
          </div>
        </div>

        <div v-if="updateInfo?.releaseNotes" class="space-y-2">
          <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">What's New</div>
          <div class="max-h-32 overflow-y-auto rounded-xl bg-white/[0.03] border border-white/5 p-3 text-xs text-[#B8C0CC] leading-relaxed whitespace-pre-wrap">
            {{ updateInfo.releaseNotes }}
          </div>
        </div>

        <!-- Progress Area -->
        <div v-if="updating" class="space-y-3 pt-2">
          <div class="flex items-center justify-between text-xs mb-1">
            <span class="font-medium text-white">{{ statusMessage }}</span>
            <span class="font-mono font-bold text-[#0A84FF]">{{ progress }}%</span>
          </div>
          <div class="h-2 w-full bg-white/5 rounded-full overflow-hidden border border-white/10 relative">
            <div
              class="h-full bg-gradient-to-r from-[#0A84FF] to-[#5E5CE6] rounded-full transition-all duration-300"
              :style="{ width: `${progress}%` }"
            ></div>
            <div class="absolute inset-0 bg-gradient-to-r from-transparent via-white/10 to-transparent animate-shimmer"></div>
          </div>
        </div>
      </div>

      <!-- Footer -->
      <div class="p-6 border-t border-white/10 bg-white/[0.01]">
        <button
          v-if="!updating"
          @click="startUpdate"
          class="w-full py-3.5 rounded-2xl bg-[#0A84FF] hover:bg-[#0071E3] text-white font-bold text-sm shadow-xl shadow-[#0A84FF]/20 transition-all active:scale-[0.98]"
        >
          Update and Restart Now
        </button>
        <div v-else class="text-center text-[11px] text-[#8E8E93] italic animate-pulse">
          Please do not close the application while the update is in progress.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { WailsService } from '../services/wails'
import { useModalStore } from '../stores/modal'
import type { UpdateInfo } from '../types'

const visible = ref(false)
const updating = ref(false)
const progress = ref(0)
const updateInfo = ref<UpdateInfo | null>(null)
const statusMessage = ref('Downloading update...')
const modalStore = useModalStore()


onMounted(async () => {
  try {
    const info = await WailsService.checkForUpdate()
    if (info && info.available && info.mandatory) {
      updateInfo.value = info
      visible.value = true
    }
  } catch (err) {
    console.error('Failed to check for updates:', err)
  }

  WailsService.onEvent('update:progress', (p: number) => {
    progress.value = p
    if (p >= 100) {
      statusMessage.value = 'Applying update...'
    }
  })
})

async function startUpdate() {
  if (!updateInfo.value?.downloadUrl) return
  updating.value = true
  try {
    await WailsService.applyUpdate(updateInfo.value.downloadUrl)
  } catch (err: any) {
    modalStore.alert('Update Failed', err.message || String(err), 'error')
    updating.value = false
  }
}

</script>
