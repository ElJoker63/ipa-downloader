# Local Library & Account Purchases

---

## Local Library

The **Library** tab maintains a persistent database of all `.ipa` and `.pkg` packages downloaded through the application:

* **Quick File Access**: Reveal the package in Finder (macOS) or File Explorer (Windows/Linux) with one click.
* **Direct USB Deployment**: If an iOS device is connected, install any app in your library with a single action.
* **Favorites**: Pin essential apps you frequently reinstall or archive.
* **Local Search**: Filter your collection by app name, version number, or download date.
* **Secure Storage**: All metadata is stored locally in an embedded zero-CGO SQLite database (`modernc.org/sqlite`). No usage metrics or history are sent to external servers.

---

## Account Purchases

The **Purchases** tab queries the historical license register for your Apple ID on Apple's servers:

* **Historical Archive**: Discover applications you purchased or claimed years ago, even if no longer prominently featured in the store.
* **Platform Categorization**: Switch easily between iPhone, iPad, Mac, and Apple TV purchases.
* **Instant Re-download**: Download any previously acquired application without searching for it in the general store catalog.
* **Local Caching**: Purchased apps are cached in the background database for instantaneous loading when opening the app.
