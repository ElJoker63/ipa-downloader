# IPA Downloader Desktop - Architecture & Developer Guide

## Overview

**IPA Downloader Desktop** is a modern, cross-platform desktop application and CLI suite that enables searching, querying, version-browsing, and downloading encrypted iOS and tvOS `.ipa` packages directly from Apple's App Store.

It is built with **Go 1.26+**, **Wails v2**, **Vue 3**, **TypeScript**, **TailwindCSS**, **Pinia**, and **pure-Go SQLite (`modernc.org/sqlite`)**.

---

## Clean Architecture Layers

```
ipa-downloader-desktop/
├── backend/
│   ├── app/         # Wails application lifecycle (Startup, DomReady, Shutdown, Native Dialogs)
│   ├── apple/       # Apple Storefront, Bag, Plist HTTP protocol wrappers
│   ├── auth/        # Apple ID authentication, 2FA code processing, keychain sessions
│   ├── config/      # Application preferences, download folders, cache management
│   ├── download/    # Concurrent download manager, FairPlay SINF signature injection
│   ├── events/      # Wails real-time event dispatcher (streaming logs, progress, auth state)
│   ├── library/     # Favorites, download history, and native file manager integrations
│   ├── models/      # Data structures shared across Go and TypeScript
│   ├── search/      # Search service with debounce, platform filters, metadata enrichment
│   ├── services/    # Unified AppService binding layer for Wails IPC
│   ├── storage/     # Thread-safe SQLite persistence layer
│   └── utils/       # Native file openers, clipboard, explorer integration
├── frontend/
│   ├── src/
│   │   ├── pages/       # Home, Search, Downloads, Favorites, History, Settings, Logs
│   │   ├── components/  # TwoFactorModal, AppDetailsModal, TitleBar, ToastContainer
│   │   ├── layouts/     # MainLayout with glass sidebar and title bar
│   │   ├── stores/      # Pinia stores (auth, search, downloads, favorites, history, settings, logs)
│   │   ├── composables/ # useTheme, useNotifications, useKeyboardShortcuts
│   │   └── types/       # TypeScript interfaces matching backend models
├── cmd/             # Cobra CLI frontend (100% backward compatible)
└── main.go          # Dual-mode launcher (Desktop GUI / CLI)
```

---

## Build & Development Commands

### 1. Run Development Mode
```bash
# Start Wails live development with Vite hot module replacement:
wails dev
```

### 2. Build Production Desktop Application
```bash
# Build desktop binary:
wails build

# Or directly using Go:
cd frontend && npm run build && cd ..
go build -o ipa-downloader-desktop.exe .
```

### 3. Run CLI Mode
```bash
# CLI commands continue to work seamlessly:
./ipa-downloader-desktop.exe auth login -e "user@icloud.com"
./ipa-downloader-desktop.exe search "Spotify"
./ipa-downloader-desktop.exe download -b "com.spotify.client" --purchase
```

---

## Key Features

1. **Dual Frontend Architecture**: Shared backend services power both the rich desktop GUI and the Cobra CLI.
2. **Zero-CGO SQLite Database**: Thread-safe storage with automated migrations using `modernc.org/sqlite`.
3. **Chunked Streaming & FairPlay SINF**: Streaming download with exponential moving average speed, pause/resume, cancel via context, and FairPlay signature patching.
4. **2FA Apple ID Verification**: Smooth in-app modal for 2-factor authentication.
5. **Real-Time Event Streaming**: Live download progress, background logs, and native desktop notifications via Wails events.
6. **Dark & Light Mode**: Tailored glassmorphic UI with automatic system preference detection.
