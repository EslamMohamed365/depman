package pypi

import "time"

// HTTP client configuration constants
const (
	// HTTPTimeout is the default timeout for HTTP requests
	HTTPTimeout = 10 * time.Second

	// MaxDisplayVersions is the maximum number of versions to display in package details
	MaxDisplayVersions = 20

	// MaxLicenseLength is the maximum length of license text before truncation
	MaxLicenseLength = 40
)

// HTTP retry configuration constants
const (
	// MaxRetries is the maximum number of retry attempts for failed requests
	MaxRetries = 3

	// RetryDelay is the initial delay before the first retry
	RetryDelay = 100 * time.Millisecond

	// BackoffMultiplier is the exponential backoff multiplier for retries
	BackoffMultiplier = 2.0
)

// HTTP status code constants
const (
	// StatusOKMin is the minimum HTTP status code for successful responses
	StatusOKMin = 200

	// StatusOKMax is the maximum HTTP status code for successful responses
	StatusOKMax = 300

	// StatusNotFound represents HTTP 404 Not Found
	StatusNotFound = 404

	// StatusOK represents HTTP 200 OK
	StatusOK = 200

	// StatusInternalServerError represents HTTP 500 Internal Server Error
	StatusInternalServerError = 500

	// StatusBadGateway represents HTTP 502 Bad Gateway
	StatusBadGateway = 502

	// StatusServiceUnavailable represents HTTP 503 Service Unavailable
	StatusServiceUnavailable = 503

	// StatusGatewayTimeout represents HTTP 504 Gateway Timeout
	StatusGatewayTimeout = 504
)

// Search configuration constants
const (
	// MaxSearchResults is the maximum number of search results to return
	MaxSearchResults = 10
)
