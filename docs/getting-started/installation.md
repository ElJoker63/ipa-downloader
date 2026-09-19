# Guía de Instalación

**IPA Downloader** está disponible para macOS, Windows y Linux mediante paquetes precompilados en [GitHub Releases](https://github.com/ElJoker63/ipa-downloader/releases/latest).

---

## macOS

### Requisitos
* macOS 12 (Monterey) o superior (compatible con Intel y Apple Silicon M1/M2/M3/M4).

### Pasos de Instalación

1. Descarga el paquete `ipa-downloader-*-macos.zip` desde la sección de **Releases**.
2. Descomprime el archivo `.zip`. Obtendrás la aplicación `ipa-downloader-desktop.app`.
3. Arrastra la aplicación a tu carpeta `/Applications` (Aplicaciones).

!!! warning "Aviso de Seguridad en macOS (Gatekeeper / Cuarentena)"
    Dado que las versiones de GitHub se distribuyen sin firma pagada de Apple Developer ID, macOS marcará el archivo descargado desde el navegador con el atributo de cuarentena (`com.apple.quarantine`).

    Para permitir que abra sin advertencias, abre la Terminal y ejecuta:
    ```bash
    xattr -cr /Applications/ipa-downloader-desktop.app
    ```
    Una vez ejecutado, podrás abrir la aplicación con normalidad desde Spotlight o Launchpad. Consulta la [guía de solución de problemas en macOS](../troubleshooting/macos.md) para más detalles.

---

## Windows

### Requisitos
* Windows 10 (versión 1809 o superior) o Windows 11 (64-bit).
* WebView2 Runtime (incluido por defecto en Windows 10/11 actualizados).

### Pasos de Instalación
1. Descarga el ejecutable `ipa-downloader-*-windows.exe` desde las **Releases**.
2. Guarda el archivo en tu carpeta de preferencia (ej. `C:\Archivos de Programa\IPA Downloader` o `Descargas`).
3. Haz doble clic en el ejecutable para iniciar la aplicación.
4. Si Windows SmartScreen muestra un aviso de editor desconocido, haz clic en **Más información** -> **Ejecutar de todas formas**.

---

## Linux

### Requisitos
* Distribuciones basadas en Debian/Ubuntu (20.04+), Arch Linux, Fedora u OpenSUSE.
* Bibliotecas `libgtk-3` y `webkit2gtk-4.0`.

### Instalación de Dependencias
En Ubuntu / Debian:
```bash
sudo apt-get update
sudo apt-get install -y libgtk-3-0 libwebkit2gtk-4.0-37
```

En Arch Linux:
```bash
sudo pacman -S gtk3 webkit2gtk
```

### Ejecución
1. Descarga el binario `ipa-downloader-*-linux`.
2. Asigna permisos de ejecución:
   ```bash
   chmod +x ipa-downloader-*-linux
   ```
3. Ejecuta la aplicación:
   ```bash
   ./ipa-downloader-*-linux
   ```
