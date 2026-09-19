# Explorador de Firmwares IPSW

La aplicación incluye un explorador integrado para consultar imágenes oficiales de restauración de firmware de Apple (`.ipsw`) para todos los modelos de iPhone, iPad, iPod touch, Apple TV y Apple Silicon Mac.

---

## Funcionalidades del Módulo

1. **Catálogo de Dispositivos**: Navega por categorías y familias de dispositivos Apple organizados por generaciones.
2. **Listado de Versiones**:
   * Identifica la versión exacta de iOS / iPadOS / macOS.
   * Fecha de lanzamiento oficial.
   * Tamaño del archivo de instalación en gigabytes.
   * Checksum SHA-1 / MD5 oficial para verificar la integridad de la descarga.
3. **Estado de Firma (Signing Status)**:
   * Muestra claramente si la versión de firmware **sigue siendo firmada por Apple** (permitiendo restaurar o hacer downgrade con iTunes/Finder).
   * Advierte cuando una versión ya no está firmada y Apple no permite su activación.
4. **Descarga Directa**: Enlaces de descarga directos a los Content Delivery Networks (CDNs) oficiales de Apple (`appldnld.apple.com`).
