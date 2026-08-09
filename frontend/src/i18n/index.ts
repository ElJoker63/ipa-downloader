import { ref, computed } from 'vue'
import { es } from './es'
import { en } from './en'
import type { Translations } from './types'

export type LanguageCode = 'es' | 'en'

const currentLanguage = ref<LanguageCode>('en') // Default English, will be updated by init

export function useI18n() {
  const t = computed<Translations>(() => {
    return currentLanguage.value === 'es' ? es : en
  })

  function setLanguage(lang: LanguageCode) {
    currentLanguage.value = lang
    try {
      localStorage.setItem('ipa_downloader_lang', lang)
    } catch {
      // ignore
    }
  }

  function getLanguage(): LanguageCode {
    return currentLanguage.value
  }

  function initLanguage(initial?: string) {
    if (initial === 'es' || initial === 'en') {
      currentLanguage.value = initial as LanguageCode
      return
    }
    try {
      const saved = localStorage.getItem('ipa_downloader_lang')
      if (saved === 'es' || saved === 'en') {
        currentLanguage.value = saved as LanguageCode
      } else {
        const browserLang = navigator.language.toLowerCase()
        if (browserLang.startsWith('es')) {
          currentLanguage.value = 'es'
        } else {
          currentLanguage.value = 'en'
        }
      }
    } catch {
      currentLanguage.value = 'en'
    }
  }

  return {
    t,
    currentLanguage,
    setLanguage,
    initLanguage,
    getLanguage,
  }
}
