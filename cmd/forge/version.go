package forge

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/AARCSX/AARCSX_Forge/pkg/version"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"gopkg.in/yaml.v3"
)

// NewVersionCommand creates the version command
func NewVersionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long: `Print Forge CLI version and Forge runtime version.
If run within a Forge project, also shows project metadata.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runVersion(cmd)
		},
	}
}

func runVersion(_ *cobra.Command) error {
	// Print CLI and runtime versions
	fmt.Printf("Forge CLI Version: %s\n", version.CLIVersion)
	fmt.Printf("Forge Runtime Version: %s\n", version.RuntimeVersion)

	// Check if we're in a Forge project
	if metadata, err := readProjectMetadata(); err == nil {
		fmt.Println("\nCurrent Project:")
		fmt.Printf("  Name: %s\n", metadata.ProjectName)
		fmt.Printf("  Forge Runtime Version: %s\n", metadata.ForgeVersion)
		fmt.Printf("  Created At: %s\n", metadata.CreatedAt.Format("2006-01-02T15:04:05Z"))
	} else if !os.IsNotExist(err) {
		// Only log error if it's not just "file not found"
		log.Warn("unable to read project metadata", zap.Error(err))
	}

	return nil
}

// ProjectMetadata represents the .forge/project.yaml file
type ProjectMetadata struct {
	ProjectName  string    `yaml:"project_name"`
	ForgeVersion string    `yaml:"forge_version"`
	CreatedAt    time.Time `yaml:"created_at"`
}

func readProjectMetadata() (*ProjectMetadata, error) {
	// Look for .forge/project.yaml in current directory or parent directories
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("get current directory: %w", err)
	}

	for {
		path := filepath.Join(dir, ".forge", "project.yaml")
		if _, err := os.Stat(path); err == nil {
			// File exists, read it
			data, err := os.ReadFile(path)
			if err != nil {
				return nil, fmt.Errorf("read project metadata: %w", err)
			}

			var metadata ProjectMetadata
			if err := yaml.Unmarshal(data, &metadata); err != nil {
				return nil, fmt.Errorf("unmarshal project metadata: %w", err)
			}

			return &metadata, nil
		}

		// Check if we've reached the root directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return nil, os.ErrNotExist
}
