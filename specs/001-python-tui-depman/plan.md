# Implementation Plan: depman — TUI Python Dependency Manager

**Branch**: `001-python-tui-depman` | **Date**: 2026-02-23 | **Spec**: [spec.md](file:///home/eslam/coding/projects/depman/specs/001-python-tui-depman/spec.md)  
**Input**: Feature specification from `/specs/001-python-tui-depman/spec.md`

## Summary

Build a terminal user interface (TUI) application in Go for managing Python project dependencies. The tool detects project files, virtual environments, and package managers (`uv`/`pip`), then provides a keyboard-driven dashboard to view, add, remove, update, and search packages. Every mutation triggers a full atomic rewrite of the dependency file. The UI uses the Bubble Tea (Elm architecture) framework with Tokyo Night theming via Lip Gloss.

## Technical Context

**Language/Version**: Go 1.22+  
**Primary Dependencies**: Bubble Tea (TUI framework), Lip Gloss (styling), Bubbles (input/list components), BurntSushi/toml (TOML parsing)  
**Storage**: Local files only — `requirements.txt`, `pyproject.toml`, `~/.config/depman/config.toml`  
**Testing**: `go test` with table-driven tests  
**Target Platform**: Linux, macOS (Unix-like systems)  
**Project Type**: CLI / TUI application  
**Performance Goals**: Dashboard renders in <1s, UI interactions <100ms perceived latency, PyPI search results within 2s  
**Constraints**: No external database, no daemon process, single-process, stdin/stdout TUI  
**Scale/Scope**: Single Python project at a time, personal/open-source tool

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

Constitution is unconfigured (template placeholders only). No project-specific gates to enforce. Proceeding with standard software engineering best practices:

- ✅ Single-purpose tool (dependency management only)
- ✅ No unnecessary abstractions — direct file I/O, subprocess calls, HTTP to PyPI
- ✅ Atomic file writes for safety
- ✅ Clean separation: detection → parsing → execution → file sync

## Project Structure

### Documentation (this feature)

```text
specs/001-python-tui-depman/
├── plan.md              # This file
├── research.md          # Phase 0: Technology research
├── data-model.md        # Phase 1: Entity definitions
├── quickstart.md        # Phase 1: Test scenarios
├── contracts/           # Phase 1: CLI interface contract
└── tasks.md             # Task list (already generated)
```

### Source Code (repository root)

```text
depman/
├── main.go                    # Entrypoint — bootstraps Bubble Tea
├── cmd/
│   └── root.go                # CLI flags, startup orchestration
├── tui/
│   ├── model.go               # Root Bubble Tea model, app state
│   ├── dashboard.go           # Main dashboard (two-panel layout)
│   ├── init.go                # Project initialization screen
│   ├── search.go              # Full PyPI search panel
│   ├── help.go                # Help overlay
│   ├── keymap.go              # Centralized vim keybinding handler
│   ├── add_package.go         # Add package flow
│   ├── remove_package.go      # Remove package flow
│   ├── update_package.go      # Single update flow
│   ├── update_all.go          # Bulk update flow
│   ├── venv_prompt.go         # Virtualenv creation prompt
│   └── components/
│       ├── statusbar.go       # Status bar component
│       ├── package_list.go    # Scrollable package list
│       ├── outdated_list.go   # Outdated packages with diff
│       ├── panel.go           # Bordered panel wrapper
│       ├── search_input.go    # Inline search input
│       ├── search_results.go  # Search results popup
│       ├── confirm.go         # Confirmation prompt
│       ├── progress.go        # Progress indicator
│       ├── notification.go    # Toast notifications
│       └── error.go           # Inline error display
├── pkg/
│   ├── detector/
│   │   ├── project.go         # Dependency file detection
│   │   └── init_project.go    # pyproject.toml / requirements.txt creation
│   ├── parser/
│   │   ├── requirements.go    # requirements.txt parser
│   │   ├── pyproject.go       # pyproject.toml parser (read + section-safe write)
│   │   ├── writer.go          # Atomic file writer (temp + rename)
│   │   └── sync.go            # Post-operation sync orchestrator
│   ├── pypi/
│   │   ├── client.go          # PyPI JSON API client
│   │   ├── search.go          # Debounced search
│   │   ├── version.go         # Semver diff calculator
│   │   └── cache.go           # In-memory response cache
│   ├── pip/
│   │   ├── runner.go          # pip/uv subprocess wrapper
│   │   ├── packages.go        # Parse `pip list` output
│   │   └── outdated.go        # Parse `pip list --outdated` output
│   └── env/
│       ├── virtualenv.go      # Virtualenv detection
│       ├── manager.go         # Package manager detection
│       └── create.go          # Virtualenv creation
└── config/
    ├── config.go              # Config struct + TOML loader
    ├── theme.go               # Tokyo Night color palette
    └── defaults.go            # Default configuration values
```

**Structure Decision**: Single-project CLI/TUI layout. All Go source under root `main.go` + `cmd/`, `tui/`, `pkg/`, `config/` packages. No backend/frontend split — this is a single-binary terminal application. `pkg/` contains reusable, testable library code; `tui/` contains Bubble Tea-specific UI models and components.

## Complexity Tracking

No constitution violations to justify — all design choices follow standard Go project patterns.
