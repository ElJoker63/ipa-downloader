<template>
  <div class="max-w-4xl mx-auto space-y-6 animate-slide-up font-sans">
    <!-- Header Section -->
    <div>
      <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF]">
        {{ t.settings.title }}
      </h1>
      <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
        {{ t.settings.subtitle }}
      </p>
    </div>

    <div class="space-y-5 pb-8">
      <!-- General & Downloads Preferences -->
      <div class="glass-card p-6 rounded-[18px] space-y-5">
        <h2 class="text-sm font-semibold uppercase tracking-wider text-[#B8C0CC]">{{ t.settings.generalSection }}</h2>

        <!-- Language Selector -->
        <div class="space-y-2">
          <label class="text-xs font-medium text-[#FFFFFF]">{{ t.settings.languageLabel }}</label>
          <div class="grid grid-cols-2 gap-3 max-w-sm">
            <button
              type="button"
              class="p-3 rounded-[12px] border text-xs font-semibold flex items-center justify-center space-x-2 transition-all duration-150"
              :class="currentLanguage === 'es' ? 'bg-[#0A84FF]/20 border-[#0A84FF] text-white shadow-sm' : 'bg-white/[0.04] border-white/[0.08] text-[#B8C0CC] hover:text-white'"
              @click="changeLang('es')"
            >
              <span class="px-1.5 py-0.5 rounded text-[10px] font-mono font-bold bg-white/[0.1] text-[#64D2FF]">ES</span>
              <span>{{ t.common.spanish }}</span>
            </button>
            <button
              type="button"
              class="p-3 rounded-[12px] border text-xs font-semibold flex items-center justify-center space-x-2 transition-all duration-150"
              :class="currentLanguage === 'en' ? 'bg-[#0A84FF]/20 border-[#0A84FF] text-white shadow-sm' : 'bg-white/[0.04] border-white/[0.08] text-[#B8C0CC] hover:text-white'"
              @click="changeLang('en')"
            >
              <span class="px-1.5 py-0.5 rounded text-[10px] font-mono font-bold bg-white/[0.1] text-[#64D2FF]">EN</span>
              <span>{{ t.common.english }}</span>
            </button>
          </div>
        </div>

        <!-- Default Download Folder -->
        <div class="space-y-2.5">
          <div class="flex items-center justify-between">
            <label class="text-xs font-medium text-[#FFFFFF]">{{ t.settings.defaultFolder }}</label>
            <button
              type="button"
              class="text-xs text-[#0A84FF] hover:underline font-medium flex items-center space-x-1"
              @click="openFolder"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              <span>{{ t.history?.showInFolder || 'Abrir Carpeta' }}</span>
            </button>
          </div>
          <div class="flex items-center space-x-2.5">
            <input
              v-model="settingsStore.settings.defaultDownloadFolder"
              type="text"
              readonly
              class="glass-input flex-1 px-3.5 py-2.5 text-xs font-mono select-all"
            />
            <button
              type="button"
              class="btn-secondary text-xs px-4 py-2.5 shrink-0 flex items-center space-x-1.5"
              @click="browseFolder"
            >
              <svg class="w-3.5 h-3.5 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
              </svg>
              <span>{{ t.settings.browseFolder }}</span>
            </button>
          </div>

          <!-- Informative Subdirectory Structure Banner -->
          <div class="p-3.5 rounded-[12px] bg-white/[0.03] border border-white/[0.06] space-y-1 text-[#B8C0CC]">
            <div class="font-semibold text-white flex items-center space-x-1.5 text-[11px] uppercase tracking-wider">
              <svg class="w-3.5 h-3.5 text-[#0A84FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <span>Organización de Archivos Descargados</span>
            </div>
            <ul class="list-disc list-inside space-y-0.5 text-[11px] font-mono text-[#7D8592] pl-1">
              <li><strong class="text-[#B8C0CC]">Aplicaciones (.ipa):</strong> se descargarán en <span class="text-[#64D2FF]">.../ipa/</span></li>
              <li><strong class="text-[#B8C0CC]">Firmwares (.ipsw):</strong> se descargarán en <span class="text-[#64D2FF]">.../ipsw/</span></li>
            </ul>
          </div>
        </div>

        <!-- Max Concurrent Downloads Slider -->
        <div class="space-y-2 pt-1">
          <div class="flex items-center justify-between">
            <label class="text-xs font-medium text-[#FFFFFF]">{{ t.settings.concurrentLimit }}</label>
            <span class="text-xs font-mono font-semibold text-[#0A84FF]">
              {{ settingsStore.settings.maxConcurrentDownloads }} {{ t.settings.simultaneousJobs }}
            </span>
          </div>
          <input
            v-model.number="settingsStore.settings.maxConcurrentDownloads"
            type="range"
            min="1"
            max="10"
            class="w-full accent-[#0A84FF] cursor-pointer"
            @change="saveSettings"
          />
          <div class="flex justify-between text-[10px] text-[#7D8592] font-mono">
            <span>{{ t.settings.conservative }}</span>
            <span>{{ t.settings.recommended }}</span>
            <span>{{ t.settings.highSpeed }}</span>
          </div>
        </div>

        <!-- Automation Toggles -->
        <div class="space-y-3 pt-2">
          <!-- Auto Acquire License -->
          <div class="flex items-center justify-between p-3.5 rounded-[14px] bg-white/[0.04] border border-white/[0.08]">
            <div class="pr-4">
              <div class="text-xs font-semibold text-[#FFFFFF]">{{ t.settings.autoLicenseTitle }}</div>
              <p class="text-[11px] text-[#B8C0CC] mt-0.5">{{ t.settings.autoLicenseDesc }}</p>
            </div>
            <input
              v-model="settingsStore.settings.autoAcquireLicense"
              type="checkbox"
              class="w-5 h-5 rounded-md text-[#0A84FF] bg-white/[0.08] border-white/[0.18] focus:ring-0 cursor-pointer shrink-0"
              @change="saveSettings"
            />
          </div>

          <!-- Remember Credentials -->
          <div class="flex items-center justify-between p-3.5 rounded-[14px] bg-white/[0.04] border border-white/[0.08]">
            <div class="pr-4">
              <div class="text-xs font-semibold text-[#FFFFFF]">{{ t.settings.keychainTitle }}</div>
              <p class="text-[11px] text-[#B8C0CC] mt-0.5">{{ t.settings.keychainDesc }}</p>
            </div>
            <input
              v-model="settingsStore.settings.rememberCredentials"
              type="checkbox"
              class="w-5 h-5 rounded-md text-[#0A84FF] bg-white/[0.08] border-white/[0.18] focus:ring-0 cursor-pointer shrink-0"
              @change="saveSettings"
            />
          </div>
        </div>
      </div>

      <!-- Storage & Diagnostics Glass Card -->
      <div class="glass-card p-6 rounded-[18px] space-y-5">
        <h2 class="text-sm font-semibold uppercase tracking-wider text-[#B8C0CC]">{{ t.settings.storageSection }}</h2>

        <div class="flex items-center justify-between p-3.5 rounded-[14px] bg-white/[0.04] border border-white/[0.08]">
          <div>
            <div class="text-xs font-semibold text-[#FFFFFF]">{{ t.settings.cacheTitle }}</div>
            <p class="text-[11px] text-[#B8C0CC] mt-0.5">{{ t.settings.cacheDesc }} <span class="font-mono text-[#64D2FF]">{{ settingsStore.cacheSize }}</span></p>
          </div>
          <button
            type="button"
            class="btn-secondary text-xs px-3.5 py-1.5"
            @click="clearCache"
          >
            {{ t.settings.clearCache }}
          </button>
        </div>

        <div class="flex items-center justify-between p-3.5 rounded-[14px] bg-white/[0.04] border border-white/[0.08]">
          <div>
            <div class="text-xs font-semibold text-[#FFFFFF]">{{ t.settings.exportLogsTitle }}</div>
            <p class="text-[11px] text-[#B8C0CC] mt-0.5">{{ t.settings.exportLogsDesc }}</p>
          </div>
          <button
            type="button"
            class="btn-primary text-xs px-4 py-1.5 shadow-sm"
            @click="exportLogs"
          >
            {{ t.settings.exportLogsButton }}
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useSettingsStore } from '../stores/settings'
import { useLogsStore } from '../stores/logs'
import { useI18n, type LanguageCode } from '../i18n'
import { useNotifications } from '../composables/useNotifications'
import { WailsService } from '../services/wails'

const settingsStore = useSettingsStore()
const logsStore = useLogsStore()
const { t, currentLanguage, setLanguage } = useI18n()
const { showToast } = useNotifications()

onMounted(async () => {
  await settingsStore.fetchSettings()
})

function changeLang(lang: LanguageCode) {
  setLanguage(lang)
  settingsStore.settings.language = lang
  settingsStore.updateSettings({ language: lang })
}

function openFolder() {
  const folder = settingsStore.settings.defaultDownloadFolder
  WailsService.openFolder(folder)
}

async function browseFolder() {
  await settingsStore.browseFolder()
  showToast(t.value.downloads.destinationFolder, settingsStore.settings.defaultDownloadFolder, 'info')
}

async function saveSettings() {
  await settingsStore.updateSettings(settingsStore.settings)
  showToast(t.value.settings.savedToast, '', 'info')
}

async function clearCache() {
  await settingsStore.clearCache()
  showToast('Cache Cleared', 'Offline app cache purged', 'info')
}

async function exportLogs() {
  const path = await logsStore.exportLogs()
  if (path) {
    showToast(t.value.logs.exportedToast, path, 'success')
  }
}
</script>
