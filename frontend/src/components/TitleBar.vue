<template>
  <header class="wails-drag h-11 flex items-center justify-between px-4 border-b border-white/[0.08] bg-[#0F1115]/80 backdrop-blur-[30px] select-none z-50">
    <!-- Left: App Icon & Brand -->
    <div class="flex items-center space-x-3">
      <img src="/logo.png" alt="IPA Downloader" class="w-6 h-6 rounded-lg object-contain shadow-sm shrink-0" />
      <span class="text-[13px] font-bold tracking-tight text-white font-sans">IPA Downloader</span>
      <span class="px-2 py-0.5 text-[10px] font-mono font-medium rounded-full bg-white/[0.08] text-[#B8C0CC] border border-white/[0.12]">v2.0</span>
    </div>

    <!-- Center: Live Connection Status Pill (macOS Capsule) -->
    <div class="flex items-center space-x-2">
      <div
        class="wails-no-drag flex items-center space-x-2 px-3 py-1 rounded-full text-xs font-medium border backdrop-blur-md transition cursor-default shadow-sm"
        :class="statusBadgeClass"
      >
        <span class="w-2 h-2 rounded-full" :class="statusDotClass"></span>
        <span>{{ authStore.status }}</span>
        <span v-if="authStore.isLoggedIn" class="text-[#B8C0CC] font-medium text-[11px] pl-1">
          ({{ authStore.account.name || 'Apple ID' }})
        </span>
      </div>
    </div>

    <!-- Right: Window Controls (Minimalist & Smooth) -->
    <div class="wails-no-drag flex items-center space-x-1.5">
      <button
        type="button"
        title="Minimize"
        class="w-7 h-7 flex items-center justify-center rounded-lg text-[#B8C0CC] hover:text-white hover:bg-white/[0.08] transition duration-150"
        @click="minimize"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
        </svg>
      </button>
      <button
        type="button"
        title="Toggle Maximize"
        class="w-7 h-7 flex items-center justify-center rounded-lg text-[#B8C0CC] hover:text-white hover:bg-white/[0.08] transition duration-150"
        @click="toggleMaximize"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
        </svg>
      </button>
      <button
        type="button"
        title="Close"
        class="w-7 h-7 flex items-center justify-center rounded-lg text-[#B8C0CC] hover:text-white hover:bg-[#FF453A] transition duration-150"
        @click="closeApp"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { WailsService } from '../services/wails'

const authStore = useAuthStore()

const statusBadgeClass = computed(() => {
  if (authStore.isLoggedIn) {
    return 'bg-[#30D158]/15 border-[#30D158]/30 text-[#30D158]'
  }
  if (authStore.isLoading) {
    return 'bg-[#FFD60A]/15 border-[#FFD60A]/30 text-[#FFD60A] animate-pulse'
  }
  return 'bg-white/[0.08] border-white/[0.18] text-[#B8C0CC]'
})

const statusDotClass = computed(() => {
  if (authStore.isLoggedIn) {
    return 'bg-[#30D158] shadow-[0_0_8px_rgba(48,209,88,0.6)] animate-pulse-subtle'
  }
  if (authStore.isLoading) {
    return 'bg-[#FFD60A] animate-ping'
  }
  return 'bg-[#7D8592]'
})

function minimize() {
  WailsService.minimizeWindow()
}

function toggleMaximize() {
  WailsService.toggleMaximizeWindow()
}

function closeApp() {
  WailsService.closeWindow()
}
</script>
