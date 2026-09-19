# Biblioteca Local y Compras de la Cuenta

---

## Biblioteca Local

La pestaña **Biblioteca** almacena un registro persistente de todos los paquetes `.ipa` y `.pkg` que has descargado a través de la aplicación:

* **Acceso Rápido al Archivo**: Abre la carpeta contenedora en el Finder (macOS) o Explorador de Archivos (Windows/Linux) con un solo clic.
* **Instalación Directa**: Si tienes un dispositivo USB conectado, puedes enviar cualquier app de tu biblioteca a instalar directamente en el dispositivo.
* **Marcador de Favoritos**: Marca con una estrella las aplicaciones esenciales que sueles reinstalar o respaldar con frecuencia.
* **Búsqueda en la Biblioteca**: Filtra por nombre, versión o fecha de descarga en tu colección local.
* **Almacenamiento Seguro**: Todos los metadatos se guardan localmente en la base de datos SQLite embebida (`modernc.org/sqlite`), sin enviar registros a servidores externos.

---

## Compras de la Cuenta (Purchased Apps)

La sección **Compras** consulta directamente el registro histórico de licencias asociadas a tu cuenta de Apple ID en los servidores de Apple:

* **Exploración Histórica**: Visualiza aplicaciones que compraste u obtuviste en el pasado, incluso si ya no aparecen destacadas en la tienda.
* **Filtrado por Plataforma**: Revisa compras separadas para iPhone, iPad, Mac o Apple TV.
* **Re-descarga Inmediata**: Descarga cualquier aplicación adquirida previamente sin necesidad de buscarla manualmente en el catálogo general.
* **Caché Local**: Las compras se cachean en segundo plano en la base de datos local para que la navegación sea instantánea cada vez que abras la aplicación.
