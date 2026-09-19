# Solución de Problemas en macOS

---

## 1. Error «La aplicación está dañada y no se puede abrir» o «No es compatible»

### Causa
Al descargar archivos `.zip` o aplicaciones directamente desde internet a través de navegadores como Safari o Google Chrome, macOS les asigna automáticamente el atributo extendido de seguridad **`com.apple.quarantine`**.

Dado que este proyecto es de código abierto y no utiliza un certificado de pago de Apple Developer ($99/año) para la notarización oficial en los servidores de Apple, el sistema **Gatekeeper** bloquea la apertura de la app o muestra un icono de prohibición.

### Solución
Para remover la marca de cuarentena y permitir que la app se abra sin problemas:

1. Abre la **Terminal** en tu Mac.
2. Ejecuta el siguiente comando (asumiendo que moviste la aplicación a tu carpeta de Aplicaciones):
   ```bash
   xattr -cr /Applications/ipa-downloader-desktop.app
   ```
3. Vuelve a hacer doble clic sobre la aplicación. Abrirá de inmediato.

---

## 2. Error al Actualizar desde la Aplicación

Si alguna versión previa falló durante la auto-actualización y dejó de abrir:
1. Asegúrate de tener instalada la versión **1.4.3 o superior**, la cual ya incluye el parche que descomprime correctamente el binario de los archivos `.zip` en macOS.
2. Si tu app quedó dañada por una versión antigua, descarga el `.zip` limpio más reciente desde [GitHub Releases](https://github.com/ElJoker63/ipa-downloader/releases/latest), reemplaza el archivo en `/Applications` y ejecuta el comando `xattr -cr` indicado arriba.

---

## 3. Dispositivos USB no detectados

Si tu iPhone o iPad no aparece en la pestaña **Dispositivos**:
1. Asegúrate de desbloquear la pantalla del dispositivo.
2. Si aparece el diálogo **«¿Confiar en esta computadora?»**, pulsa **Confiar** e introduce tu código de bloqueo.
3. Desconecta y vuelve a conectar el cable USB.
4. En macOS, asegúrate de que ninguna otra herramienta exclusiva (como máquinas virtuales que capturen los puertos USB) tenga bloqueado el dispositivo.
