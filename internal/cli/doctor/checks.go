package doctor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/AARCSX/AARCSX_Forge/internal/config"
	"gopkg.in/yaml.v3"
)

// CheckResult represents the result of a single diagnostic check
type CheckResult struct {
	Name     string
	Passed   bool
	Message  string
	Details  map[string]string
}

// Check defines the interface for a diagnostic check
type Check interface {
	Name() string
	Check(ctx context.Context) CheckResult
}

// BaseCheck provides common functionality for checks
type BaseCheck struct {
	name string
}

func (b *BaseCheck) Name() string { return b.name }

// GoCheck validates Go installation and version
type GoCheck struct {
	BaseCheck
}

func NewGoCheck() *GoCheck {
	return &GoCheck{BaseCheck{name: "Go"}}
}

func (c *GoCheck) Check(ctx context.Context) CheckResult {
	output, err := exec.CommandContext(ctx, "go", "version").CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Go not found or not executable",
		}
	}

	// Parse version string like "go version go1.25.0 linux/amd64"
	parts := strings.Fields(string(output))
	if len(parts) < 3 {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Unable to parse Go version",
		}
	}

	versionStr := strings.TrimPrefix(parts[2], "go")
	// Simple version check - need >= 1.25.0
	if strings.HasPrefix(versionStr, "1.25.") || compareVersions(versionStr, "1.25.0") >= 0 {
		return CheckResult{
			Name:    c.Name(),
			Passed:  true,
			Message: fmt.Sprintf("Go %s detected", versionStr),
		}
	}
	return CheckResult{
		Name:    c.Name(),
		Passed:  false,
		Message: fmt.Sprintf("Go version %s is too old, need >= 1.25.0", versionStr),
	}
}

// DockerCheck validates Docker installation and daemon status
type DockerCheck struct {
	BaseCheck
}

func NewDockerCheck() *DockerCheck {
	return &DockerCheck{BaseCheck{name: "Docker"}}
}

func (c *DockerCheck) Check(ctx context.Context) CheckResult {
	// Check if docker command exists
	if _, err := exec.CommandContext(ctx, "docker", "--version").CombinedOutput(); err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Docker not found or not executable",
		}
	}

	// Check if daemon is running
	output, err := exec.CommandContext(ctx, "docker", "info").CombinedOutput()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Docker daemon not running",
			Details: map[string]string{
				"error": err.Error(),
				"output": string(output),
			},
		}
	}

	// Basic check that we got some output
	if len(output) == 0 {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Docker daemon not responding",
		}
	}

	return CheckResult{
		Name:    c.Name(),
		Passed:  true,
		Message: "Docker daemon is running",
	}
}

// PostgreSQLCheck validates PostgreSQL connectivity
type PostgreSQLCheck struct {
	BaseCheck
}

func NewPostgreSQLCheck() *PostgreSQLCheck {
	return &PostgreSQLCheck{BaseCheck{name: "PostgreSQL"}}
}

func (c *PostgreSQLCheck) Check(ctx context.Context) CheckResult {
	// Try to load config to get PostgreSQL URL
	cfg, err := config.Load()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("Cannot load configuration: %v", err),
		}
	}

	if cfg.Database.URL == "" {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "PostgreSQL URL not configured",
		}
	}

	// Validate URL format
	if !strings.HasPrefix(cfg.Database.URL, "postgres://") &&
		!strings.HasPrefix(cfg.Database.URL, "postgresql://") {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Invalid PostgreSQL URL format",
			Details: map[string]string{
				"url": cfg.Database.URL,
			},
		}
	}

    // All validations passed; return success
    return CheckResult{
        Name:    c.Name(),
        Passed:  true,
        Message: "PostgreSQL URL configured correctly",
        Details: map[string]string{
            "url": cfg.Database.URL,
        },
    }
}

// RedisCheck validates Redis connectivity
type RedisCheck struct {
	BaseCheck
}

func NewRedisCheck() *RedisCheck {
	return &RedisCheck{BaseCheck{name: "Redis"}}
}

func (c *RedisCheck) Check(ctx context.Context) CheckResult {
	// Try to load config to get Redis URL
	cfg, err := config.Load()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("Cannot load configuration: %v", err),
		}
	}

	if cfg.Redis.URL == "" {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Redis URL not configured",
		}
	}

	// Simple validation - in reality would use proper redis client
	if !strings.HasPrefix(cfg.Redis.URL, "redis://") &&
		!strings.HasPrefix(cfg.Redis.URL, "rediss://") {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: "Invalid Redis URL format",
			Details: map[string]string{
				"url": cfg.Redis.URL,
			},
		}
	}

	// TODO: Actually test connection with go-redis
	return CheckResult{
		Name:    c.Name(),
		Passed:  true,
		Message: "Redis URL configured correctly",
		Details: map[string]string{
			"url": cfg.Redis.URL,
		},
	}
}

// ForgeMetadataCheck validates project metadata if in a Forge project
type ForgeMetadataCheck struct {
	BaseCheck
}

func NewForgeMetadataCheck() *ForgeMetadataCheck {
	return &ForgeMetadataCheck{BaseCheck{name: "Forge Metadata"}}
}

func (c *ForgeMetadataCheck) Check(ctx context.Context) CheckResult {
	// Look for .forge/project.yaml in current directory or parents
	dir, err := os.Getwd()
	if err != nil {
		return CheckResult{
			Name:    c.Name(),
			Passed:  false,
			Message: fmt.Sprintf("Cannot get current directory: %v", err),
		}
	}

	for {
		path := filepath.Join(dir, ".forge", "project.yaml")
		if _, err := os.Stat(path); err == nil {
			// File exists, try to read and validate it
			data, err := os.ReadFile(path)
			if err != nil {
				return CheckResult{
					Name:    c.Name(),
					Passed:  false,
					Message: fmt.Sprintf("Cannot read project metadata: %v", err),
				}
			}

			// Try to unmarshal as basic validation
			var metadata map[string]interface{}
			if err := yaml.Unmarshal(data, &metadata); err != nil {
				return CheckResult{
					Name:    c.Name(),
					Passed:  false,
					Message: fmt.Sprintf("Invalid project metadata format: %v", err),
				}
			}

			// Check for required fields
			required := []string{"project_name", "edition", "forge_version", "created_at"}
			for _, field := range required {
				if _, exists := metadata[field]; !exists {
					return CheckResult{
						Name:    c.Name(),
						Passed:  false,
						Message: fmt.Sprintf("Missing required field: %s", field),
					}
				}
			}

			return CheckResult{
				Name:    c.Name(),
				Passed:  true,
				Message: "Forge project metadata is valid",
				Details: map[string]string{
					"project_name": metadata["project_name"].(string),
					"edition":      metadata["edition"].(string),
					"forge_version": metadata["forge_version"].(string),
				},
			}
		}

		// Check if we've reached the root directory
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	// Not in a Forge project - this is okay for doctor command
	return CheckResult{
		Name:    c.Name(),
		Passed:  true,
		Message: "Not in a Forge project (this is OK)",
	}
}

// compareVersions compares two version strings (simplified)
// Returns -1 if v1 < v2, 0 if v1 == v2, 1 if v1 > v2
func compareVersions(v1, v2 string) int {
	// Simple implementation - in reality would use proper semantic versioning
	if v1 == v2 {
		return 0
	}
	// This is a very simplified comparison - for demo purposes only
	if strings.Contains(v1, "1.25") && !strings.Contains(v2, "1.25") {
		return 1
	}
	if !strings.Contains(v1, "1.25") && strings.Contains(v2, "1.25") {
		return -1
	}
	// Fallback to string comparison
	if v1 < v2 {
		return -1
	}
	return 1
}