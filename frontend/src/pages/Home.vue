<template>
  <div class="max-w-5xl mx-auto space-y-8 animate-slide-up">
    <!-- Header Section -->
    <div class="flex items-center justify-between">
      <div>
        <div class="flex items-center space-x-2 text-blue-500 text-xs font-bold uppercase tracking-wider mb-1">
          <span> App Store Direct Integration</span>
        </div>
        <h1 class="text-3xl font-extrabold tracking-tight text-white dark:text-white light:text-slate-900">
          IPA Downloader
        </h1>
        <p class="text-sm text-slate-400 dark:text-slate-400 light:text-slate-600 mt-1">
          Search, inspect version history, and download signed iOS, iPadOS, and tvOS packages directly from Apple.
        </p>
      </div>

      <!-- Live Connection Status Pill -->
      <div class="flex items-center space-x-2 px-4 py-2 rounded-2xl glass-panel border border-white/10">
        <span class="w-2.5 h-2.5 rounded-full" :class="statusDotClass"></span>
        <span class="text-xs font-bold">{{ authStore.status }}</span>
      </div>
    </div>

    <!-- Quick Stats Cards Row -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
      <div class="glass-card-interactive p-4 rounded-2xl border border-white/10">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-400">Connection</span>
          <span class="text-base">🔐</span>
        </div>
        <div class="text-lg font-bold mt-1 text-white truncate">{{ authStore.isLoggedIn ? (authStore.account.name || 'Connected') : 'Guest Mode' }}</div>
        <div class="text-[11px] text-emerald-400 font-medium">{{ authStore.isLoggedIn ? '● Active Session' : '○ Sign-in for downloads' }}</div>
      </div>

      <div class="glass-card-interactive p-4 rounded-2xl border border-white/10">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-400">Storefront</span>
          <span class="text-base">🌐</span>
        </div>
        <div class="text-lg font-bold mt-1 text-white">{{ authStore.account.storeFrontCountry || 'US' }}</div>
        <div class="text-[11px] text-slate-400 font-mono">Region: {{ authStore.account.storeFront || '143441-1,29' }}</div>
      </div>

      <div class="glass-card-interactive p-4 rounded-2xl border border-white/10">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-400">Active Queue</span>
          <span class="text-base">⬇</span>
        </div>
        <div class="text-lg font-bold mt-1 text-white">{{ downloadsStore.activeCount }}</div>
        <div class="text-[11px] text-blue-400">{{ downloadsStore.totalSpeedFormatted }}</div>
      </div>

      <div class="glass-card-interactive p-4 rounded-2xl border border-white/10">
        <div class="flex items-center justify-between">
          <span class="text-xs font-medium text-slate-400">Saved Favorites</span>
          <span class="text-base">⭐</span>
        </div>
        <div class="text-lg font-bold mt-1 text-white">{{ favoritesStore.favorites.length }}</div>
        <div class="text-[11px] text-amber-400">Bookmarked apps</div>
      </div>
    </div>

    <!-- Main Grid: Account Status / Login Form -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-6 items-start">
      <!-- Left Column: Authenticated Account or Sign-in Form -->
      <div class="md:col-span-7 space-y-6">
        <!-- Logged In Card -->
        <div v-if="authStore.isLoggedIn" class="glass-panel p-6 rounded-2xl border border-white/10 space-y-6">
          <div class="flex items-center space-x-4">
            <img src="/icon.png" alt="IPA Downloader" class="w-16 h-16 rounded-2xl object-contain bg-slate-900/80 border border-white/10 shadow-xl shadow-blue-500/25 shrink-0 p-1" />
            <div class="min-w-0 flex-1">
              <h2 class="text-xl font-extrabold text-white truncate">{{ authStore.account.name || 'Apple ID User' }}</h2>
              <div class="flex items-center gap-2 mt-1">
                <span class="px-2 py-0.5 text-[10px] font-bold rounded-full bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">
                  Storefront: {{ authStore.account.storeFrontCountry || 'US' }}
                </span>
                <span v-if="authStore.account.directoryServicesId" class="text-[10px] font-mono text-slate-400">
                  DSID: {{ authStore.account.directoryServicesId }}
                </span>
              </div>
            </div>
          </div>

          <div class="p-4 rounded-xl bg-white/5 dark:bg-white/5 light:bg-black/5 border border-white/5 space-y-2.5 text-xs text-slate-300 dark:text-slate-300 light:text-slate-700">
            <div class="flex justify-between">
              <span class="text-slate-400">Session Security:</span>
              <span class="text-emerald-400 font-semibold">● Encrypted Keychain</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">FairPlay SINF Replication:</span>
              <span class="text-blue-400 font-semibold">Enabled</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">Automatic Free Licensing:</span>
              <span class="text-indigo-400 font-semibold">Active</span>
            </div>
          </div>

          <div class="flex items-center justify-end space-x-3 pt-2">
            <router-link to="/search" class="btn-primary text-sm px-5 py-2.5 flex items-center space-x-2">
              <span>🔍 Search App Store</span>
            </router-link>
            <button
              type="button"
              class="btn-secondary text-sm px-4 py-2.5 text-rose-400 hover:text-rose-300 hover:border-rose-500/40"
              :disabled="authStore.isLoading"
              @click="logout"
            >
              Sign Out / Revoke
            </button>
          </div>
        </div>

        <!-- Login Form Card -->
        <div v-else class="glass-panel p-6 rounded-2xl border border-white/10 space-y-6">
          <div class="flex items-center space-x-3">
            <img src="/logo.png" alt="IPA Downloader" class="w-10 h-10 object-contain rounded-xl" />
            <div>
              <h2 class="text-lg font-bold text-white">Sign In to Apple ID</h2>
              <p class="text-xs text-slate-400">
                Sign in with your Apple ID credentials to acquire licenses and download signed .ipa packages.
              </p>
            </div>
          </div>

          <form class="space-y-4" @submit.prevent="handleLogin">
            <div class="space-y-1.5">
              <label class="text-xs font-semibold text-slate-300">Apple ID (Email)</label>
              <input
                v-model="email"
                type="email"
                required
                placeholder="name@icloud.com"
                class="glass-input w-full px-3.5 py-2.5 rounded-xl text-sm outline-none"
              />
            </div>

            <div class="space-y-1.5">
              <label class="text-xs font-semibold text-slate-300">Password</label>
              <input
                v-model="password"
                type="password"
                required
                placeholder="••••••••••••"
                class="glass-input w-full px-3.5 py-2.5 rounded-xl text-sm outline-none"
              />
            </div>

            <div class="flex items-center justify-between pt-1">
              <label class="flex items-center space-x-2 cursor-pointer select-none text-xs text-slate-300">
                <input
                  v-model="rememberMe"
                  type="checkbox"
                  class="rounded bg-slate-800 border-white/10 text-blue-600 focus:ring-0"
                />
                <span>Remember in secure local storage</span>
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
                  <span>Authenticating...</span>
                </span>
                <span v-else>Sign In</span>
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Right Column: Features & Quick Tips -->
      <div class="md:col-span-5 space-y-4">
        <div class="glass-panel p-5 rounded-2xl border border-white/10 space-y-3.5">
          <h3 class="text-sm font-bold uppercase tracking-wider text-slate-300">Features & Capabilities</h3>
          
          <div class="space-y-3 text-xs text-slate-300">
            <div class="flex items-start space-x-3">
              <span class="text-base">🚀</span>
              <div>
                <span class="font-bold text-white">Live Search & Version Query</span>
                <p class="text-slate-400 mt-0.5">Explore iOS, iPadOS, and tvOS apps with complete version build numbers.</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <span class="text-base">⚡</span>
              <div>
                <span class="font-bold text-white">Chunked Streaming Transfers</span>
                <p class="text-slate-400 mt-0.5">Concurrent queue with live speed calculation, pause/resume, and auto-licensing.</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <span class="text-base">🔐</span>
              <div>
                <span class="font-bold text-white">FairPlay SINF DRM Signing</span>
                <p class="text-slate-400 mt-0.5">Patches downloaded .ipa packages with your personal FairPlay DRM signatures for sideloading.</p>
              </div>
            </div>

            <div class="flex items-start space-x-3">
              <span class="text-base">💾</span>
              <div>
                <span class="font-bold text-white">Offline SQLite Persistence</span>
                <p class="text-slate-400 mt-0.5">Zero-CGO local database for favorites, downloads queue, search history, and settings.</p>
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
import { useNotifications } from '../composables/useNotifications'

const authStore = useAuthStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
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
  if (authStore.isLoggedIn) return 'bg-emerald-400 shadow-sm shadow-emerald-400/50 animate-pulse-subtle'
  if (authStore.isLoading) return 'bg-amber-400 animate-ping'
  return 'bg-slate-500'
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
