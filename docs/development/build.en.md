# Building from Source

Guide for developers looking to build **IPA Downloader** locally or contribute to the project.

---

## Prerequisites

* **Go**: Version 1.25 or higher (installed and on your `PATH`).
* **Node.js**: Version 18 or higher with `npm`.
* **Wails CLI**: Version v2.12.0 or higher:
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

---

## 1. Clone the Repository

```bash
git clone https://github.com/ElJoker63/ipa-downloader.git
cd ipa-downloader
```

---

## 2. Live Development Mode (Hot-Reload)

To develop with Vite hot module replacement:

```bash
wails dev
```

This starts the Go backend and Vite frontend development server at `http://localhost:5173`, binding them automatically.

---

## 3. Build Production Desktop Binary

### Using Wails CLI:
```bash
wails build
```
The output binary will be placed in `build/bin/`.

### For Specific Target Platforms:
```bash
# macOS Universal (Intel + Apple Silicon)
wails build -platform darwin/universal

# Windows 64-bit
wails build -platform windows/amd64

# Linux 64-bit
wails build -platform linux/amd64
```

---

## 4. Build CLI Mode Only

If you only need the command-line interface:

```bash
cd frontend && npm install && npm run build && cd ..
go build -o ipa-downloader .
```

---

## 5. Running Automated Tests

```bash
# Generate mocks if necessary
go generate ./...

# Run unit tests
go test -v ./...
```
