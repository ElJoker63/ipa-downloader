<template>
  <div class="h-full flex flex-col p-6 space-y-6 overflow-hidden">
    <!-- Header -->
    <div class="flex items-center justify-between shrink-0">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-white flex items-center space-x-3">
          <span>{{ t.firmwares.title }}</span>
          <span class="px-2.5 py-0.5 text-[10px] font-bold rounded-full bg-[#0A84FF]/20 text-[#0A84FF] border border-[#0A84FF]/30 uppercase">
            IPSW.me API
          </span>
        </h1>
        <p class="text-sm text-[#8E8E93] mt-1">
          {{ t.firmwares.subtitle }}
        </p>
      </div>

      <div class="flex items-center space-x-3">
        <div class="relative w-64">
           <svg class="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-[#8E8E93]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            :placeholder="t.firmwares.searchPlaceholder"
            class="w-full bg-white/[0.06] border border-white/[0.1] rounded-xl pl-9 pr-4 py-2 text-sm text-white focus:outline-none focus:border-[#0A84FF] transition"
          />
        </div>
        <button
          @click="fetchDevices"
          :disabled="loading"
          class="p-2 rounded-xl bg-white/[0.06] border border-white/[0.1] text-[#8E8E93] hover:text-white hover:bg-white/[0.12] transition"
        >
          <svg class="w-5 h-5" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 flex space-x-6 min-h-0">
      <!-- Device Sidebar / Categories -->
      <div class="w-72 flex flex-col space-y-4 shrink-0 overflow-y-auto pr-2 scrollbar-hide">
        <div class="space-y-1">
          <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-widest px-4 mb-2">{{ t.firmwares.categories }}</div>
          <button
            v-for="cat in categories"
            :key="cat.id"
            @click="activeCategory = cat.id"
            class="w-full px-4 py-2.5 rounded-xl text-left text-sm font-medium transition flex items-center justify-between"
            :class="activeCategory === cat.id ? 'bg-[#0A84FF] text-white shadow-lg shadow-[#0A84FF]/20' : 'text-[#8E8E93] hover:bg-white/[0.05] hover:text-white'"
          >
            <div class="flex items-center space-x-3">
              <span v-html="cat.icon"></span>
              <span>{{ cat.name }}</span>
            </div>
            <span class="text-[10px] opacity-60 bg-black/20 px-1.5 rounded-md">{{ getCategoryCount(cat.id) }}</span>
          </button>
        </div>

        <div class="flex-1 overflow-y-auto space-y-1">
           <div class="text-[10px] font-bold text-[#8E8E93] uppercase tracking-widest px-4 mb-2 mt-4">{{ t.firmwares.devices }}</div>
           <div v-if="loading" class="px-4 py-8 flex flex-col items-center space-y-2">
             <div class="w-5 h-5 border-2 border-[#0A84FF] border-t-transparent rounded-full animate-spin"></div>
             <span class="text-[10px] text-[#8E8E93]">{{ t.firmwares.loadingModels }}</span>
           </div>
           <button
            v-for="dev in filteredDevices"
            :key="dev.identifier"
            @click="selectDevice(dev)"
            class="w-full px-4 py-2 rounded-xl text-left transition group"
            :class="selectedDevice?.identifier === dev.identifier ? 'bg-white/[0.08] border border-white/10' : 'hover:bg-white/[0.04]'"
          >
            <div class="text-sm font-semibold truncate" :class="selectedDevice?.identifier === dev.identifier ? 'text-[#0A84FF]' : 'text-white group-hover:text-[#0A84FF]'">{{ dev.name }}</div>
            <div class="text-[10px] text-[#8E8E93] font-mono">{{ dev.identifier }}</div>
          </button>
        </div>
      </div>

      <!-- Firmware List (Right Side) -->
      <div class="flex-1 bg-[#171A21]/50 border border-white/[0.05] rounded-3xl overflow-hidden flex flex-col">
        <div v-if="!selectedDevice" class="flex-1 flex flex-col items-center justify-center space-y-4 opacity-40">
           <div class="w-20 h-20 rounded-full bg-white/[0.03] flex items-center justify-center">
             <svg class="w-10 h-10 text-[#8E8E93]" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M9 3v2m6-2v2M9 19v2m6-2v2M5 9H3m2 6H3m18-6h-2m2 6h-2M7 19h10a2 2 0 002-2V7a2 2 0 00-2-2H7a2 2 0 00-2 2v10a2 2 0 002 2zM9 9h6v6H9V9z" />
             </svg>
           </div>
           <div class="text-center">
             <div class="text-base font-bold text-white">{{ t.firmwares.selectDevice }}</div>
             <p class="text-xs text-[#8E8E93]">{{ t.firmwares.selectDeviceDesc }}</p>
           </div>
        </div>

        <template v-else>
          <!-- Firmware Header -->
          <div class="p-6 border-b border-white/10 flex items-center justify-between shrink-0 bg-white/[0.02]">
            <div class="flex items-center space-x-4">
              <div class="w-12 h-12 rounded-2xl bg-gradient-to-br from-[#0A84FF] to-[#5E5CE6] flex items-center justify-center text-white shadow-xl shadow-[#0A84FF]/20">
                <svg class="w-7 h-7" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
                </svg>
              </div>
              <div>
                <h2 class="text-lg font-bold text-white">{{ selectedDevice.name }}</h2>
                <div class="flex items-center space-x-2 text-[10px] font-mono text-[#8E8E93]">
                  <span>{{ selectedDevice.identifier }}</span>
                  <span>•</span>
                  <span>{{ selectedDevice.platform }}</span>
                </div>
              </div>
            </div>

            <div class="flex p-1 bg-white/[0.06] border border-white/[0.08] rounded-xl">
               <button
                @click="filterType = 'all'"
                class="px-3 py-1.5 rounded-lg text-[10px] font-bold uppercase transition"
                :class="filterType === 'all' ? 'bg-white/10 text-white' : 'text-[#8E8E93] hover:text-white'"
               >{{ t.firmwares.all }}</button>
               <button
                @click="filterType = 'signed'"
                class="px-3 py-1.5 rounded-lg text-[10px] font-bold uppercase transition flex items-center space-x-1.5"
                :class="filterType === 'signed' ? 'bg-[#30D158]/20 text-[#30D158]' : 'text-[#8E8E93] hover:text-[#30D158]'"
               >
                 <span>{{ t.firmwares.signed }}</span>
                 <div class="w-1.5 h-1.5 rounded-full bg-[#30D158] animate-pulse"></div>
               </button>
            </div>
          </div>

          <!-- Firmware List -->
          <div class="flex-1 overflow-y-auto p-6 scrollbar-hide">
            <div class="space-y-3">
              <div
                v-for="fw in filteredFirmwares"
                :key="fw.buildid"
                class="group p-4 rounded-2xl bg-white/[0.03] border border-white/[0.05] hover:border-white/10 hover:bg-white/[0.06] transition-all flex items-center justify-between"
              >
                <div class="flex items-center space-x-4 min-w-0">
                  <!-- Signing Status Icon -->
                  <div
                    class="w-10 h-10 rounded-full flex items-center justify-center shrink-0 border"
                    :class="fw.signed ? 'bg-[#30D158]/10 border-[#30D158]/30 text-[#30D158]' : 'bg-[#FF453A]/10 border-[#FF453A]/30 text-[#FF453A]'"
                  >
                    <svg v-if="fw.signed" class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M5 13l4 4L19 7" />
                    </svg>
                    <svg v-else class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2.5" d="M6 18L18 6M6 6l12 12" />
                    </svg>
                  </div>

                  <div class="min-w-0">
                    <div class="flex items-center space-x-2">
                      <span class="text-sm font-bold text-white">OS {{ fw.version }}</span>
                      <span class="text-[10px] font-mono px-1.5 py-0.5 bg-white/5 rounded text-[#8E8E93]">{{ fw.buildid }}</span>
                    </div>
                    <div class="flex items-center space-x-3 mt-1 text-[10px] text-[#8E8E93]">
                       <span class="flex items-center space-x-1">
                         <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" /></svg>
                         <span>{{ formatDate(fw.releasedate) }}</span>
                       </span>
                       <span class="flex items-center space-x-1">
                         <svg class="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" /></svg>
                         <span>{{ formatSize(fw.size) }}</span>
                       </span>
                    </div>
                  </div>
                </div>

                <button
                  @click="handleDownload(fw)"
                  class="px-4 py-2 rounded-xl bg-white/[0.08] hover:bg-[#0A84FF] text-white text-xs font-bold transition-all shadow-lg hover:shadow-[#0A84FF]/25 flex items-center space-x-2"
                >
                  <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2.5">
                    <path stroke-linecap="round" stroke-linejoin="round" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                  <span>{{ t.firmwares.download }}</span>
                </button>
              </div>
            </div>
          </div>
        </template>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { WailsService } from '../services/wails'
import { useNotifications } from '../composables/useNotifications'
import { useI18n } from '../i18n'
import type { AppleHardware, Firmware } from '../types'

const loading = ref(false)
const devices = ref<AppleHardware[]>([])
const selectedDevice = ref<AppleHardware | null>(null)
const searchQuery = ref('')
const activeCategory = ref('iPhone')
const filterType = ref<'all' | 'signed'>('all')

const { showToast } = useNotifications()
const { t } = useI18n()


const categories = [
  { id: 'iPhone', name: 'iPhone', icon: '<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>' },
  { id: 'iPad', name: 'iPad', icon: '<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 18h.01M7 21h10a2 2 0 002-2V5a2 2 0 00-2-2H7a2 2 0 00-2 2v14a2 2 0 002 2z" /></svg>' },
  { id: 'Mac', name: 'Mac', icon: '<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.75 17L9 21h6l-.75-4M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" /></svg>' },
  { id: 'AppleTV', name: 'Apple TV', icon: '<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 7l-8-4-8 4m16 0l-8 4m8-4v10l-8 4m0-10L4 7m8 4v10M4 7v10l8 4" /></svg>' },
  { id: 'AudioAccessory', name: 'HomePod', icon: '<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11a7 7 0 01-7 7m0 0a7 7 0 01-7-7m7 7v4m0 0H8m4 0h4m-4-8a3 3 0 01-3-3V5a3 3 0 116 0v6a3 3 0 01-3 3z" /></svg>' },
  { id: 'Watch', name: 'Apple Watch', icon: '<svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" /></svg>' },
]

async function fetchDevices() {
  loading.value = true
  try {
    devices.value = await WailsService.getAppleDevices()
  } catch (err) {
    showToast('API Error', 'Failed to fetch firmware list', 'error')
  } finally {
    loading.value = false
  }
}

async function selectDevice(dev: AppleHardware) {
  selectedDevice.value = dev
  // Optionally fetch fresh detail if needed, but the bulk API includes them
}

const filteredDevices = computed(() => {
  let list = devices.value.filter(d => d.identifier.startsWith(activeCategory.value) || d.name.toLowerCase().includes(activeCategory.value.toLowerCase()))

  // Custom mapping for categories that don't match identifier prefix directly
  if (activeCategory.value === 'iPhone') list = devices.value.filter(d => d.identifier.startsWith('iPhone'))
  if (activeCategory.value === 'iPad') list = devices.value.filter(d => d.identifier.startsWith('iPad'))
  if (activeCategory.value === 'Mac') list = devices.value.filter(d => d.platform === 'macOS' || d.identifier.startsWith('Mac'))
  if (activeCategory.value === 'AppleTV') list = devices.value.filter(d => d.identifier.startsWith('AppleTV'))
  if (activeCategory.value === 'AudioAccessory') list = devices.value.filter(d => d.identifier.startsWith('AudioAccessory'))
  if (activeCategory.value === 'Watch') list = devices.value.filter(d => d.identifier.startsWith('Watch'))

  if (searchQuery.value) {
    const q = searchQuery.value.toLowerCase()
    list = list.filter(d => d.name.toLowerCase().includes(q) || d.identifier.toLowerCase().includes(q))
  }
  return list
})

const filteredFirmwares = computed(() => {
  if (!selectedDevice.value?.firmwares) return []
  let list = [...selectedDevice.value.firmwares]
  if (filterType.value === 'signed') {
    list = list.filter(f => f.signed)
  }
  return list
})

function getCategoryCount(catId: string) {
    if (catId === 'iPhone') return devices.value.filter(d => d.identifier.startsWith('iPhone')).length
    if (catId === 'iPad') return devices.value.filter(d => d.identifier.startsWith('iPad')).length
    if (catId === 'Mac') return devices.value.filter(d => d.platform === 'macOS' || d.identifier.startsWith('Mac')).length
    if (catId === 'AppleTV') return devices.value.filter(d => d.identifier.startsWith('AppleTV')).length
    if (catId === 'AudioAccessory') return devices.value.filter(d => d.identifier.startsWith('AudioAccessory')).length
    if (catId === 'Watch') return devices.value.filter(d => d.identifier.startsWith('Watch')).length
    return 0
}

async function handleDownload(fw: Firmware) {
  if (!selectedDevice.value) return
  try {
    showToast('Download Queued', `${fw.filename} added to transfers`, 'success')
    await WailsService.downloadFirmware(selectedDevice.value.name, fw)
  } catch (err: any) {
    showToast('Download Error', err.message || String(err), 'error')
  }
}

function formatDate(dateStr: string) {
  if (!dateStr) return 'N/A'
  try {
    return new Date(dateStr).toLocaleDateString(undefined, { year: 'numeric', month: 'short', day: 'numeric' })
  } catch {
    return dateStr
  }
}

function formatSize(bytes: number) {
  const gb = bytes / (1024 * 1024 * 1024)
  return `${gb.toFixed(2)} GB`
}

onMounted(fetchDevices)
</script>

<style scoped>
.scrollbar-hide::-webkit-scrollbar {
  display: none;
}
.scrollbar-hide {
  -ms-overflow-style: none;
  scrollbar-width: none;
}
</style>
