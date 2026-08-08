import { ref } from 'vue'
import { WailsService } from '../services/wails'

export interface ToastNotification {
  id: string
  title: string
  message: string
  severity: 'info' | 'success' | 'warning' | 'error'
  timestamp: number
}

const toasts = ref<ToastNotification[]>([])

export function useNotifications() {
  function showToast(title: string, message: string, severity: 'info' | 'success' | 'warning' | 'error' = 'info', duration: number = 4000) {
    const id = Math.random().toString(36).substring(2, 9)
    const toast: ToastNotification = {
      id,
      title,
      message,
      severity,
      timestamp: Date.now(),
    }
    toasts.value.unshift(toast)

    if (duration > 0) {
      setTimeout(() => {
        dismissToast(id)
      }, duration)
    }
  }

  function dismissToast(id: string) {
    toasts.value = toasts.value.filter((t) => t.id !== id)
  }

  function initListeners() {
    WailsService.onEvent('notification:show', (payload: any) => {
      if (payload?.title && payload?.message) {
        showToast(payload.title, payload.message, payload.severity || 'info')
      }
    })
  }

  return {
    toasts,
    showToast,
    dismissToast,
    initListeners,
  }
}
