# Agent Guidelines for depman

This file provides guidelines for agentic coding agents working on the depman project.

## Project Overview

**depman** is a Go-based TUI application for managing Python dependencies using Bubble Tea.
- Go 1.25+ required
- Main packages: `pkg/pypi` (PyPI API client), `pkg/pip` (package operations), `tui` (UI components)
- Dependencies: charmbracelet/bubbletea, charmbracelet/lipgloss, pelletier/go-toml

---

## Build, Lint & Test Commands

### Build
```bash
go build -o depman .
```

### Run
```bash
./depman
```

### Test All
```bash
go test ./...
```

### Test Single Package
```bash
go test ./pkg/pypi/...
go test ./pkg/pip/...
go test ./config/...
```

### Test Single File
```bash
go test ./pkg/pypi/client_test.go ./pkg/pypi/client.go
```

### Test Single Function
```bash
go test -run TestIsStableVersion ./pkg/pypi/...
go test -run TestNewClient ./pkg/pypi/...
go test -run TestValidatePackageName ./pkg/pip/...
```

### Verbose Test Output
```bash
go test -v ./pkg/pypi/...
```

### Test Coverage
```bash
go test -cover ./...
```

### Lint (using go vet)
```bash
go vet ./...
```

### Format Code
```bash
go fmt ./...
```

---

## Code Style Guidelines

### Imports

**Standard library first, then external packages, then project packages.**
Group imports with blank lines between groups:

```go
import (
    "context"
    "fmt"
    "strings"
    "sync"

    "github.com/eslam/depman/pkg/log"
    "github.com/eslam/depman/pkg/pip"

    tea "github.com/charmbracelet/bubbletea"
    "github.com/charmbracelet/lipgloss"
)
```

### Naming Conventions

- **Types**: PascalCase (e.g., `Client`, `PackageManager`, `SearchResult`)
- **Functions/Methods**: PascalCase (e.g., `GetPackage`, `NewClient`, `isStableVersion`)
- **Variables/Constants**: camelCase or PascalCase for exported, camelCase for unexported
- **Constants**: PascalCase for groups, or SCREAMING_SNAKE_CASE for values (see constants below)
- **Packages**: short, lowercase, no underscores (e.g., `pkg/pypi`, `pkg/env`)

### Constants

Group related constants in blocks with comments:

```go
// HTTP Configuration
const (
    HTTPTimeout      = 10 * time.Second
    MaxRetries       = 3
    RetryDelay       = 500 // milliseconds
    BackoffMultiplier = 2.0
)

// UI Layout Constants
const (
    MinPanelWidth  = 20
    MinPanelHeight = 5
)
```

### Error Handling

- Return errors with context using `fmt.Errorf` with `%w` for wrapped errors
- Log errors before returning when appropriate
- Use sentinel errors for known error conditions

```go
// Good
if err != nil {
    return nil, fmt.Errorf("pypi: fetch package: %w", err)
}

// With logging
if err != nil {
    log.Error("failed to fetch package", "package", name, "error", err)
    return nil, fmt.Errorf("pypi: fetch package: %w", err)
}
```

### Logging

Use the structured logger from `pkg/log`:

```go
log.Info("operation started", "key", value)
log.Warn("operation failed", "error", err)
log.Debug("detail info", "debug", info)
```

### Types and Structs

- Use struct tags for JSON serialization
- Document exported types with comments
- Keep structs focused and cohesive

```go
// SearchResult represents a PyPI package from search.
type SearchResult struct {
    Name        string `json:"name"`
    Version     string `json:"version"`
    Description string `json:"summary"`
}
```

### Goroutines and Concurrency

- Always use `sync.WaitGroup` for goroutine synchronization
- Use channels for communication, avoid shared memory
- Close channels from the sending side
- Use buffered channels when the number of results is known

```go
resultChan := make(chan *SearchResult, len(pendingVariations))
var wg sync.WaitGroup

for _, v := range variations {
    wg.Add(1)
    go func(name string) {
        defer wg.Done()
        // work...
    }(v)
}

go func() {
    wg.Wait()
    close(resultChan)
}()
```

### Input Validation

Validate all external input before processing. Use the validation module:

```go
// In pkg/pip/validation.go
if err := ValidatePackageName(pkg); err != nil {
    log.Warn("package validation failed", "package", pkg, "error", err)
    return RunResult{Err: fmt.Errorf("invalid package: %w", err)}
}
```

### HTTP Clients

- Reuse HTTP clients for connection pooling
- Set timeouts on clients
- Handle retries for transient errors (429, 5xx)

```go
var defaultHTTPClient = &http.Client{
    Timeout: HTTPTimeout,
}
```

### TUI Patterns

- Use Bubble Tea's model-view-update pattern
- Handle tea.WindowSizeMsg to set dimensions
- Return tea.Cmd for async operations
- Use constants for magic numbers (see `tui/constants.go`)

```go
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    switch msg := msg.(type) {
    case tea.WindowSizeMsg:
        m.width = msg.Width
        m.height = msg.Height
    }
    return m, nil
}
```

### Testing

- Use table-driven tests for multiple cases
- Use `httptest` for HTTP testing
- Test both success and failure paths
- Keep tests in `*_test.go` files alongside the code

```go
func TestClient_GetPackage_Success(t *testing.T) {
    server := httptest.NewServer(...)
    defer server.Close()

    client := NewClient(server.URL)
    pkg, err := client.GetPackage("requests")
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    // assertions...
}
```

---

## Key Files

- `pkg/pypi/client.go` - PyPI API client with search, package details
- `pkg/pip/runner.go` - Package installation/uninstallation operations
- `pkg/pip/validation.go` - PEP 508 package name validation
- `pkg/env/manager.go` - Package manager detection (uv, pip)
- `tui/dashboard.go` - Main dashboard view with panels
- `tui/constants.go` - UI layout constants
- `pkg/log/logger.go` - Structured logging
