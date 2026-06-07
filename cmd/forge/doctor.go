package forge

import (
	"fmt"

	"github.com/AARCSX/AARCSX_Forge/internal/cli/doctor"
	"github.com/spf13/cobra"
)

// NewDoctorCommand creates the doctor command
func NewDoctorCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check development environment",
		Long: `Perform diagnostic checks on the local development environment
to verify all required tools and services are available and configured correctly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDoctor(cmd)
		},
	}
}

func runDoctor(cmd *cobra.Command) error {
	// Create checks
	checks := []doctor.Check{
		doctor.NewGoCheck(),
		doctor.NewDockerCheck(),
		doctor.NewPostgreSQLCheck(),
		doctor.NewRedisCheck(),
		doctor.NewForgeMetadataCheck(),
	}

	// Create runner and run checks
	ctx := cmd.Context()
	runner := doctor.NewRunner(checks...)
	results := runner.Run(ctx)

	// Print results
	allPassed := true
	for _, result := range results {
		if result.Passed {
			fmt.Printf("✓ %s\n", result.Name)
		} else {
			fmt.Printf("✗ %s: %s\n", result.Name, result.Message)
			allPassed = false
		}
	}

	if !allPassed {
		return fmt.Errorf("some checks failed")
	}

	return nil
}