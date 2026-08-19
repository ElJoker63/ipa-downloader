import { defineStore } from 'pinia'
import { ref } from 'vue'
import type { DownloadedIPA } from '../types'
import { WailsService } from '../services/wails'

export const useDownloadedAppsStore = defineStore('downloadedApps', () => {
  const downloadedIPAs = ref<DownloadedIPA[]>([])
  const isLoading = ref(false)

  const isCheckingUpdates = ref(false)
  const storeMetadata = ref<Record<string, { latestVersion: string; appMetadata: any }>>({})

  async function fetchDownloadedIPAs(downloadDir: string = '') {
    isLoading.value = true
    try {
      const list = await WailsService.getDownloadedIPAs(downloadDir)
      downloadedIPAs.value = list || []
    } catch (err) {
      console.error('Failed to fetch downloaded IPAs:', err)
    } finally {
      isLoading.value = false
    }
  }

  async function checkUpdatesForDownloaded() {
    if (downloadedIPAs.value.length === 0) return
    isCheckingUpdates.value = true
    try {
      const promises = downloadedIPAs.value.map(async (ipa) => {
        if (!ipa.bundleId) return
        try {
          const appMeta = await WailsService.lookupApp(ipa.bundleId, 'ios')
          if (appMeta && appMeta.version) {
            storeMetadata.value[ipa.bundleId.toLowerCase()] = {
              latestVersion: appMeta.version,
              appMetadata: appMeta,
            }
          }
        } catch {
          // Ignore individual lookup errors
        }
      })
      await Promise.all(promises)
    } finally {
      isCheckingUpdates.value = false
    }
  }

  function getDownloadedByBundleId(bundleId: string): DownloadedIPA | undefined {
    if (!bundleId) return undefined
    return downloadedIPAs.value.find((item) => item.bundleId.toLowerCase() === bundleId.toLowerCase())
  }

  function compareVersions(v1: string, v2: string): number {
    const cleanV1 = (v1 || '').replace(/^v/i, '').trim()
    const cleanV2 = (v2 || '').replace(/^v/i, '').trim()
    if (cleanV1 === cleanV2) return 0

    const p1 = cleanV1.split('.').map((n) => parseInt(n, 10) || 0)
    const p2 = cleanV2.split('.').map((n) => parseInt(n, 10) || 0)
    const maxLen = Math.max(p1.length, p2.length)
    for (let i = 0; i < maxLen; i++) {
      const num1 = p1[i] || 0
      const num2 = p2[i] || 0
      if (num1 > num2) return 1
      if (num1 < num2) return -1
    }
    return 0
  }

  function isUpdateAvailable(bundleId: string, storeVersion: string): boolean {
    const localIPA = getDownloadedByBundleId(bundleId)
    if (!localIPA || !storeVersion || !localIPA.version) return false
    return compareVersions(storeVersion, localIPA.version) > 0
  }

  function getUpdateInfo(bundleId: string, localVersion: string) {
    if (!bundleId) return null
    const data = storeMetadata.value[bundleId.toLowerCase()]
    if (!data || !data.latestVersion) return null
    if (compareVersions(data.latestVersion, localVersion) > 0) {
      return data
    }
    return null
  }

  return {
    downloadedIPAs,
    isLoading,
    isCheckingUpdates,
    storeMetadata,
    fetchDownloadedIPAs,
    checkUpdatesForDownloaded,
    getDownloadedByBundleId,
    isUpdateAvailable,
    getUpdateInfo,
    compareVersions,
  }
})
