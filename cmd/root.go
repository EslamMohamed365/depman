package cmd

import (
	"fmt"
	"os"

	"github.com/eslam/depman/config"
	"github.com/eslam/depman/pkg/detector"
	"github.com/eslam/depman/pkg/env"
	"github.com/eslam/depman/tui"

	tea "github.com/charmbracelet/bubbletea"
)

// Execute is the main entrypoint called from main.go.
func Execute() error {
	// Load user config
	cfg := config.Load()

	// Detect project dependency file
	project := detector.DetectProject(".")

	// Detect virtualenv
	venv := env.DetectVirtualenv(".")

	// Detect package manager
	mgr := env.DetectPackageManager(cfg.PackageManager.Preferred)

	// Build initial app state
	state := tui.NewAppState(project, venv, mgr, cfg)

	// Create and run the Bubble Tea program
	p := tea.NewProgram(tui.NewModel(state), tea.WithAltScreen())
	m, err := p.Run()
	if err != nil {
		return fmt.Errorf("failed to run TUI: %w", err)
	}

	// Check for any exit error from the model
	if model, ok := m.(tui.Model); ok {
		if model.Err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", model.Err)
			return model.Err
		}
	}

	return nil
}
