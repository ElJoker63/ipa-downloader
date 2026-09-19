# Contributing Guide

Thank you for contributing to **IPA Downloader**! Please follow these guidelines when opening issues, pull requests, or proposing improvements.

---

## Pull Request Workflow

1. **Fork the Repository**: Fork the official repository to your GitHub account.
2. **Create a Branch**: Use a descriptive branch name:
   ```bash
   git checkout -b fix/your-bug-fix
   # or
   git checkout -b feat/your-feature-name
   ```
3. **Develop & Test**:
   * Format code with standard Go formatting (`gofmt`).
   * Keep changes focused and atomic.
   * Run tests before submitting:
     ```bash
     go test ./...
     ```
4. **Clear Commit Messages**: Follow the [Conventional Commits](https://www.conventionalcommits.org/) convention:
   * `feat: ...` for new capabilities.
   * `fix: ...` for bug fixes.
   * `docs: ...` for documentation improvements.
   * `refactor: ...` for non-behavioral refactoring.
5. **Open a Pull Request**: Submit your PR targeting `main` with a clear explanation of what changed and test results.

---

## Code Style

* **Go**: Idiomatic Go, explicit error checks (`if err != nil`), avoid unnecessary heavy external dependencies.
* **Frontend**: Vue 3 Composition API with `<script setup lang="ts">`, modular components, and TailwindCSS utility classes.
* **Compatibility**: Preserve backward compatibility for CLI flags and JSON output formats.
