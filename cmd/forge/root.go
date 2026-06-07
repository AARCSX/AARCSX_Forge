package forge

import (
	"context"
	"os"
	"os/signal"
	"syscall"

    "github.com/AARCSX/AARCSX_Forge/internal/logger"
    "github.com/spf13/cobra"
)

var log *logger.Logger

// NewForgeCommand creates and returns the Forge CLI root command
func NewForgeCommand() *cobra.Command {
	var cfgFile string

	rootCmd := &cobra.Command{
		Use:   "forge",
		Short: "Forge CLI - AARCSX Engineering Platform",
		Long: `Forge CLI is a tool for creating and managing projects
that conform to the AARCSX Forge architecture standards.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// Load configuration
			if cfgFile != "" {
				// TODO: support config file loading if needed
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no subcommand is provided, show help
			return cmd.Help()
		},
	}

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.forge.yaml)")

	// Add subcommands
	rootCmd.AddCommand(NewVersionCommand())
	rootCmd.AddCommand(NewCreateCommand())
	rootCmd.AddCommand(NewDoctorCommand())

	return rootCmd
}

// Execute runs the Forge CLI command
func Execute() {
    var err error
    log, err = logger.New("info", "forge-cli")
    if err != nil {
        _, _ = os.Stderr.WriteString("failed to initialize logger: " + err.Error() + "\n")
        os.Exit(1)
    }

    rootCmd := NewForgeCommand()

    // Set up context for graceful shutdown
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

	// Handle OS signals for graceful termination
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	// Execute command with context
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		log.Error(ctx, "command execution failed", err, nil)
		os.Exit(1)
	}
}