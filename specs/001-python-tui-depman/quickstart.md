# Quickstart Test Scenarios: depman

**Branch**: `001-python-tui-depman` | **Date**: 2026-02-23

## Prerequisites

- Go 1.22+ installed
- Python 3.8+ installed
- `pip` available in PATH
- A test directory for running scenarios

## Scenario 1: First Launch — No Project (US1)

```bash
mkdir /tmp/depman-test && cd /tmp/depman-test
depman
```

**Expected**:
1. Screen shows: "No Python project found in current directory."
2. Three options: Create pyproject.toml (recommended), Create requirements.txt, Exit
3. Select "Create pyproject.toml" → prompts for project name (default: `depman-test`), version (default: `0.1.0`)
4. File created, tool proceeds to dashboard with empty package list
5. Status bar shows: `depman-test · system · pip · 0 packages · 0 outdated`

---

## Scenario 2: Launch with Existing Project (US1 + US2)

```bash
cd /tmp/depman-test
python -m venv .venv
source .venv/bin/activate
pip install requests flask
depman
```

**Expected**:
1. Dashboard appears with:
   - Installed panel: `requests`, `flask`, and their dependencies with versions
   - Outdated panel: any packages with available updates (color-coded)
2. Status bar shows: `.venv · pip · X packages · Y outdated`
3. Navigation with `j`/`k` moves cursor in the list

---

## Scenario 3: Add a Package (US3)

From the dashboard:

1. Press `a`
2. Type `httpx`
3. See search results from PyPI with package name, version, description
4. Select `httpx` with Enter
5. Confirm installation

**Expected**:
- Package installs successfully
- `httpx` appears in installed packages list
- `pyproject.toml` or `requirements.txt` is updated with `httpx==<version>`
- Success notification shown

---

## Scenario 4: Remove a Package (US4)

From the dashboard:

1. Navigate to `httpx` with `j`/`k`
2. Press `d`
3. Confirmation prompt appears: "Remove httpx? [y/N]"
4. Press `y`

**Expected**:
- Package uninstalled
- `httpx` removed from installed list
- Dependency file updated (httpx removed)
- Success notification shown

---

## Scenario 5: Update a Package (US5)

From the dashboard (with an outdated package):

1. Press `Tab` to switch to outdated panel
2. Navigate to an outdated package
3. Press `u`
4. See version diff: `1.0.0 → 1.1.0`
5. Confirm

**Expected**:
- Package updated to latest version
- Dependency file reflects new version
- Package moves from outdated to up-to-date

---

## Scenario 6: Search PyPI (US6)

From any screen:

1. Press `/`
2. Type `django`
3. See real-time search results

**Expected**:
- Results appear within 2 seconds
- Each result shows: name, version, description
- Selecting a result and confirming installs it

---

## Scenario 7: Vim Navigation (US7)

From the dashboard with packages:

1. Press `gg` → cursor jumps to first item
2. Press `G` → cursor jumps to last item
3. Press `Ctrl+d` → half-page down
4. Press `Ctrl+u` → half-page up
5. Press `Tab` → focus switches between installed/outdated panels
6. Press `?` → help screen appears
7. Press `q` → application exits

---

## Cleanup

```bash
deactivate
rm -rf /tmp/depman-test
```
