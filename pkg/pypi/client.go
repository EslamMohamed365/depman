package pypi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

// SearchResult represents a PyPI package from search.
type SearchResult struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"summary"`
}

// PackageDetail has full info about a package including all versions.
type PackageDetail struct {
	Name       string
	Version    string // latest
	Summary    string
	Author     string
	License    string
	HomePage   string
	Versions   []string // sorted newest first
	RequiresPy string
}

// packageInfo matches the PyPI JSON API response for a single package.
type packageInfo struct {
	Info struct {
		Name           string `json:"name"`
		Version        string `json:"version"`
		Summary        string `json:"summary"`
		Author         string `json:"author"`
		License        string `json:"license"`
		HomePage       string `json:"home_page"`
		RequiresPython string `json:"requires_python"`
		ProjectURL     string `json:"project_url"`
	} `json:"info"`
	Releases map[string]json.RawMessage `json:"releases"`
}

// Client is a PyPI API client.
type Client struct {
	BaseURL    string
	httpClient *http.Client
}

// NewClient creates a new PyPI client.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://pypi.org"
	}
	return &Client{
		BaseURL: strings.TrimRight(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetPackage fetches info for a single package from PyPI.
func (c *Client) GetPackage(name string) (*SearchResult, error) {
	url := fmt.Sprintf("%s/pypi/%s/json", c.BaseURL, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("PyPI returned status %d for %s", resp.StatusCode, name)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response for %s: %w", name, err)
	}

	var pkg packageInfo
	if err := json.Unmarshal(body, &pkg); err != nil {
		return nil, fmt.Errorf("parsing response for %s: %w", name, err)
	}

	return &SearchResult{
		Name:        pkg.Info.Name,
		Version:     pkg.Info.Version,
		Description: pkg.Info.Summary,
	}, nil
}

// GetPackageDetail fetches full details including all available versions.
func (c *Client) GetPackageDetail(name string) (*PackageDetail, error) {
	url := fmt.Sprintf("%s/pypi/%s/json", c.BaseURL, name)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching %s: %w", name, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 404 {
		return nil, nil
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("PyPI returned status %d for %s", resp.StatusCode, name)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response: %w", err)
	}

	var pkg packageInfo
	if err := json.Unmarshal(body, &pkg); err != nil {
		return nil, fmt.Errorf("parsing response: %w", err)
	}

	// Extract versions from releases map and sort newest-first
	versions := make([]string, 0, len(pkg.Releases))
	for v := range pkg.Releases {
		// Skip pre-release versions (contain a, b, rc, dev, post)
		if isStableVersion(v) {
			versions = append(versions, v)
		}
	}
	sortVersionsDesc(versions)

	// Keep max 20 versions
	if len(versions) > 20 {
		versions = versions[:20]
	}

	license := pkg.Info.License
	if len(license) > 40 {
		license = license[:40] + "…"
	}

	return &PackageDetail{
		Name:       pkg.Info.Name,
		Version:    pkg.Info.Version,
		Summary:    pkg.Info.Summary,
		Author:     pkg.Info.Author,
		License:    license,
		HomePage:   pkg.Info.HomePage,
		Versions:   versions,
		RequiresPy: pkg.Info.RequiresPython,
	}, nil
}

// Search queries PyPI for packages matching the given query.
func (c *Client) Search(query string) ([]SearchResult, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return nil, nil
	}

	var results []SearchResult
	seen := make(map[string]bool)

	// Try exact match first
	if pkg, err := c.GetPackage(query); err == nil && pkg != nil {
		results = append(results, *pkg)
		seen[strings.ToLower(pkg.Name)] = true
	}

	// Try common variations
	variations := []string{
		"python-" + query,
		"py" + query,
		query + "-python",
		query + "lib",
		query + "-py",
	}

	for _, v := range variations {
		if seen[v] {
			continue
		}
		if pkg, err := c.GetPackage(v); err == nil && pkg != nil {
			name := strings.ToLower(pkg.Name)
			if !seen[name] {
				seen[name] = true
				results = append(results, *pkg)
			}
		}
		if len(results) >= 10 {
			break
		}
	}

	return results, nil
}

func isStableVersion(v string) bool {
	lower := strings.ToLower(v)
	for _, pre := range []string{"a", "b", "rc", "dev", "alpha", "beta", "post"} {
		if strings.Contains(lower, pre) {
			return false
		}
	}
	return true
}

// sortVersionsDesc sorts semver-like strings in descending order.
func sortVersionsDesc(versions []string) {
	sort.Slice(versions, func(i, j int) bool {
		return compareVersions(versions[i], versions[j]) > 0
	})
}

func compareVersions(a, b string) int {
	pa := strings.Split(a, ".")
	pb := strings.Split(b, ".")
	maxLen := len(pa)
	if len(pb) > maxLen {
		maxLen = len(pb)
	}
	for i := 0; i < maxLen; i++ {
		var na, nb int
		if i < len(pa) {
			na, _ = strconv.Atoi(pa[i])
		}
		if i < len(pb) {
			nb, _ = strconv.Atoi(pb[i])
		}
		if na != nb {
			return na - nb
		}
	}
	return 0
}
