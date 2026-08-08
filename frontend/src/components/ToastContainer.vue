<template>
  <div class="fixed bottom-5 right-5 z-50 flex flex-col space-y-2 pointer-events-none font-sans max-w-sm w-full">
    <transition-group name="toast-slide">
      <div
        v-for="t in toasts"
        :key="t.id"
        class="pointer-events-auto p-4 rounded-[16px] shadow-[0_12px_40px_rgba(0,0,0,0.35)] flex items-start space-x-3 backdrop-blur-[30px] border transition-all duration-200"
        :class="toastStyleClass(t.severity)"
      >
        <div class="text-base shrink-0 mt-0.5">
          <span v-if="t.severity === 'success'">✅</span>
          <span v-else-if="t.severity === 'error'">❌</span>
          <span v-else-if="t.severity === 'warning'">⚠️</span>
          <span v-else>ℹ️</span>
        </div>
        <div class="min-w-0 flex-1">
          <h4 class="text-xs font-semibold text-[#FFFFFF] truncate">{{ t.title }}</h4>
          <p class="text-[11px] text-[#B8C0CC] mt-0.5 break-words font-normal">{{ t.message }}</p>
        </div>
        <button
          type="button"
          class="text-[#7D8592] hover:text-white text-xs shrink-0 pl-2"
          @click="dismissToast(t.id)"
        >
          ✕
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
      return 'bg-[#171A21]/90 border-[#30D158]/40 text-[#FFFFFF]'
    case 'error':
      return 'bg-[#171A21]/90 border-[#FF453A]/40 text-[#FFFFFF]'
    case 'warning':
      return 'bg-[#171A21]/90 border-[#FFD60A]/40 text-[#FFFFFF]'
    default:
      return 'bg-[#171A21]/90 border-[#64D2FF]/40 text-[#FFFFFF]'
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
