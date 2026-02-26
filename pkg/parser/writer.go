package parser

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/eslam/depman/pkg/detector"
	"github.com/eslam/depman/pkg/pip"
)

// WriteDependencyFile performs an atomic full rewrite of the dependency file
// based on the currently installed packages.
func WriteDependencyFile(project detector.Project, packages []pip.Package) error {
	deps := packagesToDeps(packages)

	// Sort alphabetically
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	var content string

	switch project.FileType {
	case detector.FileRequirementsTXT:
		content = FormatRequirementsTxt(deps)

	case detector.FilePyprojectTOML:
		// Read existing file to preserve non-dependency sections
		existing, err := os.ReadFile(project.FilePath)
		if err != nil {
			return fmt.Errorf("reading %s: %w", project.FilePath, err)
		}
		content = RewritePyprojectDependencies(string(existing), deps)

	default:
		return fmt.Errorf("unknown file type: %v", project.FileType)
	}

	return atomicWrite(project.FilePath, []byte(content))
}

// atomicWrite writes content to a temp file then renames to the target path.
func atomicWrite(target string, content []byte) error {
	dir := filepath.Dir(target)
	tmp, err := os.CreateTemp(dir, ".depman-*")
	if err != nil {
		return fmt.Errorf("creating temp file: %w", err)
	}
	tmpPath := tmp.Name()

	// Ensure cleanup on failure
	defer func() {
		if err != nil {
			os.Remove(tmpPath)
		}
	}()

	if _, err = tmp.Write(content); err != nil {
		tmp.Close()
		return fmt.Errorf("writing temp file: %w", err)
	}
	if err = tmp.Close(); err != nil {
		return fmt.Errorf("closing temp file: %w", err)
	}

	// Preserve original file permissions
	if info, statErr := os.Stat(target); statErr == nil {
		os.Chmod(tmpPath, info.Mode())
	}

	if err = os.Rename(tmpPath, target); err != nil {
		return fmt.Errorf("renaming temp file to %s: %w", target, err)
	}

	return nil
}

func packagesToDeps(packages []pip.Package) []Dep {
	deps := make([]Dep, len(packages))
	for i, p := range packages {
		deps[i] = Dep{
			Name:    p.Name,
			Version: p.InstalledVersion,
		}
	}
	return deps
}
