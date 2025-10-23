# sweeperd Interface Specification (MVP)

I/O conventions, exit codes, and minimal programmatic hooks. It is intended as the reference
for users and as a contract for implementation.

---

## 1. Invocation Model

`sweeperd` is a CLI-first tool. All subcommands are invoked as:

```
sweeperd <subcommand> [flags] [args]
```

Global flags apply to all subcommands and can be placed before or after the subcommand.

---

## 2. Global Flags

| Flag | Type | Default | Description |
|------|------|---------|-------------|
| `--config` | string | `~/.autoorg/config.yaml` | Path to config file. |
| `--rules`  | string | `~/.autoorg/rules.d/`    | Directory containing rule files. |
| `--log`    | string | `~/.autoorg/logs/`       | Directory for logs. |
| `--dry-run`| bool   | `false`                  | Force dry-run regardless of rule or config settings. |
| `--json`   | bool   | `false`                  | Emit machine-readable JSON on stdout for successful operations. |
| `--quiet`  | bool   | `false`                  | Suppress non-error stdout messages. |
| `--verbose`| bool   | `false`                  | Increase log verbosity on stderr. |
| `--version`| bool   | `false`                  | Print version and exit. |

Notes:
- Paths supporting `~` expansion are normalized at startup.
- `--dry-run` overrides rule-level and config-level settings.

---

## 3. Subcommands (MVP)

### 3.1 `run`
Runs all enabled rules against their configured folders once and exits.

```
sweeperd run [--rule <name>] [--match-folder <path>] [--concurrency N]
```

Flags:
- `--rule <name>`: Limit execution to a single rule by name.
- `--match-folder <path>`: Override the rule's `match.folder` at runtime (useful for testing).
- `--concurrency N`: Limit number of concurrent file operations (defaults to `parallelism` from config).

Behavior:
- Loads and validates config and rules.
- Executes in filename order of rule files by default (MVP behavior).
- Prints a human-readable summary unless `--json` is set.

Exit codes: see Section 7.

---

### 3.2 `watch`
Starts watchers for all folders referenced by enabled rules and applies actions on file system events.

```
sweeperd watch [--foreground] [--pidfile <path>] [--socket <path>]
```

Flags:
- `--foreground`: Run in the foreground (do not daemonize).
- `--pidfile <path>`: Write PID to path if supported by platform.
- `--socket <path>`: Optional control socket path for local IPC (see 6. Programmatic Control).

Behavior:
- Debounces events per file (default debounce from config).
- Logs all actions.
- Terminates on SIGINT/SIGTERM, performing best-effort shutdown.

---

### 3.3 `rules list`
Lists discovered rules and basic metadata.

```
sweeperd rules list [--long]
```

Flags:
- `--long`: Show extended details (status, action type, folder, options).

Output:
- Human-readable table by default; JSON when `--json` is set.

---

### 3.4 `rules test`
Validates rules without executing actions.

```
sweeperd rules test [--rule <name>] [--sample <path>] [--explain]
```

Flags:
- `--rule <name>`: Validate a single rule by name.
- `--sample <path>`: Test matching logic against a single file path.
- `--explain`: Print detailed reasoning about why a file would or would not match.

Exit with non-zero if any rule fails validation.

---

### 3.5 `config show`
Displays the effective, merged configuration after defaults and normalization.

```
sweeperd config show [--raw]
```

Flags:
- `--raw`: Show on-disk configuration without defaults applied.

Output:
- YAML by default; JSON if `--json` is set.

---

### 3.6 `log tail`
Streams recent log entries.

```
sweeperd log tail [--follow] [--lines N]
```

Flags:
- `--follow`: Continue streaming as new entries arrive.
- `--lines N`: Number of lines to show from the end of the log (default 200).

---

## 4. I/O and Output Conventions

- **Stdout**: User-facing results (summaries, tables, JSON with `--json`).
- **Stderr**: Progress, warnings, and errors.
- **Logs**: Structured log files are written under `--log` (default `~/.autoorg/logs/`).

### 4.1 JSON Output Schema (MVP)
When `--json` is enabled, `run` and `watch` emit newline-delimited JSON objects with the following shape:

```json
{
  "ts": "2025-10-22T18:30:01Z",
  "event": "action",
  "rule": "Move PDFs to Documents",
  "path": "/Users/u/Downloads/file.pdf",
  "action": "move",
  "target": "/Users/u/Documents/PDFs",
  "dry_run": false,
  "status": "ok"
}
```

Possible `event` values: `action`, `skip`, `error`, `start`, `stop`.

---

## 5. Configuration Interaction

- CLI flags override environment variables and configuration file values.
- Config keys of interest for the interface:
  - `rules_path`, `log_path`, `trash_path`
  - `watch_enabled`, `watch_folders`, `debounce_ms`
  - `parallelism`
  - `dry_run_global`

---

## 6. Programmatic Control (MVP placeholder)

For MVP, `sweeperd` is CLI-only. A minimal local control channel may be added via a POSIX domain socket when `watch` runs with `--socket`:

- Request format: newline-delimited JSON with `{"cmd": "status"}` or `{"cmd": "shutdown"}`.
- Response format: newline-delimited JSON with `{"ok": true, "data": ...}`.
- This interface is considered unstable in MVP; no compatibility guarantees yet.

This section can be omitted in the initial cut if not implemented.

---

## 7. Exit Codes

| Code | Meaning |
|------|---------|
| `0`  | Success. |
| `1`  | Invalid CLI usage (unknown subcommand or flags). |
| `2`  | Configuration or rule validation failed. |
| `3`  | Runtime errors occurred (partial success). |
| `4`  | Fatal error (could not start watchers or perform required action). |

---

## 8. Environment Variables

| Variable | Purpose |
|----------|---------|
| `SWEEPERD_CONFIG` | Overrides `--config`. |
| `SWEEPERD_RULES`  | Overrides `--rules`. |
| `SWEEPERD_LOG_DIR`| Overrides `--log`. |

Environment variables are lower priority than explicit CLI flags.

---

## 9. Compatibility and Versioning

- The interface follows semantic versioning once `v1.0.0` is reached.
- Until then, subcommands and flags may change between minor versions.
- `--version` prints the version string and build metadata.

---

## 10. Usage Examples

Run once with defaults:
```
sweeperd run
```

Run a single rule with JSON output:
```
sweeperd run --rule "Move PDFs to Documents" --json
```

Start watchers in foreground:
```
sweeperd watch --foreground --verbose
```

Validate rules and explain a sample file:
```
sweeperd rules test --sample ~/Downloads/report.pdf --explain
```

Show effective config:
```
sweeperd config show
```

Tail logs:
```
sweeperd log tail --follow --lines 500
```

---

## 11. Non-Goals (MVP)

- No GUI or TUI in MVP.
- No remote control API beyond optional local socket.
- No cross-machine rule sync.
- No pluggable rule marketplace.

