package config

import (
	"os"
	"path/filepath"

	toml "github.com/pelletier/go-toml/v2"
)

// Config holds user configuration loaded from ~/.config/depman/config.toml.
type Config struct {
	PackageManager PackageManagerConfig `toml:"package_manager"`
	PyPI           PyPIConfig           `toml:"pypi"`
	Theme          ThemeConfig          `toml:"theme"`
}

// PackageManagerConfig specifies the preferred package manager.
type PackageManagerConfig struct {
	Preferred string `toml:"preferred"` // "uv" | "pip" | "pip3" | "" (auto)
}

// PyPIConfig specifies PyPI connection settings.
type PyPIConfig struct {
	Mirror string `toml:"mirror"` // default: "https://pypi.org"
}

// ThemeConfig specifies the color theme.
type ThemeConfig struct {
	Name string `toml:"name"` // default: "tokyo-night"
}

// DefaultConfig returns the default configuration values.
func DefaultConfig() Config {
	return Config{
		PackageManager: PackageManagerConfig{
			Preferred: "", // auto-detect
		},
		PyPI: PyPIConfig{
			Mirror: "https://pypi.org",
		},
		Theme: ThemeConfig{
			Name: "tokyo-night",
		},
	}
}

// configPath returns the path to the config file.
func configPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "depman", "config.toml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".config", "depman", "config.toml")
}

// Load reads the config file if it exists, otherwise returns defaults.
func Load() Config {
	cfg := DefaultConfig()

	path := configPath()
	if path == "" {
		return cfg
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return cfg // file doesn't exist or is unreadable, use defaults
	}

	_ = toml.Unmarshal(data, &cfg)

	// Apply defaults for empty values
	if cfg.PyPI.Mirror == "" {
		cfg.PyPI.Mirror = "https://pypi.org"
	}
	if cfg.Theme.Name == "" {
		cfg.Theme.Name = "tokyo-night"
	}

	return cfg
}
