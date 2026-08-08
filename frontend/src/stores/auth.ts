import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { AccountProfile } from '../types'
import { WailsService } from '../services/wails'

export const useAuthStore = defineStore('auth', () => {
  const account = ref<AccountProfile>({
    name: '',
    email: '',
    storeFront: '143441-1,29',
    storeFrontCountry: 'US',
    directoryServicesId: '',
    pod: '',
    isLoggedIn: false,
  })

  const status = ref<string>('Not Connected')
  const isLoading = ref(false)
  const errorMessage = ref<string | null>(null)
  const is2FAModalOpen = ref(false)

  const pendingLoginData = ref<{ email: string; pass: string; remember: boolean } | null>(null)

  const isLoggedIn = computed(() => account.value.isLoggedIn && status.value === 'Connected')

  async function checkAccount() {
    try {
      isLoading.value = true
      const acc = await WailsService.getAccount()
      account.value = acc
      if (acc.isLoggedIn) {
        status.value = 'Connected'
      } else {
        status.value = 'Not Connected'
      }
    } catch (err: any) {
      status.value = 'Not Connected'
    } finally {
      isLoading.value = false
    }
  }

  async function login(email: string, pass: string, authCode: string = '', remember: boolean = true) {
    isLoading.value = true
    errorMessage.value = null
    status.value = 'Connecting...'

    try {
      const acc = await WailsService.login(email, pass, authCode, remember)
      account.value = acc
      status.value = 'Connected'
      is2FAModalOpen.value = false
      pendingLoginData.value = null
      return acc
    } catch (err: any) {
      const msg = err?.message || String(err)
      if (msg.includes('2FA') || msg.includes('auth code') || msg.includes('code is required')) {
        pendingLoginData.value = { email, pass, remember }
        is2FAModalOpen.value = true
        errorMessage.value = 'Two-Factor Authentication required. Enter the code sent to your Apple device.'
      } else {
        errorMessage.value = msg
        status.value = 'Not Connected'
      }
      throw err
    } finally {
      isLoading.value = false
    }
  }

  async function submit2FACode(code: string) {
    if (!pendingLoginData.value) return
    return login(
      pendingLoginData.value.email,
      pendingLoginData.value.pass,
      code,
      pendingLoginData.value.remember
    )
  }

  async function logout() {
    isLoading.value = true
    try {
      await WailsService.logout()
      account.value = {
        name: '',
        email: '',
        storeFront: '143441-1,29',
        storeFrontCountry: 'US',
        directoryServicesId: '',
        pod: '',
        isLoggedIn: false,
      }
      status.value = 'Not Connected'
    } finally {
      isLoading.value = false
    }
  }

  function initListeners() {
    WailsService.onEvent('auth:status', (data: any) => {
      if (data?.status) {
        status.value = data.status
      }
    })

    WailsService.onEvent('auth:account', (acc: AccountProfile) => {
      if (acc) {
        account.value = acc
        status.value = acc.isLoggedIn ? 'Connected' : 'Not Connected'
      }
    })
  }

  return {
    account,
    status,
    isLoading,
    errorMessage,
    isLoggedIn,
    is2FAModalOpen,
    checkAccount,
    login,
    submit2FACode,
    logout,
    initListeners,
  }
})
