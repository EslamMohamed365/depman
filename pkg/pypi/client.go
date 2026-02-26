package pypi

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/eslam/depman/pkg/log"
)

const maxDisplayVersions = MaxDisplayVersions

// defaultHTTPClient is the shared HTTP client for all PyPI requests.
// Using a single client enables connection pooling for better performance.
var defaultHTTPClient = &http.Client{
	Timeout: HTTPTimeout,
}

var retryStatuses = map[int]bool{
	StatusInternalServerError: true,
	StatusBadGateway:          true,
	StatusServiceUnavailable:  true,
	StatusGatewayTimeout:      true,
	StatusTooManyRequests:     true,
}

// HTTPClient wraps http.Client with retry support
type HTTPClient struct {
	client *http.Client
}

// NewHTTPClient creates a new HTTP client with retry support.
// It uses a shared HTTP client for connection pooling.
func NewHTTPClient() *HTTPClient {
	return &HTTPClient{
		client: defaultHTTPClient, // Use shared client for connection pooling
	}
}

// DoWithRetry performs an HTTP request with retry logic
func (c *HTTPClient) DoWithRetry(ctx context.Context, req *http.Request) (*http.Response, error) {
	var lastErr error
	for attempt := 0; attempt <= MaxRetries; attempt++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		if attempt > 0 {
			backoff := float64(RetryDelay) * math.Pow(BackoffMultiplier, float64(attempt-1))
			sleepDuration := time.Duration(int(backoff)) * time.Millisecond
			log.Info("retrying pypi request", "attempt", attempt+1, "backoff_ms", sleepDuration.Milliseconds(), "url", req.URL.String())
			timer := time.NewTimer(sleepDuration)
			select {
			case <-ctx.Done():
				timer.Stop()
				return nil, ctx.Err()
			case <-timer.C:
			}
		}

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = err
			log.Warn("pypi request failed", "error", err, "url", req.URL.String())
			continue
		}

		log.Info("pypi response received", "status", resp.StatusCode, "url", req.URL.String())

		if respOK(resp.StatusCode) {
			return resp, nil
		}

		if !retryStatuses[resp.StatusCode] {
			log.Warn("pypi non-retryable error", "status", resp.StatusCode, "url", req.URL.String())
			return resp, nil
		}

		resp.Body.Close()
		lastErr = fmt.Errorf("pypi: server error: %d", resp.StatusCode)
		log.Warn("pypi retryable error", "status", resp.StatusCode, "url", req.URL.String())
	}

	log.Error("pypi request exhausted retries", "error", lastErr, "url", req.URL.String())
	return nil, lastErr
}

func respOK(status int) bool {
	return status >= StatusOKMin && status < StatusOKMax
}

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
	httpClient *HTTPClient
}

// NewClient creates a new PyPI client.
func NewClient(baseURL string) *Client {
	if baseURL == "" {
		baseURL = "https://pypi.org"
	}
	return &Client{
		BaseURL:    strings.TrimRight(baseURL, "/"),
		httpClient: NewHTTPClient(),
	}
}

// GetPackage fetches info for a single package from PyPI.
func (c *Client) GetPackage(name string) (*SearchResult, error) {
	return c.GetPackageWithContext(context.Background(), name)
}

// GetPackageWithContext fetches info for a single package with context support.
func (c *Client) GetPackageWithContext(ctx context.Context, name string) (*SearchResult, error) {
	url := fmt.Sprintf("%s/pypi/%s/json", c.BaseURL, name)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("pypi: create request: %w", err)
	}
	resp, err := c.httpClient.DoWithRetry(ctx, req)
	if err != nil {
		log.Error("failed to fetch package from pypi", "package", name, "error", err)
		return nil, fmt.Errorf("pypi: fetch package: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != StatusOK {
		return nil, fmt.Errorf("pypi: fetch package: status %d", resp.StatusCode)
	}

	var pkg packageInfo
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("pypi: parse response: %w", err)
	}

	return &SearchResult{
		Name:        pkg.Info.Name,
		Version:     pkg.Info.Version,
		Description: pkg.Info.Summary,
	}, nil
}

// GetPackageDetail fetches full details including all available versions.
func (c *Client) GetPackageDetail(name string) (*PackageDetail, error) {
	return c.GetPackageDetailWithContext(context.Background(), name)
}

// GetPackageDetailWithContext fetches full details with context support.
func (c *Client) GetPackageDetailWithContext(ctx context.Context, name string) (*PackageDetail, error) {
	url := fmt.Sprintf("%s/pypi/%s/json", c.BaseURL, name)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("pypi: create request: %w", err)
	}
	resp, err := c.httpClient.DoWithRetry(ctx, req)
	if err != nil {
		log.Error("failed to fetch package detail from pypi", "package", name, "error", err)
		return nil, fmt.Errorf("pypi: fetch package detail: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != StatusOK {
		return nil, fmt.Errorf("pypi: fetch package detail: status %d", resp.StatusCode)
	}

	var pkg packageInfo
	if err := json.NewDecoder(resp.Body).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("pypi: parse response: %w", err)
	}

	// Extract versions from releases map and sort newest-first
	var versions []string
	if pkg.Releases == nil {
		versions = []string{}
	} else {
		versions = make([]string, 0, len(pkg.Releases))
		for v := range pkg.Releases {
			// Skip pre-release versions (contain a, b, rc, dev, post)
			if isStableVersion(v) {
				versions = append(versions, v)
			}
		}
	}
	sortVersionsDesc(versions)

	// Keep max 20 versions
	if len(versions) > maxDisplayVersions {
		versions = versions[:MaxDisplayVersions]
	}

	license := pkg.Info.License
	if len(license) > MaxLicenseLength {
		license = license[:MaxLicenseLength] + "…"
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
	return c.SearchWithContext(context.Background(), query)
}

// SearchWithContext searches with context support.
func (c *Client) SearchWithContext(ctx context.Context, query string) ([]SearchResult, error) {
	query = strings.TrimSpace(strings.ToLower(query))
	if query == "" {
		return nil, nil
	}

	var results []SearchResult
	seen := make(map[string]bool)

	// Try exact match first (sequential, as required)
	if pkg, err := c.GetPackageWithContext(ctx, query); err == nil && pkg != nil {
		results = append(results, *pkg)
		seen[strings.ToLower(pkg.Name)] = true
	}

	// If we already have enough results, return early
	if len(results) >= MaxSearchResults {
		return results, nil
	}

	// Try common variations in parallel
	variations := []string{
		"python-" + query,
		"py" + query,
		query + "-python",
		query + "lib",
		query + "-py",
	}

	// Filter out already seen variations
	var pendingVariations []string
	for _, v := range variations {
		if !seen[v] {
			pendingVariations = append(pendingVariations, v)
		}
	}

	if len(pendingVariations) == 0 {
		return results, nil
	}

	// Use buffered channel to collect results
	resultChan := make(chan *SearchResult, len(pendingVariations))
	var wg sync.WaitGroup

	// Launch goroutines for each variation
	for _, v := range pendingVariations {
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			if pkg, err := c.GetPackageWithContext(ctx, name); err == nil && pkg != nil {
				resultChan <- pkg
			}
		}(v)
	}

	// Close channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results with early termination
	for pkg := range resultChan {
		name := strings.ToLower(pkg.Name)
		if !seen[name] {
			seen[name] = true
			results = append(results, *pkg)
		}
		if len(results) >= MaxSearchResults {
			break
		}
	}

	return results, nil
}

// isStableVersion checks if a version string represents a stable release.
// Uses regex to match pre-release markers at appropriate positions.
var preReleasePattern = regexp.MustCompile(`(^|[.\-_\d])(a|alpha|b|beta|rc|dev|post)(\d|$|[.\-_])`)

func isStableVersion(v string) bool {
	// First check for obvious pre-release patterns
	lower := strings.ToLower(v)
	for _, pre := range []string{"rc", "dev", "alpha", "beta", "post"} {
		if strings.Contains(lower, pre) {
			// Check if it's a proper pre-release marker (not just a word containing it)
			if preReleasePattern.MatchString(lower) {
				return false
			}
		}
	}
	// Check single letter pre-release markers with boundary
	if strings.Contains(lower, "a") || strings.Contains(lower, "b") {
		if preReleasePattern.MatchString(lower) {
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
	// Strip 'v' prefix if present
	a = strings.TrimPrefix(a, "v")
	b = strings.TrimPrefix(b, "v")

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
