# Tasks: depman — TUI Python Dependency Manager

**Input**: Design documents from `/specs/001-python-tui-depman/`  
**Prerequisites**: spec.md (user stories), depman-prd.md (architecture & tech stack)  
**Tests**: Not explicitly requested — test tasks omitted  
**Organization**: Tasks grouped by user story priority for independent implementation

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1–US8)
- Exact file paths included in descriptions

## Path Conventions

Based on PRD project structure:

```
depman/
├── main.go
├── cmd/              # CLI entrypoint
├── tui/              # Bubble Tea models and views
│   ├── dashboard.go
│   ├── search.go
│   └── components/
├── pkg/
│   ├── detector/     # Project & dep file detection
│   ├── parser/       # requirements.txt + pyproject.toml parsing
│   ├── pypi/         # PyPI API client
│   ├── pip/          # pip/uv subprocess wrapper
│   └── env/          # Virtualenv & package manager detection
└── config/           # User config
```

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Go project initialization and basic structure

- [x] T001 Initialize Go module with `go mod init` and add Bubble Tea + Lip Gloss dependencies in `go.mod`
- [x] T002 Create project directory structure: `cmd/`, `tui/`, `tui/components/`, `pkg/detector/`, `pkg/parser/`, `pkg/pypi/`, `pkg/pip/`, `pkg/env/`, `config/`
- [x] T003 [P] Create `main.go` with CLI entrypoint that bootstraps the Bubble Tea program
- [x] T004 [P] Create `config/theme.go` with Tokyo Night color palette constants (all Lip Gloss color definitions from PRD §7.1)
- [x] T005 [P] Create `config/config.go` with config struct and TOML loader for `~/.config/depman/config.toml`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T006 Create `pkg/detector/project.go` — project dependency file detection logic: scan CWD for `pyproject.toml` → `requirements.txt` → `requirements/*.txt` in priority order
- [x] T007 Create `pkg/env/virtualenv.go` — virtualenv detection logic: check `$VIRTUAL_ENV` → `.venv/bin/python` → `venv/bin/python` → system Python, in priority order
- [x] T008 Create `pkg/env/manager.go` — package manager detection logic: check for `uv` in `$PATH` → fallback to `pip`/`pip3`, with config override support
- [x] T009 [P] Create `pkg/pip/runner.go` — subprocess runner for `pip`/`uv` commands: install, uninstall, upgrade, list, list-outdated with command mapping for both tools
- [x] T010 [P] Create `pkg/parser/requirements.go` — `requirements.txt` parser: read package==version lines, skip comments and blank lines
- [x] T011 [P] Create `pkg/parser/pyproject.go` — `pyproject.toml` parser: read `[project.dependencies]` array, preserve all other sections in memory
- [x] T012 Create `pkg/parser/writer.go` — dependency file writer: full rewrite strategy with atomic writes (temp file + rename), format rules per file type, `depman` header comment
- [x] T013 Create `pkg/pip/packages.go` — package list resolver: parse output of `pip list` / `uv pip list` into structured package data (name, version)
- [x] T014 Create `pkg/pip/outdated.go` — outdated package resolver: parse output of `pip list --outdated` / `uv pip list --outdated`, compute semver diff type (patch/minor/major)
- [x] T015 Create `tui/model.go` — root Bubble Tea model with app state, message types, and `Init`/`Update`/`View` scaffolding

**Checkpoint**: Foundation ready — user story implementation can now begin

---

## Phase 3: User Story 1 — Launch and Environment Setup (Priority: P1) 🎯 MVP

**Goal**: User launches `depman`, project is detected (or initialized), virtualenv is resolved (or created), package manager is detected, and the status bar displays all environment info.

**Independent Test**: Launch the tool in directories with/without Python projects and verify correct detection, prompts, and status bar display.

### Implementation for User Story 1

- [x] T016 [US1] Create `tui/components/statusbar.go` — status bar component showing project path, file type, venv name, package manager indicator (⚡ uv / pip), package counts, outdated count, help hint; styled with Tokyo Night palette
- [x] T017 [US1] Create `tui/init.go` — initialization screen model: "No Python project found" prompt with options to create `pyproject.toml` (recommended), `requirements.txt` (simple), or exit
- [x] T018 [US1] Create `pkg/detector/init_project.go` — `pyproject.toml` creation flow: prompt for project name (default CWD name), version (default 0.1.0), Python version; generate minimal valid file. Also handle `requirements.txt` creation (empty file with header)
- [x] T019 [US1] Create `tui/venv_prompt.go` — virtualenv creation prompt: "[C] Create .venv  [I] Ignore" status bar prompt; on C: run `uv venv .venv` or `python -m venv .venv`; set `$VIRTUAL_ENV` for subprocesses; update status bar
- [x] T020 [US1] Create `pkg/env/create.go` — virtualenv creation logic: prefer `uv venv .venv`, fallback to `python -m venv .venv`; detect and warn on broken virtualenvs (missing interpreter) with recreate/fallback options
- [x] T021 [US1] Implement startup orchestration in `cmd/root.go` — wire detection → init prompt (if needed) → venv resolution → venv prompt (if needed) → package manager detection → launch dashboard
- [x] T022 [US1] Create `tui/components/error.go` — inline error display component for showing errors like "No pip or uv found" with install hints

**Checkpoint**: User Story 1 fully functional — tool launches, detects or creates project, resolves environment, shows status bar

---

## Phase 4: User Story 2 — View Installed and Outdated Packages (Priority: P1) 🎯 MVP

**Goal**: Dashboard displays installed packages panel and outdated packages panel with color-coded semver diff, plus status bar with counts.

**Independent Test**: Launch in a project with known installed/outdated packages, verify both panels render correctly with proper color coding.

### Implementation for User Story 2

- [x] T023 [US2] Create `tui/dashboard.go` — main dashboard Bubble Tea model with two-panel layout: installed packages (left) and outdated packages (right), with focused panel tracking and Tab switching
- [x] T024 [P] [US2] Create `tui/components/package_list.go` — scrollable package list component: renders package name (purple) + version (cyan), supports selection highlighting (background #2e3250, left indicator ▶ in blue)
- [x] T025 [P] [US2] Create `tui/components/outdated_list.go` — outdated packages list component: columns for package name, current version, latest version, diff type badge (patch=teal, minor=yellow, major=red)
- [x] T026 [US2] Create `tui/components/panel.go` — bordered panel wrapper with title, border color (#414868 default, #7aa2f7 focused), background (#1a1b26), and content area
- [x] T027 [US2] Create `pkg/pypi/version.go` — semver diff calculator: compare current vs latest version, classify as patch/minor/major update
- [x] T028 [US2] Integrate package data loading into dashboard: on startup, run `pip list` and `pip list --outdated`, populate both panels, update status bar counts

**Checkpoint**: User Story 2 fully functional — dashboard shows installed and outdated packages with color coding

---

## Phase 5: User Story 3 — Add a Package (Priority: P2)

**Goal**: User presses `a`, types a package name, sees real-time PyPI search results, selects one, confirms, and the package is installed + dependency file rewritten.

**Independent Test**: Press `a`, search for "requests", select it, confirm, verify package installed and dependency file updated.

### Implementation for User Story 3

- [x] T029 [P] [US3] Create `pkg/pypi/client.go` — PyPI JSON API client: fetch package info from `https://pypi.org/pypi/<pkg>/json`, search packages, return name, version, description, download count
- [x] T030 [P] [US3] Create `pkg/pypi/search.go` — PyPI search with debounced requests: accept query string, return list of matching packages with metadata
- [x] T031 [US3] Create `tui/components/search_input.go` — inline search input component: text input with PyPI blue border, orange cursor, dim placeholder text; debounced onChange handler
- [x] T032 [US3] Create `tui/components/search_results.go` — search results popup list: shows package name, latest version, short description; selectable with Enter to confirm
- [x] T033 [US3] Create `tui/add_package.go` — add package flow model: open search input on `a`, show results, on confirm: run install command, trigger dependency file rewrite, show success/error notification, return to dashboard

**Checkpoint**: User Story 3 fully functional — users can add packages via inline search

---

## Phase 6: User Story 4 — Remove a Package (Priority: P2)

**Goal**: User navigates to a package, presses `d`/`x`, confirms, package is uninstalled and removed from dependency file.

**Independent Test**: Navigate to a known package, press `d`, confirm, verify uninstalled and removed from file.

### Implementation for User Story 4

- [x] T034 [P] [US4] Create `tui/components/confirm.go` — confirmation prompt component: "Remove <package>? [y/N]" inline prompt with accept/cancel keybindings
- [x] T035 [US4] Create `tui/remove_package.go` — remove package flow model: on `d`/`x` show confirm prompt, on confirm: run uninstall command, trigger file rewrite, show notification, refresh dashboard

**Checkpoint**: User Story 4 fully functional — users can remove packages with confirmation

---

## Phase 7: User Story 5 — Update Packages (Priority: P2)

**Goal**: User presses `u` to update one package or `U` to update all outdated, with version diff preview, confirmation, and dependency file sync.

**Independent Test**: With known outdated packages, press `u` on one, verify it updates. Press `U`, verify all update sequentially.

### Implementation for User Story 5

- [x] T036 [US5] Create `tui/update_package.go` — single update flow: on `u` show "current → latest" diff, on confirm: run upgrade, trigger file rewrite, refresh outdated panel
- [x] T037 [US5] Create `tui/update_all.go` — update all flow: on `U` show summary "X packages will be updated", on confirm: run upgrades sequentially with progress indicator, trigger file rewrite after all complete
- [x] T038 [P] [US5] Create `tui/components/progress.go` — progress indicator component for sequential update operations: shows current/total, package name being updated, success/fail status per package

**Checkpoint**: User Story 5 fully functional — single and bulk updates with progress tracking

---

## Phase 8: User Story 6 — Search PyPI Online (Priority: P3)

**Goal**: User presses `/` or `s` from anywhere, types a query, sees real-time PyPI results, can install directly.

**Independent Test**: Press `/`, type a query, verify results appear from PyPI, select and install one.

### Implementation for User Story 6

- [x] T039 [US6] Create `tui/search.go` — full search panel model: opens on `/` or `s` from any screen, text input with real-time debounced PyPI queries, results list with install action, offline indicator when PyPI unreachable
- [x] T040 [US6] Create `pkg/pypi/cache.go` — simple response cache for PyPI searches: store recent results in memory, serve cached data when offline, show offline indicator

**Checkpoint**: User Story 6 fully functional — global PyPI search with install and offline fallback

---

## Phase 9: User Story 7 — Navigate with Vim Motions (Priority: P3)

**Goal**: Full vim-style navigation: `j`/`k`, `gg`/`G`, `Ctrl+d`/`Ctrl+u`, `Tab`, `?` help, `q` quit.

**Independent Test**: Press each vim keybinding and verify expected navigation behavior.

### Implementation for User Story 7

- [x] T041 [US7] Create `tui/keymap.go` — centralized keymap handler: map all vim keybindings (`j`, `k`, `gg`, `G`, `Ctrl+d`, `Ctrl+u`, `Tab`, `q`, `?`, `:`) to TUI actions; handle multi-key sequences (`gg`)
- [x] T042 [US7] Create `tui/help.go` — help screen model: full-screen overlay showing all keybindings in a formatted table, dismiss with `?` or `Esc`
- [x] T043 [US7] Integrate vim keybindings into all existing list components (`package_list.go`, `outdated_list.go`, `search_results.go`): half-page scroll, jump-to-top/bottom, panel switching with Tab

**Checkpoint**: User Story 7 fully functional — full vim navigation across all screens

---

## Phase 10: User Story 8 — Dependency File Sync (Priority: P2)

**Goal**: Every package operation triggers a full rewrite of the dependency file using atomic writes. `pyproject.toml` non-dependency sections are preserved.

**Independent Test**: Perform add/remove/update operations, verify file content matches expected format and preserves pyproject.toml sections.

### Implementation for User Story 8

- [x] T044 [US8] Create `pkg/parser/sync.go` — sync orchestrator: after any successful pip/uv command, query full installed list, reconstruct dependency file from scratch, write atomically
- [x] T045 [US8] Integrate sync into all package action flows (`add_package.go`, `remove_package.go`, `update_package.go`, `update_all.go`): call sync after successful operations, skip on failure, show error if file write fails
- [x] T046 [US8] Add `pyproject.toml` section preservation in `pkg/parser/pyproject.go`: on rewrite, only replace `[project.dependencies]` array, preserve all other TOML content exactly

**Checkpoint**: User Story 8 fully functional — dependency file always reflects true installed state

---

## Phase 11: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [x] T047 [P] Create `tui/components/notification.go` — toast notification component (bottom-right): success (green on elevated bg), error (red on elevated bg), auto-dismiss after 3 seconds
- [x] T048 [P] Add comprehensive error handling across all flows per PRD §10: no pip/uv found → install hints, broken venv → warn + recreate, parse error → show with line number, uv not found when configured → warn + auto-fallback
- [x] T049 [P] Create `config/defaults.go` — default configuration with all configurable values: package manager preference, theme, PyPI mirror URL, virtualenv search paths
- [x] T050 Add `requirements/*.txt` support in `pkg/detector/project.go` — detect subdirectory requirement files (e.g., `requirements/base.txt`)
- [x] T051 [P] Create `README.md` with project description, installation instructions, usage guide, keybinding reference, configuration options, and screenshots
- [x] T052 Final integration pass: verify all keybindings work from all screens, status bar updates on every state change, no blocking UI during network calls (async PyPI fetches)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion — BLOCKS all user stories
- **User Stories (Phases 3–10)**: All depend on Foundational phase completion
  - US1 + US2 (both P1): Priority, implement first
  - US3, US4, US5, US8 (P2): Implement after P1 stories
  - US6, US7 (P3): Implement last
- **Polish (Phase 11)**: Depends on all user stories being complete

### User Story Dependencies

- **US1 (Launch & Environment)**: Can start after Phase 2. No other story dependencies. **MVP baseline.**
- **US2 (View Packages)**: Can start after Phase 2. Requires US1's status bar component. **Complete MVP.**
- **US3 (Add Package)**: Can start after Phase 2. Depends on US8 (file sync) for dependency file writes. Shares PyPI client with US6.
- **US4 (Remove Package)**: Can start after Phase 2. Depends on US8 for file sync. Shares confirm component with US5.
- **US5 (Update Packages)**: Can start after Phase 2. Depends on US2 (outdated list) and US8 (file sync).
- **US6 (Search PyPI)**: Can start after Phase 2. Reuses components from US3 (search input/results). Shares PyPI client.
- **US7 (Vim Motions)**: Can start after Phase 2. Integrates into all list components from other stories.
- **US8 (File Sync)**: Can start after Phase 2. **Foundational for US3–US6** — implement early despite P2 label.

### Recommended Execution Order

```
Phase 1 → Phase 2 → US1 → US2 → US8 → US3 → US4 → US5 → US6 → US7 → Polish
```

> **Note**: US8 (File Sync) is promoted ahead of other P2 stories because US3–US5 depend on it.

### Within Each User Story

- Models/types before services
- Services before TUI components
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

**Phase 1** (all [P] tasks):
```
T003 (main.go) | T004 (theme.go) | T005 (config.go)
```

**Phase 2** (parallel pairs):
```
T010 (requirements parser) | T011 (pyproject parser) | T009 (pip runner)
```

**Phase 4 — US2** (parallel components):
```
T024 (package_list.go) | T025 (outdated_list.go)
```

**Phase 5 — US3** (parallel core):
```
T029 (pypi client) | T030 (pypi search)
```

**Phase 11 — Polish** (all [P] tasks):
```
T047 (notifications) | T048 (error handling) | T049 (defaults) | T051 (README)
```

---

## Implementation Strategy

### MVP First (User Stories 1 + 2 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational
3. Complete Phase 3: US1 — Launch & Environment Setup
4. Complete Phase 4: US2 — View Installed & Outdated Packages
5. **STOP and VALIDATE**: Tool launches, detects environment, shows packages with outdated highlighting
6. Deploy/demo if ready — this is a read-only MVP that already delivers value

### Incremental Delivery

1. Setup + Foundational → Foundation ready
2. US1 + US2 → **Read-only MVP** (launch, detect, view)
3. US8 → File sync infrastructure (required for mutations)
4. US3 → Add packages (first mutation)
5. US4 → Remove packages
6. US5 → Update packages (single + bulk)
7. US6 → Online PyPI search
8. US7 → Full vim navigation
9. Polish → Error handling, notifications, README, integration pass

### Solo Developer Strategy (Recommended)

Since this is a personal project:

1. Work sequentially through phases in recommended order
2. Use [P] markers to identify parallelizable tasks within each phase
3. Commit after each completed user story — each is a deployable increment
4. Test each story independently before starting the next

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story is independently completable and testable
- US8 (File Sync) is cross-cutting but implemented as its own phase since US3–US5 depend on it
- Commit after each task or logical group
- Stop at any checkpoint to validate the story independently
- No test tasks generated (tests not explicitly requested in spec)

---

## Summary

| Metric | Value |
|--------|-------|
| **Total tasks** | 52 |
| **Phase 1 (Setup)** | 5 tasks |
| **Phase 2 (Foundational)** | 10 tasks |
| **US1 (Launch)** | 7 tasks |
| **US2 (View Packages)** | 6 tasks |
| **US3 (Add Package)** | 5 tasks |
| **US4 (Remove Package)** | 2 tasks |
| **US5 (Update Packages)** | 3 tasks |
| **US6 (Search PyPI)** | 2 tasks |
| **US7 (Vim Motions)** | 3 tasks |
| **US8 (File Sync)** | 3 tasks |
| **Polish** | 6 tasks |
| **Parallel opportunities** | 5 groups identified |
| **MVP scope** | US1 + US2 (18 tasks through Phase 4) |
