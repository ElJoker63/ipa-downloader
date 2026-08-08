<template>
  <div class="fixed bottom-5 right-5 z-50 flex flex-col space-y-2 pointer-events-none font-sans max-w-sm w-full">
    <transition-group name="toast-slide">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="pointer-events-auto p-4 rounded-[16px] shadow-[0_12px_40px_rgba(0,0,0,0.35)] flex items-start space-x-3 backdrop-blur-[30px] border transition-all duration-200"
        :class="toastStyleClass(t.severity)"
      >
        <!-- Modern SVG Icons Instead of Emojis -->
        <div class="shrink-0 mt-0.5">
          <!-- Success: CheckCircle -->
          <svg v-if="t.severity === 'success'" class="w-5 h-5 text-[#30D158]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>

          <!-- Error: XCircle -->
          <svg v-else-if="t.severity === 'error'" class="w-5 h-5 text-[#FF453A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>

          <!-- Warning: ExclamationTriangle -->
          <svg v-else-if="t.severity === 'warning'" class="w-5 h-5 text-[#FFD60A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>

          <!-- Info / Default: InformationCircle -->
          <svg v-else class="w-5 h-5 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>

        <div class="min-w-0 flex-1">
          <h4 class="text-xs font-semibold text-[#FFFFFF] truncate">{{ t.title }}</h4>
          <p class="text-[11px] text-[#B8C0CC] mt-0.5 break-words font-normal">{{ t.message }}</p>
        </div>
        <button
          type="button"
          class="text-[#7D8592] hover:text-white transition duration-150 p-1 -mr-1"
          @click="dismissToast(t.id)"
        >
          <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { useNotifications } from '../composables/useNotifications'

const { toasts, dismissToast } = useNotifications()

function toastStyleClass(severity: string) {
  switch (severity) {
    case 'success':
      return 'bg-[#171A21]/95 border-[#30D158]/40 text-[#FFFFFF]'
    case 'error':
      return 'bg-[#171A21]/95 border-[#FF453A]/40 text-[#FFFFFF]'
    case 'warning':
      return 'bg-[#171A21]/95 border-[#FFD60A]/40 text-[#FFFFFF]'
    default:
      return 'bg-[#171A21]/95 border-[#64D2FF]/40 text-[#FFFFFF]'
  }
}
</script>

<style scoped>
.toast-slide-enter-active,
.toast-slide-leave-active {
  transition: all 200ms cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-slide-enter-from {
  opacity: 0;
  transform: translateY(12px) scale(0.95);
}

.toast-slide-leave-to {
  opacity: 0;
  transform: translateY(12px) scale(0.95);
}
</style>
