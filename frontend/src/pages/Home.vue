<template>
  <div class="max-w-5xl mx-auto space-y-8 animate-slide-up font-sans">
    <!-- Header Section -->
    <div class="flex items-center justify-between">
      <div class="space-y-2">
        <h1 class="text-4xl font-bold tracking-tight text-[#FFFFFF]">
          {{ t.home.title }}
        </h1>
        <p class="text-base text-[#B8C0CC] font-normal max-w-2xl leading-relaxed">
          Professional suite for Apple device management, firmware deployment, and App Store IPA acquisition with FairPlay DRM signing.
        </p>
      </div>

      <!-- Live Connection Status Pill (macOS Capsule) -->
      <div class="flex items-center space-x-2 px-4 py-2 rounded-[14px] glass-panel shadow-lg">
        <span class="w-2.5 h-2.5 rounded-full" :class="statusDotClass"></span>
        <span class="text-xs font-semibold text-[#FFFFFF]">{{ authStore.isLoggedIn ? t.common.connected : t.common.notConnected }}</span>
      </div>
    </div>

    <!-- Quick Stats Cards Row (8px Grid & Glass Design) -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
      <div class="glass-card p-6 rounded-[24px] hover:bg-white/[0.04] transition-all cursor-pointer group" @click="router.push('/downloads')">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-widest text-[#7D8592]">{{ t.home.activeQueue }}</span>
          <div class="w-8 h-8 rounded-xl bg-[#0A84FF]/10 flex items-center justify-center text-[#0A84FF] group-hover:scale-110 transition-transform">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
            </svg>
          </div>
        </div>
        <div class="text-3xl font-bold mt-4 text-[#FFFFFF]">{{ downloadsStore.activeCount }}</div>
        <div class="text-xs text-[#0A84FF] font-medium mt-1">{{ downloadsStore.totalSpeedFormatted }}</div>
      </div>

      <div class="glass-card p-6 rounded-[24px] hover:bg-white/[0.04] transition-all cursor-pointer group" @click="router.push('/favorites')">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-widest text-[#7D8592]">{{ t.home.savedFavorites }}</span>
          <div class="w-8 h-8 rounded-xl bg-[#FFD60A]/10 flex items-center justify-center text-[#FFD60A] group-hover:scale-110 transition-transform">
            <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
              <path stroke-linecap="round" stroke-linejoin="round" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
            </svg>
          </div>
        </div>
        <div class="text-3xl font-bold mt-4 text-[#FFFFFF]">{{ favoritesStore.favorites.length }}</div>
        <div class="text-xs text-[#7D8592] font-medium mt-1">{{ t.home.bookmarkedApps }}</div>
      </div>

      <div class="glass-card p-6 rounded-[24px] hover:bg-white/[0.04] transition-all cursor-pointer group" @click="router.push('/apps')">
        <div class="flex items-center justify-between">
          <span class="text-xs font-bold uppercase tracking-widest text-[#7D8592]">Devices</span>
          <div class="w-8 h-8 rounded-xl bg-[#30D158]/10 flex items-center justify-center text-[#30D158] group-hover:scale-110 transition-transform">
             <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
             </svg>
          </div>
        </div>
        <div class="text-3xl font-bold mt-4 text-[#FFFFFF]">{{ deviceStore.devices.length }}</div>
        <div class="text-xs text-[#30D158] font-medium mt-1">{{ deviceStore.isConnected ? 'Ready for deployment' : 'Connect via USB' }}</div>
      </div>
    </div>

    <!-- Main Grid: Features -->
    <div class="grid grid-cols-1 md:grid-cols-12 gap-6 items-start">
      <div class="md:col-span-12 space-y-4">
        <div class="glass-card p-8 rounded-[32px] space-y-8">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-bold text-white">{{ t.home.featuresTitle }}</h3>
            <div class="flex space-x-2">
              <div v-for="i in 4" :key="i" class="w-1.5 h-1.5 rounded-full bg-white/10"></div>
            </div>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-8">
            <div class="space-y-3 group">
              <div class="w-12 h-12 rounded-2xl bg-[#0A84FF]/10 flex items-center justify-center text-[#0A84FF] group-hover:bg-[#0A84FF] group-hover:text-white transition-all">
                <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
              <div class="font-bold text-white">{{ t.home.feat1Title }}</div>
              <p class="text-xs text-[#7D8592] leading-relaxed">{{ t.home.feat1Desc }}</p>
            </div>

            <div class="space-y-3 group">
              <div class="w-12 h-12 rounded-2xl bg-[#64D2FF]/10 flex items-center justify-center text-[#64D2FF] group-hover:bg-[#64D2FF] group-hover:text-white transition-all">
                <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                </svg>
              </div>
              <div class="font-bold text-white">{{ t.home.feat2Title }}</div>
              <p class="text-xs text-[#7D8592] leading-relaxed">{{ t.home.feat2Desc }}</p>
            </div>

            <div class="space-y-3 group">
              <div class="w-12 h-12 rounded-2xl bg-[#30D158]/10 flex items-center justify-center text-[#30D158] group-hover:bg-[#30D158] group-hover:text-white transition-all">
                <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z" />
                </svg>
              </div>
              <div class="font-bold text-white">{{ t.home.feat3Title }}</div>
              <p class="text-xs text-[#7D8592] leading-relaxed">{{ t.home.feat3Desc }}</p>
            </div>

            <div class="space-y-3 group">
              <div class="w-12 h-12 rounded-2xl bg-[#FFD60A]/10 flex items-center justify-center text-[#FFD60A] group-hover:bg-[#FFD60A] group-hover:text-white transition-all">
                <svg class="w-6 h-6" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                  <path stroke-linecap="round" stroke-linejoin="round" d="M4 7v10c0 2.21 3.582 4 8 4s8-1.79 8-4V7M4 7c0 2.21 3.582 4 8 4s8-1.79 8-4M4 7c0-2.21 3.582-4 8-4s8 1.79 8 4" />
                </svg>
              </div>
              <div class="font-bold text-white">{{ t.home.feat4Title }}</div>
              <p class="text-xs text-[#7D8592] leading-relaxed">{{ t.home.feat4Desc }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
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

const statusDotClass = computed(() => {
  if (authStore.isLoggedIn) return 'bg-[#30D158] shadow-[0_0_8px_rgba(48,209,88,0.8)] animate-pulse-subtle'
  if (authStore.isLoading) return 'bg-[#FFD60A] animate-ping'
  return 'bg-[#7D8592]'
})
</script>
