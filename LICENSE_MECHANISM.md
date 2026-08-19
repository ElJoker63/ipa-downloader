# Mecanismo de licencias ByteCopy — spec para reimplementar en Go

Sistema de activación offline por código, firmado con Ed25519, sin servidor.
Basado en `lib/services/license_service.dart` y `registry_anchor_service.dart`.

## 1. Formato del código de licencia

Código = 10 bytes de payload + 64 bytes de firma Ed25519 = **74 bytes**, codificados en **Base32** (RFC 4648, sin padding `=`), en mayúsculas.

Layout del payload (10 bytes, big-endian):

| Offset | Bytes | Campo | Descripción |
|---|---|---|---|
| 0 | 1 | `version` | Debe ser `1` |
| 1 | 1 | `tipo` | `0`=normal, `1`=master, `2`=trial |
| 2 | 4 | `machineHash` | primeros 4 bytes de SHA-256(machineId), solo relevante si tipo=normal |
| 6 | 4 | `expiryDays` | días desde epoch 1970-01-01 (uint32 big-endian). `0xFFFFFFFF` = sin expiración (solo válido si master) |

Firma: Ed25519 sobre los 10 bytes del payload (no sobre el texto Base32).

```
payload  = version(1) | tipo(1) | machineHash(4) | expiryDays(4)
signature = Ed25519_sign(privateKey, payload)   // 64 bytes
code_raw  = payload | signature                  // 74 bytes
code_txt  = Base32(code_raw)                      // sin '='
```

## 2. Generación del código (lado emisor / vendedor)

Con la clave privada Ed25519 (nunca va en la app):
1. Calcular `machineHash` = primeros 4 bytes de SHA-256(machineIdHex del cliente), o ceros si es master/trial.
2. Calcular `expiryDays` = (fecha expiración - 1970-01-01).days, o `0xFFFFFFFF` si master sin vencimiento.
3. Armar payload de 10 bytes, firmar con Ed25519, concatenar firma, codificar en Base32.

## 3. Verificación del código (lado cliente, dentro de la app)

1. Normalizar: mayúsculas, quitar todo carácter fuera de `[A-Z2-7]`.
2. Base32-decode → debe dar exactamente 74 bytes; si no, código inválido.
3. Parsear campos según el layout de arriba.
4. Verificar `version == 1`.
5. Verificar firma Ed25519 de los 10 bytes de payload contra la clave **pública** embebida en el binario.
6. Según `tipo`:
   - `1` (master): válido en cualquier equipo.
   - `2` (trial): válido en cualquier equipo, sin chequeo de machine hash.
   - `0` (normal): calcular `SHA-256(machineId)` del equipo local, tomar primeros 4 bytes, comparar con `machineHash` del código. Si no coincide → "pertenece a otro equipo".
7. `expiryDate = 1970-01-01 + expiryDays` (nil si master con `0xFFFFFFFF`).

### machineId (fingerprint del equipo)

```
raw = MachineGuid (HKLM\SOFTWARE\Microsoft\Cryptography\MachineGuid) + "|" + COMPUTERNAME
machineId = hex(SHA-256(raw))
```
En Go: leer `MachineGuid` vía `golang.org/x/sys/windows/registry` (HKLM, puede requerir el mismo acceso que usa `reg query` — de 32/64 bits, ver nota abajo), concatenar con `os.Getenv("COMPUTERNAME")`, SHA-256, hex.

## 4. Almacenamiento local tras activar

Archivo `license.dat` (JSON, en carpeta de datos de la app — equivalente a `storage_paths.dart`):
```json
{ "code": "TEXTO_BASE32_DEL_CODIGO", "last_seen_ord": 738900 }
```
`last_seen_ord` = ordinal de día (días desde 0001-01-01, estilo `time.Date(...).Sub(epoch).Hours()/24 + 1`) de la última vez que se vio la app activa. Sirve como ancla anti-retroceso de reloj barata.

## 5. Ancla anti-rollback en el Registro de Windows (protección fuerte)

Problema que resuelve: borrar `license.dat` + atrasar el reloj del sistema + reactivar el mismo código revivía una licencia vencida, porque `last_seen_ord` volvía a cero.

Solución: guardar el ordinal de día máximo visto en una clave aparte del Registro (HKCU, sin admin), con nombre neutro para no delatar su propósito, firmada con HMAC-SHA256 atada al `machineId` para que no sea editable a mano.

```
Key:   HKCU\Software\Classes\ByteCopySvc
Value: SyncTag (REG_SZ)
Valor codificado = hex(ordDay) + "." + hex(HMAC_SHA256(key=pepper|machineId, msg="ord:"+ordDay))
pepper = "bc-anchor-7fQ1-k9F2"   // cambiar este secreto en la reimplementación
```

Lectura: parsear `hex.hex`, recomputar el HMAC con el `machineId` local; si no coincide, tratar como inexistente (valor `nil`, no error fatal — nunca debe romper la app si el registro no está disponible o `reg.exe` falla).

Escritura: siempre "silenciosa" (no falla la app si el registro no es escribible).

En Go usar `golang.org/x/sys/windows/registry` directamente en vez de invocar `reg.exe` (más robusto que parsear stdout).

## 6. Flujo de activación (`activate(code)`)

1. `decodeAndVerify(code)` → si falla, devolver el mensaje de error tal cual.
2. Cargar `license.dat` (o vacío si no existe).
3. `todayOrd` = ordinal de hoy.
4. `anchorOrd` = leer ancla del Registro para este `machineId` (0 si no hay).
5. `prevOrd` = `last_seen_ord` guardado en el archivo (0 si no hay).
6. `maxOrd = max(todayOrd, prevOrd, anchorOrd)`.
7. Guardar `code` y `last_seen_ord = maxOrd` en `license.dat`.
8. Escribir `maxOrd` en el ancla del Registro.

## 7. Flujo de verificación de estado (`getStatus()`), a correr en cada arranque / periódicamente

1. `machineId`, `todayOrd` = hoy.
2. `anchorOrd` = leer ancla del Registro.
3. `anchorRollback = todayOrd < anchorOrd`. Si **no** hay rollback, actualizar el ancla a `max(todayOrd, anchorOrd)` (avanza el ancla con el uso normal).
4. Cargar `license.dat`. Si no hay `code` guardado → no licenciado, razón "Sin código activado".
5. `decodeAndVerify(code, machineId)`. Si falla → no licenciado con el mensaje de error.
6. `lastSeenOrd = max(last_seen_ord del archivo, anchorOrd)`.
7. Si `anchorRollback || todayOrd < lastSeenOrd` → no licenciado, razón "Se detectó un cambio de fecha/hora sospechoso".
8. Si hay `expiryDate` y `hoy > expiryDate` → no licenciado, razón "La suscripción venció", `daysLeft = 0`.
9. Si todo OK: actualizar `last_seen_ord = max(todayOrd, lastSeenOrd)` en el archivo, devolver licenciado=true con `daysLeft = expiryDate - hoy` (o null si no expira).

## 8. Equivalentes Go sugeridos

| Dart | Go |
|---|---|
| `cryptography` (Ed25519) | `crypto/ed25519` (stdlib) |
| `crypto/sha256`, `crypto/hmac` | `crypto/sha256`, `crypto/hmac` (stdlib) |
| Base32 custom | `encoding/base32` con `base32.StdEncoding.WithPadding(base32.NoPadding)` |
| `Process.runSync('reg', ...)` | `golang.org/x/sys/windows/registry` |
| JSON local | `encoding/json` + carpeta `%APPDATA%\ByteCopy\` o similar |

## 9. Puntos críticos a no perder en el port

- La clave pública Ed25519 va embebida en el binario; la privada **nunca** se distribuye (queda en tu generador de códigos, separado del repo de la app).
- El anti-rollback depende de **dos** fuentes de verdad (archivo + Registro) y siempre toma el máximo de ambas — no es solo "guardar la fecha", es "nunca dejar que el ordinal retroceda".
- Todos los fallos de Registro deben ser silenciosos (try/catch amplio): si `reg.exe` o la clave no existen, la app debe seguir funcionando, solo pierde la protección extra.
- El pepper del HMAC y el nombre/ruta de la clave del Registro son secretos de ofuscación menores; cambiarlos en el port a Go (no reusar los mismos strings expuestos aquí).
- Trial y Master no verifican `machineHash`; Normal sí.
