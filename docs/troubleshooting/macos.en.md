# macOS Troubleshooting

---

## 1. Error: «App is damaged and can’t be opened» or «Not compatible»

### Cause
When downloading `.zip` files or apps directly via web browsers such as Safari or Google Chrome, macOS automatically assigns the security extended attribute **`com.apple.quarantine`**.

Because open-source releases do not use an annual paid Apple Developer certificate ($99/year) for official Apple notarization, **Gatekeeper** prevents execution or shows a prohibition sign over the app.

### Solution
To clear the quarantine flag and open the application normally:

1. Open **Terminal** on your Mac.
2. Execute the following command (assuming you moved the application to `/Applications`):
   ```bash
   xattr -cr /Applications/ipa-downloader-desktop.app
   ```
3. Double-click the application again. It will launch immediately.

---

## 2. Error After In-App Update

If an older version failed during auto-update and refused to launch:
1. Ensure you have installed version **1.4.3 or higher**, which includes the fix that unpacks the Mach-O binary from the `.zip` archive on macOS.
2. If your app was damaged by an earlier release, download the clean `.zip` from [GitHub Releases](https://github.com/ElJoker63/ipa-downloader/releases/latest), replace the app bundle in `/Applications`, and run the `xattr -cr` command above.

---

## 3. USB Devices Not Detected

If your iPhone or iPad does not appear in the **Devices** tab:
1. Unlock your device screen.
2. If the **«Trust This Computer?»** prompt appears, tap **Trust** and enter your passcode.
3. Disconnect and reconnect the USB cable.
4. On macOS, ensure no virtual machines or other software have captured exclusive USB control.
