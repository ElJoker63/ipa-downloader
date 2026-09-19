# Search & Download

---

## Search Engine

The search module interfaces with iTunes and App Store APIs:

* **Search by Keyword**: Type search terms to find relevant applications.
* **Search by Bundle Identifier**: If you know the exact Bundle ID (e.g., `org.videolan.vlc-ios`), search directly for pinpoint accuracy.
* **Platform Filters**:
  * iPhone
  * iPad
  * Apple TV (tvOS)
  * Apple Vision Pro (visionOS)
  * Mac App Store (macOS)

---

## Detailed App Inspector

Clicking any search result card displays full application details:

1. **General Metadata**: App title, developer name, category, star rating, price, download size, and release date.
2. **Screenshot Gallery**: High-resolution screenshots provided by the developer for iPhone and iPad form factors.
3. **Historical Version Selector**:
   * IPA Downloader queries Apple's server records for all historical version IDs (`externalVersionId`) registered for that title.
   * Select any previous build to download older releases instead of only the latest available version.

---

## Concurrent Download Manager

* **Chunked Streaming**: Optimized for large package downloads with resilient connection handling.
* **Live Metrics**: Moving average throughput speed, precise percentage, bytes transferred, and dynamic ETA calculation.
* **Flow Control**: Pause, resume, or cancel active transfers at any time.
* **Automatic Recovery**: If Apple CDN endpoints encounter temporary throttling or expired links, the download manager initiates license refresh and resume requests without losing existing progress.

---

## FairPlay SINF Signature Injection

Applications downloaded from the official App Store contain binaries encrypted with Apple's FairPlay DRM. In order for iOS devices to decrypt and run them, the package must contain a valid signature structure (`SC_Info/*.sinf`).

IPA Downloader:
1. Obtains license metadata returned by Apple's `buyProduct` and `downloadProduct` endpoints.
2. Generates the necessary `sinf` data blocks associated with the authenticated account.
3. Patches and injects the `SC_Info` folder directly into the `.ipa` ZIP archive.
4. Produces a fully authentic `.ipa` package ready for installation or archiving.
