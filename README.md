# IPA Downloader

[![Documentation](https://img.shields.io/badge/docs-GitHub_Pages-blue.svg)](https://eljoker63.github.io/ipa-downloader/)
[![Release](https://img.shields.io/github/v/release/ElJoker63/ipa-downloader)](https://github.com/ElJoker63/ipa-downloader/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**IPA Downloader** is a desktop application by **UDYAT** for searching, downloading, and managing iOS/iPadOS/tvOS/visionOS/macOS app packages from the App Store using your own Apple ID, with direct USB sideloading to connected Apple devices and IPSW firmware browsing — all in one native app, with a CLI mode for scripting.

> **Documentación Completa**: Consulta la documentación oficial y guías paso a paso en **[https://eljoker63.github.io/ipa-downloader/](https://eljoker63.github.io/ipa-downloader/)**.

## Overview

The app talks directly to Apple's App Store protocol (search, licensing, and download) and wraps it in a modern desktop GUI built with Wails, Vue 3, and TailwindCSS, sharing a single Go backend with a full command-line interface. Everything happens locally: your Apple ID credentials, session, and download history stay on your machine.

### Core Capabilities

- **App Store Search & Download**: Search by name or bundle ID across iOS, iPadOS, tvOS, visionOS, and macOS, with artwork, screenshots, and historical build listings, and download the original signed `.ipa`/`.pkg` package.
- **Apple ID Authentication**: Sign in with 2FA support and SAP-based signing, with credentials kept in the OS keychain.
- **Your Purchases**: Browse and re-download apps already owned by the authenticated account, with platform filtering.
- **Device Manager**: Pair a USB-connected iPhone/iPad, browse its installed apps, install or uninstall an `.ipa` directly on the device, and validate an `.ipa` before installing it — no Xcode or third-party sideloading tool required.
- **Local Library**: Keep track of downloaded apps, mark favorites, and review download history.
- **Firmware Browser**: Look up available IPSW firmware builds per device model.
- **Concurrent Download Queue**: Chunked streaming downloads with live speed, ETA, pause/resume, and automatic retry/recovery when Apple's download endpoints are flaky.
- **FairPlay DRM SINF Replication**: Automatically replicates and injects the FairPlay signature into downloaded packages.
- **Auto-Update**: Checks GitHub releases and applies updates in place.
- **Multi-Language UI**: Full interface localization in English and Spanish.
- **Live Logs**: Streaming log panel with severity filtering and one-click export, useful for diagnosing failed downloads.
- **CLI Mode**: The same binary runs as a scriptable CLI (`auth`, `search`, `purchase`, `list-purchases`, `list-versions`, `download`, ...) when invoked with arguments, for automation.

---

## Tech Stack

### Backend (Go)
- **Framework**: Wails v2 for the native desktop window and Go↔TypeScript bindings.
- **App Store Engine**: Custom iTunes Storefront / GrandSlam protocol implementation with SAP signing.
- **Device Communication**: `github.com/electricbubble/gidevice` for USB device pairing and app install/uninstall.
- **Storage**: Zero-CGO SQLite (`modernc.org/sqlite`) for local persistence (library, favorites, history, settings).
- **Security**: OS keychain integration with a persistent cookie jar for the Apple ID session.

### Frontend (TypeScript / Vue 3)
- **Framework**: Vue 3 (Composition API, `<script setup>`) with Vue Router 4.
- **State Management**: Pinia stores — auth, search, downloads, downloaded apps, purchases, favorites, history, settings, logs, device.
- **Styling**: TailwindCSS with a glassmorphism design and SF Pro typography.

---

## Building from Source

### Prerequisites

- **Go**: 1.25 or higher
- **Node.js**: 18 or higher (with npm)
- **Wails CLI**:
  ```shell
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### Build the Desktop Application

```shell
wails build
```

This builds the frontend and produces the native desktop binary in `build/bin`.

### Build the CLI Binary Manually

```shell
cd frontend && npm install && npm run build && cd ..
go build -o ipa-downloader.exe .
```

### Development Mode

```shell
wails dev
```

Runs the desktop app with Vite hot module replacement for the frontend.

### Run in CLI Mode

Passing any argument runs the CLI instead of launching the GUI:

```shell
./ipa-downloader.exe auth login -e "user@icloud.com"
./ipa-downloader.exe search "Spotify"
./ipa-downloader.exe download -b "com.spotify.client" --purchase
```

---

## CLI Reference

### Authentication

```shell
ipa-downloader auth login --email "name@icloud.com" --password "secret"
ipa-downloader auth info
ipa-downloader auth revoke
```

### Search

```shell
# platform: iphone, ipad, appletv, visionos, macos
ipa-downloader search "Telegram" --limit 10 --platform iphone
```

### License Purchase

```shell
ipa-downloader purchase --bundle-identifier "ph.telegra.Telegraph"
```

### List Purchases

```shell
ipa-downloader list-purchases --max-results 10 --page 1
```

### Version Listing

```shell
ipa-downloader list-versions --bundle-identifier "ph.telegra.Telegraph"
```

### Download

```shell
# Latest version
ipa-downloader download --bundle-identifier "ph.telegra.Telegraph" --output "./Telegram.ipa"

# A specific historical build
ipa-downloader download --bundle-identifier "ph.telegra.Telegraph" --external-version-id "854000123" --output "./Telegram_v10.ipa"

# macOS package
ipa-downloader download --bundle-identifier "ph.telegra.Telegraph" --platform macos --output "./Telegram.pkg"
```

---

## Testing

```shell
go generate ./...
go test -v ./...

cd frontend
npm run build
```

---

## License

This project is released under the [MIT License](LICENSE).

Built on top of the App Store protocol work from [ipatool](https://github.com/majd/ipatool).
