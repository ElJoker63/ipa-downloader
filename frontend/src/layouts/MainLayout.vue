<template>
  <div class="h-screen w-screen flex flex-col overflow-hidden bg-[#080C16] dark:bg-[#080C16] light:bg-[#F4F6F9] text-slate-100 dark:text-slate-100 light:text-slate-900 select-none">
    <!-- Top Custom Drag Title Bar -->
    <TitleBar />

    <!-- Application Body (Sidebar + Routed View) -->
    <div class="flex-1 flex min-h-0 overflow-hidden relative">
      <!-- Sidebar Navigation -->
      <aside class="w-64 bg-[#0B1120]/90 dark:bg-[#0B1120]/90 light:bg-white/90 backdrop-blur-2xl border-r border-white/5 dark:border-white/5 light:border-black/5 flex flex-col justify-between p-4 shrink-0 z-30 shadow-xl">
        <!-- Main Navigation Links -->
        <nav class="space-y-1.5 flex-1 overflow-y-auto">
          <router-link
            to="/"
            class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 group relative"
            :class="isActive('/') ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-slate-900 hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">🏠</span>
            <span class="flex-1">Account & Home</span>
            <span v-if="authStore.isLoggedIn" class="w-2 h-2 rounded-full bg-emerald-400 shadow-sm shadow-emerald-400/50"></span>
          </router-link>

          <router-link
            to="/search"
            class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 group relative"
            :class="isActive('/search') ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-slate-900 hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">🔍</span>
            <span class="flex-1">App Store Search</span>
          </router-link>

          <router-link
            to="/downloads"
            class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 group relative"
            :class="isActive('/downloads') ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-slate-900 hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">⬇</span>
            <span class="flex-1">Downloads</span>
            <span
              v-if="downloadsStore.activeCount > 0"
              class="px-2 py-0.5 text-[11px] font-bold rounded-full bg-blue-500 text-white shadow-md shadow-blue-500/40 animate-pulse"
            >
              {{ downloadsStore.activeCount }}
            </span>
          </router-link>

          <router-link
            to="/favorites"
            class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 group relative"
            :class="isActive('/favorites') ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-slate-900 hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">⭐</span>
            <span class="flex-1">Favorites</span>
            <span
              v-if="favoritesStore.favorites.length > 0"
              class="px-2 py-0.5 text-[10px] font-bold rounded-full bg-white/10 text-slate-300 font-mono"
            >
              {{ favoritesStore.favorites.length }}
            </span>
          </router-link>

          <router-link
            to="/history"
            class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 group relative"
            :class="isActive('/history') ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-slate-900 hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">📜</span>
            <span class="flex-1">History</span>
          </router-link>

          <div class="pt-2 pb-1">
            <div class="h-px bg-white/5 dark:bg-white/5 light:bg-black/5"></div>
          </div>

          <router-link
            to="/settings"
            class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 group relative"
            :class="isActive('/settings') ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-slate-900 hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">⚙</span>
            <span class="flex-1">Settings</span>
          </router-link>

          <router-link
            to="/logs"
            class="flex items-center space-x-3 px-3.5 py-2.5 rounded-xl text-sm font-semibold transition-all duration-200 group relative"
            :class="isActive('/logs') ? 'bg-gradient-to-r from-blue-600 to-indigo-600 text-white shadow-lg shadow-blue-600/30' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-slate-900 hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">📄</span>
            <span class="flex-1">Diagnostic Logs</span>
          </router-link>
        </nav>

        <!-- Bottom User & Theme Strip -->
        <div class="pt-3 border-t border-white/5 dark:border-white/5 light:border-black/5 space-y-2">
          <!-- Account Mini Card (Clean with Only User Name) -->
          <div v-if="authStore.isLoggedIn" class="p-2.5 rounded-xl bg-white/5 dark:bg-white/5 light:bg-black/5 border border-white/5 flex items-center space-x-2.5">
            <div class="w-7 h-7 rounded-lg bg-gradient-to-tr from-blue-600 to-indigo-500 flex items-center justify-center text-xs font-bold text-white shrink-0">
              {{ authStore.account.name ? authStore.account.name.charAt(0).toUpperCase() : '' }}
            </div>
            <div class="min-w-0 flex-1">
              <div class="text-xs font-bold truncate text-white dark:text-white light:text-slate-900">{{ authStore.account.name || 'Apple User' }}</div>
              <div class="text-[10px] text-emerald-400 font-medium">● Connected</div>
            </div>
          </div>

          <div class="flex items-center justify-between px-1">
            <button
              type="button"
              class="text-xs font-medium text-slate-400 hover:text-white dark:hover:text-white light:hover:text-black flex items-center space-x-1.5 py-1 px-2 rounded-lg hover:bg-white/5 transition"
              @click="toggleTheme"
            >
              <span>{{ isDark ? '🌙 Dark' : '☀️ Light' }}</span>
            </button>
            <span class="text-[10px] text-slate-500 font-mono">IPA Downloader</span>
          </div>
        </div>
      </aside>

      <!-- Main Router Content -->
      <main class="flex-1 flex flex-col min-w-0 overflow-y-auto bg-[#080C16] dark:bg-[#080C16] light:bg-[#F4F6F9] p-7">
        <router-view v-slot="{ Component }">
          <transition name="fade-fast" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>

    <!-- Modals & Overlays -->
    <TwoFactorModal />
    <AppDetailsModal />
    <ToastContainer />
  </div>
</template>

<script setup lang="ts">
import { useRoute } from 'vue-router'
import TitleBar from '../components/TitleBar.vue'
import TwoFactorModal from '../components/TwoFactorModal.vue'
import AppDetailsModal from '../components/AppDetailsModal.vue'
import ToastContainer from '../components/ToastContainer.vue'
import { useAuthStore } from '../stores/auth'
import { useDownloadsStore } from '../stores/downloads'
import { useFavoritesStore } from '../stores/favorites'
import { useTheme } from '../composables/useTheme'
import { useKeyboardShortcuts } from '../composables/useKeyboardShortcuts'

const route = useRoute()
const authStore = useAuthStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const { isDark, toggleTheme } = useTheme()

// Initialize keyboard shortcuts (Ctrl+1..7, Ctrl+K, Esc)
useKeyboardShortcuts()

function isActive(path: string): boolean {
  return route.path === path
}
</script>

<style scoped>
.fade-fast-enter-active,
.fade-fast-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
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
