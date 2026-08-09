<template>
  <div class="h-full flex flex-col space-y-8 animate-slide-up font-sans p-2">
    <!-- Hero Header -->
    <div class="flex items-center justify-between">
      <div class="space-y-1.5">
        <h1 class="text-4xl font-extrabold tracking-tight text-white flex items-center space-x-3">
          <span>{{ t.home.title }}</span>
          <div class="h-2 w-2 rounded-full bg-[#0A84FF] animate-pulse"></div>
        </h1>
        <p class="text-sm text-[#8E8E93] font-medium max-w-xl">{{ t.home.subtitle }}</p>
      </div>
    </div>

    <!-- Widgets Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 auto-rows-fr">
      <!-- Search Widget -->
      <div
        @click="router.push('/search')"
        class="group relative overflow-hidden p-6 rounded-[32px] bg-white/[0.03] border border-white/[0.08] hover:border-[#0A84FF]/40 hover:bg-white/[0.06] transition-all cursor-pointer shadow-xl"
      >
        <div class="flex flex-col h-full justify-between space-y-4">
          <div class="flex items-center justify-between">
            <div class="w-12 h-12 rounded-2xl bg-[#0A84FF]/10 flex items-center justify-center text-[#0A84FF] group-hover:scale-110 transition-transform duration-300">
              <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </div>
            <svg class="w-5 h-5 text-white/10 group-hover:text-white/30 transition-colors" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M9 5l7 7-7 7" />
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">{{ t.home.searchWidget }}</h3>
            <p class="text-xs text-[#8E8E93] mt-1">{{ t.home.searchWidgetDesc }}</p>
          </div>
        </div>
      </div>

      <!-- Transfers Widget -->
      <div
        @click="router.push('/downloads')"
        class="group relative overflow-hidden p-6 rounded-[32px] bg-white/[0.03] border border-white/[0.08] hover:border-[#64D2FF]/40 hover:bg-white/[0.06] transition-all cursor-pointer shadow-xl"
      >
        <div class="flex flex-col h-full justify-between space-y-4">
          <div class="flex items-center justify-between">
            <div class="w-12 h-12 rounded-2xl bg-[#64D2FF]/10 flex items-center justify-center text-[#64D2FF] group-hover:scale-110 transition-transform duration-300">
              <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
              </svg>
            </div>
            <div v-if="downloadsStore.activeCount > 0" class="px-2.5 py-1 rounded-full bg-[#64D2FF]/20 text-[#64D2FF] text-[10px] font-bold animate-pulse">
              {{ downloadsStore.activeCount }} {{ t.common.active.toUpperCase() }}
            </div>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">{{ t.home.downloadsWidget }}</h3>
            <p class="text-xs text-[#8E8E93] mt-1">{{ t.home.downloadsWidgetDesc }}</p>
          </div>
        </div>
      </div>

      <!-- Devices Widget -->
      <div
        @click="router.push('/apps')"
        class="group relative overflow-hidden p-6 rounded-[32px] bg-white/[0.03] border border-white/[0.08] hover:border-[#30D158]/40 hover:bg-white/[0.06] transition-all cursor-pointer shadow-xl"
      >
        <div class="flex flex-col h-full justify-between space-y-4">
          <div class="flex items-center justify-between">
            <div class="w-12 h-12 rounded-2xl bg-[#30D158]/10 flex items-center justify-center text-[#30D158] group-hover:scale-110 transition-transform duration-300">
              <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
            </div>
            <div v-if="deviceStore.devices.length > 0" class="px-2.5 py-1 rounded-full bg-[#30D158]/20 text-[#30D158] text-[10px] font-bold">
              {{ deviceStore.devices.length }} {{ t.common.connectedStatus.toUpperCase() }}
            </div>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">{{ t.home.appsWidget }}</h3>
            <p class="text-xs text-[#8E8E93] mt-1">{{ t.home.appsWidgetDesc }}</p>
          </div>
        </div>
      </div>

      <!-- Firmwares Widget -->
      <div
        @click="router.push('/firmwares')"
        class="group relative overflow-hidden p-6 rounded-[32px] bg-white/[0.03] border border-white/[0.08] hover:border-[#5E5CE6]/40 hover:bg-white/[0.06] transition-all cursor-pointer shadow-xl"
      >
        <div class="flex flex-col h-full justify-between space-y-4">
          <div class="flex items-center justify-between">
            <div class="w-12 h-12 rounded-2xl bg-[#5E5CE6]/10 flex items-center justify-center text-[#5E5CE6] group-hover:scale-110 transition-transform duration-300">
              <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
              </svg>
            </div>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">{{ t.home.firmwareWidget }}</h3>
            <p class="text-xs text-[#8E8E93] mt-1">{{ t.home.firmwareWidgetDesc }}</p>
          </div>
        </div>
      </div>

      <!-- Favorites Widget -->
      <div
        @click="router.push('/favorites')"
        class="group relative overflow-hidden p-6 rounded-[32px] bg-white/[0.03] border border-white/[0.08] hover:border-[#FFD60A]/40 hover:bg-white/[0.06] transition-all cursor-pointer shadow-xl"
      >
        <div class="flex flex-col h-full justify-between space-y-4">
          <div class="flex items-center justify-between">
            <div class="w-12 h-12 rounded-2xl bg-[#FFD60A]/10 flex items-center justify-center text-[#FFD60A] group-hover:scale-110 transition-transform duration-300">
              <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
              </svg>
            </div>
            <div class="text-[10px] font-bold text-[#FFD60A]">
              {{ favoritesStore.favorites.length }} {{ t.common.saved.toUpperCase() }}
            </div>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">{{ t.home.favoritesWidget }}</h3>
            <p class="text-xs text-[#8E8E93] mt-1">{{ t.home.favoritesWidgetDesc }}</p>
          </div>
        </div>
      </div>

      <!-- Logs Widget -->
      <div
        @click="router.push('/logs')"
        class="group relative overflow-hidden p-6 rounded-[32px] bg-white/[0.03] border border-white/[0.08] hover:border-white/20 hover:bg-white/[0.06] transition-all cursor-pointer shadow-xl"
      >
        <div class="flex flex-col h-full justify-between space-y-4">
          <div class="flex items-center justify-between">
            <div class="w-12 h-12 rounded-2xl bg-white/5 flex items-center justify-center text-white group-hover:scale-110 transition-transform duration-300">
              <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
          </div>
          <div>
            <h3 class="text-lg font-bold text-white">{{ t.home.logsWidget }}</h3>
            <p class="text-xs text-[#8E8E93] mt-1">{{ t.home.logsWidgetDesc }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '../stores/auth'
import { useDownloadsStore } from '../stores/downloads'
import { useFavoritesStore } from '../stores/favorites'
import { useDeviceStore } from '../stores/device'
import { useI18n } from '../i18n'

const router = useRouter()
const authStore = useAuthStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const deviceStore = useDeviceStore()
const { t } = useI18n()

onMounted(async () => {
  await downloadsStore.fetchDownloads()
  await favoritesStore.fetchFavorites()
})
</script>
