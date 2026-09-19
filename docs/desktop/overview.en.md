# Desktop Application - Overview

The **IPA Downloader** desktop application provides an intuitive, responsive, and modern interface for searching, downloading, signing, and sideloading Apple applications.

---

## User Interface

Built with **Vue 3**, **TailwindCSS**, and **Wails v2**, the desktop client features:

* **Glassmorphic Styling**: Translucent panels with subtle borders inspired by modern macOS design principles.
* **Dark & Light Themes**: Automatically respects system theme preferences or can be toggled manually in Settings.
* **Sidebar Navigation**: Instant access to primary application areas:
  * :material-magnify: **Search**: Global search across App Store catalogs.
  * :material-download: **Downloads**: Live download monitor and queue manager.
  * :material-cellphone: **Devices**: USB sideloading and app management for connected iPhones & iPads.
  * :material-bookshelf: **Library**: Local repository of downloaded packages and favorites.
  * :material-shopping: **Purchases**: Account license history across Apple platforms.
  * :material-chip: **Firmwares**: IPSW restore image explorer categorized by device model.
  * :material-cog: **Settings**: Directory preferences, download limits, and account sessions.
  * :material-clipboard-text: **Logs**: Live system log viewer with severity filtering.

---

## System Events & State

* **Real-Time Event Dispatching**: Download progress bars, authentication state changes, and completion notifications are emitted through Wails bidirectionally (`events.Emit`).
* **Desktop Notifications**: Native OS notification banners when long downloads finish or app installations complete.
* **Streaming Log Console**: Real-time diagnostic view of HTTP requests sent to Apple endpoints, network recovery attempts, and one-click log exports.
