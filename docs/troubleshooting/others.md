# Solución de Problemas en Windows y Linux

---

## Windows

### Aviso de Windows Defender / SmartScreen
* **Síntoma**: Windows muestra una ventana azul con el mensaje *«Windows protegió su PC. Microsoft Defender SmartScreen impidió el inicio de una aplicación no reconocida»*.
* **Solución**: Haz clic en el enlace **«Más información»** y posteriormente en el botón **«Ejecutar de todas formas»**. Esto ocurre porque el ejecutable no tiene un certificado EV de firma de código de pago.

### Error de WebView2 ausente
* **Síntoma**: La aplicación no abre o muestra un error relacionado con WebView2.
* **Solución**: Wails utiliza el runtime de Microsoft Edge WebView2 en Windows. La mayoría de versiones de Windows 10 y 11 ya lo incluyen. Si falta, descárgalo directamente desde la página oficial de [Microsoft Edge WebView2 Evergreen Runtime](https://developer.microsoft.com/en-us/microsoft-edge/webview2/).

---

## Linux

### Error al iniciar: biblioteca `libwebkit2gtk` no encontrada
* **Síntoma**: El ejecutable falla con un error `error while loading shared libraries: libwebkit2gtk-4.0.so.37`.
* **Solución**: Instala las bibliotecas requeridas según tu distribución:
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

## Errores Comunes de la App Store

### Error de Autenticación (`BadCredentials` o `2FA requerido`)
* Asegúrate de introducir correctamente el código de 6 dígitos en cuanto aparezca el modal de verificación de doble factor en la aplicación.
* Si Apple requiere verificación adicional en la web (como aceptar nuevos términos y condiciones de iCloud), inicia sesión una vez en [appleid.apple.com](https://appleid.apple.com) desde tu navegador para aceptar los avisos pendientes.

### Error `This item is not available in your storefront`
* Las aplicaciones de la App Store tienen licencias por país. Si intentas descargar una app exclusiva de la App Store de Estados Unidos con un Apple ID registrado en España o México, Apple rechazará la descarga.
* Asegúrate de que tu cuenta de Apple ID corresponda a la región donde la app está disponible.
