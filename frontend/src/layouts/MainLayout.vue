<template>
  <div class="h-screen w-screen flex overflow-hidden bg-[#0F1115] text-[#FFFFFF] select-none font-sans">
    <!-- macOS / visionOS Style Glass Sidebar -->
    <aside class="w-64 bg-[#171A21]/80 backdrop-blur-[30px] border-r border-white/[0.08] flex flex-col justify-between p-3.5 shrink-0 z-30 shadow-[0_12px_40px_rgba(0,0,0,0.25)]">
      <!-- Top Brand Header in Sidebar -->
      <div class="px-3 py-2 flex items-center justify-between border-b border-white/[0.08] pb-3 mb-2">
        <div class="flex items-center space-x-2.5">
          <img src="/logo.png" alt="IPA Downloader" class="w-20 h-20 rounded-lg object-contain shadow-sm shrink-0" />
          <!--div>
            <div class="text-[13px] font-bold tracking-tight text-white leading-tight">IPA Downloader</div>
          </div-->
        </div>
      </div>

      <!-- Main Navigation Links -->
      <nav class="space-y-1 flex-1 overflow-y-auto">
        <!-- Section Label -->
        <div class="px-3 pt-1 pb-1.5 text-[11px] font-medium tracking-wider uppercase text-[#7D8592]">
          {{ t.nav.menu }}
        </div>

        <!-- Home / Account -->
        <router-link
          to="/"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          <span class="flex-1">{{ t.nav.home }}</span>
          <span v-if="authStore.isLoggedIn" class="w-1.5 h-1.5 rounded-full bg-[#30D158] shadow-[0_0_6px_rgba(48,209,88,0.8)]"></span>
        </router-link>

        <!-- Search -->
        <router-link
          to="/search"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/search') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <span class="flex-1">{{ t.nav.search }}</span>
        </router-link>

        <!-- Downloads -->
        <router-link
          to="/downloads"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/downloads') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          <span class="flex-1">{{ t.nav.downloads }}</span>
          <span
            v-if="downloadsStore.activeCount > 0"
            class="px-2 py-0.5 text-[11px] font-medium rounded-full bg-[#0A84FF] text-white shadow-sm shadow-[#0A84FF]/40 animate-pulse"
          >
            {{ downloadsStore.activeCount }}
          </span>
        </router-link>

        <!-- Apps (Device Management) -->
        <router-link
          to="/apps"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/apps') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
          </svg>
          <span class="flex-1">{{ t.nav.apps || 'Apps' }}</span>
          <span
            v-if="deviceStore.isConnected"
            class="w-2 h-2 rounded-full bg-[#30D158] shadow-[0_0_6px_rgba(48,209,88,0.8)]"
          ></span>
        </router-link>

        <!-- Firmwares -->
        <router-link
          to="/firmwares"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/firmwares') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
          </svg>
          <span class="flex-1">Firmwares</span>
        </router-link>

        <!-- Favorites -->

        <router-link
          to="/favorites"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/favorites') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
          <span class="flex-1">{{ t.nav.favorites }}</span>
          <span
            v-if="favoritesStore.favorites.length > 0"
            class="px-2 py-0.5 text-[10px] font-mono font-medium rounded-full bg-white/[0.08] text-[#B8C0CC]"
          >
            {{ favoritesStore.favorites.length }}
          </span>
        </router-link>

        <!-- History -->
        <router-link
          to="/history"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/history') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          <span class="flex-1">{{ t.nav.history }}</span>
        </router-link>

        <!-- Separator -->
        <div class="py-2">
          <div class="h-px bg-white/[0.08]"></div>
        </div>

        <!-- Section Label -->
        <div class="px-3 pb-1.5 text-[11px] font-medium tracking-wider uppercase text-[#7D8592]">
          {{ t.nav.system }}
        </div>

        <!-- Settings -->
        <router-link
          to="/settings"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/settings') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          <span class="flex-1">{{ t.nav.settings }}</span>
        </router-link>

        <!-- Logs -->
        <router-link
          to="/logs"
          class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-medium transition-all duration-200 group"
          :class="isActive('/logs') ? 'nav-item-active' : 'nav-item-inactive'"
        >
          <svg class="w-4 h-4 shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
            <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <span class="flex-1">{{ t.nav.logs }}</span>
        </router-link>
      </nav>

      <!-- Bottom User Profile, Language Switcher & Branding Strip -->
      <div class="pt-3 border-t border-white/[0.08] space-y-2.5">
        <!-- Account Mini Card with Dynamic Initials Avatar -->
        <div v-if="authStore.isLoggedIn" class="p-2.5 rounded-[14px] bg-white/[0.06] border border-white/[0.12] flex items-center space-x-2.5 backdrop-blur-md">
          <div class="w-8 h-8 rounded-[10px] bg-gradient-to-tr from-[#0071E3] via-[#0A84FF] to-[#64D2FF] flex items-center justify-center text-xs font-bold text-white shadow-sm border border-white/20 shrink-0 select-none">
            {{ userInitials }}
          </div>
          <div class="min-w-0 flex-1">
            <div class="text-xs font-semibold truncate text-[#FFFFFF]">{{ authStore.account.name || 'Apple User' }}</div>
            <div class="text-[10px] text-[#30D158] font-medium flex items-center space-x-1">
              <span class="w-1.5 h-1.5 rounded-full bg-[#30D158]"></span>
              <span>{{ t.common.connected }}</span>
            </div>
          </div>
        </div>

        <!-- Language & Theme Switcher Bar -->
        <div class="flex items-center justify-between px-1 gap-1">
          <!-- Language Selector -->
          <div class="flex items-center rounded-lg bg-white/[0.06] border border-white/[0.08] p-0.5">
            <button
              type="button"
              class="px-2 py-0.5 text-[10px] font-semibold rounded transition"
              :class="currentLanguage === 'es' ? 'bg-[#0A84FF] text-white shadow-xs' : 'text-[#B8C0CC] hover:text-white'"
              @click="setLanguage('es')"
            >
              ES
            </button>
            <button
              type="button"
              class="px-2 py-0.5 text-[10px] font-semibold rounded transition"
              :class="currentLanguage === 'en' ? 'bg-[#0A84FF] text-white shadow-xs' : 'text-[#B8C0CC] hover:text-white'"
              @click="setLanguage('en')"
            >
              EN
            </button>
          </div>

          <!-- Theme Toggle Button with Clean SVG Icons -->
          <button
            type="button"
            class="text-xs font-medium text-[#B8C0CC] hover:text-white flex items-center space-x-1.5 py-1 px-2.5 rounded-lg hover:bg-white/[0.06] transition"
            @click="toggleTheme"
          >
            <svg v-if="isDark" class="w-3.5 h-3.5 text-[#FFD60A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
            </svg>
            <svg v-else class="w-3.5 h-3.5 text-[#FFD60A]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
            </svg>
            <span class="text-[10px] font-semibold">{{ isDark ? t.common.dark : t.common.light }}</span>
          </button>
        </div>
      </div>
    </aside>

    <!-- Main Canvas Content Area -->
    <main class="flex-1 flex flex-col min-w-0 overflow-y-auto bg-[#0F1115] p-8">
      <router-view v-slot="{ Component }">
        <transition name="fade-fast" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </main>

    <!-- Modals & Overlays -->
    <TwoFactorModal />
    <AppDetailsModal />
    <ToastContainer />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import TwoFactorModal from '../components/TwoFactorModal.vue'
import AppDetailsModal from '../components/AppDetailsModal.vue'
import ToastContainer from '../components/ToastContainer.vue'
import { useAuthStore } from '../stores/auth'
import { useDownloadsStore } from '../stores/downloads'
import { useFavoritesStore } from '../stores/favorites'
import { useDeviceStore } from '../stores/device'
import { useTheme } from '../composables/useTheme'
import { useI18n } from '../i18n'
import { useKeyboardShortcuts } from '../composables/useKeyboardShortcuts'

const route = useRoute()
const authStore = useAuthStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const deviceStore = useDeviceStore()
const { isDark, toggleTheme } = useTheme()
const { t, currentLanguage, setLanguage } = useI18n()

// Initialize global keyboard shortcuts (Ctrl+1..7, Ctrl+K, Esc)
useKeyboardShortcuts()

function isActive(path: string): boolean {
  return route.path === path
}

const statusBadgeClass = computed(() => {
  if (authStore.isLoggedIn) {
    return 'bg-[#30D158]/15 border-[#30D158]/30 text-[#30D158]'
  }
  if (authStore.isLoading) {
    return 'bg-[#FFD60A]/15 border-[#FFD60A]/30 text-[#FFD60A] animate-pulse'
  }
  return 'bg-white/[0.08] border-white/[0.18] text-[#B8C0CC]'
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
</script>

<style scoped>
.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: opacity 150ms ease, transform 150ms ease;
}

.fade-fast-enter-from {
  opacity: 0;
  transform: translateY(4px);
}

.fade-fast-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
</style>
