# CLI Mode - Overview

The **IPA Downloader** binary features a dual architecture: when invoked with arguments from a terminal, it executes in CLI mode using the [Cobra](https://github.com/spf13/cobra) command framework.

This allows seamless integration into bash scripts, continuous integration pipelines (CI/CD), and headless servers without graphical environments.

---

## General Syntax

```bash
ipa-downloader [command] [subcommand] [flags]
```

### Global Flags

| Flag | Shorthand | Description |
|---|---|---|
| `--format` | | Output format: `text` or `json` (default: `text`) |
| `--non-interactive` | | Disables interactive prompts (suitable for automated scripts) |
| `--verbose` | | Enables verbose debug logging to stdout |
| `--keychain-passphrase` | | Passphrase for unlocking the OS keychain if required |
| `--help` | `-h` | Help for any command or subcommand |
| `--version` | `-v` | Prints current application version |

---

## JSON Output Formatting

To simplify integration with tools like `jq` or Python/Node.js automation scripts, all commands support the `--format json` flag:

```bash
ipa-downloader search "Telegram" --format json | jq '.[0].bundleId'
```
