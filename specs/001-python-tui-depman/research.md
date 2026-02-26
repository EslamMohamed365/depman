# Research: depman — TUI Python Dependency Manager

**Branch**: `001-python-tui-depman` | **Date**: 2026-02-23

## R1: Bubble Tea Architecture Patterns

**Decision**: Use the standard Elm architecture (Model → Update → View) with nested models for each screen/flow.

**Rationale**: Bubble Tea enforces the Elm pattern natively. Each screen (dashboard, search, help, init) is a separate model with its own `Init`/`Update`/`View`. The root model delegates to the current active screen. This aligns with the PRD's screen-based design and keeps each flow independently testable.

**Alternatives considered**:
- Single monolithic model — rejected, would become unmaintainable with 10+ flows
- External state management — rejected, Bubble Tea's built-in model pattern is sufficient for a single-user TUI

**Key patterns**:
- Use `tea.Cmd` for async operations (PyPI fetches, pip commands) — never block the UI
- Use `tea.Msg` types for inter-component communication
- Use `tea.WindowSizeMsg` for responsive layout
- Bubbles library provides reusable input, list, spinner, and viewport components

---

## R2: PyPI JSON API

**Decision**: Use the PyPI JSON API at `https://pypi.org/pypi/<package>/json` for package info and `https://pypi.org/simple/` index for search.

**Rationale**: The PyPI JSON API is the official, stable, unauthenticated API. It returns package metadata, versions, and descriptions. No API key required. Rate limits are generous for single-user CLI tools.

**Alternatives considered**:
- PyPI XML-RPC API — deprecated, being phased out
- PyPI Simple API (PEP 503) — good for package resolution but doesn't include descriptions or metadata
- Third-party search APIs — unnecessary dependency

**Key patterns**:
- For individual package lookup: `GET https://pypi.org/pypi/<name>/json`
- For search: Use `https://pypi.org/search/?q=<query>` with HTML parsing, or query multiple individual packages. Note: PyPI has no official search JSON API — will need to scrape search results or use the `xmlrpc` `search` method (still available but deprecated). Best approach: query individual packages by prefix match against a cached package list.
- Debounce search input (300ms) to avoid excessive requests
- Cache responses in memory with TTL for offline fallback

---

## R3: TOML Parsing — Preserving Non-Dependency Sections

**Decision**: Use `BurntSushi/toml` (or `pelletier/go-toml/v2`) for reading TOML. For writing, use a surgical approach: parse the full file, modify only the `[project.dependencies]` array, then serialize back.

**Rationale**: The PRD explicitly requires preserving all non-dependency sections of `pyproject.toml` exactly as-is. A full parse → modify → serialize approach ensures other sections (build-system, tool configs, comments outside dependencies) are maintained.

**Alternatives considered**:
- Regex-based line replacement — fragile, breaks on multi-line values or nested tables
- Full TOML rewrite — would lose formatting and comments in other sections
- `pelletier/go-toml/v2` with `SetPath` — promising for surgical edits but may reformat

**Key decision**: Use `pelletier/go-toml/v2` which supports `toml.Marshal` with comment preservation better than `BurntSushi/toml`. For maximum safety, read the raw file, find the `[project.dependencies]` section boundaries, and replace only that slice of bytes, leaving everything else byte-identical.

---

## R4: Atomic File Writes

**Decision**: Write to a temp file in the same directory, then `os.Rename` to the target path.

**Rationale**: `os.Rename` on the same filesystem is atomic on POSIX systems. Writing to a temp file first ensures a crash or error during write never corrupts the original file.

**Alternatives considered**:
- Direct overwrite with `os.WriteFile` — not atomic, risks corruption on crash
- Write to temp dir then move — may fail across filesystem boundaries
- fsync before rename — added safety for `requirements.txt`, but Go's `os.Rename` is sufficient for most cases

**Implementation**:
```go
tmpFile, _ := os.CreateTemp(filepath.Dir(target), ".depman-*")
tmpFile.Write(content)
tmpFile.Close()
os.Rename(tmpFile.Name(), target)
```

---

## R5: Subprocess Execution for pip/uv

**Decision**: Use `os/exec.Command` with environment variable injection for virtualenv scoping.

**Rationale**: Both `pip` and `uv` are external CLI tools. All package operations (install, uninstall, upgrade, list) are executed as subprocesses. The virtualenv is scoped by setting `VIRTUAL_ENV` and prepending `<venv>/bin` to `PATH` in the subprocess environment.

**Alternatives considered**:
- Importing pip as a library — not applicable (Go, not Python)
- Using pip's internal APIs — explicitly unsupported by pip maintainers
- Shelling out via bash — unnecessary indirection, `exec.Command` is cleaner

**Key patterns**:
- Capture both stdout and stderr separately
- Check exit code for success/failure
- Parse stdout for structured output (package lists)
- Show stderr inline on failure
- Never block the TUI — run commands via `tea.Cmd` returning `tea.Msg`

---

## R6: Semver Diff Classification

**Decision**: Parse versions as `MAJOR.MINOR.PATCH` and compare component-by-component.

**Rationale**: The PRD requires color-coding outdated packages by severity (patch/minor/major). Simple numeric comparison of version components is sufficient. Python packages follow PEP 440, which is close to semver for most popular packages.

**Alternatives considered**:
- Using a Go semver library — overkill for simple `x.y.z` comparison
- PEP 440 full parser — complex (epochs, pre-releases, post-releases), but for diff classification, major/minor/patch comparison on the first three segments is good enough

**Implementation**:
- Split version by `.`, compare first 3 segments
- If major differs → major update (red)
- If minor differs → minor update (yellow)
- If patch differs → patch update (teal)
- Handle non-standard versions gracefully (treat as "unknown" severity)

---

## R7: Package Manager Detection and Configuration

**Decision**: Check `$PATH` for `uv` first, fall back to `pip`/`pip3`. Allow override via `~/.config/depman/config.toml`.

**Rationale**: `uv` is significantly faster (10-100x for installs) and is becoming the standard. Preferring it when available gives the best user experience. Config override supports corporate environments or user preferences.

**Alternatives considered**:
- Always use pip — misses the performance benefits of uv
- Always use uv — breaks for users without uv installed
- Auto-install uv — too invasive for a dependency manager

**Config format** (`~/.config/depman/config.toml`):
```toml
[package_manager]
preferred = "uv"  # "uv" | "pip" | "pip3"

[pypi]
mirror = "https://pypi.org"  # Override for corporate mirrors

[theme]
name = "tokyo-night"  # Future: "catppuccin", "gruvbox", "nord"
```
