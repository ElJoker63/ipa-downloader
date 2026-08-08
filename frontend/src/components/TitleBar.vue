<template>
  <header class="wails-drag h-10 flex items-center justify-between px-3 border-b border-white/5 dark:border-white/5 light:border-black/5 bg-[#080C16]/85 dark:bg-[#080C16]/85 light:bg-white/85 backdrop-blur-xl select-none z-40">
    <!-- Left: App Icon & Brand -->
    <div class="flex items-center space-x-2.5">
      <img src="/logo.png" alt="IPA Downloader" class="w-6 h-6 rounded-lg object-contain shadow-md shrink-0" />
      <span class="text-xs font-extrabold tracking-tight text-white dark:text-white light:text-slate-900 font-sans">IPA Downloader</span>
      <span class="px-1.5 py-0.2 text-[9px] font-mono font-bold rounded bg-blue-500/20 text-blue-400 border border-blue-500/30">v2.0</span>
    </div>

    <!-- Center: Live Connection Status Pill -->
    <div class="flex items-center space-x-2">
      <div
        class="wails-no-drag flex items-center space-x-2 px-3 py-0.5 rounded-full text-[11px] font-medium border transition cursor-default shadow-sm"
        :class="statusBadgeClass"
      >
        <span class="w-2 h-2 rounded-full" :class="statusDotClass"></span>
        <span>{{ authStore.status }}</span>
        <span v-if="authStore.isLoggedIn" class="text-slate-400 dark:text-slate-400 light:text-slate-500 font-medium text-[10px] pl-1">
          ({{ authStore.account.name || 'Apple ID' }})
        </span>
      </div>
    </div>

    <!-- Right: Window Controls (Frameless Wails Actions) -->
    <div class="wails-no-drag flex items-center space-x-1">
      <button
        type="button"
        title="Minimize"
        class="w-7 h-7 flex items-center justify-center rounded-lg text-slate-400 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/10 dark:hover:bg-white/10 light:hover:bg-black/5 transition"
        @click="minimize"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 12H4" />
        </svg>
      </button>
      <button
        type="button"
        title="Toggle Maximize"
        class="w-7 h-7 flex items-center justify-center rounded-lg text-slate-400 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/10 dark:hover:bg-white/10 light:hover:bg-black/5 transition"
        @click="toggleMaximize"
      >
        <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
        </svg>
      </button>
      <button
        type="button"
        title="Close"
        class="w-7 h-7 flex items-center justify-center rounded-lg text-slate-400 hover:text-white hover:bg-rose-600 transition"
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
    return 'bg-emerald-500/15 border-emerald-500/30 text-emerald-400'
  }
  if (authStore.isLoading) {
    return 'bg-amber-500/15 border-amber-500/30 text-amber-400 animate-pulse'
  }
  return 'bg-slate-500/15 border-white/10 text-slate-400'
})

const statusDotClass = computed(() => {
  if (authStore.isLoggedIn) {
    return 'bg-emerald-400 shadow-sm shadow-emerald-400/50 animate-pulse-subtle'
  }
  if (authStore.isLoading) {
    return 'bg-amber-400 animate-ping'
  }
  return 'bg-slate-500'
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
