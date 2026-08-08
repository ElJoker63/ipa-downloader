import { onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { useSearchStore } from '../stores/search'

export function useKeyboardShortcuts() {
  const router = useRouter()
  const searchStore = useSearchStore()

  function handleKeyDown(e: KeyboardEvent) {
    // Check if Ctrl or Cmd key is pressed
    const isModifier = e.ctrlKey || e.metaKey

    if (isModifier) {
      switch (e.key) {
        case '1':
          e.preventDefault()
          router.push('/')
          break
        case '2':
          e.preventDefault()
          router.push('/search')
          break
        case '3':
          e.preventDefault()
          router.push('/downloads')
          break
        case '4':
          e.preventDefault()
          router.push('/favorites')
          break
        case '5':
          e.preventDefault()
          router.push('/history')
          break
        case '6':
          e.preventDefault()
          router.push('/settings')
          break
        case '7':
          e.preventDefault()
          router.push('/logs')
          break
        case 'k':
        case 'K':
        case 'f':
        case 'F':
          e.preventDefault()
          router.push('/search')
          setTimeout(() => {
            const searchInput = document.getElementById('main-search-input')
            if (searchInput) {
              searchInput.focus()
            }
          }, 50)
          break
        case ',':
          e.preventDefault()
          router.push('/settings')
          break
      }
    } else if (e.key === 'Escape') {
      if (searchStore.isDetailsModalOpen) {
        searchStore.closeAppDetails()
      }
    }
  }

  onMounted(() => {
    window.addEventListener('keydown', handleKeyDown)
  })

  onUnmounted(() => {
    window.removeEventListener('keydown', handleKeyDown)
  })
}
