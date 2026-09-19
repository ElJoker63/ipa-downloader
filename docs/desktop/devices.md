# Gestor USB y Sideloading

El módulo de dispositivos de **IPA Downloader** permite interactuar directamente con iPhones y iPads conectados por cable USB mediante el protocolo nativo `usbmuxd` (utilizando la biblioteca en Go `gidevice`), sin necesidad de instalar Xcode, iTunes o herramientas de terceros.

---

## Detección y Emparejamiento

1. Conecta tu dispositivo iOS / iPadOS mediante cable USB a tu ordenador.
2. Si es la primera vez que lo conectas, desbloquea la pantalla de tu iPhone o iPad y presiona **«Confiar en esta computadora»**, introduciendo tu código de desbloqueo.
3. En la sección **Dispositivos** de la aplicación, el dispositivo aparecerá automáticamente mostrando:
   * Nombre del dispositivo (ej. *"iPhone de Usuario"*).
   * Modelo exacto e identificador de hardware (ej. `iPhone14,2`).
   * Versión de iOS / iPadOS instalada.
   * UDID (Identificador único de dispositivo).
   * Estado de emparejamiento.

---

## Instalación Directa de Paquetes (.ipa)

Para instalar una aplicación en tu dispositivo:

1. Ve a la pestaña **Dispositivos**.
2. Selecciona el dispositivo conectado de la lista.
3. Haz clic en **Instalar IPA** y selecciona cualquier archivo `.ipa` descargado o presente en tu disco duro (o haz clic en el botón de instalar directamente desde tu biblioteca local).
4. El gestor validará la estructura del archivo y transmitirá la aplicación al demonio de instalación del dispositivo (`installation_proxy`).
5. La aplicación aparecerá en la pantalla de inicio de tu dispositivo una vez completada la barra de progreso.

!!! note "Requisito de Licencia"
    Para que una aplicación descargada de la App Store oficial se ejecute en tu dispositivo sin cerrarse de inmediato, el Apple ID que descargó el paquete debe ser el mismo que está activo en el dispositivo, o el paquete debe estar firmado con un perfil de aprovisionamiento de desarrollador compatible.

---

## Gestión de Aplicaciones Instaladas

* **Listar Aplicaciones**: Visualiza todas las aplicaciones de usuario instaladas en el dispositivo junto a su Bundle Identifier y versión.
* **Desinstalación Rápida**: Puedes eliminar aplicaciones directamente desde la computadora haciendo clic en el icono de papelera junto a cada app.
* **Información del Sistema**: Consulta el espacio disponible, versión de baseband y estado de activación del dispositivo.
