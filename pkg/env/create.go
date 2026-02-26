package env

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// CreateVirtualenv creates a new .venv in the given project directory.
// It prefers `uv venv` if available, falling back to `python -m venv`.
func CreateVirtualenv(dir string) (Virtualenv, error) {
	venvPath := filepath.Join(dir, ".venv")

	// Try uv first
	if uvPath, err := exec.LookPath("uv"); err == nil {
		cmd := exec.Command(uvPath, "venv", venvPath)
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			return Virtualenv{}, fmt.Errorf("uv venv failed: %s: %w", string(output), err)
		}
	} else {
		// Fall back to python -m venv
		pythonBin := findPython()
		if pythonBin == "" {
			return Virtualenv{}, fmt.Errorf("neither uv nor python found")
		}
		cmd := exec.Command(pythonBin, "-m", "venv", venvPath)
		cmd.Dir = dir
		if output, err := cmd.CombinedOutput(); err != nil {
			return Virtualenv{}, fmt.Errorf("python -m venv failed: %s: %w", string(output), err)
		}
	}

	pythonBin := filepath.Join(venvPath, "bin", "python")
	if !fileExecutable(pythonBin) {
		return Virtualenv{}, fmt.Errorf("created .venv but python binary not found at %s", pythonBin)
	}

	return Virtualenv{
		Type:      EnvVirtualenv,
		Path:      venvPath,
		PythonBin: pythonBin,
	}, nil
}

// RecreateVirtualenv removes and recreates the virtualenv.
func RecreateVirtualenv(dir string, venvPath string) (Virtualenv, error) {
	if err := os.RemoveAll(venvPath); err != nil {
		return Virtualenv{}, fmt.Errorf("removing broken venv: %w", err)
	}
	return CreateVirtualenv(dir)
}

func findPython() string {
	for _, name := range []string{"python3", "python"} {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}
