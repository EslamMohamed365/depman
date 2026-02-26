package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eslam/depman/config"
	"github.com/eslam/depman/pkg/detector"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// InitModel handles the project initialization screen.
type InitModel struct {
	cursor      int
	creating    bool
	createType  detector.FileType
	projectName string
	version     string
	step        int // 0=name, 1=version
	input       string
}

// NewInitModel creates the init screen model.
func NewInitModel(state AppState) InitModel {
	cwd, _ := os.Getwd()
	return InitModel{
		projectName: filepath.Base(cwd),
		version:     "0.1.0",
	}
}

// ProjectCreatedMsg signals that a project file was created.
type ProjectCreatedMsg struct {
	Project detector.Project
	Err     error
}

func (m InitModel) Update(msg tea.Msg, state *AppState) (InitModel, tea.Cmd) {
	if m.creating {
		return m.updateCreating(msg, state)
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j", "down":
			m.cursor++
			if m.cursor > InitMaxCursorPosition {
				m.cursor = InitMaxCursorPosition
			}
		case "k", "up":
			m.cursor--
			if m.cursor < 0 {
				m.cursor = 0
			}
		case "enter":
			switch m.cursor {
			case 0:
				m.creating = true
				m.createType = detector.FilePyprojectTOML
				m.step = 0
				m.input = m.projectName
			case 1:
				return m, m.createRequirementsTxt(state)
			case 2:
				return m, tea.Quit
			}
		case "q":
			return m, tea.Quit
		}

	case ProjectCreatedMsg:
		if msg.Err != nil {
			state.StatusMsg = "Failed to create project: " + msg.Err.Error()
		} else {
			state.Project = msg.Project
			state.Screen = ScreenDashboard
		}
		return m, nil
	}

	return m, nil
}

func (m InitModel) updateCreating(msg tea.Msg, state *AppState) (InitModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "esc":
			m.creating = false
			return m, nil
		case "enter":
			switch m.step {
			case 0:
				if m.input != "" {
					m.projectName = m.input
				}
				m.step = 1
				m.input = m.version
			case 1:
				if m.input != "" {
					m.version = m.input
				}
				m.creating = false
				return m, m.createPyprojectToml(state)
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
		}
	}
	return m, nil
}

func (m InitModel) View(state AppState) string {
	w := state.Width
	h := state.Height
	if w == 0 {
			w = DefaultWidth
	}
	if h == 0 {
			h = DefaultHeight
	}

	container := lipgloss.NewStyle().
		Width(w).
		Height(h).
		Padding(2, 4)

	if m.creating {
		return container.Render(m.viewCreating())
	}

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(config.ColorBlue).
		MarginBottom(1)
	dimStyle := lipgloss.NewStyle().Foreground(config.ColorFGDim)
	normalStyle := lipgloss.NewStyle().Foreground(config.ColorFG)

	var b strings.Builder
	b.WriteString(titleStyle.Render("No Python project found in current directory."))
	b.WriteString("\n")
	b.WriteString(normalStyle.Render("Would you like to initialize one?"))
	b.WriteString("\n\n")

	options := []string{
		"Create pyproject.toml        (recommended)",
		"Create requirements.txt      (simple)",
		"Exit",
	}

	for i, opt := range options {
		if i == m.cursor {
			indicator := lipgloss.NewStyle().Foreground(config.ColorBlue).Render("▶ ")
			text := lipgloss.NewStyle().
				Foreground(config.ColorFG).
				Background(config.ColorBGHighlight).
				Render(opt)
			b.WriteString(fmt.Sprintf("  %s%s\n", indicator, text))
		} else {
			b.WriteString(fmt.Sprintf("    %s\n", dimStyle.Render(opt)))
		}
	}

	return container.Render(b.String())
}

func (m InitModel) viewCreating() string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(config.ColorBlue)
	labelStyle := lipgloss.NewStyle().Foreground(config.ColorFG)
	inputStyle := lipgloss.NewStyle().Foreground(config.ColorCyan)
	cursorChar := lipgloss.NewStyle().Foreground(config.ColorOrange).Render("█")
	dimStyle := lipgloss.NewStyle().Foreground(config.ColorFGDim)

	var b strings.Builder
	b.WriteString(titleStyle.Render("Create pyproject.toml"))
	b.WriteString("\n\n")

	switch m.step {
	case 0:
		b.WriteString(labelStyle.Render("Project name: "))
		b.WriteString(inputStyle.Render(m.input))
		b.WriteString(cursorChar)
		b.WriteString("\n")
		b.WriteString(dimStyle.Render(fmt.Sprintf("  (default: %s)", m.projectName)))
	case 1:
		b.WriteString(dimStyle.Render(fmt.Sprintf("Project name: %s", m.projectName)))
		b.WriteString("\n")
		b.WriteString(labelStyle.Render("Version: "))
		b.WriteString(inputStyle.Render(m.input))
		b.WriteString(cursorChar)
		b.WriteString("\n")
		b.WriteString(dimStyle.Render("  (default: 0.1.0)"))
	}

	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("Press Enter to confirm, Esc to cancel"))
	return b.String()
}

func (m InitModel) createPyprojectToml(state *AppState) tea.Cmd {
	name := m.projectName
	version := m.version
	dir := state.Project.Dir

	return func() tea.Msg {
		content := fmt.Sprintf(`[project]
name = "%s"
version = "%s"
requires-python = ">=3.8"
dependencies = [
    # Generated by depman
]
`, name, version)

		path := filepath.Join(dir, "pyproject.toml")
		if err := os.WriteFile(path, []byte(content), DefaultFilePermissions); err != nil {
			return ProjectCreatedMsg{Err: err}
		}

		return ProjectCreatedMsg{
			Project: detector.Project{
				FilePath: path,
				FileType: detector.FilePyprojectTOML,
				Dir:      dir,
			},
		}
	}
}

func (m InitModel) createRequirementsTxt(state *AppState) tea.Cmd {
	dir := state.Project.Dir
	return func() tea.Msg {
		content := "# Generated by depman — do not edit manually\n"
		path := filepath.Join(dir, "requirements.txt")
		if err := os.WriteFile(path, []byte(content), DefaultFilePermissions); err != nil {
			return ProjectCreatedMsg{Err: err}
		}
		return ProjectCreatedMsg{
			Project: detector.Project{
				FilePath: path,
				FileType: detector.FileRequirementsTXT,
				Dir:      dir,
			},
		}
	}
}
