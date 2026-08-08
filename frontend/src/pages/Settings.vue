<template>
  <div class="max-w-4xl mx-auto space-y-6 animate-slide-up font-sans">
    <!-- Header Section -->
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
        Settings & Preferences
      </h1>
      <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
        Configure transfer queues, default storage folders, FairPlay DRM automation, and diagnostics.
      </p>
    </div>

    <div class="space-y-5 pb-8">
      <!-- General & Downloads Preferences -->
      <div class="glass-card p-6 rounded-[18px] space-y-5">
        <h2 class="text-sm font-semibold uppercase tracking-wider text-[#B8C0CC]">General & Downloads</h2>

        <!-- Default Download Folder -->
        <div class="space-y-2">
          <label class="text-xs font-medium text-[#FFFFFF]">Default Download Directory</label>
          <div class="flex items-center space-x-2.5">
            <input
              v-model="settingsStore.settings.defaultDownloadFolder"
              type="text"
              readonly
              class="glass-input flex-1 px-3.5 py-2.5 text-xs font-mono select-all"
            />
            <button
              type="button"
              class="btn-secondary text-xs px-4 py-2.5 shrink-0 flex items-center space-x-1.5"
              @click="browseFolder"
            >
              <svg class="w-3.5 h-3.5 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
              </svg>
              <span>Browse Folder</span>
            </button>
          </div>
        </div>

        <!-- Max Concurrent Downloads Slider -->
        <div class="space-y-2 pt-1">
          <div class="flex items-center justify-between">
            <label class="text-xs font-medium text-[#FFFFFF]">Concurrent Transfers Limit</label>
            <span class="text-xs font-mono font-semibold text-[#0A84FF]">
              {{ settingsStore.settings.maxConcurrentDownloads }} simultaneous jobs
            </span>
          </div>
          <input
            v-model.number="settingsStore.settings.maxConcurrentDownloads"
            type="range"
            min="1"
            max="10"
            class="w-full accent-[#0A84FF] cursor-pointer"
            @change="saveSettings"
          />
          <div class="flex justify-between text-[10px] text-[#7D8592] font-mono">
            <span>1 (Conservative)</span>
            <span>3 (Recommended)</span>
            <span>10 (High Speed)</span>
          </div>
        </div>

        <!-- Automation Toggles -->
        <div class="space-y-3 pt-2">
          <!-- Auto Acquire License -->
          <div class="flex items-center justify-between p-3.5 rounded-[14px] bg-white/[0.04] border border-white/[0.08]">
            <div class="pr-4">
              <div class="text-xs font-semibold text-[#FFFFFF]">Auto-Acquire Free App Licenses</div>
              <p class="text-[11px] text-[#B8C0CC] mt-0.5">Automatically trigger iTunes license purchase if your Apple ID has not purchased the app before.</p>
            </div>
            <input
              v-model="settingsStore.settings.autoAcquireLicense"
              type="checkbox"
              class="w-5 h-5 rounded-md text-[#0A84FF] bg-white/[0.08] border-white/[0.18] focus:ring-0 cursor-pointer shrink-0"
              @change="saveSettings"
            />
          </div>

          <!-- Remember Credentials -->
          <div class="flex items-center justify-between p-3.5 rounded-[14px] bg-white/[0.04] border border-white/[0.08]">
            <div class="pr-4">
              <div class="text-xs font-semibold text-[#FFFFFF]">Persist Session in Encrypted Keychain</div>
              <p class="text-[11px] text-[#B8C0CC] mt-0.5">Store Apple Storefront authentication tokens securely in the local operating system keychain.</p>
            </div>
            <input
              v-model="settingsStore.settings.rememberCredentials"
              type="checkbox"
              class="w-5 h-5 rounded-md text-[#0A84FF] bg-white/[0.08] border-white/[0.18] focus:ring-0 cursor-pointer shrink-0"
              @change="saveSettings"
            />
          </div>
        </div>
      </div>

      <!-- Storage & Diagnostics Glass Card -->
      <div class="glass-card p-6 rounded-[18px] space-y-5">
        <h2 class="text-sm font-semibold uppercase tracking-wider text-[#B8C0CC]">Storage & Diagnostics</h2>

        <div class="flex items-center justify-between p-3.5 rounded-[14px] bg-white/[0.04] border border-white/[0.08]">
          <div>
            <div class="text-xs font-semibold text-[#FFFFFF]">Metadata & Search Cache</div>
            <p class="text-[11px] text-[#B8C0CC] mt-0.5">Current cached artwork and version history: <span class="font-mono text-[#64D2FF]">{{ settingsStore.cacheSize }}</span></p>
          </div>
          <button
            type="button"
            class="btn-secondary text-xs px-3.5 py-1.5"
            @click="clearCache"
          >
            Clear Cache
          </button>
        </div>

        <div class="flex items-center justify-between p-3.5 rounded-[14px] bg-white/[0.04] border border-white/[0.08]">
          <div>
            <div class="text-xs font-semibold text-[#FFFFFF]">Diagnostic Log Export</div>
            <p class="text-[11px] text-[#B8C0CC] mt-0.5">Export real-time FairPlay SINF DRM logs and network events to a text file.</p>
          </div>
          <button
            type="button"
            class="btn-primary text-xs px-4 py-1.5 shadow-sm"
            @click="exportLogs"
          >
            Export Logs
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSettingsStore } from '../stores/settings'
import { useLogsStore } from '../stores/logs'
import { useNotifications } from '../composables/useNotifications'

const settingsStore = useSettingsStore()
const logsStore = useLogsStore()
const { showToast } = useNotifications()

onMounted(async () => {
  await settingsStore.fetchSettings()
})

async function browseFolder() {
  await settingsStore.browseFolder()
  showToast('Download Folder Updated', settingsStore.settings.defaultDownloadFolder, 'info')
}

async function saveSettings() {
  await settingsStore.updateSettings(settingsStore.settings)
  showToast('Settings Saved', 'Preferences updated successfully', 'info')
}

async function clearCache() {
  await settingsStore.clearCache()
  showToast('Cache Cleared', 'Offline app cache purged', 'info')
}

async function exportLogs() {
  const path = await logsStore.exportLogs()
  if (path) {
    showToast('Logs Exported', `Saved to ${path}`, 'success')
  }
}
</script>
