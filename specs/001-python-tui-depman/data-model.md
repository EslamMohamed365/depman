# Data Model: depman

**Branch**: `001-python-tui-depman` | **Date**: 2026-02-23

## Entities

### Package

Represents a Python package (installed or search result).

| Field | Type | Description |
|-------|------|-------------|
| Name | string | PyPI package name (normalized lowercase) |
| InstalledVersion | string | Currently installed version (empty if not installed) |
| LatestVersion | string | Latest version available on PyPI |
| Description | string | Short package description from PyPI |
| DiffType | enum | `patch` / `minor` / `major` / `none` / `unknown` |
| IsOutdated | bool | Whether a newer version is available |

**Relationships**: Belongs to a DependencyFile (when installed). Displayed in PackageList and OutdatedList.

---

### DependencyFile

Represents the project's dependency specification file.

| Field | Type | Description |
|-------|------|-------------|
| Path | string | Absolute path to the file |
| Type | enum | `pyproject_toml` / `requirements_txt` |
| Packages | []Package | List of declared dependencies |
| RawContent | []byte | Raw file content (for pyproject.toml section preservation) |

**Validation rules**:
- Path must exist and be readable/writable
- Type determined by filename detection
- For `pyproject.toml`: only `[project.dependencies]` section is managed

---

### Environment

Represents the resolved Python execution environment.

| Field | Type | Description |
|-------|------|-------------|
| Type | enum | `virtualenv` / `system` |
| Path | string | Path to the virtualenv or system Python |
| PythonBin | string | Path to the Python binary |
| IsActive | bool | Whether this was the shell's active `$VIRTUAL_ENV` |
| IsBroken | bool | True if virtualenv exists but interpreter is missing |

**State transitions**:
- `not_found` → `created` (user chooses to create `.venv`)
- `not_found` → `ignored` (user chooses system Python)
- `broken` → `recreated` (user chooses to recreate)
- `broken` → `fallback` (user falls back to system Python)

---

### PackageManager

Represents the detected or configured package management tool.

| Field | Type | Description |
|-------|------|-------------|
| Type | enum | `uv` / `pip` / `pip3` |
| BinPath | string | Absolute path to the binary |
| IsPreferred | bool | True if explicitly set in config |

**Relationships**: Used by all package operations. Displayed in StatusBar.

---

### AppConfig

Represents user configuration from `~/.config/depman/config.toml`.

| Field | Type | Description |
|-------|------|-------------|
| PreferredManager | string | `"uv"` / `"pip"` / `"pip3"` (default: auto-detect) |
| PyPIMirror | string | PyPI base URL (default: `https://pypi.org`) |
| ThemeName | string | Color theme name (default: `"tokyo-night"`) |

---

### AppState (TUI Root State)

Represents the full application state for the Bubble Tea model.

| Field | Type | Description |
|-------|------|-------------|
| Screen | enum | `init` / `dashboard` / `search` / `help` |
| DependencyFile | DependencyFile | Current project's dep file |
| Environment | Environment | Resolved Python environment |
| PackageManager | PackageManager | Detected package manager |
| InstalledPackages | []Package | All installed packages |
| OutdatedPackages | []Package | Only outdated packages |
| Config | AppConfig | User configuration |
| ActivePanel | enum | `installed` / `outdated` |
| StatusMessage | string | Current status/notification text |
| IsLoading | bool | Whether an async operation is running |

## Relationships Diagram

```
AppState
├── DependencyFile
│   └── []Package (declared in file)
├── Environment
├── PackageManager
├── InstalledPackages → []Package (from pip list)
├── OutdatedPackages → []Package (from pip list --outdated)
└── AppConfig
```
