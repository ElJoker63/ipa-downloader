<template>
  <div
    v-if="authStore.is2FAModalOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md transition-all duration-200"
  >
    <div class="glass-card w-full max-w-md rounded-[22px] border border-white/[0.18] p-6 shadow-[0_16px_48px_rgba(0,0,0,0.45)] space-y-5 animate-modal font-sans">
      <div class="text-center space-y-1">
        <div class="w-12 h-12 rounded-[14px] bg-[#0A84FF]/15 border border-[#0A84FF]/30 flex items-center justify-center mx-auto text-xl">
          🔐
        </div>
        <h2 class="text-lg font-bold text-[#FFFFFF]">{{ t.twoFactor.title }}</h2>
        <p class="text-xs text-[#B8C0CC]">
          {{ t.twoFactor.subtitle }}
        </p>
      </div>

      <form class="space-y-4" @submit.prevent="submit2FA">
        <div class="space-y-1.5">
          <label class="text-xs font-medium text-[#B8C0CC]">{{ t.twoFactor.label }}</label>
          <input
            ref="inputRef"
            v-model="code"
            type="text"
            required
            maxlength="6"
            placeholder="123456"
            class="glass-input w-full px-4 py-3 text-center text-xl font-mono tracking-widest"
          />
        </div>

        <div class="flex items-center space-x-3 pt-2">
          <button
            type="button"
            class="btn-secondary text-xs px-4 py-2.5 flex-1"
            @click="authStore.is2FAModalOpen = false"
          >
            {{ t.common.cancel }}
          </button>
          <button
            type="submit"
            class="btn-primary text-xs px-5 py-2.5 flex-1"
            :disabled="authStore.isLoading"
          >
            <span v-if="authStore.isLoading">{{ t.twoFactor.verifying }}</span>
            <span v-else>{{ t.twoFactor.verifyButton }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, nextTick, watch } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'

const authStore = useAuthStore()
const { t } = useI18n()
const { showToast } = useNotifications()

const code = ref('')
const inputRef = ref<HTMLInputElement | null>(null)

watch(
  () => authStore.is2FAModalOpen,
  (open) => {
    if (open) {
      code.value = ''
      nextTick(() => {
        inputRef.value?.focus()
      })
    }
  }
)

async function submit2FA() {
  if (!code.value || code.value.length < 6) {
    showToast('Invalid Code', t.value.twoFactor.invalidCode, 'error')
    return
  }

  try {
    await authStore.submit2FACode(code.value)
    showToast('Verified', t.value.twoFactor.verifiedToast, 'success')
  } catch (err: any) {
    showToast('Verification Failed', err?.message || 'Invalid 2FA code', 'error')
  }
}
</script>
