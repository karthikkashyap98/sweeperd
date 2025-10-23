# sweeperd 

`sweeperd` is a lightweight, rule-based file organization daemon written in Go.
It watches your folders (or runs on demand) and automatically moves, renames, or deletes files based on declarative YAML rules.

---

## Features

- Declarative **YAML rule files** (`~/.sweeperd/rules.d/`)
- Simple **config.yaml** for global settings
- CLI-first design with subcommands:
  - `run` — Execute all rules once
  - `watch` — Daemon mode for real-time folder monitoring
  - `rules list` / `rules test` — Inspect and validate rules
  - `config show` — View effective configuration
  - `log tail` — Stream live log output
- Supports **dry-run**, **JSON output**, and **structured logs**
- Cross-platform (macOS / Linux)
- No dependencies beyond Go standard library

---

## Project Structure

```
sweeperd/
├── cmd/
│   └── sweeperd/           # CLI entrypoint
│       └── main.go
├── internal/
│   ├── cli/                # Cobra commands & CLI logic
│   ├── config/             # Config loader & validator
│   ├── rule/               # Rule parsing & validation
│   ├── watcher/            # Watch mode implementation (fsnotify)
│   ├── executor/           # Rule action executor
│   ├── app/                # Shared runtime context
│   └── utils/              # Path and common helpers
├── pkg/
│   └── constants/          # Default paths and global constants
├── docs/
│   ├── CONFIG_SPEC_MVP.md
│   ├── RULE_SPEC_MVP.md
│   └── SWEEPERD_INTERFACE_SPEC_MVP.md
└── go.mod
```

---

## Configuration

The configuration lives in `~/.autoorg/config.yaml` by default.

```yaml
rules_path: "~/.autoorg/rules.d"
log_path: "~/.autoorg/logs"
trash_path: "~/.autoorg/.trash"
watch_enabled: true
debounce_ms: 1000
parallelism: 4
dry_run_global: false
```

---

## Example Rule

Each rule is standalone and lives under `~/.autoorg/rules.d/*.yaml`.

```yaml
name: "Move PDFs to Documents"
enabled: true

match:
  folder: "~/Downloads"
  extensions: [".pdf"]
  older_than_days: 3

action:
  type: "move"
  target: "~/Documents/PDFs"
```

---

## Usage

### Run once
```bash
sweeperd run
```

### Start as a daemon (watch mode)
```bash
sweeperd watch --foreground
```

---

## Documentation

All internal specifications are documented in the `/docs` folder:

- [Rule Specification](docs/RULE_SPEC_MVP.md)
- [Configuration Specification](docs/CONFIG_SPEC_MVP.md)
- [Interface Specification](docs/SWEEPERD_INTERFACE_SPEC_MVP.md)

---


## Future Scope

- Live reload when rule files change
- Advanced filters (regex, metadata, file size)
- Rule groups and priorities
- Native macOS / Linux background service integration
- Optional TUI dashboard for monitoring
- Plugin system for custom actions

---

## 📜 License

MIT © 2025 — Karthik Sudhakar Kashyap

---


