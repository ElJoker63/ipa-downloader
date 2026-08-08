import { ref, computed } from 'vue'
import { es } from './es'
import { en } from './en'
import type { Translations } from './types'

export type LanguageCode = 'es' | 'en'

const currentLanguage = ref<LanguageCode>('es') // Default Spanish as requested or loaded from settings

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

  function initLanguage(initial?: string) {
    if (initial === 'es' || initial === 'en') {
      currentLanguage.value = initial
      return
    }
    try {
      const saved = localStorage.getItem('ipa_downloader_lang')
      if (saved === 'es' || saved === 'en') {
        currentLanguage.value = saved
      } else {
        const browserLang = navigator.language.toLowerCase()
        if (browserLang.startsWith('es')) {
          currentLanguage.value = 'es'
        } else {
          currentLanguage.value = 'en'
        }
      }
    } catch {
      currentLanguage.value = 'es'
    }
  }

  return {
    t,
    currentLanguage,
    setLanguage,
    initLanguage,
  }
}
