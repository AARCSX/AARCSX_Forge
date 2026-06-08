package forge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const gitignoreContent = `# Forge generated project files
.env
*.log
bin/
dist/
tmp/
.DS_Store
*.swp
*.swo
*~
.idea/
.vscode/
*.o
*.a
*.so
*.dylib
*.test
.env.local
vendor/
`

// CopyEmbeddedTemplate creates the basic directory structure for a new Forge project.
func CopyEmbeddedTemplate(destDir string) error {
	dirs := []string{
		"internal/app",
		"internal/config",
		"internal/database",
		"internal/health",
		"internal/identity",
		"internal/logger",
		"internal/notifications/provider",
		"internal/observability",
		"internal/platform/contextx/middleware",
		"internal/platform/events",
		"internal/platform/httpx",
		"internal/platform/middleware",
		"internal/platform/security",
		"internal/shared",
		"internal/storage/provider",
		"internal/tenants",
		"cmd/api",
		"cmd/worker",
		"migrations",
	}

	for _, dir := range dirs {
		dirPath := filepath.Join(destDir, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			return fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	placeholders := map[string]string{
		"internal/app/bootstrap.go": `package app

func BuildRuntimeDeps(cfg interface{}) interface{} {
	return nil
}
`,
		"internal/config/config.go": `package config

func Load() (interface{}, error) {
	return nil, nil
}
`,
		"internal/health/handler.go": `package health

func Handler() {
}
`,
		"internal/logger/logger.go": `package logger

func NewLogger() interface{} {
	return nil
}
`,
		"internal/shared/README.md": `# Shared Package

Common utilities and helpers.
`,
		"cmd/api/main.go": `package main

import "fmt"

func main() {
	fmt.Println("Forge API starting...")
}
`,
		"cmd/worker/main.go": `package main

import "fmt"

func main() {
	fmt.Println("Forge Worker starting...")
}
`,
	}

	for filePath, content := range placeholders {
		fullPath := filepath.Join(destDir, filePath)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			return fmt.Errorf("create directory: %w", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			return fmt.Errorf("write file: %w", err)
		}
	}

	return nil
}

// WriteProjectGitignore writes the .gitignore file.
func WriteProjectGitignore(destDir string) error {
	return os.WriteFile(filepath.Join(destDir, ".gitignore"), []byte(gitignoreContent), 0644)
}

// GenerateGoMod creates the go.mod file for the project.
func GenerateGoMod(destDir, moduleName string) error {
	goModContent := fmt.Sprintf(`module %s

go 1.25.0

require (
	github.com/gin-gonic/gin v1.12.0
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/google/uuid v1.6.0
	github.com/hibiken/asynq v0.26.0
	github.com/jackc/pgx/v5 v5.7.6
	github.com/minio/minio-go/v7 v7.0.72
	github.com/redis/go-redis/v9 v9.19.0
	github.com/spf13/cobra v1.10.2
	github.com/stretchr/testify v1.11.1
	go.uber.org/zap v1.28.0
	golang.org/x/crypto v0.52.0
	golang.org/x/time v0.15.0
	gopkg.in/yaml.v3 v3.0.1
)
`, moduleName)

	return os.WriteFile(filepath.Join(destDir, "go.mod"), []byte(goModContent), 0644)
}

// replacePlaceholders replaces templated strings in generated project files.
func replacePlaceholders(dir string, replacements map[string]string) error {
	return filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if shouldSkipFile(path) {
			return nil
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("read file: %w", err)
		}

		newContent := string(content)
		for placeholder, value := range replacements {
			newContent = strings.ReplaceAll(newContent, placeholder, value)
		}

		if newContent != string(content) {
			if err := os.WriteFile(path, []byte(newContent), 0644); err != nil {
				return fmt.Errorf("write file: %w", err)
			}
		}

		return nil
	})
}

// shouldSkipFile determines if a file should be skipped.
func shouldSkipFile(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".bmp", ".tiff", ".ico":
		return true
	case ".zip", ".tar", ".gz", ".bz2", ".xz", ".7z":
		return true
	case ".exe", ".dll", ".so", ".dylib":
		return true
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx":
		return true
	}

	if strings.HasPrefix(filepath.Base(path), ".") &&
		!strings.HasPrefix(filepath.Base(path), ".forge") {
		return true
	}

	return false
}
