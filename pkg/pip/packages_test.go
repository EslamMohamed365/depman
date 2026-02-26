package pip

import (
	"testing"

	"github.com/eslam/depman/config"
)

func TestParseVersion(t *testing.T) {
	tests := []struct {
		name     string
		version  string
		expected []int
	}{
		{
			name:     "standard version",
			version:  "1.2.3",
			expected: []int{1, 2, 3},
		},
		{
			name:     "v-prefix version",
			version:  "v2.0.1",
			expected: []int{2, 0, 1},
		},
		{
			name:     "pre-release version",
			version:  "1.0.0rc1",
			expected: []int{1, 0, 0},
		},
		{
			name:     "short version",
			version:  "1.0",
			expected: []int{1, 0, 0},
		},
		{
			name:     "single number version",
			version:  "2",
			expected: []int{2, 0, 0},
		},
		{
			name:     "invalid version",
			version:  "abc",
			expected: nil,
		},
		{
			name:     "mixed version",
			version:  "1.2.beta3",
			expected: nil, // dot before non-numeric fails parsing
		},
		{
			name:     "alpha version",
			version:  "3.0.0a1",
			expected: []int{3, 0, 0},
		},
		{
			name:     "beta version",
			version:  "2.5.1b2",
			expected: []int{2, 5, 1},
		},
		{
			name:     "dev version",
			version:  "1.2.3.dev456",
			expected: nil, // dot before non-numeric fails parsing
		},
		{
			name:     "post-release version",
			version:  "1.0.0.post1",
			expected: nil, // dot before non-numeric fails parsing
		},
		{
			name:     "empty string",
			version:  "",
			expected: nil,
		},
		{
			name:     "very short version",
			version:  "v1",
			expected: []int{1, 0, 0},
		},
		{
			name:     "four-part version",
			version:  "1.2.3.4",
			expected: []int{1, 2, 3, 4},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseVersion(tt.version)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("parseVersion(%q) = %v, want nil", tt.version, result)
				}
				return
			}

			if result == nil {
				t.Errorf("parseVersion(%q) = nil, want %v", tt.version, tt.expected)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("parseVersion(%q) length = %d, want %d", tt.version, len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i] != tt.expected[i] {
					t.Errorf("parseVersion(%q)[%d] = %d, want %d", tt.version, i, result[i], tt.expected[i])
				}
			}
		})
	}
}

func TestComputeDiff(t *testing.T) {
	tests := []struct {
		name     string
		current  string
		latest   string
		expected config.DiffType
	}{
		{
			name:     "major version upgrade",
			current:  "1.0.0",
			latest:   "2.0.0",
			expected: config.DiffMajor,
		},
		{
			name:     "minor version upgrade",
			current:  "1.0.0",
			latest:   "1.1.0",
			expected: config.DiffMinor,
		},
		{
			name:     "patch version upgrade",
			current:  "1.0.0",
			latest:   "1.0.1",
			expected: config.DiffPatch,
		},
		{
			name:     "no difference",
			current:  "1.0.0",
			latest:   "1.0.0",
			expected: config.DiffNone,
		},
		{
			name:     "invalid current version",
			current:  "abc",
			latest:   "1.0.0",
			expected: config.DiffUnknown,
		},
		{
			name:     "invalid latest version",
			current:  "1.0.0",
			latest:   "xyz",
			expected: config.DiffUnknown,
		},
		{
			name:     "both invalid",
			current:  "abc",
			latest:   "xyz",
			expected: config.DiffUnknown,
		},
		{
			name:     "pre-release to release",
			current:  "1.0.0rc1",
			latest:   "1.0.0",
			expected: config.DiffNone,
		},
		{
			name:     "with v-prefix",
			current:  "v1.0.0",
			latest:   "v1.1.0",
			expected: config.DiffMinor,
		},
		{
			name:     "short to full version",
			current:  "1.0",
			latest:   "1.0.1",
			expected: config.DiffPatch,
		},
		{
			name:     "major downgrade",
			current:  "2.0.0",
			latest:   "1.0.0",
			expected: config.DiffMajor,
		},
		{
			name:     "multiple major versions",
			current:  "1.0.0",
			latest:   "5.0.0",
			expected: config.DiffMajor,
		},
		{
			name:     "complex pre-release",
			current:  "2.0.0a1",
			latest:   "2.1.0",
			expected: config.DiffMinor,
		},
		{
			name:     "empty strings",
			current:  "",
			latest:   "",
			expected: config.DiffUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ComputeDiff(tt.current, tt.latest)
			if result != tt.expected {
				t.Errorf("ComputeDiff(%q, %q) = %v, want %v",
					tt.current, tt.latest, result, tt.expected)
			}
		})
	}
}

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		name     string
		v1       string
		v2       string
		expected string // "less", "equal", "greater"
	}{
		{
			name:     "v1 less than v2",
			v1:       "1.0.0",
			v2:       "2.0.0",
			expected: "less",
		},
		{
			name:     "v1 equal to v2",
			v1:       "1.0.0",
			v2:       "1.0.0",
			expected: "equal",
		},
		{
			name:     "v1 greater than v2",
			v1:       "2.0.0",
			v2:       "1.0.0",
			expected: "greater",
		},
		{
			name:     "minor version comparison",
			v1:       "1.5.0",
			v2:       "1.10.0",
			expected: "less",
		},
		{
			name:     "patch version comparison",
			v1:       "1.0.5",
			v2:       "1.0.15",
			expected: "less",
		},
		{
			name:     "pre-release stripped",
			v1:       "1.0.0rc1",
			v2:       "1.0.0",
			expected: "equal",
		},
		{
			name:     "different pre-releases",
			v1:       "1.0.0a1",
			v2:       "1.0.0b2",
			expected: "equal",
		},
		{
			name:     "short vs full version",
			v1:       "1.0",
			v2:       "1.0.0",
			expected: "equal",
		},
		{
			name:     "with v-prefix",
			v1:       "v1.0.0",
			v2:       "v1.0.1",
			expected: "less",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parts1 := parseVersion(tt.v1)
			parts2 := parseVersion(tt.v2)

			if parts1 == nil || parts2 == nil {
				t.Skip("Invalid version strings")
			}

			// Compare versions lexicographically
			minLen := len(parts1)
			if len(parts2) < minLen {
				minLen = len(parts2)
			}

			result := "equal"
			for i := 0; i < minLen; i++ {
				if parts1[i] < parts2[i] {
					result = "less"
					break
				}
				if parts1[i] > parts2[i] {
					result = "greater"
					break
				}
			}

			if result == "equal" && len(parts1) != len(parts2) {
				if len(parts1) < len(parts2) {
					result = "less"
				} else {
					result = "greater"
				}
			}

			if result != tt.expected {
				t.Errorf("compare(%q, %q) = %s, want %s",
					tt.v1, tt.v2, result, tt.expected)
			}
		})
	}
}

func TestParsePackageList(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expected    []Package
		expectError bool
	}{
		{
			name: "valid package list",
			jsonData: `[
				{"name": "requests", "version": "2.28.1"},
				{"name": "django", "version": "4.2.0"}
			]`,
			expected: []Package{
				{Name: "requests", InstalledVersion: "2.28.1"},
				{Name: "django", InstalledVersion: "4.2.0"},
			},
			expectError: false,
		},
		{
			name:        "empty list",
			jsonData:    `[]`,
			expected:    []Package{},
			expectError: false,
		},
		{
			name: "single package",
			jsonData: `[
				{"name": "flask", "version": "2.3.0"}
			]`,
			expected: []Package{
				{Name: "flask", InstalledVersion: "2.3.0"},
			},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			jsonData:    `{invalid json}`,
			expected:    nil,
			expectError: true,
		},
		{
			name:        "malformed JSON",
			jsonData:    `[{"name": "pkg", "version":}]`,
			expected:    nil,
			expectError: true,
		},
		{
			name: "package with pre-release version",
			jsonData: `[
				{"name": "pytest", "version": "7.4.0rc1"}
			]`,
			expected: []Package{
				{Name: "pytest", InstalledVersion: "7.4.0rc1"},
			},
			expectError: false,
		},
		{
			name: "package with v-prefix version",
			jsonData: `[
				{"name": "mypackage", "version": "v1.2.3"}
			]`,
			expected: []Package{
				{Name: "mypackage", InstalledVersion: "v1.2.3"},
			},
			expectError: false,
		},
		{
			name:        "empty string",
			jsonData:    ``,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParsePackageList(tt.jsonData)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParsePackageList() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParsePackageList() unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("ParsePackageList() returned %d packages, want %d",
					len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i].Name != tt.expected[i].Name {
					t.Errorf("Package[%d].Name = %q, want %q",
						i, result[i].Name, tt.expected[i].Name)
				}
				if result[i].InstalledVersion != tt.expected[i].InstalledVersion {
					t.Errorf("Package[%d].InstalledVersion = %q, want %q",
						i, result[i].InstalledVersion, tt.expected[i].InstalledVersion)
				}
			}
		})
	}
}

func TestParseOutdatedList(t *testing.T) {
	tests := []struct {
		name        string
		jsonData    string
		expected    []Package
		expectError bool
	}{
		{
			name: "outdated packages with different diff types",
			jsonData: `[
				{"name": "requests", "version": "2.28.1", "latest_version": "2.28.2", "latest_filetype": "wheel"},
				{"name": "django", "version": "4.1.0", "latest_version": "4.2.0", "latest_filetype": "wheel"},
				{"name": "flask", "version": "1.0.0", "latest_version": "2.0.0", "latest_filetype": "wheel"}
			]`,
			expected: []Package{
				{
					Name:             "requests",
					InstalledVersion: "2.28.1",
					LatestVersion:    "2.28.2",
					DiffType:         config.DiffPatch,
					IsOutdated:       true,
				},
				{
					Name:             "django",
					InstalledVersion: "4.1.0",
					LatestVersion:    "4.2.0",
					DiffType:         config.DiffMinor,
					IsOutdated:       true,
				},
				{
					Name:             "flask",
					InstalledVersion: "1.0.0",
					LatestVersion:    "2.0.0",
					DiffType:         config.DiffMajor,
					IsOutdated:       true,
				},
			},
			expectError: false,
		},
		{
			name:        "empty outdated list",
			jsonData:    `[]`,
			expected:    []Package{},
			expectError: false,
		},
		{
			name: "single outdated package",
			jsonData: `[
				{"name": "pytest", "version": "7.0.0", "latest_version": "7.4.0", "latest_filetype": "wheel"}
			]`,
			expected: []Package{
				{
					Name:             "pytest",
					InstalledVersion: "7.0.0",
					LatestVersion:    "7.4.0",
					DiffType:         config.DiffMinor,
					IsOutdated:       true,
				},
			},
			expectError: false,
		},
		{
			name:        "invalid JSON",
			jsonData:    `{invalid}`,
			expected:    nil,
			expectError: true,
		},
		{
			name: "pre-release versions",
			jsonData: `[
				{"name": "mypackage", "version": "1.0.0rc1", "latest_version": "1.0.0", "latest_filetype": "wheel"}
			]`,
			expected: []Package{
				{
					Name:             "mypackage",
					InstalledVersion: "1.0.0rc1",
					LatestVersion:    "1.0.0",
					DiffType:         config.DiffNone,
					IsOutdated:       true,
				},
			},
			expectError: false,
		},
		{
			name: "invalid version formats",
			jsonData: `[
				{"name": "badpkg", "version": "abc", "latest_version": "xyz", "latest_filetype": "wheel"}
			]`,
			expected: []Package{
				{
					Name:             "badpkg",
					InstalledVersion: "abc",
					LatestVersion:    "xyz",
					DiffType:         config.DiffUnknown,
					IsOutdated:       true,
				},
			},
			expectError: false,
		},
		{
			name: "short version formats",
			jsonData: `[
				{"name": "shortver", "version": "1.0", "latest_version": "1.1", "latest_filetype": "wheel"}
			]`,
			expected: []Package{
				{
					Name:             "shortver",
					InstalledVersion: "1.0",
					LatestVersion:    "1.1",
					DiffType:         config.DiffMinor,
					IsOutdated:       true,
				},
			},
			expectError: false,
		},
		{
			name:        "empty string",
			jsonData:    ``,
			expected:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseOutdatedList(tt.jsonData)

			if tt.expectError {
				if err == nil {
					t.Errorf("ParseOutdatedList() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseOutdatedList() unexpected error: %v", err)
				return
			}

			if len(result) != len(tt.expected) {
				t.Errorf("ParseOutdatedList() returned %d packages, want %d",
					len(result), len(tt.expected))
				return
			}

			for i := range result {
				if result[i].Name != tt.expected[i].Name {
					t.Errorf("Package[%d].Name = %q, want %q",
						i, result[i].Name, tt.expected[i].Name)
				}
				if result[i].InstalledVersion != tt.expected[i].InstalledVersion {
					t.Errorf("Package[%d].InstalledVersion = %q, want %q",
						i, result[i].InstalledVersion, tt.expected[i].InstalledVersion)
				}
				if result[i].LatestVersion != tt.expected[i].LatestVersion {
					t.Errorf("Package[%d].LatestVersion = %q, want %q",
						i, result[i].LatestVersion, tt.expected[i].LatestVersion)
				}
				if result[i].DiffType != tt.expected[i].DiffType {
					t.Errorf("Package[%d].DiffType = %v, want %v",
						i, result[i].DiffType, tt.expected[i].DiffType)
				}
				if result[i].IsOutdated != tt.expected[i].IsOutdated {
					t.Errorf("Package[%d].IsOutdated = %v, want %v",
						i, result[i].IsOutdated, tt.expected[i].IsOutdated)
				}
			}
		})
	}
}
