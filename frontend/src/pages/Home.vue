<template>
  <div class="max-w-5xl mx-auto space-y-8 animate-slide-up">
    <!-- Header Section -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-extrabold tracking-tight text-slate-100 dark:text-slate-100 light:text-slate-900">
          App Store Authentication
        </h1>
        <p class="text-sm text-slate-400 dark:text-slate-400 light:text-slate-600 mt-1">
          Securely authenticate with your Apple ID to search, query version history, and download iOS/tvOS packages.
        </p>
      </div>

      <!-- Live Connection Status Pill -->
      <div class="flex items-center space-x-2 px-4 py-2 rounded-2xl glass-panel border border-white/10">
        <span class="w-3 h-3 rounded-full" :class="statusDotClass"></span>
        <span class="text-sm font-semibold">{{ authStore.status }}</span>
      </div>
    </div>

    <!-- Main Grid: Account Status / Login Form -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-6 items-start">
      <!-- Left Column: Authenticated Account or Sign-in Form -->
      <div class="md:col-span-7 space-y-6">
        <!-- Logged In Card -->
        <div v-if="authStore.isLoggedIn" class="glass-panel p-6 rounded-2xl border border-white/10 space-y-6">
          <div class="flex items-center space-x-4">
            <div class="w-16 h-16 rounded-2xl bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center text-2xl font-bold text-white shadow-lg shadow-blue-500/25">
              {{ authStore.account.name ? authStore.account.name.charAt(0).toUpperCase() : '' }}
            </div>
            <div>
              <h2 class="text-xl font-bold">{{ authStore.account.name || 'Apple ID User' }}</h2>
              <p class="text-sm text-slate-400 dark:text-slate-400 light:text-slate-600 font-mono">{{ authStore.account.email }}</p>
              <div class="flex items-center gap-2 mt-1">
                <span class="px-2 py-0.5 text-[11px] font-semibold rounded bg-emerald-500/20 text-emerald-400 border border-emerald-500/30">
                  Storefront: {{ authStore.account.storeFrontCountry || 'US' }}
                </span>
                <span v-if="authStore.account.directoryServicesId" class="text-[11px] font-mono text-slate-500">
                  DSID: {{ authStore.account.directoryServicesId }}
                </span>
              </div>
            </div>
          </div>

          <div class="p-4 rounded-xl bg-white/5 dark:bg-white/5 light:bg-black/5 border border-white/5 space-y-2 text-xs text-slate-300 dark:text-slate-300 light:text-slate-700">
            <div class="flex justify-between">
              <span class="text-slate-400">Connection Status:</span>
              <span class="text-emerald-400 font-semibold">● Active Session</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">Storefront Code:</span>
              <span class="font-mono">{{ authStore.account.storeFront || '143441-1,29' }}</span>
            </div>
            <div class="flex justify-between">
              <span class="text-slate-400">Keychain Storage:</span>
              <span class="text-blue-400">Encrypted</span>
            </div>
          </div>

          <div class="flex items-center justify-end space-x-3 pt-2">
            <router-link to="/search" class="btn-primary text-sm px-5 py-2 flex items-center space-x-2">
              <span>🔍 Search Apps</span>
            </router-link>
            <button
              type="button"
              class="btn-secondary text-sm px-4 py-2 text-rose-400 hover:text-rose-300 hover:border-rose-500/40"
              :disabled="authStore.isLoading"
              @click="logout"
            >
              Sign Out / Revoke
            </button>
          </div>
        </div>

        <!-- Sign-In Form (When Not Logged In) -->
        <div v-else class="glass-panel p-6 rounded-2xl border border-white/10 space-y-5">
          <div>
            <h2 class="text-xl font-bold">Sign In with Apple ID</h2>
            <p class="text-xs text-slate-400 mt-1">Credentials are stored locally in your system keychain. Two-factor authentication (2FA) is supported.</p>
          </div>

          <form @submit.prevent="handleLogin" class="space-y-4">
            <!-- Email Input -->
            <div class="space-y-1.5">
              <label for="apple-id" class="text-xs font-semibold text-slate-300 dark:text-slate-300 light:text-slate-700">Apple ID (Email)</label>
              <input
                id="apple-id"
                v-model="email"
                type="email"
                required
                placeholder="name@icloud.com"
                class="glass-input w-full px-4 py-2.5 rounded-xl text-sm font-medium"
              />
            </div>

            <!-- Password Input -->
            <div class="space-y-1.5">
              <label for="apple-password" class="text-xs font-semibold text-slate-300 dark:text-slate-300 light:text-slate-700">Password</label>
              <input
                id="apple-password"
                v-model="password"
                type="password"
                required
                placeholder="••••••••••••"
                class="glass-input w-full px-4 py-2.5 rounded-xl text-sm font-medium"
              />
            </div>

            <!-- 2FA Code Input (Optional on first try) -->
            <div class="space-y-1.5">
              <label for="apple-2fa" class="text-xs font-semibold text-slate-300 dark:text-slate-300 light:text-slate-700 flex justify-between">
                <span>2FA Verification Code</span>
                <span class="text-slate-500 font-normal">Optional if prompted</span>
              </label>
              <input
                id="apple-2fa"
                v-model="authCode"
                type="text"
                placeholder="6-digit code (e.g. 123456)"
                class="glass-input w-full px-4 py-2.5 rounded-xl text-sm font-mono"
              />
            </div>

            <!-- Remember Credentials Checkbox -->
            <div class="flex items-center space-x-2 pt-1">
              <input
                id="remember"
                v-model="rememberCredentials"
                type="checkbox"
                class="w-4 h-4 rounded text-blue-600 focus:ring-blue-500 bg-slate-900 border-white/20"
              />
              <label for="remember" class="text-xs text-slate-300 dark:text-slate-300 light:text-slate-700 cursor-pointer">
                Remember credentials securely in local keychain
              </label>
            </div>

            <!-- Error Banner -->
            <div v-if="authStore.errorMessage" class="p-3 rounded-xl bg-rose-500/10 border border-rose-500/20 text-xs text-rose-400">
              {{ authStore.errorMessage }}
            </div>

            <!-- Login Action Button -->
            <div class="pt-2">
              <button
                type="submit"
                class="btn-primary w-full py-3 text-sm font-semibold flex items-center justify-center space-x-2"
                :disabled="authStore.isLoading || !email || !password"
              >
                <span v-if="authStore.isLoading" class="w-4 h-4 border-2 border-white/30 border-t-white rounded-full animate-spin"></span>
                <span>{{ authStore.isLoading ? 'Connecting...' : 'Sign In to App Store' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>

      <!-- Right Column: Quick Stats & Helpful Tips -->
      <div class="md:col-span-5 space-y-6">
        <!-- Quick Stats Card -->
        <div class="glass-panel p-6 rounded-2xl border border-white/10 space-y-4">
          <h3 class="text-sm font-semibold uppercase tracking-wider text-slate-400">Activity Overview</h3>
          <div class="grid grid-cols-2 gap-3">
            <div class="p-3 rounded-xl bg-white/5 border border-white/5">
              <div class="text-2xl font-bold text-blue-400">{{ downloadsStore.downloads.length }}</div>
              <div class="text-xs text-slate-400 mt-0.5">Total Downloads</div>
            </div>
            <div class="p-3 rounded-xl bg-white/5 border border-white/5">
              <div class="text-2xl font-bold text-amber-400">{{ downloadsStore.activeCount }}</div>
              <div class="text-xs text-slate-400 mt-0.5">Active Queue</div>
            </div>
            <div class="p-3 rounded-xl bg-white/5 border border-white/5">
              <div class="text-2xl font-bold text-rose-400">{{ favoritesStore.favorites.length }}</div>
              <div class="text-xs text-slate-400 mt-0.5">Saved Favorites</div>
            </div>
            <div class="p-3 rounded-xl bg-white/5 border border-white/5">
              <div class="text-2xl font-bold text-emerald-400">{{ downloadsStore.totalSpeedFormatted }}</div>
              <div class="text-xs text-slate-400 mt-0.5">Current Speed</div>
            </div>
          </div>
        </div>

        <!-- Privacy & Security Notice -->
        <div class="glass-panel p-6 rounded-2xl border border-white/10 space-y-3">
          <div class="flex items-center space-x-2 text-blue-400">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
            </svg>
            <h4 class="text-sm font-bold">Privacy & Direct Connection</h4>
          </div>
          <p class="text-xs text-slate-300 leading-relaxed">
            All requests communicate directly with Apple's App Store servers over HTTPS using official protocols. Your password and session tokens are never shared with third parties.
          </p>
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
const authCode = ref('')
const rememberCredentials = ref(true)

const statusDotClass = computed(() => {
  if (authStore.status === 'Connected') {
    return 'bg-emerald-400 shadow-sm shadow-emerald-400/50'
  }
  if (authStore.status.includes('Connecting')) {
    return 'bg-amber-400 animate-pulse'
  }
  return 'bg-slate-500'
})

onMounted(() => {
  authStore.checkAccount()
  downloadsStore.fetchDownloads()
  favoritesStore.fetchFavorites()
})

async function handleLogin() {
  try {
    await authStore.login(email.value, password.value, authCode.value, rememberCredentials.value)
    showToast('Authenticated', `Connected as ${authStore.account.email}`, 'success')
  } catch (err: any) {
    if (!err?.message?.includes('2FA')) {
      showToast('Authentication Failed', err?.message || 'Login failed', 'error')
    }
  }
}

async function logout() {
  await authStore.logout()
  showToast('Logged Out', 'App Store session revoked', 'info')
}
</script>
