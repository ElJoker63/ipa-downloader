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
  type: string
  url?: string
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
  checksum?: string
  checksumType?: string
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

export interface DeviceInfo {
  udid: string
  name: string
  model: string
  productType: string
  deviceClass: string
  iosVersion: string
  buildVersion: string
  serialNumber: string
  wifiAddress?: string
  imei?: string
  imei2?: string
  modelNumber?: string
  regionInfo?: string
  activationState?: string
  boardConfig?: string
  cpuArchitecture?: string
  hardwareModel?: string
  isJailbroken: boolean

  batteryHealth?: number
  chargeCycles?: number
  storageTotal: number
  storageUsed: number
  storageFree: number
  batteryLevel: number
  batteryCharging: boolean
  isPaired: boolean
  isConnected: boolean
}

export interface InstalledApp {
  name: string
  bundleId: string
  version: string
  shortVersion: string
  size: number
  dynamicSize?: number
  vendor?: string
  appType: string
  minimumOs?: string
  signerIdentity?: string
  artworkUrl?: string
}

export interface UpdateInfo {
  available: boolean
  latestVersion: string
  currentVersion: string
  releaseNotes: string
  downloadUrl: string
  mandatory: boolean
}


export interface IPAInfo {
  bundleId: string
  bundleName: string
  version: string
  shortVersion: string
  minimumOs: string
  architectures: string[]
  fileSizeBytes: number
  isValid: boolean
  error?: string
}

export interface DownloadedIPA {
  filePath: string
  fileName: string
  bundleId: string
  appName: string
  version: string
  shortVersion: string
  minimumOs: string
  fileSizeBytes: number
  formattedSize: string
  modTime: string
  artworkUrl?: string
  appId?: number
}

export interface DeviceInstallTask {
  id: string
  ipaName: string
  udid: string
  phase: 'Preparing' | 'Copying' | 'Installing' | 'Verifying' | 'Complete' | 'Failed' | string
  percent: number
  message: string
}



export interface Firmware {
  version: string
  buildid: string
  sha1sum: string
  md5sum: string
  size: number
  releasedate: string
  uploaddate: string
  url: string
  signed: boolean
  filename: string
}

export interface AppleHardware {
  identifier: string
  name: string
  platform: string
  firmwares?: Firmware[]
}

