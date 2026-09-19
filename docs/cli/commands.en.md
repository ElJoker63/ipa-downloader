# CLI Commands Reference

Syntax, options, and usage examples for every subcommand available in CLI mode.

---

## 1. Authentication (`auth`)

Manage sessions and credentials with the App Store.

### `auth login`
Authenticates with your Apple ID account.

```bash
ipa-downloader auth login [flags]
```

* `--email`, `-e`: Apple ID email address.
* `--password`, `-p`: Apple ID password.
* `--code`, `-c`: 6-digit two-factor authentication (2FA) verification code.

**Interactive Example:**
```bash
ipa-downloader auth login -e "user@icloud.com"
```

**Non-Interactive Single-Line Example:**
```bash
ipa-downloader auth login -e "user@icloud.com" -p "SecretPassword" -c "123456"
```

### `auth info`
Displays current session status and authenticated account details.

```bash
ipa-downloader auth info
```

### `auth revoke`
Signs out of the current session and purges saved tokens from the OS keychain.

```bash
ipa-downloader auth revoke
```

---

## 2. Search (`search`)

Searches the official App Store catalog.

```bash
ipa-downloader search <query> [flags]
```

* `--limit`, `-l`: Maximum number of results to return (default: `5`).
* `--platform`: Target platform (`iphone`, `ipad`, `appletv`, `visionos`, `macos`). Default: `iphone`.

**Examples:**
```bash
# Search iPhone apps
ipa-downloader search "WhatsApp" --limit 3

# Search macOS apps
ipa-downloader search "Final Cut" --platform macos --limit 1

# Export search results in JSON
ipa-downloader search "Slack" --format json
```

---

## 3. License Purchase (`purchase`)

Acquires a license for a free app on your Apple ID account (required if you have never downloaded it before).

```bash
ipa-downloader purchase --bundle-identifier <bundle-id>
```

* `--bundle-identifier`, `-b`: Application Bundle ID (e.g., `ph.telegra.Telegraph`).

**Example:**
```bash
ipa-downloader purchase -b "ph.telegra.Telegraph"
```

---

## 4. List Purchases (`list-purchases`)

Lists applications owned by the authenticated account.

```bash
ipa-downloader list-purchases [flags]
```

* `--max-results`: Number of results per page (default: `10`).
* `--page`: Page number for pagination.

**Example:**
```bash
ipa-downloader list-purchases --max-results 20 --page 1
```

---

## 5. Version History (`list-versions`)

Retrieves historical version identifiers (`externalVersionId`) registered on Apple servers for an app.

```bash
ipa-downloader list-versions --bundle-identifier <bundle-id>
```

**Example:**
```bash
ipa-downloader list-versions -b "com.spotify.client"
```

---

## 6. Version Metadata (`get-version-metadata`)

Retrieves detailed metadata for a specific historical build.

```bash
ipa-downloader get-version-metadata --bundle-identifier <bundle-id> --external-version-id <version-id>
```

---

## 7. Download (`download`)

Downloads an official `.ipa` or `.pkg` package from Apple servers.

```bash
ipa-downloader download [flags]
```

* `--bundle-identifier`, `-b`: Application Bundle ID.
* `--output`, `-o`: File path where the downloaded package will be saved.
* `--purchase`: Automatically acquires license if not already owned.
* `--external-version-id`: Specific historical version ID.
* `--platform`: Target package platform (`iphone`, `ipad`, `appletv`, `visionos`, `macos`).

**Examples:**
```bash
# Download latest version
ipa-downloader download -b "org.videolan.vlc-ios" --purchase

# Download specific historical version
ipa-downloader download -b "com.spotify.client" --external-version-id "854123456" -o "./Spotify_old.ipa"

# Download macOS package
ipa-downloader download -b "com.apple.Keynote" --platform macos -o "./Keynote.pkg"
```
