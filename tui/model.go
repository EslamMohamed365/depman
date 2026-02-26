package tui

import (
	"fmt"
	"github.com/eslam/depman/config"
	"github.com/eslam/depman/pkg/detector"
	"github.com/eslam/depman/pkg/env"
	"github.com/eslam/depman/pkg/log"
	"github.com/eslam/depman/pkg/pip"

	tea "github.com/charmbracelet/bubbletea"
)

// Screen represents the current active screen.
type Screen int

const (
	ScreenInit Screen = iota
	ScreenDashboard
	ScreenSearch
	ScreenHelp
)

// Panel represents which dashboard panel is focused.
type Panel int

const (
	PanelInstalled Panel = iota
	PanelOutdated
)

// AppState holds the full application state.
type AppState struct {
	Screen           Screen
	Project          detector.Project
	Venv             env.Virtualenv
	Manager          env.PackageManager
	Config           config.Config
	Installed        []pip.Package
	Outdated         []pip.Package
	ActivePanel      Panel
	StatusMsg        string
	IsLoading        bool
	Width            int
	Height           int
	VersionChangePkg string // set when user wants to change version of installed pkg
}

// NewAppState creates the initial application state from detection results.
func NewAppState(project detector.Project, venv env.Virtualenv, mgr env.PackageManager, cfg config.Config) AppState {
	screen := ScreenDashboard
	if !project.Detected() {
		log.Debug("no project detected, showing init screen")
		screen = ScreenInit
	} else {
		log.Debug("project detected", "file_type", project.FileType.String(), "path", project.FilePath)
	}

	return AppState{
		Screen:  screen,
		Project: project,
		Venv:    venv,
		Manager: mgr,
		Config:  cfg,
	}
}

// Model is the root Bubble Tea model.
type Model struct {
	state     AppState
	runner    *pip.Runner
	dashboard DashboardModel
	initView  InitModel
	search    SearchModel
	help      HelpModel
	Err       error
}

// NewModel creates the root model from the initial state.
func NewModel(state AppState) Model {
	runner := pip.NewRunner(state.Manager, state.Venv)
	return Model{
		state:     state,
		runner:    runner,
		dashboard: NewDashboardModel(state),
		initView:  NewInitModel(state),
		search:    NewSearchModel(),
		help:      NewHelpModel(),
	}
}

// -- Messages --

// PackagesLoadedMsg is sent when installed packages have been loaded.
type PackagesLoadedMsg struct {
	Installed []pip.Package
	Outdated  []pip.Package
	Err       error
}

// PackageActionMsg is sent when a package operation completes.
type PackageActionMsg struct {
	Action  string // "install", "uninstall", "upgrade"
	Package string
	Err     error
}

// StatusNotification sets a temporary status notification.
type StatusNotification struct {
	Text    string
	IsError bool
}

// -- Tea interface --

func (m Model) Init() tea.Cmd {
	log.Debug("tui initialized", "screen", m.state.Screen)
	if m.state.Screen == ScreenDashboard {
		return m.loadPackages()
	}
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Global keybindings — only handle when NOT in search or add mode
		if m.state.Screen != ScreenSearch && !m.dashboard.addMode {
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			case "q":
				if m.state.Screen == ScreenDashboard && !m.dashboard.showConfirm {
					return m, tea.Quit
				}
			case "?":
				if m.state.Screen == ScreenHelp {
					m.state.Screen = ScreenDashboard
					return m, nil
				}
				if m.state.Screen == ScreenDashboard {
					m.state.Screen = ScreenHelp
			log.Debug("screen changed", "from", m.state.Screen, "to", ScreenHelp)
					return m, nil
				}
			case "esc":
				switch m.state.Screen {
				case ScreenHelp:
					m.state.Screen = ScreenDashboard
				log.Debug("screen changed", "from", ScreenHelp, "to", ScreenDashboard)
					return m, nil
				}
			}
		} else if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.state.Width = msg.Width
		m.state.Height = msg.Height
		m.dashboard.SetSize(msg.Width, msg.Height)
		return m, nil

	case PackagesLoadedMsg:
		m.state.IsLoading = false
		log.Debug("packages loaded", "installed", len(msg.Installed), "outdated", len(msg.Outdated))
		if msg.Err != nil {
			m.state.StatusMsg = "Failed to load packages: " + msg.Err.Error()
		} else {
			m.state.Installed = msg.Installed
			m.state.Outdated = msg.Outdated
			m.dashboard.UpdatePackages(msg.Installed, msg.Outdated)
		}
		return m, nil

	case PackageActionMsg:
		m.state.IsLoading = false
		log.Debug("package action completed", "action", msg.Action, "package", msg.Package, "success", msg.Err == nil)
		if msg.Err != nil {
			m.state.StatusMsg = "Failed: " + msg.Err.Error()
		} else {
			m.state.StatusMsg = msg.Action + " " + msg.Package + " ✓"
			return m, m.loadPackages()
		}
		return m, nil

	case ProjectCreatedMsg:
		// Forward to init model
		m.initView, _ = m.initView.Update(msg, &m.state)
		if m.state.Screen == ScreenDashboard {
			// Project was created, reload runner and load packages
		log.Debug("project created, switching to dashboard", "manager", m.state.Manager, "venv", m.state.Venv.Path)
			m.runner = pip.NewRunner(m.state.Manager, m.state.Venv)
			return m, m.loadPackages()
		}
		return m, nil

	case SearchResultsMsg:
		// Forward to search model
		m.search, _ = m.search.Update(msg, &m.state, m.runner)
		return m, nil

	case PackageDetailMsg:
		// Forward to search model
		m.search, _ = m.search.Update(msg, &m.state, m.runner)
		return m, nil
	}

	// Delegate to active screen
	var cmd tea.Cmd
	switch m.state.Screen {
	case ScreenInit:
		m.initView, cmd = m.initView.Update(msg, &m.state)
	case ScreenDashboard:
		m.dashboard, cmd = m.dashboard.Update(msg, &m.state, m.runner)
	case ScreenSearch:
		// If VersionChangePkg is set, auto-fetch detail for that package
		if m.state.VersionChangePkg != "" {
			pkgName := m.state.VersionChangePkg
			log.Debug("version change requested", "package", pkgName)
			m.state.VersionChangePkg = "" // consume it
			m.search = SearchModel{
				phase:         PhaseInput,
				input:         pkgName,
				detailLoading: true,
			}
			cmd = m.search.fetchDetail(&m.state, pkgName)
		} else {
			m.search, cmd = m.search.Update(msg, &m.state, m.runner)
		}
	case ScreenHelp:
		m.help, cmd = m.help.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.state.Screen {
	case ScreenInit:
		return m.initView.View(m.state)
	case ScreenHelp:
		return m.help.View(m.state)
	case ScreenSearch:
		return m.search.View(m.state)
	case ScreenDashboard:
		return m.dashboard.View(m.state)
	default:
		return m.dashboard.View(m.state)
	}
}

// loadPackages returns a Cmd that fetches installed and outdated packages.
func (m Model) loadPackages() tea.Cmd {
	runner := m.runner
	return func() tea.Msg {
		listResult := runner.List()
		if listResult.Err != nil {
			return PackagesLoadedMsg{Err: listResult.Err}
		}
		installed, err := pip.ParsePackageList(listResult.Stdout)
		if err != nil {
			return PackagesLoadedMsg{Err: err}
		}

		outdatedResult := runner.Outdated()
		if outdatedResult.Err != nil {
			return PackagesLoadedMsg{Err: outdatedResult.Err}
		}
		outdated, err := pip.ParseOutdatedList(outdatedResult.Stdout)
		if err != nil {
			return PackagesLoadedMsg{Err: fmt.Errorf("pip: parse outdated list: %w", err)}
		}

		return PackagesLoadedMsg{Installed: installed, Outdated: outdated}
	}
}
