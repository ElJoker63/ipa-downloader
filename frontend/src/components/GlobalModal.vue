<template>
  <Transition name="fade-backdrop">
    <div v-if="modalStore.isOpen" class="fixed inset-0 bg-black/60 backdrop-blur-md z-[9999] flex items-center justify-center p-6">
      <Transition name="modal-scale">
        <div class="w-full max-w-sm bg-[#1C1C1E] border border-white/10 rounded-3xl shadow-2xl overflow-hidden p-6 space-y-6">
          <div class="flex flex-col items-center text-center space-y-4">
            <!-- Icon based on type -->
            <div
              class="w-16 h-16 rounded-2xl flex items-center justify-center"
              :class="iconBgClass"
            >
              <svg v-if="modalStore.options.type === 'confirm'" class="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <svg v-else-if="modalStore.options.type === 'error'" class="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <svg v-else-if="modalStore.options.type === 'success'" class="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <svg v-else class="w-10 h-10 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>

            <div class="space-y-1">
              <h3 class="text-lg font-bold text-white">{{ modalStore.options.title }}</h3>
              <p class="text-sm text-[#8E8E93] whitespace-pre-wrap leading-relaxed">
                {{ modalStore.options.message }}
              </p>
            </div>
          </div>

          <div class="flex flex-col space-y-2">
            <button
              @click="modalStore.handleConfirm"
              :disabled="modalStore.loading"
              class="w-full py-3 rounded-2xl text-white font-bold text-sm transition shadow-lg flex items-center justify-center space-x-2 disabled:opacity-50"
              :class="confirmBtnClass"
            >
              <svg v-if="modalStore.loading" class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>{{ modalStore.options.confirmText }}</span>
            </button>
            <button
              v-if="modalStore.options.cancelText"
              @click="modalStore.handleCancel"
              :disabled="modalStore.loading"
              class="w-full py-3 rounded-2xl bg-white/5 hover:bg-white/10 text-white font-semibold text-sm transition disabled:opacity-50"
            >
              {{ modalStore.options.cancelText }}
            </button>
          </div>
        </div>
      </Transition>
    </div>
  </Transition>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useModalStore } from '../stores/modal'

const modalStore = useModalStore()

const iconBgClass = computed(() => {
  switch (modalStore.options.type) {
    case 'error': return 'bg-red-500'
    case 'success': return 'bg-[#30D158]'
    case 'confirm': return 'bg-[#FF9F0A]'
    default: return 'bg-[#0A84FF]'
  }
})

const confirmBtnClass = computed(() => {
  switch (modalStore.options.type) {
    case 'error': return 'bg-red-500 hover:bg-red-600 shadow-red-500/20'
    case 'success': return 'bg-[#30D158] hover:bg-[#28B84B] shadow-[#30D158]/20'
    case 'confirm': return 'bg-[#FF9F0A] hover:bg-[#E68E00] shadow-[#FF9F0A]/20'
    default: return 'bg-[#0A84FF] hover:bg-[#0071E3] shadow-[#0A84FF]/20'
  }
})
</script>

<style scoped>
.fade-backdrop-enter-active, .fade-backdrop-leave-active {
  transition: opacity 0.3s ease;
}
.fade-backdrop-enter-from, .fade-backdrop-leave-to {
  opacity: 0;
}

.modal-scale-enter-active, .modal-scale-leave-active {
  transition: transform 0.3s cubic-bezier(0.34, 1.56, 0.64, 1), opacity 0.2s ease;
}
.modal-scale-enter-from, .modal-scale-leave-to {
  transform: scale(0.9) translateY(10px);
  opacity: 0;
}
</style>
