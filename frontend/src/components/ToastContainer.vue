<template>
  <div class="fixed bottom-5 right-5 z-50 flex flex-col space-y-2.5 pointer-events-none font-sans max-w-sm w-full">
    <transition-group name="toast-slide">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="pointer-events-auto relative overflow-hidden p-4 rounded-[18px] backdrop-blur-[35px] border flex items-start space-x-3.5 transition-all duration-200 shadow-2xl"
        :class="toastStyleClass(t.severity)"
      >
        <!-- Left Vertical Accent Bar -->
        <div
          class="absolute left-0 top-0 bottom-0 w-1.5"
          :class="accentBarClass(t.severity)"
        />

        <!-- Icon Container with Glow Ring -->
        <div
          class="shrink-0 p-2 rounded-xl flex items-center justify-center mt-0.5"
          :class="iconBgClass(t.severity)"
        >
          <!-- Success: CheckCircle -->
          <svg v-if="t.severity === 'success'" class="w-4 h-4 text-[#30D158]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 13l4 4L19 7" />
          </svg>

          <!-- Error: XCircle -->
          <svg v-else-if="t.severity === 'error'" class="w-4 h-4 text-[#FF453A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M6 18L18 6M6 6l12 12" />
          </svg>

          <!-- Warning: ExclamationTriangle -->
          <svg v-else-if="t.severity === 'warning'" class="w-4 h-4 text-[#FFD60A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>

          <!-- Info / Default: InformationCircle -->
          <svg v-else class="w-4 h-4 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>

        <!-- Content Area -->
        <div class="min-w-0 flex-1 space-y-0.5">
          <div class="flex items-center space-x-2">
            <span
              class="px-1.5 py-0.5 rounded-full text-[9px] font-extrabold uppercase tracking-wider"
              :class="badgeClass(t.severity)"
            >
              {{ badgeLabel(t.severity) }}
            </span>
            <h4 class="text-xs font-semibold text-[#FFFFFF] truncate">{{ t.title }}</h4>
          </div>
          <p class="text-[11px] text-[#B8C0CC] leading-snug break-words font-normal pl-0.5">{{ t.message }}</p>
        </div>

        <!-- Close Button -->
        <button
          type="button"
          class="text-[#7D8592] hover:text-white transition duration-150 p-1 -mr-1 rounded-lg hover:bg-white/[0.08]"
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
      return 'bg-[#0F1E14]/90 border-[#30D158]/35 shadow-[0_12px_36px_rgba(48,209,88,0.15)]'
    case 'error':
      return 'bg-[#220F11]/90 border-[#FF453A]/35 shadow-[0_12px_36px_rgba(255,69,58,0.15)]'
    case 'warning':
      return 'bg-[#221B0A]/90 border-[#FFD60A]/35 shadow-[0_12px_36px_rgba(255,214,10,0.15)]'
    default:
      return 'bg-[#0E1A29]/90 border-[#64D2FF]/35 shadow-[0_12px_36px_rgba(100,210,255,0.15)]'
  }
}

function accentBarClass(severity: string) {
  switch (severity) {
    case 'success':
      return 'bg-[#30D158] shadow-[0_0_12px_#30D158]'
    case 'error':
      return 'bg-[#FF453A] shadow-[0_0_12px_#FF453A]'
    case 'warning':
      return 'bg-[#FFD60A] shadow-[0_0_12px_#FFD60A]'
    default:
      return 'bg-[#64D2FF] shadow-[0_0_12px_#64D2FF]'
  }
}

function iconBgClass(severity: string) {
  switch (severity) {
    case 'success':
      return 'bg-[#30D158]/15 border border-[#30D158]/25'
    case 'error':
      return 'bg-[#FF453A]/15 border border-[#FF453A]/25'
    case 'warning':
      return 'bg-[#FFD60A]/15 border border-[#FFD60A]/25'
    default:
      return 'bg-[#64D2FF]/15 border border-[#64D2FF]/25'
  }
}

function badgeClass(severity: string) {
  switch (severity) {
    case 'success':
      return 'bg-[#30D158]/20 text-[#30D158] border border-[#30D158]/30'
    case 'error':
      return 'bg-[#FF453A]/20 text-[#FF453A] border border-[#FF453A]/30'
    case 'warning':
      return 'bg-[#FFD60A]/20 text-[#FFD60A] border border-[#FFD60A]/30'
    default:
      return 'bg-[#64D2FF]/20 text-[#64D2FF] border border-[#64D2FF]/30'
  }
}

function badgeLabel(severity: string) {
  switch (severity) {
    case 'success':
      return 'Éxito'
    case 'error':
      return 'Error'
    case 'warning':
      return 'Aviso'
    default:
      return 'Info'
  }
}
</script>

<style scoped>
.toast-slide-enter-active,
.toast-slide-leave-active {
  transition: all 220ms cubic-bezier(0.16, 1, 0.3, 1);
}

.toast-slide-enter-from {
  opacity: 0;
  transform: translateX(30px) scale(0.95);
}

.toast-slide-leave-to {
  opacity: 0;
  transform: translateX(30px) scale(0.95);
}
</style>
