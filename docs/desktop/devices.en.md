# USB Device Manager & Sideloading

The device module in **IPA Downloader** communicates directly with USB-connected iPhones and iPads over the native `usbmuxd` protocol (using the pure Go `gidevice` implementation), without requiring Xcode, iTunes, or external sideloading utilities.

---

## Detection and Pairing

1. Connect your iOS or iPadOS device to your computer via USB.
2. If connecting for the first time, unlock your screen and tap **«Trust This Computer»**, entering your device passcode.
3. In the **Devices** tab of the application, the device appears automatically showing:
   * Device name (e.g., *"User's iPhone"*).
   * Hardware model identifier (e.g., `iPhone14,2`).
   * Installed iOS / iPadOS version.
   * UDID (Unique Device Identifier).
   * Pairing status.

---

## Direct Package Installation (.ipa Sideloading)

To install an app package onto your device:

1. Go to the **Devices** tab.
2. Select the connected device.
3. Click **Install IPA** and choose any `.ipa` file from your computer (or click the install action directly from your local Library).
4. The manager validates package structure and streams the application payload to the device's native installation service (`installation_proxy`).
5. The application icon appears on your home screen once installation finishes.

!!! note "License Compatibility"
    For an official App Store package to launch without crashing immediately, the Apple ID used to download the `.ipa` must match the account active on the iOS device, or the package must be resigned with a compatible developer provisioning profile.

---

## Managing Installed Applications

* **List Installed Apps**: View all user-installed applications along with their Bundle Identifiers and version numbers.
* **One-Click Uninstall**: Remove applications directly from your computer by clicking the trash icon next to any app.
* **Device Information**: Review battery health, available disk capacity, baseband version, and activation status.
