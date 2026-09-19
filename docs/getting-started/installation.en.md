# Installation Guide

**IPA Downloader** is available for macOS, Windows, and Linux via prebuilt binaries in [GitHub Releases](https://github.com/ElJoker63/ipa-downloader/releases/latest).

---

## macOS

### Requirements
* macOS 12 (Monterey) or higher (compatible with both Intel and Apple Silicon M1/M2/M3/M4).

### Installation Steps

1. Download the `ipa-downloader-*-macos.zip` package from the **Releases** page.
2. Unzip the file to reveal `ipa-downloader-desktop.app`.
3. Move `ipa-downloader-desktop.app` to your `/Applications` directory.

!!! warning "macOS Security Notice (Gatekeeper & Quarantine)"
    Because prebuilt open-source binaries are not signed with a paid Apple Developer ID certificate, macOS automatically applies the `com.apple.quarantine` extended attribute to files downloaded via a web browser.

    To allow the app to launch without security warnings, open your Terminal and run:
    ```bash
    xattr -cr /Applications/ipa-downloader-desktop.app
    ```
    After running this command, launch the app normally via Spotlight or Launchpad. See the [macOS Troubleshooting Guide](../troubleshooting/macos.md) for more details.

---

## Windows

### Requirements
* Windows 10 (version 1809 or later) or Windows 11 (64-bit).
* WebView2 Runtime (pre-installed on most modern Windows systems).

### Installation Steps
1. Download `ipa-downloader-*-windows.exe` from **Releases**.
2. Place the executable in your preferred folder (e.g., `C:\Program Files\IPA Downloader` or `Downloads`).
3. Double-click the `.exe` to start the app.
4. If Windows SmartScreen displays an unknown publisher prompt, click **More info** -> **Run anyway**.

---

## Linux

### Requirements
* Debian/Ubuntu (20.04+), Arch Linux, Fedora, or OpenSUSE.
* `libgtk-3` and `webkit2gtk-4.0` libraries.

### Installing Dependencies
Ubuntu / Debian:
```bash
sudo apt-get update
sudo apt-get install -y libgtk-3-0 libwebkit2gtk-4.0-37
```

Arch Linux:
```bash
sudo pacman -S gtk3 webkit2gtk
```

### Running the App
1. Download the `ipa-downloader-*-linux` binary.
2. Make it executable:
   ```bash
   chmod +x ipa-downloader-*-linux
   ```
3. Run the binary:
   ```bash
   ./ipa-downloader-*-linux
   ```
