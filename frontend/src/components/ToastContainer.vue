<template>
  <div class="fixed bottom-6 right-6 z-50 flex flex-col space-y-3 pointer-events-none max-w-sm w-full">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      class="glass-panel p-4 rounded-xl shadow-2xl pointer-events-auto flex items-start space-x-3 transition-all duration-300 transform translate-y-0 opacity-100 animate-slide-up border border-white/10 dark:border-white/10 light:border-black/10 bg-[#131B2E]/95 dark:bg-[#131B2E]/95 light:bg-white/95 text-slate-100 dark:text-slate-100 light:text-slate-900"
    >
      <div class="shrink-0 mt-0.5">
        <svg v-if="toast.severity === 'success'" class="w-5 h-5 text-emerald-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <svg v-else-if="toast.severity === 'error'" class="w-5 h-5 text-rose-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <svg v-else-if="toast.severity === 'warning'" class="w-5 h-5 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
        </svg>
        <svg v-else class="w-5 h-5 text-blue-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
      </div>

      <div class="flex-1 min-w-0">
        <h4 class="text-sm font-semibold">{{ toast.title }}</h4>
        <p class="text-xs text-slate-400 dark:text-slate-400 light:text-slate-600 mt-0.5 line-clamp-2">{{ toast.message }}</p>
      </div>

      <button
        type="button"
        class="text-slate-400 hover:text-white dark:hover:text-white light:hover:text-black transition"
        @click="dismissToast(toast.id)"
      >
        <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
        </svg>
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useNotifications } from '../composables/useNotifications'

const { toasts, dismissToast } = useNotifications()
</script>
