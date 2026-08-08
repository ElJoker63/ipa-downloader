export type Platform = 'ios' | 'ipados' | 'tvos'

export interface AccountProfile {
  name: string
  email: string
  storeFront: string
  storeFrontCountry: string
  directoryServicesId: string
  pod: string
  isLoggedIn: boolean
}

export interface AppMetadata {
  id: number
  bundleId: string
  name: string
  developer: string
  developerId?: number
  version: string
  buildNumber?: string
  price: number
  formattedPrice: string
  currency?: string
  artworkUrl: string
  artworkUrl512?: string
  screenshots: string[]
  ipadScreenshots?: string[]
  description: string
  releaseNotes?: string
  releaseDate: string
  currentVersionDate?: string
  minimumOsVersion: string
  fileSizeBytes: number
  formattedSize: string
  averageUserRating: number
  userRatingCount: number
  contentAdvisoryRating: string
  genres: string[]
  primaryGenre: string
  supportedPlatforms: string[]
  isFavorite: boolean
}

export type DownloadStatus = 'queued' | 'downloading' | 'paused' | 'completed' | 'cancelled' | 'failed' | string

export interface DownloadTask {
  id: string
  appId: number
  bundleId: string
  appName: string
  developer: string
  version: string
  artworkUrl: string
  destinationPath: string
  status: DownloadStatus
  totalBytes: number
  downloadedBytes: number
  progress: number
  speedBytesPerSec: number
  formattedSpeed: string
  etaSeconds: number
  formattedETA: string
  externalVersionId?: string
  platform: string
  error?: string
  createdAt: string
  updatedAt: string
  completedAt?: string
}

export interface FavoriteApp {
  appId: number
  bundleId: string
  name: string
  developer: string
  version: string
  price: number
  formattedPrice: string
  artworkUrl: string
  primaryGenre?: string
  createdAt: string
}

export interface SearchHistoryItem {
  id: number
  term: string
  platform: string
  count: number
  createdAt: string
}

export interface AppSettings {
  theme: 'dark' | 'light' | 'system'
  language: string
  defaultDownloadFolder: string
  maxConcurrentDownloads: number
  autoCheckUpdates: boolean
  autoAcquireLicense: boolean
  rememberCredentials: boolean
  keychainPassphrase?: string
  searchLimit: number
}

export interface LogEntry {
  id: number
  timestamp: string
  level: 'DEBUG' | 'INFO' | 'WARN' | 'ERROR' | 'SUCCESS'
  message: string
  context?: string
}

export interface VersionInfo {
  externalVersionId: string
  displayVersion: string
  releaseDate: string
  formattedDate: string
}

export interface AppDetailsOutput {
  metadata: AppMetadata
  versionHistory: VersionInfo[]
  isFavorite: boolean
}
