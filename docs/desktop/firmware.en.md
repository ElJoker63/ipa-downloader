# IPSW Firmware Browser

The application includes a built-in firmware explorer to look up official Apple restore images (`.ipsw`) across iPhone, iPad, iPod touch, Apple TV, and Apple Silicon Mac models.

---

## Capabilities

1. **Device Catalog**: Browse device families categorized by hardware generations and model identifiers.
2. **Version Details**:
   * Exact iOS, iPadOS, or macOS build numbers.
   * Official public release dates.
   * Installation image file sizes in gigabytes.
   * Official SHA-1 and MD5 checksums for verifying download integrity.
3. **Signing Status**:
   * Clearly flags whether Apple **is currently signing** that firmware version (allowing restore or downgrade via Finder/iTunes).
   * Highlights unsigned legacy versions that cannot be activated by Apple servers.
4. **Direct Download Links**: Verified direct URLs to Apple's official Content Delivery Networks (`appldnld.apple.com`).
