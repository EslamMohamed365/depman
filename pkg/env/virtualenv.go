package env

import (
	"os"
	"os/exec"
	"path/filepath"
)

// EnvType represents the type of Python environment.
type EnvType int

const (
	EnvNotFound EnvType = iota
	EnvVirtualenv
	EnvSystem
)
const (
	FilePermExecutable = 0111
)

// Virtualenv holds information about the detected Python environment.
type Virtualenv struct {
	Type      EnvType
	Path      string // Path to the virtualenv or system Python dir
	PythonBin string // Path to the python binary
	IsActive  bool   // True if from $VIRTUAL_ENV
	IsBroken  bool   // True if venv exists but interpreter is missing
}

// Name returns a short display name for the environment.
func (v Virtualenv) Name() string {
	switch v.Type {
	case EnvVirtualenv:
		return filepath.Base(v.Path)
	case EnvSystem:
		return "system"
	default:
		return "none"
	}
}

// DetectVirtualenv resolves the active Python environment.
// Priority: $VIRTUAL_ENV → .venv/ → venv/ → system Python.
func DetectVirtualenv(dir string) Virtualenv {
	absDir, _ := filepath.Abs(dir)

	// 1. Check $VIRTUAL_ENV
	if venvPath := os.Getenv("VIRTUAL_ENV"); venvPath != "" {
		pythonBin := filepath.Join(venvPath, "bin", "python")
		if fileExecutable(pythonBin) {
			return Virtualenv{
				Type:      EnvVirtualenv,
				Path:      venvPath,
				PythonBin: pythonBin,
				IsActive:  true,
			}
		}
		// $VIRTUAL_ENV is set but broken
		return Virtualenv{
			Type:     EnvVirtualenv,
			Path:     venvPath,
			IsActive: true,
			IsBroken: true,
		}
	}

	// 2. Check .venv/
	dotVenv := filepath.Join(absDir, ".venv")
	if v := checkLocalVenv(dotVenv); v.Type != EnvNotFound {
		return v
	}

	// 3. Check venv/
	venv := filepath.Join(absDir, "venv")
	if v := checkLocalVenv(venv); v.Type != EnvNotFound {
		return v
	}

	// 4. Fall back to system Python
	return detectSystemPython()
}

func checkLocalVenv(path string) Virtualenv {
	pythonBin := filepath.Join(path, "bin", "python")
	info, err := os.Stat(path)
	if err != nil || !info.IsDir() {
		return Virtualenv{} // doesn't exist
	}

	if fileExecutable(pythonBin) {
		return Virtualenv{
			Type:      EnvVirtualenv,
			Path:      path,
			PythonBin: pythonBin,
		}
	}

	// Directory exists but interpreter is missing — broken
	return Virtualenv{
		Type:     EnvVirtualenv,
		Path:     path,
		IsBroken: true,
	}
}

func detectSystemPython() Virtualenv {
	// Try python3 first, then python
	for _, name := range []string{"python3", "python"} {
		path, err := exec.LookPath(name)
		if err == nil {
			return Virtualenv{
				Type:      EnvSystem,
				Path:      filepath.Dir(path),
				PythonBin: path,
			}
		}
	}
	return Virtualenv{Type: EnvNotFound}
}

func fileExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&FilePermExecutable != 0
}
