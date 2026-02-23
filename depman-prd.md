# PRD: depman — A TUI Dependency Manager

**Version:** 1.0  
**Author:** Personal / Open Source  
**Stack:** Go + Bubble Tea + Lip Gloss  
**Status:** Draft

---

## 1. Overview

`depman` is a terminal user interface (TUI) application for managing Python project dependencies. It allows developers to add, remove, update, and search packages — both locally and online — using vim-style keybindings. It also provides a clear, scannable view of outdated dependencies.

The tool is designed as a personal productivity tool built for daily use, intended to be open-sourced after reaching a stable state.

---

## 2. Problem Statement

Managing Python dependencies today involves switching between multiple commands (`pip install`, `pip uninstall`, `pip list --outdated`, manually editing `requirements.txt` or `pyproject.toml`), multiple tools, and multiple terminal windows. There is no single, keyboard-driven interface that lets a developer see, search, and act on dependencies in one place without leaving the terminal.

---

## 3. Goals

- Provide a fast, keyboard-first TUI for Python dependency management.
- Support vim motions for navigation and actions.
- Visualize outdated packages clearly (name, current version, latest version).
- Allow adding, removing, and updating packages interactively.
- Support searching for packages both locally (installed) and online (PyPI).
- Apply changes by both modifying the dependency file (`requirements.txt` / `pyproject.toml`) and running the underlying pip command.

---

## 4. Non-Goals (v1)

- No support for non-Python ecosystems (planned for future versions).
- No vulnerability scanning or security auditing.
- No changelog or release notes preview.
- No dependency tree or transitive dependency visualization.
- No monorepo / multi-project scanning (single project at a time).
- No GUI, web UI, or daemon process.
- No conda / poetry support (future versions).

---

## 5. Target User

The primary user is a solo Python developer (the author) who works across multiple projects, lives in the terminal, and values keyboard efficiency. Secondary users are developers who discover the open-source project and share similar workflows.

---

## 6. Core Features

### 6.1 Project Detection & Initialization

On launch, `depman` scans the CWD for a Python project dependency file in the following priority order:

1. `pyproject.toml` — **preferred format**
2. `requirements.txt`
3. `requirements/*.txt` (e.g., `requirements/base.txt`)

**If no project file is found**, `depman` does not exit with an error. Instead, it presents an interactive initialization prompt:

```
No Python project found in current directory.
Would you like to initialize one?

  > Create pyproject.toml        (recommended)
    Create requirements.txt      (simple)
    Exit
```

**Create pyproject.toml flow:**
- Prompts for: project name (defaults to CWD folder name), version (default `0.1.0`), Python version requirement.
- Generates a minimal, valid `pyproject.toml` with `[project]` and `[project.dependencies]` sections.
- Optionally prompts to create a virtualenv immediately after (see 6.2).

**Create requirements.txt flow:**
- Creates an empty `requirements.txt` in the CWD.
- Optionally prompts to create a virtualenv immediately after.

### 6.2 Virtual Environment Detection & Creation

Before any package operation, `depman` resolves the active Python environment in the following priority order:

1. **Active shell virtualenv** — if `$VIRTUAL_ENV` is set, use it.
2. **Local `.venv/` directory** — check for `.venv/bin/python` in the CWD.
3. **Local `venv/` directory** — check for `venv/bin/python` in the CWD.
4. **System Python** — fall back to the system `python` / `python3`.

The detected environment is shown in the status bar (e.g., `.venv · ⚡ uv` or `system · pip`). All `pip` and `uv` commands are scoped to the resolved environment.

**Creating a new virtualenv:**

If no virtualenv is detected, `depman` shows a prompt in the status bar:

```
No virtualenv found.  [C] Create .venv   [I] Ignore
```

Pressing `C` creates a new virtualenv in `.venv/` using:
- `uv venv .venv` if `uv` is available.
- `python -m venv .venv` otherwise.

After creation, `depman` automatically activates it internally (sets `$VIRTUAL_ENV` for subprocess calls) and updates the status bar. The user does not need to manually source the activate script.

If a virtualenv is found but appears broken (missing interpreter), `depman` warns the user with an option to recreate it or fall back to system Python.

### 6.3 Package Manager Detection

`depman` supports both `pip` and `uv` as package managers. Detection order:

1. **`uv` preferred** — if `uv` is available in `$PATH`, use it by default. `uv` is significantly faster for installs and resolves dependencies more reliably.
2. **`pip` fallback** — if `uv` is not found, fall back to `pip` / `pip3`.
3. **Manual override** — user can pin a preferred package manager in `~/.config/depman/config.toml`.

Command mapping:

| Action | uv command | pip command |
|---|---|---|
| Install | `uv pip install <pkg>` | `pip install <pkg>` |
| Uninstall | `uv pip uninstall <pkg>` | `pip uninstall <pkg> -y` |
| Upgrade | `uv pip install --upgrade <pkg>` | `pip install --upgrade <pkg>` |
| List installed | `uv pip list` | `pip list` |
| List outdated | `uv pip list --outdated` | `pip list --outdated` |

The active package manager is shown in the status bar (e.g., `⚡ uv` or `pip`).

### 6.4 Main Dashboard View

The primary screen shown on launch. It contains three sections:

**Installed Packages Panel**
- Lists all installed packages (name + current version).
- Highlights packages that have a newer version available.

**Outdated Packages Panel**
- A dedicated list showing only outdated packages.
- Columns: Package Name | Current Version | Latest Version | Diff Type (patch / minor / major — color-coded).

**Status Bar**
- Shows: project path, detected file type, total packages, outdated count.
- Shows available keybindings contextually.

### 6.5 Package Actions

All actions are triggered via keybindings from the main view or a focused package:

| Action | Key |
|---|---|
| Add package | `a` |
| Remove package | `d` or `x` |
| Update selected package | `u` |
| Update all outdated | `U` |
| Search online (PyPI) | `/` or `s` |
| Confirm action | `Enter` |
| Cancel | `Esc` or `q` |

**Add Package Flow:**
1. Press `a` to open an inline search input.
2. User types a package name; results are fetched from PyPI in real time (debounced).
3. Results shown in a popup list with package name, latest version, and short description.
4. User selects a package and confirms.
5. `depman` runs `pip install <package>` and updates the dependency file.

**Remove Package Flow:**
1. Navigate to a package, press `d`.
2. A confirmation prompt appears.
3. On confirm: runs `pip uninstall <package> -y` and removes the entry from the dependency file.

**Update Package Flow:**
1. Navigate to an outdated package, press `u`.
2. Shows: current version → latest version.
3. On confirm: runs `pip install --upgrade <package>` and updates the version pin in the dependency file.

**Update All Flow:**
1. Press `U` from the dashboard.
2. Shows a summary: "X packages will be updated."
3. On confirm: runs updates sequentially with a progress indicator.

### 6.6 Online Search (PyPI)

- Accessible via `/` or `s` from anywhere.
- Opens a search panel with a text input.
- Queries PyPI's JSON API in real time.
- Results show: package name, latest version, short description, download count (if available).
- Selecting a result allows immediate install.

### 6.7 Vim Motions & Navigation

| Motion | Action |
|---|---|
| `j` / `k` | Move down / up |
| `gg` | Jump to top |
| `G` | Jump to bottom |
| `Ctrl+d` | Half-page down |
| `Ctrl+u` | Half-page up |
| `Tab` | Switch between panels |
| `:` | Open command mode (future) |
| `q` | Quit |
| `?` | Show help |

### 6.8 Dependency File Sync — Full Rewrite on Every Change

`depman` owns the dependency file. Every time a package is added, removed, or updated, `depman` performs a **full rewrite** of the dependency file rather than patching individual lines. This ensures the file is always clean, consistently formatted, and in sync with the actual installed state.

**Rewrite strategy:**

1. Run the `pip` / `uv` command (install / uninstall / upgrade).
2. After the command succeeds, query the full list of installed packages from the environment (`uv pip list` or `pip list`).
3. Reconstruct the entire dependency file from scratch using the current in-memory package list.
4. Write atomically: write to a temp file first, then rename to the target path (prevents corruption on failure).

**Format rules per file type:**

For `requirements.txt`:
```
# Generated by depman — do not edit manually
package-a==1.2.3
package-b==4.5.6
```

For `pyproject.toml` — only the `[project.dependencies]` array is rewritten; all other sections (metadata, build system, tool configs) are preserved exactly as-is:
```toml
[project.dependencies]
# Generated by depman
"package-a==1.2.3",
"package-b==4.5.6",
```

**Important:** `depman` adds a comment header to dependency files it manages so users know the file is being maintained by the tool. It never touches other sections of `pyproject.toml`.

**On failure:** If the `pip`/`uv` command fails, the dependency file is NOT rewritten. The error is shown inline and the file remains unchanged.

---

## 7. UX & Design Principles

- **Speed first:** All interactions should feel instant. Network calls (PyPI) are async and never block the UI.
- **Minimal chrome:** No unnecessary popups or modal dialogs. Actions happen inline.
- **Clear feedback:** Every action has a visible outcome — success, failure, or progress.
- **Vim-native feel:** Users should never need to reach for the mouse.
- **Color coding:**
  - 🟢 Green: up to date
  - 🟡 Yellow: patch update available
  - 🟠 Orange: minor update available
  - 🔴 Red: major update available

---

## 7.1 Theme — Tokyo Night

The default and primary theme is **Tokyo Night**, implemented via Lip Gloss color definitions. All UI components — panels, status bar, highlights, borders, and text — use this palette.

### Color Palette

| Role | Name | Hex |
|---|---|---|
| Background | Night BG | `#1a1b26` |
| Background (elevated) | Storm BG | `#24283b` |
| Background (highlight) | Highlight BG | `#2e3250` |
| Border | Subtle | `#414868` |
| Foreground | Default text | `#c0caf5` |
| Foreground (dim) | Comments / hints | `#565f89` |
| Accent / Selection | Blue | `#7aa2f7` |
| Up to date | Green | `#9ece6a` |
| Patch update | Teal | `#2ac3de` |
| Minor update | Yellow | `#e0af68` |
| Major update | Red | `#f7768e` |
| Package name | Purple | `#bb9af7` |
| Version number | Cyan | `#7dcfff` |
| Search highlight | Orange | `#ff9e64` |
| Success notification | Green | `#9ece6a` |
| Error notification | Red | `#f7768e` |
| Warning notification | Yellow | `#e0af68` |
| uv indicator | Magenta | `#bb9af7` |
| pip indicator | Dim foreground | `#565f89` |

### Component Styling

**Status Bar** — bottom of screen, `#24283b` background:
```
 depman  ~/projects/myapp  .venv  ⚡ uv  12 packages  3 outdated  ? help
```
- Project path in `#7aa2f7`
- Venv name in `#9ece6a`
- `⚡ uv` in `#bb9af7`, `pip` in `#565f89`
- Outdated count in `#f7768e` if > 0, else `#565f89`

**Panel Borders** — `#414868` by default, `#7aa2f7` when panel is focused.

**Selected Row** — background `#2e3250`, text `#c0caf5`, left indicator `▶` in `#7aa2f7`.

**Package Name Column** — `#bb9af7`

**Version Columns** — current version in `#7dcfff`, latest version colored by update severity (see color coding above).

**Search Input** — border `#7aa2f7`, cursor `#ff9e64`, placeholder text `#565f89`.

**Notifications** (bottom-right toast):
- Success: `#9ece6a` text on `#24283b` bg
- Error: `#f7768e` text on `#24283b` bg

### Lip Gloss Reference (Go)

```go
var (
    ColorBG        = lipgloss.Color("#1a1b26")
    ColorBGElevated = lipgloss.Color("#24283b")
    ColorBGHighlight = lipgloss.Color("#2e3250")
    ColorBorder    = lipgloss.Color("#414868")
    ColorFG        = lipgloss.Color("#c0caf5")
    ColorFGDim     = lipgloss.Color("#565f89")
    ColorBlue      = lipgloss.Color("#7aa2f7")
    ColorGreen     = lipgloss.Color("#9ece6a")
    ColorTeal      = lipgloss.Color("#2ac3de")
    ColorYellow    = lipgloss.Color("#e0af68")
    ColorRed       = lipgloss.Color("#f7768e")
    ColorPurple    = lipgloss.Color("#bb9af7")
    ColorCyan      = lipgloss.Color("#7dcfff")
    ColorOrange    = lipgloss.Color("#ff9e64")
)
```

---

## 8. Technical Architecture

### Stack

| Layer | Technology |
|---|---|
| Language | Go |
| TUI Framework | Bubble Tea (Elm architecture) |
| Styling | Lip Gloss |
| PyPI API | PyPI JSON API (`https://pypi.org/pypi/<pkg>/json`) |
| File Parsing | Go stdlib + custom parsers for `.txt` and TOML |
| Package Commands | `os/exec` → `pip` subprocess calls |

### Project Structure (Proposed)

```
depman/
├── main.go
├── cmd/              # CLI entrypoint (cobra or flags)
├── tui/              # Bubble Tea models and views
│   ├── dashboard.go
│   ├── search.go
│   └── components/
├── pkg/
│   ├── detector/     # Project & dep file detection
│   ├── parser/       # requirements.txt + pyproject.toml parsing
│   ├── pypi/         # PyPI API client
│   └── pip/          # pip subprocess wrapper
└── config/           # User config (keybindings, theme, etc.)
```

### Key Constraints

- `pip` must be available in the current environment (or virtualenv).
- Network access required only for PyPI search and version checks.
- All file writes are atomic (write to temp file, then rename).

---

## 9. Configuration

A `~/.config/depman/config.toml` file will allow overriding:
- Preferred package manager (`uv` | `pip` | `pip3`) — default: `uv` if available
- Virtualenv search paths (additional custom locations)
- Color theme — default: `tokyo-night` (future themes: `catppuccin`, `gruvbox`, `nord`)
- Individual color overrides (any token from the Tokyo Night palette)
- Keybinding overrides (future)
- PyPI mirror URL (for corporate environments)

---

## 10. Error Handling

| Scenario | Behavior |
|---|---|
| No Python project in CWD | Offer to create `pyproject.toml` or `requirements.txt` |
| No pip or uv found | Show error with install hint for both tools |
| No virtualenv found | Prompt user to create `.venv` or continue with system Python |
| Broken virtualenv detected | Warn user, offer to recreate or fall back to system Python |
| PyPI unreachable | Show cached data + offline indicator |
| Dependency file parse error | Show error with line number, do not modify file |
| pip / uv command fails | Show stderr inline, skip file rewrite |
| uv not found, configured as default | Warn and auto-fall back to pip |
| File write fails (permissions, disk) | Show error, leave original file untouched |

---

## 11. Milestones

| Milestone | Scope |
|---|---|
| **M1 — Foundation** | Project detection, `pyproject.toml` / `requirements.txt` creation flow, virtualenv detection & creation, package manager detection (uv/pip) |
| **M2 — Outdated View** | Parse existing dep files, fetch latest versions from PyPI, show outdated panel with color-coded semver diff |
| **M3 — Core Actions** | Add, remove, update (single package), full dep file rewrite on every change |
| **M4 — Search** | Online PyPI search with real-time results |
| **M5 — Vim Motions** | Full vim navigation, `gg`, `G`, `Ctrl+d/u`, panel switching |
| **M6 — pyproject.toml** | Full read/write support for `pyproject.toml`, preserve non-dependency sections |
| **M7 — Polish** | Help screen, config file, full error handling, README, open source release |

---

## 12. Open Questions

- Should "update all" run updates in parallel or sequentially? — Sequential for safety in v1.
- Should there be a dry-run mode that shows what would change without applying it? — Nice to have, consider for M7.
- Should `uv sync` / `uv lock` be supported for projects using `uv` natively with a lockfile (`uv.lock`)? This is a different and more powerful workflow than `uv pip install`.
- When rewriting `pyproject.toml`, should `depman` preserve user comments inside `[project.dependencies]`? Currently the plan is to overwrite that section fully.
- Should the generated `pyproject.toml` include a `[build-system]` section by default, or keep it minimal?
- Should conda environments be detected passively (i.e., `$CONDA_DEFAULT_ENV`) and shown as a warning since conda is out of scope?

---

*This document will evolve as development begins. Non-goals and open questions will be revisited before v2 planning.*
