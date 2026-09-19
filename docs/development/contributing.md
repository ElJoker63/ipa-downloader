# Guía de Contribución

¡Agradecemos tus contribuciones a **IPA Downloader**! Sigue estas pautas para enviar mejoras, correcciones o nuevas funcionalidades.

---

## Flujo de Trabajo para Pull Requests

1. **Fork del Repositorio**: Haz un fork del repositorio oficial a tu cuenta personal de GitHub.
2. **Crear una Rama**: Crea una rama descriptiva para tu cambio:
   ```bash
   git checkout -b fix/nombre-del-arreglo
   # o
   git checkout -b feat/nombre-de-la-caracteristica
   ```
3. **Desarrollar y Probar**:
   * Asegúrate de que el código siga el formato estándar de Go (`gofmt`).
   * Mantén los cambios enfocados y atómicos.
   * Ejecuta las pruebas antes de enviar el PR:
     ```bash
     go test ./...
     ```
4. **Commits Claros**: Escribe mensajes de commit siguiendo la convención de [Conventional Commits](https://www.conventionalcommits.org/):
   * `feat: ...` para nuevas funcionalidades.
   * `fix: ...` para corrección de bugs.
   * `docs: ...` para cambios de documentación.
   * `refactor: ...` para refactorización de código sin alterar comportamiento.
5. **Abrir el Pull Request**: Envía el PR contra la rama `main` explicando detalladamente el problema que resuelve y las pruebas realizadas.

---

## Estilo de Código

* **Go**: Código idiomático, manejo explícito de errores (`if err != nil`), sin dependencias pesadas innecesarias.
* **Frontend**: Vue 3 Composition API con `<script setup lang="ts">`, componentes modulares y clases semánticas de TailwindCSS.
* **Compatibilidad**: Preserva la compatibilidad hacia atrás en los parámetros de la CLI y formatos de salida JSON.
