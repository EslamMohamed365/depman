package detector

import (
	"os"
	"path/filepath"
)

// FileType represents the type of dependency file found.
type FileType int

const (
	FileNone FileType = iota
	FilePyprojectTOML
	FileRequirementsTXT
)

// String returns the human-readable name of the file type.
func (f FileType) String() string {
	switch f {
	case FilePyprojectTOML:
		return "pyproject.toml"
	case FileRequirementsTXT:
		return "requirements.txt"
	default:
		return "none"
	}
}

// Project holds information about a detected Python project.
type Project struct {
	FilePath string   // Absolute path to the dependency file
	FileType FileType // Type of file detected
	Dir      string   // Project root directory
}

// Detected returns true if a project file was found.
func (p Project) Detected() bool {
	return p.FileType != FileNone
}

// DetectProject scans the given directory for Python dependency files.
// Detection priority: pyproject.toml → requirements.txt → requirements/*.txt
func DetectProject(dir string) Project {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return Project{Dir: dir}
	}

	// 1. Check pyproject.toml
	pyproject := filepath.Join(absDir, "pyproject.toml")
	if fileExists(pyproject) {
		return Project{
			FilePath: pyproject,
			FileType: FilePyprojectTOML,
			Dir:      absDir,
		}
	}

	// 2. Check requirements.txt
	reqtxt := filepath.Join(absDir, "requirements.txt")
	if fileExists(reqtxt) {
		return Project{
			FilePath: reqtxt,
			FileType: FileRequirementsTXT,
			Dir:      absDir,
		}
	}

	// 3. Check requirements/*.txt (use first found)
	reqDir := filepath.Join(absDir, "requirements")
	if dirExists(reqDir) {
		entries, err := os.ReadDir(reqDir)
		if err == nil {
			for _, e := range entries {
				if !e.IsDir() && filepath.Ext(e.Name()) == ".txt" {
					return Project{
						FilePath: filepath.Join(reqDir, e.Name()),
						FileType: FileRequirementsTXT,
						Dir:      absDir,
					}
				}
			}
		}
	}

	return Project{Dir: absDir}
}

func fileExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && !info.IsDir()
}

func dirExists(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}
