<template>
  <div v-if="authStore.is2FAModalOpen" class="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/60 backdrop-blur-md animate-fade-in">
    <div class="glass-panel w-full max-w-md p-6 rounded-2xl shadow-2xl border border-white/10 dark:border-white/10 light:border-black/10 bg-[#131B2E]/95 dark:bg-[#131B2E]/95 light:bg-white/95 text-slate-100 dark:text-slate-100 light:text-slate-900 animate-slide-up">
      <div class="flex items-center space-x-3 mb-4">
        <div class="w-10 h-10 rounded-xl bg-blue-600/20 text-blue-400 flex items-center justify-center">
          <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
        <div>
          <h3 class="text-lg font-bold">Two-Factor Authentication</h3>
          <p class="text-xs text-slate-400 dark:text-slate-400 light:text-slate-500">Enter the verification code sent to your Apple device</p>
        </div>
      </div>

      <div class="my-6">
        <div class="flex justify-between gap-2 mb-4">
          <input
            v-for="(digit, idx) in digits"
            :key="idx"
            :id="`2fa-input-${idx}`"
            v-model="digits[idx]"
            type="text"
            maxlength="1"
            class="w-12 h-14 text-center text-2xl font-bold font-mono rounded-xl bg-slate-900/50 dark:bg-slate-900/50 light:bg-slate-100 border border-white/10 focus:border-blue-500 focus:ring-2 focus:ring-blue-500/20 outline-none transition"
            @input="onInput(idx, $event)"
            @keydown.backspace="onBackspace(idx, $event)"
            @paste="onPaste"
          />
        </div>

        <p v-if="error" class="text-xs text-rose-400 mt-2 text-center">{{ error }}</p>
      </div>

      <div class="flex items-center justify-end space-x-3 mt-6">
        <button
          type="button"
          class="btn-secondary text-sm px-4 py-2"
          :disabled="isSubmitting"
          @click="cancel"
        >
          Cancel
        </button>
        <button
          type="button"
          class="btn-primary text-sm px-5 py-2 flex items-center space-x-2"
          :disabled="isSubmitting || fullCode.length < 6"
          @click="submit"
        >
          <span v-if="isSubmitting" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
          <span>Verify & Connect</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, nextTick, watch } from 'vue'
import { useAuthStore } from '../stores/auth'

const authStore = useAuthStore()
const digits = ref<string[]>(['', '', '', '', '', ''])
const isSubmitting = ref(false)
const error = ref<string | null>(null)

const fullCode = computed(() => digits.value.join(''))

watch(() => authStore.is2FAModalOpen, (open) => {
  if (open) {
    digits.value = ['', '', '', '', '', '']
    error.value = null
    nextTick(() => {
      document.getElementById('2fa-input-0')?.focus()
    })
  }
})

function onInput(idx: number, event: Event) {
  const val = (event.target as HTMLInputElement).value
  digits.value[idx] = val.slice(-1)

  if (val && idx < 5) {
    nextTick(() => {
      document.getElementById(`2fa-input-${idx + 1}`)?.focus()
    })
  }

  if (fullCode.value.length === 6) {
    submit()
  }
}

function onBackspace(idx: number, event: KeyboardEvent) {
  if (!digits.value[idx] && idx > 0) {
    document.getElementById(`2fa-input-${idx - 1}`)?.focus()
  }
}

function onPaste(event: ClipboardEvent) {
  event.preventDefault()
  const pasted = event.clipboardData?.getData('text') || ''
  const clean = pasted.replace(/\D/g, '').slice(0, 6)
  for (let i = 0; i < 6; i++) {
    digits.value[i] = clean[i] || ''
  }
  if (clean.length >= 6) {
    submit()
  } else {
    document.getElementById(`2fa-input-${Math.min(clean.length, 5)}`)?.focus()
  }
}

async function submit() {
  if (fullCode.value.length < 6) return
  isSubmitting.value = true
  error.value = null

  try {
    await authStore.submit2FACode(fullCode.value)
  } catch (err: any) {
    error.value = err?.message || 'Invalid 2FA code. Please try again.'
  } finally {
    isSubmitting.value = false
  }
}

function cancel() {
  authStore.is2FAModalOpen = false
}
</script>
