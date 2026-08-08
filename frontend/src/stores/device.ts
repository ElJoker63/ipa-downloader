import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { DeviceInfo, InstalledApp, DeviceInstallTask } from '../types'
import { WailsService } from '../services/wails'

export const useDeviceStore = defineStore('device', () => {
  const devices = ref<DeviceInfo[]>([])
  const selectedUdid = ref<string>('')
  const installedApps = ref<InstalledApp[]>([])
  const isLoadingApps = ref(false)
  const installTasks = ref<DeviceInstallTask[]>([])
  const activeTab = ref<'user' | 'system'>('user')

  const isInstalling = ref(false)
  const installError = ref<string | null>(null)


  const selectedDevice = computed(() => {
    return devices.value.find((d) => d.udid === selectedUdid.value) || null
  })

  const isConnected = computed(() => !!selectedDevice.value && selectedDevice.value.isConnected)

  const userApps = computed(() => {
    return installedApps.value.filter((a) => a.appType.toLowerCase() === 'user' || !a.appType)
  })

  const systemApps = computed(() => {
    return installedApps.value.filter((a) => a.appType.toLowerCase() === 'system')
  })

  async function checkDevices() {
    try {
      const list = await WailsService.listConnectedDevices()
      devices.value = list || []

      if (devices.value.length > 0 && !selectedUdid.value) {
        selectedUdid.value = devices.value[0].udid
      }

      if (selectedUdid.value) {
        await fetchApps()
      }
    } catch {
      devices.value = []
    }
  }

  async function fetchApps() {
    if (!selectedUdid.value) return
    isLoadingApps.value = true
    try {
      const apps = await WailsService.listInstalledApps(selectedUdid.value, activeTab.value)
      installedApps.value = apps || []
    } catch (err) {
      installedApps.value = []
    } finally {
      isLoadingApps.value = false
    }
  }

  async function pairDevice(udid?: string) {
    const targetUdid = udid || selectedUdid.value
    if (!targetUdid) return
    try {
      await WailsService.pairDevice(targetUdid)
      await checkDevices()
    } catch (err: any) {
      throw new Error(err || 'Failed to pair device')
    }
  }

  async function installIPA(path?: string, udid?: string) {
    const targetUdid = udid || selectedUdid.value
    if (!targetUdid) throw new Error('No device selected')

    let filePath = path
    if (!filePath) {
      filePath = await WailsService.selectIPAFile()
    }
    if (!filePath) return

    isInstalling.value = true
    installError.value = null

    try {
      await WailsService.installIPA(targetUdid, filePath)
    } catch (err: any) {
      installError.value = err?.message || String(err)
      isInstalling.value = false
      throw err
    }
  }

  async function installMultipleIPAs(paths: string[], udid?: string) {
    const targetUdid = udid || selectedUdid.value
    if (!targetUdid) throw new Error('No device selected')

    isInstalling.value = true
    installError.value = null

    try {
      await WailsService.installMultipleIPAs(targetUdid, paths)
    } catch (err: any) {
      installError.value = err?.message || String(err)
      isInstalling.value = false
      throw err
    }
  }

  function closeInstallModal() {
    isInstalling.value = false
    installError.value = null
  }


  async function uninstallApp(bundleId: string, udid?: string) {
    const targetUdid = udid || selectedUdid.value
    if (!targetUdid) return
    try {
      await WailsService.uninstallApp(targetUdid, bundleId)
      installedApps.value = installedApps.value.filter((a) => a.bundleId !== bundleId)
    } catch (err: any) {
      throw err
    }
  }

  function initListeners() {
    WailsService.onEvent('device:connected', (devInfo: DeviceInfo) => {
      if (!devices.value.find(d => d.udid === devInfo.udid)) {
        devices.value.push(devInfo)
      }
      if (!selectedUdid.value) {
        selectedUdid.value = devInfo.udid
        fetchApps()
      }
    })

    WailsService.onEvent('device:updated', (devInfo: DeviceInfo) => {
      const idx = devices.value.findIndex(d => d.udid === devInfo.udid)
      if (idx !== -1) {
        devices.value[idx] = devInfo
      }
    })

    WailsService.onEvent('device:disconnected', (udid: string) => {
      devices.value = devices.value.filter((d) => d.udid !== udid)
      if (selectedUdid.value === udid) {
        selectedUdid.value = devices.value.length > 0 ? devices.value[0].udid : ''
        installedApps.value = []
        if (selectedUdid.value) fetchApps()
      }
    })

    WailsService.onEvent('device:install_progress', (prog: DeviceInstallTask) => {
      const idx = installTasks.value.findIndex(t => t.id === prog.id)
      if (idx !== -1) {
        installTasks.value[idx] = prog
      } else {
        installTasks.value.unshift(prog)
      }
    })


    WailsService.onEvent('device:install_complete', (data: any) => {
      if (data.udid === selectedUdid.value) {
        fetchApps()
      }

      const idx = installTasks.value.findIndex(t => t.id === data.id)
      if (idx !== -1) {
        installTasks.value[idx].phase = 'Complete'
        installTasks.value[idx].percent = 100
        installTasks.value[idx].message = 'Installed successfully'

        // Auto-remove completed tasks after 5 seconds
        setTimeout(() => {
          removeInstallTask(data.id)
        }, 5000)
      }
    })

    WailsService.onEvent('device:install_failed', (data: any) => {
      const idx = installTasks.value.findIndex(t => t.id === data.id)
      if (idx !== -1) {
        installTasks.value[idx].phase = 'Failed'
        installTasks.value[idx].message = data.error
      }
    })
  }

  function removeInstallTask(id: string) {
    installTasks.value = installTasks.value.filter(t => t.id !== id)
  }

  return {
    devices,
    selectedUdid,
    selectedDevice,
    installedApps,
    isLoadingApps,
    installTasks,
    activeTab,
    isInstalling,
    installError,
    isConnected,
    userApps,
    systemApps,
    checkDevices,
    fetchApps,
    pairDevice,
    installIPA,
    installMultipleIPAs,
    uninstallApp,
    removeInstallTask,
    closeInstallModal,
    initListeners,
  }
})
