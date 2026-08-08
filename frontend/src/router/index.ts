import { createRouter, createWebHashHistory } from 'vue-router'
import Home from '../pages/Home.vue'
import Search from '../pages/Search.vue'
import Downloads from '../pages/Downloads.vue'
import Apps from '../pages/Apps.vue'
import Favorites from '../pages/Favorites.vue'
import History from '../pages/History.vue'
import Settings from '../pages/Settings.vue'
import Logs from '../pages/Logs.vue'

const routes = [
  { path: '/', name: 'Home', component: Home },
  { path: '/search', name: 'Search', component: Search },
  { path: '/downloads', name: 'Downloads', component: Downloads },
  { path: '/apps', name: 'Apps', component: Apps },
  { path: '/favorites', name: 'Favorites', component: Favorites },
  { path: '/history', name: 'History', component: History },
  { path: '/settings', name: 'Settings', component: Settings },
  { path: '/logs', name: 'Logs', component: Logs },
]

export const router = createRouter({
  history: createWebHashHistory(),
  routes,
})
