# Modo CLI - Visión General

El binario de **IPA Downloader** cuenta con arquitectura dual: si se invoca con argumentos desde la línea de comandos, se ejecuta en modo CLI utilizando la biblioteca [Cobra](https://github.com/spf13/cobra).

Esto permite integrar la herramienta en scripts bash, flujos de integración continua (CI/CD) o entornos de servidor headless sin entorno gráfico.

---

## Sintaxis General

```bash
ipa-downloader [comando] [subcomando] [opciones]
```

### Opciones Globales

| Bandera | Alias | Descripción |
|---|---|---|
| `--format` | | Formato de salida: `text` o `json` (por defecto: `text`) |
| `--non-interactive` | | Desactiva prompts interactivos (ideal para scripts automáticos) |
| `--verbose` | | Muestra registros detallados de depuración en la terminal |
| `--keychain-passphrase` | | Frase de paso para desbloquear el llavero si es necesario |
| `--help` | `-h` | Ayuda sobre cualquier comando o subcomando |
| `--version` | `-v` | Muestra la versión actual compilada de la herramienta |

---

## Salida en formato JSON

Para facilitar la integración con herramientas como `jq` o scripts en Python/Node.js, todos los comandos soportan el parámetro `--format json`:

```bash
ipa-downloader search "Telegram" --format json | jq '.[0].bundleId'
```
