import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { DownloadTask, AppMetadata } from '../types'
import { WailsService } from '../services/wails'

export const useDownloadsStore = defineStore('downloads', () => {
  const downloads = ref<DownloadTask[]>([])
  const isLoading = ref(false)

  // Active downloads queue (Must NEVER remove tasks while downloading, queued, paused, OR signing)
  const activeDownloads = computed(() => {
    return downloads.value.filter(
      (d) => d.status === 'downloading' || d.status === 'queued' || d.status === 'paused' || d.status === 'signing'
    )
  })

  // Completed or failed downloads
  const completedDownloads = computed(() => {
    return downloads.value.filter(
      (d) => d.status === 'completed' || d.status === 'failed' || d.status === 'cancelled'
    )
  })

  const activeCount = computed(() => {
    return downloads.value.filter(
      (d) => d.status === 'downloading' || d.status === 'queued' || d.status === 'signing'
    ).length
  })

  const totalSpeedFormatted = computed(() => {
    const sumBytes = downloads.value
      .filter((d) => d.status === 'downloading')
      .reduce((acc, d) => acc + (d.speedBytesPerSec || 0), 0)

    if (sumBytes <= 0) return '0 KB/s'
    const unit = 1024
    if (sumBytes < unit) return `${sumBytes} B/s`
    const div = 1024
    const exp = Math.floor(Math.log(sumBytes) / Math.log(div))
    const letter = 'KMGTPE'[exp - 1] || 'K'
    return `${(sumBytes / Math.pow(div, exp)).toFixed(1)} ${letter}B/s`
  })

  async function fetchDownloads() {
    isLoading.value = true
    try {
      const all = await WailsService.getAllDownloads()
      downloads.value = all || []
    } catch (err) {
      // ignore
    } finally {
      isLoading.value = false
    }
  }

  async function queueDownload(app: AppMetadata, platform: string = 'ios', externalVersionId: string = '', customPath: string = '') {
    try {
      const task = await WailsService.queueDownload(app, platform, externalVersionId, customPath)
      const existingIdx = downloads.value.findIndex((d) => d.id === task.id)
      if (existingIdx !== -1) {
        downloads.value[existingIdx] = task
      } else {
        downloads.value.unshift(task)
      }
      return task
    } catch (err) {
      throw err
    }
  }

  async function pauseDownload(id: string) {
    await WailsService.pauseDownload(id)
    const task = downloads.value.find((d) => d.id === id)
    if (task) task.status = 'paused'
  }

  async function resumeDownload(id: string) {
    await WailsService.resumeDownload(id)
    const task = downloads.value.find((d) => d.id === id)
    if (task) task.status = 'queued'
  }

  async function cancelDownload(id: string) {
    await WailsService.cancelDownload(id)
    const task = downloads.value.find((d) => d.id === id)
    if (task) task.status = 'cancelled'
  }

  async function removeTask(id: string) {
    await WailsService.deleteHistoryItem(id) // Reuse the same delete logic
    downloads.value = downloads.value.filter((d) => d.id !== id)
  }


  async function retryDownload(id: string) {
    await WailsService.retryDownload(id)
    const task = downloads.value.find((d) => d.id === id)
    if (task) {
      task.status = 'queued'
      task.progress = 0
      task.error = ''
    }
  }

  async function clearCompleted() {
    await WailsService.clearCompletedDownloads()
    downloads.value = downloads.value.filter(
      (d) => d.status === 'downloading' || d.status === 'queued' || d.status === 'paused' || d.status === 'signing'
    )
  }

  function initListeners() {
    WailsService.onEvent('download:progress', (updatedTask: DownloadTask) => {
      if (!updatedTask?.id) return
      const idx = downloads.value.findIndex((d) => d.id === updatedTask.id)
      if (idx !== -1) {
        downloads.value[idx] = { ...downloads.value[idx], ...updatedTask }
      } else {
        downloads.value.unshift(updatedTask)
      }
    })

    WailsService.onEvent('download:status', (updatedTask: DownloadTask) => {
      if (!updatedTask?.id) return
      const idx = downloads.value.findIndex((d) => d.id === updatedTask.id)
      if (idx !== -1) {
        downloads.value[idx] = { ...downloads.value[idx], ...updatedTask }
      } else {
        downloads.value.unshift(updatedTask)
      }
    })

    WailsService.onEvent('download:completed', (updatedTask: DownloadTask) => {
      if (!updatedTask?.id) return
      const idx = downloads.value.findIndex((d) => d.id === updatedTask.id)
      if (idx !== -1) {
        downloads.value[idx] = { ...downloads.value[idx], ...updatedTask, status: 'completed', progress: 100 }
      } else {
        downloads.value.unshift({ ...updatedTask, status: 'completed', progress: 100 })
      }
    })

    WailsService.onEvent('download:failed', (updatedTask: DownloadTask) => {
      if (!updatedTask?.id) return
      const idx = downloads.value.findIndex((d) => d.id === updatedTask.id)
      if (idx !== -1) {
        downloads.value[idx] = { ...downloads.value[idx], ...updatedTask, status: 'failed' }
      } else {
        downloads.value.unshift({ ...updatedTask, status: 'failed' })
      }
    })
  }

  return {
    downloads,
    activeDownloads,
    completedDownloads,
    activeCount,
    totalSpeedFormatted,
    isLoading,
    fetchDownloads,
    queueDownload,
    pauseDownload,
    resumeDownload,
    cancelDownload,
    removeTask,
    retryDownload,
    clearCompleted,
    initListeners,
  }
})

