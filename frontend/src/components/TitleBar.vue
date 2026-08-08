<template>
  <header class="wails-drag h-10 flex items-center justify-between px-3 border-b border-white/5 dark:border-white/5 light:border-black/5 bg-[#0B0F19]/80 dark:bg-[#0B0F19]/80 light:bg-white/80 backdrop-blur-md select-none z-40">
    <!-- Left: App Icon & Brand -->
    <div class="flex items-center space-x-2.5">
      <div class="w-5 h-5 rounded-lg bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center text-white shadow-sm shadow-blue-500/20">
        <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
          <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M15.97 6.37c.62-.75 1.04-1.8 0.92-2.85-.9.04-1.99.6-2.64 1.35-.57.65-.98 1.71-.85 2.73 1 .08 2.03-.51 2.57-1.23z"/>
        </svg>
      </div>
      <span class="text-xs font-bold tracking-wide text-slate-200 dark:text-slate-200 light:text-slate-800">IPATool</span>
      <span class="px-1.5 py-0.2 text-[10px] font-mono rounded bg-white/10 text-slate-300 dark:text-slate-300 light:text-slate-600">v2.0</span>
    </div>

    <!-- Center: Live Connection Status Pill -->
    <div class="flex items-center space-x-2">
      <div
        class="wails-no-drag flex items-center space-x-2 px-2.5 py-0.5 rounded-full text-[11px] font-medium border transition cursor-default"
        :class="statusBadgeClass"
      >
        <span class="w-2 h-2 rounded-full" :class="statusDotClass"></span>
        <span>{{ authStore.status }}</span>
        <span v-if="authStore.isLoggedIn" class="text-slate-400 dark:text-slate-400 light:text-slate-500 font-mono text-[10px] pl-1">
          ({{ authStore.account.email || authStore.account.storeFrontCountry }})
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
  if (authStore.status === 'Connected') {
    return 'bg-emerald-500/10 text-emerald-400 border-emerald-500/20'
  }
  if (authStore.status.includes('Connecting')) {
    return 'bg-amber-500/10 text-amber-300 border-amber-500/20'
  }
  return 'bg-slate-800/60 text-slate-400 border-white/5'
})

const statusDotClass = computed(() => {
  if (authStore.status === 'Connected') {
    return 'bg-emerald-400 shadow-sm shadow-emerald-400/50'
  }
  if (authStore.status.includes('Connecting')) {
    return 'bg-amber-400 animate-pulse'
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
