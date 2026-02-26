package tui

import (
	"fmt"
	"strings"

	"github.com/eslam/depman/config"
	"github.com/eslam/depman/pkg/pip"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// DashboardModel is the main dashboard view with installed and outdated panels.
type DashboardModel struct {
	installedCursor int
	outdatedCursor  int
	installedScroll int // viewport scroll offset
	outdatedScroll  int
	width           int
	height          int
	showConfirm     bool
	confirmAction   string
	confirmPkg      string
	addMode         bool
	addInput        string
	waitingForG     bool
}

// NewDashboardModel creates a new dashboard model.
func NewDashboardModel(state AppState) DashboardModel {
	return DashboardModel{}
}

// UpdatePackages refreshes the dashboard after package data loads.
func (d *DashboardModel) UpdatePackages(installed, outdated []pip.Package) {
	if d.installedCursor >= len(installed) {
		d.installedCursor = max(0, len(installed)-1)
	}
	if d.outdatedCursor >= len(outdated) {
		d.outdatedCursor = max(0, len(outdated)-1)
	}
}

// SetSize updates the terminal dimensions.
func (d *DashboardModel) SetSize(w, h int) {
	d.width = w
	d.height = h
}

// viewableHeight returns the number of package lines visible in a panel.
func (d DashboardModel) viewableHeight() int {
	// panel height minus border(2) minus title(1) minus blank(1) minus padding
	h := d.height - DashboardReservedLines - DashboardPaddingLines
	if h < MinPanelHeight {
		h = MinPanelHeight
	}
	return h
}

// ensureVisible adjusts scroll so the cursor is always visible.
func ensureVisible(cursor, scroll, viewH int) int {
	if cursor < scroll {
		return cursor
	}
	if cursor >= scroll+viewH {
		return cursor - viewH + 1
	}
	return scroll
}

// Update handles keypresses for the dashboard.
func (d DashboardModel) Update(msg tea.Msg, state *AppState, runner *pip.Runner) (DashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if d.showConfirm {
			return d.handleConfirm(msg, state, runner)
		}
		if d.addMode {
			return d.handleAddMode(msg, state, runner)
		}

		key := msg.String()

		if d.waitingForG {
			d.waitingForG = false
			if key == "g" {
				if state.ActivePanel == PanelInstalled {
					d.installedCursor = 0
				} else {
					d.outdatedCursor = 0
				}
				d.syncScroll(state)
				return d, nil
			}
		}

		switch key {
		case "j", "down":
			if state.ActivePanel == PanelInstalled {
				if d.installedCursor < len(state.Installed)-1 {
					d.installedCursor++
				}
			} else {
				if d.outdatedCursor < len(state.Outdated)-1 {
					d.outdatedCursor++
				}
			}
			d.syncScroll(state)
		case "k", "up":
			if state.ActivePanel == PanelInstalled {
				if d.installedCursor > 0 {
					d.installedCursor--
				}
			} else {
				if d.outdatedCursor > 0 {
					d.outdatedCursor--
				}
			}
			d.syncScroll(state)
		case "g":
			d.waitingForG = true
		case "G":
			if state.ActivePanel == PanelInstalled {
				d.installedCursor = max(0, len(state.Installed)-1)
			} else {
				d.outdatedCursor = max(0, len(state.Outdated)-1)
			}
			d.syncScroll(state)
		case "ctrl+d":
	half := d.viewableHeight() / HalfPageDivisor
			if state.ActivePanel == PanelInstalled {
				d.installedCursor = min(d.installedCursor+half, max(0, len(state.Installed)-1))
			} else {
				d.outdatedCursor = min(d.outdatedCursor+half, max(0, len(state.Outdated)-1))
			}
			d.syncScroll(state)
		case "ctrl+u":
	half := d.viewableHeight() / HalfPageDivisor
			if state.ActivePanel == PanelInstalled {
				d.installedCursor = max(d.installedCursor-half, 0)
			} else {
				d.outdatedCursor = max(d.outdatedCursor-half, 0)
			}
			d.syncScroll(state)
		case "tab":
			if state.ActivePanel == PanelInstalled {
				state.ActivePanel = PanelOutdated
			} else {
				state.ActivePanel = PanelInstalled
			}

		// Actions
		case "a", "/", "s":
			state.Screen = ScreenSearch
		case "enter":
			// On installed package: open search pre-filled for version change
			pkg := d.selectedPackage(state)
			if pkg != nil {
				state.Screen = ScreenSearch
				state.VersionChangePkg = pkg.Name
			}
		case "d", "x":
			pkg := d.selectedPackage(state)
			if pkg != nil {
				d.showConfirm = true
				d.confirmAction = "remove"
				d.confirmPkg = pkg.Name
			}
		case "u":
			if state.ActivePanel == PanelOutdated {
				pkg := d.selectedOutdated(state)
				if pkg != nil {
					d.showConfirm = true
					d.confirmAction = "update"
					d.confirmPkg = pkg.Name
				}
			}
		case "U":
			if len(state.Outdated) > 0 {
				d.showConfirm = true
				d.confirmAction = "update-all"
				d.confirmPkg = fmt.Sprintf("%d packages", len(state.Outdated))
			}
		}
	}
	return d, nil
}

func (d *DashboardModel) syncScroll(state *AppState) {
	vh := d.viewableHeight()
	if state.ActivePanel == PanelInstalled {
		d.installedScroll = ensureVisible(d.installedCursor, d.installedScroll, vh)
	} else {
		d.outdatedScroll = ensureVisible(d.outdatedCursor, d.outdatedScroll, vh)
	}
}

func (d DashboardModel) handleConfirm(msg tea.KeyMsg, state *AppState, runner *pip.Runner) (DashboardModel, tea.Cmd) {
	switch msg.String() {
	case "y", "enter":
		d.showConfirm = false
		state.IsLoading = true
		action := d.confirmAction
		pkg := d.confirmPkg
		outdated := make([]pip.Package, len(state.Outdated))
		copy(outdated, state.Outdated)
		switch action {
		case "remove":
			return d, func() tea.Msg {
				result := runner.Uninstall(pkg)
				return PackageActionMsg{Action: "uninstalled", Package: pkg, Err: result.Err}
			}
		case "update":
			return d, func() tea.Msg {
				result := runner.Upgrade(pkg)
				return PackageActionMsg{Action: "updated", Package: pkg, Err: result.Err}
			}
		case "update-all":
			return d, func() tea.Msg {
				var failed []string
				succeeded := 0
				for _, p := range outdated {
					result := runner.Upgrade(p.Name)
					if result.Err != nil {
						failed = append(failed, p.Name)
					} else {
						succeeded++
					}
				}
				if len(failed) > 0 {
					err := fmt.Errorf("packages: update failed: %s", strings.Join(failed, ", "))
					msg := fmt.Sprintf("updated %d, failed %d", succeeded, len(failed))
					return PackageActionMsg{Action: msg, Package: "", Err: err}
				}
				return PackageActionMsg{Action: fmt.Sprintf("updated %d", succeeded), Package: ""}
			}
		}
	case "n", "esc", "q":
		d.showConfirm = false
	}
	return d, nil
}

func (d DashboardModel) handleAddMode(msg tea.KeyMsg, state *AppState, runner *pip.Runner) (DashboardModel, tea.Cmd) {
	switch msg.String() {
	case "esc":
		d.addMode = false
		d.addInput = ""
	case "enter":
		if d.addInput != "" {
			d.addMode = false
			pkg := d.addInput
			d.addInput = ""
			state.IsLoading = true
			return d, func() tea.Msg {
				result := runner.Install(pkg)
				return PackageActionMsg{Action: "installed", Package: pkg, Err: result.Err}
			}
		}
	case "backspace":
		if len(d.addInput) > 0 {
			d.addInput = d.addInput[:len(d.addInput)-1]
		}
	default:
		if len(msg.String()) == 1 {
			d.addInput += msg.String()
		}
	}
	return d, nil
}

// View renders the dashboard.
func (d DashboardModel) View(state AppState) string {
	w := d.width
	h := d.height
	if w == 0 || h == 0 {
		// Fallback before first WindowSizeMsg
		return "Initializing..."
	}

	if state.IsLoading {
		return d.renderLoading(w, h)
	}

	// Reserve lines: 1 status bar + 1 overlay (optional)
	overlayLines := 0
	if d.showConfirm || d.addMode {
		overlayLines = 1
	}

	panelHeight := h - StatusBarLines - overlayLines - PanelBorderLines
	if panelHeight < MinPanelHeight {
		panelHeight = MinPanelHeight
	}

	panelWidth := w/2 - 1
	if panelWidth < MinPanelWidth {
		panelWidth = w - 2
	}

	installedPanel := d.renderInstalledPanel(state, panelWidth, panelHeight)
	outdatedPanel := d.renderOutdatedPanel(state, panelWidth, panelHeight)

	body := lipgloss.JoinHorizontal(lipgloss.Top, installedPanel, outdatedPanel)

	statusBar := d.renderStatusBar(state, w)

	var overlay string
	if d.showConfirm {
		overlay = d.renderConfirmDialog(w)
	} else if d.addMode {
		overlay = d.renderAddInput(w)
	}

	if overlay != "" {
		return lipgloss.JoinVertical(lipgloss.Left, body, overlay, statusBar)
	}

	return lipgloss.JoinVertical(lipgloss.Left, body, statusBar)
}

func (d DashboardModel) renderInstalledPanel(state AppState, width, height int) string {
	focused := state.ActivePanel == PanelInstalled
	borderColor := config.ColorBorder
	if focused {
		borderColor = config.ColorBlue
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Height(height)

	titleStr := fmt.Sprintf("Installed (%d)", len(state.Installed))
	title := lipgloss.NewStyle().Bold(true).Foreground(config.ColorFG).Render(titleStr)

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	// Apply viewport scrolling
	viewH := d.viewableHeight()
	scrollStart := d.installedScroll
	scrollEnd := scrollStart + viewH
	if scrollEnd > len(state.Installed) {
		scrollEnd = len(state.Installed)
	}

	// Scroll indicator at top
	if scrollStart > 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(config.ColorFGDim).Render("  ↑ more"))
	}

	for i := scrollStart; i < scrollEnd; i++ {
		p := state.Installed[i]
		name := lipgloss.NewStyle().Foreground(config.ColorPurple).Render(p.Name)
		ver := lipgloss.NewStyle().Foreground(config.ColorCyan).Render(p.InstalledVersion)

		if focused && i == d.installedCursor {
			indicator := lipgloss.NewStyle().Foreground(config.ColorBlue).Render("▶ ")
			line := lipgloss.NewStyle().Background(config.ColorBGHighlight).Foreground(config.ColorFG).
				Render(fmt.Sprintf("%s%s %s", indicator, name, ver))
			lines = append(lines, line)
		} else {
			lines = append(lines, fmt.Sprintf("  %s %s", name, ver))
		}
	}

	// Scroll indicator at bottom
	if scrollEnd < len(state.Installed) {
		lines = append(lines, lipgloss.NewStyle().Foreground(config.ColorFGDim).Render("  ↓ more"))
	}

	if len(state.Installed) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(config.ColorFGDim).Render("  No packages installed"))
	}

	return style.Render(strings.Join(lines, "\n"))
}

func (d DashboardModel) renderOutdatedPanel(state AppState, width, height int) string {
	focused := state.ActivePanel == PanelOutdated
	borderColor := config.ColorBorder
	if focused {
		borderColor = config.ColorBlue
	}

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Width(width).
		Height(height)

	titleStr := fmt.Sprintf("Outdated (%d)", len(state.Outdated))
	title := lipgloss.NewStyle().Bold(true).Foreground(config.ColorFG).Render(titleStr)

	var lines []string
	lines = append(lines, title)
	lines = append(lines, "")

	viewH := d.viewableHeight()
	scrollStart := d.outdatedScroll
	scrollEnd := scrollStart + viewH
	if scrollEnd > len(state.Outdated) {
		scrollEnd = len(state.Outdated)
	}

	if scrollStart > 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(config.ColorFGDim).Render("  ↑ more"))
	}

	for i := scrollStart; i < scrollEnd; i++ {
		p := state.Outdated[i]
		name := lipgloss.NewStyle().Foreground(config.ColorPurple).Render(p.Name)
		cur := lipgloss.NewStyle().Foreground(config.ColorCyan).Render(p.InstalledVersion)
		diffColor := config.DiffColor(p.DiffType)
		lat := lipgloss.NewStyle().Foreground(diffColor).Render(p.LatestVersion)
		badge := lipgloss.NewStyle().Foreground(diffColor).Render(config.DiffLabel(p.DiffType))

		if focused && i == d.outdatedCursor {
			indicator := lipgloss.NewStyle().Foreground(config.ColorBlue).Render("▶ ")
			line := lipgloss.NewStyle().Background(config.ColorBGHighlight).Foreground(config.ColorFG).
				Render(fmt.Sprintf("%s%s %s → %s (%s)", indicator, name, cur, lat, badge))
			lines = append(lines, line)
		} else {
			lines = append(lines, fmt.Sprintf("  %s %s → %s (%s)", name, cur, lat, badge))
		}
	}

	if scrollEnd < len(state.Outdated) {
		lines = append(lines, lipgloss.NewStyle().Foreground(config.ColorFGDim).Render("  ↓ more"))
	}

	if len(state.Outdated) == 0 {
		lines = append(lines, lipgloss.NewStyle().Foreground(config.ColorFGDim).Render("  All packages up to date"))
	}

	return style.Render(strings.Join(lines, "\n"))
}

func (d DashboardModel) renderStatusBar(state AppState, w int) string {
	style := lipgloss.NewStyle().
		Background(config.ColorBGElevated).
		Foreground(config.ColorFG).
		Width(w).
		Padding(0, 1)

	venvName := lipgloss.NewStyle().Foreground(config.ColorGreen).Render(state.Venv.Name())

	var mgrStyle lipgloss.Style
	if state.Manager.Type == 1 {
		mgrStyle = lipgloss.NewStyle().Foreground(config.ColorPurple)
	} else {
		mgrStyle = lipgloss.NewStyle().Foreground(config.ColorFGDim)
	}
	mgr := mgrStyle.Render(state.Manager.String())

	pkgCount := fmt.Sprintf("%d pkgs", len(state.Installed))

	outdatedCount := fmt.Sprintf("%d outdated", len(state.Outdated))
	if len(state.Outdated) > 0 {
		outdatedCount = lipgloss.NewStyle().Foreground(config.ColorRed).Render(outdatedCount)
	} else {
		outdatedCount = lipgloss.NewStyle().Foreground(config.ColorFGDim).Render(outdatedCount)
	}

	help := lipgloss.NewStyle().Foreground(config.ColorFGDim).Render("? help")

	status := state.StatusMsg
	if status != "" {
		status = " │ " + lipgloss.NewStyle().Foreground(config.ColorOrange).Render(status)
	}

	bar := fmt.Sprintf(" depman │ %s │ %s │ %s │ %s │ %s%s",
		venvName, mgr, pkgCount, outdatedCount, help, status)

	return style.Render(bar)
}

func (d DashboardModel) renderConfirmDialog(w int) string {
	style := lipgloss.NewStyle().
		Foreground(config.ColorYellow).
		Width(w).
		Padding(0, 1)
	return style.Render(fmt.Sprintf("  %s %s? [y/N] ", d.confirmAction, d.confirmPkg))
}

func (d DashboardModel) renderAddInput(w int) string {
	style := lipgloss.NewStyle().
		Foreground(config.ColorBlue).
		Width(w).
		Padding(0, 1)
	cursor := lipgloss.NewStyle().Foreground(config.ColorOrange).Render("█")
	return style.Render(fmt.Sprintf("  Add package: %s%s", d.addInput, cursor))
}

func (d DashboardModel) renderLoading(w, h int) string {
	style := lipgloss.NewStyle().
		Foreground(config.ColorFGDim).
		Width(w).
		Height(h).
		Align(lipgloss.Center, lipgloss.Center)
	return style.Render("Loading packages...")
}

func (d DashboardModel) selectedPackage(state *AppState) *pip.Package {
	if state.ActivePanel == PanelInstalled && len(state.Installed) > 0 && d.installedCursor < len(state.Installed) {
		return &state.Installed[d.installedCursor]
	}
	return nil
}

func (d DashboardModel) selectedOutdated(state *AppState) *pip.Package {
	if len(state.Outdated) > 0 && d.outdatedCursor < len(state.Outdated) {
		return &state.Outdated[d.outdatedCursor]
	}
	return nil
}

func (d DashboardModel) pageSize() int {
	if d.height > 6 {
		return d.height - 6
	}
	return DefaultPageSize
}
