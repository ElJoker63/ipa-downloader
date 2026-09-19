# Quickstart Guide

Learn how to configure your Apple ID and perform your first search and download in minutes.

---

## 1. Sign In with Your Apple ID

Downloading authentic App Store packages requires an authenticated Apple ID session.

1. Launch **IPA Downloader**.
2. Navigate to the **Account** tab or click the profile button in the header.
3. Enter your **Apple ID email** and **password**.
4. Click **Sign In**.
5. If two-factor authentication (2FA) is enabled on your account, an in-app prompt will appear asking for the 6-digit verification code sent to your trusted Apple devices or phone number. Enter the code to complete sign-in.

!!! tip "Security of Your Credentials"
    Credentials are submitted directly to Apple's official authentication endpoint (`gsa.apple.com`). The resulting session token is saved encrypted in your operating system's native keychain. No third-party servers ever receive your data.

---

## 2. Search for an App

1. Click on the **Search** tab in the sidebar.
2. Enter an app name (e.g., `Telegram`, `WhatsApp`) or Bundle ID (e.g., `ph.telegra.Telegraph`).
3. Optionally filter by target platform:
   * **iOS (iPhone)**
   * **iPadOS**
   * **tvOS**
   * **visionOS**
   * **macOS**
4. Click on any search card to open its detailed modal with descriptions, screenshots, current version, file size, and historical version listings.

---

## 3. Download the Package

1. In the app details view, click the **Download** button.
2. If your account hasn't acquired a license for the application yet (even for free apps), the app will automatically request a license (`purchase`) from Apple.
3. The download starts immediately. You can track progress, throughput, and ETA in the **Downloads** tab.
4. When finished, the app automatically replicates and injects the FairPlay SINF signature for your account into the archive.

---

## 4. Quick CLI Usage

If you prefer using the command line:

```bash
# Sign in
./ipa-downloader auth login --email "your-email@icloud.com"

# Search for apps
./ipa-downloader search "Spotify" --limit 5

# Download by Bundle ID
./ipa-downloader download -b "com.spotify.client" --purchase
```
