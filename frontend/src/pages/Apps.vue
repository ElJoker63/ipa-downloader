<template>
  <div class="h-full flex flex-col p-6 space-y-6 overflow-y-auto" @dragover.prevent @drop.prevent="handleFileDrop">
    <!-- Top Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-white flex items-center space-x-3">
          <span>{{ t.apps?.title || 'Device Apps' }}</span>
          <span
            v-if="deviceStore.devices.length > 0"
            class="px-2.5 py-0.5 text-xs font-medium rounded-full bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30 flex items-center space-x-1.5"
          >
            <span class="w-2 h-2 rounded-full bg-[#30D158] animate-pulse"></span>
            <span>{{ deviceStore.devices.length }} {{ t.apps.connectedCount }}</span>
          </span>
          <span
            v-else
            class="px-2.5 py-0.5 text-xs font-medium rounded-full bg-white/10 text-[#8E8E93] border border-white/10 flex items-center space-x-1.5"
          >
            <span class="w-2 h-2 rounded-full bg-[#8E8E93]"></span>
            <span>{{ t.apps.noDevice }}</span>
          </span>

        </h1>
        <p class="text-sm text-[#8E8E93] mt-1">
          {{ t.apps.subtitle }}
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
          <span>{{ t.apps.refresh }}</span>
        </button>

        <button
          @click="triggerIPAInstall"
          :disabled="!deviceStore.isConnected || deviceStore.isInstalling"
          class="px-4 py-2 rounded-xl bg-[#0A84FF] hover:bg-[#0071E3] text-sm font-medium text-white shadow-lg shadow-[#0A84FF]/25 transition flex items-center space-x-2 disabled:opacity-40 disabled:cursor-not-allowed"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
          </svg>
          <span>{{ t.apps.install }}</span>
        </button>
      </div>

    </div>

    <!-- Device Selection & Info -->
    <div v-if="deviceStore.devices.length > 0" class="space-y-4">
      <!-- Device Selector (Tabs style) -->
      <div v-if="deviceStore.devices.length > 1" class="flex items-center space-x-2 overflow-x-auto pb-2 scrollbar-hide">
        <button
          v-for="dev in deviceStore.devices"
          :key="dev.udid"
          @click="deviceStore.selectedUdid = dev.udid; deviceStore.fetchApps()"
          class="px-4 py-2 rounded-2xl border transition-all shrink-0 flex items-center space-x-2"
          :class="deviceStore.selectedUdid === dev.udid
            ? 'bg-[#0A84FF] border-[#0A84FF] text-white shadow-lg shadow-[#0A84FF]/20'
            : 'bg-white/[0.04] border-white/10 text-[#8E8E93] hover:bg-white/[0.08]'"
        >
          <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
          <span class="text-xs font-semibold">{{ dev.name }}</span>
        </button>
      </div>

      <!-- Active Device Info Card -->
      <div
        v-if="deviceStore.selectedDevice"
        @click="showDeviceDetails = true"
        class="p-5 rounded-2xl bg-[#171A21]/90 border border-white/[0.08] hover:border-white/[0.15] cursor-pointer group transition-all duration-300 backdrop-blur-xl shadow-xl flex flex-col md:flex-row md:items-center justify-between gap-4"
      >
        <div class="flex items-center space-x-4">
          <!-- Device Icon -->
          <div class="w-14 h-14 rounded-2xl bg-gradient-to-br from-[#0A84FF]/20 to-[#5E5CE6]/20 border border-white/10 flex items-center justify-center shrink-0 group-hover:scale-105 transition-transform">
            <svg class="w-7 h-7 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
            </svg>
          </div>

          <!-- Device Specs -->
          <div>
            <h2 class="text-base font-semibold text-white flex items-center space-x-2">
              <span>{{ deviceStore.selectedDevice.name }}</span>
              <span class="text-xs px-2 py-0.5 rounded-md bg-white/10 text-[#8E8E93] font-mono">{{ deviceStore.selectedDevice.model }}</span>
            </h2>
            <div class="text-xs text-[#8E8E93] mt-1 space-x-3 flex items-center">
              <span>iOS {{ deviceStore.selectedDevice.iosVersion }}</span>
              <span>•</span>
              <div class="flex items-center space-x-1">
                <svg class="w-3.5 h-3.5" :class="deviceStore.selectedDevice.batteryCharging ? 'text-[#30D158]' : 'text-[#8E8E93]'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 10.5h-1.5V9a1.5 1.5 0 00-1.5-1.5H5A1.5 1.5 0 003.5 9v6a1.5 1.5 0 001.5 1.5h13a1.5 1.5 0 001.5-1.5v-1.5H21v-3z" />
                </svg>
                <span>{{ deviceStore.selectedDevice.batteryLevel }}%</span>
              </div>
              <span>•</span>
              <span>{{ formatAppSize(deviceStore.selectedDevice.storageFree) }} Free</span>
            </div>
          </div>
        </div>

        <!-- Quick Badges & Pair Button -->
        <div class="flex items-center space-x-3 shrink-0">
          <div class="px-3 py-1.5 rounded-xl bg-white/[0.04] border border-white/[0.06] text-xs text-[#8E8E93]">
            <span class="text-white font-medium">{{ deviceStore.userApps.length }}</span> {{ t.apps.userApps }}
          </div>

          <button
            v-if="!deviceStore.selectedDevice.isPaired"
            @click.stop="handlePair(deviceStore.selectedDevice.udid)"
            class="px-3 py-1.5 rounded-xl bg-[#FF9F0A]/20 text-[#FF9F0A] border border-[#FF9F0A]/30 text-xs font-medium hover:bg-[#FF9F0A]/30 transition"
          >
            {{ t.apps.trustAndPair }}
          </button>
        </div>
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
        <h3 class="text-lg font-semibold text-white">{{ t.apps.noDeviceTitle }}</h3>
        <p class="text-sm text-[#8E8E93] mt-1 leading-relaxed">
          {{ t.apps.noDeviceDesc }}
        </p>
      </div>

      <!-- Instructions & Troubleshooting -->
      <div class="pt-2 text-left max-w-lg w-full bg-white/[0.03] border border-white/[0.06] rounded-xl p-4 space-y-2 text-xs text-[#8E8E93]">
        <div class="font-medium text-white flex items-center space-x-1.5">
          <svg class="w-4 h-4 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span>{{ t.apps.requirementsTitle }}</span>
        </div>
        <ul class="list-disc list-inside space-y-1 pl-1">
          <li>{{ t.apps.requirement1 }}</li>
          <li>{{ t.apps.requirement2 }}</li>
          <li>{{ t.apps.requirement3 }}</li>
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
      <div class="text-sm font-medium text-white">{{ t.apps.dragDropText }}</div>
      <div class="text-xs text-[#8E8E93]">{{ t.apps.orClickToSelect }}</div>
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
            :placeholder="t.apps.searchPlaceholder"
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
            {{ t.apps.userApps }} ({{ deviceStore.userApps.length }})
          </button>
          <button
            @click="switchTab('system')"
            class="px-3 py-1.5 rounded-lg transition"
            :class="deviceStore.activeTab === 'system' ? 'bg-[#0A84FF] text-white shadow-sm' : 'text-[#8E8E93] hover:text-white'"
          >
            {{ t.apps.systemApps }} ({{ deviceStore.systemApps.length }})
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
            <!-- App Icon -->
            <img
              v-if="app.artworkUrl"
              :src="app.artworkUrl"
              class="w-10 h-10 rounded-xl bg-[#171A21] border border-white/10 shrink-0 object-cover"
            />
            <div v-else class="w-10 h-10 rounded-xl bg-gradient-to-tr from-white/10 to-white/5 border border-white/10 flex items-center justify-center font-bold text-sm text-white shrink-0">
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
              {{ deviceStore.installError ? 'Installation Failed' : deviceStore.installProgress?.phase === 'Complete' ? 'Installation Complete' : 'Processing Queue' }}
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

    <!-- Device Details Modal -->
    <div
      v-if="showDeviceDetails && deviceStore.selectedDevice"
      class="fixed inset-0 bg-black/80 backdrop-blur-xl z-[60] flex items-center justify-center p-6"
      @click.self="showDeviceDetails = false"
    >
      <div class="w-full max-w-2xl bg-[#1C1C1E] border border-white/10 rounded-3xl shadow-2xl overflow-hidden flex flex-col max-h-[90vh]">
        <!-- Modal Header -->
        <div class="p-6 border-b border-white/10 flex items-center justify-between bg-white/[0.02]">
          <div class="flex items-center space-x-4">
            <div class="w-12 h-12 rounded-2xl bg-[#0A84FF] flex items-center justify-center text-white">
              <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
            <div>
              <h3 class="text-xl font-bold text-white">{{ deviceStore.selectedDevice.name }}</h3>
              <p class="text-xs text-[#8E8E93] font-mono mt-0.5">{{ deviceStore.selectedDevice.model }} • iOS {{ deviceStore.selectedDevice.iosVersion }}</p>
            </div>
          </div>
          <button @click="showDeviceDetails = false" class="p-2 rounded-full hover:bg-white/10 text-[#8E8E93] transition">
            <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <!-- Modal Content -->
        <div class="flex-1 overflow-y-auto p-6 space-y-8">
          <!-- Main Specs Grid -->
          <div class="grid grid-cols-2 gap-4">
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.iosVersion }}</div>
              <div class="text-sm text-white font-mono">{{ deviceStore.selectedDevice.iosVersion }} ({{ deviceStore.selectedDevice.buildVersion }})</div>
            </div>
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.serialNumber }}</div>
              <div class="text-sm text-white font-mono">{{ deviceStore.selectedDevice.serialNumber }}</div>
            </div>
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.imei }}</div>
              <div class="text-sm text-white font-mono">{{ deviceStore.selectedDevice.imei || 'N/A' }}</div>
            </div>
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.modelName }}</div>
              <div class="text-sm text-white font-mono">{{ deviceStore.selectedDevice.modelNumber || 'N/A' }}</div>
            </div>
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.region }}</div>
              <div class="text-sm text-white">{{ deviceStore.selectedDevice.regionInfo || 'N/A' }}</div>
            </div>
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.activation }}</div>
              <div class="text-sm" :class="deviceStore.selectedDevice.activationState === 'Activated' ? 'text-[#30D158]' : 'text-white'">
                {{ deviceStore.selectedDevice.activationState || 'N/A' }}
              </div>
            </div>
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.chipset }}</div>
              <div class="text-sm text-white font-mono">{{ deviceStore.selectedDevice.hardwareModel || 'N/A' }}</div>
            </div>
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.cpuArch }}</div>
              <div class="text-sm text-white font-mono">{{ deviceStore.selectedDevice.cpuArchitecture || 'N/A' }}</div>
            </div>
            <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-1">
              <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-wider">{{ t.apps.details.boardConfig }}</div>
              <div class="text-sm text-white font-mono">{{ deviceStore.selectedDevice.boardConfig || 'N/A' }}</div>
            </div>

          </div>


          <!-- Storage & Battery Row -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <!-- Storage Info -->
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <h4 class="text-sm font-semibold text-white">{{ t.apps.details.storageTitle }}</h4>
                <span class="text-xs text-[#8E8E93]">{{ formatAppSize(deviceStore.selectedDevice.storageFree) }} {{ t.apps.details.storageStatus }} {{ formatAppSize(deviceStore.selectedDevice.storageTotal) }}</span>
              </div>
              <div class="h-2.5 w-full bg-white/5 rounded-full overflow-hidden border border-white/10">
                <div
                  class="h-full bg-[#0A84FF] rounded-full"
                  :style="{ width: `${deviceStore.selectedDevice.storageTotal > 0 ? ((deviceStore.selectedDevice.storageTotal - deviceStore.selectedDevice.storageFree) / deviceStore.selectedDevice.storageTotal) * 100 : 0}%` }"
                ></div>
              </div>
              <div class="flex justify-between text-[10px] text-[#8E8E93] font-medium px-1">
                <div class="flex items-center space-x-1.5">
                  <span class="w-2 h-2 rounded-full bg-[#0A84FF]"></span>
                  <span>Used: {{ formatAppSize(deviceStore.selectedDevice.storageUsed) }}</span>
                </div>
                <div class="flex items-center space-x-1.5">
                  <span class="w-2 h-2 rounded-full bg-white/10"></span>
                  <span>Free: {{ formatAppSize(deviceStore.selectedDevice.storageFree) }}</span>
                </div>
              </div>
            </div>

            <!-- Battery Info -->
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <h4 class="text-sm font-semibold text-white">{{ t.apps.details.batteryTitle }}</h4>
                <div class="flex items-center space-x-2">
                  <span v-if="deviceStore.selectedDevice.batteryCharging" class="text-[10px] font-bold text-[#30D158] uppercase">{{ t.apps.details.batteryStatus }}</span>
                  <span class="text-sm font-bold text-white">{{ deviceStore.selectedDevice.batteryLevel }}%</span>
                </div>
              </div>
              <div class="h-2.5 w-full bg-white/5 rounded-full overflow-hidden border border-white/10 relative">
                <div
                  class="h-full rounded-full transition-all duration-500"
                  :class="deviceStore.selectedDevice.batteryLevel > 20 ? 'bg-[#30D158]' : 'bg-[#FF453A]'"
                  :style="{ width: `${deviceStore.selectedDevice.batteryLevel}%` }"
                ></div>
              </div>
              <div class="flex justify-between text-[10px] text-[#8E8E93] font-medium px-1 pt-1">
                <span>{{ t.apps.details.batteryHealth }}: <span class="text-white font-bold">{{ deviceStore.selectedDevice.batteryHealth || '--' }}%</span></span>
                <span>{{ t.apps.details.batteryCycles }}: <span class="text-white font-bold">{{ deviceStore.selectedDevice.chargeCycles || '--' }}</span></span>
              </div>
            </div>
          </div>

          <!-- Connection Details -->
          <div class="p-4 rounded-2xl bg-white/[0.03] border border-white/5 space-y-4">
            <h4 class="text-sm font-semibold text-white">{{ t.apps.details.connectivityTitle }}</h4>
            <div class="grid grid-cols-2 gap-y-4">
              <div class="space-y-1">
                <div class="text-[10px] font-bold text-[#8E8E93] uppercase">{{ t.apps.details.wifi }}</div>
                <div class="text-xs text-white font-mono">{{ deviceStore.selectedDevice.wifiAddress || 'N/A' }}</div>
              </div>
              <div class="space-y-1 text-right">
                <div class="text-[10px] font-bold text-[#8E8E93] uppercase">{{ t.apps.details.jailbreak }}</div>
                <div class="text-xs font-bold" :class="deviceStore.selectedDevice.isJailbroken ? 'text-[#FF453A]' : 'text-[#30D158]'">
                  {{ deviceStore.selectedDevice.isJailbroken ? 'Yes' : 'No' }}
                </div>
              </div>
              <div class="space-y-1">
                <div class="text-[10px] font-bold text-[#8E8E93] uppercase">{{ t.apps.details.udid }}</div>
                <div class="text-[10px] text-white font-mono break-all">{{ deviceStore.selectedDevice.udid }}</div>
              </div>
            </div>
          </div>
        </div>



        <!-- Modal Footer -->
        <div class="p-4 border-t border-white/10 bg-white/[0.01] flex justify-end">
          <button
            @click="showDeviceDetails = false"
            class="px-6 py-2 rounded-xl bg-white/5 hover:bg-white/10 text-sm font-semibold text-white transition"
          >
            {{ t.apps.details.done }}
          </button>
        </div>

      </div>
    </div>

    <!-- Uninstall Confirmation Modal (REMOVED: Now using GlobalModal) -->
  </div>
</template>



<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useDeviceStore } from '../stores/device'
import { useModalStore } from '../stores/modal'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'
import type { InstalledApp } from '../types'

const deviceStore = useDeviceStore()
const modalStore = useModalStore()
const { t } = useI18n()
const { showToast } = useNotifications()


const searchQuery = ref('')
const showDeviceDetails = ref(false)

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
  deviceStore.checkDevices()
})

function handleRefresh() {
  deviceStore.checkDevices()
}

function handlePair(udid: string) {
  deviceStore.pairDevice(udid).catch((err) => {
    showToast('Pairing Failed', err.message || 'Please unlock your device and trust this computer', 'error')
  })
}


function triggerIPAInstall() {
  deviceStore.installIPA().catch((err) => {
    showToast('Installation Error', err.message || String(err), 'error')
  })
}

function handleUninstall(app: InstalledApp) {
  modalStore.confirm(
    t.value.apps.uninstallTitle,
    t.value.apps.uninstallPrompt.replace('{name}', app.name),
    async () => {
      try {
        showToast('Uninstalling', `Removing ${app.name}...`, 'info')
        await deviceStore.uninstallApp(app.bundleId)
        showToast('Success', `${app.name} has been uninstalled`, 'success')
      } catch (err: any) {
        showToast('Uninstall Failed', err.message || String(err), 'error')
      }
    },
    t.value.apps.uninstall
  )
}



function handleFileDrop(e: DragEvent) {
  const files = e.dataTransfer?.files
  if (files && files.length > 0) {
    const ipaPaths: string[] = []
    for (let i = 0; i < files.length; i++) {
      if (files[i].name.endsWith('.ipa')) {
        const path = (files[i] as any).path
        if (path) ipaPaths.push(path)
      }
    }

    if (ipaPaths.length > 0) {
      if (ipaPaths.length === 1) {
        deviceStore.installIPA(ipaPaths[0])
      } else {
        deviceStore.installMultipleIPAs(ipaPaths)
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
  if (bytes === undefined || bytes === null || bytes < 0) return '--'
  if (bytes === 0) return '0 MB'

  // Use Decimal (1000) for storage - Industry standard for disk sizes and Apple UI
  const unit = 1000
  if (bytes < unit) return `${bytes} B`

  const kb = bytes / unit
  if (kb < unit) return `${kb.toFixed(1)} KB`

  const mb = kb / unit
  if (mb < unit) return `${mb.toFixed(1)} MB`

  const gb = mb / unit
  // No hardcoded rounding to "standard" sizes. Show exactly what the device reports.
  return `${gb.toFixed(2)} GB`
}





</script>
