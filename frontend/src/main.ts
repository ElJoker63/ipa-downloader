import { createApp } from 'vue'
import { createPinia } from 'pinia'
import { router } from './router'
import App from './App.vue'
import './index.css'

import { useAuthStore } from './stores/auth'
import { useDownloadsStore } from './stores/downloads'
import { useFavoritesStore } from './stores/favorites'
import { useLogsStore } from './stores/logs'
import { useSettingsStore } from './stores/settings'
import { useNotifications } from './composables/useNotifications'

const app = createApp(App)
const pinia = createPinia()

app.use(pinia)
app.use(router)

app.mount('#app')

// Initialize global event listeners
const authStore = useAuthStore()
const downloadsStore = useDownloadsStore()
const favoritesStore = useFavoritesStore()
const logsStore = useLogsStore()
const settingsStore = useSettingsStore()
const { initListeners: initNotificationListeners } = useNotifications()

authStore.initListeners()
downloadsStore.initListeners()
favoritesStore.initListeners()
logsStore.initListeners()
initNotificationListeners()

// Load initial settings and theme
settingsStore.fetchSettings()
authStore.checkAccount()
