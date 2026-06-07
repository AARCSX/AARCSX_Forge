package runtime

import (
	"fmt"

	"github.com/AARCSX/AARCSX_Forge/pkg/version"
)

// VersionManager manages version information for the CLI and runtime
type VersionManager struct {
	cliVersion      string
	runtimeVersion  string
}

// NewVersionManager creates a new version manager
func NewVersionManager() *VersionManager {
	return &VersionManager{
		cliVersion:     version.CLIVersion,
		runtimeVersion: version.RuntimeVersion,
	}
}

// CLIVersion returns the CLI version
func (vm *VersionManager) CLIVersion() string {
	return vm.cliVersion
}

// RuntimeVersion returns the Forge runtime version
func (vm *VersionManager) RuntimeVersion() string {
	return vm.runtimeVersion
}

// String returns formatted version information
func (vm *VersionManager) String() string {
	return fmt.Sprintf("Forge CLI Version: %s\nForge Runtime Version: %s", vm.cliVersion, vm.runtimeVersion)
}

// CheckCompatibility validates that the CLI version is compatible with the runtime version
// For now, we assume same major/minor versions are compatible
func (vm *VersionManager) CheckCompatibility() error {
	// In a real implementation, this would do proper semantic version comparison
	// For Sprint 6, we'll just check that both versions are set
	if vm.cliVersion == "" || vm.runtimeVersion == "" {
		return fmt.Errorf("version information incomplete")
	}
	return nil
}