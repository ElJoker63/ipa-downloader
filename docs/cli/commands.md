# Referencia de Comandos CLI

A continuación se detalla la sintaxis, argumentos y ejemplos de cada subcomando disponible en el modo CLI.

---

## 1. Autenticación (`auth`)

Gestiona la sesión con la App Store.

### `auth login`
Inicia sesión con tu cuenta de Apple ID.

```bash
ipa-downloader auth login [opciones]
```

* `--email`, `-e`: Dirección de correo electrónico del Apple ID.
* `--password`, `-p`: Contraseña del Apple ID.
* `--code`, `-c`: Código de verificación de doble factor (2FA) de 6 dígitos.

**Ejemplo interactivo:**
```bash
ipa-downloader auth login -e "usuario@icloud.com"
```

**Ejemplo en una sola línea:**
```bash
ipa-downloader auth login -e "usuario@icloud.com" -p "MiClaveSecreta" -c "123456"
```

### `auth info`
Muestra el estado de la sesión activa y el nombre de la cuenta actual.

```bash
ipa-downloader auth info
```

### `auth revoke`
Cierra la sesión actual y elimina los tokens guardados en el llavero local.

```bash
ipa-downloader auth revoke
```

---

## 2. Búsqueda (`search`)

Busca aplicaciones en el catálogo oficial de la App Store.

```bash
ipa-downloader search <término> [opciones]
```

* `--limit`, `-l`: Número máximo de resultados a devolver (por defecto: `5`).
* `--platform`: Plataforma objetivo (`iphone`, `ipad`, `appletv`, `visionos`, `macos`). Por defecto: `iphone`.

**Ejemplos:**
```bash
# Buscar apps para iPhone
ipa-downloader search "WhatsApp" --limit 3

# Buscar apps para macOS
ipa-downloader search "Final Cut" --platform macos --limit 1

# Exportar resultados en JSON
ipa-downloader search "Slack" --format json
```

---

## 3. Obtención de Licencia (`purchase`)

Adquiere la licencia de una aplicación gratuita en tu cuenta de Apple ID (necesario si nunca la has descargado previamente).

```bash
ipa-downloader purchase --bundle-identifier <bundle-id>
```

* `--bundle-identifier`, `-b`: Identificador del paquete de la app (ej. `ph.telegra.Telegraph`).

**Ejemplo:**
```bash
ipa-downloader purchase -b "ph.telegra.Telegraph"
```

---

## 4. Listar Compras (`list-purchases`)

Lista las aplicaciones adquiridas en la cuenta autenticada.

```bash
ipa-downloader list-purchases [opciones]
```

* `--max-results`: Cantidad de resultados por página (por defecto: `10`).
* `--page`: Número de página para paginación.

**Ejemplo:**
```bash
ipa-downloader list-purchases --max-results 20 --page 1
```

---

## 5. Consultar Versiones (`list-versions`)

Obtiene el historial de identificadores de versión (`externalVersionId`) de una aplicación específica.

```bash
ipa-downloader list-versions --bundle-identifier <bundle-id>
```

**Ejemplo:**
```bash
ipa-downloader list-versions -b "com.spotify.client"
```

---

## 6. Metadatos de Versión (`get-version-metadata`)

Consulta los metadatos detallados de una versión histórica puntual.

```bash
ipa-downloader get-version-metadata --bundle-identifier <bundle-id> --external-version-id <version-id>
```

---

## 7. Descarga (`download`)

Descarga el paquete oficial `.ipa` o `.pkg` desde los servidores de Apple.

```bash
ipa-downloader download [opciones]
```

* `--bundle-identifier`, `-b`: Bundle ID de la app.
* `--output`, `-o`: Ruta donde guardar el archivo resultante (por defecto: nombre automático en el directorio actual).
* `--purchase`: Si la cuenta no tiene la licencia de la app, intenta adquirirla automáticamente antes de descargar.
* `--external-version-id`: ID de versión histórica específica si no deseas la última versión.
* `--platform`: Plataforma del paquete (`iphone`, `ipad`, `appletv`, `visionos`, `macos`).

**Ejemplos:**
```bash
# Descargar la versión más reciente
ipa-downloader download -b "org.videolan.vlc-ios" --purchase

# Descargar una versión histórica específica
ipa-downloader download -b "com.spotify.client" --external-version-id "854123456" -o "./Spotify_old.ipa"

# Descargar paquete de macOS
ipa-downloader download -b "com.apple.Keynote" --platform macos -o "./Keynote.pkg"
```
