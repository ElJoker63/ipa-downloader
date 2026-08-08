# IPA Downloader

IPA Downloader is a cross-platform desktop application and command-line suite designed to search, inspect, and download signed `.ipa` packages directly from Apple's App Store with FairPlay SINF DRM signature replication.

## Overview

The project provides both a modern desktop graphical user interface and a full-featured CLI tool sharing a unified backend engine written in Go.

### Core Capabilities

- **Direct App Store Integration**: Authenticate with Apple ID (including 2FA verification) to acquire free licenses and download original encrypted `.ipa` binaries.
- **FairPlay DRM SINF Replication**: Automatically replicates and injects FairPlay DRM signatures into downloaded packages for sideloading and analysis.
- **Search and Version Inspection**: Real-time App Store search across iOS, iPadOS, and tvOS with high-resolution artwork, screenshot lightboxes, and historical build listings.
- **Concurrent Transfer Queue**: Chunked streaming downloads with live speed tracking, ETA calculations, pause/resume controls, and retry handling.
- **Multi-Language Support**: Complete interface localization in English and Spanish (Español).
- **Offline Persistence**: Embedded pure-Go SQLite storage for favorites, transfer queues, search history, and settings without external database dependencies.
- **Diagnostic Logging**: Live streaming logs with severity filtering (INFO, SUCCESS, WARN, ERROR) and one-click file export.
- **CLI Compatibility**: Full command-line interface preserving backward compatibility with scripts and automated workflows.

---

## Tech Stack

### Backend (Go)
- **Framework**: Wails v2 for native desktop window management and Go-to-TypeScript runtime bindings.
- **App Store Engine**: Custom iTunes Storefront and GrandSlam protocol implementations.
- **Storage**: Zero-CGO SQLite driver (`modernc.org/sqlite`) for local persistence.
- **Security**: OS Keychain integration via `github.com/byteness/keyring` with persistent cookie jar.

### Frontend (TypeScript / Vue 3)
- **Framework**: Vue 3 (Composition API / `<script setup>`).
- **State Management**: Pinia stores for authentication, downloads queue, favorites, history, settings, and logs.
- **Styling**: TailwindCSS with Apple / visionOS glassmorphism design tokens and San Francisco (SF Pro) typography.
- **Routing**: Vue Router 4.

---

## Building from Source

### Prerequisites

- **Go**: Version 1.21 or higher
- **Node.js**: Version 18 or higher (with npm)
- **Wails CLI** (optional for live development):
  ```shell
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

### Build Desktop Application

To compile the native desktop executable:

```shell
# Build frontend assets
cd frontend
npm install
npm run build
cd ..

# Compile desktop binary
go build -o ipa-downloader-desktop.exe .
```

### Build CLI Tool

To compile the command-line binary:

```shell
go build -o ipa-downloader.exe main.go
```

### Development Mode

Run live development with hot reload:

```shell
wails dev
```

---

## CLI Reference

### Authentication

```shell
# Authenticate with Apple ID
ipa-downloader auth login --email "name@icloud.com" --password "secret"

# Check current authentication status
ipa-downloader auth info

# Revoke credentials
ipa-downloader auth revoke
```

### Search

```shell
# Search for apps by name or bundle ID
ipa-downloader search "Telegram" --limit 10 --platform iphone
```

### License Purchase

```shell
# Acquire a license for a free application
ipa-downloader purchase --bundle-identifier "ph.telegra.Telegraph"
```

### Version Listing

```shell
# List all historical version build identifiers available from Apple
ipa-downloader list-versions --bundle-identifier "ph.telegra.Telegraph"
```

### Download

```shell
# Download latest version
ipa-downloader download --bundle-identifier "ph.telegra.Telegraph" --output "./Telegram.ipa"

# Download a specific historical build
ipa-downloader download --bundle-identifier "ph.telegra.Telegraph" --external-version-id "854000123" --output "./Telegram_v10.ipa"
```

---

## Testing

Run unit tests and verification across packages:

```shell
# Run Go unit tests
go test ./...

# Verify frontend types and production bundle
cd frontend
npm run build
```

---

## License

This project is released under the [MIT License](LICENSE).
