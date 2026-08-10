import { defineStore } from 'pinia'
import { ref } from 'vue'

export type ModalType = 'confirm' | 'alert' | 'info' | 'error' | 'success'

interface ModalOptions {
  title: string
  message: string
  confirmText?: string
  cancelText?: string
  type?: ModalType
  onConfirm?: () => void | Promise<void>
  onCancel?: () => void
}

export const useModalStore = defineStore('modal', () => {
  const isOpen = ref(false)
  const options = ref<ModalOptions>({
    title: '',
    message: '',
    confirmText: 'Confirm',
    cancelText: 'Cancel',
    type: 'confirm',
  })

  const loading = ref(false)

  function show(opts: ModalOptions) {
    options.value = {
      confirmText: opts.confirmText ?? 'Confirm',
      cancelText: opts.cancelText ?? 'Cancel',
      type: opts.type ?? 'confirm',
      title: opts.title,
      message: opts.message,
      onConfirm: opts.onConfirm,
      onCancel: opts.onCancel,
    }
    isOpen.value = true
  }

  function confirm(title: string, message: string, onConfirm: () => void | Promise<void>, confirmText?: string) {
    show({
      title,
      message,
      onConfirm,
      ...(confirmText ? { confirmText } : {}),
      type: 'confirm'
    })
  }

  function alert(title: string, message: string, type: ModalType = 'error') {
    show({
      title,
      message,
      type,
      confirmText: 'OK',
      cancelText: '' // No cancel for alert
    })
  }

  async function handleConfirm() {
    if (options.value.onConfirm) {
      loading.value = true
      try {
        await options.value.onConfirm()
      } finally {
        loading.value = false
        isOpen.value = false
      }
    } else {
      isOpen.value = false
    }
  }

  function handleCancel() {
    if (options.value.onCancel) {
      options.value.onCancel()
    }
    isOpen.value = false
  }

  return {
    isOpen,
    options,
    loading,
    show,
    confirm,
    alert,
    handleConfirm,
    handleCancel
  }
})
