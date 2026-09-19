# Windows & Linux Troubleshooting

---

## Windows

### Windows Defender / SmartScreen Prompt
* **Symptom**: A blue dialog appears saying *«Windows protected your PC. Microsoft Defender SmartScreen prevented an unrecognized app from starting»*.
* **Solution**: Click **«More info»** and then click **«Run anyway»**. This is expected for applications without an expensive Extended Validation (EV) code signing certificate.

### Missing WebView2 Error
* **Symptom**: The application fails to launch or reports that WebView2 cannot be found.
* **Solution**: Wails uses Microsoft Edge WebView2 on Windows. Most modern Windows 10/11 machines already include it. If missing, download it from [Microsoft Edge WebView2 Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/).

---

## Linux

### Shared Library `libwebkit2gtk` Not Found
* **Symptom**: The binary exits with `error while loading shared libraries: libwebkit2gtk-4.0.so.37`.
* **Solution**: Install the required packages for your distribution:
  * **Ubuntu / Debian**:
    ```bash
    sudo apt-get update && sudo apt-get install -y libwebkit2gtk-4.0-37 libgtk-3-0
    ```
  * **Fedora**:
    ```bash
    sudo dnf install webkit2gtk3
    ```
  * **Arch Linux**:
    ```bash
    sudo pacman -S webkit2gtk gtk3
    ```

---

## Common App Store Errors

### Authentication Error (`BadCredentials` or `2FA Required`)
* Make sure you enter the 6-digit code promptly when the two-factor modal appears.
* If Apple requires web verification (such as accepting updated iCloud terms of service), log in once at [appleid.apple.com](https://appleid.apple.com) in your browser to accept any pending notices.

### Error: `This item is not available in your storefront`
* App Store licenses are region-locked. If you attempt to download an app exclusive to the US App Store using an Apple ID registered in Spain or Mexico, Apple will reject the request.
* Ensure your Apple ID account region matches the storefront where the application is available.
