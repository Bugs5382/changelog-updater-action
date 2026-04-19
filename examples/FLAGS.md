# 🎛️ Flags

| Flag        | Short | Description                                                      | Default |
|-------------|-------|------------------------------------------------------------------|---------|
| `--tag`     | `-t`  | Release tag name, e.g. `v1.2.0`. **Required.**                   | —       |
| `--notes`   | `-n`  | Release notes body (markdown). **Required.**                     | —       |
| `--path`    | `-p`  | Directory (relative to the repo root) containing `CHANGELOG.md`. | `.`     |
| `--date`    |       | Release date injected into the version header (`YYYY-MM-DD`).    | today   |
| `--cleanup` |       | Fix any issues inside current CHANGELOG.md file                  |         |
| `--diff`    |       | Show the diff (if any) of changes.                               | `false` |
| `--dry`     |       | Dry run — parse and log without modifying `CHANGELOG.md`.        | `false` |
| `--verbose` | `-v`  | Enable debug level logging.                                      | `false` |
