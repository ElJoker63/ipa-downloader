<template>
  <div class="max-w-5xl mx-auto space-y-8 animate-slide-up font-sans">
    <!-- Header Section -->
    <div class="flex items-center justify-between">
      <div>
        <div class="flex items-center space-x-2 text-[#0A84FF] text-xs font-semibold uppercase tracking-wider mb-1">
          <svg class="w-3.5 h-3.5" viewBox="0 0 24 24" fill="currentColor">
            <path d="M18.71 19.5c-.83 1.24-1.71 2.45-3.05 2.47-1.34.03-1.77-.79-3.29-.79-1.53 0-2 .77-3.27.82-1.31.05-2.3-1.32-3.14-2.53C4.25 17 2.94 12.45 4.7 9.39c.87-1.52 2.43-2.48 4.12-2.51 1.28-.02 2.5.87 3.29.87.78 0 2.26-1.07 3.81-.91.65.03 2.47.26 3.64 1.98-.09.06-2.17 1.28-2.15 3.81.03 3.02 2.65 4.03 2.68 4.04-.03.07-.42 1.44-1.38 2.83M15.97 6.37c.62-.75 1.04-1.8 0.92-2.85-.9.04-1.99.6-2.64 1.35-.57.65-.98 1.71-.85 2.73 1 .08 2.03-.51 2.57-1.23z"/>
          </svg>
          <span>{{ t.home.badge }}</span>
        </div>
        <h1 class="text-3xl font-bold tracking-tight text-[#FFFFFF]">
          {{ t.home.title }}
        </h1>
        <p class="text-sm text-[#B8C0CC] mt-1 font-normal">
          {{ t.home.subtitle }}
        </p>
      </div>

      <!-- Live Connection Status Pill (macOS Capsule) -->
      <div class="flex items-center space-x-2 px-4 py-2 rounded-[14px] glass-panel">
        <span class="w-2.5 h-2.5 rounded-full" :class="statusDotClass"></span>
        <span class="text-xs font-semibold text-[#FFFFFF]">{{ authStore.isLoggedIn ? t.common.connected : t.common.notConnected }}</span>
      </div>
    </div>

    <!-- Quick Stats Cards Row (8px Grid & Glass Design) -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="glass-card p-5 rounded-[18px]">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-[#B8C0CC]">{{ t.home.session }}</span>
          <svg class="w-4 h-4 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
        <div class="text-lg font-bold mt-2 text-[#FFFFFF] truncate">{{ authStore.isLoggedIn ? (authStore.account.name || t.common.connected) : t.home.guestMode }}</div>
        <div class="text-[11px] text-[#30D158] font-medium mt-0.5">{{ authStore.isLoggedIn ? t.home.active : t.home.signInPrompt }}</div>
      </div>

      <div class="glass-card p-5 rounded-[18px]">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-[#B8C0CC]">{{ t.home.storefront }}</span>
          <svg class="w-4 h-4 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3.055 11H5a2 2 0 012 2v1a2 2 0 002 2 2 2 0 012 2v2.945M8 3.935V5.5A2.5 2.5 0 0010.5 8h.5a2 2 0 012 2 2 2 0 104 0 2 2 0 012-2h1.064M15 20.488V18a2 2 0 012-2h3.064M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
        </div>
        <div class="text-lg font-bold mt-2 text-[#FFFFFF]">{{ authStore.account.storeFrontCountry || 'US' }}</div>
        <div class="text-[11px] text-[#7D8592] font-mono mt-0.5">{{ t.home.region }}: {{ authStore.account.storeFront || '143441-1,29' }}</div>
      </div>

      <div class="glass-card p-5 rounded-[18px]">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-[#B8C0CC]">{{ t.home.activeQueue }}</span>
          <svg class="w-4 h-4 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
        </div>
        <div class="text-lg font-bold mt-2 text-[#FFFFFF]">{{ downloadsStore.activeCount }}</div>
        <div class="text-[11px] text-[#0A84FF] font-medium mt-0.5">{{ downloadsStore.totalSpeedFormatted }}</div>
      </div>

      <div class="glass-card p-5 rounded-[18px]">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-[#B8C0CC]">{{ t.home.savedFavorites }}</span>
          <svg class="w-4 h-4 text-[#FFD60A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
          </svg>
        </div>
        <div class="text-lg font-bold mt-2 text-[#FFFFFF]">{{ favoritesStore.favorites.length }}</div>
        <div class="text-[11px] text-[#FFD60A] font-medium mt-0.5">{{ t.home.bookmarkedApps }}</div>
      </div>
    </div>

    <!-- Main Grid: Account Status / Login Form -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-6 items-start">
      <!-- Left Column: Authenticated Account or Sign-in Form -->
      <div class="md:col-span-7 space-y-6">
        <!-- Logged In Card with Dynamic Initials Avatar -->
        <div v-if="authStore.isLoggedIn" class="glass-card p-6 rounded-[18px] space-y-6">
          <div class="flex items-center space-x-4">
            <div class="w-16 h-16 rounded-[18px] bg-gradient-to-tr from-[#0071E3] via-[#0A84FF] to-[#64D2FF] flex items-center justify-center text-2xl font-bold text-white shadow-xl shadow-[#0A84FF]/25 border border-white/20 shrink-0 select-none">
              {{ userInitials }}
            </div>
            <div class="min-w-0 flex-1">
              <h2 class="text-xl font-bold text-[#FFFFFF] truncate">{{ authStore.account.name || 'Apple ID User' }}</h2>
              <div class="flex items-center gap-2 mt-1">
                <span class="px-2.5 py-0.5 text-[11px] font-semibold rounded-full bg-[#30D158]/15 text-[#30D158] border border-[#30D158]/30">
                  {{ t.home.storefront }}: {{ authStore.account.storeFrontCountry || 'US' }}
                </span>
                <span v-if="authStore.account.directoryServicesId" class="text-[11px] font-mono text-[#7D8592]">
                  DSID: {{ authStore.account.directoryServicesId }}
                </span>
              </div>
            </div>
          </div>

          <div class="p-4 rounded-[14px] bg-white/[0.04] border border-white/[0.08] space-y-2.5 text-xs text-[#B8C0CC]">
            <div class="flex justify-between">
              <span class="text-[#7D8592]">{{ t.home.sessionSecurity }}</span>
              <span class="text-[#30D158] font-medium">{{ t.home.encryptedKeychain }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-[#7D8592]">{{ t.home.fairplaySigning }}</span>
              <span class="text-[#0A84FF] font-medium">{{ t.home.active }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-[#7D8592]">{{ t.home.autoLicensing }}</span>
              <span class="text-[#64D2FF] font-medium">{{ t.home.enabled }}</span>
            </div>
          </div>

          <div class="flex items-center justify-end space-x-3 pt-2">
            <router-link to="/search" class="btn-primary text-sm px-5 py-2.5 flex items-center space-x-2">
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
              <span>{{ t.home.searchAppStore }}</span>
            </router-link>
            <button
              type="button"
              class="btn-secondary text-sm px-4 py-2.5 text-[#FF453A] hover:text-white hover:bg-[#FF453A] hover:border-[#FF453A]"
              :disabled="authStore.isLoading"
              @click="logout"
            >
              {{ t.home.signOut }}
            </button>
          </div>
        </div>

        <!-- Login Form Glass Card -->
        <div v-else class="glass-card p-6 rounded-[18px] space-y-6">
          <div class="flex items-center space-x-3.5">
            <img src="/logo.png" alt="IPA Downloader" class="w-10 h-10 object-contain rounded-xl shrink-0" />
            <div>
              <h2 class="text-lg font-bold text-[#FFFFFF]">{{ t.home.signInTitle }}</h2>
              <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
                {{ t.home.signInSubtitle }}
              </p>
            </div>
          </div>

          <form class="space-y-4" @submit.prevent="handleLogin">
            <div class="space-y-1.5">
              <label class="text-xs font-medium text-[#B8C0CC]">{{ t.home.emailLabel }}</label>
              <input
                v-model="email"
                type="email"
                required
                placeholder="name@icloud.com"
                class="glass-input w-full px-3.5 py-2.5 text-sm"
              />
            </div>

            <div class="space-y-1.5">
              <label class="text-xs font-medium text-[#B8C0CC]">{{ t.home.passwordLabel }}</label>
              <input
                v-model="password"
                type="password"
                required
                placeholder="••••••••••••"
                class="glass-input w-full px-3.5 py-2.5 text-sm"
              />
            </div>

            <div class="flex items-center justify-between pt-1">
              <label class="flex items-center space-x-2 cursor-pointer select-none text-xs text-[#B8C0CC]">
                <input
                  v-model="rememberMe"
                  type="checkbox"
                  class="rounded bg-white/[0.08] border-white/[0.18] text-[#0A84FF] focus:ring-0"
                />
                <span>{{ t.home.rememberMe }}</span>
              </label>
            </div>

            <div class="pt-2">
              <button
                type="submit"
                class="btn-primary w-full py-2.5 text-sm"
                :disabled="authStore.isLoading"
              >
                <span v-if="authStore.isLoading" class="flex items-center space-x-2">
                  <svg class="animate-spin h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8H4z"></path>
                  </svg>
                  <span>{{ t.home.authenticating }}</span>
                </span>
                <span v-else>{{ t.home.signInButton }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Right Column: Features & Quick Tips -->
      <div class="md:col-span-5 space-y-4">
        <div class="glass-card p-6 rounded-[18px] space-y-4">
          <h3 class="text-sm font-semibold uppercase tracking-wider text-[#B8C0CC]">{{ t.home.featuresTitle }}</h3>
          
          <div class="space-y-3.5 text-xs text-[#B8C0CC]">
            <div class="flex items-start space-x-3">
              <svg class="w-4 h-4 text-[#0A84FF] shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 10V3L4 14h7v7l9-11h-7z" />
              </svg>
              <div>
                <span class="font-semibold text-[#FFFFFF]">{{ t.home.feat1Title }}</span>
                <p class="text-[#7D8592] mt-0.5">{{ t.home.feat1Desc }}</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <svg class="w-4 h-4 text-[#64D2FF] shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
              <div>
                <span class="font-semibold text-[#FFFFFF]">{{ t.home.feat2Title }}</span>
                <p class="text-[#7D8592] mt-0.5">{{ t.home.feat2Desc }}</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <svg class="w-4 h-4 text-[#30D158] shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
              </svg>
              <div>
                <span class="font-semibold text-[#FFFFFF]">{{ t.home.feat3Title }}</span>
                <p class="text-[#7D8592] mt-0.5">{{ t.home.feat3Desc }}</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <svg class="w-4 h-4 text-[#FFD60A] shrink-0 mt-0.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
              </svg>
              <div>
                <span class="font-semibold text-[#FFFFFF]">{{ t.home.feat4Title }}</span>
                <p class="text-[#7D8592] mt-0.5">{{ t.home.feat4Desc }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useAuthStore } from '../stores/auth'
import { useDownloadsStore } from '../stores/downloads'
import { useFavoritesStore } from '../stores/favorites'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'

const authStore = useAuthStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const { t } = useI18n()
const { showToast } = useNotifications()

const email = ref('')
const password = ref('')
const rememberMe = ref(true)

onMounted(async () => {
  await authStore.checkAccount()
  await downloadsStore.fetchDownloads()
  await favoritesStore.fetchFavorites()
})

const statusDotClass = computed(() => {
  if (authStore.isLoggedIn) return 'bg-[#30D158] shadow-[0_0_8px_rgba(48,209,88,0.8)] animate-pulse-subtle'
  if (authStore.isLoading) return 'bg-[#FFD60A] animate-ping'
  return 'bg-[#7D8592]'
})

const userInitials = computed(() => {
  const name = authStore.account?.name || authStore.account?.email || ''
  if (!name.trim()) return ''
  const parts = name.trim().split(/\s+/)
  if (parts.length === 1) {
    return parts[0].substring(0, 2).toUpperCase()
  }
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase()
})

async function handleLogin() {
  try {
    await authStore.login(email.value, password.value, '', rememberMe.value)
    showToast('Signed In', `Welcome back ${authStore.account.name || 'User'}!`, 'success')
  } catch (err: any) {
    if (err?.message?.includes('2FA') || err?.message?.includes('verification')) {
      authStore.is2FAModalOpen = true
      return
    }
    showToast('Sign In Failed', err?.message || 'Invalid credentials', 'error')
  }
}

async function logout() {
  await authStore.logout()
  showToast('Signed Out', 'Your Apple Storefront session has been revoked', 'info')
}
</script>
