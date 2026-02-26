# Feature Specification: depman — TUI Python Dependency Manager

**Feature Branch**: `001-python-tui-depman`  
**Created**: 2026-02-23  
**Status**: Draft  
**Input**: User description: "A TUI dependency manager for Python projects with vim-style keybindings, package management via pip/uv, PyPI search, and dependency file sync"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Launch and Environment Setup (Priority: P1)

A developer opens their terminal inside a Python project directory and launches `depman`. The tool automatically detects the project's dependency file (`pyproject.toml` or `requirements.txt`), finds or prompts creation of a virtual environment, and determines whether to use `uv` or `pip` as the package manager. The status bar displays the detected environment, package manager, and package counts. If no project file exists, the tool guides the developer through creating one interactively.

**Why this priority**: Without environment detection and initialization, no other feature can function. This is the core bootstrapping flow that every subsequent interaction depends on.

**Independent Test**: Can be fully tested by launching the tool in directories with and without Python project files and verifying the correct environment is detected, created, or prompted.

**Acceptance Scenarios**:

1. **Given** a directory with `pyproject.toml`, **When** the user launches `depman`, **Then** the tool detects the file, resolves the virtualenv, detects the package manager, and displays the main dashboard with the status bar showing project path, file type, package counts, and environment info.
2. **Given** a directory with `requirements.txt` but no `pyproject.toml`, **When** the user launches `depman`, **Then** `requirements.txt` is detected and used as the dependency file.
3. **Given** a directory with no Python project files, **When** the user launches `depman`, **Then** an interactive prompt offers to create `pyproject.toml` (recommended), `requirements.txt` (simple), or exit.
4. **Given** `$VIRTUAL_ENV` is set in the shell, **When** the user launches `depman`, **Then** the active shell virtualenv is used and shown in the status bar.
5. **Given** no virtualenv is found, **When** the user launches `depman`, **Then** a prompt appears offering to create `.venv` or ignore.
6. **Given** `uv` is available in `$PATH`, **When** the user launches `depman`, **Then** `uv` is used as the package manager and indicated in the status bar with `⚡ uv`.
7. **Given** `uv` is not available but `pip` is, **When** the user launches `depman`, **Then** `pip` is used as fallback and shown in the status bar.

---

### User Story 2 - View Installed and Outdated Packages (Priority: P1)

A developer launches `depman` in a project with existing dependencies and sees a dashboard with two panels: one listing all installed packages with their versions, and another showing only outdated packages with current version, latest version, and a color-coded severity indicator (patch, minor, or major update).

**Why this priority**: Seeing the current state of dependencies is the primary read-only use case and the foundation for all package actions. This delivers immediate value even before add/remove/update are implemented.

**Independent Test**: Can be tested by launching the tool in a project with known installed and outdated packages and verifying the correct display of package data and color coding.

**Acceptance Scenarios**:

1. **Given** a project with 10 installed packages (3 outdated), **When** the user views the dashboard, **Then** the installed panel shows all 10 packages with names and versions, and the outdated panel shows only the 3 outdated packages with current version, latest version, and diff type.
2. **Given** a package has a patch update available, **When** it appears in the outdated panel, **Then** it is color-coded with the patch severity color.
3. **Given** a package has a minor update available, **When** it appears in the outdated panel, **Then** it is color-coded with the minor severity color.
4. **Given** a package has a major update available, **When** it appears in the outdated panel, **Then** it is color-coded with the major severity color.
5. **Given** all packages are up to date, **When** the user views the dashboard, **Then** the outdated panel is empty and the status bar shows "0 outdated".

---

### User Story 3 - Add a Package (Priority: P2)

A developer wants to add a new Python package. They press `a` to open an inline search input, type a package name, see real-time results from PyPI with package names, versions, and descriptions, select one, and confirm. The tool installs the package and updates the dependency file.

**Why this priority**: Adding packages is one of the most common dependency management actions and the first mutating operation users need.

**Independent Test**: Can be tested by pressing `a`, searching for a known package (e.g., "requests"), selecting it, confirming, and verifying the package is installed and the dependency file is updated.

**Acceptance Scenarios**:

1. **Given** the user is on the dashboard, **When** they press `a`, **Then** an inline search input appears.
2. **Given** the search input is open, **When** the user types "requests", **Then** matching results are fetched from PyPI and displayed with package name, latest version, and short description.
3. **Given** search results are displayed, **When** the user selects a result and presses Enter, **Then** the package is installed via the detected package manager and the dependency file is rewritten to include the new package.
4. **Given** a package install fails, **When** the error occurs, **Then** the error message is shown inline and the dependency file is not modified.

---

### User Story 4 - Remove a Package (Priority: P2)

A developer navigates to an installed package, presses `d` or `x`, sees a confirmation prompt, and on confirm the package is uninstalled and removed from the dependency file.

**Why this priority**: Removing packages completes the basic CRUD for dependency management and is essential for project hygiene.

**Independent Test**: Can be tested by navigating to a known installed package, pressing `d`, confirming, and verifying the package is uninstalled and the dependency file no longer lists it.

**Acceptance Scenarios**:

1. **Given** a package is selected in the installed panel, **When** the user presses `d`, **Then** a confirmation prompt appears showing the package name.
2. **Given** the confirmation prompt is shown, **When** the user confirms, **Then** the package is uninstalled and the dependency file is rewritten without it.
3. **Given** the confirmation prompt is shown, **When** the user presses Esc, **Then** the action is cancelled and nothing changes.

---

### User Story 5 - Update Packages (Priority: P2)

A developer navigates to an outdated package, presses `u` to update it to the latest version. Alternatively, they press `U` to update all outdated packages at once. In both cases, the dependency file is updated to reflect new versions.

**Why this priority**: Keeping dependencies up to date is a core workflow and a key differentiator for the tool.

**Independent Test**: Can be tested by navigating to a known outdated package, pressing `u`, confirming, and verifying the package version is updated in both the environment and the dependency file.

**Acceptance Scenarios**:

1. **Given** an outdated package is selected, **When** the user presses `u`, **Then** a prompt shows the current version → latest version and asks for confirmation.
2. **Given** the user confirms a single update, **When** the update runs, **Then** the package is upgraded and the dependency file is rewritten with the new version.
3. **Given** multiple packages are outdated, **When** the user presses `U`, **Then** a summary shows "X packages will be updated" and asks for confirmation.
4. **Given** the user confirms "update all", **When** updates run, **Then** packages are updated sequentially with a progress indicator, and the dependency file is rewritten after all updates complete.
5. **Given** one update fails during "update all", **When** the error occurs, **Then** the error is shown inline, successfully updated packages are reflected in the file, and the failed package retains its previous version.

---

### User Story 6 - Search PyPI Online (Priority: P3)

A developer presses `/` or `s` from anywhere in the application to open a search panel. They type a query, see real-time results from PyPI including package name, latest version, description, and download count (if available), and can install a selected result directly.

**Why this priority**: Online search enhances discoverability but is not essential for managing existing packages. Users can add packages by name without search.

**Independent Test**: Can be tested by pressing `/`, typing a query, verifying results appear from PyPI, and optionally installing a result.

**Acceptance Scenarios**:

1. **Given** the user is on any screen, **When** they press `/` or `s`, **Then** a search panel with text input opens.
2. **Given** the search panel is open, **When** the user types a query, **Then** results are fetched from PyPI in real time (debounced) and displayed.
3. **Given** search results are displayed, **When** the user selects a result and confirms, **Then** the package is installed and added to the dependency file.
4. **Given** PyPI is unreachable, **When** the user searches, **Then** an offline indicator is shown and cached data is used if available.

---

### User Story 7 - Navigate with Vim Motions (Priority: P3)

A developer navigates the entire interface using familiar vim-style keybindings: `j`/`k` for up/down, `gg`/`G` for top/bottom, `Ctrl+d`/`Ctrl+u` for half-page scrolling, `Tab` to switch between panels, and `?` for help.

**Why this priority**: Vim motions are the ergonomic foundation of the tool's UX identity, but basic arrow-key or default Bubble Tea navigation can serve as a temporary fallback.

**Independent Test**: Can be tested by pressing each vim keybinding and verifying the expected navigation behavior in the TUI.

**Acceptance Scenarios**:

1. **Given** a list of packages, **When** the user presses `j`, **Then** the selection moves down one item.
2. **Given** a list of packages, **When** the user presses `k`, **Then** the selection moves up one item.
3. **Given** a list of packages, **When** the user presses `gg`, **Then** the selection jumps to the first item.
4. **Given** a list of packages, **When** the user presses `G`, **Then** the selection jumps to the last item.
5. **Given** the installed panel is focused, **When** the user presses `Tab`, **Then** focus moves to the outdated panel.
6. **Given** any screen, **When** the user presses `?`, **Then** a help screen with all keybindings is displayed.

---

### User Story 8 - Dependency File Sync (Priority: P2)

Every time a package is added, removed, or updated, the tool rewrites the entire dependency file from scratch based on the current installed state. The write is atomic (temp file then rename). For `pyproject.toml`, only the `[project.dependencies]` section is rewritten — all other sections are preserved exactly as-is.

**Why this priority**: File sync ensures the dependency file always reflects the true state. Without it, the tool's changes would be invisible outside the TUI.

**Independent Test**: Can be tested by performing package operations and verifying the resulting dependency file content matches expected format and includes all installed packages.

**Acceptance Scenarios**:

1. **Given** a `requirements.txt` project, **When** a package is added, **Then** the entire file is rewritten with all installed packages in `package==version` format, with a `depman` header comment.
2. **Given** a `pyproject.toml` project with metadata and build-system sections, **When** a package is added, **Then** only the `[project.dependencies]` array is rewritten and all other sections remain unchanged.
3. **Given** a package operation fails (e.g., install fails), **When** the error occurs, **Then** the dependency file is not modified.
4. **Given** a successful package operation, **When** the file is written, **Then** the write is atomic (written to temp file first, then renamed).

---

### Edge Cases

- What happens when the dependency file is modified externally while `depman` is running? — Assumption: `depman` does a full rewrite on every change based on the installed packages, so external changes to the file are overwritten on the next depman operation but installed packages remain unaffected.
- What happens when `pip`/`uv` commands hang or take extremely long? — Assumption: show a progress indicator; no timeout in v1 (user can cancel with Ctrl+C).
- What happens when the user runs `depman` in a directory with both `pyproject.toml` and `requirements.txt`? — `pyproject.toml` takes priority per detection order.
- How does the system handle packages installed outside of `depman` (e.g., manual pip install)? — On next rewrite, all currently installed packages appear in the dependency file since it rebuilds from the full installed list.
- What happens when disk is full during atomic file write? — Error is shown inline, original file is left untouched.
- What happens when a broken virtualenv is detected? — Warning is shown with options to recreate it or fall back to system Python.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST detect Python project dependency files in the current working directory in priority order: `pyproject.toml` → `requirements.txt` → `requirements/*.txt`.
- **FR-002**: System MUST offer interactive project initialization when no dependency file is found, with options to create `pyproject.toml` (recommended) or `requirements.txt`.
- **FR-003**: System MUST detect the active Python environment in priority order: active shell virtualenv (`$VIRTUAL_ENV`) → local `.venv/` → local `venv/` → system Python.
- **FR-004**: System MUST prompt the user to create a `.venv` virtualenv when none is detected, preferring `uv venv` if available, falling back to `python -m venv`.
- **FR-005**: System MUST detect and prefer `uv` as the package manager when available in `$PATH`, falling back to `pip`/`pip3`.
- **FR-006**: System MUST display a main dashboard with an installed packages panel (name + version), an outdated packages panel (name, current version, latest version, diff type), and a status bar.
- **FR-007**: System MUST allow adding packages by pressing `a`, showing real-time PyPI search results with name, version, and description, installing on confirmation, and rewriting the dependency file.
- **FR-008**: System MUST allow removing packages by pressing `d` or `x` with a confirmation prompt, uninstalling the package, and rewriting the dependency file.
- **FR-009**: System MUST allow updating a single package by pressing `u` with a version diff preview and confirmation, and updating all outdated packages by pressing `U` sequentially with a progress indicator.
- **FR-010**: System MUST perform a full rewrite of the dependency file after every successful package operation, based on the current installed package list.
- **FR-011**: System MUST write dependency files atomically (write to temp file, then rename to target path).
- **FR-012**: System MUST preserve all non-dependency sections of `pyproject.toml` during rewrites (metadata, build system, tool configs).
- **FR-013**: System MUST support online search of PyPI via `/` or `s` with real-time, debounced results showing package name, version, description, and download count.
- **FR-014**: System MUST support vim-style navigation: `j`/`k` (up/down), `gg`/`G` (top/bottom), `Ctrl+d`/`Ctrl+u` (half-page scroll), `Tab` (panel switch), `?` (help), `q` (quit).
- **FR-015**: System MUST display the detected environment, package manager, project path, file type, total packages, and outdated count in the status bar.
- **FR-016**: System MUST color-code outdated packages by semver diff type: patch, minor, or major.
- **FR-017**: System MUST NOT rewrite the dependency file when a package operation (install/uninstall/upgrade) fails, and MUST show the error inline.
- **FR-018**: System MUST add a comment header (e.g., `# Generated by depman`) to dependency files it manages.
- **FR-019**: System MUST support a user configuration file at `~/.config/depman/config.toml` for overriding preferred package manager, theme, and PyPI mirror URL.
- **FR-020**: System MUST warn the user when a broken virtualenv is detected and offer to recreate it or fall back to system Python.
- **FR-021**: System MUST show an error with install hints when neither `pip` nor `uv` is found.
- **FR-022**: System MUST show an offline indicator and use cached data when PyPI is unreachable.

### Key Entities

- **Package**: Represents an installed Python package. Key attributes: name, installed version, latest available version, semver diff type (patch/minor/major).
- **Dependency File**: The project's dependency specification file (`requirements.txt` or `pyproject.toml`). Key attributes: file path, file type, list of declared dependencies.
- **Environment**: The Python execution environment. Key attributes: type (virtualenv/system), path, Python version, associated package manager.
- **Package Manager**: The tool used to install/uninstall/upgrade packages. Key attributes: type (`uv` or `pip`), path, availability status.
- **Configuration**: User preferences loaded from `~/.config/depman/config.toml`. Key attributes: preferred package manager, theme, PyPI mirror URL.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can launch the tool and see their project's full dependency status (installed + outdated) within 5 seconds of startup.
- **SC-002**: Users can add, remove, or update a single package in 3 interactions or fewer (open action → select/confirm → done).
- **SC-003**: All package operations correctly update the dependency file to reflect the true installed state with zero manual file editing required.
- **SC-004**: The tool correctly detects the project environment (dependency file, virtualenv, package manager) on first launch without manual configuration in 95% of standard Python project setups.
- **SC-005**: Online search results from PyPI appear within 2 seconds of the user stopping typing.
- **SC-006**: All navigation actions respond instantly (under 100ms perceived latency) and never block the interface during network operations.
- **SC-007**: The tool handles all error scenarios (no project, no virtualenv, no package manager, offline PyPI, failed operations) gracefully without crashing or corrupting files.
- **SC-008**: Users familiar with vim keybindings can navigate the full interface without consulting help within their first session.
- **SC-009**: The dependency file is never left in a corrupted or partial state after any operation, including interrupted ones.

## Assumptions

- The user has Python installed and accessible via `python` or `python3` on their system.
- The target operating system is macOS or Linux (Unix-like systems with standard paths like `/bin/python`).
- Network access is available for PyPI search and version checking, but the tool degrades gracefully without it.
- Only one project is managed at a time (the one in the current working directory).
- `pip` is available inside the detected virtualenv or system Python.
- The tool does not manage transitive dependencies — only directly declared dependencies appear in the file.
- Sequential updates are used for "update all" to ensure safety and clear error attribution.
- The `pyproject.toml` initialization creates a minimal `[project]` section without a `[build-system]` section by default.
- Conda environments are out of scope for v1.
- Comments inside `[project.dependencies]` in `pyproject.toml` are not preserved during rewrites.
