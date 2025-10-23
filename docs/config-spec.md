# Configuration Specification (MVP)

The configuration file defines global settings that control rule loading, logging,
and watcher behavior. The configuration is written in YAML and located by default at:

```
~/.sweeperd/config.yaml
```

## Configuration Structure

```
Config
 ├─ rules_path: string (default: ~/.sweeperd/rules.d/)
 ├─ log_path: string (default: ~/.sweeperd/logs/)
 ├─ trash_path: string (default: ~/.sweeperd/.trash/)
 ├─ watch_enabled: bool (default: false)
 ├─ watch_folders: [string]
 ├─ debounce_ms: int (default: 1000)
 ├─ parallelism: int (default: 4)
 └─ dry_run_global: bool (default: false)
```

## Example

```yaml
rules_path: ~/.sweeperd/rules.d/
log_path: ~/.sweeperd/logs/
trash_path: ~/.sweeperd/.trash/

watch_enabled: true
watch_folders:
  - ~/Downloads
  - ~/Desktop

debounce_ms: 1000
parallelism: 4
dry_run_global: false
```

## Behavior

1. The application reads `config.yaml` on startup.
2. All relative paths are expanded to absolute paths.
3. Global options override rule-level options where applicable.
4. `watch_enabled` determines if folder watchers are started automatically.
5. Logs are written under `log_path`, and deleted files are moved to `trash_path`.

## Default Values

| Field | Default | Description |
|--------|----------|-------------|
| `rules_path` | `~/.sweeperd/rules.d/` | Directory containing rule files. |
| `log_path` | `~/.sweeperd/logs/` | Directory where logs are stored. |
| `trash_path` | `~/.sweeperd/.trash/` | Directory for safe delete operations. |
| `watch_enabled` | `false` | Enables automatic file watching. |
| `debounce_ms` | `1000` | Debounce interval for file system events in milliseconds. |
| `parallelism` | `4` | Number of concurrent rule executions. |
| `dry_run_global` | `false` | Run all operations in dry-run mode. |

## Notes

- All paths may use `~` or environment variables.
- Missing optional fields are automatically set to defaults.
- The configuration file must be valid YAML and UTF-8 encoded.

## Future Scope

Planned extensions for the configuration file include:

- **Environment Profiles:** Allow multiple configurations (e.g., work, personal).
- **Logging Levels:** Add structured logging with verbosity control.
- **Rule Caching:** Cache matched file results for faster re-runs.
- **Metrics Integration:** Emit Prometheus-style metrics for automation or CI.
- **Cloud Sync Settings:** Synchronize configuration between devices.

