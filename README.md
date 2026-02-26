# depman

A keyboard-first TUI application for managing Python project dependencies directly from your terminal.

## Installation

### From Binary

Download a release from the [releases page](https://github.com/eslam/depman/releases):

```bash
# Linux
curl -L -o depman https://github.com/eslam/depman/releases/latest/download/depman-linux-amd64
chmod +x depman
./depman

# macOS (Intel)
curl -L -o depman https://github.com/eslam/depman/releases/latest/download/depman-darwin-amd64
chmod +x depman
./depman

# macOS (Apple Silicon)
curl -L -o depman https://github.com/eslam/depman/releases/latest/download/depman-darwin-arm64
chmod +x depman
./depman
```

### From Source

```bash
# Requires Go 1.25+
go install github.com/eslam/depman@latest

# Or build manually
git clone https://github.com/eslam/depman.git
cd depman
go build -o depman .
./depman
```

### Requirements

- Python 3.x
- `pip` or `uv` (uv is preferred)
- A `pyproject.toml` or `requirements.txt` in your project

## What

`depman` provides a visual, keyboard-driven interface for adding, removing, updating, and searching Python packages — without leaving your terminal. Built for developers who live in the CLI and want a unified workflow for dependency management.

## Quick Start

```bash
# Run the binary
./depman

# Or install globally
go install github.com/eslam/depman@latest
```

On launch, `depman` detects your Python project file (`pyproject.toml` or `requirements.txt`) and virtual environment automatically.

## Features

| Feature | Description |
|---------|-------------|
| **Project Detection** | Auto-detects `pyproject.toml`, `requirements.txt`, or `requirements/*.txt` |
| **Vim Navigation** | `j/k` to move, `gg`/`G` jump, `/` to search |
| **Outdated View** | Color-coded: 🟢 patch · 🟡 minor · 🔴 major |
| **Package Actions** | Add, remove, update — all via keyboard |
| **PyPI Search** | Real-time search with package descriptions |
| **Dual Package Managers** | Works with both `uv` (preferred) and `pip` |
| **Tokyo Night Theme** | Eye-friendly dark theme out of the box |

## Keybindings

| Key | Action |
|-----|--------|
| `a` | Add package |
| `d` / `x` | Remove package |
| `u` | Update selected package |
| `U` | Update all outdated |
| `/` / `s` | Search PyPI |
| `Tab` | Switch panels |
| `?` | Show help |
| `q` | Quit |

## Tech Stack

- **Language**: Go
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Elm architecture)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)

## Project Structure

```
depman/
├── main.go          # Entry point
├── cmd/             # CLI commands
├── tui/             # Bubble Tea models & views
│   ├── dashboard.go
│   ├── search.go
│   └── model.go
├── pkg/
│   ├── detector/    # Project & file detection
│   ├── parser/      # requirements.txt / pyproject.toml
│   ├── pypi/        # PyPI API client
│   └── pip/         # pip/uv subprocess wrapper
└── config/          # User configuration
```

## Configuration

Config file: `~/.config/depman/config.toml`

```toml
package_manager = "uv"   # "uv" or "pip"
```

## Why Another Tool?

Managing Python dependencies today means juggling `pip install`, `pip uninstall`, `pip list --outdated`, and manually editing files. `depman` brings everything into one keyboard-driven interface — see your outdated packages, search PyPI, and apply changes without switching contexts.

## Status

This is a personal productivity tool in active use. It will be open-sourced once it reaches a stable state.

---

Built with Bubble Tea + Lip Gloss
