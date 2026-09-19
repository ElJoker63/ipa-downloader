# Búsqueda y Descarga de Aplicaciones

---

## Motor de Búsqueda

El módulo de búsqueda interactúa con las APIs del catálogo de iTunes y la App Store:

* **Búsqueda por Término**: Escribe palabras clave para encontrar aplicaciones relevantes.
* **Búsqueda por Identificador**: Si conoces el Bundle Identifier exacto (ej. `org.videolan.vlc-ios`), búscalo directamente para obtener el resultado preciso.
* **Filtro de Plataforma**:
  * iPhone
  * iPad
  * Apple TV (tvOS)
  * Apple Vision Pro (visionOS)
  * Mac App Store (macOS)

---

## Ficha Detallada de la Aplicación

Al hacer clic en cualquier resultado, se despliega el modal de detalles con información completa:

1. **Metadatos Generales**: Nombre, desarrollador, categoría, calificación de estrellas, precio, tamaño en megabytes y fecha de última actualización.
2. **Galería de Capturas**: Visualiza las capturas de pantalla oficiales enviadas por el desarrollador.
3. **Selector de Versiones Históricas**:
   * IPA Downloader consulta la lista completa de identificadores de versiones (`externalVersionId`) registradas en los servidores de Apple para esa aplicación.
   * Puedes seleccionar cualquier versión previa para descargar una versión histórica en lugar de la última versión publicada.

---

## Gestor de Descargas Concurrentes

* **Descarga por bloques (Chunked streaming)**: Optimizado para descargar archivos grandes a alta velocidad.
* **Métricas en tiempo real**: Velocidad de descarga promedio móvil, porcentaje exacto, bytes transferidos y cálculo dinámico de tiempo restante (ETA).
* **Control de flujo**: Pausa, reanuda o cancela cualquier descarga activa.
* **Mecanismo de Recuperación (Fallback)**: Si los servidores de volumen de Apple devuelven un error temporal de descarga o informan que el enlace expiró, el motor activa automáticamente la ruta de recuperación de licencia y reintentos sin perder el progreso previo.

---

## Inyección de Firma FairPlay (SINF)

Cuando descargas una aplicación desde la App Store oficial, los binarios ejecutables vienen cifrados mediante el DRM FairPlay de Apple. Para que el dispositivo del usuario pueda descifrarlo al ejecutarlo, el paquete requiere un archivo de firma (`SC_Info/*.sinf`).

IPA Downloader:
1. Extrae los metadatos de licencia devueltos por el endpoint `buyProduct` / `downloadProduct` de Apple.
2. Genera los bloques `sinf` correspondientes a la cuenta autenticada.
3. Parcha e inyecta la estructura `SC_Info` directamente en el archivo `.ipa`.
4. Entrega un archivo `.ipa` 100% válido y listo para ser instalado o preservado.
