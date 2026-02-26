package parser

import (
	"github.com/eslam/depman/pkg/detector"
	"github.com/eslam/depman/pkg/pip"
)

// SyncDependencyFile runs the full sync cycle after a package operation:
// 1. Query the full list of installed packages
// 2. Rewrite the dependency file from scratch
func SyncDependencyFile(project detector.Project, runner *pip.Runner) error {
	listResult := runner.List()
	if listResult.Err != nil {
		return listResult.Err
	}

	packages, err := pip.ParsePackageList(listResult.Stdout)
	if err != nil {
		return err
	}

	return WriteDependencyFile(project, packages)
}
