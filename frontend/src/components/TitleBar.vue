<template>
  <header
    class="wails-drag h-10 w-full flex items-center justify-between px-4 border-b border-white/[0.08] bg-[#12151B]/85 backdrop-blur-[40px] backdrop-saturate-150 select-none z-50 shrink-0 shadow-sm"
    style="--wails-draggable: drag;"
    @dblclick="toggleMaximize"
  >
    <!-- Left: macOS Traffic Lights & Brand -->
    <div class="flex items-center space-x-3.5">
      <!-- Traffic Lights (macOS Control Buttons) -->
      <div class="wails-no-drag flex items-center space-x-2 group cursor-default py-1" style="--wails-draggable: no-drag;">
        <!-- Close (Red) -->
        <button
          type="button"
          :title="t.common.close || 'Cerrar'"
          class="w-3 h-3 rounded-full bg-[#FF5F56] border border-[#E0443E]/60 flex items-center justify-center transition-all duration-150 hover:brightness-105 active:brightness-90 shadow-[inset_0_1px_0_rgba(255,255,255,0.25)]"
          @click.stop="closeApp"
        >
          <svg class="w-1.5 h-1.5 text-[#4D0000] opacity-0 group-hover:opacity-100 transition-opacity duration-150" viewBox="0 0 24 24" stroke="currentColor" stroke-width="3.5" fill="none">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>

        <!-- Minimize (Yellow) -->
        <button
          type="button"
          :title="t.common.minimize || 'Minimizar'"
          class="w-3 h-3 rounded-full bg-[#FFBD2E] border border-[#DEA123]/60 flex items-center justify-center transition-all duration-150 hover:brightness-105 active:brightness-90 shadow-[inset_0_1px_0_rgba(255,255,255,0.25)]"
          @click.stop="minimize"
        >
          <svg class="w-1.5 h-1.5 text-[#5C3C00] opacity-0 group-hover:opacity-100 transition-opacity duration-150" viewBox="0 0 24 24" stroke="currentColor" stroke-width="4" fill="none">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 12h16" />
          </svg>
        </button>

        <!-- Maximize / Zoom (Green) -->
        <button
          type="button"
          :title="t.common.maximize || 'Maximizar'"
          class="w-3 h-3 rounded-full bg-[#27C93F] border border-[#1AAB29]/60 flex items-center justify-center transition-all duration-150 hover:brightness-105 active:brightness-90 shadow-[inset_0_1px_0_rgba(255,255,255,0.25)]"
          @click.stop="toggleMaximize"
        >
          <svg class="w-1.5 h-1.5 text-[#004D00] opacity-0 group-hover:opacity-100 transition-opacity duration-150" viewBox="0 0 10 10" fill="currentColor">
            <polygon points="1,1 1,5.5 5.5,1" />
            <polygon points="9,9 9,4.5 4.5,9" />
          </svg>
        </button>
      </div>

      <!-- Subtle vertical divider -->
      <div class="h-3 w-px bg-white/[0.1]"></div>

      <!-- Logo & App Name -->
      <div class="flex items-center space-x-2">
        <img src="/logo.png" alt="IPA Downloader" class="w-4 h-4 rounded object-contain opacity-90 shrink-0" />
        <span class="text-xs font-semibold tracking-tight text-white/90 font-sans">{{ t.common.appName }}</span>
        <span class="px-1.5 py-0.2 text-[9px] font-mono font-medium rounded-full bg-white/[0.06] text-[#B8C0CC]/80 border border-white/[0.1]">v1.5.0</span>
      </div>
    </div>

    <!-- Center: Live Connection Status Pill (macOS Capsule) -->
    <div class="flex items-center space-x-2">
      <div
        class="wails-no-drag flex items-center space-x-2 px-3 py-0.5 rounded-full text-[11px] font-medium border backdrop-blur-xl backdrop-saturate-150 transition-all duration-300 ease-liquid cursor-default shadow-specular-soft"
        style="--wails-draggable: no-drag;"
        :class="statusBadgeClass"
      >
        <span class="w-1.5 h-1.5 rounded-full shrink-0" :class="statusDotClass"></span>
        <span class="truncate">{{ translatedStatus }}</span>
        <span v-if="authStore.isLoggedIn" class="text-[#B8C0CC] font-medium text-[10px] pl-1 max-w-[160px] truncate">
          ({{ authStore.account.name || authStore.account.email || t.common.appleId }})
        </span>
      </div>
    </div>

    <!-- Right: Language Pill & Draggable Region -->
    <div class="wails-no-drag flex items-center space-x-2" style="--wails-draggable: no-drag;">
      <!-- Language Quick Switcher -->
      <div class="glass-input flex items-center !rounded-full p-0.5 text-[10px]">
        <button
          type="button"
          class="px-2 py-0.5 rounded-full transition-all duration-200 font-semibold"
          :class="currentLanguage === 'es' ? 'bg-[#0A84FF] text-white shadow-sm' : 'text-[#B8C0CC] hover:text-white'"
          @click="setLanguage('es')"
        >
          ES
        </button>
        <button
          type="button"
          class="px-2 py-0.5 rounded-full transition-all duration-200 font-semibold"
          :class="currentLanguage === 'en' ? 'bg-[#0A84FF] text-white shadow-sm' : 'text-[#B8C0CC] hover:text-white'"
          @click="setLanguage('en')"
        >
          EN
        </button>
      </div>
    </div>
  </header>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useAuthStore } from '../stores/auth'
import { WailsService } from '../services/wails'
import { useI18n } from '../i18n'

const authStore = useAuthStore()
const { t, currentLanguage, setLanguage } = useI18n()

const translatedStatus = computed(() => {
  if (authStore.status === 'Connected') return t.value.common.connected
  if (authStore.status === 'Not Connected') return t.value.common.notConnected
  if (authStore.status === 'Connecting...') return t.value.common.connecting
  return authStore.status
})

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
