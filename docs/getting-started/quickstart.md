# Inicio Rápido

Aprende a configurar tu cuenta de Apple ID y realizar tu primera búsqueda y descarga en cuestión de minutos.

---

## 1. Inicio de Sesión con Apple ID

Para descargar aplicaciones legítimas de la App Store, Apple requiere una sesión autenticada.

1. Abre **IPA Downloader**.
2. Dirígete a la pestaña **Cuenta** o haz clic en el botón de perfil superior.
3. Introduce tu **correo electrónico de Apple ID** y tu **contraseña**.
4. Haz clic en **Iniciar Sesión**.
5. Si tu cuenta tiene activada la autenticación de doble factor (2FA), aparecerá un modal pidiendo el código de 6 dígitos enviado a tus dispositivos Apple de confianza o por SMS. Ingrésalo para completar la autenticación.

!!! tip "Seguridad de tus credenciales"
    Tus credenciales se transmiten directamente a los servidores oficiales de Apple (`gsa.apple.com`). El token de sesión resultante se guarda cifrado en el llavero de tu sistema operativo (Keychain / Credential Manager). Nunca se envían a ningún servidor de terceros.

---

## 2. Buscar una Aplicación

1. Haz clic en la pestaña **Buscar** en la barra lateral.
2. Escribe el nombre de la app (ej. `Telegram`, `WhatsApp`) o el Bundle ID (ej. `ph.telegra.Telegraph`).
3. Puedes filtrar los resultados por plataforma:
   * **iOS (iPhone)**
   * **iPadOS**
   * **tvOS**
   * **visionOS**
   * **macOS**
4. Haz clic en la tarjeta de cualquier aplicación para ver su ficha técnica: descripción, capturas de pantalla, versión actual, tamaño y lista de versiones históricas disponibles.

---

## 3. Descargar el Paquete

1. En la ficha de la aplicación, haz clic en el botón **Descargar**.
2. Si nunca has obtenido una licencia de esa app en tu cuenta (incluso si es gratuita), la aplicación solicitará automáticamente la licencia (`purchase`) a Apple.
3. La descarga comenzará de inmediato y podrás monitorear el progreso, velocidad y tiempo estimado en la pestaña **Descargas**.
4. Al completarse, la app inyectará automáticamente la firma FairPlay SINF correspondiente a tu cuenta para que el paquete sea completamente funcional.

---

## 4. Uso Rápido en Modo CLI

Si prefieres la terminal, el mismo binario incluye la interfaz de línea de comandos:

```bash
# Iniciar sesión
./ipa-downloader auth login --email "tu-correo@icloud.com"

# Buscar una aplicación
./ipa-downloader search "Spotify" --limit 5

# Descargar por Bundle ID
./ipa-downloader download -b "com.spotify.client" --purchase
```
