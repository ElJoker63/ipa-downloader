import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { DeviceInfo, InstalledApp, DeviceInstallProgress } from '../types'
import { WailsService } from '../services/wails'

export const useDeviceStore = defineStore('device', () => {
  const device = ref<DeviceInfo | null>(null)
  const installedApps = ref<InstalledApp[]>([])
  const isLoadingApps = ref(false)
  const installProgress = ref<DeviceInstallProgress | null>(null)
  const activeTab = ref<'user' | 'system'>('user')
  const isInstalling = ref(false)
  const installError = ref<string | null>(null)

  const isConnected = computed(() => !!device.value && device.value.isConnected)

  const userApps = computed(() => {
    return installedApps.value.filter((a) => a.appType.toLowerCase() === 'user' || !a.appType)
  })

  const systemApps = computed(() => {
    return installedApps.value.filter((a) => a.appType.toLowerCase() === 'system')
  })

  async function checkDevice() {
    try {
      const dev = await WailsService.getConnectedDevice()
      device.value = dev
      if (dev && dev.isConnected) {
        await fetchApps()
      }
    } catch {
      device.value = null
    }
  }

  async function fetchApps() {
    if (!isConnected.value) return
    isLoadingApps.value = true
    try {
      const apps = await WailsService.listInstalledApps(activeTab.value)
      installedApps.value = apps || []
    } catch (err) {
      installedApps.value = []
    } finally {
      isLoadingApps.value = false
    }
  }

  async function pairDevice() {
    try {
      await WailsService.pairDevice()
      await checkDevice()
    } catch (err: any) {
      throw new Error(err || 'Failed to pair device')
    }
  }

  async function installIPA(path?: string) {
    let filePath = path
    if (!filePath) {
      filePath = await WailsService.selectIPAFile()
    }
    if (!filePath) return

    isInstalling.value = true
    installError.value = null
    installProgress.value = { phase: 'Preparing', percent: 5, message: 'Initializing installation...' }

    try {
      await WailsService.installIPA(filePath)
      await fetchApps()
    } catch (err: any) {
      installError.value = err?.message || String(err)
      installProgress.value = { phase: 'Failed', percent: 0, message: installError.value || 'Installation failed' }
      throw err
    } finally {
      setTimeout(() => {
        if (!installError.value) {
          isInstalling.value = false
          installProgress.value = null
        }
      }, 3000)
    }
  }

  async function uninstallApp(bundleId: string) {
    try {
      await WailsService.uninstallApp(bundleId)
      installedApps.value = installedApps.value.filter((a) => a.bundleId !== bundleId)
    } catch (err: any) {
      throw err
    }
  }

  function initListeners() {
    // Listen to real-time device connection events
    WailsService.onEvent('device:connected', (devInfo: DeviceInfo) => {
      device.value = devInfo
      fetchApps()
    })

    WailsService.onEvent('device:disconnected', () => {
      device.value = null
      installedApps.value = []
    })

    WailsService.onEvent('device:install_progress', (prog: DeviceInstallProgress) => {
      installProgress.value = prog
      if (prog.phase === 'Complete') {
        isInstalling.value = false
      } else if (prog.phase === 'Failed') {
        installError.value = prog.message
      }
    })

    WailsService.onEvent('device:install_complete', () => {
      fetchApps()
    })
  }

  return {
    device,
    installedApps,
    isLoadingApps,
    installProgress,
    activeTab,
    isInstalling,
    installError,
    isConnected,
    userApps,
    systemApps,
    checkDevice,
    fetchApps,
    pairDevice,
    installIPA,
    uninstallApp,
    initListeners,
  }
})
