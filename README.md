# depman 🐍

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/github/go-mod/go-version/eslam/depman)](https://go.dev/)
[![TUI: Bubble Tea](https://img.shields.io/badge/TUI-Bubble%20Tea-00ADD8)](https://github.com/charmbracelet/bubbletea)
[![Release](https://img.shields.io/github/v/release/eslam/depman)](https://github.com/eslam/depman/releases/latest)
[![CI](https://github.com/eslam/depman/actions/workflows/test.yml/badge.svg)](https://github.com/eslam/depman/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/eslam/depman)](https://goreportcard.com/report/github.com/eslam/depman)

A keyboard-first Terminal User Interface (TUI) for managing Python project dependencies. Built with Go, Bubble Tea, and Lip Gloss.

![depman demo](demo.gif)

## Why depman?

Managing Python dependencies often feels like a context-switching marathon. You're jumping between `pip install`, `pip list --outdated`, and manually editing `pyproject.toml` or `requirements.txt`.

**depman** brings everything into one unified, keyboard-driven workflow. See what's outdated, search PyPI, and update your environment without ever reaching for your mouse or leaving your terminal.

## Key Features

- **Auto-Detection** - Instantly finds `pyproject.toml` or `requirements.txt` and your virtual environment
- **Lightning Fast** - Powered by `uv` (falls back to `pip`) for near-instant package operations
- **Vim-Native** - Navigate with `h/j/k/l`, jump with `gg/G`, and search with `/`
- **Visual Semver** - Color-coded updates (🟢 patch, 🟡 minor, 🔴 major) let you assess risk at a glance
- **Real-time Search** - Search PyPI with live results and package descriptions
- **Tokyo Night Theme** - A beautiful, eye-friendly dark theme out of the box

## Installation

### Binary (Recommended)

Download the latest binary for your platform from the [Releases](https://github.com/eslam/depman/releases) page.

```bash
# Linux (x86_64)
curl -L -o depman https://github.com/eslam/depman/releases/latest/download/depman-linux-amd64
chmod +x depman
sudo mv depman /usr/local/bin/

# macOS (Apple Silicon)
curl -L -o depman https://github.com/eslam/depman/releases/latest/download/depman-darwin-arm64
chmod +x depman
sudo mv depman /usr macOS (Intel)
curl -L/local/bin/

# -o depman https://github.com/eslam/depman/releases/latest/download/depman-darwin-amd64
chmod +x depman
sudo mv depman /usr/local/bin/
```

### From Source

Requires [Go](https://go.dev/) 1.25 or later.

```bash
# Clone the repository
git clone https://github.com/eslam/depman.git
cd depman

# Build
go build -o depman .

# Or install directly
go install .
```

## Quick Start

Just run `depman` in the root of your Python project:

```bash
depman
```

If no project is found, `depman` will help you initialize one.

## Keybindings

| Key | Action |
|-----|--------|
| `j` / `down` | Move down |
| `k` / `up` | Move up |
| `g` | Go to first item |
| `G` | Go to last item |
| `Ctrl+d` | Page down |
| `Ctrl+u` | Page up |
| `a` | Add a new package |
| `d` / `x` | Remove selected package |
| `u` | Update selected package |
| `U` | Update all outdated packages |
| `/` / `s` | Search PyPI online |
| `Tab` | Switch between panels |
| `?` | Show help menu |
| `q` / `Esc` | Quit |

## Configuration

`depman` looks for a configuration file at `~/.config/depman/config.toml`.

```toml
# Preferred package manager: "uv" or "pip"
package_manager = "uv"

# Log level: "debug", "info", "warn", "error"
log_level = "info"
```

## Development

### Prerequisites

- Go 1.25+
- Python 3.8+ with a virtual environment

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run a specific test
go test -run TestFunctionName ./...
```

### Building

```bash
go build -o depman .
```

### Linting

```bash
go vet ./...
go fmt ./...
```

## Tech Stack

- **Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Language**: [Go](https://go.dev/)

## Contributing

Contributions are welcome! Whether it's a bug report, a feature request, or a pull request, feel free to get involved.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and development process.

## Roadmap

- [ ] Package version pinning support
- [ ] Requirements.txt editor
- [ ] Virtual environment creation/deletion
- [ ] Dependency tree visualization
- [ ] Audit security vulnerabilities

## License

Distributed under the MIT License. See [LICENSE](LICENSE) for more information.

---

Built with ❤️ by [Eslam Mohamed](https://github.com/EslamMohamed365)
