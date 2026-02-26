package pip

import (
	"encoding/json"
	"strconv"
	"strings"

	"github.com/eslam/depman/config"
)

// Package represents an installed Python package.
type Package struct {
	Name             string          `json:"name"`
	InstalledVersion string          `json:"version"`
	LatestVersion    string          `json:"latest_version,omitempty"`
	Description      string          `json:"-"`
	DiffType         config.DiffType `json:"-"`
	IsOutdated       bool            `json:"-"`
}

// pipListEntry matches the JSON output of `pip list --format json`.
type pipListEntry struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// pipOutdatedEntry matches the JSON output of `pip list --outdated --format json`.
type pipOutdatedEntry struct {
	Name          string `json:"name"`
	Version       string `json:"version"`
	LatestVersion string `json:"latest_version"`
	LatestType    string `json:"latest_filetype"`
}

// ParsePackageList parses the JSON output from `pip list --format json`.
func ParsePackageList(jsonData string) ([]Package, error) {
	var entries []pipListEntry
	if err := json.Unmarshal([]byte(jsonData), &entries); err != nil {
		return nil, err
	}

	packages := make([]Package, len(entries))
	for i, e := range entries {
		packages[i] = Package{
			Name:             e.Name,
			InstalledVersion: e.Version,
		}
	}
	return packages, nil
}

// ParseOutdatedList parses the JSON output from `pip list --outdated --format json`.
func ParseOutdatedList(jsonData string) ([]Package, error) {
	var entries []pipOutdatedEntry
	if err := json.Unmarshal([]byte(jsonData), &entries); err != nil {
		return nil, err
	}

	packages := make([]Package, len(entries))
	for i, e := range entries {
		diff := ComputeDiff(e.Version, e.LatestVersion)
		packages[i] = Package{
			Name:             e.Name,
			InstalledVersion: e.Version,
			LatestVersion:    e.LatestVersion,
			DiffType:         diff,
			IsOutdated:       true,
		}
	}
	return packages, nil
}

// ComputeDiff classifies the semver difference between two version strings.
func ComputeDiff(current, latest string) config.DiffType {
	curParts := parseVersion(current)
	latParts := parseVersion(latest)

	if len(curParts) < 3 || len(latParts) < 3 {
		return config.DiffUnknown
	}

	if curParts[0] != latParts[0] {
		return config.DiffMajor
	}
	if curParts[1] != latParts[1] {
		return config.DiffMinor
	}
	if curParts[2] != latParts[2] {
		return config.DiffPatch
	}
	return config.DiffNone
}

func parseVersion(v string) []int {
	// Strip pre-release suffixes like "1.2.3rc1" → "1.2.3"
	for i, c := range v {
		if c != '.' && (c < '0' || c > '9') {
			v = v[:i]
			break
		}
	}

	parts := strings.Split(v, ".")
	nums := make([]int, len(parts))
	for i, p := range parts {
		n, err := strconv.Atoi(p)
		if err != nil {
			return nil
		}
		nums[i] = n
	}
	return nums
}
