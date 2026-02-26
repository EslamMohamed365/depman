package pip

import (
	"fmt"
	"regexp"
	"strings"
)

// Valid package name pattern: starts with alphanumeric, followed by alphanumerics, dots, hyphens, or underscores
// This matches PEP 508 naming conventions
var packageNamePattern = regexp.MustCompile(`^[a-zA-Z0-9]([a-zA-Z0-9._-]*[a-zA-Z0-9])?$|^[a-zA-Z0-9]$`)

// Invalid characters that could be used for command injection
var invalidCharsPattern = regexp.MustCompile(`[;&|$\` + "`" + `<>]`)

// ValidatePackageName checks if a package name is valid according to PEP 508.
// Returns an error if the name is invalid.
func ValidatePackageName(name string) error {
	if name == "" {
		return fmt.Errorf("package name cannot be empty")
	}

	// Trim whitespace
	name = strings.TrimSpace(name)

	if len(name) > 214 {
		return fmt.Errorf("package name too long (max 214 characters)")
	}

	// Check for invalid characters (command injection prevention)
	if invalidCharsPattern.MatchString(name) {
		return fmt.Errorf("package name contains invalid characters: %s", name)
	}

	// Validate against PEP 508 pattern
	if !packageNamePattern.MatchString(name) {
		return fmt.Errorf("invalid package name: %s (must match PEP 508)", name)
	}

	// Check for common malicious patterns
	lowerName := strings.ToLower(name)
	if strings.Contains(lowerName, "..") ||
		strings.HasPrefix(name, "-") ||
		strings.HasPrefix(name, "_") ||
		strings.HasPrefix(name, ".") {
		return fmt.Errorf("invalid package name: %s", name)
	}

	return nil
}

// SanitizePackageName attempts to sanitize a package name by removing invalid characters.
// Returns the sanitized name or an empty string if the name cannot be sanitized.
func SanitizePackageName(name string) string {
	// Trim whitespace
	name = strings.TrimSpace(name)

	// Remove any shell metacharacters
	sanitized := invalidCharsPattern.ReplaceAllString(name, "")

	// Remove leading hyphens, underscores, or dots
	sanitized = strings.TrimLeft(sanitized, "-_.")

	// Check if sanitized name is valid
	if ValidatePackageName(sanitized) != nil {
		return ""
	}

	return sanitized
}

// ValidateVersionSpec validates a version specifier (e.g., ">=1.0.0", "==2.0.0").
// Returns an error if the specifier is invalid.
func ValidateVersionSpec(spec string) error {
	if spec == "" {
		return nil // Empty is valid (means any version)
	}

	// Valid operators: ==, !=, <=, <, >=, >, ~=
	validOperators := []string{"==", "!=", "<=", "<", ">=", ">", "~="}

	hasValidOperator := false
	for _, op := range validOperators {
		if strings.HasPrefix(spec, op) {
			hasValidOperator = true
			break
		}
	}

	if !hasValidOperator {
		return fmt.Errorf("invalid version specifier: %s", spec)
	}

	return nil
}

// ValidatePackageSpec validates a complete package specification (e.g., "package>=1.0.0").
// Returns an error if the specification is invalid.
func ValidatePackageSpec(spec string) error {
	if spec == "" {
		return fmt.Errorf("package spec cannot be empty")
	}

	// Split on common separators to get package name
	// Handle: package>=1.0.0, package==1.0.0, package~=1.0.0, etc.
	parts := splitPackageSpec(spec)
	if len(parts) == 0 {
		return fmt.Errorf("invalid package spec: %s", spec)
	}

	packageName := parts[0]
	if err := ValidatePackageName(packageName); err != nil {
		return fmt.Errorf("invalid package name in spec: %w", err)
	}

	// Validate version spec if present
	if len(parts) > 1 {
		if err := ValidateVersionSpec(parts[1]); err != nil {
			return err
		}
	}

	return nil
}

// splitPackageSpec splits a package spec into name and version parts.
func splitPackageSpec(spec string) []string {
	// Try common operators
	operators := []string{"==", "!=", "<=", "<", ">=", ">", "~=", "@"}

	for _, op := range operators {
		if idx := strings.Index(spec, op); idx > 0 {
			return []string{spec[:idx], spec[idx:]}
		}
	}

	// No operator found, treat entire string as name
	return []string{spec}
}
