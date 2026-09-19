# Aplicación de Escritorio - Vista General

La aplicación de escritorio de **IPA Downloader** proporciona una interfaz fluida, moderna y accesible para gestionar todo el ciclo de vida de búsqueda, descarga, firma y sideloading de aplicaciones Apple.

---

## Interfaz de Usuario

Construida con **Vue 3**, **TailwindCSS** y **Wails v2**, la interfaz cuenta con:

* **Efecto Glassmorphism**: Paneles traslúcidos con bordes sutiles inspirados en macOS y diseño moderno.
* **Modo Oscuro y Modo Claro**: Adaptación automática a las preferencias del sistema o selección manual en ajustes.
* **Barra Lateral Intuitiva**: Navegación rápida entre las secciones principales:
  * :material-magnify: **Buscar**: Búsqueda global en el catálogo de la App Store.
  * :material-download: **Descargas**: Monitor de transferencias en tiempo real.
  * :material-cellphone: **Dispositivos**: Gestor de sideloading por USB para iPhone y iPad.
  * :material-bookshelf: **Biblioteca**: Archivo local de apps descargadas e historial.
  * :material-shopping: **Compras**: Historial de licencias asociadas a tu Apple ID.
  * :material-chip: **Firmwares**: Navegador de imágenes IPSW por modelo de dispositivo.
  * :material-cog: **Ajustes**: Preferencias de descargas, carpetas y autenticación.
  * :material-clipboard-text: **Logs**: Consola en vivo de eventos del sistema.

---

## Notificaciones y Estado del Sistema

* **Eventos en tiempo real**: Las barras de progreso, cambios de estado de autenticación y notificaciones de descarga completada se transmiten bidireccionalmente a través de eventos de Wails (`events.Emit`).
* **Integración con el sistema**: Soporta notificaciones nativas de escritorio al terminar descargas largas o completarse una instalación en dispositivo.
* **Consola de Logs en Vivo**: Permite inspeccionar qué peticiones HTTP se están enviando a los servidores de Apple, diagnosticar fallos de red y exportar registros con un solo clic.
