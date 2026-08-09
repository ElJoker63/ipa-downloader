<template>
  <div class="max-w-md mx-auto h-full flex flex-col justify-center animate-slide-up font-sans">
    <div class="glass-card p-8 rounded-[24px] space-y-8 shadow-2xl">
      <!-- Login Header -->
      <div class="flex flex-col items-center text-center space-y-4">
        <div class="w-20 h-20 rounded-3xl bg-gradient-to-tr from-[#0071E3] via-[#0A84FF] to-[#64D2FF] flex items-center justify-center shadow-xl shadow-[#0A84FF]/20 border border-white/20 shrink-0">
          <img src="/logo.png" alt="IPA Downloader" class="w-14 h-14 object-contain invert brightness-0" />
        </div>
        <div>
          <h1 class="text-2xl font-bold text-[#FFFFFF]">{{ t.auth.signIn }}</h1>
          <p class="text-sm text-[#B8C0CC] mt-1 font-normal max-w-[280px]">
            {{ t.auth.signInSubtitle }}
          </p>
        </div>
      </div>

      <!-- Login Form -->
      <form class="space-y-5" @submit.prevent="handleLogin">
        <div class="space-y-1.5">
          <label class="text-xs font-semibold uppercase tracking-wider text-[#7D8592] ml-1">{{ t.auth.emailLabel }}</label>
          <div class="relative">
            <svg class="w-4 h-4 absolute left-3.5 top-3 text-[#7D8592]" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" /></svg>
            <input
              v-model="email"
              type="email"
              required
              placeholder="name@icloud.com"
              class="glass-input w-full pl-11 pr-4 py-2.5 text-sm"
              :disabled="authStore.isLoading"
            />
          </div>
        </div>

        <div class="space-y-1.5">
          <label class="text-xs font-semibold uppercase tracking-wider text-[#7D8592] ml-1">{{ t.auth.passwordLabel }}</label>
          <div class="relative">
            <svg class="w-4 h-4 absolute left-3.5 top-3 text-[#7D8592]" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" /></svg>
            <input
              v-model="password"
              type="password"
              required
              placeholder="••••••••••••"
              class="glass-input w-full pl-11 pr-4 py-2.5 text-sm"
              :disabled="authStore.isLoading"
            />
          </div>
        </div>

        <div class="flex items-center justify-between px-1">
          <label class="flex items-center space-x-2.5 cursor-pointer select-none text-xs text-[#B8C0CC]">
            <input
              v-model="rememberMe"
              type="checkbox"
              class="w-4 h-4 rounded-md bg-white/[0.08] border-white/[0.18] text-[#0A84FF] focus:ring-0"
            />
            <span>{{ t.auth.rememberMe }}</span>
          </label>
        </div>

        <div class="pt-2">
          <button
            type="submit"
            class="btn-primary w-full py-3.5 text-sm font-bold shadow-lg shadow-[#0A84FF]/25"
            :disabled="authStore.isLoading"
          >
            <span v-if="authStore.isLoading" class="flex items-center justify-center space-x-2">
              <svg class="animate-spin h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>{{ t.auth.authenticating }}</span>
            </span>
            <span v-else>{{ t.auth.signInButton }}</span>
          </button>
        </div>
      </form>

      <!-- Footer Info -->
      <div class="pt-2 text-center">
        <p class="text-[10px] text-[#7D8592] px-4">
          {{ t.auth.privacyNotice }}
        </p>
      </div>
    </div>

    <button
      @click="router.push('/')"
      class="mt-6 mx-auto flex items-center space-x-2 text-[#7D8592] hover:text-white transition-colors text-sm font-medium"
    >
      <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path d="M10 19l-7-7m0 0l7-7m-7 7h18" /></svg>
      <span>{{ t.auth.backToHome }}</span>
    </button>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'

const router = useRouter()
const authStore = useAuthStore()
const { t } = useI18n()
const { showToast } = useNotifications()

const email = ref('')
const password = ref('')
const rememberMe = ref(true)

async function handleLogin() {
  try {
    await authStore.login(email.value, password.value, '', rememberMe.value)
    showToast(t.value.auth.signedInSuccess, t.value.auth.welcomeBack.replace('{name}', authStore.account.name || 'User'), 'success')
    router.push('/')
  } catch (err: any) {
    if (err?.message?.includes('2FA') || err?.message?.includes('verification')) {
      authStore.is2FAModalOpen = true
      return
    }
    showToast(t.value.auth.signInFailed, err?.message || t.value.auth.invalidCredentials, 'error')
  }
}
</script>
