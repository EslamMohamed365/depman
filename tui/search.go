package tui

import (
	"fmt"
	"strings"

	"github.com/eslam/depman/config"
	"github.com/eslam/depman/pkg/pip"
	"github.com/eslam/depman/pkg/pypi"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SearchPhase tracks the current step of the search flow.
type SearchPhase int

const (
	PhaseInput   SearchPhase = iota // typing search query
	PhaseResults                    // browsing results list
	PhaseDetail                     // viewing package detail + version select
)

// SearchModel handles the search → select → version → install flow.
type SearchModel struct {
	phase         SearchPhase
	input         string
	results       []pypi.SearchResult
	resultsCursor int
	loading       bool
	err           error

	// Detail phase
	detail        *pypi.PackageDetail
	detailLoading bool
	versionCursor int
}

// SearchResultsMsg is sent when PyPI search results arrive.
type SearchResultsMsg struct {
	Results []pypi.SearchResult
	Err     error
}

// PackageDetailMsg is sent when full package detail arrives.
type PackageDetailMsg struct {
	Detail *pypi.PackageDetail
	Err    error
}

// NewSearchModel creates the search screen model.
func NewSearchModel() SearchModel {
	return SearchModel{phase: PhaseInput}
}

func (s SearchModel) Update(msg tea.Msg, state *AppState, runner *pip.Runner) (SearchModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch s.phase {
		case PhaseInput:
			return s.updateInput(msg, state)
		case PhaseResults:
			return s.updateResults(msg, state, runner)
		case PhaseDetail:
			return s.updateDetail(msg, state, runner)
		}

	case SearchResultsMsg:
		s.loading = false
		s.err = msg.Err
		if msg.Err == nil {
			s.results = msg.Results
			s.resultsCursor = 0
			if len(msg.Results) > 0 {
				s.phase = PhaseResults
			}
		}

	case PackageDetailMsg:
		s.detailLoading = false
		if msg.Err == nil && msg.Detail != nil {
			s.detail = msg.Detail
			s.versionCursor = 0
			s.phase = PhaseDetail
		} else if msg.Err != nil {
			s.err = msg.Err
		}
	}

	return s, nil
}

func (s SearchModel) updateInput(msg tea.KeyMsg, state *AppState) (SearchModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		state.Screen = ScreenDashboard
		return NewSearchModel(), nil
	case "enter":
		if len(s.input) >= MinSearchLength {
			s.loading = true
			return s, s.doSearch(state)
		}
	case "backspace":
		if len(s.input) > 0 {
			s.input = s.input[:len(s.input)-1]
		}
	default:
		if len(msg.String()) == 1 {
			s.input += msg.String()
		}
	}
	return s, nil
}

func (s SearchModel) updateResults(msg tea.KeyMsg, state *AppState, runner *pip.Runner) (SearchModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Go back to input
		s.phase = PhaseInput
		return s, nil
	case "j", "down":
		if s.resultsCursor < len(s.results)-1 {
			s.resultsCursor++
		}
	case "k", "up":
		if s.resultsCursor > 0 {
			s.resultsCursor--
		}
	case "enter":
		if len(s.results) > 0 && s.resultsCursor < len(s.results) {
			// Fetch full detail for the selected package
			name := s.results[s.resultsCursor].Name
			s.detailLoading = true
			return s, s.fetchDetail(state, name)
		}
	}
	return s, nil
}

func (s SearchModel) updateDetail(msg tea.KeyMsg, state *AppState, runner *pip.Runner) (SearchModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		// Go back to results
		s.phase = PhaseResults
		s.detail = nil
		return s, nil
	case "j", "down":
		if s.detail != nil && s.versionCursor < len(s.detail.Versions)-1 {
			s.versionCursor++
		}
	case "k", "up":
		if s.versionCursor > 0 {
			s.versionCursor--
		}
	case "enter":
		if s.detail != nil && len(s.detail.Versions) > 0 {
			pkg := s.detail.Name
			ver := s.detail.Versions[s.versionCursor]
			installStr := pkg + "==" + ver
			state.Screen = ScreenDashboard
			state.IsLoading = true
			return NewSearchModel(), func() tea.Msg {
				result := runner.Install(installStr)
				return PackageActionMsg{Action: "installed", Package: pkg + "@" + ver, Err: result.Err}
			}
		}
	}
	return s, nil
}

func (s SearchModel) doSearch(state *AppState) tea.Cmd {
	query := s.input
	mirror := state.Config.PyPI.Mirror
	return func() tea.Msg {
		client := pypi.NewClient(mirror)
		results, err := client.Search(query)
		return SearchResultsMsg{Results: results, Err: err}
	}
}

func (s SearchModel) fetchDetail(state *AppState, name string) tea.Cmd {
	mirror := state.Config.PyPI.Mirror
	return func() tea.Msg {
		client := pypi.NewClient(mirror)
		detail, err := client.GetPackageDetail(name)
		return PackageDetailMsg{Detail: detail, Err: err}
	}
}

// View renders the current search phase.
func (s SearchModel) View(state AppState) string {
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
		Padding(1, 2)

	switch s.phase {
	case PhaseInput:
		return container.Render(s.viewInput(w))
	case PhaseResults:
		return container.Render(s.viewResults(w, h))
	case PhaseDetail:
		return container.Render(s.viewDetail(w, h))
	default:
		return container.Render(s.viewInput(w))
	}
}

func (s SearchModel) viewInput(w int) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(config.ColorBlue)
	inputStyle := lipgloss.NewStyle().Foreground(config.ColorCyan)
	cursorChar := lipgloss.NewStyle().Foreground(config.ColorOrange).Render("█")
	dimStyle := lipgloss.NewStyle().Foreground(config.ColorFGDim)

	var b strings.Builder
	b.WriteString(titleStyle.Render("🔍 Search PyPI"))
	b.WriteString("\n\n")
	b.WriteString("  Package name: ")
	b.WriteString(inputStyle.Render(s.input))
	b.WriteString(cursorChar)
	b.WriteString("\n\n")

	if s.loading {
		b.WriteString(dimStyle.Render("  Searching..."))
	} else if s.err != nil {
		b.WriteString(lipgloss.NewStyle().Foreground(config.ColorRed).Render(
			fmt.Sprintf("  Error: %v", s.err)))
	} else {
		b.WriteString(dimStyle.Render("  Type a package name and press Enter to search"))
	}

	b.WriteString("\n\n")
	b.WriteString(dimStyle.Render("  Enter to search  │  Esc to cancel"))

	return b.String()
}

func (s SearchModel) viewResults(w, h int) string {
	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(config.ColorBlue)
	dimStyle := lipgloss.NewStyle().Foreground(config.ColorFGDim)

	var b strings.Builder
	b.WriteString(titleStyle.Render(fmt.Sprintf("🔍 Results for \"%s\"", s.input)))
	b.WriteString(dimStyle.Render(fmt.Sprintf("  (%d found)", len(s.results))))
	b.WriteString("\n\n")

	if s.detailLoading {
		b.WriteString(dimStyle.Render("  Loading package details..."))
		b.WriteString("\n")
	} else if len(s.results) == 0 {
		b.WriteString(dimStyle.Render("  No packages found"))
		b.WriteString("\n")
	} else {
		maxVisible := h - ViewportHeaderLines
		if maxVisible < MinVisibleResults {
			maxVisible = MinVisibleResults
		}
		visible := min(len(s.results), maxVisible)

		for i := 0; i < visible; i++ {
			r := s.results[i]
			name := lipgloss.NewStyle().Foreground(config.ColorPurple).Bold(true).Render(r.Name)
			ver := lipgloss.NewStyle().Foreground(config.ColorCyan).Render("v" + r.Version)

			descMaxLen := w - 10
		if descMaxLen < MinDescriptionLength {
				descMaxLen = MinDescriptionLength
			}
			desc := dimStyle.Render(truncate(r.Description, descMaxLen))

			if i == s.resultsCursor {
				indicator := lipgloss.NewStyle().Foreground(config.ColorBlue).Render("▶ ")
				nameVer := lipgloss.NewStyle().Background(config.ColorBGHighlight).
					Render(fmt.Sprintf("%s%s  %s", indicator, name, ver))
				b.WriteString(fmt.Sprintf("  %s\n", nameVer))
				b.WriteString(fmt.Sprintf("      %s\n", desc))
			} else {
				b.WriteString(fmt.Sprintf("    %s  %s\n", name, ver))
				b.WriteString(fmt.Sprintf("      %s\n", desc))
			}
		}
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  Enter to view details  │  j/k navigate  │  Esc to go back"))

	return b.String()
}

func (s SearchModel) viewDetail(w, h int) string {
	if s.detail == nil {
		return ""
	}
	d := s.detail

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(config.ColorPurple)
		labelStyle := lipgloss.NewStyle().Foreground(config.ColorFGDim).Width(KeyColumnWidth)
	valueStyle := lipgloss.NewStyle().Foreground(config.ColorFG)
	verStyle := lipgloss.NewStyle().Foreground(config.ColorCyan)
	dimStyle := lipgloss.NewStyle().Foreground(config.ColorFGDim)
	sectionStyle := lipgloss.NewStyle().Bold(true).Foreground(config.ColorYellow).MarginTop(1)

	var b strings.Builder

	// Package header
	b.WriteString(titleStyle.Render(d.Name))
	b.WriteString("  ")
	b.WriteString(verStyle.Render("v" + d.Version))
	b.WriteString("\n\n")

	// Info fields
	if d.Summary != "" {
		descMaxLen := w - 20
		if descMaxLen < MinDescWidth {
			descMaxLen = MinDescWidth
		}
		b.WriteString(labelStyle.Render("Description"))
		b.WriteString(valueStyle.Render(truncate(d.Summary, descMaxLen)))
		b.WriteString("\n")
	}
	if d.Author != "" {
		b.WriteString(labelStyle.Render("Author"))
		b.WriteString(valueStyle.Render(d.Author))
		b.WriteString("\n")
	}
	if d.License != "" {
		b.WriteString(labelStyle.Render("License"))
		b.WriteString(valueStyle.Render(d.License))
		b.WriteString("\n")
	}
	if d.RequiresPy != "" {
		b.WriteString(labelStyle.Render("Requires"))
		b.WriteString(valueStyle.Render("Python " + d.RequiresPy))
		b.WriteString("\n")
	}
	if d.HomePage != "" {
		b.WriteString(labelStyle.Render("Homepage"))
		b.WriteString(dimStyle.Render(d.HomePage))
		b.WriteString("\n")
	}

	// Version selection
	b.WriteString("\n")
	b.WriteString(sectionStyle.Render("Select Version"))
	b.WriteString("\n\n")

	maxVersions := h - MaxVersionsDisplay
	if maxVersions < MinVisibleResults {
		maxVersions = MinVisibleResults
	}
	visible := min(len(d.Versions), maxVersions)

	for i := 0; i < visible; i++ {
		ver := d.Versions[i]
		isLatest := (i == 0)

		verText := verStyle.Render(ver)
		if isLatest {
			verText += dimStyle.Render(" (latest)")
		}

		if i == s.versionCursor {
			indicator := lipgloss.NewStyle().Foreground(config.ColorBlue).Render("▶ ")
			line := lipgloss.NewStyle().Background(config.ColorBGHighlight).
				Render(fmt.Sprintf("%s%s", indicator, verText))
			b.WriteString(fmt.Sprintf("  %s\n", line))
		} else {
			b.WriteString(fmt.Sprintf("    %s\n", verText))
		}
	}

	if len(d.Versions) > visible {
		b.WriteString(dimStyle.Render(fmt.Sprintf("    ... and %d more versions", len(d.Versions)-visible)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(dimStyle.Render("  Enter to install  │  j/k select version  │  Esc to go back"))

	return b.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
