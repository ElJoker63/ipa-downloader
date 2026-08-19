<template>
  <div class="max-w-6xl mx-auto space-y-6 flex flex-col h-full animate-slide-up font-sans">
    <!-- Top Header -->
    <div class="flex flex-col md:flex-row items-center justify-between gap-4 shrink-0">
      <div>
        <h1 class="text-2xl font-bold tracking-tight text-[#FFFFFF] flex items-center space-x-3">
          <span>{{ t.downloadedApps?.title || 'IPAs Descargadas' }}</span>
          <span
            v-if="downloadedAppsStore.downloadedIPAs.length > 0"
            class="px-2.5 py-0.5 text-xs font-semibold rounded-full bg-[#0A84FF]/20 text-[#0A84FF] border border-[#0A84FF]/30"
          >
            {{ downloadedAppsStore.downloadedIPAs.length }}
          </span>
        </h1>
        <p class="text-xs text-[#B8C0CC] mt-0.5 font-normal">
          {{ t.downloadedApps?.subtitle || 'Biblioteca de archivos .ipa detectados localmente en tu carpeta de descargas' }}
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center space-x-3">
        <!-- Destination Folder Badge -->
        <div class="hidden lg:flex items-center space-x-2 px-3.5 py-2 rounded-xl bg-white/[0.04] border border-white/[0.08]">
          <svg class="w-4 h-4 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
          </svg>
          <div class="flex flex-col">
            <span class="text-[9px] uppercase font-bold text-[#7D8592] leading-none">{{ t.downloads?.destinationFolder || 'Carpeta de Descargas' }}</span>
            <span class="text-[11px] text-[#B8C0CC] font-mono truncate max-w-[200px]">{{ settingsStore.settings.defaultDownloadFolder || 'Downloads' }}</span>
          </div>
        </div>

        <button
          type="button"
          @click="openDownloadFolder"
          class="btn-secondary text-xs px-3.5 py-2 flex items-center space-x-1.5"
        >
          <svg class="w-4 h-4 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 19a2 2 0 01-2-2V7a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M5 19h14a2 2 0 002-2v-5a2 2 0 00-2-2H9a2 2 0 00-2 2v5a2 2 0 01-2 2z" />
          </svg>
          <span>{{ t.history?.showInFolder || 'Abrir Carpeta' }}</span>
        </button>

        <button
          type="button"
          @click="refreshIPAs"
          :disabled="downloadedAppsStore.isLoading"
          class="px-3.5 py-2 rounded-xl bg-white/[0.06] hover:bg-white/[0.12] border border-white/[0.08] text-xs font-medium text-white transition flex items-center space-x-2 disabled:opacity-50"
        >
          <svg class="w-4 h-4 text-[#8E8E93]" :class="{ 'animate-spin': downloadedAppsStore.isLoading }" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
            <path stroke-linecap="round" stroke-linejoin="round" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          <span>{{ t.apps?.refresh || 'Actualizar' }}</span>
        </button>
      </div>
    </div>

    <!-- Search & Filter Controls -->
    <div class="flex items-center space-x-3 shrink-0">
      <div class="relative flex-1">
        <svg class="w-4 h-4 text-[#7D8592] absolute left-3.5 top-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
        </svg>
        <input
          v-model="searchQuery"
          type="text"
          :placeholder="t.downloadedApps?.searchPlaceholder || 'Filtrar por nombre, Bundle ID o nombre de archivo...'"
          class="glass-input w-full pl-10 pr-4 py-2.5 text-xs"
        />
      </div>

      <div class="text-xs text-[#7D8592] font-mono px-3 py-2 rounded-xl bg-white/[0.03] border border-white/[0.06]">
        {{ totalStorageSizeFormatted }}
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 min-h-0 overflow-y-auto">
      <!-- IPAs Grid -->
      <div v-if="filteredIPAs.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 pb-6">
        <div
          v-for="ipa in filteredIPAs"
          :key="ipa.filePath"
          class="glass-card p-5 rounded-[20px] flex flex-col justify-between space-y-4 border border-white/[0.08] hover:border-white/20 transition-all shadow-xl group"
        >
          <!-- Header Info -->
          <div class="flex items-start space-x-3.5">
            <img
              :src="ipa.artworkUrl || defaultAppIcon"
              :alt="ipa.appName"
              class="w-14 h-14 rounded-[14px] object-cover bg-[#171A21] border border-white/[0.18] shadow-md shrink-0 group-hover:scale-105 transition-transform"
            />
            <div class="min-w-0 flex-1">
              <h3 class="text-sm font-semibold truncate text-[#FFFFFF]" :title="ipa.appName">
                {{ ipa.appName || 'Aplicación desconocida' }}
              </h3>
              <p class="text-xs text-[#B8C0CC] font-mono truncate mt-0.5" :title="ipa.bundleId">
                {{ ipa.bundleId }}
              </p>

              <!-- Version & Size Badges -->
              <div class="flex flex-wrap items-center gap-1.5 mt-2">
                <span class="px-2 py-0.5 text-[10px] font-mono font-semibold rounded-md bg-[#0A84FF]/20 text-[#0A84FF] border border-[#0A84FF]/30">
                  v{{ ipa.version || ipa.shortVersion || '1.0' }}
                </span>
                <span class="px-2 py-0.5 text-[10px] font-mono rounded-md bg-white/[0.06] text-[#B8C0CC] border border-white/[0.08]">
                  {{ ipa.formattedSize }}
                </span>
                <span v-if="ipa.minimumOs" class="text-[10px] text-[#7D8592] font-mono">
                  iOS {{ ipa.minimumOs }}+
                </span>
              </div>
            </div>
          </div>

          <!-- File Name & Path Details -->
          <div class="px-3 py-2 rounded-xl bg-white/[0.03] border border-white/[0.05] text-[11px] text-[#7D8592] space-y-1">
            <div class="flex items-center justify-between font-mono">
              <span class="truncate text-[#B8C0CC]" :title="ipa.fileName">📄 {{ ipa.fileName }}</span>
            </div>
            <div class="text-[10px] text-[#7D8592]">
              Modificado: {{ formatDate(ipa.modTime) }}
            </div>
          </div>

          <!-- Action Buttons Footer -->
          <div class="flex items-center justify-between pt-2 border-t border-white/[0.08] gap-2">
            <button
              v-if="deviceStore.isConnected"
              type="button"
              class="btn-primary text-xs px-3 py-1.5 flex-1 flex items-center justify-center space-x-1.5"
              @click="installToDevice(ipa.filePath)"
            >
              <svg class="w-3.5 h-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
                <path stroke-linecap="round" stroke-linejoin="round" d="M12 18h.01M8 21h8a2 2 0 002-2V5a2 2 0 00-2-2H8a2 2 0 00-2 2v14a2 2 0 002 2z" />
              </svg>
              <span>{{ t.common?.install || 'Instalar' }}</span>
            </button>

            <button
              type="button"
              class="btn-secondary text-xs px-3 py-1.5 flex items-center space-x-1"
              @click="revealFile(ipa.filePath)"
              title="Mostrar en carpeta"
            >
              <svg class="w-3.5 h-3.5 text-[#64D2FF]" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
              </svg>
              <span class="hidden sm:inline">{{ t.history?.showInFolder || 'Ver' }}</span>
            </button>

            <button
              type="button"
              class="p-1.5 rounded-lg text-[#7D8592] hover:text-[#FF453A] hover:bg-[#FF453A]/15 transition duration-150"
              title="Eliminar IPA"
              @click="deleteIPA(ipa.filePath)"
            >
              <svg class="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.8">
                <path stroke-linecap="round" stroke-linejoin="round" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State Panel -->
      <div v-else-if="!downloadedAppsStore.isLoading" class="glass-card p-12 rounded-[22px] text-center space-y-4 max-w-lg mx-auto mt-12">
        <div class="w-16 h-16 rounded-full bg-white/[0.05] border border-white/10 flex items-center justify-center text-[#0A84FF] mx-auto shadow-inner">
          <svg class="w-8 h-8" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="1.5">
            <path stroke-linecap="round" stroke-linejoin="round" d="M5 8h14M5 8a2 2 0 01-2-2V5a2 2 0 012-2h14a2 2 0 012 2v1a2 2 0 01-2 2M5 8v10a2 2 0 002 2h14a2 2 0 002-2V8m-9 4h4" />
          </svg>
        </div>
        <div class="space-y-1">
          <h3 class="text-base font-semibold text-[#FFFFFF]">{{ t.downloadedApps?.emptyTitle || 'No se encontraron archivos .ipa' }}</h3>
          <p class="text-xs text-[#B8C0CC] leading-relaxed">
            {{ t.downloadedApps?.emptyDesc || 'Descarga aplicaciones desde la pestaña Buscar o coloca archivos .ipa en tu carpeta de descargas.' }}
          </p>
        </div>
        <button
          type="button"
          @click="openDownloadFolder"
          class="btn-secondary text-xs px-4 py-2"
        >
          {{ t.history?.showInFolder || 'Abrir Carpeta de Descargas' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useDownloadedAppsStore } from '../stores/downloadedApps'
import { useSettingsStore } from '../stores/settings'
import { useDeviceStore } from '../stores/device'
import { useI18n } from '../i18n'
import { useNotifications } from '../composables/useNotifications'
import { WailsService } from '../services/wails'

const downloadedAppsStore = useDownloadedAppsStore()
const settingsStore = useSettingsStore()
const deviceStore = useDeviceStore()
const { t } = useI18n()
const { showToast } = useNotifications()

const searchQuery = ref('')
const defaultAppIcon = 'https://is1-ssl.mzstatic.com/image/thumb/Purple126/v4/app_icon.png/512x512bb.png'

onMounted(() => {
  downloadedAppsStore.fetchDownloadedIPAs()
})

const filteredIPAs = computed(() => {
  const q = searchQuery.value.trim().toLowerCase()
  if (!q) return downloadedAppsStore.downloadedIPAs

  return downloadedAppsStore.downloadedIPAs.filter((ipa) => {
    return (
      ipa.appName.toLowerCase().includes(q) ||
      ipa.bundleId.toLowerCase().includes(q) ||
      ipa.fileName.toLowerCase().includes(q) ||
      ipa.version.toLowerCase().includes(q)
    )
  })
})

const totalStorageSizeFormatted = computed(() => {
  const totalBytes = downloadedAppsStore.downloadedIPAs.reduce((acc, ipa) => acc + (ipa.fileSizeBytes || 0), 0)
  if (totalBytes <= 0) return '0 B'
  const unit = 1024
  if (totalBytes < unit) return `${totalBytes} B`
  const div = 1024
  const exp = Math.floor(Math.log(totalBytes) / Math.log(div))
  const letter = 'KMGTPE'[exp - 1] || 'K'
  return `${(totalBytes / Math.pow(div, exp)).toFixed(1)} ${letter}B Total`
})

function refreshIPAs() {
  downloadedAppsStore.fetchDownloadedIPAs()
}

function openDownloadFolder() {
  const folder = settingsStore.settings.defaultDownloadFolder
  WailsService.openFolder(folder)
}

function revealFile(filePath: string) {
  WailsService.revealInExplorer(filePath)
}

async function installToDevice(filePath: string) {
  if (!deviceStore.selectedUdid) {
    showToast(t.value.apps.noDeviceTitle || 'No hay dispositivo seleccionado', '', 'error')
    return
  }

  try {
    await WailsService.installIPA(deviceStore.selectedUdid, filePath)
    showToast(t.value.common.installing || 'Instalación encolada en el dispositivo', '', 'success')
  } catch (err: any) {
    showToast(err?.message || 'Error al instalar IPA', '', 'error')
  }
}

async function deleteIPA(filePath: string) {
  try {
    await WailsService.deleteFile(filePath)
    showToast('success', 'Archivo IPA eliminado')
    downloadedAppsStore.fetchDownloadedIPAs()
  } catch (err: any) {
    showToast('error', err?.message || 'Error al borrar el archivo')
  }
}

function formatDate(dateStr: string): string {
  if (!dateStr) return ''
  try {
    const d = new Date(dateStr)
    return d.toLocaleDateString()
  } catch {
    return dateStr
  }
}
</script>
