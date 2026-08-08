<template>
  <div class="h-full flex flex-col p-6 space-y-6 overflow-y-auto" @dragover.prevent @drop.prevent="handleFileDrop">
    <!-- Top Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-white flex items-center space-x-3">
          <span>{{ t.apps?.title || 'Device Apps' }}</span>
          <span
            v-if="deviceStore.isConnected"
            class="px-2.5 py-0.5 text-xs font-medium rounded-full bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30 flex items-center space-x-1.5"
          >
            <span class="w-2 h-2 rounded-full bg-[#30D158] animate-pulse"></span>
            <span>Connected</span>
          </span>
          <span
            v-else
            class="px-2.5 py-0.5 text-xs font-medium rounded-full bg-white/10 text-[#8E8E93] border border-white/10 flex items-center space-x-1.5"
          >
            <span class="w-2 h-2 rounded-full bg-[#8E8E93]"></span>
            <span>Disconnected</span>
          </span>
        </h1>
        <p class="text-sm text-[#8E8E93] mt-1">
          {{ t.apps?.subtitle || 'Manage iOS applications and install IPAs on USB connected devices' }}
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center space-x-3">
        <button
          @click="handleRefresh"
          :disabled="deviceStore.isLoadingApps"
          class="px-3.5 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/[0.08] text-sm font-medium text-white transition flex items-center space-x-2 disabled:opacity-50"
        >
          <svg class="w-4 h-4 text-[#8E8E93]" :class="{ 'animate-spin': deviceStore.isLoadingApps }" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span>Refresh</span>
        </button>

        <button
          @click="triggerIPAInstall"
          :disabled="!deviceStore.isConnected || deviceStore.isInstalling"
          class="px-4 py-2 rounded-xl bg-[#0A84FF] hover:bg-[#0071E3] text-sm font-medium text-white shadow-lg shadow-[#0A84FF]/25 transition flex items-center space-x-2 disabled:opacity-40 disabled:cursor-not-allowed"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
          </svg>
          <span>Install IPA</span>
        </button>
      </div>
    </div>

    <!-- Device Info Card (When Connected) -->
    <div v-if="deviceStore.device && deviceStore.device.isConnected" class="p-5 rounded-2xl bg-[#171A21]/90 border border-white/[0.08] backdrop-blur-xl shadow-xl flex flex-col md:flex-row md:items-center justify-between gap-4">
      <div class="flex items-center space-x-4">
        <!-- Device Icon -->
        <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-[#0A84FF]/20 to-[#5E5CE6]/20 border border-white/10 flex items-center justify-center shrink-0">
          <svg class="w-7 h-7 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
        </div>

        <!-- Device Specs -->
        <div>
          <h2 class="text-base font-semibold text-white flex items-center space-x-2">
            <span>{{ deviceStore.device.name }}</span>
            <span class="text-xs px-2 py-0.5 rounded-md bg-white/10 text-[#8E8E93] font-mono">{{ deviceStore.device.model }}</span>
          </h2>
          <div class="text-xs text-[#8E8E93] mt-1 space-x-3 flex items-center">
            <span>iOS {{ deviceStore.device.iosVersion }} (Build {{ deviceStore.device.buildVersion }})</span>
            <span>•</span>
            <span class="font-mono text-[11px]">UDID: {{ truncateUDID(deviceStore.device.udid) }}</span>
          </div>
        </div>
      </div>

      <!-- Quick Badges & Pair Button -->
      <div class="flex items-center space-x-3 shrink-0">
        <div class="px-3 py-1.5 rounded-xl bg-white/[0.04] border border-white/[0.06] text-xs text-[#8E8E93]">
          <span class="text-white font-medium">{{ deviceStore.userApps.length }}</span> User Apps
        </div>

        <button
          v-if="!deviceStore.device.isPaired"
          @click="handlePair"
          class="px-3 py-1.5 rounded-xl bg-[#FF9F0A]/20 text-[#FF9F0A] border border-[#FF9F0A]/30 text-xs font-medium hover:bg-[#FF9F0A]/30 transition"
        >
          Trust & Pair Device
        </button>
      </div>
    </div>

    <!-- No Device Connected Banner -->
    <div v-else class="p-8 rounded-2xl bg-[#171A21]/60 border border-white/[0.08] backdrop-blur-xl text-center flex flex-col items-center justify-center space-y-4">
      <div class="w-16 h-16 rounded-full bg-white/[0.05] border border-white/10 flex items-center justify-center text-[#8E8E93]">
        <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
        </svg>
      </div>
      <div class="max-w-md">
        <h3 class="text-lg font-semibold text-white">No iOS Device Connected</h3>
        <p class="text-sm text-[#8E8E93] mt-1 leading-relaxed">
          Connect your iPhone or iPad via USB cable to inspect installed apps and install .ipa packages directly.
        </p>
      </div>

      <!-- Instructions & Troubleshooting -->
      <div class="pt-2 text-left max-w-lg w-full bg-white/[0.03] border border-white/[0.06] rounded-xl p-4 space-y-2 text-xs text-[#8E8E93]">
        <div class="font-medium text-white flex items-center space-x-1.5">
          <svg class="w-4 h-4 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>Connection Requirements:</span>
        </div>
        <ul class="list-disc list-inside space-y-1 pl-1">
          <li>Unlock your device screen and tap <strong>"Trust This Computer"</strong> when prompted.</li>
          <li>On Windows, ensure <strong>Apple Mobile Device Support</strong> (iTunes / Apple Devices app) is installed.</li>
          <li>For iOS 17+, lockdown querying works directly; tunnel service may be required for full dev debug tools.</li>
        </ul>
      </div>
    </div>

    <!-- Drag & Drop Zone for IPA -->
    <div
      v-if="deviceStore.isConnected"
      @click="triggerIPAInstall"
      class="border-2 border-dashed border-white/10 hover:border-[#0A84FF]/50 hover:bg-[#0A84FF]/5 rounded-2xl p-6 text-center cursor-pointer transition flex flex-col items-center justify-center space-y-2"
    >
      <svg class="w-8 h-8 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
        <path stroke-linecap="round" stroke-linejoin="round" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
      </svg>
      <div class="text-sm font-medium text-white">Drag & Drop .ipa file here to install</div>
      <div class="text-xs text-[#8E8E93]">or click to select file from file manager</div>
    </div>

    <!-- Installed Apps Section (When Connected) -->
    <div v-if="deviceStore.isConnected" class="space-y-4">
      <!-- Search & Tab Filter Bar -->
      <div class="flex items-center justify-between gap-4">
        <!-- Search Input -->
        <div class="relative flex-1 max-w-md">
          <svg class="w-4 h-4 absolute left-3.5 top-1/2 -translate-y-1/2 text-[#8E8E93]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search installed apps by name or bundle ID..."
            class="w-full bg-white/[0.06] border border-white/[0.1] rounded-xl pl-10 pr-4 py-2 text-sm text-white placeholder-[#8E8E93] focus:outline-none focus:border-[#0A84FF] transition"
          />
        </div>

        <!-- Filter Tabs: User vs System -->
        <div class="flex p-1 bg-white/[0.06] border border-white/[0.08] rounded-xl text-xs font-medium">
          <button
            @click="switchTab('user')"
            class="px-3 py-1.5 rounded-lg transition"
            :class="deviceStore.activeTab === 'user' ? 'bg-[#0A84FF] text-white shadow-sm' : 'text-[#8E8E93] hover:text-white'"
          >
            User Apps ({{ deviceStore.userApps.length }})
          </button>
          <button
            @click="switchTab('system')"
            class="px-3 py-1.5 rounded-lg transition"
            :class="deviceStore.activeTab === 'system' ? 'bg-[#0A84FF] text-white shadow-sm' : 'text-[#8E8E93] hover:text-white'"
          >
            System Apps ({{ deviceStore.systemApps.length }})
          </button>
        </div>
      </div>

      <!-- Apps Grid / List -->
      <div v-if="deviceStore.isLoadingApps" class="py-12 flex justify-center items-center">
        <svg class="w-8 h-8 text-[#0A84FF] animate-spin" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <div v-else-if="filteredApps.length === 0" class="py-12 text-center text-[#8E8E93] text-sm">
        No installed apps match your search filter.
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-3">
        <div
          v-for="app in filteredApps"
          :key="app.bundleId"
          class="p-4 rounded-xl bg-[#171A21]/70 border border-white/[0.08] hover:border-white/20 transition flex items-center justify-between group"
        >
          <div class="flex items-center space-x-3 overflow-hidden">
            <!-- App Icon Placeholder -->
            <div class="w-10 h-10 rounded-xl bg-gradient-to-tr from-white/10 to-white/5 border border-white/10 flex items-center justify-center font-bold text-sm text-white shrink-0">
              {{ app.name.charAt(0).toUpperCase() }}
            </div>
            <div class="overflow-hidden">
              <div class="text-sm font-semibold text-white truncate" :title="app.name">{{ app.name }}</div>
              <div class="text-xs text-[#8E8E93] font-mono truncate" :title="app.bundleId">{{ app.bundleId }}</div>
              <div class="text-[11px] text-[#8E8E93] mt-0.5">
                v{{ app.version || app.shortVersion || '1.0' }}
                <span v-if="app.size > 0">• {{ formatAppSize(app.size) }}</span>
              </div>
            </div>
          </div>

          <!-- Uninstall Button (For User Apps) -->
          <button
            v-if="deviceStore.activeTab === 'user'"
            @click="handleUninstall(app)"
            class="px-2.5 py-1.5 rounded-lg bg-red-500/10 text-red-400 opacity-0 group-hover:opacity-100 hover:bg-red-500/20 transition text-xs font-medium shrink-0 ml-2"
          >
            Uninstall
          </button>
        </div>
      </div>
    </div>

    <!-- Installation Progress Modal Overlay -->
    <div
      v-if="deviceStore.isInstalling || deviceStore.installProgress"
      class="fixed inset-0 bg-black/70 backdrop-blur-md z-50 flex items-center justify-center p-4"
    >
      <div class="w-full max-w-md p-6 rounded-2xl bg-[#171A21] border border-white/10 shadow-2xl space-y-4 relative">
        <!-- Header -->
        <div class="flex items-center justify-between">
          <h3 class="text-base font-semibold text-white flex items-center space-x-2">
            <svg
              class="w-5 h-5"
              :class="deviceStore.installError ? 'text-red-400' : deviceStore.installProgress?.phase === 'Complete' ? 'text-[#30D158]' : 'text-[#0A84FF]'"
              fill="none"
              viewBox="0 0 24 24"
              stroke="currentColor"
              stroke-width="2"
            >
              <path v-if="deviceStore.installError" stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              <path v-else-if="deviceStore.installProgress?.phase === 'Complete'" stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              <path v-else stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
            </svg>
            <span>
              {{ deviceStore.installError ? 'Installation Failed' : deviceStore.installProgress?.phase === 'Complete' ? 'Installation Complete' : 'Installing Application' }}
            </span>
          </h3>

          <div class="flex items-center space-x-2">
            <span class="text-xs font-mono font-bold" :class="deviceStore.installError ? 'text-red-400' : 'text-[#0A84FF]'">
              {{ deviceStore.installProgress?.percent || 0 }}%
            </span>

            <!-- Close 'X' Button -->
            <button
              @click="deviceStore.closeInstallModal()"
              class="w-7 h-7 rounded-lg bg-white/5 hover:bg-white/10 flex items-center justify-center text-[#8E8E93] hover:text-white transition"
              title="Close modal"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Progress Bar -->
        <div class="w-full bg-white/10 rounded-full h-2 overflow-hidden">
          <div
            class="h-full rounded-full transition-all duration-300"
            :class="deviceStore.installError ? 'bg-red-500' : deviceStore.installProgress?.phase === 'Complete' ? 'bg-[#30D158]' : 'bg-[#0A84FF]'"
            :style="{ width: `${deviceStore.installProgress?.percent || 0}%` }"
          ></div>
        </div>

        <!-- Phase & Message -->
        <div class="text-xs text-[#8E8E93] flex justify-between items-center">
          <span class="font-medium text-white">{{ deviceStore.installProgress?.phase }}</span>
          <span class="truncate ml-2 max-w-[240px]">{{ deviceStore.installProgress?.message }}</span>
        </div>

        <!-- Error Alert if Failed -->
        <div v-if="deviceStore.installError" class="p-3 rounded-xl bg-red-500/10 border border-red-500/20 text-xs text-red-300 space-y-1">
          <div class="font-semibold text-red-400">Error details:</div>
          <div>{{ deviceStore.installError }}</div>
        </div>

        <!-- Bottom Actions (Dismiss / Confirm Button) -->
        <div v-if="deviceStore.installError || deviceStore.installProgress?.phase === 'Complete' || !deviceStore.isInstalling" class="pt-2 flex justify-end">
          <button
            @click="deviceStore.closeInstallModal()"
            class="px-4 py-2 rounded-xl bg-white/10 hover:bg-white/20 text-xs font-semibold text-white transition"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useDeviceStore } from '../stores/device'
import { useI18n } from '../i18n'
import type { InstalledApp } from '../types'

const deviceStore = useDeviceStore()
const { t } = useI18n()

const searchQuery = ref('')

const filteredApps = computed(() => {
  const query = searchQuery.value.toLowerCase().trim()
  const list = deviceStore.activeTab === 'user' ? deviceStore.userApps : deviceStore.systemApps
  if (!query) return list
  return list.filter(
    (app) => app.name.toLowerCase().includes(query) || app.bundleId.toLowerCase().includes(query)
  )
})

onMounted(() => {
  deviceStore.initListeners()
  deviceStore.checkDevice()
})

function handleRefresh() {
  deviceStore.checkDevice()
}

function handlePair() {
  deviceStore.pairDevice().catch((err) => {
    alert(err.message || 'Pairing failed')
  })
}

function triggerIPAInstall() {
  deviceStore.installIPA().catch((err) => {
    // Error handled in store state
  })
}

function handleUninstall(app: InstalledApp) {
  if (confirm(`Are you sure you want to uninstall "${app.name}" (${app.bundleId}) from your device?`)) {
    deviceStore.uninstallApp(app.bundleId).catch((err) => {
      alert(`Failed to uninstall: ${err.message || err}`)
    })
  }
}

function handleFileDrop(e: DragEvent) {
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    const file = files[0]
    if (file.name.endsWith('.ipa')) {
      // In Wails, file.path contains the absolute local OS path
      const path = (file as any).path
      if (path) {
        deviceStore.installIPA(path)
      }
    }
  }
}

function switchTab(tab: 'user' | 'system') {
  deviceStore.activeTab = tab
  deviceStore.fetchApps()
}

function truncateUDID(udid: string) {
  if (!udid || udid.length <= 12) return udid
  return `${udid.slice(0, 6)}...${udid.slice(-6)}`
}

function formatAppSize(bytes: number) {
  if (!bytes) return ''
  const mb = bytes / (1024 * 1024)
  if (mb >= 1024) {
    return `${(mb / 1024).toFixed(1)} GB`
  }
  return `${mb.toFixed(1)} MB`
}
</script>
