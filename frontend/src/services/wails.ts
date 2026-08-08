import type {
  AccountProfile,
  AppMetadata,
  AppDetailsOutput,
  DownloadTask,
  FavoriteApp,
  SearchHistoryItem,
  AppSettings,
  LogEntry,
} from '../types'

import * as AppService from '../../wailsjs/go/services/AppService'
import * as WailsRuntime from '../../wailsjs/runtime/runtime'

export const WailsService = {
  // Auth
  async getAccount(): Promise<AccountProfile> {
    try {
      return await AppService.GetAccount()
    } catch {
      return { name: '', email: '', storeFront: '143441-1,29', storeFrontCountry: 'US', directoryServicesId: '', pod: '', isLoggedIn: false }
    }
  },

  async login(email: string, pass: string, code: string = '', remember: boolean = true): Promise<AccountProfile> {
    return await AppService.Login(email, pass, code, remember)
  },

  async logout(): Promise<void> {
    return await AppService.Logout()
  },

  async getAuthStatus(): Promise<string> {
    try {
      return await AppService.GetAuthStatus()
    } catch {
      return 'Not Connected'
    }
  },

  // Search
  async searchApps(term: string, platform: string = 'ios', limit: number = 15): Promise<AppMetadata[]> {
    try {
      const results = await AppService.SearchApps(term, platform, limit)
      return results || []
    } catch (err) {
      console.error('searchApps error:', err)
      return []
    }
  },

  async lookupApp(bundleId: string, platform: string = 'ios'): Promise<AppMetadata> {
    return await AppService.LookupApp(bundleId, platform)
  },

  async getAppDetails(appId: number, bundleId: string, platform: string = 'ios'): Promise<AppDetailsOutput> {
    return await AppService.GetAppDetails(appId, bundleId, platform)
  },

  async getSearchHistory(limit: number = 15): Promise<SearchHistoryItem[]> {
    try {
      const history = await AppService.GetSearchHistory(limit)
      return history || []
    } catch {
      return []
    }
  },

  async clearSearchHistory(): Promise<void> {
    return await AppService.ClearSearchHistory()
  },

  // Downloads
  async queueDownload(app: AppMetadata, platform: string = 'ios', externalVersionId: string = '', customPath: string = ''): Promise<DownloadTask> {
    return await AppService.QueueDownload(app as any, platform, externalVersionId, customPath)
  },

  async pauseDownload(id: string): Promise<void> {
    return await AppService.PauseDownload(id)
  },

  async resumeDownload(id: string): Promise<void> {
    return await AppService.ResumeDownload(id)
  },

  async cancelDownload(id: string): Promise<void> {
    return await AppService.CancelDownload(id)
  },

  async retryDownload(id: string): Promise<void> {
    return await AppService.RetryDownload(id)
  },

  async getActiveDownloads(): Promise<DownloadTask[]> {
    try {
      const active = await AppService.GetActiveDownloads()
      return (active || []) as DownloadTask[]
    } catch {
      return []
    }
  },

  async getAllDownloads(): Promise<DownloadTask[]> {
    try {
      const all = await AppService.GetAllDownloads()
      return (all || []) as DownloadTask[]
    } catch {
      return []
    }
  },

  async clearCompletedDownloads(): Promise<void> {
    return await AppService.ClearCompletedDownloads()
  },

  // Favorites
  async getFavorites(): Promise<FavoriteApp[]> {
    try {
      const favs = await AppService.GetFavorites()
      return (favs || []) as FavoriteApp[]
    } catch {
      return []
    }
  },

  async searchFavorites(query: string): Promise<FavoriteApp[]> {
    try {
      const results = await AppService.SearchFavorites(query)
      return (results || []) as FavoriteApp[]
    } catch {
      return []
    }
  },

  async addFavorite(app: FavoriteApp): Promise<void> {
    return await AppService.AddFavorite(app as any)
  },

  async removeFavorite(appId: number): Promise<void> {
    return await AppService.RemoveFavorite(appId)
  },

  async toggleFavorite(app: FavoriteApp): Promise<boolean> {
    return await AppService.ToggleFavorite(app as any)
  },

  // History & Files
  async getDownloadHistory(): Promise<DownloadTask[]> {
    try {
      const hist = await AppService.GetDownloadHistory()
      return (hist || []) as DownloadTask[]
    } catch {
      return []
    }
  },

  async deleteHistoryItem(id: string): Promise<void> {
    return await AppService.DeleteHistoryItem(id)
  },

  async clearDownloadHistory(): Promise<void> {
    return await AppService.ClearDownloadHistory()
  },

  async openFolder(path: string): Promise<void> {
    return await AppService.OpenFolder(path)
  },

  async openFile(path: string): Promise<void> {
    return await AppService.OpenFile(path)
  },

  async revealInExplorer(path: string): Promise<void> {
    return await AppService.RevealInExplorer(path)
  },

  // Settings
  async getSettings(): Promise<AppSettings> {
    try {
      const s = await AppService.GetSettings()
      return s as AppSettings
    } catch {
      return {
        theme: 'system',
        language: 'en',
        defaultDownloadFolder: '',
        maxConcurrentDownloads: 3,
        autoCheckUpdates: true,
        autoAcquireLicense: true,
        rememberCredentials: true,
        searchLimit: 15,
      }
    }
  },

  async saveSettings(settings: AppSettings): Promise<void> {
    return await AppService.SaveSettings(settings as any)
  },

  async selectDownloadDirectory(defaultPath: string = ''): Promise<string> {
    return await AppService.SelectDownloadDirectory(defaultPath)
  },

  async clearAppCache(): Promise<void> {
    return await AppService.ClearAppCache()
  },

  async getCacheSize(): Promise<string> {
    try {
      return await AppService.GetCacheSize()
    } catch {
      return '0 KB'
    }
  },

  // Logs
  async getLogs(limit: number = 100): Promise<LogEntry[]> {
    try {
      const logs = await AppService.GetLogs(limit)
      return (logs || []) as LogEntry[]
    } catch {
      return []
    }
  },

  async clearLogs(): Promise<void> {
    return await AppService.ClearLogs()
  },

  async exportLogs(path: string = ''): Promise<string> {
    return await AppService.ExportLogs(path)
  },

  async addLog(level: string, message: string, context: string = ''): Promise<LogEntry | null> {
    return (await AppService.AddLog(level, message, context)) as LogEntry
  },

  // Window Controls
  minimizeWindow() {
    try {
      WailsRuntime.WindowMinimise()
    } catch {
      // ignore
    }
  },

  toggleMaximizeWindow() {
    try {
      WailsRuntime.WindowToggleMaximise()
    } catch {
      // ignore
    }
  },

  closeWindow() {
    try {
      WailsRuntime.Quit()
    } catch {
      // ignore
    }
  },

  // Event Subscription
  onEvent(eventName: string, callback: (...data: any[]) => void): () => void {
    try {
      return WailsRuntime.EventsOn(eventName, callback)
    } catch {
      return () => {}
    }
  },
}
