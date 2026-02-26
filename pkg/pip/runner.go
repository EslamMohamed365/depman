package pip

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/eslam/depman/pkg/env"
	"github.com/eslam/depman/pkg/log"
)

// Runner executes pip/uv commands scoped to a specific environment.
type Runner struct {
	Manager env.PackageManager
	Venv    env.Virtualenv
}

// NewRunner creates a runner for the given manager and virtualenv.
func NewRunner(mgr env.PackageManager, venv env.Virtualenv) *Runner {
	return &Runner{Manager: mgr, Venv: venv}
}

// RunResult holds the result of a pip/uv command.
type RunResult struct {
	Stdout string
	Stderr string
	Err    error
}

// Run executes a pip/uv command with the resolved environment.
func (r *Runner) Run(bin string, args ...string) RunResult {
	log.Info("executing package manager command", "bin", bin, "args", strings.Join(args, " "))
	cmd := exec.Command(bin, args...)
	cmd.Env = r.buildEnv()

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Warn("package manager command failed", "bin", bin, "error", err, "stderr", stderr.String())
	} else {
		log.Info("package manager command succeeded", "bin", bin)
	}
	return RunResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
}

// Install installs a package.
func (r *Runner) Install(pkg string) RunResult {
	bin, args := r.Manager.InstallCmd(pkg)
	return r.Run(bin, args...)
}

// Uninstall removes a package.
func (r *Runner) Uninstall(pkg string) RunResult {
	bin, args := r.Manager.UninstallCmd(pkg)
	return r.Run(bin, args...)
}

// Upgrade upgrades a package to its latest version.
func (r *Runner) Upgrade(pkg string) RunResult {
	bin, args := r.Manager.UpgradeCmd(pkg)
	return r.Run(bin, args...)
}

// List returns the raw JSON output of installed packages.
func (r *Runner) List() RunResult {
	bin, args := r.Manager.ListCmd()
	return r.Run(bin, args...)
}

// Outdated returns the raw JSON output of outdated packages.
func (r *Runner) Outdated() RunResult {
	bin, args := r.Manager.OutdatedCmd()
	return r.Run(bin, args...)
}

// buildEnv constructs the environment variables for subprocess calls.
// It sets VIRTUAL_ENV and prepends the venv bin directory to PATH.
func (r *Runner) buildEnv() []string {
	environ := os.Environ()

	if r.Venv.Type == env.EnvVirtualenv && r.Venv.Path != "" {
		venvBin := filepath.Join(r.Venv.Path, "bin")
		environ = setEnv(environ, "VIRTUAL_ENV", r.Venv.Path)
		environ = prependPath(environ, venvBin)
	}

	return environ
}

func setEnv(environ []string, key, value string) []string {
	prefix := key + "="
	for i, e := range environ {
		if len(e) > len(prefix) && e[:len(prefix)] == prefix {
			environ[i] = fmt.Sprintf("%s=%s", key, value)
			return environ
		}
	}
	return append(environ, fmt.Sprintf("%s=%s", key, value))
}

func prependPath(environ []string, dir string) []string {
	prefix := "PATH="
	for i, e := range environ {
		if len(e) > len(prefix) && e[:len(prefix)] == prefix {
			environ[i] = fmt.Sprintf("PATH=%s:%s", dir, e[len(prefix):])
			return environ
		}
	}
	return append(environ, fmt.Sprintf("PATH=%s", dir))
}
