<template>
  <div class="max-w-4xl mx-auto space-y-6 animate-slide-up">
    <!-- Page Header -->
    <div>
      <h1 class="text-2xl font-bold">Preferences & Settings</h1>
      <p class="text-xs text-slate-400 mt-0.5">Customize application appearance, network concurrency, and storage locations.</p>
    </div>

    <!-- Main Settings Form -->
    <div class="glass-panel p-6 rounded-2xl border border-white/10 space-y-6">
      <!-- Theme Setting -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-2 pb-5 border-b border-white/5">
        <div>
          <h3 class="text-sm font-bold">Appearance Theme</h3>
          <p class="text-xs text-slate-400">Choose between dark, light, or automatic system theme.</p>
        </div>
        <div class="flex p-1 rounded-xl bg-slate-900/60 border border-white/5 shrink-0">
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-lg text-xs font-semibold transition"
            :class="settings.theme === 'dark' ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-white'"
            @click="setTheme('dark')"
          >
            🌙 Dark
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-lg text-xs font-semibold transition"
            :class="settings.theme === 'light' ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-white'"
            @click="setTheme('light')"
          >
            ☀️ Light
          </button>
          <button
            type="button"
            class="px-3.5 py-1.5 rounded-lg text-xs font-semibold transition"
            :class="settings.theme === 'system' ? 'bg-blue-600 text-white shadow-sm' : 'text-slate-400 hover:text-white'"
            @click="setTheme('system')"
          >
            💻 System
          </button>
        </div>
      </div>

      <!-- Default Download Folder -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 pb-5 border-b border-white/5">
        <div class="space-y-0.5">
          <h3 class="text-sm font-bold">Default Destination Folder</h3>
          <p class="text-xs text-slate-400">Directory where downloaded .ipa packages will be saved.</p>
          <p class="text-xs font-mono text-blue-400 truncate max-w-md pt-1">{{ settings.defaultDownloadFolder || 'Default ~/Downloads' }}</p>
        </div>
        <button
          type="button"
          class="btn-secondary text-xs px-4 py-2 shrink-0 flex items-center space-x-1.5"
          @click="browseFolder"
        >
          <span>📁 Browse Folder</span>
        </button>
      </div>

      <!-- Concurrent Downloads Limit -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-4 pb-5 border-b border-white/5">
        <div>
          <h3 class="text-sm font-bold">Concurrent Downloads</h3>
          <p class="text-xs text-slate-400">Maximum number of parallel downloads active at the same time (1–10).</p>
        </div>
        <div class="flex items-center space-x-3">
          <input
            v-model.number="settings.maxConcurrentDownloads"
            type="range"
            min="1"
            max="10"
            class="w-32 accent-blue-600"
            @change="saveSettings"
          />
          <span class="text-sm font-mono font-bold w-6 text-center">{{ settings.maxConcurrentDownloads }}</span>
        </div>
      </div>

      <!-- Automatic License Acquisition -->
      <div class="flex items-center justify-between pb-5 border-b border-white/5">
        <div>
          <h3 class="text-sm font-bold">Auto-Acquire Free License</h3>
          <p class="text-xs text-slate-400">Automatically purchase free App Store license if not yet registered to your Apple ID.</p>
        </div>
        <input
          v-model="settings.autoAcquireLicense"
          type="checkbox"
          class="w-5 h-5 rounded text-blue-600 focus:ring-blue-500 bg-slate-900 border-white/20"
          @change="saveSettings"
        />
      </div>

      <!-- Automatic Updates -->
      <div class="flex items-center justify-between pb-5 border-b border-white/5">
        <div>
          <h3 class="text-sm font-bold">Automatic Update Checks</h3>
          <p class="text-xs text-slate-400">Check for new IPATool Desktop versions on launch.</p>
        </div>
        <input
          v-model="settings.autoCheckUpdates"
          type="checkbox"
          class="w-5 h-5 rounded text-blue-600 focus:ring-blue-500 bg-slate-900 border-white/20"
          @change="saveSettings"
        />
      </div>

      <!-- Offline Metadata Cache Management -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3 pb-5 border-b border-white/5">
        <div>
          <h3 class="text-sm font-bold">Local App Cache</h3>
          <p class="text-xs text-slate-400">Current SQLite database size: <span class="text-slate-200 font-mono">{{ settingsStore.cacheSize }}</span></p>
        </div>
        <button
          type="button"
          class="btn-secondary text-xs px-4 py-2 text-rose-400 hover:text-rose-300"
          @click="clearCache"
        >
          Clear Cache
        </button>
      </div>

      <!-- Export System Diagnostic Logs -->
      <div class="flex flex-col md:flex-row md:items-center justify-between gap-3">
        <div>
          <h3 class="text-sm font-bold">Diagnostic Logs</h3>
          <p class="text-xs text-slate-400">Export internal backend runtime logs to a file for troubleshooting.</p>
        </div>
        <button
          type="button"
          class="btn-primary text-xs px-4 py-2 flex items-center space-x-1.5"
          @click="exportLogs"
        >
          <span>Export Logs</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useSettingsStore } from '../stores/settings'
import { useNotifications } from '../composables/useNotifications'

const settingsStore = useSettingsStore()
const { showToast } = useNotifications()

const settings = computed(() => settingsStore.settings)

onMounted(() => {
  settingsStore.fetchSettings()
})

function setTheme(theme: 'dark' | 'light' | 'system') {
  settingsStore.updateSettings({ theme })
  showToast('Theme Updated', `Switched theme to ${theme}`, 'info')
}

async function browseFolder() {
  await settingsStore.browseFolder()
  showToast('Saved', 'Download directory updated', 'success')
}

async function saveSettings() {
  await settingsStore.updateSettings(settings.value)
}

async function clearCache() {
  await settingsStore.clearCache()
  showToast('Cache Emptied', 'App cache has been cleared', 'info')
}

async function exportLogs() {
  const path = await settingsStore.exportLogs()
  if (path) {
    showToast('Logs Exported', `Saved to ${path}`, 'success')
  }
}
</script>
