package config

import "github.com/charmbracelet/lipgloss"

// Tokyo Night color palette
var (
	ColorBG          = lipgloss.Color("#1a1b26")
	ColorBGElevated  = lipgloss.Color("#24283b")
	ColorBGHighlight = lipgloss.Color("#2e3250")
	ColorBorder      = lipgloss.Color("#414868")
	ColorFG          = lipgloss.Color("#c0caf5")
	ColorFGDim       = lipgloss.Color("#565f89")
	ColorBlue        = lipgloss.Color("#7aa2f7")
	ColorGreen       = lipgloss.Color("#9ece6a")
	ColorTeal        = lipgloss.Color("#2ac3de")
	ColorYellow      = lipgloss.Color("#e0af68")
	ColorRed         = lipgloss.Color("#f7768e")
	ColorPurple      = lipgloss.Color("#bb9af7")
	ColorCyan        = lipgloss.Color("#7dcfff")
	ColorOrange      = lipgloss.Color("#ff9e64")
)

// DiffType represents the semver update severity.
type DiffType int

const (
	DiffNone DiffType = iota
	DiffPatch
	DiffMinor
	DiffMajor
	DiffUnknown
)

// DiffColor returns the Lip Gloss color for a given DiffType.
func DiffColor(d DiffType) lipgloss.Color {
	switch d {
	case DiffPatch:
		return ColorTeal
	case DiffMinor:
		return ColorYellow
	case DiffMajor:
		return ColorRed
	default:
		return ColorFGDim
	}
}

// DiffLabel returns a human-readable label for a DiffType.
func DiffLabel(d DiffType) string {
	switch d {
	case DiffPatch:
		return "patch"
	case DiffMinor:
		return "minor"
	case DiffMajor:
		return "major"
	case DiffNone:
		return "up to date"
	default:
		return "unknown"
	}
}
