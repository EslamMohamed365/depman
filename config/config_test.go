package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.PackageManager.Preferred != "" {
		t.Errorf("PackageManager.Preferred = %q; want %q", cfg.PackageManager.Preferred, "")
	}

	if cfg.PyPI.Mirror != "https://pypi.org" {
		t.Errorf("PyPI.Mirror = %q; want %q", cfg.PyPI.Mirror, "https://pypi.org")
	}

	if cfg.Theme.Name != "tokyo-night" {
		t.Errorf("Theme.Name = %q; want %q", cfg.Theme.Name, "tokyo-night")
	}

	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q; want %q", cfg.LogLevel, "info")
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected Config
	}{
		{
			name: "all fields specified",
			content: `[package_manager]
preferred = "uv"

[pypi]
mirror = "https://custom-pypi.org"

[theme]
name = "dracula"`,
			expected: Config{
				PackageManager: PackageManagerConfig{
					Preferred: "uv",
				},
				PyPI: PyPIConfig{
					Mirror: "https://custom-pypi.org",
				},
				Theme: ThemeConfig{
					Name: "dracula",
				},
				LogLevel: "info",
			},
		},
		{
			name: "pip as package manager",
			content: `[package_manager]
preferred = "pip"

[pypi]
mirror = "https://pypi.org"

[theme]
name = "tokyo-night"`,
			expected: Config{
				PackageManager: PackageManagerConfig{
					Preferred: "pip",
				},
				PyPI: PyPIConfig{
					Mirror: "https://pypi.org",
				},
				Theme: ThemeConfig{
					Name: "tokyo-night",
				},
				LogLevel: "info",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			depmanDir := filepath.Join(tmpDir, "depman")
			if err := os.MkdirAll(depmanDir, 0755); err != nil {
				t.Fatalf("failed to create depman dir: %v", err)
			}

			configFile := filepath.Join(depmanDir, "config.toml")
			if err := os.WriteFile(configFile, []byte(tt.content), 0644); err != nil {
				t.Fatalf("failed to write test config: %v", err)
			}

			t.Setenv("XDG_CONFIG_HOME", tmpDir)

			cfg, err := Load()
			if err != nil {
				t.Fatalf("Load() error = %v; want nil", err)
			}

			if cfg.PackageManager.Preferred != tt.expected.PackageManager.Preferred {
				t.Errorf("PackageManager.Preferred = %q; want %q", cfg.PackageManager.Preferred, tt.expected.PackageManager.Preferred)
			}
			if cfg.PyPI.Mirror != tt.expected.PyPI.Mirror {
			t.Errorf("PyPI.Mirror = %q; want %q", cfg.PyPI.Mirror, tt.expected.PyPI.Mirror)
			}
			if cfg.Theme.Name != tt.expected.Theme.Name {
				t.Errorf("Theme.Name = %q; want %q", cfg.Theme.Name, tt.expected.Theme.Name)
			}
		})
	}
}

func TestLoadConfig_MissingFile(t *testing.T) {
	tmpDir := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v; want nil", err)
	}

	// Should return defaults
	if cfg.PyPI.Mirror != "https://pypi.org" {
		t.Errorf("PyPI.Mirror = %q; want %q", cfg.PyPI.Mirror, "https://pypi.org")
	}
	if cfg.Theme.Name != "tokyo-night" {
		t.Errorf("Theme.Name = %q; want %q", cfg.Theme.Name, "tokyo-night")
	}
	if cfg.LogLevel != "info" {
		t.Errorf("LogLevel = %q; want %q", cfg.LogLevel, "info")
	}
}

func TestLoadConfig_CorruptedTOML(t *testing.T) {
	tmpDir := t.TempDir()
	depmanDir := filepath.Join(tmpDir, "depman")
	if err := os.MkdirAll(depmanDir, 0755); err != nil {
		t.Fatalf("failed to create depman dir: %v", err)
	}

	configFile := filepath.Join(depmanDir, "config.toml")
	corruptedContent := `[package_manager]
preferred = "uv"
[pypi`
	if err := os.WriteFile(configFile, []byte(corruptedContent), 0644); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	t.Setenv("XDG_CONFIG_HOME", tmpDir)

	cfg, err := Load()
	if err == nil {
		t.Fatalf("Load() error = nil; want error for corrupted TOML")
	}

	// Should still return defaults on error
	if cfg.PyPI.Mirror != "https://pypi.org" {
		t.Errorf("PyPI.Mirror = %q; want %q", cfg.PyPI.Mirror, "https://pypi.org")
	}
	_ = cfg // suppress unused warning
}

func TestConfigPath(t *testing.T) {
	tests := []struct {
		name        string
		envValue    string
		expectedDir string
	}{
		{
			name:        "XDG_CONFIG_HOME set",
			envValue:    "/custom/config",
			expectedDir: "/custom/config/depman/config.toml",
		},
		{
			name:        "XDG_CONFIG_HOME empty uses home",
			envValue:    "",
			expectedDir: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv("XDG_CONFIG_HOME", tt.envValue)
			} else {
				os.Unsetenv("XDG_CONFIG_HOME")
			}

			path := configPath()
			if tt.expectedDir != "" && path != tt.expectedDir {
				t.Errorf("configPath() = %q; want %q", path, tt.expectedDir)
			}
		})
	}
}
