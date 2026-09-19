# Compilación desde el Código Fuente

Guía para desarrolladores que desean compilar **IPA Downloader** localmente o contribuir al proyecto.

---

## Requisitos Previos

* **Go**: Versión 1.24 o superior (instalado y accesible en el `PATH`).
* **Node.js**: Versión 18 o superior con `npm`.
* **Wails CLI**: Versión v2.12.0 o superior:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

---

## 1. Clonar el Repositorio

```bash
git clone https://github.com/ElJoker63/ipa-downloader.git
cd ipa-downloader
```

---

## 2. Modo de Desarrollo en Vivo (Hot-Reload)

Para desarrollar la interfaz con recarga rápida de Vite:

```bash
wails dev
```

Este comando inicia el backend en Go y el servidor de desarrollo de Vite para el frontend en `http://localhost:5173`, conectando ambos automáticamente.

---

## 3. Compilar la Aplicación de Escritorio (Producción)

### Con la CLI de Wails:
```bash
wails build
```
El binario resultante se generará en `build/bin/`.

### Para plataformas específicas:
```bash
# macOS Universal (Intel + Apple Silicon)
wails build -platform darwin/universal

# Windows 64-bit
wails build -platform windows/amd64

# Linux 64-bit
wails build -platform linux/amd64
```

---

## 4. Compilar en Modo CLI Únicamente

Si solo necesitas el binario ejecutable con la interfaz de consola:

```bash
cd frontend && npm install && npm run build && cd ..
go build -o ipa-downloader .
```

---

## 5. Ejecutar Pruebas Automatizadas

```bash
# Generar mocks e interfaces si es necesario
go generate ./...

# Ejecutar suite de pruebas unitarias
go test -v ./...
```
