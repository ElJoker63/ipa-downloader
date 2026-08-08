<template>
  <div class="h-screen w-screen flex flex-col overflow-hidden bg-[#0B0F19] dark:bg-[#0B0F19] light:bg-[#F5F5F7] text-slate-100 dark:text-slate-100 light:text-slate-900 font-sans">
    <!-- Top Native Title Bar -->
    <TitleBar />

    <div class="flex-1 flex overflow-hidden">
      <!-- Sleek Glass Sidebar -->
      <aside class="w-56 shrink-0 flex flex-col justify-between border-r border-white/5 dark:border-white/5 light:border-black/5 bg-[#131B2E]/60 dark:bg-[#131B2E]/60 light:bg-white/70 backdrop-blur-xl p-3 z-30">
        <!-- Navigation Links -->
        <nav class="space-y-1.5">
          <router-link
            to="/"
            class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-sm font-medium transition duration-150 group"
            :class="isActive('/') ? 'bg-blue-600 text-white shadow-md shadow-blue-500/25' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">🏠</span>
            <span class="flex-1">Home</span>
          </router-link>

          <router-link
            to="/search"
            class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-sm font-medium transition duration-150 group"
            :class="isActive('/search') ? 'bg-blue-600 text-white shadow-md shadow-blue-500/25' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">🔍</span>
            <span class="flex-1">Search</span>
            <span class="text-[10px] font-mono opacity-50 border border-current px-1 rounded">⌘K</span>
          </router-link>

          <router-link
            to="/downloads"
            class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-sm font-medium transition duration-150 group"
            :class="isActive('/downloads') ? 'bg-blue-600 text-white shadow-md shadow-blue-500/25' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">⬇</span>
            <span class="flex-1">Downloads</span>
            <span
              v-if="downloadsStore.activeCount > 0"
              class="px-1.5 py-0.2 text-[11px] font-bold rounded-full bg-blue-400 text-slate-950 animate-pulse"
            >
              {{ downloadsStore.activeCount }}
            </span>
          </router-link>

          <router-link
            to="/favorites"
            class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-sm font-medium transition duration-150 group"
            :class="isActive('/favorites') ? 'bg-blue-600 text-white shadow-md shadow-blue-500/25' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">⭐</span>
            <span class="flex-1">Favorites</span>
            <span v-if="favoritesStore.favorites.length > 0" class="text-xs font-mono text-slate-500">
              {{ favoritesStore.favorites.length }}
            </span>
          </router-link>

          <router-link
            to="/history"
            class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-sm font-medium transition duration-150 group"
            :class="isActive('/history') ? 'bg-blue-600 text-white shadow-md shadow-blue-500/25' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">📜</span>
            <span class="flex-1">History</span>
          </router-link>

          <router-link
            to="/settings"
            class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-sm font-medium transition duration-150 group"
            :class="isActive('/settings') ? 'bg-blue-600 text-white shadow-md shadow-blue-500/25' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">⚙</span>
            <span class="flex-1">Settings</span>
          </router-link>

          <router-link
            to="/logs"
            class="flex items-center space-x-3 px-3 py-2.5 rounded-xl text-sm font-medium transition duration-150 group"
            :class="isActive('/logs') ? 'bg-blue-600 text-white shadow-md shadow-blue-500/25' : 'text-slate-400 dark:text-slate-400 light:text-slate-600 hover:text-white dark:hover:text-white light:hover:text-black hover:bg-white/5 dark:hover:bg-white/5 light:hover:bg-black/5'"
          >
            <span class="text-base">📄</span>
            <span class="flex-1">Logs</span>
          </router-link>
        </nav>

        <!-- Bottom User & Theme Strip -->
        <div class="pt-3 border-t border-white/5 dark:border-white/5 light:border-black/5 space-y-2">
          <!-- Account Mini Card -->
          <div v-if="authStore.isLoggedIn" class="p-2.5 rounded-xl bg-white/5 dark:bg-white/5 light:bg-black/5 border border-white/5">
            <div class="text-xs font-semibold truncate">{{ authStore.account.name || 'Apple User' }}</div>
          </div>

          <div class="flex items-center justify-between px-1">
            <button
              type="button"
              class="text-xs text-slate-400 hover:text-white dark:hover:text-white light:hover:text-black flex items-center space-x-1.5 py-1"
              @click="toggleTheme"
            >
              <span>{{ isDark ? '🌙 Dark' : '☀️ Light' }}</span>
            </button>
            <span class="text-[10px] text-slate-500 font-mono">v2.0</span>
          </div>
        </div>
      </aside>

      <!-- Main Router Content -->
      <main class="flex-1 flex flex-col min-w-0 overflow-y-auto bg-[#0B0F19] dark:bg-[#0B0F19] light:bg-[#F5F5F7] p-6">
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

// Initialize keyboard shortcuts
useKeyboardShortcuts()

function isActive(path: string) {
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
