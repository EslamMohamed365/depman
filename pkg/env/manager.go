package env

import (
	"os/exec"
)

// ManagerType represents the package manager to use.
type ManagerType int

const (
	ManagerNone ManagerType = iota
	ManagerUV
	ManagerPip
)

// PackageManager holds info about the detected package manager.
type PackageManager struct {
	Type    ManagerType
	BinPath string
}

// String returns the display name for the package manager.
func (m PackageManager) String() string {
	switch m.Type {
	case ManagerUV:
		return "⚡ uv"
	case ManagerPip:
		return "pip"
	default:
		return "none"
	}
}

// InstallCmd returns the install command for a package.
func (m PackageManager) InstallCmd(pkg string) (string, []string) {
	switch m.Type {
	case ManagerUV:
		return m.BinPath, []string{"pip", "install", pkg}
	default: // pip
		return m.BinPath, []string{"install", pkg}
	}
}

// UninstallCmd returns the uninstall command for a package.
func (m PackageManager) UninstallCmd(pkg string) (string, []string) {
	switch m.Type {
	case ManagerUV:
		return m.BinPath, []string{"pip", "uninstall", pkg}
	default:
		return m.BinPath, []string{"uninstall", pkg, "-y"}
	}
}

// UpgradeCmd returns the upgrade command for a package.
func (m PackageManager) UpgradeCmd(pkg string) (string, []string) {
	switch m.Type {
	case ManagerUV:
		return m.BinPath, []string{"pip", "install", "--upgrade", pkg}
	default:
		return m.BinPath, []string{"install", "--upgrade", pkg}
	}
}

// ListCmd returns the command to list installed packages.
func (m PackageManager) ListCmd() (string, []string) {
	switch m.Type {
	case ManagerUV:
		return m.BinPath, []string{"pip", "list", "--format", "json"}
	default:
		return m.BinPath, []string{"list", "--format", "json"}
	}
}

// OutdatedCmd returns the command to list outdated packages.
func (m PackageManager) OutdatedCmd() (string, []string) {
	switch m.Type {
	case ManagerUV:
		return m.BinPath, []string{"pip", "list", "--outdated", "--format", "json"}
	default:
		return m.BinPath, []string{"list", "--outdated", "--format", "json"}
	}
}

// DetectPackageManager finds the available package manager.
// Priority: uv (preferred) → pip → pip3.
// If preferred is set and available, use it regardless.
func DetectPackageManager(preferred string) PackageManager {
	// If user has a preference, try it first
	if preferred != "" {
		if path, err := exec.LookPath(preferred); err == nil {
			mgrType := ManagerPip
			if preferred == "uv" {
				mgrType = ManagerUV
			}
			return PackageManager{Type: mgrType, BinPath: path}
		}
	}

	// Auto-detect: uv first
	if path, err := exec.LookPath("uv"); err == nil {
		return PackageManager{Type: ManagerUV, BinPath: path}
	}

	// Try pip
	if path, err := exec.LookPath("pip"); err == nil {
		return PackageManager{Type: ManagerPip, BinPath: path}
	}

	// Try pip3
	if path, err := exec.LookPath("pip3"); err == nil {
		return PackageManager{Type: ManagerPip, BinPath: path}
	}

	return PackageManager{Type: ManagerNone}
}
