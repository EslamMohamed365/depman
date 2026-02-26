package tui

// Default UI Dimensions
const (
	// DefaultWidth is the fallback width when terminal size is unavailable
	DefaultWidth = 80

	// DefaultHeight is the fallback height when terminal size is unavailable
	DefaultHeight = 24
)

// Search Configuration
const (
	// MaxInputLength is the maximum length for search input (prevents injection attacks)
	MaxInputLength = 100

	// MinSearchLength is the minimum search query length
	MinSearchLength = 1
)

// UI Layout Constants
const (
	// MinDescriptionLength is the minimum space reserved for package descriptions
	MinDescriptionLength = 20

	// MinDescWidth is the minimum width for description display
	MinDescWidth = 30

	// MinVisibleResults is the minimum number of visible items in lists
	MinVisibleResults = 3

	// KeyColumnWidth is the width reserved for keybinding labels in help
	KeyColumnWidth = 14
)

// Panel Sizing
const (
	// MinPanelWidth is the minimum width for a panel
	MinPanelWidth = 20

	// MinPanelHeight is the minimum height for a panel
	MinPanelHeight = 5

	// PanelBorderLines is the number of lines reserved for panel borders
	PanelBorderLines = 2

	// PanelTitleLines is the number of lines reserved for panel titles
	PanelTitleLines = 1

	// PanelPaddingLines is the number of blank lines in a panel
	PanelPaddingLines = 1

	// DashboardOverlayLines is the number of lines for overlays
	DashboardOverlayLines = 1

	// StatusBarLines is the number of lines for the status bar
	StatusBarLines = 1

	// DashboardReservedLines is the total reserved lines in dashboard
	DashboardReservedLines = 4

	// InitPanelPadding is the padding for init screen panels
	InitPanelPadding = 2

	// InitMaxCursorPosition is the max cursor position in init screen
	InitMaxCursorPosition = 2

	// DashboardPaddingLines is the dashboard padding
	DashboardPaddingLines = 4
)

const (
	// ViewportHeaderLines is the number of lines before content starts
	ViewportHeaderLines = 8

	// ViewportMinHeight is the minimum viewport height
	ViewportMinHeight = 3

	// MaxVersionsDisplay is the number of versions to display before "more" indicator
	MaxVersionsDisplay = 16

	// HalfPageDivisor is used for ctrl+d/ctrl+u navigation
	HalfPageDivisor = 2

	// DefaultPageSize is the default page size for scrolling when height is unknown
	DefaultPageSize = 10
)
// File Permissions
const (
	// DefaultFilePermissions is the default permission for created files (0644)
	DefaultFilePermissions = 0644
)
