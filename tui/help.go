package tui

import (
	"strings"

	"github.com/eslam/depman/config"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// HelpModel renders the help screen overlay.
type HelpModel struct{}

// NewHelpModel creates a new help model.
func NewHelpModel() HelpModel {
	return HelpModel{}
}

func (h HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	return h, nil
}

func (hm HelpModel) View(state AppState) string {
	tw := state.Width
	th := state.Height
	if tw == 0 {
		tw = DefaultWidth
	}
	if th == 0 {
		th = DefaultHeight
	}

	container := lipgloss.NewStyle().
		Width(tw).
		Height(th).
		Padding(1, 2)

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(config.ColorBlue).
		MarginBottom(1)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(config.ColorYellow).
		MarginTop(1)

	keyStyle := lipgloss.NewStyle().
		Foreground(config.ColorCyan).
		Width(14)

	descStyle := lipgloss.NewStyle().
		Foreground(config.ColorFG)

	dimStyle := lipgloss.NewStyle().
		Foreground(config.ColorFGDim)

	var b strings.Builder

	b.WriteString(titleStyle.Render("depman — Keyboard Reference"))
	b.WriteString("\n\n")

	b.WriteString(headerStyle.Render("Navigation"))
	b.WriteString("\n")
	nav := []struct{ key, desc string }{
		{"j / ↓", "Move down"},
		{"k / ↑", "Move up"},
		{"gg", "Jump to top"},
		{"G", "Jump to bottom"},
		{"Ctrl+d", "Half-page down"},
		{"Ctrl+u", "Half-page up"},
		{"Tab", "Switch panel"},
	}
	for _, bind := range nav {
		b.WriteString(keyStyle.Render(bind.key))
		b.WriteString(descStyle.Render(bind.desc))
		b.WriteString("\n")
	}

	b.WriteString(headerStyle.Render("Package Actions"))
	b.WriteString("\n")
	actions := []struct{ key, desc string }{
		{"a", "Add package"},
		{"d / x", "Remove selected package"},
		{"u", "Update selected package"},
		{"U", "Update all outdated"},
		{"/ or s", "Search PyPI"},
		{"Enter", "Confirm action"},
		{"Esc", "Cancel / go back"},
	}
	for _, act := range actions {
		b.WriteString(keyStyle.Render(act.key))
		b.WriteString(descStyle.Render(act.desc))
		b.WriteString("\n")
	}

	b.WriteString(headerStyle.Render("General"))
	b.WriteString("\n")
	general := []struct{ key, desc string }{
		{"?", "Toggle help"},
		{"q", "Quit"},
		{"Ctrl+c", "Force quit"},
	}
	for _, g := range general {
		b.WriteString(keyStyle.Render(g.key))
		b.WriteString(descStyle.Render(g.desc))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("Press ? or Esc to close"))

	return container.Render(b.String())
}
