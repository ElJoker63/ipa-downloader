<template>
  <header
    class="wails-drag relative h-10 w-full flex items-center justify-between px-4 border-b border-white/[0.08] bg-[#12151B]/85 backdrop-blur-[40px] backdrop-saturate-150 select-none z-50 shrink-0 shadow-sm"
    style="--wails-draggable: drag;"
    @dblclick="toggleMaximize"
  >
    <!-- Left: macOS Traffic Lights Only -->
    <div class="flex items-center">
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
    </div>

    <!-- Center: Application Name Only (Centered) -->
    <div class="pointer-events-none absolute inset-x-0 flex items-center justify-center">
      <span class="text-xs font-semibold tracking-tight text-white/90 font-sans">
        {{ t.common.appName }}
      </span>
    </div>

    <!-- Right: Spacer for balance and dragging -->
    <div class="w-14"></div>
  </header>
</template>

<script setup lang="ts">
import { WailsService } from '../services/wails'
import { useI18n } from '../i18n'

const { t } = useI18n()

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
