# Rule Specification v0.1

## Overview
- Rules are yaml files tha
Rules describe the follwing:
- files to match
- actions to perform
Rules are written in YAML and loaded from `~/.autoorg/rules.d/`.



## Rule structure
```
Rule
 ├─ name: string (required)
 ├─ enabled: bool (default: true)
 ├─ match: Match
 │   ├─ folder: string (required)
 │   ├─ extensions: [string]
 │   ├─ older_than_days: int
 │   └─ name_contains: [string]
 ├─ action: Action
 │   ├─ type: string (required)
 │   └─ target: string (optional)
 └─ options: Options
     ├─ dry_run: bool (default: false)
     └─ log: bool (default: true)
```

## Example

```yaml
name: "Move PDFs to Documents"
enabled: true

match:
  folder: "~/Downloads"
  extensions: [".pdf"]
  older_than_days: 7

action:
  type: move
  target: "~/Documents/PDFs"

options:
  dry_run: false
  log: true
```

## Behavior

1. Each rule operates on the specified `folder`.
2. Files matching all criteria under `match` are selected.
3. The action is applied in order of rule load.
4. If `dry_run` is true, no files are changed.
5. Operations are logged if `log` is enabled.

## Supported Actions

- **move** – Move file to `target` folder.
- **copy** – Copy file to `target` folder.
- **delete** – Remove file (moves to `.autoorg_trash`).
- **rename** – Rename file using pattern (future extension).

## Example Rules Directory

```
~/.autoorg/rules.d/
 ├─ move_pdfs.yaml
 ├─ archive_zips.yaml
 └─ delete_temp.yaml
```

## Notes

- Rules are evaluated sequentially by filename order.
- Conflicting actions on the same file are resolved by the first matching rule.
- Paths may use `~` or absolute paths.

# Future Scope

The following features are planned for future releases beyond the MVP:

- **Rule Priorities:** Control execution order when multiple rules match a file.
- **Regex and Metadata Matching:** Allow matching based on regex patterns or file metadata (e.g., EXIF, MIME type).
- **File Watchers:** Automatically apply rules when new files appear in watched directories.
- **Scheduling:** Enable periodic cleanup via cron-like scheduling.
- **Rollback and History:** Maintain an action log to undo previous changes.
- **Compression and Archiving:** Support batch zipping or archiving rules.
- **CLI and TUI Enhancements:** Interactive terminal interface for live rule execution and log viewing.
- **Cross-Platform Sync:** Optionally sync rules and logs across machines.

