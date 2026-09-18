<template>
  <div
    v-if="authStore.is2FAModalOpen"
    class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-xl transition-all duration-200"
  >
    <div class="glass-card w-full max-w-md !rounded-[26px] p-6 space-y-5 animate-modal font-sans">
      <div class="text-center space-y-2">
        <div class="w-12 h-12 rounded-[14px] bg-[#0A84FF]/15 border border-[#0A84FF]/30 flex items-center justify-center mx-auto text-[#0A84FF] shadow-sm">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
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

        <div class="rounded-xl bg-white/[0.04] border border-white/10 p-3 space-y-1 text-left">
          <p class="text-[11px] font-semibold text-[#0A84FF] flex items-center gap-1.5">
            <svg class="w-3.5 h-3.5 flex-shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
            {{ t.twoFactor.manualCodeTipTitle }}
          </p>
          <p class="text-[11px] text-[#86868B] leading-relaxed">
            {{ t.twoFactor.manualCodeTip }}
          </p>
        </div>

        <div class="flex items-center space-x-3 pt-2">
          <button
            type="button"
            class="btn-secondary text-xs px-4 py-2.5 flex-1"
            @click="handleCancel"
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
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'

const router = useRouter()
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

function handleCancel() {
  authStore.cancel2FA()
}

async function submit2FA() {
  const cleanCode = code.value.trim().replace(/\s+/g, '')
  if (!cleanCode || cleanCode.length < 6) {
    showToast('Invalid Code', t.value.twoFactor.invalidCode, 'error')
    return
  }

  try {
    await authStore.submit2FACode(cleanCode)
    showToast('Verified', t.value.twoFactor.verifiedToast, 'success')
    router.push('/')
  } catch (err: any) {
    code.value = ''
    nextTick(() => {
      inputRef.value?.focus()
    })
    showToast('Verification Failed', err?.message || 'Invalid 2FA code', 'error')
  }
}
</script>
