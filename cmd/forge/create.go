package forge

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/AARCSX/AARCSX_Forge/internal/cli/scaffold"
	"github.com/AARCSX/AARCSX_Forge/pkg/version"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// NewCreateCommand creates the create command
func NewCreateCommand() *cobra.Command {
	var postgresURL string
	var redisURL string
	var storageProvider string
	var skipPrompts bool

	cmd := &cobra.Command{
		Use:   "create [project-name]",
		Short: "Create a new Forge project",
		Long: `Create a new Forge project from versioned templates.
The command will prompt for project details and scaffold a new project
that conforms to Forge architecture standards.`,
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			var projectName string
			if len(args) == 1 {
				projectName = args[0]
			}
			return runCreate(cmd.Context(), projectName, postgresURL, redisURL, storageProvider, skipPrompts)
		},
	}

	// Flags
	cmd.Flags().StringVar(&postgresURL, "postgres-url", "", "PostgreSQL connection URL")
	cmd.Flags().StringVar(&redisURL, "redis-url", "", "Redis connection URL")
	cmd.Flags().StringVar(&storageProvider, "storage-provider", "", "Storage provider: minio or s3")
	cmd.Flags().BoolVar(&skipPrompts, "skip-prompts", false, "Skip interactive prompts (requires all values via flags)")

	return cmd
}

func runCreate(_ context.Context, projectName, postgresURL, redisURL, storageProvider string, skipPrompts bool) error {
	// Validate project name if not provided via args
	if projectName == "" && !skipPrompts {
		// Will prompt for it
	} else if projectName == "" {
		return fmt.Errorf("project name is required")
	}

	// Validate project name
	if projectName != "" {
		if err := validateProjectName(projectName); err != nil {
			return fmt.Errorf("invalid project name: %w", err)
		}
	}

	// Interactive prompts if needed
	if !skipPrompts {
		var err error
		projectName, postgresURL, redisURL, storageProvider, err = promptForProjectSettings(projectName, postgresURL, redisURL, storageProvider)
		if err != nil {
			return fmt.Errorf("failed to prompt for project details: %w", err)
		}
	}

	// Validate all required fields are present
	if projectName == "" {
		return fmt.Errorf("project name is required")
	}
	if storageProvider == "" {
		return fmt.Errorf("storage provider is required")
	}

	// Validate storage provider
	if storageProvider != "minio" && storageProvider != "s3" {
		return fmt.Errorf("storage provider must be 'minio' or 's3'")
	}

	// Create project directory
	projectDir := filepath.Clean(projectName)
	if _, err := os.Stat(projectDir); err == nil {
		return fmt.Errorf("directory %s already exists", projectDir)
	}

	if err := os.MkdirAll(projectDir, 0755); err != nil {
		return fmt.Errorf("create project directory: %w", err)
	}

	// Download and extract template
	templatePath, err := downloadTemplate(version.RuntimeVersion)
	if err != nil {
		return fmt.Errorf("download template: %w", err)
	}
	defer os.RemoveAll(templatePath) // Clean up extracted template

	// Copy template to project directory
	if err := copyDirectory(templatePath, projectDir); err != nil {
		return fmt.Errorf("copy template: %w", err)
	}

	// Replace placeholders in files
	if err := scaffold.ReplacePlaceholders(projectDir, map[string]string{
		"{{PROJECT_NAME}}":  projectName,
		"{{MODULE_NAME}}":   toModuleName(projectName),
		"{{FORGE_VERSION}}": version.RuntimeVersion,
	}); err != nil {
		return fmt.Errorf("replace placeholders: %w", err)
	}

	// Generate supporting files
	if err := generateEnvExample(projectDir, postgresURL, redisURL, storageProvider); err != nil {
		return fmt.Errorf("generate .env.example: %w", err)
	}

	if err := generateReadme(projectDir, projectName); err != nil {
		return fmt.Errorf("generate README.md: %w", err)
	}

	if err := generateProjectMetadata(projectDir, projectName, version.RuntimeVersion); err != nil {
		return fmt.Errorf("generate project metadata: %w", err)
	}

	// Print success message
	printSuccessMessage(projectDir, projectName)

	return nil
}

func validateProjectName(name string) error {
	if name == "" {
		return fmt.Errorf("project name cannot be empty")
	}
	// Basic validation - alphanumeric, hyphens, underscores allowed
	if !regexp.MustCompile(`^[a-zA-Z0-9_-]+$`).MatchString(name) {
		return fmt.Errorf("project name can only contain letters, numbers, hyphens, and underscores")
	}
	return nil
}

func promptForProjectSettings(projectName, postgresURL, redisURL, storageProvider string) (string, string, string, string, error) {
	reader := bufio.NewReader(os.Stdin)

	if projectName == "" {
		fmt.Print("Project name: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", "", "", "", fmt.Errorf("read project name: %w", err)
		}
		projectName = strings.TrimSpace(input)
	}

	if postgresURL == "" {
		fmt.Print("PostgreSQL URL (optional): ")
		input, err := reader.ReadString('\n')
		if err == nil {
			postgresURL = strings.TrimSpace(input)
		}
	}

	if redisURL == "" {
		fmt.Print("Redis URL (optional): ")
		input, err := reader.ReadString('\n')
		if err == nil {
			redisURL = strings.TrimSpace(input)
		}
	}

	if storageProvider == "" {
		fmt.Print("Storage provider (minio/s3) [minio]: ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", "", "", "", fmt.Errorf("read storage provider: %w", err)
		}
		storageProvider = strings.TrimSpace(input)
		if storageProvider == "" {
			storageProvider = "minio"
		}
	}

	return projectName, postgresURL, redisURL, storageProvider, nil
}

// toModuleName converts project name to a suitable module name (snake_case)
func toModuleName(projectName string) string {
	return strings.ToLower(projectName)
}

func downloadTemplate(version string) (string, error) {
	// Create temp directory for template
	tempDir, err := os.MkdirTemp("", "forge-template-*")
	if err != nil {
		return "", fmt.Errorf("create temp directory: %w", err)
	}

	templateURL, err := resolveTemplateURL(version)
	if err != nil {
		return "", err
	}

	resp, err := http.Get(templateURL)
	if err != nil {
		return "", fmt.Errorf("download template from %s: %w", templateURL, err)
	}
	defer resp.Body.Close()

	// Save to file
	templatePath := filepath.Join(tempDir, "template.tar.gz")
	out, err := os.Create(templatePath)
	if err != nil {
		return "", fmt.Errorf("create template file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("save template: %w", err)
	}

	// Extract template
	extractedPath := filepath.Join(tempDir, "extracted")
	if err := os.MkdirAll(extractedPath, 0755); err != nil {
		return "", fmt.Errorf("create extract directory: %w", err)
	}

	// Extract tar.gz
	cmd := exec.Command("tar", "-xzf", templatePath, "-C", extractedPath)
	if output, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("extract template: %w, output: %s", err, string(output))
	}

	// Find the extracted directory
	entries, err := os.ReadDir(extractedPath)
	if err != nil {
		return "", fmt.Errorf("read extracted directory: %w", err)
	}
	if len(entries) == 0 {
		return "", fmt.Errorf("no content found in extracted template")
	}

	// If the archive contains a single folder wrapping everything, return that.
	// Otherwise, return the extracted path root itself.
	if len(entries) == 1 && entries[0].IsDir() {
		return filepath.Join(extractedPath, entries[0].Name()), nil
	}
	return extractedPath, nil
}

func resolveTemplateURL(requestedVersion string) (string, error) {
	client := &http.Client{Timeout: 15 * time.Second}
	attemptedURLs := buildTemplateCandidateURLs(requestedVersion)

	for _, templateURL := range attemptedURLs {
		resp, err := client.Head(templateURL)
		if err != nil {
			continue
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return templateURL, nil
		}
	}

	return "", fmt.Errorf(
		"failed to resolve community template URL. expected a release asset named forge-template-<version>.tar.gz in AARCSX/AARCSX_Forge. tried: %s",
		strings.Join(attemptedURLs, ", "),
	)
}

func buildTemplateCandidateURLs(requestedVersion string) []string {
	candidateVersions := make([]string, 0, 2)
	for _, candidate := range []string{requestedVersion, version.CLIVersion} {
		if candidate != "" && !slices.Contains(candidateVersions, candidate) {
			candidateVersions = append(candidateVersions, candidate)
		}
	}

	urls := make([]string, 0, len(candidateVersions)*2)
	for _, candidateVersion := range candidateVersions {
		for _, tag := range []string{candidateVersion, "v" + candidateVersion} {
			urls = append(urls, fmt.Sprintf(
				"https://github.com/AARCSX/AARCSX_Forge/releases/download/%s/forge-template-%s.tar.gz",
				tag,
				candidateVersion,
			))
		}
	}

	return urls
}

func copyDirectory(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("read source directory: %w", err)
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			if err := os.MkdirAll(dstPath, entry.Type().Perm()); err != nil {
				return fmt.Errorf("create directory %s: %w", dstPath, err)
			}
			if err := copyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			if err := copyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("copy file %s: %w", srcPath, err)
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("open source file: %w", err)
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)
	return err
}

func generateEnvExample(projectDir, postgresURL, redisURL, storageProvider string) error {
	envExample := fmt.Sprintf(`# Forge Environment Configuration
# Copy this file to .env and fill in the values

# Database Configuration
POSTGRES_URL=%s
REDIS_URL=%s

# Security
JWT_SIGNING_KEY=your-secret-key-here-change-in-production

# Storage
STORAGE_PROVIDER=%s
`, postgresURL, redisURL, storageProvider)

	switch storageProvider {
	case "minio":
		envExample += `
# MinIO Configuration (required if STORAGE_PROVIDER=minio)
MINIO_ENDPOINT=localhost:9000
MINIO_ACCESS_KEY=minioadmin
MINIO_SECRET_KEY=minioadmin
MINIO_BUCKET=forge-bucket
`
	case "s3":
		envExample += `
# S3 Configuration (required if STORAGE_PROVIDER=s3)
S3_ENDPOINT=
S3_REGION=us-east-1
S3_BUCKET=
`
	}

	envExample += `
# Application
APP_ENV=local
HTTP_ADDR=:8080
`

	path := filepath.Join(projectDir, ".env.example")
	return os.WriteFile(path, []byte(envExample), 0644)
}

func generateReadme(projectDir, projectName string) error {
	readme := fmt.Sprintf(`# %s
A Forge-powered community project.

## Getting Started

1. Copy .env.example to .env and configure your environment variables
2. Install dependencies: go mod tidy
3. Apply database migrations: migrate -path ./migrations -database "${POSTGRES_URL}" up
4. Run the application: go run ./cmd/api/main.go
5. In another terminal, run the worker: go run ./cmd/worker/main.go

## Project Structure

This project follows the Forge architecture standard:

- cmd/ - Application entry points (API server and worker)
- internal/ - Domain-specific code and shared libraries
- migrations/ - Database migrations
- pkg/ - Shared packages

## Documentation

For more information about Forge architecture and development practices,
see the Forge documentation.

## License

This project is open source software.
`, projectName)

	path := filepath.Join(projectDir, "README.md")
	return os.WriteFile(path, []byte(readme), 0644)
}

func generateProjectMetadata(projectDir, projectName, forgeVersion string) error {
	metadata := map[string]interface{}{
		"project_name":  projectName,
		"forge_version": forgeVersion,
		"created_at":    time.Now().UTC().Format(time.RFC3339),
	}

	forgeDir := filepath.Join(projectDir, ".forge")
	if err := os.MkdirAll(forgeDir, 0755); err != nil {
		return fmt.Errorf("create .forge directory: %w", err)
	}

	yamlData, err := yaml.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("marshal project metadata: %w", err)
	}

	path := filepath.Join(forgeDir, "project.yaml")
	return os.WriteFile(path, yamlData, 0644)
}

func printSuccessMessage(projectDir, projectName string) {
	fmt.Printf("\n🎉 Successfully created Forge project '%s'!\n\n", projectName)
	fmt.Printf("Project details:\n")
	fmt.Printf("  • Name: %s\n", projectName)
	fmt.Printf("  • Location: %s\n", projectDir)
	fmt.Printf("\nNext steps:\n")
	fmt.Printf("  1. cd %s\n", projectName)
	fmt.Printf("  2. cp .env.example .env && edit .env to configure your environment\n")
	fmt.Printf("  3. go mod tidy\n")
	fmt.Printf("  4. migrate -path ./migrations -database \"$POSTGRES_URL\" up\n")
	fmt.Printf("  5. go run ./cmd/api/main.go\n")
	fmt.Printf("  6. In another terminal: go run ./cmd/worker/main.go\n")
	fmt.Printf("\nHappy coding! 🚀\n")
}
